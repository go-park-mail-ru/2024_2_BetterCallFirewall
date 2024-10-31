package models

import (
	"time"
)

type Chat struct {
	LastMessage Message `json:"last_message"`
	Receiver    Header  `json:"receiver"`
}

type Message struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
