package model

import "github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"

type UserModel struct {
	ID       uint         `gorm:"primaryKey"`
	Name     string       `gorm:"not null"`
	Email    string       `gorm:"not null"`
	Contact  string       `gorm:"not null"`
	Address  AddressModel `gorm:"embedded"`
	Password string       `gorm:"not null"`
	Role     string       `gorm:"not null"`
	Active   bool         `gorm:"default:true"`
}

func (u *UserModel) TableName() string {
	return "Users"
}

func NewUserModelFromDomain(domain domain.User) UserModel {
	return UserModel{
		ID:       domain.ID,
		Name:     domain.Name,
		Email:    domain.Email,
		Contact:  domain.Contact,
		Role:     domain.Role,
		Password: domain.Password.GetHashed(),
		Address:  FromDomainAddress(domain.Address),
		Active:   domain.Active,
	}
}

func (u *UserModel) ToDomain() domain.User {
	pass := domain.NewPasswordFromHash(u.Password)

	return domain.User{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Contact:  u.Contact,
		Role:     u.Role,
		Password: pass,
		Active:   u.Active,
	}
}
