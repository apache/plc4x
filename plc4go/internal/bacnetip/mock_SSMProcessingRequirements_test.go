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

// MockSSMProcessingRequirements is an autogenerated mock type for the SSMProcessingRequirements type
type MockSSMProcessingRequirements struct {
	mock.Mock
}

type MockSSMProcessingRequirements_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSSMProcessingRequirements) EXPECT() *MockSSMProcessingRequirements_Expecter {
	return &MockSSMProcessingRequirements_Expecter{mock: &_m.Mock}
}

// processTask provides a mock function with given fields:
func (_m *MockSSMProcessingRequirements) processTask() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockSSMProcessingRequirements_processTask_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'processTask'
type MockSSMProcessingRequirements_processTask_Call struct {
	*mock.Call
}

// processTask is a helper method to define mock.On call
func (_e *MockSSMProcessingRequirements_Expecter) processTask() *MockSSMProcessingRequirements_processTask_Call {
	return &MockSSMProcessingRequirements_processTask_Call{Call: _e.mock.On("processTask")}
}

func (_c *MockSSMProcessingRequirements_processTask_Call) Run(run func()) *MockSSMProcessingRequirements_processTask_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockSSMProcessingRequirements_processTask_Call) Return(_a0 error) *MockSSMProcessingRequirements_processTask_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockSSMProcessingRequirements_processTask_Call) RunAndReturn(run func() error) *MockSSMProcessingRequirements_processTask_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockSSMProcessingRequirements interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockSSMProcessingRequirements creates a new instance of MockSSMProcessingRequirements. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockSSMProcessingRequirements(t mockConstructorTestingTNewMockSSMProcessingRequirements) *MockSSMProcessingRequirements {
	mock := &MockSSMProcessingRequirements{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
