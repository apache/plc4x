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

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type RecurringTask struct {
	*Task

	taskInterval       *time.Duration
	taskIntervalOffset *time.Duration

	log zerolog.Logger
}

func NewRecurringTask(localLog zerolog.Logger, taskRequirements TaskRequirements, opts ...func(*RecurringTask)) *RecurringTask {
	r := &RecurringTask{
		log: localLog,
	}
	for _, opt := range opts {
		opt(r)
	}
	r.Task = NewTask(taskRequirements)
	return r
}

func WithRecurringTaskInterval(interval time.Duration) func(task *RecurringTask) {
	return func(task *RecurringTask) {
		task.taskInterval = &interval
	}
}

func WithRecurringTaskOffset(offset time.Duration) func(task *RecurringTask) {
	return func(task *RecurringTask) {
		task.taskIntervalOffset = &offset
	}
}

func (r *RecurringTask) InstallTask(options InstallTaskOptions) {
	interval := options.Interval
	offset := options.Offset

	// set the interval if it hasn't already been set
	if interval != nil {
		r.taskInterval = interval
	}
	if offset != nil {
		r.taskIntervalOffset = offset
	}

	if r.taskInterval == nil {
		panic("interval unset, use ctor or install_task parameter")
	}
	if *r.taskInterval <= 0 {
		panic("interval must be greater than zero")
	}

	// if there is no task manager, postpone the install
	if _taskManager == nil {
		r.log.Trace().Msg("No task manager")
		_unscheduledTasks = append(_unscheduledTasks, r.TaskRequirements)
	} else {
		// get ready for the next interval plus a jitter
		now := _taskManager.GetTime().Add(10 + time.Nanosecond)

		interval := *r.taskInterval
		offset := 0 * time.Nanosecond
		if r.taskIntervalOffset != nil {
			offset = *r.taskIntervalOffset
		}
		r.log.Debug().
			Time("now", now).
			Dur("interval", interval).
			Dur("offset", offset).
			Msg("Now, interval, offset:")

		// compute the time
		_taskTime := now.Add(-offset).Add(interval).Add(-(time.Duration((now.Add(-offset).UnixNano() - time.Time{}.UnixNano()) % int64(interval)))).Add(offset)
		r.taskTime = &_taskTime
		r.log.Debug().Time("taskTime", _taskTime).Msg("task time")

		// install it
		_taskManager.InstallTask(r.TaskRequirements)
	}
}

func (r *RecurringTask) IsRecurringTask() bool {
	return true
}

func (r *RecurringTask) String() string {
	return fmt.Sprintf("RecurringTask(%v, taskInterval: %v, taskIntervalOffset: %v)", r.Task, r.taskInterval, r.taskIntervalOffset)
}
