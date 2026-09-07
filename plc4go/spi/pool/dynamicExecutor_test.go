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
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_newDynamicExecutor(t *testing.T) {
	type args struct {
		queueDepth         int
		maxNumberOfWorkers int
		log                zerolog.Logger
	}
	tests := []struct {
		name       string
		args       args
		wantAssert func(*testing.T, *dynamicExecutor) bool
	}{
		{
			name: "just create it",
			wantAssert: func(t *testing.T, d *dynamicExecutor) bool {
				return assert.NotNil(t, d.executor)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newDynamicExecutor(tt.args.queueDepth, tt.args.maxNumberOfWorkers, tt.args.log)
			assert.Truef(t, tt.wantAssert(t, got), "newDynamicExecutor(%v, %v, %v)", tt.args.queueDepth, tt.args.maxNumberOfWorkers, tt.args.log)
		})
	}
}

func Test_dynamicExecutor_Start(t *testing.T) {
	type fields struct {
		executor           *executor
		maxNumberOfWorkers int
	}
	tests := []struct {
		name       string
		fields     fields
		setup      func(t *testing.T, fields *fields)
		startTwice bool
	}{
		{
			name: "just start",
			fields: fields{
				executor: &executor{
					workItems:    make(chan workItem, 1),
					worker:       make([]*worker, 0),
					traceWorkers: true,
					ctxCancel:    func() {},
					ctx:          t.Context(),
				},
				maxNumberOfWorkers: 100,
			},
			setup: func(t *testing.T, fields *fields) {
				fields.executor.log = produceTestingLogger(t)
				fields.executor.workItems <- workItem{1, func(context.Context) {}, &future{}}
			},
		},
		{
			name: "start twice",
			fields: fields{
				executor: &executor{
					workItems:    make(chan workItem, 1),
					worker:       make([]*worker, 0),
					traceWorkers: true,
					ctxCancel:    func() {},
					ctx:          t.Context(),
				},
				maxNumberOfWorkers: 100,
			},
			setup: func(t *testing.T, fields *fields) {
				fields.executor.log = produceTestingLogger(t)
				fields.executor.workItems <- workItem{1, func(context.Context) {}, &future{}}
			},
			startTwice: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			e := &dynamicExecutor{
				executor:           tt.fields.executor,
				maxNumberOfWorkers: tt.fields.maxNumberOfWorkers,
			}
			e.Start()
			if tt.startTwice {
				e.Start()
			}
			// Let it work a bit
			time.Sleep(20 * time.Millisecond)
			t.Log("done with test")
			t.Cleanup(e.Stop)
		})
	}
}

func Test_dynamicExecutor_Stop(t *testing.T) {
	type fields struct {
		executor           *executor
		maxNumberOfWorkers int
		interrupter        chan struct{}
	}
	tests := []struct {
		name      string
		fields    fields
		setup     func(t *testing.T, fields *fields)
		startIt   bool
		stopTwice bool
	}{
		{
			name: "just stop",
			fields: fields{
				executor: &executor{
					workItems:    make(chan workItem, 1),
					worker:       make([]*worker, 0),
					traceWorkers: true,
				},
				maxNumberOfWorkers: 100,
			},
			setup: func(t *testing.T, fields *fields) {
				fields.executor.log = produceTestingLogger(t)
				fields.executor.workItems <- workItem{1, func(context.Context) {}, &future{}}
			},
		},
		{
			name: "stop started",
			fields: fields{
				executor: &executor{
					workItems:    make(chan workItem, 1),
					worker:       make([]*worker, 0),
					traceWorkers: true,
				},
				maxNumberOfWorkers: 100,
			},
			setup: func(t *testing.T, fields *fields) {
				fields.executor.log = produceTestingLogger(t)
				fields.executor.workItems <- workItem{1, func(context.Context) {}, &future{}}
			},
		},
		{
			name: "stop twice",
			fields: fields{
				executor: &executor{
					workItems:    make(chan workItem, 1),
					worker:       make([]*worker, 0),
					traceWorkers: true,
				},
				maxNumberOfWorkers: 100,
			},
			setup: func(t *testing.T, fields *fields) {
				fields.executor.log = produceTestingLogger(t)
				fields.executor.workItems <- workItem{1, func(context.Context) {}, &future{}}
			},
			stopTwice: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			e := &dynamicExecutor{
				executor:           tt.fields.executor,
				maxNumberOfWorkers: tt.fields.maxNumberOfWorkers,
				interrupter:        tt.fields.interrupter,
			}
			if tt.startIt {
				e.Start()
			}
			e.Stop()
			if tt.stopTwice {
				e.Stop()
			}
		})
	}
}

func Test_dynamicExecutor_String(t *testing.T) {
	type fields struct {
		executor           *executor
		maxNumberOfWorkers int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			fields: fields{
				executor: &executor{
					worker: []*worker{
						{
							lastReceived: func() atomic.Value {
								value := atomic.Value{}
								value.Store(time.Time{})
								return value
							}(),
						},
					},
				},
				maxNumberOfWorkers: 3,
			},
			want: `
╔═dynamicExecutor═════════════════════════════════════════════════════════════════════════════════════════╗
║╔═executor══════════════════════════════════════════════════════════════════════════════════════════════╗║
║║╔═running╗╔═shutdown╗╔═worker/value/worker════════════════════════════════════════════════════════════╗║║
║║║b0 false║║b0 false ║║╔═lastReceived════════════════╗╔═running╗╔═shutdown╗╔═interrupted╗╔═interrupter╗║║║
║║╚════════╝╚═════════╝║║0001-01-01 00:00:00 +0000 UTC║║b0 false║║b0 false ║║  b0 false  ║║0 element(s)║║║║
║║                     ║╚═════════════════════════════╝╚════════╝╚═════════╝╚════════════╝╚════════════╝║║║
║║                     ╚════════════════════════════════════════════════════════════════════════════════╝║║
║║╔═workerNumber╗╔═workItems══╗╔═traceWorkers╗                                                           ║║
║║║0x00000000 0 ║║0 element(s)║║  b0 false   ║                                                           ║║
║║╚═════════════╝╚════════════╝╚═════════════╝                                                           ║║
║╚═══════════════════════════════════════════════════════════════════════════════════════════════════════╝║
║╔═maxNumberOfWorkers═╗╔═currentNumberOfWorkers╗╔═interrupter╗                                            ║
║║0x0000000000000003 3║║     0x00000000 0      ║║0 element(s)║                                            ║
║╚════════════════════╝╚═══════════════════════╝╚════════════╝                                            ║
╚═════════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &dynamicExecutor{
				executor:           tt.fields.executor,
				maxNumberOfWorkers: tt.fields.maxNumberOfWorkers,
			}
			assert.Equalf(t, tt.want, e.String(), "String()")
		})
	}
}
