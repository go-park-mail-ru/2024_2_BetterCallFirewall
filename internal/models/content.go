package models

import (
	"time"
)

type Content struct {
	ID        uint32    `json:"id"`
	Text      string    `json:"text"`
	File      *Picture  `json:"file"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
