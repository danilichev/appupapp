package handlers

import (
	"fmt"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/labstack/echo/v4"
)

func DocsHandler(c echo.Context) error {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./internal/api/openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "AppUpApp API",
		},
	})

	if err != nil {
		fmt.Printf("%v", err)
		return c.String(
			http.StatusInternalServerError,
			"Error generating API reference",
		)
	}

	return c.HTML(http.StatusOK, htmlContent)
}
