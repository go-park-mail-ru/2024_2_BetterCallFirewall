package service

import (
	"fmt"

	"github.com/2024_2_BetterCallFirewall/internal/auth"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type SessionManagerImpl struct {
	DB auth.SessionRepository
}

func NewSessionManager(DB auth.SessionRepository) *SessionManagerImpl {
	return &SessionManagerImpl{
		DB: DB,
	}
}

func (sm *SessionManagerImpl) Check(cookie string) (*models.Session, error) {
	sess, err := sm.DB.FindSession(cookie)
	if err != nil {
		return nil, fmt.Errorf("session check: %w", err)
	}

	return sess, nil
}

func (sm *SessionManagerImpl) Create(userID uint32) (*models.Session, error) {
	sess, err := models.NewSession(userID)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	err = sm.DB.CreateSession(sess)
	if err != nil {
		return nil, fmt.Errorf("session creation: %w", err)
	}
	return sess, nil
}

func (sm *SessionManagerImpl) Destroy(sess *models.Session) error {
	err := sm.DB.DestroySession(sess.ID)
	if err != nil {
		return fmt.Errorf("session destroy: %w", err)
	}

	return nil
}
