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

package cache

import (
	context "context"

	model "github.com/apache/plc4x/plc4go/pkg/api/model"
	mock "github.com/stretchr/testify/mock"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"

	spi "github.com/apache/plc4x/plc4go/spi"
)

// mockTracedPlcConnection is an autogenerated mock type for the tracedPlcConnection type
type mockTracedPlcConnection struct {
	mock.Mock
}

type mockTracedPlcConnection_Expecter struct {
	mock *mock.Mock
}

func (_m *mockTracedPlcConnection) EXPECT() *mockTracedPlcConnection_Expecter {
	return &mockTracedPlcConnection_Expecter{mock: &_m.Mock}
}

// BlockingClose provides a mock function with given fields:
func (_m *mockTracedPlcConnection) BlockingClose() {
	_m.Called()
}

// mockTracedPlcConnection_BlockingClose_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BlockingClose'
type mockTracedPlcConnection_BlockingClose_Call struct {
	*mock.Call
}

// BlockingClose is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) BlockingClose() *mockTracedPlcConnection_BlockingClose_Call {
	return &mockTracedPlcConnection_BlockingClose_Call{Call: _e.mock.On("BlockingClose")}
}

func (_c *mockTracedPlcConnection_BlockingClose_Call) Run(run func()) *mockTracedPlcConnection_BlockingClose_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_BlockingClose_Call) Return() *mockTracedPlcConnection_BlockingClose_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockTracedPlcConnection_BlockingClose_Call) RunAndReturn(run func()) *mockTracedPlcConnection_BlockingClose_Call {
	_c.Call.Return(run)
	return _c
}

// BrowseRequestBuilder provides a mock function with given fields:
func (_m *mockTracedPlcConnection) BrowseRequestBuilder() model.PlcBrowseRequestBuilder {
	ret := _m.Called()

	var r0 model.PlcBrowseRequestBuilder
	if rf, ok := ret.Get(0).(func() model.PlcBrowseRequestBuilder); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.PlcBrowseRequestBuilder)
		}
	}

	return r0
}

// mockTracedPlcConnection_BrowseRequestBuilder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BrowseRequestBuilder'
type mockTracedPlcConnection_BrowseRequestBuilder_Call struct {
	*mock.Call
}

// BrowseRequestBuilder is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) BrowseRequestBuilder() *mockTracedPlcConnection_BrowseRequestBuilder_Call {
	return &mockTracedPlcConnection_BrowseRequestBuilder_Call{Call: _e.mock.On("BrowseRequestBuilder")}
}

func (_c *mockTracedPlcConnection_BrowseRequestBuilder_Call) Run(run func()) *mockTracedPlcConnection_BrowseRequestBuilder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_BrowseRequestBuilder_Call) Return(_a0 model.PlcBrowseRequestBuilder) *mockTracedPlcConnection_BrowseRequestBuilder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_BrowseRequestBuilder_Call) RunAndReturn(run func() model.PlcBrowseRequestBuilder) *mockTracedPlcConnection_BrowseRequestBuilder_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields:
func (_m *mockTracedPlcConnection) Close() <-chan plc4go.PlcConnectionCloseResult {
	ret := _m.Called()

	var r0 <-chan plc4go.PlcConnectionCloseResult
	if rf, ok := ret.Get(0).(func() <-chan plc4go.PlcConnectionCloseResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan plc4go.PlcConnectionCloseResult)
		}
	}

	return r0
}

// mockTracedPlcConnection_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type mockTracedPlcConnection_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) Close() *mockTracedPlcConnection_Close_Call {
	return &mockTracedPlcConnection_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *mockTracedPlcConnection_Close_Call) Run(run func()) *mockTracedPlcConnection_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_Close_Call) Return(_a0 <-chan plc4go.PlcConnectionCloseResult) *mockTracedPlcConnection_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_Close_Call) RunAndReturn(run func() <-chan plc4go.PlcConnectionCloseResult) *mockTracedPlcConnection_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Connect provides a mock function with given fields:
func (_m *mockTracedPlcConnection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	ret := _m.Called()

	var r0 <-chan plc4go.PlcConnectionConnectResult
	if rf, ok := ret.Get(0).(func() <-chan plc4go.PlcConnectionConnectResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan plc4go.PlcConnectionConnectResult)
		}
	}

	return r0
}

// mockTracedPlcConnection_Connect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Connect'
type mockTracedPlcConnection_Connect_Call struct {
	*mock.Call
}

// Connect is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) Connect() *mockTracedPlcConnection_Connect_Call {
	return &mockTracedPlcConnection_Connect_Call{Call: _e.mock.On("Connect")}
}

