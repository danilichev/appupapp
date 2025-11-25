package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5/pgxpool"

	"apps/api/internal/models"
	"apps/api/internal/utils"
)

type FolderRepo struct {
	db *pgxpool.Pool
}

func NewFolderRepo(db *pgxpool.Pool) *FolderRepo {
	return &FolderRepo{db: db}
}

var folderStruct = sqlbuilder.NewStruct(new(models.Folder)).
	For(sqlbuilder.PostgreSQL)

func (r *FolderRepo) CreateFolder(
	ctx context.Context,
	folderCreate models.FolderCreate,
) (*models.Folder, error) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("folders")
	ib.Cols("name", "parent_id", "user_id")
	ib.Values(folderCreate.Name, folderCreate.ParentId, folderCreate.UserId)
	ib.Returning(strings.Join(folderStruct.Columns(), ","))
	sql, args := ib.Build()

	var folder models.Folder
	err := r.db.QueryRow(ctx, sql, args...).Scan(folderStruct.Addr(&folder)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to create folder: %w", err)
	}
	return &folder, nil
}

func (r *FolderRepo) DeleteFolder(
	ctx context.Context,
	id string,
) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	db.DeleteFrom("folders")
	db.Where(db.Equal("id", id))
	sql, args := db.Build()

	_, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("Failed to delete folder: %w", err)
	}

	return nil
}

func (r *FolderRepo) GetFolderById(
	ctx context.Context,
	id string,
) (*models.Folder, error) {
	sb := folderStruct.SelectFrom("folders")
	sb.Where(sb.Equal("id", id))
	sql, args := sb.Build()

	var folder models.Folder
	err := r.db.QueryRow(ctx, sql, args...).Scan(folderStruct.Addr(&folder)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to get folder by id: %w", err)
	}

	return &folder, nil
}

func (r *FolderRepo) GetFoldersByUserId(
	ctx context.Context,
	userId string,
) ([]*models.Folder, error) {
	sb := folderStruct.SelectFrom("folders")
	sb.Where(sb.Equal("user_id", userId))
	sql, args := sb.Build()

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to get folders by user id: %w", err)
	}
	defer rows.Close()

	var folders []*models.Folder
	for rows.Next() {
		var folder models.Folder
		if err := rows.Scan(folderStruct.Addr(&folder)...); err != nil {
			return nil, fmt.Errorf("Failed to scan folder: %w", err)
		}
		folders = append(folders, &folder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %w", err)
	}

	return folders, nil
}

func (r *FolderRepo) UpdateFolder(
	ctx context.Context,
	id string,
	folderUpdate models.FolderUpdate,
) (*models.Folder, error) {
	ub := folderStruct.WithoutTag("pk").Update("folders", models.Folder{})
	assignments := utils.GetNotNilAssignments(folderUpdate, ub)
	if len(assignments) == 0 {
		return nil, fmt.Errorf("No fields to update")
	}
	ub.Set(assignments...)
	ub.Where(ub.Equal("id", id))
	ub.SQL("RETURNING " + strings.Join(folderStruct.Columns(), ","))
	sql, args := ub.Build()

	var folder models.Folder
	folderRow := r.db.QueryRow(ctx, sql, args...)
	err := folderRow.Scan(folderStruct.Addr(&folder)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to update folder: %w", err)
	}

	return &folder, nil
}

func (r *FolderRepo) MoveFolderItems(
	ctx context.Context,
	targetFolderId string,
	items []models.FolderItem,
) error {
	if len(items) == 0 {
		return fmt.Errorf("No items to move")
	}

	links := make([]interface{}, 0, len(items))
	folders := make([]interface{}, 0, len(items))
	for _, item := range items {
		if item.Type == "link" {
			links = append(links, item.ID)
		} else if item.Type == "folder" {
			folders = append(folders, item.ID)
		}
	}

	if len(links) == 0 && len(folders) == 0 {
		return fmt.Errorf("No valid items to move")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	if len(links) > 0 {
		ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
		ub.Update("links")
		ub.Set(ub.Assign("folder_id", targetFolderId))
		ub.Where(ub.In("id", links...))
		sql, args := ub.Build()

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf(
				"Failed to update links with new folder_id: %w",
				err,
			)
		}
	}

	if len(folders) > 0 {
		ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
		ub.Update("folders")
		ub.Set(ub.Assign("parent_id", targetFolderId))
		ub.Where(ub.In("id", folders...))
		sql, args := ub.Build()

		_, err = tx.Exec(ctx, sql, args...)
		if err != nil {
			return fmt.Errorf(
				"Failed to update folders with new parent_id: %w",
				err,
			)
		}
	}

	return nil
}
