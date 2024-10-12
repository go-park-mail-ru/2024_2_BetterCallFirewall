package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type AuthService interface {
	Register(user models.User) (uint32, error)
	Auth(user models.User) (uint32, error)
}

type SessionManager interface {
	Check(r *http.Request) (*models.Session, error)
	Create(w http.ResponseWriter, userID uint32) (*models.Session, error)
	Destroy(w http.ResponseWriter, r *http.Request) error
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any)

	ErrorWrongMethod(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type AuthController struct {
	responder      Responder
	serviceAuth    AuthService
	SessionManager SessionManager
}

func NewAuthController(responder Responder, serviceAuth AuthService, sessionManager SessionManager) *AuthController {
	return &AuthController{
		responder:      responder,
		serviceAuth:    serviceAuth,
		SessionManager: sessionManager,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router register: %w", err))
		return
	}

	user.ID, err = c.serviceAuth.Register(user)
	if errors.Is(err, myErr.ErrUserAlreadyExists) || errors.Is(err, myErr.ErrNonValidEmail) || errors.Is(err, bcrypt.ErrPasswordTooLong) {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router register: %w", err))
		return
	}

	_, err = c.SessionManager.Create(w, user.ID)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router register: %w", err))
		return
	}

	c.responder.OutputJSON(w, "user create successful")
}

func (c *AuthController) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router auth: %w", err))
		return
	}

	id, err := c.serviceAuth.Auth(user)

	if errors.Is(err, myErr.ErrWrongEmailOrPassword) || errors.Is(err, myErr.ErrNonValidEmail) {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router auth: %w", err))
		return
	}

	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router auth: %w", err))
		return
	}

	_, err = c.SessionManager.Create(w, id)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router auth: %w", err))
		return
	}

	c.responder.OutputJSON(w, "user auth")
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.responder.ErrorWrongMethod(w, errors.New("method not allowed"))
		return
	}

	err := c.SessionManager.Destroy(w, r)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router logout: %w", err))
		return
	}

	c.responder.OutputJSON(w, "user logout")
}
