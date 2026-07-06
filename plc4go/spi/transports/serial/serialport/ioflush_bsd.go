//go:build darwin || freebsd

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

// Classic BSD sys/file.h direction flags for TIOCFLUSH; not in x/sys.
const (
	bsdFREAD  = 0x1
	bsdFWRITE = 0x2
)

func (p *unixPort) FlushInput() error {
	return p.control(func(fd int) error {
		return os.NewSyscallError("ioctl TIOCFLUSH", unix.IoctlSetPointerInt(fd, unix.TIOCFLUSH, bsdFREAD))
	})
}

func (p *unixPort) FlushOutput() error {
	return p.control(func(fd int) error {
		return os.NewSyscallError("ioctl TIOCFLUSH", unix.IoctlSetPointerInt(fd, unix.TIOCFLUSH, bsdFWRITE))
	})
}

func (p *unixPort) Drain() error {
	return p.control(func(fd int) error {
		return os.NewSyscallError("ioctl TIOCDRAIN", unix.IoctlSetInt(fd, unix.TIOCDRAIN, 0))
	})
}
