package registry

import "context"

type redisClient interface {
	SetStrWithExpire(ctx context.Context, key, value string, ttl int64) error
}
