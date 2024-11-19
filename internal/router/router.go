package router

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	_ "github.com/swaggo/http-swagger"

	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type PostController interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetBatchPosts(w http.ResponseWriter, r *http.Request)
}

type ProfileController interface {
	GetProfile(w http.ResponseWriter, r *http.Request)
	GetProfileById(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	DeleteProfile(w http.ResponseWriter, r *http.Request)
	GetHeader(w http.ResponseWriter, r *http.Request)

	SendFriendReq(w http.ResponseWriter, r *http.Request)
	AcceptFriendReq(w http.ResponseWriter, r *http.Request)
	Unsubscribe(w http.ResponseWriter, r *http.Request)
	RemoveFromFriends(w http.ResponseWriter, r *http.Request)
	GetAllFriends(w http.ResponseWriter, r *http.Request)
	GetAllSubs(w http.ResponseWriter, r *http.Request)
	GetAllSubscriptions(w http.ResponseWriter, r *http.Request)

	GetCommunitySubs(w http.ResponseWriter, r *http.Request)
}

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

type CommunityController interface {
	GetOne(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}

func NewRouter(
	profileControl ProfileController,
	postControl PostController,
	sm SessionManager,
	communityController CommunityController,
	logger *logrus.Logger,
) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/profile/header", profileControl.GetHeader).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile", profileControl.GetProfile).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}", profileControl.GetProfileById).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profiles", profileControl.GetAll).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile", profileControl.UpdateProfile).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("api/v1/profile", profileControl.DeleteProfile).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}/friend/subscribe", profileControl.SendFriendReq).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}/friend/accept", profileControl.AcceptFriendReq).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}/friend/unsubscribe", profileControl.RemoveFromFriends).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}/friend/remove", profileControl.Unsubscribe).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}/friends", profileControl.GetAllFriends).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}/subscribers", profileControl.GetAllSubs).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}/subscriptions", profileControl.GetAllSubscriptions).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/v1/feed", postControl.Create).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}", postControl.GetOne).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}", postControl.Update).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}", postControl.Delete).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/feed", postControl.GetBatchPosts).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/v1/community", communityController.Create).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}", communityController.GetOne).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}", communityController.Update).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}", communityController.Delete).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/community", communityController.GetAll).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}/subs", profileControl.GetCommunitySubs).Methods(http.MethodGet, http.MethodOptions)

	res := middleware.Auth(sm, router)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)

	return res
}
