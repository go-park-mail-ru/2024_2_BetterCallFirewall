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
	Data    any    `json:"data"`
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

// TODO add requestID
func (r *Respond) OutputJSON(w http.ResponseWriter, data any) {
	writeHeaders(w)
	w.WriteHeader(http.StatusOK)

	r.logger.Info("success request")
	if err := json.NewEncoder(w).Encode(&Response{Success: true, Data: data}); err != nil {
		r.logger.Error(err)
	}
}

func (r *Respond) OutputNoMoreContentJSON(w http.ResponseWriter) {
	writeHeaders(w)
	w.WriteHeader(http.StatusNoContent)

	r.logger.Info("success request")
}

func (r *Respond) ErrorWrongMethod(w http.ResponseWriter, err error) {
	r.logger.Warn(err)
	writeHeaders(w)
	w.WriteHeader(http.StatusMethodNotAllowed)

	if err := json.NewEncoder(w).Encode(&Response{Success: false, Data: fullUnwrap(err).Error(), Message: "method not allowed"}); err != nil {
		r.logger.Error(err)
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.logger.Warn(err)
	writeHeaders(w)
	w.WriteHeader(http.StatusBadRequest)

	if err := json.NewEncoder(w).Encode(&Response{Success: false, Data: fullUnwrap(err).Error(), Message: "bad request"}); err != nil {
		r.logger.Error(err)
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	r.logger.Error(err)
	writeHeaders(w)
	if errors.Is(err, context.Canceled) {
		return
	}

	w.WriteHeader(http.StatusInternalServerError)

	if err := json.NewEncoder(w).Encode(&Response{Success: false, Data: myErr.ErrInternal, Message: "internal server error"}); err != nil {
		r.logger.Error(err)
	}
}
