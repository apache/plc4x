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

package bacgopes

import "fmt"

type OneShotFunctionTask struct {
	*OneShotDeleteTask
	fn     func(args Args, kwargs KWArgs) error
	args   Args
	kwargs KWArgs
}

func OneShotFunction(fn func(args Args, kwargs KWArgs) error, args Args, kwargs KWArgs) *OneShotFunctionTask {
	task := &OneShotFunctionTask{fn: fn, args: args, kwargs: kwargs}
	task.OneShotDeleteTask = NewOneShotDeleteTask(task, nil)
	task.InstallTask(WithInstallTaskOptionsDelta(0))
	return task
}

func FunctionTask(fn func(args Args, kwargs KWArgs) error, args Args, kwargs KWArgs) *OneShotFunctionTask {
	task := &OneShotFunctionTask{fn: fn, args: args, kwargs: kwargs}
	task.OneShotDeleteTask = NewOneShotDeleteTask(task, nil)
	return task
}

func (m *OneShotFunctionTask) ProcessTask() error {
	return m.fn(m.args, m.kwargs)
}

func (m *OneShotFunctionTask) String() string {
	return fmt.Sprintf("OneShotFunctionTask(%v, fn: %t, args: %s, kwargs: %s)", m.OneShotDeleteTask, m.fn != nil, m.args, m.kwargs)
}
