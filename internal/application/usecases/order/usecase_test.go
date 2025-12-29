package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	orderMaintenanceModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_maintenance"
	orderProductModel "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
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
	mockCust := new(mocks.MockCustomerService)
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockOrderProd := new(mocks.MockOrderProductRepository)
	mockOrderMaint := new(mocks.MockOrderMaintenanceRepository)
	mockEmail := new(mocks.MockEmail)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, mockOrderProd, mockOrderMaint, "")

	input := ports.CreateOrderInput{
		VehicleKilometers: 123,
		Note:              ptrString("note"),
		UserID:            10,
		CustomerVehicleID: 20,
		CompanyID:         30,
	}

	mockOrder.On("Create", mock.Anything, mock.MatchedBy(func(o domain.Order) bool {
		if o.VehicleKilometers != input.VehicleKilometers {
			return false
		}
		if o.CustomerVehicleID != input.CustomerVehicleID {
			return false
		}
		if o.CompanyID != input.CompanyID {
			return false
		}
		if o.Note == nil || *o.Note != *input.Note {
			return false
		}
		return true
	})).Return(&domain.Order{ID: 1}, nil)

	userId := uint(1)
	created, err := uc.CreateOrder(ctx, userId, input)

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
	mockOrderProd := new(mocks.MockOrderProductRepository)
	mockOrderMaint := new(mocks.MockOrderMaintenanceRepository)
	input := ports.CreateCompleteOrderAnalysisInput{}

	userID := uint(1)
	orderID := uint(1)
	mockEmail := new(mocks.MockEmail)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, mockOrderProd, mockOrderMaint, "")

	products := []domain.Product{p(1, 100), p(2, 150)}
	maints := []domain.Maintenance{ma(3, 250)}
	existingOrder := &domain.Order{
		ID:     orderID,
		Status: domain.OrderStatuses.RECEIVED,
	}

	mockOrder.On("GetByID", ctx, orderID).Return(existingOrder, nil)
	mockProd.On("GetByIds", ctx, input.Products).Return(products, nil)
	mockMaint.On("GetByIDs", ctx, input.Maintenances).Return(maints, nil)

	err := uc.CompleteOrderAnalysis(ctx, orderID, userID, input)

	assert.Error(t, err)
}

func TestCreateOrder_MaintenanceServiceError(t *testing.T) {
	ctx := context.Background()
	mockOrder := new(mocks.MockOrderService)
	mockProd := new(mocks.MockProductService)
	mockMaint := new(mocks.MockMaintenanceService)
	mockCust := new(mocks.MockCustomerService)
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockOrderProd := new(mocks.MockOrderProductRepository)
	mockOrderMaint := new(mocks.MockOrderMaintenanceRepository)
	mockEmail := new(mocks.MockEmail)
	uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, mockOrderProd, mockOrderMaint, "")

	input := ports.CreateCompleteOrderAnalysisInput{}

	userID := uint(1)
	orderID := uint(1)

	products := []domain.Product{p(1, 100), p(2, 150)}
	maints := []domain.Maintenance{ma(3, 250)}
	existingOrder := &domain.Order{
		ID:     orderID,
		Status: domain.OrderStatuses.RECEIVED,
	}

	mockOrder.On("GetByID", ctx, orderID).Return(existingOrder, nil)
	mockProd.On("GetByIds", ctx, input.Products).Return(products, nil)
	mockMaint.On("GetByIDs", ctx, input.Maintenances).Return(maints, nil)

	err := uc.CompleteOrderAnalysis(ctx, orderID, userID, input)

	assert.Error(t, err)
}

