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
	"os"

	"golang.org/x/sys/unix"
)

// applySettings configures the port via the kernel termios2 interface
// (TCSETS2), which accepts arbitrary baud rates through BOTHER + explicit
// Ispeed/Ospeed — so standard and non-standard rates share one code path.
// Approach derived from github.com/jacobsa/go-serial (Apache License 2.0);
// unlike that library we use the per-architecture unix.BOTHER constant
// (e.g. 0x1000 on amd64 but 0x1f on ppc64) instead of a hardcoded value.
func applySettings(fd int, t *unix.Termios, baud uint) error {
	t.Cflag &^= unix.CBAUD
	t.Cflag |= unix.BOTHER
	t.Ispeed = uint32(baud)
	t.Ospeed = uint32(baud)
	if err := unix.IoctlSetTermios(fd, unix.TCSETS2, t); err != nil {
		return os.NewSyscallError("ioctl TCSETS2", err)
	}
	return nil
}
