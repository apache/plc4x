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
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/debugging"
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

func NewAddressTuple[L any, R any](l L, r R) *AddressTuple[L, R] {
	return &AddressTuple[L, R]{l, r}
}

func (a *AddressTuple[L, R]) deepCopy() *AddressTuple[L, R] {
	if a == nil {
		return nil
	}
	// TODO: check if that works like intended (might just fail for pointer types)
	return &AddressTuple[L, R]{*CopyPtr[L](&a.Left), *CopyPtr[R](&a.Right)}
}

func (a *AddressTuple[L, R]) DeepCopy() any {
	return a.deepCopy()
}

func (a *AddressTuple[L, R]) Format(s fmt.State, v rune) {
	switch v {
	case 's', 'v', 'r':
		_, _ = fmt.Fprint(s, a.String())
	}
}

func (a *AddressTuple[L, R]) String() string {
	return fmt.Sprintf("(%v, %v)", a.Left, a.Right)
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

var CombinedPattern = regexp.MustCompile(`^(?:(?:([0-9]+)|([*])):)?(?:([*])|` + _field_address.String() + `|` + _ip_address_mask_port.String() + `)` + _at_route.String() + `$`)

//go:generate go tool plc4xGenerator -type=Address -prefix=pdu_
type Address struct {
	AddrType    AddressType
	AddrNet     *uint16
	AddrAddress []byte
	AddrLen     *uint8
	AddrRoute   *Address

	AddrIP             *uint32
	AddrMask           *uint32
	AddrHost           *uint32
	AddrSubnet         *uint32
	AddrPort           *uint16
	AddrTuple          *AddressTuple[string, uint16]
	AddrBroadcastTuple *AddressTuple[string, uint16]

	_leafName string
}

func NewAddress(args Args) (*Address, error) {
	if _debug != nil {
		_debug("__init__ %r", args)
	}
	a := &Address{
		_leafName: "Address", // TODO: usually this is done by extract, we leaf address for now as this would imply changing type to an interface and that is a big change.
	}
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
			var addrNet *uint16
			switch arg0 := args[0].(type) {
			case uint16:
				addrNet = &arg0
			case *uint16:
				addrNet = arg0
			}
			a.AddrNet = addrNet
		case LOCAL_BROADCAST_ADDRESS:
			a.AddrType = REMOTE_BROADCAST_ADDRESS
			var addrNet *uint16
			switch arg0 := args[0].(type) {
			case uint16:
				addrNet = &arg0
			case *uint16:
				addrNet = arg0
			}
			a.AddrNet = addrNet
		default:
			return nil, errors.New("unrecognized address ctor form")
		}
	}
	return a, nil
}

