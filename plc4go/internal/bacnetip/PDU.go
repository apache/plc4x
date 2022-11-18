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
	"encoding/binary"
	"fmt"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
	"reflect"
	"regexp"
)

type AddressType int

const (
	NULL_ADDRESS AddressType = iota
	LOCAL_BROADCAST_ADDRESS
	LOCAL_STATION_ADDRESS
	REMOTE_BROADCAST_ADDRESS
	REMOTE_STATION_ADDRESS
	GLOBAL_BROADCAST_ADDRESS
)

func (a AddressType) String() string {
	switch a {
	case NULL_ADDRESS:
		return "NULL_ADDRESS"
	case LOCAL_BROADCAST_ADDRESS:
		return "LOCAL_BROADCAST_ADDRESS"
	case LOCAL_STATION_ADDRESS:
		return "LOCAL_STATION_ADDRESS"
	case REMOTE_BROADCAST_ADDRESS:
		return "REMOTE_BROADCAST_ADDRESS"
	case REMOTE_STATION_ADDRESS:
		return "REMOTE_STATION_ADDRESS"
	case GLOBAL_BROADCAST_ADDRESS:
		return "GLOBAL_BROADCAST_ADDRESS"
	default:
		return "Unknown"
	}
}

type AddressTuple[L any, R any] struct {
	Left  L
	Right R
}

