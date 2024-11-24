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

type mocks struct {
	postRepo      *MockDB
	communityRepo *MockCommunityRepo
	profileRepo   *MockProfileRepo
}

func getService(ctrl *gomock.Controller) (*PostServiceImpl, *mocks) {
	m := &mocks{
		postRepo:      NewMockDB(ctrl),
		communityRepo: NewMockCommunityRepo(ctrl),
		profileRepo:   NewMockProfileRepo(ctrl),
	}

	return NewPostServiceImpl(m.postRepo, m.profileRepo, m.communityRepo), m
}

func TestNewPostService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _ := getService(ctrl)
	assert.NotNil(t, service)
}

var errMock = errors.New("mock error")

func TestCreate(t *testing.T) {
	tests := []TableTest[uint32, models.Post]{
		{
			name: "1",
			SetupInput: func() (*models.Post, error) {
				return &models.Post{PostContent: models.Content{Text: "new post"}}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request models.Post) (uint32, error) {
				return implementation.Create(ctx, &request)
			},
			ExpectedResult: func() (uint32, error) {
				return 0, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request models.Post, m *mocks) {
				m.postRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*models.Post, error) {
				return &models.Post{PostContent: models.Content{Text: "new real post"}}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request models.Post) (uint32, error) {
				return implementation.Create(ctx, &request)
			},
			ExpectedResult: func() (uint32, error) {
				return 1, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request models.Post, m *mocks) {
				m.postRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

type userAndPostIDs struct {
	userID uint32
	postId uint32
}

func TestGet(t *testing.T) {
	tests := []TableTest[*models.Post, userAndPostIDs]{
		{
			name: "1",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{postId: 0, userID: 0}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (*models.Post, error) {
				return implementation.Get(ctx, request.postId, request.userID)
			},
			ExpectedResult: func() (*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{postId: 1, userID: 1}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (*models.Post, error) {
				return implementation.Get(ctx, request.postId, request.userID)
			},
			ExpectedResult: func() (*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&models.Post{
						Header: models.Header{
							CommunityID: 0,
							AuthorID:    1,
						},
					},
					nil)
				m.profileRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "3",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{postId: 1, userID: 1}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (*models.Post, error) {
				return implementation.Get(ctx, request.postId, request.userID)
			},
			ExpectedResult: func() (*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&models.Post{
						Header: models.Header{
							CommunityID: 1,
							AuthorID:    0,
						},
					},
					nil)
				m.communityRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "4",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{postId: 1, userID: 1}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (*models.Post, error) {
				return implementation.Get(ctx, request.postId, request.userID)
			},
			ExpectedResult: func() (*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&models.Post{
						Header: models.Header{
							CommunityID: 0,
							AuthorID:    1,
						},
					},
					nil)
				m.profileRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(&models.Header{
					CommunityID: 0,
					AuthorID:    1,
					Author:      "user",
				}, nil)
				m.postRepo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
			},
		},
		{
			name: "5",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{postId: 1, userID: 1}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (*models.Post, error) {
				return implementation.Get(ctx, request.postId, request.userID)
			},
			ExpectedResult: func() (*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&models.Post{
						Header: models.Header{
							CommunityID: 1,
							AuthorID:    0,
						},
					},
					nil)
				m.communityRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(&models.Header{
					CommunityID: 1,
					AuthorID:    0,
					Author:      "community",
				}, nil)
				m.postRepo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postRepo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errMock)
			},
		},
		{
			name: "6",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{postId: 1, userID: 1}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (*models.Post, error) {
				return implementation.Get(ctx, request.postId, request.userID)
			},
			ExpectedResult: func() (*models.Post, error) {
				return &models.Post{
					Header: models.Header{
						CommunityID: 1,
						AuthorID:    0,
						Author:      "community",
					},
					LikesCount: 1,
					IsLiked:    true,
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					&models.Post{
						Header: models.Header{
							CommunityID: 1,
							AuthorID:    0,
						},
					},
					nil)
				m.communityRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(&models.Header{
					CommunityID: 1,
					AuthorID:    0,
					Author:      "community",
				}, nil)
				m.postRepo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postRepo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []TableTest[struct{}, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				in := uint32(0)
				return &in, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request uint32) (struct{}, error) {
				err := implementation.Delete(ctx, request)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request uint32, m *mocks) {
				m.postRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				in := uint32(1)
				return &in, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request uint32) (struct{}, error) {
				err := implementation.Delete(ctx, request)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request uint32, m *mocks) {
				m.postRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []TableTest[struct{}, models.Post]{
		{
			name: "1",
			SetupInput: func() (*models.Post, error) {
				return &models.Post{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request models.Post) (struct{}, error) {
				err := implementation.Update(ctx, &request)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request models.Post, m *mocks) {
				m.postRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*models.Post, error) {
				return &models.Post{ID: 1, PostContent: models.Content{Text: "New post"}}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request models.Post) (struct{}, error) {
				err := implementation.Update(ctx, &request)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request models.Post, m *mocks) {
				m.postRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

type userAndLastIDs struct {
	UserID uint32
	LastId uint32
}

func TestGetBatch(t *testing.T) {
	tests := []TableTest[[]*models.Post, userAndLastIDs]{
		{
			name: "1",
			SetupInput: func() (*userAndLastIDs, error) {
				return &userAndLastIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndLastIDs) ([]*models.Post, error) {
				return implementation.GetBatch(ctx, request.LastId, request.UserID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndLastIDs, m *mocks) {
				m.postRepo.EXPECT().GetPosts(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*userAndLastIDs, error) {
				return &userAndLastIDs{UserID: 1, LastId: 2}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndLastIDs) ([]*models.Post, error) {
				return implementation.GetBatch(ctx, request.LastId, request.UserID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndLastIDs, m *mocks) {
				m.postRepo.EXPECT().GetPosts(gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{ID: 1, Header: models.Header{CommunityID: 1}},
					}, nil)
				m.communityRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "3",
			SetupInput: func() (*userAndLastIDs, error) {
				return &userAndLastIDs{UserID: 1, LastId: 2}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndLastIDs) ([]*models.Post, error) {
				return implementation.GetBatch(ctx, request.LastId, request.UserID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return []*models.Post{
					{
						ID:         1,
						Header:     models.Header{AuthorID: 1},
						IsLiked:    true,
						LikesCount: 1,
					},
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request userAndLastIDs, m *mocks) {
				m.postRepo.EXPECT().GetPosts(gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{ID: 1, Header: models.Header{AuthorID: 1}},
					}, nil)
				m.profileRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(&models.Header{AuthorID: 1}, nil)
				m.postRepo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postRepo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

func TestGetBatchFromFriend(t *testing.T) {
	tests := []TableTest[[]*models.Post, userAndLastIDs]{
		{
			name: "1",
			SetupInput: func() (*userAndLastIDs, error) {
				return &userAndLastIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndLastIDs) ([]*models.Post, error) {
				return implementation.GetBatchFromFriend(ctx, request.LastId, request.UserID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndLastIDs, m *mocks) {
				m.profileRepo.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*userAndLastIDs, error) {
				return &userAndLastIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndLastIDs) ([]*models.Post, error) {
				return implementation.GetBatchFromFriend(ctx, request.LastId, request.UserID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: my_err.ErrNoMoreContent,
			SetupMock: func(request userAndLastIDs, m *mocks) {
				m.profileRepo.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
		},
		{
			name: "3",
			SetupInput: func() (*userAndLastIDs, error) {
				return &userAndLastIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndLastIDs) ([]*models.Post, error) {
				return implementation.GetBatchFromFriend(ctx, request.LastId, request.UserID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndLastIDs, m *mocks) {
				m.profileRepo.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).Return([]uint32{1, 2, 3}, nil)
				m.postRepo.EXPECT().GetFriendsPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "4",
			SetupInput: func() (*userAndLastIDs, error) {
				return &userAndLastIDs{UserID: 1, LastId: 2}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndLastIDs) ([]*models.Post, error) {
				return implementation.GetBatchFromFriend(ctx, request.LastId, request.UserID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndLastIDs, m *mocks) {
				m.profileRepo.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).Return([]uint32{1, 2, 3}, nil)
				m.postRepo.EXPECT().GetFriendsPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{ID: 1, Header: models.Header{CommunityID: 1}},
					}, nil)
				m.communityRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "5",
			SetupInput: func() (*userAndLastIDs, error) {
				return &userAndLastIDs{UserID: 1, LastId: 2}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndLastIDs) ([]*models.Post, error) {
				return implementation.GetBatchFromFriend(ctx, request.LastId, request.UserID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return []*models.Post{
					{
						ID:         1,
						Header:     models.Header{AuthorID: 1},
						IsLiked:    true,
						LikesCount: 1,
					},
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request userAndLastIDs, m *mocks) {
				m.profileRepo.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).Return([]uint32{1, 2, 3}, nil)
				m.postRepo.EXPECT().GetFriendsPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{ID: 1, Header: models.Header{AuthorID: 1}},
					}, nil)
				m.profileRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(&models.Header{AuthorID: 1}, nil)
				m.postRepo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postRepo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

func TestGetPostAuthor(t *testing.T) {
	tests := []TableTest[uint32, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				in := uint32(0)
				return &in, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request uint32) (uint32, error) {
				return implementation.GetPostAuthorID(ctx, request)
			},
			ExpectedResult: func() (uint32, error) {
				return uint32(0), nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request uint32, m *mocks) {
				m.postRepo.EXPECT().GetPostAuthor(gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				in := uint32(1)
				return &in, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request uint32) (uint32, error) {
				return implementation.GetPostAuthorID(ctx, request)
			},
			ExpectedResult: func() (uint32, error) {
				return uint32(3), nil
			},
			ExpectedErr: nil,
			SetupMock: func(request uint32, m *mocks) {
				m.postRepo.EXPECT().GetPostAuthor(gomock.Any(), gomock.Any()).Return(uint32(3), nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

func TestCreateCommunityPost(t *testing.T) {
	tests := []TableTest[uint32, models.Post]{
		{
			name: "1",
			SetupInput: func() (*models.Post, error) {
				return &models.Post{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request models.Post) (uint32, error) {
				return implementation.CreateCommunityPost(ctx, &request)
			},
			ExpectedResult: func() (uint32, error) {
				return uint32(0), nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request models.Post, m *mocks) {
				m.postRepo.EXPECT().CreateCommunityPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*models.Post, error) {
				return &models.Post{PostContent: models.Content{Text: "new post"}}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request models.Post) (uint32, error) {
				return implementation.CreateCommunityPost(ctx, &request)
			},
			ExpectedResult: func() (uint32, error) {
				return uint32(1), nil
			},
			ExpectedErr: nil,
			SetupMock: func(request models.Post, m *mocks) {
				m.postRepo.EXPECT().CreateCommunityPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(uint32(1), nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

type IDs struct {
	communityID uint32
	lastID      uint32
	userID      uint32
}

func TestGetCommunityPost(t *testing.T) {
	tests := []TableTest[[]*models.Post, IDs]{
		{
			name: "1",
			SetupInput: func() (*IDs, error) {
				return &IDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request IDs) ([]*models.Post, error) {
				return implementation.GetCommunityPost(ctx, request.communityID, request.userID, request.lastID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request IDs, m *mocks) {
				m.postRepo.EXPECT().GetCommunityPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*IDs, error) {
				return &IDs{userID: 1, lastID: 2, communityID: 1}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request IDs) ([]*models.Post, error) {
				return implementation.GetCommunityPost(ctx, request.lastID, request.userID, request.communityID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request IDs, m *mocks) {
				m.postRepo.EXPECT().GetCommunityPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{ID: 1, Header: models.Header{CommunityID: 1}},
					}, nil)
				m.communityRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "3",
			SetupInput: func() (*IDs, error) {
				return &IDs{userID: 1, lastID: 2, communityID: 3}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request IDs) ([]*models.Post, error) {
				return implementation.GetCommunityPost(ctx, request.lastID, request.userID, request.communityID)
			},
			ExpectedResult: func() ([]*models.Post, error) {
				return []*models.Post{
					{
						ID:         1,
						Header:     models.Header{AuthorID: 1},
						IsLiked:    true,
						LikesCount: 1,
					},
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request IDs, m *mocks) {
				m.postRepo.EXPECT().GetCommunityPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]*models.Post{
						{ID: 1, Header: models.Header{AuthorID: 1}},
					}, nil)
				m.profileRepo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(&models.Header{AuthorID: 1}, nil)
				m.postRepo.EXPECT().GetLikesOnPost(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postRepo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

type userAndCommunityIDs struct {
	userID      uint32
	communityID uint32
}

func TestCheckAccessToCommunity(t *testing.T) {
	tests := []TableTest[bool, userAndCommunityIDs]{
		{
			name: "1",
			SetupInput: func() (*userAndCommunityIDs, error) {
				return &userAndCommunityIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndCommunityIDs) (bool, error) {
				res := implementation.CheckAccessToCommunity(ctx, request.communityID, request.userID)
				return res, nil
			},
			ExpectedResult: func() (bool, error) {
				return false, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request userAndCommunityIDs, m *mocks) {
				m.communityRepo.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
			},
		},
		{
			name: "1",
			SetupInput: func() (*userAndCommunityIDs, error) {
				return &userAndCommunityIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndCommunityIDs) (bool, error) {
				res := implementation.CheckAccessToCommunity(ctx, request.communityID, request.userID)
				return res, nil
			},
			ExpectedResult: func() (bool, error) {
				return true, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request userAndCommunityIDs, m *mocks) {
				m.communityRepo.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

func TestSetLikeOnPost(t *testing.T) {
	tests := []TableTest[struct{}, userAndPostIDs]{
		{
			name: "1",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (struct{}, error) {
				err := implementation.SetLikeToPost(ctx, request.postId, request.userID)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().SetLikeToPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (struct{}, error) {
				err := implementation.SetLikeToPost(ctx, request.postId, request.userID)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().SetLikeToPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

func TestDeleteLikeFromPost(t *testing.T) {
	tests := []TableTest[struct{}, userAndPostIDs]{
		{
			name: "1",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (struct{}, error) {
				err := implementation.DeleteLikeFromPost(ctx, request.postId, request.userID)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().DeleteLikeFromPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (struct{}, error) {
				err := implementation.DeleteLikeFromPost(ctx, request.postId, request.userID)
				return struct{}{}, err
			},
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().DeleteLikeFromPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

func TestCheckLikes(t *testing.T) {
	tests := []TableTest[bool, userAndPostIDs]{
		{
			name: "1",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (bool, error) {
				return implementation.CheckLikes(ctx, request.postId, request.userID)
			},
			ExpectedResult: func() (bool, error) {
				return false, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*userAndPostIDs, error) {
				return &userAndPostIDs{}, nil
			},
			Run: func(ctx context.Context, implementation *PostServiceImpl, request userAndPostIDs) (bool, error) {
				return implementation.CheckLikes(ctx, request.postId, request.userID)
			},
			ExpectedResult: func() (bool, error) {
				return true, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request userAndPostIDs, m *mocks) {
				m.postRepo.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getService(ctrl)
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
		})
	}
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *PostServiceImpl, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
