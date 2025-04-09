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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/ads/model"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	adsModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Driver struct {
	_default.DefaultDriver

	discoverer *Discoverer

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	driver := &Driver{
		discoverer: NewDiscoverer(_options...),
		log:        customLogger,
		_options:   _options,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "ads", "Beckhoff TwinCat ADS", "tcp", NewTagHandler())
	return driver
}

func (d *Driver) GetConnectionWithContext(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	d.log.Debug().
		Stringer("transportUrl", &transportUrl).
		Int("nTransports", len(transports)).
		Int("nDriverOptions", len(driverOptions)).
		Msg("Get connection for transport url with nTransports transport(s) and nDriverOptions option(s)")
	// Get the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		d.log.Error().
			Stringer("transportUrl", &transportUrl).
			Str("scheme", transportUrl.Scheme).
			Msg("We couldn't find a transport for scheme")
		ch := make(chan plc4go.PlcConnectionConnectResult, 1)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl))
		return ch
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	driverOptions["defaultTcpPort"] = []string{strconv.Itoa(int(adsModel.AdsConstants_ADSTCPDEFAULTPORT))}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(
		transportUrl,
		driverOptions,
		append(d._options, options.WithCustomLogger(d.log))...,
	)
	if err != nil {
		d.log.Error().
			Stringer("transportUrl", &transportUrl).
			Strs("defaultTcpPort", driverOptions["defaultTcpPort"]).
			Msg("We couldn't create a transport instance for port")
		ch := make(chan plc4go.PlcConnectionConnectResult, 1)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.New("couldn't initialize transport configuration for given transport url "+transportUrl.String()))
		return ch
	}

	// Create a new codec for taking care of encoding/decoding of messages
	codec := NewMessageCodec(
		transportInstance,
		append(d._options, options.WithCustomLogger(d.log))...,
	)
	d.log.Debug().Stringer("codec", codec).Msg("working with codec")

	configuration, err := model.ParseFromOptions(d.log, driverOptions)
	if err != nil {
		d.log.Error().Err(err).Msg("Invalid driverOptions")
		ch := make(chan plc4go.PlcConnectionConnectResult, 1)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "invalid configuration"))
		return ch
	}

	// Create the new connection
	connection, err := NewConnection(codec, configuration, driverOptions)
	if err != nil {
		ch := make(chan plc4go.PlcConnectionConnectResult, 1)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "couldn't create connection"))
		return ch
	}
	d.log.Debug().Stringer("connection", connection).Msg("created connection, connecting now")
	return connection.ConnectWithContext(ctx)
}

func (d *Driver) SupportsDiscovery() bool {
	return true
}

func (d *Driver) DiscoverWithContext(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return d.discoverer.Discover(ctx, callback, discoveryOptions...)
}

func (d *Driver) Close() error {
	defer utils.StopWarn(d.log)()
	d.log.Trace().Msg("Closing driver")
	d.log.Trace().Msg("Closing discoverer")
	return d.discoverer.Close()
}
