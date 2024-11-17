package middleware

import (
	"net/http"
)

func Preflite(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "http://vilka.online")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
			w.Header().Set("Content-Type", "application/json:charset=UTF-8")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			w.WriteHeader(http.StatusOK)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
