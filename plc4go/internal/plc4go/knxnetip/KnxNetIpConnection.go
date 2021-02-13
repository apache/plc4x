//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package knxnetip

import (
	"bytes"
	"errors"
	"fmt"
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/interceptors"
	internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports/udp"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	values2 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	log "github.com/sirupsen/logrus"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ConnectionMetadata struct {
	KnxMedium         driverModel.KnxMedium
	GatewayName       string
	GatewayKnxAddress string
	ClientKnxAddress  string

	ProjectNumber          uint8
	InstallationNumber     uint8
	DeviceSerialNumber     []int8
	DeviceMulticastAddress []int8
	DeviceMacAddress       []int8
	SupportedServices      []string

	apiModel.PlcConnectionMetadata
}

func (m ConnectionMetadata) GetConnectionAttributes() map[string]string {
	return map[string]string{
		"KnxMedium":         m.KnxMedium.String(),
		"GatewayName":       m.GatewayName,
		"GatewayKnxAddress": m.GatewayKnxAddress,
		"ClientKnxAddress":  m.ClientKnxAddress,

		"ProjectNumber":          strconv.Itoa(int(m.ProjectNumber)),
		"InstallationNumber":     strconv.Itoa(int(m.InstallationNumber)),
		"DeviceSerialNumber":     utils.Int8ArrayToString(m.DeviceSerialNumber, " "),
		"DeviceMulticastAddress": utils.Int8ArrayToString(m.DeviceSerialNumber, "."),
		"DeviceMacAddress":       utils.Int8ArrayToString(m.DeviceSerialNumber, ":"),
		"SupportedServices":      strings.Join(m.SupportedServices, ", "),
	}
}

func (m ConnectionMetadata) CanRead() bool {
	return true
}

func (m ConnectionMetadata) CanWrite() bool {
	return true
}

func (m ConnectionMetadata) CanSubscribe() bool {
	return true
}

func (m ConnectionMetadata) CanBrowse() bool {
	return true
}

type KnxDeviceConnection struct {
	counter          uint8
	deviceDescriptor uint16
	maxApdu          uint16
}

type KnxMemoryReadFragment struct {
	numElements     uint8
	startingAddress uint16
}

type KnxNetIpConnection struct {
	messageCodec             spi.MessageCodec
	options                  map[string][]string
	fieldHandler             spi.PlcFieldHandler
	valueHandler             spi.PlcValueHandler
	connectionStateTimer     *time.Ticker
	quitConnectionStateTimer chan struct{}
	subscribers              []*KnxNetIpSubscriber
	leve3AddressCache        map[uint16]*driverModel.KnxGroupAddress3Level
	leve2AddressCache        map[uint16]*driverModel.KnxGroupAddress2Level
	leve1AddressCache        map[uint16]*driverModel.KnxGroupAddressFreeLevel

	valueCache      map[uint16][]int8
	valueCacheMutex sync.RWMutex
	metadata        *ConnectionMetadata
	defaultTtl      time.Duration

	GatewayKnxAddress             *driverModel.KnxAddress
	ClientKnxAddress              *driverModel.KnxAddress
	CommunicationChannelId        uint8
	SequenceCounter               int32
	TunnelingRequestExpectationId int32
	DeviceConnections             map[driverModel.KnxAddress]*KnxDeviceConnection

	requestInterceptor internalModel.RequestInterceptor
	plc4go.PlcConnection
}

type KnxReadResult struct {
	returnCode apiModel.PlcResponseCode
	value      *values.PlcValue
	numItems   uint8
}

type InternalResult struct {
	responseMessage interface{}
	err             error
}

func NewKnxNetIpConnection(transportInstance transports.TransportInstance, options map[string][]string, fieldHandler spi.PlcFieldHandler) *KnxNetIpConnection {
	connection := &KnxNetIpConnection{
		options:            options,
		fieldHandler:       fieldHandler,
		valueHandler:       NewValueHandler(),
		requestInterceptor: interceptors.NewSingleItemRequestInterceptor(),
		subscribers:        []*KnxNetIpSubscriber{},
		leve3AddressCache:  map[uint16]*driverModel.KnxGroupAddress3Level{},
		leve2AddressCache:  map[uint16]*driverModel.KnxGroupAddress2Level{},
		leve1AddressCache:  map[uint16]*driverModel.KnxGroupAddressFreeLevel{},
		valueCache:         map[uint16][]int8{},
		valueCacheMutex:    sync.RWMutex{},
		metadata:           &ConnectionMetadata{},
		defaultTtl:         time.Second * 10,
		DeviceConnections:  map[driverModel.KnxAddress]*KnxDeviceConnection{},
	}
	connection.messageCodec = NewKnxNetIpMessageCodec(transportInstance, connection.interceptIncomingMessage)
	return connection
}

