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

// pacer enforces a minimum silence gap before writes on a serial line,
// measured from the later of the last write end and the last observed
// read activity (RTU t3.5-style semantics). A nil pacer or a zero delay
// disables pacing entirely. A pacer only tracks timestamps — it does not
// serialize concurrent writers; callers must hold their own write lock
// around waitTurn+write for the gap guarantee to hold.
type pacer struct {
	delay time.Duration
	mu    sync.Mutex
	last  time.Time
}

func newPacer(delay time.Duration) *pacer {
	return &pacer{delay: delay}
}

// noteActivity records line activity (a completed write or received data).
func (p *pacer) noteActivity() {
	if p == nil || p.delay <= 0 {
		return
	}
	p.mu.Lock()
	if now := time.Now(); now.After(p.last) {
		p.last = now
	}
	p.mu.Unlock()
}

// waitTurn sleeps until the configured gap since the last activity has
// elapsed. A non-zero deadline bounds the wait: if it would expire before
// the gap does, waitTurn returns os.ErrDeadlineExceeded without sleeping
// the gap out.
func (p *pacer) waitTurn(deadline time.Time) error {
	if p == nil || p.delay <= 0 {
		return nil
	}
	for {
		p.mu.Lock()
		ready := p.last.Add(p.delay)
		p.mu.Unlock()
		wait := time.Until(ready)
		if wait <= 0 {
			return nil
		}
		if !deadline.IsZero() && ready.After(deadline) {
			return os.ErrDeadlineExceeded
		}
		time.Sleep(wait)
	}
}
