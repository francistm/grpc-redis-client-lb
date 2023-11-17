package resolver

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/grpclog"
	grpcResolver "google.golang.org/grpc/resolver"
)

//go:generate go run github.com/vektra/mockery/v2 --name ClientConn --srcpkg google.golang.org/grpc/resolver --output ../../mocks --exported

const schema = "redis"

type builder struct {
	rds              Redis
	logger           grpclog.LoggerV2
	whitelistSubnets []string
}

func NewBuilder(rds Redis, opts ...BuilderOptApplyFn) grpcResolver.Builder {
	b := &builder{
		rds:              rds,
		whitelistSubnets: []string{},
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *builder) Build(target grpcResolver.Target, cc grpcResolver.ClientConn, opts grpcResolver.BuildOptions) (grpcResolver.Resolver, error) {
	if !strings.EqualFold(schema, target.URL.Scheme) {
		return nil, errors.Errorf("unexpected schema: %s", target.URL.Scheme)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	res := &schemaResolver{
		ctx:            ctx,
		cancel:         cancelFunc,
		resolveChan:    make(chan struct{}),
		exitChan:       make(chan struct{}),
		lookupInterval: 5 * time.Second,
		rdb:            b.rds,
		clientConn:     cc,
		serviceName:    target.URL.Host,
	}

	if b.logger == nil {
		res.logger = grpclog.Component(b.Scheme())
	} else {
		res.logger = b.logger
	}

	go res.watch()

	return res, nil
}

func (b *builder) Scheme() string {
	return schema
}