// decodeAddress Initialize the address from a string.  Lots of different forms are supported
func (a *Address) decodeAddress(addr any) error {
	if _debug != nil {
		_debug("decode_address %r %T", addr, addr)
	}

	// start out assuming this is a local station and didn't get routed
	a.AddrType = LOCAL_STATION_ADDRESS
	a.AddrNet = nil
	a.AddrAddress = nil
	a.AddrLen = nil
	a.AddrRoute = nil

	switch {
	case addr == "*":
		if _debug != nil {
			_debug("    - localBroadcast")
		}
		a.AddrType = LOCAL_BROADCAST_ADDRESS
	case addr == "*:*":
		if _debug != nil {
			_debug("    - globalBroadcast")
		}
		a.AddrType = GLOBAL_BROADCAST_ADDRESS
	default:
		switch addr := addr.(type) {
		case net.Addr:
			// TODO: hacked in udp support
			udpAddr := addr.(*net.UDPAddr)
			a.AddrAddress = udpAddr.IP.To4()
			if a.AddrAddress == nil {
				a.AddrAddress = udpAddr.IP.To16()
			}
			length := uint8(len(a.AddrAddress))
			a.AddrLen = &length
			port := uint16(udpAddr.Port)
			a.AddrPort = &port
			addr.String()
		case int, int32, int64, uint, uint32, uint64:
			iaddr, err := strconv.ParseInt(fmt.Sprintf("%v", addr), 10, 64) // TODO: bit ugly but better than repeating all of it
			if err != nil {
				panic(err)
			}
			if _debug != nil {
				_debug("    - int")
			}
			if iaddr < 0 || iaddr > 255 {
				return errors.New("address out of range")
			}
			a.AddrAddress = []byte{byte(iaddr)}
			length := uint8(1)
			a.AddrLen = &length
		case []byte:
			if _debug != nil {
				_debug("    - bytes or bytearray")
			}
			a.AddrAddress = addr
			length := uint8(len(addr))
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
			if _debug != nil {
				_debug("    - str")
			}

			m := CombinedPattern.MatchString(addr)
			if m {
				if _debug != nil {
					_debug("    - combined pattern")
				}
				groups := CombinedPattern.FindStringSubmatch(addr)
				_net := groups[1]
				globalBroadcast := groups[2]
				localBroadcast := groups[3]
				localAddr := groups[4]
				localIpAddr := groups[5]
				localIpNet := groups[6]
				localIpPort := groups[7]
				routeAddr := groups[8]
				routeIpAddr := groups[9]
				routeIpPort := groups[10]

				if globalBroadcast != "" && localBroadcast != "" {
					if _debug != nil {
						_debug("    - global broadcast")
					}
					a.AddrType = GLOBAL_BROADCAST_ADDRESS
				} else if _net != "" && localBroadcast != "" {
					if _debug != nil {
						_debug("    - remote broadcast")
					}
					netAddr, err := strconv.ParseUint(_net, 10, 16)
					if err != nil {
						return errors.Wrap(err, "can't parse net")
					}
					a.AddrType = REMOTE_BROADCAST_ADDRESS
					var netAddr16 = uint16(netAddr)
					a.AddrNet = &netAddr16
				} else if localBroadcast != "" {
					if _debug != nil {
						_debug("    - local broadcast")
					}
					a.AddrType = LOCAL_BROADCAST_ADDRESS
				} else if _net != "" {
					if _debug != nil {
						_debug("    - remote station")
					}
					netAddr, err := strconv.ParseUint(_net, 10, 16)
					if err != nil {
						return errors.Wrap(err, "can't parse net")
					}
					a.AddrType = REMOTE_STATION_ADDRESS
					var netAddr16 = uint16(netAddr)
					a.AddrNet = &netAddr16
				}

				if localAddr != "" {
					if _debug != nil {
						_debug("    - simple address")
					}
					if strings.HasPrefix(localAddr, "0x") {
						var err error
						a.AddrAddress, err = Xtob(localAddr[2:])
						if err != nil {
							return errors.Wrap(err, "can't parse local address")
						}
						addrLen := uint8(len(a.AddrAddress))
						a.AddrLen = &addrLen
					} else {
						localAddr, err := strconv.ParseUint(localAddr, 10, 8)
						if err != nil {
							return errors.Wrap(err, "can't parse local addr")
						}
						a.AddrAddress = []byte{byte(localAddr)}
						addrLen := uint8(1)
						a.AddrLen = &addrLen
					}
				}

				if localIpAddr != "" {
					if _debug != nil {
						_debug("    - ip address")
					}
					if localIpPort == "" {
						localIpPort = "47808"
					}
					if localIpNet == "" {
						localIpNet = "32"
					}

					localIpPortParse, err := strconv.ParseUint(localIpPort, 10, 16)
					if err != nil {
						return errors.Wrap(err, "can't parse local addr")
					}
					localIpPort16 := uint16(localIpPortParse)
					a.AddrPort = &localIpPort16
					a.AddrTuple = &AddressTuple[string, uint16]{localIpAddr, *a.AddrPort}
					if _debug != nil {
						_debug("    - addrTuple: %r", a.AddrTuple)
					}

					parseAddr, err := netip.ParseAddr(localIpAddr)
					if err != nil {
						return errors.Wrap(err, "can't parse local addr")
					}
					addrIp := binary.BigEndian.Uint32(parseAddr.AsSlice())
					a.AddrIP = &addrIp
					localIpNetParse, err := strconv.ParseUint(localIpNet, 10, 16)
					if err != nil {
						return errors.Wrap(err, "can't parse local addr")
					}
					localNetPort16 := uint16(localIpNetParse)
					mask := uint32((_longMask << (32 - localNetPort16)) & _longMask)
					a.AddrMask = &mask
					addrHost := *a.AddrIP &^ (*a.AddrMask)
					a.AddrHost = &addrHost
					addrSubnet := *a.AddrIP & (*a.AddrMask)
					a.AddrSubnet = &addrSubnet

					bcast := *a.AddrSubnet | ^(*a.AddrMask)
					bcastIpReverse := make(net.IP, 4)
					binary.BigEndian.PutUint32(bcastIpReverse, bcast&uint32(_longMask))
					a.AddrBroadcastTuple = &AddressTuple[string, uint16]{bcastIpReverse.String(), *a.AddrPort}
					if _debug != nil {
						_debug("    - addrBroadcastTuple: %r", a.AddrBroadcastTuple)
					}

					portReverse := make([]byte, 2)
					binary.BigEndian.PutUint16(portReverse, *a.AddrPort&uint16(_shortMask))
					a.AddrAddress = append(parseAddr.AsSlice(), portReverse...)
					addrLen := uint8(6)
					a.AddrLen = &addrLen
				}

				if !Settings.RouteAware && (routeAddr != "" || routeIpAddr != "") {
					if _debug != nil {
						_debug("route provided but not route aware: %v", addr)
					}
				}

				if routeAddr != "" {
					if strings.HasPrefix(routeAddr, "0x") {
						xtob, err := Xtob(routeAddr[2:])
						if err != nil {
							return errors.Wrap(err, "can't parse route addr")
						}
						a.AddrRoute, err = NewAddress(NA(xtob))
						if err != nil {
							return errors.Wrap(err, "can't parse route")
						}
					} else {
						routeAddr, err := strconv.ParseUint(routeAddr, 10, 32)
						if err != nil {
							return errors.Wrap(err, "can't parse route addr")
						}
						a.AddrRoute, err = NewAddress(NA(routeAddr))
						if err != nil {
							return errors.Wrap(err, "can't create route")
						}
						if _debug != nil {
							_debug("    - addrRoute: %r", a.AddrRoute)
						}
					}
				} else if routeIpAddr != "" {
					if routeIpPort == "" {
						routeIpPort = "47808"
					}
					var err error
					tuple := &AddressTuple[string, string]{routeIpAddr, routeIpPort}
					a.AddrRoute, err = NewAddress(NA(tuple))
					if err != nil {
						return errors.Wrap(err, "can't create route")
					}
				}

				return nil
			}

			if ethernet_re.MatchString(addr) {
				if _debug != nil {
					_debug("    - ethernet")
				}
				var err error
				a.AddrAddress, err = Xtob(addr)
				if err != nil {
					return errors.Wrap(err, "can't parse address")
				}
				addrLen := uint8(len(a.AddrAddress))
				a.AddrLen = &addrLen
				return nil
			}

			intR := regexp.MustCompile(`^\d+$`)
			if intR.MatchString(addr) {
				if _debug != nil {
					_debug("    - int")
				}

				parseUint, err := strconv.ParseUint(addr, 10, 8)
				if err != nil {
					return errors.Wrap(err, "can't parse int")
				}
				a.AddrAddress = []byte{byte(parseUint)}
				addrLen := uint8(len(a.AddrAddress))
				a.AddrLen = &addrLen
				return nil
			}

			remoteBroadcast := regexp.MustCompile(`^\d+:[*]$`)
			if remoteBroadcast.MatchString(addr) {
				if _debug != nil {
					_debug("    - remote broadcast")
				}

				parseUint, err := strconv.ParseUint(addr[:len(addr)-2], 10, 16)
				if err != nil {
					return errors.Wrap(err, "can't parse int")
				}

				a.AddrType = REMOTE_BROADCAST_ADDRESS
				addrNet := uint16(parseUint)
				a.AddrNet = &addrNet
				a.AddrAddress = nil
				a.AddrLen = nil
				return nil
			}

			remoteStation := regexp.MustCompile(`^\d+:[*]$`)
			if remoteStation.MatchString(addr) {
				if _debug != nil {
					_debug("    - remote station")
				}

				split := strings.Split(addr, ":")
				_net, _addr := split[0], split[1]
				parseNetUint, err := strconv.ParseUint(_net, 10, 16)
				if err != nil {
					return errors.Wrap(err, "can't parse int")
				}
				parseAddrUint, err := strconv.ParseUint(_addr, 10, 8)
				if err != nil {
					return errors.Wrap(err, "can't parse int")
				}

				a.AddrType = REMOTE_STATION_ADDRESS
				addrNet := uint16(parseNetUint)
				a.AddrNet = &addrNet
				a.AddrAddress = []byte{byte(parseAddrUint)}
				addrLen := uint8(len(a.AddrAddress))
				a.AddrLen = &addrLen
				return nil
			}

			modernHexString := regexp.MustCompile(`^0x([0-9A-Fa-f][0-9A-Fa-f])+$`)
			if modernHexString.MatchString(addr) {
				if _debug != nil {
					_debug("    - modern hex string")
				}

				var err error
				a.AddrAddress, err = Xtob(addr[2:])
				if err != nil {
					return errors.Wrap(err, "can't parse address")
				}
				addrLen := uint8(len(a.AddrAddress))
				a.AddrLen = &addrLen
				return nil
			}

			oldSchoolHexString := regexp.MustCompile(`^X'([0-9A-Fa-f][0-9A-Fa-f])+'$`)
			if oldSchoolHexString.MatchString(addr) {
				if _debug != nil {
					_debug("    - modern hex string")
				}

				var err error
				a.AddrAddress, err = Xtob(addr[2 : len(addr)-1])
				if err != nil {
					return errors.Wrap(err, "can't parse address")
				}
				addrLen := uint8(len(a.AddrAddress))
				a.AddrLen = &addrLen
				return nil
			}

			remoteStationWithModernHexString := regexp.MustCompile(`^\d+:0x([0-9A-Fa-f][0-9A-Fa-f])+$`)
			if remoteStationWithModernHexString.MatchString(addr) {
				if _debug != nil {
					_debug("    - remote station with modern hex string")
				}

				split := strings.Split(addr, ":")
				_net, _addr := split[0], split[1]
				parseNetUint, err := strconv.ParseUint(_net, 10, 16)
				if err != nil {
					return errors.Wrap(err, "can't parse int")
				}

				a.AddrType = REMOTE_STATION_ADDRESS
				addrNet := uint16(parseNetUint)
				a.AddrNet = &addrNet
				a.AddrAddress, err = Xtob(_addr[2:])
				if err != nil {
					return errors.Wrap(err, "can't parse addr")
				}
				addrLen := uint8(len(a.AddrAddress))
				a.AddrLen = &addrLen
				return nil
			}

			remoteStationWithOldHexString := regexp.MustCompile(`^\d+:X'([0-9A-Fa-f][0-9A-Fa-f])+'$`)
			if remoteStationWithOldHexString.MatchString(addr) {
				if _debug != nil {
					_debug("    - remote station with modern hex string")
				}

				split := strings.Split(addr, ":")
				_net, _addr := split[0], split[1]
				parseNetUint, err := strconv.ParseUint(_net, 10, 16)
				if err != nil {
					return errors.Wrap(err, "can't parse int")
				}

				a.AddrType = REMOTE_STATION_ADDRESS
				addrNet := uint16(parseNetUint)
				a.AddrNet = &addrNet
				a.AddrAddress, err = Xtob(_addr[2 : len(_addr)-1])
				if err != nil {
					return errors.Wrap(err, "can't parse addr")
				}
				addrLen := uint8(len(a.AddrAddress))
				a.AddrLen = &addrLen
				return nil
			}

			if interface_re.MatchString(addr) {
				if _debug != nil {
					_debug("    - interface name with optional port")
				}

				groups := interface_re.FindStringSubmatch(addr)
				_interface := groups[1]
				_port := groups[2]
				if _port != "" {
					parseUint, err := strconv.ParseUint(_port, 10, 16)
					if err != nil {
						return errors.Wrap(err, "can't parse port")
					}
					port := uint16(parseUint)
					a.AddrPort = &port
				} else {
					port := uint16(47808)
					a.AddrPort = &port
				}

				_ = _interface
				_ = _port
				panic("implement me")
				return nil
			}

			return errors.New("unrecognized format")
		case *AddressTuple[string, uint16]:
			uaddr, port := addr.Left, addr.Right
			a.AddrPort = &port

			var addrstr []byte
			if uaddr == "" {
				// when ('', n) is passed it is the local host address, but that could be more than one on a multi homed machine,
				//                    the empty string # means "any".
				addrstr = make([]byte, 4)
			} else {
				addrstr = net.ParseIP(uaddr).To4()
			}
			a.AddrTuple = &AddressTuple[string, uint16]{uaddr, *a.AddrPort}
			if _debug != nil {
				_debug("    - addrstr: %v", addrstr)
			}

			ip := ipv4ToUint32(addrstr)
			a.AddrIP = &ip
			mask := uint32(0xFFFFFFFF)
			a.AddrMask = &mask
			host := uint32(0)
			a.AddrHost = &host
			subnet := uint32(0)
			a.AddrSubnet = &subnet
			a.AddrBroadcastTuple = a.AddrTuple

			a.AddrAddress = append(addrstr, Uint16ToPort(*a.AddrPort)...)
			length := uint8(6)
			a.AddrLen = &length
		case *AddressTuple[string, string]:
			uaddr, port := addr.Left, addr.Right
			portParse, err := strconv.ParseUint(port, 10, 16)
			if err != nil {
				return errors.Wrap(err, "can't parse port")
			}
			portInt := uint16(portParse)
			a.AddrPort = &portInt

			var addrstr []byte
			if uaddr == "" {
				// when ('', n) is passed it is the local host address, but that could be more than one on a multi homed machine,
				//                    the empty string # means "any".
				addrstr = make([]byte, 4)
			} else {
				addrstr = net.ParseIP(uaddr).To4()
			}
			a.AddrTuple = &AddressTuple[string, uint16]{uaddr, *a.AddrPort}
			if _debug != nil {
				_debug("    - addrstr: %r", addrstr)
			}

			ip := ipv4ToUint32(addrstr)
			a.AddrIP = &ip
			mask := uint32(0xFFFFFFFF)
			a.AddrMask = &mask
			host := uint32(0)
			a.AddrHost = &host
			subnet := uint32(0)
			a.AddrSubnet = &subnet
			a.AddrBroadcastTuple = a.AddrTuple

			a.AddrAddress = append(addrstr, Uint16ToPort(*a.AddrPort)...)
			length := uint8(6)
			a.AddrLen = &length
		case *AddressTuple[int, uint16]:
			uaddr, port := addr.Left, addr.Right
			a.AddrPort = &port

			addrstr := uint32ToIpv4(uint32(uaddr))
			a.AddrTuple = &AddressTuple[string, uint16]{addrstr.String(), *a.AddrPort}
			if _debug != nil {
				_debug("    - addrstr: %r", addrstr)
			}

			ip := ipv4ToUint32(addrstr)
			a.AddrIP = &ip
			mask := uint32(0xFFFFFFFF)
			a.AddrMask = &mask
			host := uint32(0)
			a.AddrHost = &host
			subnet := uint32(0)
			a.AddrSubnet = &subnet
			a.AddrBroadcastTuple = a.AddrTuple

			a.AddrAddress = append(addrstr, Uint16ToPort(*a.AddrPort)...)
			length := uint8(6)
			a.AddrLen = &length
		default:
			return errors.Errorf("integer, string or tuple required (Actual %T)", addr)
		}
	}
	return nil
}

