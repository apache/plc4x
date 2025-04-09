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
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	_default "github.com/apache/plc4x/plc4go/spi/default"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transactions"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Driver struct {
	_default.DefaultDriver

	discoverer              *Discoverer
	applicationManager      ApplicationManager
	tm                      transactions.RequestTransactionManager
	awaitSetupComplete      bool
	awaitDisconnectComplete bool

	log zerolog.Logger
}

func NewDriver(_options ...options.WithOption) plc4go.PlcDriver {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	driver := &Driver{
		discoverer: NewDiscoverer(_options...),
		applicationManager: ApplicationManager{
			applications: map[string]*ApplicationLayerMessageCodec{},
			log:          customLogger,
		},
		tm:                      transactions.NewRequestTransactionManager(math.MaxInt),
		awaitSetupComplete:      true,
		awaitDisconnectComplete: true,

		log: customLogger,
	}
	driver.DefaultDriver = _default.NewDefaultDriver(driver, "bacnet-ip", "BACnet/IP", "udp", NewTagHandler())
	return driver
}

func (d *Driver) GetConnectionWithContext(ctx context.Context, transportUrl url.URL, transports map[string]transports.Transport, driverOptions map[string][]string) <-chan plc4go.PlcConnectionConnectResult {
	d.log.Debug().
		Stringer("transportUrl", &transportUrl).
		Int("nTransports", len(transports)).
		Int("nDriverOptions", len(driverOptions)).
		Msg("Get connection for transport url with nTransports transport(s) and nDriverOptions option(s)")
	// Get the transport specified in the url
	transport, ok := transports[transportUrl.Scheme]
	if !ok {
		d.log.Error().
			Stringer("transportUrl", &transportUrl).
			Str("scheme", transportUrl.Scheme).
			Msg("We couldn't find a transport for scheme")
		ch := make(chan plc4go.PlcConnectionConnectResult, 1)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl))
		return ch
	}
	// Provide a default-port to the transport, which is used, if the user doesn't provide on in the connection string.
	driverOptions["defaultUdpPort"] = []string{strconv.Itoa(int(model.BacnetConstants_BACNETUDPDEFAULTPORT))}
	// Set so_reuse by default
	if _, ok := driverOptions["so-reuse"]; !ok {
		driverOptions["so-reuse"] = []string{"true"}
	}
	var udpTransport *udp.Transport
	switch transport := transport.(type) {
	case *udp.Transport:
		udpTransport = transport
	default:
		d.log.Error().Stringer("transportUrl", &transportUrl).Msg("Only udp supported at the moment")
		ch := make(chan plc4go.PlcConnectionConnectResult, 1)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("couldn't find transport for given transport url %#v", transportUrl))
		return ch
	}

	codec, err := d.applicationManager.getApplicationLayerMessageCodec(udpTransport, transportUrl, driverOptions)
	if err != nil {
		ch := make(chan plc4go.PlcConnectionConnectResult, 1)
		ch <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Wrap(err, "error getting application layer message codec"))
		return ch
	}
	d.log.Debug().Stringer("codec", codec).Msg("working with codec")

	// Create the new connection
	connection := NewConnection(codec, d.GetPlcTagHandler(), d.tm, driverOptions)
	d.log.Debug().Msg("created connection, connecting now")
	return connection.ConnectWithContext(ctx)
}

func (d *Driver) SupportsDiscovery() bool {
	return true
}

func (d *Driver) DiscoverWithContext(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	return d.discoverer.Discover(ctx, callback, discoveryOptions...)
}

func (d *Driver) Close() error {
	defer utils.StopWarn(d.log)()
	d.log.Trace().Msg("Closing driver")
	finalErr := new(utils.MultiError)
	d.log.Trace().Msg("Closing discoverer")
	if err := d.discoverer.Close(); err != nil {
		finalErr.Append(errors.Wrap(err, "failed to close discoverer"))
	}
	d.log.Trace().Msg("Closing transaction manager")
	if err := d.tm.Close(); err != nil {
		finalErr.Append(errors.Wrap(err, "error closing transaction manager"))
	}
	return finalErr.ToErrorIfAny()
}

type ApplicationManager struct {
	sync.Mutex
	applications map[string]*ApplicationLayerMessageCodec

	log zerolog.Logger
}

func (a *ApplicationManager) getApplicationLayerMessageCodec(transport *udp.Transport, transportUrl url.URL, options map[string][]string) (*ApplicationLayerMessageCodec, error) {
	var localAddress *net.UDPAddr
	var remoteAddr *net.UDPAddr
	// Find out the remote and the local ip address by opening an UPD port (which is instantly closed)
	{
		host := transportUrl.Host
		port := transportUrl.Port()
		if transportUrl.Port() == "" {
			port = options["defaultUdpPort"][0]
		}
		if resolvedRemoteAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", host, port)); err != nil {
			panic(err)
		} else {
			remoteAddr = resolvedRemoteAddr
		}
		// TODO: Possibly do with with ip-address matching similar to the raw-socket impl in Java.
		if dial, err := net.DialUDP("udp", nil, remoteAddr); err != nil {
			return nil, errors.Errorf("couldn't dial to host %#v", transportUrl.Host)
		} else {
			localAddress = dial.LocalAddr().(*net.UDPAddr)
			localAddress.Port, _ = strconv.Atoi(port)
			_ = dial.Close()
		}
	}
	a.Lock()
	defer a.Unlock()
	messageCodec, ok := a.applications[localAddress.String()]
	if !ok {
		newMessageCodec, err := NewApplicationLayerMessageCodec(a.log, transport, transportUrl, options, localAddress, remoteAddr)
		if err != nil {
			return nil, errors.Wrap(err, "error creating application layer code")
		}
		a.applications[localAddress.String()] = newMessageCodec
		return newMessageCodec, nil
	}
	return messageCodec, nil
}
