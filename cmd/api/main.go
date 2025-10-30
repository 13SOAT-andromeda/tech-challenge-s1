package main

import (
	"context"
	"log"
	"time"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/repository"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/handlers"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/usecases/session"
	"github.com/13SOAT-andromeda/tech-challenge-s1/pkg/jwt"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	ctx := context.Background()
	db, err := database.Init(ctx, *cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	log.Printf("Connecting to database")

	err = db.AutoMigrate(
		&customer.Model{},
		&company.Model{},
		&maintenance.Model{},
		&product.Model{},
		&user.Model{},
		&model.SessionModel{},
	)

	if err != nil {
		log.Fatalf("Error to executing migration: %s", err)
	}

	dbase := db.GetDB()

	customerRepository := repository.NewCustomerRepository(dbase)
	companyRepository := repository.NewCompanyRepository(dbase)
	maintenanceRepository := repository.NewMaintenenceRepository(dbase)
	productRepository := repository.NewProductRepository(dbase)
	userRepository := repository.NewUserRepository(dbase)
	sessionRepository := repository.NewSessionRepository(dbase)

	customerService := services.NewCustomerService(customerRepository)
	companyService := services.NewCompanyService(companyRepository)
	maintenanceService := services.NewMaintenanceService(maintenanceRepository)
	productService := services.NewProductService(productRepository)
	userService := services.NewUserService(userRepository)
	sessionService := services.NewSessionService(sessionRepository)

	customerHandler := handlers.NewCustomerHandler(customerService)
	companyHandler := handlers.NewCompanyHandler(companyService)
	maintenanceHandler := handlers.NewMaintenanceHandler(maintenanceService)
	productHandler := handlers.NewProductHandler(productService)
	userHandler := handlers.NewUserHandler(userService)

	// JWT Service
	accessExpiry, _ := time.ParseDuration(cfg.JWT.AccessTokenExpiry)
	refreshExpiry, _ := time.ParseDuration(cfg.JWT.RefreshTokenExpiry)
	jwtService := jwt.NewService(cfg.JWT.Secret, accessExpiry, refreshExpiry)

	// Session UseCases
	loginUseCase := session.NewLoginUseCase(userService, sessionService, jwtService, cfg)
	validateUseCase := session.NewValidateUseCase(userService, sessionService, jwtService)
	refreshUseCase := session.NewRefreshUseCase(userService, sessionService, jwtService, cfg)
	logoutUseCase := session.NewLogoutUseCase(sessionService)

	// Session Handler
	sessionHandler := handlers.NewSessionHandler(
		loginUseCase,
		validateUseCase,
		refreshUseCase,
		logoutUseCase,
	)

	router := http.NewRouter(*cfg, *customerHandler, *companyHandler, *maintenanceHandler, *productHandler, *userHandler, *sessionHandler, sessionService)
	log.Printf("Starting HTTP server on port %s", cfg.Http.Port)

	if err = userService.CreateAdminUser(ctx, cfg.AdminUser.Email, cfg.AdminUser.Password); err != nil {
		log.Fatalf("failed on create admin user: %v", err)
	}

	if err = router.Server(":" + cfg.Http.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
