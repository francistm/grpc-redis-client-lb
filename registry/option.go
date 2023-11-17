package registry

import (
	"time"

	"github.com/francistm/grpc-redis-client-lb/registry/provider"
)

type optionKind int

const (
	registryOptionProvider optionKind = iota
	registryOptionRegistryTTL
	registryOptionRegistryKeepAliveDuration
)

type Option interface {
	kind() optionKind
	value() any
}

type providerOption struct {
	provider Provider
}

func (o providerOption) kind() optionKind {
	return registryOptionProvider
}

func (o providerOption) value() any {
	return o.provider
}

func WithECSProvider() Option {
	return providerOption{
		provider: &provider.ECSProvider{},
	}
}

func WithStaticProvider(addr string) Option {
	return providerOption{
		provider: &provider.StaticProvider{
			Addr: addr,
		},
	}
}

type registryTTLOption time.Duration

func (o registryTTLOption) kind() optionKind {
	return registryOptionRegistryTTL
}

func (o registryTTLOption) value() any {
	return time.Duration(o)
}

func WithRegistryTTL(d time.Duration) Option {
	return registryTTLOption(d)
}

type registryKeepAliveDurationOption time.Duration

func (o registryKeepAliveDurationOption) kind() optionKind {
	return registryOptionRegistryKeepAliveDuration
}

func (o registryKeepAliveDurationOption) value() any {
	return time.Duration(o)
}

func WithKeepaliveDuration(d time.Duration) Option {
	return registryKeepAliveDurationOption(d)
}
