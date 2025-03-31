package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"apps/api/internal/api"
	"apps/api/internal/handlers"
	"apps/api/internal/storage"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	s.registerMiddleware(e)
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
	e.GET("/docs", handlers.DocsHandler)

	db := s.db.GetDB()

	postStore := storage.NewPostStorage(db)

	commentHandler := handlers.NewCommentHandler()
	postHandler := handlers.NewPostHandler(postStore)
	userHandler := handlers.NewUserHandler()

	combinedHandler := struct {
		*handlers.CommentHandler
		*handlers.PostHandler
		*handlers.UserHandler
	}{
		commentHandler,
		postHandler,
		userHandler,
	}

	api.RegisterHandlersWithBaseURL(e, combinedHandler, "api/v1")
}
