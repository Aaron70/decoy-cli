# Decoy CLI

A CLI tool for generating and ingesting mock data using Go templates and configurable runners.

---

## Features

* Template Engine - Generate dynamic data using Go's `text/template`. See [Template Engine](#template-engine) section.
    * Use built-in functions like: random generation function, probability, counters and more.
* Runners Engine - Ingest generated data using implemented Runners. See [Runners Engine](#runners-engine) section.
    * Concurrent Execution - You can set multiple goroutines to execute `n` times the runner.
* Persistance - You can save Templates and Runners to reuse later.

---

## Installation

### Install from source

```bash
go install github.com/aaron70/decoy-cli/decoy@latest
```

### Build locally

```bash
git clone https://github.com/aaron70/decoy-cli.git
cd decoy-cli
go build -o decoy ./decoy/main.go
```

---

## Quick Start

### Store a template

```bash
decoy template store greet -t 'Hello, {{ Coalesce .Name "World" }}!'
```

### Parse the template

```bash
decoy template parse greet --data '{ "Name": "Doe" }'
```

### Store a runner

```bash
decoy runner store echo -c 'echo User said: "{{ .Template }}"' 
```

### Execute the runner

```bash
decoy runner run cmd echo greet -v Name=Doe
```

### Explore the commands

```bash
decoy template --help
decory template parse --help
decoy runner --help
decory runner run --help
```

---

## Template Engine

Templates use Go's [`text/template`](https://pkg.go.dev/text/template) syntax with additional functions from the [decoy](https://github.com/aaron70/decoy) library.

### Built-in template functions

#### Random generation

| Function | Description | Example |
|---|---|---|
| `RandomInt min max` | Random int in `[min, max)` | `{{RandomInt 1 100}}` |
| `RandomFloat min max` | Random float64 in `[min, max)` | `{{RandomFloat 0.0 1.0}}` |
| `RandomBoolean` | Random true/false | `{{RandomBoolean}}` |
| `RandomChoice args...` | Random item from args | `{{RandomChoice "a" "b" "c"}}` |
| `RandomText maxWords` | Markov chain random text | `{{RandomText 50}}` |

#### Probability

| Function | Description | Example |
|---|---|---|
| `Probability p` | Returns true with probability `p` | `{{Probability 0.75}}` |

#### Incremental counters

| Function | Description | Example |
|---|---|---|
| `NextIncrementalInt id start step` | Next value of named counter | `{{NextIncrementalInt "c" 1 1}}` |
| `CurrentIncrementalInt id default` | Current value (non-advancing) | `{{CurrentIncrementalInt "c" 0}}` |

#### Environment & I/O

| Function | Description | Example |
|---|---|---|
| `EnvVariable key` | Read environment variable | `{{EnvVariable "HOME"}}` |
| `ReadFileString path` | Read file as string | `{{ReadFileString "data.txt"}}` |
| `ReadFileBytes path` | Read file as bytes | `{{ReadFileBytes "img.png"}}` |
| `ReadFileBase64 path` | Read file as base64 | `{{ReadFileBase64 "img.png"}}` |

#### Serialization

| Function | Description | Example |
|---|---|---|
| `JsonUnmarshalString data` | Parse JSON string to map | `{{(JsonUnmarshalString "{\"a\":1}").a}}` |
| `JsonUnmarshalBytes data` | Parse JSON bytes to map | `{{$d := ReadFileBytes "cfg.json"}}{{JsonUnmarshalBytes $d}}` |

#### List builders

| Function | Description | Example |
|---|---|---|
| `List args...` | Mixed-type slice | `{{range List "a" 1 true}}{{.}} {{end}}` |
| `ListString args...` | String slice | `{{range ListString "x" "y"}}{{.}} {{end}}` |
| `ListInt args...` | Int slice | `{{range ListInt 10 20}}{{.}} {{end}}` |
| `ListFloat64 args...` | Float64 slice | `{{range ListFloat64 1.5 2.5}}{{.}} {{end}}` |
| `ListBool args...` | Bool slice | `{{range ListBool true false}}{{.}} {{end}}` |

#### Coalesce

| Function | Description | Example |
|---|---|---|
| `Coalesce args...` | First non-zero value | `{{Coalesce .Title "Default"}}` |
| `CoalesceString args...` | String variant | `{{CoalesceString .Name "unknown"}}` |
| `CoalesceInt args...` | Int variant | `{{CoalesceInt .Count 0}}` |
| `CoalesceFloat64 args...` | Float64 variant | `{{CoalesceFloat64 .Price 9.99}}` |

#### Error

| Function | Description | Example |
|---|---|---|
| `NewError msg args...` | Halts template with formatted error | `{{NewError "bad value: %v" .Val}}` |


See [decoy's FUNCTIONS.md](https://github.com/aaron70/decoy/blob/main/FUNCTIONS.md) for full details.

---

## Runner Engine

Decoy provides some Runners to allow you to parse and ingest the data generated to your applications.

You can execute a Runner `n` times with `g` concurrent Goroutines, to generate multiple records with the given runner and given template.
The runner's configuration support the [Template Engine](#template-engine) so you can use the built-in functions within the configuration.
Additionally, Runners provide the following data available to the Template Engine:

| Key | Description | Example |
| --- | ----------- | ------- |
| Times | The number of times that the runner will be executed | `{{ .Times }}` |
| Goroutines | The number of concurrent goroutines that will execute the runner | `{{ .Goroutines }}` |
| Template | The contents of the already parsed template | `{{ .Template }}` |


### Runner types

#### `http` — HTTP requests

Make HTTP requests with configurable method, URL, query parameters, headers, and body.

Configuration:

The configuration for the `http` runner is a JSON object with the HTTP request configuration.

```json
{
  "method": "<HTTP Method>",
  "url": "<URL>",
  "queryParameters": <JSON Object>,
  "headers": <JSON Object>,
  "body": <Any JSON Node>
}
```

Configuration example:
```json
{
  "method": "POST",
  "url": "http://api.example.com/items",
  "queryParameters": {"source": "decoy"},
  "headers": {"Authorization": "Bearer token123"},
  "body": {{.Template}}
}
```

#### `cmd` — Shell commands

Executes the shell command and prints the results from stdout and stderr. 

Configuration:

The configuration for the `cmd` runner is the command itself to run. 

Configuration example:
```text
echo "{{.Template}}"
```

---
