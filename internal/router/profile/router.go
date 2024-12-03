package profile

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ProfileController interface {
	GetProfile(w http.ResponseWriter, r *http.Request)
	GetProfileById(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	DeleteProfile(w http.ResponseWriter, r *http.Request)
	GetHeader(w http.ResponseWriter, r *http.Request)
	SearchProfile(w http.ResponseWriter, r *http.Request)

	SendFriendReq(w http.ResponseWriter, r *http.Request)
	AcceptFriendReq(w http.ResponseWriter, r *http.Request)
	Unsubscribe(w http.ResponseWriter, r *http.Request)
	RemoveFromFriends(w http.ResponseWriter, r *http.Request)
	GetAllFriends(w http.ResponseWriter, r *http.Request)
	GetAllSubs(w http.ResponseWriter, r *http.Request)
	GetAllSubscriptions(w http.ResponseWriter, r *http.Request)

	GetCommunitySubs(w http.ResponseWriter, r *http.Request)
	ChangePassword(w http.ResponseWriter, r *http.Request)
}

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

func NewRouter(
	profileControl ProfileController,
	sm SessionManager,
	logger *logrus.Logger,
	httpMetric *metrics.HttpMetrics,
) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/profile/header", profileControl.GetHeader).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile", profileControl.GetProfile).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}", profileControl.GetProfileById).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profiles", profileControl.GetAll).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/profile", profileControl.UpdateProfile).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("api/v1/profile", profileControl.DeleteProfile).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}/friend/subscribe", profileControl.SendFriendReq).Methods(
		http.MethodPost, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/profile/{id}/friend/accept", profileControl.AcceptFriendReq).Methods(
		http.MethodPost, http.MethodOptions,
	)
	router.HandleFunc(
		"/api/v1/profile/{id}/friend/unsubscribe", profileControl.RemoveFromFriends,
	).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/profile/{id}/friend/remove", profileControl.Unsubscribe).Methods(
		http.MethodDelete, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/profile/{id}/friends", profileControl.GetAllFriends).Methods(
		http.MethodGet, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/profile/{id}/subscribers", profileControl.GetAllSubs).Methods(
		http.MethodGet, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/profile/{id}/subscriptions", profileControl.GetAllSubscriptions).Methods(
		http.MethodGet, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/profile/community/{id}/subs", profileControl.GetCommunitySubs).Methods(
		http.MethodGet, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/profile/search/", profileControl.SearchProfile).Methods(
		http.MethodGet, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/profile/password", profileControl.ChangePassword).Methods(
		http.MethodPut, http.MethodOptions,
	)

	router.Handle("/api/v1/metrics", promhttp.Handler())

	res := middleware.Auth(sm, router)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)
	res = middleware.HttpMetricsMiddleware(httpMetric, res)

	return res
}
