package services

import (
	"context"
	"errors"

	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/converters"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/encryption"

	employeeModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/employee"
	personModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/person"
	userModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type userService struct {
	repo         ports.UserRepository
	personRepo   ports.PersonRepository
	employeeRepo ports.EmployeeRepository
}

func NewUserService(repo ports.UserRepository, personRepo ports.PersonRepository, employeeRepo ports.EmployeeRepository) *userService {
	return &userService{repo: repo, personRepo: personRepo, employeeRepo: employeeRepo}
}

func (s *userService) Create(ctx context.Context, u domain.User) (*domain.User, error) {
	if u.Person == nil {
		return nil, errors.New("person data is needed")
	}

	if existing, err := s.GetByEmail(ctx, u.Person.Email); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, errors.New("email já existe")
	}

	pm := &personModel.Model{}
	pm.FromDomain(u.Person)
	createdPerson, err := s.personRepo.Create(ctx, pm)
	if err != nil {
		return nil, err
	}

	if u.IsEmployee() {
		e := &employeeModel.Model{}
		e.Person = *createdPerson
		e.FromDomain(u.Employee)

		_, err := s.employeeRepo.Create(ctx, e)
		if err != nil {
			return nil, err
		}
	}

	if err := u.Password.Hash(); err != nil {
		return nil, err
	}
	u.PersonID = createdPerson.ID
	u.Person = createdPerson.ToDomain()

	um := &userModel.Model{}
	um.FromDomain(&u)

	_, err = s.repo.Create(ctx, um)
	if err != nil {
		return nil, err
	}

	um.Person = *createdPerson
	return um.ToDomain(), nil
}

func (s *userService) CreateAdminUser(ctx context.Context, email, password, document string) error {

	if user, err := s.GetByEmail(ctx, email); err != nil {
		return err
	} else if user != nil {
		return nil
	}

	doc, err := domain.NewDocument(document)
	if err != nil {
		return err
	}

	p, err := domain.NewPassword(password, encryption.NewBcryptHasher())
	if err != nil {
		return err
	}

	if err := p.Hash(); err != nil {
		return err
	}

	person := domain.Person{
		Name:     "Administrador",
		Email:    email,
		Document: doc,
		Contact:  "+5511954945277",
	}

	pm := &personModel.Model{}
	pm.FromDomain(&person)
	createdPerson, err := s.personRepo.Create(ctx, pm)
	if err != nil {
		return err
	}

	u := domain.User{
		Password: p,
		Role:     "administrator",
		PersonID: createdPerson.ID,
	}

	um := &userModel.Model{}
	um.FromDomain(&u)

	if _, err = s.repo.Create(ctx, um); err != nil {
		return err
	}

	emp := &employeeModel.Model{
		Position: "Administrador",
		PersonID: createdPerson.ID,
	}

	if _, err = s.employeeRepo.Create(ctx, emp); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

func (s *userService) Search(ctx context.Context, params map[string]interface{}) (*[]domain.User, error) {
	uSearch := ports.UserSearch{Name: "", Email: "", Contact: ""}
	if params["name"] != nil {
		uSearch.Name = params["name"].(string)
	}
	if params["email"] != nil {
		uSearch.Email = params["email"].(string)
	}
	if params["contact"] != nil {
		uSearch.Contact = params["contact"].(string)
	}

	users := s.repo.Search(ctx, uSearch)
	usersD := make([]domain.User, 0, len(users))
	for _, u := range users {
		usersD = append(usersD, *u.ToDomain())
	}

	return &usersD, nil
}

func (s *userService) Update(ctx context.Context, u domain.User) (*domain.User, error) {
	existingUser, err := s.repo.FindByID(ctx, u.ID)
	if err != nil || existingUser == nil {
		return nil, errors.New("user not found")
	}

	existingDomain := existingUser.ToDomain()

	if u.Person != nil && u.Person.Email != "" && existingDomain.Person != nil && u.Person.Email != existingDomain.Person.Email {
		if uMail, err := s.GetByEmail(ctx, u.Person.Email); err != nil {
			return nil, err
		} else if uMail != nil && existingUser.ID != uMail.ID {
			return nil, errors.New("email já existe")
		}
	}

	if u.Password != nil && u.Password.GetValue() != "" {
		return nil, errors.New("senha de usuário não pode ser atualizada")
	}

	if u.Person != nil && existingDomain.Person != nil {
		existingPerson := *existingDomain.Person
		mergedPerson := converters.MergeStructs(existingPerson, *u.Person).(domain.Person)
		mergedPerson.ID = existingDomain.PersonID

		pm := &personModel.Model{}
		pm.FromDomain(&mergedPerson)
		if err := s.personRepo.Update(ctx, pm); err != nil {
			return nil, err
		}
	}

	mergedUser := converters.MergeStructs(existingDomain, u).(domain.User)
	mergedUser.Password = domain.NewPasswordFromHash(existingUser.Password, encryption.NewBcryptHasher())

	um := &userModel.Model{}
	um.FromDomain(&mergedUser)

	if err = s.repo.Update(ctx, um); err != nil {
		return nil, err
	}

	return um.ToDomain(), nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.New("ocorreu um erro ao excluir o usuário")
	}
	return nil
}

func (s *userService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, nil
	}
	return u.ToDomain(), nil
}
