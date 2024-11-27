package auth

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

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Auth(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

func NewRouter(
	authControl AuthController, sm SessionManager, logger *logrus.Logger, httpMetrics *metrics.HttpMetrics,
) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/auth/register", authControl.Register).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/auth/login", authControl.Auth).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/auth/logout", authControl.Logout).Methods(http.MethodPost, http.MethodOptions)

	router.Handle("/api/v1/metrics", promhttp.Handler())
	router.Handle(
		"/", http.HandlerFunc(
			func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)

	res := middleware.Preflite(router)
	res = middleware.AccessLog(logger, res)
	res = middleware.HttpMetricsMiddleware(httpMetrics, res)
	return res
}
