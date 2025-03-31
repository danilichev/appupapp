package handlers

import "github.com/labstack/echo/v4"

type CommentHandler struct{}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{}
}

func (h *CommentHandler) CreateComment(ctx echo.Context, postId string) error {
	return nil
}

func (h *CommentHandler) DeleteComment(ctx echo.Context, postId string, commentId string) error {
	return nil
}

func (h *CommentHandler) GetComment(ctx echo.Context, postId string, commentId string) error {
	return nil
}

func (h *CommentHandler) ListComments(ctx echo.Context, postId string) error {
	return nil
}

func (h *CommentHandler) UpdateComment(ctx echo.Context, postId string, commentId string) error {
	return nil
}
