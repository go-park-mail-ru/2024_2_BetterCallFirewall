package post

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/post_api"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mocks struct {
	client *MockPostServiceClient
}

func getAdapter(ctrl *gomock.Controller) (*GrpcSender, *mocks) {
	m := mocks{
		client: NewMockPostServiceClient(ctrl),
	}

	return &GrpcSender{client: m.client}, &m
}

var errMock = errors.New("mock error")

func TestGetAuthorsPost(t *testing.T) {
	createTime := time.Now()
	tests := []TableTest[[]*models.Post, *models.Header]{
		{
			name: "1",
			SetupInput: func() (*models.Header, error) {
				return &models.Header{}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *models.Header) ([]*models.Post, error) {
				res, err := implementation.GetAuthorsPosts(ctx, request, 0)
				return res, err
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request *models.Header, m *mocks) {
				m.client.EXPECT().GetAuthorsPosts(gomock.Any(), gomock.Any()).
					Return(&post_api.Response{}, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*models.Header, error) {
				return &models.Header{AuthorID: 1}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *models.Header) ([]*models.Post, error) {
				res, err := implementation.GetAuthorsPosts(ctx, request, 1)
				return res, err
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return []*models.Post{
					{
						ID: 1,
						PostContent: models.Content{
							Text:      "new post",
							CreatedAt: time.Unix(createTime.Unix(), 0),
							UpdatedAt: time.Unix(createTime.Unix(), 0),
							File:      []models.Picture{},
						},
						Header: models.Header{
							AuthorID: 1,
						},
					},
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request *models.Header, m *mocks) {
				m.client.EXPECT().GetAuthorsPosts(gomock.Any(), gomock.Any()).
					Return(
						&post_api.Response{
							Posts: []*post_api.Post{
								{
									ID: 1,
									PostContent: &post_api.Content{
										Text:      "new post",
										CreatedAt: createTime.Unix(),
										UpdatedAt: createTime.Unix(),
									},
									Head: &post_api.Header{
										AuthorID: 1,
									},
								},
							},
						}, nil,
					)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
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
			},
		)
	}
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (In, error)
	Run            func(context.Context, *GrpcSender, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
