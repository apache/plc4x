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

package transports

import (
	"github.com/apache/plc4x/plc4go/spi/options"
	"io"
	"net/url"
)

type Transport interface {
	io.Closer
	// GetTransportCode Get the short code used to identify this transport (As used in the connection string)
	GetTransportCode() string
	// GetTransportName Get a human-readable name for this transport
	GetTransportName() string
	// CreateTransportInstance creates transport instance
	CreateTransportInstance(transportUrl url.URL, options map[string][]string, _options ...options.WithOption) (TransportInstance, error)
}
