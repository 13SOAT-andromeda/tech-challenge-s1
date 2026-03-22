package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_ValidateRole(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		wantErr bool
	}{
		{
			name:    "valid role - customer",
			role:    "customer",
			wantErr: false,
		},
		{
			name:    "valid role - attendant",
			role:    "attendant",
			wantErr: false,
		},
		{
			name:    "valid role - mechanic",
			role:    "mechanic",
			wantErr: false,
		},
		{
			name:    "valid role - administrator",
			role:    "administrator",
			wantErr: false,
		},
		{
			name:    "invalid role - empty string",
			role:    "",
			wantErr: true,
		},
		{
			name:    "invalid role - manager",
			role:    "manager",
			wantErr: true,
		},
		{
			name:    "invalid role - admin (wrong case)",
			role:    "admin",
			wantErr: true,
		},
		{
			name:    "invalid role - Customer (uppercase)",
			role:    "Customer",
			wantErr: true,
		},
		{
			name:    "invalid role - random string",
			role:    "xyz123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				Role: tt.role,
			}

			err := user.ValidateRole()

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRole() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil && tt.wantErr {
				t.Logf("Expected error received: %v", err)
			}
		})
	}
}

func TestUser_IsCustomer(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{name: "customer role returns true", role: "customer", want: true},
		{name: "attendant role returns false", role: "attendant", want: false},
		{name: "mechanic role returns false", role: "mechanic", want: false},
		{name: "administrator role returns false", role: "administrator", want: false},
		{name: "empty role returns false", role: "", want: false},
		{name: "Customer uppercase returns false", role: "Customer", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Role: tt.role}
			assert.Equal(t, tt.want, u.IsCustomer())
		})
	}
}

func TestUser_IsEmployee(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{name: "attendant role returns true", role: "attendant", want: true},
		{name: "mechanic role returns true", role: "mechanic", want: true},
		{name: "administrator role returns true", role: "administrator", want: true},
		{name: "customer role returns false", role: "customer", want: false},
		{name: "empty role returns false", role: "", want: false},
		{name: "Attendant uppercase returns false", role: "Attendant", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Role: tt.role}
			assert.Equal(t, tt.want, u.IsEmployee())
		})
	}
}
