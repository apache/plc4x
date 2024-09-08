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

// Code generated by mockery v2.42.2. DO NOT EDIT.

package options

import mock "github.com/stretchr/testify/mock"

// MockDiscoveryOptionRemoteAddress is an autogenerated mock type for the DiscoveryOptionRemoteAddress type
type MockDiscoveryOptionRemoteAddress struct {
	mock.Mock
}

type MockDiscoveryOptionRemoteAddress_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDiscoveryOptionRemoteAddress) EXPECT() *MockDiscoveryOptionRemoteAddress_Expecter {
	return &MockDiscoveryOptionRemoteAddress_Expecter{mock: &_m.Mock}
}

// GetRemoteAddress provides a mock function with given fields:
func (_m *MockDiscoveryOptionRemoteAddress) GetRemoteAddress() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetRemoteAddress")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRemoteAddress'
type MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call struct {
	*mock.Call
}

// GetRemoteAddress is a helper method to define mock.On call
func (_e *MockDiscoveryOptionRemoteAddress_Expecter) GetRemoteAddress() *MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call {
	return &MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call{Call: _e.mock.On("GetRemoteAddress")}
}

func (_c *MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call) Run(run func()) *MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call) Return(_a0 string) *MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call) RunAndReturn(run func() string) *MockDiscoveryOptionRemoteAddress_GetRemoteAddress_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDiscoveryOptionRemoteAddress creates a new instance of MockDiscoveryOptionRemoteAddress. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDiscoveryOptionRemoteAddress(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDiscoveryOptionRemoteAddress {
	mock := &MockDiscoveryOptionRemoteAddress{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
