package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/2024_2_BetterCallFirewall/internal/models"
)

var noAuthUrls = map[string]struct{}{
	"/api/v1/auth/register": {},
	"/api/v1/auth/login":    {},
}

type SessionManager interface {
	Check(string) (*models.Session, error)
	Create(userID uint32) (*models.Session, error)
	Destroy(sess *models.Session) error
}

func Auth(sm SessionManager, next http.Handler) http.Handler {
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
		}

		if _, ok := noAuthUrls[r.URL.Path]; ok {
			logout(w, r, sm)
			next.ServeHTTP(w, r)
			return
		}

		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			unauthorized(w, r, err)
			return
		}

		sess, err := sm.Check(sessionCookie.Value)
		if err != nil {
			unauthorized(w, r, err)
			return
		}

		if sess.CreatedAt <= time.Now().Add(-time.Hour).Unix() {
			if err := sm.Destroy(sess); err != nil {
				log.Println(r.Context().Value("requestID"), err)
				internalErr(w)
				return
			}

			sess, err = sm.Create(sess.UserID)
			if err != nil {
				log.Println(r.Context().Value("requestID"), err)
				internalErr(w)
				return
			}

			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    sess.ID,
				Path:     "/",
				Domain:   "vilka.online",
				HttpOnly: true,
				Expires:  time.Now().AddDate(0, 0, 1),
			}
			http.SetCookie(w, cookie)
		}

		ctx := models.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func logout(w http.ResponseWriter, r *http.Request, sm SessionManager) {
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		return
	}

	sess, err := sm.Check(sessionCookie.Value)
	if err != nil {
		log.Println(err)
		return
	}

	err = sm.Destroy(sess)
	if err != nil {
		log.Println(err)
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sess.ID,
		Path:     "/",
		Domain:   "vilka.online",
		HttpOnly: true,
		Expires:  time.Now().AddDate(0, 0, -1),
	}

	http.SetCookie(w, cookie)
}

func unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "http://vilka.online")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusUnauthorized)

	_, _ = w.Write([]byte("not authorized"))

	log.Println(r.Context().Value("requestID"), err)
}

func internalErr(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json:charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "http://vilka.online")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(http.StatusInternalServerError)

	_, _ = w.Write([]byte("internal server error"))
}
