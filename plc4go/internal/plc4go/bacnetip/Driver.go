/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/internal/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
)

type Driver struct {
	fieldHandler            spi.PlcFieldHandler
	awaitSetupComplete      bool
	awaitDisconnectComplete bool
}

func NewDriver() plc4go.PlcDriver {
	return &Driver{
		fieldHandler:            NewFieldHandler(),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,
	}
}

func (m *Driver) GetProtocolCode() string {
	return "bacnet-ip"
}

func (m *Driver) GetProtocolName() string {
	return "BACnet/IP"
}

func (m *Driver) GetDefaultTransport() string {
	return "udp"
}

func (m *Driver) CheckQuery(query string) error {
	_, err := m.fieldHandler.ParseQuery(query)
	return err
}

func (m *Driver) GetConnection(transportUrl url.URL, transports map[string]transports.Transport, options map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	log.Debug().Stringer("transportUrl", &transportUrl).Msgf("Get connection for transport url with %d transport(s) and %d option(s)", len(transports), len(options))
	// Get an the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		log.Error().Stringer("transportUrl", &transportUrl).Msgf("We couldn't find a transport for scheme %s", transportUrl.Scheme)
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl))
		}()
		return ch
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	options["defaultUdpPort"] = []string{"47808"}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(transportUrl, options)
	if err != nil {
		log.Error().Stringer("transportUrl", &transportUrl).Msgf("We couldn't create a transport instance for port %#v", options["defaultUdpPort"])
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.New("couldn't initialize transport configuration for given transport url "+transportUrl.String()))
		}()
		return ch
	}

	codec := NewMessageCodec(transportInstance)
	log.Debug().Msgf("working with codec %#v", codec)

	// Create the new connection
	connection := NewConnection(codec, m.fieldHandler)
	log.Debug().Msg("created connection, connecting now")
	return connection.Connect()
}

func (m *Driver) SupportsDiscovery() bool {
	return false
}

func (m *Driver) Discover(callback func(event apiModel.PlcDiscoveryEvent), options ...apiModel.WithDiscoveryOption) error {
	panic("implement me")
}
