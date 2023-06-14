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
	"testing"
	"time"

	"github.com/apache/plc4x/plc4go/spi/pool"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_requestTransaction_EndRequest(t1 *testing.T) {
	type fields struct {
		parent           *requestTransactionManager
		transactionId    int32
		operation        pool.Runnable
		completionFuture pool.CompletionFuture
		transactionLog   zerolog.Logger
		completed        bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "just end it",
			fields: fields{
				parent: &requestTransactionManager{},
			},
			wantErr: true,
		},
		{
			name: "end it completed",
			fields: fields{
				parent:    &requestTransactionManager{},
				completed: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &requestTransaction{
				parent:           tt.fields.parent,
				transactionId:    tt.fields.transactionId,
				operation:        tt.fields.operation,
				completionFuture: tt.fields.completionFuture,
				transactionLog:   tt.fields.transactionLog,
				completed:        tt.fields.completed,
			}
			if err := t.EndRequest(); (err != nil) != tt.wantErr {
				t1.Errorf("EndRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_requestTransaction_FailRequest(t1 *testing.T) {
	type fields struct {
		parent           *requestTransactionManager
		transactionId    int32
		operation        pool.Runnable
		completionFuture pool.CompletionFuture
		transactionLog   zerolog.Logger
		completed        bool
	}
	type args struct {
		err error
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		mockSetup func(t *testing.T, fields *fields, args *args)
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "just fail it",
			fields: fields{
				parent: &requestTransactionManager{},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				completionFuture := NewMockCompletionFuture(t)
				expect := completionFuture.EXPECT()
				expect.Cancel(true, nil).Return()
				fields.completionFuture = completionFuture
			},
			wantErr: assert.Error,
		},
		{
			name: "just fail it (completed)",
			args: args{
				err: errors.New("nope"),
			},
			fields: fields{
				parent:    &requestTransactionManager{},
				completed: true,
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			r := &requestTransaction{
				parent:           tt.fields.parent,
				transactionId:    tt.fields.transactionId,
				operation:        tt.fields.operation,
				completionFuture: tt.fields.completionFuture,
				transactionLog:   tt.fields.transactionLog,
				completed:        tt.fields.completed,
			}
			tt.wantErr(t, r.FailRequest(tt.args.err), "FailRequest() error = %v", tt.args.err)
		})
	}
}

func Test_requestTransaction_String(t1 *testing.T) {
	type fields struct {
		parent           *requestTransactionManager
		transactionId    int32
		operation        pool.Runnable
		completionFuture pool.CompletionFuture
		transactionLog   zerolog.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "give a string",
			want: `
╔═requestTransaction═════════╗
║╔═transactionId╗╔═completed╗║
║║ 0x00000000 0 ║║ b0 false ║║
║╚══════════════╝╚══════════╝║
╚════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &requestTransaction{
				parent:           tt.fields.parent,
				transactionId:    tt.fields.transactionId,
				operation:        tt.fields.operation,
				completionFuture: tt.fields.completionFuture,
				transactionLog:   tt.fields.transactionLog,
			}
			if got := t.String(); got != tt.want {
				t1.Errorf("String() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_requestTransaction_Submit(t1 *testing.T) {
	type fields struct {
		parent           *requestTransactionManager
		transactionId    int32
		operation        pool.Runnable
		completionFuture pool.CompletionFuture
		transactionLog   zerolog.Logger
		completed        bool
	}
	type args struct {
		operation RequestTransactionRunnable
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "submit something",
			fields: fields{
				parent: &requestTransactionManager{},
			},
			args: args{
				operation: func(_ RequestTransaction) {
					// NOOP
				},
			},
		},
		{
			name: "submit something again",
			fields: fields{
				parent: &requestTransactionManager{},
				operation: func() {
					// NOOP
				},
			},
			args: args{
				operation: func(_ RequestTransaction) {
					// NOOP
				},
			},
		},
		{
			name: "submit completed",
			fields: fields{
				parent: &requestTransactionManager{},
				operation: func() {
					// NOOP
				},
				completed: true,
			},
			args: args{
				operation: func(_ RequestTransaction) {
					// NOOP
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &requestTransaction{
				parent:           tt.fields.parent,
				transactionId:    tt.fields.transactionId,
				operation:        tt.fields.operation,
				completionFuture: tt.fields.completionFuture,
				transactionLog:   tt.fields.transactionLog,
				completed:        tt.fields.completed,
			}
			t.Submit(tt.args.operation)
			t.operation()
		})
	}
}

func Test_requestTransaction_AwaitCompletion(t1 *testing.T) {
	type fields struct {
		parent           *requestTransactionManager
		transactionId    int32
		operation        pool.Runnable
		completionFuture pool.CompletionFuture
		transactionLog   zerolog.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		mockSetup   func(t *testing.T, fields *fields, args *args)
		manipulator func(t *testing.T, transaction *requestTransaction)
		wantErr     bool
	}{
		{
			name: "just wait",
			fields: fields{
				parent: &requestTransactionManager{
					runningRequests: []*requestTransaction{},
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return ctx
				}(),
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				completionFuture := NewMockCompletionFuture(t)
				expect := completionFuture.EXPECT()
				expect.AwaitCompletion(mock.Anything).Return(nil)
				fields.completionFuture = completionFuture
			},
			manipulator: func(t *testing.T, transaction *requestTransaction) {
				go func() {
					time.Sleep(100 * time.Millisecond)
					r := transaction.parent
					r.workLogMutex.RLock()
					defer r.workLogMutex.RUnlock()
					r.runningRequests = append(r.runningRequests, &requestTransaction{transactionId: 1})
				}()
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t1, &tt.fields, &tt.args)
			}
			t := &requestTransaction{
				parent:           tt.fields.parent,
				transactionId:    tt.fields.transactionId,
				operation:        tt.fields.operation,
				completionFuture: tt.fields.completionFuture,
				transactionLog:   tt.fields.transactionLog,
			}
			if err := t.AwaitCompletion(tt.args.ctx); (err != nil) != tt.wantErr {
				t1.Errorf("AwaitCompletion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
