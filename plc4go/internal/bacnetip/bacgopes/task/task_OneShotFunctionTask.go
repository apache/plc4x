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

package task

import . "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"

//go:generate go tool plc4xGenerator -type=OneShotFunctionTask -prefix=task_
type OneShotFunctionTask struct {
	*OneShotDeleteTask
	fn     GenericFunction `ignore:"true"`
	args   Args
	kwArgs KWArgs
}

func OneShotFunction(fn GenericFunction, args Args, kwArgs KWArgs) *OneShotFunctionTask {
	task := &OneShotFunctionTask{fn: fn, args: args, kwArgs: kwArgs}
	task.OneShotDeleteTask = NewOneShotDeleteTask(task)
	task.InstallTask(WithInstallTaskOptionsDelta(0))
	return task
}

func FunctionTask(fn GenericFunction, args Args, kwArgs KWArgs) *OneShotFunctionTask {
	task := &OneShotFunctionTask{fn: fn, args: args, kwArgs: kwArgs}
	task.OneShotDeleteTask = NewOneShotDeleteTask(task)
	return task
}

func (m *OneShotFunctionTask) ProcessTask() error {
	return m.fn(m.args, m.kwArgs)
}
