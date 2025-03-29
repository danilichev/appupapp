package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HelloHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello!",
	}

	return c.JSON(http.StatusOK, resp)
}
