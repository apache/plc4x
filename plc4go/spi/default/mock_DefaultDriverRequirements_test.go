/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Code generated by mockery v2.28.1. DO NOT EDIT.

package _default

import (
	context "context"

	model "github.com/apache/plc4x/plc4go/pkg/api/model"
	mock "github.com/stretchr/testify/mock"

	options "github.com/apache/plc4x/plc4go/spi/options"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"

	transports "github.com/apache/plc4x/plc4go/spi/transports"

	url "net/url"
)

// MockDefaultDriverRequirements is an autogenerated mock type for the DefaultDriverRequirements type
type MockDefaultDriverRequirements struct {
	mock.Mock
}

type MockDefaultDriverRequirements_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDefaultDriverRequirements) EXPECT() *MockDefaultDriverRequirements_Expecter {
	return &MockDefaultDriverRequirements_Expecter{mock: &_m.Mock}
}

// DiscoverWithContext provides a mock function with given fields: callback, event, discoveryOptions
func (_m *MockDefaultDriverRequirements) DiscoverWithContext(callback context.Context, event func(model.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	_va := make([]interface{}, len(discoveryOptions))
	for _i := range discoveryOptions {
		_va[_i] = discoveryOptions[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, callback, event)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(model.PlcDiscoveryItem), ...options.WithDiscoveryOption) error); ok {
		r0 = rf(callback, event, discoveryOptions...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDefaultDriverRequirements_DiscoverWithContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DiscoverWithContext'
type MockDefaultDriverRequirements_DiscoverWithContext_Call struct {
	*mock.Call
}

// DiscoverWithContext is a helper method to define mock.On call
//   - callback context.Context
//   - event func(model.PlcDiscoveryItem)
//   - discoveryOptions ...options.WithDiscoveryOption
func (_e *MockDefaultDriverRequirements_Expecter) DiscoverWithContext(callback interface{}, event interface{}, discoveryOptions ...interface{}) *MockDefaultDriverRequirements_DiscoverWithContext_Call {
	return &MockDefaultDriverRequirements_DiscoverWithContext_Call{Call: _e.mock.On("DiscoverWithContext",
		append([]interface{}{callback, event}, discoveryOptions...)...)}
}

func (_c *MockDefaultDriverRequirements_DiscoverWithContext_Call) Run(run func(callback context.Context, event func(model.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption)) *MockDefaultDriverRequirements_DiscoverWithContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]options.WithDiscoveryOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(options.WithDiscoveryOption)
			}
		}
		run(args[0].(context.Context), args[1].(func(model.PlcDiscoveryItem)), variadicArgs...)
	})
	return _c
}

func (_c *MockDefaultDriverRequirements_DiscoverWithContext_Call) Return(_a0 error) *MockDefaultDriverRequirements_DiscoverWithContext_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDefaultDriverRequirements_DiscoverWithContext_Call) RunAndReturn(run func(context.Context, func(model.PlcDiscoveryItem), ...options.WithDiscoveryOption) error) *MockDefaultDriverRequirements_DiscoverWithContext_Call {
	_c.Call.Return(run)
	return _c
}

// GetConnectionWithContext provides a mock function with given fields: ctx, transportUrl, _a2, _a3
func (_m *MockDefaultDriverRequirements) GetConnectionWithContext(ctx context.Context, transportUrl url.URL, _a2 map[string]transports.Transport, _a3 map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	ret := _m.Called(ctx, transportUrl, _a2, _a3)

	var r0 <-chan plc4go.PlcConnectionConnectResult
	if rf, ok := ret.Get(0).(func(context.Context, url.URL, map[string]transports.Transport, map[string][]string) <-chan plc4go.PlcConnectionConnectResult); ok {
		r0 = rf(ctx, transportUrl, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan plc4go.PlcConnectionConnectResult)
		}
	}

	return r0
}

// MockDefaultDriverRequirements_GetConnectionWithContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConnectionWithContext'
type MockDefaultDriverRequirements_GetConnectionWithContext_Call struct {
	*mock.Call
}

// GetConnectionWithContext is a helper method to define mock.On call
//   - ctx context.Context
//   - transportUrl url.URL
//   - _a2 map[string]transports.Transport
//   - _a3 map[string][]string
func (_e *MockDefaultDriverRequirements_Expecter) GetConnectionWithContext(ctx interface{}, transportUrl interface{}, _a2 interface{}, _a3 interface{}) *MockDefaultDriverRequirements_GetConnectionWithContext_Call {
	return &MockDefaultDriverRequirements_GetConnectionWithContext_Call{Call: _e.mock.On("GetConnectionWithContext", ctx, transportUrl, _a2, _a3)}
}

func (_c *MockDefaultDriverRequirements_GetConnectionWithContext_Call) Run(run func(ctx context.Context, transportUrl url.URL, _a2 map[string]transports.Transport, _a3 map[string][]string)) *MockDefaultDriverRequirements_GetConnectionWithContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(url.URL), args[2].(map[string]transports.Transport), args[3].(map[string][]string))
	})
	return _c
}

func (_c *MockDefaultDriverRequirements_GetConnectionWithContext_Call) Return(_a0 <-chan plc4go.PlcConnectionConnectResult) *MockDefaultDriverRequirements_GetConnectionWithContext_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDefaultDriverRequirements_GetConnectionWithContext_Call) RunAndReturn(run func(context.Context, url.URL, map[string]transports.Transport, map[string][]string) <-chan plc4go.PlcConnectionConnectResult) *MockDefaultDriverRequirements_GetConnectionWithContext_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDefaultDriverRequirements interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDefaultDriverRequirements creates a new instance of MockDefaultDriverRequirements. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDefaultDriverRequirements(t mockConstructorTestingTNewMockDefaultDriverRequirements) *MockDefaultDriverRequirements {
	mock := &MockDefaultDriverRequirements{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
