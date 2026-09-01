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

package ads

import (
	"context"
	"net/url"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	adsModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Driver struct {
	_default.DefaultDriver

	discoverer *Discoverer

	awaitSetupComplete      bool
	awaitDisconnectComplete bool

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	driver := &Driver{
		discoverer:              NewDiscoverer(_options...),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,
		log:                     customLogger,
		_options:                _options,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "ads", "Beckhoff TwinCat ADS", "tcp", NewTagHandler())
	return driver
}

func (d *Driver) GetConnection(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) (plc4go.PlcConnection, error) {
	connectionLog := d.log.With().Ctx(ctx).Str("transportUrl", transportUrl.String()).Logger()
	connectionLog.Debug().
		Stringer("transportUrl", &transportUrl).
		Int("nTransports", len(transports)).
		Int("nDriverOptions", len(driverOptions)).
		Msg("Get connection for transport url with nTransports transport(s) and nDriverOptions option(s)")
	// Get the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		connectionLog.Error().
			Stringer("transportUrl", &transportUrl).
			Str("scheme", transportUrl.Scheme).
			Msg("We couldn't find a transport for scheme")
		return nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl)
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	driverOptions["defaultTcpPort"] = []string{strconv.Itoa(int(adsModel.AdsConstants_ADSTCPDEFAULTPORT))}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(
		transportUrl,
		driverOptions,
		append(d._options, options.WithCustomLogger(connectionLog))...,
	)
	if err != nil {
		connectionLog.Error().
			Stringer("transportUrl", &transportUrl).
			Strs("defaultTcpPort", driverOptions["defaultTcpPort"]).
			Msg("We couldn't create a transport instance for port")
		return nil, errors.New("couldn't initialize transport configuration for given transport url " + transportUrl.String())
	}

	// Create a new codec for taking care of encoding/decoding of messages
	codec := NewMessageCodec(
		transportInstance,
		append(d._options, options.WithCustomLogger(connectionLog))...,
	)
	connectionLog.Debug().Interface("codec", codec).Msg("working with codec")

	configuration, err := model.ParseFromOptions(connectionLog, driverOptions)
	if err != nil {
		connectionLog.Error().Err(err).Msg("Invalid driverOptions")
		return nil, errors.Wrap(err, "invalid configuration")
	}
	// TODO: check if we want to attach the configuration to the logger

	// Create the new connection
	connection, err := NewConnection(codec, configuration, driverOptions)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create connection")
	}
	connection.driverContext.awaitSetupComplete = d.awaitSetupComplete
	connection.driverContext.awaitDisconnectComplete = d.awaitDisconnectComplete
	connectionLog.Debug().Interface("connection", connection).Msg("created connection, connecting now")
	if err := connection.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "Error connecting connection")
	}
	return connection, nil
}

// SetAwaitSetupComplete lets the driver testsuite runner return from Connect before the
// connection handshake ran, so the testsuite can feed the handshake's responses itself.
func (d *Driver) SetAwaitSetupComplete(awaitComplete bool) {
	d.awaitSetupComplete = awaitComplete
}

// SetAwaitDisconnectComplete is the disconnect-side counterpart of SetAwaitSetupComplete.
func (d *Driver) SetAwaitDisconnectComplete(awaitComplete bool) {
	d.awaitDisconnectComplete = awaitComplete
}

func (d *Driver) SupportsDiscovery() bool {
	return true
}

func (d *Driver) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return d.discoverer.Discover(ctx, callback, discoveryOptions...)
}

func (d *Driver) Close() error {
	defer utils.StopWarn(d.log)()
	d.log.Trace().Msg("Closing driver")
	d.log.Trace().Msg("Closing discoverer")
	return d.discoverer.Close()
}
