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

package utils

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockLengthAware is an autogenerated mock type for the LengthAware type
type MockLengthAware struct {
	mock.Mock
}

type MockLengthAware_Expecter struct {
	mock *mock.Mock
}

func (_m *MockLengthAware) EXPECT() *MockLengthAware_Expecter {
	return &MockLengthAware_Expecter{mock: &_m.Mock}
}

// GetLengthInBits provides a mock function with given fields: ctx
func (_m *MockLengthAware) GetLengthInBits(ctx context.Context) uint16 {
	ret := _m.Called(ctx)

	var r0 uint16
	if rf, ok := ret.Get(0).(func(context.Context) uint16); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// MockLengthAware_GetLengthInBits_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLengthInBits'
type MockLengthAware_GetLengthInBits_Call struct {
	*mock.Call
}

// GetLengthInBits is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockLengthAware_Expecter) GetLengthInBits(ctx interface{}) *MockLengthAware_GetLengthInBits_Call {
	return &MockLengthAware_GetLengthInBits_Call{Call: _e.mock.On("GetLengthInBits", ctx)}
}

func (_c *MockLengthAware_GetLengthInBits_Call) Run(run func(ctx context.Context)) *MockLengthAware_GetLengthInBits_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockLengthAware_GetLengthInBits_Call) Return(_a0 uint16) *MockLengthAware_GetLengthInBits_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLengthAware_GetLengthInBits_Call) RunAndReturn(run func(context.Context) uint16) *MockLengthAware_GetLengthInBits_Call {
	_c.Call.Return(run)
	return _c
}

// GetLengthInBytes provides a mock function with given fields: ctx
func (_m *MockLengthAware) GetLengthInBytes(ctx context.Context) uint16 {
	ret := _m.Called(ctx)

	var r0 uint16
	if rf, ok := ret.Get(0).(func(context.Context) uint16); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(uint16)
	}

	return r0
}

// MockLengthAware_GetLengthInBytes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLengthInBytes'
type MockLengthAware_GetLengthInBytes_Call struct {
	*mock.Call
}

// GetLengthInBytes is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockLengthAware_Expecter) GetLengthInBytes(ctx interface{}) *MockLengthAware_GetLengthInBytes_Call {
	return &MockLengthAware_GetLengthInBytes_Call{Call: _e.mock.On("GetLengthInBytes", ctx)}
}

func (_c *MockLengthAware_GetLengthInBytes_Call) Run(run func(ctx context.Context)) *MockLengthAware_GetLengthInBytes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockLengthAware_GetLengthInBytes_Call) Return(_a0 uint16) *MockLengthAware_GetLengthInBytes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockLengthAware_GetLengthInBytes_Call) RunAndReturn(run func(context.Context) uint16) *MockLengthAware_GetLengthInBytes_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockLengthAware interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockLengthAware creates a new instance of MockLengthAware. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockLengthAware(t mockConstructorTestingTNewMockLengthAware) *MockLengthAware {
	mock := &MockLengthAware{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
