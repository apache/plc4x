//go:build windows

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
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// winPort implements Port on Windows with synchronous I/O. Deadlines are
// emulated with COMMTIMEOUTS recomputed before every blocking attempt:
// each ReadFile/WriteFile call is bounded by the remaining time to the
// deadline (capped at maxIOSliceMillis so Close and deadline changes are
// observed even on a port with no traffic).
//
// Reads use the documented COMMTIMEOUTS pattern for "return as soon as at
// least one byte is available, or after the total timeout": interval and
// multiplier MAXDWORD with a bounded positive total constant.
//
// COMMTIMEOUTS is one flat struct per handle governing both directions, so
// every SetCommTimeouts call (whether issued by Read or Write) repopulates
// BOTH directions from a single snapshot of both deadlines - leaving the
// other direction's fields untouched or zeroed would silently corrupt it
// (a zero write timeout means "wait forever"). Because the two directions
// recompute independently, a deadline change made while the other
// direction is mid-attempt can leave that attempt's view of the change
// stale by at most one slice iteration; it self-corrects on that
// direction's next attempt.
//
// Close() drains in-flight I/O attempts (Read/Write/Drain, tracked via ioCount)
// before releasing the handle, repeatedly cancelling outstanding I/O with
// CancelIoEx, so a syscall can never start on, or run past, a closed/
// recycled handle.
//
// Locking discipline: short, bounded comm calls (escape codes, modem
// status, purge, SetCommState) run under mu via withHandle, since they
// complete quickly and never block on the wire. Read, Write, and Drain
// instead run on the in-flight I/O path (beginIO/endIO) and hold no lock
// while their syscall is outstanding, because those syscalls can block for
// an unbounded time (Drain's FlushFileBuffers in particular can stall
// indefinitely under a hardware flow-control deadlock); Close's CancelIoEx
// drain loop needs the syscall to actually be running, not queued up
// behind mu, in order to cancel it.
type winPort struct {
	mu            sync.Mutex
	handle        windows.Handle
	closed        bool
	ioCount       int // number of I/O attempts (Read/Write/Drain) currently between beginIO and endIO
	readDeadline  time.Time
	writeDeadline time.Time
}

const winMaxDword = 0xFFFFFFFF

