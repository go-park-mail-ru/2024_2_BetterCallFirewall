package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

var errMock = errors.New("mock error")

type mockPostService struct{}

func (m *mockPostService) Create(ctx context.Context, post *models.Post) (uint32, error) {
	if post.PostContent.Text == "wrong post" {
		return 0, errMock
	}

	return 1, nil
}

func (m *mockPostService) Get(ctx context.Context, postID uint32) (*models.Post, error) {
	if postID == 100 {
		return nil, myErr.ErrPostNotFound
	}

	if postID == 200 {
		return nil, errMock
	}

	return &models.Post{ID: postID}, nil
}

func (m *mockPostService) Update(ctx context.Context, post *models.Post) error {
	if post.ID == 2 {
		return myErr.ErrPostNotFound
	}
	if post.PostContent.Text == "bad text in post" {
		return errMock
	}

	return nil
}

func (m *mockPostService) Delete(ctx context.Context, postID uint32) error {
	if postID == 300 {
		return myErr.ErrPostNotFound
	}
	if postID == 400 {
		return errMock
	}

	return nil
}

func (m *mockPostService) GetBatch(ctx context.Context, lastID uint32) ([]*models.Post, error) {
	data := ctx.Value("query")
	if data == "internalError" {
		return nil, errMock
	}
	if data == "0 post" {
		return nil, myErr.ErrNoMoreContent
	}
	if data == "another err" {
		return nil, myErr.ErrAnotherService
	}

	return []*models.Post{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}}, nil
}

func (m *mockPostService) GetBatchFromFriend(ctx context.Context, userID uint32, lastID uint32) ([]*models.Post, error) {
	data := ctx.Value("query")
	if data == "internalError" {
		return nil, errMock
	}
	if data == "1 post" {
		return []*models.Post{{ID: 1}}, myErr.ErrNoMoreContent
	}

	return []*models.Post{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}, {ID: 5}, {ID: 6}, {ID: 7}, {ID: 8}}, nil
}

func (m *mockPostService) GetPostAuthorID(ctx context.Context, postID uint32) (uint32, error) {
	if postID == 100 {
		return 0, myErr.ErrPostNotFound
	}
	if postID == 200 {
		return 0, errMock
	}

	return 1, nil
}

type mockResponder struct {
}

func (m *mockResponder) OutputJSON(w http.ResponseWriter, data any, _ string) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
}

func (m *mockResponder) OutputNoMoreContentJSON(w http.ResponseWriter, _ string) {
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("no more content"))
}

func (m *mockResponder) ErrorInternal(w http.ResponseWriter, err error, _ string) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("internal server error"))
}

func (m *mockResponder) ErrorBadRequest(w http.ResponseWriter, err error, _ string) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("bad request"))
}

func (m *mockResponder) LogError(err error, _ string) {}

type mockFileService struct{}

func (m *mockFileService) Download(ctx context.Context, file multipart.File, postID, profileID uint32) error {
	return nil
}

func (m *mockFileService) GetPostPicture(ctx context.Context, postID uint32) *models.Picture {
	return nil
}

func (m *mockFileService) UpdatePostFile(ctx context.Context, file multipart.File, postID uint32) error {
	return nil
}

type TestCase struct {
	r        *http.Request
	w        *httptest.ResponseRecorder
	wantBody string
	wantCode int
}

var (
	controller = NewPostController(&mockPostService{}, &mockResponder{}, &mockFileService{})
)

