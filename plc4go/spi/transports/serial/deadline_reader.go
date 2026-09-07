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

package serial

import (
	"sync/atomic"
	"time"

	"github.com/apache/plc4x/plc4go/spi/transports/serial/serialport"
)

// deadlineReader arms the port's read deadline before every Read. An
// explicit deadline set via setExplicitDeadline (the context deadline of
// the current transport operation) always wins; an explicit deadline
// observed as already expired is honored for exactly one read — so the
// bounded caller gets its timeout error — and then auto-clears. Without an
// active explicit deadline the configured fallback applies (0 = blocking),
// which keeps bare buffered-layer polls bounded instead of hanging forever
// on a silent line.
type deadlineReader struct {
	port     serialport.Port
	fallback time.Duration
	activity *pacer       // notified on successful reads; may be nil
	explicit atomic.Value // time.Time; zero value = none
}

func newDeadlineReader(port serialport.Port, fallback time.Duration, activity *pacer) *deadlineReader {
	r := &deadlineReader{port: port, fallback: fallback, activity: activity}
	r.explicit.Store(time.Time{})
	return r
}

// setExplicitDeadline sets (or, with a zero time, clears) the explicit
// read deadline. Safe for concurrent use with Read.
func (r *deadlineReader) setExplicitDeadline(t time.Time) {
	r.explicit.Store(t)
}

func (r *deadlineReader) Read(p []byte) (int, error) {
	deadline, _ := r.explicit.Load().(time.Time)
	switch {
	case !deadline.IsZero():
		if err := r.port.SetReadDeadline(deadline); err != nil {
			return 0, err
		}
		if !time.Now().Before(deadline) {
			// Expired: honored for this read, cleared for the next. CAS so
			// a concurrent setExplicitDeadline is never overwritten. The old
			// value MUST be the exact Load()ed time.Time (never reconstructed
			// or derived): CompareAndSwap uses interface ==, which includes
			// time.Time's monotonic-clock reading.
			r.explicit.CompareAndSwap(deadline, time.Time{})
		}
	case r.fallback > 0:
		if err := r.port.SetReadDeadline(time.Now().Add(r.fallback)); err != nil {
			return 0, err
		}
	default:
		if err := r.port.SetReadDeadline(time.Time{}); err != nil {
			return 0, err
		}
	}
	n, err := r.port.Read(p)
	if n > 0 {
		r.activity.noteActivity()
	}
	return n, err
}
