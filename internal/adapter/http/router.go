package http

import (
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/handlers"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/middlewares"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
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
	vehicleHandler handlers.VehicleHandler,
	orderHandler handlers.OrderHandler,
	sessionHandler handlers.SessionHandler,
	sessionService ports.SessionService,
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

	// Initialize auth middleware
	authMiddleware := middlewares.NewAuthMiddleware(&config, sessionService)

	// Public routes (no authentication required)
	sessionGroup := router.Group("/sessions")
	{
		sessionGroup.POST("", sessionHandler.Login)            // POST /sessions (login)
		sessionGroup.GET("/validate", sessionHandler.Validate) // GET /sessions/validate
		sessionGroup.POST("/refresh", sessionHandler.Refresh)  // POST /sessions/refresh
		sessionGroup.DELETE("/logout", sessionHandler.Logout)  // DELETE /sessions/logout
	}

	// Protected routes (authentication required)
	protected := router.Group("/")
	protected.Use(authMiddleware.AuthRequired())
	{
		customerGroup := protected.Group("/customers")
		{
			customerGroup.GET("", customerHandler.Search)
			customerGroup.POST("", customerHandler.CreateCustomer)
			customerGroup.GET("/:id", customerHandler.GetCustomerByID)
			customerGroup.PUT("/:id", customerHandler.UpdateCustomer)
			customerGroup.DELETE("/:id", customerHandler.DeleteCustomer)
			customerGroup.GET("/:id/vehicles", customerHandler.GetCustomerVehicles)
			customerGroup.POST("/:id/vehicles/:vehicleId", customerHandler.AddVehicleToCustomer)
			customerGroup.DELETE("/:id/vehicles/:vehicleId", customerHandler.RemoveVehicleFromCustomer)
		}

		companyGroup := protected.Group("/companies")
		{
			companyGroup.POST("", companyHandler.CreateCompany)
			companyGroup.GET("/:id", companyHandler.GetCompanyByID)
			companyGroup.PUT("/:id", companyHandler.UpdateCompany)
			companyGroup.DELETE("/:id", companyHandler.DeleteCompany)
		}

		maintenances := protected.Group("/maintenances")
		{
			maintenances.POST("", maintenanceHandler.CreateMaintenance)
			maintenances.GET("/:id", maintenanceHandler.GetMaintenanceByID)
			maintenances.PUT("/:id", maintenanceHandler.UpdateMaintenance)
			maintenances.DELETE("/:id", maintenanceHandler.DeleteMaintenance)
		}

		productGroup := protected.Group("/products")
		{
			productGroup.POST("", productHandler.CreateProduct)
			productGroup.GET("", productHandler.GetAllProducts) // Todo: implementar search
			productGroup.GET("/:id", productHandler.GetProductByID)
			productGroup.DELETE("/:id", productHandler.DeleteProduct)
			productGroup.PUT("/:id", productHandler.UpdateProduct)

			productGroup.PATCH("/:id/stock", productHandler.ManageStockItem)
		}

		userGroup := protected.Group("/user")
		{
			userGroup.GET("", userHandler.GetAll)
			userGroup.POST("", userHandler.Create)
			userGroup.GET("/:id", userHandler.GetByID)
			userGroup.GET("/search", userHandler.Search)
			userGroup.PUT("/:id", userHandler.Update)
			userGroup.DELETE("/:id", userHandler.Delete)
		}

		vehicleGroup := protected.Group("/vehicles")
		{
			vehicleGroup.GET("", vehicleHandler.GetAll)
			vehicleGroup.POST("", vehicleHandler.Create)
			vehicleGroup.GET("/:id", vehicleHandler.GetByID)
			vehicleGroup.PUT("/:id", vehicleHandler.Update)
			vehicleGroup.DELETE("/:id", vehicleHandler.Delete)
		}

		orderGroup := protected.Group("/orders")
		{
			orderGroup.GET("", orderHandler.GetAll)
			orderGroup.GET("/:id", orderHandler.GetByID)
			orderGroup.POST("", orderHandler.Create)
			orderGroup.POST("/:id/assign", orderHandler.Assign)
			orderGroup.DELETE("/:id", orderHandler.Delete)
		}
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
