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
	UserHandler handlers.UserHandler,
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
		customerGroup := router.Group("/customer")
		{
			customerGroup.GET("", customerHandler.GetAllCustomers)
			customerGroup.POST("", customerHandler.CreateCustomer)
			customerGroup.GET("/:id", customerHandler.GetCustomerByID)
		}

		userGroup := router.Group("/user")
		{
			userGroup.GET("", UserHandler.GetAll)
			userGroup.POST("", UserHandler.Create)
			userGroup.GET("/:id", UserHandler.GetByID)
			userGroup.GET("/search", UserHandler.Search)
			userGroup.PUT("/:id", UserHandler.Update)
			userGroup.DELETE("/:id", UserHandler.Delete)
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
