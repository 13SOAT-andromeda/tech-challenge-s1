package order

import (
	"context"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/monetary"
)

type UseCase struct {
	productService ports.ProductService
}

func New(productService ports.ProductService) *UseCase {
	return &UseCase{productService: productService}
}

// @TODO: receber parametros necessarios para coletar as informacoes para criacao do order, deve vir do handle
func (uc *UseCase) Create(ctx context.Context) (*domain.Order, error) {
	var price float64

	status := domain.RECEIVED
	kilometers := 40.000
	note := "anything else"
	diagnosticNote := "anything else"

	// @TODO: get products ids
	productIDs := make([]uint, len([]domain.Product{}))
	for i, item := range []domain.Product{} {
		productIDs[i] = item.ID
	}

	productsPrice, err := uc.productService.CheckProductPrice(ctx, productIDs)
	if err != nil {
		// @TODO: adicionar log e tratamento aqui
		return nil, err
	}

	price += monetary.SumPrices(productsPrice)

	order := &domain.Order{
		ID:                1,
		DateIn:            &time.Time{},
		Number:            "ABC-123", // @TODO trocar number para outro nome tipo order-code order..
		VehicleKilometers: &kilometers,
		Status:            status,
		Note:              &note,
		DiagnosticNote:    &diagnosticNote,
		UserId:            1,
		Price:             &price,
		CustomerVehicle: domain.CustomerVehicle{
			ID:         1,
			CustomerId: 1,
			VehicleId:  1,
		},
		Company: domain.Company{
			ID: 1,
		},
		User: domain.User{
			ID: 1,
		},
	}

}
