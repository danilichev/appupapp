package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"apps/api/internal/services"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (h *PingHandler) GetPing(c echo.Context) error {
	logger := c.Logger()

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*services.JwtClaims)

	logger.Infof("User ID: %s", claims.UserId)

	return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
}
