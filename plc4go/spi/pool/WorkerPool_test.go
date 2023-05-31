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
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestExecutor_Start(t *testing.T) {
	type fields struct {
		running      bool
		shutdown     bool
		worker       []*worker
		queue        chan workItem
		traceWorkers bool
	}
	tests := []struct {
		name      string
		fields    fields
		shouldRun bool
	}{
		{
			name:      "Start fresh",
			shouldRun: true,
		},
		{
			name: "Start running",
			fields: fields{
				running: true,
			},
			shouldRun: true,
		},
		{
			name: "Start stopping",
			fields: fields{
				running:  true,
				shutdown: true,
			},
			shouldRun: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				running:      tt.fields.running,
				shutdown:     tt.fields.shutdown,
				worker:       tt.fields.worker,
				workItems:    tt.fields.queue,
				traceWorkers: tt.fields.traceWorkers,
			}
			e.Start()
			assert.Equal(t, tt.shouldRun, e.IsRunning(), "should be running")
		})
	}
}

func TestExecutor_Stop(t *testing.T) {
	type fields struct {
		running      bool
		shutdown     bool
		worker       []*worker
		queue        chan workItem
		traceWorkers bool
	}
	tests := []struct {
		name      string
		fields    fields
		shouldRun bool
	}{
		{
			name:      "Stop stopped",
			shouldRun: false,
		},
		{
			name: "Stop running",
			fields: fields{
				running: true,
				queue:   make(chan workItem),
				worker: []*worker{
					func() *worker {
						w := &worker{}
						w.initialize()
						return w
					}(),
				},
			},
			shouldRun: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				running:      tt.fields.running,
				shutdown:     tt.fields.shutdown,
				worker:       tt.fields.worker,
				workItems:    tt.fields.queue,
				traceWorkers: tt.fields.traceWorkers,
			}
			e.Stop()
		})
	}
}

func TestExecutor_Submit(t *testing.T) {
	type fields struct {
		running      bool
		shutdown     bool
		worker       []*worker
		queue        chan workItem
		traceWorkers bool
	}
	type args struct {
		workItemId int32
		runnable   Runnable
		context    context.Context
	}
	tests := []struct {
		name                      string
		fields                    fields
		args                      args
		completionFutureValidator func(t *testing.T, future CompletionFuture) bool
		waitForCompletion         bool
	}{
		{
			name: "submitting nothing",
			completionFutureValidator: func(t *testing.T, completionFuture CompletionFuture) bool {
				return assert.Error(t, completionFuture.(*future).err.Load().(error))
			},
		},
		{
			name: "submit canceled",
			fields: fields{
				queue: make(chan workItem, 0),
			},
			args: args{
				workItemId: 13,
				runnable: func() {
					// We do something for 3 seconds
					<-time.NewTimer(3 * time.Second).C
				},
				context: func() context.Context {
					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return ctx
				}(),
			},
			completionFutureValidator: func(t *testing.T, completionFuture CompletionFuture) bool {
				err := completionFuture.(*future).err.Load().(error)
				return assert.Error(t, err)
			},
		},
		{
			name: "Submit something which doesn't complete",
			fields: fields{
				queue: make(chan workItem, 1),
			},
			args: args{
				workItemId: 13,
				runnable: func() {
					// We do something for 3 seconds
					<-time.NewTimer(3 * time.Second).C
				},
				context: context.TODO(),
			},
			completionFutureValidator: func(t *testing.T, completionFuture CompletionFuture) bool {
				completed := completionFuture.(*future).completed.Load()
				return assert.False(t, completed)
			},
		},
		{
			name: "Submit something which does complete",
			fields: func() fields {
				var executor = NewFixedSizeExecutor(1, 1).(*executor)
				return fields{
					running:      executor.running,
					shutdown:     executor.shutdown,
					worker:       executor.worker,
					queue:        executor.workItems,
					traceWorkers: true,
				}
			}(),
			args: args{
				workItemId: 13,
				runnable: func() {
					// NOOP
				},
				context: context.TODO(),
			},
			completionFutureValidator: func(t *testing.T, completionFuture CompletionFuture) bool {
				completed := completionFuture.(*future).completed.Load()
				return assert.True(t, completed)
			},
			waitForCompletion: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				running:      tt.fields.running,
				shutdown:     tt.fields.shutdown,
				worker:       tt.fields.worker,
				workItems:    tt.fields.queue,
				traceWorkers: tt.fields.traceWorkers,
			}
			e.Start()
			completionFuture := e.Submit(tt.args.context, tt.args.workItemId, tt.args.runnable)
			if tt.waitForCompletion {
				assert.NoError(t, completionFuture.AwaitCompletion(testContext(t)))
			}
			assert.True(t, tt.completionFutureValidator(t, completionFuture), "Submit(%v, %v)", tt.args.workItemId, tt.args.runnable)
		})
	}
}