func TestCreate(t *testing.T) {
	var (
		post1, _        = json.Marshal(&models.Post{PostContent: models.Content{Text: "post 1"}})
		badPost, _      = json.Marshal(&models.Post{PostContent: models.Content{Text: "wrong post"}})
		sessGoodUser, _ = models.NewSession(1)
		ctxSess         = models.ContextWithSession(context.Background(), sessGoodUser)
	)

	tests := []TestCase{
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/api/v1/feed", bytes.NewBuffer([]byte("wrong json"))),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/api/v1/feed", bytes.NewBuffer(post1)),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/api/v1/feed", bytes.NewBuffer(post1)).WithContext(ctxSess),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/api/v1/feed", bytes.NewBuffer(badPost)).WithContext(ctxSess),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
	}

	for _, tt := range tests {
		controller.Create(tt.w, tt.r)
		if tt.w.Code != tt.wantCode {
			t.Errorf("Create() code = %d, want %d", tt.w.Code, tt.wantCode)
		}
		if strings.TrimSpace(tt.w.Body.String()) != tt.wantBody {
			t.Errorf("Create() body = %s, want %s", tt.w.Body.String(), tt.wantBody)
		}
	}
}

func TestGetOne(t *testing.T) {
	var (
		badID         = map[string]string{"id": "-1"}
		badIDNotFound = map[string]string{"id": "100"}
		badIDInternal = map[string]string{"id": "200"}
		goodID        = map[string]string{"id": "1"}
	)

	tests := []TestCase{
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "/api/v1/feed/", nil),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/api/v1/feed/-1", nil), badID),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/api/v1/feed/100", nil), badIDNotFound),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/api/v1/feed/200", nil), badIDInternal),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal server error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/api/v1/feed/1", nil), goodID),
			wantCode: http.StatusOK,
			wantBody: "Ok",
		},
	}

	for _, tt := range tests {
		controller.GetOne(tt.w, tt.r)
		if tt.w.Code != tt.wantCode {
			t.Errorf("GetOne() code = %d, want %d", tt.w.Code, http.StatusOK)
		}
		if strings.TrimSpace(tt.w.Body.String()) != tt.wantBody {
			t.Errorf("GetOne() body = %s, want %s", tt.w.Body.String(), "Ok")
		}
	}
}

func TestUpdate(t *testing.T) {
	var (
		post1, _             = json.Marshal(&models.Post{ID: 1, PostContent: models.Content{Text: "post 1"}})
		post2, _             = json.Marshal(&models.Post{ID: 2, PostContent: models.Content{Text: "post 2"}})
		post3, _             = json.Marshal(&models.Post{ID: 1, PostContent: models.Content{Text: "bad text in post"}})
		badPost, _           = json.Marshal(&models.Post{PostContent: models.Content{Text: "wrong post"}})
		notFoundPost, _      = json.Marshal(&models.Post{ID: 100, PostContent: models.Content{Text: "not found"}})
		internalErrorPost, _ = json.Marshal(&models.Post{ID: 200, PostContent: models.Content{Text: "internal error"}})
		sessGoodUser, _      = models.NewSession(1)
		ctxSess              = models.ContextWithSession(context.Background(), sessGoodUser)
		sessBadUser, _       = models.NewSession(2)
		ctxSessBad           = models.ContextWithSession(context.Background(), sessBadUser)
	)

	tests := []TestCase{
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPut, "/api/v1/feed/1", bytes.NewBuffer(badPost)),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPut, "/api/v1/feed/1", bytes.NewBuffer(post1)),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPut, "/api/v1/feed/100", bytes.NewBuffer(notFoundPost)).WithContext(ctxSessBad),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPut, "/api/v1/feed/200", bytes.NewBuffer(internalErrorPost)).WithContext(ctxSessBad),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal server error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPut, "/api/v1/feed/1", bytes.NewBuffer(post1)).WithContext(ctxSessBad),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPut, "/api/v1/feed/2", bytes.NewBuffer(post2)).WithContext(ctxSess),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPut, "/api/v1/feed/1", bytes.NewBuffer(post3)).WithContext(ctxSess),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal server error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPut, "/api/v1/feed/1", bytes.NewBuffer(post1)).WithContext(ctxSess),
			wantCode: http.StatusOK,
			wantBody: "Ok",
		},
	}

	for _, tt := range tests {
		controller.Update(tt.w, tt.r)
		if tt.w.Code != tt.wantCode {
			t.Errorf("Update() code = %d, want %d", tt.w.Code, tt.wantCode)
		}
		if strings.TrimSpace(tt.w.Body.String()) != tt.wantBody {
			t.Errorf("Update() body = %s, want %s", tt.w.Body.String(), "Ok")
		}
	}
}

