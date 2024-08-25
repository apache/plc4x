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
	"io"
	"net"
	"net/netip"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/globals"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var _shortMask = 0xFFFF
var _longMask = 0xFFFF_FFFF

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

var combined_pattern = regexp.MustCompile(`^(?:(?:([0-9]+)|([*])):)?(?:([*])|` + _field_address.String() + `|` + _ip_address_mask_port.String() + `)` + _at_route.String() + `$`)

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

	log zerolog.Logger
}

func NewAddress(localLog zerolog.Logger, args ...any) (*Address, error) {
	a := &Address{
		log: localLog,
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
			var addrNet = (args[0]).(uint16)
			a.AddrNet = &addrNet
		case LOCAL_BROADCAST_ADDRESS:
			a.AddrType = REMOTE_BROADCAST_ADDRESS
			var addrNet = (args[0]).(uint16)
			a.AddrNet = &addrNet
		default:
			return nil, errors.New("unrecognized address ctor form")
		}
	}
	return a, nil
}

// decodeAddress Initialize the address from a string.  Lots of different forms are supported
func (a *Address) decodeAddress(addr any) error {
	a.log.Debug().Type("addrType", addr).Interface("addr", addr).Msg("decodeAddress")

	// start out assuming this is a local station and didn't get routed
	a.AddrType = LOCAL_STATION_ADDRESS
	a.AddrNet = nil
	a.AddrAddress = nil
	a.AddrLen = nil
	a.AddrRoute = nil

	switch {
	case addr == "*":
		a.log.Debug().Msg("localBroadcast")
		a.AddrType = LOCAL_BROADCAST_ADDRESS
	case addr == "*:*":
		a.log.Debug().Msg("globalBroadcast")
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
			a.log.Debug().Msg("int")
			if iaddr < 0 || iaddr > 255 {
				return errors.New("address out of range")
			}
			a.AddrAddress = []byte{byte(iaddr)}
			length := uint8(1)
			a.AddrLen = &length
		case []byte:
			a.log.Debug().Msg("byte array")
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
			a.log.Trace().Msg("str")

			m := combined_pattern.MatchString(addr)
			if m {
				a.log.Trace().Msg("combined pattern")
				groups := combined_pattern.FindStringSubmatch(addr)
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
					a.log.Trace().Msg("global broadcast")
					a.AddrType = GLOBAL_BROADCAST_ADDRESS
				} else if _net != "" && localBroadcast != "" {
					a.log.Trace().Msg("remote broadcast")
					netAddr, err := strconv.ParseUint(_net, 10, 16)
					if err != nil {
						return errors.Wrap(err, "can't parse net")
					}
					a.AddrType = REMOTE_BROADCAST_ADDRESS
					var netAddr16 = uint16(netAddr)
					a.AddrNet = &netAddr16
				} else if localBroadcast != "" {
					a.log.Trace().Msg("local broadcast")
					a.AddrType = LOCAL_BROADCAST_ADDRESS
				} else if _net != "" {
					a.log.Trace().Msg("remote station")
					netAddr, err := strconv.ParseUint(_net, 10, 16)
					if err != nil {
						return errors.Wrap(err, "can't parse net")
					}
					a.AddrType = REMOTE_STATION_ADDRESS
					var netAddr16 = uint16(netAddr)
					a.AddrNet = &netAddr16
				}

				if localAddr != "" {
					a.log.Trace().Msg("simple address")
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
					a.log.Trace().Msg("ip address")
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
					a.log.Debug().Stringer("addrTuple", a.AddrTuple).Msg("addrTuple")

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
					a.log.Debug().Stringer("addrBroadcastTuple", a.AddrBroadcastTuple).Msg("addrBroadcastTuple")

					portReverse := make([]byte, 2)
					binary.BigEndian.PutUint16(portReverse, *a.AddrPort&uint16(_shortMask))
					a.AddrAddress = append(parseAddr.AsSlice(), portReverse...)
					addrLen := uint8(6)
					a.AddrLen = &addrLen
				}

				if !Settings.RouteAware && (routeAddr != "" || routeIpAddr != "") {
					a.log.Warn().Msgf("route provided but not route aware: %v", addr)
				}

				if routeAddr != "" {
					if strings.HasPrefix(routeAddr, "0x") {
						xtob, err := Xtob(routeAddr[2:])
						if err != nil {
							return errors.Wrap(err, "can't parse route addr")
						}
						a.AddrRoute, err = NewAddress(a.log, xtob)
						if err != nil {
							return errors.Wrap(err, "can't parse route")
						}
					} else {
						routeAddr, err := strconv.ParseUint(routeAddr, 10, 32)
						if err != nil {
							return errors.Wrap(err, "can't parse route addr")
						}
						a.AddrRoute, err = NewAddress(a.log, routeAddr)
						if err != nil {
							return errors.Wrap(err, "can't create route")
						}
						a.log.Debug().Interface("addRoute", a.AddrRoute).Msg("addRoute")
					}
				} else if routeIpAddr != "" {
					if routeIpPort == "" {
						routeIpPort = "47808"
					}
					var err error
					tuple := &AddressTuple[string, string]{routeIpAddr, routeIpPort}
					a.AddrRoute, err = NewAddress(a.log, tuple)
					if err != nil {
						return errors.Wrap(err, "can't create route")
					}
				}

				return nil
			}

			if ethernet_re.MatchString(addr) {
				a.log.Trace().Msg("ethernet")
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
				a.log.Trace().Msg("int")

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
				a.log.Trace().Msg("remote broadcast")

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
				a.log.Trace().Msg("remote station")

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
				a.log.Trace().Msg("modern hex string")

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
				a.log.Trace().Msg("modern hex string")

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
				a.log.Trace().Msg("remote station with modern hex string")

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
				a.log.Trace().Msg("remote station with modern hex string")

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
				a.log.Trace().Msg("interface name with optional port")

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
			a.log.Debug().Hex("addrstr", addrstr).Msg("addrstr:")

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
			a.log.Debug().Hex("addrstr", addrstr).Msg("addrstr:")

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
			length := uint8(6)
			a.AddrLen = &length
		case *AddressTuple[int, uint16]:
			uaddr, port := addr.Left, addr.Right
			a.AddrPort = &port

			addrstr := uint32ToIpv4(uint32(uaddr))
			a.AddrTuple = &AddressTuple[string, uint16]{addrstr.String(), *a.AddrPort}
			a.log.Debug().Hex("addrstr", addrstr).Msg("addrstr:")

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
			a.log.Debug().Str("thisString", thisString).Str("otherString", otherString).Msg("Mismatch")
		}
		return equals
	case Address:
		thisString := a.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			a.log.Debug().Str("thisString", thisString).Str("otherString", otherString).Msg("Mismatch")
		}
		return equals
	case AddressTuple[string, uint16]:
		thisString := a.AddrTuple.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			a.log.Debug().Str("thisString", thisString).Str("otherString", otherString).Msg("Mismatch")
		}
		return equals
	case *AddressTuple[string, uint16]:
		thisString := a.AddrTuple.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			a.log.Debug().Str("thisString", thisString).Str("otherString", otherString).Msg("Mismatch")
		}
		return equals
	case AddressTuple[string, int]:
		thisString := a.AddrTuple.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			a.log.Debug().Str("thisString", thisString).Str("otherString", otherString).Msg("Mismatch")
		}
		return equals
	case *AddressTuple[string, int]:
		thisString := a.AddrTuple.String()
		otherString := other.String()
		equals := thisString == otherString
		if !equals {
			a.log.Debug().Str("thisString", thisString).Str("otherString", otherString).Msg("Mismatch")
		}
		return equals
	default:
		a.log.Debug().Type("otherType", other).Msg("Unmapped type comparison")
		return false
	}
}

