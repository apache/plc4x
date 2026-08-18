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

package eip

import (
	"context"
	"encoding/binary"
	"net/url"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

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
		tm:                      transactions.NewRequestTransactionManager(1, _options...),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,

		log:      customLogger,
		_options: _options,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "eip", "EthernetIP", "tcp", NewTagHandler())
	return driver
}

func (d *Driver) GetConnection(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) (plc4go.PlcConnection, error) {
	connectionLog := d.log.With().Ctx(ctx).Str("transportUrl", transportUrl.String()).Logger()
	connectionLog.Debug().
		Int("nTransports", len(transports)).
		Int("nDriverOptions", len(driverOptions)).
		Msg("Get connection for transport url with nTransports transport(s) and nDriverOptions option(s)")
	// Get an the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		connectionLog.Error().
			Stringer("transportUrl", &transportUrl).
			Str("scheme", transportUrl.Scheme).
			Msg("We couldn't find a transport for scheme")
		return nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl)
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	driverOptions["defaultTcpPort"] = []string{"44818"}
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

	configuration, err := ParseFromOptions(connectionLog, driverOptions)
	if err != nil {
		connectionLog.Error().Err(err).Msg("Invalid driverOptions")
		return nil, errors.Wrap(err, "Invalid driverOptions")
	}

	byteOrder := binary.ByteOrder(binary.BigEndian)
	if !configuration.bigEndian {
		byteOrder = binary.LittleEndian
	}
	codec := NewMessageCodec(
		transportInstance,
		byteOrder,
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

	// Create the new connection
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
	return d.tm.Close()
}
