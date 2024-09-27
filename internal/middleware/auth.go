package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/2024_2_BetterCallFirewall/internal/auth/models"
)

var noAuthUrls = map[string]struct{}{
	"/auth/register": {},
	"/auth/login":    {},
}

type SessionManager interface {
	Check(r *http.Request) (*models.Session, error)
}

func Auth(sm SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		sess, err := sm.Check(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(fmt.Errorf("not authorized: %w", err).Error()))
			log.Println(err)
			return
		}

		ctx := models.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
