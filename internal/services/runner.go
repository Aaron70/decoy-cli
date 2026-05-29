package services

import (
	"bytes"
	"context"
	"fmt"
	"io"

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

	pool, err := concurrency.NewPool(ctx,
		concurrency.NewPoolWithMaxWorkers[string, error](workers),
	)
	if err != nil {
		return err
	}

	tasks := make(chan string, n)
	errors := make(chan error, 1)

	go func() {
		defer close(tasks)
		for range n {

			parsedTemplate := bytes.NewBufferString("")
			err := svc.Decoy.ParseTemplate(parsedTemplate, tmpl,
				decoy.WithTemplateNamed("template"),
				decoy.WithData(data),
			)
			if err != nil {
				channels.Send(ctx, errors, fmt.Errorf("TemplateParseError: %w", err))
				return
			}

			runnerData := map[string]any{
				"Template": parsedTemplate,
			}
			parsedConfiguration := bytes.NewBufferString("")

			err = svc.Decoy.ParseTemplate(parsedConfiguration, config,
				decoy.WithTemplateNamed("runner"),
				decoy.WithData(runnerData),
			)
			if err != nil {
				channels.Send(ctx, errors, fmt.Errorf("RunnerConfigurationParseError: %w", err))
				return
			}

			err = channels.Send(ctx, tasks, parsedConfiguration.String())
			if err != nil {
				channels.Send(ctx, errors, err)
				return
			}
		}
	}()

	pool.PushTasks(tasks, func(ctx context.Context, config string) {
		res, err := runner.Run(config)
		if err != nil {
			channels.Send(ctx, errors, err)
			return
		}

		if w != nil {
			fmt.Fprintf(w, "%s", res)
		}
	})
	go func() {
		pool.Wait()
		close(errors)
	}()

	var errCtx error
	open := true
	for open {
		err, open, errCtx = channels.Recv(ctx, errors)
		if errCtx != nil {
			close(errors)
			return errCtx
		}
		if err != nil {
			close(errors)
			return err
		}
	}

	return nil
}
