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

package knxnetip

import (
	model "github.com/apache/plc4x/plc4go/pkg/api/model"
	mock "github.com/stretchr/testify/mock"

	readwritemodel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"

	values "github.com/apache/plc4x/plc4go/pkg/api/values"
)

// MockGroupAddressTag is an autogenerated mock type for the GroupAddressTag type
type MockGroupAddressTag struct {
	mock.Mock
}

type MockGroupAddressTag_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGroupAddressTag) EXPECT() *MockGroupAddressTag_Expecter {
	return &MockGroupAddressTag_Expecter{mock: &_m.Mock}
}

// GetAddressString provides a mock function with given fields:
func (_m *MockGroupAddressTag) GetAddressString() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockGroupAddressTag_GetAddressString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAddressString'
type MockGroupAddressTag_GetAddressString_Call struct {
	*mock.Call
}

// GetAddressString is a helper method to define mock.On call
func (_e *MockGroupAddressTag_Expecter) GetAddressString() *MockGroupAddressTag_GetAddressString_Call {
	return &MockGroupAddressTag_GetAddressString_Call{Call: _e.mock.On("GetAddressString")}
}

func (_c *MockGroupAddressTag_GetAddressString_Call) Run(run func()) *MockGroupAddressTag_GetAddressString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockGroupAddressTag_GetAddressString_Call) Return(_a0 string) *MockGroupAddressTag_GetAddressString_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGroupAddressTag_GetAddressString_Call) RunAndReturn(run func() string) *MockGroupAddressTag_GetAddressString_Call {
	_c.Call.Return(run)
	return _c
}

// GetArrayInfo provides a mock function with given fields:
func (_m *MockGroupAddressTag) GetArrayInfo() []model.ArrayInfo {
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

// MockGroupAddressTag_GetArrayInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetArrayInfo'
type MockGroupAddressTag_GetArrayInfo_Call struct {
	*mock.Call
}

// GetArrayInfo is a helper method to define mock.On call
func (_e *MockGroupAddressTag_Expecter) GetArrayInfo() *MockGroupAddressTag_GetArrayInfo_Call {
	return &MockGroupAddressTag_GetArrayInfo_Call{Call: _e.mock.On("GetArrayInfo")}
}

func (_c *MockGroupAddressTag_GetArrayInfo_Call) Run(run func()) *MockGroupAddressTag_GetArrayInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockGroupAddressTag_GetArrayInfo_Call) Return(_a0 []model.ArrayInfo) *MockGroupAddressTag_GetArrayInfo_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGroupAddressTag_GetArrayInfo_Call) RunAndReturn(run func() []model.ArrayInfo) *MockGroupAddressTag_GetArrayInfo_Call {
	_c.Call.Return(run)
	return _c
}

// GetTagType provides a mock function with given fields:
func (_m *MockGroupAddressTag) GetTagType() *readwritemodel.KnxDatapointType {
	ret := _m.Called()

	var r0 *readwritemodel.KnxDatapointType
	if rf, ok := ret.Get(0).(func() *readwritemodel.KnxDatapointType); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*readwritemodel.KnxDatapointType)
		}
	}

	return r0
}

// MockGroupAddressTag_GetTagType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTagType'
type MockGroupAddressTag_GetTagType_Call struct {
	*mock.Call
}

// GetTagType is a helper method to define mock.On call
func (_e *MockGroupAddressTag_Expecter) GetTagType() *MockGroupAddressTag_GetTagType_Call {
	return &MockGroupAddressTag_GetTagType_Call{Call: _e.mock.On("GetTagType")}
}

func (_c *MockGroupAddressTag_GetTagType_Call) Run(run func()) *MockGroupAddressTag_GetTagType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockGroupAddressTag_GetTagType_Call) Return(_a0 *readwritemodel.KnxDatapointType) *MockGroupAddressTag_GetTagType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGroupAddressTag_GetTagType_Call) RunAndReturn(run func() *readwritemodel.KnxDatapointType) *MockGroupAddressTag_GetTagType_Call {
	_c.Call.Return(run)
	return _c
}

// GetValueType provides a mock function with given fields:
func (_m *MockGroupAddressTag) GetValueType() values.PlcValueType {
	ret := _m.Called()

	var r0 values.PlcValueType
	if rf, ok := ret.Get(0).(func() values.PlcValueType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(values.PlcValueType)
	}

	return r0
}

// MockGroupAddressTag_GetValueType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetValueType'
type MockGroupAddressTag_GetValueType_Call struct {
	*mock.Call
}

// GetValueType is a helper method to define mock.On call
func (_e *MockGroupAddressTag_Expecter) GetValueType() *MockGroupAddressTag_GetValueType_Call {
	return &MockGroupAddressTag_GetValueType_Call{Call: _e.mock.On("GetValueType")}
}

