package community_api

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
	communityService *MockCommunityService
}

func getAdapter(ctrl *gomock.Controller) (*Adapter, *mocks) {
	m := &mocks{
		communityService: NewMockCommunityService(ctrl),
	}

	return New(m.communityService), m
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	a, _ := getAdapter(ctrl)
	assert.NotNil(t, a)
}

func TestCheckAccess(t *testing.T) {
	tests := []TableTest[CheckAccessResponse, CheckAccessRequest]{
		{
			name: "1",
			SetupInput: func() (*CheckAccessRequest, error) {
				res := &CheckAccessRequest{UserID: 1, CommunityID: 1}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *CheckAccessRequest) (*CheckAccessResponse, error) {
				return implementation.CheckAccess(ctx, request)
			},
			ExpectedResult: func() (*CheckAccessResponse, error) {
				return &CheckAccessResponse{Access: false}, nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *CheckAccessRequest, m *mocks) {
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(false)
			},
		},
		{
			name: "2",
			SetupInput: func() (*CheckAccessRequest, error) {
				res := &CheckAccessRequest{UserID: 1, CommunityID: 10}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *CheckAccessRequest) (*CheckAccessResponse, error) {
				return implementation.CheckAccess(ctx, request)
			},
			ExpectedResult: func() (*CheckAccessResponse, error) {
				return &CheckAccessResponse{Access: true}, nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *CheckAccessRequest, m *mocks) {
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(true)
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

func TestGetHeader(t *testing.T) {
	tests := []TableTest[GetHeaderResponse, GetHeaderRequest]{
		{
			name: "1",
			SetupInput: func() (*GetHeaderRequest, error) {
				res := &GetHeaderRequest{CommunityID: 1}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *GetHeaderRequest) (*GetHeaderResponse, error) {
				return implementation.GetHeader(ctx, request)
			},
			ExpectedResult: func() (*GetHeaderResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.NotFound,
			SetupMock: func(request *GetHeaderRequest, m *mocks) {
				m.communityService.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(nil, my_err.ErrWrongCommunity)
			},
		},
		{
			name: "2",
			SetupInput: func() (*GetHeaderRequest, error) {
				res := &GetHeaderRequest{CommunityID: 1}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *GetHeaderRequest) (*GetHeaderResponse, error) {
				return implementation.GetHeader(ctx, request)
			},
			ExpectedResult: func() (*GetHeaderResponse, error) {
				return nil, nil
			},
			ExpectedErrCode: codes.Internal,
			SetupMock: func(request *GetHeaderRequest, m *mocks) {
				m.communityService.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
			},
		},
		{
			name: "3",
			SetupInput: func() (*GetHeaderRequest, error) {
				res := &GetHeaderRequest{CommunityID: 1}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Adapter, request *GetHeaderRequest) (*GetHeaderResponse, error) {
				return implementation.GetHeader(ctx, request)
			},
			ExpectedResult: func() (*GetHeaderResponse, error) {
				return &GetHeaderResponse{Head: &Header{
						AuthorID:    0,
						CommunityID: 1,
						Author:      "community",
						Avatar:      "some avatar",
					}},
					nil
			},
			ExpectedErrCode: codes.OK,
			SetupMock: func(request *GetHeaderRequest, m *mocks) {
				m.communityService.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(&models.Header{
						AuthorID:    0,
						CommunityID: 1,
						Author:      "community",
						Avatar:      "some avatar",
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
