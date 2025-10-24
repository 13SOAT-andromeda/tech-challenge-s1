package handlers

import (
	"strconv"

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
	Name     string  `json:"name" binding:"required"`
	Quantity uint    `json:"quantity" binding:"required"`
	Price    float64 `json:"price" binding:"required,gt=0"`
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	var json createProductRequest
	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	p := domain.Product{
		Name:     json.Name,
		Quantity: json.Quantity,
		Price:    monetary.ConvertToMinorUnitInt64(json.Price, 2),
	}

	if _, err := h.service.Create(ctx.Request.Context(), p); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Product created successfully"})
}

func (h *ProductHandler) GetAllProducts(ctx *gin.Context) {
	products, err := h.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, products)
}

func (h *ProductHandler) GetProductByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	product, err := h.service.GetById(ctx.Request.Context(), uint(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, product)
}

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)

	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), uint(id)); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Product deleted successfully"})
}
