package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		startTime := time.Now()
		log.Printf("Request: %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log the response time
		duration := time.Since(startTime)
		log.Printf("Response: %s %s %s %s", r.Method, r.URL.Path, r.RemoteAddr, duration)
	})
}
