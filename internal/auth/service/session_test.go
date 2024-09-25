package service

import (
	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type MocDB struct {
}

type Test struct {
	testReq  *http.Request
	testResp *httptest.ResponseRecorder
	testRes  *models.Session
	err      error
}

func (m MocDB) CreateSession(session *models.Session) error {
	//TODO implement me
	panic("implement me")
}

func (m MocDB) FindSession(sessID string) (*models.Session, error) {
	if sessID == "2ht4k8s6v0m4hgl1" {
		return &models.Session{
			ID:     "2ht4k8s6v0m4hgl1",
			UserID: 1,
		}, nil
	}
	return nil, nil
}

func (m MocDB) DestroySession(sessID string) error {
	//TODO implement me
	panic("implement me")
}

func TestCheck(t *testing.T) {
	sm := NewSessionManager(MocDB{})
	testCookie := &http.Cookie{
		Name:    "session_id",
		Value:   "2ht4k8s6v0m4hgl1",
		Expires: time.Now().Add(24 * time.Hour),
	}
	wrongTestCookie := &http.Cookie{
		Name:    "sessionId",
		Value:   "2ht4k8s6v0m4hgl1",
		Expires: time.Now().Add(24 * time.Hour),
	}

	testReq := httptest.NewRequest(http.MethodGet, "/", nil)
	testReq.AddCookie(testCookie)
	wrongTestReq := httptest.NewRequest(http.MethodGet, "/", nil)
	wrongTestReq.AddCookie(wrongTestCookie)

	tests := []Test{
		{
			testReq: httptest.NewRequest(http.MethodGet, "/", nil),
			err:     myErr.ErrNoAuth,
		},
		{
			testReq: testReq,
			testRes: &models.Session{
				ID:     "2ht4k8s6v0m4hgl1",
				UserID: 1,
			},
			err: nil,
		},
		{
			testReq: wrongTestReq,
			err:     myErr.ErrNoAuth,
		},
	}

	for caseNum, test := range tests {
		res, err := sm.Check(test.testReq)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(res, test.testRes) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.testRes, res)
		}
	}

}
