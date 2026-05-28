package services

import "github.com/aaron70/goaty/repositories"

type Runner struct {
	repo repositories.Repository[string, string]
}


func NewRunner(repo repositories.Repository[string, string]) *Runner {
	return &Runner{
		repo: repo,
	}
}

func (svc Runner) Save(name, contents string) (string, error) {
	return svc.repo.Save(name, contents)
}

func (svc Runner) Get(name string) (string, error) {
	return svc.repo.Get(name)
}

func (svc Runner) GetAll() ([]string, error) {
	return svc.repo.GetAll()
}

func (svc Runner) Delete(id string) (string, error) {
	return svc.repo.Delete(id)
}
