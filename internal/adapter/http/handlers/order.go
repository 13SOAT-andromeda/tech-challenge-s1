package handlers

import (
	"net/http"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/converters"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service ports.OrderService
	usecase ports.OrderUseCase
}

func NewOrderHandler(service ports.OrderService, usecase ports.OrderUseCase) *OrderHandler {
	return &OrderHandler{service: service, usecase: usecase}
}

type createOrderRequest struct {
	VehicleKilometers int     `json:"vehicle_kilometers" binding:"required"`
	Note              *string `json:"note"`
	CustomerVehicleID uint    `json:"customer_vehicle_id" binding:"required"`
	CompanyID         uint    `json:"company_id" binding:"required"`
	ProductIDs        []uint  `json:"product_ids" binding:"required"`
	MaintenanceIDs    []uint  `json:"maintenance_ids" binding:"required"`
}

func (h *OrderHandler) Create(ctx *gin.Context) {
	var request createOrderRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, exists := ctx.Get("user_id")
	if !exists {
		response.RespondError(ctx, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	input := ports.CreateOrderInput{
		VehicleKilometers: request.VehicleKilometers,
		Note:              request.Note,
		UserID:            userId.(uint),
		CustomerVehicleID: request.CustomerVehicleID,
		CompanyID:         request.CompanyID,
		ProductIDs:        request.ProductIDs,
		MaintenanceIDs:    request.MaintenanceIDs,
	}

	created, err := h.usecase.CreateOrder(ctx, input)
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated(ctx, created, "Order created successfully")
}

func (h *OrderHandler) Assign(ctx *gin.Context) {
	orderID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userId, exists := ctx.Get("user_id")
	if !exists {
		response.RespondError(ctx, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	if err := h.usecase.AssignOrder(ctx, uint(orderID), userId.(uint)); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess(ctx, "", "Order assigned successfully")
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

func (h *OrderHandler) ApproveOrder(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecase.ApproveOrder(ctx, uint(id)); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess(ctx, id, "Order approved successfully")
}

func (h *OrderHandler) RejectOrder(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.usecase.RejectOrder(ctx, uint(id)); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess(ctx, id, "Order rejected successfully")
}
