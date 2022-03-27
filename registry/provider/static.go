package provider

import "context"

type StaticProvider struct {
	Addr string
}

func (p *StaticProvider) DetectHostAddr(ctx context.Context) (string, error) {
	return p.Addr, nil
}
