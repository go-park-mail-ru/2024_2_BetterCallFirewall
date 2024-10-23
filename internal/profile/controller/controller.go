package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/auth/controller"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/profile"
)

type ProfileHandlerImplementation struct {
	ProfileManager profile.ProfileUsecase
	Responder      controller.Responder
}

func NewProfileController(manager profile.ProfileUsecase, responder controller.Responder) *ProfileHandlerImplementation {
	return &ProfileHandlerImplementation{
		ProfileManager: manager,
		Responder:      responder,
	}
}

func (h *ProfileHandlerImplementation) GetProfile(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	userId := sess.UserID
	userProfile, err := h.ProfileManager.GetProfileById(r.Context(), userId)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	h.Responder.OutputJSON(w, userProfile, reqID)
}

func (h *ProfileHandlerImplementation) GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	userId := sess.UserID
	profiles, err := h.ProfileManager.GetAll(r.Context(), userId)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	h.Responder.OutputJSON(w, profiles, reqID)
}

func (h *ProfileHandlerImplementation) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)

	newProfile := models.FullProfile{}
	err := json.NewDecoder(r.Body).Decode(&newProfile)
	r.Body.Close()
	if err != nil {
		h.Responder.ErrorBadRequest(w, fmt.Errorf("update error:%w", err), reqID)
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorBadRequest(w, myErr.ErrSessionNotFound, reqID)
	}
	userId := sess.UserID

	err = h.ProfileManager.UpdateProfile(userId, &newProfile)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
	}
	h.Responder.OutputJSON(w, newProfile, reqID)
}

func (h *ProfileHandlerImplementation) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)
	sess, err := models.SessionFromContext(r.Context())

	if err != nil {
		h.Responder.ErrorBadRequest(w, myErr.ErrSessionNotFound, reqID)
	}
	userId := sess.UserID
	err = h.ProfileManager.DeleteProfile(userId)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
	}
	http.Redirect(w, r, "/api/v1/auth/logout", http.StatusContinue)
	return
}

func GetIdFromQuery(r *http.Request) (uint32, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return 0, errors.New("id is empty")
	}
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(uid), nil
}

func (h *ProfileHandlerImplementation) GetProfileById(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)
	id, err := GetIdFromQuery(r)

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
	}

	profile, err := h.ProfileManager.GetProfileById(r.Context(), id)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
	}
	h.Responder.OutputJSON(w, profile, reqID)
}

func (h *ProfileHandlerImplementation) GetAll(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)
	sess, err := models.SessionFromContext(r.Context())

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
	}
	uid := sess.UserID
	profiles, err := h.ProfileManager.GetAll(r.Context(), uid)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
	}
	h.Responder.OutputJSON(w, profiles, reqID)
}

func GetReceiverAndSender(r *http.Request) (uint32, uint32, error) {
	id, err := GetIdFromQuery(r)
	if err != nil {
		return 0, 0, err
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return 0, 0, err
	}

	return id, sess.UserID, nil
}

func (h *ProfileHandlerImplementation) SendFriendReq(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)
	receiver, sender, err := GetReceiverAndSender(r)

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
	}

	err = h.ProfileManager.SendFriendReq(receiver, sender)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
	}
	h.Responder.OutputJSON(w, "success", reqID)

}

func (h *ProfileHandlerImplementation) AcceptFriendReq(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)
	whose, who, err := GetReceiverAndSender(r)

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
	}
	err = h.ProfileManager.AcceptFriendReq(who, whose)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
	}
	h.Responder.OutputJSON(w, "success", reqID)
}

func (h *ProfileHandlerImplementation) RemoveFromFriends(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)
	whose, who, err := GetReceiverAndSender(r)

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
	}
	err = h.ProfileManager.RemoveFromFriends(r.Context(), who, whose)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
	}
	h.Responder.OutputJSON(w, "success", reqID)
}

func (h *ProfileHandlerImplementation) GetAllFriends(w http.ResponseWriter, r *http.Request) {
	reqID := r.Context().Value("requestID").(string)
	id, err := GetIdFromQuery(r)

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
	}
	profiles, err := h.ProfileManager.GetAllFriends(r.Context(), id)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
	}
	h.Responder.OutputJSON(w, profiles, reqID)
}
