package server

import (
	"fmt"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	// Register middleware
	s.registerMiddleware(e)

	// Register routes
	s.registerRoutes(e)

	return e
}

func (s *Server) registerMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func (s *Server) registerRoutes(e *echo.Echo) {
	e.GET("/", s.HelloWorldHandler)
	e.GET("/docs", s.docsHandler)
	e.GET("/health", s.healthHandler)
}

func (s *Server) docsHandler(c echo.Context) error {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "Blog API",
		},
	})

	if err != nil {
		fmt.Printf("%v", err)
		return c.String(http.StatusInternalServerError, "Error generating API reference")
	}

	return c.HTML(http.StatusOK, htmlContent)
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World!",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
