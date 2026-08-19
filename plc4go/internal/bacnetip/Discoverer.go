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
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type Discoverer struct {
	wg sync.WaitGroup // use to track spawned go routines

	// discoveryTimeout caps both the WhoIs broadcast loop and the wait for IAm
	// responses. Originally hardcoded to 60s; configurable via
	// Configuration.DiscoveryTimeoutSeconds passed through NewDiscovererWithTimeout.
	discoveryTimeout time.Duration

	passLogToModel bool
	log            zerolog.Logger
}

func NewDiscoverer(_options ...options.WithOption) *Discoverer {
	passLoggerToModel, _ := options.ExtractPassLoggerToModel(_options...)
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &Discoverer{
		discoveryTimeout: 5 * time.Second,
		passLogToModel:   passLoggerToModel,
		log:              customLogger,
	}
}

// SetDiscoveryTimeout overrides the default 5-second window for discovery.
// 0 falls back to the default. Intended to be called from Driver setup once
// Configuration.DiscoveryTimeoutSeconds is known.
func (d *Discoverer) SetDiscoveryTimeout(t time.Duration) {
	if t <= 0 {
		t = 5 * time.Second
	}
	d.discoveryTimeout = t
}

func (d *Discoverer) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {
	// TODO: handle ctx
	interfaces, err := extractInterfaces(discoveryOptions)
	if err != nil {
		return errors.Wrap(err, "error extracting interfaces")
	}

	specificOptions, err := extractProtocolSpecificOptions(discoveryOptions)
	if err != nil {
		return errors.Wrap(err, "error extracting protocol specific options")
	}

	communicationChannels, err := d.buildupCommunicationChannels(ctx, interfaces, specificOptions.bacNetPort)
	if err != nil {
		return errors.Wrap(err, "error building communication channels")
	}

	timeout := d.discoveryTimeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	ctx, cancelFunc := utils.WithNamedTimeout(ctx, "discovery timeout", timeout)
	defer cancelFunc()
	incomingBVLCChannel, err := d.broadcastAndDiscover(ctx, communicationChannels, specificOptions)
	if err != nil {
		return errors.Wrap(err, "error broadcasting and discovering")
	}
	// Dispatch received IAm/IHave to the callback in the background so this
	// function can enforce the discovery window and tear down cleanly. Running
	// it synchronously here would block until ctx cancellation AND its select
	// had no ctx.Done branch, so Discover never returned.
	d.wg.Go(func() {
		d.handleIncomingBVLCs(ctx, callback, incomingBVLCChannel)
	})
	// Wait for the discovery window to elapse OR for the caller to cancel.
	select {
	case <-time.After(timeout):
	case <-ctx.Done():
	}
	for _, channel := range communicationChannels {
		_ = channel.Close()
	}
	return nil
}

