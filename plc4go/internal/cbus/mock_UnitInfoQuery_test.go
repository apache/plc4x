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

package cbus

import (
	model "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	mock "github.com/stretchr/testify/mock"
)

// MockUnitInfoQuery is an autogenerated mock type for the UnitInfoQuery type
type MockUnitInfoQuery struct {
	mock.Mock
}

type MockUnitInfoQuery_Expecter struct {
	mock *mock.Mock
}

func (_m *MockUnitInfoQuery) EXPECT() *MockUnitInfoQuery_Expecter {
	return &MockUnitInfoQuery_Expecter{mock: &_m.Mock}
}

// GetAttribute provides a mock function with given fields:
func (_m *MockUnitInfoQuery) GetAttribute() *model.Attribute {
	ret := _m.Called()

	var r0 *model.Attribute
	if rf, ok := ret.Get(0).(func() *model.Attribute); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Attribute)
		}
	}

	return r0
}

// MockUnitInfoQuery_GetAttribute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAttribute'
type MockUnitInfoQuery_GetAttribute_Call struct {
	*mock.Call
}

// GetAttribute is a helper method to define mock.On call
func (_e *MockUnitInfoQuery_Expecter) GetAttribute() *MockUnitInfoQuery_GetAttribute_Call {
	return &MockUnitInfoQuery_GetAttribute_Call{Call: _e.mock.On("GetAttribute")}
}

func (_c *MockUnitInfoQuery_GetAttribute_Call) Run(run func()) *MockUnitInfoQuery_GetAttribute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUnitInfoQuery_GetAttribute_Call) Return(_a0 *model.Attribute) *MockUnitInfoQuery_GetAttribute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUnitInfoQuery_GetAttribute_Call) RunAndReturn(run func() *model.Attribute) *MockUnitInfoQuery_GetAttribute_Call {
	_c.Call.Return(run)
	return _c
}

// GetQueryString provides a mock function with given fields:
func (_m *MockUnitInfoQuery) GetQueryString() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockUnitInfoQuery_GetQueryString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQueryString'
type MockUnitInfoQuery_GetQueryString_Call struct {
	*mock.Call
}

// GetQueryString is a helper method to define mock.On call
func (_e *MockUnitInfoQuery_Expecter) GetQueryString() *MockUnitInfoQuery_GetQueryString_Call {
	return &MockUnitInfoQuery_GetQueryString_Call{Call: _e.mock.On("GetQueryString")}
}

func (_c *MockUnitInfoQuery_GetQueryString_Call) Run(run func()) *MockUnitInfoQuery_GetQueryString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUnitInfoQuery_GetQueryString_Call) Return(_a0 string) *MockUnitInfoQuery_GetQueryString_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUnitInfoQuery_GetQueryString_Call) RunAndReturn(run func() string) *MockUnitInfoQuery_GetQueryString_Call {
	_c.Call.Return(run)
	return _c
}

// GetUnitAddress provides a mock function with given fields:
func (_m *MockUnitInfoQuery) GetUnitAddress() model.UnitAddress {
	ret := _m.Called()

	var r0 model.UnitAddress
	if rf, ok := ret.Get(0).(func() model.UnitAddress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.UnitAddress)
		}
	}

	return r0
}

// MockUnitInfoQuery_GetUnitAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUnitAddress'
type MockUnitInfoQuery_GetUnitAddress_Call struct {
	*mock.Call
}

// GetUnitAddress is a helper method to define mock.On call
func (_e *MockUnitInfoQuery_Expecter) GetUnitAddress() *MockUnitInfoQuery_GetUnitAddress_Call {
	return &MockUnitInfoQuery_GetUnitAddress_Call{Call: _e.mock.On("GetUnitAddress")}
}

func (_c *MockUnitInfoQuery_GetUnitAddress_Call) Run(run func()) *MockUnitInfoQuery_GetUnitAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockUnitInfoQuery_GetUnitAddress_Call) Return(_a0 model.UnitAddress) *MockUnitInfoQuery_GetUnitAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockUnitInfoQuery_GetUnitAddress_Call) RunAndReturn(run func() model.UnitAddress) *MockUnitInfoQuery_GetUnitAddress_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockUnitInfoQuery interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockUnitInfoQuery creates a new instance of MockUnitInfoQuery. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockUnitInfoQuery(t mockConstructorTestingTNewMockUnitInfoQuery) *MockUnitInfoQuery {
	mock := &MockUnitInfoQuery{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
