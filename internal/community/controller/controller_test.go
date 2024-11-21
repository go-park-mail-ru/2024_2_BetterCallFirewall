package controller

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mocks struct {
	communityService *MockcommunityService
	responder        *Mockresponder
}

func getController(ctrl *gomock.Controller) (*Controller, *mocks) {
	m := &mocks{
		communityService: NewMockcommunityService(ctrl),
		responder:        NewMockresponder(ctrl),
	}

	return NewCommunityController(m.responder, m.communityService), m
}

func TestNewController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	res, _ := getController(ctrl)
	assert.NotNil(t, res)
}

func TestGetOne(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/community/n", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.GetOne(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.GetOne(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().GetOne(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "jovbn"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.GetOne(request.w, request.r)
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
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.GetOne(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().GetOne(gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(func(w, data, req any) {
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

func TestGetAll(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/community?id=ljkl", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.GetAll(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/community", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.GetAll(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/community?id=1", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.GetAll(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(func(w, req any) {
					request.w.WriteHeader(http.StatusNoContent)
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/community", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.GetAll(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(
					[]*models.CommunityCard{{ID: 1}},
					nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(func(w, data, req any) {
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

func TestUpdate(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/community/n", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPut, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
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
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/community/1", bytes.NewBuffer([]byte(`{"id":1}`)))
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/community/1", bytes.NewBuffer([]byte(`{"id":1}`)))
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.communityService.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/community/1", bytes.NewBuffer([]byte(`{"id":1}`)))
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.communityService.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(func(w, data, req any) {
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

func TestDelete(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/community/n", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.communityService.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.communityService.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(func(w, data, req any) {
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

func TestCreate(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/n", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Create(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/1", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Create(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/1", bytes.NewBuffer([]byte(`{"id":1}`)))
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Create(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/1", bytes.NewBuffer([]byte(`{"id":1}`)))
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.Create(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(func(w, data, req any) {
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

func TestJoinToCommunity(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/n/join", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.JoinToCommunity(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/n/join", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.JoinToCommunity(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/1/join", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.JoinToCommunity(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().JoinCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/2/join", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "10", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.JoinToCommunity(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().JoinCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
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

func TestLeaveFromCommunity(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/n/leave", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.LeaveFromCommunity(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/n/leave", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.LeaveFromCommunity(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/1/leave", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.LeaveFromCommunity(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().LeaveCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/2/leave", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "10", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.LeaveFromCommunity(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().LeaveCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
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

func TestAddAdmin(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/n/add_admin", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.AddAdmin(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/n/add_admin", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.AddAdmin(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/1/add_admin", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.AddAdmin(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/1/add_admin", bytes.NewBuffer([]byte(`{kj`)))
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.AddAdmin(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)

				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/1/add_admin", bytes.NewBuffer([]byte(`1`)))
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.AddAdmin(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.communityService.EXPECT().AddAdmin(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/community/1/add_admin", bytes.NewBuffer([]byte(`1`)))
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *Controller, request Request) (Response, error) {
				implementation.AddAdmin(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.communityService.EXPECT().CheckAccess(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.communityService.EXPECT().AddAdmin(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
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
	Run            func(context.Context, *Controller, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
