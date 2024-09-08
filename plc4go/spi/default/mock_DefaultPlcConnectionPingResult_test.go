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

package _default

import mock "github.com/stretchr/testify/mock"

// MockDefaultPlcConnectionPingResult is an autogenerated mock type for the DefaultPlcConnectionPingResult type
type MockDefaultPlcConnectionPingResult struct {
	mock.Mock
}

type MockDefaultPlcConnectionPingResult_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDefaultPlcConnectionPingResult) EXPECT() *MockDefaultPlcConnectionPingResult_Expecter {
	return &MockDefaultPlcConnectionPingResult_Expecter{mock: &_m.Mock}
}

// GetErr provides a mock function with given fields:
func (_m *MockDefaultPlcConnectionPingResult) GetErr() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetErr")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDefaultPlcConnectionPingResult_GetErr_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetErr'
type MockDefaultPlcConnectionPingResult_GetErr_Call struct {
	*mock.Call
}

// GetErr is a helper method to define mock.On call
func (_e *MockDefaultPlcConnectionPingResult_Expecter) GetErr() *MockDefaultPlcConnectionPingResult_GetErr_Call {
	return &MockDefaultPlcConnectionPingResult_GetErr_Call{Call: _e.mock.On("GetErr")}
}

func (_c *MockDefaultPlcConnectionPingResult_GetErr_Call) Run(run func()) *MockDefaultPlcConnectionPingResult_GetErr_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDefaultPlcConnectionPingResult_GetErr_Call) Return(_a0 error) *MockDefaultPlcConnectionPingResult_GetErr_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDefaultPlcConnectionPingResult_GetErr_Call) RunAndReturn(run func() error) *MockDefaultPlcConnectionPingResult_GetErr_Call {
	_c.Call.Return(run)
	return _c
}

// String provides a mock function with given fields:
func (_m *MockDefaultPlcConnectionPingResult) String() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for String")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockDefaultPlcConnectionPingResult_String_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'String'
type MockDefaultPlcConnectionPingResult_String_Call struct {
	*mock.Call
}

// String is a helper method to define mock.On call
func (_e *MockDefaultPlcConnectionPingResult_Expecter) String() *MockDefaultPlcConnectionPingResult_String_Call {
	return &MockDefaultPlcConnectionPingResult_String_Call{Call: _e.mock.On("String")}
}

func (_c *MockDefaultPlcConnectionPingResult_String_Call) Run(run func()) *MockDefaultPlcConnectionPingResult_String_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDefaultPlcConnectionPingResult_String_Call) Return(_a0 string) *MockDefaultPlcConnectionPingResult_String_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDefaultPlcConnectionPingResult_String_Call) RunAndReturn(run func() string) *MockDefaultPlcConnectionPingResult_String_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDefaultPlcConnectionPingResult creates a new instance of MockDefaultPlcConnectionPingResult. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDefaultPlcConnectionPingResult(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDefaultPlcConnectionPingResult {
	mock := &MockDefaultPlcConnectionPingResult{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
