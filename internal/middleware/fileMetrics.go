package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/metrics"
)

type fileResponseWriter struct {
	http.ResponseWriter
	statusCode int
	file       []byte
}

func NewFileResponseWriter(w http.ResponseWriter) *fileResponseWriter {
	return &fileResponseWriter{w, http.StatusOK, make([]byte, 0)}
}

func (rw *fileResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}

	return h.Hijack()
}

func (rw *fileResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *fileResponseWriter) Write(b []byte) (int, error) {
	rw.file = append(rw.file, b...)
	return rw.ResponseWriter.Write(b)
}

func FileMetricsMiddleware(metr *metrics.FileMetrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		respWithCode := NewFileResponseWriter(w)
		next.ServeHTTP(respWithCode, r)
		statusCode := respWithCode.statusCode
		path := r.URL.Path
		method := r.Method
		var (
			err    error
			format string
			size   int64
		)
		if r.Method == http.MethodPost {
			format, size, err = getFormatAndSize(r)
		} else if r.Method == http.MethodGet {
			file := respWithCode.file
			format = http.DetectContentType(file[:512])
			size = int64(len(file))
		}
		if err != nil {
			format = "error"
			size = 0
		}
		if statusCode != http.StatusOK && statusCode != http.StatusNoContent {
			metr.IncErrors(path, strconv.Itoa(statusCode), method, format, size)
		}
		metr.IncHits(path, strconv.Itoa(statusCode), method, format, size)
		metr.ObserveTiming(path, strconv.Itoa(statusCode), method, format, size, time.Since(start).Seconds())
	})
}

func getFormatAndSize(r *http.Request) (string, int64, error) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", 0, err
	}
	defer file.Close()
	size := fileHeader.Size
	format := fileHeader.Header.Get("Content-Type")
	return format, size, nil

}
