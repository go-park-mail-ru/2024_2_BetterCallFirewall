package models

import (
	"time"
)

type Post struct {
	ID        uint32    `json:"id"`
	Header    string    `json:"header"`
	Body      string    `json:"body"`
	FilesPath []string  `json:"files_path"`
	UserID    uint32    `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type PostDB struct {
	ID        uint32 `json:"id"`
	AuthorID  uint32 `json:"author_id"`
	ContentID uint32 `json:"content_id"`
}
