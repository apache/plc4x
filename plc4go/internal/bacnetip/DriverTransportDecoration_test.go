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

package bacnetip

import (
	"context"
	"net"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
)

// decoratedUdpTransport stands in for any decorator placed in front of the udp transport --
// frame capture, tracing, throttling. It is not a *udp.Transport, but it offers everything the
// driver needs, including the local-address bind.
type decoratedUdpTransport struct {
	transports.Transport

	inner *udp.Transport
	used  bool
}

func (d *decoratedUdpTransport) CreateTransportInstanceForLocalAddress(transportUrl url.URL, opts map[string][]string, localAddress *net.UDPAddr, _options ...options.WithOption) (transports.TransportInstance, error) {
	d.used = true
	return d.inner.CreateTransportInstanceForLocalAddress(transportUrl, opts, localAddress, _options...)
}

// TestDriverAcceptsDecoratedUdpTransport pins that the driver selects a transport by the
// capability it needs rather than by concrete type. A decorator satisfies the transport
// interface but is not a *udp.Transport, so a concrete type assertion rejects it and no
// wrapped transport can ever be used.
func TestDriverAcceptsDecoratedUdpTransport(t *testing.T) {
	inner := udp.NewTransport()
	decorated := &decoratedUdpTransport{Transport: inner, inner: inner}

	driver := NewDriver()
	t.Cleanup(func() { _ = driver.Close() })

	connection, err := driver.GetConnection(
		context.Background(),
		url.URL{Scheme: "udp", Host: "127.0.0.1:47808"},
		map[string]transports.Transport{"udp": decorated},
		// An ephemeral local port keeps the test off the well-known BACnet port.
		map[string][]string{"local-port": {"0"}},
	)
	require.NoError(t, err)
	require.NotNil(t, connection)
	t.Cleanup(func() { _ = connection.Close() })

	assert.True(t, decorated.used, "the driver must go through the decorator, not around it")
}
