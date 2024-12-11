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

func getController(ctrl *gomock.Controller) (*PostController, *mocks) {
	m := &mocks{
		postService:    NewMockPostService(ctrl),
		responder:      NewMockResponder(ctrl),
		commentService: NewMockCommentService(ctrl),
	}

	return NewPostController(m.postService, m.commentService, m.responder), m
}

func TestNewPostController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	handler, _ := getController(ctrl)
	assert.NotNil(t, handler)
}

func TestCreate(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/feed",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
					http.MethodPost, "/api/v1/feed",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				m.postService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint32(0), errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
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
					http.MethodPost, "/api/v1/feed",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				m.postService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(uint32(2), nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/feed?community=ljkhkg",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/feed?community=10",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				m.postService.EXPECT().CheckAccessToCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
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
			name: "7",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/feed?community=10",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				m.postService.EXPECT().CheckAccessToCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.postService.EXPECT().CreateCommunityPost(gomock.Any(), gomock.Any()).Return(
					uint32(0), errors.New("error"),
				)
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "8",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/feed?community=10",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				m.postService.EXPECT().CheckAccessToCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.postService.EXPECT().CreateCommunityPost(gomock.Any(), gomock.Any()).Return(uint32(10), nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "9",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/feed?community=10",
					bytes.NewBuffer([]byte(`{"post_content":{"text":"new post Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed tellus arcu, vulputate rutrum enim vitae, tincidunt imperdiet tellus. Aenean vulputate elit consequat lorem pellentesque bibendum. Donec sed mi posuere dolor semper mollis eu eget dolor. Proin et eleifend magna. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Curabitur tempus ultricies mi, eget malesuada metus. Nam sit amet felis nec dolor vehicula dapibus gravida in nunc. Mauris turpis et. "}}`)),
				)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						if _, err1 := request.w.Write([]byte("bad request")); err1 != nil {
							panic(err1)
						}
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

func TestGetOne(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "jhg"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				m.postService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, my_err.ErrPostNotFound)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/10", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				m.postService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/10", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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
				m.postService.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
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

