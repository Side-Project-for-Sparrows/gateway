package middlewaretype

import (
	"log"
	"net/http"
	"time"
)

// 미들웨어 단위 (개별 미들웨어 함수 타입 정의)
// type Middleware func(http.ResponseWriter, *http.Request) error
type Middleware func(input MiddlewareInput) (*HeaderPatch, error)

type HeaderPatch struct {
	RequestAdd         http.Header
	RequestDelete      []string
	ResponseAdd        http.Header
	ResponseDelete     []string
	ContextAdd         map[any]any
	ResponseStatusCode int
	ResponseBody       []byte
}

func mergePatches(patches []*HeaderPatch) *HeaderPatch {
	merged := &HeaderPatch{
		RequestAdd:  http.Header{},
		ResponseAdd: http.Header{},
		ContextAdd:  map[any]any{},
	}

	for _, p := range patches {
		for k, vals := range p.RequestAdd {
			for _, v := range vals {
				merged.RequestAdd.Add(k, v)
			}
		}
		for k, vals := range p.ResponseAdd {
			for _, v := range vals {
				merged.ResponseAdd.Add(k, v)
			}
		}
		for k, v := range p.ContextAdd {
			merged.ContextAdd[k] = v
		}
		merged.RequestDelete = append(merged.RequestDelete, p.RequestDelete...)
		merged.ResponseDelete = append(merged.ResponseDelete, p.ResponseDelete...)
	}
	return merged
}

type ReadonlyContext interface {
	Value(key any) any
	Done() <-chan struct{}
	Err() error
	Deadline() (deadline time.Time, ok bool)
}

type MiddlewareInput struct {
	ctx        ReadonlyContext
	method     string
	path       string
	headers    http.Header
	remoteAddr string
}

func (mi MiddlewareInput) Ctx() ReadonlyContext { return mi.ctx }
func (mi MiddlewareInput) Headers() http.Header { return mi.headers }
func (mi MiddlewareInput) Method() string       { return mi.method }
func (mi MiddlewareInput) Path() string         { return mi.path }
func (mi MiddlewareInput) RemoteAddr() string   { return mi.remoteAddr }

type ctxKey string

const (
	RequesterIDKey ctxKey = "X-Requester-ID"
	PathKey        ctxKey = "X-Request-Path"
	AuthHeaderKey  ctxKey = "Authorization"
)

func GenerateMiddlewareInput(r *http.Request) MiddlewareInput {

	// Authorization 헤더 가져오기
	log.Printf("[JWT] input.Path=%q", r.URL.Path)
	log.Printf("[JWT] input.Path=%q", r.Header)

	return MiddlewareInput{
		ctx:        r.Context(),
		method:     r.Method,
		path:       r.URL.Path,
		headers:    r.Header.Clone(), // 안전하게 복사
		remoteAddr: r.RemoteAddr,
	}
}
