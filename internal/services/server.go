package services

import (
	"os"
	"path"

	"github.com/aaron70/decoy-cli/internal/model"
	"github.com/aaron70/goaty/errors"
	"github.com/aaron70/goaty/repositories"
)

type Server struct {
	repo     repositories.Repository[string, model.ServerSpec]
	BasePath string
}

func NewServer(path string, repo repositories.Repository[string, model.ServerSpec]) (*Server, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, nil
	}
	return &Server{
		repo:     repo,
		BasePath: path,
	}, nil
}

func (svc Server) Update(name, _type, contents string) (model.ServerSpec, error) {
	path := svc.filePath(name)
	err := os.WriteFile(path, []byte(contents), os.ModePerm)
	if err != nil {
		return model.ServerSpec{}, nil
	}
	return svc.repo.Update(name, model.ServerSpec{
		Name: name,
		Type: _type,
		Spec: path,
	})
}

func (svc Server) Get(name string) (model.ServerSpec, error) {
	spec, err := svc.repo.Get(name)
	if err != nil {
		return model.ServerSpec{}, err
	}

	path := svc.filePath(name)
	contents, err := os.ReadFile(path)
	if err != nil {
		return model.ServerSpec{}, err
	}
	spec.Spec = string(contents)

	return spec, nil
}

func (svc Server) GetAll() ([]model.ServerSpec, error) {
	return svc.repo.GetAll()
}

func (svc Server) Delete(name string) (model.ServerSpec, error) {
	path := svc.filePath(name)
	err := os.Remove(path)
	if err != nil {
		return model.ServerSpec{}, errors.NewError(nil, err, "Couldn't remove the %s file", path)
	}
	return svc.repo.Delete(name)
}

func (svc Server) filePath(name string) string { return path.Join(svc.BasePath, name) }
