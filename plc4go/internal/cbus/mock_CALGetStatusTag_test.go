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

// Code generated by mockery v2.32.0. DO NOT EDIT.

package cbus

import (
	model "github.com/apache/plc4x/plc4go/pkg/api/model"
	mock "github.com/stretchr/testify/mock"

	readwritemodel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"

	values "github.com/apache/plc4x/plc4go/pkg/api/values"
)

// MockCALGetStatusTag is an autogenerated mock type for the CALGetStatusTag type
type MockCALGetStatusTag struct {
	mock.Mock
}

type MockCALGetStatusTag_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCALGetStatusTag) EXPECT() *MockCALGetStatusTag_Expecter {
	return &MockCALGetStatusTag_Expecter{mock: &_m.Mock}
}

// GetAddressString provides a mock function with given fields:
func (_m *MockCALGetStatusTag) GetAddressString() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockCALGetStatusTag_GetAddressString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAddressString'
type MockCALGetStatusTag_GetAddressString_Call struct {
	*mock.Call
}

// GetAddressString is a helper method to define mock.On call
func (_e *MockCALGetStatusTag_Expecter) GetAddressString() *MockCALGetStatusTag_GetAddressString_Call {
	return &MockCALGetStatusTag_GetAddressString_Call{Call: _e.mock.On("GetAddressString")}
}

func (_c *MockCALGetStatusTag_GetAddressString_Call) Run(run func()) *MockCALGetStatusTag_GetAddressString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCALGetStatusTag_GetAddressString_Call) Return(_a0 string) *MockCALGetStatusTag_GetAddressString_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCALGetStatusTag_GetAddressString_Call) RunAndReturn(run func() string) *MockCALGetStatusTag_GetAddressString_Call {
	_c.Call.Return(run)
	return _c
}

// GetArrayInfo provides a mock function with given fields:
func (_m *MockCALGetStatusTag) GetArrayInfo() []model.ArrayInfo {
	ret := _m.Called()

	var r0 []model.ArrayInfo
	if rf, ok := ret.Get(0).(func() []model.ArrayInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.ArrayInfo)
		}
	}

	return r0
}

// MockCALGetStatusTag_GetArrayInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetArrayInfo'
type MockCALGetStatusTag_GetArrayInfo_Call struct {
	*mock.Call
}

// GetArrayInfo is a helper method to define mock.On call
func (_e *MockCALGetStatusTag_Expecter) GetArrayInfo() *MockCALGetStatusTag_GetArrayInfo_Call {
	return &MockCALGetStatusTag_GetArrayInfo_Call{Call: _e.mock.On("GetArrayInfo")}
}

func (_c *MockCALGetStatusTag_GetArrayInfo_Call) Run(run func()) *MockCALGetStatusTag_GetArrayInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCALGetStatusTag_GetArrayInfo_Call) Return(_a0 []model.ArrayInfo) *MockCALGetStatusTag_GetArrayInfo_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCALGetStatusTag_GetArrayInfo_Call) RunAndReturn(run func() []model.ArrayInfo) *MockCALGetStatusTag_GetArrayInfo_Call {
	_c.Call.Return(run)
	return _c
}

// GetBridgeAddresses provides a mock function with given fields:
func (_m *MockCALGetStatusTag) GetBridgeAddresses() []readwritemodel.BridgeAddress {
	ret := _m.Called()

	var r0 []readwritemodel.BridgeAddress
	if rf, ok := ret.Get(0).(func() []readwritemodel.BridgeAddress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]readwritemodel.BridgeAddress)
		}
	}

	return r0
}

// MockCALGetStatusTag_GetBridgeAddresses_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBridgeAddresses'
type MockCALGetStatusTag_GetBridgeAddresses_Call struct {
	*mock.Call
}

