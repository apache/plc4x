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

package bacnetip

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type TaskRequirements interface {
	ProcessTask() error
	InstallTask(when *time.Time, delta *time.Duration)
	GetTaskTime() *time.Time
	GetIsScheduled() bool
	SetIsScheduled(isScheduled bool)
}

type Task struct {
	taskRequirements TaskRequirements
	taskTime         *time.Time
	isScheduled      bool
}

func NewTask(taskRequirements TaskRequirements) *Task {
	return &Task{taskRequirements: taskRequirements}
}

func (t *Task) InstallTask(when *time.Time, delta *time.Duration) {
	// check for delta from now
	if when == nil && delta != nil {
		_when := _taskManager.GetTime().Add(*delta)
		when = &_when
	}

	// fallback to the initial value
	if when == nil {
		_when := t.taskTime
		when = _when
	}
	if when == nil {
		panic("schedule missing, use zero for 'now'")
	}
	t.taskTime = when

	// pass along to the task manager
	_taskManager.InstallTask(t.taskRequirements)
}

func (t *Task) SuspendTask() {
	_taskManager.SuspendTask(t.taskRequirements)
}

func (t *Task) Resume() {
	_taskManager.ResumeTask(t.taskRequirements)
}

func (t *Task) GetTaskTime() *time.Time {
	return t.taskTime
}

func (t *Task) GetIsScheduled() bool {
	return t.isScheduled
}

func (t *Task) SetIsScheduled(isScheduled bool) {
	t.isScheduled = isScheduled
}

func (t *Task) String() string {
	return fmt.Sprintf("Task(taskTime: %v, isScheduled: %v)", t.taskTime, t.isScheduled)
}

type OneShotTaskRequirements interface {
	ProcessTask() error
}

type OneShotTask struct {
	*Task
	OneShotTaskRequirements
}

func NewOneShotTask(oneShotTaskRequirements OneShotTaskRequirements, when *time.Time) *OneShotTask {
	o := &OneShotTask{
		OneShotTaskRequirements: oneShotTaskRequirements,
	}
	o.Task = NewTask(o)
	if when != nil {
		o.taskTime = when
	}
	return o
}

type OneShotDeleteTask struct {
	*Task
	OneShotTaskRequirements
}

func NewOneShotDeleteTask(oneShotTaskRequirements OneShotTaskRequirements, when *time.Time) *OneShotDeleteTask {
	o := &OneShotDeleteTask{OneShotTaskRequirements: oneShotTaskRequirements}
	o.Task = NewTask(o)
	if when != nil {
		o.taskTime = when
	}
	return o
}

func (r *OneShotDeleteTask) IsOneShotDeleteTask() bool {
	return true
}

type OneShotFunctionTask struct {
	*OneShotDeleteTask
	fn func() error
}

func (m *OneShotFunctionTask) ProcessTask() error {
	return m.fn()
}

func OneShotFunction(fn func() error) *OneShotFunctionTask {
	task := &OneShotFunctionTask{fn: fn}
	task.OneShotDeleteTask = NewOneShotDeleteTask(task, nil)
	var delta time.Duration = 0
	task.InstallTask(nil, &delta)
	return task
}

func FunctionTask(fn func() error) *OneShotFunctionTask {
	task := &OneShotFunctionTask{fn: fn}
	task.OneShotDeleteTask = NewOneShotDeleteTask(task, nil)
	return task
}

type RecurringTaskRequirements interface {
	ProcessTask() error
}

type RecurringTask struct {
	*Task
	RecurringTaskRequirements
	taskInterval       *time.Duration
	taskIntervalOffset *time.Duration
}

func NewRecurringTask(localLog zerolog.Logger, recurringTaskRequirements RecurringTaskRequirements, interval *time.Duration, offset *time.Duration) *RecurringTask {
	r := &RecurringTask{RecurringTaskRequirements: recurringTaskRequirements}
	r.Task = NewTask(r)
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
	if *r.taskInterval <= 0.0 {
		panic("interval must be greater than zero")
	}

	// get ready for the next interval plus a jitter
	now := _taskManager.GetTime().Add(10 + time.Nanosecond)

	if r.taskIntervalOffset != nil {
		_offset := *r.taskIntervalOffset
		offset = &_offset
	} else {
		_offset := time.Duration(0)
		offset = &_offset
	}
	localLog.Debug().
		Interface("now", now).
		Interface("interval", interval).
		Interface("offset", offset).
		Msg("Now, interval, offset:")

	// compute the time
	_taskTime := now.Add(-*offset).Add(*interval) // TODO: check why upstream is doing the modulo operation (missing code here)
	r.taskTime = &_taskTime

	// install it
	_taskManager.InstallTask(r)

	return r
}

