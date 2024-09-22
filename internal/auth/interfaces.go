package auth

import (
	"errors"
	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"net/http"
)

var (
	ErrNoAuth = errors.New("no session found")
)

type UserManager interface {
	Authorize(string, string) (*models.User, error)
	Register(string, string, string) (*models.User, error)
}

type SessionManager interface {
	Check(r *http.Request) (*models.Session, error)                       //TODO сделать через контексты
	Create(w http.ResponseWriter, userID uint32) (*models.Session, error) //TODO
}