// GetBridgeAddresses is a helper method to define mock.On call
func (_e *MockCALGetStatusTag_Expecter) GetBridgeAddresses() *MockCALGetStatusTag_GetBridgeAddresses_Call {
	return &MockCALGetStatusTag_GetBridgeAddresses_Call{Call: _e.mock.On("GetBridgeAddresses")}
}

func (_c *MockCALGetStatusTag_GetBridgeAddresses_Call) Run(run func()) *MockCALGetStatusTag_GetBridgeAddresses_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCALGetStatusTag_GetBridgeAddresses_Call) Return(_a0 []readwritemodel.BridgeAddress) *MockCALGetStatusTag_GetBridgeAddresses_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCALGetStatusTag_GetBridgeAddresses_Call) RunAndReturn(run func() []readwritemodel.BridgeAddress) *MockCALGetStatusTag_GetBridgeAddresses_Call {
	_c.Call.Return(run)
	return _c
}

// GetCount provides a mock function with given fields:
func (_m *MockCALGetStatusTag) GetCount() uint8 {
	ret := _m.Called()

	var r0 uint8
	if rf, ok := ret.Get(0).(func() uint8); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint8)
	}

	return r0
}

// MockCALGetStatusTag_GetCount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCount'
type MockCALGetStatusTag_GetCount_Call struct {
	*mock.Call
}

// GetCount is a helper method to define mock.On call
func (_e *MockCALGetStatusTag_Expecter) GetCount() *MockCALGetStatusTag_GetCount_Call {
	return &MockCALGetStatusTag_GetCount_Call{Call: _e.mock.On("GetCount")}
}

func (_c *MockCALGetStatusTag_GetCount_Call) Run(run func()) *MockCALGetStatusTag_GetCount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCALGetStatusTag_GetCount_Call) Return(_a0 uint8) *MockCALGetStatusTag_GetCount_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCALGetStatusTag_GetCount_Call) RunAndReturn(run func() uint8) *MockCALGetStatusTag_GetCount_Call {
	_c.Call.Return(run)
	return _c
}

// GetParameter provides a mock function with given fields:
func (_m *MockCALGetStatusTag) GetParameter() readwritemodel.Parameter {
	ret := _m.Called()

	var r0 readwritemodel.Parameter
	if rf, ok := ret.Get(0).(func() readwritemodel.Parameter); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(readwritemodel.Parameter)
	}

	return r0
}

// MockCALGetStatusTag_GetParameter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetParameter'
type MockCALGetStatusTag_GetParameter_Call struct {
	*mock.Call
}

// GetParameter is a helper method to define mock.On call
func (_e *MockCALGetStatusTag_Expecter) GetParameter() *MockCALGetStatusTag_GetParameter_Call {
	return &MockCALGetStatusTag_GetParameter_Call{Call: _e.mock.On("GetParameter")}
}

func (_c *MockCALGetStatusTag_GetParameter_Call) Run(run func()) *MockCALGetStatusTag_GetParameter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCALGetStatusTag_GetParameter_Call) Return(_a0 readwritemodel.Parameter) *MockCALGetStatusTag_GetParameter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCALGetStatusTag_GetParameter_Call) RunAndReturn(run func() readwritemodel.Parameter) *MockCALGetStatusTag_GetParameter_Call {
	_c.Call.Return(run)
	return _c
}

// GetTagType provides a mock function with given fields:
func (_m *MockCALGetStatusTag) GetTagType() TagType {
	ret := _m.Called()

	var r0 TagType
	if rf, ok := ret.Get(0).(func() TagType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(TagType)
	}

	return r0
}

// MockCALGetStatusTag_GetTagType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTagType'
type MockCALGetStatusTag_GetTagType_Call struct {
	*mock.Call
}

// GetTagType is a helper method to define mock.On call
func (_e *MockCALGetStatusTag_Expecter) GetTagType() *MockCALGetStatusTag_GetTagType_Call {
	return &MockCALGetStatusTag_GetTagType_Call{Call: _e.mock.On("GetTagType")}
}

