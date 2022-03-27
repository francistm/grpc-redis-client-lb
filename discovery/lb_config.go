package discovery

import "encoding/json"

type loadBalancingConfig struct {
	RoundRobin struct{} `json:"round_robin"`
}

type serviceConfig struct {
	LoadBalancingConfig []loadBalancingConfig `json:"loadBalancingConfig"`
}

func newServiceConfigJSON() string {
	c := serviceConfig{
		LoadBalancingConfig: []loadBalancingConfig{
			{},
		},
	}

	b, _ := json.Marshal(c)

	return string(b)
}
