package order

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

// @TODO: receber parametros necessarios para coletar as informacoes para criacao do order, deve vir do handle
//func (uc *UseCase) Create(ctx context.Context) (*domain.Order, error) {
//	var price float64
//
//	status := domain.RECEIVED
//	kilometers := 40.000
//	note := "anything else"
//	diagnosticNote := "anything else"
//
//	// @TODO: get products ids
//	productIDs := make([]uint, len([]domain.Product{}))
//	for i, item := range []domain.Product{} {
//		productIDs[i] = item.ID
//	}
//
//	productsPrice, err := uc.productService.CheckProductPrice(ctx, productIDs)
//
//	if err != nil {
//		// @TODO: adicionar log e tratamento aqui
//		return nil, err
//	}
//
//	values := make([]float64, 0, len(productsPrice))
//	for _, v := range productsPrice {
//		values = append(values, v)
//	}
//
//	price += monetary.SumPrices(values)
//
//	order := &domain.Order{
//		ID:                1,
//		DateIn:            &time.Time{},
//		Number:            "ABC-123", // @TODO trocar number para outro nome tipo order-code order..
//		VehicleKilometers: &kilometers,
//		Status:            status,
//		Note:              &note,
//		DiagnosticNote:    &diagnosticNote,
//		UserId:            1,
//		Price:             price,
//		CustomerVehicle: domain.CustomerVehicle{
//			ID:         1,
//			CustomerId: 1,
//			VehicleId:  1,
//		},
//		Company: domain.Company{
//			ID: 1,
//		},
//		User: domain.User{
//			ID: 1,
//		},
//	}
//
//	return order, nil
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
