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

type LinkRepo struct {
	db *pgxpool.Pool
}

func NewLinkRepo(db *pgxpool.Pool) *LinkRepo {
	return &LinkRepo{db: db}
}

var linkStruct = sqlbuilder.NewStruct(new(models.Link)).
	For(sqlbuilder.PostgreSQL)

func (r *LinkRepo) CreateLink(
	ctx context.Context,
	params models.LinkCreate,
) (*models.Link, error) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("links")
	ib.Cols("description", "folder_id", "name", "url", "user_id")
	ib.Values(
		params.Description,
		params.FolderId,
		params.Name,
		params.Url,
		params.UserId,
	)
	ib.Returning(strings.Join(linkStruct.Columns(), ","))
	sql, args := ib.Build()
	fmt.Printf("SQL: %s, Args: %v\n", sql, args)

	var link models.Link
	err := r.db.QueryRow(ctx, sql, args...).Scan(linkStruct.Addr(&link)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to create link: %w", err)
	}
	return &link, nil
}

func (r *LinkRepo) DeleteLink(
	ctx context.Context,
	id string,
) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	db.DeleteFrom("links")
	db.Where(db.Equal("id", id))
	sql, args := db.Build()

	_, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("Failed to delete link: %w", err)
	}

	return nil
}

func (r *LinkRepo) GetLinkById(
	ctx context.Context,
	id string,
) (*models.Link, error) {
	sb := linkStruct.SelectFrom("links")
	sb.Where(sb.Equal("id", id))
	sql, args := sb.Build()

	var link models.Link
	err := r.db.QueryRow(ctx, sql, args...).Scan(linkStruct.Addr(&link)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to get link by id: %w", err)
	}

	return &link, nil
}

func (r *LinkRepo) GetLinksByFolderId(
	ctx context.Context,
	folderId string,
	limit int,
	offset int,
) ([]*models.Link, int, error) {
	sb := linkStruct.SelectFrom("links")
	sb.Where(sb.Equal("folder_id", folderId))
	sb.OrderBy("created_at").Desc()
	sb.Limit(limit)
	sb.Offset(offset)
	sql, args := sb.Build()

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to query links by folder_id: %w", err)
	}
	defer rows.Close()

	var links []*models.Link
	for rows.Next() {
		var link models.Link
		err := rows.Scan(linkStruct.Addr(&link)...)
		if err != nil {
			return nil, 0, fmt.Errorf("Failed to scan link: %w", err)
		}
		links = append(links, &link)
	}

	var total int
	sql = `SELECT COUNT(*) FROM links WHERE folder_id = $1`
	err = r.db.QueryRow(ctx, sql, folderId).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to count links by folder_id: %w", err)
	}

	return links, total, nil
}

func (r *LinkRepo) UpdateLink(
	ctx context.Context,
	id string,
	params models.LinkUpdate,
) (*models.Link, error) {
	ub := linkStruct.WithoutTag("pk").Update("links", models.Link{})
	assignments := utils.GetNotNilAssignments(params, ub)
	if len(assignments) == 0 {
		return nil, fmt.Errorf("No fields to update")
	}
	ub.Set(assignments...)
	ub.Where(ub.Equal("id", id))
	ub.SQL("RETURNING " + strings.Join(linkStruct.Columns(), ","))
	sql, args := ub.Build()

	var link models.Link
	linkRow := r.db.QueryRow(ctx, sql, args...)
	err := linkRow.Scan(linkStruct.Addr(&link)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to update link: %w", err)
	}

	return &link, nil
}
