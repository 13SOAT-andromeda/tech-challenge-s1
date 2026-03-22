package ports

import (
	"context"

	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type PersonRepository interface {
	Repository[personModel.Model]
}

type PersonService interface {
	Create(ctx context.Context, p domain.Person) (*domain.Person, error)
	GetByID(ctx context.Context, id uint) (*domain.Person, error)
	Update(ctx context.Context, p domain.Person) (*domain.Person, error)
	Delete(ctx context.Context, id uint) error
}
