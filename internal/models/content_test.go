package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

//easyjson:skip
type TestCase struct {
	content    Content
	contentDto ContentDto
}

func TestFromDto(t *testing.T) {
	tests := []TestCase{
		{content: Content{}, contentDto: ContentDto{}},
		{content: Content{Text: "text"}, contentDto: ContentDto{Text: "text"}},
		{content: Content{File: []Picture{"image"}}, contentDto: ContentDto{File: "image"}},
		{
			content:    Content{File: []Picture{"image", "second image"}},
			contentDto: ContentDto{File: "image||;||second image"},
		},
	}

	for _, test := range tests {
		res := test.contentDto.FromDto()
		assert.Equal(t, test.content, res)
	}
}

func TestToDto(t *testing.T) {
	tests := []TestCase{
		{content: Content{}, contentDto: ContentDto{}},
		{content: Content{Text: "text"}, contentDto: ContentDto{Text: "text"}},
		{content: Content{File: []Picture{"image"}}, contentDto: ContentDto{File: "image"}},
		{
			content:    Content{File: []Picture{"image", "second image"}},
			contentDto: ContentDto{File: "image||;||second image"},
		},
		{
			content:    Content{File: []Picture{""}},
			contentDto: ContentDto{File: ""},
		},
	}

	for _, test := range tests {
		res := test.content.ToDto()
		assert.Equal(t, test.contentDto, res)
	}
}

func TestMarshalJson(t *testing.T) {
	c := &Content{
		Text:      "comment",
		File:      []Picture{Picture("image")},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	want := []byte(`{"text":"comment","file":["image"],"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`)

	res, err := easyjson.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, want, res)

	res, err = json.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, want, res)
}

func TestUnmarshallJson(t *testing.T) {
	sl := []byte(`{"text":"comment","file":["image"],"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`)
	c := &Content{}
	want := &Content{
		Text:      "comment",
		File:      []Picture{Picture("image")},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	err := easyjson.Unmarshal(sl, c)
	assert.NoError(t, err)
	assert.Equal(t, want, c)

	err = json.Unmarshal(sl, want)
	assert.NoError(t, err)
	assert.Equal(t, want, c)
}
