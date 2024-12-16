package controller

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

var errMock = errors.New("mock error")

func getController(ctrl *gomock.Controller) (*StickersHandlerImplementation, *mocks) {
	m := &mocks{
		stickerService: NewMockUsecase(ctrl),
		responder:      NewMockResponder(ctrl),
	}

	return NewStickerController(m.stickerService, m.responder), m
}

func TestNewPostController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler, _ := getController(ctrl)
	assert.NotNil(t, handler)
}

func TestAddNewSticker(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/stickers", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.AddNewSticker(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						if _, err1 := request.w.Write([]byte("bad request")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "2",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/stickers", bytes.NewBuffer([]byte(`"files"`)))
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.AddNewSticker(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						if _, err1 := request.w.Write([]byte("bad request")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/stickers", bytes.NewBuffer([]byte(`"/image/someimage"`)),
				)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.AddNewSticker(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						if _, err1 := request.w.Write([]byte("bad request")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/stickers", bytes.NewBuffer([]byte(`{"file":"/image/someimage"}`)),
				)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.AddNewSticker(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.stickerService.EXPECT().AddNewSticker(gomock.Any(), gomock.Any(), gomock.Any()).Return(errMock)
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						_, _ = request.w.Write([]byte("error"))
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/stickers", bytes.NewBuffer([]byte(`{"file":"/image/someimage"}`)),
				)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.AddNewSticker(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.stickerService.EXPECT().AddNewSticker(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusOK)
						_, _ = request.w.Write([]byte("OK"))
					},
				)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
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
			},
		)
	}
}

func TestGetAllSticker(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/stickers/all", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.GetAllStickers(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.stickerService.EXPECT().GetAllStickers(gomock.Any()).Return(nil, my_err.ErrNoStickers)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(
					func(w, req any) {
						request.w.WriteHeader(http.StatusNoContent)
					},
				)
			},
		},
		{
			name: "2",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/stickers/all", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.GetAllStickers(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.stickerService.EXPECT().GetAllStickers(gomock.Any()).Return(nil, errMock)
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						_, _ = request.w.Write([]byte("error"))
					},
				)
			},
		},
		{
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/stickers/all", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.GetAllStickers(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				pic := models.Picture("/image/sticker1")
				m.stickerService.EXPECT().GetAllStickers(gomock.Any()).Return([]*models.Picture{&pic}, nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusOK)
						_, _ = request.w.Write([]byte("OK"))
					},
				)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
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
			},
		)
	}
}

func TestGetMineSticker(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/stickers", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.GetMineStickers(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						_, _ = request.w.Write([]byte("bad request"))
					},
				)
			},
		},
		{
			name: "2",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/stickers", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.GetMineStickers(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.stickerService.EXPECT().GetMineStickers(gomock.Any(), gomock.Any()).Return(nil, errMock)
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						_, _ = request.w.Write([]byte("error"))
					},
				)
			},
		},
		{
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/stickers", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.GetMineStickers(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				pic := models.Picture("/image/sticker1")
				m.stickerService.EXPECT().GetMineStickers(gomock.Any(), gomock.Any()).Return(
					[]*models.Picture{&pic}, nil,
				)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusOK)
						_, _ = request.w.Write([]byte("OK"))
					},
				)
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/stickers", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(
				ctx context.Context, implementation *StickersHandlerImplementation, request Request,
			) (Response, error) {
				implementation.GetMineStickers(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.stickerService.EXPECT().GetMineStickers(gomock.Any(), gomock.Any()).Return(
					nil, my_err.ErrNoStickers,
				)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(
					func(w, req any) {
						request.w.WriteHeader(http.StatusNoContent)
					},
				)
			},
		},
	}

	for _, v := range tests {
		t.Run(
			v.name, func(t *testing.T) {
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
			},
		)
	}
}

type mocks struct {
	stickerService *MockUsecase
	responder      *MockResponder
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
	Run            func(context.Context, *StickersHandlerImplementation, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
