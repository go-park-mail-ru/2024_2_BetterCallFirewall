package controller

import (
	"github.com/2024_2_BetterCallFirewall/internal/auth/controller"
	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
	"github.com/2024_2_BetterCallFirewall/internal/profile/service"
	"net/http"
)

type ProfileHandler struct {
	Repo      service.ProfileUsecase
	Responder controller.Responder
}

func NewProfileController(repo service.ProfileUsecase, responder controller.Responder) *ProfileHandler {
	return &ProfileHandler{
		Repo:      repo,
		Responder: responder,
	}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	sess := models.SessionFromContext(r.Context())
	userId := sess.UserID
	userProfile, err := h.Repo.GetProfileById(userId)
	if err != nil {
		h.Responder.
	}
}
