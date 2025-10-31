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

type updateStockRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Stock     uint `json:"stock" binding:"required,min=1"`
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

func (h *ProductHandler) AddStockItem(c *gin.Context) {
	var req updateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.AddStockItem(c.Request.Context(), req.ProductID, req.Stock); err != nil {
		response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	stockResponse := gin.H{
		"product_id": req.ProductID,
		"Stock":      req.Stock,
	}

	response.RespondSuccess[any](c, stockResponse, "Estoque adicionado com sucesso")
}

func (h *ProductHandler) RemoveStockItem(c *gin.Context) {

	var req updateStockRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.RemoveStockItem(c.Request.Context(), req.ProductID, req.Stock); err != nil {
		response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	stockResponse := gin.H{
		"product_id": req.ProductID,
		"Stock":      req.Stock,
	}

	response.RespondSuccess[any](c, stockResponse, "Estoque removido com sucesso")
}

func (h *ProductHandler) SetStockItem(c *gin.Context) {
	var req setStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.SetStock(c.Request.Context(), req.ProductID, req.Stock); err != nil {
		response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	stockResponse := gin.H{
		"product_id": req.ProductID,
		"Stock":      req.Stock,
	}

	response.RespondSuccess[any](c, stockResponse, "Stock Setted Successfully")
}

func (h *ProductHandler) GetCurrentStock(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		response.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	stock, err := h.service.GetCurrentStock(c.Request.Context(), uint(productID))

	if err != nil {
		response.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	stockResponse := gin.H{
		"product_id": uint(productID),
		"stock":      stock,
	}

	response.RespondSuccess[any](c, stockResponse, "")
}
