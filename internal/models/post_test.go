package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
)

//easyjson:skip
type TestCasePost struct {
	post    Post
	postDto PostDto
}

func TestPostFromDto(t *testing.T) {
	tests := []TestCasePost{
		{post: Post{}, postDto: PostDto{}},
		{post: Post{PostContent: Content{Text: "text"}}, postDto: PostDto{PostContent: ContentDto{Text: "text"}}},
		{
			post:    Post{PostContent: Content{File: []Picture{"image"}}},
			postDto: PostDto{PostContent: ContentDto{File: "image"}},
		},
		{
			post:    Post{PostContent: Content{File: []Picture{"image", "second image"}}},
			postDto: PostDto{PostContent: ContentDto{File: "image||;||second image"}},
		},
	}

	for _, test := range tests {
		res := test.postDto.FromDto()
		assert.Equal(t, test.post, res)
	}
}

func TestPostToDto(t *testing.T) {
	tests := []TestCasePost{
		{post: Post{}, postDto: PostDto{}},
		{post: Post{PostContent: Content{Text: "text"}}, postDto: PostDto{PostContent: ContentDto{Text: "text"}}},
		{
			post:    Post{PostContent: Content{File: []Picture{"image"}}},
			postDto: PostDto{PostContent: ContentDto{File: "image"}},
		},
		{
			post:    Post{PostContent: Content{File: []Picture{"image", "second image"}}},
			postDto: PostDto{PostContent: ContentDto{File: "image||;||second image"}},
		},
	}

	for _, test := range tests {
		res := test.post.ToDto()
		assert.Equal(t, test.postDto, res)
	}
}

//easyjson:skip
type TestCaseComment struct {
	comment    Comment
	commentDto CommentDto
}

func TestCommentFromDto(t *testing.T) {
	tests := []TestCaseComment{
		{comment: Comment{}, commentDto: CommentDto{}},
		{
			comment:    Comment{Content: Content{File: []Picture{"image"}}},
			commentDto: CommentDto{Content: ContentDto{File: "image"}},
		},
	}

	for _, test := range tests {
		res := test.commentDto.FromDto()
		assert.Equal(t, test.comment, res)
	}
}

func TestCommentToDto(t *testing.T) {
	tests := []TestCaseComment{
		{comment: Comment{}, commentDto: CommentDto{}},
		{
			comment:    Comment{Content: Content{File: []Picture{"image"}}},
			commentDto: CommentDto{Content: ContentDto{File: "image"}},
		},
	}

	for _, test := range tests {
		res := test.comment.ToDto()
		assert.Equal(t, test.commentDto, res)
	}
}

func TestMarshalPost(t *testing.T) {
	p := &Post{
		ID: 1,
		Header: Header{
			AuthorID:    10,
			CommunityID: 0,
			Author:      "Alexey",
			Avatar:      "/image",
		},
		PostContent: Content{
			Text:      "text",
			File:      []Picture{"image"},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		LikesCount:   1,
		IsLiked:      true,
		CommentCount: 10,
	}
	want := []byte(`{"id":1,"header":{"author_id":10,"community_id":0,"author":"Alexey","avatar":"/image"},"post_content":{"text":"text","file":["image"],"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"likes_count":1,"is_liked":true,"comment_count":10}`)

	res, err := easyjson.Marshal(p)
	assert.NoError(t, err)
	assert.Equal(t, want, res)

	res, err = json.Marshal(p)
	assert.NoError(t, err)
	assert.Equal(t, want, res)
}

func TestUnmarshallPost(t *testing.T) {
	want := &Post{
		ID: 1,
		Header: Header{
			AuthorID:    10,
			CommunityID: 0,
			Author:      "Alexey",
			Avatar:      "/image",
		},
		PostContent: Content{
			Text:      "text",
			File:      []Picture{"image"},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		LikesCount:   1,
		IsLiked:      true,
		CommentCount: 10,
	}
	sl := []byte(`{"id":1,"header":{"author_id":10,"community_id":0,"author":"Alexey","avatar":"/image"},"post_content":{"text":"text","file":["image"],"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"likes_count":1,"is_liked":true,"comment_count":10}`)
	p := &Post{}

	err := easyjson.Unmarshal(sl, p)
	assert.NoError(t, err)
	assert.Equal(t, want, p)

	err = json.Unmarshal(sl, p)
	assert.NoError(t, err)
	assert.Equal(t, want, p)
}

func TestMarshallComment(t *testing.T) {
	c := &Comment{
		ID: 10,
		Header: Header{
			AuthorID:    10,
			CommunityID: 0,
			Author:      "Alexey",
			Avatar:      "/image",
		},
		Content: Content{
			Text:      "comment",
			File:      nil,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		LikesCount: 0,
		IsLiked:    false,
	}
	want := []byte(`{"id":10,"header":{"author_id":10,"community_id":0,"author":"Alexey","avatar":"/image"},"content":{"text":"comment","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"likes_count":0,"is_liked":false}`)

	res, err := easyjson.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, res, want)

	res, err = json.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, want, res)
}

func TestUnmarshallComment(t *testing.T) {
	want := &Comment{
		ID: 10,
		Header: Header{
			AuthorID:    10,
			CommunityID: 0,
			Author:      "Alexey",
			Avatar:      "/image",
		},
		Content: Content{
			Text:      "comment",
			File:      nil,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		LikesCount: 0,
		IsLiked:    false,
	}
	sl := []byte(`{"id":10,"header":{"author_id":10,"community_id":0,"author":"Alexey","avatar":"/image"},"content":{"text":"comment","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"},"likes_count":0,"is_liked":false}`)
	c := &Comment{}

	err := easyjson.Unmarshal(sl, c)
	assert.NoError(t, err)
	assert.Equal(t, want, c)

	err = json.Unmarshal(sl, c)
	assert.NoError(t, err)
	assert.Equal(t, want, c)
}
