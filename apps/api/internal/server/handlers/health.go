package handlers

import (
	"net/http"

	"apps/api/internal/database"

	"github.com/labstack/echo/v4"
)

func HealthHandler(database database.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, database.Health())
	}
}
