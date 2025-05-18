package limiter

import "context"

type RateLimiterStrategy interface {
	IsAllowed(ctx context.Context, key string, limit int, blockDurationSeconds int) (bool, error)
}
