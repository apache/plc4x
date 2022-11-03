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

package ads

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/protocols/ads/discovery/readwrite/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/ads/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	internalModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/utils"
	values2 "github.com/apache/plc4x/plc4go/spi/values"
	"github.com/rs/zerolog/log"
)

type Discoverer struct {
	messageCodec spi.MessageCodec
}

func NewDiscoverer() *Discoverer {
	return &Discoverer{}
}

func (d *Discoverer) Discover(ctx context.Context, callback func(event apiModel.PlcDiscoveryItem), discoveryOptions ...options.WithDiscoveryOption) error {

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Set up a listening socket on all devices for processing the responses to any search requests

	// Open a listening udp socket for the incoming responses
	responseAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", model.AdsDiscoveryConstants_ADSDISCOVERYUDPDEFAULTPORT))
	if err != nil {
		panic(err)
	}
	socket, err := net.ListenUDP("udp4", responseAddr)
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	// Start a worker to receive responses
	go func() {
		buf := make([]byte, 1024)
		for {
			length, fromAddr, err := socket.ReadFromUDP(buf)
			if length == 0 {
				continue
			}
			readBuffer := utils.NewLittleEndianReadBufferByteBased(buf[0:length])
			discoveryResponse, err := model.AdsDiscoveryParse(readBuffer)
			if err != nil {
				log.Error().Err(err).Str("src-ip", fromAddr.String()).Msg("error decoding response")
				continue
			}

			if (discoveryResponse.GetRequestId() == 0) &&
				(discoveryResponse.GetPortNumber() == model.AdsPortNumbers_SYSTEM_SERVICE) &&
				(discoveryResponse.GetOperation() == model.Operation_DISCOVERY_RESPONSE) {
				remoteAmsNetId := discoveryResponse.GetAmsNetId()
				var hostNameBlock model.AdsDiscoveryBlockHostName
				//var osDataBlock model.AdsDiscoveryBlockOsData
				var versionBlock model.AdsDiscoveryBlockVersion
				var fingerprintBlock model.AdsDiscoveryBlockFingerprint
				for _, block := range discoveryResponse.GetBlocks() {
					switch block.GetBlockType() {
					case model.AdsDiscoveryBlockType_HOST_NAME:
						hostNameBlock = block.(model.AdsDiscoveryBlockHostName)
						/*									case model.AdsDiscoveryBlockType_OS_DATA:
															osDataBlock = block.(model.AdsDiscoveryBlockOsData)*/
					case model.AdsDiscoveryBlockType_VERSION:
						versionBlock = block.(model.AdsDiscoveryBlockVersion)
					case model.AdsDiscoveryBlockType_FINGERPRINT:
						fingerprintBlock = block.(model.AdsDiscoveryBlockFingerprint)
					}
				}

				if hostNameBlock != nil {
					opts := make(map[string][]string)
					//					opts["sourceAmsNetId"] = []string{localIpV4Address.String() + ".1.1"}
					opts["sourceAmsPort"] = []string{"65534"}
					opts["targetAmsNetId"] = []string{strconv.Itoa(int(remoteAmsNetId.GetOctet1())) + "." +
						strconv.Itoa(int(remoteAmsNetId.GetOctet2())) + "." +
						strconv.Itoa(int(remoteAmsNetId.GetOctet3())) + "." +
						strconv.Itoa(int(remoteAmsNetId.GetOctet4())) + "." +
						strconv.Itoa(int(remoteAmsNetId.GetOctet5())) + "." +
						strconv.Itoa(int(remoteAmsNetId.GetOctet6()))}
					// TODO: Check if this is legit, or if we can get the information from somewhere.
					opts["targetAmsPort"] = []string{"851"}

					attributes := make(map[string]values.PlcValue)
					attributes["hostName"] = values2.NewPlcSTRING(hostNameBlock.GetHostName().GetText())
					if versionBlock != nil {
						versionData := versionBlock.GetVersionData()
						patchVersion := (int(versionData[3])&0xFF)<<8 | (int(versionData[2]) & 0xFF)
						attributes["twinCatVersion"] = values2.NewPlcSTRING(fmt.Sprintf("%d.%d.%d", int(versionData[0])&0xFF, int(versionData[1])&0xFF, patchVersion))
					}
					if fingerprintBlock != nil {
						attributes["fingerprint"] = values2.NewPlcSTRING(string(fingerprintBlock.GetData()))
					}
					// TODO: Find out how to handle the OS Data

					// Add an entry to the results.
					remoteAddress, err2 := url.Parse("udp://" + strconv.Itoa(int(remoteAmsNetId.GetOctet1())) + "." +
						strconv.Itoa(int(remoteAmsNetId.GetOctet2())) + "." +
						strconv.Itoa(int(remoteAmsNetId.GetOctet3())) + "." +
						strconv.Itoa(int(remoteAmsNetId.GetOctet4())) + ":" +
						strconv.Itoa(int(driverModel.AdsConstants_ADSTCPDEFAULTPORT)))
					if err2 == nil {
						plcDiscoveryItem := &internalModel.DefaultPlcDiscoveryItem{
							ProtocolCode:  "ads",
							TransportCode: "tcp",
							TransportUrl:  *remoteAddress,
							Options:       opts,
							Name:          hostNameBlock.GetHostName().GetText(),
							Attributes:    attributes,
						}

						// Pass the event back to the callback
						callback(plcDiscoveryItem)
					}
				}
			}
		}
	}()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Find out which interfaces to use for sending out search requests

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

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Send out search requests on all selected interfaces

	// Iterate over all network devices of this system.
	for _, interf := range interfaces {
		addrs, err := interf.Addrs()
		if err != nil {
			return err
		}
		// Iterate over all addresses the current interface has configured
		// For ADS we're only interested in IPv4 addresses, as it doesn't
		// seem to work with IPv6.
		for _, addr := range addrs {
			var ipv4Addr net.IP
			switch addr.(type) {
			// If the device is configured to communicate with a subnet
			case *net.IPNet:
				ipv4Addr = addr.(*net.IPNet).IP.To4()

			// If the device is configured for a point-to-point connection
			case *net.IPAddr:
				ipv4Addr = addr.(*net.IPAddr).IP.To4()
			}

			// If we found an IPv4 address and this is not a loopback address,
			// add it to the list of devices we will open ports and send discovery
			// messages from.
			if ipv4Addr != nil && !ipv4Addr.IsLoopback() {
				// Calculate the broadcast address for this interface
				broadcastAddress := make(net.IP, len(ipv4Addr))
				binary.BigEndian.PutUint32(broadcastAddress, binary.BigEndian.Uint32(ipv4Addr)|^binary.BigEndian.Uint32(net.IP(addr.(*net.IPNet).Mask).To4()))

				// Prepare the discovery packet data
				// Create the discovery request message for this device.
				amsNetId := model.NewAmsNetId(ipv4Addr[0], ipv4Addr[1], ipv4Addr[2], ipv4Addr[3], uint8(1), uint8(1))
				discoveryRequestMessage := model.NewAdsDiscovery(0, model.Operation_DISCOVERY_REQUEST, amsNetId, model.AdsPortNumbers_SYSTEM_SERVICE, []model.AdsDiscoveryBlock{})

				// Serialize the message
				buffer := new(bytes.Buffer)
				buffer.Grow(int(discoveryRequestMessage.GetLengthInBytes()))
				writeBuffer := utils.NewCustomWriteBufferByteBased(buffer, binary.LittleEndian)
				discoveryRequestMessage.Serialize(writeBuffer)

				// Create a not-connected UDP connection to the broadcast address
				requestAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", broadcastAddress.String(), model.AdsDiscoveryConstants_ADSDISCOVERYUDPDEFAULTPORT))
				if err != nil {
					log.Error().Err(err).Str("broadcast-ip", broadcastAddress.String()).Msg("Error resolving target socket for broadcast search")
					continue
				}
				/*localAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ipv4Addr.String(), model.AdsDiscoveryConstants_ADSDISCOVERYUDPDEFAULTPORT))
				if err != nil {
					log.Error().Err(err).Str("local-ip", ipv4Addr.String()).Msg("Error resolving local address for broadcast search")
					continue
				}
				udp, err := net.DialUDP("udp4", localAddr, requestAddr)
				if err != nil {
					log.Error().Err(err).Str("local-ip", ipv4Addr.String()).Str("broadcast-ip", broadcastAddress.String()).
						Msg("Error creating sending udp socket for broadcast search")
					continue
				}*/

				// Send out the message.
				_, err = socket.WriteTo(writeBuffer.GetBytes(), requestAddr)
				if err != nil {
					log.Error().Err(err).Str("broadcast-ip", broadcastAddress.String()).Msg("Error sending request for broadcast search")
					continue
				}
			}
		}
	}

	time.Sleep(time.Second * 10)
	return nil
}
