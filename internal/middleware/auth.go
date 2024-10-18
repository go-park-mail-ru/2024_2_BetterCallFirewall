package middleware

import (
	"fmt"
	"github.com/2024_2_BetterCallFirewall/internal/models"
	"log"
	"net/http"
	"time"
)

var noAuthUrls = map[string]struct{}{
	"/api/v1/auth/register": {},
	"/api/v1/auth/login":    {},
}

type SessionManager interface {
	Check(r *http.Request) (*models.Session, error)
	Create(w http.ResponseWriter, userID uint32) (*models.Session, error)
	Destroy(w http.ResponseWriter, r *http.Request) error
}

func Auth(sm SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "http://185.241.194.197:8000")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
			w.Header().Set("Content-Type", "application/json:charset=UTF-8")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			w.WriteHeader(http.StatusOK)
			return
		}

		sess, err := sm.Check(r)

		if _, ok := noAuthUrls[r.URL.Path]; ok {
			if err == nil {
				err := sm.Destroy(w, r.WithContext(models.ContextWithSession(r.Context(), sess)))
				if err != nil {
					log.Println(err)
				}
			}
			next.ServeHTTP(w, r)
			return
		}

		if err != nil {
			w.Header().Set("Content-Type", "application/json:charset=UTF-8")
			w.Header().Set("Access-Control-Allow-Origin", "http://185.241.194.197:8000")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(fmt.Errorf("not authorized: %w", err).Error()))
			log.Println(err)
			return
		}

		if sess.CreatedAt <= time.Now().Add(-12*time.Hour).Unix() {
			sess, err = sm.Create(w, sess.UserID)
			if err != nil {
				log.Println(err)
			}
		}

		ctx := models.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
