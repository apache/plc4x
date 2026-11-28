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

// Package rawsocket is a raw Ethernet transport over Linux AF_PACKET sockets, in pure Go (no cgo,
// no libpcap), so it works in CGO_ENABLED=0 builds. It is what Profinet runs on.
//
//	raw-socket://eth0
//
// Because Ethernet frames vary in length, the instance delivers them as records on its byte
// stream: a 2 byte big-endian length, then the frame from the destination MAC address on (no
// preamble, no FCS). Writes take the same records. The kernel filters what arrives, from these
// options (normally set by the driver):
//
//	raw-socket.ethertypes  comma separated EtherTypes, e.g. 0x8892,0x88cc
//	raw-socket.udp-ports   comma separated UDP ports of IPv4 datagrams, e.g. 34964
//
// Without either, every frame on the interface arrives. The process needs CAP_NET_RAW. On anything
// but Linux the transport is still registered and fails at connect time with a clear error.
package rawsocket

import (
	"encoding/binary"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

const (
	// TransportCode is the scheme of a connection string using this transport, the code plc4j's
	// raw socket transport uses.
	TransportCode = "raw-socket"

	// EtherTypesOption and UDPPortsOption are the filter options.
	EtherTypesOption = "raw-socket.ethertypes"
	UDPPortsOption   = "raw-socket.udp-ports"

	// RecordHeaderSize is the length prefix of a frame record.
	RecordHeaderSize = 2
	// MaxFrameSize is the largest frame accepted: a VLAN tagged 1500 byte payload. Jumbo frames are
	// no use to the protocols this serves.
	MaxFrameSize = 1522
	// MinFrameSize is a destination, a source and an EtherType.
	MinFrameSize = 14
)

var interfaceNamePattern = regexp.MustCompile(`^[A-Za-z0-9_.\-:@]{1,15}$`)

// EncodeRecord prefixes a frame with its length.
func EncodeRecord(frame []byte) ([]byte, error) {
	if len(frame) < MinFrameSize || len(frame) > MaxFrameSize {
		return nil, errors.Errorf("an Ethernet frame is %d to %d bytes, got %d", MinFrameSize, MaxFrameSize, len(frame))
	}
	record := make([]byte, RecordHeaderSize+len(frame))
	binary.BigEndian.PutUint16(record, uint16(len(frame)))
	copy(record[RecordHeaderSize:], frame)
	return record, nil
}

// RecordFrameLength reads the length prefix of a record, checking it is a plausible frame.
func RecordFrameLength(header []byte) (int, error) {
	if len(header) < RecordHeaderSize {
		return 0, errors.New("short record header")
	}
	length := int(binary.BigEndian.Uint16(header))
	if length < MinFrameSize || length > MaxFrameSize {
		return 0, errors.Errorf("record announces a %d byte frame", length)
	}
	return length, nil
}

type Transport struct {
	log zerolog.Logger
}

var _ transports.Transport = (*Transport)(nil)

func NewTransport(_options ...options.WithOption) *Transport {
	return &Transport{log: options.ExtractCustomLoggerOrDefaultToGlobal(_options...)}
}

func (m *Transport) GetTransportCode() string {
	return TransportCode
}

func (m *Transport) GetTransportName() string {
	return "Linux raw Ethernet (AF_PACKET) Transport"
}

func (m *Transport) CreateTransportInstance(transportUrl url.URL, transportOptions map[string][]string, _options ...options.WithOption) (transports.TransportInstance, error) {
	interfaceName := transportUrl.Host
	if interfaceName == "" {
		interfaceName = transportUrl.Opaque
	}
	if !interfaceNamePattern.MatchString(interfaceName) {
		return nil, errors.Errorf("raw-socket needs a network interface name such as eth0, got %q", interfaceName)
	}
	filter, err := ParseFilter(transportOptions)
	if err != nil {
		return nil, err
	}
	return NewTransportInstance(interfaceName, filter, m, _options...), nil
}

// ParseFilter reads the filter options.
func ParseFilter(transportOptions map[string][]string) (Filter, error) {
	var filter Filter
	for _, value := range transportOptions[EtherTypesOption] {
		for _, field := range strings.Split(value, ",") {
			etherType, err := parseUint16(field)
			if err != nil {
				return Filter{}, errors.Wrapf(err, "bad EtherType %q in %s", field, EtherTypesOption)
			}
			filter.EtherTypes = append(filter.EtherTypes, etherType)
		}
	}
	for _, value := range transportOptions[UDPPortsOption] {
		for _, field := range strings.Split(value, ",") {
			port, err := parseUint16(field)
			if err != nil {
				return Filter{}, errors.Wrapf(err, "bad UDP port %q in %s", field, UDPPortsOption)
			}
			filter.UDPPorts = append(filter.UDPPorts, port)
		}
	}
	return filter, nil
}

func parseUint16(text string) (uint16, error) {
	text = strings.TrimSpace(text)
	value, err := strconv.ParseUint(text, 0, 16)
	return uint16(value), err
}

func (m *Transport) Close() error {
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}
