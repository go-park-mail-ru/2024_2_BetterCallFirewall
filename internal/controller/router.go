package controller

import (
	"net/http"

	"github.com/2024_2_BetterCallFirewall/internal/middleware"
)

func NewAuthRouter(controller *AuthController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", controller.Register)
	mux.HandleFunc("/auth/login", controller.Auth)
	res := middleware.Auth(controller.sessionManager, mux)
	return res
}