func openPort(portName string, cfg Config) (Port, error) {
	// The \\.\ prefix is required for COM10 and above and harmless below.
	path := portName
	if !strings.HasPrefix(path, `\\.\`) {
		path = `\\.\` + path
	}
	path16, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, fmt.Errorf("serialport: invalid port name %q: %w", portName, err)
	}
	handle, err := windows.CreateFile(path16,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		0,   // exclusive access
		nil, // default security attributes
		windows.OPEN_EXISTING,
		0, // synchronous I/O
		0)
	if err != nil {
		return nil, fmt.Errorf("serialport: opening %s: %w", portName, os.NewSyscallError("CreateFile", err))
	}
	p := &winPort{handle: handle}
	if err := p.configure(cfg); err != nil {
		_ = windows.CloseHandle(handle)
		return nil, fmt.Errorf("serialport: configuring %s: %w", portName, err)
	}
	return p, nil
}

func (p *winPort) configure(cfg Config) error {
	if err := windows.SetupComm(p.handle, 4096, 4096); err != nil {
		return os.NewSyscallError("SetupComm", err)
	}
	var dcb windows.DCB
	dcb.DCBlength = uint32(unsafe.Sizeof(dcb))
	if err := windows.GetCommState(p.handle, &dcb); err != nil {
		return os.NewSyscallError("GetCommState", err)
	}
	s := makeDCBSettings(cfg)
	dcb.BaudRate = s.BaudRate
	dcb.ByteSize = s.ByteSize
	dcb.Parity = s.Parity
	dcb.StopBits = s.StopBits
	dcb.Flags = s.Flags
	dcb.XonLim = s.XonLim
	dcb.XoffLim = s.XoffLim
	dcb.XonChar = s.XonChar
	dcb.XoffChar = s.XoffChar
	if err := windows.SetCommState(p.handle, &dcb); err != nil {
		return os.NewSyscallError("SetCommState", err)
	}
	return nil
}

// commTimeouts builds a COMMTIMEOUTS value covering both directions from
// the two slice lengths already computed for this instant. The read
// fields always use the MAXDWORD/MAXDWORD/constant "return on first byte
// or after the constant" pattern; the write multiplier stays 0 so the
// constant alone bounds the write.
func commTimeouts(readSliceMs, writeSliceMs uint32) windows.CommTimeouts {
	return windows.CommTimeouts{
		ReadIntervalTimeout:        winMaxDword,
		ReadTotalTimeoutMultiplier: winMaxDword,
		ReadTotalTimeoutConstant:   readSliceMs,
		WriteTotalTimeoutConstant:  writeSliceMs,
	}
}

// sliceMillisForOtherDirection is like ioSliceMillis but never reports an
// expired deadline as 0: a 0 total timeout constant means "wait forever"
// to COMMTIMEOUTS (for the write side) or is at best ambiguous (for the
// read side), which is exactly the cross-contamination this helper exists
// to avoid. An already-expired deadline instead gets the smallest positive
// slice so the other direction times out promptly rather than blocking;
// that direction's own attempt independently detects the expiry as an
// error.
func sliceMillisForOtherDirection(now, deadline time.Time) uint32 {
	ms, ok := ioSliceMillis(now, deadline)
	if !ok {
		return 1
	}
	return ms
}

// beginIO validates the port is still open, registers an in-flight
// attempt (so Close knows to wait for it), and returns a single
// consistent snapshot of the handle and both deadlines. Every attempt
// that successfully calls beginIO must call endIO exactly once when it is
// done, regardless of outcome.
func (p *winPort) beginIO() (handle windows.Handle, readDeadline, writeDeadline time.Time, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return 0, time.Time{}, time.Time{}, os.ErrClosed
	}
	p.ioCount++
	return p.handle, p.readDeadline, p.writeDeadline, nil
}

// endIO retires an in-flight attempt registered by beginIO.
func (p *winPort) endIO() {
	p.mu.Lock()
	p.ioCount--
	p.mu.Unlock()
}

func (p *winPort) isClosed() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.closed
}

func (p *winPort) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}
	for {
		n, done, err := p.readAttempt(buf)
		if done {
			return n, err
		}
	}
}

// readAttempt performs one bounded ReadFile call. done is true when the
// result should be returned to the caller as-is; false means the read
// deadline has not yet expired and Read should loop for another attempt.
func (p *winPort) readAttempt(buf []byte) (n int, done bool, err error) {
	handle, readDeadline, writeDeadline, err := p.beginIO()
	if err != nil {
		return 0, true, err
	}
	defer p.endIO()

	now := time.Now()
	readSliceMs, ok := ioSliceMillis(now, readDeadline)
	if !ok {
		return 0, true, os.ErrDeadlineExceeded
	}
	timeouts := commTimeouts(readSliceMs, sliceMillisForOtherDirection(now, writeDeadline))
	if err := windows.SetCommTimeouts(handle, &timeouts); err != nil {
		return 0, true, os.NewSyscallError("SetCommTimeouts", err)
	}
	var got uint32
	if err := windows.ReadFile(handle, buf, &got, nil); err != nil {
		if errors.Is(err, windows.ERROR_OPERATION_ABORTED) && p.isClosed() {
			return int(got), true, os.ErrClosed
		}
		return int(got), true, os.NewSyscallError("ReadFile", err)
	}
	if got > 0 {
		return int(got), true, nil
	}
	// Timed out with no data: expired deadline is an error, absence of
	// a deadline means keep waiting in bounded slices.
	if !readDeadline.IsZero() && !time.Now().Before(readDeadline) {
		return 0, true, os.ErrDeadlineExceeded
	}
	return 0, false, nil
}

func (p *winPort) Write(buf []byte) (int, error) {
	written := 0
	for {
		n, done, err := p.writeAttempt(buf, written)
		written += n
		if done {
			return written, err
		}
	}
}

// writeAttempt performs one bounded WriteFile call for the remainder of
// buf starting at offset written. n is the number of bytes written by
// this attempt alone. done is true when the result should be returned to
// the caller; false means the write deadline has not yet expired and
// Write should loop for another attempt.
func (p *winPort) writeAttempt(buf []byte, written int) (n int, done bool, err error) {
	handle, readDeadline, writeDeadline, err := p.beginIO()
	if err != nil {
		return 0, true, err
	}
	defer p.endIO()

	now := time.Now()
	writeSliceMs, ok := ioSliceMillis(now, writeDeadline)
	if !ok {
		return 0, true, os.ErrDeadlineExceeded
	}
	timeouts := commTimeouts(sliceMillisForOtherDirection(now, readDeadline), writeSliceMs)
	if err := windows.SetCommTimeouts(handle, &timeouts); err != nil {
		return 0, true, os.NewSyscallError("SetCommTimeouts", err)
	}
	var got uint32
	err = windows.WriteFile(handle, buf[written:], &got, nil)
	if err != nil {
		if errors.Is(err, windows.ERROR_OPERATION_ABORTED) && p.isClosed() {
			return int(got), true, os.ErrClosed
		}
		return int(got), true, os.NewSyscallError("WriteFile", err)
	}
	if written+int(got) == len(buf) {
		return int(got), true, nil
	}
	// Partial write: the slice timed out before the buffer drained.
	if !writeDeadline.IsZero() && !time.Now().Before(writeDeadline) {
		return int(got), true, os.ErrDeadlineExceeded
	}
	return int(got), false, nil
}

func (p *winPort) SetReadDeadline(t time.Time) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return os.ErrClosed
	}
	p.readDeadline = t
	return nil
}

func (p *winPort) SetWriteDeadline(t time.Time) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return os.ErrClosed
	}
	p.writeDeadline = t
	return nil
}

func (p *winPort) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	handle := p.handle
	// Drain attempts already registered via beginIO before releasing the
	// handle: closed is now true, so no new attempt can increment
	// ioCount past this point (beginIO checks closed under mu), but an
	// attempt that snapshotted the handle just before we set closed may
	// not have issued its ReadFile/WriteFile yet. Repeatedly cancel
	// outstanding I/O until every registered attempt has reported back
	// via endIO; CancelIoEx is a no-op if nothing is in flight yet, so we
	// retry on a short poll until the racing attempt's syscall has
	// actually started (and been cancelled) or completed.
	for p.ioCount > 0 {
		p.mu.Unlock()
		_ = windows.CancelIoEx(handle, nil)
		time.Sleep(time.Millisecond)
		p.mu.Lock()
	}
	err := windows.CloseHandle(handle)
	p.mu.Unlock()
	if err != nil {
		return os.NewSyscallError("CloseHandle", err)
	}
	return nil
}

var _ ControlPort = (*winPort)(nil)

// withHandle runs a short, non-blocking comm call under the state lock.
// Only bounded, buffer-level calls may use this (escape codes, status
// reads, purge); blocking I/O keeps using beginIO/endIO.
func (p *winPort) withHandle(name string, op func(h windows.Handle) error) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return os.ErrClosed
	}
	return os.NewSyscallError(name, op(p.handle))
}

func (p *winPort) SetDTR(assert bool) error {
	return p.withHandle("EscapeCommFunction", func(h windows.Handle) error {
		return windows.EscapeCommFunction(h, winEscapeFor(true, assert))
	})
}

func (p *winPort) SetRTS(assert bool) error {
	return p.withHandle("EscapeCommFunction", func(h windows.Handle) error {
		return windows.EscapeCommFunction(h, winEscapeFor(false, assert))
	})
}

func (p *winPort) ModemStatus() (ModemStatus, error) {
	var bits uint32
	err := p.withHandle("GetCommModemStatus", func(h windows.Handle) error {
		return windows.GetCommModemStatus(h, &bits)
	})
	if err != nil {
		return ModemStatus{}, err
	}
	return modemStatusFromWinBits(bits), nil
}

func (p *winPort) SendBreak(d time.Duration) error {
	if d <= 0 {
		return errors.New("serialport: break duration must be greater than 0")
	}
	if err := p.withHandle("SetCommBreak", windows.SetCommBreak); err != nil {
		return err
	}
	time.Sleep(d) // deliberately outside the lock
	return p.withHandle("ClearCommBreak", windows.ClearCommBreak)
}

func (p *winPort) FlushInput() error {
	return p.withHandle("PurgeComm", func(h windows.Handle) error {
		return windows.PurgeComm(h, windows.PURGE_RXCLEAR)
	})
}

func (p *winPort) FlushOutput() error {
	return p.withHandle("PurgeComm", func(h windows.Handle) error {
		return windows.PurgeComm(h, windows.PURGE_TXCLEAR)
	})
}

// Drain blocks until all written data has been transmitted. It runs on the
// in-flight I/O path (beginIO/endIO) rather than under the state lock:
// FlushFileBuffers is unbounded when hardware flow control stalls the
// transmitter, and Close's CancelIoEx drain loop must be able to abort it.
func (p *winPort) Drain() error {
	handle, _, _, err := p.beginIO()
	if err != nil {
		return err
	}
	defer p.endIO()
	if err := windows.FlushFileBuffers(handle); err != nil {
		if errors.Is(err, windows.ERROR_OPERATION_ABORTED) && p.isClosed() {
			return os.ErrClosed
		}
		return os.NewSyscallError("FlushFileBuffers", err)
	}
	return nil
}

func (p *winPort) SetConfig(cfg Config) error {
	normalized, err := cfg.normalize()
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return os.ErrClosed
	}
	return p.configure(normalized)
}
