package models

import (
	"time"
)

type Chat struct {
	LastMessage string    `json:"last_message"`
	LastDate    time.Time `json:"last_date"`
	Receiver    Header    `json:"receiver"`
}

type Message struct {
	Sender    uint32         `json:"sender"`
	Receiver  uint32         `json:"receiver"`
	Content   MessageContent `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
}

type MessageContent struct {
	Text        string `json:"text"`
	FilePath    string `json:"file_path"`
	StickerPath string `json:"sticker_path"`
}
