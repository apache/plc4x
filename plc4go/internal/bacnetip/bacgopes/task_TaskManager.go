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
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

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

//go:generate plc4xGenerator -type=taskItem -prefix=task_
type taskItem struct {
	taskTime *time.Time
	id       int
	task     TaskRequirements
}

type taskManager struct {
	sync.Mutex

	tasks PriorityQueue[int64, *taskItem]
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
	heap.Push(&m.tasks, &PriorityItem[int64, *taskItem]{
		value:    &taskItem{taskTime: task.GetTaskTime(), id: m.count, task: task},
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
	pqi := heap.Pop(&m.tasks).(*PriorityItem[int64, *taskItem])
	return pqi.value.task
}

func (m *taskManager) CountTasks() int {
	return m.tasks.Len()
}

func (m *taskManager) String() string {
	return fmt.Sprintf("TaskManager{tasks: %v}", m.tasks)
}
