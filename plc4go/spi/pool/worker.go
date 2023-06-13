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
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

type worker struct {
	id          int
	shutdown    atomic.Bool
	interrupted atomic.Bool
	interrupter chan struct{}
	executor    interface {
		isTraceWorkers() bool
		getWorksItems() chan workItem
		getWorkerWaitGroup() *sync.WaitGroup
	}
	hasEnded     atomic.Bool
	lastReceived time.Time

	log zerolog.Logger
}

func (w *worker) initialize() {
	w.shutdown.Store(false)
	w.interrupted.Store(false)
	w.interrupter = make(chan struct{}, 1)
	w.hasEnded.Store(false)
	w.lastReceived = time.Now()
}

func (w *worker) work() {
	w.executor.getWorkerWaitGroup().Add(1)
	defer w.executor.getWorkerWaitGroup().Done()
	defer func() {
		if err := recover(); err != nil {
			w.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
		}
		if !w.shutdown.Load() {
			// if we are not in shutdown we continue
			w.work()
		}
	}()
	workerLog := w.log.With().Int("Worker id", w.id).Logger()
	if !w.executor.isTraceWorkers() {
		workerLog = zerolog.Nop()
	}
	workerLog.Debug().Msgf("current ended state: %t", w.hasEnded.Load())
	w.hasEnded.Store(false)
	workerLog.Debug().Msgf("setting to not ended")

	for !w.shutdown.Load() {
		workerLog.Debug().Msg("Working")
		select {
		case _workItem := <-w.executor.getWorksItems():
			w.lastReceived = time.Now()
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
	w.hasEnded.Store(true)
	workerLog.Debug().Msg("setting to ended")
}

func (w *worker) String() string {
	return fmt.Sprintf("worker{\n"+
		"\tid: %d,\n"+
		"\tshutdown: %v,\n"+
		"\tinterrupted: %t,\n"+
		"\thasEnded: %t,\n"+
		"\tlastReceived: %s,\n"+
		"}",
		w.id,
		w.shutdown.Load(),
		w.interrupted.Load(),
		w.hasEnded.Load(),
		w.lastReceived,
	)
}
