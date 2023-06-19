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
		os, _ := GetUserOS(r.Context())
		al := accessLog{
			Timestamp: start,
			Latency:   int64(time.Since(start).Milliseconds()),
			Path:      r.URL.Path,
			OS:        os,
		}
		bytes, err := json.Marshal(al)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(bytes))
	}
	return http.HandlerFunc(fn)
}
