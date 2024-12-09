package file

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

type FileController interface {
	Upload(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
}

func NewRouter(
	fc FileController, sm SessionManager, logger *logrus.Logger, fileMetric *metrics.FileMetrics,
) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/image/{name}", fc.Upload).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/image", fc.Download).Methods(http.MethodPost, http.MethodOptions)

	router.Handle("/api/v1/metrics", promhttp.Handler())

	res := middleware.Auth(sm, router)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)
	res = middleware.FileMetricsMiddleware(fileMetric, res)

	return res
}
