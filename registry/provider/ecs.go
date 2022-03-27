package provider

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

const (
	awsECSMetadataUriName   = "ECS_CONTAINER_METADATA_URI"
	awsECSMetadataUriV4Name = "ECS_CONTAINER_METADATA_URI_V4"
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

type ECSProvider struct {
}

func (p *ECSProvider) detectEndpointUri() (string, error) {
	if os.Getenv(awsECSMetadataUriName) != "" {
		return os.Getenv(awsECSMetadataUriName), nil
	}

	if os.Getenv(awsECSMetadataUriV4Name) != "" {
		return os.Getenv(awsECSMetadataUriV4Name), nil
	}

	return "", ErrNotInECS
}

func (p *ECSProvider) DetectHostAddr(ctx context.Context) (string, error) {
	metadata := &awsMetadata{}
	endpointUri, err := p.detectEndpointUri()

	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointUri, nil)

	if err != nil {
		return "", err
	}

	response, err := new(http.Client).Do(request)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(metadata); err != nil {
		return "", err
	}

	for _, network := range metadata.Networks {
		if len(network.IPv4Addresses) > 0 {
			return network.IPv4Addresses[0], nil
		}
	}

	return "", errors.New("ecs container doesn't have any ipv4 addr")
}
