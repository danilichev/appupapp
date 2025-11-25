package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"apps/api/internal/models"
	"apps/api/internal/utils"
)

func getTestLinkRepo() *LinkRepo {
	return NewLinkRepo(testDbService.GetDB())
}

func TestLinkRepo_CreateLink(t *testing.T) {
	ctx := context.Background()

	t.Run("should create link successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)
		folder := createTestFolder(t, ctx, user.ID)

		linkRepo := getTestLinkRepo()

		linkCreate := models.LinkCreate{
			Description: "Test Link Description",
			FolderId:    folder.ID,
			Name:        "Test Link",
			Url:         "https://example.com",
			UserId:      user.ID,
		}
		link, err := linkRepo.CreateLink(ctx, linkCreate)

		require.NoError(t, err)
		require.NotNil(t, link)

		assert.NotEmpty(t, link.ID)
		assert.Equal(t, linkCreate.Description, link.Description)
		assert.Equal(t, linkCreate.FolderId, link.FolderId)
		assert.Equal(t, linkCreate.Name, link.Name)
		assert.Equal(t, linkCreate.Url, link.Url)
		assert.Equal(t, user.ID, link.UserId)
	})
	t.Run(
		"should return error when creating link with invalid folder ID",
		func(t *testing.T) {
			cleanupTestDatabase()
			user := createTestUSer(t, ctx)
			linkRepo := getTestLinkRepo()
			linkCreate := models.LinkCreate{
				Description: "Test Link Description",
				FolderId:    "invalid-folder-id",
				Name:        "Test Link",
				Url:         "https://example.com",
				UserId:      user.ID,
			}
			link, err := linkRepo.CreateLink(ctx, linkCreate)
			require.Error(t, err)
			require.Nil(t, link)
		},
	)
	t.Run(
		"should return error when creating link with invalid user ID",
		func(t *testing.T) {
			cleanupTestDatabase()
			linkRepo := getTestLinkRepo()
			linkCreate := models.LinkCreate{
				Description: "Test Link Description",
				FolderId:    "valid-folder-id",
				Name:        "Test Link",
				Url:         "https://example.com",
				UserId:      "invalid-user-id",
			}
			link, err := linkRepo.CreateLink(ctx, linkCreate)
			require.Error(t, err)
			require.Nil(t, link)
		},
	)
	t.Run(
		"should return error when creating link with invalid data",
		func(t *testing.T) {
			cleanupTestDatabase()
			linkRepo := getTestLinkRepo()
			linkCreate := models.LinkCreate{
				Description: "",
				FolderId:    "",
				Name:        "",
				Url:         "",
				UserId:      "",
			}
			link, err := linkRepo.CreateLink(ctx, linkCreate)
			require.Error(t, err)
			require.Nil(t, link)
		},
	)
}
func TestLinkRepo_DeleteLink(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete link successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)
		folder := createTestFolder(t, ctx, user.ID)

		linkRepo := getTestLinkRepo()

		linkCreate := models.LinkCreate{
			Description: "Test Link Description",
			FolderId:    folder.ID,
			Name:        "Test Link",
			Url:         "https://example.com",
			UserId:      user.ID,
		}
		link, err := linkRepo.CreateLink(ctx, linkCreate)

		require.NoError(t, err)
		require.NotNil(t, link)

		err = linkRepo.DeleteLink(ctx, link.ID)
		require.NoError(t, err)

		deletedLink, err := linkRepo.GetLinkById(ctx, link.ID)
		require.Error(t, err)
		require.Nil(t, deletedLink)
	})
	t.Run(
		"should return error when deleting link with invalid ID",
		func(t *testing.T) {
			cleanupTestDatabase()
			linkRepo := getTestLinkRepo()
			err := linkRepo.DeleteLink(ctx, "invalid-link-id")
			require.Error(t, err)
		},
	)
}
func TestLinkRepo_GetLinkById(t *testing.T) {
	ctx := context.Background()

	t.Run("should get link by ID successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)
		folder := createTestFolder(t, ctx, user.ID)

		linkRepo := getTestLinkRepo()

		linkCreate := models.LinkCreate{
			Description: "Test Link Description",
			FolderId:    folder.ID,
			Name:        "Test Link",
			Url:         "https://example.com",
			UserId:      user.ID,
		}
		link, err := linkRepo.CreateLink(ctx, linkCreate)

		require.NoError(t, err)
		require.NotNil(t, link)

		retrievedLink, err := linkRepo.GetLinkById(ctx, link.ID)
		require.NoError(t, err)
		require.NotNil(t, retrievedLink)

		assert.Equal(t, link.ID, retrievedLink.ID)
		assert.Equal(t, link.Description, retrievedLink.Description)
		assert.Equal(t, link.FolderId, retrievedLink.FolderId)
		assert.Equal(t, link.Name, retrievedLink.Name)
		assert.Equal(t, link.Url, retrievedLink.Url)
		assert.Equal(t, user.ID, retrievedLink.UserId)
	})
	t.Run(
		"should return error when getting link with invalid ID",
		func(t *testing.T) {
			cleanupTestDatabase()
			linkRepo := getTestLinkRepo()
			link, err := linkRepo.GetLinkById(ctx, "invalid-link-id")
			require.Error(t, err)
			require.Nil(t, link)
		},
	)
}
func TestLinkRepo_GetLinksByFolderId(t *testing.T) {
	ctx := context.Background()

	t.Run("should get links by folder ID successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)
		folder := createTestFolder(t, ctx, user.ID)

		linkRepo := getTestLinkRepo()

		linkCreate := models.LinkCreate{
			Description: "Test Link Description",
			FolderId:    folder.ID,
			Name:        "Test Link",
			Url:         "https://example.com",
			UserId:      user.ID,
		}
		link, err := linkRepo.CreateLink(ctx, linkCreate)

		require.NoError(t, err)
		require.NotNil(t, link)

		links, total, err := linkRepo.GetLinksByFolderId(ctx, folder.ID, 10, 0)
		require.NoError(t, err)
		require.NotNil(t, links)
		require.NotNil(t, total)

		assert.Equal(t, 1, total)
		assert.Len(t, links, 1)
		assert.Equal(t, link.ID, links[0].ID)
		assert.Equal(t, link.Description, links[0].Description)
		assert.Equal(t, link.FolderId, links[0].FolderId)
		assert.Equal(t, link.Name, links[0].Name)
		assert.Equal(t, link.Url, links[0].Url)
		assert.Equal(t, user.ID, links[0].UserId)
	})
	t.Run(
		"should return empty slice when no links found for folder ID",
		func(t *testing.T) {
			cleanupTestDatabase()

			user := createTestUSer(t, ctx)
			folder := createTestFolder(t, ctx, user.ID)

			linkRepo := getTestLinkRepo()
			links, total, err := linkRepo.GetLinksByFolderId(
				ctx,
				folder.ID,
				10,
				0,
			)
			require.NoError(t, err)
			require.Empty(t, links)
			require.Equal(t, 0, total)
		},
	)
}
func TestLinkRepo_UpdateLink(t *testing.T) {
	ctx := context.Background()

	t.Run("should update link successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)
		folder := createTestFolder(t, ctx, user.ID)

		linkRepo := getTestLinkRepo()

		linkCreate := models.LinkCreate{
			Description: "Test Link Description",
			FolderId:    folder.ID,
			Name:        "Test Link",
			Url:         "https://example.com",
			UserId:      user.ID,
		}
		link, err := linkRepo.CreateLink(ctx, linkCreate)

		require.NoError(t, err)
		require.NotNil(t, link)

		updateParams := models.LinkUpdate{
			Description: utils.StringPtr("Updated Link Description"),
			Name:        utils.StringPtr("Updated Link Name"),
		}
		updatedLink, err := linkRepo.UpdateLink(ctx, link.ID, updateParams)
		require.NoError(t, err)
		require.NotNil(t, updatedLink)

		assert.Equal(t, link.ID, updatedLink.ID)
		assert.Equal(t, *updateParams.Description, updatedLink.Description)
		assert.Equal(t, link.FolderId, updatedLink.FolderId)
		assert.Equal(t, *updateParams.Name, updatedLink.Name)
		assert.Equal(t, link.Url, updatedLink.Url)
		assert.Equal(t, user.ID, updatedLink.UserId)
	})
	t.Run(
		"should return error when updating link with invalid ID",
		func(t *testing.T) {
			cleanupTestDatabase()
			linkRepo := getTestLinkRepo()
			updateParams := models.LinkUpdate{
				Description: utils.StringPtr("Updated Link Description"),
				Name:        utils.StringPtr("Updated Link Name"),
			}
			link, err := linkRepo.UpdateLink(
				ctx,
				"invalid-link-id",
				updateParams,
			)
			require.Error(t, err)
			require.Nil(t, link)
		},
	)
}
