package resolver

import "context"

//go:generate go run github.com/vektra/mockery/v2 --name Redis --output ../../mocks
type Redis interface {
	Keys(ctx context.Context, pattern string) ([]string, error)
	GetStr(ctx context.Context, key string) (string, error)
}
