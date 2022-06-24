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
	"github.com/apache/plc4x/plc4go/internal/spi"
	"reflect"
	"time"

	"github.com/apache/plc4x/plc4go/internal/spi/plcerrors"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/pkg/errors"
)

///////////////////////////////////////////////////////////////////////////////////////////////////////
// KnxNetIpConnection internal operations
//
// These are used internally by functions of the KnxNetIpConnection.
//
// All of the sendXYZ functions take care of sending a request and waiting for the matching
// response. They don't actually process the data in the response, they just handle receiving
// it and returning it to the calling function.
//
// They all assume the connection is checked and is available.
///////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *Connection) sendGatewaySearchRequest() (driverModel.SearchResponse, error) {
	localAddress, err := m.getLocalAddress()
	if err != nil {
		return nil, errors.Wrap(err, "error getting local address")
	}

	localAddr := driverModel.NewIPAddress(localAddress.IP)
	discoveryEndpoint := driverModel.NewHPAIDiscoveryEndpoint(
		driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port),
	)
	searchRequest := driverModel.NewSearchRequest(discoveryEndpoint)

	result := make(chan driverModel.SearchResponse)
	errorResult := make(chan error)
	err = m.messageCodec.SendRequest(searchRequest,
		func(message spi.Message) bool {
			searchResponse := driverModel.CastSearchResponse(message)
			return searchResponse != nil
		},
		func(message spi.Message) error {
			searchResponse := driverModel.CastSearchResponse(message)
			result <- searchResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending search request")
	}

	ttlTimer := time.NewTimer(m.defaultTtl)
	select {
	case response := <-result:
		if !ttlTimer.Stop() {
			<-ttlTimer.C
		}
		return response, nil
	case errorResponse := <-errorResult:
		if !ttlTimer.Stop() {
			<-ttlTimer.C
		}
		return nil, errorResponse
		// For search requests there is no timeout handler running, so we have to do it manually.
	case <-ttlTimer.C:
		ttlTimer.Stop()
		return nil, errors.New("timeout")
	}
}

