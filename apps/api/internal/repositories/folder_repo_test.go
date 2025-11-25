package repositories

import (
	"apps/api/internal/models"
	"apps/api/internal/utils"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestFolderRepo() *FolderRepo {
	return NewFolderRepo(testDbService.GetDB())
}

func TestFolderRepo_CreateFolder(t *testing.T) {
	ctx := context.Background()

	t.Run("should create folder successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)

		folderRepo := getTestFolderRepo()

		folderCreate := models.FolderCreate{
			Name:     "Test Folder",
			ParentId: nil,
			UserId:   user.ID,
		}
		folder, err := folderRepo.CreateFolder(ctx, folderCreate)

		require.NoError(t, err)
		require.NotNil(t, folder)

		assert.NotEmpty(t, folder.ID)
		assert.Equal(t, folderCreate.Name, folder.Name)
		assert.Nil(t, folder.ParentId)
		assert.Equal(t, user.ID, folder.UserId)
	})

	t.Run(
		"should return error when creating folder with invalid data",
		func(t *testing.T) {
			cleanupTestDatabase()
			folderRepo := getTestFolderRepo()

			folderCreate := models.FolderCreate{
				Name:     "",
				ParentId: nil,
				UserId:   "",
			}
			folder, err := folderRepo.CreateFolder(ctx, folderCreate)

			require.Error(t, err)
			require.Nil(t, folder)
		},
	)
}

func TestFolderRepo_DeleteFolder(t *testing.T) {
	ctx := context.Background()

	t.Run("should delete folder successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)

		folderRepo := getTestFolderRepo()

		folderCreate := models.FolderCreate{
			Name:     "Test Folder",
			ParentId: nil,
			UserId:   user.ID,
		}
		createdFolder, err := folderRepo.CreateFolder(ctx, folderCreate)
		require.NoError(t, err)

		err = folderRepo.DeleteFolder(ctx, createdFolder.ID)
		require.NoError(t, err)

		deletedFolder, err := folderRepo.GetFolderById(ctx, createdFolder.ID)
		require.Error(t, err)
		require.Nil(t, deletedFolder)
	})

	t.Run(
		"should return error when deleting link with invalid ID",
		func(t *testing.T) {
			cleanupTestDatabase()
			folderRepo := getTestFolderRepo()

			err := folderRepo.DeleteFolder(
				ctx,
				"invalid-folder-id",
			)
			require.Error(t, err)
		},
	)
}

func TestFolderRepo_GetFolderById(t *testing.T) {
	ctx := context.Background()

	t.Run("should get folder by ID successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)

		folderRepo := getTestFolderRepo()

		folderCreate := models.FolderCreate{
			Name:     "Test Folder",
			ParentId: nil,
			UserId:   user.ID,
		}
		createdFolder, err := folderRepo.CreateFolder(ctx, folderCreate)
		require.NoError(t, err)

		folder, err := folderRepo.GetFolderById(ctx, createdFolder.ID)

		require.NoError(t, err)
		require.NotNil(t, folder)

		assert.Equal(t, createdFolder.ID, folder.ID)
		assert.Equal(t, createdFolder.Name, folder.Name)
		assert.Equal(t, createdFolder.ParentId, folder.ParentId)
		assert.Equal(t, createdFolder.UserId, folder.UserId)
	})

	t.Run("should return error when folder not found", func(t *testing.T) {
		cleanupTestDatabase()
		folderRepo := getTestFolderRepo()

		nonExistentFolderId := "00000000-0000-0000-0000-000000000000"
		folder, err := folderRepo.GetFolderById(ctx, nonExistentFolderId)
		require.Error(t, err)
		require.Nil(t, folder)
	})

	t.Run("should handle invalid UUID format", func(t *testing.T) {
		cleanupTestDatabase()
		folderRepo := getTestFolderRepo()

		user, err := folderRepo.GetFolderById(ctx, "invalid-uuid")

		assert.Error(t, err)
		assert.Nil(t, user)
	})
}

func TestFolderRepo_GetFoldersByUserId(t *testing.T) {
	ctx := context.Background()

	t.Run("should get folders by user ID successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)

		folderRepo := getTestFolderRepo()

		folderCreate1 := models.FolderCreate{
			Name:     "Test Folder 1",
			ParentId: user.FolderId,
			UserId:   user.ID,
		}
		folderCreate2 := models.FolderCreate{
			Name:     "Test Folder 2",
			ParentId: user.FolderId,
			UserId:   user.ID,
		}

		_, err := folderRepo.CreateFolder(ctx, folderCreate1)
		require.NoError(t, err)
		_, err = folderRepo.CreateFolder(ctx, folderCreate2)
		require.NoError(t, err)

		folders, err := folderRepo.GetFoldersByUserId(ctx, user.ID)

		require.NoError(t, err)
		// One folder is created by default when a user is created
		require.Len(t, folders, 3)

		assert.Equal(t, user.ID, folders[0].UserId)
		assert.Equal(t, user.ID, folders[1].UserId)
		assert.Equal(t, user.ID, folders[2].UserId)
	})

	t.Run(
		"should return empty slice when no folders found for user",
		func(t *testing.T) {
			cleanupTestDatabase()
			folderRepo := getTestFolderRepo()

			userId := "00000000-0000-0000-0000-000000000000"
			folders, err := folderRepo.GetFoldersByUserId(ctx, userId)

			require.NoError(t, err)
			require.Empty(t, folders)
		},
	)
}

