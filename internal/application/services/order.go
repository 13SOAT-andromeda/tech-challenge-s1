package services

import (
	"context"
	"strconv"
	"strings"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/email"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/errors"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/monetary"
)

var (
	ErrOrderIdInvalid = &errors.ValidationError{Message: "Order Id invalid"}
	ErrOrderNotFound  = &errors.ValidationError{Message: "Order not found"}
	ErrOrderDelete    = &errors.ValidationError{Message: "An error occurred while trying to delete the order"}
)

type OrderService struct {
	repo ports.OrderRepository
}

func NewOrderService(repo ports.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, o domain.Order) (*domain.Order, error) {
	model := order.Model{}
	model.FromDomain(&o)

	_, err := s.repo.Create(ctx, &model)

	if err != nil {
		return nil, err
	}

	created := model.ToDomain()

	return created, nil
}

func (s *OrderService) GetByID(ctx context.Context, id uint) (*domain.Order, error) {
	result, err := s.repo.FindOrderByID(ctx, id)

	if err != nil {
		return nil, err
	}
	o := result.ToDomain()

	return o, nil
}

func (s *OrderService) GetAll(ctx context.Context, params map[string]interface{}) (*[]domain.Order, error) {
	oSearch := ports.OrderSearch{Status: "", Enabled: true, OrderBy: "", SortDesc: false}

	if params["status"] != nil {
		oSearch.Status = params["status"].(string)
	}

	if params["enabled"] != nil {
		if enabledStr, ok := params["enabled"].(string); ok {
			oSearch.Enabled = enabledStr == "true" || enabledStr == "1"
		}
	}

	if params["orderby"] != nil {
		oSearch.OrderBy = params["orderby"].(string)
	}

	if params["sortdesc"] != nil {
		if sortDescStr, ok := params["sortdesc"].(string); ok {
			oSearch.SortDesc = sortDescStr == "true" || sortDescStr == "1"
		}
	}

	orders, err := s.repo.Search(ctx, oSearch)

	if err != nil {
		return nil, err
	}

	ordersD := make([]domain.Order, 0, len(orders))

	for _, item := range orders {
		ordersD = append(ordersD, *item.ToDomain())
	}

	return &ordersD, nil
}

func (s *OrderService) Delete(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return ErrOrderDelete
	}

	return nil
}

func (s *OrderService) Update(ctx context.Context, o domain.Order) error {
	model := order.Model{}
	model.FromDomain(&o)

	err := s.repo.Update(ctx, &model)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetApprovalTemplate(order domain.Order, customer domain.Customer, apiUrl string) (string, error) {
	html, err := email.LoadTemplate(email.ORDER_APPROVAL)

	if err != nil {
		return "", err
	}

	dateIn := order.DateIn.Format("02/01/2006")
	price := ""
	diagnosticNote := ""
	note := ""
	id := strconv.FormatUint(uint64(order.ID), 8)

	if order.DiagnosticNote != nil {
		diagnosticNote = *order.DiagnosticNote
	}

	if order.Price != nil {
		price = monetary.FormatPtBrCurrency(*order.Price)
	}

	if order.Note != nil {
		note = *order.Note
	}

	html = strings.ReplaceAll(html, "$Name", customer.Name)
	html = strings.ReplaceAll(html, "$ID", id)
	html = strings.ReplaceAll(html, "$DateIn", dateIn)
	html = strings.ReplaceAll(html, "$DiagnosticNote", diagnosticNote)
	html = strings.ReplaceAll(html, "$Value", price)
	html = strings.ReplaceAll(html, "$Note", note)
	html = strings.ReplaceAll(html, "$Approval_url", apiUrl+"/orders/"+id+"/approve")
	html = strings.ReplaceAll(html, "$Repprove_url", apiUrl+"/orders/"+id+"/reject")

	return html, nil
}
