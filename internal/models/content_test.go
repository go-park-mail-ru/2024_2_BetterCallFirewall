package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
