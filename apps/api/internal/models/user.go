package models

import (
	"time"
)

type User struct {
	ID           string    `db:"id"            fieldtag:"pk" json:"id"`
	Email        string    `db:"email"                       json:"email"`
	PasswordHash string    `db:"password_hash"               json:"-"`
	FolderId     *string   `db:"folder_id"                   json:"folderId"`
	CreatedAt    time.Time `db:"created_at"                  json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at"                  json:"updatedAt"`
}

type UserCreate struct {
	Email        string `db:"email"         json:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
}
