package internalhttp

import (
	"log"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Вызов handler
		next.ServeHTTP(w, r)

		latency := time.Since(start)

		log.Printf(
			"%s [%s] %s %s %s %dms \"%s\"\n",
			r.RemoteAddr,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			latency.Milliseconds(),
			r.UserAgent(),
		)
	})
}