func (a *Address) Equals(other any) bool {
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
		thisString := a.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			if _debug != nil {
				_debug("Mismatch %v != %v", thisString, otherString)
			}
		}
		return equals
	case Address:
		thisString := a.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			if _debug != nil {
				_debug("Mismatch %v != %v", thisString, otherString)
			}
		}
		return equals
	case AddressTuple[string, uint16]:
		thisString := a.AddrTuple.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			if _debug != nil {
				_debug("Mismatch %v != %v", thisString, otherString)
			}
		}
		return equals
	case *AddressTuple[string, uint16]:
		thisString := a.AddrTuple.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			if _debug != nil {
				_debug("Mismatch %v != %v", thisString, otherString)
			}
		}
		return equals
	case AddressTuple[string, int]:
		thisString := a.AddrTuple.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			if _debug != nil {
				_debug("Mismatch %v != %v", thisString, otherString)
			}
		}
		return equals
	case *AddressTuple[string, int]:
		thisString := a.AddrTuple.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			if _debug != nil {
				_debug("Mismatch %v != %v", thisString, otherString)
			}
		}
		return equals
	default:
		if _debug != nil {
			_debug("Unmapped type comparison %T", other)
		}
		return false
	}
}

func (a *Address) Format(s fmt.State, v rune) {
	if a == nil {
		_, _ = fmt.Fprint(s, "<nil>")
		return
	}
	switch v {
	case 'r':
		_, _ = fmt.Fprintf(s, "<%s %s>", a._leafName, a.String())
	case 'v', 's':
		_, _ = fmt.Fprint(s, a.String())
	}
}

