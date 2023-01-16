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
	"math"
	"time"
)

var running bool
var spin = 10 * time.Millisecond
var sleepTime = 0 * time.Nanosecond
var deferredFunctions []func() error

func init() {
	running = true
	go func() {
		for running {
			// get the next task
			task, delta := _taskManager.getNextTask()
			if task != nil {
				_taskManager.processTask(task)
			}

			// if delta is None, there are no tasks, default to spinning
			if delta == 0 {
				delta = spin
			}

			// there may be threads around, sleep for a bit
			if sleepTime > 0 && delta > sleepTime {
				time.Sleep(sleepTime)
				delta -= sleepTime
			}

			// delta should be no more than the spin value
			delta = time.Duration(math.Min(float64(delta), float64(spin)))

			// if there are deferred functions, use a small delta
			if len(deferredFunctions) > 0 {
				delta = time.Duration(math.Min(float64(delta), float64(1*time.Millisecond)))
			}

			// wait for socket
			time.Sleep(delta)

			// check for deferred functions
			fnlist := deferredFunctions
			// empty list
			deferredFunctions = nil
			for _, fn := range fnlist {
				if err := fn(); err != nil {
					log.Debug().Err(err).Msg("error executing deferred function")
				}
			}
		}
	}()
}

func Deferred(fn func() error) {
	log.Debug().Msg("Deferred")

	// append it to the list
	deferredFunctions = append(deferredFunctions, fn)

	// trigger the task manager event
	// TODO: there is no trigger
}
