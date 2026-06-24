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

package bacnetip

import (
	"syscall"

	"golang.org/x/sys/unix"
)

// controlDiscoverySocket configures the discovery UDP socket. SO_BINDTODEVICE
// ties the wildcard-bound socket to one NIC so it reliably receives unicast
// IAm replies arriving on that interface (notably on virtual/test interfaces,
// where a plain wildcard socket may not be handed the unicast datagrams).
// SO_REUSEADDR/SO_REUSEPORT allow rebinding across discovery sweeps and
// co-existing with other BACnet listeners; SO_BROADCAST permits sending the
// WhoIs to the subnet broadcast address. This mirrors gobacnet's datalink.
func controlDiscoverySocket(c syscall.RawConn, interfaceName string) error {
	var sockErr error
	ctrlErr := c.Control(func(fd uintptr) {
		if sockErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); sockErr != nil {
			return
		}
		if sockErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); sockErr != nil {
			return
		}
		if sockErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_BROADCAST, 1); sockErr != nil {
			return
		}
		if interfaceName != "" {
			sockErr = unix.SetsockoptString(int(fd), unix.SOL_SOCKET, unix.SO_BINDTODEVICE, interfaceName)
		}
	})
	if ctrlErr != nil {
		return ctrlErr
	}
	return sockErr
}