func (r *RecurringTask) IsRecurringTask() bool {
	return true
}

type recurringFunctionTask struct {
	*RecurringTask
	fn func() error
}

func newRecurringFunctionTask(localLog zerolog.Logger, interval *time.Duration, fn func() error) *recurringFunctionTask {
	r := &recurringFunctionTask{fn: fn}
	r.RecurringTask = NewRecurringTask(localLog, r, interval, nil)
	return r
}

func (r recurringFunctionTask) ProcessTask() error {
	return r.fn()
}

func RecurringFunctionTask(localLog zerolog.Logger, interval *time.Duration, fn func() error) *RecurringTask {
	return newRecurringFunctionTask(localLog, interval, fn).RecurringTask
}

var _taskManager = TaskManager{}

type TaskManager struct {
	sync.Mutex
	Tasks []TaskRequirements

	log zerolog.Logger
}

func NewTaskManager(localLog zerolog.Logger) *TaskManager {
	return &TaskManager{
		log: localLog,
	}
}

func (m *TaskManager) GetTime() time.Time {
	return time.Now()
}

func (m *TaskManager) InstallTask(task TaskRequirements) {
	m.Lock()
	defer m.Unlock()
	m.log.Debug().Interface("task", task).Msg("InstallTask")

	// if the taskTime is None is hasn't been computed correctly
	if task.GetTaskTime() == nil {
		panic("task time is None")
	}

	// if this is already installed, suspend it
	if task.GetIsScheduled() {
		m.SuspendTask(task)
	}

	// save this in the task list
	// TODO: we might need to insert it at the right place
	m.Tasks = append(m.Tasks, task)

	task.SetIsScheduled(true)
}

func (m *TaskManager) SuspendTask(task TaskRequirements) {
	m.log.Debug().Interface("task", task).Msg("SuspendTask ")
	m.Lock()
	defer m.Unlock()

	iToDelete := -1
	for i, _task := range m.Tasks {
		if _task == task {
			m.log.Debug().Msg("task found")
			iToDelete = i
			task.SetIsScheduled(false)
			break
		}
	}
	if iToDelete > 0 {
		m.Tasks = append(m.Tasks[:iToDelete], m.Tasks[iToDelete+1:]...)
	} else {
		m.log.Debug().Msg("task not found")
	}
}

func (m *TaskManager) ResumeTask(task TaskRequirements) {
	m.log.Debug().Interface("task", task).Msg("ResumeTask")
	m.Lock()
	defer m.Unlock()

	// just re-install it
	m.InstallTask(task)
}

func (m *TaskManager) GetNextTask() (TaskRequirements, time.Duration) {
	//log.Trace().Msg("GetNextTask")
	m.Lock()
	defer m.Unlock()

	now := time.Now()

	var task TaskRequirements
	var delta time.Duration

	if len(m.Tasks) > 0 {
		nextTask := m.Tasks[0]
		when := nextTask.GetTaskTime()
		if when.Before(now) {
			// pull it off the list and mark that it's no longer scheduled
			m.Tasks = m.Tasks[1:] // TODO: guard against empty list
			task = nextTask
			task.SetIsScheduled(false)

			if len(m.Tasks) > 0 {
				nextTask = m.Tasks[0]
				when = nextTask.GetTaskTime()
				// peek at the next task, return how long to wait
				delta = when.Sub(now) // TODO: avoid negative
			}
		} else {
			delta = when.Sub(now)
		}
	}

	// return the task to run and how long to wait for the next one
	return task, delta
}

func (m *TaskManager) ProcessTask(task TaskRequirements) {
	m.log.Debug().Interface("task", task).Msg("ProcessTask")

	// process the task
	if err := task.ProcessTask(); err != nil {
		m.log.Error().Err(err).Msg("Error processing Task")
	}

	switch task.(type) {
	case interface{ IsRecurringTask() bool }:
		task.InstallTask(nil, nil)
	case interface{ IsOneShotDeleteTask() bool }:
		// TODO: Delete? How?
	}
}
