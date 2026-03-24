package services

import (
	"context"

	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type personService struct {
	repo ports.PersonRepository
}

func NewPersonService(repo ports.PersonRepository) ports.PersonService {
	return &personService{repo: repo}
}

func (s *personService) Create(ctx context.Context, p domain.Person) (*domain.Person, error) {
	m := &personModel.Model{}
	m.FromDomain(&p)

	created, err := s.repo.Create(ctx, m)
	if err != nil {
		return nil, err
	}

	return created.ToDomain(), nil
}

func (s *personService) GetByID(ctx context.Context, id uint) (*domain.Person, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return m.ToDomain(), nil
}

func (s *personService) Update(ctx context.Context, p domain.Person) (*domain.Person, error) {
	m := &personModel.Model{}
	m.FromDomain(&p)

	if err := s.repo.Update(ctx, m); err != nil {
		return nil, err
	}

	return m.ToDomain(), nil
}

func (s *personService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
