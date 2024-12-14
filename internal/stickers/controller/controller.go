package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/internal/stickers"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

const imagePrefix = "/image/"

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type Responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestID string)
	OutputNoMoreContentJSON(w http.ResponseWriter, requestID string)

	ErrorBadRequest(w http.ResponseWriter, err error, requestID string)
	ErrorInternal(w http.ResponseWriter, err error, requestID string)
	LogError(err error, requestID string)
}

type StickersHandlerImplementation struct {
	StickersManager stickers.Usecase
	Responder       Responder
}

func NewStickerController(manager stickers.Usecase, responder Responder) *StickersHandlerImplementation {
	return &StickersHandlerImplementation{
		StickersManager: manager,
		Responder:       responder,
	}
}

func (s StickersHandlerImplementation) AddNewSticker(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)

	if !ok {
		s.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	filePath := ""
	if err := json.NewDecoder(r.Body).Decode(&filePath); err != nil {
		s.Responder.ErrorBadRequest(w, my_err.ErrNoFile, reqID)
		fmt.Println(err)
		return
	}

	if !validate(filePath) {
		s.Responder.ErrorBadRequest(w, my_err.ErrNoImage, reqID)
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		s.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	err = s.StickersManager.AddNewSticker(r.Context(), filePath, sess.UserID)
	if err != nil {
		s.Responder.ErrorInternal(w, err, reqID)
		return
	}

	s.Responder.OutputJSON(w, "New sticker is added", reqID)
}

func (s StickersHandlerImplementation) GetAllStickers(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)

	if !ok {
		s.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	res, err := s.StickersManager.GetAllStickers(r.Context())
	if err != nil {
		if errors.Is(err, my_err.ErrNoStickers) {
			s.Responder.OutputNoMoreContentJSON(w, reqID)
			return
		}
		s.Responder.ErrorInternal(w, err, reqID)
		return
	}

	s.Responder.OutputJSON(w, res, reqID)
}

func (s StickersHandlerImplementation) GetMineStickers(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)

	if !ok {
		s.Responder.LogError(my_err.ErrInvalidContext, "")
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		s.Responder.ErrorBadRequest(w, err, reqID)
		return
	}

	res, err := s.StickersManager.GetMineStickers(r.Context(), sess.UserID)
	if err != nil {
		if errors.Is(err, my_err.ErrNoStickers) {
			s.Responder.OutputNoMoreContentJSON(w, reqID)
			return
		}
		s.Responder.ErrorInternal(w, err, reqID)
		return
	}

	s.Responder.OutputJSON(w, res, reqID)
}

func validate(filepath string) bool {
	return len(filepath) < 100 && strings.HasPrefix(filepath, imagePrefix)
}