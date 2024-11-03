package controller

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/myErr"
)

type fileService interface {
	Upload(ctx context.Context, name string) ([]byte, error)
}

type responder interface {
	OutputBytes(w http.ResponseWriter, data []byte, requestID string)
	LogError(err error, requestID string)
	ErrorBadRequest(w http.ResponseWriter, err error, requestID string)
}

type FileController struct {
	fileService fileService
	responder   responder
}

func NewFileController(fileService fileService, responder responder) *FileController {
	return &FileController{
		fileService: fileService,
		responder:   responder,
	}
}

func (fc *FileController) Upload(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		vars      = mux.Vars(r)
		name      = vars["name"]
	)

	if !ok {
		fc.responder.LogError(myErr.ErrInvalidContext, "")
	}

	res, err := fc.fileService.Upload(r.Context(), name)
	if err != nil {
		fc.responder.ErrorBadRequest(w, myErr.ErrWrongFile, reqID)
		return
	}

	fc.responder.OutputBytes(w, res, reqID)
}
