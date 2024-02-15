package main

import (
	conf "example1/config"
	"example1/internal/handler"
	"example1/internal/repository"
	serv "example1/internal/server"
	"example1/internal/service"
	log "example1/pkg/logger"
	"example1/pkg/postgres"
)

func main() {
	config := conf.GetConfig()
	server := serv.New()
	logger := log.GetLogger(config.LevelDebug)
	cl, err := postgres.NewDatabaseClient(config.Storage)
	if err != nil {
		logger.Fatal(err)
	}

	productRepository := repository.NewProductRepository(cl)
	warehouseRepository := repository.NewWarehouseRepository(cl)

	productService := service.NewProductService(productRepository, warehouseRepository)
	warehouseService := service.NewWarehouseService(warehouseRepository)

	warehouseHandler := handler.NewWarehouseHandler(warehouseService)
	productHandler := handler.NewProductHandler(productService)

	server.ListenAndServe(productHandler, warehouseHandler)
}
