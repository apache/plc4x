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

package bacnetip

import mock "github.com/stretchr/testify/mock"

// mock_IOQController is an autogenerated mock type for the _IOQController type
type mock_IOQController struct {
	mock.Mock
}

type mock_IOQController_Expecter struct {
	mock *mock.Mock
}

func (_m *mock_IOQController) EXPECT() *mock_IOQController_Expecter {
	return &mock_IOQController_Expecter{mock: &_m.Mock}
}

// ProcessIO provides a mock function with given fields: iocb
func (_m *mock_IOQController) ProcessIO(iocb _IOCB) error {
	ret := _m.Called(iocb)

	var r0 error
	if rf, ok := ret.Get(0).(func(_IOCB) error); ok {
		r0 = rf(iocb)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mock_IOQController_ProcessIO_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProcessIO'
type mock_IOQController_ProcessIO_Call struct {
	*mock.Call
}

// ProcessIO is a helper method to define mock.On call
//   - iocb _IOCB
func (_e *mock_IOQController_Expecter) ProcessIO(iocb interface{}) *mock_IOQController_ProcessIO_Call {
	return &mock_IOQController_ProcessIO_Call{Call: _e.mock.On("ProcessIO", iocb)}
}

func (_c *mock_IOQController_ProcessIO_Call) Run(run func(iocb _IOCB)) *mock_IOQController_ProcessIO_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(_IOCB))
	})
	return _c
}

func (_c *mock_IOQController_ProcessIO_Call) Return(_a0 error) *mock_IOQController_ProcessIO_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mock_IOQController_ProcessIO_Call) RunAndReturn(run func(_IOCB) error) *mock_IOQController_ProcessIO_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMock_IOQController interface {
	mock.TestingT
	Cleanup(func())
}

// newMock_IOQController creates a new instance of mock_IOQController. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMock_IOQController(t mockConstructorTestingTnewMock_IOQController) *mock_IOQController {
	mock := &mock_IOQController{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
