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

package tcp

import (
	"io"
	"sync/atomic"
	"time"
)

// deadlineConn is the subset of net.Conn the deadlineReader needs.
type deadlineConn interface {
	io.Reader
	SetReadDeadline(t time.Time) error
}

// deadlineReader arms the connection's read deadline before every Read - the
// TCP analogue of the serial transport's deadlineReader. An explicit deadline
// set via setExplicitDeadline (the context deadline of the current transport
// operation) always wins; an explicit deadline observed as already expired is
// honored for exactly one read - so the bounded caller gets its timeout error
// - and then auto-clears. Without an active explicit deadline the configured
// fallback applies (0 = blocking), which keeps bare buffered-layer polls
// (DefaultCodec drives Receive with a deadline-less long-lived context)
// bounded instead of hanging forever on a silent connection.
type deadlineReader struct {
	conn     deadlineConn
	fallback time.Duration
	explicit atomic.Value // time.Time; zero value = none
}

func newDeadlineReader(conn deadlineConn, fallback time.Duration) *deadlineReader {
	r := &deadlineReader{conn: conn, fallback: fallback}
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
		if err := r.conn.SetReadDeadline(deadline); err != nil {
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
		if err := r.conn.SetReadDeadline(time.Now().Add(r.fallback)); err != nil {
			return 0, err
		}
	default:
		if err := r.conn.SetReadDeadline(time.Time{}); err != nil {
			return 0, err
		}
	}
	return r.conn.Read(p)
}
