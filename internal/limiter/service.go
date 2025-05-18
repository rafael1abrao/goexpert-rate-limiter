package limiter

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

type LimiterService struct {
	strategy             RateLimiterStrategy
	rateLimitIP          int
	rateLimitToken       int
	blockDurationSeconds int
}

func NewLimiterService(strategy RateLimiterStrategy) *LimiterService {
	return &LimiterService{
		strategy:             strategy,
		rateLimitIP:          mustGetEnvAsInt("RATE_LIMIT_IP"),
		rateLimitToken:       mustGetEnvAsInt("RATE_LIMIT_TOKEN_DEFAULT"),
		blockDurationSeconds: mustGetEnvAsInt("BLOCK_TIME_SECONDS"),
	}
}

func mustGetEnvAsInt(key string) int {
	v, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		panic(fmt.Sprintf("Invalid env var: %s", key))
	}
	return v
}

func (l *LimiterService) IsRequestAllowed(ctx context.Context, ip, token string) (bool, error) {
	key := ""
	limit := l.rateLimitIP

	if token != "" {
		key = fmt.Sprintf("token:%s", token)
		limit = l.rateLimitToken
	} else {
		key = fmt.Sprintf("ip:%s", ip)
	}

	return l.strategy.IsAllowed(ctx, key, limit, l.blockDurationSeconds)
}
