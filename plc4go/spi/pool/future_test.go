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
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

func Test_future_AwaitCompletion(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		completer func(*sync.WaitGroup, *future)
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "completes with error",
			args: args{ctx: t.Context()},
			completer: func(wg *sync.WaitGroup, f *future) {
				defer wg.Done()
				f.Cancel(false, errors.New("Uh oh"))
			},
			wantErr: assert.Error,
		},
		{
			name: "completes regular",
			args: args{ctx: t.Context()},
			completer: func(wg *sync.WaitGroup, f *future) {
				defer wg.Done()
				time.Sleep(30 * time.Millisecond)
				f.complete()
			},
			wantErr: assert.NoError,
		},
		{
			name: "completes not int time",
			args: args{ctx: func() context.Context {
				deadline, cancel := context.WithDeadline(t.Context(), time.Now().Add(30*time.Millisecond))
				t.Cleanup(cancel)
				return deadline
			}()},
			completer: func(wg *sync.WaitGroup, f *future) {
				defer wg.Done()
				time.Sleep(300 * time.Millisecond)
			},
			wantErr: assert.Error,
		},
		{
			name: "completes canceled without error",
			args: args{ctx: t.Context()},
			completer: func(wg *sync.WaitGroup, f *future) {
				defer wg.Done()
				time.Sleep(300 * time.Millisecond)
				f.Cancel(true, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				assert.Same(t, Canceled, err)
				return true
			},
		},
		{
			name: "completes canceled with particular error",
			args: args{ctx: t.Context()},
			completer: func(wg *sync.WaitGroup, f *future) {
				defer wg.Done()
				time.Sleep(300 * time.Millisecond)
				f.Cancel(true, errors.New("Uh oh"))
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				assert.Equal(t, "Uh oh", err.Error())
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &future{}
			wg := sync.WaitGroup{}
			wg.Add(1)
			go tt.completer(&wg, f)
			tt.wantErr(t, f.AwaitCompletion(tt.args.ctx), fmt.Sprintf("AwaitCompletion(%v)", tt.args.ctx))
			wg.Wait()
		})
	}
}

func Test_future_Cancel(t *testing.T) {
	type args struct {
		interrupt bool
		err       error
	}
	tests := []struct {
		name     string
		args     args
		verifier func(*testing.T, *future)
	}{
		{
			name: "cancel cancels",
			verifier: func(t *testing.T, f *future) {
				assert.True(t, f.cancelRequested.Load())
			},
		},
		{
			name: "cancel with interrupt",
			args: args{
				interrupt: true,
				err:       nil,
			},
			verifier: func(t *testing.T, f *future) {
				assert.True(t, f.cancelRequested.Load())
				assert.False(t, f.errored.Load())
				assert.Nil(t, f.err.Load())
			},
		},
		{
			name: "cancel with err",
			args: args{
				interrupt: true,
				err:       errors.New("Uh Oh"),
			},
			verifier: func(t *testing.T, f *future) {
				assert.True(t, f.cancelRequested.Load())
				assert.True(t, f.errored.Load())
				assert.NotNil(t, f.err.Load())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &future{}
			f.Cancel(tt.args.interrupt, tt.args.err)
			tt.verifier(t, f)
		})
	}
}

