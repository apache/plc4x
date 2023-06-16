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
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net/url"
	"strconv"
)

type TransportType string

const (
	UDP  TransportType = "udp"
	TCP  TransportType = "tcp"
	PCAP TransportType = "pcap"
)

type Transport struct {
	log zerolog.Logger
}

func NewTransport(_options ...options.WithOption) *Transport {
	return &Transport{
		log: options.ExtractCustomLogger(_options...),
	}
}

func (m *Transport) GetTransportCode() string {
	return "pcap"
}

func (m *Transport) GetTransportName() string {
	return "PCAP(NG) Playback Transport"
}

func (m *Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string, _options ...options.WithOption) (transports.TransportInstance, error) {
	var transportType = PCAP
	if val, ok := options["transport-type"]; ok {
		transportType = TransportType(val[0])
	}
	var portRange = ""
	if val, ok := options["transport-port-range"]; ok {
		portRange = val[0]
	}
	var speedFactor float32 = 1.0
	if val, ok := options["speed-factor"]; ok {
		if parsedSpeedFactory, err := strconv.ParseFloat(val[0], 32); err != nil {
			return nil, errors.Wrap(err, "error parsing speed-factor")
		} else {
			speedFactor = float32(parsedSpeedFactory)
		}
	}

	return NewPcapTransportInstance(transportUrl.Path, transportType, portRange, speedFactor, m, _options...), nil
}

func (m *Transport) Close() error {
	m.log.Trace().Msg("Closing")
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}
