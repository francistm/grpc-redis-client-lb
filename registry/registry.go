package registry

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func RegisterService(ctx context.Context, rdsClient RedisClient, serviceName string, listenPort int, opts ...Option) (string, error) {
	var (
		provider                  Provider
		registryTTL               time.Duration
		registryKeepAliveDuration time.Duration
	)

	for _, opt := range opts {
		switch opt.kind() {
		case registryOptionProvider:
			provider = opt.value().(Provider)

		case registryOptionRegistryTTL:
			registryTTL = opt.value().(time.Duration)

		case registryOptionRegistryKeepAliveDuration:
			registryKeepAliveDuration = opt.value().(time.Duration)
		}
	}

	if provider == nil {
		return "", errors.New("missing host addr provider")
	}

	if registryTTL <= registryKeepAliveDuration {
		return "", errors.New("registry TTL should be greater than registry keepalive duration")
	}

	instanceUUID := uuid.NewString()
	instanceID := strings.Join([]string{"", serviceName, instanceUUID}, ":")

	hostAddr, err := provider.DetectHostAddr(ctx)

	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve container addr")
	}

	listenAddr := fmt.Sprintf("%s:%d", hostAddr, listenPort)

	go func() {
		var t *time.Timer

		for {
			t = time.NewTimer(registryKeepAliveDuration)
			_ = rdsClient.SetStrWithExpire(ctx, instanceID, listenAddr, int64(registryTTL.Seconds()))

			select {
			case <-ctx.Done():
				t.Stop()
				return

			case <-t.C:
			}
		}
	}()

	return listenAddr, nil
}
