package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type mockPostService struct{}

func (m *mockPostService) GetAll() []*models.Post {
	return []*models.Post{{Header: "Header", Body: "Body", CreatedAt: "2012-10-12"}}
}

type mockResponder struct{}

func (m *mockResponder) OutputJSON(w http.ResponseWriter, _ any) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("data send to client"))
}

func (m *mockResponder) ErrorWrongMethod(w http.ResponseWriter, _ error) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte("send err to client"))
}

type TestCase struct {
	r        *http.Request
	w        *httptest.ResponseRecorder
	wantBody string
	wantCode int
}

func TestPostController(t *testing.T) {
	controller := NewPostController(&mockPostService{}, &mockResponder{})
	tests := []TestCase{
		{
			r:        httptest.NewRequest("POST", "/api/v1/post", nil),
			w:        httptest.NewRecorder(),
			wantBody: "send err to client",
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			r:        httptest.NewRequest("GET", "/api/v1/post", nil),
			w:        httptest.NewRecorder(),
			wantBody: "data send to client",
			wantCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		controller.GetAll(test.w, test.r)
		if test.wantCode != test.w.Code {
			t.Errorf("Wanted response code: %d, got: %d", test.w.Code, test.w.Code)
		}
		if test.w.Body.String() != test.wantBody {
			t.Errorf("Wanted response body: %s, got: %s", test.wantBody, test.w.Body.String())
		}
	}
}
