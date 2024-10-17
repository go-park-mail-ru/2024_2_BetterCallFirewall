package router

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/middleware"
)

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type PostController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
}

func NewAuthRouter(authControl AuthController, postControl PostController, sm middleware.SessionManager) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/auth/register", authControl.Register)
	mux.HandleFunc("/api/v1/auth/login", authControl.Auth)
	mux.HandleFunc("/api/v1/auth/logout", authControl.Logout)
	mux.HandleFunc("/api/v1/post", postControl.GetAll)
	res := middleware.Auth(sm, mux)
	res = middleware.AccessLog(logrus.New(), res)

	return res
}