func (_c *mockTracedPlcConnection_Connect_Call) Run(run func()) *mockTracedPlcConnection_Connect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_Connect_Call) Return(_a0 <-chan plc4go.PlcConnectionConnectResult) *mockTracedPlcConnection_Connect_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_Connect_Call) RunAndReturn(run func() <-chan plc4go.PlcConnectionConnectResult) *mockTracedPlcConnection_Connect_Call {
	_c.Call.Return(run)
	return _c
}

// ConnectWithContext provides a mock function with given fields: ctx
func (_m *mockTracedPlcConnection) ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult {
	ret := _m.Called(ctx)

	var r0 <-chan plc4go.PlcConnectionConnectResult
	if rf, ok := ret.Get(0).(func(context.Context) <-chan plc4go.PlcConnectionConnectResult); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan plc4go.PlcConnectionConnectResult)
		}
	}

	return r0
}

// mockTracedPlcConnection_ConnectWithContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConnectWithContext'
type mockTracedPlcConnection_ConnectWithContext_Call struct {
	*mock.Call
}

// ConnectWithContext is a helper method to define mock.On call
//   - ctx context.Context
func (_e *mockTracedPlcConnection_Expecter) ConnectWithContext(ctx interface{}) *mockTracedPlcConnection_ConnectWithContext_Call {
	return &mockTracedPlcConnection_ConnectWithContext_Call{Call: _e.mock.On("ConnectWithContext", ctx)}
}

func (_c *mockTracedPlcConnection_ConnectWithContext_Call) Run(run func(ctx context.Context)) *mockTracedPlcConnection_ConnectWithContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *mockTracedPlcConnection_ConnectWithContext_Call) Return(_a0 <-chan plc4go.PlcConnectionConnectResult) *mockTracedPlcConnection_ConnectWithContext_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_ConnectWithContext_Call) RunAndReturn(run func(context.Context) <-chan plc4go.PlcConnectionConnectResult) *mockTracedPlcConnection_ConnectWithContext_Call {
	_c.Call.Return(run)
	return _c
}

// GetConnectionId provides a mock function with given fields:
func (_m *mockTracedPlcConnection) GetConnectionId() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// mockTracedPlcConnection_GetConnectionId_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConnectionId'
type mockTracedPlcConnection_GetConnectionId_Call struct {
	*mock.Call
}

// GetConnectionId is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) GetConnectionId() *mockTracedPlcConnection_GetConnectionId_Call {
	return &mockTracedPlcConnection_GetConnectionId_Call{Call: _e.mock.On("GetConnectionId")}
}

func (_c *mockTracedPlcConnection_GetConnectionId_Call) Run(run func()) *mockTracedPlcConnection_GetConnectionId_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_GetConnectionId_Call) Return(_a0 string) *mockTracedPlcConnection_GetConnectionId_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_GetConnectionId_Call) RunAndReturn(run func() string) *mockTracedPlcConnection_GetConnectionId_Call {
	_c.Call.Return(run)
	return _c
}

// GetMetadata provides a mock function with given fields:
func (_m *mockTracedPlcConnection) GetMetadata() model.PlcConnectionMetadata {
	ret := _m.Called()

	var r0 model.PlcConnectionMetadata
	if rf, ok := ret.Get(0).(func() model.PlcConnectionMetadata); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.PlcConnectionMetadata)
		}
	}

	return r0
}

// mockTracedPlcConnection_GetMetadata_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMetadata'
type mockTracedPlcConnection_GetMetadata_Call struct {
	*mock.Call
}

// GetMetadata is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) GetMetadata() *mockTracedPlcConnection_GetMetadata_Call {
	return &mockTracedPlcConnection_GetMetadata_Call{Call: _e.mock.On("GetMetadata")}
}

func (_c *mockTracedPlcConnection_GetMetadata_Call) Run(run func()) *mockTracedPlcConnection_GetMetadata_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_GetMetadata_Call) Return(_a0 model.PlcConnectionMetadata) *mockTracedPlcConnection_GetMetadata_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_GetMetadata_Call) RunAndReturn(run func() model.PlcConnectionMetadata) *mockTracedPlcConnection_GetMetadata_Call {
	_c.Call.Return(run)
	return _c
}

// GetTracer provides a mock function with given fields:
func (_m *mockTracedPlcConnection) GetTracer() *spi.Tracer {
	ret := _m.Called()

	var r0 *spi.Tracer
	if rf, ok := ret.Get(0).(func() *spi.Tracer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*spi.Tracer)
		}
	}

	return r0
}

// mockTracedPlcConnection_GetTracer_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTracer'
type mockTracedPlcConnection_GetTracer_Call struct {
	*mock.Call
}

// GetTracer is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) GetTracer() *mockTracedPlcConnection_GetTracer_Call {
	return &mockTracedPlcConnection_GetTracer_Call{Call: _e.mock.On("GetTracer")}
}

