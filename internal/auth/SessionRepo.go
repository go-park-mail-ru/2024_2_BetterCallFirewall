package auth

import (
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type SessionRepository interface {
	CreateSession(*models.Session) error
	FindSession(sessID string) (*models.Session, error)
	DestroySession(sessID string) error
}
