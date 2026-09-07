//go:build linux

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
	"context"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"

	"github.com/apache/plc4x/plc4go/spi/transports"
)

// openPTY allocates a pseudo-terminal pair; the slave path stands in for a
// real serial device (same pattern as the serialport package's own tests).
func openPTY(t *testing.T) (master *os.File, slavePath string) {
	t.Helper()
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|unix.O_NOCTTY, 0)
	require.NoError(t, err, "opening /dev/ptmx")
	t.Cleanup(func() { _ = master.Close() })
	ptn, err := unix.IoctlGetInt(int(master.Fd()), unix.TIOCGPTN)
	require.NoError(t, err, "TIOCGPTN")
	require.NoError(t, unix.IoctlSetPointerInt(int(master.Fd()), unix.TIOCSPTLCK, 0), "unlocking pty slave")
	return master, fmt.Sprintf("/dev/pts/%d", ptn)
}

func TestTransportInstance_EndToEndOnPTY(t *testing.T) {
	master, slavePath := openPTY(t)

	instance := NewTransportInstance(slavePath, 115200, 1, NewTransport())
	require.NoError(t, instance.Connect(context.Background()))
	t.Cleanup(func() { _ = instance.Close() })
	assert.True(t, instance.IsConnected())

	// Write goes out on the wire (fixes the historic "Not connected" bug:
	// Connect previously never flipped the connected flag).
	require.NoError(t, instance.Write(context.Background(), []byte{0x01, 0x02, 0x03}))
	wire := make([]byte, 3)
	_, err := master.Read(wire)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, wire)

	// Reads work and honor context deadlines through the buffered layer.
	_, err = master.Write([]byte{0xCA, 0xFE})
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data, err := instance.Read(ctx, 2)
	require.NoError(t, err)
	assert.Equal(t, []byte{0xCA, 0xFE}, data)

	// A read into silence now actually times out instead of hanging:
	// SetReadDeadline is no longer a no-op.
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer shortCancel()
	start := time.Now()
	_, err = instance.Read(shortCtx, 1)
	require.Error(t, err)
	assert.Less(t, time.Since(start), 5*time.Second)
}

func TestTransportInstance_CloseIsIdempotentOnPTY(t *testing.T) {
	_, slavePath := openPTY(t)
	instance := NewTransportInstance(slavePath, 9600, 1, NewTransport())
	require.NoError(t, instance.Connect(context.Background()))
	require.NoError(t, instance.Close())
	assert.False(t, instance.IsConnected())
	require.NoError(t, instance.Close(), "second close must be a no-op")
}

func TestTransportInstance_WriteDeadlineDoesNotStick(t *testing.T) {
	_, slavePath := openPTY(t)
	instance := NewTransportInstance(slavePath, 9600, 1, NewTransport())
	require.NoError(t, instance.Connect(context.Background()))
	t.Cleanup(func() { _ = instance.Close() })

	shortCtx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	require.NoError(t, instance.Write(shortCtx, []byte{0x01}))

	time.Sleep(60 * time.Millisecond) // let the first write's deadline lapse

	// A deadline-less write must not inherit the lapsed deadline.
	require.NoError(t, instance.Write(context.Background(), []byte{0x02}))
}

func TestTransportInstance_OptionsAppliedOnPTY(t *testing.T) {
	master, slavePath := openPTY(t)

	transport := NewTransport()
	instance, err := transport.CreateTransportInstance(
		url.URL{Scheme: "serial", Path: slavePath},
		map[string][]string{
			"serial.data-bits": {"7"},
			"serial.parity":    {"EVEN"}, // deliberately non-canonical case
			"serial.stop-bits": {"2"},
			"serial.dtr":       {"true"}, // must warn, not fail, on a pty
		},
	)
	require.NoError(t, err)
	require.NoError(t, instance.Connect(context.Background()))
	t.Cleanup(func() { _ = instance.Close() })

	// The pty pair shares one termios; reading it back on the MASTER side
	// proves the options reached the kernel.
	readBack, err := unix.IoctlGetTermios(int(master.Fd()), unix.TCGETS)
	require.NoError(t, err)
	assert.NotZero(t, readBack.Cflag&unix.CS7, "CS7")
	assert.NotZero(t, readBack.Cflag&unix.CSTOPB, "CSTOPB")
	// PARENB deliberately not asserted (known pty parity quirk).
}