func (_c *MockGroupAddressTag_GetValueType_Call) Run(run func()) *MockGroupAddressTag_GetValueType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockGroupAddressTag_GetValueType_Call) Return(_a0 values.PlcValueType) *MockGroupAddressTag_GetValueType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGroupAddressTag_GetValueType_Call) RunAndReturn(run func() values.PlcValueType) *MockGroupAddressTag_GetValueType_Call {
	_c.Call.Return(run)
	return _c
}

// IsPatternTag provides a mock function with given fields:
func (_m *MockGroupAddressTag) IsPatternTag() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockGroupAddressTag_IsPatternTag_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsPatternTag'
type MockGroupAddressTag_IsPatternTag_Call struct {
	*mock.Call
}

// IsPatternTag is a helper method to define mock.On call
func (_e *MockGroupAddressTag_Expecter) IsPatternTag() *MockGroupAddressTag_IsPatternTag_Call {
	return &MockGroupAddressTag_IsPatternTag_Call{Call: _e.mock.On("IsPatternTag")}
}

func (_c *MockGroupAddressTag_IsPatternTag_Call) Run(run func()) *MockGroupAddressTag_IsPatternTag_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockGroupAddressTag_IsPatternTag_Call) Return(_a0 bool) *MockGroupAddressTag_IsPatternTag_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGroupAddressTag_IsPatternTag_Call) RunAndReturn(run func() bool) *MockGroupAddressTag_IsPatternTag_Call {
	_c.Call.Return(run)
	return _c
}

// matches provides a mock function with given fields: knxGroupAddress
func (_m *MockGroupAddressTag) matches(knxGroupAddress readwritemodel.KnxGroupAddress) bool {
	ret := _m.Called(knxGroupAddress)

	var r0 bool
	if rf, ok := ret.Get(0).(func(readwritemodel.KnxGroupAddress) bool); ok {
		r0 = rf(knxGroupAddress)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockGroupAddressTag_matches_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'matches'
type MockGroupAddressTag_matches_Call struct {
	*mock.Call
}

// matches is a helper method to define mock.On call
//   - knxGroupAddress readwritemodel.KnxGroupAddress
func (_e *MockGroupAddressTag_Expecter) matches(knxGroupAddress interface{}) *MockGroupAddressTag_matches_Call {
	return &MockGroupAddressTag_matches_Call{Call: _e.mock.On("matches", knxGroupAddress)}
}

func (_c *MockGroupAddressTag_matches_Call) Run(run func(knxGroupAddress readwritemodel.KnxGroupAddress)) *MockGroupAddressTag_matches_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(readwritemodel.KnxGroupAddress))
	})
	return _c
}

func (_c *MockGroupAddressTag_matches_Call) Return(_a0 bool) *MockGroupAddressTag_matches_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGroupAddressTag_matches_Call) RunAndReturn(run func(readwritemodel.KnxGroupAddress) bool) *MockGroupAddressTag_matches_Call {
	_c.Call.Return(run)
	return _c
}

// toGroupAddress provides a mock function with given fields:
func (_m *MockGroupAddressTag) toGroupAddress() readwritemodel.KnxGroupAddress {
	ret := _m.Called()

	var r0 readwritemodel.KnxGroupAddress
	if rf, ok := ret.Get(0).(func() readwritemodel.KnxGroupAddress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(readwritemodel.KnxGroupAddress)
		}
	}

	return r0
}

// MockGroupAddressTag_toGroupAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'toGroupAddress'
type MockGroupAddressTag_toGroupAddress_Call struct {
	*mock.Call
}

// toGroupAddress is a helper method to define mock.On call
func (_e *MockGroupAddressTag_Expecter) toGroupAddress() *MockGroupAddressTag_toGroupAddress_Call {
	return &MockGroupAddressTag_toGroupAddress_Call{Call: _e.mock.On("toGroupAddress")}
}

func (_c *MockGroupAddressTag_toGroupAddress_Call) Run(run func()) *MockGroupAddressTag_toGroupAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockGroupAddressTag_toGroupAddress_Call) Return(_a0 readwritemodel.KnxGroupAddress) *MockGroupAddressTag_toGroupAddress_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockGroupAddressTag_toGroupAddress_Call) RunAndReturn(run func() readwritemodel.KnxGroupAddress) *MockGroupAddressTag_toGroupAddress_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockGroupAddressTag interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockGroupAddressTag creates a new instance of MockGroupAddressTag. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockGroupAddressTag(t mockConstructorTestingTNewMockGroupAddressTag) *MockGroupAddressTag {
	mock := &MockGroupAddressTag{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
