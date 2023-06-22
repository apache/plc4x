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

package udp

import (
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net"
	"net/url"
	"regexp"
	"strconv"
)

type Transport struct {
	log zerolog.Logger
}

func NewTransport(_options ...options.WithOption) *Transport {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Transport{
		log: customLogger,
	}
}

func (m *Transport) GetTransportCode() string {
	return "udp"
}

func (m *Transport) GetTransportName() string {
	return "UDP Datagram Transport"
}

func (m *Transport) CreateTransportInstance(transportUrl url.URL, options map[string][]string, _options ...options.WithOption) (transports.TransportInstance, error) {
	return m.CreateTransportInstanceForLocalAddress(transportUrl, options, nil, _options...)
}

func (m *Transport) CreateTransportInstanceForLocalAddress(transportUrl url.URL, options map[string][]string, localAddress *net.UDPAddr, _options ...options.WithOption) (transports.TransportInstance, error) {
	connectionStringRegexp := regexp.MustCompile(`^((?P<ip>[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3})|(?P<hostname>[a-zA-Z0-9.\-]+))(:(?P<port>[0-9]{1,5}))?`)
	var remoteAddressString string
	var remotePort int
	if match := utils.GetSubgroupMatches(connectionStringRegexp, transportUrl.Host); match != nil {
		if val, ok := match["ip"]; ok && len(val) > 0 {
			remoteAddressString = val
		} else if val, ok := match["hostname"]; ok && len(val) > 0 {
			remoteAddressString = val
		} else {
			return nil, errors.New("missing hostname or ip to connect")
		}

		if val, ok := match["port"]; ok && len(val) > 0 {
			portVal, err := strconv.Atoi(val)
			if err != nil {
				return nil, errors.Wrap(err, "error setting port")
			}
			remotePort = portVal
		} else if val, ok := options["defaultUdpPort"]; ok && len(val) > 0 {
			portVal, err := strconv.Atoi(val[0])
			if err != nil {
				return nil, errors.Wrap(err, "error setting default udp port")
			}
			remotePort = portVal
		} else {
			return nil, errors.New("error setting port. No explicit or default port provided")
		}
	}
	var connectTimeout uint32 = 1000
	if val, ok := options["connect-timeout"]; ok {
		if parsedConnectTimeout, err := strconv.ParseUint(val[0], 10, 32); err != nil {
			return nil, errors.Wrap(err, "error setting connect-timeout")
		} else {
			connectTimeout = uint32(parsedConnectTimeout)
		}
	}

	var soReUse bool
	if val, ok := options["so-reuse"]; ok {
		if parseBool, err := strconv.ParseBool(val[0]); err != nil {
			return nil, errors.Wrap(err, "error setting so-reuse")
		} else {
			soReUse = parseBool
		}
	}

	// Potentially resolve the ip address, if a hostname was provided
	remoteAddress, err := net.ResolveUDPAddr("udp", remoteAddressString+":"+strconv.Itoa(remotePort))
	if err != nil {
		return nil, errors.Wrap(err, "error resolving typ address")
	}

	return NewTransportInstance(localAddress, remoteAddress, connectTimeout, soReUse, m, _options...), nil
}

func (m *Transport) Close() error {
	m.log.Trace().Msg("Closing")
	return nil
}

func (m *Transport) String() string {
	return m.GetTransportCode() + "(" + m.GetTransportName() + ")"
}
