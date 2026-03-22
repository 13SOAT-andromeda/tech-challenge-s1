package repository

import (
	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"gorm.io/gorm"
)

type personRepository struct {
	*BaseRepository[personModel.Model]
}

func NewPersonRepository(db *gorm.DB) ports.PersonRepository {
	return &personRepository{
		BaseRepository: NewBaseRepository[personModel.Model](db),
	}
}
