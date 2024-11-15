package controller

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mocks struct {
	profileManager *MockProfileUsecase
	responder      *MockResponder
}

func getController(ctrl *gomock.Controller) (*ProfileHandlerImplementation, *mocks) {
	m := &mocks{
		profileManager: NewMockProfileUsecase(ctrl),
		responder:      NewMockResponder(ctrl),
	}

	return NewProfileController(m.profileManager, m.responder), m
}

func TestNewProfileHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler, _ := getController(ctrl)
	assert.NotNil(t, handler)
}

func TestGetHeader(t *testing.T) {

	tests := []TableTest[Request, Request]{
		{},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serv, mock := getController(ctrl)
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

type Request struct {
	w http.ResponseWriter
	r *http.Request
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *ProfileHandlerImplementation, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
