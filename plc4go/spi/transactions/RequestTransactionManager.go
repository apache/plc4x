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

package transactions

import (
	"container/list"
	"context"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/pool"
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi/utils"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var sharedExecutorInstance pool.Executor // shared instance

func init() {
	sharedExecutorInstance = pool.NewFixedSizeExecutor(runtime.NumCPU(), 100, pool.WithExecutorOptionTracerWorkers(config.TraceTransactionManagerWorkers))
	sharedExecutorInstance.Start()
}

type RequestTransactionRunnable func(RequestTransaction)

// RequestTransactionManager handles transactions
type RequestTransactionManager interface {
	io.Closer
	// CloseGraceful gives some time opposed to io.Closer
	CloseGraceful(timeout time.Duration) error
	// SetNumberOfConcurrentRequests sets the number of concurrent requests that will be sent out to a device
	SetNumberOfConcurrentRequests(numberOfConcurrentRequests int)
	// StartTransaction starts a RequestTransaction
	StartTransaction() RequestTransaction
}

// NewRequestTransactionManager creates a new RequestTransactionManager
func NewRequestTransactionManager(numberOfConcurrentRequests int, _options ...options.WithOption) RequestTransactionManager {
	_requestTransactionManager := &requestTransactionManager{
		numberOfConcurrentRequests: numberOfConcurrentRequests,
		currentTransactionId:       0,
		workLog:                    *list.New(),
		executor:                   sharedExecutorInstance,

		traceTransactionManagerTransactions: config.TraceTransactionManagerTransactions,

		log: options.ExtractCustomLogger(_options...),
	}
	for _, option := range _options {
		switch option := option.(type) {
		case *withCustomExecutor:
			_requestTransactionManager.executor = option.executor
		}
	}
	return _requestTransactionManager
}

// WithCustomExecutor sets a custom Executor for the RequestTransactionManager
func WithCustomExecutor(executor pool.Executor) options.WithOption {
	return &withCustomExecutor{executor: executor}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type withCustomExecutor struct {
	options.Option
	executor pool.Executor
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=requestTransactionManager
type requestTransactionManager struct {
	runningRequests []*requestTransaction
	// How many transactions are allowed to run at the same time?
	numberOfConcurrentRequests int
	// Assigns each request a Unique Transaction Id, especially important for failure handling
	currentTransactionId int32
	transactionMutex     sync.RWMutex
	// Important, this is a FIFO Queue for Fairness!
	workLog      list.List `ignore:"true"` // TODO: no support for list yet
	workLogMutex sync.RWMutex
	executor     pool.Executor

	// Indicates it this rtm is in shutdown
	shutdown bool

	// flag set to true if it should trace transactions
	traceTransactionManagerTransactions bool

	log zerolog.Logger `ignore:"true"`
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (r *requestTransactionManager) SetNumberOfConcurrentRequests(numberOfConcurrentRequests int) {
	r.log.Info().Msgf("Setting new number of concurrent requests %d", numberOfConcurrentRequests)
	// If we reduced the number of concurrent requests and more requests are in-flight
	// than should be, at least log a warning.
	if numberOfConcurrentRequests < len(r.runningRequests) {
		r.log.Warn().Msg("The number of concurrent requests was reduced and currently more requests are in flight.")
	}

	r.numberOfConcurrentRequests = numberOfConcurrentRequests

	// As we might have increased the number, try to send some more requests.
	r.processWorklog()
}

func (r *requestTransactionManager) submitTransaction(transaction *requestTransaction) {
	// Add this Request with the transaction i the work log
	// Put Transaction into work log
	r.workLogMutex.Lock()
	r.workLog.PushFront(transaction)
	r.workLogMutex.Unlock()
	// Try to Process the work log
	r.processWorklog()
}

func (r *requestTransactionManager) processWorklog() {
	r.workLogMutex.RLock()
	defer r.workLogMutex.RUnlock()
	r.log.Debug().Msgf("Processing work log with size of %d (%d concurrent requests allowed)", r.workLog.Len(), r.numberOfConcurrentRequests)
	for len(r.runningRequests) < r.numberOfConcurrentRequests && r.workLog.Len() > 0 {
		front := r.workLog.Front()
		next := front.Value.(*requestTransaction)
		r.log.Debug().Msgf("Handling next %v. (Adding to running requests (length: %d))", next, len(r.runningRequests))
		r.runningRequests = append(r.runningRequests, next)
		completionFuture := r.executor.Submit(context.Background(), next.transactionId, next.operation)
		next.completionFuture = completionFuture
		r.workLog.Remove(front)
	}
}

func (r *requestTransactionManager) StartTransaction() RequestTransaction {
	r.transactionMutex.Lock()
	defer r.transactionMutex.Unlock()
	currentTransactionId := r.currentTransactionId
	r.currentTransactionId += 1
	transactionLogger := r.log.With().Int32("currentTransactionId", currentTransactionId).Logger()
	if !r.traceTransactionManagerTransactions {
		transactionLogger = zerolog.Nop()
	}
	transaction := &requestTransaction{
		parent:         r,
		transactionId:  currentTransactionId,
		transactionLog: transactionLogger,
	}
	if r.shutdown {
		transaction.completed = true
		transaction.completionFuture = &completedFuture{errors.New("request transaction manager in shutdown")}
	}
	return transaction
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

func (r *requestTransactionManager) Close() error {
	return r.CloseGraceful(0)
}

func (r *requestTransactionManager) CloseGraceful(timeout time.Duration) error {
	r.shutdown = true
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		defer utils.CleanupTimer(timer)
		signal := make(chan struct{})
		go func() {
			for {
				if len(r.runningRequests) == 0 {
					break
				}
				time.Sleep(10 * time.Millisecond)
			}
			close(signal)
		}()
		select {
		case <-timer.C:
			r.log.Warn().Msgf("timout after %d", timeout)
		case <-signal:
		}
	}
	r.transactionMutex.Lock()
	defer r.transactionMutex.Unlock()
	r.workLogMutex.RLock()
	defer r.workLogMutex.RUnlock()
	r.runningRequests = nil
	return r.executor.Close()
}
