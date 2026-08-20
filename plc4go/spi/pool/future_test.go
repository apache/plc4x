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

// Test_future_Cancel_publishesErrorBeforeReleasingTheWaiter guards the store ordering in
// future.Cancel: settled leaves the poll loop of AwaitCompletion as soon as it observes one of the
// terminal flags and result reads the error only afterwards, so an error published after those flags
// can be missed entirely, making the waiter report the Canceled sentinel instead of the error the
// caller passed in. That used to flake Test_future_AwaitCompletion/completes_canceled_with_particular_error
// under load.
//
// Racing a spinning reader against Cancel does not pin this down: the window between two adjacent
// atomic stores is a handful of instructions, and a reader loop only lands inside it when the race
// detector instruments every store (measured on the broken order: zero misses in 20 000 rounds
// without -race, dozens with it). So instead of racing, this steps Cancel through its stores and
// checks the invariant after every single one of them - every intermediate state a waiter could ever
// observe is visited, deterministically and without -race: once the future reports itself settled,
// result has to be the final answer already.
func Test_future_Cancel_publishesErrorBeforeReleasingTheWaiter(t *testing.T) {
	cancelErr := errors.New("Uh oh")
	tests := []struct {
		name string
		err  error
		want error
	}{
		{name: "cancel carrying an error", err: cancelErr, want: cancelErr},
		{name: "plain cancel", err: nil, want: Canceled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &future{}
			stores := 0
			f.cancel(true, tt.err, func() {
				stores++
				if !f.settled() {
					return
				}
				assert.Equal(t, tt.want, f.result(),
					"store %d released AwaitCompletion with the wrong result", stores)
			})
			assert.Greater(t, stores, 1, "the observer has to see every store")
			assert.True(t, f.settled(), "a cancelled future has to be settled")
			assert.Equal(t, tt.want, f.result())
			assert.True(t, f.interruptRequested.Load())
			assert.True(t, f.cancelRequested.Load())
		})
	}
}

// Test_future_Cancel_isSettledWithItsErrorWhenItReturns is the plain end to end shape of the above:
// whatever Cancel does internally, a waiter which shows up after it returned gets the error.
func Test_future_Cancel_isSettledWithItsErrorWhenItReturns(t *testing.T) {
	cancelErr := errors.New("Uh oh")
	f := &future{}
	f.Cancel(true, cancelErr)
	assert.True(t, f.settled())
	assert.Equal(t, cancelErr, f.result())
	assert.Equal(t, cancelErr, f.AwaitCompletion(context.Background()))
}
