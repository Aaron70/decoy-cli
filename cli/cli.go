package cli

import (
	"io"
	"path"

	"github.com/aaron70/decoy"
	"github.com/aaron70/decoy-cli/internal/services"
	"github.com/aaron70/goaty/errors"
	"github.com/aaron70/goaty/repositories"
)

type CLI struct {
	TemplateSvc *services.Template
	RunnerSvc   *services.Runner
	Decoy       *decoy.Decoy
}

func NewCLI(basePath string) (*CLI, error) {
	marshal := func(s string) ([]byte, error) {
		return []byte(s), nil
	}
	unmarshal := func(b []byte) (string, error) {
		return string(b), nil
	}

	templateRepo, err := repositories.NewFSRepositoryWithSerializer[string,](path.Join(basePath, "templates"), marshal, unmarshal)
	if err != nil {
		return nil, errors.NewError(nil, err, "Couldn't create the templates repository")
	}
	runnerRepo, err := repositories.NewFSRepositoryWithSerializer[string, ](path.Join(basePath, "runners"), marshal, unmarshal)
	if err != nil {
		return nil, errors.NewError(nil, err, "Couldn't create the runners repository")
	}
	decoy, err := decoy.NewDecoyWithSeed(0)
	if err != nil {
		return nil, errors.NewError(nil, err, "Couldn't create the decoy instance")
	}
	return &CLI{
		TemplateSvc: services.NewTemplate(templateRepo),
		RunnerSvc:   services.NewRunner(runnerRepo),
		Decoy:       decoy,
	}, nil
}

func (c CLI) ReadStringFrom(r io.Reader) (string, error) {
	contents, err := io.ReadAll(r)
	return string(contents), err
}
