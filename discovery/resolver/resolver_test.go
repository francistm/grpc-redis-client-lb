package resolver

import (
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/resolver"
	grpcResolver "google.golang.org/grpc/resolver"

	"github.com/francistm/grpc-redis-client-lb/discovery/resolver/logger"
	"github.com/francistm/grpc-redis-client-lb/mocks"
)

func Test_schemaResolver_ResolveNow(t *testing.T) {
	type args struct {
		options grpcResolver.ResolveNowOptions
	}

	rds := &mocks.Redis{}
	clientConn := &mocks.ClientConn{}
	nopLogger := &logger.NopLogger{}

	// case 2
	rds.On("Keys", mock.Anything, ":case-2:*").Return([]string{}, errors.Errorf("test error"))
	clientConn.
		On("ReportError", mock.MatchedBy(func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "test error")
		})).
		Return(nil)

	// case 3
	rds.On("Keys", mock.Anything, ":case-3:*").Return([]string{}, nil)
	clientConn.
		On("ReportError", mock.MatchedBy(func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "no registered service")
		})).
		Return(nil)

	// case 4
	rds.On("Keys", mock.Anything, ":case-4:*").Return([]string{"a", "b"}, nil)
	rds.On("GetStr", mock.Anything, "a").Return("", errors.Errorf("test error"))
	rds.On("GetStr", mock.Anything, "b").Return("", errors.Errorf("test error"))
	clientConn.
		On("ReportError", mock.MatchedBy(func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "no valid service registeration")
		})).
		Return(nil)

	// case 5
	rds.On("Keys", mock.Anything, ":case-5:*").Return([]string{"a", "b"}, nil)
	rds.On("GetStr", mock.Anything, "a").Return("c", nil)

	// case 6
	rds.On("Keys", mock.Anything, ":case-6:*").Return([]string{"service:a", "service:b", "service:c", "service:d", "service:e"}, nil)
	rds.On("GetStr", mock.Anything, "service:a").Return("172.17.0.1:8888", nil)
	rds.On("GetStr", mock.Anything, "service:b").Return("172.17.0.2:8888", nil)
	rds.On("GetStr", mock.Anything, "service:c").Return("abcdefg", nil)
	rds.On("GetStr", mock.Anything, "service:d").Return("10.0.0.5:8888", nil)
	rds.On("GetStr", mock.Anything, "service:e").Return("abcdefg:8888", nil)
	clientConn.
		On("UpdateState", mock.MatchedBy(func(state resolver.State) bool {
			expectedList := [][2]string{
				{"a", "172.17.0.1:8888"},
				{"b", "172.17.0.2:8888"},
			}

			if len(state.Addresses) != len(expectedList) {
				return false
			}

			for _, addr := range state.Addresses {
				isMatched := false
				for _, expectedAddr := range expectedList {
					if addr.Addr == expectedAddr[1] && addr.Attributes.Value("id") == expectedAddr[0] {
						isMatched = true
						break
					}
				}

				if !isMatched {
					return false
				}
			}

			return true
		})).
		Return(nil)

	_, whiteList, _ := net.ParseCIDR("172.17.0.0/23")

	tests := []struct {
		name string
		r    *schemaResolver
		args args
	}{
		{
			name: "case 1",
			r: &schemaResolver{
				serviceName: "case-1",
				mu:          &sync.RWMutex{},
				isClosed:    true,
			},
		},
		{
			name: "case 2",
			r: &schemaResolver{
				mu:          &sync.RWMutex{},
				serviceName: "case-2",
				clientConn:  clientConn,
				rdb:         rds,
				logger:      nopLogger,
			},
			args: args{},
		},
		{
			name: "case 3",
			r: &schemaResolver{
				serviceName: "case-3",
				mu:          &sync.RWMutex{},
				clientConn:  clientConn,
				rdb:         rds,
				logger:      nopLogger,
			},
			args: args{},
		},
		{
			name: "case 4",
			r: &schemaResolver{
				mu:          &sync.RWMutex{},
				serviceName: "case-4",
				clientConn:  clientConn,
				rdb:         rds,
				logger:      nopLogger,
			},
			args: args{},
		},
		{
			name: "case 5",
			r: &schemaResolver{
				mu:          &sync.RWMutex{},
				serviceName: "case-5",
				clientConn:  clientConn,
				rdb:         rds,
				logger:      nopLogger,
			},
			args: args{},
		},
		{
			name: "case 6",
			r: &schemaResolver{
				mu:            &sync.RWMutex{},
				serviceName:   "case-6",
				clientConn:    clientConn,
				whitelistNets: []*net.IPNet{whiteList},
				rdb:           rds,
				logger:        nopLogger,
			},
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.ResolveNow(tt.args.options)
		})
	}

	rds.AssertExpectations(t)
	clientConn.AssertExpectations(t)
}

func Test_schemaResolver_Close(t *testing.T) {
	tests := []struct {
		name string
		r    *schemaResolver
	}{
		{
			name: "case 1",
			r: &schemaResolver{
				mu:          &sync.RWMutex{},
				watchTicker: time.NewTicker(10 * time.Second),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.Close()
			assert.True(t, tt.r.isClosed)
		})
	}
}

func Test_schemaResolver_watch(t *testing.T) {
	rdb := &mocks.Redis{}
	clientConn := &mocks.ClientConn{}

	// case 1
	rdb.On("Keys", mock.Anything, ":test-service:*").Return(([]string)(nil), nil)
	clientConn.
		On("ReportError", mock.MatchedBy(func(err error) bool {
			return err != nil && strings.Contains(err.Error(), "no registered service for")
		})).
		Return(nil)

	r := &schemaResolver{
		rdb:         rdb,
		mu:          &sync.RWMutex{},
		logger:      &logger.NopLogger{},
		serviceName: "test-service",
		clientConn:  clientConn,
		watchTicker: time.NewTicker(100 * time.Millisecond),
	}

	time.AfterFunc(500*time.Millisecond, func() {
		r.Close()
	})

	tests := []struct {
		name string
		r    *schemaResolver
	}{
		{
			name: "case 1",
			r:    r,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.watch()
		})
	}

	rdb.AssertExpectations(t)
	clientConn.AssertExpectations(t)
}
