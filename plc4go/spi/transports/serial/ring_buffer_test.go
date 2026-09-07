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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

func TestByteRing_FIFOAndWraparound(t *testing.T) {
	r := newByteRing(8)
	assert.Zero(t, r.append([]byte("abcde")))
	buf := make([]byte, 3)
	n, err := r.read(buf, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, "abc", string(buf[:n]))
	// wraps around the 8-byte capacity
	assert.Zero(t, r.append([]byte("fghi")))
	got := make([]byte, 0, 6)
	for len(got) < 6 {
		n, err = r.read(buf, time.Time{})
		require.NoError(t, err)
		got = append(got, buf[:n]...)
	}
	assert.Equal(t, "defghi", string(got))
}

func TestByteRing_BlockingReadWokenByAppend(t *testing.T) {
	r := newByteRing(8)
	go func() {
		time.Sleep(20 * time.Millisecond)
		r.append([]byte{0x42})
	}()
	buf := make([]byte, 1)
	n, err := r.read(buf, time.Now().Add(5*time.Second))
	require.NoError(t, err)
	require.Equal(t, 1, n)
	assert.Equal(t, byte(0x42), buf[0])
}

func TestByteRing_DeadlineExpiry(t *testing.T) {
	r := newByteRing(8)
	start := time.Now()
	_, err := r.read(make([]byte, 1), time.Now().Add(40*time.Millisecond))
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	assert.GreaterOrEqual(t, time.Since(start), 40*time.Millisecond)

	// already-expired deadline errors immediately, even with data buffered
	r.append([]byte{0x01})
	_, err = r.read(make([]byte, 1), time.Now().Add(-time.Second))
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
}

func TestByteRing_DropOldestOnOverflow(t *testing.T) {
	r := newByteRing(4)
	assert.Zero(t, r.append([]byte("ab")))
	assert.Equal(t, 2, r.append([]byte("cdef")), "two oldest bytes dropped")
	assert.EqualValues(t, 2, r.droppedTotal())
	buf := make([]byte, 4)
	n, err := r.read(buf, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, "cdef", string(buf[:n]))

	// a chunk larger than capacity keeps only its tail
	dropped := r.append([]byte("0123456789"))
	assert.Equal(t, 6, dropped)
	n, err = r.read(buf, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, "6789", string(buf[:n]))
}

func TestByteRing_CloseBeatsBufferedData(t *testing.T) {
	r := newByteRing(8)
	r.append([]byte("data"))
	r.close()
	assert.True(t, r.isClosed())
	_, err := r.read(make([]byte, 4), time.Time{})
	require.ErrorIs(t, err, os.ErrClosed)
}

func TestByteRing_FailDrainsDataThenSurfacesError(t *testing.T) {
	r := newByteRing(8)
	r.append([]byte("ok"))
	portErr := errors.New("port exploded")
	r.fail(portErr)
	buf := make([]byte, 8)
	n, err := r.read(buf, time.Time{})
	require.NoError(t, err, "buffered data drains before the error")
	assert.Equal(t, "ok", string(buf[:n]))
	_, err = r.read(buf, time.Time{})
	require.ErrorContains(t, err, "port exploded")
}

func TestByteRing_ClosedBeatsExpiredDeadline(t *testing.T) {
	r := newByteRing(8)
	r.close()
	_, err := r.read(make([]byte, 1), time.Now().Add(-time.Second))
	require.ErrorIs(t, err, os.ErrClosed, "closed must take precedence over an expired deadline")
}