func (m *KnxNetIpConnection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	result := make(chan plc4go.PlcConnectionConnectResult)
	sendResult := func(connection plc4go.PlcConnection, err error) {
		result <- plc4go.NewPlcConnectionConnectResult(connection, err)
	}

	go func() {
		err := m.messageCodec.Connect()
		if err != nil {
			sendResult(nil, errors.New("error opening connection"))
			return
		}

		searchResponseChannel := m.sendGatewaySearchRequest()
		select {
		case searchResponse := <-searchResponseChannel:
			if searchResponse.err != nil {
				sendResult(nil, errors.New("error discovering device capabilities"))
				return
			}

			searchResponseMessage := driverModel.CastSearchResponse(searchResponse.responseMessage)
			// Save some important information
			m.metadata.KnxMedium = searchResponseMessage.DibDeviceInfo.KnxMedium
			m.metadata.GatewayName = string(bytes.Trim(utils.Int8ArrayToByteArray(
				searchResponseMessage.DibDeviceInfo.DeviceFriendlyName), "\x00"))
			m.GatewayKnxAddress = searchResponseMessage.DibDeviceInfo.KnxAddress
			m.metadata.GatewayKnxAddress = KnxAddressToString(m.GatewayKnxAddress)
			m.metadata.ProjectNumber = searchResponseMessage.DibDeviceInfo.ProjectInstallationIdentifier.ProjectNumber
			m.metadata.InstallationNumber = searchResponseMessage.DibDeviceInfo.ProjectInstallationIdentifier.InstallationNumber
			m.metadata.DeviceSerialNumber = searchResponseMessage.DibDeviceInfo.KnxNetIpDeviceSerialNumber
			m.metadata.DeviceMulticastAddress = searchResponseMessage.DibDeviceInfo.KnxNetIpDeviceMulticastAddress.Addr
			m.metadata.DeviceMacAddress = searchResponseMessage.DibDeviceInfo.KnxNetIpDeviceMacAddress.Addr
			m.metadata.SupportedServices = []string{}
			supportsTunneling := false
			for _, serviceId := range searchResponseMessage.DibSuppSvcFamilies.ServiceIds {
				m.metadata.SupportedServices = append(m.metadata.SupportedServices, serviceId.Child.GetTypeName())
				// If this is an instance of the "tunneling", service, this connection supports tunneling
				_, ok := serviceId.Child.(*driverModel.KnxNetIpTunneling)
				if ok {
					supportsTunneling = true
					break
				}
			}

			// If the current device supports tunneling, create a tunneling connection.
			// Via this connection we then get access to the entire KNX network this Gateway is connected to.
			if supportsTunneling {
				// As soon as we got a successful search-response back, send a connection request.
				connectionResponseChannel := m.sendGatewayConnectionRequest()
				select {
				case connectionResponse := <-connectionResponseChannel:
					if connectionResponse.err != nil {
						sendResult(nil, errors.New("error connecting to device"))
						return
					}

					connectionResponseMessage := driverModel.CastConnectionResponse(connectionResponse.responseMessage)
					// Save the communication channel id
					m.CommunicationChannelId = connectionResponseMessage.CommunicationChannelId

					// Reset the sequence counter
					m.SequenceCounter = -1

					// If the connection was successful, the gateway will now forward any packets
					// on the KNX bus that are broadcast packets to us, so we have to setup things
					// to handle these incoming messages.
					switch connectionResponseMessage.Status {
					case driverModel.Status_NO_ERROR:
						// Save the KNX Address the Gateway assigned to us for this connection.
						tunnelConnectionDataBlock :=
							driverModel.CastConnectionResponseDataBlockTunnelConnection(
								connectionResponseMessage.ConnectionResponseDataBlock)
						m.ClientKnxAddress = tunnelConnectionDataBlock.KnxAddress

						// Start a timer that sends connection-state requests every 60 seconds
						/*log.Infof("Starting Keep-Alive Timer")
												m.connectionStateTimer = time.NewTicker(60 * time.Second)
						                        m.quitConnectionStateTimer = make(chan struct{})
												go func() {
													for {
														select {
														case <-m.connectionStateTimer.C:
															log.Infof("Executing Keep-Alive action")
															// We're using the connection-state-request as ping operation ...
															ping := m.Ping()
															select {
															case pingResult := <-ping:
																if pingResult.Err != nil {
																	// TODO: Do some error handling here ...
																	m.connectionStateTimer.Stop()
																}
															case <-time.After(m.defaultTtl * 2):
																// Close the connection
																m.Close()
															}

														// If externally a request to stop the timer was issued, stop the timer.
														case <-m.quitConnectionStateTimer:
															// TODO: Do some error handling here ...
															if m.connectionStateTimer != nil {
						                                        m.connectionStateTimer.Stop()
						                                    }
															return
														}
													}
												}()*/

						// Create a go routine to handle incoming tunneling-requests which haven't been
						// handled by any other handler. This is where usually the GroupValueWrite messages
						// are being handled.
						go func() {
							defaultIncomingMessageChannel := m.messageCodec.GetDefaultIncomingMessageChannel()
							for {
								incomingMessage := <-defaultIncomingMessageChannel
								tunnelingRequest := driverModel.CastTunnelingRequest(incomingMessage)
								if tunnelingRequest == nil {
									tunnelingResponse := driverModel.CastTunnelingResponse(incomingMessage)
									if tunnelingResponse != nil {
										log.Warnf("Got an unhandled TunnelingResponse message %v\n", tunnelingResponse)
									} else {
										log.Warnf("Not a TunnelingRequest or TunnelingResponse message %v\n", incomingMessage)
									}
								} else {
									if tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
										log.Warnf("Not for this connection %v\n", tunnelingRequest)
										continue
									}

									lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
									if lDataInd != nil {
										// Get APDU, source and target address
										lDataFrameData := driverModel.CastLDataExtended(lDataInd.DataFrame)
										sourceAddress := lDataFrameData.SourceAddress

										// If this is not an APDU, there is no need to further handle it.
										if lDataFrameData.Apdu == nil {
											continue
										}

										// If this is an incoming disconnect request, remove the device
										// from the device connections, otherwise handle it as normal
										// incoming message.
										apduControlContainer := driverModel.CastApduControlContainer(lDataFrameData.Apdu)
										if apduControlContainer != nil {
											disconnectApdu := driverModel.CastApduControlDisconnect(apduControlContainer.ControlApdu)
											if disconnectApdu != nil {
												if m.DeviceConnections[*sourceAddress] != nil /* && m.ClientKnxAddress == Int8ArrayToKnxAddress(targetAddress)*/ {
													// Remove the connection
													delete(m.DeviceConnections, *sourceAddress)
												}
											}
										} else {
											m.handleIncomingTunnelingRequest(tunnelingRequest)
										}
									}
								}
							}
						}()

						// Fire the "connected" event
						sendResult(m, nil)
					case driverModel.Status_NO_MORE_CONNECTIONS:
						sendResult(nil, errors.New("no more connections"))
					}
				case <-time.After(m.defaultTtl * 2):
					log.Debug("Timeout receiving connection response from gateway.")
				}
			}
		case <-time.After(m.defaultTtl * 3):
			log.Debug("Timeout receiving search request from gateway.")
		}
	}()
	return result
}

func (m *KnxNetIpConnection) Close() <-chan plc4go.PlcConnectionCloseResult {
	result := make(chan plc4go.PlcConnectionCloseResult)

	go func() {
		// Stop the connection-state checker.
		if m.connectionStateTimer != nil {
			m.connectionStateTimer.Stop()
		}

		// Disconnect from all knx devices we are still connected to.
		for targetAddress := range m.DeviceConnections {
			disconnects := m.DisconnectFromDevice(targetAddress)
			select {
			case _ = <-disconnects:
			case <-time.After(m.defaultTtl * 2):
				log.Debugf("Timeout disconnecting from device %s.", KnxAddressToString(&targetAddress))
			}
		}

		// Send a disconnect request from the gateway.
		disconnectionResponseChannel := m.sendGatewayDisconnectionRequest()
		select {
		case disconnectResponse := <-disconnectionResponseChannel:
			if disconnectResponse.err != nil {
				result <- plc4go.NewPlcConnectionCloseResult(m, errors.New(
					fmt.Sprintf("got an error when disconnecting: %s", disconnectResponse.err.Error())))
			} else {
				result <- plc4go.NewPlcConnectionCloseResult(m, nil)
			}

		case <-time.After(m.defaultTtl * 2):
			log.Debug("Timeout disconnecting from gateway.")
			result <- plc4go.NewPlcConnectionCloseResult(m, nil)
		}
	}()

	return result
}

func (m *KnxNetIpConnection) IsConnected() bool {
	if m.messageCodec != nil {
		pingChannel := m.Ping()
		select {
		case pingResponse := <-pingChannel:
			return pingResponse.Err == nil
		case <-time.After(m.defaultTtl * 2):
			log.Debug("Timeout checking if the connection is alive.")
			return false
		}
	}
	return false
}

