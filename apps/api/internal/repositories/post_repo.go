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

type PostRepo struct {
	db *pgxpool.Pool
}

func NewPostRepo(db *pgxpool.Pool) *PostRepo {
	return &PostRepo{db: db}
}

var postStruct = sqlbuilder.NewStruct(new(models.Post)).
	For(sqlbuilder.PostgreSQL)

func (r *PostRepo) CreatePost(
	ctx context.Context,
	params models.PostCreate,
) (*models.Post, error) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("posts")
	ib.Cols("author_id", "content", "title")
	ib.Values(
		params.AuthorId,
		params.Content,
		params.Title,
	)
	ib.Returning(strings.Join(postStruct.Columns(), ","))
	sql, args := ib.Build()
	fmt.Printf("SQL: %s, Args: %v\n", sql, args)

	var post models.Post
	err := r.db.QueryRow(ctx, sql, args...).Scan(postStruct.Addr(&post)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to create post: %w", err)
	}
	return &post, nil
}

func (r *PostRepo) DeletePost(
	ctx context.Context,
	id string,
) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()
	db.DeleteFrom("posts")
	db.Where(db.Equal("id", id))
	sql, args := db.Build()

	_, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("Failed to delete post: %w", err)
	}

	return nil
}

func (r *PostRepo) GetPostById(
	ctx context.Context,
	id string,
) (*models.Post, error) {
	sb := postStruct.SelectFrom("posts")
	sb.Where(sb.Equal("id", id))
	sql, args := sb.Build()

	var post models.Post
	err := r.db.QueryRow(ctx, sql, args...).Scan(postStruct.Addr(&post)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to get post by id: %w", err)
	}

	return &post, nil
}

func (r *PostRepo) GetPosts(
	ctx context.Context,
	limit int,
	offset int,
) ([]*models.Post, int, error) {
	sb := postStruct.SelectFrom("posts")
	sb.OrderBy("created_at").Desc()
	sb.Limit(limit)
	sb.Offset(offset)
	sql, args := sb.Build()

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to query posts by author_id: %w", err)
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(postStruct.Addr(&post)...)
		if err != nil {
			return nil, 0, fmt.Errorf("Failed to scan post: %w", err)
		}
		posts = append(posts, &post)
	}

	var total int
	sql = `SELECT COUNT(*) FROM posts`
	err = r.db.QueryRow(ctx, sql).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to count posts: %w", err)
	}

	return posts, total, nil
}

func (r *PostRepo) UpdatePost(
	ctx context.Context,
	id string,
	params models.PostUpdate,
) (*models.Post, error) {
	ub := postStruct.WithoutTag("pk").Update("posts", models.Post{})
	assignments := utils.GetNotNilAssignments(params, ub)
	if len(assignments) == 0 {
		return nil, fmt.Errorf("No fields to update")
	}
	ub.Set(assignments...)
	ub.Where(ub.Equal("id", id))
	ub.SQL("RETURNING " + strings.Join(postStruct.Columns(), ","))
	sql, args := ub.Build()

	var post models.Post
	postRow := r.db.QueryRow(ctx, sql, args...)
	err := postRow.Scan(postStruct.Addr(&post)...)
	if err != nil {
		return nil, fmt.Errorf("Failed to update post: %w", err)
	}

	return &post, nil
}
