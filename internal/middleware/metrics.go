package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/2024_2_BetterCallFirewall/internal/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
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
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		if statusCode != http.StatusOK && statusCode != http.StatusNoContent {
			metr.IncErrors(path, strconv.Itoa(statusCode))
		}
		metr.IncHits(path, strconv.Itoa(statusCode))
		metr.ObserveTiming(path, strconv.Itoa(statusCode), time.Since(start).Seconds())
	})
}