func createMockOrder(id uint, status string) *order.Model {
	now := time.Now()
	m := &order.Model{}
	m.ID = id
	m.CreatedAt = now
	m.UpdatedAt = now
	m.DateIn = now
	m.Status = status
	m.VehicleKilometers = 50000
	m.UserID = 1
	m.CustomerVehicleID = 1
	m.CompanyID = 1

	// User
	m.User = user.Model{}
	m.User.ID = 1
	m.User.Name = "Test User"
	m.User.Email = "test@example.com"

	// CustomerVehicle
	m.CustomerVehicle = customer_vehicle.Model{}
	m.CustomerVehicle.ID = 1
	m.CustomerVehicle.CustomerID = 1

	// Company
	m.Company = company.Model{}
	m.Company.ID = 1
	m.Company.Name = "Test Company"

	return m
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
			return o.ID == orderID && o.UserID == userID && o.Status == domain.OrderStatuses.IN_ANALYSIS
		})).Return(nil)

		err := useCase.AssignOrder(ctx, orderID, userID)

		assert.NoError(t, err)
		assert.Equal(t, userID, existingOrder.UserID)
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
		existingOrder.CustomerVehicle.CustomerID = customerID

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
		mockOrderService.On("GetApprovalTemplate", mock.Anything, mock.Anything, "http://localhost:8080").Return("<h1>template</h1>", nil)
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
		existingOrder := createMockOrder(orderID, string(domain.OrderStatuses.ANALYSIS_FINISHED))

		customerID := uint(1)
		customer := &domain.Customer{
			ID:    customerID,
			Name:  "Test Customer",
			Email: "customer@test.com",
		}

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)
		mockCustomerService.On("GetByID", ctx, customerID).Return(customer, nil)
		mockOrderService.On("GetApprovalTemplate", mock.Anything, mock.Anything, "http://localhost:8080").Return("<h1>template</h1>", nil)
		mockEmail.On("Send", customer.Name, customer.Email, "Aprovação de Ordem de Serviço", "<h1>template</h1>").Return(nil)
		mockRepo.On("Update", ctx, mock.Anything).Return(errors.New("database error"))

		err := useCase.RequestApproval(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update order status:")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when customer not found", func(t *testing.T) {
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
		existingOrder.CustomerVehicle.CustomerID = customerID

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)
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
		existingOrder.CustomerVehicle.CustomerID = customerID

		customer := &domain.Customer{
			ID:    customerID,
			Name:  "Test Customer",
			Email: "customer@test.com",
		}

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)
		mockCustomerService.On("GetByID", ctx, customerID).Return(customer, nil)
		mockOrderService.On("GetApprovalTemplate", mock.Anything, mock.Anything, "http://localhost:8080").Return("", errors.New("template error"))

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
		existingOrder.CustomerVehicle.CustomerID = customerID

		customer := &domain.Customer{
			ID:    customerID,
			Name:  "Test Customer",
			Email: "customer@test.com",
		}

		mockRepo.On("FindOrderByID", ctx, orderID).Return(existingOrder, nil)
		mockCustomerService.On("GetByID", ctx, customerID).Return(customer, nil)
		mockOrderService.On("GetApprovalTemplate", mock.Anything, mock.Anything, "http://localhost:8080").Return("<h1>template</h1>", nil)
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

