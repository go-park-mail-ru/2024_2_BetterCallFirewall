package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

type commentMocks struct {
	repo        *MockdbI
	profileRepo *MockprofileRepoI
}

func getCommentService(ctrl *gomock.Controller) (*CommentService, *commentMocks) {
	m := &commentMocks{
		repo:        NewMockdbI(ctrl),
		profileRepo: NewMockprofileRepoI(ctrl),
	}

	return NewCommentService(m.repo, m.profileRepo), m
}

type inputCreate struct {
	userID  uint32
	postID  uint32
	comment *models.Content
}

func TestNewCommentService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _ := getCommentService(ctrl)
	assert.NotNil(t, service)
}

func TestCommentService_Comment(t *testing.T) {
	tests := []TableTest3[*models.Comment, inputCreate]{
		{
			name: "1",
			SetupInput: func() (*inputCreate, error) {
				return &inputCreate{}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputCreate,
			) (*models.Comment, error) {
				return implementation.Comment(ctx, request.userID, request.postID, request.comment)
			},
			ExpectedResult: func() (*models.Comment, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request inputCreate, m *commentMocks) {
				m.repo.EXPECT().CreateComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(uint32(0), errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*inputCreate, error) {
				return &inputCreate{}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputCreate,
			) (*models.Comment, error) {
				return implementation.Comment(ctx, request.userID, request.postID, request.comment)
			},
			ExpectedResult: func() (*models.Comment, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request inputCreate, m *commentMocks) {
				m.repo.EXPECT().CreateComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(uint32(1), nil)
				m.profileRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "3",
			SetupInput: func() (*inputCreate, error) {
				return &inputCreate{
					userID: 1,
					postID: 1,
					comment: &models.Content{
						Text: "new comment",
					},
				}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputCreate,
			) (*models.Comment, error) {
				return implementation.Comment(ctx, request.userID, request.postID, request.comment)
			},
			ExpectedResult: func() (*models.Comment, error) {
				return &models.Comment{
					ID: 1,
					Content: models.Content{
						Text: "new comment",
					},
					Header: models.Header{
						AuthorID: 1,
						Author:   "Alexey Zemliakov",
					},
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request inputCreate, m *commentMocks) {
				m.repo.EXPECT().CreateComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(uint32(1), nil)
				m.profileRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(
						&models.Header{
							AuthorID: 1,
							Author:   "Alexey Zemliakov",
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

				serv, mock := getCommentService(ctrl)
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

type inputDelete struct {
	userID    uint32
	commentID uint32
}

func TestCommentService_Delete(t *testing.T) {
	tests := []TableTest3[struct{}, inputDelete]{
		{
			name: "1",
			SetupInput: func() (*inputDelete, error) {
				return &inputDelete{}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputDelete,
			) (struct{}, error) {
				return struct{}{}, implementation.DeleteComment(ctx, request.userID, request.commentID)
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request inputDelete, m *commentMocks) {
				m.repo.EXPECT().GetCommentAuthor(gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*inputDelete, error) {
				return &inputDelete{
					userID: 1,
				}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputDelete,
			) (struct{}, error) {
				return struct{}{}, implementation.DeleteComment(ctx, request.commentID, request.userID)
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: my_err.ErrAccessDenied,
			SetupMock: func(request inputDelete, m *commentMocks) {
				m.repo.EXPECT().GetCommentAuthor(gomock.Any(), gomock.Any()).Return(uint32(0), nil)
			},
		},
		{
			name: "3",
			SetupInput: func() (*inputDelete, error) {
				return &inputDelete{
					userID: 1,
				}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputDelete,
			) (struct{}, error) {
				return struct{}{}, implementation.DeleteComment(ctx, request.commentID, request.userID)
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request inputDelete, m *commentMocks) {
				m.repo.EXPECT().GetCommentAuthor(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.repo.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(errMock)
			},
		},
		{
			name: "4",
			SetupInput: func() (*inputDelete, error) {
				return &inputDelete{
					userID: 1,
				}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputDelete,
			) (struct{}, error) {
				return struct{}{}, implementation.DeleteComment(ctx, request.commentID, request.userID)
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request inputDelete, m *commentMocks) {
				m.repo.EXPECT().GetCommentAuthor(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.repo.EXPECT().DeleteComment(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				serv, mock := getCommentService(ctrl)
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

type inputEdit struct {
	userID    uint32
	commentID uint32
	comment   *models.Content
}

func TestCommentService_Edit(t *testing.T) {
	tests := []TableTest3[struct{}, inputEdit]{
		{
			name: "1",
			SetupInput: func() (*inputEdit, error) {
				return &inputEdit{}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputEdit,
			) (struct{}, error) {
				return struct{}{}, implementation.EditComment(ctx, request.commentID, request.userID, request.comment)
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request inputEdit, m *commentMocks) {
				m.repo.EXPECT().GetCommentAuthor(gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*inputEdit, error) {
				return &inputEdit{
					userID: 10,
				}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputEdit,
			) (struct{}, error) {
				return struct{}{}, implementation.EditComment(ctx, request.commentID, request.userID, request.comment)
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: my_err.ErrAccessDenied,
			SetupMock: func(request inputEdit, m *commentMocks) {
				m.repo.EXPECT().GetCommentAuthor(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
			},
		},
		{
			name: "3",
			SetupInput: func() (*inputEdit, error) {
				return &inputEdit{
					userID: 10,
				}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputEdit,
			) (struct{}, error) {
				return struct{}{}, implementation.EditComment(ctx, request.commentID, request.userID, request.comment)
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request inputEdit, m *commentMocks) {
				m.repo.EXPECT().GetCommentAuthor(gomock.Any(), gomock.Any()).Return(uint32(10), nil)
				m.repo.EXPECT().UpdateComment(gomock.Any(), gomock.Any(), gomock.Any()).Return(errMock)
			},
		},
		{
			name: "4",
			SetupInput: func() (*inputEdit, error) {
				return &inputEdit{
					userID: 10,
				}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputEdit,
			) (struct{}, error) {
				return struct{}{}, implementation.EditComment(ctx, request.commentID, request.userID, request.comment)
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request inputEdit, m *commentMocks) {
				m.repo.EXPECT().GetCommentAuthor(gomock.Any(), gomock.Any()).Return(uint32(10), nil)
				m.repo.EXPECT().UpdateComment(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				serv, mock := getCommentService(ctrl)
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

type inputGet struct {
	postID uint32
	lastID uint32
	newest bool
}

func TestCommentService_GetComments(t *testing.T) {
	tests := []TableTest3[[]*models.Comment, inputGet]{
		{
			name: "1",
			SetupInput: func() (*inputGet, error) {
				return &inputGet{}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputGet,
			) ([]*models.Comment, error) {
				return implementation.GetComments(ctx, request.postID, request.lastID, request.newest)
			},
			ExpectedResult: func() ([]*models.Comment, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request inputGet, m *commentMocks) {
				m.repo.EXPECT().GetComments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*inputGet, error) {
				return &inputGet{}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputGet,
			) ([]*models.Comment, error) {
				return implementation.GetComments(ctx, request.postID, request.lastID, request.newest)
			},
			ExpectedResult: func() ([]*models.Comment, error) {
				return []*models.Comment{
					{
						ID:     1,
						Header: models.Header{AuthorID: 1, Author: "Alexey Zemliakov"},
					},
					{
						ID:     2,
						Header: models.Header{AuthorID: 1, Author: "Alexey Zemliakov"},
					},
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request inputGet, m *commentMocks) {
				m.repo.EXPECT().GetComments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(
						[]*models.Comment{
							{ID: 1},
							{ID: 2},
						},
						nil,
					)
				m.profileRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(
						&models.Header{
							AuthorID: 1,
							Author:   "Alexey Zemliakov",
						},
						nil,
					).AnyTimes()
			},
		},
		{
			name: "3",
			SetupInput: func() (*inputGet, error) {
				return &inputGet{}, nil
			},
			Run: func(
				ctx context.Context, implementation *CommentService, request inputGet,
			) ([]*models.Comment, error) {
				return implementation.GetComments(ctx, request.postID, request.lastID, request.newest)
			},
			ExpectedResult: func() ([]*models.Comment, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request inputGet, m *commentMocks) {
				m.repo.EXPECT().GetComments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(
						[]*models.Comment{
							{ID: 1},
							{ID: 2},
						},
						nil,
					)
				m.profileRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(nil, errMock).AnyTimes()
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				serv, mock := getCommentService(ctrl)
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

type TableTest3[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *CommentService, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *commentMocks)
}
