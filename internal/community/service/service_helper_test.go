package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mocksHelper struct {
	repo *MockrepoHelper
}

func getHelper(ctrl *gomock.Controller) (*ServiceHelper, *mocksHelper) {
	m := &mocksHelper{
		repo: NewMockrepoHelper(ctrl),
	}

	return NewServiceHelper(m.repo), m
}

func TestNewServiceHelper(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	res, _ := getHelper(ctrl)
	assert.NotNil(t, res)
}

func TestCheckAccessHelper(t *testing.T) {
	tests := []TableTest2[bool, InputCheckAccess]{
		{
			name: "1",
			SetupInput: func() (*InputCheckAccess, error) {
				input := InputCheckAccess{userID: 0, communityID: 0}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *ServiceHelper, input InputCheckAccess) (bool, error) {
				res := implementation.CheckAccess(ctx, input.userID, input.communityID)
				return res, nil
			},
			ExpectedResult: func() (bool, error) {
				return false, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input InputCheckAccess, m *mocksHelper) {
				m.repo.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
			},
		},
		{
			name: "2",
			SetupInput: func() (*InputCheckAccess, error) {
				input := InputCheckAccess{userID: 1, communityID: 10}
				return &input, nil
			},
			Run: func(ctx context.Context, implementation *ServiceHelper, input InputCheckAccess) (bool, error) {
				res := implementation.CheckAccess(ctx, input.userID, input.communityID)
				return res, nil
			},
			ExpectedResult: func() (bool, error) {
				return true, nil
			},
			ExpectedErr: nil,
			SetupMock: func(input InputCheckAccess, m *mocksHelper) {
				m.repo.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getHelper(ctrl)
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

type TableTest2[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *ServiceHelper, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocksHelper)
}