func (d *Discoverer) broadcastAndDiscover(ctx context.Context, communicationChannels []communicationChannel, specificOptions *protocolSpecificOptions) (chan receivedBvlcMessage, error) {
	incomingBVLCChannel := make(chan receivedBvlcMessage, 32)
	for _, communicationChannelInstance := range communicationChannels {
		if err := ctx.Err(); err != nil {
			return incomingBVLCChannel, err
		}
		// Prepare the discovery packet data
		{
			var lowLimit driverModel.BACnetContextTagUnsignedInteger
			var highLimit driverModel.BACnetContextTagUnsignedInteger
			if whoIsOptions := specificOptions.whoIsOptions; whoIsOptions != nil && whoIsOptions.limits != nil {
				lowLimit = driverModel.CreateBACnetContextTagUnsignedInteger(0, whoIsOptions.limits.low)
				highLimit = driverModel.CreateBACnetContextTagUnsignedInteger(1, whoIsOptions.limits.high)
			}

			// serializeWhoIs frames a WhoIs BVLC: an Original-Broadcast-NPDU for a
			// broadcast (so peers reply with a broadcast IAm), an
			// Original-Unicast-NPDU for a directed WhoIs.
			serializeWhoIs := func(directed bool) ([]byte, error) {
				requestWhoIs := driverModel.NewBACnetUnconfirmedServiceRequestWhoIs(lowLimit, highLimit)
				apdu := driverModel.NewAPDUUnconfirmedRequest(requestWhoIs)
				control := driverModel.NewNPDUControl(false, false, false, false, driverModel.NPDUNetworkPriority_NORMAL_MESSAGE)
				npdu := driverModel.NewNPDU(1, control, nil, nil, nil, nil, nil, nil, nil, nil, apdu)
				var bvlc driverModel.BVLC
				if directed {
					bvlc = driverModel.NewBVLCOriginalUnicastNPDU(npdu)
				} else {
					bvlc = driverModel.NewBVLCOriginalBroadcastNPDU(npdu)
				}
				return bvlc.Serialize()
			}

			// Always broadcast on the interface (standard discovery on a local
			// subnet).
			if theBytes, err := serializeWhoIs(false); err != nil {
				return nil, err
			} else if _, err := communicationChannelInstance.broadcastConnection.WriteTo(theBytes, communicationChannelInstance.broadcastTarget); err != nil {
				d.log.Debug().Err(err).Msg("Error sending broadcast WhoIs")
			}

			// Additionally send a directed unicast WhoIs to each explicit target.
			// See remoteAddress / remoteAddresses for why this is needed where a
			// broadcast IAm is not routed back to the sender.
			var directedTargets []string
			if specificOptions.remoteAddress != "" {
				directedTargets = append(directedTargets, specificOptions.remoteAddress)
			}
			directedTargets = append(directedTargets, specificOptions.remoteAddresses...)
			if len(directedTargets) > 0 {
				theBytes, err := serializeWhoIs(true)
				if err != nil {
					return nil, err
				}
				for _, ra := range directedTargets {
					udpAddr, rerr := resolveBacnetUDPAddr(ra, specificOptions.bacNetPort)
					if rerr != nil {
						d.log.Warn().Err(rerr).Str("remoteAddress", ra).Msg("invalid directed WhoIs target; skipping")
						continue
					}
					if _, err := communicationChannelInstance.unicastConnection.WriteTo(theBytes, udpAddr); err != nil {
						d.log.Debug().Err(err).Stringer("target", udpAddr).Msg("Error sending directed WhoIs")
					}
				}
			}
		}
		if whoHasOptions := specificOptions.whoHasOptions; whoHasOptions != nil {
			var lowLimit driverModel.BACnetContextTagUnsignedInteger
			var highLimit driverModel.BACnetContextTagUnsignedInteger
			if limits := whoHasOptions.limits; limits != nil {
				lowLimit = driverModel.CreateBACnetContextTagUnsignedInteger(0, limits.deviceInstanceRangeLow)
				highLimit = driverModel.CreateBACnetContextTagUnsignedInteger(1, limits.deviceInstanceRangeHigh)
			}
			var object driverModel.BACnetUnconfirmedServiceRequestWhoHasObject
			if identifier := whoHasOptions.object.identifier; identifier != nil {
				var objectType uint16
				objectTypeByName, ok := driverModel.BACnetObjectTypeByName(identifier.type_)
				if ok {
					parseUint, err := strconv.ParseUint(identifier.type_, 10, 16)
					if err != nil {
						return nil, err
					}
					objectType = uint16(parseUint)
				} else {
					objectType = uint16(objectTypeByName)
				}
				objectIdentifier := driverModel.CreateBACnetContextTagObjectIdentifier(2, objectType, uint32(identifier.instance))
				object = driverModel.NewBACnetUnconfirmedServiceRequestWhoHasObjectIdentifier(objectIdentifier.GetHeader(), objectIdentifier)
			} else if name := whoHasOptions.object.name; name != nil {
				characterString := driverModel.CreateBACnetContextTagCharacterString(3, driverModel.BACnetCharacterEncoding_ISO_10646, *name)
				object = driverModel.NewBACnetUnconfirmedServiceRequestWhoHasObjectName(characterString.GetHeader(), characterString)
			} else {
				// Neither identifier nor name supplied — the caller wired the
				// who-has options together inconsistently. Skip this discovery
				// loop entry rather than crashing the entire driver.
				d.log.Warn().Msg("WhoHas options missing both identifier and name; skipping")
				continue
			}
			requestWhoHas := driverModel.NewBACnetUnconfirmedServiceRequestWhoHas(lowLimit, highLimit, object)
			apdu := driverModel.NewAPDUUnconfirmedRequest(requestWhoHas)

			control := driverModel.NewNPDUControl(false, false, false, false, driverModel.NPDUNetworkPriority_NORMAL_MESSAGE)
			npdu := driverModel.NewNPDU(1, control, nil, nil, nil, nil, nil, nil, nil, nil, apdu)
			bvlc := driverModel.NewBVLCOriginalUnicastNPDU(npdu)

			// Send the search request.
			theBytes, err := bvlc.Serialize()
			if err != nil {
				return nil, err
			}
			if _, err := communicationChannelInstance.broadcastConnection.WriteTo(theBytes, communicationChannelInstance.broadcastTarget); err != nil {
				d.log.Debug().Err(err).Msg("Error sending broadcast")
			}
		}

		d.wg.Go(func() {
			for {
				if err := ctx.Err(); err != nil {
					d.log.Debug().Err(err).Msg("ending")
					return
				}
				blockingReadChan := make(chan bool)
				d.wg.Go(func() {
					buf := make([]byte, 4096)
					n, addr, err := communicationChannelInstance.unicastConnection.ReadFrom(buf)
					if err != nil {
						d.log.Debug().Err(err).Msg("Ending unicast receive")
						blockingReadChan <- false
						return
					}
					d.log.Debug().Stringer("addr", addr).Msg("Received unicast bvlc")
					ctxForModel := options.GetLoggerContextForModel(ctx, d.log, options.WithPassLoggerToModel(d.passLogToModel))
					incomingBvlc, err := driverModel.BVLCParse[driverModel.BVLC](ctxForModel, buf[:n])
					if err != nil {
						d.log.Warn().Err(err).Msg("Could not parse bvlc")
						blockingReadChan <- true
						return
					}
					incomingBVLCChannel <- receivedBvlcMessage{incomingBvlc, addr}
					blockingReadChan <- true
				})
				select {
				case ok := <-blockingReadChan:
					if !ok {
						d.log.Debug().Msg("Ending unicast reading")
						return
					}
					d.log.Trace().Msg("Received something unicast")
				case <-ctx.Done():
					d.log.Debug().Err(ctx.Err()).Msg("Ending unicast receive")
					return
				}
			}
		})

		d.wg.Go(func() {
			for {
				if err := ctx.Err(); err != nil {
					d.log.Debug().Err(err).Msg("ending")
					return
				}
				blockingReadChan := make(chan bool)
				d.wg.Go(func() {
					buf := make([]byte, 4096)
					n, addr, err := communicationChannelInstance.broadcastConnection.ReadFrom(buf)
					if err != nil {
						d.log.Debug().Err(err).Msg("Ending broadcast receive")
						blockingReadChan <- false
						return
					}
					d.log.Debug().Stringer("addr", addr).Msg("Received broadcast bvlc")
					ctxForModel := options.GetLoggerContextForModel(ctx, d.log, options.WithPassLoggerToModel(d.passLogToModel))
					incomingBvlc, err := driverModel.BVLCParse[driverModel.BVLC](ctxForModel, buf[:n])
					if err != nil {
						d.log.Warn().Err(err).Msg("Could not parse bvlc")
						blockingReadChan <- true
					}
					incomingBVLCChannel <- receivedBvlcMessage{incomingBvlc, addr}
					blockingReadChan <- true
				})
				select {
				case ok := <-blockingReadChan:
					if !ok {
						d.log.Debug().Msg("Ending broadcast reading")
						return
					}
					d.log.Trace().Msg("Received something broadcast")
				case <-ctx.Done():
					d.log.Debug().Err(ctx.Err()).Msg("Ending broadcast receive")
					return
				}
			}
		})
	}
	return incomingBVLCChannel, nil
}

