package controller

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/like"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestID string)
	OutputNoMoreContentJSON(w http.ResponseWriter, requestID string)

	ErrorBadRequest(w http.ResponseWriter, err error, requestID string)
	ErrorInternal(w http.ResponseWriter, err error, requestID string)
	LogError(err error, requestID string)
}

type LikeController struct {
	LikeManager like.ReactionUsecase
	Responder   Responder
}

func NewLikeController(likeManager like.ReactionUsecase) *LikeController {
	return &LikeController{
		LikeManager: likeManager,
	}
}

func getIdFromReq(r *http.Request) (uint32, error) {
	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		return 0, err
	}
	return sess.UserID, nil
}

func getIdFromURL(r *http.Request) (uint32, error) {
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

func getId(r *http.Request) (uint32, uint32, error) {
	userID, err := getIdFromReq(r)
	if err != nil {
		return 0, 0, err
	}
	postID, err := getIdFromURL(r)
	if err != nil {
		return 0, 0, err
	}
	return userID, postID, nil
}

func (l LikeController) SetLikeToPost(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		l.Responder.LogError(myErr.ErrInvalidContext, "")
	}

	userID, postID, err := getId(r)
	if err != nil {
		l.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	err = l.LikeManager.SetLikeToPost(r.Context(), postID, userID)
	if err != nil {
		l.Responder.ErrorInternal(w, err, reqID)
		return
	}

	l.Responder.OutputJSON(w, "like is set on post", reqID)
}

func (l LikeController) SetLikeToComment(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		l.Responder.LogError(myErr.ErrInvalidContext, "")
	}

	userID, commentID, err := getId(r)
	if err != nil {
		l.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	err = l.LikeManager.SetLikeToComment(r.Context(), commentID, userID)
	if err != nil {
		l.Responder.ErrorInternal(w, err, reqID)
		return
	}

	l.Responder.OutputJSON(w, "like is set on comment", reqID)
}

func (l LikeController) SetLikeToFile(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		l.Responder.LogError(myErr.ErrInvalidContext, "")
	}

	userID, fileID, err := getId(r)
	if err != nil {
		l.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	err = l.LikeManager.SetLikeToFile(r.Context(), fileID, userID)
	if err != nil {
		l.Responder.ErrorInternal(w, err, reqID)
		return
	}

	l.Responder.OutputJSON(w, "like is set on file", reqID)
}

func (l LikeController) DeleteLikeFromPost(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		l.Responder.LogError(myErr.ErrInvalidContext, "")
	}

	userID, postID, err := getId(r)
	if err != nil {
		l.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	err = l.LikeManager.DeleteLikeFromPost(r.Context(), postID, userID)
	if err != nil {
		l.Responder.ErrorInternal(w, err, reqID)
		return
	}

	l.Responder.OutputJSON(w, "like is unset from post", reqID)
}

func (l LikeController) DeleteLikeFromComment(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		l.Responder.LogError(myErr.ErrInvalidContext, "")
	}

	userID, commentID, err := getId(r)
	if err != nil {
		l.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	err = l.LikeManager.DeleteLikeFromComment(r.Context(), commentID, userID)
	if err != nil {
		l.Responder.ErrorInternal(w, err, reqID)
		return
	}

	l.Responder.OutputJSON(w, "like is unset from comment", reqID)
}

func (l LikeController) DeleteLikeFromFile(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		l.Responder.LogError(myErr.ErrInvalidContext, "")
	}

	userID, fileID, err := getId(r)
	if err != nil {
		l.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	err = l.LikeManager.SetLikeToPost(r.Context(), fileID, userID)
	if err != nil {
		l.Responder.ErrorInternal(w, err, reqID)
		return
	}

	l.Responder.OutputJSON(w, "like is unset from file", reqID)
}

func (l LikeController) GetLikesOnPost(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		l.Responder.LogError(myErr.ErrInvalidContext, "")
	}

	postID, err := getIdFromURL(r)
	if err != nil {
		l.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	likes, err := l.LikeManager.GetLikesOnPost(r.Context(), postID)
	if err != nil {
		if errors.Is(err, myErr.ErrWrongPost) {
			l.Responder.ErrorBadRequest(w, err, reqID)
			return
		}
		l.Responder.ErrorInternal(w, err, reqID)
		return
	}
	l.Responder.OutputJSON(w, likes, reqID)
}
