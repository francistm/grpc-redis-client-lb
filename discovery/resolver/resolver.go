package resolver

import (
	"context"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/grpclog"
	grpcResolver "google.golang.org/grpc/resolver"
)

var _ grpcResolver.Resolver = (*schemaResolver)(nil)

type schemaResolver struct {
	isClosed uint32

	mu            *sync.RWMutex
	watchTicker   *time.Ticker
	rdb           Redis
	logger        grpclog.LoggerV2
	clientConn    grpcResolver.ClientConn
	serviceName   string
	whitelistNets []*net.IPNet
}

func buildResolverAddr(id, addr string) grpcResolver.Address {
	attrs := attributes.New("id", id)

	return grpcResolver.Address{
		Addr:       addr,
		Attributes: attrs,
	}
}

func (r *schemaResolver) watch() {
	for atomic.LoadUint32(&r.isClosed) == 0 {
		r.ResolveNow(grpcResolver.ResolveNowOptions{})
		<-r.watchTicker.C
	}
}

func (r *schemaResolver) ResolveNow(options grpcResolver.ResolveNowOptions) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if atomic.LoadUint32(&r.isClosed) > 0 {
		return
	}

	ctx := context.Background()
	redisNS := strings.Join([]string{"", r.serviceName, "*"}, ":")
	originalKeys, err := r.rdb.Keys(ctx, redisNS)

	if err != nil {
		r.logger.Errorln(err)
		r.clientConn.ReportError(errors.Wrap(err, "failed to get keys from redis"))
		return
	}

	if len(originalKeys) == 0 {
		r.logger.Errorf("no registered service for %s://%s", schema, r.serviceName)
		r.clientConn.ReportError(errors.Errorf("no registered service for %s://%s", schema, r.serviceName))
		return
	}

	addrs := make([]grpcResolver.Address, 0, len(originalKeys))

	for _, originalKey := range originalKeys {
		originalValue, err := r.rdb.GetStr(ctx, originalKey)
		originalKeyParts := strings.Split(strings.TrimPrefix(originalKey, ":"), ":")

		if err != nil {
			r.logger.Errorln(err)
			continue
		}

		if len(originalKeyParts) != 2 {
			r.logger.Errorf("redis key (%s) has invalid format", originalKey)
			r.clientConn.ReportError(errors.Errorf("redis key (%s) has invalid format", originalKey))
			continue
		}

		// split the ip and port number
		originalIPAddr := strings.Split(originalValue, ":")

		if len(originalIPAddr) != 2 {
			continue
		}

		originalIP := net.ParseIP(originalIPAddr[0])

		// is not a valid ipv4 address
		if originalIP == nil {
			continue
		}

		// skip if not in whitelist
		if len(r.whitelistNets) > 0 {
			for _, i := range r.whitelistNets {
				if i.Contains(originalIP) {
					goto AppendToAddr
				}
			}

			continue
		}

	AppendToAddr:
		addr := buildResolverAddr(originalKeyParts[1], originalValue)
		addrs = append(addrs, addr)
	}

	if len(addrs) == 0 {
		r.clientConn.ReportError(errors.Errorf("no valid service registeration for %s://%s", schema, r.serviceName))
	} else {
		err := r.clientConn.UpdateState(grpcResolver.State{
			Addresses: addrs,
		})

		if err != nil {
			r.clientConn.ReportError(errors.Wrapf(err, "failed to update conn state"))
		}
	}
}

func (r *schemaResolver) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	atomic.StoreUint32(&r.isClosed, 1)
	r.watchTicker.Stop()
}
