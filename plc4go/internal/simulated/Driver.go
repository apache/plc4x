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
	"github.com/apache/plc4x/plc4go/pkg/api"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"net/url"
)

type Driver struct {
	_default.DefaultDriver
	valueHandler ValueHandler
}

func NewDriver() plc4go.PlcDriver {
	return &Driver{
		DefaultDriver: _default.NewDefaultDriver("simulated", "Simulated PLC4X Datasource", "none", NewFieldHandler()),
		valueHandler:  NewValueHandler(),
	}
}

// GetConnection Establishes a connection to a given PLC using the information in the connectionString
func (d *Driver) GetConnection(_ url.URL, _ map[string]transports.Transport, options map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	connection := NewConnection(NewDevice("test"), d.GetPlcFieldHandler(), d.valueHandler, options)
	return connection.Connect()
}

// SupportsDiscovery returns true if this driver supports discovery
// TODO: Actually the connection could support discovery to list up all fields in the Device
func (d *Driver) SupportsDiscovery() bool {
	return false
}

func (d *Driver) Discover(_ func(event model.PlcDiscoveryItem), _ ...options.WithDiscoveryOption) error {
	return errors.New("unsupported operation")
}
