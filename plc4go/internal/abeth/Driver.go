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

package abeth

import (
	"context"
	"net/url"
	"slices"
	"strconv"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

const (
	// defaultTransportCode is the transport a connection string without a scheme gets, mirroring
	// plc4j AbEthDriver.getDefaultTransportCode.
	defaultTransportCode = "tcp"
	// AbEthPort is the TCP port an Allen-Bradley ETH interface listens on, mirroring plc4j's
	// AbEthDriver.AB_ETH_PORT. Despite the protocol's name this is plain TCP - no raw ethernet
	// access is needed anywhere.
	AbEthPort = 2222
	// maxConcurrentRequests is what plc4j's AbEthConnection.getMaxConcurrentRequests returns: the
	// PLC answers one request at a time.
	maxConcurrentRequests = 1
)

// supportedTransportCodes are the transports this driver speaks, mirroring plc4j
// AbEthDriver.getSupportedTransportCodes:
//
//	tcp:  the real thing
//	test: in-process loopback for the driver testsuite framework
var supportedTransportCodes = []string{"tcp", "test"}

type Driver struct {
	_default.DefaultDriver

	tm                      transactions.RequestTransactionManager
	awaitSetupComplete      bool
	awaitDisconnectComplete bool

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	driver := &Driver{
		tm:                      transactions.NewRequestTransactionManager(maxConcurrentRequests, _options...),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,

		log:      customLogger,
		_options: _options,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "ab-eth", "Allen Bradley ETH", defaultTransportCode, NewTagHandler())
	return driver
}

// GetSupportedTransportCodes is the list of transports an ab-eth connection string may name.
func (d *Driver) GetSupportedTransportCodes() []string {
	return append([]string(nil), supportedTransportCodes...)
}

func (d *Driver) GetConnection(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) (plc4go.PlcConnection, error) {
	connectionLog := d.log.With().Ctx(ctx).Str("transportUrl", transportUrl.String()).Logger()
	connectionLog.Debug().
		Int("nTransports", len(transports)).
		Int("nDriverOptions", len(driverOptions)).
		Msg("Get connection for transport url with nTransports transport(s) and nDriverOptions option(s)")

	// CIP encapsulation is a stream protocol, which rules out the datagram transports; refusing
	// them here rather than at the first read keeps a mistyped connection string from looking like a
	// PLC which never answers.
	if !slices.Contains(supportedTransportCodes, transportUrl.Scheme) {
		connectionLog.Error().
			Str("scheme", transportUrl.Scheme).
			Strs("supportedTransportCodes", supportedTransportCodes).
			Msg("The ab-eth driver doesn't support this transport")
		return nil, errors.Errorf("the ab-eth driver doesn't support the transport %s, it supports %v", transportUrl.Scheme, supportedTransportCodes)
	}
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		connectionLog.Error().
			Stringer("transportUrl", &transportUrl).
			Str("scheme", transportUrl.Scheme).
			Msg("We couldn't find a transport for scheme")
		return nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl)
	}
	// Provide a default-port to the transport, which is used if the user doesn't provide one in the
	// connection string.
	driverOptions["defaultTcpPort"] = []string{strconv.Itoa(AbEthPort)}
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

	configuration, err := ParseFromOptions(connectionLog, driverOptions)
	if err != nil {
		connectionLog.Error().Err(err).Msg("Invalid driverOptions")
		return nil, errors.Wrap(err, "Invalid driverOptions")
	}

	codec := NewMessageCodec(
		transportInstance,
		append(d._options, options.WithCustomLogger(connectionLog))...,
	)
	connectionLog.Debug().Interface("codec", codec).Msg("working with codec")

	driverContext, err := NewDriverContext(configuration)
	if err != nil {
		connectionLog.Error().Err(err).Msg("Invalid driverOptions")
		return nil, errors.Wrap(err, "Invalid driverOptions")
	}
	driverContext.awaitSetupComplete = d.awaitSetupComplete
	driverContext.awaitDisconnectComplete = d.awaitDisconnectComplete

	connection := NewConnection(
		codec,
		configuration,
		driverContext,
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

func (d *Driver) SetAwaitSetupComplete(awaitComplete bool) {
	d.awaitSetupComplete = awaitComplete
}

func (d *Driver) SetAwaitDisconnectComplete(awaitComplete bool) {
	d.awaitDisconnectComplete = awaitComplete
}

func (d *Driver) Close() error {
	defer utils.StopWarn(d.log)()
	if err := d.tm.Close(); err != nil {
		return errors.Wrap(err, "error closing transaction manager")
	}
	return nil
}
