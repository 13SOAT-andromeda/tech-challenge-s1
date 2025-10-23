package main

import (
	"context"
	"log"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/repository"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/handlers"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
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
	)

	if err != nil {
		log.Fatalf("Error to executing migration: %s", err)
	}

	dbase := db.GetDB()

	customerRepository := repository.NewCustomerRepository(dbase)
	companyRepository := repository.NewCompanyRepository(dbase)
	maintenanceRepository := repository.NewMaintenenceRepository(dbase)
	productRepository := repository.NewProductRepository(dbase)

	customerService := services.NewCustomerService(customerRepository)
	companyService := services.NewCompanyService(companyRepository)
	maintenanceService := services.NewMaintenanceService(maintenanceRepository)
	productService := services.NewProductService(productRepository)

	customerHandler := handlers.NewCustomerHandler(customerService)
	companyHandler := handlers.NewCompanyHandler(companyService)
	maintenanceHandler := handlers.NewMaintenanceHandler(maintenanceService)
	productHandler := handlers.NewProductHandler(productService)

	router := http.NewRouter(*cfg, *customerHandler, *companyHandler, *maintenanceHandler, *productHandler)
	log.Printf("Starting HTTP server on port %s", cfg.Http.Port)

	if err = router.Server(":" + cfg.Http.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
