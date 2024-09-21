package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type AuthService interface {
	Register(user models.User) error
	Auth(user models.User) error
	VerifyToken(token string) error
}

type AuthController struct {
	responder   Responder
	serviceAuth AuthService
}

func NewAuthController(responder Responder, serviceAuth AuthService) *AuthController {
	return &AuthController{
		responder:   responder,
		serviceAuth: serviceAuth,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("controller register: %w", err))
	}
	err = c.serviceAuth.Register(user)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("controller register: %w", err))
	}
	c.responder.OutputJSON(w, Message{msg: "user create successful"})
}