func (d *Discoverer) handleIncomingBVLCs(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), incomingBVLCChannel chan receivedBvlcMessage) {
	for {
		select {
		case <-ctx.Done():
			d.log.Debug().Err(ctx.Err()).Msg("Ending incoming BVLC handling")
			return
		case receivedBvlc := <-incomingBVLCChannel:
			var npdu driverModel.NPDU
			if bvlc, ok := receivedBvlc.bvlc.(interface{ GetNpdu() driverModel.NPDU }); ok {
				npdu = bvlc.GetNpdu()
			}
			_ = npdu
			if apdu := npdu.GetApdu(); apdu == nil {
				nlm := npdu.GetNlm()
				d.log.Debug().Interface("nlm", nlm).Msg("Got nlm")
				continue
			}
			apdu := npdu.GetApdu()
			if _, ok := apdu.(driverModel.APDUConfirmedRequest); ok {
				d.log.Debug().Interface("apdu", apdu).Msg("Got apdu")
				continue
			}
			apduUnconfirmedRequest := apdu.(driverModel.APDUUnconfirmedRequest)
			serviceRequest := apduUnconfirmedRequest.GetServiceRequest()
			switch serviceRequest := serviceRequest.(type) {
			case driverModel.BACnetUnconfirmedServiceRequestIAm:
				iAm := serviceRequest
				remoteUrl, err := url.Parse("udp://" + receivedBvlc.addr.String())
				if err != nil {
					d.log.Debug().Err(err).Msg("Error parsing url")
				}
				// A routed I-Am carries the device's origin as an NPDU source
				// specifier (SNET/SADR, ASHRAE 135 clause 6.2.4). Surface it as
				// ready-to-use connection options (the same RemoteNetwork /
				// RemoteAddress keys the connection URL accepts) so consumers can
				// reach the device through the relaying router: transport URL =
				// the router (the frame's UDP source), options = the routed hop.
				discoveryOptions := routedOriginOptions(npdu)
				discoveryEvent := spiModel.NewDefaultPlcDiscoveryItem(
					"bacnet-ip",
					"udp",
					*remoteUrl,
					discoveryOptions,
					fmt.Sprintf("device %v:%v", iAm.GetDeviceIdentifier().GetObjectType(), iAm.GetDeviceIdentifier().GetInstanceNumber()),
					nil,
				)

				// Pass the event back to the callback
				callback(discoveryEvent)
			case driverModel.BACnetUnconfirmedServiceRequestIHave:
				iHave := serviceRequest
				remoteUrl, err := url.Parse("udp://" + receivedBvlc.addr.String())
				if err != nil {
					d.log.Debug().Err(err).Msg("Error parsing url")
				}
				discoveryEvent := spiModel.NewDefaultPlcDiscoveryItem(
					"bacnet-ip",
					"udp",
					*remoteUrl,
					nil,
					fmt.Sprintf("device %v:%v with %v:%v and %v", iHave.GetDeviceIdentifier().GetObjectType(), iHave.GetDeviceIdentifier().GetInstanceNumber(), iHave.GetObjectIdentifier().GetObjectType(), iHave.GetObjectIdentifier().GetInstanceNumber(), iHave.GetObjectName().GetValue()),
					nil,
				)

				// Pass the event back to the callback
				callback(discoveryEvent)
			}
		}
	}
}

