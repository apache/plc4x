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
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"sync/atomic"
	"time"
)

type Runnable func()

type Worker struct {
	id          int
	shutdown    atomic.Bool
	runnable    Runnable
	interrupted atomic.Bool
	executor    *Executor
}

func (w *Worker) work() {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Error().Msgf("Recovering from panic()=%v", recovered)
		}
		if !w.shutdown.Load() {
			// if we are not in shutdown we continue
			w.work()
		}
	}()
	workerLog := log.With().Int("Worker id", w.id).Logger()
	if !w.executor.traceWorkers {
		workerLog = zerolog.Nop()
	}

	for !w.shutdown.Load() {
		workerLog.Debug().Msg("Working")
		select {
		case workItem := <-w.executor.queue:
			workerLog.Debug().Msgf("Got work item %v", workItem)
			if workItem.completionFuture.cancelRequested.Load() || (w.shutdown.Load() && w.interrupted.Load()) {
				workerLog.Debug().Msg("We need to stop")
				// TODO: do we need to complete with a error?
			} else {
				workerLog.Debug().Msgf("Running work item %v", workItem)
				workItem.runnable()
				workItem.completionFuture.complete()
				workerLog.Debug().Msgf("work item %v completed", workItem)
			}
		default:
			workerLog.Debug().Msgf("Idling")
			time.Sleep(time.Millisecond * 10)
		}
	}
}

type WorkItem struct {
	workItemId       int32
	runnable         Runnable
	completionFuture *future
}

func (w *WorkItem) String() string {
	return fmt.Sprintf("Workitem{wid:%d}", w.workItemId)
}

type Executor struct {
	running      bool
	shutdown     bool
	stateChange  sync.Mutex
	worker       []*Worker
	queue        chan WorkItem
	traceWorkers bool
}

func NewFixedSizeExecutor(numberOfWorkers, queueDepth int, options ...ExecutorOption) *Executor {
	workers := make([]*Worker, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		workers[i] = &Worker{
			id: i,
		}
	}
	executor := &Executor{
		queue:  make(chan WorkItem, queueDepth),
		worker: workers,
	}
	for _, option := range options {
		option(executor)
	}
	for i := 0; i < numberOfWorkers; i++ {
		workers[i].executor = executor
	}
	return executor
}

type ExecutorOption func(*Executor)

func WithExecutorOptionTracerWorkers(traceWorkers bool) ExecutorOption {
	return func(executor *Executor) {
		executor.traceWorkers = traceWorkers
	}
}

func (e *Executor) Submit(ctx context.Context, workItemId int32, runnable Runnable) CompletionFuture {
	log.Trace().Int32("workItemId", workItemId).Msg("Submitting runnable")
	completionFuture := &future{}
	select {
	case e.queue <- WorkItem{
		workItemId:       workItemId,
		runnable:         runnable,
		completionFuture: completionFuture,
	}:
		log.Trace().Msg("Item added")
	case <-ctx.Done():
		completionFuture.Cancel(false, ctx.Err())
	}

	log.Trace().Int32("workItemId", workItemId).Msg("runnable queued")
	return completionFuture
}

func (e *Executor) Start() {
	e.stateChange.Lock()
	defer e.stateChange.Unlock()
	if e.running || e.shutdown {
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
		worker.shutdown.Store(true)
		worker.interrupted.Store(true)
	}
	e.running = false
	e.shutdown = false
}

func (e *Executor) IsRunning() bool {
	return e.running && !e.shutdown
}

type CompletionFuture interface {
	AwaitCompletion(ctx context.Context) error
	Cancel(interrupt bool, err error)
}

type future struct {
	cancelRequested    atomic.Bool
	interruptRequested atomic.Bool
	completed          atomic.Bool
	errored            atomic.Bool
	err                atomic.Value
}

func (f *future) Cancel(interrupt bool, err error) {
	f.cancelRequested.Store(true)
	f.interruptRequested.Store(interrupt)
	f.errored.Store(true)
	f.err.Store(err)
}

func (f *future) complete() {
	f.completed.Store(true)
}

func (f *future) AwaitCompletion(ctx context.Context) error {
	for !f.completed.Load() && !f.errored.Load() && ctx.Err() != nil {
		time.Sleep(time.Millisecond * 10)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return f.err.Load().(error)
}
