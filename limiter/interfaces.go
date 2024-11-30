package limiter

import "context"

type Limiter interface {
	Allow(ctx context.Context, key string, limit int) (bool, error)
	Block(ctx context.Context, key string) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}
