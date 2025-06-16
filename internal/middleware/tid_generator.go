package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	TIDHeader    = "X_TRACE_ID"
	SpanIDHeader = "X_SPAN_ID"
	MethodKey    = "METHOD"
)

// TIDMiddleware sets up tracing headers: TID, PID, CID (aka SPAN ID)
func TIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract from headers or create new ones
		tid := uuid.New().String()

		span := "ROOT"

		r.Header.Set(TIDHeader, tid)
		r.Header.Set(SpanIDHeader, span)

		w.Header().Set(TIDHeader, tid)
		w.Header().Set(SpanIDHeader, span)

		ctx := r.Context()
		ctx = withTraceContext(ctx, tid, span)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func withTraceContext(ctx context.Context, tid string, span string) context.Context {
	ctx = context.WithValue(ctx, TIDHeader, tid)
	ctx = context.WithValue(ctx, SpanIDHeader, span)
	return ctx
}
