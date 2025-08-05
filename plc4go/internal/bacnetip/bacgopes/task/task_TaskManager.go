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
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
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

// Deprecated: should only be used by core... find a better way
func GetTaskManager() TaskManager {
	return _taskManager
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

//go:generate go tool plc4xGenerator -type=taskItem -prefix=task_
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
	if _debug != nil {
		_debug("__init__")
	}
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

func (t *taskManager) GetTime() time.Time {
	if _debug != nil {
		_debug("get_time")
	}
	return time.Now()
}

func GetTaskManagerTime() time.Time {
	if _taskManager == nil {
		return time.Now()
	}
	return _taskManager.GetTime()
}

func (t *taskManager) ClearTasks() {
	t.tasks.Clear()
}

func (t *taskManager) InstallTask(task TaskRequirements) {
	if _debug != nil {
		_debug("install_task %r @ %r", task, task.GetTaskTime())
	}
	t.log.Debug().Stringer("task", task).Msg("InstallTask")
	t.Lock()
	defer t.Unlock()

	// if the taskTime is None it hasn't been computed correctly
	if task.GetTaskTime() == nil {
		panic("task time is None")
	}

	// if this is already installed, suspend it
	if task.GetIsScheduled() {
		t.Unlock()
		t.SuspendTask(task)
		t.Lock()
	}

	// save this in the task list
	t.count++
	heap.Push(&t.tasks, &PriorityItem[int64, *taskItem]{
		Value:    &taskItem{taskTime: task.GetTaskTime(), id: t.count, task: task},
		Priority: task.GetTaskTime().UnixNano() - time.Time{}.UnixNano(),
	})
	if _debug != nil {
		_debug("    - tasks: %r", t.tasks)
	}
	t.log.Debug().Stringer("tasks", t.tasks).Msg("tasks")

	task.SetIsScheduled(true)
}

func (t *taskManager) SuspendTask(task TaskRequirements) {
	if _debug != nil {
		_debug("suspend_task %r", task)
	}
	t.log.Debug().Stringer("task", task).Msg("SuspendTask ")
	t.Lock()
	defer t.Unlock()

	deleted := false
	for i, pqi := range t.tasks {
		//when := _value.taskTime
		//n := _value.id
		curtask := pqi.Value.task
		if task == curtask {
			if _debug != nil {
				_debug("    - task found")
			}
			t.log.Debug().Msg("task found")
			heap.Remove(&t.tasks, i)
			deleted = true

			task.SetIsScheduled(false)
			break
		}
	}
	if !deleted {
		if _debug != nil {
			_debug("    - task not found")
		}
		t.log.Debug().Msg("task not found")
	}
}

func (t *taskManager) ResumeTask(task TaskRequirements) {
	if _debug != nil {
		_debug("resume_task %r", task)
	}
	t.log.Debug().Stringer("task", task).Msg("ResumeTask")
	t.Lock()
	defer t.Unlock()

	// just re-install it
	t.InstallTask(task)
}

func (t *taskManager) GetNextTask() (TaskRequirements, *time.Duration) {
	if _debug != nil {
		_debug("get_next_task")
	}
	t.log.Trace().Msg("GetNextTask")
	t.Lock()
	defer t.Unlock()

	now := t.GetTime()

	var task TaskRequirements
	var delta *time.Duration

	if len(t.tasks) > 0 {
		pqi := t.tasks[0]
		when := pqi.Value.taskTime
		id := pqi.Value.id
		_ = id
		nextTask := pqi.Value.task
		if when.Before(now) {
			// pull it off the list and mark that it's no longer scheduled
			heap.Pop(&t.tasks)
			task = nextTask
			task.SetIsScheduled(false)

			if len(t.tasks) > 0 {
				pqi := t.tasks[0]
				when := pqi.Value.taskTime
				id := pqi.Value.id
				_ = id
				nextTask := pqi.Value.task
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

func (t *taskManager) ProcessTask(task TaskRequirements) {
	if _debug != nil {
		_debug("process_task %r", task)
	}
	t.log.Debug().Stringer("task", task).Msg("ProcessTask")

	// process the task
	if err := task.ProcessTask(); err != nil {
		t.log.Error().Err(err).Msg("Error processing Task")
	}

	switch task.(type) {
	case interface{ IsRecurringTask() bool }:
		task.InstallTask(WithInstallTaskOptionsNone())
	case interface{ IsOneShotDeleteTask() bool }:
		// TODO: Delete? How?
	}
}

func (t *taskManager) GetTasks() []TaskRequirements {
	t.Lock()
	defer t.Unlock()
	tasks := make([]TaskRequirements, len(t.tasks))
	for i, pqi := range t.tasks {
		tasks[i] = pqi.Value.task
	}
	return tasks
}

func (t *taskManager) PopTask() TaskRequirements {
	t.log.Trace().Msg("pop task")
	t.Lock()
	defer t.Unlock()
	pqi := heap.Pop(&t.tasks).(*PriorityItem[int64, *taskItem])
	return pqi.Value.task
}

func (t *taskManager) CountTasks() int {
	return t.tasks.Len()
}

func (t *taskManager) String() string {
	return fmt.Sprintf("TaskManager{tasks: %v}", t.tasks)
}
