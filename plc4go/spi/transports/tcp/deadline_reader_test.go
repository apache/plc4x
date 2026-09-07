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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recordingConn implements deadlineConn and records SetReadDeadline calls.
type recordingConn struct {
	readDeadlines []time.Time
}

func (r *recordingConn) Read(p []byte) (int, error) {
	return len(p), nil
}

func (r *recordingConn) SetReadDeadline(t time.Time) error {
	r.readDeadlines = append(r.readDeadlines, t)
	return nil
}

func TestDeadlineReader_FallbackArmsRelativeDeadline(t *testing.T) {
	conn := &recordingConn{}
	reader := newDeadlineReader(conn, 200*time.Millisecond)

	before := time.Now()
	_, err := reader.Read(make([]byte, 1))
	require.NoError(t, err)

	require.Len(t, conn.readDeadlines, 1)
	armed := conn.readDeadlines[0]
	assert.False(t, armed.Before(before.Add(200*time.Millisecond)), "deadline must be >= now+fallback")
	assert.True(t, armed.Before(before.Add(2*time.Second)), "deadline must be near now+fallback")
}

func TestDeadlineReader_ZeroFallbackClearsDeadline(t *testing.T) {
	conn := &recordingConn{}
	reader := newDeadlineReader(conn, 0)

	_, err := reader.Read(make([]byte, 1))
	require.NoError(t, err)

	require.Len(t, conn.readDeadlines, 1)
	assert.True(t, conn.readDeadlines[0].IsZero(), "no fallback means the deadline is cleared")
}

func TestDeadlineReader_ExplicitDeadlineWins(t *testing.T) {
	conn := &recordingConn{}
	reader := newDeadlineReader(conn, time.Hour)

	explicit := time.Now().Add(30 * time.Millisecond)
	reader.setExplicitDeadline(explicit)
	_, err := reader.Read(make([]byte, 1))
	require.NoError(t, err)

	require.Len(t, conn.readDeadlines, 1)
	assert.Equal(t, explicit, conn.readDeadlines[0])
}

func TestDeadlineReader_ExpiredExplicitHonoredOnceThenFallback(t *testing.T) {
	conn := &recordingConn{}
	reader := newDeadlineReader(conn, time.Hour)

	expired := time.Now().Add(-time.Second)
	reader.setExplicitDeadline(expired)

	// First read still sees the expired explicit deadline (the ctx-bounded
	// caller must get its timeout), ...
	_, _ = reader.Read(make([]byte, 1))
	require.Len(t, conn.readDeadlines, 1)
	assert.Equal(t, expired, conn.readDeadlines[0])

	// ... the next read falls back.
	_, _ = reader.Read(make([]byte, 1))
	require.Len(t, conn.readDeadlines, 2)
	assert.False(t, conn.readDeadlines[1].Equal(expired), "expired explicit deadline must auto-clear")
	assert.False(t, conn.readDeadlines[1].IsZero(), "fallback must be armed")
}

func TestDeadlineReader_ClearingExplicitRestoresFallback(t *testing.T) {
	conn := &recordingConn{}
	reader := newDeadlineReader(conn, time.Minute)

	reader.setExplicitDeadline(time.Now().Add(time.Hour))
	_, _ = reader.Read(make([]byte, 1))
	reader.setExplicitDeadline(time.Time{})
	_, _ = reader.Read(make([]byte, 1))

	require.Len(t, conn.readDeadlines, 2)
	assert.False(t, conn.readDeadlines[1].Equal(conn.readDeadlines[0]), "cleared explicit must not persist")
}
