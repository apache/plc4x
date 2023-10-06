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
	"github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net/url"
	"runtime/debug"
	"strconv"
)

type TcpDriver struct {
	_default.DefaultDriver

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewModbusTcpDriver(_options ...options.WithOption) *TcpDriver {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	driver := &TcpDriver{
		log:      customLogger,
		_options: _options,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "modbus-tcp", "Modbus TCP", "tcp", NewTagHandler())
	return driver
}

func (d TcpDriver) GetConnectionWithContext(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	d.log.Debug().
		Stringer("transportUrl", &transportUrl).
		Int("nTransports", len(transports)).
		Int("nDriverOptions", len(driverOptions)).
		Msg("Get connection for transport url with nTransports transport(s) and nDriverOptions option(s)")
	// Get an the transport specified in the url
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
	driverOptions["defaultTcpPort"] = []string{"502"}
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
	// TODO: the code below looks strange: where is defaultChanel being used?
	defaultChanel := make(chan any)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				d.log.Error().
					Str("stack", string(debug.Stack())).
					Interface("err", err).
					Msg("panic-ed")
			}
		}()
		for {
			msg := <-defaultChanel
			adu := msg.(model.ModbusTcpADU)
			d.log.Debug().Stringer("adu", adu).Msg("got message in the default handler")
		}
	}()
	codec := NewMessageCodec(
		transportInstance,
		append(d._options, options.WithCustomLogger(d.log))...,
	)
	d.log.Debug().Stringer("codec", codec).Msg("working with codec")

	// If a unit-identifier was provided in the connection string use this, otherwise use the default of 1
	unitIdentifier := uint8(1)
	if value, ok := driverOptions["unit-identifier"]; ok {
		var intValue uint64
		intValue, err = strconv.ParseUint(value[0], 10, 8)
		if err == nil {
			unitIdentifier = uint8(intValue)
		}
	}
	d.log.Debug().Uint8("unitIdentifier", unitIdentifier).Msg("using unit identifier")

	// Create the new connection
	connection := NewConnection(
		unitIdentifier,
		codec,
		driverOptions,
		d.GetPlcTagHandler(),
		append(d._options, options.WithCustomLogger(d.log))...,
	)
	d.log.Debug().Stringer("connection", connection).Msg("created connection, connecting now")
	return connection.ConnectWithContext(ctx)
}