func TestNewFixedSizeExecutor(t *testing.T) {
	type args struct {
		numberOfWorkers int
		queueDepth      int
		options         []options.WithOption
	}
	tests := []struct {
		name              string
		args              args
		executorValidator func(*testing.T, *executor) bool
	}{
		{
			name: "new Executor",
			args: args{
				numberOfWorkers: 13,
				queueDepth:      14,
				options:         []options.WithOption{WithExecutorOptionTracerWorkers(true)},
			},
			executorValidator: func(t *testing.T, e *executor) bool {
				return !e.running && !e.shutdown && len(e.worker) == 13 && cap(e.workItems) == 14
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixedSizeExecutor := NewFixedSizeExecutor(tt.args.numberOfWorkers, tt.args.queueDepth, tt.args.options...)
			defer fixedSizeExecutor.Stop()
			assert.True(t, tt.executorValidator(t, fixedSizeExecutor.(*executor)), "NewFixedSizeExecutor(%v, %v, %v)", tt.args.numberOfWorkers, tt.args.queueDepth, tt.args.options)
		})
	}
}

func TestNewDynamicExecutor(t *testing.T) {
	type args struct {
		numberOfWorkers int
		queueDepth      int
		options         []options.WithOption
	}
	tests := []struct {
		name              string
		args              args
		manipulator       func(*testing.T, *executor)
		executorValidator func(*testing.T, *executor) bool
	}{
		{
			name: "new Executor",
			args: args{
				numberOfWorkers: 13,
				queueDepth:      14,
				options:         []options.WithOption{WithExecutorOptionTracerWorkers(true)},
			},
			executorValidator: func(t *testing.T, e *executor) bool {
				assert.False(t, e.running)
				assert.False(t, e.shutdown)
				assert.Len(t, e.worker, 1)
				assert.Equal(t, cap(e.workItems), 14)
				return true
			},
		},
		{
			name: "test scaling",
			args: args{
				numberOfWorkers: 2,
				queueDepth:      2,
				options:         []options.WithOption{WithExecutorOptionTracerWorkers(true)},
			},
			manipulator: func(t *testing.T, e *executor) {
				{
					oldUpScaleInterval := upScaleInterval
					t.Cleanup(func() {
						t.Logf("Ressetting up scale interval to %v", oldUpScaleInterval)
						upScaleInterval = oldUpScaleInterval
					})
					upScaleInterval = 10 * time.Millisecond
					t.Logf("Changed up scale interval to %v", upScaleInterval)
				}
				{
					oldDownScaleInterval := downScaleInterval
					t.Cleanup(func() {
						t.Logf("Ressetting down scale interval to %v", oldDownScaleInterval)
						downScaleInterval = oldDownScaleInterval
					})
					downScaleInterval = 10 * time.Millisecond
					t.Logf("Changed down scale interval to %v", downScaleInterval)
				}
				{
					oldTimeToBecomeUnused := timeToBecomeUnused
					t.Cleanup(func() {
						t.Logf("Ressetting time to be become unused to %v", oldTimeToBecomeUnused)
						timeToBecomeUnused = oldTimeToBecomeUnused
					})
					timeToBecomeUnused = 100 * time.Millisecond
				}
				t.Log("fill some jobs")
				go func() {
					for i := 0; i < 500; i++ {
						e.workItems <- workItem{
							workItemId: int32(i),
							runnable: func() {
								max := 100
								min := 10
								sleepTime := time.Duration(rand.Intn(max-min)+min) * time.Millisecond
								t.Logf("Sleeping for %v", sleepTime)
								time.Sleep(sleepTime)
							},
							completionFuture: &future{},
						}
					}
				}()
			},
			executorValidator: func(t *testing.T, e *executor) bool {
				time.Sleep(500 * time.Millisecond)
				assert.False(t, e.running)
				assert.False(t, e.shutdown)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixedSizeExecutor := NewDynamicExecutor(tt.args.numberOfWorkers, tt.args.queueDepth, tt.args.options...)
			defer fixedSizeExecutor.Stop()
			if tt.manipulator != nil {
				tt.manipulator(t, fixedSizeExecutor.(*executor))
			}
			assert.True(t, tt.executorValidator(t, fixedSizeExecutor.(*executor)), "NewFixedSizeExecutor(%v, %v, %v)", tt.args.numberOfWorkers, tt.args.queueDepth, tt.args.options)
		})
	}
}