func (a *Address) deepCopy() *Address {
	if a == nil {
		return nil
	}
	return &Address{
		a.AddrType,
		CopyPtr(a.AddrNet),
		bytes.Clone(a.AddrAddress),
		CopyPtr(a.AddrLen),
		a.AddrRoute.deepCopy(),
		CopyPtr(a.AddrIP),
		CopyPtr(a.AddrMask),
		CopyPtr(a.AddrHost),
		CopyPtr(a.AddrSubnet),
		CopyPtr(a.AddrPort),
		a.AddrTuple.deepCopy(),
		a.AddrBroadcastTuple.deepCopy(),
		a._leafName,
	}
}

func (a *Address) DeepCopy() any {
	return a.deepCopy()
}

func (a *Address) AlternateString() (string, bool) {
	if IsDebuggingActive() || true { // TODO: figure out what to do when we want the below string in testing etc...
		if a == nil {
			return "<nil>", true
		}
		result := ""
		if a.AddrType == NULL_ADDRESS {
			result = "Null"
		} else if a.AddrType == LOCAL_BROADCAST_ADDRESS {
			result = "*"
		} else if a.AddrType == LOCAL_STATION_ADDRESS {
			if a.AddrLen != nil && *a.AddrLen == 1 {
				result += fmt.Sprintf("%v", a.AddrAddress[0])
			} else {
				port := binary.BigEndian.Uint16(a.AddrAddress[len(a.AddrAddress)-2:])
				if len(a.AddrAddress) == 6 && port >= 47808 && port <= 47823 {
					var octests = make([]string, 4)
					for i, address := range a.AddrAddress[0:4] {
						octests[i] = fmt.Sprintf("%d", address)
					}
					result += fmt.Sprintf("%v", strings.Join(octests, "."))
					if port != 47808 {
						result += fmt.Sprintf(":%v", port)
					}
				} else {
					result += fmt.Sprintf("0x%x", a.AddrAddress)
				}
			}
		} else if a.AddrType == REMOTE_BROADCAST_ADDRESS {
			result = fmt.Sprintf("%d:*", *a.AddrNet)
		} else if a.AddrType == REMOTE_STATION_ADDRESS {
			result = fmt.Sprintf("%d:", *a.AddrNet)
			if a.AddrLen != nil && *a.AddrLen == 1 {
				result += fmt.Sprintf("%v", a.AddrAddress[0])
			} else {
				port := binary.BigEndian.Uint16(a.AddrAddress[len(a.AddrAddress)-2:])
				if len(a.AddrAddress) == 6 && port >= 47808 && port <= 47823 {
					var octests = make([]string, 4)
					for i, address := range a.AddrAddress[0:4] {
						octests[i] = fmt.Sprintf("%d", address)
					}
					result += fmt.Sprintf("%v", strings.Join(octests, "."))
					if port != 47808 {
						result += fmt.Sprintf(":%v", port)
					}
				} else {
					result += fmt.Sprintf("0x%x", a.AddrAddress)
				}
			}
		} else if a.AddrType == GLOBAL_BROADCAST_ADDRESS {
			result = "*:*"
		} else {
			panic("Unknown address type: " + a.AddrType.String())
		}

		if a.AddrRoute != nil {
			result += fmt.Sprintf("@%s", a.AddrRoute)
		}
		return result, true
	}
	return "", false
}