func (d *Discoverer) buildupCommunicationChannels(ctx context.Context, interfaces []net.Interface, bacNetPort int) (communicationChannels []communicationChannel, err error) {
	// Iterate over all network devices of this system.
	for _, networkInterface := range interfaces {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		unicastInterfaceAddress, err := networkInterface.Addrs()
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting Addresses for %v", networkInterface)
		}
		// Iterate over all addresses the current interface has configured
		for _, unicastAddress := range unicastInterfaceAddress {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			var ipAddr net.IP
			switch addr := unicastAddress.(type) {
			// If the device is configured to communicate with a subnet
			case *net.IPNet:
				ipAddr = addr.IP.To4()
				if ipAddr == nil {
					// BACnet/IPv6 (Annex U) uses a different BVLC framing
					// (BVLC6) that isn't in the generated model yet, so we
					// skip IPv6 addresses until that lands.
					continue
				}

			// If the device is configured for a point-to-point connection
			case *net.IPAddr:
				ipAddr = addr.IP.To4()
				if ipAddr == nil {
					continue
				}
			default:
				continue
			}

			if !ipAddr.IsGlobalUnicast() {
				continue
			}

			// Plain net.ListenPacket. SO_REUSEPORT/REUSEADDR were attempted in
			// an earlier iteration but they don't actually solve the
			// "multiple BACnet stacks on one host" problem on UDP broadcast —
			// the kernel hashes each broadcast packet to ONE socket, so the
			// second stack still misses traffic. Co-locating BACnet stacks
			// requires a userland demultiplexer, not socket options. If you
			// hit "address already in use" here, stop whatever else is bound
			// to the BACnet/IP UDP port.
			// Bind to the wildcard address (0.0.0.0) rather than the specific
			// interface IP. A socket bound to a specific IP fails to receive the
			// unicast IAm replies on some interfaces (notably the virtual
			// interfaces used in tests) even though the packets arrive at that IP
			// — which is why the conventional BACnet/IP stack (and gobacnet) bind
			// the wildcard. SO_BROADCAST lets the same socket send the WhoIs to
			// the subnet broadcast address.
			ifName := networkInterface.Name
			lc := net.ListenConfig{Control: func(_ string, _ string, c syscall.RawConn) error {
				return controlDiscoverySocket(c, ifName)
			}}
			conn, err := lc.ListenPacket(ctx, "udp4", fmt.Sprintf("0.0.0.0:%d", bacNetPort))
			if err != nil {
				d.log.Debug().Err(err).Msg("Error building discovery socket")
				continue
			}

			_, cidr, _ := net.ParseCIDR(unicastAddress.String())
			broadcastIP := make(net.IP, len(cidr.IP))
			for i := range broadcastIP {
				broadcastIP[i] = cidr.IP[i] | ^cidr.Mask[i]
			}
			broadcastTarget := &net.UDPAddr{IP: broadcastIP, Port: bacNetPort}
			d.log.Debug().
				Str("interface", networkInterface.Name).
				Stringer("local", conn.LocalAddr()).
				Stringer("broadcastTarget", broadcastTarget).
				Msg("discovery channel bound")
			communicationChannels = append(communicationChannels, communicationChannel{
				networkInterface:    networkInterface,
				unicastConnection:   conn,
				broadcastConnection: conn,
				broadcastTarget:     broadcastTarget,
				log:                 d.log,
			})
		}
	}
	return
}

