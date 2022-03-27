package registry

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func RegisterService(ctx context.Context, rdsClient RedisClient, serviceName string, listenPort int, opts ...RegistryOptionApplyFn) (string, error) {
	option := &RegistryOption{
		registryTTL:               10 * time.Second,
		registryKeepAliveDuration: 3 * time.Second,
	}

	for _, applyFn := range opts {
		applyFn(option)
	}

	if option.provider == nil {
		return "", errors.New("missing host addr provider")
	}

	if option.registryTTL <= option.registryKeepAliveDuration {
		return "", errors.New("registry TTL should be greater than registry keepalive duration")
	}

	instanceUUID := uuid.NewString()
	instanceID := strings.Join([]string{"", serviceName, instanceUUID}, ":")

	hostAddr, err := option.provider.DetectHostAddr(ctx)

	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve container addr")
	}

	listenAddr := fmt.Sprintf("%s:%d", hostAddr, listenPort)

	go func() {
		t := time.NewTicker(option.registryKeepAliveDuration)

		select {
		case <-t.C:
			_ = rdsClient.SetStrWithExpire(ctx, instanceID, listenAddr, int64(option.registryTTL.Seconds()))

		case <-ctx.Done():
			t.Stop()
		}
	}()

	return listenAddr, nil
}