func TestUpdate(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
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
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/10", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
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
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/10", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetPostAuthorID(gomock.Any(), gomock.Any()).Return(
					uint32(0), errors.New("error"),
				)
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
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/10", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetPostAuthorID(gomock.Any(), gomock.Any()).Return(uint32(10), nil)
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
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPut, "/api/v1/feed/10",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetPostAuthorID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(my_err.ErrPostNotFound)
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
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPut, "/api/v1/feed/10",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetPostAuthorID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "7",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPut, "/api/v1/feed/10",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetPostAuthorID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "8",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPut, "/api/v1/feed/10?community=nkljbkvhj",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
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
			name: "9",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPut, "/api/v1/feed/10?community=10",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().CheckAccessToCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(false)
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
			name: "10",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPut, "/api/v1/feed/10?community=10",
					bytes.NewBuffer([]byte(`{"id":1}`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().CheckAccessToCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.postService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "11",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPut, "/api/v1/feed/10?community=10",
					bytes.NewBuffer([]byte(`{"id"`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().CheckAccessToCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
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
			name: "12",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPut, "/api/v1/feed/10?community=1",
					bytes.NewBuffer([]byte(`{"post_content":{"text":"new post Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed tellus arcu, vulputate rutrum enim vitae, tincidunt imperdiet tellus. Aenean vulputate elit consequat lorem pellentesque bibendum. Donec sed mi posuere dolor semper mollis eu eget dolor. Proin et eleifend magna. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Curabitur tempus ultricies mi, eget malesuada metus. Nam sit amet felis nec dolor vehicula dapibus gravida in nunc. Mauris turpis et. "}}`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "10"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Update(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().CheckAccessToCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
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

func TestDelete(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
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
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
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
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetPostAuthorID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postService.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(my_err.ErrPostNotFound)
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
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetPostAuthorID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postService.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetPostAuthorID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
				m.postService.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1?community=lojhiuk", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, error, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						if _, err1 := request.w.Write([]byte("bad request")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "7",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1?community=10", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, error, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						if _, err1 := request.w.Write([]byte("bad request")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "8",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1?community=10", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Delete(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().CheckAccessToCommunity(gomock.Any(), gomock.Any(), gomock.Any()).Return(true)
				m.postService.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, error, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
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

func TestGetBatchPost(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed?id=jvhjh", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=nbn", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, my_err.ErrNoMoreContent,
				)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(
					func(w, req any) {
						request.w.WriteHeader(http.StatusNoContent)
					},
				)
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, errors.New("error"),
				)
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetBatch(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=friend", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
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
			name: "7",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=friend", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetBatchFromFriend(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "8",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed?community=jklh", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
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
			name: "9",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed?community=10&id=10", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
				m.postService.EXPECT().GetCommunityPost(
					gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(nil, nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "11",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=vbfvib", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetBatchPosts(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any()).Do(func(err, req any) {})
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

func TestSetLikeOnPost(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/like", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.SetLikeOnPost(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/like", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.SetLikeOnPost(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/like", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.SetLikeOnPost(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.postService.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/like", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.SetLikeOnPost(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.postService.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					false, errors.New("error"),
				)
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/like", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.SetLikeOnPost(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.postService.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
				m.postService.EXPECT().SetLikeToPost(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/like", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.SetLikeOnPost(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.postService.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
				m.postService.EXPECT().SetLikeToPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
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

func TestDeleteLikeFromPost(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/unlike", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteLikeFromPost(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/unlike", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteLikeFromPost(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/unlike", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteLikeFromPost(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.postService.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/unlike", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteLikeFromPost(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.postService.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					false, errors.New("error"),
				)
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/unlike", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteLikeFromPost(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.postService.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				m.postService.EXPECT().DeleteLikeFromPost(
					gomock.Any(), gomock.Any(), gomock.Any(),
				).Return(errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
							panic(err1)
						}
					},
				)
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2/unlike", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteLikeFromPost(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.postService.EXPECT().CheckLikes(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				m.postService.EXPECT().DeleteLikeFromPost(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusOK)
						if _, err1 := request.w.Write([]byte("OK")); err1 != nil {
							panic(err1)
						}
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

func TestComment(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Comment(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Comment(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Comment(request.w, request.r)
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
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2", bytes.NewBuffer([]byte(`{"id":1}`)))
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Comment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().Comment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
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
				req := httptest.NewRequest(http.MethodPost, "/api/v1/feed/2", bytes.NewBuffer([]byte(`{"id":1}`)))
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Comment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().Comment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(
						&models.Comment{
							Content: models.Content{
								Text: "New comment",
							},
						}, nil,
					)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusOK)
						_, _ = request.w.Write([]byte("OK"))
					},
				)
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPost, "/api/v1/feed/2",
					bytes.NewBuffer([]byte(`{"file":"очень большой текст, написанный в поле файл, как он сюда попал - честно говоря хз, надо проверить валидацию на файл длину"}`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.Comment(request.w, request.r)
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

func TestGetComment(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed/2/comment", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetComments(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed/2/comment?id=fnf", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetComments(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed/2/comment?id=4", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetComments(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().GetComments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						_, _ = request.w.Write([]byte("error"))
					},
				)
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed/2/comment?sort=old", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetComments(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().GetComments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusOK)
						_, _ = request.w.Write([]byte("OK"))
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/feed/2/comment", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{"id": "2"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.GetComments(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().GetComments(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, my_err.ErrNoMoreContent)
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

func TestEditComment(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/2/", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.EditComment(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/2/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.EditComment(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/2/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.EditComment(request.w, request.r)
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
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/2/1", bytes.NewBuffer([]byte(`{"text":"1"}`)))
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.EditComment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().EditComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(my_err.ErrAccessDenied)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						_, _ = request.w.Write([]byte("bad request"))
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/2/1", bytes.NewBuffer([]byte(`{"text":"1"}`)))
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.EditComment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().EditComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(my_err.ErrWrongComment)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						_, _ = request.w.Write([]byte("bad request"))
					},
				)
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/2/1", bytes.NewBuffer([]byte(`{"text":"1"}`)))
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.EditComment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().EditComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						_, _ = request.w.Write([]byte("error"))
					},
				)
			},
		},
		{
			name: "7",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodPut, "/api/v1/feed/2/1", bytes.NewBuffer([]byte(`{"text":"1"}`)))
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.EditComment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().EditComment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusOK)
						_, _ = request.w.Write([]byte("OK"))
					},
				)
			},
		},
		{
			name: "8",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(
					http.MethodPut, "/api/v1/feed/2/1",
					bytes.NewBuffer([]byte(`{"file":"очень большой текст, написанный в поле файл, как он сюда попал - честно говоря хз, надо проверить валидацию на файл длину"}`)),
				)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.EditComment(request.w, request.r)
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

func TestDeleteComment(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/2/", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteComment(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/2/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteComment(request.w, request.r)
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
			name: "3",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/2/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteComment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().DeleteComment(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(my_err.ErrAccessDenied)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						_, _ = request.w.Write([]byte("bad request"))
					},
				)
			},
		},
		{
			name: "4",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/2/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteComment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusBadRequest, Body: "bad request"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().DeleteComment(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(my_err.ErrWrongComment)
				m.responder.EXPECT().ErrorBadRequest(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusBadRequest)
						_, _ = request.w.Write([]byte("bad request"))
					},
				)
			},
		},
		{
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/2/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteComment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().DeleteComment(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, err, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						_, _ = request.w.Write([]byte("error"))
					},
				)
			},
		},
		{
			name: "6",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodDelete, "/api/v1/feed/2/1", nil)
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, map[string]string{commentIDKey: "1"})
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *PostController, request Request) (Response, error) {
				implementation.DeleteComment(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.commentService.EXPECT().DeleteComment(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				m.responder.EXPECT().OutputJSON(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
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

type mocks struct {
	postService    *MockPostService
	responder      *MockResponder
	commentService *MockCommentService
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
	Run            func(context.Context, *PostController, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
