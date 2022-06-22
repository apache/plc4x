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
	"github.com/IBM/netaddr"
	internalModel "github.com/apache/plc4x/plc4go/internal/spi/model"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/libp2p/go-reuseport"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/options"
	"github.com/apache/plc4x/plc4go/internal/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/spi/transports/udp"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

type Discoverer struct {
	messageCodec spi.MessageCodec
}

func NewDiscoverer() *Discoverer {
	return &Discoverer{}
}

func (d *Discoverer) Discover(callback func(event apiModel.PlcDiscoveryEvent), discoveryOptions ...options.WithDiscoveryOption) error {
	interfaces, err := extractInterfaces(discoveryOptions)
	if err != nil {
		return errors.Wrap(err, "error extracting interfaces")
	}

	whoIsLowLimit, whoIsHighLimit, bacNetPort, err := extractProtocolSpecificOptions(discoveryOptions)
	if err != nil {
		return errors.Wrap(err, "error extracting protocol specific options")
	}

	communicationChannels, err := buildupCommunicationChannels(interfaces, bacNetPort)
	if err != nil {
		return errors.Wrap(err, "error building communication channels")
	}

	// TODO: make adjustable
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Second*60)
	defer func() {
		cancelFunc()
	}()
	incomingBVLCChannel, err := broadcastAndDiscover(ctx, communicationChannels, whoIsLowLimit, whoIsHighLimit)
	if err != nil {
		return errors.Wrap(err, "error broadcasting and discovering")
	}
	handleIncomingBVLCs(ctx, callback, incomingBVLCChannel)
	// TODO: make adjustable
	time.Sleep(time.Second * 60)
	for _, channel := range communicationChannels {
		_ = channel.Close()
	}
	return nil
}

func broadcastAndDiscover(ctx context.Context, communicationChannels []communicationChannel, whoIsLowLimit *uint, whoIsHighLimit *uint) (chan driverModel.BVLC, error) {
	incomingBVLCChannel := make(chan driverModel.BVLC, 0)
	for _, communicationChannelInstance := range communicationChannels {
		// Create a codec for sending and receiving messages.
		codec := NewMessageCodec(communicationChannelInstance.unicastTransport)
		// Explicitly start the worker
		if err := codec.Connect(); err != nil {
			log.Warn().Err(err).Msg("Error connecting")
			continue
		}

		// Prepare the discovery packet data
		var lowLimit driverModel.BACnetContextTagUnsignedInteger
		if whoIsLowLimit != nil {
			lowLimit = driverModel.CreateBACnetContextTagUnsignedInteger(0, *whoIsLowLimit)
		}
		var highLimit driverModel.BACnetContextTagUnsignedInteger
		if whoIsHighLimit != nil {
			highLimit = driverModel.CreateBACnetContextTagUnsignedInteger(1, *whoIsHighLimit)
		}
		requestWhoIs := driverModel.NewBACnetUnconfirmedServiceRequestWhoIs(lowLimit, highLimit, 0)
		apdu := driverModel.NewAPDUUnconfirmedRequest(requestWhoIs, 0)

		control := driverModel.NewNPDUControl(false, false, false, false, driverModel.NPDUNetworkPriority_NORMAL_MESSAGE)
		npdu := driverModel.NewNPDU(1, control, nil, nil, nil, nil, nil, nil, nil, nil, apdu, 0)
		bvlc := driverModel.NewBVLCOriginalUnicastNPDU(npdu, 0)

		// Send the search request.
		if err := codec.Send(bvlc); err != nil {
			log.Debug().Err(err).Msg("Error sending broadcast")
		}
		go func(communicationChannelInstance communicationChannel) {
			for {
				select {
				case message := <-codec.GetDefaultIncomingMessageChannel():
					if incomingBvlc, ok := message.(driverModel.BVLC); ok {
						incomingBVLCChannel <- incomingBvlc
					}
				case <-ctx.Done():
					log.Debug().Err(ctx.Err()).Msg("Ending unicast receive")
					return
				}
			}
		}(communicationChannelInstance)

		go func(communicationChannelInstance communicationChannel) {
			for {
				blockingReadChan := make(chan bool, 0)
				go func() {
					buf := make([]byte, 4096)
					n, addr, err := communicationChannelInstance.broadcastConnection.ReadFrom(buf)
					if err != nil {
						panic(err)
					}
					log.Debug().Stringer("addr", addr).Msg("Received broadcast bvlc")
					incomingBvlc, err := driverModel.BVLCParse(utils.NewReadBufferByteBased(buf[:n]))
					if err != nil {
						panic(err)
					}
					incomingBVLCChannel <- incomingBvlc
					blockingReadChan <- true
				}()
				select {
				case <-blockingReadChan:
					log.Trace().Msg("Received something")
				case <-ctx.Done():
					log.Debug().Err(ctx.Err()).Msg("Ending unicast receive")
					return
				}
			}
		}(communicationChannelInstance)
	}
	return incomingBVLCChannel, nil
}

