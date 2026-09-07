//go:build darwin

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
	"os"

	"golang.org/x/sys/unix"
)

// darwinStandardBaudRates are the rates TIOCSETA accepts directly
// (Darwin speed_t values are numeric, so Bxxx == the rate itself).
var darwinStandardBaudRates = map[uint]bool{
	50: true, 75: true, 110: true, 134: true, 150: true, 200: true,
	300: true, 600: true, 1200: true, 1800: true, 2400: true, 4800: true,
	7200: true, 9600: true, 14400: true, 19200: true, 28800: true,
	38400: true, 57600: true, 76800: true, 115200: true, 230400: true,
}

// ioctlIOSSIOSPEED is IOSSIOSPEED from IOKit's serial/ioss.h, which sets
// arbitrary baud rates. Not defined in x/sys/unix; the value (with its
// 4-byte payload encoding) is the one established in the wild by
// github.com/jacobsa/go-serial (Apache License 2.0).
const ioctlIOSSIOSPEED = 0x80045402

// applySettings derived from github.com/jacobsa/go-serial (Apache License 2.0).
func applySettings(fd int, t *unix.Termios, baud uint) error {
	if darwinStandardBaudRates[baud] {
		t.Ispeed = uint64(baud)
		t.Ospeed = uint64(baud)
		if err := unix.IoctlSetTermios(fd, unix.TIOCSETA, t); err != nil {
			return os.NewSyscallError("ioctl TIOCSETA", err)
		}
		return nil
	}
	// Non-standard rate: configure everything else with a placeholder
	// standard rate first, then set the real rate via IOSSIOSPEED.
	t.Ispeed = 19200
	t.Ospeed = 19200
	if err := unix.IoctlSetTermios(fd, unix.TIOCSETA, t); err != nil {
		return os.NewSyscallError("ioctl TIOCSETA", err)
	}
	if err := unix.IoctlSetPointerInt(fd, ioctlIOSSIOSPEED, int(baud)); err != nil {
		return os.NewSyscallError("ioctl IOSSIOSPEED", err)
	}
	return nil
}
