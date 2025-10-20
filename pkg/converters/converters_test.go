package converters

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

func TestMergeStructs(t *testing.T) {
	type TestStruct struct {
		Name   string
		Email  string
		Age    int
		Active bool
		Height float64
		Weight float64
		Type   uint
	}

	existing := TestStruct{
		Name:   "João",
		Email:  "joao@email.com",
		Age:    30,
		Active: true,
		Height: 1.80,
		Weight: 70.0,
		Type:   1,
	}

	update := TestStruct{
		Name:   "João Silva",
		Email:  "",
		Age:    0,
		Active: false,
		Height: 1.80,
		Type:   0,
	}

	result := MergeStructs(existing, update).(TestStruct)

	expected := TestStruct{
		Name:   "João Silva",
		Email:  "joao@email.com",
		Age:    30,
		Active: false,
		Height: 1.80,
		Weight: 70.0,
		Type:   1,
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

func TestParamsToMap(t *testing.T) {
	params := url.Values{
		"name":  {"Jon Snow"},
		"email": {"jon@winterfell.com"},
		"age":   {"30"},
	}

	result := ParamsToMap(params)

	expected := map[string]interface{}{
		"name":  "Jon Snow",
		"email": "jon@winterfell.com",
		"age":   "30",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}
