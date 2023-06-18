package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		defer func() {
			//nilが返ってきた場合はパニックが起こっていない
			if err := recover(); err != nil {
				jsonBody, _ := json.Marshal(map[string]string{
					"error": fmt.Sprintf("%v", err),
				})
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
