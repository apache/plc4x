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
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

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
	return f, nil
}
