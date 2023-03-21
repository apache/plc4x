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

package utils

import (
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type Runnable func()

type Worker struct {
	id          int
	shutdown    bool
	runnable    Runnable
	interrupted bool
	executor    *Executor
}

func (w *Worker) work() {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Error().Msgf("Recovering from panic()=%v", recovered)
		}
		if !w.shutdown {
			// if we are not in shutdown we continue
			w.work()
		}
	}()
	workerLog := log.With().Int("Worker id", w.id).Logger()
	if !config.TraceTransactionManagerWorkers {
		workerLog = zerolog.Nop()
	}

	for !w.shutdown {
		workerLog.Debug().Msg("Working")
		select {
		case workItem := <-w.executor.queue:
			workerLog.Debug().Msgf("Got work item %v", workItem)
			if workItem.completionFuture.cancelRequested || (w.shutdown && w.interrupted) {
				workerLog.Debug().Msg("We need to stop")
				// TODO: do we need to complete with a error?
			} else {
				workerLog.Debug().Msgf("Running work item %v", workItem)
				workItem.runnable()
				workItem.completionFuture.Complete()
				workerLog.Debug().Msgf("work item %v completed", workItem)
			}
		default:
			workerLog.Debug().Msgf("Idling")
			time.Sleep(time.Millisecond * 10)
		}
	}
}

type WorkItem struct {
	transactionId    int32
	runnable         Runnable
	completionFuture *CompletionFuture
}

func (w *WorkItem) String() string {
	return fmt.Sprintf("Workitem{tid:%d}", w.transactionId)
}

type Executor struct {
	running     bool
	shutdown    bool
	stateChange sync.Mutex
	worker      []*Worker
	queue       chan WorkItem
}

func NewFixedSizeExecutor(numberOfWorkers int) *Executor {
	workers := make([]*Worker, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		workers[i] = &Worker{
			id:          i,
			shutdown:    false,
			runnable:    nil,
			interrupted: false,
			executor:    nil,
		}
	}
	executor := Executor{
		queue:  make(chan WorkItem, 100),
		worker: workers,
	}
	for i := 0; i < numberOfWorkers; i++ {
		worker := workers[i]
		worker.executor = &executor
	}
	return &executor
}

func (e *Executor) Submit(transactionId int32, runnable Runnable) *CompletionFuture {
	log.Trace().Int32("transactionId", transactionId).Msg("Submitting runnable")
	completionFuture := &CompletionFuture{}
	// TODO: add select and timeout if queue is full
	e.queue <- WorkItem{
		transactionId:    transactionId,
		runnable:         runnable,
		completionFuture: completionFuture,
	}
	log.Trace().Int32("transactionId", transactionId).Msg("runnable queued")
	return completionFuture
}

func (e *Executor) Start() {
	e.stateChange.Lock()
	defer e.stateChange.Unlock()
	if e.running {
		return
	}
	e.running = true
	e.shutdown = false
	for i := 0; i < len(e.worker); i++ {
		worker := e.worker[i]
		go worker.work()
	}
}

func (e *Executor) Stop() {
	e.stateChange.Lock()
	defer e.stateChange.Unlock()
	if !e.running {
		return
	}
	e.shutdown = true
	close(e.queue)
	for i := 0; i < len(e.worker); i++ {
		worker := e.worker[i]
		worker.shutdown = true
		worker.interrupted = true
	}
	e.running = false
}

type CompletionFuture struct {
	cancelRequested    bool
	interruptRequested bool
	completed          bool
	errored            bool
	err                error
}

func (f *CompletionFuture) Cancel(interrupt bool, err error) {
	f.cancelRequested = true
	f.interruptRequested = interrupt
	f.errored = true
	f.err = err
}

func (f *CompletionFuture) Complete() {
	f.completed = true
}

func (f *CompletionFuture) AwaitCompletion() error {
	for !f.completed && !f.errored {
		time.Sleep(time.Millisecond * 10)
	}
	return f.err
}
