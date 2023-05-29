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

package model

import mock "github.com/stretchr/testify/mock"

// MockPlcWriteRequestBuilder is an autogenerated mock type for the PlcWriteRequestBuilder type
type MockPlcWriteRequestBuilder struct {
	mock.Mock
}

type MockPlcWriteRequestBuilder_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPlcWriteRequestBuilder) EXPECT() *MockPlcWriteRequestBuilder_Expecter {
	return &MockPlcWriteRequestBuilder_Expecter{mock: &_m.Mock}
}

// AddTag provides a mock function with given fields: tagName, tag, value
func (_m *MockPlcWriteRequestBuilder) AddTag(tagName string, tag PlcTag, value interface{}) PlcWriteRequestBuilder {
	ret := _m.Called(tagName, tag, value)

	var r0 PlcWriteRequestBuilder
	if rf, ok := ret.Get(0).(func(string, PlcTag, interface{}) PlcWriteRequestBuilder); ok {
		r0 = rf(tagName, tag, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PlcWriteRequestBuilder)
		}
	}

	return r0
}

// MockPlcWriteRequestBuilder_AddTag_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddTag'
type MockPlcWriteRequestBuilder_AddTag_Call struct {
	*mock.Call
}

// AddTag is a helper method to define mock.On call
//   - tagName string
//   - tag PlcTag
//   - value interface{}
func (_e *MockPlcWriteRequestBuilder_Expecter) AddTag(tagName interface{}, tag interface{}, value interface{}) *MockPlcWriteRequestBuilder_AddTag_Call {
	return &MockPlcWriteRequestBuilder_AddTag_Call{Call: _e.mock.On("AddTag", tagName, tag, value)}
}

func (_c *MockPlcWriteRequestBuilder_AddTag_Call) Run(run func(tagName string, tag PlcTag, value interface{})) *MockPlcWriteRequestBuilder_AddTag_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(PlcTag), args[2].(interface{}))
	})
	return _c
}

func (_c *MockPlcWriteRequestBuilder_AddTag_Call) Return(_a0 PlcWriteRequestBuilder) *MockPlcWriteRequestBuilder_AddTag_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcWriteRequestBuilder_AddTag_Call) RunAndReturn(run func(string, PlcTag, interface{}) PlcWriteRequestBuilder) *MockPlcWriteRequestBuilder_AddTag_Call {
	_c.Call.Return(run)
	return _c
}

// AddTagAddress provides a mock function with given fields: tagName, tagAddress, value
func (_m *MockPlcWriteRequestBuilder) AddTagAddress(tagName string, tagAddress string, value interface{}) PlcWriteRequestBuilder {
	ret := _m.Called(tagName, tagAddress, value)

	var r0 PlcWriteRequestBuilder
	if rf, ok := ret.Get(0).(func(string, string, interface{}) PlcWriteRequestBuilder); ok {
		r0 = rf(tagName, tagAddress, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PlcWriteRequestBuilder)
		}
	}

	return r0
}

// MockPlcWriteRequestBuilder_AddTagAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddTagAddress'
type MockPlcWriteRequestBuilder_AddTagAddress_Call struct {
	*mock.Call
}

// AddTagAddress is a helper method to define mock.On call
//   - tagName string
//   - tagAddress string
//   - value interface{}
func (_e *MockPlcWriteRequestBuilder_Expecter) AddTagAddress(tagName interface{}, tagAddress interface{}, value interface{}) *MockPlcWriteRequestBuilder_AddTagAddress_Call {
	return &MockPlcWriteRequestBuilder_AddTagAddress_Call{Call: _e.mock.On("AddTagAddress", tagName, tagAddress, value)}
}

func (_c *MockPlcWriteRequestBuilder_AddTagAddress_Call) Run(run func(tagName string, tagAddress string, value interface{})) *MockPlcWriteRequestBuilder_AddTagAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *MockPlcWriteRequestBuilder_AddTagAddress_Call) Return(_a0 PlcWriteRequestBuilder) *MockPlcWriteRequestBuilder_AddTagAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcWriteRequestBuilder_AddTagAddress_Call) RunAndReturn(run func(string, string, interface{}) PlcWriteRequestBuilder) *MockPlcWriteRequestBuilder_AddTagAddress_Call {
	_c.Call.Return(run)
	return _c
}

// Build provides a mock function with given fields:
func (_m *MockPlcWriteRequestBuilder) Build() (PlcWriteRequest, error) {
	ret := _m.Called()

	var r0 PlcWriteRequest
	var r1 error
	if rf, ok := ret.Get(0).(func() (PlcWriteRequest, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() PlcWriteRequest); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PlcWriteRequest)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockPlcWriteRequestBuilder_Build_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Build'
type MockPlcWriteRequestBuilder_Build_Call struct {
	*mock.Call
}

// Build is a helper method to define mock.On call
func (_e *MockPlcWriteRequestBuilder_Expecter) Build() *MockPlcWriteRequestBuilder_Build_Call {
	return &MockPlcWriteRequestBuilder_Build_Call{Call: _e.mock.On("Build")}
}

func (_c *MockPlcWriteRequestBuilder_Build_Call) Run(run func()) *MockPlcWriteRequestBuilder_Build_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockPlcWriteRequestBuilder_Build_Call) Return(_a0 PlcWriteRequest, _a1 error) *MockPlcWriteRequestBuilder_Build_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockPlcWriteRequestBuilder_Build_Call) RunAndReturn(run func() (PlcWriteRequest, error)) *MockPlcWriteRequestBuilder_Build_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockPlcWriteRequestBuilder interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockPlcWriteRequestBuilder creates a new instance of MockPlcWriteRequestBuilder. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockPlcWriteRequestBuilder(t mockConstructorTestingTNewMockPlcWriteRequestBuilder) *MockPlcWriteRequestBuilder {
	mock := &MockPlcWriteRequestBuilder{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
