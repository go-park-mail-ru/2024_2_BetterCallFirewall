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

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}

	return h.Hijack()
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func HttpMetricsMiddleware(metr *metrics.HttpMetrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		respWithCode := NewResponseWriter(w)
		next.ServeHTTP(respWithCode, r)
		statusCode := respWithCode.statusCode
		path := r.URL.Path
		method := r.Method
		if statusCode != http.StatusOK && statusCode != http.StatusNoContent {
			metr.IncErrors(path, strconv.Itoa(statusCode), method)
		}
		metr.IncHits(path, strconv.Itoa(statusCode), method)
		metr.ObserveTiming(path, strconv.Itoa(statusCode), method, time.Since(start).Seconds())
	})
}
