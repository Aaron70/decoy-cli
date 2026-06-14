package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/aaron70/decoy"
	"github.com/aaron70/goaty/validations"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

func decoyResponse(w http.ResponseWriter, status int, message string, args ...any) {
	w.WriteHeader(999)
	w.Header().Add("Content-Type", "application/json")
	res := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  status,
		Message: fmt.Sprintf(message, args...),
	}
	bytes, _ := json.MarshalIndent(res, "", "  ")
	w.Write(bytes)
}

func responseToStatusCode(response string) int {
	status, err := strconv.Atoi(response)
	if err != nil {
		switch strings.ToLower(response) {
		case "1xx":
			return 100
		case "2xx":
			return 200
		case "3xx":
			return 300
		case "4xx":
			return 400
		case "5xx":
			return 500
		default:
			return http.StatusMultiStatus
		}
	}
	return status
}

func (s RestServer) mockHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			decoyResponse(w, http.StatusInternalServerError, "The endpoint expects a well formed and complete OpenAPI v3 Specification. If the specification is missing something and the endpoint is trying to access it, the endpoint will panic. Please read carefully your server spec. Endpoint has panic with: %v", r)
		}
	}()

	r.URL.Path = strings.TrimPrefix(r.URL.Path, "/mock")
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	var err error
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			decoyResponse(w, http.StatusBadRequest, "Couldn't read request body: %v", err)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	selectedResponse := r.URL.Query().Get("decoy-response")
	selectedContentType := r.URL.Query().Get("decoy-content-type")
	if validations.StrIsBlank(selectedContentType) {
		selectedContentType = r.Header.Get("Content-Type")
	}
	selectedExample := r.URL.Query().Get("decoy-example")
	decoyParse := r.URL.Query().Get("decoy-parse")
	shouldParse := false
	if validations.StrIsBlank(decoyParse) || strings.HasPrefix(selectedContentType, "text/") {
		shouldParse = true
	}

	route, pathParams, err := s.specRouter.FindRoute(r)
	if err != nil {
		decoyResponse(w, http.StatusNotFound, "Unknown endpoint: %s", err)
		return
	}

	input := &openapi3filter.RequestValidationInput{
		Request:    r,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(r.Context(), input); err != nil {
		decoyResponse(w, http.StatusBadRequest, "Invalid request: %s", err)
		return
	}

	op := route.Operation

	responses := op.Responses.Map()
	var response *openapi3.ResponseRef
	if validations.StrIsBlank(selectedResponse) {
		for selectedResponse, response = range responses {
			break
		}
	} else {
		var ok bool
		response, ok = responses[selectedResponse]
		if !ok {
			decoyResponse(w, http.StatusNoContent, "No responses in the spec for the given response: %s", selectedResponse)
			return
		}
	}
	if response == nil {
		decoyResponse(w, http.StatusNoContent, "No responses in the spec for the given requests")
		return
	}

	var contentType *openapi3.MediaType
	if validations.StrIsBlank(selectedContentType) {
		for selectedContentType, contentType = range response.Value.Content {
			break
		}
	} else {
		var ok bool
		contentType, ok = response.Value.Content[selectedContentType]
		if !ok {
			decoyResponse(w, http.StatusNoContent, "No media type in the spec for the response %q with media type: %s", selectedResponse, selectedContentType)
			return
		}
	}
	if contentType == nil {
		decoyResponse(w, http.StatusNoContent, "No media type in the spec for the response %q with media type: %s", selectedResponse, selectedContentType)
		return
	}

	examples := contentType.Examples
	var example *openapi3.ExampleRef
	if validations.StrIsBlank(selectedExample) {
		for selectedExample, example = range examples {
			break
		}
	} else {
		var ok bool
		example, ok = examples[selectedExample]
		if !ok {
			decoyResponse(w, http.StatusNoContent, "No examples in the spec for the response %q and Content-Type %q with the given example: %s", selectedResponse, selectedContentType, selectedExample)
			return
		}
	}
	if example == nil {
		decoyResponse(w, http.StatusNoContent, "No examples in the spec for the given response %q with Content-Type %q", selectedResponse, selectedContentType)
		return
	}

	statusCode := responseToStatusCode(selectedResponse)
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", selectedContentType)
	if shouldParse {
		requestContentType := r.URL.Query().Get("Content-Type")
		var parsedBody any
		if len(bodyBytes) > 0 {
			requestContentType := r.Header.Get("Content-Type")
			if strings.HasPrefix(requestContentType, "application/json") {
				if err := json.Unmarshal(bodyBytes, &parsedBody); err != nil {
					parsedBody = string(bodyBytes)
				}
			} else {
				parsedBody = string(bodyBytes)
			}
		}

		s.Logger.Printf("Body Parsed: %+v", parsedBody)
		data := map[string]any{
			"Request": map[string]any{
				"Method":       r.Method,
				"Path":         r.URL.String(),
				"Header":       r.Header,
				"Body":         parsedBody,
				"Content-Type": requestContentType,
			},
			"Response": map[string]any{
				"ContentType": selectedContentType,
				"StatusCode":  statusCode,
				"Example":     selectedExample,
			},
		}
		str, ok := example.Value.Value.(string)
		if !ok {
			decoyResponse(w, 500, "couldn't parse template example: template example is not a string")
			return
		}
		tmpl, err := s.Decoy.Decoy.ParseTemplateString(str, decoy.WithData(data))
		if err != nil {
			decoyResponse(w, 500, "couldn't parse template example: %s", err)
			return
		}
		writeBody(w, selectedContentType, tmpl)
	} else {
		writeBody(w, selectedContentType, example.Value.Value)
	}
}

func writeBody(w http.ResponseWriter, contentType string, body any) {
	if body == nil {
		return
	}

	switch contentType {
	case "text/plain":
		fmt.Fprintf(w, "%s", body.(string))
	case "application/json":
		writeJsonBody(w, body)
	default:
		fmt.Fprintf(w, "%+v", body)
	}
}

func writeJsonBody(w http.ResponseWriter, body any) {
	str, ok := body.(string)
	if ok {
		w.Write([]byte(str))
		return
	}

	bytes, ok := body.([]byte)
	if ok {
		w.Write(bytes)
		return
	}

	bytes, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		decoyResponse(w, 500, "Invalid body, couldn't deserialize the example. Please fix the example in the spec and try again. Error: %s", err)
		return
	}
	w.Write(bytes)
}