func TestTransportInstance_FallbackReadDeadlineBoundsSilentRead(t *testing.T) {
	_, slavePath := openPTY(t)
	transport := NewTransport()
	instance, err := transport.CreateTransportInstance(
		url.URL{Scheme: "serial", Path: slavePath},
		map[string][]string{"serial.read-timeout-ms": {"200"}},
	)
	require.NoError(t, err)
	require.NoError(t, instance.Connect(context.Background()))
	t.Cleanup(func() { _ = instance.Close() })

	start := time.Now()
	_, err = instance.Read(context.Background(), 1) // no ctx deadline
	elapsed := time.Since(start)

	require.Error(t, err, "silent line must time out via the fallback deadline")
	assert.GreaterOrEqual(t, elapsed, 200*time.Millisecond)
	assert.Less(t, elapsed, 5*time.Second)
}

func TestTransportInstance_ExplicitCtxDeadlineBeatsFallback(t *testing.T) {
	_, slavePath := openPTY(t)
	transport := NewTransport()
	instance, err := transport.CreateTransportInstance(
		url.URL{Scheme: "serial", Path: slavePath},
		map[string][]string{"serial.read-timeout-ms": {"60000"}},
	)
	require.NoError(t, err)
	require.NoError(t, instance.Connect(context.Background()))
	t.Cleanup(func() { _ = instance.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err = instance.Read(ctx, 1)
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.Less(t, elapsed, 5*time.Second, "the 100ms ctx deadline must win over the 60s fallback")
}

func TestTransportInstance_ReusePortSharesOnePTY(t *testing.T) {
	master, slavePath := openPTY(t)
	transport := NewTransport()
	options := map[string][]string{"serial.reuse-port": {"true"}}

	first, err := transport.CreateTransportInstance(url.URL{Scheme: "serial", Path: slavePath}, options)
	require.NoError(t, err)
	second, err := transport.CreateTransportInstance(url.URL{Scheme: "serial", Path: slavePath}, options)
	require.NoError(t, err)

	require.NoError(t, first.Connect(context.Background()))
	t.Cleanup(func() { _ = first.Close() })
	require.NoError(t, second.Connect(context.Background()))
	t.Cleanup(func() { _ = second.Close() })

	_, err = master.Write([]byte{0xCA, 0xFE})
	require.NoError(t, err)

	for _, instance := range []transports.TransportInstance{first, second} {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		data, err := instance.Read(ctx, 2)
		cancel()
		require.NoError(t, err)
		assert.Equal(t, []byte{0xCA, 0xFE}, data, "both connections see the broadcast")
	}
}

func TestTransportInstance_ReusePortConfigMismatch(t *testing.T) {
	_, slavePath := openPTY(t)
	transport := NewTransport()

	first, err := transport.CreateTransportInstance(url.URL{Scheme: "serial", Path: slavePath},
		map[string][]string{"serial.reuse-port": {"true"}, "serial.baud-rate": {"9600"}})
	require.NoError(t, err)
	require.NoError(t, first.Connect(context.Background()))
	t.Cleanup(func() { _ = first.Close() })

	second, err := transport.CreateTransportInstance(url.URL{Scheme: "serial", Path: slavePath},
		map[string][]string{"serial.reuse-port": {"true"}, "serial.baud-rate": {"19200"}})
	require.NoError(t, err)
	err = second.Connect(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), slavePath)
}

func TestTransportInstance_DedicatedInterframeDelayOnPTY(t *testing.T) {
	master, slavePath := openPTY(t)
	transport := NewTransport()
	instance, err := transport.CreateTransportInstance(url.URL{Scheme: "serial", Path: slavePath},
		map[string][]string{"serial.interframe-delay": {"60"}})
	require.NoError(t, err)
	require.NoError(t, instance.Connect(context.Background()))
	t.Cleanup(func() { _ = instance.Close() })

	require.NoError(t, instance.Write(context.Background(), []byte{0x01}))
	start := time.Now()
	require.NoError(t, instance.Write(context.Background(), []byte{0x02}))
	assert.GreaterOrEqual(t, time.Since(start), 50*time.Millisecond)

	buf := make([]byte, 2)
	total := 0
	for total < 2 {
		n, err := master.Read(buf[total:])
		require.NoError(t, err)
		total += n
	}
	assert.Equal(t, []byte{0x01, 0x02}, buf)
}
