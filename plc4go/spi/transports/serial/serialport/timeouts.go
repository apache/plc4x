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

package serialport

import "time"

// maxIOSliceMillis caps how long a single blocking I/O attempt may last on
// platforms where blocking calls cannot be interrupted (Windows), so that
// Close and deadline changes are observed within this bound.
const maxIOSliceMillis = 60_000

// ioSliceMillis returns how many milliseconds the next blocking I/O attempt
// may last, given the current time and an absolute deadline (zero = none).
// ok is false when the deadline has already expired; ms is then 0.
// A non-expired deadline always yields ms >= 1 (sub-millisecond remainders
// round up) because a 0 would mean "no timeout" to COMMTIMEOUTS.
func ioSliceMillis(now, deadline time.Time) (ms uint32, ok bool) {
	if deadline.IsZero() {
		return maxIOSliceMillis, true
	}
	remaining := deadline.Sub(now)
	if remaining <= 0 {
		return 0, false
	}
	millis := (remaining + time.Millisecond - 1) / time.Millisecond
	if millis > maxIOSliceMillis {
		return maxIOSliceMillis, true
	}
	return uint32(millis), true
}
