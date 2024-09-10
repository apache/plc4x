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

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
)

type RecurringFunctionTask struct {
	*RecurringTask
	fn     GenericFunction
	args   Args
	kwargs KWArgs
}

func NewRecurringFunctionTask(localLog zerolog.Logger, fn GenericFunction, args Args, kwargs KWArgs, opts ...func(*RecurringFunctionTask)) *RecurringFunctionTask {
	r := &RecurringFunctionTask{fn: fn, args: args, kwargs: kwargs}
	for _, opt := range opts {
		opt(r)
	}
	r.RecurringTask = NewRecurringTask(localLog, r)
	return r
}

func WithRecurringFunctionTaskInterval(interval time.Duration) func(*RecurringFunctionTask) {
	return func(r *RecurringFunctionTask) {
		r.taskInterval = &interval
	}
}

func (r *RecurringFunctionTask) ProcessTask() error {
	return r.fn(r.args, r.kwargs)
}

func (r *RecurringFunctionTask) String() string {
	return fmt.Sprintf("RecurringFunctionTask(%v, fn: %t, args: %s, kwargs: %s)", r.RecurringTask, r.fn != nil, r.args, r.kwargs)
}
