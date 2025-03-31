package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"apps/api/internal/api"
)

type PostStorage struct {
	db *pgxpool.Pool
}

func NewPostStorage(db *pgxpool.Pool) *PostStorage {
	return &PostStorage{db: db}
}

func (s *PostStorage) CreatePost(ctx context.Context, post *api.Post) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO posts (title, content, author_id)
		VALUES ($1, $2, $3)
	`, post.Title, post.Content, post.AuthorId)
	if err != nil {
		return err
	}
	return nil
}
