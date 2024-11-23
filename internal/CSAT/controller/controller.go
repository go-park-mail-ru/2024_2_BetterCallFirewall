package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

type Service interface {
	SaveMetrics(ctx context.Context, csat *models.CSAT, userID uint32) error
	GetMetrics(ctx context.Context, since, before time.Time) (*models.CSATResult, error)
	CheckExperience(userID uint32) bool
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestId string)

	LogError(err error, requestID string)
	ErrorBadRequest(w http.ResponseWriter, err error, requestID string)
	ErrorInternal(w http.ResponseWriter, err error, requestID string)
}

type Controller struct {
	service   Service
	responder Responder
}

func NewCSATController(service Service, responder Responder) *Controller {
	return &Controller{
		service:   service,
		responder: responder,
	}
}

func (cs *Controller) CheckExperience(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
	)

	if !ok {
		cs.responder.LogError(my_err.ErrInvalidContext, "")
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		cs.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	res := cs.service.CheckExperience(sess.UserID)
	cs.responder.OutputJSON(w, res, reqID)
}

func (cs *Controller) SaveMetrics(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
	)

	if !ok {
		cs.responder.LogError(my_err.ErrInvalidContext, "")
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		cs.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	csat := models.CSAT{}
	err = json.NewDecoder(r.Body).Decode(&csat)
	if err != nil {
		cs.responder.ErrorBadRequest(w, err, reqID)
	}

	err = cs.service.SaveMetrics(r.Context(), &csat, sess.UserID)
	if err != nil {
		cs.responder.ErrorInternal(w, err, reqID)
		return
	}

	cs.responder.OutputJSON(w, "metrics are saved", reqID)
}

func (cs *Controller) GetMetrics(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
	)

	if !ok {
		cs.responder.LogError(my_err.ErrInvalidContext, "")
	}
	since, before := getTimeFromQuery(r)
	metrics, err := cs.service.GetMetrics(r.Context(), since, before)
	if err != nil {
		cs.responder.ErrorInternal(w, err, reqID)
		return
	}

	cs.responder.OutputJSON(w, metrics, reqID)

}

func getTimeFromQuery(r *http.Request) (time.Time, time.Time) {
	since, before := r.URL.Query().Get("since"), r.URL.Query().Get("before")
	sinceTime, err := time.Parse(time.RFC3339, since)
	if err != nil {
		sinceTime = time.Time{}
	}
	beforeTime, err := time.Parse(time.RFC3339, before)
	if err != nil {
		beforeTime = time.Time{}
	}
	return sinceTime, beforeTime
}
