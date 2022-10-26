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

package spi

import "github.com/apache/plc4x/plc4go/spi/transports"

type TransportAware interface {
	// RegisterTransport Manually register a new driver
	RegisterTransport(transport transports.Transport)
	// ListTransportNames List the names of all drivers registered in the system
	ListTransportNames() []string
	// GetTransport Get access to a driver instance for a given driver-name
	GetTransport(transportName string, connectionString string, options map[string][]string) (transports.Transport, error)
}
