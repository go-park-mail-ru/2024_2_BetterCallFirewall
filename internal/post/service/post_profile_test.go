package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mocksHelper struct {
	repo *MockPostProfileDB
}

func getServiceHelper(ctrl *gomock.Controller) (*PostProfileImpl, *mocksHelper) {
	m := &mocksHelper{
		repo: NewMockPostProfileDB(ctrl),
	}

	return NewPostProfileImpl(m.repo), m
}

type input struct {
	header *models.Header
	userID uint32
}

func TestGetAuthorsPost(t *testing.T) {
	tests := []TableTest2[[]*models.Post, input]{
		{
			name: "1",
			SetupInput: func() (*input, error) {
				return &input{}, nil
			},
			Run: func(ctx context.Context, implementation *PostProfileImpl, request input) ([]*models.Post, error) {
				return implementation.GetAuthorsPosts(ctx, request.header, request.userID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request input, m *mocksHelper) {
				m.repo.EXPECT().GetAuthorPosts(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*input, error) {
				return &input{}, nil
			},
			Run: func(ctx context.Context, implementation *PostProfileImpl, request input) ([]*models.Post, error) {
				return implementation.GetAuthorsPosts(ctx, request.header, request.userID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request input, m *mocksHelper) {
				m.repo.EXPECT().GetAuthorPosts(gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{
							ID: 1,
						},
					}, nil,
				)
				m.repo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
			},
		},
		{
			name: "3",
			SetupInput: func() (*input, error) {
				return &input{}, nil
			},
			Run: func(ctx context.Context, implementation *PostProfileImpl, request input) ([]*models.Post, error) {
				return implementation.GetAuthorsPosts(ctx, request.header, request.userID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request input, m *mocksHelper) {
				m.repo.EXPECT().GetAuthorPosts(gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{
							ID: 1,
						},
					}, nil,
				)
				m.repo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.repo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errMock)
			},
		},
		{
			name: "4",
			SetupInput: func() (*input, error) {
				return &input{}, nil
			},
			Run: func(ctx context.Context, implementation *PostProfileImpl, request input) ([]*models.Post, error) {
				return implementation.GetAuthorsPosts(ctx, request.header, request.userID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return []*models.Post{
					{
						ID:           1,
						IsLiked:      true,
						LikesCount:   1,
						CommentCount: 1,
					},
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request input, m *mocksHelper) {
				m.repo.EXPECT().GetAuthorPosts(gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{
							ID: 1,
						},
					}, nil,
				)
				m.repo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(1), nil).AnyTimes()
				m.repo.EXPECT().GetCommentCount(gomock.Any(), gomock.Any()).Return(uint32(1), nil).AnyTimes()
				m.repo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
		{
			name: "5",
			SetupInput: func() (*input, error) {
				return &input{}, nil
			},
			Run: func(ctx context.Context, implementation *PostProfileImpl, request input) ([]*models.Post, error) {
				return implementation.GetAuthorsPosts(ctx, request.header, request.userID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request input, m *mocksHelper) {
				m.repo.EXPECT().GetAuthorPosts(gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{
							ID: 1,
						},
					}, nil,
				)
				m.repo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.repo.EXPECT().GetCommentCount(gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
				m.repo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				serv, mock := getServiceHelper(ctrl)
				ctx := context.Background()

				input, err := v.SetupInput()
				if err != nil {
					t.Error(err)
				}

				v.SetupMock(*input, mock)

				res, err := v.ExpectedResult()
				if err != nil {
					t.Error(err)
				}

				actual, err := v.Run(ctx, serv, *input)
				assert.Equal(t, res, actual)
				if !errors.Is(err, v.ExpectedErr) {
					t.Errorf("expect %v, got %v", v.ExpectedErr, err)
				}
			},
		)
	}
}

func TestNewPostProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _ := getServiceHelper(ctrl)
	assert.NotNil(t, service)
}

type TableTest2[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *PostProfileImpl, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocksHelper)
}
