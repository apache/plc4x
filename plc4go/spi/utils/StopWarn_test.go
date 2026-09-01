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

package utils

import (
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestStopWarn(t *testing.T) {
	t.Run("couple of warnings", func(t *testing.T) {
		logger := produceTestingLogger(t)
		logHook := new(logHook)
		logger = logger.Hook(logHook)
		logger.Debug().Msg("Debug test message")
		t.Log("Starting func using stopWarn")
		func() {
			defer StopWarn(logger, WithStopWarnInterval(10*time.Millisecond), WithStopWarnProcessId("TestStopWarn"), WithStopWarnIncludeGoroutinesStack())()
			// Wait for the warnings instead of sleeping a fixed time: coarse timer
			// granularity (~15.6ms on Windows) and dropped ticks make a tick count
			// per wall-clock window unpredictable on loaded CI runners.
			deadline := time.Now().Add(10 * time.Second)
			for logHook.countMessage("TestStopWarn still in progress") <= 3 && time.Now().Before(deadline) {
				time.Sleep(5 * time.Millisecond)
			}
			t.Log("wait done")
		}()
		t.Log("Doing assertions")
		foundMessages := 0
		for _, message := range logHook.messages {
			if message == "TestStopWarn still in progress" {
				foundMessages++
			}
		}
		assert.Greaterf(t, foundMessages, 3, "%s should contain at least three warning. Found %d times", logHook.messages, foundMessages)
	})
	t.Run("shutdown in time", func(t *testing.T) {
		logger := produceTestingLogger(t)
		logHook := new(logHook)
		logger = logger.Hook(logHook)
		logger.Debug().Msg("Debug test message")
		t.Log("Starting func using stopWarn")
		func() {
			defer StopWarn(logger, WithStopWarnInterval(100*time.Millisecond), WithStopWarnProcessId("TestStopWarn"))()
			time.Sleep(10 * time.Millisecond)
			t.Log("sleep done")
		}()
		t.Log("Doing assertions")
		foundMessages := 0
		for _, message := range logHook.messages {
			if message == "TestStopWarn still in progress" {
				foundMessages++
			}
		}
		assert.Equalf(t, 0, foundMessages, "%s should contain at least three warning. Found %d times", logHook.messages, foundMessages)
	})
	t.Run("registrar receives ticks and end", func(t *testing.T) {
		logger := produceTestingLogger(t)
		var mu sync.Mutex
		var initial StopWarnTick
		var ticks []StopWarnTick
		ended := false
		register := func(in StopWarnTick) (func(StopWarnTick), func()) {
			mu.Lock()
			initial = in
			mu.Unlock()
			return func(t StopWarnTick) {
					mu.Lock()
					ticks = append(ticks, t)
					mu.Unlock()
				}, func() {
					mu.Lock()
					ended = true
					mu.Unlock()
				}
		}
		func() {
			defer StopWarn(logger,
				WithStopWarnInterval(10*time.Millisecond),
				WithStopWarnProcessId("TestStopWarn"),
				WithStopWarnRegistrar(register),
			)()
			// Condition-based wait: see "couple of warnings" above.
			deadline := time.Now().Add(10 * time.Second)
			for time.Now().Before(deadline) {
				mu.Lock()
				enough := len(ticks) >= 3
				mu.Unlock()
				if enough {
					break
				}
				time.Sleep(5 * time.Millisecond)
			}
		}()
		mu.Lock()
		defer mu.Unlock()
		assert.Equal(t, "TestStopWarn", initial.ProcessID, "registrar should receive initial tick at arm time")
		assert.False(t, initial.StartTime.IsZero(), "initial StartTime should be set")
		assert.Equal(t, initial.StartTime, initial.WarnTime, "initial WarnTime equals StartTime")
		assert.GreaterOrEqual(t, len(ticks), 3, "registrar should observe multiple ticks")
		assert.True(t, ended, "registrar onEnd should fire when stop func returns")
		for _, tk := range ticks {
			assert.Equal(t, "TestStopWarn", tk.ProcessID)
			assert.Equal(t, initial.StartTime, tk.StartTime, "StartTime should be stable across ticks")
			assert.False(t, tk.WarnTime.Before(initial.StartTime), "WarnTime should be at or after StartTime")
		}
	})
	t.Run("registrar onEnd only when nil onTick", func(t *testing.T) {
		logger := produceTestingLogger(t)
		var mu sync.Mutex
		ended := false
		register := func(StopWarnTick) (func(StopWarnTick), func()) {
			return nil, func() { mu.Lock(); ended = true; mu.Unlock() }
		}
		func() {
			defer StopWarn(logger,
				WithStopWarnInterval(10*time.Millisecond),
				WithStopWarnRegistrar(register),
			)()
			time.Sleep(35 * time.Millisecond)
		}()
		mu.Lock()
		defer mu.Unlock()
		assert.True(t, ended, "nil onTick should not prevent onEnd from firing")
	})
}

type logHook struct {
	// Run is invoked concurrently: the StopWarn ticker goroutine logs its
	// warnings while the stop func logs "done" from the test goroutine.
	mu        sync.Mutex
	logEvents []zerolog.Event
	level     []zerolog.Level
	messages  []string
}

func (logHook *logHook) countMessage(message string) int {
	logHook.mu.Lock()
	defer logHook.mu.Unlock()
	count := 0
	for _, m := range logHook.messages {
		if m == message {
			count++
		}
	}
	return count
}

func (logHook *logHook) Run(logEvent *zerolog.Event, level zerolog.Level, message string) {
	logHook.mu.Lock()
	defer logHook.mu.Unlock()
	logHook.logEvents = append(logHook.logEvents, *logEvent)
	logHook.level = append(logHook.level, level)
	logHook.messages = append(logHook.messages, message)
}
