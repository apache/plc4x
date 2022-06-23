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

func broadcastAndDiscover(ctx context.Context, communicationChannels []communicationChannel, whoIsLowLimit *uint, whoIsHighLimit *uint) (chan receivedBvlcMessage, error) {
	incomingBVLCChannel := make(chan receivedBvlcMessage, 0)
	for _, communicationChannelInstance := range communicationChannels {
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
		wbbb := utils.NewWriteBufferByteBased()
		if err := bvlc.Serialize(wbbb); err != nil {
			panic(err)
		}
		if _, err := communicationChannelInstance.broadcastConnection.WriteTo(wbbb.GetBytes(), communicationChannelInstance.broadcastConnection.LocalAddr()); err != nil {
			log.Debug().Err(err).Msg("Error sending broadcast")
		}

		go func(communicationChannelInstance communicationChannel) {
			for {
				blockingReadChan := make(chan bool, 0)
				go func() {
					buf := make([]byte, 4096)
					n, addr, err := communicationChannelInstance.unicastConnection.ReadFrom(buf)
					if err != nil {
						log.Debug().Err(err).Msg("Ending unicast receive")
						blockingReadChan <- false
						return
					}
					log.Debug().Stringer("addr", addr).Msg("Received broadcast bvlc")
					incomingBvlc, err := driverModel.BVLCParse(utils.NewReadBufferByteBased(buf[:n]))
					if err != nil {
						log.Warn().Err(err).Msg("Could not parse bvlc")
						blockingReadChan <- true
						return
					}
					incomingBVLCChannel <- receivedBvlcMessage{incomingBvlc, addr}
					blockingReadChan <- true
				}()
				select {
				case ok := <-blockingReadChan:
					if !ok {
						log.Debug().Msg("Ending unicast reading")
						return
					}
					log.Trace().Msg("Received something unicast")
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
						log.Debug().Err(err).Msg("Ending broadcast receive")
						blockingReadChan <- false
						return
					}
					log.Debug().Stringer("addr", addr).Msg("Received broadcast bvlc")
					incomingBvlc, err := driverModel.BVLCParse(utils.NewReadBufferByteBased(buf[:n]))
					if err != nil {
						log.Warn().Err(err).Msg("Could not parse bvlc")
						blockingReadChan <- true
					}
					incomingBVLCChannel <- receivedBvlcMessage{incomingBvlc, addr}
					blockingReadChan <- true
				}()
				select {
				case ok := <-blockingReadChan:
					if !ok {
						log.Debug().Msg("Ending broadcast reading")
						return
					}
					log.Trace().Msg("Received something broadcast")
				case <-ctx.Done():
					log.Debug().Err(ctx.Err()).Msg("Ending broadcast receive")
					return
				}
			}
		}(communicationChannelInstance)
	}
	return incomingBVLCChannel, nil
}

func handleIncomingBVLCs(ctx context.Context, callback func(event apiModel.PlcDiscoveryEvent), incomingBVLCChannel chan receivedBvlcMessage) {
	for {
		select {
		case receivedBvlc := <-incomingBVLCChannel:
			var npdu driverModel.NPDU
			if bvlc, ok := receivedBvlc.bvlc.(interface{ GetNpdu() driverModel.NPDU }); ok {
				npdu = bvlc.GetNpdu()
			}
			_ = npdu
			if apdu := npdu.GetApdu(); apdu == nil {
				nlm := npdu.GetNlm()
				log.Debug().Msgf("Got nlm\n%v", nlm)
				continue
			}
			apdu := npdu.GetApdu()
			if _, ok := apdu.(driverModel.APDUConfirmedRequestExactly); ok {
				log.Debug().Msgf("Got apdu \n%v", apdu)
				continue
			}
			apduUnconfirmedRequest := apdu.(driverModel.APDUUnconfirmedRequestExactly)
			serviceRequest := apduUnconfirmedRequest.GetServiceRequest()
			if _, ok := serviceRequest.(driverModel.BACnetUnconfirmedServiceRequestIAmExactly); !ok {
				log.Debug().Msgf("Got serviceRequest \n%v", serviceRequest)
				continue
			}
			iam := serviceRequest.(driverModel.BACnetUnconfirmedServiceRequestIAm)
			remoteUrl, err := url.Parse("udp://" + receivedBvlc.addr.String())
			if err != nil {
				log.Debug().Err(err).Msg("Error parsing url")
			}
			discoveryEvent := &internalModel.DefaultPlcDiscoveryEvent{
				ProtocolCode:  "bacnet-ip",
				TransportCode: "udp",
				TransportUrl:  *remoteUrl,
				Name:          fmt.Sprintf("device %v", iam.GetDeviceIdentifier().GetInstanceNumber()),
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

			// Handle undirected
			unicastConnection, err := reuseport.ListenPacket("udp4", fmt.Sprintf("%v:%d", ipAddr, bacNetPort))
			if err != nil {
				log.Debug().Err(err).Msg("Error building unicast Port")
				continue
			}

			_, cidr, _ := net.ParseCIDR(unicastAddress.String())
			broadcastAddr := netaddr.BroadcastAddr(cidr)
			// Handle undirected
			broadcastConnection, err := reuseport.ListenPacket("udp4", fmt.Sprintf("%v:%d", broadcastAddr, bacNetPort))
			if err != nil {
				if err := unicastConnection.Close(); err != nil {
					log.Debug().Err(err).Msg("Error closing transport instance")
				}
				log.Debug().Err(err).Msg("Error building broadcast Port")
				continue
			}
			communicationChannels = append(communicationChannels, communicationChannel{
				networkInterface:    networkInterface,
				unicastConnection:   unicastConnection,
				broadcastConnection: broadcastConnection,
			})
		}
	}
	return
}

type receivedBvlcMessage struct {
	bvlc driverModel.BVLC
	addr net.Addr
}

type communicationChannel struct {
	networkInterface    net.Interface
	unicastConnection   net.PacketConn
	broadcastConnection net.PacketConn
}

func (c communicationChannel) Close() error {
	_ = c.unicastConnection.Close()
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