func Test_future_String(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "string it",
			want: `
╔═future══════════════════════════════════════════════════════╗
║╔═cancelRequested╗╔═interruptRequested╗╔═completed╗╔═errored╗║
║║    b0 false    ║║     b0 false      ║║ b0 false ║║b0 false║║
║╚════════════════╝╚═══════════════════╝╚══════════╝╚════════╝║
╚═════════════════════════════════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &future{}
			assert.Equalf(t, tt.want, f.String(), "String()")
		})
	}
}

func Test_future_complete(t *testing.T) {
	tests := []struct {
		name     string
		verifier func(*testing.T, *future)
	}{
		{
			name: "complete completes",
			verifier: func(t *testing.T, f *future) {
				assert.True(t, f.completed.Load())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &future{}
			f.complete()
			tt.verifier(t, f)
		})
	}
}

// Test_future_Cancel_isSettledWithItsErrorWhenItReturns is the plain end to end shape of a
// cancellation: whatever Cancel does internally, a waiter which shows up after it returned gets the
// error.
func Test_future_Cancel_isSettledWithItsErrorWhenItReturns(t *testing.T) {
	cancelErr := errors.New("Uh oh")
	f := &future{}
	f.Cancel(true, cancelErr)
	assert.True(t, f.settled())
	assert.Equal(t, cancelErr, f.result())
	assert.Equal(t, cancelErr, f.AwaitCompletion(context.Background()))
}

// Test_future_settleReleasesEveryWaiter pins the two properties the channel buys us: one settle
// releases every waiter (a poll loop did that too, each on its own tick), and every one of them sees
// the state the settling goroutine published, not an intermediate step of it.
//
// The second half is what used to break: while AwaitCompletion left its loop on a terminal flag and
// only then read the error, an error published after that flag was invisible and the waiter reported
// the Canceled sentinel instead - a real flake, and one that needed the store order of several
// separate atomics to stay correct. Closing the channel is a single happens-before edge, so a waiter
// which wakes up at all sees every store that preceded the close. Waiters registered both before and
// after the settle to cover the lazy channel creation racing with the close.
func Test_future_settleReleasesEveryWaiter(t *testing.T) {
	const waiters = 32
	cancelErr := errors.New("Uh oh")

	f := &future{}
	results := make(chan error, 2*waiters)
	waiting := sync.WaitGroup{}
	for range waiters {
		waiting.Go(func() {
			results <- f.AwaitCompletion(t.Context())
		})
	}

	f.Cancel(true, cancelErr)

	// The late waiters exercise the same channel after it got closed.
	for range waiters {
		waiting.Go(func() {
			results <- f.AwaitCompletion(t.Context())
		})
	}
	waiting.Wait()

	close(results)
	seen := 0
	for err := range results {
		seen++
		assert.Same(t, cancelErr, err, "every waiter has to see the error the cancellation carried")
	}
	assert.Equal(t, 2*waiters, seen)
}

// Test_future_settleWakesTheWaiterImmediately guards the reason for the channel: a waiter used to
// notice a settle only on the next tick of a 10ms poll loop, so a batch of futures paid that latency
// one after the other. A closed channel hands the waiter back right away, and being generous by two
// orders of magnitude keeps this from being a benchmark.
func Test_future_settleWakesTheWaiterImmediately(t *testing.T) {
	f := &future{}
	done := make(chan time.Duration, 1)
	started := make(chan struct{})
	go func() {
		close(started)
		start := time.Now()
		assert.NoError(t, f.AwaitCompletion(t.Context()))
		done <- time.Since(start)
	}()
	<-started
	// Give the waiter a moment to actually park on the channel, so we measure the wake up and not the
	// go routine start up.
	time.Sleep(20 * time.Millisecond)

	f.complete()

	select {
	case waited := <-done:
		assert.Less(t, waited, 100*time.Millisecond, "the waiter has to wake up on the settle")
	case <-time.After(5 * time.Second):
		t.Fatal("AwaitCompletion did not return after the future completed")
	}
}

// Test_future_settleTwiceDoesNotPanic covers the double settle: Cancel is part of the public
// CompletionFuture, so a caller outside the pool can cancel a work item which a worker then runs and
// completes anyway (and the other way round for a future which is completed before somebody cancels
// it). Closing an already closed channel panics, so both orders have to be no-ops for the second
// settle.
func Test_future_settleTwiceDoesNotPanic(t *testing.T) {
	t.Run("cancel then complete", func(t *testing.T) {
		cancelErr := errors.New("Uh oh")
		f := &future{}
		f.Cancel(true, cancelErr)
		assert.NotPanics(t, f.complete)
		assert.True(t, f.completed.Load())
		assert.Equal(t, cancelErr, f.AwaitCompletion(t.Context()), "the cancellation error survives a late completion")
	})
	t.Run("complete then cancel", func(t *testing.T) {
		f := &future{}
		f.complete()
		assert.NoError(t, f.AwaitCompletion(t.Context()))
		assert.NotPanics(t, func() { f.Cancel(true, errors.New("Uh oh")) })
	})
	t.Run("complete twice", func(t *testing.T) {
		f := &future{}
		f.complete()
		assert.NotPanics(t, f.complete)
	})
	t.Run("cancel twice", func(t *testing.T) {
		f := &future{}
		f.Cancel(false, nil)
		assert.NotPanics(t, func() { f.Cancel(false, nil) })
	})
	t.Run("concurrent settles", func(t *testing.T) {
		f := &future{}
		settling := sync.WaitGroup{}
		for i := range 16 {
			settling.Go(func() {
				if i%2 == 0 {
					f.complete()
					return
				}
				f.Cancel(false, nil)
			})
		}
		settling.Wait()
		assert.True(t, f.settled())
	})
}

// Test_future_AwaitCompletion_ctxCancellationWhileWaiting pins what a waiter gets when its own
// context goes away while the future is still in flight: the context error, right away, and the
// future stays unsettled - the work item is still queued, nobody cancelled it.
func Test_future_AwaitCompletion_ctxCancellationWhileWaiting(t *testing.T) {
	f := &future{}
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() {
		done <- f.AwaitCompletion(ctx)
	}()
	time.Sleep(20 * time.Millisecond)

	cancel()

	select {
	case err := <-done:
		assert.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("AwaitCompletion did not return after its context was canceled")
	}
	assert.False(t, f.settled(), "a waiter giving up must not settle the future")

	// A waiter with a healthy context still gets the result afterwards.
	f.complete()
	assert.NoError(t, f.AwaitCompletion(t.Context()))
}

// Test_future_AwaitCompletion_alreadyDoneCtx keeps the precedence the poll loop had: a context which
// is already done wins over a settled future, because the caller stopped being interested.
func Test_future_AwaitCompletion_alreadyDoneCtx(t *testing.T) {
	f := &future{}
	f.complete()
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	assert.ErrorIs(t, f.AwaitCompletion(ctx), context.Canceled)
}

// Test_future_AwaitCompletion_doesNotLeakWaiters checks that no waiter stays parked on the channel,
// whichever way it left: released by a settle, or given up on its own context. A nil channel would
// have parked all of them forever, which is exactly the failure mode the lazy creation has to avoid.
func Test_future_AwaitCompletion_doesNotLeakWaiters(t *testing.T) {
	baseline := runtime.NumGoroutine()

	settled := &future{}
	abandoned := &future{}
	ctx, cancel := context.WithCancel(context.Background())

	waiting := sync.WaitGroup{}
	for range 32 {
		waiting.Go(func() {
			assert.NoError(t, settled.AwaitCompletion(t.Context()))
		})
		waiting.Go(func() {
			assert.ErrorIs(t, abandoned.AwaitCompletion(ctx), context.Canceled)
		})
	}

	settled.complete()
	cancel()
	waiting.Wait()

	assertNoWaiterLeak(t, baseline)
}

// assertNoWaiterLeak checks that the go routine count settles back at (or below) the baseline. Go
// routines the runtime tears down lazily make a single sample flaky, hence the retries.
func assertNoWaiterLeak(t *testing.T, baseline int) {
	t.Helper()
	current := runtime.NumGoroutine()
	for i := 0; i < 100 && current > baseline; i++ {
		time.Sleep(20 * time.Millisecond)
		current = runtime.NumGoroutine()
	}
	if current > baseline {
		buf := make([]byte, 1<<20)
		n := runtime.Stack(buf, true)
		t.Fatalf("goroutine leak detected: baseline %d, now %d\n%s", baseline, current, buf[:n])
	}
}
