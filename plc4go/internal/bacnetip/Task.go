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
	"container/heap"
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

func WithInstallTaskOptionsNone() InstallTaskOptions {
	return InstallTaskOptions{}
}

func WithInstallTaskOptionsWhen(when time.Time) InstallTaskOptions {
	return InstallTaskOptions{
		When: &when,
	}
}

func (i InstallTaskOptions) WithWhen(when time.Time) InstallTaskOptions {
	i.When = &when
	return i
}

func WithInstallTaskOptionsDelta(delta time.Duration) InstallTaskOptions {
	return InstallTaskOptions{
		Delta: &delta,
	}
}

func (i InstallTaskOptions) WithDelta(delta time.Duration) InstallTaskOptions {
	i.Delta = &delta
	return i
}

func WithInstallTaskOptionsInterval(interval time.Duration) InstallTaskOptions {
	return InstallTaskOptions{
		Interval: &interval,
	}
}

func (i InstallTaskOptions) WithInterval(interval time.Duration) InstallTaskOptions {
	i.Interval = &interval
	return i
}

func WithInstallTaskOptionsOffset(offset time.Duration) InstallTaskOptions {
	return InstallTaskOptions{
		Offset: &offset,
	}
}

func (i InstallTaskOptions) WithOffset(offset time.Duration) InstallTaskOptions {
	i.Offset = &offset
	return i
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

func NewTask(taskRequirements TaskRequirements, opts ...func(*Task)) *Task {
	t := &Task{TaskRequirements: taskRequirements}
	for _, opt := range opts {
		opt(t)
	}
	return t
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
	o.Task = NewTask(taskRequirements, func(task *Task) {
		task.taskTime = when
	})
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
	o.Task = NewTask(taskRequirements, func(task *Task) {
		task.taskTime = when
	})
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

type RecurringFunctionTask struct {
	*RecurringTask
	fn     func(args Args, kwargs KWArgs) error
	args   Args
	kwargs KWArgs
}

func NewRecurringFunctionTask(localLog zerolog.Logger, fn func(args Args, kwargs KWArgs) error, args Args, kwargs KWArgs, opts ...func(*RecurringFunctionTask)) *RecurringFunctionTask {
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

var _taskManager TaskManager
var _taskManagerMutex sync.Mutex
var _unscheduledTasks []TaskRequirements

type TaskManager interface {
	fmt.Stringer
	GetTime() time.Time
	InstallTask(task TaskRequirements)
	SuspendTask(task TaskRequirements)
	ResumeTask(task TaskRequirements)
	GetNextTask() (TaskRequirements, *time.Duration)
	ProcessTask(task TaskRequirements)
	GetTasks() []TaskRequirements
	PopTask() TaskRequirements
	CountTasks() int
	ClearTasks()
}

func OverwriteTaskManager(localLog zerolog.Logger, manager TaskManager) (oldManager TaskManager) {
	_taskManagerMutex.Lock()
	defer _taskManagerMutex.Unlock()
	if _taskManager != nil {
		oldManager = _taskManager
		if oldManager.CountTasks() > 0 {
			localLog.Warn().Stringer("oldManager", oldManager).Msg("Overwriting task manager with pending tasks")
		}
		_taskManager.ClearTasks()
	}
	_taskManager = manager
	return
}

func ClearTaskManager(localLog zerolog.Logger) {
	_taskManagerMutex.Lock()
	defer _taskManagerMutex.Unlock()
	if _taskManager == nil {
		localLog.Warn().Msg("No task manager to clear ")
		return
	}
	if _taskManager.CountTasks() > 0 {
		localLog.Warn().Stringer("taskManager", _taskManager).Msg("Clearing task manager with pending tasks")
	}
	_taskManager.ClearTasks()
}

type taskItem struct {
	taskTime *time.Time
	id       int
	task     TaskRequirements
}

func (t taskItem) String() string {
	return fmt.Sprintf("taskItem(taskTime:%v, id:%d, %v)", t.taskTime, t.id, t.task)
}

type taskManager struct {
	sync.Mutex

	tasks PriorityQueue[int64, taskItem]
	count int

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

	// task manager is this instance
	_taskManager = t

	// unique sequence counter
	t.count = 0

	// there may be tasks created that couldn't be scheduled
	// because a task manager wasn't created yet.
	for _, task := range _unscheduledTasks {
		task.InstallTask(WithInstallTaskOptionsNone())
	}

	return t
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
	m.tasks.clear()
}

func (m *taskManager) InstallTask(task TaskRequirements) {
	m.log.Debug().Stringer("task", task).Msg("InstallTask")
	m.Lock()
	defer m.Unlock()

	// if the taskTime is None it hasn't been computed correctly
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
	m.count++
	heap.Push(&m.tasks, &PriorityItem[int64, taskItem]{
		value:    taskItem{taskTime: task.GetTaskTime(), id: m.count, task: task},
		priority: task.GetTaskTime().UnixNano() - time.Time{}.UnixNano(),
	})
	m.log.Debug().Stringer("tasks", m.tasks).Msg("tasks")

	task.SetIsScheduled(true)
}

func (m *taskManager) SuspendTask(task TaskRequirements) {
	m.log.Debug().Stringer("task", task).Msg("SuspendTask ")
	m.Lock()
	defer m.Unlock()

	deleted := false
	for i, pqi := range m.tasks {
		//when := _task.value.taskTime
		//n := _task.value.id
		curtask := pqi.value.task
		if task == curtask {
			m.log.Debug().Msg("task found")
			heap.Remove(&m.tasks, i)
			deleted = true

			task.SetIsScheduled(false)
			break
		}
	}
	if !deleted {
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

	now := m.GetTime()

	var task TaskRequirements
	var delta *time.Duration

	if len(m.tasks) > 0 {
		pqi := m.tasks[0]
		when := pqi.value.taskTime
		id := pqi.value.id
		_ = id
		nextTask := pqi.value.task
		if when.Before(now) {
			// pull it off the list and mark that it's no longer scheduled
			heap.Pop(&m.tasks)
			task = nextTask
			task.SetIsScheduled(false)

			if len(m.tasks) > 0 {
				pqi := m.tasks[0]
				when := pqi.value.taskTime
				id := pqi.value.id
				_ = id
				nextTask := pqi.value.task
				_ = nextTask

				// peek at the next task, return how long to wait
				newDelta := max(when.Sub(now), 0)
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
		task.InstallTask(WithInstallTaskOptionsNone())
	case interface{ IsOneShotDeleteTask() bool }:
		// TODO: Delete? How?
	}
}

func (m *taskManager) GetTasks() []TaskRequirements {
	m.Lock()
	defer m.Unlock()
	tasks := make([]TaskRequirements, len(m.tasks))
	for i, pqi := range m.tasks {
		tasks[i] = pqi.value.task
	}
	return tasks
}

func (m *taskManager) PopTask() TaskRequirements {
	m.log.Trace().Msg("pop task")
	m.Lock()
	defer m.Unlock()
	pqi := heap.Pop(&m.tasks).(*PriorityItem[int64, taskItem])
	return pqi.value.task
}

func (m *taskManager) CountTasks() int {
	return m.tasks.Len()
}

func (m *taskManager) String() string {
	return fmt.Sprintf("TaskManager{tasks: %v}", m.tasks)
}
