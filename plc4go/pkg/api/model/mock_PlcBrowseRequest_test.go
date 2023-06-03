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

package model

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockPlcBrowseRequest is an autogenerated mock type for the PlcBrowseRequest type
type MockPlcBrowseRequest struct {
	mock.Mock
}

type MockPlcBrowseRequest_Expecter struct {
	mock *mock.Mock
}

func (_m *MockPlcBrowseRequest) EXPECT() *MockPlcBrowseRequest_Expecter {
	return &MockPlcBrowseRequest_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields:
func (_m *MockPlcBrowseRequest) Execute() <-chan PlcBrowseRequestResult {
	ret := _m.Called()

	var r0 <-chan PlcBrowseRequestResult
	if rf, ok := ret.Get(0).(func() <-chan PlcBrowseRequestResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan PlcBrowseRequestResult)
		}
	}

	return r0
}

// MockPlcBrowseRequest_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockPlcBrowseRequest_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
func (_e *MockPlcBrowseRequest_Expecter) Execute() *MockPlcBrowseRequest_Execute_Call {
	return &MockPlcBrowseRequest_Execute_Call{Call: _e.mock.On("Execute")}
}

func (_c *MockPlcBrowseRequest_Execute_Call) Run(run func()) *MockPlcBrowseRequest_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockPlcBrowseRequest_Execute_Call) Return(_a0 <-chan PlcBrowseRequestResult) *MockPlcBrowseRequest_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcBrowseRequest_Execute_Call) RunAndReturn(run func() <-chan PlcBrowseRequestResult) *MockPlcBrowseRequest_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteWithContext provides a mock function with given fields: ctx
func (_m *MockPlcBrowseRequest) ExecuteWithContext(ctx context.Context) <-chan PlcBrowseRequestResult {
	ret := _m.Called(ctx)

	var r0 <-chan PlcBrowseRequestResult
	if rf, ok := ret.Get(0).(func(context.Context) <-chan PlcBrowseRequestResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan PlcBrowseRequestResult)
		}
	}

	return r0
}

// MockPlcBrowseRequest_ExecuteWithContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteWithContext'
type MockPlcBrowseRequest_ExecuteWithContext_Call struct {
	*mock.Call
}

// ExecuteWithContext is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockPlcBrowseRequest_Expecter) ExecuteWithContext(ctx interface{}) *MockPlcBrowseRequest_ExecuteWithContext_Call {
	return &MockPlcBrowseRequest_ExecuteWithContext_Call{Call: _e.mock.On("ExecuteWithContext", ctx)}
}

func (_c *MockPlcBrowseRequest_ExecuteWithContext_Call) Run(run func(ctx context.Context)) *MockPlcBrowseRequest_ExecuteWithContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockPlcBrowseRequest_ExecuteWithContext_Call) Return(_a0 <-chan PlcBrowseRequestResult) *MockPlcBrowseRequest_ExecuteWithContext_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcBrowseRequest_ExecuteWithContext_Call) RunAndReturn(run func(context.Context) <-chan PlcBrowseRequestResult) *MockPlcBrowseRequest_ExecuteWithContext_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteWithInterceptor provides a mock function with given fields: interceptor
func (_m *MockPlcBrowseRequest) ExecuteWithInterceptor(interceptor func(PlcBrowseItem) bool) <-chan PlcBrowseRequestResult {
	ret := _m.Called(interceptor)

	var r0 <-chan PlcBrowseRequestResult
	if rf, ok := ret.Get(0).(func(func(PlcBrowseItem) bool) <-chan PlcBrowseRequestResult); ok {
		r0 = rf(interceptor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan PlcBrowseRequestResult)
		}
	}

	return r0
}

// MockPlcBrowseRequest_ExecuteWithInterceptor_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteWithInterceptor'
type MockPlcBrowseRequest_ExecuteWithInterceptor_Call struct {
	*mock.Call
}

