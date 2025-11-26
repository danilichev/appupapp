package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"apps/api/internal/api"
	"apps/api/internal/errors"
	"apps/api/internal/models"
	"apps/api/internal/repositories"
	"apps/api/internal/schemas"
	"apps/api/internal/utils"
)

type PostHandler struct {
	postRepo *repositories.PostRepo
	userRepo *repositories.UserRepo
}

func NewPostHandler(
	postRepo *repositories.PostRepo,
	userRepo *repositories.UserRepo,
) *PostHandler {
	return &PostHandler{
		postRepo,
		userRepo,
	}
}

func (h *PostHandler) DeletePostsPostId(
	c echo.Context,
	postId string,
) error {
	userId := c.Get("userId").(string)

	post, err := h.postRepo.GetPostById(
		c.Request().Context(),
		postId,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"Post not found",
		)
	}
	if post.AuthorId != userId {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"You do not have permission to delete this post",
		)
	}

	if err := h.postRepo.DeletePost(
		c.Request().Context(),
		postId,
	); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to delete post",
		)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *PostHandler) GetPosts(
	c echo.Context,
	params api.GetPostsParams,
) error {
	if errs := schemas.GetPostsParamsSchema.Validate(&params); errs != nil {
		return errors.NewValidationError(&errs)
	}

	limit := 20
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}
	offset := 0
	if params.Offset != nil && *params.Offset >= 0 {
		offset = *params.Offset
	}
	fmt.Print("Fetching posts with params: ", limit, offset)

	posts, total, err := h.postRepo.GetPosts(
		c.Request().Context(),
		limit,
		offset,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to retrieve links",
		)
	}
	if posts == nil {
		fmt.Println("No posts found")
		posts = []*models.Post{}
	}

	return c.JSON(
		http.StatusOK,
		api.PaginatedPosts{
			Items:  utils.MapSlice(posts, mapModelPostToApi),
			Limit:  &limit,
			Offset: &offset,
			Total:  total,
		},
	)
}

func (h *PostHandler) GetPostsPostId(
	c echo.Context,
	postId string,
) error {
	post, err := h.postRepo.GetPostById(
		c.Request().Context(),
		postId,
	)
	if err != nil || post == nil {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"Post not found",
		)
	}

	return c.JSON(
		http.StatusOK,
		mapModelPostToApi(post),
	)
}

func (h *PostHandler) PatchPostsPostId(
	c echo.Context,
	postId string,
) error {
	var req api.UpdatePostRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}

	if errs := schemas.UpdatePostRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	post, err := h.postRepo.UpdatePost(
		c.Request().Context(),
		postId,
		models.PostUpdate{
			Content: req.Content,
			Title:   req.Title,
		},
	)

	if err != nil || post == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to update post",
		)
	}

	return c.JSON(http.StatusOK, mapModelPostToApi(post))
}

func (h *PostHandler) PostPosts(c echo.Context) error {
	var req api.CreatePostRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}

	if errs := schemas.CreatePostRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	post, err := h.postRepo.CreatePost(
		c.Request().Context(),
		models.PostCreate{
			AuthorId: c.Get("userId").(string),
			Title:    req.Title,
			Content:  req.Content,
		},
	)

	if err != nil || post == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to create post")
	}

	return c.JSON(http.StatusCreated, mapModelPostToApi(post))
}

func mapModelPostToApi(post *models.Post) api.Post {
	if post == nil {
		return api.Post{}
	}
	return api.Post{
		Id:       post.ID,
		AuthorId: post.AuthorId,
		Content:  post.Content,
		Title:    post.Title,
	}
}
