package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type SessionRepository interface {
	CreateSession(*models.Session) error
	FindSession(sessID string) (*models.Session, error)
	DestroySession(sessID string) error
}

type SessionManagerImpl struct {
	DB SessionRepository
}

func NewSessionManager(DB SessionRepository) *SessionManagerImpl {
	return &SessionManagerImpl{
		DB: DB,
	}
}

func (sm *SessionManagerImpl) Check(r *http.Request) (*models.Session, error) {
	sessionCookie, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		return nil, myErr.ErrNoAuth
	}

	sess, err := sm.DB.FindSession(sessionCookie.Value)
	if err != nil {
		return nil, fmt.Errorf("session check: %w", err)
	}

	return sess, nil
}

func (sm *SessionManagerImpl) Create(w http.ResponseWriter, userID uint32) (*models.Session, error) {
	sess, err := models.NewSession(userID)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	err = sm.DB.CreateSession(sess)
	if err != nil {
		return nil, fmt.Errorf("session creation: %w", err)
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	return sess, nil
}

func (sm *SessionManagerImpl) Destroy(w http.ResponseWriter, r *http.Request) error {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return fmt.Errorf("session destroy: %w", err)
	}
	err = sm.DB.DestroySession(sess.ID)
	if err != nil {
		return fmt.Errorf("session destroy: %w", err)
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    "/",
	}
	http.SetCookie(w, cookie)

	return nil
}
