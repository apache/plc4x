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
	"fmt"
	"math"
	"net"
	"net/url"
	"strconv"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Driver struct {
	_default.DefaultDriver
	tm                      spi.RequestTransactionManager
	awaitSetupComplete      bool
	awaitDisconnectComplete bool
	DeviceInventory         DeviceInventory
}

func NewDriver() plc4go.PlcDriver {
	return &Driver{
		DefaultDriver:           _default.NewDefaultDriver("bacnet-ip", "BACnet/IP", "udp", NewTagHandler()),
		tm:                      *spi.NewRequestTransactionManager(math.MaxInt),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,
	}
}

func (m *Driver) GetConnection(transportUrl url.URL, transports map[string]transports.Transport, options map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	log.Debug().Stringer("transportUrl", &transportUrl).Msgf("Get connection for transport url with %d transport(s) and %d option(s)", len(transports), len(options))
	// Get an the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		log.Error().Stringer("transportUrl", &transportUrl).Msgf("We couldn't find a transport for scheme %s", transportUrl.Scheme)
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl))
		}()
		return ch
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	options["defaultUdpPort"] = []string{strconv.Itoa(int(model.BacnetConstants_BACNETUDPDEFAULTPORT))}
	// Set so_reuse by default
	if _, ok := options["so-reuse"]; !ok {
		options["so-reuse"] = []string{"true"}
	}
	var udpTransport *udp.Transport
	switch transport := transport.(type) {
	case *udp.Transport:
		udpTransport = transport
	default:
		log.Error().Stringer("transportUrl", &transportUrl).Msg("Only udp supported at the moment")
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl))
		}()
		return ch
	}

	var localAddr *net.UDPAddr
	{
		host := transportUrl.Host
		port := transportUrl.Port()
		if transportUrl.Port() == "" {
			port = options["defaultUdpPort"][0]
		}
		var remoteAddr *net.UDPAddr
		if resolvedRemoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", host, port)); err != nil {
			panic(err)
		} else {
			remoteAddr = resolvedRemoteAddr
		}
		if dial, err := net.DialUDP("udp", nil, remoteAddr); err != nil {
			log.Error().Stringer("transportUrl", &transportUrl).Msg("host unreachable")
			ch := make(chan plc4go.PlcConnectionConnectResult)
			go func() {
				ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't dial to host %#v", transportUrl.Host))
			}()
			return ch
		} else {
			localAddr = dial.LocalAddr().(*net.UDPAddr)
			localAddr.Port, _ = strconv.Atoi(port)
			_ = dial.Close()
		}
	}
	// Have the transport create a new transport-instance.
	transportInstance, err := udpTransport.CreateTransportInstanceForLocalAddress(transportUrl, options, localAddr)
	if err != nil {
		log.Error().Stringer("transportUrl", &transportUrl).Msgf("We couldn't create a transport instance for port %#v", options["defaultUdpPort"])
		ch := make(chan plc4go.PlcConnectionConnectResult)
		go func() {
			ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't initialize transport configuration for given transport url %v", transportUrl))
		}()
		return ch
	}

	codec := NewApplicationLayerMessageCodec(transportInstance, &m.DeviceInventory)
	log.Debug().Msgf("working with codec %#v", codec)

	// Create the new connection
	connection := NewConnection(codec, m.GetPlcTagHandler(), &m.tm, options)
	log.Debug().Msg("created connection, connecting now")
	return connection.Connect()
}

func (m *Driver) SupportsDiscovery() bool {
	return true
}

func (m *Driver) Discover(callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return m.DiscoverWithContext(context.TODO(), callback, discoveryOptions...)
}

func (m *Driver) DiscoverWithContext(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return NewDiscoverer().Discover(ctx, callback, discoveryOptions...)
}
