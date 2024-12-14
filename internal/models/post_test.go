package models

import (
	"testing"

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
