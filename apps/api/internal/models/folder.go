package models

import (
	"time"
)

type Folder struct {
	ID        string    `db:"id"         fieldtag:"pk" json:"id"`
	Name      string    `db:"name"                     json:"name"`
	ParentId  *string   `db:"parent_id"                json:"parent_id,omitempty"`
	UserId    string    `db:"user_id"                  json:"user_id"`
	CreatedAt time.Time `db:"created_at"               json:"created_at"`
	UpdatedAt time.Time `db:"updated_at"               json:"updated_at"`
}

type FolderCreate struct {
	Name     string  `db:"name"      json:"name"`
	ParentId *string `db:"parent_id" json:"parent_id,omitempty"`
	UserId   string  `db:"user_id"   json:"user_id"`
}

type FolderUpdate struct {
	Name     *string `db:"name"      json:"name"`
	ParentId *string `db:"parent_id" json:"parent_id,omitempty"`
}

type FolderItem struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
