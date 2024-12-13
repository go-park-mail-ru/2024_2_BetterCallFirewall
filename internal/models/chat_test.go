package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCaseMessageContent struct {
	content    MessageContent
	contentDto MessageContentDto
}

func TestFromDtoMessageContent(t *testing.T) {
	tests := []TestCaseMessageContent{
		{content: MessageContent{}, contentDto: MessageContentDto{}},
		{content: MessageContent{Text: "text"}, contentDto: MessageContentDto{Text: "text"}},
		{content: MessageContent{FilePath: []string{"image"}}, contentDto: MessageContentDto{FilePath: "image"}},
		{
			content:    MessageContent{FilePath: []string{"image", "second image"}},
			contentDto: MessageContentDto{FilePath: "image||;||second image"},
		},
	}

	for _, test := range tests {
		res := test.contentDto.FromDto()
		assert.Equal(t, test.content, res)
	}
}

func TestToDtoMessageContent(t *testing.T) {
	tests := []TestCaseMessageContent{
		{content: MessageContent{}, contentDto: MessageContentDto{}},
		{content: MessageContent{Text: "text"}, contentDto: MessageContentDto{Text: "text"}},
		{content: MessageContent{FilePath: []string{"image"}}, contentDto: MessageContentDto{FilePath: "image"}},
		{
			content:    MessageContent{FilePath: []string{"image", "second image"}},
			contentDto: MessageContentDto{FilePath: "image||;||second image"},
		},
		{
			content:    MessageContent{FilePath: []string{""}},
			contentDto: MessageContentDto{FilePath: ""},
		},
	}

	for _, test := range tests {
		res := test.content.ToDto()
		assert.Equal(t, test.contentDto, res)
	}
}

type TestCaseMessage struct {
	message    Message
	messageDto MessageDto
}

func TestMessageFromDto(t *testing.T) {
	tests := []TestCaseMessage{
		{message: Message{}, messageDto: MessageDto{}},
		{
			message:    Message{Content: MessageContent{FilePath: []string{"image"}}},
			messageDto: MessageDto{Content: MessageContentDto{FilePath: "image"}},
		},
	}

	for _, test := range tests {
		res := test.messageDto.FromDto()
		assert.Equal(t, test.message, res)
	}
}

func TestMessageToDto(t *testing.T) {
	tests := []TestCaseMessage{
		{message: Message{}, messageDto: MessageDto{}},
		{
			message:    Message{Content: MessageContent{FilePath: []string{"image"}}},
			messageDto: MessageDto{Content: MessageContentDto{FilePath: "image"}},
		},
	}

	for _, test := range tests {
		res := test.message.ToDto()
		assert.Equal(t, test.messageDto, res)
	}
}
