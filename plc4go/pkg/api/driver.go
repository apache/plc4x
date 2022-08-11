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

package plc4go

import (
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"net/url"
)

type PlcDriver interface {
	// GetProtocolCode Get the short code used to identify this driver (As used in the connection string)
	GetProtocolCode() string
	// GetProtocolName Get a human readable name for this driver
	GetProtocolName() string

	// GetDefaultTransport If the driver has a default form of transport, provide this and make
	// providing the transport code optional in the connection string
	GetDefaultTransport() string

	// CheckQuery Have the driver parse the query string and provide feedback if it's not a valid one
	CheckQuery(query string) error

	// GetConnection Establishes a connection to a given PLC using the information in the connectionString
	// FIXME: this leaks spi in the signature move to spi driver or create interfaces. Can also be done by moving spi in a proper module
	GetConnection(transportUrl url.URL, transports map[string]transports.Transport, options map[string][]string) <-chan PlcConnectionConnectResult

	// SupportsDiscovery returns true if this driver supports discovery
	SupportsDiscovery() bool

	// Discover TODO: document me
	// FIXME: this leaks spi in the signature move to spi driver or create interfaces. Can also be done by moving spi in a proper module
	Discover(callback func(event model.PlcDiscoveryEvent), discoveryOptions ...options.WithDiscoveryOption) error
}
