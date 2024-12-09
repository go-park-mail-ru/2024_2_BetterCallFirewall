package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type requestID string

var requestKey requestID = "requestID"

func AccessLog(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestKey, id)
		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))
		logger.Infof("New request:%s\n \tMethod: %v\n\tRemote addr: %v\n\tURL: %v\n\tTime: %v", id, r.Method, r.RemoteAddr, r.URL.String(), time.Since(start))
	})
}
