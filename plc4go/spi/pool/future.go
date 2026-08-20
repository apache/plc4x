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
	"sync/atomic"
	"time"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

//go:generate go tool plc4xGenerator -type=future
type future struct {
	cancelRequested    atomic.Bool
	interruptRequested atomic.Bool
	completed          atomic.Bool
	errored            atomic.Bool
	err                atomic.Value
}

func (f *future) Cancel(interrupt bool, err error) {
	f.cancel(interrupt, err, func() {})
}

// cancel publishes a cancellation store by store, calling observe after each of them.
//
// The store order matters: settled leaves the poll loop of AwaitCompletion as soon as it observes
// errored or cancelRequested, and result reads err only afterwards. So err has to be published
// before the flags which release the waiter, otherwise the waiter can wake up on a cancel that
// carried an error and still find no error to report, returning the Canceled sentinel instead of the
// error the caller passed in. sync/atomic operations are sequentially consistent, so a waiter which
// sees a flag set here is guaranteed to also see every store preceding it below.
//
// observe is what makes that order testable: getting a waiter to observe the window between two
// adjacent atomic stores is a matter of luck (in practice it only happens under the race detector,
// which instruments every one of them), while stepping through the stores hits every intermediate
// state by construction. Cancel passes a no-op, so production code pays nothing for it - a func
// literal without captures is a static value, not an allocation.
func (f *future) cancel(interrupt bool, err error, observe func()) {
	if err != nil {
		f.err.Store(err)
		observe()
		f.errored.Store(true)
		observe()
	}
	f.interruptRequested.Store(interrupt)
	observe()
	f.cancelRequested.Store(true)
	observe()
}

func (f *future) complete() {
	f.completed.Store(true)
}

// Canceled is returned on CompletionFuture.AwaitCompletion when a CompletionFuture was canceled
var Canceled = errors.New("Canceled")

// settled reports whether the future reached a terminal state, so AwaitCompletion can stop waiting.
func (f *future) settled() bool {
	return f.completed.Load() || f.errored.Load() || f.cancelRequested.Load()
}

// result maps the terminal state of a settled future onto the error AwaitCompletion reports.
func (f *future) result() error {
	if err, ok := f.err.Load().(error); ok {
		return err
	}
	if f.cancelRequested.Load() {
		return Canceled
	}
	return nil
}

func (f *future) AwaitCompletion(ctx context.Context) error {
	for !f.settled() && ctx.Err() == nil {
		time.Sleep(10 * time.Millisecond)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return f.result()
}
