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
	"github.com/rs/zerolog/log"
	"sync"
	"time"
	"unsafe"
)

type _TaskRequirements interface {
	processTask()
}

type _Task struct {
	_TaskRequirements
	taskTime    *time.Time
	isScheduled bool
}

func _New_Task(_TaskRequirements _TaskRequirements) *_Task {
	return &_Task{_TaskRequirements: _TaskRequirements}
}

func (t *_Task) InstallTask(when *time.Time, delta *time.Duration) {
	// check for delta from now
	if when == nil && delta != nil {
		_when := _taskManager.getTime().Add(*delta)
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
	_taskManager.installTask(t)
}

func (t *_Task) processTask() {
	panic("processTask must be overridden")
}

func (t *_Task) SuspendTask() {
	_taskManager.suspendTask(t)
}

func (t *_Task) Resume() {
	_taskManager.resumeTask(t)
}

func (t *_Task) LowerThan(other *_Task) bool {
	return *unsafe.Pointer(t) < *unsafe.Pointer(other)
}

type OneShotTask struct {
	*_Task
}

func NewOneShotTask(when *time.Time) *OneShotTask {
	o := &OneShotTask{}
	o._TaskRequirements = o
	if when != nil {
		o.taskTime = when
	}
	return o
}

type OneShotDeleteTask struct {
	*_Task
}

func NewOneShotDeleteTask(when *time.Time) *OneShotDeleteTask {
	o := &OneShotDeleteTask{}
	o._TaskRequirements = o
	if when != nil {
		o.taskTime = when
	}
	return o
}

type OneShotFunctionTask struct {
	*OneShotDeleteTask
	fn func()
}

func (m *OneShotFunctionTask) processTask() {
	m.fn()
}

func OneShotFunction(fn func()) *OneShotFunctionTask {
	task := &OneShotFunctionTask{NewOneShotDeleteTask(nil), fn}

	var delta time.Duration = 0
	task.InstallTask(nil, &delta)
	return task
}

func FunctionTask(fn func()) *OneShotFunctionTask {
	task := &OneShotFunctionTask{NewOneShotDeleteTask(nil), fn}

	log.Debug().Msgf("task: %v", task)
	return task
}

type RecurringTask struct {
	*_Task
	taskInterval       *time.Duration
	taskIntervalOffset *time.Duration
}

func NewRecurringTask(_TaskRequirements _TaskRequirements, interval *time.Duration, offset *time.Duration) *RecurringTask {
	r := &RecurringTask{}
	r._TaskRequirements = _New_Task(_TaskRequirements)
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
	now := _taskManager.getTime().Add(10 + time.Nanosecond)

	if r.taskIntervalOffset != nil {
		_offset := *r.taskIntervalOffset
		offset = &_offset
	} else {
		_offset := time.Duration(0)
		offset = &_offset
	}
	log.Debug().Msgf("Now, interval, offset: %v, %v, %v", now, interval, offset)

	// compute the time
	_taskTime := now.Add(-*offset).Add(*interval) // TODO: check why upstream is doing the modulo operation (missing code here)
	r.taskTime = &_taskTime

	// install it
	_taskManager.installTask(r._Task)

	return r
}

type _RecurringFunctionTask struct {
	*RecurringTask
	fn func()
}

func _New_RecurringFunctionTask(interval *time.Duration) *_RecurringFunctionTask {
	r := &_RecurringFunctionTask{}
	r.RecurringTask = NewRecurringTask(r, interval, nil)
	return r
}

func (r _RecurringFunctionTask) processTask() {
	r.fn()
}

func RecurringFunctionTask(interval *time.Duration) *RecurringTask {
	return _New_RecurringFunctionTask(interval).RecurringTask
}

var _taskManager = TaskManager{}

func init() {
	go func() {
		for {
			task, delta := _taskManager.getNextTask()
			_taskManager.processTask(task)
			time.Sleep(delta)
		}
	}()
}

type TaskManager struct {
	sync.Mutex
	tasks []*_Task
}

func (m *TaskManager) getTime() time.Time {
	return time.Now()
}

func (m *TaskManager) installTask(task *_Task) {
	m.Lock()
	defer m.Unlock()
	log.Debug().Msgf("installTask %v@%v", task, task.taskTime)

	// if the taskTime is None is hasn't been computed correctly
	if task.taskTime == nil {
		panic("task time is None")
	}

	// if this is already installed, suspend it
	if task.isScheduled {
		m.suspendTask(task)
	}

	// save this in the task list
	// TODO: we might need to insert it at the right place
	m.tasks = append(m.tasks, task)

	task.isScheduled = true
}

func (m *TaskManager) suspendTask(task *_Task) {
	log.Debug().Msgf("suspendTask %v", task)
	m.Lock()
	defer m.Unlock()

	iToDelete := -1
	for i, _task := range m.tasks {
		if _task == task {
			log.Debug().Msgf("task found")
			iToDelete = i
			task.isScheduled = false
			break
		}
	}
	if iToDelete > 0 {
		m.tasks = append(m.tasks[:iToDelete], m.tasks[iToDelete+1:]...)
	} else {
		log.Debug().Msgf("task not found")
	}
}

func (m *TaskManager) resumeTask(task *_Task) {
	log.Debug().Msgf("resumeTask %v", task)
	m.Lock()
	defer m.Unlock()

	// just re-install it
	m.installTask(task)
}

func (m *TaskManager) getNextTask() (*_Task, time.Duration) {
	log.Debug().Msgf("getNextTask")
	m.Lock()
	defer m.Unlock()

	now := time.Now()

	var task *_Task
	var delta time.Duration

	if m.tasks != nil {
		nextTask := m.tasks[0]
		when := nextTask.taskTime
		if when.Before(now) {
			// pull it off the list and mark that it's no longer scheduled
			m.tasks = m.tasks[1:] // TODO: guard against empty list
			task = nextTask
			task.isScheduled = false

			if m.tasks != nil {
				nextTask = m.tasks[0]
				when = nextTask.taskTime
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

func (m *TaskManager) processTask(task *_Task) {
	log.Debug().Msgf("processTask %v", task)

	// process the task
	task.processTask()

	switch task._TaskRequirements.(type) {
	case *RecurringTask:
		task.InstallTask(nil, nil)
	case OneShotDeleteTask:
		// TODO: Delete? How?
	}
}
