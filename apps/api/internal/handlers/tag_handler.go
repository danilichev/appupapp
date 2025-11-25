package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"apps/api/internal/api"
	"apps/api/internal/utils"
)

type TagHandler struct{}

func NewTagHandler() *TagHandler {
	return &TagHandler{}
}

// GetTags handles retrieving all available tags
func (h *TagHandler) GetTags(c echo.Context) error {
	// TODO: Implement actual tag retrieval logic

	// Placeholder response
	tags := []api.Tag{
		{
			Id:   utils.StringPtr("tag-1"),
			Name: utils.StringPtr("development"),
		},
		{
			Id:   utils.StringPtr("tag-2"),
			Name: utils.StringPtr("productivity"),
		},
		{
			Id:   utils.StringPtr("tag-3"),
			Name: utils.StringPtr("design"),
		},
	}

	return c.JSON(http.StatusOK, tags)
}
