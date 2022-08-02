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

package cbus

import (
	"github.com/apache/plc4x/plc4go/internal/spi"
	_default "github.com/apache/plc4x/plc4go/internal/spi/default"
	"github.com/apache/plc4x/plc4go/internal/spi/transports"
	"github.com/apache/plc4x/plc4go/pkg/api"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
	"strconv"
)

type Driver struct {
	_default.DefaultDriver
	tm                      spi.RequestTransactionManager
	awaitSetupComplete      bool
	awaitDisconnectComplete bool
}

func NewDriver() plc4go.PlcDriver {
	return &Driver{
		DefaultDriver:           _default.NewDefaultDriver("c-bus", "Clipsal Bus", "tcp", NewFieldHandler()),
		tm:                      spi.NewRequestTransactionManager(1),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,
	}
}

func (m *Driver) GetConnection(transportUrl url.URL, transports map[string]transports.Transport, options map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	log.Debug().Stringer("transportUrl", &transportUrl).Msgf("Get connection for transport url with %d transport(s) and %d option(s)", len(transports), len(options))
	// Get the transport specified in the url
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
	options["defaultTcpPort"] = []string{strconv.FormatUint(uint64(readWriteModel.CBusConstants_CBUSTCPDEFAULTPORT), 10)}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(transportUrl, options)
	if err != nil {
		log.Error().Stringer("transportUrl", &transportUrl).Msgf("We couldn't create a transport instance for port %#v", options["defaultTcpPort"])
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.New("couldn't initialize transport configuration for given transport url "+transportUrl.String()))
		}()
		return ch
	}

	configuration, err := ParseFromOptions(options)
	if err != nil {
		log.Error().Err(err).Msgf("Invalid options")
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "Invalid options"))
		}()
		return ch
	}

	codec := NewMessageCodec(transportInstance, configuration.srchk)
	log.Debug().Msgf("working with codec %#v", codec)

	driverContext, err := NewDriverContext(configuration)
	if err != nil {
		log.Error().Err(err).Msgf("Invalid options")
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "Invalid options"))
		}()
		return ch
	}
	driverContext.awaitSetupComplete = m.awaitSetupComplete
	driverContext.awaitDisconnectComplete = m.awaitDisconnectComplete

	// Create the new connection
	connection := NewConnection(codec, configuration, driverContext, m.GetPlcFieldHandler(), &m.tm, options)
	log.Debug().Msg("created connection, connecting now")
	return connection.Connect()
}

func (m *Driver) SetAwaitSetupComplete(awaitComplete bool) {
	m.awaitSetupComplete = awaitComplete
}

func (m *Driver) SetAwaitDisconnectComplete(awaitComplete bool) {
	m.awaitDisconnectComplete = awaitComplete
}
