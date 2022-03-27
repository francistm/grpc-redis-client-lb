package registry

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	instanceIDSeparate       = ":"
	registeredServiceTTL     = 10
	registerHeartbeatSeconds = 3
)

var (
	ErrNotInECS = errors.New("server didn't running in ECS")
)

type awsMetadata struct {
	Networks []*awsMetadataNetwork `json:"Networks,omitempty"`
}

type awsMetadataNetwork struct {
	NetworkMode   string   `json:"NetworkMode,omitempty"`
	IPv4Addresses []string `json:"IPv4Addresses,omitempty"`
}

func RegisterService(ctx context.Context, rdsClient redisClient, serviceName, listenPort string) (string, error) {
	instanceUUID := uuid.NewString()
	instanceID := strings.Join([]string{"", serviceName, instanceUUID}, instanceIDSeparate)

	addr, err := GetECSContainerAddr()

	if errors.Is(err, ErrNotInECS) {
		return "", nil
	}

	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve container addr")
	}

	port, err := normalizeListenPort(listenPort)

	if err != nil {
		return "", errors.Wrapf(err, "%s is not an valid listenPort", listenPort)
	}

	listenAddr := fmt.Sprintf("%s:%d", addr, port)

	go func() {
		t := time.NewTicker(registerHeartbeatSeconds * time.Second)

		for ; true; <-t.C {
			_ = rdsClient.SetStrWithExpire(ctx, instanceID, listenAddr, registeredServiceTTL)
		}
	}()

	return listenAddr, nil
}

func normalizeListenPort(portS string) (port int, err error) {
	portS = strings.TrimPrefix(portS, ":")
	portInt64, err := strconv.ParseInt(portS, 10, 64)

	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int(portInt64), nil
}
