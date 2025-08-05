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

package state_machine

import (
	"fmt"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
)

//go:generate go tool plc4xGenerator -type=TimeoutTask -prefix=state_machine_
type TimeoutTask struct {
	*OneShotTask

	fn     GenericFunction `ignore:"true"`
	args   Args
	kwArgs KWArgs
}

func NewTimeoutTask(fn GenericFunction, args Args, kwArgs KWArgs, options ...Option) *TimeoutTask {
	if _debug != nil {
		_debug("__init__ %r %r %r", fn, args, kwArgs)
	}
	t := &TimeoutTask{
		fn:     fn,
		args:   args,
		kwArgs: kwArgs,
	}
	ApplyAppliers(options, t)
	optionsForParent := AddLeafTypeIfAbundant(options, t)
	t.OneShotTask = NewOneShotTask(t, optionsForParent...)
	return t
}

func (t *TimeoutTask) ProcessTask() error {
	if _debug != nil {
		_debug("process_task %r", t.fn)
	}
	return t.fn(t.args, t.kwArgs)
}

func (t *TimeoutTask) Format(s fmt.State, c rune) {
	switch c {
	case 's', 'v', 'r':
		_, _ = s.Write([]byte(fmt.Sprintf("process_task %p", t.fn)))
	}
}
