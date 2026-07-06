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

// Package serialport provides access to serial ports (RS-232/RS-485 style
// UART devices) with net.Conn-like read/write deadline semantics.
//
// It is an internal replacement for the unmaintained
// github.com/jacobsa/go-serial dependency, informed by that library
// (Apache License 2.0) and the underlying platform specifications
// (POSIX termios, Windows DCB/COMMTIMEOUTS).
package serialport

import (
	"errors"
	"fmt"
	"io"
	"time"
)

// Parity configures the parity bit generated and checked by the UART.
type Parity int

const (
	ParityNone Parity = iota
	ParityOdd
	ParityEven
)

// StopBits configures the number of stop bits.
type StopBits int

const (
	StopBitsOne          StopBits = iota
	StopBitsOnePointFive          // only supported on Windows, and only with 5 data bits
	StopBitsTwo
)

// Config describes how a serial port is configured on Open. The zero value
// of every field except BaudRate is a valid default (8 data bits, 1 stop
// bit, no parity, no flow control).
//
// There are deliberately no VMIN/VTIME-style read-threshold options: ports
// are opened non-blocking and Read returns as soon as at least one byte is
// available, bounded by SetReadDeadline — the same semantics as net.Conn.
type Config struct {
	BaudRate          uint
	DataBits          uint // 5..8; 0 defaults to 8
	StopBits          StopBits
	Parity            Parity
	RTSCTSFlowControl bool
}

// Port is an open serial port. Read blocks until at least one byte is
// available, the read deadline expires (os.ErrDeadlineExceeded), or the
// port is closed.
type Port interface {
	io.ReadWriteCloser
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

// ErrUnsupportedPlatform is returned by Open on operating systems without
// a serial port implementation.
var ErrUnsupportedPlatform = errors.New("serial ports are not supported on this platform")

// Open opens and configures the named serial port (e.g. "/dev/ttyUSB0" or
// "COM3").
func Open(portName string, cfg Config) (Port, error) {
	if portName == "" {
		return nil, errors.New("serialport: port name must not be empty")
	}
	normalized, err := cfg.normalize()
	if err != nil {
		return nil, err
	}
	return openPort(portName, normalized)
}

// normalize validates cfg and fills in defaults.
func (c Config) normalize() (Config, error) {
	if c.BaudRate == 0 {
		return c, errors.New("serialport: baud rate must be greater than 0")
	}
	if c.DataBits == 0 {
		c.DataBits = 8
	}
	if c.DataBits < 5 || c.DataBits > 8 {
		return c, fmt.Errorf("serialport: invalid data bits %d (must be 5..8)", c.DataBits)
	}
	switch c.StopBits {
	case StopBitsOne, StopBitsOnePointFive, StopBitsTwo:
	default:
		return c, fmt.Errorf("serialport: invalid stop bits value %d", c.StopBits)
	}
	switch c.Parity {
	case ParityNone, ParityOdd, ParityEven:
	default:
		return c, fmt.Errorf("serialport: invalid parity value %d", c.Parity)
	}
	if c.StopBits == StopBitsOnePointFive && c.DataBits != 5 {
		return c, errors.New("serialport: 1.5 stop bits require 5 data bits")
	}
	return c, nil
}