func (d *Discoverer) Close() error {
	defer utils.StopWarn(d.log)()
	d.log.Trace().Msg("Waiting for goroutines to stop")
	d.wg.Wait()
	return nil
}

type receivedBvlcMessage struct {
	bvlc driverModel.BVLC
	addr net.Addr
}

type communicationChannel struct {
	networkInterface    net.Interface
	unicastConnection   net.PacketConn
	broadcastConnection net.PacketConn
	// broadcastTarget is the subnet broadcast address (e.g. 192.168.100.255:47808)
	// to which WhoIs is sent. The socket itself is bound to the wildcard address.
	broadcastTarget net.Addr
	log             zerolog.Logger
}

func (c communicationChannel) Close() error {
	defer utils.StopWarn(c.log)()
	_ = c.unicastConnection.Close()
	// unicastConnection and broadcastConnection are the same socket; closing
	// twice is harmless (the second returns an already-closed error).
	if c.broadcastConnection != c.unicastConnection {
		_ = c.broadcastConnection.Close()
	}
	return nil
}

func extractInterfaces(discoveryOptions []options.WithDiscoveryOption) ([]net.Interface, error) {
	allInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// If no device is explicitly selected via option, simply use all of them
	// However if a discovery option is present to select a device by name, only
	// add those devices matching any of the given names.
	var interfaces []net.Interface
	deviceNames := options.FilterDiscoveryOptionsDeviceName(discoveryOptions)
	if len(deviceNames) > 0 {
		for _, curInterface := range allInterfaces {
			for _, deviceNameOption := range deviceNames {
				if curInterface.Name == deviceNameOption.GetDeviceName() {
					interfaces = append(interfaces, curInterface)
					break
				}
			}
		}
	} else {
		interfaces = allInterfaces
	}
	return interfaces, nil
}

