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

package time_machine

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/core"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

var globalTimeMachine *TimeMachine
var globalTimeMachineMutex sync.Mutex

func IsGlobalTimeMachineSet() bool {
	return globalTimeMachine != nil
}

// NewGlobalTimeMachine creates a new TimeMachine and set it as global.
// Usually it is sufficient to use ExclusiveGlobalTimeMachine
func NewGlobalTimeMachine(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)
	if globalTimeMachine != nil {
		testingLogger.Warn().Msg("global time machine set, overwriting")
	}
	testingLogger.Trace().Msg("creating new global time machine")
	globalTimeMachine = NewTimeMachine(testingLogger)
	testingLogger.Trace().Msg("overwriting global task manager")
	oldManager := OverwriteTaskManager(testingLogger, globalTimeMachine)
	t.Cleanup(func() {
		testingLogger.Trace().Msg("clearing task manager")
		ClearTaskManager(testingLogger)
		testingLogger.Trace().Msg("Restoring old manager")
		OverwriteTaskManager(testingLogger, oldManager)
	})
}

// ClearGlobalTimeMachine clears the global time machine during the test duration.
// Usually it is sufficient to use ExclusiveGlobalTimeMachine.
// Attention: Use in combination with LockGlobalTimeMachine to avoid side effects.
func ClearGlobalTimeMachine(t *testing.T) {
	testingLogger := testutils.ProduceTestingLogger(t)
	if globalTimeMachine == nil {
		testingLogger.Warn().Msg("global time machine not set")
	}
	globalTimeMachine = nil
	ClearTaskManager(testingLogger)
}

// LockGlobalTimeMachine locks the global time machine during the test duration.
// Usually it is sufficient to use ExclusiveGlobalTimeMachine
func LockGlobalTimeMachine(t *testing.T) {
	globalTimeMachineMutex.Lock()
	t.Cleanup(globalTimeMachineMutex.Unlock)
}

// ExclusiveGlobalTimeMachine is a combination of LockGlobalTimeMachine, NewGlobalTimeMachine and ClearGlobalTimeMachine
func ExclusiveGlobalTimeMachine(t *testing.T) {
	LockGlobalTimeMachine(t)
	NewGlobalTimeMachine(t)
	t.Cleanup(func() {
		ClearGlobalTimeMachine(t)
	})
}

type TimeMachine struct {
	TaskManager

	currentTime time.Time
	timeLimit   time.Time
	startTime   time.Time

	log zerolog.Logger
}

func NewTimeMachine(localLog zerolog.Logger) *TimeMachine {
	t := &TimeMachine{
		log: localLog,
	}
	if _debug != nil {
		_debug("__init__")
	}
	t.TaskManager = NewTaskManager(localLog)
	return t
}

func (t *TimeMachine) GetTime() time.Time {
	if _debug != nil {
		_debug("get_time @ %r:", t.currentTime)
	}
	t.log.Trace().Time("currentTime", t.currentTime).Msg("GetTime")
	return t.currentTime
}

func (t *TimeMachine) InstallTask(task TaskRequirements) {
	if _debug != nil {
		_debug("install_task @ %r: %r @ %r", t.currentTime, task, task.GetTaskTime())
	}
	t.log.Debug().Time("currentTime", t.currentTime).Stringer("task", task).Msg("InstallTask")
	t.TaskManager.InstallTask(task)
}

func (t *TimeMachine) SuspendTask(task TaskRequirements) {
	if _debug != nil {
		_debug("suspend_task @ %r: %r", t.currentTime, task)
	}
	t.log.Debug().Time("currentTime", t.currentTime).Stringer("task", task).Msg("SuspendTask")
	t.TaskManager.SuspendTask(task)
}

func (t *TimeMachine) ResumeTask(task TaskRequirements) {
	if _debug != nil {
		_debug("resume_task @ %r: %r", t.currentTime, task)
	}
	t.log.Debug().Time("currentTime", t.currentTime).Stringer("task", task).Msg("ResumeTask")
	t.TaskManager.ResumeTask(task)
}

func (t *TimeMachine) MoreToDo() bool {
	if _debug != nil {
		_debug("more_to_do @ %r:", t.currentTime)
	}
	t.log.Debug().Time("currentTime", t.currentTime).Msg("MoreToDo")
	if len(DeferredFunctions) > 0 {
		if _debug != nil {
			_debug("    - deferred functions")
		}
		t.log.Trace().Msg("deferredFunctions")
		return true
	}

	if _debug != nil {
		_debug("    - time_limit: %r", t.timeLimit)
	}
	if _debug != nil {
		_debug("    - tasks: %r", t.GetTasks())
	}

	t.log.Debug().Time("timeLimit", t.timeLimit).Msg("timeLimit")
	if t.log.Debug().Enabled() {
		stringers := make([]fmt.Stringer, len(t.GetTasks()))
		for i, task := range t.GetTasks() { //TODO: check if there is something more efficient
			stringers[i] = task
		}
		t.log.Debug().Stringers("tasks", stringers).Msg("tasks")
	}

	if !t.timeLimit.IsZero() && t.currentTime.After(t.timeLimit) {
		if _debug != nil {
			_debug("    - time limit reached or exceeded")
		}
		t.log.Trace().Msg("time limit reached")
		return false
	}

	if len(t.GetTasks()) == 0 {
		if _debug != nil {
			_debug("    - no more tasks")
		}
		t.log.Trace().Msg("no more tasks")
		return false
	}

	task := t.GetTasks()[0]
	when := task.GetTaskTime()
	if when.After(t.timeLimit) {
		if _debug != nil {
			_debug("    - next task at or exceeds time limit")
		}
		t.log.Debug().Msg("New tasks exceeded time limit")
		return false
	}
	if _debug != nil {
		_debug("    - task: %r", task)
	}
	t.log.Debug().Stringer("task", task).Msg("task")
	return true
}

