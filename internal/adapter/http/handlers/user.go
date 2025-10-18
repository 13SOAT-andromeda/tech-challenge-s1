package handlers

import (
	"net/http"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type createUserRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Password      string `json:"password" binding:"required"`
	Contact       string `json:"contact" binding:"required"`
	Role          string `json:"role" binding:"required"`
	Address       string `json:"address" binding:"required"`
	AddressNumber string `json:"address_number" binding:"required"`
	City          string `json:"city" binding:"required"`
	Neighborhood  string `json:"neighborhood" binding:"required"`
	Country       string `json:"country" binding:"required"`
	ZipCode       string `json:"zip_code" binding:"required"`
}

func (h *UserHandler) Create(ctx *gin.Context) {
	var json createUserRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	p, err := domain.NewPassword(json.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := domain.User{
		Name:     json.Name,
		Email:    json.Email,
		Password: p,
		Role:     json.Role,
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

	if _, err := h.service.Create(ctx, u); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Customer created successfully"})
}

func (h *UserHandler) GetAll(ctx *gin.Context) {
	users, err := h.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, users)
}

func (h *UserHandler) GetByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid customer ID"})
		return
	}

	customer, err := h.service.GetByID(ctx, uint(id))
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

func (h *UserHandler) Search(ctx *gin.Context) {
	params := ParamsToMap(ctx.Params)

	users, err := h.service.Search(ctx, params)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, users)
}

func ParamsToMap(params gin.Params) map[string]interface{} {
	paramsMap := make(map[string]interface{})
	for _, param := range params {
		paramsMap[param.Key] = param.Value
	}

	return paramsMap
}
