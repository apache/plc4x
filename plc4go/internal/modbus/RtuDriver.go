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

package modbus

import (
	"context"
	"net/url"
	"sync"

	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

type RtuDriver struct {
	_default.DefaultDriver

	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewModbusRtuDriver(_options ...options.WithOption) *RtuDriver {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	driver := &RtuDriver{
		log:      customLogger,
		_options: _options,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "modbus-rtu", "Modbus RTU", "serial", NewTagHandler())
	return driver
}

func (d *RtuDriver) GetConnection(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) (plc4go.PlcConnection, error) {
	connectionLog := d.log.With().Ctx(ctx).Str("transportUrl", transportUrl.String()).Logger()
	configuration, err := ParseFromOptions(connectionLog, driverOptions)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse config")
	}
	connectionLog = connectionLog.With().Uint8("unitIdentifier", configuration.unitIdentifier).Logger()
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
	driverOptions["defaultTcpPort"] = []string{"502"}
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

	// Create a new codec for taking care of encoding/decoding of messages. RTU frames its
	// messages with an address and a CRC rather than with an MBAP header, so it needs a codec
	// (and, further down, an ADU flavor) of its own.
	codec := NewMessageCodecRtu(
		transportInstance,
		append(d._options, options.WithCustomLogger(connectionLog))...,
	)
	connectionLog.Debug().Interface("codec", codec).Msg("working with codec")

	// Create the new connection
	connection := NewConnection(
		configuration.withFlavor(flavorRtu),
		codec, driverOptions,
		d.GetPlcTagHandler(),
		append(d._options, options.WithCustomLogger(connectionLog))...,
	)
	connectionLog.Debug().Interface("connection", connection).Msg("created connection, connecting now")
	if err := connection.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "Error connecting connection")
	}
	return connection, nil
}
