package post_api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mocks struct {
	postService *MockPostService
}

func getAdapter(ctrl *gomock.Controller) (*Adapter, *mocks) {
	m := &mocks{
		postService: NewMockPostService(ctrl),
	}

	return New(m.postService), m
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a, _ := getAdapter(ctrl)
	assert.NotNil(t, a)
}

func TestGetAuthorsPosts(t *testing.T) {
	errMock := errors.New("mock error")
	createTime := time.Now()
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				res := &Request{Head: &Header{}}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *Request) (*Response, error) {
				return implementation.GetAuthorsPosts(ctx, request)
			},
			ExpectedResult: func() (*Response, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request *Request, m *mocks) {
				m.postService.EXPECT().GetAuthorsPosts(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*Request, error) {
				res := &Request{Head: &Header{AuthorID: 1, Author: "Alexey Zemliakov"}}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *Request) (*Response, error) {
				return implementation.GetAuthorsPosts(ctx, request)
			},
			ExpectedResult: func() (*Response, error) {
				return &Response{
						Posts: []*Post{
							{
								ID:          1,
								PostContent: &Content{Text: "New Post", CreatedAt: createTime.Unix(), UpdatedAt: createTime.Unix()},
								Head:        &Header{AuthorID: 1, Author: "Alexey Zemliakov"},
							},
						},
					},
					nil
			},
			ExpectedErr: nil,
			SetupMock: func(request *Request, m *mocks) {
				m.postService.EXPECT().GetAuthorsPosts(gomock.Any(), gomock.Any()).
					Return(
						[]*models.Post{
							{
								ID:          1,
								PostContent: models.Content{Text: "New Post", CreatedAt: createTime, UpdatedAt: createTime},
								Header:      models.Header{AuthorID: 1, Author: "Alexey Zemliakov"},
							},
						},
						nil,
					)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			adapter, mock := getAdapter(ctrl)
			ctx := context.Background()

			input, err := v.SetupInput()
			if err != nil {
				t.Error(err)
			}

			v.SetupMock(input, mock)

			res, err := v.ExpectedResult()
			if err != nil {
				t.Error(err)
			}

			actual, err := v.Run(ctx, adapter, input)
			assert.Equal(t, res, actual)
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expect %v, got %v", v.ExpectedErr, err)
			}
		})
	}
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *Adapter, *In) (*T, error)
	ExpectedResult func() (*T, error)
	ExpectedErr    error
	SetupMock      func(*In, *mocks)
}
