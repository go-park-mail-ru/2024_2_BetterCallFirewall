package profile_api

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

type mocks struct {
	profileService *MockprofileService
}

func getAdapter(ctrl *gomock.Controller) (*Adapter, *mocks) {
	m := &mocks{
		profileService: NewMockprofileService(ctrl),
	}

	return New(m.profileService), m
}

var errMock = errors.New("mock error")

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a, _ := getAdapter(ctrl)
	assert.NotNil(t, a)
}

func TestGetHeader(t *testing.T) {
	tests := []TableTest[HeaderResponse, HeaderRequest]{
		{
			name: "1",
			SetupInput: func() (*HeaderRequest, error) {
				res := &HeaderRequest{UserID: 0}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *HeaderRequest) (*HeaderResponse, error) {
				return implementation.GetHeader(ctx, request)
			},
			ExpectedResult: func() (*HeaderResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.Internal,
			SetupMock: func(request *HeaderRequest, m *mocks) {
				m.profileService.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*HeaderRequest, error) {
				res := &HeaderRequest{UserID: 1}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *HeaderRequest) (*HeaderResponse, error) {
				return implementation.GetHeader(ctx, request)
			},
			ExpectedResult: func() (*HeaderResponse, error) {
				return &HeaderResponse{
					Head: &Header{
						AuthorID: 1,
						Author:   "Alexey Zemliakov",
					},
				}, nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *HeaderRequest, m *mocks) {
				m.profileService.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(&models.Header{AuthorID: 1, Author: "Alexey Zemliakov"}, nil)
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

func TestGetFriendID(t *testing.T) {
	tests := []TableTest[FriendsResponse, FriendsRequest]{
		{
			name: "1",
			SetupInput: func() (*FriendsRequest, error) {
				res := &FriendsRequest{UserID: 0}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *FriendsRequest) (*FriendsResponse, error) {
				return implementation.GetFriendsID(ctx, request)
			},
			ExpectedResult: func() (*FriendsResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.Internal,
			SetupMock: func(request *FriendsRequest, m *mocks) {
				m.profileService.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*FriendsRequest, error) {
				res := &FriendsRequest{UserID: 1}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *FriendsRequest) (*FriendsResponse, error) {
				return implementation.GetFriendsID(ctx, request)
			},
			ExpectedResult: func() (*FriendsResponse, error) {
				return &FriendsResponse{UserID: []uint32{2, 5, 4, 3}}, nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *FriendsRequest, m *mocks) {
				m.profileService.EXPECT().GetFriendsID(gomock.Any(), gomock.Any()).
					Return([]uint32{2, 5, 4, 3}, nil)
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

func TestCreate(t *testing.T) {
	tests := []TableTest[CreateResponse, CreateRequest]{
		{
			name: "1",
			SetupInput: func() (*CreateRequest, error) {
				res := &CreateRequest{User: &User{ID: 0}}
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
				m.profileService.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(uint32(0), errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*CreateRequest, error) {
				res := &CreateRequest{User: &User{ID: 1, FirstName: "Alexey", LastName: "Zemliakov"}}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *CreateRequest) (*CreateResponse, error) {
				return implementation.Create(ctx, request)
			},
			ExpectedResult: func() (*CreateResponse, error) {
				return &CreateResponse{ID: 1}, nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *CreateRequest, m *mocks) {
				m.profileService.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(uint32(1), nil)
			},
		},
		{
			name: "3",
			SetupInput: func() (*CreateRequest, error) {
				res := &CreateRequest{User: &User{ID: 0}}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *CreateRequest) (*CreateResponse, error) {
				return implementation.Create(ctx, request)
			},
			ExpectedResult: func() (*CreateResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.AlreadyExists,
			SetupMock: func(request *CreateRequest, m *mocks) {
				m.profileService.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(uint32(0), my_err.ErrUserAlreadyExists)
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

func TestGetByEmail(t *testing.T) {
	tests := []TableTest[GetByEmailResponse, GetByEmailRequest]{
		{
			name: "1",
			SetupInput: func() (*GetByEmailRequest, error) {
				res := &GetByEmailRequest{Email: ""}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *GetByEmailRequest) (*GetByEmailResponse, error) {
				return implementation.GetUserByEmail(ctx, request)
			},
			ExpectedResult: func() (*GetByEmailResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.Internal,
			SetupMock: func(request *GetByEmailRequest, m *mocks) {
				m.profileService.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*GetByEmailRequest, error) {
				res := &GetByEmailRequest{Email: "alex.zem@gigamail.com"}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *GetByEmailRequest) (*GetByEmailResponse, error) {
				return implementation.GetUserByEmail(ctx, request)
			},
			ExpectedResult: func() (*GetByEmailResponse, error) {
				return &GetByEmailResponse{User: &User{ID: 1, Email: "alex.zem@gigamail.com"}}, nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *GetByEmailRequest, m *mocks) {
				m.profileService.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).
					Return(&models.User{ID: 1, Email: "alex.zem@gigamail.com"}, nil)
			},
		},
		{
			name: "3",
			SetupInput: func() (*GetByEmailRequest, error) {
				res := &GetByEmailRequest{Email: ""}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *GetByEmailRequest) (*GetByEmailResponse, error) {
				return implementation.GetUserByEmail(ctx, request)
			},
			ExpectedResult: func() (*GetByEmailResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.NotFound,
			SetupMock: func(request *GetByEmailRequest, m *mocks) {
				m.profileService.EXPECT().GetByEmail(gomock.Any(), gomock.Any()).
					Return(nil, my_err.ErrUserNotFound)
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
