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
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
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
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/header", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetHeader(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/header", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetHeader(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(nil, my_err.ErrProfileNotFound)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/header", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "0", UserID: 0}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetHeader(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "internal error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("internal error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/header", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "10", UserID: 10}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetHeader(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetHeader(gomock.Any(), gomock.Any()).Return(&models.Header{}, nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(func(w, header, req any) {
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

func TestGetProfile(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfile(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfile(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetProfileById(gomock.Any(), gomock.Any()).Return(&models.FullProfile{}, my_err.ErrProfileNotFound)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfile(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetProfileById(gomock.Any(), gomock.Any()).Return(&models.FullProfile{}, my_err.ErrNoMoreContent)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(func(w, req any) {
					request.w.WriteHeader(http.StatusNoContent)
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfile(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetProfileById(gomock.Any(), gomock.Any()).Return(&models.FullProfile{}, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "10", UserID: 100}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfile(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetProfileById(gomock.Any(), gomock.Any()).Return(&models.FullProfile{}, nil)
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

func TestUpdateProfile(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.UpdateProfile(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.UpdateProfile(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPut, "/api/v1/profile",
					bytes.NewBuffer([]byte(`{"id":0, "first_name":"Alexey"}`)))
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.UpdateProfile(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/profile",
					bytes.NewBuffer([]byte(`{"id":1, "first_name":"Alexey"}`)))
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.UpdateProfile(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(nil)
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

func TestDeleteProfile(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.DeleteProfile(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.DeleteProfile(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().DeleteProfile(gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.DeleteProfile(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().DeleteProfile(gomock.Any()).Return(nil)
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

func TestGetProfileByID(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/a", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfileById(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/a", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "a"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfileById(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/a", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "5000000000"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfileById(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/a", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfileById(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetProfileById(gomock.Any(), gomock.Any()).Return(&models.FullProfile{}, my_err.ErrProfileNotFound)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/a", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfileById(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetProfileById(gomock.Any(), gomock.Any()).Return(&models.FullProfile{}, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile/a", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetProfileById(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetProfileById(gomock.Any(), gomock.Any()).Return(&models.FullProfile{}, nil)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=bivhub", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
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
				m.profileManager.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
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
				m.profileManager.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(func(w, req any) {
					request.w.WriteHeader(http.StatusNoContent)
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=10", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
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
				m.profileManager.EXPECT().GetAll(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]*models.ShortProfile{{ID: 1}},
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

func TestSendFriendReq(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.SendFriendReq(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.SendFriendReq(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.SendFriendReq(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().SendFriendReq(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusBadRequest)
					request.w.Write([]byte("bad request"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.SendFriendReq(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().SendFriendReq(gomock.Any(), gomock.Any()).Return(nil)
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

func TestAcceptFriendReq(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.AcceptFriendReq(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.AcceptFriendReq(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.AcceptFriendReq(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().AcceptFriendReq(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.AcceptFriendReq(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().AcceptFriendReq(gomock.Any(), gomock.Any()).Return(nil)
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

func TestRemoveFromFriends(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.RemoveFromFriends(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.RemoveFromFriends(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.RemoveFromFriends(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().RemoveFromFriends(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.RemoveFromFriends(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().RemoveFromFriends(gomock.Any(), gomock.Any()).Return(nil)
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

func TestUnsubscribe(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.Unsubscribe(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.Unsubscribe(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.Unsubscribe(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().Unsubscribe(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.Unsubscribe(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().Unsubscribe(gomock.Any(), gomock.Any()).Return(nil)
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

func TestGetAllFriends(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllFriends(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=bjk", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllFriends(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllFriends(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetAllFriends(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllFriends(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetAllFriends(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(func(w, req any) {
					request.w.WriteHeader(http.StatusNoContent)
				})
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllFriends(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetAllFriends(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]*models.ShortProfile{{ID: 1}},
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

func TestGetAllSubs(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubs(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=bjk", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubs(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubs(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetAllSubs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubs(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetAllSubs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(func(w, req any) {
					request.w.WriteHeader(http.StatusNoContent)
				})
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubs(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetAllSubs(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]*models.ShortProfile{{ID: 1}},
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

func TestGetAllSubscriptions(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubscriptions(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=bjk", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubscriptions(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubscriptions(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetAllSubscriptions(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubscriptions(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetAllSubscriptions(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(func(w, req any) {
					request.w.WriteHeader(http.StatusNoContent)
				})
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetAllSubscriptions(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetAllSubscriptions(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]*models.ShortProfile{{ID: 1}},
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

func TestGetCommunitySubs(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetCommunitySubs(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=bjk", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetCommunitySubs(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetCommunitySubs(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetCommunitySubs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(func(w, err, req any) {
					request.w.WriteHeader(http.StatusInternalServerError)
					request.w.Write([]byte("error"))
				})
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetCommunitySubs(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetCommunitySubs(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(func(w, req any) {
					request.w.WriteHeader(http.StatusNoContent)
				})
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/profile?last_id=1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ProfileHandlerImplementation, request Request) (Response, error) {
				implementation.GetCommunitySubs(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.profileManager.EXPECT().GetCommunitySubs(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					[]*models.ShortProfile{{ID: 1}},
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
	Run            func(context.Context, *ProfileHandlerImplementation, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
