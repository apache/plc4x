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
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"reflect"
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
			log.Fatal().Msgf("Recovering from panic()=%v", recovered)
		}
		if !w.shutdown {
			// TODO: if we are not in shutdown we continue
			w.work()
		}
	}()
	for !w.shutdown {
		log.Debug().Int("Worker id", w.id).Msg("Working")
		select {
		case workItem := <-w.executor.queue:
			log.Debug().Int("Worker id", w.id).Msgf("Got work item %v", workItem)
			if workItem.completionFuture.cancelRequested || (w.shutdown && w.interrupted) {
				log.Debug().Int("Worker id", w.id).Msg("We need to stop")
				// TODO: do we need to complete with a error?
			} else {
				log.Debug().Int("Worker id", w.id).Msgf("Running work item %v", workItem)
				workItem.runnable()
				workItem.completionFuture.complete()
				log.Debug().Int("Worker id", w.id).Msgf("work item %v completed", workItem)
			}
		default:
			log.Debug().Int("Worker id", w.id).Msgf("Idling")
			time.Sleep(time.Millisecond * 10)
		}
	}
}

type WorkItem struct {
	runnable         Runnable
	completionFuture *CompletionFuture
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

func (e *Executor) submit(runnable Runnable) *CompletionFuture {
	completionFuture := &CompletionFuture{}
	e.queue <- WorkItem{
		runnable:         runnable,
		completionFuture: completionFuture,
	}
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
}

func (f CompletionFuture) cancel(interrupt bool) {
	f.cancelRequested = true
	f.interruptRequested = interrupt
	f.errored = true
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
	// If we reduced the number of concurrent requests and more requests are in-flight
	// than should be, at least log a warning.
	if numberOfConcurrentRequests < len(r.runningRequests) {
		log.Warn().Msg("The number of concurrent requests was reduced and currently more requests are in flight.")
	}

	r.numberOfConcurrentRequests = numberOfConcurrentRequests

	// As we might have increased the number, try to send some more requests.
	r.processWorklog()
}

func (r *RequestTransactionManager) submit(context func(RequestTransaction)) {
	transaction := r.StartRequest()
	context(*transaction)
	// r.submitHandle(transaction);
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
	for len(r.runningRequests) < r.getNumberOfConcurrentRequests() && r.worklog.Len() > 0 {
		next := r.worklog.Front().Value.(*RequestTransaction)
		r.runningRequests = append(r.runningRequests, next)
		completionFuture := executor.submit(next.operation)
		next.completionFuture = completionFuture
	}
}

func (r *RequestTransactionManager) StartRequest() *RequestTransaction {
	r.transationMutex.Lock()
	defer r.transationMutex.Unlock()
	currentTransactionId := r.transactionId
	r.transactionId += 1
	return &RequestTransaction{r, currentTransactionId, nil, nil}
}

func (r *RequestTransactionManager) getNumberOfActiveRequests() int {
	return len(r.runningRequests)
}

func (r *RequestTransactionManager) failRequest(transaction *RequestTransaction) error {
	// Try to fail it!
	transaction.completionFuture.cancel(true)
	// End it
	return r.endRequest(transaction)
}

func (r *RequestTransactionManager) endRequest(transaction *RequestTransaction) error {
	found := false
	index := -1
	for i, runningRequest := range r.runningRequests {
		// TODO: check this implementation
		if (&runningRequest) == (&transaction) {
			found = true
			index = i
			break
		}
	}
	if !found {
		return errors.New("Unknown Transaction or Transaction already finished!")
	}
	r.runningRequests = append(r.runningRequests[:index], r.runningRequests[index+1:]...)
	// Process the worklog, a slot should be free now
	r.processWorklog()
	return nil
}

func (t *RequestTransaction) start() {
}

func (t *RequestTransaction) failRequest(err error) error {
	return t.parent.failRequest(t)
}

func (t *RequestTransaction) EndRequest() error {
	// Remove it from Running Requests
	return t.parent.endRequest(t)
}

func (t *RequestTransaction) setOperation(operation Runnable) {
	t.operation = operation
}

func (t *RequestTransaction) getCompletionFuture() *CompletionFuture {
	return t.completionFuture
}

func (t *RequestTransaction) setCompletionFuture(completionFuture *CompletionFuture) {
	t.completionFuture = completionFuture
}

func (t *RequestTransaction) Submit(operation Runnable) {
	log.Trace().Msgf("Submission of transaction %d", t.transactionId)
	t.setOperation(NewTransactionOperation(t.transactionId, operation))
	t.parent.submitHandle(t)
}

func (t *RequestTransaction) equals(o *RequestTransaction) bool {
	if t == o {
		return true
	}
	if o == nil || reflect.TypeOf(t).Kind() != reflect.TypeOf(o).Kind() {
		return false
	}
	that := o
	return t.transactionId == that.transactionId
}

func (t *RequestTransaction) hashCode() int32 {
	return t.transactionId
}

func NewTransactionOperation(transactionId int32, delegate Runnable) Runnable {
	return func() {
		log.Trace().Int32("transactionId", transactionId).Msgf("Start execution of transaction %d", transactionId)
		delegate()
		log.Trace().Int32("transactionId", transactionId).Msgf("Completed execution of transaction %d", transactionId)
	}
}
