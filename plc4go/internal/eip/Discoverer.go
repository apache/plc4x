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

package eip

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// EipUdpDiscoveryDefaultPort is the UDP port EIP/CIP devices listen on for
// ListIdentity broadcast discovery requests (mirrors plc4j's
// Constants.EIPUDPDISCOVERYDEFAULTPORT).
const EipUdpDiscoveryDefaultPort = 44818

// discoveryReadBufSize is a sensible upper bound for a ListIdentity response
// datagram.
const discoveryReadBufSize = 4096

// discoveryWindow is how long each per-interface scan waits for ListIdentity
// responses to arrive after sending the broadcast request.
const discoveryWindow = 5 * time.Second

type Discoverer struct {
	wg sync.WaitGroup // use to track spawned go routines

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewDiscoverer(_options ...options.WithOption) *Discoverer {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Discoverer{
		log:      customLogger,
		_options: _options,
	}
}

func buildListIdentityRequest() readWriteModel.EipListIdentityRequest {
	return readWriteModel.NewEipListIdentityRequest(0, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0)
}

// subnetBroadcast computes the subnet broadcast address (all host bits set)
// for the given IPv4 network, e.g. 192.168.1.10/24 -> 192.168.1.255.
// Returns nil if ipNet isn't (or can't be normalized to) an IPv4 network.
func subnetBroadcast(ipNet *net.IPNet) net.IP {
	if ipNet == nil {
		return nil
	}
	ip := ipNet.IP.To4()
	if ip == nil {
		return nil
	}
	mask := ipNet.Mask
	if len(mask) != len(ip) {
		// Normalize a possibly 16-byte mask down to 4 bytes for an IPv4 address.
		if ones, bits := ipNet.Mask.Size(); bits == 32 || bits == 128 {
			mask = net.CIDRMask(ones, 32)
		} else {
			return nil
		}
	}
	broadcast := make(net.IP, len(ip))
	for i := range ip {
		broadcast[i] = ip[i] | ^mask[i]
	}
	return broadcast
}

// identityFromPacket extracts the CipIdentity item from a ListIdentity
// response packet, if present.
func identityFromPacket(packet readWriteModel.EipPacket) (readWriteModel.CipIdentity, bool) {
	listIdentityResponse, ok := packet.(readWriteModel.EipListIdentityResponse)
	if !ok {
		return nil, false
	}
	for _, item := range listIdentityResponse.GetItems() {
		if identity, ok := item.(readWriteModel.CipIdentity); ok {
			return identity, true
		}
	}
	return nil, false
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

	listIdentityRequest := buildListIdentityRequest()
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := listIdentityRequest.SerializeWithWriteBuffer(ctx, wb); err != nil {
		return errors.Wrap(err, "error serializing ListIdentity request")
	}
	requestBytes := wb.GetBytes()

	// Iterate over all network devices of this system.
	for _, netInterface := range interfaces {
		if err := ctx.Err(); err != nil {
			return err
		}
		addrs, err := netInterface.Addrs()
		if err != nil {
			return err
		}
		// Iterate over all addresses the current interface has configured.
		// We're only interested in IPv4 subnets, as we need to compute a
		// broadcast address to send the ListIdentity request to.
		for _, addr := range addrs {
			if err := ctx.Err(); err != nil {
				return err
			}
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ipv4Addr := ipNet.IP.To4()
			if ipv4Addr == nil || ipv4Addr.IsLoopback() {
				continue
			}
			broadcastIP := subnetBroadcast(ipNet)
			if broadcastIP == nil {
				continue
			}

			ifName := netInterface.Name
			broadcastTarget := &net.UDPAddr{IP: broadcastIP, Port: EipUdpDiscoveryDefaultPort}
			d.wg.Go(func() {
				defer func() {
					if err := recover(); err != nil {
						d.log.Error().
							Str("stack", string(debug.Stack())).
							Interface("err", err).
							Msg("panic-ed")
					}
				}()
				d.scanInterface(ctx, ifName, broadcastTarget, requestBytes, callback)
			})
		}
	}
	return nil
}

// scanInterface opens a per-interface UDP socket with SO_BROADCAST set, sends
// the ListIdentity broadcast request to the interface's subnet broadcast
// address, and reads responses until ctx is cancelled or the discovery
// window elapses.
func (d *Discoverer) scanInterface(ctx context.Context, ifName string, broadcastTarget *net.UDPAddr, requestBytes []byte, callback func(event apiModel.PlcDiscoveryItem)) {
	// Bind to an ephemeral port. Unlike BACnet, EIP replies are sent back
	// unicast to the sender's source port, so we don't need to bind the
	// well-known discovery port locally.
	lc := net.ListenConfig{Control: func(_ string, _ string, c syscall.RawConn) error {
		return controlDiscoverySocket(c)
	}}
	conn, err := lc.ListenPacket(ctx, "udp4", "0.0.0.0:0")
	if err != nil {
		d.log.Debug().Err(err).Str("interface", ifName).Msg("error opening discovery socket")
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	d.log.Debug().
		Str("interface", ifName).
		Stringer("local", conn.LocalAddr()).
		Stringer("broadcastTarget", broadcastTarget).
		Msg("sending ListIdentity broadcast")
	if _, err := conn.WriteTo(requestBytes, broadcastTarget); err != nil {
		d.log.Debug().Err(err).Str("interface", ifName).Stringer("broadcastTarget", broadcastTarget).Msg("error sending ListIdentity broadcast")
		return
	}

	buf := make([]byte, discoveryReadBufSize)
	deadline := time.Now().Add(discoveryWindow)
	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			d.log.Debug().Err(err).Msg("done")
			return
		}
		// Tick the read deadline so we periodically re-check ctx.Done() and
		// the overall discovery window even if nothing arrives.
		tick := 500 * time.Millisecond
		if remaining := time.Until(deadline); remaining < tick {
			tick = remaining
		}
		if err := conn.SetReadDeadline(time.Now().Add(tick)); err != nil {
			d.log.Debug().Err(err).Msg("error setting read deadline")
			return
		}
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			d.log.Debug().Err(err).Msg("error reading discovery response")
			return
		}
		rb := utils.NewReadBufferByteBased(buf[:n], utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
		packet, err := readWriteModel.EipPacketParseWithBuffer[readWriteModel.EipPacket](ctx, rb, true)
		if err != nil {
			d.log.Debug().Err(err).Msg("error parsing discovery response")
			continue
		}
		identity, ok := identityFromPacket(packet)
		if !ok {
			continue
		}
		// The responder's IP is taken from the UDP source address of the
		// reply, not from the identity's socket-address field (which some
		// devices leave zeroed or report inconsistently).
		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			host = addr.String()
		}
		remoteUrl, err := url.Parse(fmt.Sprintf("eip://%s:%d", host, EipUdpDiscoveryDefaultPort))
		if err != nil {
			d.log.Debug().Err(err).Msg("error building discovery url")
			continue
		}
		discoveryEvent := spiModel.NewDefaultPlcDiscoveryItem(
			"eip",
			"udp",
			*remoteUrl,
			nil,
			identity.GetProductName(),
			nil,
		)
		// Pass the event back to the callback
		callback(discoveryEvent)
	}
}

func (d *Discoverer) Close() error {
	defer utils.StopWarn(d.log)()
	d.log.Trace().Msg("Closing discoverer")
	d.log.Trace().Msg("waiting for wait group")
	d.wg.Wait()
	return nil
}
