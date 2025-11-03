package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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
	mockOrderRepo := new(mocks.MockOrderRepository)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockOrderRepo)

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
	mockOrderRepo := new(mocks.MockOrderRepository)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockOrderRepo)
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
	mockOrderRepo := new(mocks.MockOrderRepository)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockOrderRepo)

	input := ports.CreateOrderInput{ProductIDs: []uint{}, MaintenanceIDs: []uint{1}}

	mockProd.On("GetByIds", ctx, input.ProductIDs).Return([]domain.Product{}, nil)
	mockMaint.On("GetByIDs", ctx, input.MaintenanceIDs).Return(nil, errors.New("maint error"))

	created, err := uc.CreateOrder(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, created)
	mockProd.AssertExpectations(t)
	mockMaint.AssertExpectations(t)
}

func createMockOrder(id uint, status string) *order.Model {
	now := time.Now()
	return &order.Model{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		DateIn:            now,
		Status:            status,
		VehicleKilometers: 50000,
		UserID:            1,
		CustomerVehicleID: 1,
		CompanyID:         1,
		User: user.Model{
			ID:    1,
			Name:  "Test User",
			Email: "test@example.com",
		},
		CustomerVehicle: customer_vehicle.Model{
			Model: gorm.Model{ID: 1},
		},
		Company: company.Model{
			Model: gorm.Model{ID: 1},
			Name:  "Test Company",
		},
	}
}

func TestUseCase_ApproveOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("should approve order successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.AWAITING_APPROVAL))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(o *order.Model) bool {
			return o.ID == orderID && o.Status == string(domain.OrderStatuses.APPROVED)
		})).Return(nil)

		err := useCase.ApproveOrder(ctx, orderID)

		assert.NoError(t, err)
		assert.Equal(t, string(domain.OrderStatuses.APPROVED), existingOrder.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(999)

		mockRepo.On("FindByID", ctx, orderID).Return(nil, errors.New("not found"))

		err := useCase.ApproveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Order with Id 999 not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when order status is not AWAITING_APPROVAL", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.APPROVED))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)

		err := useCase.ApproveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Order cannot be approved. Current status:")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.AWAITING_APPROVAL))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("database error"))

		err := useCase.ApproveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to approve order:")
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_RejectOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("should reject order successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.AWAITING_APPROVAL))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(o *order.Model) bool {
			return o.ID == orderID && o.Status == string(domain.OrderStatuses.FINISHED)
		})).Return(nil)

		err := useCase.RejectOrder(ctx, orderID)

		assert.NoError(t, err)
		assert.Equal(t, string(domain.OrderStatuses.FINISHED), existingOrder.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(999)

		mockRepo.On("FindByID", ctx, orderID).Return(nil, errors.New("not found"))

		err := useCase.RejectOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Order with Id 999 not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when order status is not AWAITING_APPROVAL", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.FINISHED))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)

		err := useCase.RejectOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Order cannot be reject. Current status:")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.AWAITING_APPROVAL))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("database error"))

		err := useCase.RejectOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to reject order:")
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_ArchiveOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("should archive order successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.FINISHED))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(o *order.Model) bool {
			return o.ID == orderID && o.Status == string(domain.OrderStatuses.DELIVERED)
		})).Return(nil)

		err := useCase.ArchiveOrder(ctx, orderID)

		assert.NoError(t, err)
		assert.Equal(t, string(domain.OrderStatuses.DELIVERED), existingOrder.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(999)

		mockRepo.On("FindByID", ctx, orderID).Return(nil, errors.New("not found"))

		err := useCase.ArchiveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Order with Id 999 not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when order status is not FINISHED", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.AWAITING_APPROVAL))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)

		err := useCase.ArchiveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Order cannot be archived. Current status:")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when order status is APPROVED", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.APPROVED))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)

		err := useCase.ArchiveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Order cannot be archived. Current status:")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when order status is DELIVERED", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.DELIVERED))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)

		err := useCase.ArchiveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Order cannot be archived. Current status:")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.FINISHED))

		mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("database error"))

		err := useCase.ArchiveOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to archive order:")
		mockRepo.AssertExpectations(t)
	})
}
