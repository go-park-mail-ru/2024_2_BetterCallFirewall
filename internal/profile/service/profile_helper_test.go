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
	repo *Mockrepository
}

func getServiceHelper(ctrl *gomock.Controller) (*ProfileHelper, *mocksHelper) {
	m := &mocksHelper{
		repo: NewMockrepository(ctrl),
	}

	return NewProfileHelper(m.repo), m
}

var errMock = errors.New("mock error")

func TestNewProfileHelper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	service, _ := getServiceHelper(ctrl)
	assert.NotNil(t, service)
}

func TestCreate(t *testing.T) {
	tests := []TableTest[uint32, models.User]{
		{
			name: "1",
			SetupInput: func() (*models.User, error) {
				return &models.User{}, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHelper, request models.User) (uint32, error) {
				return implementation.Create(ctx, &request)
			},
			ExpectedResult: func() (uint32, error) {
				return 0, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request models.User, m *mocksHelper) {
				m.repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint32(0), errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*models.User, error) {
				return &models.User{}, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHelper, request models.User) (uint32, error) {
				return implementation.Create(ctx, &request)
			},
			ExpectedResult: func() (uint32, error) {
				return 1, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request models.User, m *mocksHelper) {
				m.repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
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
		})
	}
}

func TestGetByEmail(t *testing.T) {
	tests := []TableTest[*models.User, string]{
		{
			name: "1",
			SetupInput: func() (*string, error) {
				r := "email"
				return &r, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHelper, request string) (*models.User, error) {
				return implementation.GetByEmail(ctx, request)
			},
			ExpectedResult: func() (*models.User, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request string, m *mocksHelper) {
				m.repo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*string, error) {
				r := "email@email.ru"
				return &r, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHelper, request string) (*models.User, error) {
				return implementation.GetByEmail(ctx, request)
			},
			ExpectedResult: func() (*models.User, error) {
				return &models.User{
					ID:    1,
					Email: "email@email.ru",
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request string, m *mocksHelper) {
				m.repo.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).Return(&models.User{ID: 1, Email: "email@email.ru"}, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
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
		})
	}
}

func TestGetHeaderHelper(t *testing.T) {
	tests := []TableTest[*models.Header, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				r := uint32(0)
				return &r, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHelper, request uint32) (*models.Header, error) {
				return implementation.GetHeader(ctx, request)
			},
			ExpectedResult: func() (*models.Header, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request uint32, m *mocksHelper) {
				m.repo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				r := uint32(1)
				return &r, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHelper, request uint32) (*models.Header, error) {
				return implementation.GetHeader(ctx, request)
			},
			ExpectedResult: func() (*models.Header, error) {
				return &models.Header{
					AuthorID: 1,
					Author:   "Alexey Zemliakov",
				}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request uint32, m *mocksHelper) {
				m.repo.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(&models.Header{AuthorID: 1, Author: "Alexey Zemliakov"}, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
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
		})
	}
}

func TestGetFriendsID(t *testing.T) {
	tests := []TableTest[[]uint32, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				r := uint32(0)
				return &r, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHelper, request uint32) ([]uint32, error) {
				return implementation.GetFriendsID(ctx, request)
			},
			ExpectedResult: func() ([]uint32, error) {
				return nil, nil
			},
			ExpectedErr: errMock,
			SetupMock: func(request uint32, m *mocksHelper) {
				m.repo.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				r := uint32(2)
				return &r, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHelper, request uint32) ([]uint32, error) {
				return implementation.GetFriendsID(ctx, request)
			},
			ExpectedResult: func() ([]uint32, error) {
				return []uint32{1, 3, 4, 5}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request uint32, m *mocksHelper) {
				m.repo.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).Return([]uint32{1, 3, 4, 5}, nil)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
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
		})
	}
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *ProfileHelper, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocksHelper)
}
