package resolver

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	grpcResolver "google.golang.org/grpc/resolver"

	"github.com/francistm/grpc-redis-client-lb/discovery/resolver/logger"
	"github.com/francistm/grpc-redis-client-lb/mocks"
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
	type args struct {
		target grpcResolver.Target
		cc     grpcResolver.ClientConn
		opts   grpcResolver.BuildOptions
	}

	assertion := assert.New(t)
	nopLogger := &logger.NopLogger{}
	rdsCalled := make(chan struct{})
	reportErrCalled := make(chan struct{})

	rds := &mocks.Redis{}
	clientConn := &mocks.ClientConn{}

	rds.On("Keys", mock.Anything, ":abde:*").
		Return([]string{}, nil).
		Run(func(args mock.Arguments) {
			rdsCalled <- struct{}{}
		})
	clientConn.
		On("ReportError", mock.MatchedBy(func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "no registered service for")
		})).
		Run(func(args mock.Arguments) {
			reportErrCalled <- struct{}{}
		})

	b := &builder{
		rds:              rds,
		logger:           nopLogger,
		whitelistSubnets: []string{"172.17.0.0/22"},
	}
	tests := []struct {
		name    string
		b       *builder
		args    args
		wantErr bool
	}{
		{
			name: "case 1",
			b:    b,
			args: args{
				target: grpcResolver.Target{
					URL: url.URL{
						Scheme: "",
						Host:   "abcd",
					},
				},
				cc: clientConn,
			},
			wantErr: true,
		},
		{
			name: "case 2",
			b:    b,
			args: args{
				target: grpcResolver.Target{
					URL: url.URL{
						Scheme: "redis",
						Host:   "abde",
					},
				},
				cc: clientConn,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver, err := tt.b.Build(tt.args.target, tt.args.cc, tt.args.opts)

			if tt.wantErr {
				assertion.Error(err)
			} else {
				assertion.NoError(err)
				assertion.NotNil(resolver)
			}
		})
	}

	_, _ = <-rdsCalled, <-reportErrCalled

	rds.AssertExpectations(t)
	clientConn.AssertExpectations(t)
}
