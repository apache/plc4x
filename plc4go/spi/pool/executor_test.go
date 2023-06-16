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
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_executor_Close(t *testing.T) {
	type fields struct {
		running      bool
		shutdown     bool
		worker       []*worker
		queueDepth   int
		workItems    chan workItem
		traceWorkers bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "close it",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				running:      tt.fields.running,
				shutdown:     tt.fields.shutdown,
				worker:       tt.fields.worker,
				queueDepth:   tt.fields.queueDepth,
				workItems:    tt.fields.workItems,
				traceWorkers: tt.fields.traceWorkers,
				log:          produceTestingLogger(t),
			}
			tt.wantErr(t, e.Close(), fmt.Sprintf("Close()"))
		})
	}
}

func Test_executor_IsRunning(t *testing.T) {
	type fields struct {
		running      bool
		shutdown     bool
		worker       []*worker
		queueDepth   int
		workItems    chan workItem
		traceWorkers bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "no",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				running:      tt.fields.running,
				shutdown:     tt.fields.shutdown,
				worker:       tt.fields.worker,
				queueDepth:   tt.fields.queueDepth,
				workItems:    tt.fields.workItems,
				traceWorkers: tt.fields.traceWorkers,
				log:          produceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, e.IsRunning(), "IsRunning()")
		})
	}
}

func Test_executor_Start(t *testing.T) {
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

func Test_executor_Stop(t *testing.T) {
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

func Test_executor_Submit(t *testing.T) {
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

func Test_executor_getWorkerWaitGroup(t *testing.T) {
	type fields struct {
		running      bool
		shutdown     bool
		worker       []*worker
		queueDepth   int
		workItems    chan workItem
		traceWorkers bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *sync.WaitGroup
	}{
		{
			name: "get it",
			want: &sync.WaitGroup{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				running:      tt.fields.running,
				shutdown:     tt.fields.shutdown,
				worker:       tt.fields.worker,
				queueDepth:   tt.fields.queueDepth,
				workItems:    tt.fields.workItems,
				traceWorkers: tt.fields.traceWorkers,
				log:          produceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, e.getWorkerWaitGroup(), "getWorkerWaitGroup()")
		})
	}
}

func Test_executor_getWorksItems(t *testing.T) {
	type fields struct {
		running      bool
		shutdown     bool
		worker       []*worker
		queueDepth   int
		workItems    chan workItem
		traceWorkers bool
		log          zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   chan workItem
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				running:      tt.fields.running,
				shutdown:     tt.fields.shutdown,
				worker:       tt.fields.worker,
				queueDepth:   tt.fields.queueDepth,
				workItems:    tt.fields.workItems,
				traceWorkers: tt.fields.traceWorkers,
				log:          tt.fields.log,
			}
			assert.Equalf(t, tt.want, e.getWorksItems(), "getWorksItems()")
		})
	}
}

func Test_executor_isTraceWorkers(t *testing.T) {
	type fields struct {
		running      bool
		shutdown     bool
		worker       []*worker
		queueDepth   int
		workItems    chan workItem
		traceWorkers bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "it is not",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				running:      tt.fields.running,
				shutdown:     tt.fields.shutdown,
				worker:       tt.fields.worker,
				queueDepth:   tt.fields.queueDepth,
				workItems:    tt.fields.workItems,
				traceWorkers: tt.fields.traceWorkers,
				log:          produceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, e.isTraceWorkers(), "isTraceWorkers()")
		})
	}
}

func Test_executor_String(t *testing.T) {
	type fields struct {
		running      bool
		shutdown     bool
		worker       []*worker
		queueDepth   int
		workItems    chan workItem
		traceWorkers bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			fields: fields{
				running:  true,
				shutdown: true,
				worker: []*worker{
					{
						id:          1,
						shutdown:    atomic.Bool{},
						interrupted: atomic.Bool{},
						hasEnded:    atomic.Bool{},
						lastReceived: func() atomic.Value {
							value := atomic.Value{}
							value.Store(time.Time{})
							return value
						}(),
					},
				},
				queueDepth:   2,
				traceWorkers: true,
			},
			want: `
╔═executor════════════════════════════════════════════════════════════════════════════════════════════════╗
║╔═running╗╔═shutdown╗                                                                                    ║
║║b1 true ║║ b1 true ║                                                                                    ║
║╚════════╝╚═════════╝                                                                                    ║
║╔═worker/value/worker═══════════════════════════════════════════════════════════════════════════════════╗║
║║╔═id═════════════════╗╔═shutdown╗╔═interrupted╗╔═interrupter╗╔═hasEnded╗╔═lastReceived════════════════╗║║
║║║0x0000000000000001 1║║b0 false ║║  b0 false  ║║0 element(s)║║b0 false ║║0001-01-01 00:00:00 +0000 UTC║║║
║║╚════════════════════╝╚═════════╝╚════════════╝╚════════════╝╚═════════╝╚═════════════════════════════╝║║
║╚═══════════════════════════════════════════════════════════════════════════════════════════════════════╝║
║╔═queueDepth═════════╗╔═workItems══╗╔═traceWorkers╗                                                      ║
║║0x0000000000000002 2║║0 element(s)║║   b1 true   ║                                                      ║
║╚════════════════════╝╚════════════╝╚═════════════╝                                                      ║
╚═════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &executor{
				running:      tt.fields.running,
				shutdown:     tt.fields.shutdown,
				worker:       tt.fields.worker,
				queueDepth:   tt.fields.queueDepth,
				workItems:    tt.fields.workItems,
				traceWorkers: tt.fields.traceWorkers,
				log:          produceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, e.String(), "String()")
		})
	}
}
