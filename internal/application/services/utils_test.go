package services

import (
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

func TestMergeStructs(t *testing.T) {
	type TestStruct struct {
		Name   string
		Email  string
		Age    int
		Active bool
	}

	existing := TestStruct{
		Name:   "João",
		Email:  "joao@email.com",
		Age:    30,
		Active: true,
	}

	update := TestStruct{
		Name:   "João Silva",
		Email:  "",
		Age:    0,
		Active: false,
	}

	result := MergeStructs(existing, update).(TestStruct)

	expected := TestStruct{
		Name:   "João Silva",
		Email:  "joao@email.com",
		Age:    30,
		Active: false,
	}

	if result != expected {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestMergeStructsWithDomainUser(t *testing.T) {
	existing := domain.User{
		ID:      1,
		Name:    "João",
		Email:   "joao@email.com",
		Contact: "123456789",
		Role:    "user",
		Active:  true,
	}

	update := domain.User{
		ID:      1,
		Name:    "João Silva",
		Email:   "",
		Contact: "",
		Role:    "admin",
		Active:  false,
	}

	result := MergeStructs(existing, update).(domain.User)

	expected := domain.User{
		ID:      1,
		Name:    "João Silva",
		Email:   "joao@email.com",
		Contact: "123456789",
		Role:    "admin",
		Active:  false,
	}

	if result.Name != expected.Name {
		t.Errorf("Name: expected %s, got %s", expected.Name, result.Name)
	}
	if result.Email != expected.Email {
		t.Errorf("Email: expected %s, got %s", expected.Email, result.Email)
	}
	if result.Contact != expected.Contact {
		t.Errorf("Contact: expected %s, got %s", expected.Contact, result.Contact)
	}
	if result.Role != expected.Role {
		t.Errorf("Role: expected %s, got %s", expected.Role, result.Role)
	}
	if result.Active != expected.Active {
		t.Errorf("Active: expected %t, got %t", expected.Active, result.Active)
	}
}
