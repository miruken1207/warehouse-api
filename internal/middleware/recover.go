package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/miruken1207/warehouse-api/internal/handler"
)

func Recover(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered",
					"error", rec,
					"method", r.Method,
					"route", r.URL.Path,
					"stack", string(debug.Stack()),
				)
				handler.WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
