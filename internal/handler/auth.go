package handler

import (
	"context"
	"errors"
	"example1/internal/service"
	"example1/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	productWorker = iota
	warehouseWorker
	admin
)

type authHandler struct {
	Service service.Auth
	Logger  logger.Logger
}

func NewAuthHandler(s service.Auth) AuthHandler {
	return &authHandler{
		Service: s,
		Logger:  logger.Get(),
	}
}

type AuthHandler interface {
	Authorize(c *gin.Context)
}

func (h *authHandler) Authorize(c *gin.Context) {
	h.Logger.Info("starting Authorize method")
	access := c.GetHeader("jwt")

	authInfo, err := h.Service.Authorize(context.Background(), access)
	if err != nil {
		if errors.Is(err, service.TokenTimeOutErr) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"messege": "login again"})
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Header("jwt", access)
	c.Set("id", authInfo.ID)
	c.Set("login", authInfo.Login)
	c.Set("role", authInfo.Role)

	c.Next()
}
