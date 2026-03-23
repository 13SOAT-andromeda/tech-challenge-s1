package handlers

import (
	"net/http"
	"strconv"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/response"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/domain"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service ports.ProductService
}

func NewProductHandler(service ports.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

type CreateProductRequest struct {
	Name  string `json:"name" binding:"required"`
	Price int64  `json:"price" binding:"required,gt=0"`
	Stock int64  `json:"stock" binding:"required"`
}

type UpdateProductRequest struct {
	Name  string `json:"name" binding:"required"`
	Price int64  `json:"price" binding:"required,gt=0"`
}

type StockItem struct {
	ID        uint   `json:"id" binding:"required"`
	Quantity  uint   `json:"quantity" binding:"required"`
	Operation string `json:"operation" binding:"required, oneof=INCREASE DECREASE"`
}

type updateStockRequest struct {
	Items []StockItem `json:"items" binding:"required"`
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var json CreateProductRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	stockReq := uint(json.Stock)
	p := domain.Product{
		Name:  json.Name,
		Price: json.Price,
		Stock: &stockReq,
	}

	if _, err := h.service.Create(ctx.Request.Context(), p); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondCreated[any](ctx, nil, "Product created successfully")
}

func (h *ProductHandler) GetProducts(ctx *gin.Context) {
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

	product, err := h.service.Delete(ctx.Request.Context(), uint(productId))

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

	var json UpdateProductRequest

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
		Price: json.Price,
	}

	if _, err := h.service.Update(ctx.Request.Context(), p); err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[any](ctx, nil, "Product updated successfully")
}

func (h *ProductHandler) UpdateStockBatch(ctx *gin.Context) {
	var request updateStockRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	domainItems := make([]domain.StockItem, len(request.Items))
	for i, item := range request.Items {
		op := domain.StockOperation(item.Operation)
		domainItems[i] = domain.StockItem{
			ID:        item.ID,
			Quantity:  item.Quantity,
			Operation: &op,
		}
	}

	err := h.service.UpdateStock(ctx.Request.Context(), domainItems)
	if err != nil {
		response.RespondError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	response.RespondSuccess[any](ctx, nil, "Product stock updated successfully")
}
