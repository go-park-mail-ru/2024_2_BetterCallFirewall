package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

var (
	errMock = errors.New("err")
	posts   = []*models.Post{{PostContent: models.Content{Text: "my post"}}}
)

type mockProfileDB struct{}

func (m mockProfileDB) GetAuthorPosts(ctx context.Context, header *models.Header) ([]*models.Post, error) {
	if header == nil {
		return nil, errMock
	}

	return posts, nil
}

type TestCase struct {
	header   *models.Header
	wantPost []*models.Post
	wantErr  error
}

func TestGetAuthorsPosts(t *testing.T) {
	tests := []TestCase{
		{header: nil, wantPost: nil, wantErr: errMock},
		{header: &models.Header{}, wantPost: posts, wantErr: nil},
	}

	serv := NewPostProfileImpl(mockProfileDB{})

	for _, test := range tests {
		post, err := serv.GetAuthorsPosts(context.Background(), test.header)
		assert.Equal(t, test.wantPost, post)
		if !errors.Is(err, test.wantErr) {
			t.Errorf("wrong error, expected: %#v, got: %#v", test.wantErr, err)
		}
	}
}