func (t *TimeMachine) GetNextTask() (TaskRequirements, *time.Duration) {
	if _debug != nil {
		_debug("get_next_task @ %r:", t.currentTime)
	}
	if _debug != nil {
		_debug("    - time_limit: %r", t.timeLimit)
	}
	if _debug != nil {
		_debug("    - tasks: %r", t.GetTasks())
	}
	t.log.Debug().Time("currentTime", t.currentTime).Msg("GetNextTask")
	t.log.Debug().Time("timeLimit", t.timeLimit).Msg("timeLimit")
	if t.log.Debug().Enabled() {
		stringers := make([]fmt.Stringer, len(t.GetTasks()))
		for i, task := range t.GetTasks() { //TODO: check if there is something more efficient
			stringers[i] = task
		}
		t.log.Debug().Stringers("tasks", stringers).Msg("tasks")
	}

	var task TaskRequirements
	var delta *time.Duration

	if !t.timeLimit.IsZero() && t.currentTime.After(t.timeLimit) {
		if _debug != nil {
			_debug("    - time limit reached")
		}
		t.log.Trace().Msg("time limit reached")
	} else if len(t.GetTasks()) == 0 {
		if _debug != nil {
			_debug("    - no more tasks")
		}
		t.log.Trace().Msg("no more tasks")
	} else {
		// peek at the next task and see when it is supposed to run
		task = t.GetTasks()[0]
		if taskTime := task.GetTaskTime(); taskTime != nil && task.GetTaskTime().After(t.timeLimit) {
			if _debug != nil {
				_debug("    - time limit reached")
			}

			t.currentTime = *taskTime
		} else {
			// pull it off the list
			task = t.PopTask()
			if _debug != nil {
				_debug("    - when, task: %r, %s", task.GetTaskTime(), task)
			}
			t.log.Debug().Stringer("task", task).Msg("when, task")

			// mark that it is no longer scheduled
			task.SetIsScheduled(false)

			// advance the time
			if taskTime := task.GetTaskTime(); taskTime != nil {
				t.currentTime = *taskTime
			}

			// do not wait, time has moved
			var newDelta time.Duration
			delta = &newDelta
		}
	}
	return task, delta
}

func (t *TimeMachine) ProcessTask(task TaskRequirements) {
	if _debug != nil {
		_debug("process_task @ %r: %r", t.currentTime, task)
	}
	t.log.Debug().Time("currentTime", t.currentTime).Stringer("task", task).Msg("ProcessTask")
	t.TaskManager.ProcessTask(task)
}

func (t *TimeMachine) String() string {
	return fmt.Sprintf("TimeMachine(%s, currentTime:%s, timeLimit: %s, startTime, %s)", t.TaskManager, t.currentTime, t.timeLimit, t.startTime)
}

// ResetTimeMachine This function is called to reset the clock before running a set of tests.
func ResetTimeMachine(startTime time.Time) {
	if globalTimeMachine == nil {
		panic("no time machine")
	}

	if _debug != nil {
		_debug("reset_time_machine %r", startTime)
	}
	globalTimeMachine.ClearTasks()
	globalTimeMachine.currentTime = startTime
	globalTimeMachine.timeLimit = time.Time{}
}

// RunTimeMachine This function is called after a set of tasks have been installed
//
//	and they should Run.  The machine will stop when the stop time has been
//	reached (maybe the middle of some tests) and can be called again to
//	continue running.
func RunTimeMachine(localLog zerolog.Logger, duration time.Duration, stopTime time.Time) {
	if globalTimeMachine == nil {
		panic("no time machine")
	}
	if _debug != nil {
		_debug("run_time_machine %r %r", duration, stopTime)
	}
	localLog.Debug().Dur("duration", duration).Time("stopTime", stopTime).Msg("RunTimeMachine")

	/* TODO: we don't have a proper tristate, maybe we change currentTime to a pointer
	if !globalTimeMachine.currentTime.IsZero() {
		panic("Reset the time machine before running")
	}*/

	if duration != 0 {
		globalTimeMachine.timeLimit = globalTimeMachine.currentTime.Add(duration)
	} else if !stopTime.IsZero() {
		if _debug != nil {
			_debug("    - stop_time: %r", stopTime)
		}
		globalTimeMachine.timeLimit = stopTime
	} else {
		panic("duration or stopTime is required")
	}

	if len(DeferredFunctions) > 0 {
		if _debug != nil {
			_debug("    - deferred functions!")
		}
		localLog.Debug().Msg("deferredFunctions")
	}

	for {
		RunOnce(localLog)
		localLog.Trace().Msg("ran once")
		if _debug != nil {
			_debug("    - ran once")
		}
		if !globalTimeMachine.MoreToDo() {
			localLog.Trace().Msg("no more to do")
			if _debug != nil {
				_debug("    - no more to do")
			}
			break
		}
	}

	globalTimeMachine.currentTime = globalTimeMachine.timeLimit
}

// GlobalTimeMachineCurrentTime Return the current time from the time machine.
func GlobalTimeMachineCurrentTime() time.Time {
	return globalTimeMachine.currentTime
}
