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
	"os"
	"sync"
	"time"
)

// byteRing is a bounded FIFO byte buffer connecting the shared-port reader
// (append side) to one subscription (read side). On overflow the OLDEST
// bytes are dropped — a slow subscriber loses history but can never block
// the shared reader; protocol codecs resynchronize on the resulting gap.
//
// Contract: exactly one goroutine blocks in read at a time (wake-ups are
// coalesced through a one-slot channel).
type byteRing struct {
	mu      sync.Mutex
	buf     []byte
	start   int
	length  int
	dropped uint64
	closed  bool
	err     error
	notify  chan struct{}
}

func newByteRing(capacity int) *byteRing {
	return &byteRing{
		buf:    make([]byte, capacity),
		notify: make(chan struct{}, 1),
	}
}

func (r *byteRing) signal() {
	select {
	case r.notify <- struct{}{}:
	default:
	}
}

// append copies p into the ring, dropping the oldest bytes when capacity
// is exceeded, and returns how many bytes were dropped. Appends after
// close are discarded.
func (r *byteRing) append(p []byte) int {
	r.mu.Lock()
	if r.closed || len(p) == 0 {
		r.mu.Unlock()
		return 0
	}
	dropped := 0
	if len(p) >= len(r.buf) {
		dropped = r.length + len(p) - len(r.buf)
		p = p[len(p)-len(r.buf):]
		r.start, r.length = 0, 0
	} else if over := r.length + len(p) - len(r.buf); over > 0 {
		dropped = over
		r.start = (r.start + over) % len(r.buf)
		r.length -= over
	}
	for i := range p {
		r.buf[(r.start+r.length+i)%len(r.buf)] = p[i]
	}
	r.length += len(p)
	r.dropped += uint64(dropped)
	r.mu.Unlock()
	r.signal()
	return dropped
}

// fail records a fatal error. Buffered data remains readable; once it is
// drained, reads return err.
func (r *byteRing) fail(err error) {
	r.mu.Lock()
	if r.err == nil {
		r.err = err
	}
	r.mu.Unlock()
	r.signal()
}

// close makes all subsequent reads fail with os.ErrClosed (buffered data
// is NOT drained first — matching a closed port's behavior).
func (r *byteRing) close() {
	r.mu.Lock()
	r.closed = true
	r.mu.Unlock()
	r.signal()
}

func (r *byteRing) isClosed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.closed
}

func (r *byteRing) droppedTotal() uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.dropped
}

// read blocks until data is available, the deadline expires
// (os.ErrDeadlineExceeded; immediately if already past), the ring is
// closed (os.ErrClosed), or a recorded failure surfaces after the buffer
// drains. A zero deadline blocks indefinitely.
func (r *byteRing) read(p []byte, deadline time.Time) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	var timeout <-chan time.Time
	if !deadline.IsZero() {
		wait := time.Until(deadline)
		if wait <= 0 {
			r.mu.Lock()
			closed := r.closed
			r.mu.Unlock()
			if closed {
				return 0, os.ErrClosed
			}
			return 0, os.ErrDeadlineExceeded
		}
		timer := time.NewTimer(wait)
		defer timer.Stop()
		timeout = timer.C
	}
	for {
		r.mu.Lock()
		if r.closed {
			r.mu.Unlock()
			return 0, os.ErrClosed
		}
		if r.length > 0 {
			n := min(len(p), r.length)
			for i := 0; i < n; i++ {
				p[i] = r.buf[(r.start+i)%len(r.buf)]
			}
			r.start = (r.start + n) % len(r.buf)
			r.length -= n
			r.mu.Unlock()
			return n, nil
		}
		if r.err != nil {
			err := r.err
			r.mu.Unlock()
			return 0, err
		}
		r.mu.Unlock()
		select {
		case <-r.notify:
		case <-timeout:
			return 0, os.ErrDeadlineExceeded
		}
	}
}
