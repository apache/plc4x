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

package simulated

import (
	"context"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/rs/zerolog"
	"net/url"

	"github.com/apache/plc4x/plc4go/pkg/api"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/transports"
)

type Driver struct {
	_default.DefaultDriver
	valueHandler ValueHandler

	log zerolog.Logger
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	driver := &Driver{
		valueHandler: NewValueHandler(),

		log: options.ExtractCustomLogger(_options...),
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "simulated", "Simulated PLC4X Datasource", "none", NewTagHandler())
	return driver
}

func (d *Driver) GetConnectionWithContext(ctx context.Context, _ url.URL, _ map[string]transports.Transport, driverOptions map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	connection := NewConnection(NewDevice("test", options.WithCustomLogger(d.log)), d.GetPlcTagHandler(), d.valueHandler, driverOptions)
	d.log.Debug().Msgf("Connecting and returning connection %v", connection)
	return connection.ConnectWithContext(ctx)
}

// SupportsDiscovery returns true if this driver supports discovery
// TODO: Actually the connection could support discovery to list up all tags in the Device
func (d *Driver) SupportsDiscovery() bool {
	return false
}
