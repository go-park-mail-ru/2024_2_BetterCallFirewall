package router

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

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

func writeHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "http://185.241.194.197:8000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func (r *Respond) OutputJSON(w http.ResponseWriter, data any) {
	writeHeaders(w)
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(&Response{Success: true, Data: data})
	if err != nil {
		r.logger.Println(err)
	}
	dataJ, err := json.Marshal(&Response{Success: true, Data: data})
	if err != nil {
		r.logger.Println(err)
	}
	r.logger.Println(string(dataJ))
}

func (r *Respond) ErrorWrongMethod(w http.ResponseWriter, err error) {
	r.logger.Println(err)
	writeHeaders(w)
	w.WriteHeader(http.StatusMethodNotAllowed)

	errJ := json.NewEncoder(w).Encode(&Response{Success: false, Data: err.Error(), Message: "method not allowed"})
	if errJ != nil {
		r.logger.Println(errJ)
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.logger.Println(err)
	writeHeaders(w)
	w.WriteHeader(http.StatusBadRequest)

	errJ := json.NewEncoder(w).Encode(&Response{Success: false, Data: err.Error(), Message: "bad request"})
	if errJ != nil {
		r.logger.Println(errJ)
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	r.logger.Println(err)
	writeHeaders(w)
	if errors.Is(err, context.Canceled) {
		return
	}

	w.WriteHeader(http.StatusInternalServerError)

	errJ := json.NewEncoder(w).Encode(&Response{Success: false, Data: err.Error(), Message: "internal server error"})
	if errJ != nil {
		r.logger.Println(errJ)
	}
}
