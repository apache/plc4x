//go:build linux || darwin || freebsd

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
	"time"

	"golang.org/x/sys/unix"
)

// unixPort is an open POSIX serial port: a runtime-poller-managed *os.File
// (providing Read/Write/Close and deadlines) plus ioctl-backed control
// operations. Control ioctls go through SyscallConn().Control — never
// .Fd(), which would de-register the file from the poller.
type unixPort struct {
	*os.File
}

var _ ControlPort = (*unixPort)(nil)

// control runs op on the raw fd while the runtime holds it stable.
//
// SyscallConn() and rawConn.Control() only fail for an invalid or closing
// *os.File — internal/poll surfaces that as "use of closed file", which
// does not satisfy errors.Is(err, os.ErrClosed) and would otherwise diverge
// from Read/Write-after-close (which do) and from the Windows ControlPort
// implementation (which returns os.ErrClosed). Normalize to os.ErrClosed
// here; op's own error is returned separately as opErr.
func (p *unixPort) control(op func(fd int) error) error {
	rawConn, err := p.SyscallConn()
	if err != nil {
		return os.ErrClosed
	}
	var opErr error
	if err := rawConn.Control(func(fd uintptr) { opErr = op(int(fd)) }); err != nil {
		return os.ErrClosed
	}
	return opErr
}

func (p *unixPort) setModemBits(bits int, assert bool) error {
	return p.control(func(fd int) error {
		if assert {
			return os.NewSyscallError("ioctl TIOCMBIS", unix.IoctlSetPointerInt(fd, unix.TIOCMBIS, bits))
		}
		return os.NewSyscallError("ioctl TIOCMBIC", unix.IoctlSetPointerInt(fd, unix.TIOCMBIC, bits))
	})
}

func (p *unixPort) SetDTR(assert bool) error { return p.setModemBits(unix.TIOCM_DTR, assert) }
func (p *unixPort) SetRTS(assert bool) error { return p.setModemBits(unix.TIOCM_RTS, assert) }

func (p *unixPort) ModemStatus() (ModemStatus, error) {
	var status ModemStatus
	err := p.control(func(fd int) error {
		bits, err := unix.IoctlGetInt(fd, unix.TIOCMGET)
		if err != nil {
			return os.NewSyscallError("ioctl TIOCMGET", err)
		}
		status = ModemStatus{
			CTS: bits&unix.TIOCM_CTS != 0,
			DSR: bits&unix.TIOCM_DSR != 0,
			DCD: bits&unix.TIOCM_CAR != 0,
			RI:  bits&unix.TIOCM_RNG != 0,
		}
		return nil
	})
	return status, err
}

func (p *unixPort) SendBreak(d time.Duration) error {
	if d <= 0 {
		return errors.New("serialport: break duration must be greater than 0")
	}
	if err := p.control(func(fd int) error {
		return os.NewSyscallError("ioctl TIOCSBRK", unix.IoctlSetInt(fd, unix.TIOCSBRK, 0))
	}); err != nil {
		return err
	}
	time.Sleep(d)
	return p.control(func(fd int) error {
		return os.NewSyscallError("ioctl TIOCCBRK", unix.IoctlSetInt(fd, unix.TIOCCBRK, 0))
	})
}

func (p *unixPort) SetConfig(cfg Config) error {
	normalized, err := cfg.normalize()
	if err != nil {
		return err
	}
	t, err := makeTermios(normalized)
	if err != nil {
		return err
	}
	return p.control(func(fd int) error {
		return applySettings(fd, t, normalized.BaudRate)
	})
}

// openPort opens the device non-blocking so the Go runtime poller manages
// it — that is what makes SetReadDeadline/SetWriteDeadline work on a plain
// *os.File (which therefore is the Port implementation on POSIX).
//
// The termios configuration is applied through SyscallConn().Control and
// never via Fd(): calling Fd() would switch the file to blocking mode and
// de-register it from the poller, silently breaking deadline support.
func openPort(portName string, cfg Config) (Port, error) {
	t, err := makeTermios(cfg)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(portName, os.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0)
	if err != nil {
		return nil, fmt.Errorf("serialport: opening %s: %w", portName, err)
	}
	rawConn, err := f.SyscallConn()
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("serialport: accessing raw connection of %s: %w", portName, err)
	}
	var applyErr error
	if err := rawConn.Control(func(fd uintptr) {
		applyErr = applySettings(int(fd), t, cfg.BaudRate)
	}); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("serialport: raw control on %s: %w", portName, err)
	}
	if applyErr != nil {
		_ = f.Close()
		return nil, fmt.Errorf("serialport: configuring %s: %w", portName, applyErr)
	}
	// Deadlines are part of the Port contract: fail fast if the runtime
	// poller could not adopt this fd instead of degrading silently.
	if err := f.SetReadDeadline(time.Time{}); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("serialport: %s does not support I/O deadlines: %w", portName, err)
	}
	return &unixPort{File: f}, nil
}
