package community_api

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
			ExpectedErr: nil,
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
			ExpectedErr: nil,
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
			if !errors.Is(err, v.ExpectedErr) {
				t.Errorf("expect %v, got %v", v.ExpectedErr, err)
			}
		})
	}
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *Adapter, *In) (*T, error)
	ExpectedResult func() (*T, error)
	ExpectedErr    error
	SetupMock      func(*In, *mocks)
}
