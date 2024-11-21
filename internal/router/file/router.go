package file

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

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

func NewRouter(fc FileController, sm SessionManager, logger *logrus.Logger) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/image/{name}", fc.Upload).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/image", fc.Upload).Methods(http.MethodPost, http.MethodOptions)

	res := middleware.Auth(sm, router)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)

	return res
}
