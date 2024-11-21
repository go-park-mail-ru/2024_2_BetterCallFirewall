package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

type MocSessDB struct {
	Storage map[string]*models.Session
}

type Test struct {
	testCookie  string
	testId      uint32
	testSession *models.Session
	testRes     *models.Session
	err         error
}

func (m *MocSessDB) CreateSession(session *models.Session) error {
	for _, val := range m.Storage {
		if val.UserID == session.UserID {
			return my_err.ErrSessionAlreadyExists
		}
	}
	m.Storage[session.ID] = session
	return nil
}

func (m *MocSessDB) FindSession(sessID string) (*models.Session, error) {
	session, ok := m.Storage[sessID]
	if !ok {
		return nil, my_err.ErrNoAuth
	}
	return session, nil
}

func (m *MocSessDB) DestroySession(sessID string) error {
	if _, ok := m.Storage[sessID]; !ok {
		return my_err.ErrSessionNotFound
	}
	return nil
}

const (
	IdInBase        = 1
	IdNotInBase     = 2
	CookieInBase    = "2ht4k8s6v0m4hgl1"
	CookieNotInBase = "2ht4k8sg30m4hgl1"
)

var (
	activeSession = &models.Session{
		ID:     CookieInBase,
		UserID: IdInBase,
	}
	inactiveSession = &models.Session{
		ID:     CookieNotInBase,
		UserID: IdNotInBase,
	}
	db = &MocSessDB{
		Storage: map[string]*models.Session{
			activeSession.ID: activeSession,
		},
	}
	sm = NewSessionManager(db)
)

func TestCheck(t *testing.T) {
	tests := []Test{
		{
			testCookie: CookieNotInBase,
			err:        my_err.ErrNoAuth,
		},
		{
			testCookie: CookieInBase,
			testRes: &models.Session{
				ID:     CookieInBase,
				UserID: IdInBase,
			},
			err: nil,
		},
		{
			testCookie: "",
			err:        my_err.ErrNoAuth,
		},
	}

	for caseNum, test := range tests {
		res, err := sm.Check(test.testCookie)
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

func TestCreateSession(t *testing.T) {
	tests := []Test{
		{
			testId: IdNotInBase,
			testRes: &models.Session{
				UserID: IdNotInBase,
			},
			err: nil,
		},
		{
			testId:  IdInBase,
			testRes: nil,
			err:     my_err.ErrSessionAlreadyExists,
		},
	}

	for caseNum, test := range tests {
		res, err := sm.Create(test.testId)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
		if res != nil && res.UserID != test.testRes.UserID {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, test.testRes, res)
		}
	}
}

func TestDestroy(t *testing.T) {
	tests := []Test{
		{
			testSession: activeSession,
			err:         nil,
		},
		{
			testSession: inactiveSession,
			err:         my_err.ErrSessionNotFound,
		},
		{
			testSession: nil,
			err:         my_err.ErrNoAuth,
		},
	}

	for caseNum, test := range tests {
		err := sm.Destroy(test.testSession)
		if err != nil && test.err == nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && test.err != nil {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !errors.Is(err, test.err) {
			t.Errorf("[%d] wrong error, expected: %#v, got: %#v", caseNum, test.err, err)
		}
	}
}
