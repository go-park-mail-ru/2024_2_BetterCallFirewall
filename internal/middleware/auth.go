package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/auth"
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

var noAuthUrls = map[string]struct{}{
	"/api/v1/auth/register": {},
	"/api/v1/auth/login":    {},
}

func Auth(sm auth.SessionManager, next http.Handler) http.Handler {
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

		sessionCookie, err := r.Cookie("session_id")
		if errors.Is(err, http.ErrNoCookie) {
			log.Println(http.ErrNoCookie)
			return
		}
		sess, err := sm.Check(sessionCookie.Value)

		if _, ok := noAuthUrls[r.URL.Path]; ok {
			if err == nil {
				//TODO подумать над использованием /logout
				err := sm.Destroy(sess)
				if err != nil {
					log.Println(r.Context().Value("requestID"), err)
				}
				cookie := &http.Cookie{
					Name:     "session_id",
					Value:    sess.ID,
					Path:     "/",
					HttpOnly: true,
					Expires:  time.Now().AddDate(0, 0, -1),
				}
				http.SetCookie(w, cookie)
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
			log.Println(r.Context().Value("requestID"), err)
			return
		}

		if sess.CreatedAt <= time.Now().Add(-12*time.Hour).Unix() {
			sess, err = sm.Create(w, sess.UserID)
			if err != nil {
				log.Println(r.Context().Value("requestID"), err)
			}
		}

		ctx := models.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
