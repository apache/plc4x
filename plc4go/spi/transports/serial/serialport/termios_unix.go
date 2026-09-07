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

	"golang.org/x/sys/unix"
)

// makeTermios translates an already-normalized Config into a raw-mode
// termios structure per POSIX: non-canonical, no echo, no signal or
// flow-control character processing, receiver enabled, modem control
// lines ignored (CLOCAL). Baud rate is set separately by the per-OS
// applySettings, because speed handling differs between kernels.
func makeTermios(cfg Config) (*unix.Termios, error) {
	t := &unix.Termios{}
	t.Cflag |= unix.CLOCAL | unix.CREAD

	switch cfg.DataBits {
	case 5:
		t.Cflag |= unix.CS5
	case 6:
		t.Cflag |= unix.CS6
	case 7:
		t.Cflag |= unix.CS7
	case 8:
		t.Cflag |= unix.CS8
	}

	switch cfg.StopBits {
	case StopBitsOne:
		// CSTOPB cleared
	case StopBitsTwo:
		t.Cflag |= unix.CSTOPB
	case StopBitsOnePointFive:
		return nil, errors.New("serialport: 1.5 stop bits are not supported on this platform")
	}

	switch cfg.Parity {
	case ParityNone:
	case ParityOdd:
		t.Cflag |= unix.PARENB | unix.PARODD
	case ParityEven:
		t.Cflag |= unix.PARENB
	case ParityMark:
		if err := setMarkSpaceParity(t, true); err != nil {
			return nil, err
		}
	case ParitySpace:
		if err := setMarkSpaceParity(t, false); err != nil {
			return nil, err
		}
	}

	if cfg.RTSCTSFlowControl {
		t.Cflag |= unix.CRTSCTS
	}

	if cfg.XONXOFFFlowControl {
		t.Iflag |= unix.IXON | unix.IXOFF
		t.Cc[unix.VSTART] = 0x11 // DC1
		t.Cc[unix.VSTOP] = 0x13  // DC3
	}

	// The port is opened O_NONBLOCK and driven by the runtime poller, which
	// makes the kernel ignore VMIN/VTIME; set the conventional raw-mode
	// values anyway to document intent (and for any future blocking use).
	t.Cc[unix.VMIN] = 1
	t.Cc[unix.VTIME] = 0
	return t, nil
}
