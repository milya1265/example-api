package handler

import (
	"example1/internal/DTO"
	"example1/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	ErrInvalidBody   = "invalid request body"
	ErrInternalError = "internal error"
)

type warehouseHandler struct {
	Service service.WarehouseService
}

func NewWarehouseHandler(s service.WarehouseService) WarehouseHandler {
	return &warehouseHandler{
		Service: s,
	}
}

type WarehouseHandler interface {
	Register(r *gin.Engine)
}

func (h *warehouseHandler) Register(router *gin.Engine) {
	router.Handle(http.MethodPost, "/GetAllProducts", h.GetAllProducts)
}

func (h *warehouseHandler) GetAllProducts(c *gin.Context) {
	req := DTO.ReqGetProducts{}
	err := c.BindJSON(&req)
	if err != nil || req.WarehouseID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": ErrInvalidBody})
		return
	}

	res, err := h.Service.GetProducts(&req)
	if err != nil {
		if service.ErrNoProducts == err.Error() {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": ErrInternalError})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, res)
}
