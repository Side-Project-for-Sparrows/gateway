package observability

import (
	"context"
	"net/http"

	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/google/uuid"
)

const (
	TIDHeader    = "X_TRACE_ID"
	SpanIDHeader = "X_SPAN_ID"
	MethodKey    = "METHOD"
)

// TIDMiddleware sets up tracing headers: TID, PID, CID (aka SPAN ID)
func TIDMiddleware() middlewaretype.Middleware {
	return func(w http.ResponseWriter, r *http.Request) error {
		tid := uuid.New().String()
		span := "ROOT"

		r.Header.Set(TIDHeader, tid)
		r.Header.Set(SpanIDHeader, span)
		w.Header().Set(TIDHeader, tid)
		w.Header().Set(SpanIDHeader, span)

		ctx := context.WithValue(r.Context(), TIDHeader, tid)
		ctx = context.WithValue(ctx, SpanIDHeader, span)
		*r = *r.WithContext(ctx)

		return nil
	}
}
