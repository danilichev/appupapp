package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func BindRequest[T any](c echo.Context, req *T) error {
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Invalid request payload",
		)
	}
	return nil
}
