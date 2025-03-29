package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"apps/api/internal/database"
)

func HealthHandler(database database.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, database.Health())
	}
}