func (m *KnxNetIpConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
	result := make(chan plc4go.PlcConnectionPingResult)

	go func() {
		// Send the connection state request
		connectionStateResponseChannel := m.sendConnectionStateRequest()
		select {
		case connectionStateResponse := <-connectionStateResponseChannel:
			if connectionStateResponse.err != nil {
				result <- plc4go.NewPlcConnectionPingResult(errors.New(
					fmt.Sprintf("got an error: %s", connectionStateResponse.err.Error())))
			} else {
				result <- plc4go.NewPlcConnectionPingResult(nil)
			}
		case <-time.After(m.defaultTtl * 2):
			log.Debug("Timeout executing connection-state-request.")
		}
		return
	}()

	return result
}

func (m *KnxNetIpConnection) ConnectToDevice(targetAddress driverModel.KnxAddress) <-chan *KnxDeviceConnection {
	result := make(chan *KnxDeviceConnection)

	// If we're already connected, use that connection instead.
	if connection, ok := m.DeviceConnections[targetAddress]; ok {
		result <- connection
		return result
	}

	// First send a connection request
	go func() {
		controlConnectResponseChannel := m.sendDeviceConnectionRequest(targetAddress)
		select {
		case conResult := <-controlConnectResponseChannel:
			if conResult.err != nil {
				result <- nil
				return
			}

			// If the connection request was successful, try to read the device-descriptor
			deviceDescriptorResponses := m.sendDeviceDeviceDescriptorReadRequest(targetAddress)
			go func() {
				select {
				case _ = <-deviceDescriptorResponses:
					// Last, not least, read the max APDU size
					propertyValueResponses := m.sendDevicePropertyReadRequest(targetAddress, 0, 56)
					go func() {
						select {
						case propertyValueResponse := <-propertyValueResponses:
							// If we were able to read the max APDU size, then use the minimum of
							// the connection APDU size and the device APDU size, otherwise use the
							// default APDU Size of 15
							// Defined in: 03_05_01 Resources v01.09.03 AS Page 40
							deviceApduSize := uint16(15)
							if propertyValueResponse.returnCode == apiModel.PlcResponseCode_OK {
								deviceApduSize = (*propertyValueResponse.value).GetUint16()
							}
							connection := m.DeviceConnections[targetAddress]
							connection.maxApdu = uint16(math.Min(float64(deviceApduSize), 240))
							result <- connection
						case <-time.After(m.defaultTtl * 2):
							result <- nil
							log.Debugf("Timeout receiving property read response from device: %s", KnxAddressToString(&targetAddress))
						}
					}()
				}
			}()
		case <-time.After(m.defaultTtl * 3):
			result <- nil
			log.Debugf("Timeout connecting to device: %s", KnxAddressToString(&targetAddress))
		}
	}()

	return result
}

func (m *KnxNetIpConnection) DisconnectFromDevice(targetAddress driverModel.KnxAddress) <-chan *KnxDeviceConnection {
	result := make(chan *KnxDeviceConnection)

	if connection, ok := m.DeviceConnections[targetAddress]; ok {
		disconnects := m.sendDeviceDisconnectionRequest(targetAddress)
		go func() {
			select {
			case _ = <-disconnects:
				// Remove the connection from the list.
				delete(m.DeviceConnections, targetAddress)
				result <- connection
			case <-time.After(m.defaultTtl * 2):
				log.Debugf("Timeout disconnecting from device: %s", KnxAddressToString(&targetAddress))
			}
		}()
	} else {
		result <- nil
	}

	return result
}

func (m *KnxNetIpConnection) ReadDeviceProperty(targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	go func() {
		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.ConnectToDevice(targetAddress)
			select {
			case connection = <-connections:
				// If we didn't get a connect, abort
				if connection == nil {
					result <- KnxReadResult{
						returnCode: apiModel.PlcResponseCode_INVALID_ADDRESS,
						value:      nil,
						numItems:   0,
					}
				}
			case <-time.After(m.defaultTtl * 2):
				log.Debugf("Timeout connecting to device: %s", KnxAddressToString(&targetAddress))
			}
		}

		// If we successfully got a connection, read the property
		if connection != nil {
			propertyValueResponses := m.sendDevicePropertyReadRequest(targetAddress, objectId, propertyId)
			select {
			case propertyValueResponse := <-propertyValueResponses:
				result <- propertyValueResponse
			case <-time.After(m.defaultTtl * 2):
				log.Debugf("Timeout reading device property from device: %s: ObjectId %d, PropertyId %d",
					KnxAddressToString(&targetAddress), objectId, propertyId)
			}
		}
	}()

	return result
}

func (m *KnxNetIpConnection) ReadDeviceMemory(targetAddress driverModel.KnxAddress, address uint16, numElements uint8, datapointType *driverModel.KnxDatapointType) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	if datapointType == nil {
		// This is a 0-255 valued 1-byte value.
		dpt := driverModel.KnxDatapointType_DPT_DecimalFactor
		datapointType = &dpt
	}

	go func() {
		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.ConnectToDevice(targetAddress)
			select {
			case connection = <-connections:
				// If we didn't get a connect, abort
				if connection == nil {
					result <- KnxReadResult{
						returnCode: apiModel.PlcResponseCode_INVALID_ADDRESS,
						value:      nil,
						numItems:   0,
					}
				}
			case <-time.After(m.defaultTtl * 2):
				log.Debugf("Timeout connecting to device: %s", KnxAddressToString(&targetAddress))
			}
		}

		// If we successfully got a connection, read the property
		if connection != nil {
			// Depending on the gateway Max APDU and the device Max APDU, split this up into multiple requests.
			// An APDU starts with the last 6 bits of the first data byte containing the count
			// followed by the 16-bit address, so these are already used.
			elementSize := datapointType.DatapointMainType().SizeInBits() / 8
			remainingRequestElements := numElements
			curStartingAddress := address
			var results []KnxReadResult
			for remainingRequestElements > 0 {
				// As the maxApdu can change, we have to do this in the loop.
				maxNumBytes := uint8(math.Min(float64(connection.maxApdu-3), float64(63)))
				maxNumElementsPerRequest := uint8(math.Floor(float64(maxNumBytes / elementSize)))
				numElements := uint8(math.Min(float64(remainingRequestElements), float64(maxNumElementsPerRequest)))
				propertyValueResponses := m.sendDeviceMemoryReadRequest(targetAddress, curStartingAddress, numElements, *datapointType)
				select {
				case propertyValueResponse := <-propertyValueResponses:
					results = append(results, propertyValueResponse)
					// Update the reading position.
					remainingRequestElements = remainingRequestElements - propertyValueResponse.numItems
					curStartingAddress = curStartingAddress + uint16(propertyValueResponse.numItems*elementSize)
				case <-time.After(m.defaultTtl * 4):
					log.Debugf("Timeout reading device memory from device: %s: address %d, number of elements %d of type %s",
						KnxAddressToString(&targetAddress), address, numElements, datapointType.Name())
					// TODO: Return an error
				}
			}
			if len(results) > 1 {
				returnCode := apiModel.PlcResponseCode_OK
				// Merge the results
				var plcValues []values.PlcValue
				for _, curResult := range results {
					if curResult.returnCode > returnCode {
						returnCode = curResult.returnCode
					}
					plcValues = append(plcValues, (*(curResult.value)).GetList()...)
				}
				var plcList values.PlcValue
				plcList = values2.NewPlcList(plcValues)
				result <- KnxReadResult{
					returnCode: returnCode,
					value:      &plcList,
					numItems:   numElements,
				}
			} else if len(results) == 1 {
				result <- results[0]
			}
		}
	}()

	return result
}

