package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/profile"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type Responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestID string)
	OutputNoMoreContentJSON(w http.ResponseWriter, requestID string)

	ErrorBadRequest(w http.ResponseWriter, err error, requestID string)
	ErrorInternal(w http.ResponseWriter, err error, requestID string)
	LogError(err error, requestID string)
}

type ProfileHandlerImplementation struct {
	ProfileManager profile.ProfileUsecase
	Responder      Responder
}

func NewProfileController(manager profile.ProfileUsecase, responder Responder) *ProfileHandlerImplementation {
	return &ProfileHandlerImplementation{
		ProfileManager: manager,
		Responder:      responder,
	}
}

func (h *ProfileHandlerImplementation) GetHeader(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorBadRequest(w, fmt.Errorf("update profile: %w", my_err.ErrSessionNotFound), reqID)
		return
	}

	userId := sess.UserID
	header, err := h.ProfileManager.GetHeader(r.Context(), userId)
	if err != nil {
		if errors.Is(err, my_err.ErrProfileNotFound) {
			h.Responder.ErrorBadRequest(w, err, reqID)
			return
		}
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}

	h.Responder.OutputJSON(w, &header, reqID)
}

func (h *ProfileHandlerImplementation) GetProfile(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorBadRequest(w, fmt.Errorf("update profile: %w", my_err.ErrSessionNotFound), reqID)
		return
	}
	userId := sess.UserID
	userProfile, err := h.ProfileManager.GetProfileById(r.Context(), userId)
	if err != nil {
		if errors.Is(err, my_err.ErrProfileNotFound) {
			h.Responder.ErrorBadRequest(w, err, reqID)
			return
		}
		if errors.Is(err, my_err.ErrNoMoreContent) {
			h.Responder.OutputNoMoreContentJSON(w, reqID)
			return
		}
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}

	h.Responder.OutputJSON(w, userProfile, reqID)
}

func (h *ProfileHandlerImplementation) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		h.Responder.ErrorBadRequest(w, fmt.Errorf("update profile: %w", my_err.ErrSessionNotFound), reqID)
		return
	}

	newProfile, err := h.getNewProfile(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}
	newProfile.ID = sess.UserID

	err = h.ProfileManager.UpdateProfile(r.Context(), newProfile)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}

	h.Responder.OutputJSON(w, newProfile, reqID)
}

func (h *ProfileHandlerImplementation) getNewProfile(r *http.Request) (*models.FullProfile, error) {
	newProfile := models.FullProfile{}
	err := json.NewDecoder(r.Body).Decode(&newProfile)
	if err != nil {
		return nil, err
	}

	return &newProfile, nil
}

func (h *ProfileHandlerImplementation) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		sess, err = models.SessionFromContext(r.Context())
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, my_err.ErrSessionNotFound, reqID)
		return
	}

	userId := sess.UserID
	err = h.ProfileManager.DeleteProfile(userId)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}

	h.Responder.OutputJSON(w, "delete", reqID)
}

func GetIdFromURL(r *http.Request) (uint32, error) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return 0, my_err.ErrEmptyId
	}

	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, err
	}
	if uid > math.MaxInt {
		return 0, my_err.ErrBigId
	}
	return uint32(uid), nil
}

func (h *ProfileHandlerImplementation) GetProfileById(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		id, err   = GetIdFromURL(r)
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	profile, err := h.ProfileManager.GetProfileById(r.Context(), id)
	if err != nil {
		if errors.Is(err, my_err.ErrProfileNotFound) {
			h.Responder.ErrorBadRequest(w, err, reqID)
			return
		}
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}

	h.Responder.OutputJSON(w, profile, reqID)
}

func GetLastId(r *http.Request) (uint32, error) {
	strLastId := r.URL.Query().Get("last_id")
	if strLastId == "" {
		return 0, nil
	}
	lastId, err := strconv.ParseUint(strLastId, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(lastId), nil
}

func (h *ProfileHandlerImplementation) GetAll(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		sess, err = models.SessionFromContext(r.Context())
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	lastId, err := GetLastId(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}
	uid := sess.UserID
	profiles, err := h.ProfileManager.GetAll(r.Context(), uid, lastId)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	if len(profiles) == 0 {
		h.Responder.OutputNoMoreContentJSON(w, reqID)
		return
	}
	h.Responder.OutputJSON(w, profiles, reqID)
}

func GetReceiverAndSender(r *http.Request) (uint32, uint32, error) {
	id, err := GetIdFromURL(r)
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
	var (
		reqID, ok             = r.Context().Value("requestID").(string)
		receiver, sender, err = GetReceiverAndSender(r)
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	err = h.ProfileManager.SendFriendReq(receiver, sender)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	h.Responder.OutputJSON(w, "success", reqID)
}

func (h *ProfileHandlerImplementation) AcceptFriendReq(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok       = r.Context().Value("requestID").(string)
		whose, who, err = GetReceiverAndSender(r)
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}
	err = h.ProfileManager.AcceptFriendReq(who, whose)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	h.Responder.OutputJSON(w, "success", reqID)
}

func (h *ProfileHandlerImplementation) RemoveFromFriends(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok       = r.Context().Value("requestID").(string)
		whose, who, err = GetReceiverAndSender(r)
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}
	err = h.ProfileManager.RemoveFromFriends(who, whose)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	h.Responder.OutputJSON(w, "success", reqID)
}

func (h *ProfileHandlerImplementation) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	whose, who, err := GetReceiverAndSender(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}
	err = h.ProfileManager.Unsubscribe(who, whose)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	h.Responder.OutputJSON(w, "success", reqID)
}

func (h *ProfileHandlerImplementation) GetAllFriends(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		id, err   = GetIdFromURL(r)
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	lastId, err := GetLastId(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	profiles, err := h.ProfileManager.GetAllFriends(r.Context(), id, lastId)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	if len(profiles) == 0 {
		h.Responder.OutputNoMoreContentJSON(w, reqID)
		return
	}

	h.Responder.OutputJSON(w, profiles, reqID)
}

func (h *ProfileHandlerImplementation) GetAllSubs(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		id, err   = GetIdFromURL(r)
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	lastId, err := GetLastId(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}
	profiles, err := h.ProfileManager.GetAllSubs(r.Context(), id, lastId)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	if len(profiles) == 0 {
		h.Responder.OutputNoMoreContentJSON(w, reqID)
		return
	}
	h.Responder.OutputJSON(w, profiles, reqID)
}

func (h *ProfileHandlerImplementation) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		id, err   = GetIdFromURL(r)
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	lastId, err := GetLastId(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	profiles, err := h.ProfileManager.GetAllSubscriptions(r.Context(), id, lastId)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}
	if len(profiles) == 0 {
		h.Responder.OutputNoMoreContentJSON(w, reqID)
		return
	}

	h.Responder.OutputJSON(w, profiles, reqID)
}

func (h *ProfileHandlerImplementation) GetCommunitySubs(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		id, err   = GetIdFromURL(r)
	)

	if !ok {
		h.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	lastId, err := GetLastId(r)
	if err != nil {
		h.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	subs, err := h.ProfileManager.GetCommunitySubs(r.Context(), id, lastId)
	if err != nil {
		h.Responder.ErrorInternal(w, err, reqID)
		return
	}

	if len(subs) == 0 {
		h.Responder.OutputNoMoreContentJSON(w, reqID)
		return
	}

	h.Responder.OutputJSON(w, subs, reqID)
}
