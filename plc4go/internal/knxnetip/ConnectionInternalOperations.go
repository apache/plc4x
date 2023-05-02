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
	"context"
	"reflect"
	"time"

	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
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

func (m *Connection) sendGatewaySearchRequest(ctx context.Context) (driverModel.SearchResponse, error) {
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
	err = m.messageCodec.SendRequest(ctx, searchRequest,
		func(message spi.Message) bool {
			_, ok := message.(driverModel.SearchResponseExactly)
			return ok
		},
		func(message spi.Message) error {
			searchResponse := message.(driverModel.SearchResponse)
			result <- searchResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
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

func (m *Connection) sendGatewayConnectionRequest(ctx context.Context) (driverModel.ConnectionResponse, error) {
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
	err = m.messageCodec.SendRequest(ctx, connectionRequest,
		func(message spi.Message) bool {
			_, ok := message.(driverModel.ConnectionResponseExactly)
			return ok
		},
		func(message spi.Message) error {
			connectionResponse := message.(driverModel.ConnectionResponse)
			result <- connectionResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
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

func (m *Connection) sendGatewayDisconnectionRequest(ctx context.Context) (driverModel.DisconnectResponse, error) {
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
	err = m.messageCodec.SendRequest(ctx, disconnectRequest,
		func(message spi.Message) bool {
			_, ok := message.(driverModel.DisconnectResponseExactly)
			return ok
		},
		func(message spi.Message) error {
			disconnectResponse := message.(driverModel.DisconnectResponse)
			result <- disconnectResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl,
	)

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

func (m *Connection) sendConnectionStateRequest(ctx context.Context) (driverModel.ConnectionStateResponse, error) {
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
	err = m.messageCodec.SendRequest(ctx, connectionStateRequest,
		func(message spi.Message) bool {
			_, ok := message.(driverModel.ConnectionStateResponseExactly)
			return ok
		},
		func(message spi.Message) error {
			connectionStateResponse := message.(driverModel.ConnectionStateResponse)
			result <- connectionStateResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl,
	)

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

func (m *Connection) sendGroupAddressReadRequest(ctx context.Context, groupAddress []byte) (driverModel.ApduDataGroupValueResponse, error) {
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
	err := m.messageCodec.SendRequest(ctx, groupAddressReadRequest,
		func(message spi.Message) bool {
			tunnelingRequest, ok := message.(driverModel.TunnelingRequestExactly)
			if !ok || tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd, ok := tunnelingRequest.GetCemi().(driverModel.LDataIndExactly)
			if !ok {
				return false
			}
			dataFrameExt, ok := lDataInd.GetDataFrame().(driverModel.LDataExtendedExactly)
			if !ok {
				return false
			}
			dataContainer, ok := dataFrameExt.GetApdu().(driverModel.ApduDataContainerExactly)
			if !ok {
				return false
			}
			_, ok = dataContainer.GetDataApdu().(driverModel.ApduDataGroupValueResponseExactly)
			if !ok {
				return false
			}
			// Check if it's a value response for the given group address
			return dataFrameExt.GetGroupAddress() && reflect.DeepEqual(dataFrameExt.GetSourceAddress(), groupAddress)
		},
		func(message spi.Message) error {
			tunnelingRequest := message.(driverModel.TunnelingRequest)
			lDataInd := tunnelingRequest.GetCemi().(driverModel.LDataInd)
			dataFrameExt := lDataInd.GetDataFrame().(driverModel.LDataExtended)
			dataContainer := dataFrameExt.GetApdu().(driverModel.ApduDataContainer)
			groupReadResponse := dataContainer.GetDataApdu().(driverModel.ApduDataGroupValueResponse)

			result <- groupReadResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl,
	)

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

func (m *Connection) sendDeviceConnectionRequest(ctx context.Context, targetAddress driverModel.KnxAddress) (driverModel.ApduControlConnect, error) {
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
	err := m.messageCodec.SendRequest(ctx, deviceConnectionRequest,
		func(message spi.Message) bool {
			tunnelingRequest, ok := message.(driverModel.TunnelingRequestExactly)
			if !ok || tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataCon, ok := tunnelingRequest.GetCemi().(driverModel.LDataConExactly)
			if !ok {
				return false
			}
			lDataFrameExt, ok := lDataCon.GetDataFrame().(driverModel.LDataExtendedExactly)
			if !ok {
				return false
			}
			// Check if the address matches
			if ByteArrayToKnxAddress(lDataFrameExt.GetDestinationAddress()) != targetAddress {
				return false
			}
			apduControlContainer, ok := lDataFrameExt.GetApdu().(driverModel.ApduControlContainerExactly)
			if !ok {
				return false
			}
			_, ok = apduControlContainer.GetControlApdu().(driverModel.ApduControlConnectExactly)
			return ok
		},
		func(message spi.Message) error {
			tunnelingRequest := message.(driverModel.TunnelingRequest)
			lDataCon := tunnelingRequest.GetCemi().(driverModel.LDataCon)
			lDataFrameExt := lDataCon.GetDataFrame().(driverModel.LDataExtended)
			apduControlContainer := lDataFrameExt.GetApdu().(driverModel.ApduControlContainer)
			apduControlConnect := apduControlContainer.GetControlApdu().(driverModel.ApduControlConnect)

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
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl,
	)

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

func (m *Connection) sendDeviceDisconnectionRequest(ctx context.Context, targetAddress driverModel.KnxAddress) (driverModel.ApduControlDisconnect, error) {
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
	if err := m.messageCodec.SendRequest(ctx, deviceDisconnectionRequest,
		func(message spi.Message) bool {
			tunnelingRequest, ok := message.(driverModel.TunnelingRequestExactly)
			if !ok || tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataCon, ok := tunnelingRequest.GetCemi().(driverModel.LDataConExactly)
			if !ok {
				return false
			}
			dataFrameExt, ok := lDataCon.GetDataFrame().(driverModel.LDataExtendedExactly)
			if !ok {
				return false
			}
			curTargetAddress := ByteArrayToKnxAddress(dataFrameExt.GetDestinationAddress())
			// Check if the address matches
			if curTargetAddress != targetAddress {
				return false
			}
			apduControlContainer, ok := dataFrameExt.GetApdu().(driverModel.ApduControlContainerExactly)
			if !ok {
				return false
			}
			apduControlDisconnect := driverModel.ApduControlDisconnect(apduControlContainer.GetControlApdu())
			return apduControlDisconnect != nil
		},
		func(message spi.Message) error {
			tunnelingRequest := message.(driverModel.TunnelingRequest)
			lDataCon := tunnelingRequest.GetCemi().(driverModel.LDataCon)
			dataFrameExt := lDataCon.GetDataFrame().(driverModel.LDataExtended)
			apduControlContainer := dataFrameExt.GetApdu().(driverModel.ApduControlContainer)
			apduControlDisconnect := apduControlContainer.GetControlApdu().(driverModel.ApduControlDisconnect)

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
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl,
	); err != nil {
		return nil, errors.Wrap(err, "got error sending request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDeviceAuthentication(ctx context.Context, targetAddress driverModel.KnxAddress, authenticationLevel uint8, buildingKey []byte) (driverModel.ApduDataExtAuthorizeResponse, error) {
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
	if err := m.messageCodec.SendRequest(ctx, deviceAuthenticationRequest,
		func(message spi.Message) bool {
			tunnelingRequest, ok := message.(driverModel.TunnelingRequestExactly)
			if !ok || tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd := tunnelingRequest.GetCemi().(driverModel.LDataIndExactly)
			if !ok {
				return false
			}
			dataFrameExt, ok := lDataInd.GetDataFrame().(driverModel.LDataExtendedExactly)
			if !ok {
				return false
			}
			apduDataContainer, ok := dataFrameExt.GetApdu().(driverModel.ApduDataContainerExactly)
			if !ok {
				return false
			}
			apduDataOther, ok := apduDataContainer.GetDataApdu().(driverModel.ApduDataOtherExactly)
			if !ok {
				return false
			}
			_, ok = apduDataOther.GetExtendedApdu().(driverModel.ApduDataExtAuthorizeResponseExactly)
			if !ok {
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
		}, func(message spi.Message) error {
			tunnelingRequest := message.(driverModel.TunnelingRequest)
			lDataInd := tunnelingRequest.GetCemi().(driverModel.LDataInd)
			dataFrameExt := lDataInd.GetDataFrame().(driverModel.LDataExtended)
			apduDataContainer := dataFrameExt.GetApdu().(driverModel.ApduDataContainer)
			apduDataOther := apduDataContainer.GetDataApdu().(driverModel.ApduDataOther)
			apduAuthorizeResponse := apduDataOther.GetExtendedApdu().(driverModel.ApduDataExtAuthorizeResponse)

			// Acknowledge the receipt
			_ = m.sendDeviceAck(ctx, targetAddress, dataFrameExt.GetApdu().GetCounter(), func(err error) {
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
		}, func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl,
	); err != nil {
		return nil, errors.Wrap(err, "got error sending request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDeviceDeviceDescriptorReadRequest(ctx context.Context, targetAddress driverModel.KnxAddress) (driverModel.ApduDataDeviceDescriptorResponse, error) {
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
	err := m.messageCodec.SendRequest(ctx, deviceDescriptorReadRequest, func(message spi.Message) bool {
		tunnelingRequest, ok := message.(driverModel.TunnelingRequestExactly)
		if !ok || tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
			return false
		}
		lDataInd, ok := tunnelingRequest.GetCemi().(driverModel.LDataIndExactly)
		if !ok {
			return false
		}
		dataFrameExt, ok := lDataInd.GetDataFrame().(driverModel.LDataExtendedExactly)
		if !ok {
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
		dataContainer, ok := dataFrameExt.GetApdu().(driverModel.ApduDataContainerExactly)
		if !ok {
			return false
		}
		_, ok = dataContainer.GetDataApdu().(driverModel.ApduDataDeviceDescriptorResponseExactly)
		if !ok {
			return false
		}
		return true
	}, func(message spi.Message) error {
		tunnelingRequest := message.(driverModel.TunnelingRequest)
		lDataInd := tunnelingRequest.GetCemi().(driverModel.LDataInd)
		dataFrame := lDataInd.GetDataFrame().(driverModel.LDataExtended)
		dataContainer := dataFrame.GetApdu().(driverModel.ApduDataContainer)
		deviceDescriptorResponse := dataContainer.GetDataApdu().(driverModel.ApduDataDeviceDescriptorResponse)

		// Acknowledge the receipt
		_ = m.sendDeviceAck(ctx, targetAddress, dataFrame.GetApdu().GetCounter(), func(err error) {
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
	}, func(err error) error {
		// If this is a timeout, do a check if the connection requires a reconnection
		if _, isTimeout := err.(utils.TimeoutError); isTimeout {
			m.handleTimeout()
		}
		errorResult <- errors.Wrap(err, "got error processing request")
		return nil
	}, m.defaultTtl)

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

func (m *Connection) sendDevicePropertyReadRequest(ctx context.Context, targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8, propertyIndex uint16, numElements uint8) (driverModel.ApduDataExtPropertyValueResponse, error) {
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
	if err := m.messageCodec.SendRequest(ctx, propertyReadRequest,
		func(message spi.Message) bool {
			tunnelingRequest, ok := message.(driverModel.TunnelingRequestExactly)
			if !ok || tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd, ok := tunnelingRequest.GetCemi().(driverModel.LDataIndExactly)
			if !ok {
				return false
			}
			dataFrameExt, ok := lDataInd.GetDataFrame().(driverModel.LDataExtendedExactly)
			if !ok {
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
			dataContainer, ok := dataFrameExt.GetApdu().(driverModel.ApduDataContainerExactly)
			if !ok {
				return false
			}
			dataApduOther, ok := dataContainer.GetDataApdu().(driverModel.ApduDataOtherExactly)
			if !ok {
				return false
			}
			propertyValueResponse, ok := dataApduOther.GetExtendedApdu().(driverModel.ApduDataExtPropertyValueResponseExactly)
			if !ok {
				return false
			}
			return propertyValueResponse.GetObjectIndex() == objectId && propertyValueResponse.GetPropertyId() == propertyId
		},
		func(message spi.Message) error {
			tunnelingRequest := message.(driverModel.TunnelingRequest)
			lDataInd := tunnelingRequest.GetCemi().(driverModel.LDataInd)
			dataFrameExt := lDataInd.GetDataFrame().(driverModel.LDataExtended)
			dataContainer := dataFrameExt.GetApdu().(driverModel.ApduDataContainer)
			dataApduOther := dataContainer.GetDataApdu().(driverModel.ApduDataOther)
			propertyValueResponse := dataApduOther.GetExtendedApdu().(driverModel.ApduDataExtPropertyValueResponse)

			// Acknowledge the receipt
			_ = m.sendDeviceAck(ctx, targetAddress, dataFrameExt.GetApdu().GetCounter(), func(err error) {
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
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl,
	); err != nil {
		return nil, errors.Wrap(err, "got error sending device property read request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDevicePropertyDescriptionReadRequest(ctx context.Context, targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8) (driverModel.ApduDataExtPropertyDescriptionResponse, error) {
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
	err := m.messageCodec.SendRequest(ctx, propertyReadRequest,
		func(message spi.Message) bool {
			tunnelingRequest, ok := message.(driverModel.TunnelingRequestExactly)
			if !ok || tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd, ok := tunnelingRequest.GetCemi().(driverModel.LDataIndExactly)
			if !ok {
				return false
			}
			dataFrameExt, ok := lDataInd.GetDataFrame().(driverModel.LDataExtendedExactly)
			if !ok {
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
			dataContainer, ok := dataFrameExt.GetApdu().(driverModel.ApduDataContainerExactly)
			if !ok {
				return false
			}
			dataApduOther, ok := dataContainer.GetDataApdu().(driverModel.ApduDataOtherExactly)
			if !ok {
				return false
			}
			propertyDescriptionResponse, ok := dataApduOther.GetExtendedApdu().(driverModel.ApduDataExtPropertyDescriptionResponseExactly)
			if !ok {
				return false
			}
			return propertyDescriptionResponse.GetObjectIndex() == objectId && propertyDescriptionResponse.GetPropertyId() == propertyId
		},
		func(message spi.Message) error {
			tunnelingRequest := message.(driverModel.TunnelingRequest)
			lDataInd := tunnelingRequest.GetCemi().(driverModel.LDataInd)
			dataFrameExt := lDataInd.GetDataFrame().(driverModel.LDataExtended)
			dataContainer := dataFrameExt.GetApdu().(driverModel.ApduDataContainer)
			dataApduOther := dataContainer.GetDataApdu().(driverModel.ApduDataOther)
			propertyDescriptionResponse := dataApduOther.GetExtendedApdu().(driverModel.ApduDataExtPropertyDescriptionResponse)

			// Acknowledge the receipt
			_ = m.sendDeviceAck(ctx, targetAddress, dataFrameExt.GetApdu().GetCounter(), func(err error) {
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
		}, func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrapf(err, "got error processing request")
			return nil
		}, m.defaultTtl)

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

func (m *Connection) sendDeviceMemoryReadRequest(ctx context.Context, targetAddress driverModel.KnxAddress, address uint16, numBytes uint8) (driverModel.ApduDataMemoryResponse, error) {
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
	if err := m.messageCodec.SendRequest(ctx, propertyReadRequest,
		func(message spi.Message) bool {
			tunnelingRequest, ok := message.(driverModel.TunnelingRequestExactly)
			if !ok ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataInd, ok := tunnelingRequest.GetCemi().(driverModel.LDataIndExactly)
			if !ok {
				return false
			}
			dataFrameExt, ok := lDataInd.GetDataFrame().(driverModel.LDataExtendedExactly)
			if !ok {
				return false
			}
			dataContainer, ok := dataFrameExt.GetApdu().(driverModel.ApduDataContainerExactly)
			if !ok {
				return false
			}
			dataApduMemoryResponse, ok := dataContainer.GetDataApdu().(driverModel.ApduDataMemoryResponseExactly)
			if !ok {
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
			tunnelingRequest := message.(driverModel.TunnelingRequest)
			lDataInd := tunnelingRequest.GetCemi().(driverModel.LDataInd)
			dataFrameExt := lDataInd.GetDataFrame().(driverModel.LDataExtended)
			dataContainer := dataFrameExt.GetApdu().(driverModel.ApduDataContainer)
			dataApduMemoryResponse := dataContainer.GetDataApdu().(driverModel.ApduDataMemoryResponse)

			// Acknowledge the receipt
			_ = m.sendDeviceAck(ctx, targetAddress, dataFrameExt.GetApdu().GetCounter(), func(err error) {
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
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.Wrap(err, "got error processing request")
			return nil
		},
		m.defaultTtl,
	); err != nil {
		return nil, errors.Wrap(err, "got error sending memory read request")
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *Connection) sendDeviceAck(ctx context.Context, targetAddress driverModel.KnxAddress, counter uint8, callback func(err error)) error {
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

	if err := m.messageCodec.SendRequest(ctx, ack,
		func(message spi.Message) bool {
			tunnelingRequest, ok := message.(driverModel.TunnelingRequestExactly)
			if !ok ||
				tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
				return false
			}
			lDataCon, ok := tunnelingRequest.GetCemi().(driverModel.LDataConExactly)
			if !ok {
				return false
			}
			dataFrameExt, ok := lDataCon.GetDataFrame().(driverModel.LDataExtendedExactly)
			if !ok {
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
			controlContainer, ok := dataFrameExt.GetApdu().(driverModel.ApduControlContainer)
			if !ok {
				return false
			}
			_, ok = controlContainer.GetControlApdu().(driverModel.ApduControlAckExactly)
			return ok
		},
		func(message spi.Message) error {
			callback(nil)
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(utils.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			callback(errors.Wrap(err, "got error processing request"))
			return nil
		},
		m.defaultTtl,
	); err != nil {
		return errors.Wrap(err, "got error sending ack request")
	}

	return nil
}