func (m *KnxNetIpConnection) GetMetadata() apiModel.PlcConnectionMetadata {
	return m.metadata
}

func (m *KnxNetIpConnection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilder(
		m.fieldHandler, NewKnxNetIpReader(m))
}

func (m *KnxNetIpConnection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(
		m.fieldHandler, m.valueHandler, NewKnxNetIpWriter(m.messageCodec))
}

func (m *KnxNetIpConnection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return internalModel.NewDefaultPlcSubscriptionRequestBuilder(
		m.fieldHandler, m.valueHandler, NewKnxNetIpSubscriber(m))
}

func (m *KnxNetIpConnection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	return internalModel.NewDefaultPlcBrowseRequestBuilder(NewKnxNetIpBrowser(m, m.messageCodec))
}

func (m *KnxNetIpConnection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return nil /*internalModel.NewDefaultPlcUnsubscriptionRequestBuilder(
	  m.fieldHandler, m.valueHandler, NewKnxNetIpSubscriber(m.messageCodec))*/
}

func (m *KnxNetIpConnection) GetTransportInstance() transports.TransportInstance {
	if mc, ok := m.messageCodec.(spi.TransportInstanceExposer); ok {
		return mc.GetTransportInstance()
	}
	return nil
}

func (m *KnxNetIpConnection) GetPlcFieldHandler() spi.PlcFieldHandler {
	return m.fieldHandler
}

func (m *KnxNetIpConnection) GetPlcValueHandler() spi.PlcValueHandler {
	return m.valueHandler
}

func (m *KnxNetIpConnection) send(request *driverModel.KnxNetIpMessage) error {
	// If this is a tunneling request, we need to update the communicationChannelId and assign a sequenceCounter
	tunnelingRequest := driverModel.CastTunnelingRequest(request)
	if tunnelingRequest != nil {
		tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId = m.CommunicationChannelId
		tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter = m.getNewSequenceCounter()
	}
	return m.messageCodec.Send(request)
}

func (m *KnxNetIpConnection) sendGatewaySearchRequest() chan InternalResult {
	result := make(chan InternalResult)
	localAddress, err := m.getLocalAddress()
	if err != nil {
		//		close(result)
		return result
	}
	localAddr := driverModel.NewIPAddress(utils.ByteArrayToInt8Array(localAddress.IP))
	discoveryEndpoint := driverModel.NewHPAIDiscoveryEndpoint(
		driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port))
	searchRequest := driverModel.NewSearchRequest(discoveryEndpoint)
	err = m.messageCodec.SendRequest(searchRequest,
		func(message interface{}) bool {
			searchResponse := driverModel.CastSearchResponse(message)
			return searchResponse != nil
		},
		func(message interface{}) error {
			searchResponse := driverModel.CastSearchResponse(message)
			result <- InternalResult{
				responseMessage: *searchResponse,
			}
			//			close(result)
			return nil
		},
		func(err error) error {
			result <- InternalResult{
				err: errors.New(fmt.Sprintf("got error processing request: %s", err)),
			}
			return nil
		},
		m.defaultTtl)
	if err != nil {
		result <- InternalResult{
			err: errors.New(fmt.Sprintf("got error sending request: %s", err)),
		}
	}
	return result
}

func (m *KnxNetIpConnection) sendGatewayConnectionRequest() chan InternalResult {
	result := make(chan InternalResult)
	localAddress, err := m.getLocalAddress()
	if err != nil {
		//		close(result)
		result <- InternalResult{
			err: err,
		}
		return result
	}
	localAddr := driverModel.NewIPAddress(utils.ByteArrayToInt8Array(localAddress.IP)[len(localAddress.IP)-4:])
	connectionRequest := driverModel.NewConnectionRequest(
		driverModel.NewHPAIDiscoveryEndpoint(driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port)),
		driverModel.NewHPAIDataEndpoint(driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port)),
		driverModel.NewConnectionRequestInformationTunnelConnection(driverModel.KnxLayer_TUNNEL_LINK_LAYER),
	)
	err = m.messageCodec.SendRequest(connectionRequest,
		func(message interface{}) bool {
			connectionResponse := driverModel.CastConnectionResponse(message)
			return connectionResponse != nil
		},
		func(message interface{}) error {
			connectionResponse := driverModel.CastConnectionResponse(message)
			result <- InternalResult{
				responseMessage: *connectionResponse,
			}
			//			close(result)
			return nil
		},
		func(err error) error {
			result <- InternalResult{
				err: errors.New(fmt.Sprintf("got error processing request: %s", err)),
			}
			return nil
		},
		m.defaultTtl)
	if err != nil {
		result <- InternalResult{
			err: errors.New(fmt.Sprintf("got error sending request: %s", err)),
		}
	}
	return result
}

func (m *KnxNetIpConnection) sendGatewayDisconnectionRequest() chan InternalResult {
	result := make(chan InternalResult)
	localAddress, err := m.getLocalAddress()
	if err != nil {
		//		close(result)
		return result
	}
	localAddr := driverModel.NewIPAddress(utils.ByteArrayToInt8Array(localAddress.IP)[len(localAddress.IP)-4:])
	disconnectRequest := driverModel.NewDisconnectRequest(
		m.CommunicationChannelId,
		driverModel.NewHPAIControlEndpoint(
			driverModel.HostProtocolCode_IPV4_UDP,
			localAddr,
			uint16(localAddress.Port)))
	err = m.messageCodec.SendRequest(disconnectRequest,
		func(message interface{}) bool {
			disconnectResponse := driverModel.CastDisconnectResponse(message)
			return disconnectResponse != nil
		},
		func(message interface{}) error {
			disconnectResponse := driverModel.CastDisconnectResponse(message)
			result <- InternalResult{
				responseMessage: *disconnectResponse,
			}
			//			close(result)
			return nil
		},
		func(err error) error {
			result <- InternalResult{
				err: errors.New(fmt.Sprintf("got error processing request: %s", err)),
			}
			return nil
		},
		m.defaultTtl)
	if err != nil {
		result <- InternalResult{
			err: errors.New(fmt.Sprintf("got error sending request: %s", err)),
		}
	}
	return result
}

