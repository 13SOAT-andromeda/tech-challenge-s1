package order

import (
	"context"
	"fmt"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
)

type UseCase struct {
	orderService       ports.OrderService
	productService     ports.ProductService
	maintenanceService ports.MaintenanceService
	customerService    ports.CustomerService
	emailService       ports.Email
	orderRepository    ports.OrderRepository
	apiUrl             string
}

func NewOrderUseCase(
	orderService ports.OrderService,
	productsService ports.ProductService,
	maintenanceService ports.MaintenanceService,
	customerService ports.CustomerService,
	emailService ports.Email,
	orderRepository ports.OrderRepository,
	apiUrl string,
) *UseCase {
	return &UseCase{
		orderService:       orderService,
		productService:     productsService,
		maintenanceService: maintenanceService,
		customerService:    customerService,
		emailService:       emailService,
		orderRepository:    orderRepository,
		apiUrl:             apiUrl,
	}
}

func (uc *UseCase) CreateOrder(ctx context.Context, input ports.CreateOrderInput) (*domain.Order, error) {
	products, err := uc.productService.GetByIds(ctx, input.ProductIDs)
	if err != nil {
		return nil, err
	}

	maintenances, err := uc.maintenanceService.GetByIDs(ctx, input.MaintenanceIDs)
	if err != nil {
		return nil, err
	}

	totalPrice := 0.0

	for _, v := range products {
		totalPrice += float64(v.Price)
	}

	for _, v := range maintenances {
		totalPrice += float64(v.Price)
	}

	order := domain.Order{
		DateIn:            time.Now(),
		DateOut:           nil,
		Status:            domain.OrderStatuses.RECEIVED,
		VehicleKilometers: input.VehicleKilometers,
		Note:              input.Note,
		Price:             &totalPrice,
		CustomerVehicle:   domain.CustomerVehicle{ID: input.CustomerVehicleID},
		User:              domain.User{ID: input.UserID},
		Company:           domain.Company{ID: input.CompanyID},
	}

	created, err := uc.orderService.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (uc *UseCase) AssignOrder(ctx context.Context, orderID uint, userID uint) error {
	order, err := uc.orderService.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order == nil {
		return domain.ErrOrderNotFound
	}

	order.User = domain.User{ID: userID}
	order.Status = domain.OrderStatuses.IN_ANALYSIS

	err = uc.orderService.Update(ctx, *order)
	return err
}

func (s *UseCase) ApproveOrder(ctx context.Context, id uint) error {

	existentOrder, err := s.orderRepository.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.AWAITING_APPROVAL {
		return fmt.Errorf("order cannot be approved. Current status: %s", existentOrder.Status)
	}

	existentOrder.Status = string(domain.OrderStatuses.APPROVED)

	if err := s.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("failed to approve order: %w", err)
	}

	return nil
}

func (s *UseCase) RejectOrder(ctx context.Context, id uint) error {

	existentOrder, err := s.orderRepository.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.AWAITING_APPROVAL {
		return fmt.Errorf("order cannot be reject. Current status: %s", existentOrder.Status)
	}

	existentOrder.Status = string(domain.OrderStatuses.FINISHED)

	if err := s.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("failed to reject order: %w", err)
	}

	return nil
}

func (s *UseCase) ArchiveOrder(ctx context.Context, id uint) error {

	existentOrder, err := s.orderRepository.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.FINISHED {
		return fmt.Errorf("order cannot be archived. Current status: %s", existentOrder.Status)
	}

	existentOrder.Status = string(domain.OrderStatuses.DELIVERED)

	if err := s.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("failed to archive order: %w", err)
	}

	return nil
}

func (s *UseCase) RequestApproval(ctx context.Context, id uint) error {
	existentOrder, err := s.orderRepository.FindOrderByID(ctx, id)

	if err != nil {
		return fmt.Errorf("order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.ANALYSIS_FINISHED {
		return fmt.Errorf("notification cannot be sent. Order should be in %s status. Current status: %s", domain.OrderStatuses.ANALYSIS_FINISHED, existentOrder.Status)
	}

	existentOrder.Status = string(domain.OrderStatuses.AWAITING_APPROVAL)

	if err := s.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	c, err := s.customerService.GetByID(ctx, existentOrder.CustomerVehicle.CustomerId)

	if err != nil {
		return fmt.Errorf("error on find order's customer: %w", err)
	}

	html, err := s.orderService.GetApprovalTemplate(*existentOrder.ToDomain(), *c, s.apiUrl)

	if err != nil {
		return fmt.Errorf("failed to parse mail template: %w", err)
	}

	err = s.emailService.Send(c.Name, c.Email, "Aprovação de Ordem de Serviço", html)

	if err != nil {
		return fmt.Errorf("failed to send approval notification: %w", err)
	}

	return nil
}
