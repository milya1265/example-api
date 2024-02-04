package server

import (
	"example1/config"
	"example1/internal/handler"
	"example1/internal/logger"
	"github.com/gin-gonic/gin"
)

type httpServer struct {
	Router *gin.Engine
	Config *config.Config
}

func New() HttpServer {
	return &httpServer{
		Router: gin.New(),
		Config: config.GetConfig(),
	}
}

type HttpServer interface {
	ListenAndServe(handler.ProductHandler, handler.WarehouseHandler)
}

func (h *httpServer) ListenAndServe(productHandler handler.ProductHandler, warehouseHandler handler.WarehouseHandler) {
	productHandler.Register(h.Router)
	warehouseHandler.Register(h.Router)
	err := h.Router.Run(h.Config.Listen.Port)
	if err != nil {
		logger.Fatal(err)
	}
}
