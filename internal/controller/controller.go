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

type SessionService interface {
	CreateSession(firstName string) error
}

type AuthController struct {
	responder      Responder
	serviceAuth    AuthService
	sessionService SessionService
}

func NewAuthController(responder Responder, serviceAuth AuthService, sessionService SessionService) *AuthController {
	return &AuthController{
		responder:      responder,
		serviceAuth:    serviceAuth,
		sessionService: sessionService,
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

	err = c.sessionService.CreateSession(user.FirstName)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("controller register: %w", err))
	}

	c.responder.OutputJSON(w, Message{Msg: "user create successful"})
}
