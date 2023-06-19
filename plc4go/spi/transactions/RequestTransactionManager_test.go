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
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/pool"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewRequestTransactionManager(t *testing.T) {
	type args struct {
		numberOfConcurrentRequests       int
		requestTransactionManagerOptions []options.WithOption
	}
	tests := []struct {
		name  string
		args  args
		setup func(t *testing.T, args *args)
		want  RequestTransactionManager
	}{
		{
			name: "just create one",
			want: &requestTransactionManager{
				workLog:  *list.New(),
				executor: sharedExecutorInstance,
			},
		},
		{
			name: "just create one with option",
			args: args{
				numberOfConcurrentRequests: 2,
				requestTransactionManagerOptions: []options.WithOption{
					WithCustomExecutor(sharedExecutorInstance),
				},
			},
			want: &requestTransactionManager{
				numberOfConcurrentRequests: 2,
				workLog:                    *list.New(),
				executor:                   sharedExecutorInstance,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.args)
			}
			if got := NewRequestTransactionManager(tt.args.numberOfConcurrentRequests, tt.args.requestTransactionManagerOptions...); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewRequestTransactionManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCustomExecutor(t *testing.T) {
	type args struct {
		executor pool.Executor
	}
	tests := []struct {
		name string
		args args
		want options.WithOption
	}{
		{
			name: "with a option",
			args: args{
				executor: sharedExecutorInstance,
			},
			want: WithCustomExecutor(sharedExecutorInstance),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customExecutor := WithCustomExecutor(tt.args.executor)
			assert.NotNil(t, customExecutor)
		})
	}
}

func Test_requestTransactionManager_SetNumberOfConcurrentRequests(t *testing.T) {
	type fields struct {
		runningRequests            []*requestTransaction
		numberOfConcurrentRequests int
		currentTransactionId       int32
		workLog                    list.List
		executor                   pool.Executor
	}
	type args struct {
		numberOfConcurrentRequests int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set a number",
		},
		{
			name: "set a number on running requests",
			fields: fields{
				runningRequests: []*requestTransaction{
					{}, // empty one is sufficient
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &requestTransactionManager{
				runningRequests:            tt.fields.runningRequests,
				numberOfConcurrentRequests: tt.fields.numberOfConcurrentRequests,
				currentTransactionId:       tt.fields.currentTransactionId,
				workLog:                    tt.fields.workLog,
				executor:                   tt.fields.executor,
			}
			r.SetNumberOfConcurrentRequests(tt.args.numberOfConcurrentRequests)
		})
	}
}

func Test_requestTransactionManager_StartTransaction(t *testing.T) {
	type fields struct {
		runningRequests                     []*requestTransaction
		numberOfConcurrentRequests          int
		currentTransactionId                int32
		workLog                             list.List
		executor                            pool.Executor
		traceTransactionManagerTransactions bool
	}
	tests := []struct {
		name        string
		fields      fields
		setup       func(t *testing.T, fields *fields)
		manipulator func(t *testing.T, manager *requestTransactionManager)
		wantAssert  func(t *testing.T, requestTransaction RequestTransaction) bool
	}{
		{
			name: "start one",
			wantAssert: func(t *testing.T, requestTransaction RequestTransaction) bool {
				assert.False(t, requestTransaction.IsCompleted())
				return true
			},
		},
		{
			name: "start one in shutdown",
			manipulator: func(t *testing.T, manager *requestTransactionManager) {
				manager.shutdown.Store(true)
			},
			wantAssert: func(t *testing.T, requestTransaction RequestTransaction) bool {
				assert.True(t, requestTransaction.IsCompleted())
				assert.Error(t, requestTransaction.AwaitCompletion(context.Background()))
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			r := &requestTransactionManager{
				runningRequests:                     tt.fields.runningRequests,
				numberOfConcurrentRequests:          tt.fields.numberOfConcurrentRequests,
				currentTransactionId:                tt.fields.currentTransactionId,
				workLog:                             tt.fields.workLog,
				executor:                            tt.fields.executor,
				traceTransactionManagerTransactions: tt.fields.traceTransactionManagerTransactions,
				log:                                 testutils.ProduceTestingLogger(t),
			}
			if tt.manipulator != nil {
				tt.manipulator(t, r)
			}
			if got := r.StartTransaction(); !assert.True(t, tt.wantAssert(t, got)) {
				t.Errorf("StartTransaction() = %v", got)
			}
		})
	}
}

