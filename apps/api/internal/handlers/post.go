package handlers

import (
	"github.com/labstack/echo/v4"

	"apps/api/internal/api"
	"apps/api/internal/storage"
)

type PostHandler struct {
	storage *storage.PostStorage
}

func NewPostHandler(store *storage.PostStorage) *PostHandler {
	return &PostHandler{storage: store}
}

func (h *PostHandler) CreatePost(ctx echo.Context) error {
	print("CreatePost")
	return nil
}

func (h *CommentHandler) DeletePost(ctx echo.Context, postId string) error {
	return nil
}

func (h *CommentHandler) GetPost(ctx echo.Context, postId string) error {
	return nil
}

func (h *CommentHandler) ListPosts(ctx echo.Context, params api.ListPostsParams) error {
	return nil
}

func (h *PostHandler) UpdatePost(ctx echo.Context, postId string) error {
	return nil
}
