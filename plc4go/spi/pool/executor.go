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
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

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

func (e *executor) Close() error {
	e.Stop()
	return nil
}

func (e *executor) IsRunning() bool {
	return e.running && !e.shutdown
}

func (e *executor) String() string {
	return fmt.Sprintf("executor{\n"+
		"\trunning: %t,\n"+
		"\tshutdown: %t,\n"+
		"\tworker: %s,\n"+
		"\tqueueDepth: %d,\n"+
		"\tworkItems: %d elements,\n"+
		"\ttraceWorkers: %t,\n"+
		"\n}",
		e.running,
		e.shutdown,
		e.worker,
		e.queueDepth,
		len(e.workItems),
		e.traceWorkers,
	)
}
