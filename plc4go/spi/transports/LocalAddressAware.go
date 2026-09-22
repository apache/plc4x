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
	"net"
	"net/url"

	"github.com/apache/plc4x/plc4go/spi/options"
)

// LocalAddressAware is implemented by transports that can bind the local end of a connection to
// an address the caller chooses, rather than letting the system pick an ephemeral one.
//
// Some protocols require it. BACnet/IP peers reply to the well-known port rather than to the
// port a request came from, so a connection that did not bind that port never sees the answer.
//
// It is a capability rather than part of Transport because most transports have no use for it.
// Drivers must therefore ask for it by interface and not by concrete type: a transport is free
// to be wrapped -- for frame capture, tracing or throttling -- and a decorator that forwards
// this method serves the driver just as well as the transport underneath it.
type LocalAddressAware interface {
	Transport
	// CreateTransportInstanceForLocalAddress creates a transport instance bound to localAddress.
	// A nil localAddress leaves the choice to the system, as CreateTransportInstance does.
	CreateTransportInstanceForLocalAddress(transportUrl url.URL, options map[string][]string, localAddress *net.UDPAddr, _options ...options.WithOption) (TransportInstance, error)
}
