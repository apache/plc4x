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

// Package profinet is the Profinet IO controller driver. Discovery and the GSDML reader are ported
// from plc4j's profinet-ng driver; the connection (application relation and cyclic data) is built
// from IEC 61158-6-10 and captured controller traffic, because plc4j's never completes.
package profinet

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
)

// DriverCode is the protocol code, plc4j's ProfinetDriver.DRIVER_CODE.
const DriverCode = "profinet"

type Driver struct {
	_default.DefaultDriver

	discoverer *Discoverer

	log      zerolog.Logger
	_options []options.WithOption
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	driver := &Driver{
		discoverer: NewDiscoverer(_options...),
		log:        options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
		_options:   _options,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, DriverCode, "Profinet", rawsocket.TransportCode, NewTagHandler())
	return driver
}

func (d *Driver) SupportsDiscovery() bool {
	return true
}

func (d *Driver) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return d.discoverer.Discover(ctx, callback, discoveryOptions...)
}

// GetConnection connects as IO controller to the device the options name, on the interface the
// URL names: profinet:raw-socket://eth0?device-name=et200sp-1&gsd-directory=/etc/gsd
func (d *Driver) GetConnection(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) (plc4go.PlcConnection, error) {
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		return nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl)
	}
	configuration, err := ParseFromOptions(d.log, driverOptions)
	if err != nil {
		return nil, errors.Wrap(err, "invalid driver options")
	}
	// DCP and the cyclic frames share the Profinet EtherType; the connection management goes
	// through the kernel's UDP.
	transportOptions := map[string][]string{rawsocket.EtherTypesOption: {fmt.Sprintf("0x%04X", etherTypeProfinet)}}
	for key, values := range driverOptions {
		if strings.HasPrefix(key, "virtual-eth.") {
			transportOptions[key] = values
		}
	}
	instance, err := transport.CreateTransportInstance(transportUrl, transportOptions, d._options...)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't initialize the transport for %s", transportUrl.String())
	}
	st, ok := instance.(station)
	if !ok {
		return nil, errors.Errorf("the %s transport doesn't give its MAC address; Profinet needs raw-socket", transportUrl.Scheme)
	}
	connection := newConnection(st, configuration, d._options...)
	if err := connection.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "error connecting")
	}
	return connection, nil
}

func (d *Driver) Close() error {
	return d.discoverer.Close()
}
