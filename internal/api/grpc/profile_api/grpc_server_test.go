package profile_api

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
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
			ExpectedErr: errMock,
			SetupMock: func(request *HeaderRequest, m *mocks) {
				m.profileService.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(models.Header{}, errMock)
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
			ExpectedErr: nil,
			SetupMock: func(request *HeaderRequest, m *mocks) {
				m.profileService.EXPECT().GetHeader(gomock.Any(), gomock.Any()).
					Return(models.Header{AuthorID: 1, Author: "Alexey Zemliakov"}, nil)
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
			ExpectedErr: errMock,
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
			ExpectedErr: nil,
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
