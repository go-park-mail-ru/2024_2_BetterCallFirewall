package stickers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type Controller interface {
	AddNewSticker(w http.ResponseWriter, r *http.Request)
	GetAllStickers(w http.ResponseWriter, r *http.Request)
	GetMineStickers(w http.ResponseWriter, r *http.Request)
}

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

func NewRouter(controller Controller, sm SessionManager, logger *logrus.Logger) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/stickers/all", controller.GetAllStickers).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/stickers", controller.AddNewSticker).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/stickers", controller.GetMineStickers).Methods(http.MethodGet, http.MethodOptions)

	res := middleware.Auth(sm, router)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)

	return res
}
