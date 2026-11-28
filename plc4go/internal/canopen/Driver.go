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

// Package canopen is the CANopen driver, ported from plc4j's canopen driver.
//
//	canopen:socketcan://can0?node-id=10&heartbeat=true
//	canopen:virtualcan://bench
package canopen

import (
	"context"
	"net/url"
	"slices"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/socketcan"
	"github.com/apache/plc4x/plc4go/spi/transports/virtualcan"
)

// testTransportCode is the in-process loopback of the driver test suite framework. On it frames use
// the spec's compact CANOpenFrame layout, which is what the protocol's driver test suites are
// written in.
const testTransportCode = "test"

var supportedTransportCodes = []string{socketcan.TransportCode, virtualcan.TransportCode, testTransportCode}

type Driver struct {
	_default.DefaultDriver

	log      zerolog.Logger
	_options []options.WithOption
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	driver := &Driver{
		log:      options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
		_options: _options,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "canopen", "CANopen", socketcan.TransportCode, NewTagHandler())
	return driver
}

func (d *Driver) GetSupportedTransportCodes() []string {
	return append([]string(nil), supportedTransportCodes...)
}

func (d *Driver) GetConnection(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) (plc4go.PlcConnection, error) {
	connectionLog := d.log.With().Ctx(ctx).Str("transportUrl", transportUrl.String()).Logger()
	if !slices.Contains(supportedTransportCodes, transportUrl.Scheme) {
		return nil, errors.Errorf("the canopen driver doesn't support the transport %s, it supports %v", transportUrl.Scheme, supportedTransportCodes)
	}
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		return nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl)
	}
	connectionOptions := append(slices.Clip(d._options), options.WithCustomLogger(connectionLog))
	transportInstance, err := transport.CreateTransportInstance(transportUrl, driverOptions, connectionOptions...)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't initialize the transport for %s", transportUrl.String())
	}
	configuration, err := ParseFromOptions(connectionLog, driverOptions)
	if err != nil {
		return nil, errors.Wrap(err, "invalid driver options")
	}
	codec := NewMessageCodec(transportInstance, transportUrl.Scheme == testTransportCode, connectionOptions...)
	connection := NewConnection(codec, configuration, driverOptions, connectionOptions...)
	if err := connection.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "error connecting")
	}
	return connection, nil
}

func (d *Driver) Close() error {
	return nil
}