func (m *KnxNetIpConnection) sendConnectionStateRequest() chan InternalResult {
	result := make(chan InternalResult)
	localAddress, err := m.getLocalAddress()
	if err != nil {
		//		close(result)
		return result
	}
	localAddr := driverModel.NewIPAddress(utils.ByteArrayToInt8Array(localAddress.IP)[len(localAddress.IP)-4:])
	connectionStateRequest := driverModel.NewConnectionStateRequest(
		m.CommunicationChannelId,
		driverModel.NewHPAIControlEndpoint(
			driverModel.HostProtocolCode_IPV4_UDP,
			localAddr, uint16(localAddress.Port)))
	err = m.messageCodec.SendRequest(connectionStateRequest,
		func(message interface{}) bool {
			connectionStateResponse := driverModel.CastConnectionStateResponse(message)
			return connectionStateResponse != nil
		},
		func(message interface{}) error {
			connectionStateResponse := driverModel.CastConnectionStateResponse(message)
			result <- InternalResult{
				responseMessage: *connectionStateResponse,
			}
			//			close(result)
			return nil
		},
		func(err error) error {
			result <- InternalResult{
				err: errors.New(fmt.Sprintf("got error processing request: %s", err)),
			}
			return nil
		},
		m.defaultTtl)
	if err != nil {
		result <- InternalResult{
			err: errors.New(fmt.Sprintf("got error sending request: %s", err)),
		}
	}
	return result
}

func (m *KnxNetIpConnection) sendDeviceConnectionRequest(targetAddress driverModel.KnxAddress) chan InternalResult {
	result := make(chan InternalResult)

	go func() {
		// Send a connection request to the individual KNX device
		deviceConnectionRequest := driverModel.NewTunnelingRequest(
			driverModel.NewTunnelingRequestDataBlock(0, 0),
			driverModel.NewLDataReq(0, nil,
				driverModel.NewLDataExtended(false, 6, uint8(0),
					driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
					driverModel.NewApduControlContainer(driverModel.NewApduControlConnect(), false, 0),
					true, true, driverModel.CEMIPriority_SYSTEM, false, false)))
		err := m.sendRequest(
			deviceConnectionRequest,
			// The Gateway is now supposed to send an Ack to this request.
			func(message interface{}) bool {
				tunnelingRequest := driverModel.CastTunnelingRequest(message)
				if tunnelingRequest == nil ||
					tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
					return false
				}
				lDataCon := driverModel.CastLDataCon(tunnelingRequest.Cemi)
				if lDataCon == nil {
					return false
				}
				lDataFrameExt := driverModel.CastLDataExtended(lDataCon.DataFrame)
				if lDataFrameExt == nil {
					return false
				}
				// Check if the address matches
				if *Int8ArrayToKnxAddress(lDataFrameExt.DestinationAddress) != targetAddress {
					return false
				}
				apduControlContainer := driverModel.CastApduControlContainer(lDataFrameExt.Apdu)
				if apduControlContainer == nil {
					return false
				}
				apduControlConnect := driverModel.CastApduControlConnect(apduControlContainer.ControlApdu)
				return apduControlConnect != nil
			},
			func(message interface{}) error {
				tunnelingRequest := driverModel.CastTunnelingRequest(message)
				lDataCon := driverModel.CastLDataCon(tunnelingRequest.Cemi)
				// If the error flag is set, there was an error connecting
				if lDataCon.DataFrame.ErrorFlag {
					result <- InternalResult{
						err: errors.New("error flag set"),
					}
					return nil
				}
				lDataFrameExt := driverModel.CastLDataExtended(lDataCon.DataFrame)
				apduControlContainer := driverModel.CastApduControlContainer(lDataFrameExt.Apdu)
				apduControlConnect := driverModel.CastApduControlConnect(apduControlContainer.ControlApdu)

				// Create a new connection object and save it in the map
				deviceConnection := &KnxDeviceConnection{
					counter: 0,
					// I was told this value on the knx-forum.
					// Seems the max payload is 3 bytes less ...
					maxApdu: 0, // This is the default max APDU Size
				}
				m.DeviceConnections[targetAddress] = deviceConnection

				result <- InternalResult{
					responseMessage: apduControlConnect,
				}
				return nil
			},
			func(err error) error {
				result <- InternalResult{
					err: errors.New(fmt.Sprintf("got error processing request: %s", err)),
				}
				return nil
			},
			m.defaultTtl)
		if err != nil {
			result <- InternalResult{
				err: errors.New(fmt.Sprintf("got error sending request: %s", err)),
			}
		}
	}()

	return result
}

func (m *KnxNetIpConnection) sendDeviceDisconnectionRequest(targetAddress driverModel.KnxAddress) chan InternalResult {
	result := make(chan InternalResult)

	go func() {
		// Send a connection request to the individual KNX device
		deviceDisconnectionRequest := driverModel.NewTunnelingRequest(
			driverModel.NewTunnelingRequestDataBlock(0, 0),
			driverModel.NewLDataReq(0, nil,
				driverModel.NewLDataExtended(false, 6, uint8(0),
					driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
					driverModel.NewApduControlContainer(driverModel.NewApduControlDisconnect(), false, 0),
					true, true, driverModel.CEMIPriority_SYSTEM, false, false)))
		err := m.sendRequest(
			deviceDisconnectionRequest,
			// The Gateway is now supposed to send an Ack to this request.
			func(message interface{}) bool {
				tunnelingRequest := driverModel.CastTunnelingRequest(message)
				if tunnelingRequest == nil ||
					tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
					return false
				}
				lDataCon := driverModel.CastLDataCon(tunnelingRequest.Cemi)
				if lDataCon == nil {
					return false
				}
				dataFrameExt := driverModel.CastLDataExtended(lDataCon.DataFrame)
				if dataFrameExt == nil {
					return false
				}
				curTargetAddress := Int8ArrayToKnxAddress(dataFrameExt.DestinationAddress)
				// Check if the address matches
				if *curTargetAddress != targetAddress {
					return false
				}
				apduControlContainer := driverModel.CastApduControlContainer(dataFrameExt.Apdu)
				if apduControlContainer == nil {
					return false
				}
				apduControlDisconnect := driverModel.CastApduControlDisconnect(apduControlContainer.ControlApdu)
				return apduControlDisconnect != nil
			},
			func(message interface{}) error {
				tunnelingRequest := driverModel.CastTunnelingRequest(message)
				lDataCon := driverModel.CastLDataCon(tunnelingRequest.Cemi)
				// If the error flag is set, there was an error connecting
				if lDataCon.DataFrame.ErrorFlag {
					//					close(result)
				}
				dataFrameExt := driverModel.CastLDataExtended(lDataCon.DataFrame)
				apduControlContainer := driverModel.CastApduControlContainer(dataFrameExt.Apdu)
				apduControlDisconnect := driverModel.CastApduControlDisconnect(apduControlContainer.ControlApdu)
				result <- InternalResult{
					responseMessage: *apduControlDisconnect,
				}
				return nil
			},
			func(err error) error {
				result <- InternalResult{
					err: errors.New(fmt.Sprintf("got error processing request: %s", err)),
				}
				return nil
			},
			m.defaultTtl)
		if err != nil {
			result <- InternalResult{
				err: errors.New(fmt.Sprintf("got error sending request: %s", err)),
			}
		}
	}()

	return result
}

