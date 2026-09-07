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

package serialport

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

// openPTY allocates a pseudo-terminal pair and returns the master side plus
// the slave device path. A PTY slave accepts the same termios ioctls as a
// real serial device (baud is faked but accepted), which lets the full
// open→configure→read/write→deadline→close path run without hardware.
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

func TestOpenOnPTY_ReadAndWrite(t *testing.T) {
	master, slavePath := openPTY(t)

	port, err := Open(slavePath, Config{BaudRate: 115200})
	require.NoError(t, err)
	defer port.Close()

	// master -> port
	_, err = master.Write([]byte("hello"))
	require.NoError(t, err)
	require.NoError(t, port.SetReadDeadline(time.Now().Add(5*time.Second)))
	buf := make([]byte, 16)
	n, err := port.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(buf[:n]))

	// port -> master
	_, err = port.Write([]byte("world"))
	require.NoError(t, err)
	n, err = master.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, "world", string(buf[:n]))
}

func TestOpenOnPTY_FullConfigAccepted(t *testing.T) {
	// A PTY accepts arbitrary termios settings; this proves applySettings
	// round-trips a non-default configuration without error.
	_, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{
		BaudRate: 2500000, // non-standard on purpose: exercises the BOTHER path
		DataBits: 7, StopBits: StopBitsTwo, Parity: ParityEven, RTSCTSFlowControl: true,
	})
	require.NoError(t, err)
	require.NoError(t, port.Close())
}

func TestOpenOnPTY_ReadDeadlineExpires(t *testing.T) {
	_, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	defer port.Close()

	require.NoError(t, port.SetReadDeadline(time.Now().Add(50*time.Millisecond)))
	start := time.Now()
	_, err = port.Read(make([]byte, 1))
	elapsed := time.Since(start)

	require.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrDeadlineExceeded), "want os.ErrDeadlineExceeded, got %v", err)
	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond)
	assert.Less(t, elapsed, 5*time.Second)
}

func TestOpenOnPTY_ClearingDeadlineRestoresBlockingRead(t *testing.T) {
	master, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	defer port.Close()

	// Expired deadline: read fails immediately.
	require.NoError(t, port.SetReadDeadline(time.Now().Add(-time.Second)))
	_, err = port.Read(make([]byte, 1))
	require.True(t, errors.Is(err, os.ErrDeadlineExceeded))

	// Cleared deadline: read blocks until data arrives.
	require.NoError(t, port.SetReadDeadline(time.Time{}))
	go func() {
		time.Sleep(20 * time.Millisecond)
		_, _ = master.Write([]byte{0x42})
	}()
	buf := make([]byte, 1)
	n, err := port.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	assert.Equal(t, byte(0x42), buf[0])
}

func TestOpenOnPTY_UseAfterCloseFails(t *testing.T) {
	_, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	require.NoError(t, port.Close())

	_, err = port.Read(make([]byte, 1))
	require.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrClosed), "read after close: want os.ErrClosed, got %v", err)

	_, err = port.Write([]byte{0x01})
	require.Error(t, err)
	assert.True(t, errors.Is(err, os.ErrClosed), "write after close: want os.ErrClosed, got %v", err)
}

func TestOpenOnPTY_ControlOpsAfterCloseReturnErrClosed(t *testing.T) {
	_, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	cp := port.(ControlPort)
	require.NoError(t, port.Close())

	require.ErrorIs(t, cp.FlushInput(), os.ErrClosed)
	require.ErrorIs(t, cp.Drain(), os.ErrClosed)
	require.ErrorIs(t, cp.SetDTR(true), os.ErrClosed)
	_, err = cp.ModemStatus()
	require.ErrorIs(t, err, os.ErrClosed)
	require.ErrorIs(t, cp.SetConfig(Config{BaudRate: 19200}), os.ErrClosed)
	require.ErrorIs(t, cp.SendBreak(10*time.Millisecond), os.ErrClosed)
}

func TestOpenNonexistentDeviceFails(t *testing.T) {
	_, err := Open("/dev/plc4x-does-not-exist", Config{BaudRate: 9600})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "plc4x-does-not-exist")
}

func TestOpenOnPTY_ImplementsControlPort(t *testing.T) {
	_, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	defer port.Close()
	_, ok := port.(ControlPort)
	require.True(t, ok, "Open's result must implement ControlPort")
}

func TestOpenOnPTY_SetConfigReconfiguresLive(t *testing.T) {
	_, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	defer port.Close()
	cp := port.(ControlPort)

	require.NoError(t, cp.SetConfig(Config{BaudRate: 19200, DataBits: 7, StopBits: StopBitsTwo, Parity: ParityEven}))

	// Read back the termios state to prove the reconfiguration landed.
	rawConn, err := port.(*unixPort).SyscallConn()
	require.NoError(t, err)
	var readBack *unix.Termios
	var ioctlErr error
	require.NoError(t, rawConn.Control(func(fd uintptr) {
		readBack, ioctlErr = unix.IoctlGetTermios(int(fd), unix.TCGETS)
	}))
	require.NoError(t, ioctlErr)
	assert.NotZero(t, readBack.Cflag&unix.CS7, "CS7")
	assert.NotZero(t, readBack.Cflag&unix.CSTOPB, "CSTOPB")
	// PARENB is deliberately not asserted: pty drivers may drop or reject
	// parity flags (observed on Linux 7.x), while CS7/CSTOPB reliably stick
	// and prove the reconfiguration reached the kernel.

	// Invalid reconfiguration is rejected by validation, not the kernel.
	require.ErrorContains(t, cp.SetConfig(Config{}), "baud rate")
}

func TestOpenOnPTY_FlushInputDiscardsPendingData(t *testing.T) {
	master, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	defer port.Close()
	cp := port.(ControlPort)

	_, err = master.Write([]byte("stale"))
	require.NoError(t, err)
	time.Sleep(50 * time.Millisecond) // let the pty deliver into the input queue
	require.NoError(t, cp.FlushInput())

	require.NoError(t, port.SetReadDeadline(time.Now().Add(100*time.Millisecond)))
	_, err = port.Read(make([]byte, 8))
	require.True(t, errors.Is(err, os.ErrDeadlineExceeded), "flushed data must be gone, got %v", err)
}

func TestOpenOnPTY_FlushOutputAndDrainSucceed(t *testing.T) {
	_, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	defer port.Close()
	cp := port.(ControlPort)
	require.NoError(t, cp.FlushOutput())
	require.NoError(t, cp.Drain())
}

func TestOpenOnPTY_SendBreak(t *testing.T) {
	_, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	defer port.Close()
	cp := port.(ControlPort)

	require.Error(t, cp.SendBreak(0), "non-positive duration must be rejected")

	// Linux ptys accept break as a no-op (send_break returns 0 without a
	// driver break_ctl), so success plus elapsed time is what we can pin.
	start := time.Now()
	require.NoError(t, cp.SendBreak(60*time.Millisecond))
	assert.GreaterOrEqual(t, time.Since(start), 60*time.Millisecond)
}

func TestOpenOnPTY_ModemOpsSurfaceDriverErrors(t *testing.T) {
	// Linux ptys have no tiocmget/tiocmset — the calls must fail cleanly
	// (on real UARTs they succeed; this pins error propagation, not values).
	_, slavePath := openPTY(t)
	port, err := Open(slavePath, Config{BaudRate: 9600})
	require.NoError(t, err)
	defer port.Close()
	cp := port.(ControlPort)

	_, err = cp.ModemStatus()
	require.Error(t, err)
	require.Error(t, cp.SetDTR(true))
	require.Error(t, cp.SetRTS(false))
}
