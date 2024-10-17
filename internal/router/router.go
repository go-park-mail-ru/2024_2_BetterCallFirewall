package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/swaggo/http-swagger"

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

type ProfileController interface {
	GetProfileById(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	DeleteProfile(w http.ResponseWriter, r *http.Request)

	SendFriendReq(w http.ResponseWriter, r *http.Request)
	AcceptFriendReq(w http.ResponseWriter, r *http.Request)
	RemoveFromFriends(w http.ResponseWriter, r *http.Request)
	GetAllFriends(w http.ResponseWriter, r *http.Request)
}

func NewAuthRouter(authControl AuthController, profileControl ProfileController, postControl PostController, sm middleware.SessionManager) http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/api/v1/auth/register", authControl.Register).Methods(http.MethodPost)
	mux.HandleFunc("/api/v1/auth/login", authControl.Auth).Methods(http.MethodPost)
	mux.HandleFunc("/api/v1/auth/logout", authControl.Logout).Methods(http.MethodPost)
	mux.HandleFunc("/api/v1/post", postControl.GetAll).Methods(http.MethodGet)

	mux.HandleFunc("/api/v1/profile/{id}", profileControl.GetProfileById).Methods(http.MethodGet)
	mux.HandleFunc("/api/v1/profiles", profileControl.GetAll).Methods(http.MethodGet)
	mux.HandleFunc("/api/v1/update_profile", profileControl.UpdateProfile).Methods(http.MethodPut)
	mux.HandleFunc("api/v1/delete_profile", profileControl.DeleteProfile).Methods(http.MethodDelete)
	mux.HandleFunc("/api/v1/send_friend_request/{id}", profileControl.SendFriendReq).Methods(http.MethodPost)
	mux.HandleFunc("/api/v1/accept_friend_request/{id}", profileControl.AcceptFriendReq).Methods(http.MethodPost)
	mux.HandleFunc("/api/v1/remove_friend/{id}", profileControl.RemoveFromFriends).Methods(http.MethodPost)
	mux.HandleFunc("/api/v1/get_friends/{id}", profileControl.GetAll).Methods(http.MethodGet)

	res := middleware.Auth(sm, mux)
	res = middleware.AccessLog(log.Default(), res)

	return res
}