func (m *KnxNetIpConnection) sendDeviceDeviceDescriptorReadRequest(targetAddress driverModel.KnxAddress) chan InternalResult {
	result := make(chan InternalResult)

	connection, ok := m.DeviceConnections[targetAddress]
	if !ok {
		result <- InternalResult{
			err: errors.New("not connected"),
		}
		return result
	}

	// Next, read the device descriptor so we know how we have to communicate with the device.
	counter := connection.counter
	connection.counter++
	deviceDescriptorReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(
			// This is actually set in the KnxNetIpConnection.SendMessage method
			0,
			// This is actually set in the KnxNetIpConnection.SendMessage method
			0),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, uint8(0),
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataDeviceDescriptorRead(0), true, counter),
				true, true, driverModel.CEMIPriority_LOW, false, false)))
	err := m.sendRequest(
		deviceDescriptorReadRequest,
		func(message interface{}) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
				return false
			}
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
			if lDataInd == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.DataFrame)
			if dataFrameExt == nil {
				return false
			}
			// Check if the address matches
			if *dataFrameExt.SourceAddress != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.Apdu.Counter != counter {
				return false
			}
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.Apdu)
			if dataContainer == nil {
				return false
			}
			deviceDescriptorResponse := driverModel.CastApduDataDeviceDescriptorResponse(dataContainer.DataApdu)
			if deviceDescriptorResponse == nil {
				return false
			}
			return true
		},
		func(message interface{}) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
			dataFrame := driverModel.CastLDataExtended(lDataInd.DataFrame)
			dataContainer := driverModel.CastApduDataContainer(dataFrame.Apdu)
			deviceDescriptorResponse := driverModel.CastApduDataDeviceDescriptorResponse(dataContainer.DataApdu)

			// Send an ACK and wait for it to be delivered
			go func() {
				ackResults := m.sendDeviceAck(targetAddress, dataFrame.Apdu.Counter)
				select {
				case ackResult := <-ackResults:
					if ackResult.err != nil {
						fmt.Printf("Error sending Ack: %s", ackResult.err)
					}
				case <-time.After(m.defaultTtl * 2):
					log.Debugf("Timeout sending device ack to %s.", KnxAddressToString(&targetAddress))
				}

				// Save the device-descriptor value
				deviceDescriptor := uint16(deviceDescriptorResponse.Data[0])<<8 | (uint16(deviceDescriptorResponse.Data[1]) & 0xFF)
				connection.deviceDescriptor = deviceDescriptor

				result <- InternalResult{
					responseMessage: *deviceDescriptorResponse,
				}
			}()

			return nil
		},
		func(err error) error {
			result <- InternalResult{
				err: errors.New(fmt.Sprintf("got error processing request: %s", err)),
			}
			return nil
		},
		m.defaultTtl)
	if err != nil {
		result <- InternalResult{
			err: errors.New(fmt.Sprintf("got error sending request: %s", err)),
		}
	}
	return result
}

func (m *KnxNetIpConnection) sendDevicePropertyReadRequest(targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8) chan KnxReadResult {
	result := make(chan KnxReadResult)

	connection, ok := m.DeviceConnections[targetAddress]
	if !ok {
		//		close(result)
		result <- KnxReadResult{
			returnCode: apiModel.PlcResponseCode_INTERNAL_ERROR,
			value:      nil,
			numItems:   0,
		}
		return result
	}

	// Next, read the device descriptor so we know how we have to communicate with the device.
	counter := connection.counter
	connection.counter++

	// Send the property read request and wait for a confirmation that this property is readable.
	propertyReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(0, 0),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, 0,
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataOther(
						driverModel.NewApduDataExtPropertyValueRead(objectId, propertyId, 1, 1)),
					true, counter),
				true, true, driverModel.CEMIPriority_LOW, false, false)))
	err := m.sendRequest(
		propertyReadRequest,
		func(message interface{}) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
				return false
			}
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
			if lDataInd == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.DataFrame)
			if dataFrameExt == nil {
				return false
			}
			// Check if the address matches
			if *dataFrameExt.SourceAddress != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.Apdu.Counter != counter {
				return false
			}
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.Apdu)
			if dataContainer == nil {
				return false
			}
			dataApduOther := driverModel.CastApduDataOther(dataContainer.DataApdu)
			if dataApduOther == nil {
				return false
			}
			propertyValueResponse := driverModel.CastApduDataExtPropertyValueResponse(dataApduOther.ExtendedApdu)
			if propertyValueResponse == nil {
				return false
			}
			return propertyValueResponse.ObjectIndex == objectId && propertyValueResponse.PropertyId == propertyId
		},
		func(message interface{}) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.DataFrame)
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.Apdu)
			dataApduOther := driverModel.CastApduDataOther(dataContainer.DataApdu)
			propertyValueResponse := driverModel.CastApduDataExtPropertyValueResponse(dataApduOther.ExtendedApdu)

			// Send an ACK and wait for it to be delivered
			go func() {
				ackResults := m.sendDeviceAck(targetAddress, dataFrameExt.Apdu.Counter)
				select {
				case ackResult := <-ackResults:
					if ackResult.err != nil {
						log.Errorf("Error sending Ack: %s", ackResult.err)
					}
				case <-time.After(m.defaultTtl * 2):
					log.Debugf("Timeout sending device ack to %s.", KnxAddressToString(&targetAddress))
				}

				// If the count is 0, then this property doesn't exist or the user has no permission to read it.
				if propertyValueResponse.Count == 0 {
					result <- KnxReadResult{
						returnCode: apiModel.PlcResponseCode_NOT_FOUND,
						value:      nil,
						numItems:   0,
					}
				} else {
					// Find out the type of the property
					var objectType *driverModel.KnxInterfaceObjectType
					for curObjectType := driverModel.KnxInterfaceObjectType_OT_UNKNOWN; curObjectType <= driverModel.KnxInterfaceObjectType_OT_SUNBLIND_SENSOR_BASIC; curObjectType++ {
						if curObjectType.Code() == strconv.Itoa(int(objectId)) {
							objectType = &curObjectType
							break
						}
					}
					property := driverModel.KnxInterfaceObjectProperty_PID_UNKNOWN
					if objectType != nil {
						for curProperty := driverModel.KnxInterfaceObjectProperty_PID_UNKNOWN; curProperty <= driverModel.KnxInterfaceObjectProperty_PID_SUNBLIND_SENSOR_BASIC_ENABLE_TOGGLE_MODE; curProperty++ {
							if curProperty.PropertyId() == propertyId {
								if curProperty.ObjectType() == driverModel.KnxInterfaceObjectType_OT_GENERAL || curProperty.ObjectType() == *objectType {
									property = curProperty
									break
								}
							}
						}
					}

					dataLength := uint8(len(propertyValueResponse.Data))
					data := propertyValueResponse.Data
					rb := utils.NewReadBuffer(data)
					plcValue, err := driverModel.KnxPropertyParse(rb, property.PropertyDataType(), dataLength)

					// Return the result
					if err != nil {
						result <- KnxReadResult{
							returnCode: apiModel.PlcResponseCode_INTERNAL_ERROR,
							value:      nil,
							numItems:   0,
						}
					} else {
						result <- KnxReadResult{
							returnCode: apiModel.PlcResponseCode_OK,
							value:      &plcValue,
							numItems:   1,
						}
					}
				}

			}()
			return nil
		},
		func(err error) error {
			result <- KnxReadResult{
				returnCode: apiModel.PlcResponseCode_RESPONSE_PENDING,
				value:      nil,
				numItems:   0,
			}
			return nil
		},
		m.defaultTtl)

	if err != nil {
		//		close(result)
	}
	return result
}

