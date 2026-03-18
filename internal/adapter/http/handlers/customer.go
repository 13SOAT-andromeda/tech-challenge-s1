package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain/filter"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomerHandler struct {
	service ports.CustomerService
	useCase ports.CustomerUseCase
}

func NewCustomerHandler(service ports.CustomerService, useCase ports.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{
		service: service,
		useCase: useCase,
	}
}

type CreateCustomerRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Document      string `json:"document" binding:"required"`
	Type          string `json:"type" binding:"required"`
	Contact       string `json:"contact" binding:"required"`
	Address       string `json:"address" binding:"required"`
	AddressNumber string `json:"address_number" binding:"required"`
	City          string `json:"city" binding:"required"`
	Neighborhood  string `json:"neighborhood" binding:"required"`
	Country       string `json:"country" binding:"required"`
	ZipCode       string `json:"zip_code" binding:"required"`
}

func (h *CustomerHandler) CreateCustomer(ctx *gin.Context) {

	var json CreateCustomerRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldError := range validationErrors {
				if fieldError.Field() == "Type" && fieldError.Tag() == "oneof" {
					response.RespondError(ctx, http.StatusBadRequest, "type must be one of: pf, pj'")

					return
				}
			}
		}

		response.RespondError(ctx, http.StatusBadRequest, err.Error())

		return
	}

	doc, err := domain.NewDocument(json.Document)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
	}

	c := domain.Customer{
		Name:     json.Name,
		Email:    json.Email,
		Document: doc,
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

	customer, err := h.service.Create(ctx.Request.Context(), c)

	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated(ctx, customer, "Customer created successfully")
}

func (h *CustomerHandler) Search(ctx *gin.Context) {

	customerFilter := &filter.CustomerFilter{}

	if doc := ctx.Query("document"); doc != "" {
		customerFilter.Document = &doc
	}

	if name := ctx.Query("name"); name != "" {
		customerFilter.Name = &name
	}

	if status := ctx.Query("status"); status != "" {
		customerFilter.Status = status == "true" || status == "1"
	}

	if email := ctx.Query("email"); email != "" {
		customerFilter.Email = &email
	}

	customers, err := h.service.Search(ctx.Request.Context(), customerFilter)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.RespondSuccess[[]domain.Customer](ctx, customers, "")
}

func (h *CustomerHandler) GetCustomerByID(ctx *gin.Context) {
	customerIdStr := ctx.Param("id")

	customerId, err := strconv.ParseUint(customerIdStr, 10, 64)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.service.GetByID(ctx.Request.Context(), uint(customerId))
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
	customerIdStr := ctx.Param("id")

	customerId, err := strconv.ParseUint(customerIdStr, 10, 64)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.service.DeleteByID(ctx.Request.Context(), uint(customerId))

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
	customerIdStr := ctx.Param("id")

	customerId, err := strconv.ParseUint(customerIdStr, 10, 64)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var json CreateCustomerRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {

		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, fieldError := range validationErrors {
				if fieldError.Field() == "Type" && fieldError.Tag() == "oneof" {
					response.RespondError(ctx, http.StatusBadRequest, "type must be one of: pf, pj")

					return
				}
			}
		}

		if err != nil {
			response.RespondError(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}

	doc, err := domain.NewDocument(json.Document)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	c := domain.Customer{
		ID:       uint(customerId),
		Name:     json.Name,
		Email:    json.Email,
		Document: doc,
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

	if err := h.service.UpdateByID(ctx.Request.Context(), uint(customerId), c); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[any](ctx, nil, "Customer updated successfully")
}

func (h *CustomerHandler) AddVehicleToCustomer(ctx *gin.Context) {
	customerIDStr := ctx.Param("id")
	vehicleIDStr := ctx.Param("vehicleId")

	customerID, err := strconv.ParseUint(customerIDStr, 10, 64)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 64)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	if err := h.useCase.AddVehicleToCustomer(ctx.Request.Context(), uint(customerID), uint(vehicleID)); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated[any](ctx, nil, "Vehicle associated with customer successfully")
}

func (h *CustomerHandler) RemoveVehicleFromCustomer(ctx *gin.Context) {
	customerIDStr := ctx.Param("id")
	vehicleIDStr := ctx.Param("vehicleId")

	customerID, err := strconv.ParseUint(customerIDStr, 10, 64)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	vehicleID, err := strconv.ParseUint(vehicleIDStr, 10, 64)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid vehicle ID")
		return
	}

	if err := h.useCase.RemoveVehicleFromCustomer(ctx.Request.Context(), uint(customerID), uint(vehicleID)); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[any](ctx, nil, "Vehicle removed from customer successfully")
}

func (h *CustomerHandler) GetCustomerVehicles(ctx *gin.Context) {
	customerIDStr := ctx.Param("id")

	customerID, err := strconv.ParseUint(customerIDStr, 10, 64)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	vehicles, err := h.useCase.GetCustomerVehicles(ctx.Request.Context(), uint(customerID))
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[[]domain.CustomerVehicle](ctx, vehicles, "")
}
