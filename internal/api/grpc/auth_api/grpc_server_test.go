package auth_api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mocks struct {
	sessionManager *MockSessionManager
}

func getAdapter(ctrl *gomock.Controller) (*Adapter, *mocks) {
	m := &mocks{
		sessionManager: NewMockSessionManager(ctrl),
	}

	return New(m.sessionManager), m
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a, _ := getAdapter(ctrl)
	assert.NotNil(t, a)
}

var errMock = errors.New("mock error")

func TestCreate(t *testing.T) {
	createTime := time.Now().Unix()
	tests := []TableTest[CreateResponse, CreateRequest]{
		{
			name: "1",
			SetupInput: func() (*CreateRequest, error) {
				res := &CreateRequest{}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *CreateRequest) (*CreateResponse, error) {
				return implementation.Create(ctx, request)
			},
			ExpectedResult: func() (*CreateResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.Internal,
			SetupMock: func(request *CreateRequest, m *mocks) {
				m.sessionManager.EXPECT().Create(gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*CreateRequest, error) {
				res := &CreateRequest{UserID: 1}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *CreateRequest) (*CreateResponse, error) {
				return implementation.Create(ctx, request)
			},
			ExpectedResult: func() (*CreateResponse, error) {
				return &CreateResponse{Sess: &Session{ID: "1", UserID: 1, CreatedAt: createTime}}, nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *CreateRequest, m *mocks) {
				m.sessionManager.EXPECT().Create(gomock.Any()).Return(&models.Session{ID: "1", UserID: request.UserID, CreatedAt: createTime}, nil)
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
			assert.Equal(t, status.Code(err), v.ExpectedErrCode)
		})
	}
}

func TestCheck(t *testing.T) {
	createTime := time.Now().Unix()
	tests := []TableTest[CheckResponse, CheckRequest]{
		{
			name: "1",
			SetupInput: func() (*CheckRequest, error) {
				res := &CheckRequest{}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *CheckRequest) (*CheckResponse, error) {
				return implementation.Check(ctx, request)
			},
			ExpectedResult: func() (*CheckResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.Internal,
			SetupMock: func(request *CheckRequest, m *mocks) {
				m.sessionManager.EXPECT().Check(gomock.Any()).Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*CheckRequest, error) {
				res := &CheckRequest{Cookie: "1"}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *CheckRequest) (*CheckResponse, error) {
				return implementation.Check(ctx, request)
			},
			ExpectedResult: func() (*CheckResponse, error) {
				return &CheckResponse{Sess: &Session{ID: "1", UserID: 1, CreatedAt: createTime}}, nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *CheckRequest, m *mocks) {
				m.sessionManager.EXPECT().Check(gomock.Any()).Return(&models.Session{ID: "1", UserID: 1, CreatedAt: createTime}, nil)
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
			assert.Equal(t, status.Code(err), v.ExpectedErrCode)
		})
	}
}

func TestDestroy(t *testing.T) {
	createTime := time.Now().Unix()
	tests := []TableTest[EmptyResponse, DestroyRequest]{
		{
			name: "1",
			SetupInput: func() (*DestroyRequest, error) {
				res := &DestroyRequest{Sess: &Session{ID: "1", UserID: 0, CreatedAt: createTime}}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *DestroyRequest) (*EmptyResponse, error) {
				return implementation.Destroy(ctx, request)
			},
			ExpectedResult: func() (*EmptyResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.Internal,
			SetupMock: func(request *DestroyRequest, m *mocks) {
				m.sessionManager.EXPECT().Destroy(gomock.Any()).Return(errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*DestroyRequest, error) {
				res := &DestroyRequest{Sess: &Session{ID: "1", UserID: 1, CreatedAt: createTime}}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *DestroyRequest) (*EmptyResponse, error) {
				return implementation.Destroy(ctx, request)
			},
			ExpectedResult: func() (*EmptyResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *DestroyRequest, m *mocks) {
				m.sessionManager.EXPECT().Destroy(gomock.Any()).Return(nil)
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
			assert.Equal(t, status.Code(err), v.ExpectedErrCode)
		})
	}
}

type TableTest[T, In any] struct {
	name            string
	SetupInput      func() (*In, error)
	Run             func(context.Context, *Adapter, *In) (*T, error)
	ExpectedResult  func() (*T, error)
	ExpectedErrCode codes.Code
	SetupMock       func(*In, *mocks)
}
