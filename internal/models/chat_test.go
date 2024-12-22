package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

//easyjson:skip
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

//easyjson:skip
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

func TestMarshal(t *testing.T) {
	createTime := time.Time{}
	m := &Message{
		Content: MessageContent{Text: "new message", FilePath: []string{"image"}}, Sender: 1, Receiver: 2,
		CreatedAt: createTime,
	}

	want := []byte(`{"sender":1,"receiver":2,"content":{"text":"new message","file_path":["image"],"sticker_path":""},"created_at":"0001-01-01T00:00:00Z"}`)
	res, err := easyjson.Marshal(m)
	assert.NoError(t, err)
	assert.Equal(t, string(want), string(res))

	res, err = json.Marshal(m)
	assert.NoError(t, err)
	assert.Equal(t, string(want), string(res))

}

func TestUnmarshal(t *testing.T) {
	sl := []byte(`{"sender":1,"receiver":2,"content":{"text":"new message","file_path":["image"],"sticker_path":""},"created_at":"0001-01-01T00:00:00Z"}`)
	m := &Message{}
	err := easyjson.Unmarshal(sl, m)
	assert.NoError(t, err)
	createTime := time.Time{}
	want := &Message{
		Content: MessageContent{Text: "new message", FilePath: []string{"image"}}, Sender: 1, Receiver: 2,
		CreatedAt: createTime,
	}

	assert.Equal(t, want, m)
}

func TestMarshalChat(t *testing.T) {
	c := &Chat{
		LastMessage: "message",
		LastDate:    time.Time{},
		Receiver: Header{
			Author: "Andrew Savvateev",
		},
	}
	want := []byte(`{"last_message":"message","last_date":"0001-01-01T00:00:00Z","receiver":{"author_id":0,"community_id":0,"author":"Andrew Savvateev","avatar":""}}`)
	res, err := easyjson.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, string(want), string(res))

	res, err = json.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, string(want), string(res))
}

func TestUnmarshalChat(t *testing.T) {
	sl := []byte(`{"last_message":"message","last_date":"0001-01-01T00:00:00Z","receiver":{"author_id":0,"community_id":0,"author":"Andrew Savvateev","avatar":""}}`)
	c := &Chat{}
	want := &Chat{
		LastMessage: "message",
		LastDate:    time.Time{},
		Receiver: Header{
			Author: "Andrew Savvateev",
		},
	}

	err := easyjson.Unmarshal(sl, c)
	assert.NoError(t, err)
	assert.Equal(t, want, c)

	err = json.Unmarshal(sl, c)
	assert.NoError(t, err)
	assert.Equal(t, want, c)
}
