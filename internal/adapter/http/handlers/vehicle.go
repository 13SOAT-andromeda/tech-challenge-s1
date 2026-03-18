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

type VehicleHandler struct {
	service ports.VehicleService
}

func NewVehicleHandler(service ports.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

type CreateVehicleRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
	Brand string `json:"brand" binding:"required"`
	Plate string `json:"plate" binding:"required"`
	Year  int    `json:"year" binding:"required"`
}

type UpdateVehicleRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Brand string `json:"brand"`
	Plate string `json:"plate"`
	Year  int    `json:"year"`
}

func (h *VehicleHandler) Create(ctx *gin.Context) {
	var json CreateVehicleRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	p, err := domain.NewPlate(json.Plate)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	v := domain.Vehicle{
		Name:  json.Name,
		Plate: p,
		Year:  json.Year,
		Brand: json.Brand,
		Color: json.Color,
	}

	vehicle, err := h.service.Create(ctx.Request.Context(), v)

	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated(ctx, vehicle, "vehicle created successfully")
}

func (h *VehicleHandler) GetAll(ctx *gin.Context) {
	u := ctx.Request.URL.Query()
	params := converters.ParamsToMap(u)

	vehicles, err := h.service.GetAll(ctx.Request.Context(), params)
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	response.RespondSuccess(ctx, vehicles, "")
}

func (h *VehicleHandler) GetByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	vehicle, err := h.service.GetByID(ctx.Request.Context(), uint(id))
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if vehicle == nil {
		response.RespondError(ctx, http.StatusNotFound, err.Error())
		return
	}
	response.RespondSuccess(ctx, vehicle, "")
}

func (h *VehicleHandler) Update(ctx *gin.Context) {
	var json UpdateVehicleRequest
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	p, err := domain.NewPlate(json.Plate)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	v := domain.Vehicle{
		ID:    uint(id),
		Name:  json.Name,
		Plate: p,
		Year:  json.Year,
		Brand: json.Brand,
		Color: json.Color,
	}

	updated, err := h.service.Update(ctx.Request.Context(), v)

	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess(ctx, updated, "User updated successfully")
}

func (h *VehicleHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess(ctx, id, "Vehicle deleted")
}
