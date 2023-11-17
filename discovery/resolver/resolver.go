package resolver

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/resolver"
	grpcResolver "google.golang.org/grpc/resolver"
)

var _ grpcResolver.Resolver = (*schemaResolver)(nil)

type schemaResolver struct {
	ctx            context.Context
	cancel         context.CancelFunc
	resolveChan    chan struct{}
	exitChan       chan struct{}
	lookupInterval time.Duration
	rdb            Redis
	logger         grpclog.LoggerV2
	clientConn     grpcResolver.ClientConn
	serviceName    string
}

func buildResolverAddr(id, addr string) grpcResolver.Address {
	attrs := attributes.New("id", id)

	return grpcResolver.Address{
		Addr:       addr,
		Attributes: attrs,
	}
}

func (r *schemaResolver) watch() {
	defer func() {
		r.exitChan <- struct{}{}
	}()

	for {
		state, err := r.lookup()

		if err != nil {
			r.clientConn.ReportError(err)
		} else if err := r.clientConn.UpdateState(state); err != nil {
			r.logger.Errorf("update state failed for service %s due %s", r.serviceName, err.Error())
		}

		timer := time.NewTimer(r.lookupInterval)

		select {
		case <-r.ctx.Done():
			timer.Stop()
			return

		case <-timer.C:
			// block until interval timer emitted

		case <-r.resolveChan:
			// block until ResolvedNow called
		}
	}
}

func (r *schemaResolver) lookup() (resolver.State, error) {
	var (
		ctx     = r.ctx
		state   = resolver.State{}
		redisNS = strings.Join([]string{"", r.serviceName, "*"}, ":")
	)

	r.logger.Infof("looking up the address for service %s", r.serviceName)

	originalKeys, err := r.rdb.Keys(ctx, redisNS)

	if err != nil {
		return state, fmt.Errorf("failed to get keys from redis when discover service %s %w", r.serviceName, err)
	}

	if len(originalKeys) == 0 {
		return state, nil
	}

	state.Addresses = make([]resolver.Address, 0, len(originalKeys))

	for _, originalKey := range originalKeys {
		originalValue, err := r.rdb.GetStr(ctx, originalKey)

		if err != nil {
			r.logger.Errorf("failed to get value of key %s due %s", originalKey, err)
			continue
		}

		// ensure the value will follow the format of "host:port"
		if !strings.Contains(originalValue, ":") {
			continue
		}

		state.Addresses = append(state.Addresses, resolver.Address{
			Addr:       originalValue,
			ServerName: r.serviceName,
		})
	}

	return state, nil
}

func (r *schemaResolver) ResolveNow(options grpcResolver.ResolveNowOptions) {
	select {
	case r.resolveChan <- struct{}{}:
		// begin to resolve

	default:
		// resolving in progress
	}
}

func (r *schemaResolver) Close() {
	r.cancel()
	<-r.exitChan
}
