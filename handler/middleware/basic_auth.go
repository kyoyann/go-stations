package middleware

import (
	"net/http"
	"os"
)

func BasicAuth(h http.Handler) http.Handler {
	uid := os.Getenv("BASIC_AUTH_USER_ID")
	pw := os.Getenv("BASIC_AUTH_PASSWORD")
	fn := func(w http.ResponseWriter, r *http.Request) {
		clientID, clientSecret, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if clientID != uid || clientSecret != pw {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
