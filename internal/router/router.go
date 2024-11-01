package router

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	_ "github.com/swaggo/http-swagger"

	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type PostController interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetBatchPosts(w http.ResponseWriter, r *http.Request)
}

type ProfileController interface {
	GetProfileById(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	DeleteProfile(w http.ResponseWriter, r *http.Request)

	SendFriendReq(w http.ResponseWriter, r *http.Request)
	AcceptFriendReq(w http.ResponseWriter, r *http.Request)
	Unsubscribe(w http.ResponseWriter, r *http.Request)
	RemoveFromFriends(w http.ResponseWriter, r *http.Request)
	GetAllFriends(w http.ResponseWriter, r *http.Request)
	GetAllSubs(w http.ResponseWriter, r *http.Request)
	GetAllSubscriptions(w http.ResponseWriter, r *http.Request)
}

func NewRouter(
	authControl AuthController,
	profileControl ProfileController,
	postControl PostController,
	sm middleware.SessionManager,
	logger *logrus.Logger,
) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/auth/register", authControl.Register).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/auth/login", authControl.Auth).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/auth/logout", authControl.Logout).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/api/v1/profile/{id}", profileControl.GetProfileById).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profiles", profileControl.GetAll).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/update", profileControl.UpdateProfile).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("api/v1/profile/delete", profileControl.DeleteProfile).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/friend/subscribe/{id}", profileControl.SendFriendReq).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/friend/accept/{id}", profileControl.AcceptFriendReq).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/friend/unsubscribe/{id}", profileControl.Unsubscribe).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/friend/remove/{id}", profileControl.RemoveFromFriends).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/friends/{id}", profileControl.GetAllFriends).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/subscribers/{id}", profileControl.GetAllSubs).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/subscriptions/{id}", profileControl.GetAllSubscriptions).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/v1/feed", postControl.Create).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}", postControl.GetOne).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}", postControl.Update).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}", postControl.Delete).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/feed", postControl.GetBatchPosts).Methods(http.MethodGet, http.MethodOptions)

	res := middleware.Auth(sm, router)
	res = middleware.AccessLog(logger, res)

	return res
}
