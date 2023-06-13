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
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestNewFixedSizeExecutor(t *testing.T) {
	type args struct {
		numberOfWorkers int
		queueDepth      int
		options         []options.WithOption
	}
	tests := []struct {
		name              string
		args              args
		setup             func(t *testing.T, args *args)
		executorValidator func(*testing.T, *executor) bool
	}{
		{
			name: "new Executor",
			args: args{
				numberOfWorkers: 13,
				queueDepth:      14,
				options:         []options.WithOption{WithExecutorOptionTracerWorkers(true)},
			},
			setup: func(t *testing.T, args *args) {
				args.options = append(args.options, options.WithCustomLogger(produceTestLogger(t)))
			},
			executorValidator: func(t *testing.T, e *executor) bool {
				return !e.running && !e.shutdown && len(e.worker) == 13 && cap(e.workItems) == 14
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.args)
			}
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
		setup             func(*testing.T, *args)
		manipulator       func(*testing.T, *dynamicExecutor)
		executorValidator func(*testing.T, *dynamicExecutor) bool
	}{
		{
			name: "new Executor",
			args: args{
				numberOfWorkers: 13,
				queueDepth:      14,
				options:         []options.WithOption{WithExecutorOptionTracerWorkers(true)},
			},
			setup: func(t *testing.T, args *args) {
				args.options = append(args.options, options.WithCustomLogger(produceTestLogger(t)))
			},
			executorValidator: func(t *testing.T, e *dynamicExecutor) bool {
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
			setup: func(t *testing.T, args *args) {
				args.options = append(args.options, options.WithCustomLogger(produceTestLogger(t)))
			},
			manipulator: func(t *testing.T, e *dynamicExecutor) {
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
			executorValidator: func(t *testing.T, e *dynamicExecutor) bool {
				time.Sleep(500 * time.Millisecond)
				assert.False(t, e.running)
				assert.False(t, e.shutdown)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.args)
			}
			dynamicSizedExecutor := NewDynamicExecutor(tt.args.numberOfWorkers, tt.args.queueDepth, tt.args.options...)
			defer dynamicSizedExecutor.Stop()
			if tt.manipulator != nil {
				tt.manipulator(t, dynamicSizedExecutor.(*dynamicExecutor))
			}
			assert.True(t, tt.executorValidator(t, dynamicSizedExecutor.(*dynamicExecutor)), "NewFixedSizeExecutor(%v, %v, %v)", tt.args.numberOfWorkers, tt.args.queueDepth, tt.args.options)
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

// from: https://github.com/golang/go/issues/36532#issuecomment-575535452
func testContext(t *testing.T) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return ctx
}

// note: we can't use testutils here due to import cycle
func produceTestLogger(t *testing.T) zerolog.Logger {
	return zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(t),
		func(w *zerolog.ConsoleWriter) {
			// TODO: this is really an issue with go-junit-report not sanitizing output before dumping into xml...
			onJenkins := os.Getenv("JENKINS_URL") != ""
			onGithubAction := os.Getenv("GITHUB_ACTIONS") != ""
			onCI := os.Getenv("CI") != ""
			if onJenkins || onGithubAction || onCI {
				w.NoColor = true
			}
		}))
}
