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

package socketcan

import (
	"net"
	"os"

	"golang.org/x/sys/unix"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

// openSocket opens a non-blocking CAN_RAW socket bound to the interface. Non-blocking matters: it is
// what lets os.NewFile hand the descriptor to the runtime poller, so deadlines work and Close wakes
// up a pending read.
func openSocket(interfaceName string) (*os.File, error) {
	networkInterface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, errors.Wrapf(err, "no network interface %s", interfaceName)
	}
	fd, err := unix.Socket(unix.AF_CAN, unix.SOCK_RAW|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, unix.CAN_RAW)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a CAN_RAW socket")
	}
	if err := unix.Bind(fd, &unix.SockaddrCAN{Ifindex: networkInterface.Index}); err != nil {
		_ = unix.Close(fd)
		return nil, errors.Wrapf(err, "error binding the CAN socket to %s", interfaceName)
	}
	return os.NewFile(uintptr(fd), "socketcan:"+interfaceName), nil
}
