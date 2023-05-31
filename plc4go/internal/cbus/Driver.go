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
	"context"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/rs/zerolog"
	"net/url"
	"strconv"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
)

type Driver struct {
	_default.DefaultDriver
	tm                      transactions.RequestTransactionManager
	awaitSetupComplete      bool
	awaitDisconnectComplete bool

	log zerolog.Logger
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	driver := &Driver{
		tm:                      transactions.NewRequestTransactionManager(1, _options...),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,
		log:                     options.ExtractCustomLogger(_options...),
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "c-bus", "Clipsal Bus", "tcp", NewTagHandler())
	return driver
}

func (m *Driver) GetConnectionWithContext(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	m.log.Debug().Stringer("transportUrl", &transportUrl).Msgf("Get connection for transport url with %d transport(s) and %d option(s)", len(transports), len(driverOptions))
	// Get the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		m.log.Error().Stringer("transportUrl", &transportUrl).Msgf("We couldn't find a transport for scheme %s", transportUrl.Scheme)
		return m.reportError(errors.Errorf("couldn't find transport for given transport url %v", transportUrl))
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	driverOptions["defaultTcpPort"] = []string{strconv.FormatUint(uint64(readWriteModel.CBusConstants_CBUSTCPDEFAULTPORT), 10)}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(transportUrl, driverOptions)
	if err != nil {
		m.log.Error().Err(err).Stringer("transportUrl", &transportUrl).Msgf("We couldn't create a transport instance for port %#v", driverOptions["defaultTcpPort"])
		return m.reportError(errors.Wrapf(err, "couldn't initialize transport configuration for given transport url %s", transportUrl.String()))
	}

	configuration, err := ParseFromOptions(m.log, driverOptions)
	if err != nil {
		m.log.Error().Err(err).Msgf("Invalid options")
		return m.reportError(errors.Wrap(err, "Invalid options"))
	}

	codec := NewMessageCodec(transportInstance, options.WithCustomLogger(m.log))
	m.log.Debug().Msgf("working with codec %#v", codec)

	driverContext := NewDriverContext(configuration)
	driverContext.awaitSetupComplete = m.awaitSetupComplete
	driverContext.awaitDisconnectComplete = m.awaitDisconnectComplete

	// Create the new connection
	connection := NewConnection(codec, configuration, driverContext, m.GetPlcTagHandler(), m.tm, driverOptions, options.WithCustomLogger(m.log))
	m.log.Debug().Msg("created connection, connecting now")
	return connection.ConnectWithContext(ctx)
}

func (m *Driver) reportError(err error) <-chan plc4go.PlcConnectionConnectResult {
	ch := make(chan plc4go.PlcConnectionConnectResult, 1)
	ch <- _default.NewDefaultPlcConnectionConnectResult(nil, err)
	return ch
}

func (m *Driver) SetAwaitSetupComplete(awaitComplete bool) {
	m.awaitSetupComplete = awaitComplete
}

func (m *Driver) SetAwaitDisconnectComplete(awaitComplete bool) {
	m.awaitDisconnectComplete = awaitComplete
}

func (m *Driver) SupportsDiscovery() bool {
	return true
}

func (m *Driver) DiscoverWithContext(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return NewDiscoverer(options.WithCustomLogger(m.log)).Discover(ctx, callback, discoveryOptions...)
}
