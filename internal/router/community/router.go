package community

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type CommunityController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	SearchCommunity(w http.ResponseWriter, r *http.Request)
}

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

func NewRouter(communityController CommunityController, sm SessionManager, logger *logrus.Logger) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/community", communityController.Create).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}", communityController.GetOne).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}", communityController.Update).Methods(http.MethodPut, http.MethodOptions)
	router.HandleFunc("/api/v1/community/{id}", communityController.Delete).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/api/v1/community", communityController.GetAll).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/community/search", communityController.SearchCommunity).Methods(http.MethodGet, http.MethodOptions)

	res := middleware.Auth(sm, router)
	res = middleware.Preflite(res)
	res = middleware.AccessLog(logger, res)

	return res
}
