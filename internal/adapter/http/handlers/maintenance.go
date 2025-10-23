package handlers

import (
	"strconv"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/gin-gonic/gin"
)

type MaintenanceHandler struct {
	service ports.MaintenanceService
}

func NewMaintenanceHandler(service ports.MaintenanceService) *MaintenanceHandler {
	return &MaintenanceHandler{service: service}
}

type createMaintenanceRequest struct {
	Name                string    `json:"name" binding:"required"`
	DefaultPrice        *float64  `json:"default_price" binding:"required"`
	CategoryId          uint      `json:"category_id" binding:"required"`
	Number              string    `json:"number" binding:"required"`
	MaintenanceCategory string    `json:"maintenanceCategory" binding:"required"`
	CreatedAt           time.Time `json:"created_at" binding:"required"`
	UpdatedAt           time.Time `json:"updated_at" binding:"required"`
}

func (h *MaintenanceHandler) CreateMaintenance(ctx *gin.Context) {
	var json createMaintenanceRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c := domain.Maintenance{
		Name:         json.Name,
		DefaultPrice: json.DefaultPrice,
		CategoryId:   json.CategoryId,
		Number:       json.Number,
		MaintenanceCategory: domain.MaintenanceCategory{
			CreatedAt: json.CreatedAt,
			UpdatedAt: json.UpdatedAt,
		},
	}

	if _, err := h.service.Create(ctx.Request.Context(), c); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Maintenance created successfully"})
}

func (h *MaintenanceHandler) GetMaintenanceByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	maintenance, err := h.service.GetByID(ctx.Request.Context(), uint(idUint))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, maintenance)
}

func (h *MaintenanceHandler) UpdateMaintenance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	var json createMaintenanceRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c := domain.Maintenance{
		ID:           uint(idUint),
		Name:         json.Name,
		DefaultPrice: json.DefaultPrice,
		CategoryId:   json.CategoryId,
		Number:       json.Number,
		MaintenanceCategory: domain.MaintenanceCategory{
			ID:        uint(idUint),
			CreatedAt: json.CreatedAt,
			UpdatedAt: json.UpdatedAt,
		},
	}

	if err := h.service.UpdateByID(ctx.Request.Context(), uint(idUint), c); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Maintenance updated successfully"})
}

func (h *MaintenanceHandler) DeleteMaintenance(ctx *gin.Context) {
	idStr := ctx.Param("id")
	idUint, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	if _, err := h.service.DeleteByID(ctx.Request.Context(), uint(idUint)); err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to delete maintenance"})
		return
	}
	ctx.JSON(200, gin.H{"message": "Maintenance deleted successfully"})
}
