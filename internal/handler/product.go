package handler

import (
	"errors"
	"example1/internal/DTO"
	"example1/internal/service"
	"example1/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

var LackOfDataError = errors.New("lack of data")
var InvalidBodyError = errors.New("invalid body")

type productHandler struct {
	ProductService service.ProductService
	Middleware     AuthHandler
	Logger         logger.Logger
}

func NewProductHandler(s service.ProductService, auth AuthHandler) ProductHandler {
	return &productHandler{s, auth, logger.Get()}
}

type ProductHandler interface {
	Register(r *gin.Engine)
}

func (h *productHandler) Register(router *gin.Engine) {
	router.Handle(http.MethodPost, "/ReserveProduct", h.Middleware.Authorize, h.ReserveProducts)
	router.Handle(http.MethodPost, "/FreeReservation", h.Middleware.Authorize, h.FreeReservation)
}

func (h *productHandler) ReserveProducts(c *gin.Context) {
	h.Logger.Info("start handler ReserveProducts")

	temp, ok := c.Get("role")
	h.Logger.Info(temp)
	if !ok {
		h.Logger.Error("unauthorized")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userRole := temp.(int32)
	if userRole != productWorker && userRole != admin {
		h.Logger.Error("forbidden")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	req := &DTO.ReqReserveProduct{}
	err := c.BindJSON(req)
	if err != nil {
		h.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": InvalidBodyError.Error()})
		return
	}
	if req.WarehouseID == 0 || len(req.UniqueCodes) == 0 || len(req.Counts) == 0 ||
		len(req.UniqueCodes) != len(req.Counts) {
		h.Logger.Error(LackOfDataError)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": LackOfDataError.Error()})
		return
	}

	res, err := h.ProductService.Reserve(req)
	if err != nil {
		h.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(res.Unsuccessful) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"successful": res.Successful})
		return
	}

	if len(res.Successful) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"unsuccessful": res.Unsuccessful, "errors": res.Errors})
		return
	}

	c.AbortWithStatusJSON(http.StatusMultiStatus, res)
}

func (h *productHandler) FreeReservation(c *gin.Context) {
	h.Logger.Info("start handler FreeReservation")

	temp, ok := c.Get("role")
	if !ok {
		h.Logger.Error("unauthorized")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userRole := temp.(int32)
	if userRole != productWorker && userRole != admin {
		h.Logger.Error("forbidden")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	req := &DTO.ReqFreeReservation{}
	err := c.BindJSON(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": InvalidBodyError.Error()})
		return
	}

	if len(req.ID) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": LackOfDataError.Error()})
		return
	}

	res, err := h.ProductService.FreeReservation(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(res.Unsuccessful) == 0 {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{"successful": res.Successful})
		return
	}
	if len(res.Successful) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"unsuccessful": res.Unsuccessful, "errors": res.Errors})
		return
	}

	c.AbortWithStatusJSON(http.StatusMultiStatus, res)
}
