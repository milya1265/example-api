package handler

import (
	"errors"
	"example1/internal/DTO"
	"example1/internal/service"
	"example1/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ErrInvalidBody = errors.New("invalid request body")
var ErrInternalError = errors.New("internal error")

type warehouseHandler struct {
	Service service.WarehouseService
	Logger  logger.Logger
}

func NewWarehouseHandler(s service.WarehouseService) WarehouseHandler {
	return &warehouseHandler{
		Service: s,
		Logger:  logger.Get(),
	}
}

type WarehouseHandler interface {
	Register(r *gin.Engine)
}

func (h *warehouseHandler) Register(router *gin.Engine) {
	router.Handle(http.MethodPost, "/GetAllProducts", h.GetAllProducts)
}

func (h *warehouseHandler) GetAllProducts(c *gin.Context) {
	h.Logger.Info("start handler GetAllProducts")

	req := DTO.ReqGetProducts{}
	err := c.BindJSON(&req)
	if err != nil || req.WarehouseID == 0 {
		h.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ErrInvalidBody.Error()})
		return
	}

	res, err := h.Service.GetProducts(&req)
	if err != nil {
		if errors.Is(service.ErrNoProducts, err) {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": ErrInternalError.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}