// ExecuteWithInterceptor is a helper method to define mock.On call
//   - interceptor func(PlcBrowseItem) bool
func (_e *MockPlcBrowseRequest_Expecter) ExecuteWithInterceptor(interceptor interface{}) *MockPlcBrowseRequest_ExecuteWithInterceptor_Call {
	return &MockPlcBrowseRequest_ExecuteWithInterceptor_Call{Call: _e.mock.On("ExecuteWithInterceptor", interceptor)}
}

func (_c *MockPlcBrowseRequest_ExecuteWithInterceptor_Call) Run(run func(interceptor func(PlcBrowseItem) bool)) *MockPlcBrowseRequest_ExecuteWithInterceptor_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(func(PlcBrowseItem) bool))
	})
	return _c
}

func (_c *MockPlcBrowseRequest_ExecuteWithInterceptor_Call) Return(_a0 <-chan PlcBrowseRequestResult) *MockPlcBrowseRequest_ExecuteWithInterceptor_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcBrowseRequest_ExecuteWithInterceptor_Call) RunAndReturn(run func(func(PlcBrowseItem) bool) <-chan PlcBrowseRequestResult) *MockPlcBrowseRequest_ExecuteWithInterceptor_Call {
	_c.Call.Return(run)
	return _c
}

// ExecuteWithInterceptorWithContext provides a mock function with given fields: ctx, interceptor
func (_m *MockPlcBrowseRequest) ExecuteWithInterceptorWithContext(ctx context.Context, interceptor func(PlcBrowseItem) bool) <-chan PlcBrowseRequestResult {
	ret := _m.Called(ctx, interceptor)

	var r0 <-chan PlcBrowseRequestResult
	if rf, ok := ret.Get(0).(func(context.Context, func(PlcBrowseItem) bool) <-chan PlcBrowseRequestResult); ok {
		r0 = rf(ctx, interceptor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan PlcBrowseRequestResult)
		}
	}

	return r0
}

// MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExecuteWithInterceptorWithContext'
type MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call struct {
	*mock.Call
}

// ExecuteWithInterceptorWithContext is a helper method to define mock.On call
//   - ctx context.Context
//   - interceptor func(PlcBrowseItem) bool
func (_e *MockPlcBrowseRequest_Expecter) ExecuteWithInterceptorWithContext(ctx interface{}, interceptor interface{}) *MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call {
	return &MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call{Call: _e.mock.On("ExecuteWithInterceptorWithContext", ctx, interceptor)}
}

func (_c *MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call) Run(run func(ctx context.Context, interceptor func(PlcBrowseItem) bool)) *MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(PlcBrowseItem) bool))
	})
	return _c
}

func (_c *MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call) Return(_a0 <-chan PlcBrowseRequestResult) *MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call) RunAndReturn(run func(context.Context, func(PlcBrowseItem) bool) <-chan PlcBrowseRequestResult) *MockPlcBrowseRequest_ExecuteWithInterceptorWithContext_Call {
	_c.Call.Return(run)
	return _c
}

// GetQuery provides a mock function with given fields: queryName
func (_m *MockPlcBrowseRequest) GetQuery(queryName string) PlcQuery {
	ret := _m.Called(queryName)

	var r0 PlcQuery
	if rf, ok := ret.Get(0).(func(string) PlcQuery); ok {
		r0 = rf(queryName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(PlcQuery)
		}
	}

	return r0
}

// MockPlcBrowseRequest_GetQuery_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQuery'
type MockPlcBrowseRequest_GetQuery_Call struct {
	*mock.Call
}

// GetQuery is a helper method to define mock.On call
//   - queryName string
func (_e *MockPlcBrowseRequest_Expecter) GetQuery(queryName interface{}) *MockPlcBrowseRequest_GetQuery_Call {
	return &MockPlcBrowseRequest_GetQuery_Call{Call: _e.mock.On("GetQuery", queryName)}
}

