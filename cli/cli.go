package cli

import (
	"io"
	"path"

	"github.com/aaron70/decoy"
	"github.com/aaron70/decoy-cli/internal/model"
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
	templateRepo, err := repositories.NewFSRepository[string, model.Template](path.Join(basePath, "templates"))
	if err != nil {
		return nil, errors.NewError(nil, err, "Couldn't create the templates repository")
	}
	runnerRepo, err := repositories.NewFSRepository[string, model.Runner](path.Join(basePath, "runners"))
	if err != nil {
		return nil, errors.NewError(nil, err, "Couldn't create the runners repository")
	}
	decoy, err := decoy.NewDecoyWithSeed(0)
	if err != nil {
		return nil, errors.NewError(nil, err, "Couldn't create the decoy instance")
	}
	return &CLI{
		TemplateSvc: services.NewTemplate(templateRepo),
		RunnerSvc:   services.NewRunner(runnerRepo, decoy),
		Decoy:       decoy,
	}, nil
}

func (c CLI) ReadStringFrom(r io.Reader) (string, error) {
	contents, err := io.ReadAll(r)
	return string(contents), err
}
