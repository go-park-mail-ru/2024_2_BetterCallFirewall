package post

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

type Controller interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetBatchPosts(w http.ResponseWriter, r *http.Request)

	SetLikeOnPost(w http.ResponseWriter, r *http.Request)
	DeleteLikeFromPost(w http.ResponseWriter, r *http.Request)
}

func NewRouter(
	contr Controller, sm SessionManager, logger *logrus.Logger, postMetric *metrics.HttpMetrics,
) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/feed", contr.Create).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}", contr.GetOne).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}", contr.Update).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}", contr.Delete).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/feed", contr.GetBatchPosts).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/api/v1/feed/{id}/like", contr.SetLikeOnPost).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/feed/{id}/unlike", contr.DeleteLikeFromPost).Methods(http.MethodPost, http.MethodOptions)

	router.Handle("/api/v1/metrics", promhttp.Handler())
	router.Handle(
		"/", http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)

	res := middleware.Auth(sm, router)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)
	res = middleware.HttpMetricsMiddleware(postMetric, res)

	return res
}
