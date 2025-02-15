package models

import (
	"strings"
	"time"
)

//easyjson:json
type Chat struct {
	LastMessage string    `json:"last_message"`
	LastDate    time.Time `json:"last_date"`
	Receiver    Header    `json:"receiver"`
}

//easyjson:json
type Message struct {
	Sender    uint32         `json:"sender"`
	Receiver  uint32         `json:"receiver"`
	Content   MessageContent `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
}

//easyjson:skip
type MessageDto struct {
	Sender    uint32
	Receiver  uint32
	Content   MessageContentDto
	CreatedAt time.Time
}

func (m *Message) ToDto() MessageDto {
	return MessageDto{
		Sender:    m.Sender,
		Receiver:  m.Receiver,
		Content:   m.Content.ToDto(),
		CreatedAt: m.CreatedAt,
	}
}

func (m *MessageDto) FromDto() Message {
	return Message{
		Sender:    m.Sender,
		Receiver:  m.Receiver,
		Content:   m.Content.FromDto(),
		CreatedAt: m.CreatedAt,
	}
}

//easyjson:json
type MessageContent struct {
	Text        string   `json:"text"`
	FilePath    []string `json:"file_path"`
	StickerPath string   `json:"sticker_path"`
}

//easyjson:skip
type MessageContentDto struct {
	Text        string
	FilePath    string
	StickerPath string
}

func (mc *MessageContent) ToDto() MessageContentDto {
	files := make([]string, 0, len(mc.FilePath))
	for _, f := range mc.FilePath {
		if f == "" {
			continue
		}
		files = append(files, f)
	}

	return MessageContentDto{
		Text:        mc.Text,
		FilePath:    strings.Join(files, "||;||"),
		StickerPath: mc.StickerPath,
	}
}

func (mc *MessageContentDto) FromDto() MessageContent {
	files := strings.Split(mc.FilePath, "||;||")
	contentFiles := make([]string, 0, len(files))
	for _, f := range files {
		if f == "" {
			continue
		}
		contentFiles = append(contentFiles, f)
	}

	if len(contentFiles) == 0 {
		contentFiles = nil
	}

	return MessageContent{
		Text:        mc.Text,
		FilePath:    contentFiles,
		StickerPath: mc.StickerPath,
	}
}
