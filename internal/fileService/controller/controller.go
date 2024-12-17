package controller

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"

	"github.com/2024_2_BetterCallFirewall/internal/middleware"
	"github.com/2024_2_BetterCallFirewall/pkg/my_err"
)

var fileFormat = map[string]struct{}{
	"jpeg": {},
	"jpg":  {},
	"png":  {},
	"webp": {},
	"gif":  {},
}

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type fileService interface {
	Upload(ctx context.Context, name string) ([]byte, error)
	Download(ctx context.Context, file multipart.File, format string) (string, error)
	DownloadNonImage(ctx context.Context, file multipart.File, format, realName string) (string, error)
	UploadNonImage(ctx context.Context, name string) ([]byte, error)
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

func sanitize(input string) string {
	sanitizer := bluemonday.UGCPolicy()
	cleaned := sanitizer.Sanitize(input)
	return cleaned
}

func (fc *FileController) UploadNonImage(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value("requestID").(string)
		vars      = mux.Vars(r)
		name      = sanitize(vars["name"])
	)

	if !ok {
		fc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	if name == "" {
		fc.responder.ErrorBadRequest(w, errors.New("name is empty"), reqID)
		return
	}

	res, err := fc.fileService.UploadNonImage(r.Context(), name)
	if err != nil {
		fc.responder.ErrorBadRequest(w, fmt.Errorf("%w: %w", err, my_err.ErrWrongFile), reqID)
		return
	}

	fc.responder.OutputBytes(w, res, reqID)
}

func (fc *FileController) Upload(w http.ResponseWriter, r *http.Request) {
	var (
		reqID, ok = r.Context().Value(middleware.RequestKey).(string)
		vars      = mux.Vars(r)
		name      = sanitize(vars["name"])
	)

	if !ok {
		fc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	if name == "" {
		fc.responder.ErrorBadRequest(w, errors.New("name is empty"), reqID)
		return
	}

	res, err := fc.fileService.Upload(r.Context(), name)
	if err != nil {
		fc.responder.ErrorBadRequest(w, fmt.Errorf("%w: %w", err, my_err.ErrWrongFile), reqID)
		return
	}

	fc.responder.OutputBytes(w, res, reqID)
}

func (fc *FileController) Download(w http.ResponseWriter, r *http.Request) {
	reqID, ok := r.Context().Value(middleware.RequestKey).(string)
	if !ok {
		fc.responder.LogError(my_err.ErrInvalidContext, "")
	}

	err := r.ParseMultipartForm(10 * (10 << 20)) // 100Mbyte
	if err != nil {
		fc.responder.ErrorBadRequest(w, my_err.ErrToLargeFile, reqID)
		return
	}
	defer func() {
		err = r.MultipartForm.RemoveAll()
		if err != nil {
			fc.responder.LogError(err, reqID)
		}
	}()

	file, header, err := r.FormFile("file")
	if err != nil {
		fc.responder.ErrorBadRequest(w, err, reqID)
		return
	}

	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			fc.responder.LogError(err, reqID)
		}
	}(file)

	formats := strings.Split(header.Header.Get("Content-Type"), "/")
	if len(formats) != 2 {
		fc.responder.ErrorBadRequest(w, my_err.ErrWrongFiletype, reqID)
		return
	}
	var url string
	format := formats[1]

	if _, ok := fileFormat[format]; ok {
		url, err = fc.fileService.Download(r.Context(), file, format)
	} else {
		name := header.Filename
		if len(name+format) > 55 {
			fc.responder.ErrorBadRequest(w, errors.New("file name is too big"), reqID)
			return
		}
		url, err = fc.fileService.DownloadNonImage(r.Context(), file, format, name)
	}

	if err != nil {
		fc.responder.ErrorBadRequest(w, err, reqID)
		return
	}
	fc.responder.OutputJSON(w, url, reqID)
}
