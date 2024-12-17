package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/mailru/easyjson"
	"golang.org/x/crypto/bcrypt"

	"github.com/2024_2_BetterCallFirewall/internal/auth"
	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"

	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

var emailRegex = regexp.MustCompile(`^[\w-.]+@([\w-]+\.)\w{2,4}$`)

type AuthService interface {
	Register(user models.User, ctx context.Context) (uint32, error)
	Auth(user models.User, ctx context.Context) (uint32, error)
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestID string)

	ErrorBadRequest(w http.ResponseWriter, err error, requestID string)
	ErrorInternal(w http.ResponseWriter, err error, requestID string)
	LogError(err error, requestID string)
}

type AuthController struct {
	responder      Responder
	serviceAuth    AuthService
	SessionManager auth.SessionManager
}

func NewAuthController(
	responder Responder, serviceAuth AuthService, sessionManager auth.SessionManager,
) *AuthController {
	return &AuthController{
		responder:      responder,
		serviceAuth:    serviceAuth,
		SessionManager: sessionManager,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value(middleware.RequestKey).(string)
	if !ok {
		c.responder.LogError(my_err.ErrInvalidContext, "")
	}

	user := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router register: %w", err), reqID)
		return
	}
	if !validate(user) {
		c.responder.ErrorBadRequest(w, my_err.ErrBadUserInfo, reqID)
		return
	}

	user.ID, err = c.serviceAuth.Register(user, r.Context())
	if errors.Is(err, my_err.ErrUserAlreadyExists) || errors.Is(err, my_err.ErrNonValidEmail) || errors.Is(
		err, bcrypt.ErrPasswordTooLong,
	) {
		c.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router register: %w", err), reqID)
		return
	}

	sess, err := c.SessionManager.Create(user.ID)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router register: %w", err), reqID)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().AddDate(0, 0, 1),
	}

	http.SetCookie(w, cookie)

	c.responder.OutputJSON(w, "user create successful", reqID)
}

func (c *AuthController) Auth(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value(middleware.RequestKey).(string)
	if !ok {
		c.responder.LogError(my_err.ErrInvalidContext, "")
	}

	user := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router auth: %w", err), reqID)
		return
	}
	if !validateAuth(user) {
		c.responder.ErrorBadRequest(w, my_err.ErrBadUserInfo, reqID)
		return
	}

	id, err := c.serviceAuth.Auth(user, r.Context())

	if errors.Is(err, my_err.ErrWrongEmailOrPassword) || errors.Is(err, my_err.ErrNonValidEmail) {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router auth: %w", err), reqID)
		return
	}

	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router auth: %w", err), reqID)
		return
	}

	sess, err := c.SessionManager.Create(id)
	if err != nil {
		c.responder.ErrorInternal(w, fmt.Errorf("router auth: %w", err), reqID)
		return
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().AddDate(0, 0, 1),
	}
	http.SetCookie(w, cookie)

	c.responder.OutputJSON(w, "user auth", reqID)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value(middleware.RequestKey).(string)
	if !ok {
		c.responder.LogError(my_err.ErrInvalidContext, "")
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		c.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, err := c.SessionManager.Check(sessionCookie.Value)
	if err != nil {
		c.responder.ErrorBadRequest(w, my_err.ErrNoAuth, reqID)
		return
	}

	err = c.SessionManager.Destroy(sess)
	if err != nil {
		c.responder.ErrorBadRequest(w, fmt.Errorf("router logout: %w", err), reqID)
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().AddDate(0, 0, -1),
	}
	http.SetCookie(w, cookie)

	c.responder.OutputJSON(w, "user logout", reqID)
}

func validate(user models.User) bool {
	if len([]rune(user.FirstName)) < 3 || len([]rune(user.LastName)) < 3 || len([]rune(user.Password)) < 6 ||
		len([]rune(user.FirstName)) > 30 || len([]rune(user.LastName)) > 30 {
		return false
	}
	return true
}

func validateAuth(user models.User) bool {
	if len([]rune(user.Password)) < 6 || !emailRegex.MatchString(user.Email) {
		return false
	}

	return true
}
