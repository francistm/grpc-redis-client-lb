package registry

import (
	"time"

	"github.com/francistm/grpc-redis-client-lb/registry/provider"
)

type RegistryOption struct {
	provider                  Provider
	registryTTL               time.Duration
	registryKeepAliveDuration time.Duration
}

type RegistryOptionApplyFn func(opts *RegistryOption)

func WithECSProvider() RegistryOptionApplyFn {
	return func(opts *RegistryOption) {
		opts.provider = &provider.ECSProvider{}
	}
}

func WithStaticProvider(addr string) RegistryOptionApplyFn {
	return func(opts *RegistryOption) {
		opts.provider = &provider.StaticProvider{
			Addr: addr,
		}
	}
}

func WithRegistryTTL(d time.Duration) RegistryOptionApplyFn {
	return func(opts *RegistryOption) {
		opts.registryTTL = d
	}
}

func WithKeepaliveDuration(d time.Duration) RegistryOptionApplyFn {
	return func(opts *RegistryOption) {
		opts.registryKeepAliveDuration = d
	}
}
