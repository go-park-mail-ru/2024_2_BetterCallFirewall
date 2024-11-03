package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

var (
	mockErrorInternal = errors.New("mock internal error")
)

type MockAuthService struct{}

func (m MockAuthService) Register(user models.User, ctx context.Context) (uint32, error) {
	if user.ID == 1 {
		return user.ID, myErr.ErrUserAlreadyExists
	}

	if user.ID == 0 {
		return user.ID, mockErrorInternal
	}

	return user.ID, nil
}

func (m MockAuthService) Auth(user models.User, ctx context.Context) (uint32, error) {
	if user.ID == 1 {
		return user.ID, myErr.ErrWrongEmailOrPassword
	}

	if user.ID == 0 {
		return user.ID, mockErrorInternal
	}

	return user.ID, nil
}

type MockSessionManager struct{}

func (m MockSessionManager) Check(str string) (*models.Session, error) {
	if len(str) > 0 {
		return models.NewSession(10)
	}
	return nil, mockErrorInternal
}

func (m MockSessionManager) Create(userID uint32) (*models.Session, error) {
	if userID == 2 {
		return nil, mockErrorInternal
	}
	return models.NewSession(10)
}

func (m MockSessionManager) Destroy(sess *models.Session) error {
	if sess.UserID == 0 {
		return mockErrorInternal
	}
	return nil
}

type MockResponder struct{}

func (r *MockResponder) OutputJSON(w http.ResponseWriter, data any, _ string) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}

func (r *MockResponder) ErrorUnAuthorized(w http.ResponseWriter, _ error, _ string) {
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte("unauthorized error"))
}

func (r *MockResponder) ErrorBadRequest(w http.ResponseWriter, _ error, _ string) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("bad request error"))
}

func (r *MockResponder) ErrorInternal(w http.ResponseWriter, _ error, _ string) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("internal error"))

}

func (r *MockResponder) ErrorForbidden(w http.ResponseWriter, _ error) {
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write([]byte("forbidden error"))

}

func (r *MockResponder) LogError(err error, requestID string) {}

type TestCase struct {
	w        *httptest.ResponseRecorder
	r        *http.Request
	wantCode int
	wantBody string
}

func TestRegister(t *testing.T) {
	controller := NewAuthController(&MockResponder{}, MockAuthService{}, MockSessionManager{})
	jsonUser0, _ := json.Marshal(models.User{ID: 0})
	jsonUser1, _ := json.Marshal(models.User{ID: 1})
	jsonUser2, _ := json.Marshal(models.User{ID: 2})
	jsonUser3, _ := json.Marshal(models.User{ID: 3})

	testCases := []TestCase{
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("wrong json"))),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonUser1)),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonUser0)),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonUser2)),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonUser3)),
			wantCode: http.StatusOK,
			wantBody: `"user create successful"`,
		},
	}

	for _, tt := range testCases {
		controller.Register(tt.w, tt.r)
		if tt.w.Code != tt.wantCode {
			t.Errorf("Register() code = %d, want %d", tt.w.Code, tt.wantCode)
		}
		if strings.TrimSpace(tt.w.Body.String()) != tt.wantBody {
			t.Errorf("Register() body = %s, want %s", tt.w.Body.String(), tt.wantBody)
		}
	}
}

func TestAuth(t *testing.T) {
	controller := NewAuthController(&MockResponder{}, MockAuthService{}, MockSessionManager{})
	jsonUser0, _ := json.Marshal(models.User{ID: 0})
	jsonUser1, _ := json.Marshal(models.User{ID: 1})
	jsonUser2, _ := json.Marshal(models.User{ID: 2})
	jsonUser3, _ := json.Marshal(models.User{ID: 3})

	testCases := []TestCase{
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte("wrong json"))),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonUser1)),
			wantCode: http.StatusBadRequest,
			wantBody: "bad request error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonUser0)),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonUser2)),
			wantCode: http.StatusInternalServerError,
			wantBody: "internal error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonUser3)),
			wantCode: http.StatusOK,
			wantBody: `"user auth"`,
		},
		{
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonUser3)),
			wantCode: http.StatusOK,
			wantBody: `"user auth"`,
		},
	}

	for _, tt := range testCases {
		controller.Auth(tt.w, tt.r)
		if tt.w.Code != tt.wantCode {
			t.Errorf("Auth() code = %d, want %d", tt.w.Code, tt.w.Code)
		}
		if strings.TrimSpace(tt.w.Body.String()) != tt.wantBody {
			t.Errorf("Auth() body = %s, want %s", tt.w.Body.String(), tt.wantBody)
		}
	}
}

var (
	session, _        = models.NewSession(1)
	ctxWithSession    = models.ContextWithSession(context.Background(), session)
	reqWithSession    = httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctxWithSession)
	reqWithoutSession = httptest.NewRequest(http.MethodPost, "/", nil)
	badSession, _     = models.NewSession(0)
	ctxBadSession     = models.ContextWithSession(context.Background(), badSession)
	reqWithBadSession = httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctxBadSession)
)

func TestLogout(t *testing.T) {
	controller := NewAuthController(&MockResponder{}, MockAuthService{}, MockSessionManager{})

	testCases := []TestCase{
		{
			w:        httptest.NewRecorder(),
			r:        reqWithoutSession,
			wantCode: http.StatusBadRequest,
			wantBody: "bad request error",
		},
		{
			w:        httptest.NewRecorder(),
			r:        reqWithSession,
			wantCode: http.StatusOK,
			wantBody: `"user logout"`,
		},
		{
			w:        httptest.NewRecorder(),
			r:        reqWithBadSession,
			wantCode: http.StatusBadRequest,
			wantBody: "bad request error",
		},
	}

	for caseNum, tt := range testCases {
		controller.Logout(tt.w, tt.r)

		if tt.w.Code != tt.wantCode {
			t.Errorf("[%d] Logout() code = %d, want %d", caseNum, tt.w.Code, tt.wantCode)
		}
		if strings.TrimSpace(tt.w.Body.String()) != tt.wantBody {
			t.Errorf("[%d] Logout() body = %s, want %s", caseNum, tt.w.Body.String(), tt.wantBody)
		}
	}
}
