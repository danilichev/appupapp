package server

import (
	"net/http"
	"slices"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"apps/api/internal/api"
	"apps/api/internal/errors"
	"apps/api/internal/handlers"
	"apps/api/internal/repositories"
	"apps/api/internal/services"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Logger.SetLevel(log.INFO)
	e.Validator = services.NewValidator()

	e.HTTPErrorHandler = errors.HTTPErrorHandler

	s.registerMiddleware(e)
	s.registerRoutes(e)

	return e
}

func (s *Server) registerMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://*", "http://*"},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"PATCH",
		},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	e.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(services.JwtClaims)
		},
		SigningKey: []byte(s.config.Jwt.SecretKey),
		Skipper: func(c echo.Context) bool {
			notRestrictedPathes := []string{
				"/api/v1/auth/login",
				"/api/v1/auth/refresh",
				"/api/v1/auth/register",
				"/api/v1/ping",
				"/docs",
			}
			return slices.Contains(notRestrictedPathes, c.Path())
		},
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user")
			if user == nil {
				return next(c)
			}

			token, ok := user.(*jwt.Token)
			if !ok {
				return next(c)
			}

			claims, ok := token.Claims.(*services.JwtClaims)
			if !ok || !token.Valid {
				return next(c)
			}

			c.Set("userId", claims.UserId)
			return next(c)
		}
	})
}

func (s *Server) registerRoutes(e *echo.Echo) {
	e.GET("/docs", handlers.DocsHandler)

	db := s.db.GetDB()

	postRepo := repositories.NewPostRepo(db)
	userRepo := repositories.NewUserRepo(db)

	jwtService := services.NewJWTService(s.config.Jwt)

	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	pingHandler := handlers.NewPingHandler()
	postHandler := handlers.NewPostHandler(postRepo, userRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	combinedHandler := struct {
		*handlers.AuthHandler
		*handlers.PingHandler
		*handlers.PostHandler
		*handlers.UserHandler
	}{
		authHandler,
		pingHandler,
		postHandler,
		userHandler,
	}

	api.RegisterHandlersWithBaseURL(e, combinedHandler, "api/v1")
}
