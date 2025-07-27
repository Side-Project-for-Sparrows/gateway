package observability

import (
	"net/http"

	"github.com/Side-Project-for-Sparrows/gateway/internal/middleware/middlewaretype"
	"github.com/google/uuid"
)

const (
	TIDHeader    = "X_TRACE_ID"
	SpanIDHeader = "X_SPAN_ID"
	MethodKey    = "METHOD"
)

func TIDMiddleware() middlewaretype.Middleware {
	return func(input middlewaretype.MiddlewareInput) (*middlewaretype.HeaderPatch, error) {
		tid := uuid.New().String()
		span := "ROOT"

		return &middlewaretype.HeaderPatch{
			RequestAdd: http.Header{
				TIDHeader:    []string{tid},
				SpanIDHeader: []string{span},
			},
			ResponseAdd: http.Header{
				TIDHeader:    []string{tid},
				SpanIDHeader: []string{span},
			},
			ContextAdd: map[any]any{
				TIDHeader:    tid,
				SpanIDHeader: span,
			},
		}, nil
	}
}