func TestUseCase_CompleteOrderAnalysis(t *testing.T) {
	ctx := context.Background()

	t.Run("should complete analysis successfully", func(t *testing.T) {
		mockOrder := new(mocks.MockOrderService)
		mockProd := new(mocks.MockProductService)
		mockMaint := new(mocks.MockMaintenanceService)
		mockCust := new(mocks.MockCustomerService)
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockOrderProd := new(mocks.MockOrderProductRepository)
		mockOrderMaint := new(mocks.MockOrderMaintenanceRepository)
		mockEmail := new(mocks.MockEmail)
		uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, mockOrderProd, mockOrderMaint, "")

		orderID := uint(1)
		userID := uint(2)
		existingOrder := &domain.Order{ID: orderID, Status: domain.OrderStatuses.IN_ANALYSIS}

		products := []domain.Product{p(1, 100), p(2, 150)}
		maints := []domain.Maintenance{ma(3, 250)}

		input := ports.CreateCompleteOrderAnalysisInput{
			DiagnosticNote: ptrString("diagnostic"),
			Products: []domain.StockItem{
				{ID: 1, Quantity: 2},
				{ID: 2, Quantity: 1},
			},
			Maintenances: []uint{3},
		}

		productsIds := []uint{1, 2}
		maintenanceIds := []uint{3}
		mockOrder.On("GetByID", ctx, orderID).Return(existingOrder, nil)
		mockProd.On("GetByIds", ctx, productsIds).Return(products, nil)
		mockMaint.On("GetByIDs", ctx, maintenanceIds).Return(maints, nil)
		mockOrderProd.On("Create", ctx, mock.Anything).Return(&orderProductModel.Model{}, nil).Times(2)
		mockOrderMaint.On("Create", ctx, mock.Anything).Return(&orderMaintenanceModel.Model{}, nil).Times(1)
		mockOrder.On("Update", ctx, mock.MatchedBy(func(o domain.Order) bool {
			// Price should be set and status updated to AWAITING_APPROVAL and diagnostic note set
			if o.ID != orderID {
				return false
			}
			if o.Status != domain.OrderStatuses.ANALYSIS_FINISHED {
				return false
			}
			if o.DiagnosticNote == nil || *o.DiagnosticNote != *input.DiagnosticNote {
				return false
			}
			if o.UserID != userID {
				return false
			}
			if o.Price == nil {
				return false
			}
			// total price = 100 + 150 + 250 = 500
			if *o.Price != 600 {
				return false
			}
			return true
		})).Return(nil)

		err := uc.CompleteOrderAnalysis(ctx, orderID, userID, input)
		assert.NoError(t, err)

		mockOrder.AssertExpectations(t)
		mockProd.AssertExpectations(t)
		mockMaint.AssertExpectations(t)
	})

	t.Run("should return error when GetByID fails", func(t *testing.T) {
		mockOrder := new(mocks.MockOrderService)
		mockProd := new(mocks.MockProductService)
		mockMaint := new(mocks.MockMaintenanceService)
		mockCust := new(mocks.MockCustomerService)
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockOrderProd := new(mocks.MockOrderProductRepository)
		mockOrderMaint := new(mocks.MockOrderMaintenanceRepository)
		mockEmail := new(mocks.MockEmail)
		uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, mockOrderProd, mockOrderMaint, "")
		orderID := uint(1)
		userID := uint(2)
		input := ports.CreateCompleteOrderAnalysisInput{}

		mockOrder.On("GetByID", ctx, orderID).Return(nil, errors.New("not found"))

		err := uc.CompleteOrderAnalysis(ctx, orderID, userID, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order with Id 1 not found")

		mockOrder.AssertExpectations(t)
	})

	t.Run("should return error when order status is not IN_ANALYSIS", func(t *testing.T) {
		mockOrder := new(mocks.MockOrderService)
		mockProd := new(mocks.MockProductService)
		mockMaint := new(mocks.MockMaintenanceService)
		mockCust := new(mocks.MockCustomerService)
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockOrderProd := new(mocks.MockOrderProductRepository)
		mockOrderMaint := new(mocks.MockOrderMaintenanceRepository)
		mockEmail := new(mocks.MockEmail)
		uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, mockOrderProd, mockOrderMaint, "")
		orderID := uint(1)
		userID := uint(2)
		existingOrder := &domain.Order{ID: orderID, Status: domain.OrderStatuses.RECEIVED}
		input := ports.CreateCompleteOrderAnalysisInput{}

		mockOrder.On("GetByID", ctx, orderID).Return(existingOrder, nil)

		err := uc.CompleteOrderAnalysis(ctx, orderID, userID, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order cannot complete analysis")

		mockOrder.AssertExpectations(t)
	})

	t.Run("should return error when ProductService.GetByIds fails", func(t *testing.T) {
		mockOrder := new(mocks.MockOrderService)
		mockProd := new(mocks.MockProductService)
		mockMaint := new(mocks.MockMaintenanceService)
		mockCust := new(mocks.MockCustomerService)
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockOrderProd := new(mocks.MockOrderProductRepository)
		mockOrderMaint := new(mocks.MockOrderMaintenanceRepository)
		mockEmail := new(mocks.MockEmail)
		uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, mockOrderProd, mockOrderMaint, "")
		orderID := uint(1)
		userID := uint(2)
		existingOrder := &domain.Order{ID: orderID, Status: domain.OrderStatuses.IN_ANALYSIS}
		input := ports.CreateCompleteOrderAnalysisInput{Products: []domain.StockItem{
			{ID: 1, Quantity: 2},
		}}

		productIds := []uint{1}
		mockOrder.On("GetByID", ctx, orderID).Return(existingOrder, nil)
		mockProd.On("GetByIds", ctx, productIds).Return(nil, errors.New("product error"))

		err := uc.CompleteOrderAnalysis(ctx, orderID, userID, input)
		assert.Error(t, err)

		mockOrder.AssertExpectations(t)
		mockProd.AssertExpectations(t)
	})

	t.Run("should return error when MaintenanceService.GetByIDs fails", func(t *testing.T) {
		mockOrder := new(mocks.MockOrderService)
		mockProd := new(mocks.MockProductService)
		mockMaint := new(mocks.MockMaintenanceService)
		mockCust := new(mocks.MockCustomerService)
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockOrderProd := new(mocks.MockOrderProductRepository)
		mockOrderMaint := new(mocks.MockOrderMaintenanceRepository)
		mockEmail := new(mocks.MockEmail)
		uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, mockOrderProd, mockOrderMaint, "")
		orderID := uint(1)
		userID := uint(2)
		existingOrder := &domain.Order{ID: orderID, Status: domain.OrderStatuses.IN_ANALYSIS}
		input := ports.CreateCompleteOrderAnalysisInput{Maintenances: []uint{1}}

		maintenanceIds := []uint{1}

		mockOrder.On("GetByID", ctx, orderID).Return(existingOrder, nil)
		mockProd.On("GetByIds", ctx, mock.Anything).Return([]domain.Product{}, nil)
		mockMaint.On("GetByIDs", ctx, maintenanceIds).Return(nil, errors.New("maintenance error"))

		err := uc.CompleteOrderAnalysis(ctx, orderID, userID, input)
		assert.Error(t, err)

		mockOrder.AssertExpectations(t)
		mockProd.AssertExpectations(t)
		mockMaint.AssertExpectations(t)
	})

	t.Run("should return error when Update fails", func(t *testing.T) {
		mockOrder := new(mocks.MockOrderService)
		mockProd := new(mocks.MockProductService)
		mockMaint := new(mocks.MockMaintenanceService)
		mockCust := new(mocks.MockCustomerService)
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockOrderProd := new(mocks.MockOrderProductRepository)
		mockOrderMaint := new(mocks.MockOrderMaintenanceRepository)
		mockEmail := new(mocks.MockEmail)
		uc := NewOrderUseCase(mockOrder, mockProd, mockMaint, mockCust, mockEmail, mockOrderRepo, mockOrderProd, mockOrderMaint, "")
		orderID := uint(1)
		userID := uint(2)
		existingOrder := &domain.Order{ID: orderID, Status: domain.OrderStatuses.IN_ANALYSIS}

		products := []domain.Product{p(1, 100)}
		maints := []domain.Maintenance{ma(2, 200)}

		input := ports.CreateCompleteOrderAnalysisInput{
			DiagnosticNote: ptrString("diag"),
			Products:       []domain.StockItem{{ID: 1, Quantity: 1}},
			Maintenances:   []uint{2},
		}

		productIds := []uint{1}
		maintenanceIds := []uint{2}

		mockOrder.On("GetByID", ctx, orderID).Return(existingOrder, nil)
		mockProd.On("GetByIds", ctx, productIds).Return(products, nil)
		mockMaint.On("GetByIDs", ctx, maintenanceIds).Return(maints, nil)
		mockOrderProd.On("Create", ctx, mock.Anything).Return(&orderProductModel.Model{}, nil).Times(1)
		mockOrderMaint.On("Create", ctx, mock.Anything).Return(&orderMaintenanceModel.Model{}, nil).Times(1)
		mockOrder.On("Update", ctx, mock.Anything).Return(errors.New("update error"))

		err := uc.CompleteOrderAnalysis(ctx, orderID, userID, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to complete order analysis")

		mockOrder.AssertExpectations(t)
		mockProd.AssertExpectations(t)
		mockMaint.AssertExpectations(t)
	})
}

