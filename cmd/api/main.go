package main

import (
	"context"
	"log"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model"
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
		&model.CustomerModel{},
	)

	if err != nil {
		log.Fatalf("Error to executing migration: %s", err)
	}

	dbase := db.GetDB()

	customerRepository := repository.NewCustomerRepository(dbase)

	customerService := services.NewCustomerService(customerRepository)

	customerHandler := handlers.NewCustomerHandler(customerService)

	router := http.NewRouter(*cfg, *customerHandler)
	log.Printf("Starting HTTP server on port %s", cfg.Http.Port)

	if err = router.Server(":" + cfg.Http.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
