package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/auth_api"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mocks struct {
	client *MockAuthServiceClient
}

func getAdapter(ctrl *gomock.Controller) (*GrpcSender, *mocks) {
	m := mocks{
		client: NewMockAuthServiceClient(ctrl),
	}

	return &GrpcSender{client: m.client}, &m
}

var errMock = errors.New("mock error")

func TestCreate(t *testing.T) {
	tests := []TableTest[*models.Session, uint32]{
		{
			name: "1",
			SetupInput: func() (*uint32, error) {
				res := uint32(0)
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *uint32) (*models.Session, error) {
				res, err := implementation.Create(*request)
				return res, err
			},
			ExpectedErr: errMock,
			ExpectedResult: func() (*models.Session, error) {
				return nil, nil
			},
			SetupMock: func(request *uint32, m *mocks) {
				m.client.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*uint32, error) {
				res := uint32(0)
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *uint32) (*models.Session, error) {
				res, err := implementation.Create(*request)
				return res, err
			},
			ExpectedErr: nil,
			ExpectedResult: func() (*models.Session, error) {
				return &models.Session{
					ID:     "session",
					UserID: 1,
				}, nil
			},
			SetupMock: func(request *uint32, m *mocks) {
				m.client.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(&auth_api.CreateResponse{
						Sess: &auth_api.Session{
							ID:     "session",
							UserID: 1,
						},
					}, nil)
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
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expected error %v, got %v", v.ExpectedErr, err)
			}
			assert.Equal(t, res, actual)
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []TableTest[*models.Session, string]{
		{
			name: "1",
			SetupInput: func() (*string, error) {
				res := ""
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *string) (*models.Session, error) {
				res, err := implementation.Check(*request)
				return res, err
			},
			ExpectedErr: errMock,
			ExpectedResult: func() (*models.Session, error) {
				return nil, nil
			},
			SetupMock: func(request *string, m *mocks) {
				m.client.EXPECT().Check(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*string, error) {
				res := "session"
				return &res, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *string) (*models.Session, error) {
				res, err := implementation.Check(*request)
				return res, err
			},
			ExpectedErr: nil,
			ExpectedResult: func() (*models.Session, error) {
				return &models.Session{
					ID:     "session",
					UserID: 1,
				}, nil
			},
			SetupMock: func(request *string, m *mocks) {
				m.client.EXPECT().Check(gomock.Any(), gomock.Any()).
					Return(&auth_api.CheckResponse{
						Sess: &auth_api.Session{
							ID:     "session",
							UserID: 1,
						},
					}, nil)
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
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expected error %v, got %v", v.ExpectedErr, err)
			}
			assert.Equal(t, res, actual)
		})
	}
}

func TestDestroy(t *testing.T) {
	tests := []TableTest[struct{}, models.Session]{
		{
			name: "1",
			SetupInput: func() (*models.Session, error) {
				return &models.Session{}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *models.Session) (struct{}, error) {
				err := implementation.Destroy(request)
				return struct{}{}, err
			},
			ExpectedErr: errMock,
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			SetupMock: func(request *models.Session, m *mocks) {
				m.client.EXPECT().Destroy(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*models.Session, error) {
				return &models.Session{ID: "session", UserID: 1}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *models.Session) (struct{}, error) {
				err := implementation.Destroy(request)
				return struct{}{}, err
			},
			ExpectedErr: nil,
			ExpectedResult: func() (struct{}, error) {
				return struct{}{}, nil
			},
			SetupMock: func(request *models.Session, m *mocks) {
				m.client.EXPECT().Destroy(gomock.Any(), gomock.Any()).
					Return(nil, nil)
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
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expected error %v, got %v", v.ExpectedErr, err)
			}
			assert.Equal(t, res, actual)
		})
	}
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *GrpcSender, *In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(*In, *mocks)
}