func TestUseCase_StartWorkOrder(t *testing.T) {
	ctx := context.Background()
	orderID := uint(1)
	stockMock := uint(100)
	quantity1 := uint(2)
	quantity2 := uint(1)
	products := []domain.Product{
		{ID: 10, Name: "Product A", Price: 100, Stock: &stockMock, Quantity: &quantity1},
		{ID: 11, Name: "Product B", Price: 50, Stock: &stockMock, Quantity: &quantity2},
	}

	t.Run("should start work order successfully", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		mockProductService := new(mocks.MockProductService)
		uc := &UseCase{
			orderService:   mockOrderService,
			productService: mockProductService,
		}
		existingOrder := &domain.Order{
			ID:       orderID,
			Status:   domain.OrderStatuses.APPROVED,
			Products: &products,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil).Once()
		mockProductService.On("CheckAvailability", ctx, uint(10), uint(2)).Return(true, nil).Once()
		mockProductService.On("CheckAvailability", ctx, uint(11), uint(1)).Return(true, nil).Once()
		mockProductService.On("UpdateStock", ctx, mock.AnythingOfType("[]domain.StockItem")).Return(nil, nil).Once()
		mockOrderService.On("Update", ctx, mock.MatchedBy(func(o domain.Order) bool {
			return o.ID == orderID && o.Status == domain.OrderStatuses.IN_PROGRESS
		})).Return(nil).Once()

		err := uc.StartWorkOrder(ctx, orderID)

		assert.NoError(t, err)
		mockOrderService.AssertExpectations(t)
		mockProductService.AssertExpectations(t)
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		uc := &UseCase{orderService: mockOrderService}
		mockOrderService.On("GetByID", ctx, orderID).Return(nil, errors.New("not found")).Once()

		err := uc.StartWorkOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get order")
		mockOrderService.AssertExpectations(t)
	})

	t.Run("should return error for wrong order status", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		uc := &UseCase{orderService: mockOrderService}
		existingOrder := &domain.Order{
			ID:     orderID,
			Status: domain.OrderStatuses.IN_ANALYSIS,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil).Once()

		err := uc.StartWorkOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order cannot start work")
		mockOrderService.AssertExpectations(t)
	})

	t.Run("should return error when product is not available", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		mockProductService := new(mocks.MockProductService)
		uc := &UseCase{
			orderService:   mockOrderService,
			productService: mockProductService,
		}
		existingOrder := &domain.Order{
			ID:       orderID,
			Status:   domain.OrderStatuses.APPROVED,
			Products: &products,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil).Once()
		mockProductService.On("CheckAvailability", ctx, uint(10), uint(2)).Return(false, nil).Once()

		err := uc.StartWorkOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "is not available")
		mockOrderService.AssertExpectations(t)
		mockProductService.AssertExpectations(t)
	})

	t.Run("should return error on check availability failure", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		mockProductService := new(mocks.MockProductService)
		uc := &UseCase{
			orderService:   mockOrderService,
			productService: mockProductService,
		}
		existingOrder := &domain.Order{
			ID:       orderID,
			Status:   domain.OrderStatuses.APPROVED,
			Products: &products,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil).Once()
		mockProductService.On("CheckAvailability", ctx, uint(10), uint(2)).Return(false, errors.New("db error")).Once()

		err := uc.StartWorkOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to check availability")
		mockOrderService.AssertExpectations(t)
		mockProductService.AssertExpectations(t)
	})

	t.Run("should return error on update stock failure", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		mockProductService := new(mocks.MockProductService)
		uc := &UseCase{
			orderService:   mockOrderService,
			productService: mockProductService,
		}
		existingOrder := &domain.Order{
			ID:       orderID,
			Status:   domain.OrderStatuses.APPROVED,
			Products: &products,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil).Once()
		mockProductService.On("CheckAvailability", ctx, uint(10), uint(2)).Return(true, nil).Once()
		mockProductService.On("CheckAvailability", ctx, uint(11), uint(1)).Return(true, nil).Once()
		mockProductService.On("UpdateStock", ctx, mock.AnythingOfType("[]domain.StockItem")).Return(nil, errors.New("db error")).Once()

		err := uc.StartWorkOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decrement stock")
		mockOrderService.AssertExpectations(t)
		mockProductService.AssertExpectations(t)
	})

	t.Run("should return error on order update failure", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		mockProductService := new(mocks.MockProductService)
		uc := &UseCase{
			orderService:   mockOrderService,
			productService: mockProductService,
		}
		existingOrder := &domain.Order{
			ID:       orderID,
			Status:   domain.OrderStatuses.APPROVED,
			Products: &products,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil).Once()
		mockProductService.On("CheckAvailability", ctx, uint(10), uint(2)).Return(true, nil).Once()
		mockProductService.On("CheckAvailability", ctx, uint(11), uint(1)).Return(true, nil).Once()
		mockProductService.On("UpdateStock", ctx, mock.AnythingOfType("[]domain.StockItem")).Return(nil, nil).Once()
		mockOrderService.On("Update", ctx, mock.Anything).Return(errors.New("db error")).Once()

		err := uc.StartWorkOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update order status")
		mockOrderService.AssertExpectations(t)
		mockProductService.AssertExpectations(t)
	})
}

