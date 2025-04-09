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
	"fmt"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type stopWarnOptions struct {
	processId              string
	processInfo            string
	interval               time.Duration
	extraSkipOffset        int
	includeGoroutinesStack bool
}

// StopWarn gives out warning every interval (default 5 seconds) when a function doesn't terminate. Usage: `defer StopWarn(log)()`
func StopWarn(localLog zerolog.Logger, opts ...func(*stopWarnOptions)) func() {
	o := &stopWarnOptions{
		processInfo: "",
		interval:    5 * time.Second,
	}
	for _, opt := range opts {
		opt(o)
	}
	if o.processInfo == "" {
		_, file, line, ok := runtime.Caller(1 + o.extraSkipOffset)
		if ok {
			o.processInfo = fmt.Sprintf("%s:%d", file, line)
		}
	}
	localLog = localLog.With().Str("processInfo", o.processInfo).Dur("interval", o.interval).Logger()
	ticker := time.NewTicker(o.interval)
	wg := new(sync.WaitGroup)
	done := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		localLog.Trace().Msgf("start checking")
		startTime := time.Now()
		for {
			localLog.Trace().Msgf("check cycle")
			select {
			case <-done:
				ticker.Stop()
				return
			case warnTime := <-ticker.C:
				processId := ""
				if o.processId != "" {
					processId = o.processId + " "
				}
				stackInfo := new(strings.Builder)
				if o.includeGoroutinesStack {
					goroutines := pprof.Lookup("goroutine")
					if goroutines != nil {
						if err := goroutines.WriteTo(stackInfo, 2); err != nil {
							localLog.Warn().Err(err).Msg("could not write to stack")
							stackInfo.WriteString(err.Error())
						}
					} else {
						log.Warn().Msg("lookup goroutine failed")
					}
				}
				localLog.Warn().
					Time("startTime", startTime).
					Time("warnTime", warnTime).
					TimeDiff("inProgressFor", warnTime, startTime).
					Stringer("stackInfo", stackInfo).
					Msgf("%sstill in progress", processId)
			}
		}
	}()
	start := time.Now()
	return func() {
		localLog.Trace().TimeDiff("check duration", time.Now(), start).Msg("done")
		close(done)
		wg.Wait() // This is to avoid late logs in case when the shutdown is really fast
	}
}

// WithStopWarnProcessId sets a process id which will be prefixed to the message
func WithStopWarnProcessId(processId string) func(*stopWarnOptions) {
	return func(o *stopWarnOptions) {
		o.processId = processId
	}
}

// WithStopWarnProcessInfo set the processInfo
func WithStopWarnProcessInfo(processInfo string) func(*stopWarnOptions) {
	return func(o *stopWarnOptions) {
		o.processInfo = processInfo
	}
}

// WithStopWarnInterval sets the interval at which a warning is logged (default 5 seconds). MUST be greater 0.
func WithStopWarnInterval(interval time.Duration) func(*stopWarnOptions) {
	return func(o *stopWarnOptions) {
		o.interval = interval
	}
}

// WithStopWarnExtraSkipOffset sets an extra offset for the skip. Skip uses 1 so if you want to have 0 use -1 as arg.
func WithStopWarnExtraSkipOffset(offset int) func(*stopWarnOptions) {
	return func(o *stopWarnOptions) {
		o.extraSkipOffset = offset
	}
}

// WithStopWarnIncludeGoroutinesStack is a flag which instructs the warn log to include the list of all current goroutines
func WithStopWarnIncludeGoroutinesStack() func(*stopWarnOptions) {
	return func(o *stopWarnOptions) {
		o.includeGoroutinesStack = true
	}
}