func (_c *mockTracedPlcConnection_GetTracer_Call) Run(run func()) *mockTracedPlcConnection_GetTracer_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_GetTracer_Call) Return(_a0 *spi.Tracer) *mockTracedPlcConnection_GetTracer_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_GetTracer_Call) RunAndReturn(run func() *spi.Tracer) *mockTracedPlcConnection_GetTracer_Call {
	_c.Call.Return(run)
	return _c
}

// IsConnected provides a mock function with given fields:
func (_m *mockTracedPlcConnection) IsConnected() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// mockTracedPlcConnection_IsConnected_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsConnected'
type mockTracedPlcConnection_IsConnected_Call struct {
	*mock.Call
}

// IsConnected is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) IsConnected() *mockTracedPlcConnection_IsConnected_Call {
	return &mockTracedPlcConnection_IsConnected_Call{Call: _e.mock.On("IsConnected")}
}

func (_c *mockTracedPlcConnection_IsConnected_Call) Run(run func()) *mockTracedPlcConnection_IsConnected_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_IsConnected_Call) Return(_a0 bool) *mockTracedPlcConnection_IsConnected_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_IsConnected_Call) RunAndReturn(run func() bool) *mockTracedPlcConnection_IsConnected_Call {
	_c.Call.Return(run)
	return _c
}

// IsTraceEnabled provides a mock function with given fields:
func (_m *mockTracedPlcConnection) IsTraceEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// mockTracedPlcConnection_IsTraceEnabled_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsTraceEnabled'
type mockTracedPlcConnection_IsTraceEnabled_Call struct {
	*mock.Call
}

// IsTraceEnabled is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) IsTraceEnabled() *mockTracedPlcConnection_IsTraceEnabled_Call {
	return &mockTracedPlcConnection_IsTraceEnabled_Call{Call: _e.mock.On("IsTraceEnabled")}
}

func (_c *mockTracedPlcConnection_IsTraceEnabled_Call) Run(run func()) *mockTracedPlcConnection_IsTraceEnabled_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_IsTraceEnabled_Call) Return(_a0 bool) *mockTracedPlcConnection_IsTraceEnabled_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_IsTraceEnabled_Call) RunAndReturn(run func() bool) *mockTracedPlcConnection_IsTraceEnabled_Call {
	_c.Call.Return(run)
	return _c
}

// Ping provides a mock function with given fields:
func (_m *mockTracedPlcConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
	ret := _m.Called()

	var r0 <-chan plc4go.PlcConnectionPingResult
	if rf, ok := ret.Get(0).(func() <-chan plc4go.PlcConnectionPingResult); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan plc4go.PlcConnectionPingResult)
		}
	}

	return r0
}

// mockTracedPlcConnection_Ping_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Ping'
type mockTracedPlcConnection_Ping_Call struct {
	*mock.Call
}

// Ping is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) Ping() *mockTracedPlcConnection_Ping_Call {
	return &mockTracedPlcConnection_Ping_Call{Call: _e.mock.On("Ping")}
}

func (_c *mockTracedPlcConnection_Ping_Call) Run(run func()) *mockTracedPlcConnection_Ping_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_Ping_Call) Return(_a0 <-chan plc4go.PlcConnectionPingResult) *mockTracedPlcConnection_Ping_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_Ping_Call) RunAndReturn(run func() <-chan plc4go.PlcConnectionPingResult) *mockTracedPlcConnection_Ping_Call {
	_c.Call.Return(run)
	return _c
}

// ReadRequestBuilder provides a mock function with given fields:
func (_m *mockTracedPlcConnection) ReadRequestBuilder() model.PlcReadRequestBuilder {
	ret := _m.Called()

	var r0 model.PlcReadRequestBuilder
	if rf, ok := ret.Get(0).(func() model.PlcReadRequestBuilder); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.PlcReadRequestBuilder)
		}
	}

	return r0
}

// mockTracedPlcConnection_ReadRequestBuilder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadRequestBuilder'
type mockTracedPlcConnection_ReadRequestBuilder_Call struct {
	*mock.Call
}

// ReadRequestBuilder is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) ReadRequestBuilder() *mockTracedPlcConnection_ReadRequestBuilder_Call {
	return &mockTracedPlcConnection_ReadRequestBuilder_Call{Call: _e.mock.On("ReadRequestBuilder")}
}

func (_c *mockTracedPlcConnection_ReadRequestBuilder_Call) Run(run func()) *mockTracedPlcConnection_ReadRequestBuilder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_ReadRequestBuilder_Call) Return(_a0 model.PlcReadRequestBuilder) *mockTracedPlcConnection_ReadRequestBuilder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_ReadRequestBuilder_Call) RunAndReturn(run func() model.PlcReadRequestBuilder) *mockTracedPlcConnection_ReadRequestBuilder_Call {
	_c.Call.Return(run)
	return _c
}

