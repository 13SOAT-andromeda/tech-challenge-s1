package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomerHandler struct {
	service ports.CustomerService
}

func NewCustomerHandler(service ports.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

type createCustomerRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Document      string `json:"document" binding:"required"`
	Type          string `json:"type" binding:"required,oneof=administrator mechanic attendant"`
	Contact       string `json:"contact" binding:"required"`
	Address       string `json:"address" binding:"required"`
	AddressNumber string `json:"address_number" binding:"required"`
	City          string `json:"city" binding:"required"`
	Neighborhood  string `json:"neighborhood" binding:"required"`
	Country       string `json:"country" binding:"required"`
	ZipCode       string `json:"zip_code" binding:"required"`
}

func (h *CustomerHandler) CreateCustomer(ctx *gin.Context) {

	var json createCustomerRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				if fieldError.Field() == "Type" && fieldError.Tag() == "oneof" {
					response.RespondError(ctx, http.StatusBadRequest, "type must be one of: administrator, mechanic, attendant'")

					return
				}
			}
		}

		response.RespondError(ctx, http.StatusBadRequest, err.Error())

		return
	}

	c := domain.Customer{
		Name:     json.Name,
		Email:    json.Email,
		Document: json.Document,
		Type:     json.Type,
		Contact:  json.Contact,
		Address: &domain.Address{
			Address:       json.Address,
			AddressNumber: json.AddressNumber,
			City:          json.City,
			Neighborhood:  json.Neighborhood,
			Country:       json.Country,
			ZipCode:       json.ZipCode,
		},
	}

	if _, err := h.service.Create(ctx, c); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated[any](ctx, nil, "Customer created successfully")
}

func (h *CustomerHandler) GetAllCustomers(ctx *gin.Context) {
	customers, err := h.service.GetAll(ctx)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())

		return
	}
	response.RespondSuccess[[]domain.Customer](ctx, customers, "")
}

func (h *CustomerHandler) Search(ctx *gin.Context) {

	filter := &filter.CustomerFilter{}

	if doc := ctx.Query("document"); doc != "" {
		filter.Document = &doc
	}

	if name := ctx.Query("name"); name != "" {
		filter.Name = &name
	}

	if status := ctx.Query("status"); status != "" {
		filter.Status = status == "true" || status == "1"
	}

	if email := ctx.Query("email"); email != "" {
		filter.Email = &email
	}

	customers, err := h.service.Search(ctx, filter)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.RespondSuccess[[]domain.Customer](ctx, customers, "")
}

func (h *CustomerHandler) GetCustomerByID(ctx *gin.Context) {
	customerID := ctx.Param("id")

	var id uint
	if _, err := fmt.Sscanf(customerID, "%d", &id); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.service.GetByID(ctx, id)
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if customer == nil {
		response.RespondError(ctx, http.StatusNotFound, "Customer not found")
		return
	}

	response.RespondSuccess[domain.Customer](ctx, *customer, "")
}

func (h *CustomerHandler) DeleteCustomer(ctx *gin.Context) {
	customerId := ctx.Param("id")

	var id uint
	if _, err := fmt.Sscanf(customerId, "%d", &id); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.service.DeleteByID(ctx, id)

	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if customer == nil {
		response.RespondError(ctx, http.StatusNotFound, "Customer not found")
		return
	}

	response.RespondSuccess[domain.Customer](ctx, *customer, "")
}

func (h *CustomerHandler) UpdateCustomer(ctx *gin.Context) {
	customerId := ctx.Param("id")

	var id uint
	if _, err := fmt.Sscanf(customerId, "%d", &id); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid customer ID")

		return
	}

	var json createCustomerRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldError := range validationErrors {
				if fieldError.Field() == "Type" && fieldError.Tag() == "oneof" {
					response.RespondError(ctx, http.StatusBadRequest, "type must be one of: administrator, mechanic, attendant")

					return
				}
			}
		}

		if err != nil {
			response.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c := domain.Customer{
		ID:       id,
		Name:     json.Name,
		Email:    json.Email,
		Document: json.Document,
		Type:     json.Type,
		Contact:  json.Contact,
		Address: &domain.Address{
			Address:       json.Address,
			AddressNumber: json.AddressNumber,
			City:          json.City,
			Neighborhood:  json.Neighborhood,
			Country:       json.Country,
			ZipCode:       json.ZipCode,
		},
	}

	if err := h.service.UpdateByID(ctx, id, c); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[any](ctx, nil, "Customer updated successfully")
}
