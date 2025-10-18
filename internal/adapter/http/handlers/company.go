package handlers

import (
	"fmt"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	service ports.CompanyService
}

func NewCompanyHandler(service ports.CompanyService) *CompanyHandler {
	return &CompanyHandler{service: service}
}

type createCompanyRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Document      string `json:"document" binding:"required"`
	Contact       string `json:"contact" binding:"required"`
	Address       string `json:"address" binding:"required"`
	AddressNumber string `json:"address_number" binding:"required"`
	City          string `json:"city" binding:"required"`
	Neighborhood  string `json:"neighborhood" binding:"required"`
	Country       string `json:"country" binding:"required"`
	ZipCode       string `json:"zip_code" binding:"required"`
}

func (h *CompanyHandler) CreateCompany(ctx *gin.Context) {
	var json createCompanyRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c := domain.Company{
		Name:     json.Name,
		Email:    json.Email,
		Document: json.Document,
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

	if _, err := h.service.Create(ctx.Request.Context(), c); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Company created successfully"})
}

func (h *CompanyHandler) GetCompanyByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "ID inválido"})
		return
	}
	company, err := h.service.GetByID(ctx.Request.Context(), uint(idUint))
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Company not found"})
		return
	}
	ctx.JSON(200, company)
}

func (h *CompanyHandler) UpdateCompany(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "ID invalid"})
		return
	}

	var json createCompanyRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c := domain.Company{
		ID:       uint(idUint),
		Name:     json.Name,
		Email:    json.Email,
		Document: json.Document,
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

	if _, err := h.service.UpdateByID(ctx.Request.Context(), uint(idUint), c); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to update company"})
		return
	}

	ctx.JSON(200, gin.H{"message": "Company updated successfully"})
}

func (h *CompanyHandler) DeleteCompany(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "ID invalid"})
		return
	}
	if _, err := h.service.DeleteByID(ctx.Request.Context(), uint(idUint)); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete company"})
		return
	}
	ctx.JSON(200, gin.H{"message": fmt.Sprintf("Company with ID %d deleted successfully", idUint)})
}
