package community

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type CommunityController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	JoinToCommunity(w http.ResponseWriter, r *http.Request)
	LeaveFromCommunity(w http.ResponseWriter, r *http.Request)
	AddAdmin(w http.ResponseWriter, r *http.Request)
	SearchCommunity(w http.ResponseWriter, r *http.Request)
}

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

func NewRouter(
	communityController CommunityController, sm SessionManager, logger *logrus.Logger,
	communityMetrics *metrics.HttpMetrics,
) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/community", communityController.Create).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}", communityController.GetOne).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}", communityController.Update).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}", communityController.Delete).Methods(
		http.MethodDelete, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/community", communityController.GetAll).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}/join", communityController.JoinToCommunity).Methods(
		http.MethodPost, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/community/{id}/leave", communityController.LeaveFromCommunity).Methods(
		http.MethodPost, http.MethodOptions,
	)
	router.HandleFunc("api/v1/community/{id}/add_admin", communityController.AddAdmin).Methods(
		http.MethodPost, http.MethodOptions,
	)
	router.HandleFunc("/api/v1/community/search/", communityController.SearchCommunity).Methods(
		http.MethodGet, http.MethodOptions,
	)

	router.Handle("/api/v1/metrics", promhttp.Handler())
	router.Handle(
		"/", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)

	res := middleware.Auth(sm, router)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)
	res = middleware.HttpMetricsMiddleware(communityMetrics, res)

	return res
}
