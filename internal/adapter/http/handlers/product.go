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

type ProductHandler struct {
	service ports.ProductService
}

func NewProductHandler(service ports.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

type createProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
	Stock int64   `json:"stock" binding:"required"`
}

type updateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

type updateStockRequest struct {
	Quantity  uint   `json:"quantity" binding:"required,min=1"`
	Operation string `json:"operation" binding:"required"`
}

type setStockRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Stock     uint `json:"stock" binding:"required,min=1"`
}

type checkAvailabilityRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  uint `json:"quantity" binding:"required,min=1"`
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var json createProductRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	p := domain.Product{
		Name:  json.Name,
		Price: monetary.ConvertToMinorUnitInt64(json.Price, 2),
		Stock: uint(json.Stock),
	}

	if _, err := h.service.Create(ctx.Request.Context(), p); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated[any](ctx, nil, "Product created successfully")
}

func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	products, err := h.service.GetAll(ctx.Request.Context())
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[[]domain.Product](ctx, products, "")
}

func (h *ProductHandler) GetProductByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.GetById(ctx.Request.Context(), uint(id))
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[domain.Product](ctx, *product, "")
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	productIdStg := ctx.Param("id")

	productId, err := strconv.ParseUint(productIdStg, 10, 64)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.Delete(ctx, uint(productId))

	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if product == nil {
		response.RespondError(ctx, http.StatusNotFound, "Customer not found")
		return
	}

	response.RespondSuccess[domain.Product](ctx, *product, "")
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {

	var json updateProductRequest

	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	productIdStg := ctx.Param("id")

	productId, err := strconv.ParseUint(productIdStg, 10, 64)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	p := domain.Product{
		ID:    uint(productId),
		Name:  json.Name,
		Price: monetary.ConvertToMinorUnitInt64(json.Price, 2),
	}

	if _, err := h.service.Update(ctx.Request.Context(), p); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[any](ctx, nil, "Product updated successfully")
}

func (h *ProductHandler) ManageStockItem(ctx *gin.Context) {

	productIdStg := ctx.Param("id")

	productId, err := strconv.ParseUint(productIdStg, 10, 64)

	if err != nil {
		response.RespondError(ctx, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var req updateStockRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.service.ManageStockItem(ctx, uint(productId), req.Quantity, req.Operation)

	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if product == nil {
		response.RespondError(ctx, http.StatusNotFound, "Customer not found")
		return
	}

	response.RespondSuccess[domain.Product](ctx, *product, "")
}
