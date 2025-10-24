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
			customerGroup.GET("", customerHandler.GetAllCustomers)
			customerGroup.POST("", customerHandler.CreateCustomer)
			customerGroup.GET("/:id", customerHandler.GetCustomerByID)
		}
	}

	{
		companyGroup := router.Group("/companies")
		{
			companyGroup.POST("", companyHandler.CreateCompany)
			companyGroup.GET("/:id", companyHandler.GetCompanyByID)
			companyGroup.PUT("/:id", companyHandler.UpdateCompany)
			companyGroup.DELETE("/:id", companyHandler.DeleteCompany)
		}
	}

	{
		companyGroup := router.Group("/maintenances")
		{
			companyGroup.POST("", maintenanceHandler.CreateMaintenance)
			companyGroup.GET("/:id", maintenanceHandler.GetMaintenanceByID)
			companyGroup.PUT("/:id", maintenanceHandler.UpdateMaintenance)
			companyGroup.DELETE("/:id", maintenanceHandler.DeleteMaintenance)
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
