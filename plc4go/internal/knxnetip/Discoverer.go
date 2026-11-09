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
	"golang.org/x/net/ipv4"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// knxMulticastDiscoveryAddress is the KNXnet/IP system-setup multicast address all
// KNXnet/IP routers/interfaces listen on for SearchRequests.
var knxMulticastDiscoveryAddress = &net.UDPAddr{IP: net.IPv4(224, 0, 23, 12), Port: 3671}

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

	// discoveryWindow is how long a scan waits for SearchResponses (discoveryWindow by
	// default, shorter in the tests).
	discoveryWindow time.Duration

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewDiscoverer(_options ...options.WithOption) *Discoverer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	closeCtx, closeCancel := context.WithCancel(context.Background())
	return &Discoverer{
		closeCtx:        closeCtx,
		closeCancel:     closeCancel,
		discoveryWindow: discoveryWindow,
		log:             customLogger,
		_options:        _options,
	}
}

func (d *Discoverer) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
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
				d.scanAddress(scanCtx, knxMulticastDiscoveryAddress, ipv4Addr, callback)
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

// scanAddress opens one udp socket bound to the given local ipv4 address, sends a
// SearchRequest to the given destination (the KNXnet/IP discovery multicast address) and
// reports every SearchResponse it receives within the discovery window to callback.
//
// The socket is not connected: a gateway answers from its own unicast address to the endpoint
// the request names, and a socket connected to the multicast address would drop all of that.
// Every resource it acquires is released again before it returns, on all paths.
func (d *Discoverer) scanAddress(ctx context.Context, destination *net.UDPAddr, ipv4Addr net.IP, callback func(event apiModel.PlcDiscoveryItem)) {
	if err := ctx.Err(); err != nil {
		d.log.Debug().Err(err).Msg("done")
		return
	}
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: ipv4Addr})
	if err != nil {
		d.log.Error().Err(err).Stringer("localAddress", ipv4Addr).Msg("error opening the discovery socket")
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			d.log.Debug().Err(err).Msg("error closing the discovery socket")
		}
	}()
	localAddress := conn.LocalAddr().(*net.UDPAddr)
	if destination.IP.IsMulticast() {
		// Send the request out of the interface this address belongs to, not the one the
		// routing table picks for the multicast group.
		if netInterface := interfaceOf(ipv4Addr); netInterface != nil {
			if err := ipv4.NewPacketConn(conn).SetMulticastInterface(netInterface); err != nil {
				d.log.Debug().Err(err).Str("interface", netInterface.Name).Msg("error selecting the multicast interface")
			}
		}
	}
	d.log.Debug().Stringer("localAddress", localAddress).Stringer("destination", destination).Msg("Scanning")

	// The request names the socket it is sent from: that is where the gateways answer.
	searchRequestMessage := driverModel.NewSearchRequest(driverModel.NewHPAIDiscoveryEndpoint(
		driverModel.HostProtocolCode_IPV4_UDP, driverModel.NewIPAddress(localAddress.IP.To4()), uint16(localAddress.Port)))
	data, err := searchRequestMessage.Serialize()
	if err != nil {
		d.log.Error().Err(err).Msg("error serializing the search request")
		return
	}
	if _, err := conn.WriteToUDP(data, destination); err != nil {
		d.log.Debug().Err(err).Stringer("destination", destination).Msg("Error sending the search request")
		return
	}

	// Keep on reading responses till the discovery window is done.
	window := d.discoveryWindow
	if window <= 0 {
		window = discoveryWindow
	}
	deadline := time.Now().Add(window)
	buf := make([]byte, 1500)
	ctxForModel := options.GetLoggerContextForModel(ctx, d.log)
	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			d.log.Debug().Err(err).Msg("done")
			return
		}
		// Wake up at least every tick to notice a cancelled context.
		readDeadline := time.Now().Add(discoveryReceiveTick)
		if deadline.Before(readDeadline) {
			readDeadline = deadline
		}
		_ = conn.SetReadDeadline(readDeadline)
		n, source, err := conn.ReadFromUDP(buf)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				continue
			}
			d.log.Debug().Err(err).Msg("error receiving search responses")
			return
		}
		message, err := driverModel.KnxNetIpMessageParse[driverModel.KnxNetIpMessage](ctxForModel, buf[:n])
		if err != nil {
			d.log.Debug().Err(err).Stringer("source", source).Msg("unparseable datagram")
			continue
		}
		discoveryItem, ok := d.discoveryItemFromMessage(message, source)
		if !ok {
			continue
		}
		// Pass the event back to the callback
		callback(discoveryItem)
	}
}

// interfaceOf returns the network interface the given address is configured on.
func interfaceOf(address net.IP) *net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, netInterface := range interfaces {
		addrs, err := netInterface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.Equal(address) {
				return &netInterface
			}
		}
	}
	return nil
}

// discoveryItemFromMessage turns a received message into a discovery item. It
// returns false for anything which isn't a usable SearchResponse. A gateway names
// the control endpoint to connect to; one behind NAT names 0.0.0.0:0, which means
// the address the response came from.
func (d *Discoverer) discoveryItemFromMessage(message spi.Message, source *net.UDPAddr) (apiModel.PlcDiscoveryItem, bool) {
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
	host := fmt.Sprintf("%d.%d.%d.%d:%d", addr[0], addr[1], addr[2], addr[3], controlEndpoint.GetIpPort())
	if net.IP(addr[:4]).IsUnspecified() && controlEndpoint.GetIpPort() == 0 {
		if source == nil {
			d.log.Debug().Msg("search response without a control endpoint")
			return nil, false
		}
		host = source.String()
	}
	remoteUrl, err := url.Parse("udp://" + host)
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
