package models

import (
	"strings"
	"time"
)

type Content struct {
	Text      string    `json:"text"`
	File      []Picture `json:"file,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Content) ToDto() ContentDto {
	files := make([]string, 0, len(c.File))
	for _, f := range c.File {
		if f == "" {
			continue
		}
		files = append(files, string(f))
	}

	return ContentDto{
		Text:      c.Text,
		File:      Picture(strings.Join(files, "||;||")),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

type ContentDto struct {
	Text      string
	File      Picture
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *ContentDto) FromDto() Content {
	files := strings.Split(string(c.File), "||;||")
	contentFiles := make([]Picture, 0, len(files))
	for _, f := range files {
		if f == "" {
			continue
		}
		contentFiles = append(contentFiles, Picture(f))
	}

	if len(contentFiles) == 0 {
		contentFiles = nil
	}

	return Content{
		Text:      c.Text,
		File:      contentFiles,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