var _field_address = regexp.MustCompile(`((?:\d+)|(?:0x(?:[0-9A-Fa-f][0-9A-Fa-f])+))`)
var _ip_address_port = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)(?::(\d+))?`)
var _ip_address_mask_port = regexp.MustCompile(`(\d+\.\d+\.\d+\.\d+)(?:/(\d+))?(?::(\d+))?`)
var _net_ip_address_port = regexp.MustCompile(`(\d+):` + _ip_address_port.String())
var _at_route = regexp.MustCompile(`(?:[@](?:` + _field_address.String() + `|` + _ip_address_port.String() + `))?`)

var field_address_re = regexp.MustCompile(`^` + _field_address.String() + `$`)
var ip_address_port_re = regexp.MustCompile(`^` + _ip_address_port.String() + `$`)
var ip_address_mask_port_re = regexp.MustCompile(`^` + _ip_address_mask_port.String() + `$`)
var net_ip_address_port_re = regexp.MustCompile(`^` + _net_ip_address_port.String() + `$`)
var net_ip_address_mask_port_re = regexp.MustCompile(`^` + _net_ip_address_port.String() + `$`)

var ethernet_re = regexp.MustCompile(`^([0-9A-Fa-f][0-9A-Fa-f][:]){5}([0-9A-Fa-f][0-9A-Fa-f])$`)
var interface_re = regexp.MustCompile(`^(?:([\w]+))(?::(\d+))?$`)

var net_broadcast_route_re = regexp.MustCompile(`^([0-9])+:[*]` + _at_route.String() + `$`)
var net_station_route_re = regexp.MustCompile(`^([0-9])+:` + _field_address.String() + _at_route.String() + `$`)
var net_ip_address_route_re = regexp.MustCompile(`^([0-9])+:` + _ip_address_port.String() + _at_route.String() + `$`)

var combined_pattern = regexp.MustCompile(`^(?:(?:([0-9]+)|([*])):)?(?:([*])|` + _field_address.String() + `|` + _ip_address_mask_port.String() + `)` + _at_route.String() + `$`)

type Address struct {
	AddrType    AddressType
	AddrNet     *uint16
	AddrAddress []byte
	AddrLen     *uint32
	AddrRoute   *uint32

	AddrIP             *uint32
	AddrMask           *uint32
	AddrHost           *uint32
	AddrSubnet         *uint32
	AddrPort           *uint16
	AddrTuple          *AddressTuple[string, uint16]
	AddrBroadcastTuple *AddressTuple[string, uint16]
}

func NewAddress(args ...interface{}) (*Address, error) {
	log.Debug().Interface("args", args).Msg("NewAddress")
	a := &Address{}
	a.AddrNet = nil
	a.AddrAddress = nil
	a.AddrLen = nil
	a.AddrRoute = nil

	switch len(args) {
	case 1:
		if err := a.decodeAddress(args[0]); err != nil {
			return nil, errors.Wrap(err, "decodeAddress")
		}
	case 2:
		if err := a.decodeAddress(args[1]); err != nil {
			return nil, errors.Wrap(err, "decodeAddress")
		}
		switch a.AddrType {
		case LOCAL_STATION_ADDRESS:
			a.AddrType = REMOTE_STATION_ADDRESS
			var net = (args[0]).(uint16)
			a.AddrNet = &net
		case LOCAL_BROADCAST_ADDRESS:
			a.AddrType = REMOTE_BROADCAST_ADDRESS
			var net = (args[0]).(uint16)
			a.AddrNet = &net
		default:
			return nil, errors.New("unrecognized address ctor form")
		}
	}
	return a, nil
}

// decodeAddress Initialize the address from a string.  Lots of different forms are supported
func (a *Address) decodeAddress(addr interface{}) error {
	log.Debug().Msgf("decodeAddress %v (%T)", addr, addr)

	// start out assuming this is a local station and didn't get routed
	a.AddrType = LOCAL_STATION_ADDRESS
	a.AddrNet = nil
	a.AddrAddress = nil
	a.AddrLen = nil
	a.AddrRoute = nil

	switch {
	case addr == "*":
		log.Debug().Msg("localBroadcast")
		a.AddrType = LOCAL_BROADCAST_ADDRESS
	case addr == "*:*":
		log.Debug().Msg("globalBroadcast")
		a.AddrType = GLOBAL_BROADCAST_ADDRESS
	default:
		switch addr := addr.(type) {
		case net.Addr:
			// TODO: hacked in udp support
			udpAddr := addr.(*net.UDPAddr)
			a.AddrAddress = udpAddr.IP
			length := uint32(len(a.AddrAddress))
			a.AddrLen = &length
			port := uint16(udpAddr.Port)
			a.AddrPort = &port
			addr.String()
		case int:
			log.Debug().Msg("int")
			if addr < 0 || addr > 255 {
				return errors.New("address out of range")
			}
			a.AddrAddress = []byte{byte(addr)}
			length := uint32(1)
			a.AddrLen = &length
		case []byte:
			log.Debug().Msg("byte array")
			a.AddrAddress = addr
			length := uint32(len(addr))
			a.AddrLen = &length

			if *a.AddrLen == 6 {
				ip := ipv4ToUint32(addr[:4])
				a.AddrIP = &ip
				mask := uint32((1 << 32) - 1)
				a.AddrMask = &mask
				host := *a.AddrIP & ^(*a.AddrMask)
				a.AddrHost = &host
				subnet := *a.AddrIP & *a.AddrMask
				a.AddrSubnet = &subnet
				port := portToUint16(addr[4:])
				a.AddrPort = &port

				a.AddrTuple = &AddressTuple[string, uint16]{uint32ToIpv4(*a.AddrIP).String(), *a.AddrPort}
				a.AddrBroadcastTuple = &AddressTuple[string, uint16]{"255.255.255.255", *a.AddrPort}
			}
		case string:
			log.Debug().Msg("str")

			m := combined_pattern.MatchString(addr)
			if m {
				log.Debug().Msg("combined pattern")
				groups := combined_pattern.FindStringSubmatch(addr)
				net := groups[0]
				global_broadcast := groups[1]
				local_broadcast := groups[2]
				local_addr := groups[3]
				local_ip_addr := groups[4]
				local_ip_net := groups[5]
				local_ip_port := groups[6]
				route_addr := groups[7]
				route_ip_addr := groups[8]
				route_ip_port := groups[9]

				a := func(...interface{}) {

				}
				a(net, global_broadcast, local_broadcast, local_addr, local_ip_addr, local_ip_net, local_ip_port, route_addr, route_ip_addr, route_ip_port)
			}
			panic("parsing not yet ported")
		case AddressTuple[string, uint16]:
			uaddr, port := addr.Left, addr.Right
			a.AddrPort = &port

			var addrstr []byte
			if uaddr == "" {
				// when ('', n) is passed it is the local host address, but that could be more than one on a multi homed machine,
				//                    the empty string # means "any".
				addrstr = make([]byte, 4)
			} else {
				addrstr = net.ParseIP(uaddr)
			}
			a.AddrTuple = &AddressTuple[string, uint16]{uaddr, *a.AddrPort}
			log.Debug().Msgf("addrstr: %s", addrstr)

			ip := ipv4ToUint32(addrstr)
			a.AddrIP = &ip
			mask := uint32(0xFFFFFFFF)
			a.AddrMask = &mask
			host := uint32(0)
			a.AddrHost = &host
			subnet := uint32(0)
			a.AddrSubnet = &subnet
			a.AddrBroadcastTuple = a.AddrTuple

			a.AddrAddress = append(addrstr, uint16ToPort(*a.AddrPort)...)
			length := uint32(6)
			a.AddrLen = &length
		case AddressTuple[int, uint16]:
			uaddr, port := addr.Left, addr.Right
			a.AddrPort = &port

			addrstr := uint32ToIpv4(uint32(uaddr))
			a.AddrTuple = &AddressTuple[string, uint16]{addrstr.String(), *a.AddrPort}
			log.Debug().Msgf("addrstr: %s", addrstr)

			ip := ipv4ToUint32(addrstr)
			a.AddrIP = &ip
			mask := uint32(0xFFFFFFFF)
			a.AddrMask = &mask
			host := uint32(0)
			a.AddrHost = &host
			subnet := uint32(0)
			a.AddrSubnet = &subnet
			a.AddrBroadcastTuple = a.AddrTuple

			a.AddrAddress = append(addrstr, uint16ToPort(*a.AddrPort)...)
			length := uint32(6)
			a.AddrLen = &length
		default:
			return errors.Errorf("integer, string or tuple required (Actual %T)", addr)
		}
	}
	return nil
}

func (a *Address) Equals(other interface{}) bool {
	if a == nil && other == nil {
		return true
	} else if a == nil && other != nil {
		return false
	}
	switch other := other.(type) {
	case *Address:
		if a == other {
			return true
		}
		// TODO: don't use reflect here
		return reflect.DeepEqual(a, other)
	case Address:
		// TODO: don't use reflect here
		return reflect.DeepEqual(*a, other)
	default:
		return false
	}
}

func (a *Address) String() string {
	if a == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Address{AddrType: %s, AddrNet: %d, AddrAddress: %x, AddrLen: %d, AddrRoute: %d, AddrIP: %d, AddrMask: %d, AddrHost: %d, AddrSubnet: %d, AddrPort: %d, AddrTuple: %v, AddrBroadcastTuple: %v}", a.AddrType, a.AddrNet, a.AddrAddress, a.AddrLen, a.AddrRoute, a.AddrIP, a.AddrMask, a.AddrHost, a.AddrSubnet, a.AddrPort, a.AddrTuple, a.AddrBroadcastTuple)
}

func portToUint16(port []byte) uint16 {
	switch len(port) {
	case 2:
	default:
		panic("port must be 2 bytes")
	}
	return binary.BigEndian.Uint16(port)
}

func uint16ToPort(number uint16) []byte {
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

type LocalStation struct {
	Address
}

func NewLocalStation(addr interface{}, route *uint32) (*LocalStation, error) {
	l := &LocalStation{}
	l.AddrType = LOCAL_STATION_ADDRESS
	l.AddrRoute = route

	switch addr := addr.(type) {
	case int:
		if addr < 0 || addr > 255 {
			return nil, errors.New("address out of range")
		}
		l.AddrAddress = []byte{byte(addr)}
		length := uint32(1)
		l.AddrLen = &length
	case []byte:
		log.Debug().Msg("bytearray")
		l.AddrAddress = addr
		length := uint32(len(addr))
		l.AddrLen = &length
	default:
		return nil, errors.New("integer or byte array required")
	}
	return l, nil
}

type RemoteStation struct {
	Address
}

func NewRemoteStation(net *uint16, addr interface{}, route *uint32) (*RemoteStation, error) {
	l := &RemoteStation{}
	l.AddrType = REMOTE_STATION_ADDRESS
	l.AddrNet = net
	l.AddrRoute = route

	switch addr := addr.(type) {
	case int:
		if addr < 0 || addr > 255 {
			return nil, errors.New("address out of range")
		}
		l.AddrAddress = []byte{byte(addr)}
		length := uint32(1)
		l.AddrLen = &length
	case []byte:
		log.Debug().Msg("bytearray")
		l.AddrAddress = addr
		length := uint32(len(addr))
		l.AddrLen = &length
	default:
		return nil, errors.New("integer or byte array required")
	}
	return l, nil
}

type LocalBroadcast struct {
	Address
}

func NewLocalBroadcast(route *uint32) (*LocalBroadcast, error) {
	l := &LocalBroadcast{}
	l.AddrType = LOCAL_BROADCAST_ADDRESS
	l.AddrRoute = route
	return l, nil
}

type RemoteBroadcast struct {
	Address
}

func NewRemoteBroadcast(net *uint16, route *uint32) (*RemoteBroadcast, error) {
	r := &RemoteBroadcast{}
	r.AddrType = REMOTE_BROADCAST_ADDRESS
	r.AddrNet = net
	r.AddrRoute = route
	return r, nil
}

type GlobalBroadcast struct {
	Address
}

func NewGlobalBroadcast(route *uint32) (*GlobalBroadcast, error) {
	g := &GlobalBroadcast{}
	g.AddrType = GLOBAL_BROADCAST_ADDRESS
	g.AddrRoute = route
	return g, nil
}

type PCI struct {
	*_PCI
	expectingReply  bool
	networkPriority readWriteModel.NPDUNetworkPriority
}

func NewPCI(msg spi.Message, pduSource Address, pduDestination Address, expectingReply bool, networkPriority readWriteModel.NPDUNetworkPriority) *PCI {
	return &PCI{
		_New_PCI(msg, pduSource, pduDestination),
		expectingReply,
		networkPriority,
	}
}

type _PDU interface {
	spi.Message
	GetPDUSource() Address
	GetPDUDestination() Address
	GetExpectingReply() bool
	GetNetworkPriority() readWriteModel.NPDUNetworkPriority
}

type PDU struct {
	spi.Message
	*PCI
}

func NewPDU(msg spi.Message, pduOptions ...PDUOption) *PDU {
	nullAddress, _ := NewAddress()
	p := &PDU{
		msg,
		NewPCI(msg, *nullAddress, *nullAddress, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE),
	}
	for _, option := range pduOptions {
		option(p)
	}
	return p
}

func NewPDUFromPDU(pdu _PDU, pduOptions ...PDUOption) *PDU {
	msg := pdu.(*PDU).Message
	p := &PDU{
		msg,
		NewPCI(msg, pdu.GetPDUSource(), pdu.GetPDUDestination(), pdu.GetExpectingReply(), pdu.GetNetworkPriority()),
	}
	for _, option := range pduOptions {
		option(p)
	}
	return p
}

func NewPDUWithAllOptions(msg spi.Message, pduSource Address, pduDestination Address, expectingReply bool, networkPriority readWriteModel.NPDUNetworkPriority) *PDU {
	return &PDU{
		msg,
		NewPCI(msg, pduSource, pduDestination, expectingReply, networkPriority),
	}
}

type PDUOption func(pdu *PDU)

func WithPDUSource(pduSource Address) PDUOption {
	return func(pdu *PDU) {
		pdu.pduSource = pduSource
	}
}

func WithPDUDestination(pduDestination Address) PDUOption {
	return func(pdu *PDU) {
		pdu.pduDestination = pduDestination
	}
}

func WithPDUExpectingReply(expectingReply bool) PDUOption {
	return func(pdu *PDU) {
		pdu.expectingReply = expectingReply
	}
}

func WithPDUNetworkPriority(networkPriority readWriteModel.NPDUNetworkPriority) PDUOption {
	return func(pdu *PDU) {
		pdu.networkPriority = networkPriority
	}
}

func (p *PDU) GetPDUSource() Address {
	return p.pduSource
}

func (p *PDU) GetPDUDestination() Address {
	return p.pduDestination
}

func (p *PDU) GetExpectingReply() bool {
	return p.expectingReply
}

func (p *PDU) GetNetworkPriority() readWriteModel.NPDUNetworkPriority {
	return p.networkPriority
}
