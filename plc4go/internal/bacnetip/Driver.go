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

package bacnetip

import (
	"context"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Driver struct {
	_default.DefaultDriver

	discoverer              *Discoverer
	tm                      transactions.RequestTransactionManager
	awaitSetupComplete      bool
	awaitDisconnectComplete bool

	_options []options.WithOption
	log      zerolog.Logger
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	driver := &Driver{
		discoverer:              NewDiscoverer(_options...),
		tm:                      transactions.NewRequestTransactionManager(math.MaxInt),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,

		_options: _options,
		log:      customLogger,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "bacnet-ip", "BACnet/IP", "udp", NewTagHandler())
	return driver
}

func (d *Driver) GetConnection(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) (plc4go.PlcConnection, error) {
	connectionLog := d.log.With().Ctx(ctx).Str("transportUrl", transportUrl.String()).Logger()
	connectionLog.Debug().
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
		return nil, errors.Errorf("couldn't find transport for given transport url %v", transportUrl)
	}
	// Provide a default-port to the transport, used if the user doesn't provide one in the connection string.
	driverOptions["defaultUdpPort"] = []string{strconv.FormatUint(uint64(model.BacnetConstants_BACNETUDPDEFAULTPORT), 10)}
	// Set so_reuse by default so multiple BACnet processes can share the BACnet/IP UDP port.
	if _, ok := driverOptions["so-reuse"]; !ok {
		driverOptions["so-reuse"] = []string{"true"}
	}
	// Have the transport create a new transport-instance.
	transportInstance, err := transport.CreateTransportInstance(
		transportUrl,
		driverOptions,
		append(d._options, options.WithCustomLogger(connectionLog))...,
	)
	if err != nil {
		connectionLog.Error().
			Stringer("transportUrl", &transportUrl).
			Strs("defaultUdpPort", driverOptions["defaultUdpPort"]).
			Msg("We couldn't create a transport instance for port")
		return nil, errors.Wrapf(err, "couldn't initialize transport configuration for given transport url %s", transportUrl.String())
	}

	// Parse Configuration early so we can propagate the discovery timeout to
	// the Discoverer for the lifetime of this Driver instance. (Discovery is a
	// driver-level call, but the timeout is naturally part of Configuration.)
	if cfg, cfgErr := ParseFromOptions(connectionLog, driverOptions); cfgErr == nil {
		d.discoverer.SetDiscoveryTimeout(time.Duration(cfg.DiscoveryTimeoutSeconds) * time.Second)
	}

	codec := NewMessageCodec(
		transportInstance,
		append(d._options, options.WithCustomLogger(connectionLog))...,
	)
	connectionLog.Debug().Interface("codec", codec).Msg("working with codec")

	// Create the new connection
	connection := NewConnection(
		codec,
		d.GetPlcTagHandler(),
		d.tm,
		driverOptions,
		append(d._options, options.WithCustomLogger(connectionLog))...,
	)
	connectionLog.Debug().Msg("created connection, connecting now")
	if err := connection.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "Error connecting connection")
	}
	return connection, nil
}

func (d *Driver) SupportsDiscovery() bool {
	return true
}

func (d *Driver) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return d.discoverer.Discover(ctx, callback, discoveryOptions...)
}

func (d *Driver) SetAwaitSetupComplete(awaitComplete bool) {
	d.awaitSetupComplete = awaitComplete
}

func (d *Driver) SetAwaitDisconnectComplete(awaitComplete bool) {
	d.awaitDisconnectComplete = awaitComplete
}

func (d *Driver) Close() error {
	defer utils.StopWarn(d.log)()
	d.log.Trace().Msg("Closing driver")
	var collectedErrors []error
	d.log.Trace().Msg("Closing discoverer")
	if err := d.discoverer.Close(); err != nil {
		collectedErrors = append(collectedErrors, errors.Wrap(err, "failed to close discoverer"))
	}
	d.log.Trace().Msg("Closing transaction manager")
	if err := d.tm.Close(); err != nil {
		collectedErrors = append(collectedErrors, errors.Wrap(err, "error closing transaction manager"))
	}
	if err := errors.Join(collectedErrors...); err != nil {
		return errors.Wrap(err, "error closing driver")
	}
	return nil
}
