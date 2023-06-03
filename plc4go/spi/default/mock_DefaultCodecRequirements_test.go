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

package _default

import (
	spi "github.com/apache/plc4x/plc4go/spi"
	mock "github.com/stretchr/testify/mock"
)

// MockDefaultCodecRequirements is an autogenerated mock type for the DefaultCodecRequirements type
type MockDefaultCodecRequirements struct {
	mock.Mock
}

type MockDefaultCodecRequirements_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDefaultCodecRequirements) EXPECT() *MockDefaultCodecRequirements_Expecter {
	return &MockDefaultCodecRequirements_Expecter{mock: &_m.Mock}
}

// GetCodec provides a mock function with given fields:
func (_m *MockDefaultCodecRequirements) GetCodec() spi.MessageCodec {
	ret := _m.Called()

	var r0 spi.MessageCodec
	if rf, ok := ret.Get(0).(func() spi.MessageCodec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(spi.MessageCodec)
		}
	}

	return r0
}

// MockDefaultCodecRequirements_GetCodec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCodec'
type MockDefaultCodecRequirements_GetCodec_Call struct {
	*mock.Call
}

// GetCodec is a helper method to define mock.On call
func (_e *MockDefaultCodecRequirements_Expecter) GetCodec() *MockDefaultCodecRequirements_GetCodec_Call {
	return &MockDefaultCodecRequirements_GetCodec_Call{Call: _e.mock.On("GetCodec")}
}

func (_c *MockDefaultCodecRequirements_GetCodec_Call) Run(run func()) *MockDefaultCodecRequirements_GetCodec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDefaultCodecRequirements_GetCodec_Call) Return(_a0 spi.MessageCodec) *MockDefaultCodecRequirements_GetCodec_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDefaultCodecRequirements_GetCodec_Call) RunAndReturn(run func() spi.MessageCodec) *MockDefaultCodecRequirements_GetCodec_Call {
	_c.Call.Return(run)
	return _c
}

// Receive provides a mock function with given fields:
func (_m *MockDefaultCodecRequirements) Receive() (spi.Message, error) {
	ret := _m.Called()

	var r0 spi.Message
	var r1 error
	if rf, ok := ret.Get(0).(func() (spi.Message, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() spi.Message); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(spi.Message)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDefaultCodecRequirements_Receive_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Receive'
type MockDefaultCodecRequirements_Receive_Call struct {
	*mock.Call
}

// Receive is a helper method to define mock.On call
func (_e *MockDefaultCodecRequirements_Expecter) Receive() *MockDefaultCodecRequirements_Receive_Call {
	return &MockDefaultCodecRequirements_Receive_Call{Call: _e.mock.On("Receive")}
}

func (_c *MockDefaultCodecRequirements_Receive_Call) Run(run func()) *MockDefaultCodecRequirements_Receive_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockDefaultCodecRequirements_Receive_Call) Return(_a0 spi.Message, _a1 error) *MockDefaultCodecRequirements_Receive_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDefaultCodecRequirements_Receive_Call) RunAndReturn(run func() (spi.Message, error)) *MockDefaultCodecRequirements_Receive_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: message
func (_m *MockDefaultCodecRequirements) Send(message spi.Message) error {
	ret := _m.Called(message)

	var r0 error
	if rf, ok := ret.Get(0).(func(spi.Message) error); ok {
		r0 = rf(message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDefaultCodecRequirements_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MockDefaultCodecRequirements_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - message spi.Message
func (_e *MockDefaultCodecRequirements_Expecter) Send(message interface{}) *MockDefaultCodecRequirements_Send_Call {
	return &MockDefaultCodecRequirements_Send_Call{Call: _e.mock.On("Send", message)}
}

func (_c *MockDefaultCodecRequirements_Send_Call) Run(run func(message spi.Message)) *MockDefaultCodecRequirements_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(spi.Message))
	})
	return _c
}

func (_c *MockDefaultCodecRequirements_Send_Call) Return(_a0 error) *MockDefaultCodecRequirements_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDefaultCodecRequirements_Send_Call) RunAndReturn(run func(spi.Message) error) *MockDefaultCodecRequirements_Send_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockDefaultCodecRequirements interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDefaultCodecRequirements creates a new instance of MockDefaultCodecRequirements. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDefaultCodecRequirements(t mockConstructorTestingTNewMockDefaultCodecRequirements) *MockDefaultCodecRequirements {
	mock := &MockDefaultCodecRequirements{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
