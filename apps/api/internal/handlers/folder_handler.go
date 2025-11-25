package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"apps/api/internal/api"
	"apps/api/internal/errors"
	"apps/api/internal/models"
	"apps/api/internal/repositories"
	"apps/api/internal/schemas"
	"apps/api/internal/utils"
)

type FolderHandler struct {
	folderRepo *repositories.FolderRepo
}

func NewFolderHandler(folderRepo *repositories.FolderRepo) *FolderHandler {
	return &FolderHandler{
		folderRepo: folderRepo,
	}
}

func (h *FolderHandler) DeleteFoldersFolderId(
	c echo.Context,
	folderId string,
) error {
	if err := h.checkFolder(c, folderId); err != nil {
		return err
	}

	if err := h.folderRepo.DeleteFolder(
		c.Request().Context(),
		folderId,
	); err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to delete folder",
		)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *FolderHandler) GetFolders(
	c echo.Context,
) error {
	userId := c.Get("userId").(string)

	folders, err := h.folderRepo.GetFoldersByUserId(
		c.Request().Context(),
		userId,
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to get folders",
		)
	}

	return c.JSON(http.StatusOK, utils.MapSlice(folders, mapModelFolderToApi))
}

func (h *FolderHandler) PatchFoldersFolderId(
	c echo.Context,
	folderId string,
) error {
	var req api.UpdateFolderRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}

	if errs := schemas.UpdateFolderRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	err := h.checkFolder(c, folderId)
	if err != nil {
		return err
	}

	var folder *models.Folder
	folder, err = h.folderRepo.UpdateFolder(
		c.Request().Context(),
		folderId,
		models.FolderUpdate{
			Name:     req.Name,
			ParentId: req.ParentId,
		},
	)
	if err != nil || folder == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to update folder",
		)
	}

	return c.JSON(http.StatusOK, mapModelFolderToApi(folder))
}

func (h *FolderHandler) PostFolders(c echo.Context) error {
	var req api.CreateFolderRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}

	if errs := schemas.CreateFolderRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	userId := c.Get("userId").(string)
	folder, err := h.folderRepo.CreateFolder(
		c.Request().Context(),
		models.FolderCreate{
			Name:     req.Name,
			ParentId: utils.StringPtr(req.ParentId),
			UserId:   userId,
		},
	)

	if err != nil || folder == nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to create folder")
	}

	return c.JSON(http.StatusCreated, mapModelFolderToApi(folder))
}

func (h *FolderHandler) PostFoldersItemsMove(
	c echo.Context,
) error {
	var req api.MoveItemsRequest
	if err := utils.BindRequest(c, &req); err != nil {
		return err
	}
	if errs := schemas.MoveFolderItemsRequestSchema.Validate(&req); errs != nil {
		return errors.NewValidationError(&errs)
	}

	err := h.checkFolder(c, req.TargetFolderId)
	if err != nil {
		return err
	}

	items := mapApiFolderItemsToModel(req.Items)
	if items == nil {
		return echo.NewHTTPError(
			http.StatusBadRequest,
			"Invalid folder items",
		)
	}
	err = h.folderRepo.MoveFolderItems(
		c.Request().Context(),
		req.TargetFolderId,
		mapApiFolderItemsToModel(req.Items),
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			err,
		)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (h *FolderHandler) checkFolder(
	c echo.Context,
	folderId string,
) error {
	userId := c.Get("userId").(string)

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

	return nil
}

func mapModelFolderToApi(folder *models.Folder) api.Folder {
	if folder == nil {
		return api.Folder{}
	}
	return api.Folder{
		Id:       folder.ID,
		Name:     folder.Name,
		ParentId: folder.ParentId,
		UserId:   folder.UserId,
	}
}

func mapApiFolderItemsToModel(
	items []api.FolderItem,
) []models.FolderItem {
	if items == nil {
		return nil
	}
	return utils.MapSlice(
		items,
		func(item api.FolderItem) models.FolderItem {
			return models.FolderItem{
				ID:   item.Id,
				Type: string(item.Type),
			}
		},
	)
}
