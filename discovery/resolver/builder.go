package resolver

import (
	"net"
	"strings"
	"sync"
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

	res := &schemaResolver{
		mu:            new(sync.RWMutex),
		rdb:           b.rds,
		clientConn:    cc,
		serviceName:   target.URL.Host,
		watchTicker:   time.NewTicker(5 * time.Second),
		whitelistNets: make([]*net.IPNet, 0, len(b.whitelistSubnets)),
	}

	if b.logger == nil {
		res.logger = grpclog.Component(b.Scheme())
	} else {
		res.logger = b.logger
	}

	for _, s := range b.whitelistSubnets {
		if _, net, err := net.ParseCIDR(s); err == nil {
			res.whitelistNets = append(res.whitelistNets, net)
		}
	}

	go res.watch()

	return res, nil
}

func (b *builder) Scheme() string {
	return schema
}
