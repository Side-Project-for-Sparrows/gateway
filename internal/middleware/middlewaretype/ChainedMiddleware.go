package middlewaretype

import (
	"log"
)

// 순차 실행 미들웨어 그룹
type ChainedMiddleware struct {
	middlewares []Middleware
}

func NewMiddlewareChain() *ChainedMiddleware {
	return &ChainedMiddleware{middlewares: []Middleware{}}
}

func (mc *ChainedMiddleware) AndThen(mw Middleware) *ChainedMiddleware {
	mc.middlewares = append(mc.middlewares, mw)
	return mc
}

// 순차 실행하면서 HeaderPatch들을 수집
func (mc *ChainedMiddleware) Execute(input MiddlewareInput) ([]*HeaderPatch, error) {
	var patches []*HeaderPatch

	for i, mw := range mc.middlewares {
		patch, err := mw(input)
		if err != nil {
			log.Printf("[ChainedMiddleware] Middleware #%d failed: %v", i, err)
			return nil, err
		}
		if patch != nil {
			patches = append(patches, patch)
		}
	}
	return patches, nil
}

// 다른 체인에 조합되기 위함.
func (mc *ChainedMiddleware) AsMiddleware() Middleware {
	return func(input MiddlewareInput) (*HeaderPatch, error) {
		patches, err := mc.Execute(input)
		if err != nil {
			return nil, err
		}

		// 여러 patch를 하나로 merge
		return mergePatches(patches), nil
	}
}
