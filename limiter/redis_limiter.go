package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client        *redis.Client
	ipRate        int
	tokenRate     int
	blockDuration time.Duration
}

func NewRedisLimiter(address string, ipRate, tokenRate int, blockDuration time.Duration) *RedisLimiter {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})
	client.FlushAll(context.Background())
	return &RedisLimiter{
		client:        client,
		ipRate:        ipRate,
		tokenRate:     tokenRate,
		blockDuration: blockDuration,
	}
}

func (r *RedisLimiter) Allow(ctx context.Context, key string, limit int) (bool, error) {
	c, err := r.client.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return false, err
	}
	if c >= limit {
		return false, nil
	}
	pipe := r.client.TxPipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Second)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisLimiter) Block(ctx context.Context, key string) error {
	return r.client.Set(ctx, fmt.Sprintf("blocked:%s", key), 1, r.blockDuration*time.Second).Err()
}

func (r *RedisLimiter) IsBlocked(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, fmt.Sprintf("blocked:%s", key)).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