type protocolSpecificOptions struct {
	bacNetPort int
	// remoteAddress, when set, sends the WhoIs as a directed unicast to this
	// host (in addition to the interface broadcast), enabling targeted discovery
	// of a specific device or subnet. Host only or host:port; port defaults to
	// bacNetPort.
	remoteAddress string
	// remoteAddresses, when non-empty, sends a directed unicast WhoIs to each
	// listed host (in addition to the broadcast and remoteAddress). This lets a
	// caller sweep a list of candidate device IPs. A directed unicast forces the
	// sender's OS to ARP each target first, which seeds the reverse path so the
	// IAm reply can be delivered even where a broadcast IAm would not be routed
	// back (through a router/BBMD, or on stacks that learn MACs only from ARP).
	// Each entry is host only or host:port; port defaults to bacNetPort.
	remoteAddresses []string
	whoIsOptions    *struct {
		limits *struct {
			low  uint
			high uint
		}
	}
	whoHasOptions *struct {
		limits *struct {
			deviceInstanceRangeLow  uint
			deviceInstanceRangeHigh uint
		}
		object struct {
			identifier *struct {
				type_    string
				instance uint
			}
			name *string
		}
	}
}

func bacNetPort(port int) option {
	return func(specificOptions *protocolSpecificOptions) error {
		specificOptions.bacNetPort = port
		return nil
	}
}

// resolveBacnetUDPAddr parses a "host" or "host:port" string into a UDP address,
// defaulting the port to defaultPort when none is supplied.
func resolveBacnetUDPAddr(addr string, defaultPort int) (*net.UDPAddr, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		// No port present — treat the whole string as the host.
		host = addr
		portStr = strconv.Itoa(defaultPort)
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return nil, errors.Errorf("invalid remote-address host %q", host)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid remote-address port %q", portStr)
	}
	return &net.UDPAddr{IP: ip, Port: port}, nil
}

func remoteAddress(addr string) option {
	return func(specificOptions *protocolSpecificOptions) error {
		specificOptions.remoteAddress = addr
		return nil
	}
}

func remoteAddresses(addrs []string) option {
	return func(specificOptions *protocolSpecificOptions) error {
		specificOptions.remoteAddresses = addrs
		return nil
	}
}

func whoIsLimits(whoIsLowLimit, whoIsHighLimit uint) option {
	return func(specificOptions *protocolSpecificOptions) error {
		specificOptions.whoIsOptions = &struct {
			limits *struct {
				low  uint
				high uint
			}
		}{&struct {
			low  uint
			high uint
		}{whoIsLowLimit, whoIsHighLimit}}
		return nil
	}
}

func whoHasOption() option {
	return func(specificOptions *protocolSpecificOptions) error {
		specificOptions.whoHasOptions = &struct {
			limits *struct {
				deviceInstanceRangeLow  uint
				deviceInstanceRangeHigh uint
			}
			object struct {
				identifier *struct {
					type_    string
					instance uint
				}
				name *string
			}
		}{}
		return nil
	}
}

func whoHasLimits(whoHasDeviceInstanceRangeLowLimit, whoHasDeviceInstanceRangeHighLimit uint) option {
	return func(specificOptions *protocolSpecificOptions) error {
		if specificOptions.whoHasOptions == nil {
			return errors.New("WithDiscoveryOptionWhoHas must be passed before WhoHasLimits")
		}
		specificOptions.whoHasOptions.limits = &struct {
			deviceInstanceRangeLow  uint
			deviceInstanceRangeHigh uint
		}{whoHasDeviceInstanceRangeLowLimit, whoHasDeviceInstanceRangeHighLimit}
		return nil
	}
}

func whoHasObjectIdentifier(objectIdentifierType string, objectIdentifierInstance uint) option {
	return func(specificOptions *protocolSpecificOptions) error {
		if specificOptions.whoHasOptions == nil {
			return errors.New("WithDiscoveryOptionWhoHas must be passed before WhoHasObjectIdentifier")
		}
		specificOptions.whoHasOptions.object.identifier = &struct {
			type_    string
			instance uint
		}{objectIdentifierType, objectIdentifierInstance}
		return nil
	}
}

