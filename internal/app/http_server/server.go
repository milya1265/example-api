package http_server

import (
	"example1/config"
	"example1/internal/handler"
	"example1/pkg/logger"
	"github.com/gin-gonic/gin"
)

type httpServer struct {
	Router *gin.Engine
	Config *config.Config
	Logger logger.Logger
}

func New() HttpServer {
	return &httpServer{
		Router: gin.New(),
		Config: config.GetConfig(),
		Logger: logger.Get(),
	}
}

type HttpServer interface {
	ListenAndServe(handler.ProductHandler, handler.WarehouseHandler)
}

func (h *httpServer) ListenAndServe(productHandler handler.ProductHandler, warehouseHandler handler.WarehouseHandler) {
	h.Logger.Info("starting ListenAndServe server")
	productHandler.Register(h.Router)
	warehouseHandler.Register(h.Router)
	err := h.Router.Run(h.Config.Listen.HttpPort)
	if err != nil {
		h.Logger.Fatal(err)
	}
}
