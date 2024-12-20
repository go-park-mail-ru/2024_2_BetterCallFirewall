package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
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

const (
	charset   = "charset=utf"
	txt       = "txt"
	plain     = "plain"
	maxMemory = 100 * 1024 * 1024
)

//go:generate mockgen -destination=mock.go -source=$GOFILE -package=${GOPACKAGE}
type fileService interface {
	Upload(ctx context.Context, name string) ([]byte, error)
	Download(ctx context.Context, file io.Reader, format string) (string, error)
	DownloadNonImage(ctx context.Context, file io.Reader, format, realName string) (string, error)
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

func getFormat(buf []byte) string {
	formats := http.DetectContentType(buf)
	format := strings.Split(formats, "/")[1]

	if strings.Contains(format, charset) {
		format = strings.Split(format, ";")[0]
	}

	if format == plain {
		format = txt
	}

	return format
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

	err := r.ParseMultipartForm(maxMemory) // 100Mbyte
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

	buf := bytes.NewBuffer(make([]byte, 20))
	n, err := file.Read(buf.Bytes())
	if err != nil {
		fc.responder.ErrorBadRequest(w, err, reqID)
		return
	}
	format := getFormat(buf.Bytes()[:n])

	_, err = io.Copy(buf, file)
	if err != nil {
		fc.responder.ErrorBadRequest(w, err, reqID)
		return
	}
	var url string

	if _, ok := fileFormat[format]; ok {
		url, err = fc.fileService.Download(r.Context(), buf, format)
	} else {
		name := header.Filename
		if len([]rune(name+format)) > 55 {
			fc.responder.ErrorBadRequest(w, errors.New("file name is too big"), reqID)
			return
		}
		url, err = fc.fileService.DownloadNonImage(r.Context(), buf, format, name)
	}

	if err != nil {
		fc.responder.ErrorBadRequest(w, err, reqID)
		return
	}
	fc.responder.OutputJSON(w, url, reqID)
}
