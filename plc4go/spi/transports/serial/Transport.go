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

package serial

import (
	"net"
	"net/url"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

type Transport struct {
	registry *sharedPortRegistry

	log zerolog.Logger
}

func NewTransport(_options ...options.WithOption) *Transport {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Transport{
		registry: newSharedPortRegistry(customLogger),
		log:      customLogger,
	}
}

func (m *Transport) GetTransportCode() string {
	return "serial"
}

func (m *Transport) GetTransportName() string {
	return "Serial Transport"
}

func (m *Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string, _options ...options.WithOption) (transports.TransportInstance, error) {
	return m.CreateTransportInstanceForLocalAddress(transportUrl, options, nil, _options...)
}

func (m *Transport) CreateTransportInstanceForLocalAddress(transportUrl url.URL, options map[string][]string, _ *net.UDPAddr, _options ...options.WithOption) (transports.TransportInstance, error) {
	var serialPortName = transportUrl.Path

	cfg, err := parseSerialOptions(options)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing serial options")
	}

	return NewTransportInstanceWithConfig(serialPortName, cfg, m, _options...), nil
}

func (m *Transport) Close() error {
	m.log.Trace().Msg("Closing")
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}
