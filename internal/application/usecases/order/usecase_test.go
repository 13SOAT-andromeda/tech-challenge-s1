package order

import (
	"context"
	"errors"
	"testing"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helper to create product/maintenance with price
func p(id uint, price int64) domain.Product {
	return domain.Product{ID: id, Price: price, Name: "p"}
}
func ma(id uint, price int64) domain.Maintenance {
	return domain.Maintenance{ID: id, Price: price, Name: "m"}
}
func ptrString(s string) *string { return &s }

func TestCreateOrder_Success(t *testing.T) {
	ctx := context.Background()
	mockOrder := new(mocks.MockOrderService)
	mockProd := new(mocks.MockProductService)
	mockMaint := new(mocks.MockMaintenanceService)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint)

	input := ports.CreateOrderInput{
		VehicleKilometers: 123,
		Note:              ptrString("note"),
		UserID:            10,
		CustomerVehicleID: 20,
		CompanyID:         30,
		ProductIDs:        []uint{1, 2},
		MaintenanceIDs:    []uint{3},
	}

	products := []domain.Product{p(1, 100), p(2, 150)}
	maints := []domain.Maintenance{ma(3, 250)}

	mockProd.On("GetByIds", ctx, input.ProductIDs).Return(products, nil)
	mockMaint.On("GetByIDs", ctx, input.MaintenanceIDs).Return(maints, nil)

	// Expect order creation, check price value passed via matcher
	mockOrder.On("Create", mock.Anything, mock.MatchedBy(func(o domain.Order) bool {
		if o.VehicleKilometers != input.VehicleKilometers {
			return false
		}
		if o.User.ID != input.UserID {
			return false
		}
		if o.CustomerVehicle.ID != input.CustomerVehicleID {
			return false
		}
		if o.Company.ID != input.CompanyID {
			return false
		}
		if o.Note == nil || *o.Note != *input.Note {
			return false
		}
		if o.Price == nil {
			return false
		}
		// total: 100 + 150 + 250 = 500
		return *o.Price == 500.0
	})).Return(&domain.Order{ID: 1}, nil)

	created, err := uc.CreateOrder(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, uint(1), created.ID)

	mockProd.AssertExpectations(t)
	mockMaint.AssertExpectations(t)
	mockOrder.AssertExpectations(t)
}

func TestCreateOrder_ProductServiceError(t *testing.T) {
	ctx := context.Background()
	mockOrder := new(mocks.MockOrderService)
	mockProd := new(mocks.MockProductService)
	mockMaint := new(mocks.MockMaintenanceService)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint)

	input := ports.CreateOrderInput{ProductIDs: []uint{1}}

	mockProd.On("GetByIds", ctx, input.ProductIDs).Return(nil, errors.New("prod error"))

	created, err := uc.CreateOrder(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, created)
	mockProd.AssertExpectations(t)
}

func TestCreateOrder_MaintenanceServiceError(t *testing.T) {
	ctx := context.Background()
	mockOrder := new(mocks.MockOrderService)
	mockProd := new(mocks.MockProductService)
	mockMaint := new(mocks.MockMaintenanceService)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint)

	input := ports.CreateOrderInput{ProductIDs: []uint{}, MaintenanceIDs: []uint{1}}

	mockProd.On("GetByIds", ctx, input.ProductIDs).Return([]domain.Product{}, nil)
	mockMaint.On("GetByIDs", ctx, input.MaintenanceIDs).Return(nil, errors.New("maint error"))

	created, err := uc.CreateOrder(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, created)
	mockProd.AssertExpectations(t)
	mockMaint.AssertExpectations(t)
}