func whoHasObjectName(objectName string) option {
	return func(specificOptions *protocolSpecificOptions) error {
		if specificOptions.whoHasOptions == nil {
			return errors.New("WithDiscoveryOptionWhoHas must be passed before WhoHasObjectName")
		}
		specificOptions.whoHasOptions.object.name = &objectName
		return nil
	}
}

func newProtocolSpecificOptions(options ...option) (*protocolSpecificOptions, error) {
	var specificOptions protocolSpecificOptions
	for _, _option := range options {
		if parseErr := _option(&specificOptions); parseErr != nil {
			return nil, parseErr
		}
	}
	return &specificOptions, nil
}

type option func(specificOptions *protocolSpecificOptions) error

func extractProtocolSpecificOptions(discoveryOptions []options.WithDiscoveryOption) (*protocolSpecificOptions, error) {
	var collectedOptions []option
	filteredOptionMap := make(map[string][]any)
	for _, protocolSpecificOption := range options.FilterDiscoveryOptionProtocolSpecific(discoveryOptions) {
		key := protocolSpecificOption.GetKey()
		value := protocolSpecificOption.GetValue()
		if _, ok := filteredOptionMap[key]; !ok {
			filteredOptionMap[key] = make([]any, 0)
		}
		filteredOptionMap[key] = append(filteredOptionMap[key], value)
	}
	keyDependencies := map[string][]struct {
		key           string
		mustBePresent bool
	}{
		"who-is-low-limit":                         {{"who-is-high-limit", true}},
		"who-is-high-limit":                        {{"who-is-low-limit", true}},
		"who-has-device-instance-range-low-limit":  {{"who-has-device-instance-range-high-limit", true}, {"who-has-object*", true}},
		"who-has-device-instance-range-high-limit": {{"who-has-device-instance-range-low-limit", true}, {"who-has-object*", true}},
		"who-has-object-identifier-type":           {{"who-has-object-identifier-instance", true}, {"who-has-object-name", false}},
		"who-has-object-identifier-instance":       {{"who-has-object-identifier-type", true}, {"who-has-object-name", false}},
		"who-has-object-name":                      {{"who-has-object-identifier-instance", false}, {"who-has-object-identifier-type", false}},
	}
	for key, value := range keyDependencies {
		if _, ok := filteredOptionMap[key]; ok {
			for _, otherKey := range value {
				if before, ok0 := strings.CutSuffix(otherKey.key, "*"); ok0 {
					prefix := before
					mustBePresent := otherKey.mustBePresent
					var found bool
					for key := range filteredOptionMap {
						found = found || strings.HasPrefix(key, prefix)
					}
					if mustBePresent && !found {
						return nil, errors.Errorf("When %s is set one of %s must also be set", key, otherKey.key)
					} else if !mustBePresent && found {
						return nil, errors.Errorf("When %s is set none of %s must be set", key, otherKey.key)
					}
				} else if _, otherOk := filteredOptionMap[otherKey.key]; otherOk && !otherKey.mustBePresent {
					return nil, errors.Errorf("When %s is set %s must not be set", key, otherKey.key)
				} else if !otherOk && otherKey.mustBePresent {
					return nil, errors.Errorf("When %s is set %s must be set too", key, otherKey.key)
				}
			}
		}
	}
	if _, ok := filteredOptionMap["bacnet-port"]; ok {
		parsedInt, err := OneInt(filteredOptionMap, "bacnet-port")
		if err != nil {
			return nil, err
		}
		collectedOptions = append(collectedOptions, bacNetPort(parsedInt))
	} else {
		collectedOptions = append(collectedOptions, bacNetPort(47808))
	}

	if _, ok := filteredOptionMap["remote-address"]; ok {
		addr, err := OneString(filteredOptionMap, "remote-address")
		if err != nil {
			return nil, err
		}
		collectedOptions = append(collectedOptions, remoteAddress(addr))
	}

	if _, ok := filteredOptionMap["remote-addresses"]; ok {
		joined, err := OneString(filteredOptionMap, "remote-addresses")
		if err != nil {
			return nil, err
		}
		var addrs []string
		for _, a := range strings.Split(joined, ",") {
			if a = strings.TrimSpace(a); a != "" {
				addrs = append(addrs, a)
			}
		}
		if len(addrs) > 0 {
			collectedOptions = append(collectedOptions, remoteAddresses(addrs))
		}
	}

	if whoIsLow, whoIsHigh, ok, err := func() (whoIsLowLimit uint, whoIsHighLimit uint, ok bool, err error) {
		if _, limitPresent := filteredOptionMap["who-is-low-limit"]; !limitPresent {
			return
		}
		ok = true
		whoIsLowLimit, err = OneUint(filteredOptionMap, "who-is-low-limit")
		whoIsHighLimit, err = OneUint(filteredOptionMap, "who-is-high-limit")
		return
	}(); ok {
		collectedOptions = append(collectedOptions, whoIsLimits(whoIsLow, whoIsHigh))
	} else if err != nil {
		return nil, err
	}
	for key := range filteredOptionMap {
		if strings.HasPrefix(key, "who-has-object") {
			collectedOptions = append(collectedOptions, whoHasOption())
			break
		}
	}
	if whoHasDeviceInstanceRangeLowLimit, whoHasDeviceInstanceRangeHighLimit, ok, err := func() (whoIsLowLimit uint, whoIsHighLimit uint, ok bool, err error) {
		if _, limitPresent := filteredOptionMap["who-has-device-instance-range-low-limit"]; !limitPresent {
			return
		}
		ok = true
		whoIsLowLimit, err = OneUint(filteredOptionMap, "who-has-device-instance-range-low-limit")
		whoIsHighLimit, err = OneUint(filteredOptionMap, "who-has-device-instance-range-high-limit")
		return
	}(); ok {
		collectedOptions = append(collectedOptions, whoHasLimits(whoHasDeviceInstanceRangeLowLimit, whoHasDeviceInstanceRangeHighLimit))
	} else if err != nil {
		return nil, err
	}

	if whoHasObjectIdentifierType, objectIdentifierInstance, ok, err := func() (whoHasObjectIdentifierType string, whoHasObjectIdentifierInstance uint, ok bool, err error) {
		if _, limitPresent := filteredOptionMap["who-has-object-identifier-type"]; !limitPresent {
			return
		}
		ok = true
		whoHasObjectIdentifierType, err = OneString(filteredOptionMap, "who-has-object-identifier-type")
		whoHasObjectIdentifierInstance, err = OneUint(filteredOptionMap, "who-has-object-identifier-instance")
		return
	}(); ok {
		collectedOptions = append(collectedOptions, whoHasObjectIdentifier(whoHasObjectIdentifierType, objectIdentifierInstance))
	} else if err != nil {
		return nil, err
	}

	if _, ok := filteredOptionMap["who-has-object-name"]; ok {
		if name, err := OneString(filteredOptionMap, "who-has-object-name"); err != nil {
			return nil, err
		} else {
			collectedOptions = append(collectedOptions, whoHasObjectName(name))
		}
	}
	return newProtocolSpecificOptions(collectedOptions...)
}

func OneInt(filteredOptionMap map[string][]any, key string) (int, error) {
	value, err := One(filteredOptionMap, key)
	if err != nil {
		return 0, err
	}
	parsedInt, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing option bacnet-port")
	}
	return int(parsedInt), nil
}

func OneUint(filteredOptionMap map[string][]any, key string) (uint, error) {
	value, err := One(filteredOptionMap, key)
	if err != nil {
		return 0, err
	}
	parsedInt, err := strconv.ParseUint(fmt.Sprintf("%v", value), 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "Error parsing option bacnet-port")
	}
	return uint(parsedInt), nil
}
func OneString(filteredOptionMap map[string][]any, key string) (string, error) {
	value, err := One(filteredOptionMap, key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", value), nil
}

func One(filteredOptionMap map[string][]any, key string) (any, error) {
	values := filteredOptionMap[key]
	if len(values) != 1 {
		return nil, errors.Errorf("%s expects only one value", key)
	}
	return values[0], nil
}
