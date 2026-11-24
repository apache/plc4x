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

// Package profinet is the Profinet IO controller driver, ported from plc4j's profinet-ng driver.
// So far it discovers devices (DCP Identify); connecting to them is not built yet.
package profinet

import (
	"context"
	"net/url"

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

	log zerolog.Logger
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	driver := &Driver{
		discoverer: NewDiscoverer(_options...),
		log:        options.ExtractCustomLoggerOrDefaultToGlobal(_options...),
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, DriverCode, "Profinet", rawsocket.TransportCode, nil)
	return driver
}

func (d *Driver) SupportsDiscovery() bool {
	return true
}

func (d *Driver) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return d.discoverer.Discover(ctx, callback, discoveryOptions...)
}

func (d *Driver) GetConnection(_ context.Context, _ url.URL, _ map[string]transports.Transport, _ map[string][]string) (plc4go.PlcConnection, error) {
	return nil, errors.New("the profinet driver can discover devices but cannot connect to them yet")
}

func (d *Driver) Close() error {
	return d.discoverer.Close()
}
