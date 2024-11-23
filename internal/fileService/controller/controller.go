package controller

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

var fileFormat = map[string]struct{}{
	"image/jpeg": {},
	"image/jpg":  {},
	"image/png":  {},
	"image/webp": {},
}

type fileService interface {
	Upload(ctx context.Context, name string) ([]byte, error)
	Download(ctx context.Context, file multipart.File) (string, error)
}

type responder interface {
	OutputBytes(w http.ResponseWriter, data []byte, requestID string)
	OutputJSON(w http.ResponseWriter, data any, requestId string)

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
		fc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	res, err := fc.fileService.Upload(r.Context(), name)
	if err != nil {
		fc.responder.ErrorBadRequest(w, fmt.Errorf("%w: %w", err, my_err.ErrWrongFile), reqID)
		return
	}

	fc.responder.OutputBytes(w, res, reqID)
}

func (fc *FileController) Download(w http.ResponseWriter, r *http.Request) {
	log.Println("GET REQUEST FOR CREATE FILE")

	reqID, ok := r.Context().Value("requestID").(string)
	if !ok {
		fc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	err := r.ParseMultipartForm(10 << 20) // 10Mbyte
	defer r.MultipartForm.RemoveAll()
	if err != nil {
		fc.responder.ErrorBadRequest(w, my_err.ErrToLargeFile, reqID)
		return
	}

	log.Println("PARSE FILE")

	file, header, err := r.FormFile("file")
	if err != nil {
		file = nil
	} else {
		format := header.Header.Get("Content-Type")
		if _, ok := fileFormat[format]; !ok {
			fc.responder.ErrorBadRequest(w, my_err.ErrWrongFiletype, reqID)
			return
		}
	}
	defer file.Close()
	log.Println("DOWNLOAD FILE")

	url, err := fc.fileService.Download(r.Context(), file)
	if err != nil {
		fc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	fc.responder.OutputJSON(w, url, reqID)
}
