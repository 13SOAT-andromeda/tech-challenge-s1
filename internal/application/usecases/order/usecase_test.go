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
	mockCust := new(mocks.MockCustomerService)
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockEmail := new(mocks.MockEmail)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, "")

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
	mockCust := new(mocks.MockCustomerService)
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockEmail := new(mocks.MockEmail)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, "")
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
	mockCust := new(mocks.MockCustomerService)
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockEmail := new(mocks.MockEmail)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, "")

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

func TestUseCase_AssignOrder(t *testing.T) {
	ctx := context.Background()
	t.Run("should assign order successfully", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		useCase := &UseCase{orderService: mockOrderService}

		orderID := uint(1)
		userID := uint(2)
		existingOrder := &domain.Order{
			ID:     orderID,
			Status: domain.OrderStatuses.RECEIVED,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil)
		mockOrderService.On("Update", ctx, mock.MatchedBy(func(o domain.Order) bool {
			return o.ID == orderID && o.User.ID == userID && o.Status == domain.OrderStatuses.IN_ANALYSIS
		})).Return(nil)

		err := useCase.AssignOrder(ctx, orderID, userID)

		assert.NoError(t, err)
		assert.Equal(t, userID, existingOrder.User.ID)
		assert.Equal(t, domain.OrderStatuses.IN_ANALYSIS, existingOrder.Status)
		mockOrderService.AssertExpectations(t)
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		useCase := &UseCase{orderService: mockOrderService}

		orderID := uint(999)
		userID := uint(2)

		mockOrderService.On("GetByID", ctx, orderID).Return(nil, domain.ErrOrderNotFound)

		err := useCase.AssignOrder(ctx, orderID, userID)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrOrderNotFound, err)
		mockOrderService.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		useCase := &UseCase{orderService: mockOrderService}

		orderID := uint(1)
		userID := uint(2)
		existingOrder := &domain.Order{
			ID:     orderID,
			Status: domain.OrderStatuses.RECEIVED,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil)
		mockOrderService.On("Update", ctx, mock.Anything).Return(errors.New("update error"))

		err := useCase.AssignOrder(ctx, orderID, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update error")
		mockOrderService.AssertExpectations(t)
	})
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
		assert.Contains(t, err.Error(), "order with Id 999 not found")
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
		assert.Contains(t, err.Error(), "order cannot be approved. Current status:")
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
		assert.Contains(t, err.Error(), "failed to approve order:")
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
		assert.Contains(t, err.Error(), "order with Id 999 not found")
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
		assert.Contains(t, err.Error(), "order cannot be reject. Current status:")
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
		assert.Contains(t, err.Error(), "failed to reject order:")
		mockRepo.AssertExpectations(t)
	})
}

