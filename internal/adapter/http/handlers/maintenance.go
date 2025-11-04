package handlers

import (
	"net/http"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/monetary"
	"github.com/gin-gonic/gin"
)

type MaintenanceHandler struct {
	service ports.MaintenanceService
}

func NewMaintenanceHandler(service ports.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{service: service}
}

type CreateMaintenanceRequest struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required"`
	CategoryId string  `json:"category" binding:"required"`
}

func (h *MaintenanceHandler) CreateMaintenance(ctx *gin.Context) {
	var json CreateMaintenanceRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	c := domain.Maintenance{
		Name:       json.Name,
		Price:      monetary.ConvertToMinorUnitInt64(json.Price, 2),
		CategoryId: domain.MaintenanceCategory(json.CategoryId),
	}

	if err := c.ValidateMaintenanceCategory(); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := h.service.Create(ctx.Request.Context(), c); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated[any](ctx, nil, "Maintenance created successfully")
}

func (h *MaintenanceHandler) GetMaintenanceByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}
	maintenance, err := h.service.GetByID(ctx.Request.Context(), uint(idUint))
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if maintenance == nil {
		response.RespondError(ctx, http.StatusNotFound, "Maintenance not found")
		return
	}

	response.RespondSuccess[domain.Maintenance](ctx, *maintenance, "")
}

func (h *MaintenanceHandler) UpdateMaintenance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}

	var json CreateMaintenanceRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	c := domain.Maintenance{
		ID:         uint(idUint),
		Name:       json.Name,
		Price:      monetary.ConvertToMinorUnitInt64(json.Price, 2),
		CategoryId: domain.MaintenanceCategory(json.CategoryId),
	}

	if err := c.ValidateMaintenanceCategory(); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateByID(ctx.Request.Context(), uint(idUint), c); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[any](ctx, nil, "Maintenance updated successfully")
}

func (h *MaintenanceHandler) DeleteMaintenance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid ID")
		return
	}

	maintenance, err := h.service.DeleteByID(ctx.Request.Context(), uint(idUint))
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if maintenance == nil {
		response.RespondError(ctx, http.StatusNotFound, "Maintenance not found")
		return
	}

	response.RespondSuccess[domain.Maintenance](ctx, *maintenance, "")
}
