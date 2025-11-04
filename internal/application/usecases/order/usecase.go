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

func (uc *UseCase) CreateOrder(ctx context.Context, userID uint, input ports.CreateOrderInput) (*domain.Order, error) {
	order := domain.Order{
		DateIn:            time.Now(),
		DateOut:           nil,
		Status:            domain.OrderStatuses.RECEIVED,
		VehicleKilometers: input.VehicleKilometers,
		Note:              input.Note,
		Price:             nil,
		CustomerVehicle:   domain.CustomerVehicle{ID: input.CustomerVehicleID},
		User:              domain.User{ID: userID},
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

func (uc *UseCase) CompleteOrderAnalysis(ctx context.Context, id uint, userID uint, input ports.CreateCompleteOrderAnalysisInput) error {
	order, err := uc.orderService.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("order with Id %d not found", id)
	}

	if order.Status != domain.OrderStatuses.IN_ANALYSIS {
		return fmt.Errorf("order cannot complete analysis. Current status: %s", order.Status)
	}

	productIds := make([]uint, 0, len(input.Products))
	productQuantities := make(map[uint]int, len(input.Products))

	for _, item := range input.Products {
		productIds = append(productIds, item.ID)
		productQuantities[item.ID] = int(item.Quantity)
	}

	products, err := uc.productService.GetByIds(ctx, productIds)
	if err != nil {
		return err
	}

	maintenanceIds := make([]uint, 0, len(input.Maintenances))
	for _, v := range input.Maintenances {
		maintenanceIds = append(maintenanceIds, v.ID)
	}

	maintenances, err := uc.maintenanceService.GetByIDs(ctx, maintenanceIds)
	if err != nil {
		return err
	}

	totalPrice := 0.0

	orderProducts := make([]domain.ProductItem, 0, len(products))
	for _, product := range products {
		quantity := productQuantities[product.ID]
		totalPrice += float64(product.Price) * float64(quantity)

		orderProducts = append(orderProducts, domain.ProductItem{
			ID:       product.ID,
			Quantity: uint(quantity),
		})
	}

	orderMaintenances := make([]domain.MaintenanceItem, 0, len(maintenances))
	for _, maintenance := range maintenances {
		totalPrice += float64(maintenance.Price)

		orderMaintenances = append(orderMaintenances, domain.MaintenanceItem{
			ID: maintenance.ID,
		})
	}

	order.DiagnosticNote = input.DiagnosticNote
	order.Status = domain.OrderStatuses.AWAITING_APPROVAL
	order.Price = &totalPrice
	order.User.ID = userID
	order.Products = &orderProducts
	order.Maintenances = &orderMaintenances

	if err := uc.orderService.Update(ctx, *order); err != nil {
		return fmt.Errorf("failed to complete order analysis: %w", err)
	}

	return nil
}

func (uc *UseCase) ApproveOrder(ctx context.Context, id uint) error {

	existentOrder, err := uc.orderRepository.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.AWAITING_APPROVAL {
		return fmt.Errorf("order cannot be approved. Current status: %s", existentOrder.Status)
	}

	now := time.Now()

	existentOrder.Status = string(domain.OrderStatuses.APPROVED)
	existentOrder.DateApproved = &now

	if err := uc.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("failed to approve order: %w", err)
	}

	return nil
}

func (uc *UseCase) RejectOrder(ctx context.Context, id uint) error {

	existentOrder, err := uc.orderRepository.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.AWAITING_APPROVAL {
		return fmt.Errorf("order cannot be reject. Current status: %s", existentOrder.Status)
	}

	now := time.Now()

	existentOrder.Status = string(domain.OrderStatuses.FINISHED)
	existentOrder.DateRejected = &now

	if err := uc.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("failed to reject order: %w", err)
	}

	return nil
}

func (uc *UseCase) ArchiveOrder(ctx context.Context, id uint) error {

	existentOrder, err := uc.orderRepository.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.FINISHED {
		return fmt.Errorf("order cannot be archived. Current status: %s", existentOrder.Status)
	}

	existentOrder.Status = string(domain.OrderStatuses.DELIVERED)

	if err := uc.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("failed to archive order: %w", err)
	}

	return nil
}

func (uc *UseCase) RequestApproval(ctx context.Context, id uint) error {
	existentOrder, err := uc.orderRepository.FindOrderByID(ctx, id)

	if err != nil {
		return fmt.Errorf("order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.ANALYSIS_FINISHED {
		return fmt.Errorf("notification cannot be sent. Order should be in %s status. Current status: %s", domain.OrderStatuses.ANALYSIS_FINISHED, existentOrder.Status)
	}

	existentOrder.Status = string(domain.OrderStatuses.AWAITING_APPROVAL)

	if err := uc.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	c, err := uc.customerService.GetByID(ctx, existentOrder.CustomerVehicle.CustomerId)

	if err != nil {
		return fmt.Errorf("error on find order's customer: %w", err)
	}

	html, err := uc.orderService.GetApprovalTemplate(*existentOrder.ToDomain(), *c, uc.apiUrl)

	if err != nil {
		return fmt.Errorf("failed to parse mail template: %w", err)
	}

	err = uc.emailService.Send(c.Name, c.Email, "Aprovação de Ordem de Serviço", html)

	if err != nil {
		return fmt.Errorf("failed to send approval notification: %w", err)
	}

	return nil
}

func (uc *UseCase) StartWorkOrder(ctx context.Context, id uint) error {
	order, err := uc.orderService.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order with id %d: %w", id, err)
	}

	if order.Status != domain.OrderStatuses.APPROVED {
		return fmt.Errorf("order cannot start work. current status: %s", order.Status)
	}

	productItems := make([]domain.ProductItem, 0, len(*order.Products))

	for _, item := range *order.Products {
		available, err := uc.productService.CheckAvailability(ctx, item.ID, item.Quantity)
		if err != nil {
			return fmt.Errorf("failed to check availability for product %d: %w", item.ID, err)
		}

		if !available {
			return fmt.Errorf("cannot start work: product ID %d is not available in quantity %d", item.ID, item.Quantity)
		}

		productItems = append(productItems, domain.ProductItem{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}

	err = uc.productService.UpdateStock(ctx, productItems, domain.StockOperationRemove)
	if err != nil {
		return fmt.Errorf("failed to decrement stock: %w", err)
	}

	order.Status = domain.OrderStatuses.IN_PROGRESS

	if err = uc.orderService.Update(ctx, *order); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}
