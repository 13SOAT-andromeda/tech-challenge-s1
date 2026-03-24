package mocks

import (
	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
)

type MockPersonRepository struct {
	MockGenericRepository[personModel.Model]
}

var _ ports.PersonRepository = (*MockPersonRepository)(nil)
