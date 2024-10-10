package models

import (
	"time"
)

type Content struct {
	ID        uint32    `json:"id"`
	Text      string    `json:"text"`
	FilesPath []string  `json:"files_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
