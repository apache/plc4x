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

	"github.com/apache/plc4x/plc4go/spi/options"

	"io"
	"time"
)

type Runnable func()

type CompletionFuture interface {
	AwaitCompletion(ctx context.Context) error
	Cancel(interrupt bool, err error)
}

type Executor interface {
	io.Closer
	Start()
	Stop()
	Submit(ctx context.Context, workItemId int32, runnable Runnable) CompletionFuture
	IsRunning() bool
}

func NewFixedSizeExecutor(numberOfWorkers, queueDepth int, _options ...options.WithOption) Executor {
	workers := make([]*worker, numberOfWorkers)
	customLogger := options.ExtractCustomLogger(_options...)
	for i := 0; i < numberOfWorkers; i++ {
		w := worker{
			id:  i,
			log: customLogger,
		}
		w.lastReceived.Store(time.Time{})
		workers[i] = &w
	}
	_executor := &executor{
		queueDepth: queueDepth,
		workItems:  make(chan workItem, queueDepth),
		worker:     workers,
		log:        customLogger,
	}
	_executor.traceWorkers, _ = options.ExtractTracerWorkers(_options...)
	for i := 0; i < numberOfWorkers; i++ {
		workers[i].executor = _executor
	}
	return _executor
}

func NewDynamicExecutor(maxNumberOfWorkers, queueDepth int, _options ...options.WithOption) Executor {
	customLogger := options.ExtractCustomLogger(_options...)
	_executor := &dynamicExecutor{
		executor: &executor{
			workItems: make(chan workItem, queueDepth),
			worker:    make([]*worker, 0),
			log:       customLogger,
		},
		maxNumberOfWorkers: maxNumberOfWorkers,
	}
	_executor.traceWorkers, _ = options.ExtractTracerWorkers(_options...)
	// We spawn one initial worker
	w := worker{
		id:          0,
		interrupter: make(chan struct{}, 1),
		executor:    _executor,
		log:         customLogger,
	}
	w.lastReceived.Store(time.Now())
	_executor.worker = append(_executor.worker, &w)
	return _executor
}
