package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/metrics"
)

func FileMetricsMiddleware(metr *metrics.FileMetrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		respWithCode := NewResponseWriter(w)
		next.ServeHTTP(respWithCode, r)
		statusCode := respWithCode.statusCode
		path := r.URL.Path
		method := r.Method
		format, size, err := getFormatAndSize(r)
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
