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

package knxnetip

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports/udp"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// knxMulticastDiscoveryUrl is the KNXnet/IP system-setup multicast address all
// KNXnet/IP routers/interfaces listen on for SearchRequests.
const knxMulticastDiscoveryUrl = "udp://224.0.23.12:3671"

// discoveryWindow is how long each per-address scan waits for SearchResponses
// to arrive after sending the SearchRequest.
// TODO: Make this configurable
const discoveryWindow = 5 * time.Second

// discoveryReceiveTick bounds how long a single receive waits before the scan
// loop re-checks the context and the overall discovery window.
const discoveryReceiveTick = 1 * time.Second

type Discoverer struct {
	wg sync.WaitGroup // use to track spawned go routines

	// closeCtx is cancelled by Close, which aborts all in-flight scans instead
	// of making Close block for the remainder of their discovery window.
	closeCtx    context.Context
	closeCancel context.CancelFunc

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewDiscoverer(_options ...options.WithOption) *Discoverer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	closeCtx, closeCancel := context.WithCancel(context.Background())
	return &Discoverer{
		closeCtx:    closeCtx,
		closeCancel: closeCancel,
		log:         customLogger,
		_options:    _options,
	}
}

func (d *Discoverer) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	udpTransport := udp.NewTransport()

	// Create a connection string for the KNX broadcast discovery address.
	connectionUrl, err := url.Parse(knxMulticastDiscoveryUrl)
	if err != nil {
		return err
	}

	allInterfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	// If no device is explicitly selected via option, simply use all of them
	// However if a discovery option is present to select a device by name, only
	// add those devices matching any of the given names.
	var interfaces []net.Interface
	deviceNames := options.FilterDiscoveryOptionsDeviceName(discoveryOptions)
	if len(deviceNames) > 0 {
		for _, curInterface := range allInterfaces {
			if err := ctx.Err(); err != nil {
				return err
			}
			for _, deviceNameOption := range deviceNames {
				if err := ctx.Err(); err != nil {
					return err
				}
				if curInterface.Name == deviceNameOption.GetDeviceName() {
					interfaces = append(interfaces, curInterface)
					break
				}
			}
		}
	} else {
		interfaces = allInterfaces
	}

	// Every scan runs on a context derived from the caller's context which is
	// additionally cancelled by Close. The AfterFunc registration is released
	// again as soon as all scans of this Discover call are done, so neither the
	// caller's context nor closeCtx accumulate registrations over time.
	closeCtx := d.closeCtx
	if closeCtx == nil {
		closeCtx = context.Background()
	}
	scanCtx, cancelScans := context.WithCancel(ctx)
	stopCloseHook := context.AfterFunc(closeCtx, cancelScans)

	// scanWg only tracks the scans of this Discover call, d.wg tracks it in turn
	// so that Close waits for everything without any accounting living inside
	// the scan goroutines themselves.
	scanWg := &sync.WaitGroup{}
	// Iterate over all network devices of this system. A failure on one of them has to
	// leave the loop through scanErr instead of returning right away: the scans which
	// already were started are only covered by Close once the d.wg.Go below registered
	// the waiter, so an early return would leak them out of the "Close waits for
	// everything" guarantee.
	var scanErr error
	for _, netInterface := range interfaces {
		if err := ctx.Err(); err != nil {
			scanErr = err
			break
		}
		addrs, err := netInterface.Addrs()
		if err != nil {
			scanErr = err
			break
		}
		// Iterate over all addresses the current interface has configured
		// For KNX we're only interested in IPv4 addresses, as it doesn't
		// seem to work with IPv6.
		for _, addr := range addrs {
			var ipv4Addr net.IP
			switch typedAddr := addr.(type) {
			// If the device is configured to communicate with a subnet
			case *net.IPNet:
				ipv4Addr = typedAddr.IP.To4()

			// If the device is configured for a point-to-point connection
			case *net.IPAddr:
				ipv4Addr = typedAddr.IP.To4()
			}

			// If we found an IPv4 address and this is not a loopback address,
			// add it to the list of devices we will open ports and send discovery
			// messages from.
			if ipv4Addr == nil || ipv4Addr.IsLoopback() {
				continue
			}
			scanWg.Go(func() {
				defer func() {
					if err := recover(); err != nil {
						d.log.Error().
							Str("stack", string(debug.Stack())).
							Interface("err", err).
							Msg("panic-ed")
					}
				}()
				d.scanAddress(scanCtx, connectionUrl, ipv4Addr, udpTransport, callback)
			})
		}
	}
	d.wg.Go(func() {
		defer cancelScans()
		defer stopCloseHook()
		scanWg.Wait()
		d.log.Trace().Msg("All scans done")
	})
	if scanErr != nil {
		// Stop whatever already got started; the waiter above keeps tracking it, so
		// Close still waits for those scans to actually finish.
		cancelScans()
		return scanErr
	}
	return nil
}

