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

package core

import (
	"math"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
)

var running bool
var spin = 10 * time.Millisecond
var sleepTime = 0 * time.Nanosecond

type deferredFunctionTuple struct {
	fn     GenericFunction
	args   Args
	kwArgs KWArgs
}

var DeferredFunctions []deferredFunctionTuple

var ErrorCallback func(err error)

var wg sync.WaitGroup // use to track spawned go routines

func run() {
	running = true
	wg.Add(1)
	go func() {
		defer wg.Done()
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		running = false
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for running {
			// get the next task
			var delta time.Duration
			task, taskDelta := GetTaskManager().GetNextTask()
			if task != nil {
				GetTaskManager().ProcessTask(task)
			}

			// if delta is None, there are no Tasks, default to spinning
			if taskDelta == nil {
				delta = spin
			} else {
				delta = *taskDelta
			}

			// there may be threads around, sleep for a bit
			if sleepTime > 0 && delta > sleepTime {
				time.Sleep(sleepTime)
				delta -= sleepTime
			}

			// delta should be no more than the spin value
			delta = time.Duration(math.Min(float64(delta), float64(spin)))

			// if there are deferred functions, use a small delta
			if len(DeferredFunctions) > 0 {
				delta = time.Duration(math.Min(float64(delta), float64(1*time.Millisecond)))
			}

			// wait for socket
			time.Sleep(delta)

			// check for deferred functions
			for len(DeferredFunctions) > 0 {
				fnlist := DeferredFunctions
				// empty list
				DeferredFunctions = nil
				for _, fnTuple := range fnlist {
					fn := fnTuple.fn
					args := fnTuple.args
					kwArgs := fnTuple.kwArgs
					if err := fn(args, kwArgs); err != nil {
						if ErrorCallback != nil {
							ErrorCallback(err)
						}
					}
				}
			}
		}
	}()
}

// RunOnce makes a pass through the scheduled tasks and deferred functions just
//
//	like the run() function but without the asyncore call (so there is no
//	socket IO activity) and the timers.
func RunOnce(localLog zerolog.Logger) {
	localLog.Trace().Msg("run_once")
	taskManager := NewTaskManager(localLog)

	for {
		// get the next task
		var task TaskRequirements
		task, delta := taskManager.GetNextTask()
		var displayDelta time.Duration
		if delta != nil {
			displayDelta = *delta
		}
		localLog.Debug().Stringer("task", task).Dur("delta", displayDelta).Msg("task")

		// if there is a task to process, do it
		if task != nil {
			taskManager.ProcessTask(task)
		}

		// check for deferred functions

		for len(DeferredFunctions) > 0 {
			// get a reference to the list
			fnlist := DeferredFunctions
			DeferredFunctions = nil

			// call the functions
			for _, fnTuple := range fnlist {
				fn := fnTuple.fn
				args := fnTuple.args
				kwArgs := fnTuple.kwArgs
				if err := fn(args, kwArgs); err != nil {
					if ErrorCallback != nil {
						ErrorCallback(err)
					}
				}
			}
		}
		if delta == nil || *delta != 0 {
			break
		}
	}
}

func Deferred(fn GenericFunction, args Args, kwArgs KWArgs) {
	// append it to the list
	DeferredFunctions = append(DeferredFunctions, deferredFunctionTuple{fn, args, kwArgs})

	// trigger the task manager event
	// TODO: there is no trigger
}
