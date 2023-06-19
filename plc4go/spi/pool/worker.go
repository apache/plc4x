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
	"github.com/rs/zerolog/log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=worker
type worker struct {
	id       int
	executor interface {
		isTraceWorkers() bool
		getWorksItems() chan workItem
		getWorkerWaitGroup() *sync.WaitGroup
	}

	lastReceived atomic.Value

	stateChange sync.Mutex
	running     atomic.Bool
	shutdown    atomic.Bool
	interrupted atomic.Bool
	interrupter chan struct{}

	log zerolog.Logger `ignore:"true"`
}

func (w *worker) initialize() {
	w.shutdown.Store(false)
	w.interrupted.Store(false)
	w.interrupter = make(chan struct{}, 1)
	w.lastReceived.Store(time.Now())
}

func (w *worker) start() {
	w.stateChange.Lock()
	defer w.stateChange.Unlock()
	if w.running.Load() {
		log.Warn().Msg("Worker already started")
		return
	}
	if w.executor.isTraceWorkers() {
		w.log.Debug().Msgf("Starting worker\n%s", w)
	}
	w.executor.getWorkerWaitGroup().Add(1)
	go w.work()
}

func (w *worker) stop(interrupt bool) {
	w.stateChange.Lock()
	defer w.stateChange.Unlock()
	if !w.running.Load() {
		log.Warn().Msg("Worker not running")
		return
	}
	w.shutdown.Store(true)
	if interrupt {
		w.interrupted.Store(true)
		close(w.interrupter)
	}
}

func (w *worker) work() {
	defer w.executor.getWorkerWaitGroup().Done()
	defer func() {
		if err := recover(); err != nil {
			w.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
			if !w.shutdown.Load() {
				// if we are not in shutdown we continue
				w.start()
			}
		}
	}()
	w.running.Store(true)
	defer w.running.Store(false)
	workerLog := w.log.With().Int("Worker id", w.id).Logger()
	if !w.executor.isTraceWorkers() {
		workerLog = zerolog.Nop()
	}

	for !w.shutdown.Load() {
		workerLog.Trace().Msg("Working")
		select {
		case _workItem := <-w.executor.getWorksItems():
			w.lastReceived.Store(time.Now())
			workerLog.Debug().Msgf("Got work item %v", _workItem)
			if _workItem.completionFuture.cancelRequested.Load() || (w.shutdown.Load() && w.interrupted.Load()) {
				workerLog.Debug().Msg("We need to stop")
				// TODO: do we need to complete with a error?
			} else {
				workerLog.Debug().Msgf("Running work item %v", _workItem)
				_workItem.runnable()
				_workItem.completionFuture.complete()
				workerLog.Debug().Msgf("work item %v completed", _workItem)
			}
		case <-w.interrupter:
			workerLog.Debug().Msg("We got interrupted")
		}
	}
	workerLog.Trace().Msg("done")
}
