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
	"apps/api/internal/services"
	"apps/api/internal/utils"
)

type LinkHandler struct {
	folderRepo       *repositories.FolderRepo
	linkRepo         *repositories.LinkRepo
	userRepo         *repositories.UserRepo
	parseHtmlService *services.ParseHtmlService
}

func NewLinkHandler(
	linkRepo *repositories.LinkRepo,
	folderRepo *repositories.FolderRepo,
	userRepo *repositories.UserRepo,
	parseHtmlService *services.ParseHtmlService,
) *LinkHandler {
	return &LinkHandler{
		folderRepo,
		linkRepo,
		userRepo,
		parseHtmlService,
	}
}

func (h *LinkHandler) DeleteLinksLinkId(
	c echo.Context,
	linkId string,
) error {
	userId := c.Get("userId").(string)

	link, err := h.linkRepo.GetLinkById(
		c.Request().Context(),
		linkId,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusNotFound,
			"Link not found",
		)
	}
	if link.UserId != userId {
		return echo.NewHTTPError(
			http.StatusForbidden,
			"You do not have permission to delete this link",
		)
	}

	if err := h.linkRepo.DeleteLink(
		c.Request().Context(),
		linkId,
	); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to delete link",
		)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *LinkHandler) GetLinks(
	c echo.Context,
	params api.GetLinksParams,
) error {
	if errs := schemas.GetLinksParamsSchema.Validate(&params); errs != nil {
		return errors.NewValidationError(&errs)
	}

	userId := c.Get("userId").(string)
	limit := 20
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}
	offset := 0
	if params.Offset != nil && *params.Offset >= 0 {
		offset = *params.Offset
	}
	fmt.Print("Fetching links with params: ", limit, offset, params.FolderId)

	var folderId string
	if params.FolderId != nil {
		folderId = *params.FolderId
		folder, err := h.folderRepo.GetFolderById(
			c.Request().Context(),
			folderId,
		)
		if err != nil {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"Folder not found",
			)
		}
		if folder.UserId != userId {
			return echo.NewHTTPError(
				http.StatusForbidden,
				"You do not have permission to access this folder",
			)
		}
	} else {
		user, err := h.userRepo.GetUserById(
			c.Request().Context(),
			userId,
		)
		if err != nil || user == nil {
			return echo.NewHTTPError(
				http.StatusNotFound,
				"User not found",
			)
		}
		folderId = *user.FolderId
	}

	links, total, err := h.linkRepo.GetLinksByFolderId(
		c.Request().Context(),
		folderId,
		limit,
		offset,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to retrieve links",
		)
	}
	if links == nil {
		fmt.Println("No links found for folder:", folderId)
		links = []*models.Link{}
	}

	return c.JSON(
		http.StatusOK,
		api.PaginatedLinks{
			Items:  utils.MapSlice(links, mapModelLinkToApi),
			Limit:  &offset,
			Offset: &offset,
			Total:  total,
		},
	)
}

func (h *LinkHandler) PatchLinksLinkId(
	c echo.Context,
	urlId string,
) error {
	var req api.UpdateLinkRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}

	if errs := schemas.UpdateLinkRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	link, err := h.linkRepo.UpdateLink(
		c.Request().Context(),
		urlId,
		models.LinkUpdate{
			Description: req.Description,
			FolderId:    req.FolderId,
			IsFavorite:  req.IsFavorite,
			Name:        req.Name,
			Url:         req.Url,
		},
	)

	if err != nil || link == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to update link",
		)
	}

	return c.JSON(http.StatusOK, mapModelLinkToApi(link))
}

func (h *LinkHandler) PostLinks(c echo.Context) error {
	var req api.CreateLinkRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}

	if errs := schemas.CreateLinkRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	link, err := h.linkRepo.CreateLink(
		c.Request().Context(),
		models.LinkCreate{
			Description: *req.Description,
			FolderId:    req.FolderId,
			Name:        *req.Name,
			Url:         req.Url,
			UserId:      c.Get("userId").(string),
		},
	)

	if err != nil || link == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to create link")
	}

	return c.JSON(http.StatusCreated, mapModelLinkToApi(link))
}

func (h *LinkHandler) PostLinksParse(c echo.Context) error {
	var req api.ParseHtmlRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}

	if errs := schemas.ParseHtmlRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	parsedHtml, err := h.parseHtmlService.ParseHtml(req.Url)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to parse HTML",
		)
	}

	return c.JSON(http.StatusOK, parsedHtml)
}

func mapModelLinkToApi(link *models.Link) api.Link {
	if link == nil {
		return api.Link{}
	}
	return api.Link{
		Id:          link.ID,
		Description: &link.Description,
		FolderId:    link.FolderId,
		IsFavorite:  &link.IsFavorite,
		Name:        &link.Name,
		Url:         link.Url,
	}
}