// scanAddress opens one udp socket bound to the given local ipv4 address, sends
// a SearchRequest to the KNXnet/IP discovery multicast address and reports every
// SearchResponse it receives within the discovery window to callback.
//
// Every resource it acquires is released again before it returns, on all paths.
func (d *Discoverer) scanAddress(ctx context.Context, connectionUrl *url.URL, ipv4Addr net.IP, udpTransport *udp.Transport, callback func(event apiModel.PlcDiscoveryItem)) {
	if err := ctx.Err(); err != nil {
		d.log.Debug().Err(err).Msg("done")
		return
	}
	// Create a new "connection" (Actually open a local udp socket and target outgoing packets to that address)
	transportInstance, err := udpTransport.CreateTransportInstanceForLocalAddress(*connectionUrl, nil,
		&net.UDPAddr{IP: ipv4Addr, Port: 0})
	if err != nil {
		d.log.Error().Err(err).Msg("error creating transport instance")
		return
	}
	udpTransportInstance, ok := transportInstance.(*udp.TransportInstance)
	if !ok {
		d.log.Error().Type("transportInstanceType", transportInstance).Msg("unexpected transport instance type")
		if err := transportInstance.Close(); err != nil {
			d.log.Debug().Err(err).Msg("error closing transport instance")
		}
		return
	}

	d.log.Debug().Stringer("udpTransportInstance", udpTransportInstance).Msg("Scanning")
	// Create a codec for sending and receiving messages.
	codec := NewMessageCodec(
		udpTransportInstance,
		nil,
		append(d._options, options.WithCustomLogger(d.log))...,
	)
	// Explicitly start the worker. This also connects the transport instance.
	if err := codec.Connect(ctx); err != nil {
		d.log.Error().Err(err).Msg("Error connecting")
		// The codec never took ownership of the transport instance, so the socket
		// (if it got opened at all) has to be closed here.
		if err := udpTransportInstance.Close(); err != nil {
			d.log.Debug().Err(err).Msg("error closing transport instance")
		}
		return
	}
	// Disconnecting stops the codecs receive worker and closes the udp socket.
	defer func() {
		if err := codec.Disconnect(); err != nil {
			d.log.Debug().Err(err).Msg("error disconnecting codec")
		}
	}()

	localAddress := udpTransportInstance.LocalAddress
	localAddr := driverModel.NewIPAddress(localAddress.IP)

	// Prepare the discovery packet data
	discoveryEndpoint := driverModel.NewHPAIDiscoveryEndpoint(
		driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port))
	searchRequestMessage := driverModel.NewSearchRequest(discoveryEndpoint)
	// Send the search request.
	if err := codec.Send(ctx, "device_scan_search_request", searchRequestMessage); err != nil {
		d.log.Debug().Err(err).Interface("searchRequestMessage", searchRequestMessage).Msg("Error sending message")
		return
	}

	// Keep on reading responses till the discovery window is done.
	timeout := time.NewTimer(discoveryReceiveTick)
	defer timeout.Stop()
	deadline := time.Now().Add(discoveryWindow)
	for time.Now().Before(deadline) {
		tick := discoveryReceiveTick
		if remaining := time.Until(deadline); remaining < tick {
			tick = remaining
		}
		if !timeout.Stop() {
			select {
			case <-timeout.C:
			default:
			}
		}
		timeout.Reset(tick)
		select {
		case <-ctx.Done():
			d.log.Debug().Err(ctx.Err()).Msg("done")
			return
		case message := <-codec.GetDefaultIncomingMessageChannel():
			discoveryItem, ok := d.discoveryItemFromMessage(message)
			if !ok {
				continue
			}
			// Pass the event back to the callback
			callback(discoveryItem)
		case <-timeout.C:
			continue
		}
	}
}

// discoveryItemFromMessage turns a received message into a discovery item. It
// returns false for anything which isn't a usable SearchResponse.
func (d *Discoverer) discoveryItemFromMessage(message spi.Message) (apiModel.PlcDiscoveryItem, bool) {
	if message == nil {
		return nil, false
	}
	searchResponse, ok := message.(driverModel.SearchResponse)
	if !ok {
		d.log.Debug().Type("messageType", message).Msg("unexpected message type")
		return nil, false
	}
	controlEndpoint := searchResponse.GetHpaiControlEndpoint()
	deviceInfo := searchResponse.GetDibDeviceInfo()
	if controlEndpoint == nil || controlEndpoint.GetIpAddress() == nil || deviceInfo == nil {
		d.log.Debug().Msg("incomplete search response")
		return nil, false
	}
	addr := controlEndpoint.GetIpAddress().GetAddr()
	if len(addr) < 4 {
		d.log.Debug().Int("addrLength", len(addr)).Msg("unexpected address length")
		return nil, false
	}
	remoteUrl, err := url.Parse(fmt.Sprintf("udp://%d.%d.%d.%d:%d",
		addr[0], addr[1], addr[2], addr[3], controlEndpoint.GetIpPort()))
	if err != nil {
		d.log.Debug().Err(err).Msg("error building discovery url")
		return nil, false
	}
	deviceName := string(bytes.Trim(deviceInfo.GetDeviceFriendlyName(), "\x00"))
	return spiModel.NewDefaultPlcDiscoveryItem(
		"knxnet-ip",
		"udp",
		*remoteUrl,
		nil,
		deviceName,
		nil,
	), true
}

func (d *Discoverer) Close() error {
	defer utils.StopWarn(d.log)()
	d.log.Trace().Msg("Closing discoverer")
	if d.closeCancel != nil {
		d.log.Trace().Msg("Cancelling running scans")
		d.closeCancel()
	}
	d.log.Trace().Msg("waiting for wait group")
	d.wg.Wait()
	return nil
}
