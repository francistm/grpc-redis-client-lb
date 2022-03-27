package registry

import "context"

type RedisClient interface {
	SetStrWithExpire(ctx context.Context, key, value string, ttl int64) error
}

type Provider interface {
	DetectHostAddr(ctx context.Context) (string, error)
}
