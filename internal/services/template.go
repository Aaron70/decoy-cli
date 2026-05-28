package services

import "github.com/aaron70/goaty/repositories"


type Template struct {
	repo repositories.Repository[string, string]
}


func NewTemplate(repo repositories.Repository[string, string]) *Template {
	return &Template{
		repo: repo,
	}
}

func (svc Template) Save(name, contents string) (string, error) {
	return svc.repo.Save(name, contents)
}

func (svc Template) Update(name, contents string) (string, error) {
	return svc.repo.Update(name, contents)
}

func (svc Template) Get(name string) (string, error) {
	return svc.repo.Get(name)
}

func (svc Template) GetAll() ([]string, error) {
	return svc.repo.GetAll()
}

func (svc Template) Delete(id string) (string, error) {
	return svc.repo.Delete(id)
}
