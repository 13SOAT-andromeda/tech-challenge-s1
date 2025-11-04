package domain

import (
	"testing"
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
