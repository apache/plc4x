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

// Code generated by mockery v2.28.2. DO NOT EDIT.

package options

import mock "github.com/stretchr/testify/mock"

// MockDiscoveryOptionProtocolSpecific is an autogenerated mock type for the DiscoveryOptionProtocolSpecific type
type MockDiscoveryOptionProtocolSpecific struct {
	mock.Mock
}

type MockDiscoveryOptionProtocolSpecific_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDiscoveryOptionProtocolSpecific) EXPECT() *MockDiscoveryOptionProtocolSpecific_Expecter {
	return &MockDiscoveryOptionProtocolSpecific_Expecter{mock: &_m.Mock}
}

// GetKey provides a mock function with given fields:
func (_m *MockDiscoveryOptionProtocolSpecific) GetKey() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockDiscoveryOptionProtocolSpecific_GetKey_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetKey'
type MockDiscoveryOptionProtocolSpecific_GetKey_Call struct {
	*mock.Call
}

// GetKey is a helper method to define mock.On call
func (_e *MockDiscoveryOptionProtocolSpecific_Expecter) GetKey() *MockDiscoveryOptionProtocolSpecific_GetKey_Call {
	return &MockDiscoveryOptionProtocolSpecific_GetKey_Call{Call: _e.mock.On("GetKey")}
}

func (_c *MockDiscoveryOptionProtocolSpecific_GetKey_Call) Run(run func()) *MockDiscoveryOptionProtocolSpecific_GetKey_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDiscoveryOptionProtocolSpecific_GetKey_Call) Return(_a0 string) *MockDiscoveryOptionProtocolSpecific_GetKey_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDiscoveryOptionProtocolSpecific_GetKey_Call) RunAndReturn(run func() string) *MockDiscoveryOptionProtocolSpecific_GetKey_Call {
	_c.Call.Return(run)
	return _c
}

// GetValue provides a mock function with given fields:
func (_m *MockDiscoveryOptionProtocolSpecific) GetValue() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// MockDiscoveryOptionProtocolSpecific_GetValue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetValue'
type MockDiscoveryOptionProtocolSpecific_GetValue_Call struct {
	*mock.Call
}

// GetValue is a helper method to define mock.On call
func (_e *MockDiscoveryOptionProtocolSpecific_Expecter) GetValue() *MockDiscoveryOptionProtocolSpecific_GetValue_Call {
	return &MockDiscoveryOptionProtocolSpecific_GetValue_Call{Call: _e.mock.On("GetValue")}
}

func (_c *MockDiscoveryOptionProtocolSpecific_GetValue_Call) Run(run func()) *MockDiscoveryOptionProtocolSpecific_GetValue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDiscoveryOptionProtocolSpecific_GetValue_Call) Return(_a0 interface{}) *MockDiscoveryOptionProtocolSpecific_GetValue_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDiscoveryOptionProtocolSpecific_GetValue_Call) RunAndReturn(run func() interface{}) *MockDiscoveryOptionProtocolSpecific_GetValue_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDiscoveryOptionProtocolSpecific interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDiscoveryOptionProtocolSpecific creates a new instance of MockDiscoveryOptionProtocolSpecific. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDiscoveryOptionProtocolSpecific(t mockConstructorTestingTNewMockDiscoveryOptionProtocolSpecific) *MockDiscoveryOptionProtocolSpecific {
	mock := &MockDiscoveryOptionProtocolSpecific{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