func TestWithExecutorOptionTracerWorkers(t *testing.T) {
	type args struct {
		traceWorkers bool
	}
	tests := []struct {
		name              string
		args              args
		executorValidator options.WithOption
	}{
		{
			name:              "option should set option",
			args:              args{traceWorkers: true},
			executorValidator: &tracerWorkersOption{traceWorkers: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.executorValidator, WithExecutorOptionTracerWorkers(tt.args.traceWorkers))
		})
	}
}

func TestWorkItem_String(t *testing.T) {
	type fields struct {
		workItemId       int32
		runnable         Runnable
		completionFuture *future
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Simple test",
			want: "Workitem{wid:0}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &workItem{
				workItemId:       tt.fields.workItemId,
				runnable:         tt.fields.runnable,
				completionFuture: tt.fields.completionFuture,
			}
			assert.Equalf(t, tt.want, w.String(), "String()")
		})
	}
}

func TestWorker_work(t *testing.T) {
	type fields struct {
		id       int
		executor *executor
	}
	tests := []struct {
		name                       string
		fields                     fields
		timeBeforeFirstValidation  time.Duration
		firstValidation            func(*testing.T, *worker)
		timeBeforeManipulation     time.Duration
		manipulator                func(*worker)
		timeBeforeSecondValidation time.Duration
		secondValidation           func(*testing.T, *worker)
	}{
		{
			name: "Worker should work till shutdown (even if it panics)",
			fields: fields{
				id: 0,
				executor: func() *executor {
					e := &executor{
						workItems:    make(chan workItem),
						traceWorkers: true,
					}
					go func() {
						e.workItems <- workItem{
							workItemId: 0,
							runnable: func() {
								panic("Oh no what should I do???")
							},
							completionFuture: &future{},
						}
					}()
					return e
				}(),
			},
			timeBeforeFirstValidation: 50 * time.Millisecond,
			firstValidation: func(t *testing.T, w *worker) {
				assert.False(t, w.hasEnded.Load(), "should not be ended")
			},
			manipulator: func(w *worker) {
				w.shutdown.Store(true)
				w.interrupter <- struct{}{}
			},
			timeBeforeSecondValidation: 150 * time.Millisecond,
			secondValidation: func(t *testing.T, w *worker) {
				assert.True(t, w.hasEnded.Load(), "should be ended")
			},
		},
		{
			name: "Worker should work till shutdown",
			fields: fields{
				id: 1,
				executor: func() *executor {
					e := &executor{
						workItems:    make(chan workItem),
						traceWorkers: true,
					}
					go func() {
						e.workItems <- workItem{
							workItemId: 0,
							runnable: func() {
								time.Sleep(time.Millisecond * 70)
							},
							completionFuture: &future{},
						}
					}()
					return e
				}(),
			},
			timeBeforeFirstValidation: 50 * time.Millisecond,
			firstValidation: func(t *testing.T, w *worker) {
				assert.False(t, w.hasEnded.Load(), "should not be ended")
			},
			manipulator: func(w *worker) {
				w.shutdown.Store(true)
			},
			timeBeforeSecondValidation: 150 * time.Millisecond,
			secondValidation: func(t *testing.T, w *worker) {
				assert.True(t, w.hasEnded.Load(), "should be ended")
			},
		},
		{
			name: "Work interrupted",
			fields: fields{
				id: 1,
				executor: func() *executor {
					e := &executor{
						workItems:    make(chan workItem),
						traceWorkers: true,
					}
					return e
				}(),
			},
			timeBeforeFirstValidation: 50 * time.Millisecond,
			firstValidation: func(t *testing.T, w *worker) {
				assert.False(t, w.hasEnded.Load(), "should not be ended")
			},
			manipulator: func(w *worker) {
				w.shutdown.Store(true)
				w.interrupter <- struct{}{}
			},
			timeBeforeSecondValidation: 150 * time.Millisecond,
			secondValidation: func(t *testing.T, w *worker) {
				assert.True(t, w.hasEnded.Load(), "should be ended")
			},
		},
		{
			name: "Work on canceled",
			fields: fields{
				id: 1,
				executor: func() *executor {
					e := &executor{
						workItems:    make(chan workItem),
						traceWorkers: true,
					}
					go func() {
						completionFuture := &future{}
						completionFuture.cancelRequested.Store(true)
						e.workItems <- workItem{
							workItemId: 0,
							runnable: func() {
								time.Sleep(time.Millisecond * 70)
							},
							completionFuture: completionFuture,
						}
					}()
					return e
				}(),
			},
			timeBeforeManipulation: 50 * time.Millisecond,
			manipulator: func(w *worker) {
				w.shutdown.Store(true)
				w.interrupter <- struct{}{}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &worker{
				id:          tt.fields.id,
				interrupter: make(chan struct{}, 1),
				executor:    tt.fields.executor,
			}
			go w.work()
			if tt.firstValidation != nil {
				time.Sleep(tt.timeBeforeFirstValidation)
				t.Logf("firstValidation after %v", tt.timeBeforeFirstValidation)
				tt.firstValidation(t, w)
			}
			if tt.manipulator != nil {
				time.Sleep(tt.timeBeforeManipulation)
				t.Logf("manipulator after %v", tt.timeBeforeManipulation)
				tt.manipulator(w)
			}
			if tt.secondValidation != nil {
				time.Sleep(tt.timeBeforeSecondValidation)
				t.Logf("secondValidation after %v", tt.timeBeforeSecondValidation)
				tt.secondValidation(t, w)
			}
		})
	}
}

