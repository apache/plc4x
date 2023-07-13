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
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/plc4x/plc4go/spi/pool"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RequestTransaction represents a transaction
type RequestTransaction interface {
	fmt.Stringer
	// FailRequest signals that this transaction has failed
	FailRequest(err error) error
	// EndRequest signals that this transaction is done
	EndRequest() error
	// Submit submits a RequestTransactionRunnable to the RequestTransactionManager
	Submit(operation RequestTransactionRunnable)
	// AwaitCompletion wait for this RequestTransaction to finish. Returns an error if it finished unsuccessful
	AwaitCompletion(ctx context.Context) error
	// IsCompleted indicates that the that this RequestTransaction is completed
	IsCompleted() bool
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=requestTransaction
type requestTransaction struct {
	parent        *requestTransactionManager `ignore:"true"`
	transactionId int32

	/** The initial operation to perform to kick off the request */
	operation        pool.Runnable `ignore:"true"` // TODO: maybe we can treat this as a function some day if we are able to check the definition in gen
	completionFuture atomic.Pointer[pool.CompletionFuture]

	stateChangeMutex sync.Mutex
	completed        bool

	log zerolog.Logger `ignore:"true"`
}

func newRequestTransaction(localLog zerolog.Logger, parent *requestTransactionManager, transactionId int32) *requestTransaction {
	return &requestTransaction{
		parent:        parent,
		transactionId: transactionId,
		log:           localLog.With().Int32("transactionId", transactionId).Logger(),
	}
}

func (t *requestTransaction) setCompletionFuture(completionFuture pool.CompletionFuture) {
	t.completionFuture.Store(&completionFuture)
}

func (t *requestTransaction) getCompletionFuture() pool.CompletionFuture {
	completionFutureLoaded := t.completionFuture.Load()
	if completionFutureLoaded == nil {
		return nil
	}
	return *completionFutureLoaded
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (t *requestTransaction) FailRequest(err error) error {
	t.stateChangeMutex.Lock()
	defer t.stateChangeMutex.Unlock()
	if t.completed {
		return errors.Wrap(err, "calling fail on a already completed transaction")
	}
	t.log.Trace().Msg("Fail the request")
	t.completed = true
	return t.parent.failRequest(t, err)
}

func (t *requestTransaction) EndRequest() error {
	t.stateChangeMutex.Lock()
	defer t.stateChangeMutex.Unlock()
	if t.completed {
		return errors.New("calling end on a already completed transaction")
	}
	t.log.Trace().Msg("Ending the request")
	t.completed = true
	// Remove it from Running Requests
	return t.parent.endRequest(t)
}

func (t *requestTransaction) Submit(operation RequestTransactionRunnable) {
	t.stateChangeMutex.Lock()
	defer t.stateChangeMutex.Unlock()
	if t.completed {
		t.log.Warn().Msg("calling submit on a already completed transaction")
		return
	}
	if t.operation != nil {
		t.log.Warn().Msg("Operation already set")
	}
	t.log.Trace().Int32("transactionId", t.transactionId).Msg("Submission")
	t.operation = func() {
		t.log.Trace().Int32("transactionId", t.transactionId).Msg("Start operation")
		operation(t)
		t.log.Trace().Int32("transactionId", t.transactionId).Msg("Completed operation")
	}
	t.parent.submitTransaction(t)
}

func (t *requestTransaction) AwaitCompletion(ctx context.Context) error {
	t.log.Trace().Int32("transactionId", t.transactionId).Msg("Awaiting completion")
	timeout, cancelFunc := context.WithTimeout(ctx, time.Minute*30) // This is intentionally set very high
	defer cancelFunc()
	for t.getCompletionFuture() == nil {
		time.Sleep(time.Millisecond * 10)
		if err := timeout.Err(); err != nil {
			log.Error().Msg("Timout after a long time. This means something is very of here")
			return errors.Wrap(err, "Error waiting for completion future to be set")
		}
	}
	if err := t.getCompletionFuture().AwaitCompletion(ctx); err != nil {
		t.log.Trace().Int32("transactionId", t.transactionId).Msg("Errored")
		return err
	}
	stillActive := true
	for stillActive {
		stillActive = false
		t.parent.runningRequestMutex.RLock()
		for _, runningRequest := range t.parent.runningRequests {
			if runningRequest.transactionId == t.transactionId {
				stillActive = true
				break
			}
		}
		t.parent.runningRequestMutex.RUnlock()
	}
	t.log.Trace().Int32("transactionId", t.transactionId).Msg("Completed")
	return nil
}

func (t *requestTransaction) IsCompleted() bool {
	return t.completed
}
