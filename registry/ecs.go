package registry

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func GetECSContainerAddr() (string, error) {
	var (
		endpoint   string
		metadata   awsMetadata
		v3Endpoint = os.Getenv("ECS_CONTAINER_METADATA_URI")
		v4Endpoint = os.Getenv("ECS_CONTAINER_METADATA_URI_V4")
	)

	if len(v3Endpoint) > 0 {
		endpoint = v3Endpoint
	} else if len(v4Endpoint) > 0 {
		endpoint = v4Endpoint
	} else {
		return "", ErrNotInECS
	}

	res, err := http.Get(endpoint)

	if err != nil {
		return "", errors.WithStack(err)
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&metadata); err != nil {
		return "", errors.WithStack(err)
	}

	for _, network := range metadata.Networks {
		if len(network.IPv4Addresses) > 0 {
			return network.IPv4Addresses[0], nil
		}
	}

	return "", errors.New("ecs doesn't have any ipv4 addr")
}
