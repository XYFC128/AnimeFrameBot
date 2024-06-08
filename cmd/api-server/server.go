package main

import (
	"log"
	"net/http"
	"time"
)

func NewServer(imagepath string) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, imagepath)
	var handler http.Handler = loggingMiddleWare(mux)
	return handler
}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL, wrapped.statusCode, time.Since(start))
	})
}
