package controller

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/auth/controller"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
	"github.com/2024_2_BetterCallFirewall/internal/profile"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any)

	ErrorWrongMethod(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type ProfileHandlerImplementation struct {
	ProfileManager profile.ProfileUsecase
	Responder      controller.Responder
}

func NewProfileController(manager profile.ProfileUsecase, responder Responder) *ProfileHandlerImplementation {
	return &ProfileHandlerImplementation{
		ProfileManager: manager,
		Responder:      responder,
	}
}

func (h *ProfileHandlerImplementation) GetProfile(w http.ResponseWriter, r *http.Request) {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	userId := sess.UserID
	userProfile, err := h.ProfileManager.GetProfileById(r.Context(), userId)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, userProfile)
}

func (h *ProfileHandlerImplementation) GetAllProfiles(w http.ResponseWriter, r *http.Request) {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	userId := sess.UserID
	profiles, err := h.ProfileManager.GetAll(r.Context(), userId)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, profiles)
}

func (h *ProfileHandlerImplementation) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	newProfile := models.FullProfile{}
	err := json.NewDecoder(r.Body).Decode(&newProfile)
	r.Body.Close()
	if err != nil {
		h.Responder.ErrorBadRequest(w, fmt.Errorf("update error:%w", err))
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorBadRequest(w, myErr.ErrSessionNotFound)
		return
	}
	userId := sess.UserID

	err = h.ProfileManager.UpdateProfile(userId, &newProfile)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, newProfile)
}

func (h *ProfileHandlerImplementation) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorBadRequest(w, myErr.ErrSessionNotFound)
		return
	}
	userId := sess.UserID
	err = h.ProfileManager.DeleteProfile(userId)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	http.Redirect(w, r, "/api/v1/auth/logout", http.StatusContinue)
	return
}

func GetIdFromQuery(r *http.Request) (uint32, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return 0, myErr.ErrEmptyId
	}

	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, err
	}
	if uid > math.MaxInt {
		return 0, myErr.ErrBigId
	}
	return uint32(uid), nil
}

func (h *ProfileHandlerImplementation) GetProfileById(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromQuery(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err)
		return
	}

	profile, err := h.ProfileManager.GetProfileById(r.Context(), id)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, profile)
}

func (h *ProfileHandlerImplementation) GetAll(w http.ResponseWriter, r *http.Request) {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorBadRequest(w, err)
		return
	}
	uid := sess.UserID
	profiles, err := h.ProfileManager.GetAll(r.Context(), uid)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, profiles)
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
	receiver, sender, err := GetReceiverAndSender(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err)
		return
	}

	err = h.ProfileManager.SendFriendReq(receiver, sender)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, "success")

}

func (h *ProfileHandlerImplementation) AcceptFriendReq(w http.ResponseWriter, r *http.Request) {
	whose, who, err := GetReceiverAndSender(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err)
		return
	}
	err = h.ProfileManager.AcceptFriendReq(who, whose)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, "success")
}

func (h *ProfileHandlerImplementation) RemoveFromFriends(w http.ResponseWriter, r *http.Request) {
	whose, who, err := GetReceiverAndSender(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err)
		return
	}
	err = h.ProfileManager.RemoveFromFriends(who, whose)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, "success")
}

func (h *ProfileHandlerImplementation) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	whose, who, err := GetReceiverAndSender(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err)
		return
	}
	err = h.ProfileManager.Unsubscribe(who, whose)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, "success")
}

func (h *ProfileHandlerImplementation) GetAllFriends(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromQuery(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err)
		return
	}
	profiles, err := h.ProfileManager.GetAllFriends(r.Context(), id)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, profiles)
}

func (h *ProfileHandlerImplementation) GetAllSubs(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromQuery(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err)
		return
	}
	profiles, err := h.ProfileManager.GetAllSubs(r.Context(), id)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, profiles)
}

func (h *ProfileHandlerImplementation) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromQuery(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err)
		return
	}
	profiles, err := h.ProfileManager.GetAllSubscriptions(r.Context(), id)
	if err != nil {
		h.Responder.ErrorInternal(w, err)
		return
	}
	h.Responder.OutputJSON(w, profiles)
}
