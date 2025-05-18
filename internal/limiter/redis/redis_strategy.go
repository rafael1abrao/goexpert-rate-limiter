package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStrategy struct {
	Client *redis.Client
}

func NewRedisStrategy(client *redis.Client) *RedisStrategy {
	return &RedisStrategy{Client: client}
}

func (r *RedisStrategy) IsAllowed(ctx context.Context, key string, limit int, blockDurationSeconds int) (bool, error) {
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		err := r.Client.Expire(ctx, key, time.Duration(blockDurationSeconds)*time.Second).Err()
		if err != nil {
			return false, err
		}
	}

	if int(count) > limit {
		return false, nil
	}

	return true, nil
}
