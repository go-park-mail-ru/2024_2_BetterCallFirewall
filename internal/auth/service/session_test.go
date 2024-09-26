package service

import (
	"errors"
	"fmt"
	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type MocDB struct {
	Storage map[string]*models.Session
}

type Test struct {
	testReq  *http.Request
	testResp *httptest.ResponseRecorder
	testRes  *models.Session
	err      error
}

func (m *MocDB) CreateSession(session *models.Session) error {
	if _, ok := m.Storage[session.ID]; ok {
		return fmt.Errorf("session already exists")
	}
	m.Storage[session.ID] = session
	return nil
}

func (m *MocDB) FindSession(sessID string) (*models.Session, error) {
	session, ok := m.Storage[sessID]
	if !ok {
		return nil, myErr.ErrNoAuth
	}
	return session, nil
}

func (m *MocDB) DestroySession(sessID string) error {
	if _, ok := m.Storage[sessID]; !ok {
		return fmt.Errorf("session not found")
	}
	delete(m.Storage, sessID)
	return nil
}

var (
	db = &MocDB{
		Storage: map[string]*models.Session{
			"2ht4k8s6v0m4hgl1": &models.Session{
				ID:     "2ht4k8s6v0m4hgl1",
				UserID: 1,
			},
		},
	}
	testCookie = &http.Cookie{
		Name:    "session_id",
		Value:   "2ht4k8s6v0m4hgl1",
		Expires: time.Now().Add(24 * time.Hour),
	}
	wrongTestCookie = &http.Cookie{
		Name:    "sessionId",
		Value:   "2ht4k8sg30m4hgl1",
		Expires: time.Now().Add(24 * time.Hour),
	}
	testSession = &models.Session{
		ID:     "2ht4k8s6v0m4hgl1",
		UserID: 1,
	}
	sm            = NewSessionManager(db)
	goodCookieReq = httptest.NewRequest(http.MethodGet, "/", nil)
	badCookieReq  = httptest.NewRequest(http.MethodGet, "/", nil)
)

func TestCheck(t *testing.T) {
	goodCookieReq.AddCookie(testCookie)

	badCookieReq.AddCookie(wrongTestCookie)

	tests := []Test{
		{
			testReq: httptest.NewRequest(http.MethodGet, "/", nil),
			err:     myErr.ErrNoAuth,
		},
		{
			testReq: goodCookieReq,
			testRes: testSession,
			err:     nil,
		},
		{
			testReq: badCookieReq,
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
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
		if !reflect.DeepEqual(res, test.testRes) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.testRes, res)
		}
	}
}

func TestCreate(t *testing.T) {
	//TODO
}
