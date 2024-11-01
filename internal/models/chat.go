package models

import (
	"time"
)

type Chat struct {
	LastMessage Message `json:"last_message"`
	Receiver    Header  `json:"receiver"`
}

type Message struct {
	Sender    uint32    `json:"sender"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