func (m *KnxNetIpConnection) sendDeviceMemoryReadRequest(targetAddress driverModel.KnxAddress, address uint16, numElements uint8, datapointType driverModel.KnxDatapointType) chan KnxReadResult {
	result := make(chan KnxReadResult)

	connection, ok := m.DeviceConnections[targetAddress]
	if !ok {
		//		close(result)
		result <- KnxReadResult{
			returnCode: apiModel.PlcResponseCode_INTERNAL_ERROR,
			value:      nil,
			numItems:   0,
		}
		return result
	}

	// Next, read the device descriptor so we know how we have to communicate with the device.
	counter := connection.counter
	connection.counter++
	numBytes := numElements * uint8(math.Max(float64(1), float64(datapointType.DatapointMainType().SizeInBits()/8)))

	// Send the property read request and wait for a confirmation that this property is readable.
	propertyReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(0, 0),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, 0,
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataMemoryRead(numBytes, address),
					true, counter),
				true, true, driverModel.CEMIPriority_LOW, false, false)))
	err := m.sendRequest(
		propertyReadRequest,
		func(message interface{}) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
				return false
			}
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
			if lDataInd == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.DataFrame)
			if dataFrameExt == nil {
				return false
			}
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.Apdu)
			if dataContainer == nil {
				return false
			}
			dataApduMemoryResponse := driverModel.CastApduDataMemoryResponse(dataContainer.DataApdu)
			if dataApduMemoryResponse == nil {
				return false
			}

			// Check if the address matches
			if *dataFrameExt.SourceAddress != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.Apdu.Counter != counter {
				return false
			}
			return dataApduMemoryResponse.Address == address
		},
		func(message interface{}) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.DataFrame)
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.Apdu)
			dataApduMemoryResponse := driverModel.CastApduDataMemoryResponse(dataContainer.DataApdu)

			// Send an ACK and wait for it to be delivered
			go func() {
				ackResults := m.sendDeviceAck(targetAddress, dataFrameExt.Apdu.Counter)
				select {
				case ackResult := <-ackResults:
					if ackResult.err != nil {
						log.Errorf("Error sending Ack: %s", ackResult.err)
					}
				case <-time.After(m.defaultTtl * 2):
					log.Debugf("Timeout sending device ack to %s.", KnxAddressToString(&targetAddress))
				}

				// If the number of bytes read is less than expected,
				// Update the connection.maxApdu value. This is required
				// as some devices seem to be sending back less than the
				// number of bytes specified than the maxApdu.
				if uint8(len(dataApduMemoryResponse.Data)) < numBytes {
					connection.maxApdu = uint16(len(dataApduMemoryResponse.Data) + 3)
				}

				// Parse the data according to the property type information
				rb := utils.NewReadBuffer(dataApduMemoryResponse.Data)
				var plcValues []values.PlcValue
				for rb.HasMore(datapointType.DatapointMainType().SizeInBits()) {
					plcValue, err := driverModel.KnxDatapointParse(rb, datapointType)
					// Return the result
					if err != nil {
						result <- KnxReadResult{
							returnCode: apiModel.PlcResponseCode_INTERNAL_ERROR,
							value:      nil,
							numItems:   0,
						}
						return
					} else {
						plcValues = append(plcValues, plcValue)
					}
				}
				// If there are still remaining bytes, keep them for the next time.

				// If this is a single value, just return that directly.
				// If it's not, wrap it in a PlcList structure.
				if len(plcValues) == 1 {
					result <- KnxReadResult{
						returnCode: apiModel.PlcResponseCode_OK,
						value:      &plcValues[0],
						numItems:   1,
					}
				} else {
					var plcList values.PlcValue
					plcList = values2.NewPlcList(plcValues)
					result <- KnxReadResult{
						returnCode: apiModel.PlcResponseCode_OK,
						value:      &plcList,
						numItems:   uint8(len(plcValues)),
					}
				}
			}()
			return nil
		},
		func(err error) error {
			result <- KnxReadResult{
				returnCode: apiModel.PlcResponseCode_RESPONSE_PENDING,
				value:      nil,
				numItems:   0,
			}
			return nil
		},
		m.defaultTtl)

	if err != nil {
		//		close(result)
	}
	return result
}

func (m *KnxNetIpConnection) sendDeviceAck(targetAddress driverModel.KnxAddress, counter uint8) chan InternalResult {
	result := make(chan InternalResult)

	ack := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(0, 0),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, uint8(0),
				driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduControlContainer(driverModel.NewApduControlAck(), true, counter),
				true, true, driverModel.CEMIPriority_SYSTEM, false, false)))
	err := m.sendRequest(
		ack,
		func(message interface{}) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
				return false
			}
			lDataCon := driverModel.CastLDataCon(tunnelingRequest.Cemi)
			if lDataCon == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataCon.DataFrame)
			if dataFrameExt == nil {
				return false
			}
			// Check if the addresses match
			if *dataFrameExt.SourceAddress != *m.ClientKnxAddress {
				return false
			}
			curTargetAddress := Int8ArrayToKnxAddress(dataFrameExt.DestinationAddress)
			if *curTargetAddress != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.Apdu.Counter != counter {
				return false
			}
			controlContainer := driverModel.CastApduControlContainer(dataFrameExt.Apdu)
			if controlContainer == nil {
				return false
			}
			dataApduAck := driverModel.CastApduControlAck(controlContainer.ControlApdu)
			if dataApduAck == nil {
				return false
			}
			return true
		},
		func(message interface{}) error {
			result <- InternalResult{}
			return nil
		},
		func(err error) error {
			result <- InternalResult{
				err: errors.New(fmt.Sprintf("got error processing request: %s", err)),
			}
			return nil
		},
		m.defaultTtl)

	if err != nil {
		result <- InternalResult{
			err: errors.New(fmt.Sprintf("got error sending request: %s", err)),
		}
	}

	return result
}

