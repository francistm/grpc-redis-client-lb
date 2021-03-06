// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	resolver "google.golang.org/grpc/resolver"

	serviceconfig "google.golang.org/grpc/serviceconfig"
)

// ClientConn is an autogenerated mock type for the ClientConn type
type ClientConn struct {
	mock.Mock
}

// NewAddress provides a mock function with given fields: addresses
func (_m *ClientConn) NewAddress(addresses []resolver.Address) {
	_m.Called(addresses)
}

// NewServiceConfig provides a mock function with given fields: serviceConfig
func (_m *ClientConn) NewServiceConfig(serviceConfig string) {
	_m.Called(serviceConfig)
}

// ParseServiceConfig provides a mock function with given fields: serviceConfigJSON
func (_m *ClientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	ret := _m.Called(serviceConfigJSON)

	var r0 *serviceconfig.ParseResult
	if rf, ok := ret.Get(0).(func(string) *serviceconfig.ParseResult); ok {
		r0 = rf(serviceConfigJSON)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*serviceconfig.ParseResult)
		}
	}

	return r0
}

// ReportError provides a mock function with given fields: _a0
func (_m *ClientConn) ReportError(_a0 error) {
	_m.Called(_a0)
}

// UpdateState provides a mock function with given fields: _a0
func (_m *ClientConn) UpdateState(_a0 resolver.State) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(resolver.State) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
