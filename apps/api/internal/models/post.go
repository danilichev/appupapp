package models

import (
	"time"
)

type Post struct {
	ID        string    `db:"id"         fieldtag:"pk" json:"id"`
	AuthorId  string    `db:"author_id"                json:"authorId"`
	Content   string    `db:"content"                  json:"content"`
	Title     string    `db:"title"                    json:"title"`
	CreatedAt time.Time `db:"created_at"               json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at"               json:"updatedAt"`
}

type PostCreate struct {
	Content  string `db:"content"   json:"content"`
	Title    string `db:"title"     json:"title"`
	AuthorId string `db:"author_id" json:"authorId"`
}

type PostUpdate struct {
	AuthorId *string `db:"author_id" json:"authorId"`
	Content  *string `db:"content"   json:"content"`
	Title    *string `db:"title"     json:"title"`
}
