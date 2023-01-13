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
)

type _TaskRequirements interface {
	processTask() error
	InstallTask(when *time.Time, delta *time.Duration)
	getTaskTime() *time.Time
	getIsScheduled() bool
	setIsScheduled(isScheduled bool)
}

type _Task struct {
	taskRequirements _TaskRequirements
	taskTime         *time.Time
	isScheduled      bool
}

func _New_Task(_TaskRequirements _TaskRequirements) *_Task {
	return &_Task{taskRequirements: _TaskRequirements}
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
	_taskManager.installTask(t.taskRequirements)
}

func (t *_Task) SuspendTask() {
	_taskManager.suspendTask(t.taskRequirements)
}

func (t *_Task) Resume() {
	_taskManager.resumeTask(t.taskRequirements)
}

func (t *_Task) getTaskTime() *time.Time {
	return t.taskTime
}

func (t *_Task) getIsScheduled() bool {
	return t.isScheduled
}

func (t *_Task) setIsScheduled(isScheduled bool) {
	t.isScheduled = isScheduled
}

type OneShotTaskRequirements interface {
	processTask() error
}

type OneShotTask struct {
	*_Task
	OneShotTaskRequirements
}

func NewOneShotTask(oneShotTaskRequirements OneShotTaskRequirements, when *time.Time) *OneShotTask {
	o := &OneShotTask{
		OneShotTaskRequirements: oneShotTaskRequirements,
	}
	o._Task = _New_Task(o)
	if when != nil {
		o.taskTime = when
	}
	return o
}

type OneShotDeleteTask struct {
	*_Task
	OneShotTaskRequirements
}

func NewOneShotDeleteTask(oneShotTaskRequirements OneShotTaskRequirements, when *time.Time) *OneShotDeleteTask {
	o := &OneShotDeleteTask{OneShotTaskRequirements: oneShotTaskRequirements}
	o._Task = _New_Task(o)
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

func (m *OneShotFunctionTask) processTask() error {
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

	log.Debug().Msgf("task: %v", task)
	return task
}

type RecurringTaskRequirements interface {
	processTask() error
}

type RecurringTask struct {
	*_Task
	RecurringTaskRequirements
	taskInterval       *time.Duration
	taskIntervalOffset *time.Duration
}

func NewRecurringTask(recurringTaskRequirements RecurringTaskRequirements, interval *time.Duration, offset *time.Duration) *RecurringTask {
	r := &RecurringTask{RecurringTaskRequirements: recurringTaskRequirements}
	r._Task = _New_Task(r)
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
	_taskManager.installTask(r)

	return r
}

func (r *RecurringTask) IsRecurringTask() bool {
	return true
}

type _RecurringFunctionTask struct {
	*RecurringTask
	fn func() error
}

func _New_RecurringFunctionTask(interval *time.Duration, fn func() error) *_RecurringFunctionTask {
	r := &_RecurringFunctionTask{fn: fn}
	r.RecurringTask = NewRecurringTask(r, interval, nil)
	return r
}

func (r _RecurringFunctionTask) processTask() error {
	return r.fn()
}

func RecurringFunctionTask(interval *time.Duration, fn func() error) *RecurringTask {
	return _New_RecurringFunctionTask(interval, fn).RecurringTask
}

var _taskManager = TaskManager{}

func init() {
	go func() {
		for {
			task, delta := _taskManager.getNextTask()
			if task == nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			_taskManager.processTask(task)
			time.Sleep(delta)
		}
	}()
}

type TaskManager struct {
	sync.Mutex
	tasks []_TaskRequirements
}

func (m *TaskManager) getTime() time.Time {
	return time.Now()
}

func (m *TaskManager) installTask(task _TaskRequirements) {
	m.Lock()
	defer m.Unlock()
	log.Debug().Msgf("installTask %v@%v", task, task.getTaskTime())

	// if the taskTime is None is hasn't been computed correctly
	if task.getTaskTime() == nil {
		panic("task time is None")
	}

	// if this is already installed, suspend it
	if task.getIsScheduled() {
		m.suspendTask(task)
	}

	// save this in the task list
	// TODO: we might need to insert it at the right place
	m.tasks = append(m.tasks, task)

	task.setIsScheduled(true)
}

func (m *TaskManager) suspendTask(task _TaskRequirements) {
	log.Debug().Msgf("suspendTask %v", task)
	m.Lock()
	defer m.Unlock()

	iToDelete := -1
	for i, _task := range m.tasks {
		if _task == task {
			log.Debug().Msgf("task found")
			iToDelete = i
			task.setIsScheduled(false)
			break
		}
	}
	if iToDelete > 0 {
		m.tasks = append(m.tasks[:iToDelete], m.tasks[iToDelete+1:]...)
	} else {
		log.Debug().Msgf("task not found")
	}
}

func (m *TaskManager) resumeTask(task _TaskRequirements) {
	log.Debug().Msgf("resumeTask %v", task)
	m.Lock()
	defer m.Unlock()

	// just re-install it
	m.installTask(task)
}

func (m *TaskManager) getNextTask() (_TaskRequirements, time.Duration) {
	log.Trace().Msgf("getNextTask")
	m.Lock()
	defer m.Unlock()

	now := time.Now()

	var task _TaskRequirements
	var delta time.Duration

	if len(m.tasks) > 0 {
		nextTask := m.tasks[0]
		when := nextTask.getTaskTime()
		if when.Before(now) {
			// pull it off the list and mark that it's no longer scheduled
			m.tasks = m.tasks[1:] // TODO: guard against empty list
			task = nextTask
			task.setIsScheduled(false)

			if len(m.tasks) > 0 {
				nextTask = m.tasks[0]
				when = nextTask.getTaskTime()
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

func (m *TaskManager) processTask(task _TaskRequirements) {
	log.Debug().Msgf("processTask %v", task)

	// process the task
	if err := task.processTask(); err != nil {
		log.Error().Err(err).Msg("Error processing Task")
	}

	switch task.(type) {
	case interface{ IsRecurringTask() bool }:
		task.InstallTask(nil, nil)
	case interface{ IsOneShotDeleteTask() bool }:
		// TODO: Delete? How?
	}
}
