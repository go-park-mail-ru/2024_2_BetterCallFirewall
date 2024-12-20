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

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

type mocks struct {
	chatService *MockChatService
	responder   *MockResponder
}

func getController(ctrl *gomock.Controller) (*ChatController, *mocks) {
	m := &mocks{
		chatService: NewMockChatService(ctrl),
		responder:   NewMockResponder(ctrl),
	}

	return NewChatController(m.chatService, m.responder), m
}

func TestNewController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	res, _ := getController(ctrl)
	assert.NotNil(t, res)
}

func TestSanitize(t *testing.T) {
	test := "<script> alert(1) </script>"
	expected := ""
	res := sanitize(test)
	assert.Equal(t, expected, res)
}

func TestSanitizeFiles(t *testing.T) {
	test := []string{"<script> alert(1) </script>"}
	var expected []string
	res := sanitizeFiles(test)
	assert.Equal(t, expected, res)
}

func TestGetAllChat(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetAllChats(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat?lastTime=nlbk", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetAllChats(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat?lastTime=2006-01-02T15:04:05Z", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetAllChats(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.chatService.EXPECT().GetAllChats(gomock.Any(), gomock.Any(), gomock.Any()).Return(
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat?lastTime=2006-01-02T15:04:05Z", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetAllChats(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.chatService.EXPECT().GetAllChats(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, nil,
				)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat?lastTime=2006-01-02T15:04:05Z", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetAllChats(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.chatService.EXPECT().GetAllChats(gomock.Any(), gomock.Any(), gomock.Any()).Return(
					nil, errors.New("error"),
				)
				m.responder.EXPECT().ErrorInternal(request.w, gomock.Any(), gomock.Any()).Do(
					func(w, data, req any) {
						request.w.WriteHeader(http.StatusInternalServerError)
						if _, err1 := request.w.Write([]byte("error")); err1 != nil {
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

func TestGetChat(t *testing.T) {
	tests := []TableTest[Response, Request]{
		{
			name: "1",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat/1", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetChat(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat/1?lastTime=nlbk", nil)
				w := httptest.NewRecorder()
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetChat(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat/1?lastTime=nlbk", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetChat(request.w, request.r)
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat/1", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "kljk"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetChat(request.w, request.r)
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
			name: "5",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat/1?", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetChat(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusInternalServerError, Body: "error"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.chatService.EXPECT().GetChat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))
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
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat/1?", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetChat(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusNoContent}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.chatService.EXPECT().GetChat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, my_err.ErrNoMoreContent)
				m.responder.EXPECT().OutputNoMoreContentJSON(request.w, gomock.Any()).Do(
					func(w, req any) {
						request.w.WriteHeader(http.StatusNoContent)
					},
				)
			},
		},
		{
			name: "7",
			SetupInput: func() (*Request, error) {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/message/chat/1?", nil)
				w := httptest.NewRecorder()
				req = req.WithContext(models.ContextWithSession(req.Context(), &models.Session{ID: "1", UserID: 1}))
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				res := &Request{r: req, w: w}
				return res, nil
			},
			Run: func(ctx context.Context, implementation *ChatController, request Request) (Response, error) {
				implementation.GetChat(request.w, request.r)
				res := Response{StatusCode: request.w.Code, Body: request.w.Body.String()}
				return res, nil
			},
			ExpectedResult: func() (Response, error) {
				return Response{StatusCode: http.StatusOK, Body: "OK"}, nil
			},
			ExpectedErr: nil,
			SetupMock: func(request Request, m *mocks) {
				m.responder.EXPECT().LogError(gomock.Any(), gomock.Any())
				m.chatService.EXPECT().GetChat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, nil)
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
	Run            func(context.Context, *ChatController, In) (T, error)
	ExpectedResult func() (T, error)
	ExpectedErr    error
	SetupMock      func(In, *mocks)
}
