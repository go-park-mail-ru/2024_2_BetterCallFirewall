package middleware

import (
	"log"
	"net/http"
	"time"
)

func AccessLog(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Printf("New request:\n \tMethod: %v\n\tRemote addr: %v\n\tURL: %v\n\tTime: %v", r.Method, r.RemoteAddr, r.URL.String(), time.Since(start))
	})
}
