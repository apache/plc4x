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

package spi

import (
	"container/list"
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var sharedExecutorInstance utils.Executor // shared instance

func init() {
	sharedExecutorInstance = utils.NewFixedSizeExecutor(runtime.NumCPU(), 100, utils.WithExecutorOptionTracerWorkers(config.TraceTransactionManagerWorkers))
	sharedExecutorInstance.Start()
}

// RequestTransaction represents a transaction
type RequestTransaction interface {
	fmt.Stringer
	// FailRequest signals that this transaction has failed
	FailRequest(err error) error
	// EndRequest signals that this transaction is done
	EndRequest() error
	// Submit submits a Runnable to the RequestTransactionManager
	Submit(operation utils.Runnable)
	// AwaitCompletion wait for this RequestTransaction to finish. Returns an error if it finished unsuccessful
	AwaitCompletion(ctx context.Context) error
}

// RequestTransactionManager handles transactions
type RequestTransactionManager interface {
	// SetNumberOfConcurrentRequests sets the number of concurrent requests that will be sent out to a device
	SetNumberOfConcurrentRequests(numberOfConcurrentRequests int)
	// StartTransaction starts a RequestTransaction
	StartTransaction() RequestTransaction
}

// NewRequestTransactionManager creates a new RequestTransactionManager
func NewRequestTransactionManager(numberOfConcurrentRequests int, requestTransactionManagerOptions ...RequestTransactionManagerOption) RequestTransactionManager {
	_requestTransactionManager := &requestTransactionManager{
		numberOfConcurrentRequests: numberOfConcurrentRequests,
		transactionId:              0,
		workLog:                    *list.New(),
		executor:                   sharedExecutorInstance,
	}
	for _, requestTransactionManagerOption := range requestTransactionManagerOptions {
		requestTransactionManagerOption(_requestTransactionManager)
	}
	return _requestTransactionManager
}

type RequestTransactionManagerOption func(requestTransactionManager *requestTransactionManager)

// WithCustomExecutor sets a custom Executor for the RequestTransactionManager
func WithCustomExecutor(executor utils.Executor) RequestTransactionManagerOption {
	return func(requestTransactionManager *requestTransactionManager) {
		requestTransactionManager.executor = executor
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type requestTransaction struct {
	parent        *requestTransactionManager
	transactionId int32

	/** The initial operation to perform to kick off the request */
	operation        utils.Runnable
	completionFuture utils.CompletionFuture

	transactionLog zerolog.Logger
}

type requestTransactionManager struct {
	runningRequests []*requestTransaction
	// How many Transactions are allowed to run at the same time?
	numberOfConcurrentRequests int
	// Assigns each request a Unique Transaction Id, especially important for failure handling
	transactionId    int32
	transactionMutex sync.RWMutex
	// Important, this is a FIFO Queue for Fairness!
	workLog      list.List
	workLogMutex sync.RWMutex
	executor     utils.Executor
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (r *requestTransactionManager) SetNumberOfConcurrentRequests(numberOfConcurrentRequests int) {
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

func (r *requestTransactionManager) submitHandle(handle *requestTransaction) {
	if handle.operation == nil {
		panic("invalid handle")
	}
	// Add this Request with this handle i the Worklog
	// Put Transaction into Worklog
	r.workLogMutex.Lock()
	r.workLog.PushFront(handle)
	r.workLogMutex.Unlock()
	// Try to Process the Worklog
	r.processWorklog()
}

func (r *requestTransactionManager) processWorklog() {
	r.workLogMutex.RLock()
	defer r.workLogMutex.RUnlock()
	log.Debug().Msgf("Processing work log with size of %d (%d concurrent requests allowed)", r.workLog.Len(), r.numberOfConcurrentRequests)
	for len(r.runningRequests) < r.numberOfConcurrentRequests && r.workLog.Len() > 0 {
		front := r.workLog.Front()
		if front == nil {
			return
		}
		next := front.Value.(*requestTransaction)
		log.Debug().Msgf("Handling next %v. (Adding to running requests (length: %d))", next, len(r.runningRequests))
		r.runningRequests = append(r.runningRequests, next)
		completionFuture := r.executor.Submit(context.Background(), next.transactionId, next.operation)
		next.completionFuture = completionFuture
		r.workLog.Remove(front)
	}
}

func (r *requestTransactionManager) StartTransaction() RequestTransaction {
	r.transactionMutex.Lock()
	defer r.transactionMutex.Unlock()
	currentTransactionId := r.transactionId
	r.transactionId += 1
	transactionLogger := log.With().Int32("transactionId", currentTransactionId).Logger()
	if !config.TraceTransactionManagerTransactions {
		transactionLogger = zerolog.Nop()
	}
	return &requestTransaction{
		r,
		currentTransactionId,
		nil,
		nil,
		transactionLogger,
	}
}

func (r *requestTransactionManager) getNumberOfActiveRequests() int {
	return len(r.runningRequests)
}

func (r *requestTransactionManager) failRequest(transaction *requestTransaction, err error) error {
	// Try to fail it!
	transaction.completionFuture.Cancel(true, err)
	// End it
	return r.endRequest(transaction)
}

func (r *requestTransactionManager) endRequest(transaction *requestTransaction) error {
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
	// Process the workLog, a slot should be free now
	transaction.transactionLog.Debug().Msg("Processing the workLog")
	r.processWorklog()
	return nil
}

func (t *requestTransaction) FailRequest(err error) error {
	t.transactionLog.Trace().Msg("Fail the request")
	return t.parent.failRequest(t, err)
}

func (t *requestTransaction) EndRequest() error {
	t.transactionLog.Trace().Msg("Ending the request")
	// Remove it from Running Requests
	return t.parent.endRequest(t)
}

func (t *requestTransaction) Submit(operation utils.Runnable) {
	if t.operation != nil {
		panic("Operation already set")
	}
	t.transactionLog.Trace().Msgf("Submission of transaction %d", t.transactionId)
	t.operation = func() {
		t.transactionLog.Trace().Msgf("Start execution of transaction %d", t.transactionId)
		operation()
		t.transactionLog.Trace().Msgf("Completed execution of transaction %d", t.transactionId)
	}
	t.parent.submitHandle(t)
}

func (t *requestTransaction) AwaitCompletion(ctx context.Context) error {
	for t.completionFuture == nil {
		time.Sleep(time.Millisecond * 10)
	}
	if err := t.completionFuture.AwaitCompletion(ctx); err != nil {
		return err
	}
	stillActive := true
	for stillActive {
		stillActive = false
		for _, runningRequest := range t.parent.runningRequests {
			if runningRequest.transactionId == t.transactionId {
				stillActive = true
				break
			}
		}
	}
	return nil
}

func (t *requestTransaction) String() string {
	return fmt.Sprintf("Transaction{tid:%d}", t.transactionId)
}