func TestUseCase_CompleteWorkOrder(t *testing.T) {
	ctx := context.Background()
	orderID := uint(1)

	t.Run("should complete work order successfully", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		uc := &UseCase{orderService: mockOrderService}

		existingOrder := &domain.Order{
			ID:     orderID,
			Status: domain.OrderStatuses.IN_PROGRESS,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil).Once()
		mockOrderService.On("Update", ctx, mock.MatchedBy(func(o domain.Order) bool {
			return o.ID == orderID && o.Status == domain.OrderStatuses.FINISHED
		})).Return(nil).Once()

		err := uc.CompleteWorkOrder(ctx, orderID)

		assert.NoError(t, err)
		mockOrderService.AssertExpectations(t)
	})

	t.Run("should return error when order not found", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		uc := &UseCase{orderService: mockOrderService}

		mockOrderService.On("GetByID", ctx, orderID).Return(nil, errors.New("not found")).Once()

		err := uc.CompleteWorkOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get order")
		mockOrderService.AssertExpectations(t)
	})

	t.Run("should return error for wrong order status", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		uc := &UseCase{orderService: mockOrderService}
		existingOrder := &domain.Order{
			ID:     orderID,
			Status: domain.OrderStatuses.APPROVED,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil).Once()

		err := uc.CompleteWorkOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order cannot complete work")
		mockOrderService.AssertExpectations(t)
	})

	t.Run("should return error on order update failure", func(t *testing.T) {
		mockOrderService := new(mocks.MockOrderService)
		uc := &UseCase{orderService: mockOrderService}
		existingOrder := &domain.Order{
			ID:     orderID,
			Status: domain.OrderStatuses.IN_PROGRESS,
		}

		mockOrderService.On("GetByID", ctx, orderID).Return(existingOrder, nil).Once()
		mockOrderService.On("Update", ctx, mock.Anything).Return(errors.New("db error")).Once()

		err := uc.CompleteWorkOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to update order status")
		mockOrderService.AssertExpectations(t)
	})
}
