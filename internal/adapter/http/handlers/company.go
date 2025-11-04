package handlers

import (
	"net/http"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
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

type CreateCompanyRequest struct {
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
	var json CreateCompanyRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	c := domain.Company{
		Name:     json.Name,
		Email:    json.Email,
		Document: json.Document,
		Contact:  json.Contact,
		Address: &domain.Address{
			AddressNumber: json.AddressNumber,
			City:          json.City,
			Neighborhood:  json.Neighborhood,
			Country:       json.Country,
			ZipCode:       json.ZipCode,
		},
	}

	company, err := h.service.Create(ctx.Request.Context(), c)

	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated(ctx, company, "Company created successfully")
}

func (h *CompanyHandler) GetCompanyByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "invalid ID")
		return
	}
	company, err := h.service.GetByID(ctx.Request.Context(), uint(idUint))
	if err != nil {
		response.RespondError(ctx, http.StatusNotFound, "company not found")
		return
	}
	response.RespondSuccess(ctx, company, "")
}

func (h *CompanyHandler) UpdateCompany(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "invalid ID")
		return
	}

	var json CreateCompanyRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
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

	if err := h.service.UpdateByID(ctx.Request.Context(), uint(idUint), c); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, "Failed to update company")
		return
	}

	response.RespondSuccess(ctx, "", "Company updated successfully")
}

func (h *CompanyHandler) DeleteCompany(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "invalid ID")
		return
	}
	if _, err := h.service.DeleteByID(ctx.Request.Context(), uint(idUint)); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, "Failed to delete company")
		return
	}
	response.RespondSuccess(ctx, "", "Company deleted successfully")
}
