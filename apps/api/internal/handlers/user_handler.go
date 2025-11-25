package handlers

import (
	"apps/api/internal/api"
	"apps/api/internal/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userRepo *repositories.UserRepo
}

func NewUserHandler(userRepo *repositories.UserRepo) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) GetUsersMe(
	c echo.Context,
) error {
	userId := c.Get("userId")
	if userId == nil {
		return echo.NewHTTPError(
			http.StatusUnauthorized,
			"User not authenticated",
		)
	}

	user, err := h.userRepo.GetUserById(c.Request().Context(), userId.(string))
	if err != nil || user == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to retrieve user",
		)
	}

	return c.JSON(
		http.StatusOK,
		api.User{
			Email: user.Email,
			Id:    user.ID,
		})
}