func (m *Connection) sendGatewayConnectionRequest() (driverModel.ConnectionResponse, error) {
	localAddress, err := m.getLocalAddress()
	if err != nil {
		return nil, errors.Wrap(err, "error getting local address")
	}

	localAddr := driverModel.NewIPAddress(localAddress.IP[len(localAddress.IP)-4:])
	connectionRequest := driverModel.NewConnectionRequest(
		driverModel.NewHPAIDiscoveryEndpoint(driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port)),
		driverModel.NewHPAIDataEndpoint(driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port)),
		driverModel.NewConnectionRequestInformationTunnelConnection(driverModel.KnxLayer_TUNNEL_LINK_LAYER),
	)

	result := make(chan driverModel.ConnectionResponse)
	errorResult := make(chan error)
	err = m.messageCodec.SendRequest(connectionRequest,
		func(message spi.Message) bool {
			connectionResponse := driverModel.CastConnectionResponse(message)
			return connectionResponse != nil
		},
		func(message spi.Message) error {
			connectionResponse := driverModel.CastConnectionResponse(message)
			result <- connectionResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendGatewayDisconnectionRequest() (driverModel.DisconnectResponse, error) {
	localAddress, err := m.getLocalAddress()
	if err != nil {
		return nil, errors.Wrap(err, "error getting local address")
	}

	localAddr := driverModel.NewIPAddress(localAddress.IP[len(localAddress.IP)-4:])
	disconnectRequest := driverModel.NewDisconnectRequest(
		m.CommunicationChannelId,
		driverModel.NewHPAIControlEndpoint(
			driverModel.HostProtocolCode_IPV4_UDP,
			localAddr,
			uint16(localAddress.Port),
		),
	)

	result := make(chan driverModel.DisconnectResponse)
	errorResult := make(chan error)
	err = m.messageCodec.SendRequest(disconnectRequest,
		func(message spi.Message) bool {
			disconnectResponse := driverModel.CastDisconnectResponse(message)
			return disconnectResponse != nil
		},
		func(message spi.Message) error {
			disconnectResponse := driverModel.CastDisconnectResponse(message)
			result <- disconnectResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendConnectionStateRequest() (driverModel.ConnectionStateResponse, error) {
	localAddress, err := m.getLocalAddress()
	if err != nil {
		return nil, errors.Wrap(err, "error getting local address")
	}

	localAddr := driverModel.NewIPAddress(localAddress.IP[len(localAddress.IP)-4:])
	connectionStateRequest := driverModel.NewConnectionStateRequest(
		m.CommunicationChannelId,
		driverModel.NewHPAIControlEndpoint(
			driverModel.HostProtocolCode_IPV4_UDP,
			localAddr, uint16(localAddress.Port)))

	result := make(chan driverModel.ConnectionStateResponse)
	errorResult := make(chan error)
	err = m.messageCodec.SendRequest(connectionStateRequest,
		func(message spi.Message) bool {
			connectionStateResponse := driverModel.CastConnectionStateResponse(message)
			return connectionStateResponse != nil
		},
		func(message spi.Message) error {
			connectionStateResponse := driverModel.CastConnectionStateResponse(message)
			result <- connectionStateResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendGroupAddressReadRequest(groupAddress []byte) (driverModel.ApduDataGroupValueResponse, error) {
	// Send the property read request and wait for a confirmation that this property is readable.
	groupAddressReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(
			0,
			nil,
			driverModel.NewLDataExtended(
				true,
				6,
				0,
				m.ClientKnxAddress, groupAddress,
				driverModel.NewApduDataContainer(driverModel.NewApduDataGroupValueRead(0), false, 0, 0),
				true,
				true,
				driverModel.CEMIPriority_LOW,
				false,
				false,
			),
			0,
		),
		0,
	)

	result := make(chan driverModel.ApduDataGroupValueResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		groupAddressReadRequest,
		func(message spi.Message) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			if lDataInd == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			if dataFrameExt == nil {
				return false
			}
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			if dataContainer == nil {
				return false
			}
			groupReadResponse := driverModel.CastApduDataGroupValueResponse(dataContainer.GetDataApdu())
			if groupReadResponse == nil {
				return false
			}
			// Check if it's a value response for the given group address
			return dataFrameExt.GetGroupAddress() && reflect.DeepEqual(dataFrameExt.GetSourceAddress(), groupAddress)
		},
		func(message spi.Message) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			groupReadResponse := driverModel.CastApduDataGroupValueResponse(dataContainer.GetDataApdu())

			result <- groupReadResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDeviceConnectionRequest(targetAddress driverModel.KnxAddress) (driverModel.ApduControlConnect, error) {
	// Send a connection request to the individual KNX device
	deviceConnectionRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(
			0,
			nil,
			driverModel.NewLDataExtended(
				false,
				6,
				uint8(0),
				driverModel.NewKnxAddress(0, 0, 0), KnxAddressToByteArray(targetAddress),
				driverModel.NewApduControlContainer(driverModel.NewApduControlConnect(), false, 0, 0),
				true,
				true,
				driverModel.CEMIPriority_SYSTEM,
				false,
				false,
			),
			0,
		),
		0,
	)

	result := make(chan driverModel.ApduControlConnect)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		deviceConnectionRequest,
		// The Gateway is now supposed to send an Ack to this request.
		func(message spi.Message) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataCon := driverModel.CastLDataCon(tunnelingRequest.GetCemi())
			if lDataCon == nil {
				return false
			}
			lDataFrameExt := driverModel.CastLDataExtended(lDataCon.GetDataFrame())
			if lDataFrameExt == nil {
				return false
			}
			// Check if the address matches
			if ByteArrayToKnxAddress(lDataFrameExt.GetDestinationAddress()) != targetAddress {
				return false
			}
			apduControlContainer := driverModel.CastApduControlContainer(lDataFrameExt.GetApdu())
			if apduControlContainer == nil {
				return false
			}
			apduControlConnect := driverModel.CastApduControlConnect(apduControlContainer.GetControlApdu())
			return apduControlConnect != nil
		},
		func(message spi.Message) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataCon := driverModel.CastLDataCon(tunnelingRequest.GetCemi())
			lDataFrameExt := driverModel.CastLDataExtended(lDataCon.GetDataFrame())
			apduControlContainer := driverModel.CastApduControlContainer(lDataFrameExt.GetApdu())
			apduControlConnect := driverModel.CastApduControlConnect(apduControlContainer.GetControlApdu())

			// If the error flag is set, there was an error connecting
			if lDataCon.GetDataFrame().GetErrorFlag() {
				errorResult <- errors.Errorf("error connecting to device at: %s", KnxAddressToString(targetAddress))
			} else {
				result <- apduControlConnect
			}

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDeviceDisconnectionRequest(targetAddress driverModel.KnxAddress) (driverModel.ApduControlDisconnect, error) {
	// Send a connection request to the individual KNX device
	deviceDisconnectionRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(
			0,
			nil,
			driverModel.NewLDataExtended(
				false,
				6,
				uint8(0),
				driverModel.NewKnxAddress(0, 0, 0), KnxAddressToByteArray(targetAddress),
				driverModel.NewApduControlContainer(driverModel.NewApduControlDisconnect(), false, 0, 0),
				true,
				true,
				driverModel.CEMIPriority_SYSTEM,
				false,
				false,
			),
			0,
		),
		0,
	)

	result := make(chan driverModel.ApduControlDisconnect)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		deviceDisconnectionRequest,
		// The Gateway is now supposed to send an Ack to this request.
		func(message spi.Message) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataCon := driverModel.CastLDataCon(tunnelingRequest.GetCemi())
			if lDataCon == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataCon.GetDataFrame())
			if dataFrameExt == nil {
				return false
			}
			curTargetAddress := ByteArrayToKnxAddress(dataFrameExt.GetDestinationAddress())
			// Check if the address matches
			if curTargetAddress != targetAddress {
				return false
			}
			apduControlContainer := driverModel.CastApduControlContainer(dataFrameExt.GetApdu())
			if apduControlContainer == nil {
				return false
			}
			apduControlDisconnect := driverModel.CastApduControlDisconnect(apduControlContainer.GetControlApdu())
			return apduControlDisconnect != nil
		},
		func(message spi.Message) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataCon := driverModel.CastLDataCon(tunnelingRequest.GetCemi())
			dataFrameExt := driverModel.CastLDataExtended(lDataCon.GetDataFrame())
			apduControlContainer := driverModel.CastApduControlContainer(dataFrameExt.GetApdu())
			apduControlDisconnect := driverModel.CastApduControlDisconnect(apduControlContainer.GetControlApdu())

			// If the error flag is set, there was an error disconnecting
			if lDataCon.GetDataFrame().GetErrorFlag() {
				errorResult <- errors.Errorf("error disconnecting from device at: %s", KnxAddressToString(targetAddress))
			} else {
				result <- apduControlDisconnect
			}

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDeviceAuthentication(targetAddress driverModel.KnxAddress, authenticationLevel uint8, buildingKey []byte) (driverModel.ApduDataExtAuthorizeResponse, error) {
	// Check if there is already a connection available,
	// if not, create a new one.
	connection, ok := m.DeviceConnections[targetAddress]
	if !ok {
		return nil, errors.New("not connected")
	}

	// Send a connection request to the individual KNX device
	counter := connection.counter
	connection.counter++
	deviceAuthenticationRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(
			0,
			nil,
			driverModel.NewLDataExtended(
				false,
				6,
				uint8(0),
				driverModel.NewKnxAddress(0, 0, 0), KnxAddressToByteArray(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataOther(
						driverModel.NewApduDataExtAuthorizeRequest(authenticationLevel, utils.ByteArrayToUint8Array(buildingKey), 0),
						0,
					),
					true,
					counter,
					0,
				),
				true,
				true,
				driverModel.CEMIPriority_SYSTEM,
				false,
				false,
			),
			0,
		),
		0,
	)

	result := make(chan driverModel.ApduDataExtAuthorizeResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		deviceAuthenticationRequest,
		// The Gateway is now supposed to send an Ack to this request.
		func(message spi.Message) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			if lDataInd == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			if dataFrameExt == nil {
				return false
			}
			apduDataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			if apduDataContainer == nil {
				return false
			}
			apduDataOther := driverModel.CastApduDataOther(apduDataContainer.GetDataApdu())
			if apduDataOther == nil {
				return false
			}
			apduAuthorizeResponse := driverModel.CastApduDataExtAuthorizeResponse(apduDataOther.GetExtendedApdu())
			if apduAuthorizeResponse == nil {
				return false
			}
			curTargetAddress := ByteArrayToKnxAddress(dataFrameExt.GetDestinationAddress())
			// Check if the addresses match
			if curTargetAddress != m.ClientKnxAddress {
				return false
			}
			if dataFrameExt.GetSourceAddress() != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.GetApdu().GetCounter() != counter {
				return false
			}
			return true
		},
		func(message spi.Message) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			apduDataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			apduDataOther := driverModel.CastApduDataOther(apduDataContainer.GetDataApdu())
			apduAuthorizeResponse := driverModel.CastApduDataExtAuthorizeResponse(apduDataOther.GetExtendedApdu())

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrameExt.GetApdu().GetCounter(), func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.GetDataFrame().GetErrorFlag() {
					errorResult <- errors.New("error authenticating at device: " + KnxAddressToString(targetAddress))
				} else if err != nil {
					errorResult <- errors.Wrapf(err, "error sending ack to device: %s", KnxAddressToString(targetAddress))
				} else {
					result <- apduAuthorizeResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDeviceDeviceDescriptorReadRequest(targetAddress driverModel.KnxAddress) (driverModel.ApduDataDeviceDescriptorResponse, error) {
	// Next, read the device descriptor so we know how we have to communicate with the device.
	counter := m.getNextCounter(targetAddress)
	deviceDescriptorReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(
			0,
			nil,
			driverModel.NewLDataExtended(
				false,
				6,
				uint8(0),
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToByteArray(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataDeviceDescriptorRead(0, 0), true, counter, 0,
				),
				true,
				true,
				driverModel.CEMIPriority_LOW,
				false,
				false,
			),
			0,
		),
		0,
	)

	result := make(chan driverModel.ApduDataDeviceDescriptorResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		deviceDescriptorReadRequest,
		func(message spi.Message) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			if lDataInd == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			if dataFrameExt == nil {
				return false
			}
			// Check if the address matches
			if dataFrameExt.GetSourceAddress() != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.GetApdu().GetCounter() != counter {
				return false
			}
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			if dataContainer == nil {
				return false
			}
			deviceDescriptorResponse := driverModel.CastApduDataDeviceDescriptorResponse(dataContainer.GetDataApdu())
			if deviceDescriptorResponse == nil {
				return false
			}
			return true
		},
		func(message spi.Message) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			dataFrame := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			dataContainer := driverModel.CastApduDataContainer(dataFrame.GetApdu())
			deviceDescriptorResponse := driverModel.CastApduDataDeviceDescriptorResponse(dataContainer.GetDataApdu())

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrame.GetApdu().GetCounter(), func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.GetDataFrame().GetErrorFlag() {
					errorResult <- errors.New("error reading device descriptor from device: " + KnxAddressToString(targetAddress))
				} else if err != nil {
					errorResult <- errors.Wrapf(err, "error sending ack to device: %s", KnxAddressToString(targetAddress))
				} else {
					result <- deviceDescriptorResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending device descriptor read request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDevicePropertyReadRequest(targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8, propertyIndex uint16, numElements uint8) (driverModel.ApduDataExtPropertyValueResponse, error) {
	// Next, read the device descriptor so we know how we have to communicate with the device.
	// Send the property read request and wait for a confirmation that this property is readable.
	counter := m.getNextCounter(targetAddress)
	propertyReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(
			0,
			nil,
			driverModel.NewLDataExtended(
				false,
				6,
				0,
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToByteArray(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataOther(
						driverModel.NewApduDataExtPropertyValueRead(objectId, propertyId, numElements, propertyIndex, 0),
						0,
					),
					true,
					counter,
					0,
				),
				true,
				true,
				driverModel.CEMIPriority_LOW,
				false,
				false,
			),
			0,
		),
		0,
	)

	result := make(chan driverModel.ApduDataExtPropertyValueResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		propertyReadRequest,
		func(message spi.Message) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			if lDataInd == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			if dataFrameExt == nil {
				return false
			}
			// Check if the address matches
			if dataFrameExt.GetSourceAddress() != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.GetApdu().GetCounter() != counter {
				return false
			}
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			if dataContainer == nil {
				return false
			}
			dataApduOther := driverModel.CastApduDataOther(dataContainer.GetDataApdu())
			if dataApduOther == nil {
				return false
			}
			propertyValueResponse := driverModel.CastApduDataExtPropertyValueResponse(dataApduOther.GetExtendedApdu())
			if propertyValueResponse == nil {
				return false
			}
			return propertyValueResponse.GetObjectIndex() == objectId && propertyValueResponse.GetPropertyId() == propertyId
		},
		func(message spi.Message) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			dataApduOther := driverModel.CastApduDataOther(dataContainer.GetDataApdu())
			propertyValueResponse := driverModel.CastApduDataExtPropertyValueResponse(dataApduOther.GetExtendedApdu())

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrameExt.GetApdu().GetCounter(), func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.GetDataFrame().GetErrorFlag() {
					errorResult <- errors.New("error reading property value from device: " + KnxAddressToString(targetAddress))
				} else if err != nil {
					errorResult <- errors.Wrapf(err, "error sending ack to device: %s", KnxAddressToString(targetAddress))
				} else {
					result <- propertyValueResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending device property read request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDevicePropertyDescriptionReadRequest(targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8) (driverModel.ApduDataExtPropertyDescriptionResponse, error) {
	// Next, read the device descriptor so we know how we have to communicate with the device.
	// Send the property read request and wait for a confirmation that this property is readable.
	counter := m.getNextCounter(targetAddress)
	propertyReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(
			0,
			nil,
			driverModel.NewLDataExtended(
				false,
				6,
				0,
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToByteArray(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataOther(
						driverModel.NewApduDataExtPropertyDescriptionRead(objectId, propertyId, 1, 0),
						0,
					),
					true,
					counter,
					0,
				),
				true,
				true,
				driverModel.CEMIPriority_LOW,
				false,
				false,
			),
			0,
		),
		0,
	)

	result := make(chan driverModel.ApduDataExtPropertyDescriptionResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		propertyReadRequest,
		func(message spi.Message) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			if lDataInd == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			if dataFrameExt == nil {
				return false
			}
			// Check if the address matches
			if dataFrameExt.GetSourceAddress() != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.GetApdu().GetCounter() != counter {
				return false
			}
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			if dataContainer == nil {
				return false
			}
			dataApduOther := driverModel.CastApduDataOther(dataContainer.GetDataApdu())
			if dataApduOther == nil {
				return false
			}
			propertyDescriptionResponse := driverModel.CastApduDataExtPropertyDescriptionResponse(dataApduOther.GetExtendedApdu())
			if propertyDescriptionResponse == nil {
				return false
			}
			return propertyDescriptionResponse.GetObjectIndex() == objectId && propertyDescriptionResponse.GetPropertyId() == propertyId
		},
		func(message spi.Message) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			dataApduOther := driverModel.CastApduDataOther(dataContainer.GetDataApdu())
			propertyDescriptionResponse := driverModel.CastApduDataExtPropertyDescriptionResponse(dataApduOther.GetExtendedApdu())

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrameExt.GetApdu().GetCounter(), func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.GetDataFrame().GetErrorFlag() {
					errorResult <- errors.Errorf("error reading property description from device: %s", KnxAddressToString(targetAddress))
				} else if err != nil {
					errorResult <- errors.Wrapf(err, "error sending ack to device: %s", KnxAddressToString(targetAddress))
				} else {
					result <- propertyDescriptionResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrapf(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending property description read request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDeviceMemoryReadRequest(targetAddress driverModel.KnxAddress, address uint16, numBytes uint8) (driverModel.ApduDataMemoryResponse, error) {
	// Next, read the device descriptor so we know how we have to communicate with the device.
	counter := m.getNextCounter(targetAddress)

	// Send the property read request and wait for a confirmation that this property is readable.
	propertyReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(
			0,
			nil,
			driverModel.NewLDataExtended(
				false,
				6,
				0,
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToByteArray(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataMemoryRead(numBytes, address, 0),
					true,
					counter,
					0,
				),
				true,
				true,
				driverModel.CEMIPriority_LOW,
				false,
				false,
			),
			0,
		),
		0,
	)

	result := make(chan driverModel.ApduDataMemoryResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		propertyReadRequest,
		func(message spi.Message) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			if lDataInd == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			if dataFrameExt == nil {
				return false
			}
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			if dataContainer == nil {
				return false
			}
			dataApduMemoryResponse := driverModel.CastApduDataMemoryResponse(dataContainer.GetDataApdu())
			if dataApduMemoryResponse == nil {
				return false
			}

			// Check if the address matches
			if dataFrameExt.GetSourceAddress() != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.GetApdu().GetCounter() != counter {
				return false
			}
			return dataApduMemoryResponse.GetAddress() == address
		},
		func(message spi.Message) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.GetCemi())
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.GetDataFrame())
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.GetApdu())
			dataApduMemoryResponse := driverModel.CastApduDataMemoryResponse(dataContainer.GetDataApdu())

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrameExt.GetApdu().GetCounter(), func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.GetDataFrame().GetErrorFlag() {
					errorResult <- errors.Errorf("error reading memory from device: %s", KnxAddressToString(targetAddress))
				} else if err != nil {
					errorResult <- errors.Errorf("error sending ack to device: %s", KnxAddressToString(targetAddress))
				} else {
					result <- dataApduMemoryResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.Wrap(err, "got error sending memory read request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDeviceAck(targetAddress driverModel.KnxAddress, counter uint8, callback func(err error)) error {
	ack := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(
			0,
			nil,
			driverModel.NewLDataExtended(
				false,
				6,
				uint8(0),
				driverModel.NewKnxAddress(0, 0, 0), KnxAddressToByteArray(targetAddress),
				driverModel.NewApduControlContainer(driverModel.NewApduControlAck(), true, counter, 0),
				true,
				true,
				driverModel.CEMIPriority_SYSTEM,
				false,
				false,
			),
			0,
		),
		0,
	)

	err := m.messageCodec.SendRequest(
		ack,
		func(message spi.Message) bool {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			if tunnelingRequest == nil ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataCon := driverModel.CastLDataCon(tunnelingRequest.GetCemi())
			if lDataCon == nil {
				return false
			}
			dataFrameExt := driverModel.CastLDataExtended(lDataCon.GetDataFrame())
			if dataFrameExt == nil {
				return false
			}
			// Check if the addresses match
			if dataFrameExt.GetSourceAddress() != m.ClientKnxAddress {
				return false
			}
			curTargetAddress := ByteArrayToKnxAddress(dataFrameExt.GetDestinationAddress())
			if curTargetAddress != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.GetApdu().GetCounter() != counter {
				return false
			}
			controlContainer := driverModel.CastApduControlContainer(dataFrameExt.GetApdu())
			if controlContainer == nil {
				return false
			}
			dataApduAck := driverModel.CastApduControlAck(controlContainer.GetControlApdu())
			if dataApduAck == nil {
				return false
			}
			return true
		},
		func(message spi.Message) error {
			callback(nil)
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(plcerrors.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			callback(errors.Wrap(err, "got error processing request"))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return errors.Wrap(err, "got error sending ack request")
	}

	return nil
}
