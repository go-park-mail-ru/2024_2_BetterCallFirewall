package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/2024_2_BetterCallFirewall/internal/models"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type AuthService interface {
	Register(user models.User, ctx context.Context) (uint32, error)
	Auth(user models.User, ctx context.Context) (uint32, error)
}

type SessionManager interface {
	Check(r *http.Request) (*models.Session, error)
	Create(w http.ResponseWriter, userID uint32) (*models.Session, error)
	Destroy(w http.ResponseWriter, r *http.Request) error
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestID string)

	ErrorWrongMethod(w http.ResponseWriter, err error, requestID string)
	ErrorBadRequest(w http.ResponseWriter, err error, requestID string)
	ErrorInternal(w http.ResponseWriter, err error, requestID string)
	LogError(err error, requestID string)
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
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		c.responder.LogError(myErr.ErrInvalidContext, "")
	}

	if r.Method != http.MethodPost {
		c.responder.ErrorWrongMethod(w, errors.New("method not allowed"), reqID)
		return
	}

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router register: %w", err), reqID)
		return
	}

	user.ID, err = c.serviceAuth.Register(user, r.Context())
	if errors.Is(err, myErr.ErrUserAlreadyExists) || errors.Is(err, myErr.ErrNonValidEmail) || errors.Is(err, bcrypt.ErrPasswordTooLong) {
		c.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router register: %w", err), reqID)
		return
	}

	_, err = c.SessionManager.Create(w, user.ID)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router register: %w", err), reqID)
		return
	}

	c.responder.OutputJSON(w, "user create successful", reqID)
}

func (c *AuthController) Auth(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		c.responder.LogError(myErr.ErrInvalidContext, "")
	}

	if r.Method != http.MethodPost {
		c.responder.ErrorWrongMethod(w, errors.New("method not allowed"), reqID)
		return
	}

	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router auth: %w", err), reqID)
		return
	}

	id, err := c.serviceAuth.Auth(user, r.Context())

	if errors.Is(err, myErr.ErrWrongEmailOrPassword) || errors.Is(err, myErr.ErrNonValidEmail) {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router auth: %w", err), reqID)
		return
	}

	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router auth: %w", err), reqID)
		return
	}

	_, err = c.SessionManager.Create(w, id)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router auth: %w", err), reqID)
		return
	}

	c.responder.OutputJSON(w, "user auth", reqID)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		c.responder.LogError(myErr.ErrInvalidContext, "")
	}

	if r.Method != http.MethodPost {
		c.responder.ErrorWrongMethod(w, errors.New("method not allowed"), reqID)
		return
	}

	err := c.SessionManager.Destroy(w, r)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router logout: %w", err), reqID)
		return
	}

	c.responder.OutputJSON(w, "user logout", reqID)
}
