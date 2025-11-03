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

func (s *UseCase) ApproveOrder(ctx context.Context, id uint) error {

	existentOrder, err := s.orderRepository.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("Order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.AWAITING_APPROVAL {
		return fmt.Errorf("Order cannot be approved. Current status: %s", existentOrder.Status)
	}

	existentOrder.Status = string(domain.OrderStatuses.APPROVED)

	if err := s.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("Failed to approve order: %w", err)
	}

	return nil
}

func (s *UseCase) RejectOrder(ctx context.Context, id uint) error {

	existentOrder, err := s.orderRepository.FindByID(ctx, id)

	if err != nil {
		return fmt.Errorf("Order with Id %d not found", id)
	}

	if domain.OrderStatus(existentOrder.Status) != domain.OrderStatuses.AWAITING_APPROVAL {
		return fmt.Errorf("Order cannot be reject. Current status: %s", existentOrder.Status)
	}

	existentOrder.Status = string(domain.OrderStatuses.FINISHED)

	if err := s.orderRepository.Update(ctx, existentOrder); err != nil {
		return fmt.Errorf("Failed to reject order: %w", err)
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
