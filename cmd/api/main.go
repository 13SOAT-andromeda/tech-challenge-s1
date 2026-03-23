package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/DataDog/dd-trace-go/v2/profiler"
	"go.uber.org/zap"

	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/config"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/company"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/customer_vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_maintenance"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/order_product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/product"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/user"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/model/vehicle"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/database/repository"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/email"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/http/handlers"
	appmetrics "github.com/13SOAT-andromeda/tech-challenge-s1/internal/adapter/metrics"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/ports"
	"github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/services"
	customerUseCase "github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/usecases/customer"
	orderUsecase "github.com/13SOAT-andromeda/tech-challenge-s1/internal/application/usecases/order"
)

func main() {

	logger, err := zap.NewProduction()

	defer func() {
		tracer.Stop()
		profiler.Stop()
		logger.Sync()
	}()

	sugar := logger.Sugar()

	if err != nil {
		sugar.Fatalf("Error on start logger zap: %s", err)
	}

	cfg, err := config.Init()
	if err != nil {
		sugar.Fatalf("failed to load config: %v", err)
	}

	err = profiler.Start(
		profiler.WithEnv(cfg.Env),
		profiler.WithService(cfg.Service),
		profiler.WithVersion(cfg.Version),
		profiler.WithTags("layer:api"),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	)

	if err != nil {
		sugar.Fatalf("Error on start datadog profiler: %s", err)
	}

	err = tracer.Start(
		tracer.WithEnv(cfg.Env),
		tracer.WithService(cfg.Service),
		tracer.WithServiceVersion(cfg.Version),
	)

	if err != nil {
		sugar.Fatalf("Error on start datadog tracer: %s", err)
	}

	ctx := context.Background()
	db, err := database.Init(ctx, *cfg.Database)
	if err != nil {
		sugar.Fatalf("failed to connect database: %v", err)
	}

	sugar.Infof("Connecting to database")

	err = db.AutoMigrate(
		&customer.Model{},
		&company.Model{},
		&maintenance.Model{},
		&product.Model{},
		&user.Model{},
		&vehicle.Model{},
		&customer_vehicle.Model{},
		&order.Model{},
		&order_maintenance.Model{},
		&order_product.Model{},
	)

	if err != nil {
		sugar.Fatalf("Error to executing migration: %s", err)
	}

	dbase := db.GetDB()

	if err = database.Seed(dbase); err != nil {
		sugar.Fatalf("failed to seed database: %v", err)
	}
	apiUrl := cfg.Http.ApiUrl

	var orderMetrics ports.OrderMetrics = appmetrics.NoopOrderMetrics{}
	if !cfg.DogStatsD.Disabled && cfg.DogStatsD.Addr != "" {
		statsdClient, errStatsd := statsd.New(cfg.DogStatsD.Addr,
			statsd.WithNamespace("tech_challenge."),
			statsd.WithTags([]string{
				"env:" + cfg.Env,
				"service:" + cfg.Service,
				"version:" + cfg.Version,
			}),
		)
		if errStatsd != nil {
			sugar.Warnw("dogstatsd indisponível, métricas de ordem desativadas", "error", errStatsd)
		} else {
			defer statsdClient.Close()
			orderMetrics = appmetrics.NewOrderStatsd(statsdClient)
		}
	}

	// Repositories
	customerRepository := repository.NewCustomerRepository(dbase)
	companyRepository := repository.NewCompanyRepository(dbase)
	maintenanceRepository := repository.NewMaintenanceRepository(dbase)
	productRepository := repository.NewProductRepository(dbase)
	userRepository := repository.NewUserRepository(dbase)
	vehicleRepository := repository.NewVehicleRepository(dbase)
	orderRepository := repository.NewOrderRepository(dbase)
	customerVehicleRepository := repository.NewCustomerVehicleRepository(dbase)
	orderMaintenanceRepository := repository.NewOrderMaintenanceRepository(dbase)
	orderProductRepository := repository.NewOrderProductRepository(dbase)

	// Services
	vehicleService := services.NewVehicleService(vehicleRepository)
	customerService := services.NewCustomerService(customerRepository, customerVehicleRepository, vehicleService)
	companyService := services.NewCompanyService(companyRepository)
	maintenanceService := services.NewMaintenanceService(maintenanceRepository, orderMaintenanceRepository)
	productService := services.NewProductService(productRepository)
	userService := services.NewUserService(userRepository)
	orderService := services.NewOrderService(orderRepository)
	emailService := email.NewSendtrap(cfg.MailTrap.ApiKey, cfg.MailTrap.ApiUrl)

	// UseCases
	createCustomerUseCase := customerUseCase.NewCustomerUseCase(customerRepository, customerVehicleRepository, vehicleService)
	createOrderUseCase := orderUsecase.NewOrderUseCase(orderService, productService, maintenanceService, customerService, emailService, orderRepository, orderProductRepository, orderMaintenanceRepository, apiUrl, orderMetrics)

	// Handlers
	customerHandler := handlers.NewCustomerHandler(customerService, createCustomerUseCase)
	companyHandler := handlers.NewCompanyHandler(companyService)
	maintenanceHandler := handlers.NewMaintenanceHandler(maintenanceService)
	productHandler := handlers.NewProductHandler(productService)
	userHandler := handlers.NewUserHandler(userService)
	vehicleHandler := handlers.NewVehicleHandler(vehicleService)
	orderHandler := handlers.NewOrderHandler(orderService, createOrderUseCase)

	router := http.NewRouter(*cfg, logger, *customerHandler, *companyHandler, *maintenanceHandler, *productHandler, *userHandler, *vehicleHandler, *orderHandler, cfg.JWT.Secret)
	sugar.Info("Starting HTTP server on port %s", cfg.Http.Port)

	if err = userService.CreateAdminUser(ctx, cfg.AdminUser.Email, cfg.AdminUser.Password, cfg.AdminUser.Document); err != nil {
		sugar.Fatalf("failed on create admin user: %v", err)
	}

	if err = router.Server(":" + cfg.Http.Port); err != nil {
		sugar.Fatalf("failed to start server: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	go func() {
		<-sigChan
		tracer.Stop()
	}()
}
