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

type InstallTaskOptions struct {
	When     *time.Time
	Delta    *time.Duration
	Interval *time.Duration
	Offset   *time.Duration
}

type TaskRequirements interface {
	fmt.Stringer
	ProcessTask() error
	InstallTask(options InstallTaskOptions)
	GetTaskTime() *time.Time
	GetIsScheduled() bool
	SetIsScheduled(isScheduled bool)
}

type Task struct {
	TaskRequirements
	taskTime    *time.Time
	isScheduled bool
}

func NewTask(taskRequirements TaskRequirements) *Task {
	return &Task{TaskRequirements: taskRequirements}
}

func (t *Task) InstallTask(options InstallTaskOptions) {
	when := options.When
	delta := options.Delta
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
	_taskManager.InstallTask(t.TaskRequirements)
}

func (t *Task) SuspendTask() {
	_taskManager.SuspendTask(t.TaskRequirements)
}

func (t *Task) Resume() {
	_taskManager.ResumeTask(t.TaskRequirements)
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

type OneShotTask struct {
	*Task
}

func NewOneShotTask(taskRequirements TaskRequirements, when *time.Time) *OneShotTask {
	o := &OneShotTask{}
	o.Task = NewTask(taskRequirements)
	if when != nil {
		o.taskTime = when
	}
	return o
}

func (t *OneShotTask) String() string {
	return fmt.Sprintf("OneShotTask(%v)", t.Task)
}

type OneShotDeleteTask struct {
	*Task
}

func NewOneShotDeleteTask(taskRequirements TaskRequirements, when *time.Time) *OneShotDeleteTask {
	o := &OneShotDeleteTask{}
	o.Task = NewTask(taskRequirements)
	if when != nil {
		o.taskTime = when
	}
	return o
}

func (r *OneShotDeleteTask) IsOneShotDeleteTask() bool {
	return true
}

func (r *OneShotDeleteTask) String() string {
	return fmt.Sprintf("OneShotDeleteTask(%v)", r.Task)
}

type OneShotFunctionTask struct {
	*OneShotDeleteTask
	fn     func(args Args, kwargs KWArgs) error
	args   Args
	kwargs KWArgs
}

func OneShotFunction(fn func(args Args, kwargs KWArgs) error, args Args, kwargs KWArgs) *OneShotFunctionTask {
	task := &OneShotFunctionTask{fn: fn, args: args, kwargs: kwargs}
	task.OneShotDeleteTask = NewOneShotDeleteTask(task, nil)
	var delta time.Duration = 0
	task.InstallTask(InstallTaskOptions{Delta: &delta})
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

type RecurringTask struct {
	*Task

	taskInterval       *time.Duration
	taskIntervalOffset *time.Duration

	log zerolog.Logger
}

func NewRecurringTask(localLog zerolog.Logger, taskRequirements TaskRequirements, interval *time.Duration, offset *time.Duration) *RecurringTask {
	r := &RecurringTask{
		log: localLog,
	}
	r.Task = NewTask(taskRequirements)

	// save the interval, but do not automatically install
	r.taskInterval = interval
	r.taskIntervalOffset = offset
	return r
}

func (r *RecurringTask) InstallTask(options InstallTaskOptions) {
	interval := options.Interval
	offset := options.Offset
	if r.taskInterval == nil {
		panic("interval unset, use ctor or install_task parameter")
	}
	if *r.taskInterval <= 0.0 {
		panic("interval must be greater than zero")
	}

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
		_unscheduledTasks = append(_unscheduledTasks, r)
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
		_taskManager.InstallTask(r)
	}
}

func (r *RecurringTask) IsRecurringTask() bool {
	return true
}

func (r *RecurringTask) String() string {
	return fmt.Sprintf("RecurringTask(%v, taskInterval: %v, taskIntervalOffset: %v)", r.Task, r.taskInterval, r.taskIntervalOffset)
}

type RecurringFunctionTask struct {
	*RecurringTask
	fn     func(args Args, kwargs KWArgs) error
	args   Args
	kwargs KWArgs
}

func NewRecurringFunctionTask(localLog zerolog.Logger, interval *time.Duration, fn func(args Args, kwargs KWArgs) error, args Args, kwargs KWArgs) *RecurringFunctionTask {
	r := &RecurringFunctionTask{fn: fn, args: args, kwargs: kwargs}
	r.RecurringTask = NewRecurringTask(localLog, r, interval, nil)
	return r
}

func (r *RecurringFunctionTask) ProcessTask() error {
	return r.fn(r.args, r.kwargs)
}

func (r *RecurringFunctionTask) String() string {
	return fmt.Sprintf("RecurringFunctionTask(%v, fn: %t, args: %s, kwargs: %s)", r.RecurringTask, r.fn != nil, r.args, r.kwargs)
}

var _taskManager TaskManager
var _taskManagerMutex sync.Mutex
var _unscheduledTasks []any //TODO: check method clash in install task

type TaskManager interface {
	GetTime() time.Time
	InstallTask(task TaskRequirements)
	SuspendTask(task TaskRequirements)
	ResumeTask(task TaskRequirements)
	GetNextTask() (TaskRequirements, *time.Duration)
	ProcessTask(task TaskRequirements)
	GetTasks() []TaskRequirements
	PopTask() TaskRequirements
	ClearTasks()
}

type taskManager struct {
	sync.Mutex

	tasks []TaskRequirements

	log zerolog.Logger
}

func NewTaskManager(localLog zerolog.Logger) TaskManager {
	_taskManagerMutex.Lock()
	defer _taskManagerMutex.Unlock()
	if _taskManager != nil {
		return _taskManager
	}
	t := &taskManager{
		log: localLog,
	}

	// TODO: trigger

	// TODO: counter

	// TODO: unscheduled tasks

	_taskManager = t
	return t
}

func OverwriteTaskManager(manager TaskManager) {
	_taskManager = manager
}

func (m *taskManager) GetTime() time.Time {
	return time.Now()
}

func GetTaskManagerTime() time.Time {
	if _taskManager == nil {
		return time.Now()
	}
	return _taskManager.GetTime()
}

func (m *taskManager) ClearTasks() {
	m.tasks = nil
}

func (m *taskManager) InstallTask(task TaskRequirements) {
	m.log.Debug().Stringer("task", task).Msg("InstallTask")
	m.Lock()
	defer m.Unlock()

	// if the taskTime is None is hasn't been computed correctly
	if task.GetTaskTime() == nil {
		panic("task time is None")
	}

	// if this is already installed, suspend it
	if task.GetIsScheduled() {
		m.Unlock()
		m.SuspendTask(task)
		m.Lock()
	}

	// save this in the task list
	// TODO: we might need to insert it at the right place
	m.tasks = append(m.tasks, task)

	task.SetIsScheduled(true)
}

func (m *taskManager) SuspendTask(task TaskRequirements) {
	m.log.Debug().Stringer("task", task).Msg("SuspendTask ")
	m.Lock()
	defer m.Unlock()

	iToDelete := -1
	for i, _task := range m.tasks {
		if _task == task {
			m.log.Debug().Msg("task found")
			iToDelete = i
			task.SetIsScheduled(false)
			break
		}
	}
	if iToDelete > 0 {
		m.tasks = append(m.tasks[:iToDelete], m.tasks[iToDelete+1:]...)
	} else {
		m.log.Debug().Msg("task not found")
	}
}

func (m *taskManager) ResumeTask(task TaskRequirements) {
	m.log.Debug().Stringer("task", task).Msg("ResumeTask")
	m.Lock()
	defer m.Unlock()

	// just re-install it
	m.InstallTask(task)
}

func (m *taskManager) GetNextTask() (TaskRequirements, *time.Duration) {
	m.log.Trace().Msg("GetNextTask")
	m.Lock()
	defer m.Unlock()

	now := time.Now()

	var task TaskRequirements
	var delta *time.Duration

	if len(m.tasks) > 0 {
		nextTask := m.tasks[0]
		when := nextTask.GetTaskTime()
		if when.Before(now) {
			// pull it off the list and mark that it's no longer scheduled
			m.tasks = m.tasks[1:] // TODO: guard against empty list
			task = nextTask
			task.SetIsScheduled(false)

			if len(m.tasks) > 0 {
				nextTask = m.tasks[0]
				when = nextTask.GetTaskTime()
				// peek at the next task, return how long to wait
				newDelta := when.Sub(now) // TODO: avoid negative
				delta = &newDelta
			}
		} else {
			newDelta := when.Sub(now)
			delta = &newDelta
		}
	}

	// return the task to run and how long to wait for the next one
	return task, delta
}

func (m *taskManager) ProcessTask(task TaskRequirements) {
	m.log.Debug().Stringer("task", task).Msg("ProcessTask")

	// process the task
	if err := task.ProcessTask(); err != nil {
		m.log.Error().Err(err).Msg("Error processing Task")
	}

	switch task.(type) {
	case interface{ IsRecurringTask() bool }:
		task.InstallTask(InstallTaskOptions{})
	case interface{ IsOneShotDeleteTask() bool }:
		// TODO: Delete? How?
	}
}

func (m *taskManager) GetTasks() []TaskRequirements {
	return m.tasks
}

func (m *taskManager) PopTask() TaskRequirements {
	m.log.Trace().Msg("pop task")
	m.Lock()
	defer m.Unlock()
	task := m.tasks[0]
	m.tasks = m.tasks[1:]
	return task
}
