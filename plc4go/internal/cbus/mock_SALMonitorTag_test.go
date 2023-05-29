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
	apimodel "github.com/apache/plc4x/plc4go/pkg/api/model"
	mock "github.com/stretchr/testify/mock"

	model "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"

	values "github.com/apache/plc4x/plc4go/pkg/api/values"
)

// MockSALMonitorTag is an autogenerated mock type for the SALMonitorTag type
type MockSALMonitorTag struct {
	mock.Mock
}

type MockSALMonitorTag_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSALMonitorTag) EXPECT() *MockSALMonitorTag_Expecter {
	return &MockSALMonitorTag_Expecter{mock: &_m.Mock}
}

// GetAddressString provides a mock function with given fields:
func (_m *MockSALMonitorTag) GetAddressString() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockSALMonitorTag_GetAddressString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAddressString'
type MockSALMonitorTag_GetAddressString_Call struct {
	*mock.Call
}

// GetAddressString is a helper method to define mock.On call
func (_e *MockSALMonitorTag_Expecter) GetAddressString() *MockSALMonitorTag_GetAddressString_Call {
	return &MockSALMonitorTag_GetAddressString_Call{Call: _e.mock.On("GetAddressString")}
}

func (_c *MockSALMonitorTag_GetAddressString_Call) Run(run func()) *MockSALMonitorTag_GetAddressString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSALMonitorTag_GetAddressString_Call) Return(_a0 string) *MockSALMonitorTag_GetAddressString_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSALMonitorTag_GetAddressString_Call) RunAndReturn(run func() string) *MockSALMonitorTag_GetAddressString_Call {
	_c.Call.Return(run)
	return _c
}

// GetApplication provides a mock function with given fields:
func (_m *MockSALMonitorTag) GetApplication() *model.ApplicationIdContainer {
	ret := _m.Called()

	var r0 *model.ApplicationIdContainer
	if rf, ok := ret.Get(0).(func() *model.ApplicationIdContainer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationIdContainer)
		}
	}

	return r0
}

// MockSALMonitorTag_GetApplication_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetApplication'
type MockSALMonitorTag_GetApplication_Call struct {
	*mock.Call
}

// GetApplication is a helper method to define mock.On call
func (_e *MockSALMonitorTag_Expecter) GetApplication() *MockSALMonitorTag_GetApplication_Call {
	return &MockSALMonitorTag_GetApplication_Call{Call: _e.mock.On("GetApplication")}
}

func (_c *MockSALMonitorTag_GetApplication_Call) Run(run func()) *MockSALMonitorTag_GetApplication_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSALMonitorTag_GetApplication_Call) Return(_a0 *model.ApplicationIdContainer) *MockSALMonitorTag_GetApplication_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSALMonitorTag_GetApplication_Call) RunAndReturn(run func() *model.ApplicationIdContainer) *MockSALMonitorTag_GetApplication_Call {
	_c.Call.Return(run)
	return _c
}

// GetArrayInfo provides a mock function with given fields:
func (_m *MockSALMonitorTag) GetArrayInfo() []apimodel.ArrayInfo {
	ret := _m.Called()

	var r0 []apimodel.ArrayInfo
	if rf, ok := ret.Get(0).(func() []apimodel.ArrayInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]apimodel.ArrayInfo)
		}
	}

	return r0
}

// MockSALMonitorTag_GetArrayInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetArrayInfo'
type MockSALMonitorTag_GetArrayInfo_Call struct {
	*mock.Call
}

// GetArrayInfo is a helper method to define mock.On call
func (_e *MockSALMonitorTag_Expecter) GetArrayInfo() *MockSALMonitorTag_GetArrayInfo_Call {
	return &MockSALMonitorTag_GetArrayInfo_Call{Call: _e.mock.On("GetArrayInfo")}
}

