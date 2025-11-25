package handlers

import (
	"fmt"
	"net/http"
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"apps/api/internal/api"
	"apps/api/internal/errors"
	"apps/api/internal/models"
	"apps/api/internal/repositories"
	"apps/api/internal/services"
	"apps/api/internal/utils"
)

type AuthHandler struct {
	jwtService *services.JWTService
	userRepo   *repositories.UserRepo
}

func NewAuthHandler(
	userRepo *repositories.UserRepo,
	jwtService *services.JWTService,
) *AuthHandler {
	return &AuthHandler{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

var loginRequestSchema = z.Struct(z.Schema{
	"email": utils.EmailSchema,
	"password": z.String().
		Min(1, z.Message("Should not be empty")).
		Required(z.Message("Password is required")),
})

func (h *AuthHandler) PostAuthLogin(c echo.Context) error {
	var req api.LoginRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}

	if errs := loginRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	email := strings.TrimSpace(string(req.Email))
	password := strings.TrimSpace(req.Password)

	user, err := h.userRepo.GetUserByEmail(c.Request().Context(), email)
	if err != nil || user == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to retrieve user",
		)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"Invalid email or password",
		)
	}

	authToken, err := h.jwtService.GenerateAuthToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to generate auth token",
		)
	}

	return c.JSON(http.StatusOK, authToken)
}

var registerRequestSchema = z.Struct(z.Schema{
	"email": utils.EmailSchema,
	"password": z.String().
		Min(6, z.Message("Must be at least 6 characters")).
		Required(z.Message("Password is required")),
})

func (h *AuthHandler) PostAuthRegister(c echo.Context) error {
	var req api.RegisterRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}

	if errs := registerRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	email := strings.TrimSpace(string(req.Email))
	password := strings.TrimSpace(req.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to hash password",
		)
	}

	user, err := h.userRepo.CreateUser(
		c.Request().Context(),
		models.UserCreate{
			Email:        email,
			PasswordHash: string(hashedPassword),
		},
	)

	fmt.Println("User created:", user)
	fmt.Println("User created err:", err)

	if err != nil || user == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to create user",
		)
	}

	authToken, err := h.jwtService.GenerateAuthToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to generate auth token",
		)
	}

	return c.JSON(http.StatusCreated, authToken)
}

func (h *AuthHandler) PostAuthRefresh(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Missing Authorization header",
		)
	}

	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Invalid Authorization header format",
		)
	}

	refreshToken := authHeader[len(prefix):]

	claims, err := h.jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"Invalid refresh token",
		)
	}

	user, err := h.userRepo.GetUserById(c.Request().Context(), claims.UserId)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to retrieve user",
		)
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
	}

	newAccessToken, err := h.jwtService.GenerateAuthToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to generate new access token",
		)
	}

	return c.JSON(http.StatusOK, newAccessToken)
}
