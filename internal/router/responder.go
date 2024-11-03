package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

func fullUnwrap(err error) error {
	var last error

	for err != nil {
		last = err
		err = errors.Unwrap(err)
	}

	return last
}

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Respond struct {
	logger *log.Logger
}

func NewResponder(logger *log.Logger) *Respond {
	return &Respond{logger: logger}
}

func writeHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "http://185.241.194.197:8000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func (r *Respond) OutputJSON(w http.ResponseWriter, data any, requestID string) {
	writeHeaders(w)
	w.WriteHeader(http.StatusOK)

	r.logger.Infof("req: %s: success request", requestID)
	if err := json.NewEncoder(w).Encode(&Response{Success: true, Data: data}); err != nil {
		r.logger.Error(err)
	}
}

func (r *Respond) OutputNoMoreContentJSON(w http.ResponseWriter, requestID string) {
	writeHeaders(w)
	w.WriteHeader(http.StatusNoContent)

	r.logger.Infof("req: %s: success request", requestID)
}

func (r *Respond) OutputBytes(w http.ResponseWriter, data []byte, requestID string) {
	w.Header().Set("Content-Type", "image/avif,image/webp")
	w.Header().Set("Access-Control-Allow-Origin", "http://185.241.194.197:8000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(data)

	r.logger.Infof("req: %s: success request", requestID)
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error, requestID string) {
	r.logger.Warnf("req: %s: %v", requestID, err)
	writeHeaders(w)
	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(&Response{Success: false, Message: fullUnwrap(err).Error()}); err != nil {
		r.logger.Errorf("req: %s: %v", requestID, err)
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error, requestID string) {
	r.logger.Errorf("req: %s: %v", requestID, err)
	writeHeaders(w)
	if errors.Is(err, context.Canceled) {
		return
	}

	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(&Response{Success: false, Message: myErr.ErrInternal.Error()}); err != nil {
		r.logger.Errorf("req: %s: %v", requestID, err)
	}
}

func (r *Respond) LogError(err error, requestID string) {
	r.logger.Errorf("req: %s: %v", requestID, err)
}
