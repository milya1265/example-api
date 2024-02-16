package app

import (
	"example1/config"
	"example1/internal/app/grpc"
	servhttp "example1/internal/app/http_server"
	"example1/internal/handler"
	"example1/internal/repository"
	repo "example1/internal/repository/sqlc/generate"
	"example1/internal/service"
	log "example1/pkg/logger"
	"example1/pkg/postgres"
	"net"
	"time"
)

type App struct {
	Config     *config.Config
	Logger     log.Logger
	HttpServer servhttp.HttpServer
	GRPCServer grpcserver.GRPCServer
}

func New(config *config.Config) *App {
	l := log.GetLogger(config.LevelDebug)
	httpServer := servhttp.New()
	grpcServer := grpcserver.New(":8081", time.Second)
	app := &App{
		Logger:     l,
		Config:     config,
		HttpServer: httpServer,
		GRPCServer: grpcServer,
	}
	return app
}

func (a *App) Run() {
	logger := log.GetLogger(a.Config.LevelDebug)
	cl, err := postgres.NewDatabaseClient(a.Config.Storage)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("you")

	userRepository := repo.New(cl)
	productRepository := repository.NewProductRepository(cl)
	warehouseRepository := repository.NewWarehouseRepository(cl)

	authService := service.NewAuthService(*userRepository, *a.Config)
	productService := service.NewProductService(productRepository, warehouseRepository)
	warehouseService := service.NewWarehouseService(warehouseRepository)

	middleware := handler.NewAuthHandler(authService)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService, middleware)
	productHandler := handler.NewProductHandler(productService, middleware)

	go a.HttpServer.ListenAndServe(warehouseHandler, productHandler)

	a.Logger.Info("starting http server")

	l, err := net.Listen("tcp", a.Config.Listen.GrpcPort)
	if err != nil {
		a.Logger.Fatal(err)
	}

	a.Logger.Info("starting grpc server ", l.Addr().String())

	if err := a.GRPCServer.Serve(l, &authService); err != nil {
		a.Logger.Fatal(err)
	}

}