func (_c *MockCALGetStatusTag_GetTagType_Call) Run(run func()) *MockCALGetStatusTag_GetTagType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCALGetStatusTag_GetTagType_Call) Return(_a0 TagType) *MockCALGetStatusTag_GetTagType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCALGetStatusTag_GetTagType_Call) RunAndReturn(run func() TagType) *MockCALGetStatusTag_GetTagType_Call {
	_c.Call.Return(run)
	return _c
}

// GetUnitAddress provides a mock function with given fields:
func (_m *MockCALGetStatusTag) GetUnitAddress() readwritemodel.UnitAddress {
	ret := _m.Called()

	var r0 readwritemodel.UnitAddress
	if rf, ok := ret.Get(0).(func() readwritemodel.UnitAddress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(readwritemodel.UnitAddress)
		}
	}

	return r0
}

// MockCALGetStatusTag_GetUnitAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUnitAddress'
type MockCALGetStatusTag_GetUnitAddress_Call struct {
	*mock.Call
}

// GetUnitAddress is a helper method to define mock.On call
func (_e *MockCALGetStatusTag_Expecter) GetUnitAddress() *MockCALGetStatusTag_GetUnitAddress_Call {
	return &MockCALGetStatusTag_GetUnitAddress_Call{Call: _e.mock.On("GetUnitAddress")}
}

func (_c *MockCALGetStatusTag_GetUnitAddress_Call) Run(run func()) *MockCALGetStatusTag_GetUnitAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCALGetStatusTag_GetUnitAddress_Call) Return(_a0 readwritemodel.UnitAddress) *MockCALGetStatusTag_GetUnitAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCALGetStatusTag_GetUnitAddress_Call) RunAndReturn(run func() readwritemodel.UnitAddress) *MockCALGetStatusTag_GetUnitAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetValueType provides a mock function with given fields:
func (_m *MockCALGetStatusTag) GetValueType() values.PlcValueType {
	ret := _m.Called()

	var r0 values.PlcValueType
	if rf, ok := ret.Get(0).(func() values.PlcValueType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(values.PlcValueType)
	}

	return r0
}

// MockCALGetStatusTag_GetValueType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetValueType'
type MockCALGetStatusTag_GetValueType_Call struct {
	*mock.Call
}

// GetValueType is a helper method to define mock.On call
func (_e *MockCALGetStatusTag_Expecter) GetValueType() *MockCALGetStatusTag_GetValueType_Call {
	return &MockCALGetStatusTag_GetValueType_Call{Call: _e.mock.On("GetValueType")}
}

func (_c *MockCALGetStatusTag_GetValueType_Call) Run(run func()) *MockCALGetStatusTag_GetValueType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCALGetStatusTag_GetValueType_Call) Return(_a0 values.PlcValueType) *MockCALGetStatusTag_GetValueType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCALGetStatusTag_GetValueType_Call) RunAndReturn(run func() values.PlcValueType) *MockCALGetStatusTag_GetValueType_Call {
	_c.Call.Return(run)
	return _c
}

// String provides a mock function with given fields:
func (_m *MockCALGetStatusTag) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockCALGetStatusTag_String_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'String'
type MockCALGetStatusTag_String_Call struct {
	*mock.Call
}

// String is a helper method to define mock.On call
func (_e *MockCALGetStatusTag_Expecter) String() *MockCALGetStatusTag_String_Call {
	return &MockCALGetStatusTag_String_Call{Call: _e.mock.On("String")}
}

func (_c *MockCALGetStatusTag_String_Call) Run(run func()) *MockCALGetStatusTag_String_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCALGetStatusTag_String_Call) Return(_a0 string) *MockCALGetStatusTag_String_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCALGetStatusTag_String_Call) RunAndReturn(run func() string) *MockCALGetStatusTag_String_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockCALGetStatusTag creates a new instance of MockCALGetStatusTag. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockCALGetStatusTag(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockCALGetStatusTag {
	mock := &MockCALGetStatusTag{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
