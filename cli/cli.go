package cli

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"
	"time"

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

func (c CLI) AskForInput(r io.Reader, w io.Writer, msg string, args ...any) (string, error) {
	fmt.Fprintf(w, msg, args...)
	return readLine(r)
}

func readLine(r io.Reader) (string, error) {
	br := bufio.NewReader(r)
	line, err := br.ReadString('\n')
	return strings.TrimRight(line, "\r\n"), err
}

func (c CLI) ReadStringFrom(r io.Reader) (string, error) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Go(func() {
		select {
		case <-ctx.Done():
		case <-time.Tick(time.Second * 2):
			fmt.Println("Waiting input from stdin...")
		}
	})

	contents, err := io.ReadAll(r)
	return string(contents), err
}
