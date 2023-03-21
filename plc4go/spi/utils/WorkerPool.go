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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"sync/atomic"
	"time"
)

type Runnable func()

type worker struct {
	id          int
	shutdown    atomic.Bool
	interrupted atomic.Bool
	executor    *executor
	hasEnded    bool
}

func (w *worker) work() {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Error().Msgf("Recovering from panic()=%v", recovered)
		}
		if !w.shutdown.Load() {
			// if we are not in shutdown we continue
			w.work()
		}
	}()
	w.hasEnded = false
	workerLog := log.With().Int("Worker id", w.id).Logger()
	if !w.executor.traceWorkers {
		workerLog = zerolog.Nop()
	}

	for !w.shutdown.Load() {
		workerLog.Debug().Msg("Working")
		select {
		case _workItem := <-w.executor.queue:
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
		default:
			workerLog.Debug().Msgf("Idling")
			time.Sleep(time.Millisecond * 10)
		}
	}
	w.hasEnded = true
}

type workItem struct {
	workItemId       int32
	runnable         Runnable
	completionFuture *future
}

func (w *workItem) String() string {
	return fmt.Sprintf("Workitem{wid:%d}", w.workItemId)
}

type Executor interface {
	Start()
	Stop()
	Submit(ctx context.Context, workItemId int32, runnable Runnable) CompletionFuture
	IsRunning() bool
}

type executor struct {
	running      bool
	shutdown     bool
	stateChange  sync.Mutex
	worker       []*worker
	queue        chan workItem
	traceWorkers bool
}

func NewFixedSizeExecutor(numberOfWorkers, queueDepth int, options ...ExecutorOption) Executor {
	workers := make([]*worker, numberOfWorkers)
	for i := 0; i < numberOfWorkers; i++ {
		workers[i] = &worker{
			id: i,
		}
	}
	executor := &executor{
		queue:  make(chan workItem, queueDepth),
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

type ExecutorOption func(*executor)

func WithExecutorOptionTracerWorkers(traceWorkers bool) ExecutorOption {
	return func(executor *executor) {
		executor.traceWorkers = traceWorkers
	}
}

func (e *executor) Submit(ctx context.Context, workItemId int32, runnable Runnable) CompletionFuture {
	log.Trace().Int32("workItemId", workItemId).Msg("Submitting runnable")
	completionFuture := &future{}
	select {
	case e.queue <- workItem{
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

func (e *executor) Start() {
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

func (e *executor) Stop() {
	e.stateChange.Lock()
	defer e.stateChange.Unlock()
	if !e.running || e.shutdown {
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

func (e *executor) IsRunning() bool {
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
	if err != nil {
		f.errored.Store(true)
		f.err.Store(err)
	}
}

func (f *future) complete() {
	f.completed.Store(true)
}

// Canceled is returned on CompletionFuture.AwaitCompletion when a CompletionFuture was canceled
var Canceled = errors.New("Canceled")

func (f *future) AwaitCompletion(ctx context.Context) error {
	for !f.completed.Load() && !f.errored.Load() && !f.cancelRequested.Load() && ctx.Err() == nil {
		time.Sleep(time.Millisecond * 10)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err, ok := f.err.Load().(error); ok {
		return err
	}
	if f.cancelRequested.Load() {
		return Canceled
	}
	return nil
}