func (m *KnxNetIpConnection) sendRequest(request *driverModel.KnxNetIpMessage, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, handleError spi.HandleError, ttl time.Duration) error {
	// If this is a tunneling request, we need to update the communicationChannelId and assign a sequenceCounter
	tunnelingRequest := driverModel.CastTunnelingRequest(request)
	if tunnelingRequest != nil {
		tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId = m.CommunicationChannelId
		tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter = m.getNewSequenceCounter()
	}
	return m.messageCodec.SendRequest(request, acceptsMessage, handleMessage, handleError, ttl)
}

func (m *KnxNetIpConnection) getLocalAddress() (*net.UDPAddr, error) {
	transportInstanceExposer, ok := m.messageCodec.(spi.TransportInstanceExposer)
	if !ok {
		return nil, errors.New("used transport, is not a TransportInstanceExposer")
	}

	// Prepare a SearchReq
	udpTransportInstance, ok := transportInstanceExposer.GetTransportInstance().(*udp.UdpTransportInstance)
	if !ok {
		return nil, errors.New("used transport, is not a UdpTransportInstance")
	}

	return udpTransportInstance.LocalAddress, nil
}

func (m *KnxNetIpConnection) interceptIncomingMessage(interface{}) {
	if m.connectionStateTimer != nil {
		// Reset the timer for sending the ConnectionStateRequest
		m.connectionStateTimer.Reset(60 * time.Second)
	}
}

func (m *KnxNetIpConnection) getNewSequenceCounter() uint8 {
	sequenceCounter := atomic.AddInt32(&m.SequenceCounter, 1)
	if sequenceCounter >= math.MaxUint8 {
		atomic.StoreInt32(&m.SequenceCounter, -1)
		sequenceCounter = -1
	}
	return uint8(sequenceCounter)
}

func (m *KnxNetIpConnection) castIpToKnxAddress(ip net.IP) *driverModel.IPAddress {
	return driverModel.NewIPAddress(utils.ByteArrayToInt8Array(ip)[len(ip)-4:])
}

func (m *KnxNetIpConnection) handleIncomingTunnelingRequest(tunnelingRequest *driverModel.TunnelingRequest) {
	go func() {
		lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi.Child)
		if lDataInd != nil {
			var destinationAddress []int8
			switch lDataInd.DataFrame.Child.(type) {
			case *driverModel.LDataExtended:
				dataFrame := driverModel.CastLDataExtended(lDataInd.DataFrame)
				destinationAddress = dataFrame.DestinationAddress
				switch dataFrame.Apdu.Child.(type) {
				case *driverModel.ApduDataContainer:
					container := driverModel.CastApduDataContainer(dataFrame.Apdu)
					switch container.DataApdu.Child.(type) {
					case *driverModel.ApduDataGroupValueWrite:
						groupValueWrite := driverModel.CastApduDataGroupValueWrite(container.DataApdu)
						if destinationAddress != nil {
							addressData := uint16(destinationAddress[0])<<8 | (uint16(destinationAddress[1]) & 0xFF)
							m.valueCacheMutex.RLock()
							val, ok := m.valueCache[addressData]
							m.valueCacheMutex.RUnlock()
							changed := false

							var payload []int8
							payload = append(payload, groupValueWrite.DataFirstByte)
							payload = append(payload, groupValueWrite.Data...)
							if !ok || !m.sliceEqual(val, payload) {
								m.valueCacheMutex.Lock()
								m.valueCache[addressData] = payload
								m.valueCacheMutex.Unlock()
								// If this is a new value, we have to also provide the 3 different types of addresses.
								if !ok {
									arb := utils.NewReadBuffer(utils.Int8ArrayToUint8Array(destinationAddress))
									if address, err2 := driverModel.KnxGroupAddressParse(arb, 3); err2 == nil {
										m.leve3AddressCache[addressData] = driverModel.CastKnxGroupAddress3Level(address)
									} else {
										fmt.Printf("Error parsing Group Address %s", err2.Error())
									}
									arb.Reset()
									if address, err2 := driverModel.KnxGroupAddressParse(arb, 2); err2 == nil {
										m.leve2AddressCache[addressData] = driverModel.CastKnxGroupAddress2Level(address)
									}
									arb.Reset()
									if address, err2 := driverModel.KnxGroupAddressParse(arb, 1); err2 == nil {
										m.leve1AddressCache[addressData] = driverModel.CastKnxGroupAddressFreeLevel(address)
									}
								}
								changed = true
							}
							if m.subscribers != nil {
								for _, subscriber := range m.subscribers {
									subscriber.handleValueChange(lDataInd.DataFrame, changed)
								}
							}
						}
					default:
						// If this is an individual address and it is targeted at us, we need to ack that.
						if !dataFrame.GroupAddress {
							targetAddress := Int8ArrayToKnxAddress(dataFrame.DestinationAddress)
							if *targetAddress == *m.ClientKnxAddress {
								log.Infof("Acknowleding an unhandled data message.")
								m.sendDeviceAck(*dataFrame.SourceAddress, dataFrame.Apdu.Counter)
							}
						}
					}
				case *driverModel.ApduControlContainer:
					// If this is an individual address and it is targeted at us, we need to ack that.
					if !dataFrame.GroupAddress {
						targetAddress := Int8ArrayToKnxAddress(dataFrame.DestinationAddress)
						if *targetAddress == *m.ClientKnxAddress {
							log.Infof("Acknowleding an unhandled contol message.")
							m.sendDeviceAck(*dataFrame.SourceAddress, dataFrame.Apdu.Counter)
						}
					}
				}
			default:
				log.Infof("Unknown unhandled message.")
			}
		}
	}()
}

func (m *KnxNetIpConnection) getGroupAddressNumLevels() uint8 {
	if val, ok := m.options["group-address-num-levels"]; ok {
		groupAddressNumLevels, err := strconv.Atoi(val[0])
		if err == nil {
			return uint8(groupAddressNumLevels)
		}
	}
	return 3
}

func (m *KnxNetIpConnection) addSubscriber(subscriber *KnxNetIpSubscriber) {
	for _, sub := range m.subscribers {
		if sub == subscriber {
			return
		}
	}
	m.subscribers = append(m.subscribers, subscriber)
}

func (m *KnxNetIpConnection) removeSubscriber(subscriber *KnxNetIpSubscriber) {
	for i, sub := range m.subscribers {
		if sub == subscriber {
			m.subscribers = append(m.subscribers[:i], m.subscribers[i+1:]...)
		}
	}
}

func (m *KnxNetIpConnection) sliceEqual(a, b []int8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func KnxAddressToString(knxAddress *driverModel.KnxAddress) string {
	return strconv.Itoa(int(knxAddress.MainGroup)) + "." + strconv.Itoa(int(knxAddress.MiddleGroup)) + "." + strconv.Itoa(int(knxAddress.SubGroup))
}
