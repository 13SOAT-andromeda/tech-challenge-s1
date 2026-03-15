package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/converters"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/encryption"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

type CreateUserRequest struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	Document      string `json:"document" binding:"required"`
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

type UpdateUserRequest struct {
	Name          string `json:"name"`
	Document      string `json:"document"`
	Contact       string `json:"contact"`
	Address       string `json:"address"`
	AddressNumber string `json:"address_number"`
	City          string `json:"city"`
	Neighborhood  string `json:"neighborhood"`
	Country       string `json:"country"`
	ZipCode       string `json:"zip_code"`
}

func (h *UserHandler) Create(ctx *gin.Context) {

	var json CreateUserRequest

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
		return
	}

	p, err := domain.NewPassword(json.Password, encryption.NewBcryptHasher())

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	u := domain.User{
		Name:     json.Name,
		Email:    json.Email,
		Document: doc,
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

	if err := u.ValidateRole(); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.service.Create(ctx, u)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.RespondCreated(ctx, user, "user created successfully")
}

func (h *UserHandler) GetByID(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) Search(ctx *gin.Context) {
	u := ctx.Request.URL.Query()
	params := converters.ParamsToMap(u)

	users, err := h.service.Search(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) Update(ctx *gin.Context) {
	var json UpdateUserRequest
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var doc *domain.Document
	if json.Document != "" {
		doc, err = domain.NewDocument(json.Document)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid document: %s", err.Error())})
			return
		}
	}

	u := domain.User{
		ID:       uint(id),
		Name:     json.Name,
		Document: doc,
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

	if _, err := h.service.Update(ctx, u); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
