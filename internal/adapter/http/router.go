package http

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func NewRouter(
	config config.Config,
	customerHandler handlers.CustomerHandler,
	companyHandler handlers.CompanyHandler,
	maintenanceHandler handlers.MaintenanceHandler,
	productHandler handlers.ProductHandler,
	userHandler handlers.UserHandler,
) *Router {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.Http.AllowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), cors.New(corsConfig))

	{
		customerGroup := router.Group("/customers")
		{
			customerGroup.GET("/:id", customerHandler.GetCustomerByID)
			customerGroup.GET("", customerHandler.Search)
			customerGroup.POST("", customerHandler.CreateCustomer)
			customerGroup.PUT("/:id", customerHandler.UpdateCustomer)
			customerGroup.DELETE("/:id", customerHandler.DeleteCustomer)
		}
	}

	companyGroup := router.Group("/companies")
	{
		companyGroup.POST("", companyHandler.CreateCompany)
		companyGroup.GET("/:id", companyHandler.GetCompanyByID)
		companyGroup.PUT("/:id", companyHandler.UpdateCompany)
		companyGroup.DELETE("/:id", companyHandler.DeleteCompany)
	}

	maintenances := router.Group("/maintenances")
	{
		maintenances.POST("", maintenanceHandler.CreateMaintenance)
		maintenances.GET("/:id", maintenanceHandler.GetMaintenanceByID)
		maintenances.PUT("/:id", maintenanceHandler.UpdateMaintenance)
		maintenances.DELETE("/:id", maintenanceHandler.DeleteMaintenance)
	}

	productGroup := router.Group("/products")
	{
		productGroup.POST("", productHandler.CreateProduct)
		productGroup.GET("", productHandler.GetAllProducts) // Todo: implementar search
		productGroup.GET("/:id", productHandler.GetProductByID)
		productGroup.DELETE("/:id", productHandler.DeleteProduct)

		productGroup.POST("/stock/add", productHandler.AddStockItem)
		productGroup.POST("/stock/remove", productHandler.RemoveStockItem)
		productGroup.POST("/stock/set", productHandler.SetStockItem)
	}

	userGroup := router.Group("/user")
	{
		userGroup.GET("", userHandler.GetAll)
		userGroup.POST("", userHandler.Create)
		userGroup.GET("/:id", userHandler.GetByID)
		userGroup.GET("/search", userHandler.Search)
		userGroup.PUT("/:id", userHandler.Update)
		userGroup.DELETE("/:id", userHandler.Delete)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// serve static swagger files from ./swagger so /swagger/swagger.yaml is available
	router.Static("/swagger", "./swagger")

	// Serve the swagger UI under /docs and point it to the static spec at /swagger/swagger.yaml
	// Use a different prefix than /swagger to avoid wildcard conflicts with the static route.
	router.GET("/docs/*any", swagger.WrapHandler(swaggerFiles.Handler, swagger.URL("/swagger/swagger.yaml")))

	return &Router{router}
}

func (r *Router) Server(listenAddr string) error {
	return r.Run(listenAddr)
}
