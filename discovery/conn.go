package discovery

import (
	"sync"

	grpcresolver "google.golang.org/grpc/resolver"

	"github.com/francistm/grpc-redis-client-lb/discovery/resolver"
)

var _registryOnce sync.Once

func RegisterSchema(rds resolver.Redis, opts ...resolver.BuilderOptApplyFn) {
	_registryOnce.Do(func() {
		grpcresolver.Register(resolver.NewBuilder(rds, opts...))
	})
}
