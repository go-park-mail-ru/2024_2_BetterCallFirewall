package controller

import (
	"net/http"

	"github.com/2024_2_BetterCallFirewall/internal/middleware"
)

func NewAuthRouter(controller *AuthController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", controller.Register)
	mux.HandleFunc("/auth/login", controller.Auth)
	mux.HandleFunc("/profile", MockHandler)
	res := middleware.Auth(controller.sessionManager, mux)
	return res
}

func MockHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Don't work auth"))
}
