package models

import (
	"time"
)

type Link struct {
	ID          string    `db:"id"          fieldtag:"pk" json:"id"`
	Description string    `db:"description"               json:"description"`
	FolderId    string    `db:"folder_id"                 json:"folderId"`
	IsFavorite  bool      `db:"is_favorite"               json:"isFavorite"`
	Name        string    `db:"name"                      json:"name"`
	Url         string    `db:"url"                       json:"url"`
	UserId      string    `db:"user_id"                   json:"userId"`
	CreatedAt   time.Time `db:"created_at"                json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"                json:"updatedAt"`
}

type LinkCreate struct {
	Description string `db:"description" json:"description"`
	FolderId    string `db:"folder_id"   json:"folderId"`
	Name        string `db:"name"        json:"name"`
	Url         string `db:"url"         json:"url"`
	UserId      string `db:"user_id"     json:"userId"`
}

type LinkUpdate struct {
	Description *string `db:"description" json:"description"`
	FolderId    *string `db:"folder_id"   json:"folderId"`
	IsFavorite  *bool   `db:"is_favorite" json:"isFavorite"`
	Name        *string `db:"name"        json:"name"`
	Url         *string `db:"url"         json:"url"`
}
