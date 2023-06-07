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
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type Runnable func()

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

type workItem struct {
	workItemId       int32
	runnable         Runnable
	completionFuture *future
}

func (w workItem) String() string {
	return fmt.Sprintf("Workitem{wid:%d, runnable(%t)}, completionFuture(%v)}", w.workItemId, w.runnable != nil, w.completionFuture)
}

type Executor interface {
	io.Closer
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
	queueDepth   int
	workItems    chan workItem
	traceWorkers bool

	workerWaitGroup sync.WaitGroup

	log zerolog.Logger
}

func (e *executor) isTraceWorkers() bool {
	return e.traceWorkers
}

func (e *executor) getWorksItems() chan workItem {
	return e.workItems
}

func (e *executor) getWorkerWaitGroup() *sync.WaitGroup {
	return &e.workerWaitGroup
}

type dynamicExecutor struct {
	*executor

	maxNumberOfWorkers     int
	currentNumberOfWorkers atomic.Int32
	dynamicStateChange     sync.Mutex
	interrupter            chan struct{}

	dynamicWorkers sync.WaitGroup
}

func NewFixedSizeExecutor(numberOfWorkers, queueDepth int, _options ...options.WithOption) Executor {
	workers := make([]*worker, numberOfWorkers)
	customLogger := options.ExtractCustomLogger(_options...)
	for i := 0; i < numberOfWorkers; i++ {
		workers[i] = &worker{
			id:  i,
			log: customLogger,
		}
	}
	_executor := &executor{
		queueDepth: queueDepth,
		workItems:  make(chan workItem, queueDepth),
		worker:     workers,
		log:        customLogger,
	}
	for _, option := range _options {
		switch option := option.(type) {
		case *tracerWorkersOption:
			_executor.traceWorkers = option.traceWorkers
		}
	}
	for i := 0; i < numberOfWorkers; i++ {
		workers[i].executor = _executor
	}
	return _executor
}

var upScaleInterval = 100 * time.Millisecond
var downScaleInterval = 5 * time.Second
var timeToBecomeUnused = 5 * time.Second

func NewDynamicExecutor(maxNumberOfWorkers, queueDepth int, _options ...options.WithOption) Executor {
	customLogger := options.ExtractCustomLogger(_options...)
	_executor := &dynamicExecutor{
		executor: &executor{
			workItems: make(chan workItem, queueDepth),
			worker:    make([]*worker, 0),
			log:       customLogger,
		},
		maxNumberOfWorkers: maxNumberOfWorkers,
	}
	for _, option := range _options {
		switch option := option.(type) {
		case *tracerWorkersOption:
			_executor.traceWorkers = option.traceWorkers
		}
	}
	// We spawn one initial worker
	_executor.worker = append(_executor.worker, &worker{
		id:           0,
		interrupter:  make(chan struct{}, 1),
		executor:     _executor,
		lastReceived: time.Now(),
		log:          customLogger,
	})
	return _executor
}

func WithExecutorOptionTracerWorkers(traceWorkers bool) options.WithOption {
	return &tracerWorkersOption{traceWorkers: traceWorkers}
}

type tracerWorkersOption struct {
	options.Option
	traceWorkers bool
}

func (e *executor) Submit(ctx context.Context, workItemId int32, runnable Runnable) CompletionFuture {
	if runnable == nil {
		value := atomic.Value{}
		value.Store(errors.New("runnable must not be nil"))
		return &future{err: value}
	}
	e.log.Trace().Int32("workItemId", workItemId).Msg("Submitting runnable")
	completionFuture := &future{}
	if e.shutdown {
		completionFuture.Cancel(false, errors.New("executor in shutdown"))
		return completionFuture
	}
	select {
	case e.workItems <- workItem{
		workItemId:       workItemId,
		runnable:         runnable,
		completionFuture: completionFuture,
	}:
		e.log.Trace().Msg("Item added")
	case <-ctx.Done():
		completionFuture.Cancel(false, ctx.Err())
	}

	e.log.Trace().Int32("workItemId", workItemId).Msg("runnable queued")
	return completionFuture
}

func (e *executor) Start() {
	e.stateChange.Lock()
	defer e.stateChange.Unlock()
	if e.running || e.shutdown {
		e.log.Warn().Msg("Already started")
		return
	}
	e.running = true
	e.shutdown = false
	for i := 0; i < len(e.worker); i++ {
		worker := e.worker[i]
		worker.initialize()
		go worker.work()
	}
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

func (e *executor) Stop() {
	e.log.Trace().Msg("stopping now")
	e.stateChange.Lock()
	defer e.stateChange.Unlock()
	if !e.running || e.shutdown {
		e.log.Warn().Msg("already stopped")
		return
	}
	e.shutdown = true
	for i := 0; i < len(e.worker); i++ {
		worker := e.worker[i]
		worker.shutdown.Store(true)
		worker.interrupted.Store(true)
		close(worker.interrupter)
	}
	e.running = false
	e.shutdown = false
	e.log.Debug().Msgf("waiting for %d workers to stop", len(e.worker))
	e.workerWaitGroup.Wait()
	e.log.Trace().Msg("stopped")
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

func (e *executor) Close() error {
	e.Stop()
	return nil
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

func (f *future) String() string {
	return fmt.Sprintf("future: cancelRequested(%t), interruptRequested(%t), completed(%t), errored(%t), err(%v)",
		f.cancelRequested.Load(),
		f.interruptRequested.Load(),
		f.completed.Load(),
		f.errored.Load(),
		f.err.Load(),
	)
}
