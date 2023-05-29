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

package spi

import mock "github.com/stretchr/testify/mock"

// MockAcceptsMessage is an autogenerated mock type for the AcceptsMessage type
type MockAcceptsMessage struct {
	mock.Mock
}

type MockAcceptsMessage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAcceptsMessage) EXPECT() *MockAcceptsMessage_Expecter {
	return &MockAcceptsMessage_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: message
func (_m *MockAcceptsMessage) Execute(message Message) bool {
	ret := _m.Called(message)

	var r0 bool
	if rf, ok := ret.Get(0).(func(Message) bool); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockAcceptsMessage_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockAcceptsMessage_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - message Message
func (_e *MockAcceptsMessage_Expecter) Execute(message interface{}) *MockAcceptsMessage_Execute_Call {
	return &MockAcceptsMessage_Execute_Call{Call: _e.mock.On("Execute", message)}
}

func (_c *MockAcceptsMessage_Execute_Call) Run(run func(message Message)) *MockAcceptsMessage_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Message))
	})
	return _c
}

func (_c *MockAcceptsMessage_Execute_Call) Return(_a0 bool) *MockAcceptsMessage_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAcceptsMessage_Execute_Call) RunAndReturn(run func(Message) bool) *MockAcceptsMessage_Execute_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockAcceptsMessage interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAcceptsMessage creates a new instance of MockAcceptsMessage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAcceptsMessage(t mockConstructorTestingTNewMockAcceptsMessage) *MockAcceptsMessage {
	mock := &MockAcceptsMessage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