func TestUseCase_RequestApproval(t *testing.T) {
	ctx := context.Background()

	t.Run("should request approval successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		mockOrderService := new(mocks.MockOrderService)
		mockCustomerService := new(mocks.MockCustomerService)
		mockEmail := new(mocks.MockEmail)

		useCase := &UseCase{
			orderRepository: mockRepo,
			orderService:    mockOrderService,
			customerService: mockCustomerService,
			emailService:    mockEmail,
			apiUrl:          "http://localhost:8080",
		}

		orderID := uint(1)
		customerID := uint(10)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.ANALYSIS_FINISHED))
		existingOrder.CustomerVehicle.CustomerId = customerID

		customer := &domain.Customer{
			ID:    customerID,
			Name:  "Test Customer",
			Email: "customer@test.com",
		}

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.MatchedBy(func(o *order.Model) bool {
			return o.ID == orderID && o.Status == string(domain.OrderStatuses.AWAITING_APPROVAL)
		})).Return(nil)
		mockCustomerService.On("GetByID", ctx, customerID).Return(customer, nil)
		mockOrderService.On("GetApprovalTemplate", mock.Anything, *customer, "http://localhost:8080").Return("<h1>template</h1>", nil)
		mockEmail.On("Send", customer.Name, customer.Email, "Aprovação de Ordem de Serviço", "<h1>template</h1>").Return(nil)

		err := useCase.RequestApproval(ctx, orderID)

		assert.NoError(t, err)
		assert.Equal(t, string(domain.OrderStatuses.AWAITING_APPROVAL), existingOrder.Status)
		mockRepo.AssertExpectations(t)
		mockCustomerService.AssertExpectations(t)
		mockOrderService.AssertExpectations(t)
		mockEmail.AssertExpectations(t)
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(999)

		mockRepo.On("FindOrderByID", ctx, orderID).Return(nil, errors.New("not found"))

		err := useCase.RequestApproval(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order with Id 999 not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when order status is not ANALYSIS_FINISHED", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.RECEIVED))

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)

		err := useCase.RequestApproval(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "notification cannot be sent")
		assert.Contains(t, err.Error(), "Current status: Recebida")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		useCase := &UseCase{orderRepository: mockRepo}

		orderID := uint(1)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.ANALYSIS_FINISHED))

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("database error"))

		err := useCase.RequestApproval(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update order status:")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when customer not found", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		mockCustomerService := new(mocks.MockCustomerService)
		useCase := &UseCase{
			orderRepository: mockRepo,
			customerService: mockCustomerService,
		}

		orderID := uint(1)
		customerID := uint(10)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.ANALYSIS_FINISHED))
		existingOrder.CustomerVehicle.CustomerId = customerID

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(nil)
		mockCustomerService.On("GetByID", ctx, customerID).Return(nil, errors.New("customer not found"))

		err := useCase.RequestApproval(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error on find order's customer:")
		mockRepo.AssertExpectations(t)
		mockCustomerService.AssertExpectations(t)
	})

	t.Run("should return error when template generation fails", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		mockOrderService := new(mocks.MockOrderService)
		mockCustomerService := new(mocks.MockCustomerService)
		useCase := &UseCase{
			orderRepository: mockRepo,
			orderService:    mockOrderService,
			customerService: mockCustomerService,
			apiUrl:          "http://localhost:8080",
		}

		orderID := uint(1)
		customerID := uint(10)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.ANALYSIS_FINISHED))
		existingOrder.CustomerVehicle.CustomerId = customerID

		customer := &domain.Customer{
			ID:    customerID,
			Name:  "Test Customer",
			Email: "customer@test.com",
		}

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(nil)
		mockCustomerService.On("GetByID", ctx, customerID).Return(customer, nil)
		mockOrderService.On("GetApprovalTemplate", mock.Anything, *customer, "http://localhost:8080").Return("", errors.New("template error"))

		err := useCase.RequestApproval(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse mail template:")
		mockRepo.AssertExpectations(t)
		mockCustomerService.AssertExpectations(t)
		mockOrderService.AssertExpectations(t)
	})

	t.Run("should return error when email send fails", func(t *testing.T) {
		mockRepo := new(mocks.MockOrderRepository)
		mockOrderService := new(mocks.MockOrderService)
		mockCustomerService := new(mocks.MockCustomerService)
		mockEmail := new(mocks.MockEmail)
		useCase := &UseCase{
			orderRepository: mockRepo,
			orderService:    mockOrderService,
			customerService: mockCustomerService,
			emailService:    mockEmail,
			apiUrl:          "http://localhost:8080",
		}

		orderID := uint(1)
		customerID := uint(10)
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.ANALYSIS_FINISHED))
		existingOrder.CustomerVehicle.CustomerId = customerID

		customer := &domain.Customer{
			ID:    customerID,
			Name:  "Test Customer",
			Email: "customer@test.com",
		}

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(nil)
		mockCustomerService.On("GetByID", ctx, customerID).Return(customer, nil)
		mockOrderService.On("GetApprovalTemplate", mock.Anything, *customer, "http://localhost:8080").Return("<h1>template</h1>", nil)
		mockEmail.On("Send", customer.Name, customer.Email, "Aprovação de Ordem de Serviço", "<h1>template</h1>").Return(errors.New("email error"))

		err := useCase.RequestApproval(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to send approval notification:")
		mockRepo.AssertExpectations(t)
		mockCustomerService.AssertExpectations(t)
		mockOrderService.AssertExpectations(t)
		mockEmail.AssertExpectations(t)
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
		assert.Contains(t, err.Error(), "order with Id 999 not found")
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
		assert.Contains(t, err.Error(), "order cannot be archived. Current status:")
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
		assert.Contains(t, err.Error(), "order cannot be archived. Current status:")
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
		assert.Contains(t, err.Error(), "order cannot be archived. Current status:")
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
		assert.Contains(t, err.Error(), "failed to archive order:")
		mockRepo.AssertExpectations(t)
	})
}
