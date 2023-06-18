package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type accessLog struct {
	Timestamp time.Time
	Latency   int64
	Path      string
	OS        string
}

func AccessLogger(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		h.ServeHTTP(w, r)

		al := accessLog{
			Timestamp: start,
			Latency:   int64(time.Since(start).Milliseconds()),
			Path:      r.URL.Path,
			OS:        r.Context().Value(UserAgentKey).(string),
		}
		bytes, err := json.Marshal(al)
		if err != nil {
			log.Println(err)
		}
		log.Println(string(bytes))
	}
	return http.HandlerFunc(fn)
}
