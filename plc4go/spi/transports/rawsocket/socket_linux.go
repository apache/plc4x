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

package rawsocket

import (
	"net"
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/apache/plc4x/plc4go/spi/errors"
)

type linuxSocket struct {
	file     *os.File
	ifindex  int
	hardware net.HardwareAddr
}

// htons is for the protocol argument of AF_PACKET, which is in network byte order.
func htons(value uint16) uint16 {
	return value<<8 | value>>8
}

// openSocket opens a non-blocking AF_PACKET socket bound to the interface, with the filter attached
// before binding so no unfiltered frame is ever queued. Outgoing frames are not looped back.
func openSocket(interfaceName string, program []Instruction) (packetSocket, error) {
	networkInterface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return nil, errors.Wrapf(err, "no network interface %s", interfaceName)
	}
	fd, err := unix.Socket(unix.AF_PACKET, unix.SOCK_RAW|unix.SOCK_NONBLOCK|unix.SOCK_CLOEXEC, int(htons(unix.ETH_P_ALL)))
	if err != nil {
		return nil, errors.Wrap(err, "error creating an AF_PACKET socket (the process needs CAP_NET_RAW)")
	}
	fail := func(err error, what string) (packetSocket, error) {
		_ = unix.Close(fd)
		return nil, errors.Wrap(err, what)
	}
	filter := make([]unix.SockFilter, len(program))
	for i, instruction := range program {
		filter[i] = unix.SockFilter{Code: instruction.Code, Jt: instruction.Jt, Jf: instruction.Jf, K: instruction.K}
	}
	fprog := unix.SockFprog{Len: uint16(len(filter)), Filter: (*unix.SockFilter)(unsafe.Pointer(&filter[0]))}
	if err := unix.SetsockoptSockFprog(fd, unix.SOL_SOCKET, unix.SO_ATTACH_FILTER, &fprog); err != nil {
		return fail(err, "error attaching the packet filter")
	}
	// Linux 4.20+. Without it our own frames come back to us; the driver can live with that, so a
	// failure here is not fatal.
	_ = unix.SetsockoptInt(fd, unix.SOL_PACKET, unix.PACKET_IGNORE_OUTGOING, 1)
	if err := unix.Bind(fd, &unix.SockaddrLinklayer{Protocol: htons(unix.ETH_P_ALL), Ifindex: networkInterface.Index}); err != nil {
		return fail(err, "error binding the raw socket to "+interfaceName)
	}
	return &linuxSocket{
		file:     os.NewFile(uintptr(fd), "raw-socket:"+interfaceName),
		ifindex:  networkInterface.Index,
		hardware: networkInterface.HardwareAddr,
	}, nil
}

func (s *linuxSocket) readFrame(buffer []byte) (int, error) {
	return s.file.Read(buffer)
}

func (s *linuxSocket) writeFrame(frame []byte) error {
	n, err := s.file.Write(frame)
	if err != nil {
		return err
	}
	if n != len(frame) {
		return errors.Errorf("short raw write: %d of %d bytes", n, len(frame))
	}
	return nil
}

func (s *linuxSocket) setReadDeadline(deadline time.Time) error {
	return s.file.SetReadDeadline(deadline)
}

func (s *linuxSocket) setWriteDeadline(deadline time.Time) error {
	return s.file.SetWriteDeadline(deadline)
}

func (s *linuxSocket) hardwareAddress() net.HardwareAddr {
	return s.hardware
}

func (s *linuxSocket) close() error {
	return s.file.Close()
}
