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
	sharedExecutorInstance = *utils.NewFixedSizeExecutor(runtime.NumCPU())
	sharedExecutorInstance.Start()
}

type RequestTransaction struct {
	parent        *RequestTransactionManager
	transactionId int32

	/** The initial operation to perform to kick off the request */
	operation        utils.Runnable
	completionFuture *utils.CompletionFuture

	transactionLog zerolog.Logger
}

func (t *RequestTransaction) String() string {
	return fmt.Sprintf("Transaction{tid:%d}", t.transactionId)
}

// RequestTransactionManager handles transactions
type RequestTransactionManager struct {
	runningRequests []*RequestTransaction
	// How many Transactions are allowed to run at the same time?
	numberOfConcurrentRequests int
	// Assigns each request a Unique Transaction Id, especially important for failure handling
	transactionId    int32
	transactionMutex sync.RWMutex
	// Important, this is a FIFO Queue for Fairness!
	workLog      list.List
	workLogMutex sync.RWMutex
	executor     *utils.Executor
}

// NewRequestTransactionManager creates a new RequestTransactionManager
func NewRequestTransactionManager(numberOfConcurrentRequests int, requestTransactionManagerOptions ...RequestTransactionManagerOption) *RequestTransactionManager {
	requestTransactionManager := &RequestTransactionManager{
		numberOfConcurrentRequests: numberOfConcurrentRequests,
		transactionId:              0,
		workLog:                    *list.New(),
		executor:                   &sharedExecutorInstance,
	}
	for _, requestTransactionManagerOption := range requestTransactionManagerOptions {
		requestTransactionManagerOption(requestTransactionManager)
	}
	return requestTransactionManager
}

type RequestTransactionManagerOption func(requestTransactionManager *RequestTransactionManager)

// WithCustomExecutor sets a custom Executor for the RequestTransactionManager
func WithCustomExecutor(executor *utils.Executor) RequestTransactionManagerOption {
	return func(requestTransactionManager *RequestTransactionManager) {
		requestTransactionManager.executor = executor
	}
}

// SetNumberOfConcurrentRequests sets the number of concurrent requests that will be sent out to a device
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
	r.workLogMutex.Lock()
	r.workLog.PushFront(handle)
	r.workLogMutex.Unlock()
	// Try to Process the Worklog
	r.processWorklog()
}

func (r *RequestTransactionManager) processWorklog() {
	r.workLogMutex.RLock()
	defer r.workLogMutex.RUnlock()
	log.Debug().Msgf("Processing work log with size of %d (%d concurrent requests allowed)", r.workLog.Len(), r.numberOfConcurrentRequests)
	for len(r.runningRequests) < r.numberOfConcurrentRequests && r.workLog.Len() > 0 {
		front := r.workLog.Front()
		if front == nil {
			return
		}
		next := front.Value.(*RequestTransaction)
		log.Debug().Msgf("Handling next %v. (Adding to running requests (length: %d))", next, len(r.runningRequests))
		r.runningRequests = append(r.runningRequests, next)
		// TODO: use sharedInstance if none is present
		completionFuture := sharedExecutorInstance.Submit(next.transactionId, next.operation)
		next.completionFuture = completionFuture
		r.workLog.Remove(front)
	}
}

// StartTransaction starts a RequestTransaction
func (r *RequestTransactionManager) StartTransaction() *RequestTransaction {
	r.transactionMutex.Lock()
	defer r.transactionMutex.Unlock()
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
	transaction.completionFuture.Cancel(true, err)
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
	// Process the workLog, a slot should be free now
	transaction.transactionLog.Debug().Msg("Processing the workLog")
	r.processWorklog()
	return nil
}

// FailRequest signals that this transaction has failed
func (t *RequestTransaction) FailRequest(err error) error {
	t.transactionLog.Trace().Msg("Fail the request")
	return t.parent.failRequest(t, err)
}

// EndRequest signals that this transaction is done
func (t *RequestTransaction) EndRequest() error {
	t.transactionLog.Trace().Msg("Ending the request")
	// Remove it from Running Requests
	return t.parent.endRequest(t)
}

// Submit submits a Runnable to the RequestTransactionManager
func (t *RequestTransaction) Submit(operation utils.Runnable) {
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

// AwaitCompletion wait for this RequestTransaction to finish. Returns an error if it finished unsuccessful
func (t *RequestTransaction) AwaitCompletion() error {
	for t.completionFuture == nil {
		time.Sleep(time.Millisecond * 10)
	}
	if err := t.completionFuture.AwaitCompletion(); err != nil {
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
