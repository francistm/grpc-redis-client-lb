package resolver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	grpcResolver "google.golang.org/grpc/resolver"

	"github.com/francistm/grpc-redis-client-lb/discovery/resolver/logger"
)

func TestNewBuilder(t *testing.T) {
	nopLogger := &logger.NopLogger{}

	type args struct {
		rds    Redis
		subnet []string
	}
	tests := []struct {
		name string
		args args
		want grpcResolver.Builder
	}{
		{
			name: "case 1",
			args: args{},
			want: &builder{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBuilder(tt.args.rds, WithLogger(nopLogger))
			assert.IsType(t, tt.want, got)
		})
	}
}

func Test_builder_Scheme(t *testing.T) {
	tests := []struct {
		name string
		b    *builder
		want string
	}{
		{
			name: "case 1",
			b:    &builder{},
			want: schema,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Scheme()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_builder_Build(t *testing.T) {
	t.Skip("Not implemented yet")
}