func Test_requestTransactionManager_endRequest(t *testing.T) {
	type fields struct {
		runningRequests            []*requestTransaction
		numberOfConcurrentRequests int
		currentTransactionId       int32
		workLog                    list.List
		executor                   pool.Executor
	}
	type args struct {
		transaction *requestTransaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "end request with unknown transaction",
			args: args{
				transaction: &requestTransaction{},
			},
			wantErr: true,
		},
		{
			name: "end request",
			args: args{
				transaction: &requestTransaction{},
			},
			fields: fields{
				runningRequests: []*requestTransaction{
					{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &requestTransactionManager{
				runningRequests:            tt.fields.runningRequests,
				numberOfConcurrentRequests: tt.fields.numberOfConcurrentRequests,
				currentTransactionId:       tt.fields.currentTransactionId,
				workLog:                    tt.fields.workLog,
				executor:                   tt.fields.executor,
			}
			if err := r.endRequest(tt.args.transaction); (err != nil) != tt.wantErr {
				t.Errorf("endRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_requestTransactionManager_failRequest(t *testing.T) {
	type fields struct {
		runningRequests            []*requestTransaction
		numberOfConcurrentRequests int
		currentTransactionId       int32
		workLog                    list.List
		executor                   pool.Executor
	}
	type args struct {
		transaction *requestTransaction
		err         error
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		wantErr   bool
	}{
		{
			name: "fail a request",
			args: args{
				transaction: &requestTransaction{},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				completionFutureMock := NewMockCompletionFuture(t)
				expect := completionFutureMock.EXPECT()
				expect.Cancel(true, nil).Return()
				var completionFuture pool.CompletionFuture = completionFutureMock
				args.transaction.completionFuture.Store(&completionFuture)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			r := &requestTransactionManager{
				runningRequests:            tt.fields.runningRequests,
				numberOfConcurrentRequests: tt.fields.numberOfConcurrentRequests,
				currentTransactionId:       tt.fields.currentTransactionId,
				workLog:                    tt.fields.workLog,
				executor:                   tt.fields.executor,
				log:                        testutils.ProduceTestingLogger(t),
			}
			if err := r.failRequest(tt.args.transaction, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("failRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_requestTransactionManager_getNumberOfActiveRequests(t *testing.T) {
	type fields struct {
		runningRequests            []*requestTransaction
		numberOfConcurrentRequests int
		currentTransactionId       int32
		workLog                    list.List
		executor                   pool.Executor
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "get em",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &requestTransactionManager{
				runningRequests:            tt.fields.runningRequests,
				numberOfConcurrentRequests: tt.fields.numberOfConcurrentRequests,
				currentTransactionId:       tt.fields.currentTransactionId,
				workLog:                    tt.fields.workLog,
				executor:                   tt.fields.executor,
			}
			if got := r.getNumberOfActiveRequests(); got != tt.want {
				t.Errorf("getNumberOfActiveRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_requestTransactionManager_processWorklog(t *testing.T) {
	type fields struct {
		runningRequests            []*requestTransaction
		numberOfConcurrentRequests int
		currentTransactionId       int32
		workLog                    list.List
		executor                   pool.Executor
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "process nothing",
		},
		{
			name: "process one",
			fields: fields{
				numberOfConcurrentRequests: 100,
				workLog: func() list.List {
					l := list.New()
					l.PushBack(&requestTransaction{})
					return *l
				}(),
				executor: sharedExecutorInstance,
			},
		},
		{
			name: "process two",
			fields: fields{
				numberOfConcurrentRequests: 100,
				workLog: func() list.List {
					l := list.New()
					var completionFuture pool.CompletionFuture = NewMockCompletionFuture(t)
					r1 := &requestTransaction{
						transactionId: 1,
					}
					r1.completionFuture.Store(&completionFuture)
					l.PushBack(r1)
					r2 := &requestTransaction{
						transactionId: 2,
					}
					r2.completionFuture.Store(&completionFuture)
					l.PushBack(r2)
					return *l
				}(),
				executor: sharedExecutorInstance,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &requestTransactionManager{
				runningRequests:            tt.fields.runningRequests,
				numberOfConcurrentRequests: tt.fields.numberOfConcurrentRequests,
				currentTransactionId:       tt.fields.currentTransactionId,
				workLog:                    tt.fields.workLog,
				executor:                   tt.fields.executor,
			}
			r.processWorklog()
		})
	}
}

func Test_requestTransactionManager_submitTransaction(t *testing.T) {
	type fields struct {
		runningRequests            []*requestTransaction
		numberOfConcurrentRequests int
		currentTransactionId       int32
		workLog                    list.List
		executor                   pool.Executor
	}
	type args struct {
		handle *requestTransaction
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "submit it",
			args: args{
				handle: &requestTransaction{
					operation: func() {
						// doesn't matter
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &requestTransactionManager{
				runningRequests:            tt.fields.runningRequests,
				numberOfConcurrentRequests: tt.fields.numberOfConcurrentRequests,
				currentTransactionId:       tt.fields.currentTransactionId,
				workLog:                    tt.fields.workLog,
				executor:                   tt.fields.executor,
			}
			r.submitTransaction(tt.args.handle)
		})
	}
}

func Test_requestTransactionManager_Close(t *testing.T) {
	type fields struct {
		runningRequests                     []*requestTransaction
		numberOfConcurrentRequests          int
		currentTransactionId                int32
		workLog                             list.List
		executor                            pool.Executor
		traceTransactionManagerTransactions bool
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func(t *testing.T, fields *fields)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "close it",
			setup: func(t *testing.T, fields *fields) {
				executor := NewMockExecutor(t)
				executor.EXPECT().Close().Return(nil)
				fields.executor = executor
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			r := &requestTransactionManager{
				runningRequests:                     tt.fields.runningRequests,
				numberOfConcurrentRequests:          tt.fields.numberOfConcurrentRequests,
				currentTransactionId:                tt.fields.currentTransactionId,
				workLog:                             tt.fields.workLog,
				executor:                            tt.fields.executor,
				traceTransactionManagerTransactions: tt.fields.traceTransactionManagerTransactions,
				log:                                 testutils.ProduceTestingLogger(t),
			}
			tt.wantErr(t, r.Close(), fmt.Sprintf("Close()"))
		})
	}
}

func Test_requestTransactionManager_CloseGraceful(t *testing.T) {
	type fields struct {
		runningRequests                     []*requestTransaction
		numberOfConcurrentRequests          int
		currentTransactionId                int32
		workLog                             list.List
		executor                            pool.Executor
		traceTransactionManagerTransactions bool
		log                                 zerolog.Logger
	}
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func(t *testing.T, fields *fields)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "close it",
			setup: func(t *testing.T, fields *fields) {
				executor := NewMockExecutor(t)
				executor.EXPECT().Close().Return(nil)
				fields.executor = executor
			},
			wantErr: assert.NoError,
		},
		{
			name: "close it with timeout",
			args: args{
				timeout: 20 * time.Millisecond,
			},
			setup: func(t *testing.T, fields *fields) {
				executor := NewMockExecutor(t)
				executor.EXPECT().Close().Return(nil)
				fields.executor = executor
			},
			wantErr: assert.NoError,
		},
		{
			name: "close it with timeout fires",
			fields: fields{
				runningRequests: []*requestTransaction{
					{},
				},
			},
			args: args{
				timeout: 20 * time.Millisecond,
			},
			setup: func(t *testing.T, fields *fields) {
				executor := NewMockExecutor(t)
				executor.EXPECT().Close().Return(nil)
				fields.executor = executor
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			r := &requestTransactionManager{
				runningRequests:                     tt.fields.runningRequests,
				numberOfConcurrentRequests:          tt.fields.numberOfConcurrentRequests,
				currentTransactionId:                tt.fields.currentTransactionId,
				workLog:                             tt.fields.workLog,
				executor:                            tt.fields.executor,
				traceTransactionManagerTransactions: tt.fields.traceTransactionManagerTransactions,
				log:                                 tt.fields.log,
			}
			tt.wantErr(t, r.CloseGraceful(tt.args.timeout), fmt.Sprintf("CloseGraceful(%v)", tt.args.timeout))
		})
	}
}

func Test_requestTransactionManager_String(t *testing.T) {
	type fields struct {
		runningRequests                     []*requestTransaction
		numberOfConcurrentRequests          int
		currentTransactionId                int32
		workLog                             list.List
		executor                            pool.Executor
		traceTransactionManagerTransactions bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			fields: fields{
				runningRequests: []*requestTransaction{
					{
						transactionId: 2,
					},
				},
				numberOfConcurrentRequests: 3,
				currentTransactionId:       4,
				workLog: func() list.List {
					v := list.List{}
					v.PushBack(nil)
					return v
				}(),
				executor:                            pool.NewFixedSizeExecutor(1, 1),
				traceTransactionManagerTransactions: true,
			},
			want: `
╔═requestTransactionManager═══════════════════════════════════════════════════════════════════════════════════════════╗
║╔═runningRequests/value/requestTransaction╗╔═numberOfConcurrentRequests╗╔═currentTransactionId╗                      ║
║║      ╔═transactionId╗╔═completed╗       ║║   0x0000000000000003 3    ║║    0x00000004 4     ║                      ║
║║      ║ 0x00000002 2 ║║ b0 false ║       ║╚═══════════════════════════╝╚═════════════════════╝                      ║
║║      ╚══════════════╝╚══════════╝       ║                                                                          ║
║╚═════════════════════════════════════════╝                                                                          ║
║╔═executor/executor══════════════════════════════════════════════════════════════════════════════════════╗╔═shutdown╗║
║║╔═running╗╔═shutdown╗                                                                                   ║║b0 false ║║
║║║b0 false║║b0 false ║                                                                                   ║╚═════════╝║
║║╚════════╝╚═════════╝                                                                                   ║           ║
║║╔═worker/value/worker══════════════════════════════════════════════════════════════════════════════════╗║           ║
║║║╔═id═════════════════╗╔═lastReceived════════════════╗╔═running╗╔═shutdown╗╔═interrupted╗╔═interrupter╗║║           ║
║║║║0x0000000000000000 0║║0001-01-01 00:00:00 +0000 UTC║║b0 false║║b0 false ║║  b0 false  ║║0 element(s)║║║           ║
║║║╚════════════════════╝╚═════════════════════════════╝╚════════╝╚═════════╝╚════════════╝╚════════════╝║║           ║
║║╚══════════════════════════════════════════════════════════════════════════════════════════════════════╝║           ║
║║╔═queueDepth═════════╗╔═workItems══╗╔═traceWorkers╗                                                     ║           ║
║║║0x0000000000000001 1║║0 element(s)║║  b0 false   ║                                                     ║           ║
║║╚════════════════════╝╚════════════╝╚═════════════╝                                                     ║           ║
║╚════════════════════════════════════════════════════════════════════════════════════════════════════════╝           ║
║╔═traceTransactionManagerTransactions╗                                                                               ║
║║              b1 true               ║                                                                               ║
║╚════════════════════════════════════╝                                                                               ║
╚═════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &requestTransactionManager{
				runningRequests:                     tt.fields.runningRequests,
				numberOfConcurrentRequests:          tt.fields.numberOfConcurrentRequests,
				currentTransactionId:                tt.fields.currentTransactionId,
				workLog:                             tt.fields.workLog,
				executor:                            tt.fields.executor,
				traceTransactionManagerTransactions: tt.fields.traceTransactionManagerTransactions,
				log:                                 testutils.ProduceTestingLogger(t),
			}
			assert.Equalf(t, tt.want, r.String(), "String()")
		})
	}
}
