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

package tracer

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTracer(t *testing.T) {
	type args struct {
		connectionId string
	}
	tests := []struct {
		name string
		args args
		want Tracer
	}{
		{
			name: "create it",
			want: &tracer{
				traceEntries: []TraceEntry{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewTracer(tt.args.connectionId), "NewTracer(%v)", tt.args.connectionId)
		})
	}
}

func Test_tracer_AddTrace(t1 *testing.T) {
	type fields struct {
		connectionId string
		traceEntries []TraceEntry
	}
	type args struct {
		operation string
		message   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Add a trace",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tracer{
				connectionId: tt.fields.connectionId,
				traceEntries: tt.fields.traceEntries,
			}
			t.AddTrace(tt.args.operation, tt.args.message)
		})
	}
}

func Test_tracer_AddTransactionalStartTrace(t1 *testing.T) {
	type fields struct {
		connectionId string
		traceEntries []TraceEntry
	}
	type args struct {
		operation string
		message   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "start a trace",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tracer{
				connectionId: tt.fields.connectionId,
				traceEntries: tt.fields.traceEntries,
			}
			id := t.AddTransactionalStartTrace(tt.args.operation, tt.args.message)
			assert.NotNil(t1, id, "AddTransactionalStartTrace(%v, %v)", tt.args.operation, tt.args.message)
		})
	}
}

func Test_tracer_AddTransactionalTrace(t1 *testing.T) {
	type fields struct {
		connectionId string
		traceEntries []TraceEntry
	}
	type args struct {
		transactionId string
		operation     string
		message       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "add more trace",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tracer{
				connectionId: tt.fields.connectionId,
				traceEntries: tt.fields.traceEntries,
			}
			t.AddTransactionalTrace(tt.args.transactionId, tt.args.operation, tt.args.message)
		})
	}
}

func Test_tracer_FilterTraces(t1 *testing.T) {
	type fields struct {
		connectionId string
		traceEntries []TraceEntry
	}
	type args struct {
		traces              []TraceEntry
		connectionIdFilter  string
		transactionIdFilter string
		operationFilter     string
		messageFilter       string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []TraceEntry
	}{
		{
			name: "filter",
			args: args{
				traces: []TraceEntry{
					{
						Timestamp:     time.Time{},
						ConnectionId:  "_",
						TransactionId: "2",
						Operation:     "3",
						Message:       "4",
					},
					{
						Timestamp:     time.Time{},
						ConnectionId:  "1",
						TransactionId: "_",
						Operation:     "3",
						Message:       "4",
					},
					{
						Timestamp:     time.Time{},
						ConnectionId:  "1",
						TransactionId: "2",
						Operation:     "_",
						Message:       "4",
					},
					{
						Timestamp:     time.Time{},
						ConnectionId:  "1",
						TransactionId: "2",
						Operation:     "3",
						Message:       "_",
					},
					{
						Timestamp:     time.Time{},
						ConnectionId:  "1",
						TransactionId: "2",
						Operation:     "3",
						Message:       "4",
					},
				},
				connectionIdFilter:  "1",
				transactionIdFilter: "2",
				operationFilter:     "3",
				messageFilter:       "4",
			},
			want: []TraceEntry{
				{
					Timestamp:     time.Time{},
					ConnectionId:  "1",
					TransactionId: "2",
					Operation:     "3",
					Message:       "4",
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tracer{
				connectionId: tt.fields.connectionId,
				traceEntries: tt.fields.traceEntries,
			}
			assert.Equalf(t1, tt.want, t.FilterTraces(tt.args.traces, tt.args.connectionIdFilter, tt.args.transactionIdFilter, tt.args.operationFilter, tt.args.messageFilter), "FilterTraces(%v, %v, %v, %v, %v)", tt.args.traces, tt.args.connectionIdFilter, tt.args.transactionIdFilter, tt.args.operationFilter, tt.args.messageFilter)
		})
	}
}

func Test_tracer_GetConnectionId(t1 *testing.T) {
	type fields struct {
		connectionId string
		traceEntries []TraceEntry
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tracer{
				connectionId: tt.fields.connectionId,
				traceEntries: tt.fields.traceEntries,
			}
			assert.Equalf(t1, tt.want, t.GetConnectionId(), "GetConnectionId()")
		})
	}
}

func Test_tracer_GetTraces(t1 *testing.T) {
	type fields struct {
		connectionId string
		traceEntries []TraceEntry
	}
	tests := []struct {
		name   string
		fields fields
		want   []TraceEntry
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tracer{
				connectionId: tt.fields.connectionId,
				traceEntries: tt.fields.traceEntries,
			}
			assert.Equalf(t1, tt.want, t.GetTraces(), "GetTraces()")
		})
	}
}

func Test_tracer_ResetTraces(t1 *testing.T) {
	type fields struct {
		connectionId string
		traceEntries []TraceEntry
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "reset it",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tracer{
				connectionId: tt.fields.connectionId,
				traceEntries: tt.fields.traceEntries,
			}
			t.ResetTraces()
		})
	}
}

func Test_tracer_SetConnectionId(t1 *testing.T) {
	type fields struct {
		connectionId string
		traceEntries []TraceEntry
	}
	type args struct {
		connectionId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set it",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &tracer{
				connectionId: tt.fields.connectionId,
				traceEntries: tt.fields.traceEntries,
			}
			t.SetConnectionId(tt.args.connectionId)
		})
	}
}
