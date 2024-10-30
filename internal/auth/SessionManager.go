package auth

import (
	"net/http"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*http.Cookie, error)
	Destroy(sess *models.Session) error
}
