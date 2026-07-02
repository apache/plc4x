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
	"net"
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
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
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
		tm:                      transactions.NewRequestTransactionManager(math.MaxInt, options.WithCustomLogger(customLogger)),
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
	// BACnet/IP uses port 47808 on both sides of a conversation; spec-conformant
	// peers (bacpypes3, EcoStruxure, Niagara, ...) send unsolicited messages
	// and responses back to the well-known port regardless of the request's
	// source port. The generic transport.CreateTransportInstance dials with
	// LocalAddress=nil, which gives us an ephemeral source — fine for protocols
	// that reply to the source port, but for BACnet that means responses get
	// dropped by the kernel.
	//
	// Use CreateTransportInstanceForLocalAddress with a fixed 0.0.0.0:47808
	// bind. Callers that need to co-locate multiple BACnet connections in one
	// process can override via the "local-port" driver option (uint), or 0
	// for explicit ephemeral.
	localPort := int(model.BacnetConstants_BACNETUDPDEFAULTPORT)
	if val, ok := driverOptions["local-port"]; ok && len(val) > 0 {
		if parsed, parseErr := strconv.Atoi(val[0]); parseErr != nil {
			connectionLog.Warn().Err(parseErr).Str("local-port", val[0]).Msg("ignoring invalid local-port option")
		} else {
			localPort = parsed
		}
	}
	localAddress := &net.UDPAddr{IP: net.IPv4zero, Port: localPort}
	connectionLog.Info().Stringer("localAddress", localAddress).Msg("BACnet driver binding local UDP")

	udpTransport, ok := transport.(*udp.Transport)
	if !ok {
		return nil, errors.Errorf("BACnet/IP requires the udp transport; got %T", transport)
	}
	transportInstance, err := udpTransport.CreateTransportInstanceForLocalAddress(
		transportUrl,
		driverOptions,
		localAddress,
		append(d._options, options.WithCustomLogger(connectionLog))...,
	)
	if err != nil {
		connectionLog.Error().
			Stringer("transportUrl", &transportUrl).
			Strs("defaultUdpPort", driverOptions["defaultUdpPort"]).
			Int("localPort", localPort).
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
