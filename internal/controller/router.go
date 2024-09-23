package controller

import (
	"net/http"
)

func NewAuthRouter(controller *AuthController) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", controller.Register)
	mux.HandleFunc("/auth/login", controller.Auth)

	return mux
}
