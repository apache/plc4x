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

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIOSliceMillis(t *testing.T) {
	now := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		deadline time.Time
		wantMs   uint32
		wantOK   bool
	}{
		{"no deadline uses the max slice", time.Time{}, maxIOSliceMillis, true},
		{"future deadline: exact milliseconds", now.Add(1500 * time.Millisecond), 1500, true},
		{"sub-millisecond remainder rounds up, never 0", now.Add(100 * time.Microsecond), 1, true},
		{"far deadline capped to max slice", now.Add(time.Hour), maxIOSliceMillis, true},
		{"deadline exactly now: expired", now, 0, false},
		{"past deadline: expired", now.Add(-time.Second), 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms, ok := ioSliceMillis(now, tt.deadline)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantMs, ms)
		})
	}
}
