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
	"fmt"
	"github.com/IBM/netaddr"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/libp2p/go-reuseport"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/plc4x/plc4go/internal/spi"
	internalModel "github.com/apache/plc4x/plc4go/internal/spi/model"
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
	udpTransport := udp.NewTransport()

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

	var whoIsLowLimit *uint
	var whoIsHighLimit *uint
	var bacNetPort *int
	for _, protocolSpecificOption := range options.FilterDiscoveryOptionProtocolSpecific(discoveryOptions) {
		switch protocolSpecificOption.GetKey() {
		case "bacnet-port":
			bacNetPortParsed, err := strconv.ParseInt(fmt.Sprintf("%v", protocolSpecificOption.GetValue()), 10, 8)
			if err != nil {
				return errors.Wrap(err, "Error parsing option")
			}
			bacNetPortParsedInt := int(bacNetPortParsed)
			bacNetPort = &bacNetPortParsedInt
		case "who-is-low-limit":
			whoIsLowLimitParsed, err := strconv.ParseUint(fmt.Sprintf("%v", protocolSpecificOption.GetValue()), 10, 8)
			if err != nil {
				return errors.Wrap(err, "Error parsing option")
			}
			whoIsLowLimitParsedUint := uint(whoIsLowLimitParsed)
			whoIsLowLimit = &whoIsLowLimitParsedUint
		case "who-is-high-limit":
			whoIsHighLimitParsed, err := strconv.ParseUint(fmt.Sprintf("%v", protocolSpecificOption.GetValue()), 10, 8)
			if err != nil {
				return errors.Wrap(err, "Error parsing option")
			}
			whoIsHighLimitParsedUint := uint(whoIsHighLimitParsed)
			whoIsHighLimit = &whoIsHighLimitParsedUint
		}
	}
	if whoIsLowLimit != nil && whoIsHighLimit == nil || whoIsLowLimit == nil && whoIsHighLimit != nil {
		return errors.Errorf("who-is high-limit must be specified together")
	}
	if bacNetPort == nil {
		defaultBacNetPort := 47808
		bacNetPort = &defaultBacNetPort
	}

	var tranportInstances []transports.TransportInstance
	// Iterate over all network devices of this system.
	for _, networkInterface := range interfaces {
		unicastInterfaceAddress, err := networkInterface.Addrs()
		if err != nil {
			return errors.Wrapf(err, "Error getting Addresses for %v", networkInterface)
		}
		// Iterate over all addresses the current interface has configured
		for _, unicastAddress := range unicastInterfaceAddress {
			var ipAddr net.IP
			switch addr := unicastAddress.(type) {
			// If the device is configured to communicate with a subnet
			case *net.IPNet:
				ipAddr = addr.IP.To4()
				if ipAddr == nil {
					ipAddr = addr.IP.To16()
				}

			// If the device is configured for a point-to-point connection
			case *net.IPAddr:
				ipAddr = addr.IP.To4()
				if ipAddr == nil {
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
			udpAddr := &net.UDPAddr{IP: broadcastAddr, Port: 0xBAC0}
			connectionUrl, err := url.Parse(fmt.Sprintf("udp://%s", udpAddr))
			if err != nil {
				log.Debug().Err(err).Msg("error parsing url")
				continue
			}
			localAddr := &net.UDPAddr{IP: ipAddr, Port: 0xBAC0}
			transportInstance, err :=
				udpTransport.CreateTransportInstanceForLocalAddress(*connectionUrl, nil, localAddr)
			if err != nil {
				log.Debug().Err(err).Msg("error creating transport instance")
				return err
			}
			tranportInstances = append(tranportInstances, transportInstance)
		}
	}

	for _, transportInstance := range tranportInstances {
		// Create a codec for sending and receiving messages.
		codec := NewMessageCodec(transportInstance)
		// Explicitly start the worker
		if err := codec.Connect(); err != nil {
			log.Warn().Err(err).Msg("Error connecting")
			continue
		}

		// Cast to the UDP transport instance so we can access information on the local port.
		udpTransportInstance, ok := transportInstance.(*udp.TransportInstance)
		if !ok {
			return errors.New("couldn't cast transport instance to UDP transport instance")
		}
		_ = udpTransportInstance

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
		err = codec.Send(bvlc)
		go func() {
			// Keep on reading responses till the timeout is done.
			// TODO: Make this configurable
			timeout := time.NewTimer(time.Second * 1)
			timeout.Stop()
			for start := time.Now(); time.Since(start) < time.Second*50; {
				timeout.Reset(time.Second * 1)
				select {
				case message := <-codec.GetDefaultIncomingMessageChannel():
					{
						if !timeout.Stop() {
							<-timeout.C
						}
						//bvlc := message.(driverModel.BVLC)
						_ = message
						println(message.(fmt.Stringer).String())
						continue

						deviceName := "todo"
						remoteUrl, err := url.Parse("udp://todo")
						if err != nil {
							log.Debug().Err(err).Msg("Error parsing url")
							continue
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
					}
					continue
				case <-timeout.C:
					{
						timeout.Stop()
						continue
					}
				}
			}
		}()
	}

	// Handle undirected
	pc, err := reuseport.ListenPacket("udp4", "192.168.178.255:47808")
	if err != nil {
		panic(err)
	}

	go func() {
		defer pc.Close()
		// Keep on reading responses till the timeout is done.
		// TODO: Make this configurable
		timeout := time.NewTimer(time.Second * 1)
		timeout.Stop()
		for start := time.Now(); time.Since(start) < time.Second*50; {
			timeout.Reset(time.Second * 1)

			buf := make([]byte, 1024)
			n, addr, err := pc.ReadFrom(buf)
			if err != nil {
				panic(err)
			}
			parse, err := driverModel.BVLCParse(utils.NewReadBufferByteBased(buf[:n]))
			if err != nil {
				panic(err)
			}
			_ = parse

			fmt.Printf("%s sent this: %s\n", addr, buf[:n])
		}
	}()
	return nil
}
