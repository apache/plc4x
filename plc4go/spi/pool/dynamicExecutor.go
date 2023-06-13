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

package pool

import (
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/rs/zerolog"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

var upScaleInterval = 100 * time.Millisecond
var downScaleInterval = 5 * time.Second
var timeToBecomeUnused = 5 * time.Second

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=dynamicExecutor
type dynamicExecutor struct {
	*executor

	maxNumberOfWorkers     int
	currentNumberOfWorkers atomic.Int32
	dynamicStateChange     sync.Mutex
	interrupter            chan struct{}

	dynamicWorkers sync.WaitGroup
}

func (e *dynamicExecutor) Start() {
	e.dynamicStateChange.Lock()
	defer e.dynamicStateChange.Unlock()
	if e.running || e.shutdown {
		e.log.Warn().Msg("Already started")
		return
	}
	if e.interrupter != nil {
		e.log.Debug().Msg("Ensuring that the old spawner/killers are not running")
		close(e.interrupter)
		e.dynamicWorkers.Wait()
	}

	e.executor.Start()
	mutex := sync.Mutex{}
	e.interrupter = make(chan struct{})
	// Worker spawner
	go func() {
		e.dynamicWorkers.Add(1)
		defer e.dynamicWorkers.Done()
		defer func() {
			if err := recover(); err != nil {
				e.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
			}
		}()
		workerLog := e.log.With().Str("Worker type", "spawner").Logger()
		if !e.traceWorkers {
			workerLog = zerolog.Nop()
		}
		for e.running && !e.shutdown {
			workerLog.Trace().Msg("running")
			mutex.Lock()
			numberOfItemsInQueue := len(e.workItems)
			numberOfWorkers := len(e.worker)
			workerLog.Debug().Msgf("Checking if numberOfItemsInQueue(%d) > numberOfWorkers(%d) && numberOfWorkers(%d) < maxNumberOfWorkers(%d)", numberOfItemsInQueue, numberOfWorkers, numberOfWorkers, e.maxNumberOfWorkers)
			if numberOfItemsInQueue > numberOfWorkers && numberOfWorkers < e.maxNumberOfWorkers {
				workerLog.Trace().Msg("spawning new worker")
				_worker := &worker{
					id:           numberOfWorkers - 1,
					interrupter:  make(chan struct{}, 1),
					executor:     e,
					lastReceived: time.Now(),
					log:          e.log,
				}
				e.worker = append(e.worker, _worker)
				_worker.initialize()
				workerLog.Info().Int("Worker id", _worker.id).Msg("spawning")
				go _worker.work()
				e.currentNumberOfWorkers.Add(1)
			} else {
				workerLog.Trace().Msg("Nothing to scale")
			}
			mutex.Unlock()
			func() {
				workerLog.Debug().Msgf("Sleeping for %v", upScaleInterval)
				timer := time.NewTimer(upScaleInterval)
				defer utils.CleanupTimer(timer)
				select {
				case <-timer.C:
				case <-e.interrupter:
					workerLog.Info().Msg("interrupted")
				}
			}()
		}
		workerLog.Info().Msg("Terminated")
	}()
	// Worker killer
	go func() {
		e.dynamicWorkers.Add(1)
		defer e.dynamicWorkers.Done()
		defer func() {
			if err := recover(); err != nil {
				e.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
			}
		}()
		workerLog := e.log.With().Str("Worker type", "killer").Logger()
		if !e.traceWorkers {
			workerLog = zerolog.Nop()
		}
		for e.running && !e.shutdown {
			workerLog.Trace().Msg("running")
			mutex.Lock()
			newWorkers := make([]*worker, 0)
			for _, _worker := range e.worker {
				deadline := time.Now().Add(-timeToBecomeUnused)
				workerLog.Debug().Int("Worker id", _worker.id).Msgf("Checking if %v is before %v", _worker.lastReceived, deadline)
				if _worker.lastReceived.Before(deadline) {
					workerLog.Info().Int("Worker id", _worker.id).Msg("killing")
					_worker.interrupted.Store(true)
					close(_worker.interrupter)
					e.currentNumberOfWorkers.Add(-1)
				} else {
					workerLog.Debug().Int("Worker id", _worker.id).Msg("still ok")
					newWorkers = append(newWorkers, _worker)
				}
			}
			e.worker = newWorkers
			mutex.Unlock()
			func() {
				workerLog.Debug().Msgf("Sleeping for %v", downScaleInterval)
				timer := time.NewTimer(downScaleInterval)
				defer utils.CleanupTimer(timer)
				select {
				case <-timer.C:
				case <-e.interrupter:
					workerLog.Info().Msg("interrupted")
				}
			}()
		}
		workerLog.Info().Msg("Terminated")
	}()
}

func (e *dynamicExecutor) Stop() {
	e.log.Trace().Msg("stopping now")
	e.dynamicStateChange.Lock()
	defer e.dynamicStateChange.Unlock()
	if !e.running || e.shutdown {
		e.log.Warn().Msg("already stopped")
		return
	}
	close(e.interrupter)
	e.log.Trace().Msg("stopping inner executor")
	e.executor.Stop()
	e.log.Debug().Msgf("waiting for %d dynamic workers to stop", e.currentNumberOfWorkers.Load())
	e.dynamicWorkers.Wait()
	e.log.Trace().Msg("stopped")
}
