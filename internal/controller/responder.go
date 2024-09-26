package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, data any)

	ErrorWrongMethod(w http.ResponseWriter, err error)
	ErrorUnAuthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
}

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type Respond struct {
	logger *log.Logger
}

func NewResponder(logger *log.Logger) *Respond {
	return &Respond{logger: logger}
}

func (r *Respond) OutputJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&Response{Success: true, Data: data})
	if err != nil {
		r.logger.Println(err)
	}
	r.logger.Println(data)
}

func (r *Respond) ErrorWrongMethod(w http.ResponseWriter, err error) {
	r.logger.Println(err)
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	w.WriteHeader(http.StatusMethodNotAllowed)
	errJ := json.NewEncoder(w).Encode(&Response{Success: false, Data: err.Error(), Message: "method not allowed"})
	if errJ != nil {
		r.logger.Println(err)
	}
}

func (r *Respond) ErrorUnAuthorized(w http.ResponseWriter, err error) {
	r.logger.Println(err)
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	errJ := json.NewEncoder(w).Encode(&Response{Success: false, Data: err.Error(), Message: "unauthorized"})
	if errJ != nil {
		r.logger.Println(err)
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.logger.Println(err)
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	errJ := json.NewEncoder(w).Encode(&Response{Success: false, Data: err.Error(), Message: "bad request"})
	if errJ != nil {
		r.logger.Println(err)
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	r.logger.Println(err)
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	if errors.Is(err, context.Canceled) {
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	errJ := json.NewEncoder(w).Encode(&Response{Success: false, Data: err.Error(), Message: "internal server error"})
	if errJ != nil {
		r.logger.Println(err)
	}
}

func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	r.logger.Println(err)
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	w.WriteHeader(http.StatusForbidden)
	errJ := json.NewEncoder(w).Encode(&Response{Success: false, Data: err.Error(), Message: "error forbidden"})
	if errJ != nil {
		r.logger.Println(err)
	}
}
