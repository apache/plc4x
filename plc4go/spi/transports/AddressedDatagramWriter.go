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
	"context"
	"net"
)

// AddressedDatagramWriter is implemented by datagram transport instances that can send a datagram
// from their own socket to an address other than the remote they were created for.
//
// Some protocols need it. A BACnet/IP foreign device registers with a BBMD from the very socket
// the BBMD later forwards broadcasts to, and that socket is the one the connection to a device
// already uses.
//
// Like LocalAddressAware it is a capability rather than part of TransportInstance, so drivers ask
// for it by interface and a decorator that forwards the method serves them as well.
type AddressedDatagramWriter interface {
	// WriteToAddress sends data as one datagram to address. It fails on an instance whose socket
	// is connected to its remote, as such a socket can send nowhere else.
	WriteToAddress(ctx context.Context, data []byte, address *net.UDPAddr) error
}
