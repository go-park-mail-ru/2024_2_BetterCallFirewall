package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/2024_2_BetterCallFirewall/internal/auth"
	"github.com/2024_2_BetterCallFirewall/internal/models"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type AuthService interface {
	Register(user models.User, ctx context.Context) (uint32, error)
	Auth(user models.User, ctx context.Context) (uint32, error)
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
	SessionManager auth.SessionManager
}

func NewAuthController(responder Responder, serviceAuth AuthService, sessionManager auth.SessionManager) *AuthController {
	return &AuthController{
		responder:      responder,
		serviceAuth:    serviceAuth,
		SessionManager: sessionManager,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router register: %w", err))
		return
	}

	user.ID, err = c.serviceAuth.Register(user, r.Context())
	if errors.Is(err, myErr.ErrUserAlreadyExists) || errors.Is(err, myErr.ErrNonValidEmail) || errors.Is(err, bcrypt.ErrPasswordTooLong) {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router register: %w", err))
		return
	}

	sess, err := c.SessionManager.Create(user.ID)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router register: %w", err))
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().AddDate(0, 0, 1),
	}

	http.SetCookie(w, cookie)

	c.responder.OutputJSON(w, "user create successful")
}

func (c *AuthController) Auth(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router auth: %w", err))
		return
	}

	id, err := c.serviceAuth.Auth(user, r.Context())

	if errors.Is(err, myErr.ErrWrongEmailOrPassword) || errors.Is(err, myErr.ErrNonValidEmail) {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router auth: %w", err))
		return
	}

	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router auth: %w", err))
		return
	}

	sess, err := c.SessionManager.Create(id)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router auth: %w", err))
		return
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().AddDate(0, 0, 1),
	}
	http.SetCookie(w, cookie)

	c.responder.OutputJSON(w, "user auth")
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		c.responder.ErrorBadRequest(w, myErr.ErrNoAuth)
	}
	err = c.SessionManager.Destroy(sess)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router logout: %w", err))
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().AddDate(0, 0, -1),
	}
	http.SetCookie(w, cookie)

	c.responder.OutputJSON(w, "user logout")
}
