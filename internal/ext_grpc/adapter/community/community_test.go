package community

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/api/grpc/community_api"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mocks struct {
	client *MockCommunityServiceClient
}

func getAdapter(ctrl *gomock.Controller) (*GrpcSender, *mocks) {
	m := mocks{
		client: NewMockCommunityServiceClient(ctrl),
	}

	return &GrpcSender{client: m.client}, &m
}

type input struct {
	userID      uint32
	communityID uint32
}

func TestCheckAccess(t *testing.T) {
	tests := []TableTest[bool, *input]{
		{
			name: "1",
			SetupInput: func() (*input, error) {
				return &input{}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *input) (bool, error) {
				res := implementation.CheckAccess(ctx, request.communityID, request.userID)
				return res, nil
			},
			ExpectedResult: func() (bool, error) {
				return false, nil
			},
			SetupMock: func(request *input, m *mocks) {
				m.client.EXPECT().CheckAccess(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("mock error"))
			},
		},
		{
			name: "2",
			SetupInput: func() (*input, error) {
				return &input{userID: 1, communityID: 10}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *input) (bool, error) {
				res := implementation.CheckAccess(ctx, request.communityID, request.userID)
				return res, nil
			},
			ExpectedResult: func() (bool, error) {
				return true, nil
			},
			SetupMock: func(request *input, m *mocks) {
				m.client.EXPECT().CheckAccess(gomock.Any(), gomock.Any()).
					Return(&community_api.CheckAccessResponse{Access: true}, nil)
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
			assert.NoError(t, err)
			assert.Equal(t, res, actual)
		})
	}
}

func TestGetHeader(t *testing.T) {
	errMock := errors.New("mock error")
	tests := []TableTest[*models.Header, *input]{
		{
			name: "1",
			SetupInput: func() (*input, error) {
				return &input{}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *input) (*models.Header, error) {
				res, err := implementation.GetHeader(ctx, request.communityID)
				return res, err
			},
			ExpectedErr: errMock,
			ExpectedResult: func() (*models.Header, error) {
				return nil, nil
			},
			SetupMock: func(request *input, m *mocks) {
				m.client.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(nil, errMock)
			},
		},
		{
			name: "2",
			SetupInput: func() (*input, error) {
				return &input{}, nil
			},
			Run: func(ctx context.Context, implementation *GrpcSender, request *input) (*models.Header, error) {
				res, err := implementation.GetHeader(ctx, request.communityID)
				return res, err
			},
			ExpectedErr: nil,
			ExpectedResult: func() (*models.Header, error) {
				return &models.Header{
					AuthorID:    0,
					CommunityID: 1,
					Author:      "Community",
					Avatar:      "/avatar",
				}, nil
			},
			SetupMock: func(request *input, m *mocks) {
				m.client.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(&community_api.GetHeaderResponse{
						Head: &community_api.Header{
							CommunityID: 1,
							Author:      "Community",
							Avatar:      "/avatar",
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

func TestPostProvider(t *testing.T) {
	_, err := GetCommunityProvider("", "")
	assert.Nil(t, err)
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (In, error)
	Run            func(context.Context, *GrpcSender, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
