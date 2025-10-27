package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerInitialization(t *testing.T) {

	c := Customer{
		ID:    1,
		Name:  "Gedan",
		Email: "gedan@example.com",
		Document: &Document{
			Number: "45653421898",
		},
		Type:    "teste",
		Contact: "11999999999",
		Address: &Address{
			Address:       "Rua Teste",
			City:          "New York",
			AddressNumber: "317",
			ZipCode:       "1234",
			Neighborhood:  "New York",
			Country:       "Brasil",
		},
	}

	assert.NotNil(t, c)
	assert.Equal(t, "Gedan", c.Name)
	assert.Equal(t, "gedan@example.com", c.Email)
	assert.Equal(t, "45653421898", c.Document.GetDocumentNumber())
	assert.Equal(t, "teste", c.Type)
	assert.Equal(t, "11999999999", c.Contact)
	assert.Equal(t, "Rua Teste", c.Address.Address)
	assert.Equal(t, "317", c.Address.AddressNumber)
	assert.Equal(t, "New York", c.Address.Neighborhood)
	assert.Equal(t, "New York", c.Address.City)
	assert.Equal(t, "Brasil", c.Address.Country)
	assert.Equal(t, "1234", c.Address.ZipCode)
}

func TestValidateCustomerType(t *testing.T) {
	tests := []struct {
		name         string
		customerType string
		wantErr      bool
		errMsg       string
	}{
		{
			name:         "valid type - administrator",
			customerType: "administrator",
			wantErr:      false,
		},
		{
			name:         "valid type - attendant",
			customerType: "attendant",
			wantErr:      false,
		},
		{
			name:         "valid type - mechanic",
			customerType: "mechanic",
			wantErr:      false,
		},
		{
			name:         "invalid type - manager",
			customerType: "manager",
			wantErr:      true,
			errMsg:       "customer type 'manager' is not valid. Accepted types: [administrator attendant mechanic]",
		},
		{
			name:         "invalid type - empty string",
			customerType: "",
			wantErr:      true,
			errMsg:       "customer type '' is not valid. Accepted types: [administrator attendant mechanic]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Customer{
				Type: tt.customerType,
			}

			err := c.ValidateCustomerType()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				} else if err.Error() != tt.errMsg {
					t.Errorf("unexpected error message.\nExpected: %s\nGot: %s", tt.errMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
