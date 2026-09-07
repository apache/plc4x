//go:build !cgo && !windows

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

package pcap

import (
	"net/url"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

// The real pcap transport binds to libpcap through gopacket/pcap, which is
// cgo-only on every platform except Windows. This stub keeps the package
// (and everything importing it) compilable with CGO_ENABLED=0 — e.g. when
// cross-compiling — and reports a clear error if the transport is actually
// used in such a build.

type Transport struct {
	log zerolog.Logger
}

var _ transports.Transport = (*Transport)(nil)

func NewTransport(_options ...options.WithOption) *Transport {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Transport{
		log: customLogger,
	}
}

func (m *Transport) GetTransportCode() string {
	return "pcap"
}

func (m *Transport) GetTransportName() string {
	return "PCAP(NG) Playback Transport"
}

func (m *Transport) CreateTransportInstance(_ url.URL, _ map[string][]string, _ ...options.WithOption) (transports.TransportInstance, error) {
	return nil, errors.New("the pcap transport requires cgo (libpcap); rebuild with CGO_ENABLED=1")
}

func (m *Transport) Close() error {
	m.log.Trace().Msg("Closing")
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}
