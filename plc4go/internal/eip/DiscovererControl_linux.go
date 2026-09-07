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

package eip

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// controlDiscoverySocket configures the discovery UDP socket. SO_BROADCAST is
// required on Linux to send to a subnet broadcast address at all: without it,
// sendto() on a broadcast destination fails with EACCES, which is exactly why
// discovery over the shared udp transport (which never sets this) silently
// finds nothing.
//
// Unlike the BACnet discovery socket, EIP discovery binds an ephemeral local
// port rather than the well-known discovery port, and only ever talks to one
// interface's broadcast address per socket, so SO_REUSEADDR/SO_REUSEPORT
// (needed there to let multiple stacks share the well-known port) and
// SO_BINDTODEVICE (needed there to steer wildcard-bound unicast replies to a
// specific NIC) aren't needed here.
func controlDiscoverySocket(c syscall.RawConn) error {
	var sockErr error
	ctrlErr := c.Control(func(fd uintptr) {
		sockErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_BROADCAST, 1)
	})
	if ctrlErr != nil {
		return ctrlErr
	}
	return sockErr
}