// SubscriptionRequestBuilder provides a mock function with given fields:
func (_m *mockTracedPlcConnection) SubscriptionRequestBuilder() model.PlcSubscriptionRequestBuilder {
	ret := _m.Called()

	var r0 model.PlcSubscriptionRequestBuilder
	if rf, ok := ret.Get(0).(func() model.PlcSubscriptionRequestBuilder); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.PlcSubscriptionRequestBuilder)
		}
	}

	return r0
}

// mockTracedPlcConnection_SubscriptionRequestBuilder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SubscriptionRequestBuilder'
type mockTracedPlcConnection_SubscriptionRequestBuilder_Call struct {
	*mock.Call
}

// SubscriptionRequestBuilder is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) SubscriptionRequestBuilder() *mockTracedPlcConnection_SubscriptionRequestBuilder_Call {
	return &mockTracedPlcConnection_SubscriptionRequestBuilder_Call{Call: _e.mock.On("SubscriptionRequestBuilder")}
}

func (_c *mockTracedPlcConnection_SubscriptionRequestBuilder_Call) Run(run func()) *mockTracedPlcConnection_SubscriptionRequestBuilder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_SubscriptionRequestBuilder_Call) Return(_a0 model.PlcSubscriptionRequestBuilder) *mockTracedPlcConnection_SubscriptionRequestBuilder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_SubscriptionRequestBuilder_Call) RunAndReturn(run func() model.PlcSubscriptionRequestBuilder) *mockTracedPlcConnection_SubscriptionRequestBuilder_Call {
	_c.Call.Return(run)
	return _c
}

// UnsubscriptionRequestBuilder provides a mock function with given fields:
func (_m *mockTracedPlcConnection) UnsubscriptionRequestBuilder() model.PlcUnsubscriptionRequestBuilder {
	ret := _m.Called()

	var r0 model.PlcUnsubscriptionRequestBuilder
	if rf, ok := ret.Get(0).(func() model.PlcUnsubscriptionRequestBuilder); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.PlcUnsubscriptionRequestBuilder)
		}
	}

	return r0
}

// mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnsubscriptionRequestBuilder'
type mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call struct {
	*mock.Call
}

// UnsubscriptionRequestBuilder is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) UnsubscriptionRequestBuilder() *mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call {
	return &mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call{Call: _e.mock.On("UnsubscriptionRequestBuilder")}
}

func (_c *mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call) Run(run func()) *mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call) Return(_a0 model.PlcUnsubscriptionRequestBuilder) *mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call) RunAndReturn(run func() model.PlcUnsubscriptionRequestBuilder) *mockTracedPlcConnection_UnsubscriptionRequestBuilder_Call {
	_c.Call.Return(run)
	return _c
}

// WriteRequestBuilder provides a mock function with given fields:
func (_m *mockTracedPlcConnection) WriteRequestBuilder() model.PlcWriteRequestBuilder {
	ret := _m.Called()

	var r0 model.PlcWriteRequestBuilder
	if rf, ok := ret.Get(0).(func() model.PlcWriteRequestBuilder); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(model.PlcWriteRequestBuilder)
		}
	}

	return r0
}

// mockTracedPlcConnection_WriteRequestBuilder_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WriteRequestBuilder'
type mockTracedPlcConnection_WriteRequestBuilder_Call struct {
	*mock.Call
}

// WriteRequestBuilder is a helper method to define mock.On call
func (_e *mockTracedPlcConnection_Expecter) WriteRequestBuilder() *mockTracedPlcConnection_WriteRequestBuilder_Call {
	return &mockTracedPlcConnection_WriteRequestBuilder_Call{Call: _e.mock.On("WriteRequestBuilder")}
}

func (_c *mockTracedPlcConnection_WriteRequestBuilder_Call) Run(run func()) *mockTracedPlcConnection_WriteRequestBuilder_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockTracedPlcConnection_WriteRequestBuilder_Call) Return(_a0 model.PlcWriteRequestBuilder) *mockTracedPlcConnection_WriteRequestBuilder_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockTracedPlcConnection_WriteRequestBuilder_Call) RunAndReturn(run func() model.PlcWriteRequestBuilder) *mockTracedPlcConnection_WriteRequestBuilder_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockTracedPlcConnection interface {
	mock.TestingT
	Cleanup(func())
}

// newMockTracedPlcConnection creates a new instance of mockTracedPlcConnection. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockTracedPlcConnection(t mockConstructorTestingTnewMockTracedPlcConnection) *mockTracedPlcConnection {
	mock := &mockTracedPlcConnection{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
