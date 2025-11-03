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
	orderRepository    ports.OrderRepository
}

func NewOrderUseCase(orderService ports.OrderService, productsService ports.ProductService, maintenanceService ports.MaintenanceService, orderRepository ports.OrderRepository) *UseCase {
	return &UseCase{
		orderService:       orderService,
		productService:     productsService,
		maintenanceService: maintenanceService,
		orderRepository:    orderRepository,
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

	products, err := uc.productService.GetByIds(ctx, input.ProductIDs)
	if err != nil {
		return err
	}

	maintenances, err := uc.maintenanceService.GetByIDs(ctx, input.MaintenanceIDs)
	if err != nil {
		return err
	}

	totalPrice := 0.0

	for _, v := range products {
		totalPrice += float64(v.Price)
	}

	for _, v := range maintenances {
		totalPrice += float64(v.Price)
	}

	order.DiagnosticNote = input.DiagnosticNote
	order.Status = domain.OrderStatuses.AWAITING_APPROVAL
	order.Price = &totalPrice
	order.User.ID = userID

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

	existentOrder.Status = string(domain.OrderStatuses.APPROVED)

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

	existentOrder.Status = string(domain.OrderStatuses.FINISHED)

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

//
//import (
//	"errors"
//
//	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
//)
//
//var (
//	ErrOrderNotFound         = errors.New("Order Not found")
//	ErrOrderAlreadyApproved  = errors.New("Order Already Approved")
//	ErrOrderAlreadyCancelled = errors.New("Order already cancelled")
//	ErrInvalidOrderStatus    = errors.New("Order Status invalid")
//)
//
//type UseCase struct {
//	productService ports.ProductService
//	orderService   ports.OrderService
//}
//
//func NewUseCase(productService ports.ProductService) *UseCase {
//	return &UseCase{productService: productService}
//}

//func (uc *UseCase) ProcessPurchase(ctx context.Context, productID uint, quantity uint) error {
//	if quantity == 0 {
//		return ErrInvalidQuantity
//	}
//
//	product, err := uc.productService.GetById(ctx, productID)
//	if err != nil {
//		return err
//	}
//
//	if product == nil {
//		return ErrProductNotFound
//	}
//
//	if err := product.CanBePurchased(quantity); err != nil {
//		return err
//	}
//
//	if err := product.DecreaseStock(quantity); err != nil {
//		return err
//	}
//
//	_, err = uc.productService.Update(ctx, *product)
//	return err
//}

//func (uc *UseCase) CalculateOrderTotal(ctx context.Context, orderID uint) (float64, error) {
//	order, err := uc.orderService.GetById(ctx, orderID)
//	if err != nil {
//		return 0, err
//	}
//
//	if order == nil {
//		return 0, ErrOrderNotFound
//	}
//
//	var total float64
//	for _, item := range order.Items {
//		product, err := uc.productService.GetById(ctx, item.ProductID)
//		if err != nil {
//			return 0, err
//		}
//
//		itemTotal := float64(product.Price) * float64(item.Quantity) / 100.0
//		total += itemTotal
//	}
//
//	return total, nil
//}
