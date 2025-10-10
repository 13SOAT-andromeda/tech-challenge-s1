package handlers

import (
	"fmt"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service ports.CustomerService
}

func NewCustomerHandler(service ports.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

type createCustomerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Document string `json:"document" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=individual company"`
	Contact  string `json:"contact" binding:"required"`
	Address  struct {
		Address       string `json:"address" binding:"required"`
		AddressNumber string `json:"address_number" binding:"required"`
		City          string `json:"city" binding:"required"`
		Neighborhood  string `json:"neighborhood" binding:"required"`
		Country       string `json:"country" binding:"required"`
		ZipCode       string `json:"zip_code" binding:"required"`
	} `json:"address" binding:"required,dive"`
}

func (h *CustomerHandler) CreateCustomer(ctx *gin.Context) {
	var json createCustomerRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c := domain.Customer{
		Name:     json.Name,
		Email:    json.Email,
		Document: json.Document,
		Type:     json.Type,
		Contact:  json.Contact,
		Address: &domain.Address{
			Address:       json.Address.Address,
			AddressNumber: json.Address.AddressNumber,
			City:          json.Address.City,
			Neighborhood:  json.Address.Neighborhood,
			Country:       json.Address.Country,
			ZipCode:       json.Address.ZipCode,
		},
	}

	if _, err := h.service.Create(ctx, c); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Customer created successfully"})
}

func (h *CustomerHandler) GetAllCustomers(ctx *gin.Context) {
	customers, err := h.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, customers)
}

func (h *CustomerHandler) GetCustomerByID(ctx *gin.Context) {
	customerID := ctx.Param("id")

	var id uint
	if _, err := fmt.Sscanf(customerID, "%d", &id); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid customer ID"})
		return
	}

	customer, err := h.service.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if customer == nil {
		ctx.JSON(404, gin.H{"error": "Customer not found"})
		return
	}
	ctx.JSON(200, customer)
}
