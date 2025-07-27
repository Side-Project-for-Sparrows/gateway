package handler

import "net/http"

// dummyHandler는 실제 HTTP 호출 없이 즉시 종료함
func DummyHandler(w http.ResponseWriter, r *http.Request) {
	// no-op: do nothing
}
