package middlewaretype

import (
	"net/http"
)

// 미들웨어 단위 (개별 미들웨어 함수 타입 정의)
type Middleware func(http.ResponseWriter, *http.Request) error

func (m Middleware) AdaptWithNext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m == nil {
			next.ServeHTTP(w, r)
			return
		}
		if err := m(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
