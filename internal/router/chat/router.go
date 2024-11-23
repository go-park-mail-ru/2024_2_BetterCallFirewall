package chat

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/metrics"
	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type ChatController interface {
	SetConnection(w http.ResponseWriter, r *http.Request)
	GetAllChats(w http.ResponseWriter, r *http.Request)
	GetChat(w http.ResponseWriter, r *http.Request)
}

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

func NewRouter(cc ChatController, sm SessionManager, logger *logrus.Logger, chatMetrics *metrics.HttpMetrics) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/messages/chats", cc.GetAllChats).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/messages/chat/{id}", cc.GetChat).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/messages/ws", cc.SetConnection)
	router.Handle("/api/v1/metrics", promhttp.Handler())

	res := middleware.Auth(sm, router)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)
	res = middleware.HttpMetricsMiddleware(chatMetrics, res)

	return res
}
