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
	testId   uint32
	testRes  *models.Session
	err      error
}

func (m *MocDB) CreateSession(session *models.Session) error {
	if _, ok := m.Storage[session.ID]; ok {
		return myErr.ErrSessionAlreadyExists
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

const (
	IdInBase        = 1
	IdNotInBase     = 2
	CookieInBase    = "2ht4k8s6v0m4hgl1"
	CookieNotInBase = "2ht4k8sg30m4hgl1"
)

var (
	db = &MocDB{
		Storage: map[string]*models.Session{
			"2ht4k8s6v0m4hgl1": &models.Session{
				ID:     "2ht4k8s6v0m4hgl1",
				UserID: IdInBase,
			},
		},
	}
	baseCookie = &http.Cookie{
		Name:    "session_id",
		Value:   CookieInBase,
		Expires: time.Now().Add(24 * time.Hour),
	}
	newCookie = &http.Cookie{
		Name:    "sessionId",
		Value:   CookieNotInBase,
		Expires: time.Now().Add(24 * time.Hour),
	}
	sm                 = NewSessionManager(db)
	CookieInBaseReq    = httptest.NewRequest(http.MethodGet, "/", nil)
	CookieNotInBaseReq = httptest.NewRequest(http.MethodGet, "/", nil)
)

func TestCheck(t *testing.T) {
	CookieInBaseReq.AddCookie(baseCookie)
	CookieNotInBaseReq.AddCookie(newCookie)

	tests := []Test{
		{
			testReq: httptest.NewRequest(http.MethodGet, "/", nil),
			err:     myErr.ErrNoAuth,
		},
		{
			testReq: CookieInBaseReq,
			testRes: &models.Session{
				ID:     CookieInBase,
				UserID: IdInBase,
			},
			err: nil,
		},
		{
			testReq: CookieNotInBaseReq,
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
	tests := []Test{
		{
			testResp: httptest.NewRecorder(),
			testId:   IdNotInBase,
			testRes: &models.Session{
				UserID: IdNotInBase,
			},
			err: nil,
		},
		{
			testResp: httptest.NewRecorder(),
			testId:   IdInBase,
			testRes:  nil,
			err:      myErr.ErrSessionAlreadyExists,
		},
	}

	for caseNum, test := range tests {
		res, err := sm.Create(test.testResp, test.testId)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
		if test.testRes.UserID != res.UserID {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.testRes, res)
		}
		expCookie := &http.Cookie{
			Name:    "session_id",
			Value:   res.ID,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Second),
		}
		resCookie := test.testResp.Header().Get("Set-Cookie")
		if !reflect.DeepEqual(expCookie.String(), resCookie) {
			t.Errorf("[%d] wrong cookie, expected: %#v, got: %#v", caseNum, expCookie, resCookie)
		}
	}
}

func TestDestroy(t *testing.T) {}
