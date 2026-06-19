package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type ctxKey int

const batchIDKey ctxKey = 0

func setBatchID(ctx context.Context, id string) {
	if p, ok := ctx.Value(batchIDKey).(*string); ok {
		*p = id
	}
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		var batchID string
		ctx := context.WithValue(r.Context(), batchIDKey, &batchID)

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(rec, r.WithContext(ctx))
		duration := time.Since(start)

		attrs := []any{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Int64("duration_ms", duration.Milliseconds()),
		}
		if batchID != "" {
			attrs = append(attrs, slog.String("batch_id", batchID))
		}
		logger.Info("requête HTTP", attrs...)
	})
}
