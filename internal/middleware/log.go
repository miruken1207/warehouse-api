package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
)

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *logResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func Log(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &logResponseWriter{w, 200}
		next.ServeHTTP(lrw, r)
		logger.Info("HTTP Request", "method", r.Method, "route", r.URL.Path, "status code", strconv.Itoa(lrw.statusCode))
	})
}