func Test_future_AwaitCompletion(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		completer func(*future)
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "completes with error",
			args: args{ctx: context.TODO()},
			completer: func(f *future) {
				f.Cancel(false, errors.New("Uh oh"))
			},
			wantErr: assert.Error,
		},
		{
			name: "completes regular",
			args: args{ctx: context.TODO()},
			completer: func(f *future) {
				time.Sleep(time.Millisecond * 30)
				f.complete()
			},
			wantErr: assert.NoError,
		},
		{
			name: "completes not int time",
			args: args{ctx: func() context.Context {
				deadline, cancel := context.WithDeadline(context.Background(), time.Now().Add(30*time.Millisecond))
				t.Cleanup(cancel)
				return deadline
			}()},
			completer: func(f *future) {
				time.Sleep(time.Millisecond * 300)
			},
			wantErr: assert.Error,
		},
		{
			name: "completes canceled without error",
			args: args{ctx: context.TODO()},
			completer: func(f *future) {
				time.Sleep(time.Millisecond * 300)
				f.Cancel(true, nil)
			},
			wantErr: func(t assert.TestingT, err error, i ...any) bool {
				assert.Same(t, Canceled, err)
				return true
			},
		},
		{
			name: "completes canceled with particular error",
			args: args{ctx: context.TODO()},
			completer: func(f *future) {
				time.Sleep(time.Millisecond * 300)
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
			go tt.completer(f)
			tt.wantErr(t, f.AwaitCompletion(tt.args.ctx), fmt.Sprintf("AwaitCompletion(%v)", tt.args.ctx))
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

// from: https://github.com/golang/go/issues/36532#issuecomment-575535452
func testContext(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return ctx
}
