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

// Package virtualcan is an in-process CAN bus. Every instance naming the same bus sees the frames
// every other instance on it writes, the way nodes on a real bus do; an instance does not see its
// own frames, like a SocketCAN socket without CAN_RAW_RECV_OWN_MSGS.
//
//	canopen:virtualcan://bench
//
// It works on every platform and needs no privileges, which makes it what tests and simulators
// use: a simulated device opens a second instance on the same bus and answers the driver's frames.
// Frames use the record layout of package can.
package virtualcan

import (
	"net/url"
	"sync"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

// TransportCode is the scheme of a connection string using this transport.
const TransportCode = "virtualcan"

type Transport struct {
	log zerolog.Logger
}

var _ transports.Transport = (*Transport)(nil)

func NewTransport(_options ...options.WithOption) *Transport {
	return &Transport{
		log: options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
}

func (m *Transport) GetTransportCode() string {
	return TransportCode
}

func (m *Transport) GetTransportName() string {
	return "In-process virtual CAN bus"
}

func (m *Transport) CreateTransportInstance(transportUrl url.URL, _ map[string][]string, _options ...options.WithOption) (transports.TransportInstance, error) {
	busName := transportUrl.Host
	if busName == "" {
		busName = transportUrl.Opaque
	}
	if busName == "" {
		return nil, errors.New("virtualcan needs a bus name, as in virtualcan://bench")
	}
	return NewTransportInstance(busName, _options...), nil
}

func (m *Transport) Close() error {
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}

// bus is the set of instances connected under one name.
type bus struct {
	mutex   sync.RWMutex
	members map[*TransportInstance]struct{}
}

var (
	busesMutex sync.Mutex
	buses      = map[string]*bus{}
)

func join(name string, member *TransportInstance) *bus {
	busesMutex.Lock()
	defer busesMutex.Unlock()
	b, ok := buses[name]
	if !ok {
		b = &bus{members: map[*TransportInstance]struct{}{}}
		buses[name] = b
	}
	b.mutex.Lock()
	b.members[member] = struct{}{}
	b.mutex.Unlock()
	return b
}

func leave(name string, member *TransportInstance) {
	busesMutex.Lock()
	defer busesMutex.Unlock()
	b, ok := buses[name]
	if !ok {
		return
	}
	b.mutex.Lock()
	delete(b.members, member)
	empty := len(b.members) == 0
	b.mutex.Unlock()
	if empty {
		delete(buses, name)
	}
}

// deliver hands frames to every member except the sender.
func (b *bus) deliver(sender *TransportInstance, records []byte) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	for member := range b.members {
		if member != sender {
			member.receive(records)
		}
	}
}