func TestFolderRepo_UpdateFolder(t *testing.T) {
	ctx := context.Background()

	t.Run("should update folder successfully", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)

		folderRepo := getTestFolderRepo()

		folderCreate := models.FolderCreate{
			Name:     "Test Folder",
			ParentId: user.FolderId,
			UserId:   user.ID,
		}
		createdFolder, err := folderRepo.CreateFolder(ctx, folderCreate)
		require.NoError(t, err)

		folderCreate2 := models.FolderCreate{
			Name:     "Test Folder2",
			ParentId: user.FolderId,
			UserId:   user.ID,
		}
		createdFolder2, err2 := folderRepo.CreateFolder(ctx, folderCreate2)
		require.NoError(t, err2)

		folderUpdate := models.FolderUpdate{
			Name:     utils.StringPtr("Updated Folder"),
			ParentId: utils.StringPtr(createdFolder2.ID),
		}
		updatedFolder, err := folderRepo.UpdateFolder(
			ctx,
			createdFolder.ID,
			folderUpdate,
		)

		require.NoError(t, err)
		require.NotNil(t, updatedFolder)

		assert.Equal(t, createdFolder.ID, updatedFolder.ID)
		assert.Equal(t, "Updated Folder", updatedFolder.Name)
		assert.Equal(t, folderUpdate.ParentId, updatedFolder.ParentId)
		assert.Equal(t, user.ID, updatedFolder.UserId)
	})

	t.Run(
		"should return error when updating invalid folder",
		func(t *testing.T) {
			cleanupTestDatabase()
			folderRepo := getTestFolderRepo()

			folderUpdate := models.FolderUpdate{
				Name: utils.StringPtr("Updated Folder"),
				ParentId: utils.StringPtr(
					"invalid-parent-id",
				),
			}
			folder, err := folderRepo.UpdateFolder(
				ctx,
				"invalid-folder-id",
				folderUpdate,
			)

			require.Error(t, err)
			require.Nil(t, folder)
		},
	)
}

func TestFolderRepo_MoveFolderItems(t *testing.T) {
	ctx := context.Background()

	t.Run(
		"should move a mix of links and folders successfully",
		func(t *testing.T) {
			cleanupTestDatabase()

			user := createTestUSer(t, ctx)
			folderRepo := getTestFolderRepo()

			// Create target folder
			targetFolder := createTestFolder(t, ctx, user.ID)

			// Create source folder and items
			sourceFolder := createTestFolder(t, ctx, user.ID)
			link1 := createTestLink(t, ctx, user.ID, sourceFolder.ID)
			link2 := createTestLink(t, ctx, user.ID, sourceFolder.ID)
			folder1 := createTestFolder(t, ctx, user.ID)
			folder2 := createTestFolder(t, ctx, user.ID)

			items := []models.FolderItem{
				{ID: link1.ID, Type: "link"},
				{ID: link2.ID, Type: "link"},
				{ID: folder1.ID, Type: "folder"},
				{ID: folder2.ID, Type: "folder"},
			}

			err := folderRepo.MoveFolderItems(ctx, targetFolder.ID, items)
			require.NoError(t, err)

			// Verify links moved
			assertLinkFolder(t, ctx, link1.ID, targetFolder.ID)
			assertLinkFolder(t, ctx, link2.ID, targetFolder.ID)

			// Verify folders moved
			assertFolderParent(t, ctx, folder1.ID, targetFolder.ID)
			assertFolderParent(t, ctx, folder2.ID, targetFolder.ID)
		},
	)

	t.Run("should handle empty items list", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)
		folderRepo := getTestFolderRepo()

		targetFolder := createTestFolder(t, ctx, user.ID)

		err := folderRepo.MoveFolderItems(
			ctx,
			targetFolder.ID,
			[]models.FolderItem{},
		)
		require.NoError(t, err)
	})

	t.Run("should rollback transaction on failure", func(t *testing.T) {
		cleanupTestDatabase()

		user := createTestUSer(t, ctx)
		folderRepo := getTestFolderRepo()

		targetFolder := createTestFolder(t, ctx, user.ID)
		invalidItem := models.FolderItem{ID: "invalid-id", Type: "link"}

		err := folderRepo.MoveFolderItems(
			ctx,
			targetFolder.ID,
			[]models.FolderItem{invalidItem},
		)
		require.Error(t, err)
	})
}
