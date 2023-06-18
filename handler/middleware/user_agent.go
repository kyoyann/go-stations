package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type contexKey string

const UserAgentKey contexKey = "UserAgent"

func SetUserAgent(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), UserAgentKey, ua.OS)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
