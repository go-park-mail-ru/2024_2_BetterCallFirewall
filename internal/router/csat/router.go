package csat

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type csatController interface {
	SaveMetrics(w http.ResponseWriter, r *http.Request)
	GetMetrics(w http.ResponseWriter, r *http.Request)
	CheckExperience(w http.ResponseWriter, r *http.Request)
}

type sessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

func NewRouter(controller csatController, sm sessionManager, logger *logrus.Logger) http.Handler {
	rout := mux.NewRouter()

	rout.HandleFunc("/api/v1/csat", controller.SaveMetrics).Methods(http.MethodPost, http.MethodOptions)
	rout.HandleFunc("/api/v1/csat/metrics", controller.GetMetrics).Methods(http.MethodGet, http.MethodOptions)
	rout.HandleFunc("/api/v1/csat", controller.CheckExperience).Methods(http.MethodGet, http.MethodOptions)

	res := middleware.Auth(sm, rout)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)

	return res
}
