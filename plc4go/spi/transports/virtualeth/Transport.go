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

// Package virtualeth is an in-process Ethernet segment, the raw-socket transport's counterpart for
// tests and simulators. Every instance on the same segment sees the frames the others write (not its
// own), filtered by the same options as raw-socket, and frames use raw-socket's record layout.
//
//	virtual-eth://bench
//
// Each instance gets a locally administered MAC address of its own, or the one given with the
// virtual-eth.mac option.
package virtualeth

import (
	"crypto/rand"
	"net"
	"net/url"
	"sync"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
)

const (
	// TransportCode is the scheme of a connection string using this transport.
	TransportCode = "virtual-eth"
	// MACOption sets the instance's MAC address.
	MACOption = "virtual-eth.mac"
)

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
	return "In-process virtual Ethernet segment"
}

func (m *Transport) CreateTransportInstance(transportUrl url.URL, transportOptions map[string][]string, _options ...options.WithOption) (transports.TransportInstance, error) {
	segment := transportUrl.Host
	if segment == "" {
		segment = transportUrl.Opaque
	}
	if segment == "" {
		return nil, errors.New("virtual-eth needs a segment name, as in virtual-eth://bench")
	}
	filter, err := rawsocket.ParseFilter(transportOptions)
	if err != nil {
		return nil, err
	}
	var mac net.HardwareAddr
	if values := transportOptions[MACOption]; len(values) > 0 {
		if mac, err = net.ParseMAC(values[0]); err != nil || len(mac) != 6 {
			return nil, errors.Errorf("bad %s %q", MACOption, values[0])
		}
	}
	return NewTransportInstance(segment, filter, mac, _options...), nil
}

func (m *Transport) Close() error {
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}

// RandomMAC is a locally administered unicast address.
func RandomMAC() net.HardwareAddr {
	mac := make(net.HardwareAddr, 6)
	_, _ = rand.Read(mac)
	mac[0] = mac[0]&0xFC | 0x02
	return mac
}

type segment struct {
	mutex   sync.RWMutex
	members map[*TransportInstance]struct{}
}

var (
	segmentsMutex sync.Mutex
	segments      = map[string]*segment{}
)

func join(name string, member *TransportInstance) *segment {
	segmentsMutex.Lock()
	defer segmentsMutex.Unlock()
	s, ok := segments[name]
	if !ok {
		s = &segment{members: map[*TransportInstance]struct{}{}}
		segments[name] = s
	}
	s.mutex.Lock()
	s.members[member] = struct{}{}
	s.mutex.Unlock()
	return s
}

func leave(name string, member *TransportInstance) {
	segmentsMutex.Lock()
	defer segmentsMutex.Unlock()
	s, ok := segments[name]
	if !ok {
		return
	}
	s.mutex.Lock()
	delete(s.members, member)
	empty := len(s.members) == 0
	s.mutex.Unlock()
	if empty {
		delete(segments, name)
	}
}

func (s *segment) deliver(sender *TransportInstance, frame []byte) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for member := range s.members {
		if member != sender {
			member.receive(frame)
		}
	}
}
