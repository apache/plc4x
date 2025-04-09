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
			time.Sleep(100 * time.Millisecond)
			t.Log("sleep done")
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
}

type logHook struct {
	logEvents []zerolog.Event
	level     []zerolog.Level
	messages  []string
}

func (logHook *logHook) Run(logEvent *zerolog.Event, level zerolog.Level, message string) {
	logHook.logEvents = append(logHook.logEvents, *logEvent)
	logHook.level = append(logHook.level, level)
	logHook.messages = append(logHook.messages, message)
}