func (a *Address) String() string {
	if a == nil {
		return "<nil>"
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
	return result
}

func (a *Address) GoString() string { //TODO: not valid yet, needs adjustments to have proper output syntax
	if a == nil {
		return "<nil>"
	}
	var sb strings.Builder
	sb.WriteString(a.AddrType.String())
	if a.AddrNet != nil {
		_, _ = fmt.Fprintf(&sb, ", net: %d", *a.AddrNet)
	}
	if len(a.AddrAddress) > 0 {
		_, _ = fmt.Fprintf(&sb, ", address: %d", a.AddrAddress)
	}
	if a.AddrLen != nil {
		_, _ = fmt.Fprintf(&sb, " with len %d", *a.AddrLen)
	}
	if a.AddrRoute != nil {
		_, _ = fmt.Fprintf(&sb, ", route: %s", a.AddrRoute)
	}
	if a.AddrIP != nil {
		_, _ = fmt.Fprintf(&sb, ", ip: %d", *a.AddrIP)
	}
	if a.AddrMask != nil {
		_, _ = fmt.Fprintf(&sb, ", mask: %d", *a.AddrMask)
	}
	if a.AddrHost != nil {
		_, _ = fmt.Fprintf(&sb, ", host: %d", *a.AddrHost)
	}
	if a.AddrSubnet != nil {
		_, _ = fmt.Fprintf(&sb, ", subnet: %d", *a.AddrSubnet)
	}
	if a.AddrPort != nil {
		_, _ = fmt.Fprintf(&sb, ", port: %d", *a.AddrPort)
	}
	if a.AddrTuple != nil {
		_, _ = fmt.Fprintf(&sb, ", tuple: %s", a.AddrTuple)
	}
	if a.AddrBroadcastTuple != nil {
		_, _ = fmt.Fprintf(&sb, ", broadcast tuple: %s", a.AddrBroadcastTuple)
	}
	return sb.String()
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

// PackIpAddr Given an IP address tuple like ('1.2.3.4', 47808) return the six-octet string
// useful for a BACnet address.
func PackIpAddr(addrTuple *AddressTuple[string, uint16]) (octetString []byte) {
	addr, port := addrTuple.Left, addrTuple.Right
	octetString = append(net.ParseIP(addr).To4(), uint16ToPort(port)...)
	return
}

// UnpackIpAddr Given a six-octet BACnet address, return an IP address tuple.
func UnpackIpAddr(addr []byte) (addrTuple *AddressTuple[string, uint16]) {
	ip := ipv4ToUint32(addr[:4])
	port := portToUint16(addr[4:])

	return &AddressTuple[string, uint16]{uint32ToIpv4(ip).String(), port}
}

func NewLocalStation(localLog zerolog.Logger, addr any, route *Address) (*Address, error) {
	l := &Address{
		log: localLog,
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
		localLog.Debug().Msg("bytearray")
		l.AddrAddress = addr
		length := uint8(len(addr))
		l.AddrLen = &length
	default:
		return nil, errors.New("integer or byte array required")
	}
	return l, nil
}

func NewRemoteStation(localLog zerolog.Logger, net *uint16, addr any, route *Address) (*Address, error) {
	l := &Address{
		log: localLog,
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
		localLog.Debug().Msg("bytearray")
		l.AddrAddress = addr
		length := uint8(len(addr))
		l.AddrLen = &length
	default:
		return nil, errors.New("integer or byte array required")
	}
	return l, nil
}

func NewLocalBroadcast(route *Address) *Address {
	l := &Address{}
	l.AddrType = LOCAL_BROADCAST_ADDRESS
	l.AddrRoute = route
	return l
}

func NewRemoteBroadcast(net uint16, route *Address) *Address {
	r := &Address{}
	r.AddrType = REMOTE_BROADCAST_ADDRESS
	r.AddrNet = &net
	r.AddrRoute = route
	return r
}

func NewGlobalBroadcast(route *Address) *Address {
	g := &Address{}
	g.AddrType = GLOBAL_BROADCAST_ADDRESS
	g.AddrRoute = route
	return g
}

type PCI interface {
	IPCI
	GetExpectingReply() bool
	GetNetworkPriority() readWriteModel.NPDUNetworkPriority
}

type _PCI struct {
	*__PCI
	expectingReply  bool
	networkPriority readWriteModel.NPDUNetworkPriority
}

var _ PCI = (*_PCI)(nil)

func newPCI(rootMessage spi.Message, pduUserData spi.Message, pduSource *Address, pduDestination *Address, expectingReply bool, networkPriority readWriteModel.NPDUNetworkPriority) *_PCI {
	return &_PCI{
		new__PCI(rootMessage, pduUserData, pduSource, pduDestination),
		expectingReply,
		networkPriority,
	}
}

func (p *_PCI) Update(pci Arg) error {
	if err := p.__PCI.Update(pci); err != nil {
		return errors.Wrap(err, "error updating __PCI")
	}
	switch pci := pci.(type) {
	case PCI:
		p.expectingReply = pci.GetExpectingReply()
		p.networkPriority = pci.GetNetworkPriority()
		return nil
	default:
		return errors.Errorf("invalid PCI type %T", pci)
	}
}

func (p *_PCI) GetExpectingReply() bool {
	return p.expectingReply
}

func (p *_PCI) GetNetworkPriority() readWriteModel.NPDUNetworkPriority {
	return p.networkPriority
}

func (p *_PCI) deepCopy() *_PCI {
	__pci := p.__PCI.deepCopy()
	expectingReply := p.expectingReply
	networkPriority := p.networkPriority // Those are immutable so no copy needed
	return &_PCI{__pci, expectingReply, networkPriority}
}

func (p *_PCI) String() string {
	return fmt.Sprintf("_PCI{%s, expectingReply: %t, networkPriority: %s}", p.__PCI, p.expectingReply, p.networkPriority)
}

type PDUData interface {
	SetPduData([]byte)
	GetPduData() []byte
	Get() (byte, error)
	GetShort() (int16, error)
	GetLong() (int64, error)
	GetData(dlen int) ([]byte, error)
	Put(byte)
	PutData(...byte)
	PutShort(uint16)
	PutLong(uint32)
}

// _PDUData is basically a bridge to spi.Message
type _PDUData struct {
	data []byte
}

var _ PDUData = (*_PDUData)(nil)

func NewPDUData(args Args) PDUData {
	p := &_PDUData{}
	if len(args) == 0 {
		return p
	}
	switch data := args[0].(type) {
	case []byte:
		p.data = make([]byte, len(data))
		copy(p.data, data)
	case PDU:
		otherData := data.GetPduData()
		p.data = make([]byte, len(otherData))
		copy(p.data, otherData)
	case PDUData:
		otherData := data.GetPduData()
		p.data = make([]byte, len(otherData))
		copy(p.data, otherData)
	}
	return p
}

func (d *_PDUData) SetPduData(data []byte) {
	d.data = data
}

func (d *_PDUData) GetPduData() []byte {
	return d.data
}

func (d *_PDUData) Get() (byte, error) {
	if d.data == nil || len(d.data) == 0 {
		return 0, io.EOF
	}
	octet := d.data[0]
	d.data = d.data[1:]
	return octet, nil
}

func (d *_PDUData) GetData(dlen int) ([]byte, error) {
	if len(d.data) < dlen {
		return nil, io.EOF
	}
	data := d.data[:dlen]
	d.data = d.data[dlen:]
	if len(data) == 0 {
		return nil, nil
	}
	return data, nil
}

func (d *_PDUData) GetShort() (int16, error) {
	data, err := d.GetData(2)
	if err != nil {
		return 0, err
	}
	return int16(binary.BigEndian.Uint16(data)), nil
}

func (d *_PDUData) GetLong() (int64, error) {
	data, err := d.GetData(4)
	if err != nil {
		return 0, err
	}
	return int64(binary.BigEndian.Uint64(data)), nil
}

func (d *_PDUData) Put(n byte) {
	d.data = append(d.data, n)
}

func (d *_PDUData) PutData(n ...byte) {
	d.data = append(d.data, n...)
}

func (d *_PDUData) PutShort(n uint16) {
	ba := make([]byte, 2)
	binary.BigEndian.PutUint16(ba, n)
	d.data = append(d.data, ba...)
}

func (d *_PDUData) PutLong(n uint32) {
	ba := make([]byte, 4)
	binary.BigEndian.PutUint32(ba, n)
	d.data = append(d.data, ba...)
}

func (d *_PDUData) deepCopy() *_PDUData {
	copyPDUData := *d
	copyPDUData.data = make([]byte, len(d.data))
	copy(copyPDUData.data, d.data)
	return &copyPDUData
}

func (d *_PDUData) String() string {
	return Btox(d.data, ".")
}

type PDU interface {
	PCI
	PDUData
	DeepCopy() PDU
}

// PDUContract provides a set of functions which can be overwritten by a sub struct
type PDUContract interface {
	GetName() string
}

type _PDU struct {
	PDUContract
	*_PCI
	*_PDUData
}

func NewPDU(pdu spi.Message, pduOptions ...PDUOption) PDU {
	p := &_PDU{
		_PCI: newPCI(pdu, nil, nil, nil, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE),
	}
	p.PDUContract = p
	for _, option := range pduOptions {
		option(p)
	}
	p._PDUData = NewPDUData(NewArgs(pdu)).(*_PDUData)
	return p
}

func NewPDUFromPDUWithNewMessage(pdu PDU, pduUserData spi.Message, pduOptions ...PDUOption) PDU {
	p := &_PDU{
		_PCI: newPCI(pdu, pduUserData, pdu.GetPDUSource(), pdu.GetPDUDestination(), pdu.GetExpectingReply(), pdu.GetNetworkPriority()),
	}
	p.PDUContract = p
	for _, option := range pduOptions {
		option(p)
	}
	p._PDUData = NewPDUData(NewArgs(pdu)).(*_PDUData)
	return p
}

type PDUOption func(pdu *_PDU)

func WithPDUUserData(message spi.Message) PDUOption {
	return func(pdu *_PDU) {
		pdu.rootMessage = message
	}
}

func WithPDUSource(pduSource *Address) PDUOption {
	return func(pdu *_PDU) {
		pdu.pduSource = pduSource
	}
}

func WithPDUDestination(pduDestination *Address) PDUOption {
	return func(pdu *_PDU) {
		pdu.pduDestination = pduDestination
	}
}

func WithPDUExpectingReply(expectingReply bool) PDUOption {
	return func(pdu *_PDU) {
		pdu.expectingReply = expectingReply
	}
}

func WithPDUNetworkPriority(networkPriority readWriteModel.NPDUNetworkPriority) PDUOption {
	return func(pdu *_PDU) {
		pdu.networkPriority = networkPriority
	}
}

func (p *_PDU) GetRootMessage() spi.Message {
	return p.rootMessage
}

func (p *_PDU) deepCopy() *_PDU {
	pduCopy := &_PDU{_PCI: p._PCI.deepCopy(), _PDUData: p._PDUData.deepCopy()}
	pduCopy.PDUContract = pduCopy
	return pduCopy
}

func (p *_PDU) DeepCopy() PDU {
	return p.deepCopy()
}

func (p *_PDU) GetName() string {
	return "PDU"
}

func (p *_PDU) String() string {
	if globals.ExtendedPDUOutput {
		return fmt.Sprintf("_PDU{%s}", p._PCI)
	}
	return fmt.Sprintf("<%s %s -> %s : %s>", p.PDUContract.GetName(), p.GetPDUSource(), p.GetPDUDestination(), p._PDUData)
}
