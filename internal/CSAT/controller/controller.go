package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
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
}

type Controller struct {
	service   Service
	responder Responder
}
