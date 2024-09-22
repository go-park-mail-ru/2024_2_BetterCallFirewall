package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type AuthService interface {
	Register(user models.User) error
}

type SessionManager interface {
	Check(r *http.Request) (*models.Session, error)                       //TODO сделать через контексты
	Create(w http.ResponseWriter, userID uint32) (*models.Session, error) //TODO сделаь без w
	Destroy(w http.ResponseWriter, r *http.Request) error
}

type AuthController struct {
	responder      Responder
	serviceAuth    AuthService
	sessionManager SessionManager
}

func NewAuthController(responder Responder, serviceAuth AuthService, sessionManager SessionManager) *AuthController {
	return &AuthController{
		responder:      responder,
		serviceAuth:    serviceAuth,
		sessionManager: sessionManager,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("controller register: %w", err))
		return
	}
	err = c.serviceAuth.Register(user)
	if errors.Is(err, myErr.ErrUserAlreadyExists) {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("controller register: %w", err))
		return
	}

	_, err = c.sessionManager.Create(w, user.ID)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("controller register: %w", err))
	}

	c.responder.OutputJSON(w, Message{Msg: "user create successful"})
}
