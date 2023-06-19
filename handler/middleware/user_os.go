package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mileusna/useragent"
)

type contexKey string

const userOSKey contexKey = "UserOS"

func SetUserOS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), userOSKey, ua.OS)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetUserOS(ctx context.Context) (string, error) {
	v := ctx.Value(userOSKey)

	os, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("token not found")
	}

	return os, nil
}
