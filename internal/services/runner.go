package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/aaron70/decoy"
	"github.com/aaron70/decoy-cli/internal/model"
	"github.com/aaron70/decoy-cli/internal/runners"
	"github.com/aaron70/goaty/channels"
	"github.com/aaron70/goaty/concurrency"
	"github.com/aaron70/goaty/repositories"
)

type RunnerType string

const (
	CMD  RunnerType = "cmd"
	HTTP RunnerType = "http"
)

type Runner struct {
	repo    repositories.Repository[string, model.Runner]
	Runners map[RunnerType]runners.Runner
	Decoy   *decoy.Decoy
}

func NewRunner(repo repositories.Repository[string, model.Runner], decoy *decoy.Decoy) *Runner {
	return &Runner{
		repo:  repo,
		Decoy: decoy,
		Runners: map[RunnerType]runners.Runner{
			"http": runners.NewHttpRunner(),
			"cmd":  runners.NewCmdRunner(),
		},
	}
}

func (svc Runner) Save(name, contents string) (model.Runner, error) {
	return svc.repo.Save(name, model.Runner{
		Name:   name,
		Config: contents,
	})
}

func (svc Runner) Update(name, contents string) (model.Runner, error) {
	return svc.repo.Update(name, model.Runner{
		Name:   name,
		Config: contents,
	})
}

func (svc Runner) Get(name string) (model.Runner, error) {
	return svc.repo.Get(name)
}

func (svc Runner) GetAll() ([]model.Runner, error) {
	return svc.repo.GetAll()
}

func (svc Runner) Delete(id string) (model.Runner, error) {
	return svc.repo.Delete(id)
}

func (svc Runner) Run(w io.Writer, _type RunnerType, config, tmpl string, data any, n int, workers int) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner, exists := svc.Runners[_type]
	if !exists {
		return fmt.Errorf("Invalid Runner type %q, not registered", _type)
	}

	templateCompiled, err := svc.Decoy.CompileTemplate(tmpl,
		decoy.WithTemplateNamed("template"),
	)
	if err != nil {
		return fmt.Errorf("TemplateParseError: %w", err)
	}

	runnerCompiled, err := svc.Decoy.CompileTemplate(config,
		decoy.WithTemplateNamed("runner"),
	)
	if err != nil {
		return fmt.Errorf("RunnerConfigurationParseError: %w", err)
	}

	var wg sync.WaitGroup
	records := make(chan string, n)
	errs := make(chan error, n)

	// TODO: Consider using a pool of workers
	wg.Go(func() {
		defer close(records)

		buffer := new(bytes.Buffer)
		for range n {
			buffer.Reset()
			err := templateCompiled.Execute(buffer, data)
			if err != nil {
				channels.Send(ctx, errs, err)
				return
			}

			runnerData := map[string]any{
				"Template":   buffer.String(),
				"Times":      n,
				"Goroutines": workers,
			}
			buffer.Reset()

			err = runnerCompiled.Execute(buffer, runnerData)
			if err != nil {
				channels.Send(ctx, errs, err)
				return
			}

			err = channels.Send(ctx, records, buffer.String())
			if err != nil {
				channels.Send(ctx, errs, err)
				return
			}
		}
	})

	pool, _ := concurrency.NewPool(ctx,
		concurrency.NewPoolWithMaxWorkers[string](workers),
		concurrency.NewPoolWithBufferSize[string](n),
	)

	err = pool.PushTasks(records, func(ctx context.Context, config string) {
		res, err := runner.Run(ctx, config)
		if err != nil {
			channels.Send(ctx, errs, err)
			return
		}

		if w != nil {
			fmt.Fprintf(w, "%s", res)
		}
	})
	if err != nil {
		return err
	}

	go func() {
		pool.Wait()
		close(errs)
	}()

	var errCtx error
	open := true
	for open {
		err, open, errCtx = channels.Recv(ctx, errs)
		if errCtx != nil {
			return errCtx
		}
		if err != nil {
			return err
		}
	}

	return nil
}