func TestDelete(t *testing.T) {
	var (
		badID           = map[string]string{"id": "-1"}
		badIDNotFound   = map[string]string{"id": "100"}
		badIDInternal   = map[string]string{"id": "200"}
		badIDNotFound2  = map[string]string{"id": "300"}
		badIDInternal2  = map[string]string{"id": "400"}
		goodID          = map[string]string{"id": "1"}
		sessGoodUser, _ = models.NewSession(1)
		ctxSess         = models.ContextWithSession(context.Background(), sessGoodUser)
		sessBadUser, _  = models.NewSession(2)
		ctxSessBad      = models.ContextWithSession(context.Background(), sessBadUser)
	)

	tests := []TestCase{
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodDelete, "/api/v1/feed/", nil),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/api/v1/feed/-1", nil), badID),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1", nil), goodID),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/api/v1/feed/100", nil).WithContext(ctxSess), badIDNotFound),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/api/v1/feed/200", nil).WithContext(ctxSess), badIDInternal),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal server error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1", nil).WithContext(ctxSessBad), goodID),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/api/v1/feed/300", nil).WithContext(ctxSess), badIDNotFound2),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/api/v1/feed/400", nil).WithContext(ctxSess), badIDInternal2),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal server error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/api/v1/feed/1", nil).WithContext(ctxSess), goodID),
			wantCode: http.StatusOK,
			wantBody: "Ok",
		},
	}
	for _, tt := range tests {
		controller.Delete(tt.w, tt.r)
		if tt.w.Code != tt.wantCode {
			t.Errorf("Delete() code = %d, want %d", tt.w.Code, tt.wantCode)
		}
		if strings.TrimSpace(tt.w.Body.String()) != tt.wantBody {
			t.Errorf("Delete() body = %s, want %s", tt.w.Body.String(), tt.wantBody)
		}
	}
}

func TestGetBatch(t *testing.T) {
	var (
		sessGoodUser, _ = models.NewSession(1)
		ctxSess         = models.ContextWithSession(context.Background(), sessGoodUser)
	)

	tests := []TestCase{
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=false", nil),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil).
				WithContext(context.WithValue(context.Background(), "query", "internalError")),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal server error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=friend", nil),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=friend", nil),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil).
				WithContext(context.WithValue(context.Background(), "query", "0 post")),
			wantCode: http.StatusNoContent,
			wantBody: "no more content",
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil).
				WithContext(context.WithValue(context.Background(), "query", "another err")),
			wantCode: http.StatusOK,
			wantBody: "Ok",
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodGet, "/api/v1/feed", nil).
				WithContext(context.WithValue(context.Background(), "query", "many posts")),
			wantCode: http.StatusOK,
			wantBody: "Ok",
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=friend", nil).
				WithContext(context.WithValue(context.Background(), "query", "1 post")).
				WithContext(ctxSess),
			wantCode: http.StatusOK,
			wantBody: "Ok",
		},
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=friend", nil).
				WithContext(ctxSess),
			wantCode: http.StatusOK,
			wantBody: "Ok",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=friend&id=-1", nil).WithContext(ctxSess),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "/api/v1/feed?section=friend&id=aaaa", nil).WithContext(ctxSess),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request",
		},
	}

	for _, tt := range tests {
		controller.GetBatchPosts(tt.w, tt.r)
		if tt.w.Code != tt.wantCode {
			t.Errorf("GetBatch() code = %d, want %d", tt.w.Code, tt.wantCode)
		}
		if strings.TrimSpace(tt.w.Body.String()) != tt.wantBody {
			t.Errorf("GetBatch() body = %s, want %s", tt.w.Body.String(), "Ok")
		}
	}
}
