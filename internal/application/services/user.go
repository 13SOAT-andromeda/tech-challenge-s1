package services

import (
	"context"

	appErrors "github.com/13SOAT-andromeda/tech-challenge-s1/internal/core/errors"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, u domain.User) (*domain.User, error) {

	if u.Address == nil {
		u.Address = &domain.Address{}
	}

	userModel := model.NewUserModelFromDomain(u)

	_, err := s.repo.Create(ctx, &userModel)

	if err != nil {
		return nil, err
	}

	created := userModel.ToDomain()

	return &created, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]domain.User, error) {

	users, err := s.repo.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	usersD := make([]domain.User, 0, len(users))

	for _, user := range users {
		usersD = append(usersD, user.ToDomain())
	}

	return usersD, nil
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*domain.User, error) {

	user, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}
	customerDomain := user.ToDomain()

	return &customerDomain, nil
}

func (s *UserService) Search(ctx context.Context, params map[string]interface{}) (*[]domain.User, error) {

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

	for _, user := range users {
		usersD = append(usersD, user.ToDomain())
	}

	return &usersD, nil
}

func (s *UserService) Update(ctx context.Context, u domain.User) (*domain.User, error) {

	existingUser, err := s.repo.FindByID(ctx, u.ID)

	if err != nil || existingUser == nil {
		return nil, appErrors.ErrUserNotFound
	}

	existingDomain := existingUser.ToDomain()

	if u.Email != "" && u.Email != existingDomain.Email {
		if exists, err := s.exists(ctx, u.ID, u.Email); err != nil {
			return nil, err
		} else if exists {
			return nil, appErrors.ErrUserEmailAlreadyExists
		}
	}

	if u.Password != nil {
		return nil, appErrors.ErrUserPasswordUpdateNotAllowed
	}

	mergedUser := MergeStructs(existingDomain, u).(domain.User)

	mergedUser.Password = domain.NewPasswordFromHash(existingUser.Password)

	userModel := model.NewUserModelFromDomain(mergedUser)

	err = s.repo.Update(ctx, &userModel)
	if err != nil {
		return nil, err
	}

	updated := userModel.ToDomain()

	return &updated, nil
}

func (s *UserService) Delete(ctx context.Context, id uint) error {

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return appErrors.ErrUserDelete
	}

	return nil
}

func (s *UserService) exists(ctx context.Context, id uint, email string) (bool, error) {

	return s.repo.Exists(ctx, id, email)
}
