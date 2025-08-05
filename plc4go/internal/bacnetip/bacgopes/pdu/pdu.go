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

package pdu

import (
	"encoding/binary"
	"net"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
)

var _debug = CreateDebugPrinter()

var _shortMask = 0xFFFF

var _longMask = 0xFFFF_FFFF

func portToUint16(port []byte) uint16 {
	switch len(port) {
	case 2:
	default:
		panic("port must be 2 bytes")
	}
	return binary.BigEndian.Uint16(port)
}

func Uint16ToPort(number uint16) []byte {
	port := make([]byte, 2)
	binary.BigEndian.PutUint16(port, number)
	return port
}

func ipv4ToUint32(ip net.IP) uint32 {
	switch len(ip) {
	case 4:
	default:
		panic("ip must be either 4 bytes")
	}
	return binary.BigEndian.Uint32(ip)
}

func uint32ToIpv4(number uint32) net.IP {
	ipv4 := make(net.IP, 4)
	binary.BigEndian.PutUint32(ipv4, number)
	return ipv4
}

// PackIpAddr Given an IP address tuple like ('1.2.3.4', 47808) return the six-octet string
// useful for a BACnet address.
func PackIpAddr(addrTuple *AddressTuple[string, uint16]) (octetString []byte) {
	addr, port := addrTuple.Left, addrTuple.Right
	octetString = append(net.ParseIP(addr).To4(), Uint16ToPort(port)...)
	return
}

// UnpackIpAddr Given a six-octet BACnet address, return an IP address tuple.
func UnpackIpAddr(addr []byte) (addrTuple *AddressTuple[string, uint16]) {
	ip := ipv4ToUint32(addr[:4])
	port := portToUint16(addr[4:])

	return &AddressTuple[string, uint16]{uint32ToIpv4(ip).String(), port}
}

// TODO: convert to struct
func NewLocalStation(addr any, route *Address) (*Address, error) {
	l := &Address{
		_leafName: "LocalStation",
	}
	l.AddrType = LOCAL_STATION_ADDRESS
	l.AddrRoute = route

	switch addr := addr.(type) {
	case int:
		if addr < 0 || addr > 255 {
			return nil, errors.New("address out of range")
		}
		l.AddrAddress = []byte{byte(addr)}
		length := uint8(1)
		l.AddrLen = &length
	case []byte:
		if _debug != nil {
			_debug("    - bytes or bytearray")
		}
		l.AddrAddress = addr
		length := uint8(len(addr))
		l.AddrLen = &length
	default:
		return nil, errors.New("integer or byte array required")
	}
	return l, nil
}

// TODO: convert to struct
func NewRemoteStation(net *uint16, addr any, route *Address) (*Address, error) {
	l := &Address{
		_leafName: "RemoteStation",
	}
	l.AddrType = REMOTE_STATION_ADDRESS
	l.AddrNet = net
	l.AddrRoute = route

	switch addr := addr.(type) {
	case int:
		if addr < 0 || addr > 255 {
			return nil, errors.New("address out of range")
		}
		l.AddrAddress = []byte{byte(addr)}
		length := uint8(1)
		l.AddrLen = &length
	case []byte:
		if _debug != nil {
			_debug("    - bytes or bytearray")
		}
		l.AddrAddress = addr
		length := uint8(len(addr))
		l.AddrLen = &length
	default:
		return nil, errors.New("integer or byte array required")
	}
	return l, nil
}

func NewLocalBroadcast(route *Address) *Address {
	l := &Address{
		_leafName: "LocalBroadcast",
	}
	l.AddrType = LOCAL_BROADCAST_ADDRESS
	l.AddrRoute = route
	return l
}

// TODO: convert to struct
func NewRemoteBroadcast(net uint16, route *Address) *Address {
	r := &Address{
		_leafName: "RemoteBroadcast",
	}
	r.AddrType = REMOTE_BROADCAST_ADDRESS
	r.AddrNet = &net
	r.AddrRoute = route
	return r
}

// TODO: convert to struct
func NewGlobalBroadcast(route *Address) *Address {
	g := &Address{
		_leafName: "GlobalBroadcast",
	}
	g.AddrType = GLOBAL_BROADCAST_ADDRESS
	g.AddrRoute = route
	return g
}
