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
	"sync"
	"sync/atomic"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

//go:generate go tool plc4xGenerator -type=future
type future struct {
	cancelRequested    atomic.Bool
	interruptRequested atomic.Bool
	completed          atomic.Bool
	errored            atomic.Bool
	err                atomic.Value

	// settledChanOnce creates settledChan, settleOnce closes it. Both are lazy on purpose: futures are
	// built as composite literals (&future{}), so there is no constructor which could set the channel
	// up, and a receive on a nil channel blocks forever.
	//
	// The channel is not serialized: len of a closed signalling channel is 0 just like len of an open
	// one, so the box would claim something it cannot know, while the flags above carry the real state
	// (errored and interruptRequested are only ever read for that debug output and by the tests).
	settledChanOnce sync.Once
	settleOnce      sync.Once
	settledChan     chan struct{} `ignore:"true"`
}

// settledSignal returns the channel which is closed once the future settles, creating it on first
// use. sync.Once orders the creation against every reader, so racing waiters all get the same
// channel and never a nil one.
func (f *future) settledSignal() chan struct{} {
	f.settledChanOnce.Do(func() {
		f.settledChan = make(chan struct{})
	})
	return f.settledChan
}

// settle publishes the terminal state of the future to everyone waiting on it, by closing the
// channel AwaitCompletion parks on.
//
// Closing is the whole point: it turns the handful of stores a settling caller made into a single
// happens-before edge, so a waiter which observes the close is guaranteed to observe all of them -
// the state a future reports can no longer depend on the order those stores were issued in. It also
// wakes the waiters immediately instead of on the next tick of a poll loop.
//
// The close has to happen exactly once: Cancel is part of the public CompletionFuture, so a caller
// outside the pool can cancel a work item which a worker then completes anyway, and a second close
// panics. sync.Once makes whichever of the two arrives second a no-op, so the waiters the first
// settle released stay released - and because result reads the error before any flag, a completion
// following a cancellation still reports the cancellation error.
func (f *future) settle() {
	f.settleOnce.Do(func() {
		close(f.settledSignal())
	})
}

func (f *future) Cancel(interrupt bool, err error) {
	if err != nil {
		f.err.Store(err)
		f.errored.Store(true)
	}
	f.interruptRequested.Store(interrupt)
	f.cancelRequested.Store(true)
	f.settle()
}

func (f *future) complete() {
	f.completed.Store(true)
	f.settle()
}

// Canceled is returned on CompletionFuture.AwaitCompletion when a CompletionFuture was canceled
var Canceled = errors.New("Canceled")

// settled reports whether the future reached a terminal state.
func (f *future) settled() bool {
	select {
	case <-f.settledSignal():
		return true
	default:
		return false
	}
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
	select {
	case <-f.settledSignal():
	case <-ctx.Done():
	}
	// A context which is already done outranks a settled future, as it did while this polled: the
	// caller stopped being interested in the answer.
	if err := ctx.Err(); err != nil {
		return err
	}
	return f.result()
}
