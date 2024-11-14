package controller

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/models"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

const defaultAva = 2

type responder interface {
	OutputJSON(w http.ResponseWriter, data any, requestID string)
	OutputNoMoreContentJSON(w http.ResponseWriter, requestId string)

	ErrorBadRequest(w http.ResponseWriter, err error, requestID string)
	ErrorInternal(w http.ResponseWriter, err error, requestID string)
	LogError(err error, requestID string)
}

type communityService interface {
	Get(ctx context.Context, lastID uint32) ([]*models.CommunityCard, error)
	GetOne(ctx context.Context, id uint32) (*models.Community, error)
	Update(ctx context.Context, id uint32, community *models.Community) error
	Delete(ctx context.Context, id uint32) error
	Create(ctx context.Context, community *models.Community, authorID uint32) error
	CheckAccess(ctx context.Context, communityID, userID uint32) bool
}

type Controller struct {
	responder responder
	service   communityService
}

func NewController(responder responder, service communityService) *Controller {
	return &Controller{
		responder: responder,
		service:   service,
	}
}

func (c *Controller) GetOne(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		c.responder.LogError(my_err.ErrInvalidContext, "")
	}

	id, err := getIDFromQuery(r)
	if err != nil {
		c.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	community, err := c.service.GetOne(r.Context(), id)
	if err != nil {
		c.responder.ErrorInternal(w, err, reqID)
		return
	}

	c.responder.OutputJSON(w, community, reqID)
}

func (c *Controller) GetAll(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		lastID    = r.URL.Query().Get("id")
		intLastID uint64
		err       error
	)

	if !ok {
		c.responder.LogError(my_err.ErrInvalidContext, "")
	}
	if lastID == "" {
		intLastID = math.MaxInt32
	} else {
		intLastID, err = strconv.ParseUint(lastID, 10, 32)
		if err != nil {
			c.responder.ErrorBadRequest(w, my_err.ErrInvalidQuery, reqID)
			return
		}
	}

	res, err := c.service.Get(r.Context(), uint32(intLastID))
	if err != nil {
		c.responder.ErrorInternal(w, err, reqID)
		return
	}

	c.responder.OutputJSON(w, res, reqID)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		c.responder.LogError(my_err.ErrInvalidContext, "")
	}

	id, err := getIDFromQuery(r)
	if err != nil {
		c.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		c.responder.ErrorInternal(w, err, reqID)
		return
	}

	newCommunity, err := c.getCommunityFromBody(r)
	if err != nil {
		c.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	if !c.service.CheckAccess(r.Context(), id, sess.UserID) {
		c.responder.ErrorBadRequest(w, my_err.ErrAccessDenied, reqID)
		return
	}

	err = c.service.Update(r.Context(), id, &newCommunity)
	if err != nil {
		c.responder.ErrorInternal(w, err, reqID)
	}

	c.responder.OutputJSON(w, newCommunity, reqID)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		c.responder.LogError(my_err.ErrInvalidContext, "")
	}

	id, err := getIDFromQuery(r)
	if err != nil {
		c.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		c.responder.ErrorInternal(w, err, reqID)
		return
	}

	if !c.service.CheckAccess(r.Context(), id, sess.UserID) {
		c.responder.ErrorBadRequest(w, my_err.ErrAccessDenied, reqID)
		return
	}

	err = c.service.Delete(r.Context(), id)
	if err != nil {
		c.responder.ErrorInternal(w, err, reqID)
		return
	}

	c.responder.OutputJSON(w, id, reqID)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		c.responder.LogError(my_err.ErrInvalidContext, "")
	}

	sess, err := models.SessionFromContext(r.Context())
	if err != nil {
		c.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	newCommunity, err := c.getCommunityFromBody(r)
	if err != nil {
		c.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	err = c.service.Create(r.Context(), &newCommunity, sess.UserID)
	if err != nil {
		c.responder.ErrorInternal(w, err, reqID)
		return
	}

	c.responder.OutputJSON(w, newCommunity.ID, reqID)
}

func (c *Controller) getCommunityFromBody(r *http.Request) (models.Community, error) {
	var res models.Community

	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return res, err
	}

	return models.Community{}, nil
}

func getIDFromQuery(r *http.Request) (uint32, error) {
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