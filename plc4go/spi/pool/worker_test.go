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
	"github.com/rs/zerolog"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_worker_initialize(t *testing.T) {
	type fields struct {
		id          int
		interrupter chan struct{}
		executor    interface {
			isTraceWorkers() bool
			getWorksItems() chan workItem
			getWorkerWaitGroup() *sync.WaitGroup
		}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "do it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &worker{
				id:          tt.fields.id,
				interrupter: tt.fields.interrupter,
				executor:    tt.fields.executor,
				log:         produceTestingLogger(t),
			}
			w.initialize()
		})
	}
}

func Test_worker_start(t *testing.T) {
	type fields struct {
		id       int
		executor interface {
			isTraceWorkers() bool
			getWorksItems() chan workItem
			getWorkerWaitGroup() *sync.WaitGroup
		}
		lastReceived atomic.Value
		interrupter  chan struct{}
	}
	tests := []struct {
		name        string
		fields      fields
		manipulator func(t *testing.T, worker *worker)
	}{
		{
			name: "start it",
			fields: fields{
				executor: func() *executor {
					e := &executor{
						workItems:    make(chan workItem, 1),
						traceWorkers: true,
					}
					e.workItems <- workItem{
						workItemId: 0,
						runnable: func() {
							// No-op
						},
						completionFuture: &future{},
					}
					return e
				}(),
			},
		},
		{
			name: "start started",
			fields: fields{
				executor: func() *executor {
					e := &executor{
						workItems:    make(chan workItem, 1),
						traceWorkers: true,
					}
					e.workItems <- workItem{
						workItemId: 0,
						runnable: func() {
							// No-op
						},
						completionFuture: &future{},
					}
					return e
				}(),
			},
			manipulator: func(t *testing.T, worker *worker) {
				worker.running.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &worker{
				id:           tt.fields.id,
				executor:     tt.fields.executor,
				lastReceived: tt.fields.lastReceived,
				interrupter:  tt.fields.interrupter,
				log:          produceTestingLogger(t),
			}
			if tt.manipulator != nil {
				tt.manipulator(t, w)
			}
			w.start()
			t.Cleanup(func() {
				w.stop(false)
			})
		})
	}
}

func Test_worker_stop(t *testing.T) {
	type fields struct {
		id       int
		executor interface {
			isTraceWorkers() bool
			getWorksItems() chan workItem
			getWorkerWaitGroup() *sync.WaitGroup
		}
		lastReceived atomic.Value
		interrupter  chan struct{}
		log          zerolog.Logger
	}
	tests := []struct {
		name        string
		fields      fields
		manipulator func(t *testing.T, worker *worker)
	}{
		{
			name: "stop it",
		},
		{
			name: "stop started",
			fields: fields{
				interrupter: make(chan struct{}),
				executor: func() *executor {
					e := &executor{
						workItems:    make(chan workItem),
						traceWorkers: true,
					}
					go func() {
						e.workItems <- workItem{
							workItemId: 0,
							runnable: func() {
								// No-op
							},
							completionFuture: &future{},
						}
					}()
					return e
				}(),
			},
			manipulator: func(t *testing.T, worker *worker) {
				worker.running.Store(true)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &worker{
				id:           tt.fields.id,
				executor:     tt.fields.executor,
				lastReceived: tt.fields.lastReceived,
				interrupter:  tt.fields.interrupter,
				log:          tt.fields.log,
			}
			if tt.manipulator != nil {
				tt.manipulator(t, w)
			}
			w.stop(true)
		})
	}
}

func Test_worker_work(t *testing.T) {
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
				assert.True(t, w.running.Load(), "should be running")
			},
			manipulator: func(w *worker) {
				w.shutdown.Store(true)
				w.interrupter <- struct{}{}
			},
			timeBeforeSecondValidation: 150 * time.Millisecond,
			secondValidation: func(t *testing.T, w *worker) {
				assert.False(t, w.running.Load(), "should not be running")
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
				assert.True(t, w.running.Load(), "should be running")
			},
			manipulator: func(w *worker) {
				w.shutdown.Store(true)
			},
			timeBeforeSecondValidation: 150 * time.Millisecond,
			secondValidation: func(t *testing.T, w *worker) {
				assert.False(t, w.running.Load(), "should not be running")
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
				assert.True(t, w.running.Load(), "should be running")
			},
			manipulator: func(w *worker) {
				w.shutdown.Store(true)
				w.interrupter <- struct{}{}
			},
			timeBeforeSecondValidation: 150 * time.Millisecond,
			secondValidation: func(t *testing.T, w *worker) {
				assert.False(t, w.running.Load(), "should not be running")
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
				log:         produceTestingLogger(t),
			}
			w.executor.getWorkerWaitGroup().Add(1)
			w.running.Store(true)
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

			close(w.interrupter)
			time.Sleep(50 * time.Millisecond) // TODO: replace with worker stop... (which in turn essures the worker is not running anymore)
		})
	}
}

func Test_worker_String(t *testing.T) {
	type fields struct {
		id int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			want: `
╔═worker═══════════════════════════════════════════════════════════════════════════════════════════════╗
║╔═id═════════════════╗╔═lastReceived════════════════╗╔═running╗╔═shutdown╗╔═interrupted╗╔═interrupter╗║
║║0x0000000000000000 0║║0001-01-01 00:00:00 +0000 UTC║║b0 false║║b0 false ║║  b0 false  ║║0 element(s)║║
║╚════════════════════╝╚═════════════════════════════╝╚════════╝╚═════════╝╚════════════╝╚════════════╝║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &worker{
				id:  tt.fields.id,
				log: produceTestingLogger(t),
			}
			w.lastReceived.Store(time.Time{})
			assert.Equalf(t, tt.want, w.String(), "String()")
		})
	}
}
