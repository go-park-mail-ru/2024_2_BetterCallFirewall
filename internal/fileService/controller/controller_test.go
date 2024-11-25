package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mocks struct {
	fileService *MockfileService
	responder   *Mockresponder
}

func getController(ctrl *gomock.Controller) (*FileController, *mocks) {
	m := &mocks{
		fileService: NewMockfileService(ctrl),
		responder:   NewMockresponder(ctrl),
	}

	return NewFileController(m.fileService, m.responder), m
}

func TestNewController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	res, _ := getController(ctrl)
	assert.NotNil(t, res)
}

var errMock = errors.New("mock error")

func TestUpload(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/image/default", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *FileController, request Request) (Response, error) {
				implementation.Upload(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "2",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/image/default", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"name": "default"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *FileController, request Request) (Response, error) {
				implementation.Upload(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.fileService.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(nil, errMock)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/image/default", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"name": "default"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *FileController, request Request) (Response, error) {
				implementation.Upload(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.fileService.EXPECT().Upload(gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputBytes(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusOK)
					request.w.Write([]byte("OK"))
				})
			},
		},
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
	w *httptest.ResponseRecorder
	r *http.Request
}

type Response struct {
	StatusCode int
	Body       string
}

type TableTest[T, In any] struct {
	name           string
	SetupInput     func() (*In, error)
	Run            func(context.Context, *FileController, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
