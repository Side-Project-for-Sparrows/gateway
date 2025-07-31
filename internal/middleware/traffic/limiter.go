package traffic

import (
	"time"
)

type Limiter interface {
	IsOverLimit(key string, now time.Time) bool
}
