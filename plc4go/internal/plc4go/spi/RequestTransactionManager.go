//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package spi

import (
	"container/list"
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

/** Executor that performs all operations */
var executor Executor // shared instance

func init() {
	executor = *NewFixedSizeExecutor(4)
	executor.start()
}

type Runnable func()

type Worker struct {
	id          int
	shutdown    bool
	runnable    Runnable
	interrupted bool
	executor    *Executor
}

func (w Worker) work() {
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Error().Msgf("Recovering from panic()=%v", recovered)
		}
		if !w.shutdown {
			// TODO: if we are not in shutdown we continue
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
	transactionId    int32
	runnable         Runnable
	completionFuture *CompletionFuture
}

func (w WorkItem) String() string {
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
	executor = Executor{
		queue:  make(chan WorkItem, 100),
		worker: workers,
	}
	for i := 0; i < numberOfWorkers; i++ {
		worker := workers[i]
		worker.executor = &executor
	}
	return &executor
}

func (e *Executor) submit(transactionId int32, runnable Runnable) *CompletionFuture {
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

func (e *Executor) start() {
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

func (e *Executor) stop() {
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

func (f CompletionFuture) cancel(interrupt bool, err error) {
	f.cancelRequested = true
	f.interruptRequested = interrupt
	f.errored = true
	f.err = err
}

func (f CompletionFuture) complete() {
	f.completed = true
}

func (f CompletionFuture) AwaitCompletion() {
	for !f.completed || !f.errored {
		time.Sleep(time.Millisecond * 10)
	}
}

type RequestTransaction struct {
	parent        *RequestTransactionManager
	transactionId int32

	/** The initial operation to perform to kick off the request */
	operation        Runnable
	completionFuture *CompletionFuture

	transactionLog zerolog.Logger
}

func (t RequestTransaction) String() string {
	return fmt.Sprintf("Transaction{tid:%d}", t.transactionId)
}

type RequestTransactionManager struct {
	runningRequests []*RequestTransaction
	// How many Transactions are allowed to run at the same time?
	numberOfConcurrentRequests int
	// Assigns each request a Unique Transaction Id, especially important for failure handling
	transactionId   int32
	transationMutex sync.RWMutex
	// Important, this is a FIFO Queue for Fairness!
	worklog      list.List
	worklogMutex sync.RWMutex
}

func NewRequestTransactionManager(numberOfConcurrentRequests int) RequestTransactionManager {
	return RequestTransactionManager{
		numberOfConcurrentRequests: numberOfConcurrentRequests,
		transactionId:              0,
		worklog:                    *list.New(),
	}
}

func (r *RequestTransactionManager) getNumberOfConcurrentRequests() int {
	return r.numberOfConcurrentRequests
}

func (r *RequestTransactionManager) SetNumberOfConcurrentRequests(numberOfConcurrentRequests int) {
	log.Info().Msgf("Setting new number of concurrent requests %d", numberOfConcurrentRequests)
	// If we reduced the number of concurrent requests and more requests are in-flight
	// than should be, at least log a warning.
	if numberOfConcurrentRequests < len(r.runningRequests) {
		log.Warn().Msg("The number of concurrent requests was reduced and currently more requests are in flight.")
	}

	r.numberOfConcurrentRequests = numberOfConcurrentRequests

	// As we might have increased the number, try to send some more requests.
	r.processWorklog()
}

func (r *RequestTransactionManager) submitHandle(handle *RequestTransaction) {
	if handle.operation == nil {
		panic("invalid handle")
	}
	// Add this Request with this handle i the Worklog
	// Put Transaction into Worklog
	r.worklogMutex.Lock()
	r.worklog.PushFront(handle)
	r.worklogMutex.Unlock()
	// Try to Process the Worklog
	r.processWorklog()
}

func (r *RequestTransactionManager) processWorklog() {
	r.worklogMutex.RLock()
	defer r.worklogMutex.RUnlock()
	log.Debug().Msgf("Processing work log with size of %d (%d concurrent requests allowed)", r.worklog.Len(), r.numberOfConcurrentRequests)
	for len(r.runningRequests) < r.numberOfConcurrentRequests && r.worklog.Len() > 0 {
		front := r.worklog.Front()
		if front == nil {
			return
		}
		next := front.Value.(*RequestTransaction)
		log.Debug().Msgf("Handling next %v. (Adding to running requests (length: %d))", next, len(r.runningRequests))
		r.runningRequests = append(r.runningRequests, next)
		completionFuture := executor.submit(next.transactionId, next.operation)
		next.completionFuture = completionFuture
		r.worklog.Remove(front)
	}
}

func (r *RequestTransactionManager) StartTransaction() *RequestTransaction {
	r.transationMutex.Lock()
	defer r.transationMutex.Unlock()
	currentTransactionId := r.transactionId
	r.transactionId += 1
	transactionLogger := log.With().Int32("transactionId", currentTransactionId).Logger()
	if !config.TraceTransactionManagerTransactions {
		transactionLogger = zerolog.Nop()
	}
	return &RequestTransaction{
		r,
		currentTransactionId,
		nil,
		nil,
		transactionLogger,
	}
}

func (r *RequestTransactionManager) getNumberOfActiveRequests() int {
	return len(r.runningRequests)
}

func (r *RequestTransactionManager) failRequest(transaction *RequestTransaction, err error) error {
	// Try to fail it!
	transaction.completionFuture.cancel(true, err)
	// End it
	return r.endRequest(transaction)
}

func (r *RequestTransactionManager) endRequest(transaction *RequestTransaction) error {
	transaction.transactionLog.Debug().Msg("Trying to find a existing transaction")
	found := false
	index := -1
	for i, runningRequest := range r.runningRequests {
		if runningRequest.transactionId == transaction.transactionId {
			transaction.transactionLog.Debug().Msg("Found a existing transaction")
			found = true
			index = i
			break
		}
	}
	if !found {
		return errors.New("Unknown Transaction or Transaction already finished!")
	}
	transaction.transactionLog.Debug().Msg("Removing the existing transaction transaction")
	r.runningRequests = append(r.runningRequests[:index], r.runningRequests[index+1:]...)
	// Process the worklog, a slot should be free now
	transaction.transactionLog.Debug().Msg("Processing the worklog")
	r.processWorklog()
	return nil
}

func (t *RequestTransaction) FailRequest(err error) error {
	t.transactionLog.Trace().Msg("Fail the request")
	return t.parent.failRequest(t, err)
}

func (t *RequestTransaction) EndRequest() error {
	t.transactionLog.Trace().Msg("Ending the request")
	// Remove it from Running Requests
	return t.parent.endRequest(t)
}

func (t *RequestTransaction) Submit(operation Runnable) {
	if t.operation != nil {
		panic("Operation already set")
	}
	t.transactionLog.Trace().Msgf("Submission of transaction %d", t.transactionId)
	t.operation = t.NewTransactionOperation(operation)
	t.parent.submitHandle(t)
}

func (t *RequestTransaction) NewTransactionOperation(delegate Runnable) Runnable {
	return func() {
		t.transactionLog.Trace().Msgf("Start execution of transaction %d", t.transactionId)
		delegate()
		t.transactionLog.Trace().Msgf("Completed execution of transaction %d", t.transactionId)
	}
}