func (_c *MockPlcBrowseRequest_GetQuery_Call) Run(run func(queryName string)) *MockPlcBrowseRequest_GetQuery_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockPlcBrowseRequest_GetQuery_Call) Return(_a0 PlcQuery) *MockPlcBrowseRequest_GetQuery_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcBrowseRequest_GetQuery_Call) RunAndReturn(run func(string) PlcQuery) *MockPlcBrowseRequest_GetQuery_Call {
	_c.Call.Return(run)
	return _c
}

// GetQueryNames provides a mock function with given fields:
func (_m *MockPlcBrowseRequest) GetQueryNames() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// MockPlcBrowseRequest_GetQueryNames_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQueryNames'
type MockPlcBrowseRequest_GetQueryNames_Call struct {
	*mock.Call
}

// GetQueryNames is a helper method to define mock.On call
func (_e *MockPlcBrowseRequest_Expecter) GetQueryNames() *MockPlcBrowseRequest_GetQueryNames_Call {
	return &MockPlcBrowseRequest_GetQueryNames_Call{Call: _e.mock.On("GetQueryNames")}
}

func (_c *MockPlcBrowseRequest_GetQueryNames_Call) Run(run func()) *MockPlcBrowseRequest_GetQueryNames_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockPlcBrowseRequest_GetQueryNames_Call) Return(_a0 []string) *MockPlcBrowseRequest_GetQueryNames_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcBrowseRequest_GetQueryNames_Call) RunAndReturn(run func() []string) *MockPlcBrowseRequest_GetQueryNames_Call {
	_c.Call.Return(run)
	return _c
}

// IsAPlcMessage provides a mock function with given fields:
func (_m *MockPlcBrowseRequest) IsAPlcMessage() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockPlcBrowseRequest_IsAPlcMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsAPlcMessage'
type MockPlcBrowseRequest_IsAPlcMessage_Call struct {
	*mock.Call
}

// IsAPlcMessage is a helper method to define mock.On call
func (_e *MockPlcBrowseRequest_Expecter) IsAPlcMessage() *MockPlcBrowseRequest_IsAPlcMessage_Call {
	return &MockPlcBrowseRequest_IsAPlcMessage_Call{Call: _e.mock.On("IsAPlcMessage")}
}

func (_c *MockPlcBrowseRequest_IsAPlcMessage_Call) Run(run func()) *MockPlcBrowseRequest_IsAPlcMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockPlcBrowseRequest_IsAPlcMessage_Call) Return(_a0 bool) *MockPlcBrowseRequest_IsAPlcMessage_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcBrowseRequest_IsAPlcMessage_Call) RunAndReturn(run func() bool) *MockPlcBrowseRequest_IsAPlcMessage_Call {
	_c.Call.Return(run)
	return _c
}

// String provides a mock function with given fields:
func (_m *MockPlcBrowseRequest) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockPlcBrowseRequest_String_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'String'
type MockPlcBrowseRequest_String_Call struct {
	*mock.Call
}

// String is a helper method to define mock.On call
func (_e *MockPlcBrowseRequest_Expecter) String() *MockPlcBrowseRequest_String_Call {
	return &MockPlcBrowseRequest_String_Call{Call: _e.mock.On("String")}
}

func (_c *MockPlcBrowseRequest_String_Call) Run(run func()) *MockPlcBrowseRequest_String_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockPlcBrowseRequest_String_Call) Return(_a0 string) *MockPlcBrowseRequest_String_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockPlcBrowseRequest_String_Call) RunAndReturn(run func() string) *MockPlcBrowseRequest_String_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockPlcBrowseRequest interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockPlcBrowseRequest creates a new instance of MockPlcBrowseRequest. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockPlcBrowseRequest(t mockConstructorTestingTNewMockPlcBrowseRequest) *MockPlcBrowseRequest {
	mock := &MockPlcBrowseRequest{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
