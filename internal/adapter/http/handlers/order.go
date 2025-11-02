package handlers

import (
	"net/http"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/converters"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service ports.OrderService
}

func NewOrderHandler(service ports.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

type createOrderRequest struct {
	VehicleKilometers int      `json:"vehicle_kilometers" binding:"required"`
	Note              *string  `json:"note"`
	Price             *float64 `json:"price"`
	UserID            uint     `json:"user_id" binding:"required"`
	CustomerVehicleID uint     `json:"customer_vehicle_id" binding:"required"`
	CompanyID         uint     `json:"company_id" binding:"required"`
}

func (h *OrderHandler) Create(ctx *gin.Context) {
	var json createOrderRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	order := domain.Order{
		User: domain.User{
			ID: json.UserID,
		},
		CustomerVehicle: domain.CustomerVehicle{
			ID: json.CustomerVehicleID,
		},
		Company: domain.Company{
			ID: json.CompanyID,
		},
		VehicleKilometers: json.VehicleKilometers,
		Note:              json.Note,
		Price:             json.Price,
	}

	created, err := h.service.Create(ctx, order)
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated(ctx, created, "Order created successfully")
}

func (h *OrderHandler) GetAll(ctx *gin.Context) {
	u := ctx.Request.URL.Query()
	params := converters.ParamsToMap(u)

	orders, err := h.service.GetAll(ctx, params)
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.RespondSuccess(ctx, orders, "")
}

func (h *OrderHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if order == nil {
		response.RespondError(ctx, http.StatusNotFound, "Order not found")
		return
	}
	response.RespondSuccess(ctx, order, "")
}

func (h *OrderHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(ctx, uint(id)); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess(ctx, id, "Order deleted successfully")
}
