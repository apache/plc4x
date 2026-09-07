//go:build freebsd

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

// applySettings sets the rate directly: FreeBSD speed_t values are numeric
// and the driver rejects rates it cannot provide.
func applySettings(fd int, t *unix.Termios, baud uint) error {
	t.Ispeed = uint32(baud)
	t.Ospeed = uint32(baud)
	if err := unix.IoctlSetTermios(fd, unix.TIOCSETA, t); err != nil {
		return os.NewSyscallError("ioctl TIOCSETA", err)
	}
	return nil
}
