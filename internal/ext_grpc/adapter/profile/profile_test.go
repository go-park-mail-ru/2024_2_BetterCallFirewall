package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/profile_api"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mocks struct {
	client *MockProfileServiceClient
}

func getAdapter(ctrl *gomock.Controller) (*GrpcSender, *mocks) {
	m := mocks{
		client: NewMockProfileServiceClient(ctrl),
	}

	return &GrpcSender{client: m.client}, &m
}

var errMock = errors.New("mock error")

func TestGetHeader(t *testing.T) {
	tests := []TableTest[*models.Header, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				res := uint32(0)
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *uint32) (*models.Header, error) {
				return implementation.GetHeader(ctx, *request)
			},
			ExpectedResult: func() (*models.Header, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request *uint32, m *mocks) {
				m.client.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				res := uint32(1)
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *uint32) (*models.Header, error) {
				return implementation.GetHeader(ctx, *request)
			},
			ExpectedResult: func() (*models.Header, error) {
				return &models.Header{AuthorID: 1, Author: "Alexey Zemliakov"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request *uint32, m *mocks) {
				m.client.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(&profile_api.HeaderResponse{Head: &profile_api.Header{AuthorID: 1, Author: "Alexey Zemliakov"}}, nil)
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

func TestGetFriendsID(t *testing.T) {
	tests := []TableTest[[]uint32, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				res := uint32(1)
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *uint32) ([]uint32, error) {
				res, err := implementation.GetFriendsID(ctx, *request)
				return res, err
			},
			ExpectedResult: func() ([]uint32, error) {
				return []uint32{2, 3, 4, 5}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request *uint32, m *mocks) {
				m.client.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).
					Return(&profile_api.FriendsResponse{UserID: []uint32{2, 3, 4, 5}}, nil)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				res := uint32(0)
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *uint32) ([]uint32, error) {
				res, err := implementation.GetFriendsID(ctx, *request)
				return res, err
			},
			ExpectedResult: func() ([]uint32, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request *uint32, m *mocks) {
				m.client.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
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

func TestCreate(t *testing.T) {
	tests := []TableTest[uint32, models.User]{
		{
			name: "1",
			SetupInput: func() (*models.User, error) {
				return &models.User{ID: 0}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *models.User) (uint32, error) {
				res, err := implementation.Create(ctx, request)
				return res, err
			},
			ExpectedResult: func() (uint32, error) {
				return 0, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request *models.User, m *mocks) {
				m.client.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(&profile_api.CreateResponse{}, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*models.User, error) {
				return &models.User{ID: 1}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *models.User) (uint32, error) {
				res, err := implementation.Create(ctx, request)
				return res, err
			},
			ExpectedResult: func() (uint32, error) {
				return 1, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request *models.User, m *mocks) {
				m.client.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(&profile_api.CreateResponse{ID: 1}, nil)
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

func TestGetByEmail(t *testing.T) {
	tests := []TableTest[*models.User, string]{
		{
			name: "1",
			SetupInput: func() (*string, error) {
				res := ""
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *string) (*models.User, error) {
				res, err := implementation.GetByEmail(ctx, *request)
				return res, err
			},
			ExpectedResult: func() (*models.User, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request *string, m *mocks) {
				m.client.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Return(&profile_api.GetByEmailResponse{}, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*string, error) {
				res := "alexey@gmail.ru"
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *string) (*models.User, error) {
				res, err := implementation.GetByEmail(ctx, *request)
				return res, err
			},
			ExpectedResult: func() (*models.User, error) {
				return &models.User{ID: 1, Email: "alexey@gmail.ru"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request *string, m *mocks) {
				m.client.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Return(&profile_api.GetByEmailResponse{User: &profile_api.User{ID: 1, Email: "alexey@gmail.ru"}}, nil)
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

func TestProfileProvider(t *testing.T) {
	_, err := GetProfileProvider("", "")
	assert.Nil(t, err)
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *GrpcSender, *In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(*In, *mocks)
}