func (_c *MockSALMonitorTag_GetArrayInfo_Call) Run(run func()) *MockSALMonitorTag_GetArrayInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSALMonitorTag_GetArrayInfo_Call) Return(_a0 []apimodel.ArrayInfo) *MockSALMonitorTag_GetArrayInfo_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSALMonitorTag_GetArrayInfo_Call) RunAndReturn(run func() []apimodel.ArrayInfo) *MockSALMonitorTag_GetArrayInfo_Call {
	_c.Call.Return(run)
	return _c
}

// GetTagType provides a mock function with given fields:
func (_m *MockSALMonitorTag) GetTagType() TagType {
	ret := _m.Called()

	var r0 TagType
	if rf, ok := ret.Get(0).(func() TagType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(TagType)
	}

	return r0
}

// MockSALMonitorTag_GetTagType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTagType'
type MockSALMonitorTag_GetTagType_Call struct {
	*mock.Call
}

// GetTagType is a helper method to define mock.On call
func (_e *MockSALMonitorTag_Expecter) GetTagType() *MockSALMonitorTag_GetTagType_Call {
	return &MockSALMonitorTag_GetTagType_Call{Call: _e.mock.On("GetTagType")}
}

func (_c *MockSALMonitorTag_GetTagType_Call) Run(run func()) *MockSALMonitorTag_GetTagType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSALMonitorTag_GetTagType_Call) Return(_a0 TagType) *MockSALMonitorTag_GetTagType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSALMonitorTag_GetTagType_Call) RunAndReturn(run func() TagType) *MockSALMonitorTag_GetTagType_Call {
	_c.Call.Return(run)
	return _c
}

// GetUnitAddress provides a mock function with given fields:
func (_m *MockSALMonitorTag) GetUnitAddress() model.UnitAddress {
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

// MockSALMonitorTag_GetUnitAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUnitAddress'
type MockSALMonitorTag_GetUnitAddress_Call struct {
	*mock.Call
}

// GetUnitAddress is a helper method to define mock.On call
func (_e *MockSALMonitorTag_Expecter) GetUnitAddress() *MockSALMonitorTag_GetUnitAddress_Call {
	return &MockSALMonitorTag_GetUnitAddress_Call{Call: _e.mock.On("GetUnitAddress")}
}

func (_c *MockSALMonitorTag_GetUnitAddress_Call) Run(run func()) *MockSALMonitorTag_GetUnitAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSALMonitorTag_GetUnitAddress_Call) Return(_a0 model.UnitAddress) *MockSALMonitorTag_GetUnitAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSALMonitorTag_GetUnitAddress_Call) RunAndReturn(run func() model.UnitAddress) *MockSALMonitorTag_GetUnitAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetValueType provides a mock function with given fields:
func (_m *MockSALMonitorTag) GetValueType() values.PlcValueType {
	ret := _m.Called()

	var r0 values.PlcValueType
	if rf, ok := ret.Get(0).(func() values.PlcValueType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(values.PlcValueType)
	}

	return r0
}

// MockSALMonitorTag_GetValueType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetValueType'
type MockSALMonitorTag_GetValueType_Call struct {
	*mock.Call
}

// GetValueType is a helper method to define mock.On call
func (_e *MockSALMonitorTag_Expecter) GetValueType() *MockSALMonitorTag_GetValueType_Call {
	return &MockSALMonitorTag_GetValueType_Call{Call: _e.mock.On("GetValueType")}
}

func (_c *MockSALMonitorTag_GetValueType_Call) Run(run func()) *MockSALMonitorTag_GetValueType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSALMonitorTag_GetValueType_Call) Return(_a0 values.PlcValueType) *MockSALMonitorTag_GetValueType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSALMonitorTag_GetValueType_Call) RunAndReturn(run func() values.PlcValueType) *MockSALMonitorTag_GetValueType_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockSALMonitorTag interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockSALMonitorTag creates a new instance of MockSALMonitorTag. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockSALMonitorTag(t mockConstructorTestingTNewMockSALMonitorTag) *MockSALMonitorTag {
	mock := &MockSALMonitorTag{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