func handleIncomingBVLCs(ctx context.Context, callback func(event apiModel.PlcDiscoveryEvent), incomingBVLCChannel chan driverModel.BVLC) {
	for {
		select {
		case bvlc := <-incomingBVLCChannel:
			_ = bvlc
			deviceName := "todo"
			remoteUrl, err := url.Parse("udp://todo")
			if err != nil {
				log.Debug().Err(err).Msg("Error parsing url")
			}
			discoveryEvent := &internalModel.DefaultPlcDiscoveryEvent{
				ProtocolCode:  "bacnet-ip",
				TransportCode: "udp",
				TransportUrl:  *remoteUrl,
				Options:       nil,
				Name:          deviceName,
			}

			// Pass the event back to the callback
			callback(discoveryEvent)
		case <-ctx.Done():
			log.Debug().Err(ctx.Err()).Msg("Ending unicast receive")
			return
		}
	}
}

func buildupCommunicationChannels(interfaces []net.Interface, bacNetPort int) (communicationChannels []communicationChannel, err error) {
	udpTransport := udp.NewTransport()
	// Iterate over all network devices of this system.
	for _, networkInterface := range interfaces {
		unicastInterfaceAddress, err := networkInterface.Addrs()
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting Addresses for %v", networkInterface)
		}
		// Iterate over all addresses the current interface has configured
		for _, unicastAddress := range unicastInterfaceAddress {
			var ipAddr net.IP
			switch addr := unicastAddress.(type) {
			// If the device is configured to communicate with a subnet
			case *net.IPNet:
				ipAddr = addr.IP.To4()
				if ipAddr == nil {
					// TODO: for now we only support ipv4 (reuse doesn't like v6 address strings atm)
					continue
					ipAddr = addr.IP.To16()
				}

			// If the device is configured for a point-to-point connection
			case *net.IPAddr:
				ipAddr = addr.IP.To4()
				if ipAddr == nil {
					// TODO: for now we only support ipv4 (reuse doesn't like v6 address strings atm)
					continue
					ipAddr = addr.IP.To16()
				}
			default:
				continue
			}

			if !ipAddr.IsGlobalUnicast() {
				continue
			}

			_, cidr, _ := net.ParseCIDR(unicastAddress.String())
			broadcastAddr := netaddr.BroadcastAddr(cidr)
			udpAddr := &net.UDPAddr{IP: broadcastAddr, Port: bacNetPort}
			connectionUrl, err := url.Parse(fmt.Sprintf("udp://%s", udpAddr))
			if err != nil {
				log.Debug().Err(err).Msg("error parsing url")
				continue
			}
			localAddr := &net.UDPAddr{IP: ipAddr, Port: bacNetPort}
			transportInstance, err :=
				udpTransport.CreateTransportInstanceForLocalAddress(*connectionUrl, nil, localAddr)
			if err != nil {
				return nil, errors.Wrap(err, "error creating transport instance")
			}
			if err := transportInstance.Connect(); err != nil {
				log.Warn().Err(err).Msgf("Can't connect to %v", localAddr)
				continue
			}

			// Handle undirected
			pc, err := reuseport.ListenPacket("udp4", broadcastAddr.String()+":47808")
			if err != nil {
				if err := transportInstance.Close(); err != nil {
					log.Debug().Err(err).Msg("Error closing transport instance")
				}
				return nil, err
			}
			communicationChannels = append(communicationChannels, communicationChannel{
				networkInterface:    networkInterface,
				unicastTransport:    transportInstance,
				broadcastConnection: pc,
			})
		}
	}
	return
}

type communicationChannel struct {
	networkInterface    net.Interface
	unicastTransport    transports.TransportInstance
	broadcastConnection net.PacketConn
}

func (c communicationChannel) Close() error {
	_ = c.unicastTransport.Close()
	_ = c.broadcastConnection.Close()
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

func extractProtocolSpecificOptions(discoveryOptions []options.WithDiscoveryOption) (whoIsLowLimit *uint, whoIsHighLimit *uint, bacNetPort int, err error) {
	bacNetPort = 47808
	for _, protocolSpecificOption := range options.FilterDiscoveryOptionProtocolSpecific(discoveryOptions) {
		switch protocolSpecificOption.GetKey() {
		case "bacnet-port":
			bacNetPortParsed, parseError := strconv.ParseInt(fmt.Sprintf("%v", protocolSpecificOption.GetValue()), 10, 8)
			if parseError != nil {
				return nil, nil, 0, errors.Wrap(parseError, "Error parsing option")
			}
			bacNetPortParsedInt := int(bacNetPortParsed)
			bacNetPort = bacNetPortParsedInt
		case "who-is-low-limit":
			whoIsLowLimitParsed, parseError := strconv.ParseUint(fmt.Sprintf("%v", protocolSpecificOption.GetValue()), 10, 8)
			if parseError != nil {
				return nil, nil, 0, errors.Wrap(parseError, "Error parsing option")
			}
			whoIsLowLimitParsedUint := uint(whoIsLowLimitParsed)
			whoIsLowLimit = &whoIsLowLimitParsedUint
		case "who-is-high-limit":
			whoIsHighLimitParsed, parseError := strconv.ParseUint(fmt.Sprintf("%v", protocolSpecificOption.GetValue()), 10, 8)
			if parseError != nil {
				return nil, nil, 0, errors.Wrap(parseError, "Error parsing option")
			}
			whoIsHighLimitParsedUint := uint(whoIsHighLimitParsed)
			whoIsHighLimit = &whoIsHighLimitParsedUint
		}
	}
	if whoIsLowLimit != nil && whoIsHighLimit == nil || whoIsLowLimit == nil && whoIsHighLimit != nil {
		return nil, nil, 0, errors.Errorf("who-is high-limit must be specified together")
	}
	return
}
