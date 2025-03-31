package handlers

import (
	"apps/api/internal/api"

	"github.com/labstack/echo/v4"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) CreateUser(ctx echo.Context) error {
	return nil
}

func (h *CommentHandler) DeleteUser(ctx echo.Context, userId string) error {
	return nil
}

func (h *CommentHandler) GetUser(ctx echo.Context, userId string) error {
	return nil
}

func (h *CommentHandler) ListUsers(ctx echo.Context, params api.ListUsersParams) error {
	return nil
}

func (h *UserHandler) UpdateUser(ctx echo.Context, userId string) error {
	return nil
}
