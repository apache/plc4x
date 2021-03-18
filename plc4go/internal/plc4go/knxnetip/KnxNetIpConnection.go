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
	"encoding/hex"
	"errors"
	"fmt"
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi"
	errors2 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/interceptors"
	internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/transports/udp"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	values2 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/plc4go"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"github.com/rs/zerolog/log"
	"math"
	"net"
	"reflect"
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

	valueCache      map[uint16][]int8
	valueCacheMutex sync.RWMutex
	metadata        *ConnectionMetadata
	defaultTtl      time.Duration
	connectionTtl   time.Duration
	buildingKey     []byte

	// Used for detecting connection problems
	connectionTimeoutTimer *time.Timer

	GatewayKnxAddress             *driverModel.KnxAddress
	ClientKnxAddress              *driverModel.KnxAddress
	CommunicationChannelId        uint8
	SequenceCounter               int32
	TunnelingRequestExpectationId int32
	DeviceConnections             map[driverModel.KnxAddress]*KnxDeviceConnection

	requestInterceptor internalModel.RequestInterceptor
	plc4go.PlcConnection
	sync.Mutex
}

type KnxReadResult struct {
	value    *values.PlcValue
	numItems uint8
	err      error
}

type KnxDeviceConnectResult struct {
	connection *KnxDeviceConnection
	err        error
}

type KnxDeviceDisconnectResult struct {
	connection *KnxDeviceConnection
	err        error
}

type KnxDeviceAuthenticateResult struct {
	err error
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
		valueCache:         map[uint16][]int8{},
		valueCacheMutex:    sync.RWMutex{},
		metadata:           &ConnectionMetadata{},
		defaultTtl:         time.Second * 10,
		DeviceConnections:  map[driverModel.KnxAddress]*KnxDeviceConnection{},
	}
	connection.connectionTtl = connection.defaultTtl * 2

	// If a building key was provided, save that in a dedicated variable
	if buildingKey, ok := options["buildingKey"]; ok {
		bc, err := hex.DecodeString(buildingKey[0])
		if err == nil {
			connection.buildingKey = bc
		}
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

		searchResponse, err := m.sendGatewaySearchRequest()
		if err != nil {
			sendResult(nil, errors.New("error discovering device capabilities"))
			return
		}

		// Save some important information
		m.metadata.KnxMedium = searchResponse.DibDeviceInfo.KnxMedium
		m.metadata.GatewayName = string(bytes.Trim(utils.Int8ArrayToByteArray(
			searchResponse.DibDeviceInfo.DeviceFriendlyName), "\x00"))
		m.GatewayKnxAddress = searchResponse.DibDeviceInfo.KnxAddress
		m.metadata.GatewayKnxAddress = KnxAddressToString(m.GatewayKnxAddress)
		m.metadata.ProjectNumber = searchResponse.DibDeviceInfo.ProjectInstallationIdentifier.ProjectNumber
		m.metadata.InstallationNumber = searchResponse.DibDeviceInfo.ProjectInstallationIdentifier.InstallationNumber
		m.metadata.DeviceSerialNumber = searchResponse.DibDeviceInfo.KnxNetIpDeviceSerialNumber
		m.metadata.DeviceMulticastAddress = searchResponse.DibDeviceInfo.KnxNetIpDeviceMulticastAddress.Addr
		m.metadata.DeviceMacAddress = searchResponse.DibDeviceInfo.KnxNetIpDeviceMacAddress.Addr
		m.metadata.SupportedServices = []string{}
		supportsTunneling := false
		for _, serviceId := range searchResponse.DibSuppSvcFamilies.ServiceIds {
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
			connectionResponse, err := m.sendGatewayConnectionRequest()
			if err != nil {
				sendResult(nil, errors.New("error connecting to device"))
				return
			}

			// Save the communication channel id
			m.CommunicationChannelId = connectionResponse.CommunicationChannelId

			// Reset the sequence counter
			m.SequenceCounter = -1

			// If the connection was successful, the gateway will now forward any packets
			// on the KNX bus that are broadcast packets to us, so we have to setup things
			// to handle these incoming messages.
			switch connectionResponse.Status {
			case driverModel.Status_NO_ERROR:
				// Save the KNX Address the Gateway assigned to us for this connection.
				tunnelConnectionDataBlock :=
					driverModel.CastConnectionResponseDataBlockTunnelConnection(
						connectionResponse.ConnectionResponseDataBlock)
				m.ClientKnxAddress = tunnelConnectionDataBlock.KnxAddress

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
								log.Warn().Msgf("Got an unhandled TunnelingResponse message %v\n", tunnelingResponse)
							} else {
								log.Warn().Msgf("Not a TunnelingRequest or TunnelingResponse message %v\n", incomingMessage)
							}
						} else {
							if tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
								log.Warn().Msgf("Not for this connection %v\n", tunnelingRequest)
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
		}
	}()

	return result
}

func (m *KnxNetIpConnection) BlockingClose() {
	closeResults := m.Close()
	select {
	case <-closeResults:
		return
	case <-time.After(m.defaultTtl):
		return
	}
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
			disconnects := m.DeviceDisconnect(targetAddress)
			select {
			case _ = <-disconnects:
			case <-time.After(m.defaultTtl):
				// If we got a timeout here, well just continue the device will just auto disconnect.
				log.Debug().Msgf("Timeout disconnecting from device %s.", KnxAddressToString(&targetAddress))
			}
		}

		// Send a disconnect request from the gateway.
		_, err := m.sendGatewayDisconnectionRequest()
		if err != nil {
			result <- plc4go.NewPlcConnectionCloseResult(m, errors.New(
				fmt.Sprintf("got an error when disconnecting: %s", err.Error())))
		} else {
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
		case <-time.After(m.defaultTtl):
			m.handleTimeout()
			return false
		}
	}
	return false
}

func (m *KnxNetIpConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
	result := make(chan plc4go.PlcConnectionPingResult)

	go func() {
		// Send the connection state request
		_, err := m.sendConnectionStateRequest()
		if err != nil {
			result <- plc4go.NewPlcConnectionPingResult(errors.New(
				fmt.Sprintf("got an error: %s", err.Error())))
		} else {
			result <- plc4go.NewPlcConnectionPingResult(nil)
		}
		return
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

///////////////////////////////////////////////////////////////////////////////////////////////////////
// KNX Specific Operations used by the driver internally
//
// These functions all provide access to some of the internal KNX operations
// They provide this functionality to other parts of the KNX driver, which are
// not part of the PLC4Go API.
//
// Remarks about these functions:
// They expect the called private functions to handle timeouts, so these will not.
///////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *KnxNetIpConnection) ReadGroupAddress(groupAddress []int8, datapointType *driverModel.KnxDatapointType) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	sendResponse := func(value *values.PlcValue, numItems uint8, err error) {
		select {
		case result <- KnxReadResult{
			value:    value,
			numItems: numItems,
			err:      err,
		}:
		default:
		}
	}

	go func() {
		groupAddressReadResponse, err := m.sendGroupAddressReadRequest(groupAddress)
		if err != nil {
			sendResponse(nil, 0, errors.New(
				"error reading group address: "+err.Error()))
			return
		}

		var payload []int8
		payload = append(payload, groupAddressReadResponse.DataFirstByte)
		payload = append(payload, groupAddressReadResponse.Data...)

		// Parse the response data.
		rb := utils.NewReadBuffer(utils.Int8ArrayToByteArray(payload))
		// If the size of the field is greater than 6, we have to skip the first byte
		if datapointType.DatapointMainType().SizeInBits() > 6 {
			_, _ = rb.ReadUint8(8)
		}
		// Set a default datatype if none is provided
		if *datapointType == driverModel.KnxDatapointType_DPT_UNKNOWN {
			defaultDatapointType := driverModel.KnxDatapointType_USINT
			datapointType = &defaultDatapointType
		}
		// Parse the value
		plcValue, err := driverModel.KnxDatapointParse(rb, *datapointType)
		if err != nil {
			sendResponse(nil, 0, errors.New("error parsing group address response: "+err.Error()))
			return
		}

		// Return the value
		sendResponse(&plcValue, 1, nil)
	}()

	return result
}

func (m *KnxNetIpConnection) DeviceConnect(targetAddress driverModel.KnxAddress) <-chan KnxDeviceConnectResult {
	result := make(chan KnxDeviceConnectResult)

	sendResponse := func(connection *KnxDeviceConnection, err error) {
		select {
		case result <- KnxDeviceConnectResult{
			connection: connection,
			err:        err,
		}:
		default:
		}
	}

	go func() {
		// If we're already connected, use that connection instead.
		if connection, ok := m.DeviceConnections[targetAddress]; ok {
			sendResponse(connection, nil)
			return
		}

		// First send a connection request
		controlConnectResponse, err := m.sendDeviceConnectionRequest(targetAddress)
		if err != nil {
			sendResponse(nil, errors.New(
				"error creating device connection: "+err.Error()))
			return
		}
		if controlConnectResponse == nil {
			sendResponse(nil, errors.New("error creating device connection"))
			return
		}

		// Create the new connection object.
		connection := &KnxDeviceConnection{
			counter: 0,
			// I was told this value on the knx-forum.
			// Seems the max payload is 3 bytes less ...
			maxApdu: 0, // This is the default max APDU Size
		}
		m.DeviceConnections[targetAddress] = connection

		// If the connection request was successful, try to read the device-descriptor
		deviceDescriptorResponse, err := m.sendDeviceDeviceDescriptorReadRequest(targetAddress)
		if err != nil {
			sendResponse(nil, errors.New(
				"error reading device descriptor: "+err.Error()))
			return
		}
		// Save the device-descriptor value
		deviceDescriptor := uint16(deviceDescriptorResponse.Data[0])<<8 | (uint16(deviceDescriptorResponse.Data[1]) & 0xFF)
		connection.deviceDescriptor = deviceDescriptor

		// Last, not least, read the max APDU size
		// If we were able to read the max APDU size, then use the minimum of
		// the connection APDU size and the device APDU size, otherwise use the
		// default APDU Size of 15
		// Defined in: 03_05_01 Resources v01.09.03 AS Page 40
		deviceApduSize := uint16(15)
		propertyValueResponse, err := m.sendDevicePropertyReadRequest(targetAddress, 0, 56, 1, 1)
		if err == nil {
			// If the count is 0, then this property doesn't exist or the user has no permission to read it.
			// In all other cases we expect the response to contain the value.
			if propertyValueResponse.Count > 0 {
				dataLength := uint8(len(propertyValueResponse.Data))
				data := propertyValueResponse.Data
				rb := utils.NewReadBuffer(data)
				plcValue, err := driverModel.KnxPropertyParse(rb,
					driverModel.KnxInterfaceObjectProperty_PID_DEVICE_MAX_APDULENGTH.PropertyDataType(), dataLength)

				// Return the result
				if err == nil {
					deviceApduSize = plcValue.GetUint16()
				}
			}
		}

		// Set the max apdu size for this connection.
		connection.maxApdu = uint16(math.Min(float64(deviceApduSize), 240))

		sendResponse(connection, nil)
	}()

	return result
}

func (m *KnxNetIpConnection) DeviceDisconnect(targetAddress driverModel.KnxAddress) <-chan KnxDeviceDisconnectResult {
	result := make(chan KnxDeviceDisconnectResult)

	sendResponse := func(connection *KnxDeviceConnection, err error) {
		select {
		case result <- KnxDeviceDisconnectResult{
			connection: connection,
			err:        err,
		}:
		default:
		}
	}

	go func() {
		if connection, ok := m.DeviceConnections[targetAddress]; ok {
			_, err := m.sendDeviceDisconnectionRequest(targetAddress)

			// Remove the connection from the list.
			delete(m.DeviceConnections, targetAddress)

			sendResponse(connection, err)
		} else {
			sendResponse(connection, nil)
		}
	}()

	return result
}

func (m *KnxNetIpConnection) DeviceAuthenticate(targetAddress driverModel.KnxAddress, buildingKey []byte) <-chan KnxDeviceAuthenticateResult {
	result := make(chan KnxDeviceAuthenticateResult)

	sendResponse := func(err error) {
		select {
		case result <- KnxDeviceAuthenticateResult{
			err: err,
		}:
		default:
		}
	}

	go func() {
		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.DeviceConnect(targetAddress)
			deviceConnectionResult := <-connections
			// If we didn't get a connect, abort
			if deviceConnectionResult.err != nil {
				sendResponse(errors.New(
					"error connecting to device at: " +
						KnxAddressToString(&targetAddress)))
			}
		}

		// If we successfully got a connection, read the property
		if connection != nil {
			authenticationLevel := uint8(0)
			authenticationResponse, err := m.sendDeviceAuthentication(targetAddress, authenticationLevel, buildingKey)
			if err == nil {
				if authenticationResponse.Level == authenticationLevel {
					sendResponse(nil)
				} else {
					// We authenticated correctly but not to the level requested.
					sendResponse(errors.New("got error authenticating at device " +
						KnxAddressToString(&targetAddress)))
				}
			} else {
				sendResponse(errors.New(
					"got error authenticating at device " +
						KnxAddressToString(&targetAddress)))
			}
		} else {
			sendResponse(errors.New("unable to connect to device"))
		}
	}()

	return result
}

func (m *KnxNetIpConnection) DeviceReadProperty(targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8, propertyIndex uint16, numElements uint8) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	sendResponse := func(value *values.PlcValue, numItems uint8, err error) {
		select {
		case result <- KnxReadResult{
			value:    value,
			numItems: numItems,
			err:      err,
		}:
		default:
		}
	}

	go func() {
		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.DeviceConnect(targetAddress)
			deviceConnectionResult := <-connections
			// If we didn't get a connect, abort
			if deviceConnectionResult.err != nil {
				sendResponse(nil, 0, errors.New(
					"error connecting to device at: "+
						KnxAddressToString(&targetAddress)))
			}
		}

		// If we successfully got a connection, read the property
		if connection != nil {
			propertyValueResponse, err := m.sendDevicePropertyReadRequest(targetAddress, objectId, propertyId, propertyIndex, numElements)
			if err != nil {
				sendResponse(nil, 0, err)
				return
			}

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
			if err != nil {
				sendResponse(nil, 0, err)
			} else {
				sendResponse(&plcValue, 1, err)
			}
		} else {
			sendResponse(nil, 0, errors.New("unable to connect to device"))
		}
	}()

	return result
}

func (m *KnxNetIpConnection) DeviceReadPropertyDescriptor(targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	sendResponse := func(value *values.PlcValue, numItems uint8, err error) {
		select {
		case result <- KnxReadResult{
			value:    value,
			numItems: numItems,
			err:      err,
		}:
		default:
		}
	}

	go func() {
		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.DeviceConnect(targetAddress)
			deviceConnectionResult := <-connections
			// If we didn't get a connect, abort
			if deviceConnectionResult.err != nil {
				sendResponse(nil, 0, errors.New(
					"error connecting to device at: "+
						KnxAddressToString(&targetAddress)))
			}
		}

		// If we successfully got a connection, read the property
		if connection != nil {
			propertyDescriptionResponse, err := m.sendDevicePropertyDescriptionReadRequest(targetAddress, objectId, propertyId)
			if err != nil {
				sendResponse(nil, 0, err)
				return
			}

			val := map[string]values.PlcValue{}
			val["writable"] = values2.NewPlcBOOL(propertyDescriptionResponse.WriteEnabled)
			val["dataType"] = values2.NewPlcSTRING(propertyDescriptionResponse.PropertyDataType.Name())
			val["maxElements"] = values2.NewPlcUINT(propertyDescriptionResponse.MaxNrOfElements)
			val["readLevel"] = values2.NewPlcSTRING(propertyDescriptionResponse.ReadLevel.String())
			val["writeLevel"] = values2.NewPlcSTRING(propertyDescriptionResponse.WriteLevel.String())
			str := values2.NewPlcStruct(val)
			sendResponse(&str, 1, nil)
		} else {
			sendResponse(nil, 0, errors.New("unable to connect to device"))
		}
	}()

	return result
}

func (m *KnxNetIpConnection) DeviceReadMemory(targetAddress driverModel.KnxAddress, address uint16, numElements uint8, datapointType *driverModel.KnxDatapointType) <-chan KnxReadResult {
	result := make(chan KnxReadResult)

	sendResponse := func(value *values.PlcValue, numItems uint8, err error) {
		select {
		case result <- KnxReadResult{
			value:    value,
			numItems: numItems,
			err:      err,
		}:
		default:
		}
	}

	go func() {
		// Set a default datatype, if none is specified
		if datapointType == nil {
			dpt := driverModel.KnxDatapointType_USINT
			datapointType = &dpt
		}

		// Check if there is already a connection available,
		// if not, create a new one.
		connection, ok := m.DeviceConnections[targetAddress]
		if !ok {
			connections := m.DeviceConnect(targetAddress)
			deviceConnectionResult := <-connections
			// If we didn't get a connect, abort
			if deviceConnectionResult.err != nil {
				sendResponse(nil, 0, errors.New(
					"error connecting to device at: "+
						KnxAddressToString(&targetAddress)))
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
			var results []values.PlcValue
			for remainingRequestElements > 0 {
				// As the maxApdu can change, we have to do this in the loop.
				maxNumBytes := uint8(math.Min(float64(connection.maxApdu-3), float64(63)))
				maxNumElementsPerRequest := uint8(math.Floor(float64(maxNumBytes / elementSize)))
				numElements := uint8(math.Min(float64(remainingRequestElements), float64(maxNumElementsPerRequest)))
				numBytes := numElements * uint8(math.Max(float64(1), float64(datapointType.DatapointMainType().SizeInBits()/8)))
				memoryReadResponse, err := m.sendDeviceMemoryReadRequest(targetAddress, curStartingAddress, numBytes)
				if err != nil {
					return
				}

				// If the number of bytes read is less than expected,
				// Update the connection.maxApdu value. This is required
				// as some devices seem to be sending back less than the
				// number of bytes specified than the maxApdu.
				if uint8(len(memoryReadResponse.Data)) < numBytes {
					connection.maxApdu = uint16(len(memoryReadResponse.Data) + 3)
				}

				// Parse the data according to the property type information
				rb := utils.NewReadBuffer(memoryReadResponse.Data)
				for rb.HasMore(datapointType.DatapointMainType().SizeInBits()) {
					plcValue, err := driverModel.KnxDatapointParse(rb, *datapointType)
					// Return the result
					if err != nil {
						sendResponse(nil, 0, err)
						return
					}
					results = append(results, plcValue)

					// Update the counters and addresses.
					remainingRequestElements--
					curStartingAddress = curStartingAddress + uint16(elementSize)
				}
				// If there are still remaining bytes, keep them for the next time.
			}
			if len(results) > 1 {
				var plcList values.PlcValue
				plcList = values2.NewPlcList(results)
				sendResponse(&plcList, 1, nil)
			} else if len(results) == 1 {
				sendResponse(&results[0], 1, nil)
			}
		}
	}()

	return result
}

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

func (m *KnxNetIpConnection) sendGatewaySearchRequest() (*driverModel.SearchResponse, error) {
	localAddress, err := m.getLocalAddress()
	if err != nil {
		return nil, errors.New("error getting local address: " + err.Error())
	}

	localAddr := driverModel.NewIPAddress(utils.ByteArrayToInt8Array(localAddress.IP))
	discoveryEndpoint := driverModel.NewHPAIDiscoveryEndpoint(
		driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port))
	searchRequest := driverModel.NewSearchRequest(discoveryEndpoint)

	result := make(chan *driverModel.SearchResponse)
	errorResult := make(chan error)
	err = m.messageCodec.SendRequest(searchRequest,
		func(message interface{}) bool {
			searchResponse := driverModel.CastSearchResponse(message)
			return searchResponse != nil
		},
		func(message interface{}) error {
			searchResponse := driverModel.CastSearchResponse(message)
			result <- searchResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing search request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending search request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendGatewayConnectionRequest() (*driverModel.ConnectionResponse, error) {
	localAddress, err := m.getLocalAddress()
	if err != nil {
		return nil, errors.New("error getting local address: " + err.Error())
	}

	localAddr := driverModel.NewIPAddress(utils.ByteArrayToInt8Array(localAddress.IP)[len(localAddress.IP)-4:])
	connectionRequest := driverModel.NewConnectionRequest(
		driverModel.NewHPAIDiscoveryEndpoint(driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port)),
		driverModel.NewHPAIDataEndpoint(driverModel.HostProtocolCode_IPV4_UDP, localAddr, uint16(localAddress.Port)),
		driverModel.NewConnectionRequestInformationTunnelConnection(driverModel.KnxLayer_TUNNEL_LINK_LAYER),
	)

	result := make(chan *driverModel.ConnectionResponse)
	errorResult := make(chan error)
	err = m.messageCodec.SendRequest(connectionRequest,
		func(message interface{}) bool {
			connectionResponse := driverModel.CastConnectionResponse(message)
			return connectionResponse != nil
		},
		func(message interface{}) error {
			connectionResponse := driverModel.CastConnectionResponse(message)
			result <- connectionResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendGatewayDisconnectionRequest() (*driverModel.DisconnectResponse, error) {
	localAddress, err := m.getLocalAddress()
	if err != nil {
		return nil, errors.New("error getting local address: " + err.Error())
	}

	localAddr := driverModel.NewIPAddress(utils.ByteArrayToInt8Array(localAddress.IP)[len(localAddress.IP)-4:])
	disconnectRequest := driverModel.NewDisconnectRequest(
		m.CommunicationChannelId,
		driverModel.NewHPAIControlEndpoint(
			driverModel.HostProtocolCode_IPV4_UDP,
			localAddr,
			uint16(localAddress.Port)))

	result := make(chan *driverModel.DisconnectResponse)
	errorResult := make(chan error)
	err = m.messageCodec.SendRequest(disconnectRequest,
		func(message interface{}) bool {
			disconnectResponse := driverModel.CastDisconnectResponse(message)
			return disconnectResponse != nil
		},
		func(message interface{}) error {
			disconnectResponse := driverModel.CastDisconnectResponse(message)
			result <- disconnectResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendConnectionStateRequest() (*driverModel.ConnectionStateResponse, error) {
	localAddress, err := m.getLocalAddress()
	if err != nil {
		return nil, errors.New("error getting local address: " + err.Error())
	}

	localAddr := driverModel.NewIPAddress(utils.ByteArrayToInt8Array(localAddress.IP)[len(localAddress.IP)-4:])
	connectionStateRequest := driverModel.NewConnectionStateRequest(
		m.CommunicationChannelId,
		driverModel.NewHPAIControlEndpoint(
			driverModel.HostProtocolCode_IPV4_UDP,
			localAddr, uint16(localAddress.Port)))

	result := make(chan *driverModel.ConnectionStateResponse)
	errorResult := make(chan error)
	err = m.messageCodec.SendRequest(connectionStateRequest,
		func(message interface{}) bool {
			connectionStateResponse := driverModel.CastConnectionStateResponse(message)
			return connectionStateResponse != nil
		},
		func(message interface{}) error {
			connectionStateResponse := driverModel.CastConnectionStateResponse(message)
			result <- connectionStateResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendGroupAddressReadRequest(groupAddress []int8) (*driverModel.ApduDataGroupValueResponse, error) {
	// Send the property read request and wait for a confirmation that this property is readable.
	groupAddressReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(true, 6, 0,
				m.ClientKnxAddress, groupAddress,
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataGroupValueRead(), false, 0),
				true, true, driverModel.CEMIPriority_LOW, false, false)))

	result := make(chan *driverModel.ApduDataGroupValueResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		groupAddressReadRequest,
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
			groupReadResponse := driverModel.CastApduDataGroupValueResponse(dataContainer.DataApdu)
			if groupReadResponse == nil {
				return false
			}
			// Check if it's a value response for the given group address
			return dataFrameExt.GroupAddress && reflect.DeepEqual(dataFrameExt.DestinationAddress, groupAddress)
		},
		func(message interface{}) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.DataFrame)
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.Apdu)
			groupReadResponse := driverModel.CastApduDataGroupValueResponse(dataContainer.DataApdu)

			result <- groupReadResponse
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendDeviceConnectionRequest(targetAddress driverModel.KnxAddress) (*driverModel.ApduControlConnect, error) {
	// Send a connection request to the individual KNX device
	deviceConnectionRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, uint8(0),
				driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduControlContainer(driverModel.NewApduControlConnect(), false, 0),
				true, true, driverModel.CEMIPriority_SYSTEM, false, false)))

	result := make(chan *driverModel.ApduControlConnect)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
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
			lDataFrameExt := driverModel.CastLDataExtended(lDataCon.DataFrame)
			apduControlContainer := driverModel.CastApduControlContainer(lDataFrameExt.Apdu)
			apduControlConnect := driverModel.CastApduControlConnect(apduControlContainer.ControlApdu)

			// If the error flag is set, there was an error connecting
			if lDataCon.DataFrame.ErrorFlag {
				errorResult <- errors.New("error connecting to device at: " + KnxAddressToString(&targetAddress))
			} else {
				result <- apduControlConnect
			}

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendDeviceDisconnectionRequest(targetAddress driverModel.KnxAddress) (*driverModel.ApduControlDisconnect, error) {
	// Send a connection request to the individual KNX device
	deviceDisconnectionRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, uint8(0),
				driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduControlContainer(driverModel.NewApduControlDisconnect(), false, 0),
				true, true, driverModel.CEMIPriority_SYSTEM, false, false)))

	result := make(chan *driverModel.ApduControlDisconnect)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
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
			dataFrameExt := driverModel.CastLDataExtended(lDataCon.DataFrame)
			apduControlContainer := driverModel.CastApduControlContainer(dataFrameExt.Apdu)
			apduControlDisconnect := driverModel.CastApduControlDisconnect(apduControlContainer.ControlApdu)

			// If the error flag is set, there was an error disconnecting
			if lDataCon.DataFrame.ErrorFlag {
				errorResult <- errors.New("error disconnecting from device at: " + KnxAddressToString(&targetAddress))
			} else {
				result <- apduControlDisconnect
			}

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendDeviceAuthentication(targetAddress driverModel.KnxAddress, authenticationLevel uint8, buildingKey []byte) (*driverModel.ApduDataExtAuthorizeResponse, error) {
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
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, uint8(0),
				driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduDataContainer(driverModel.NewApduDataOther(
					driverModel.NewApduDataExtAuthorizeRequest(authenticationLevel, utils.ByteArrayToUint8Array(buildingKey))),
					true, counter),
				true, true, driverModel.CEMIPriority_SYSTEM, false, false)))

	result := make(chan *driverModel.ApduDataExtAuthorizeResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
		deviceAuthenticationRequest,
		// The Gateway is now supposed to send an Ack to this request.
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
			apduDataContainer := driverModel.CastApduDataContainer(dataFrameExt.Apdu)
			if apduDataContainer == nil {
				return false
			}
			apduDataOther := driverModel.CastApduDataOther(apduDataContainer.DataApdu)
			if apduDataOther == nil {
				return false
			}
			apduAuthorizeResponse := driverModel.CastApduDataExtAuthorizeResponse(apduDataOther.ExtendedApdu)
			if apduAuthorizeResponse == nil {
				return false
			}
			curTargetAddress := Int8ArrayToKnxAddress(dataFrameExt.DestinationAddress)
			// Check if the addresses match
			if *curTargetAddress != *m.ClientKnxAddress {
				return false
			}
			if *dataFrameExt.SourceAddress != targetAddress {
				return false
			}
			// Check if the counter matches
			if dataFrameExt.Apdu.Counter != counter {
				return false
			}
			return true
		},
		func(message interface{}) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.DataFrame)
			apduDataContainer := driverModel.CastApduDataContainer(dataFrameExt.Apdu)
			apduDataOther := driverModel.CastApduDataOther(apduDataContainer.DataApdu)
			apduAuthorizeResponse := driverModel.CastApduDataExtAuthorizeResponse(apduDataOther.ExtendedApdu)

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrameExt.Apdu.Counter, func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.DataFrame.ErrorFlag {
					errorResult <- errors.New("error authenticating at device: " + KnxAddressToString(&targetAddress))
				} else if err != nil {
					errorResult <- errors.New("error sending ack to device: " + KnxAddressToString(&targetAddress))
				} else {
					result <- apduAuthorizeResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendDeviceDeviceDescriptorReadRequest(targetAddress driverModel.KnxAddress) (*driverModel.ApduDataDeviceDescriptorResponse, error) {
	// Next, read the device descriptor so we know how we have to communicate with the device.
	counter := m.getNextCounter(targetAddress)
	deviceDescriptorReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, uint8(0),
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataDeviceDescriptorRead(0), true, counter),
				true, true, driverModel.CEMIPriority_LOW, false, false)))

	result := make(chan *driverModel.ApduDataDeviceDescriptorResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
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

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrame.Apdu.Counter, func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.DataFrame.ErrorFlag {
					errorResult <- errors.New("error reading device descriptor from device: " + KnxAddressToString(&targetAddress))
				} else if err != nil {
					errorResult <- errors.New("error sending ack to device: " + KnxAddressToString(&targetAddress))
				} else {
					result <- deviceDescriptorResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending device descriptor read request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendDevicePropertyReadRequest(targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8, propertyIndex uint16, numElements uint8) (*driverModel.ApduDataExtPropertyValueResponse, error) {
	// Next, read the device descriptor so we know how we have to communicate with the device.
	// Send the property read request and wait for a confirmation that this property is readable.
	counter := m.getNextCounter(targetAddress)
	propertyReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, 0,
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataOther(
						driverModel.NewApduDataExtPropertyValueRead(objectId, propertyId, numElements, propertyIndex)),
					true, counter),
				true, true, driverModel.CEMIPriority_LOW, false, false)))

	result := make(chan *driverModel.ApduDataExtPropertyValueResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
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

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrameExt.Apdu.Counter, func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.DataFrame.ErrorFlag {
					errorResult <- errors.New("error reading property value from device: " + KnxAddressToString(&targetAddress))
				} else if err != nil {
					errorResult <- errors.New("error sending ack to device: " + KnxAddressToString(&targetAddress))
				} else {
					result <- propertyValueResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending device property read request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendDevicePropertyDescriptionReadRequest(targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8) (*driverModel.ApduDataExtPropertyDescriptionResponse, error) {
	// Next, read the device descriptor so we know how we have to communicate with the device.
	// Send the property read request and wait for a confirmation that this property is readable.
	counter := m.getNextCounter(targetAddress)
	propertyReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, 0,
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataOther(
						driverModel.NewApduDataExtPropertyDescriptionRead(objectId, propertyId, 1)),
					true, counter),
				true, true, driverModel.CEMIPriority_LOW, false, false)))

	result := make(chan *driverModel.ApduDataExtPropertyDescriptionResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
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
			propertyDescriptionResponse := driverModel.CastApduDataExtPropertyDescriptionResponse(dataApduOther.ExtendedApdu)
			if propertyDescriptionResponse == nil {
				return false
			}
			return propertyDescriptionResponse.ObjectIndex == objectId && propertyDescriptionResponse.PropertyId == propertyId
		},
		func(message interface{}) error {
			tunnelingRequest := driverModel.CastTunnelingRequest(message)
			lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
			dataFrameExt := driverModel.CastLDataExtended(lDataInd.DataFrame)
			dataContainer := driverModel.CastApduDataContainer(dataFrameExt.Apdu)
			dataApduOther := driverModel.CastApduDataOther(dataContainer.DataApdu)
			propertyDescriptionResponse := driverModel.CastApduDataExtPropertyDescriptionResponse(dataApduOther.ExtendedApdu)

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrameExt.Apdu.Counter, func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.DataFrame.ErrorFlag {
					errorResult <- errors.New("error reading property description from device: " + KnxAddressToString(&targetAddress))
				} else if err != nil {
					errorResult <- errors.New("error sending ack to device: " + KnxAddressToString(&targetAddress))
				} else {
					result <- propertyDescriptionResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending property description read request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendDeviceMemoryReadRequest(targetAddress driverModel.KnxAddress, address uint16, numBytes uint8) (*driverModel.ApduDataMemoryResponse, error) {
	// Next, read the device descriptor so we know how we have to communicate with the device.
	counter := m.getNextCounter(targetAddress)

	// Send the property read request and wait for a confirmation that this property is readable.
	propertyReadRequest := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, 0,
				driverModel.NewKnxAddress(0, 0, 0),
				KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduDataContainer(
					driverModel.NewApduDataMemoryRead(numBytes, address),
					true, counter),
				true, true, driverModel.CEMIPriority_LOW, false, false)))

	result := make(chan *driverModel.ApduDataMemoryResponse)
	errorResult := make(chan error)
	err := m.messageCodec.SendRequest(
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

			// Acknowledge the receipt
			_ = m.sendDeviceAck(targetAddress, dataFrameExt.Apdu.Counter, func(err error) {
				// If the error flag is set, there was an error authenticating
				if lDataInd.DataFrame.ErrorFlag {
					errorResult <- errors.New("error reading memory from device: " + KnxAddressToString(&targetAddress))
				} else if err != nil {
					errorResult <- errors.New("error sending ack to device: " + KnxAddressToString(&targetAddress))
				} else {
					result <- dataApduMemoryResponse
				}
			})

			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			errorResult <- errors.New(fmt.Sprintf("got error processing request: %s", err))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("got error sending memory read request: %s", err))
	}

	select {
	case response := <-result:
		return response, nil
	case errorResponse := <-errorResult:
		return nil, errorResponse
	}
}

func (m *KnxNetIpConnection) sendDeviceAck(targetAddress driverModel.KnxAddress, counter uint8, callback func(err error)) error {
	ack := driverModel.NewTunnelingRequest(
		driverModel.NewTunnelingRequestDataBlock(m.CommunicationChannelId, m.getNewSequenceCounter()),
		driverModel.NewLDataReq(0, nil,
			driverModel.NewLDataExtended(false, 6, uint8(0),
				driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
				driverModel.NewApduControlContainer(driverModel.NewApduControlAck(), true, counter),
				true, true, driverModel.CEMIPriority_SYSTEM, false, false)))

	err := m.messageCodec.SendRequest(
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
			callback(nil)
			return nil
		},
		func(err error) error {
			// If this is a timeout, do a check if the connection requires a reconnection
			if _, isTimeout := err.(errors2.TimeoutError); isTimeout {
				m.handleTimeout()
			}
			callback(errors.New(fmt.Sprintf("got error processing request: %s", err)))
			return nil
		},
		m.defaultTtl)

	if err != nil {
		return errors.New(fmt.Sprintf("got error sending ack request: %s", err))
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////
// Internal helper functions
///////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *KnxNetIpConnection) interceptIncomingMessage(interface{}) {
	m.resetTimeout()
	if m.connectionStateTimer != nil {
		// Reset the timer for sending the ConnectionStateRequest
		m.connectionStateTimer.Reset(60 * time.Second)
	}
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
							var payload []int8
							payload = append(payload, groupValueWrite.DataFirstByte)
							payload = append(payload, groupValueWrite.Data...)

							m.handleValueCacheUpdate(destinationAddress, payload)
						}
					default:
						// If this is an individual address and it is targeted at us, we need to ack that.
						if !dataFrame.GroupAddress {
							targetAddress := Int8ArrayToKnxAddress(dataFrame.DestinationAddress)
							if *targetAddress == *m.ClientKnxAddress {
								log.Info().Msg("Acknowleding an unhandled data message.")
								_ = m.sendDeviceAck(*dataFrame.SourceAddress, dataFrame.Apdu.Counter, func(err error) {})
							}
						}
					}
				case *driverModel.ApduControlContainer:
					// If this is an individual address and it is targeted at us, we need to ack that.
					if !dataFrame.GroupAddress {
						targetAddress := Int8ArrayToKnxAddress(dataFrame.DestinationAddress)
						if *targetAddress == *m.ClientKnxAddress {
							log.Info().Msg("Acknowleding an unhandled contol message.")
							_ = m.sendDeviceAck(*dataFrame.SourceAddress, dataFrame.Apdu.Counter, func(err error) {})
						}
					}
				}
			default:
				log.Info().Msg("Unknown unhandled message.")
			}
		}
	}()
}

func (m *KnxNetIpConnection) handleValueCacheUpdate(destinationAddress []int8, payload []int8) {
	addressData := uint16(destinationAddress[0])<<8 | (uint16(destinationAddress[1]) & 0xFF)

	m.valueCacheMutex.RLock()
	val, ok := m.valueCache[addressData]
	m.valueCacheMutex.RUnlock()
	changed := false
	if !ok || !m.sliceEqual(val, payload) {
		m.valueCacheMutex.Lock()
		m.valueCache[addressData] = payload
		m.valueCacheMutex.Unlock()
		changed = true
	}
	if m.subscribers != nil {
		for _, subscriber := range m.subscribers {
			subscriber.handleValueChange(destinationAddress, payload, changed)
		}
	}
}

func (m *KnxNetIpConnection) handleTimeout() {
	// If this is the first timeout in a sequence, start the timer.
	if m.connectionTimeoutTimer == nil {
		m.connectionTimeoutTimer = time.NewTimer(m.connectionTtl)
		go func() {
			<-m.connectionTimeoutTimer.C
			m.resetConnection()
		}()
	}
}

func (m *KnxNetIpConnection) resetTimeout() {
	if m.connectionTimeoutTimer != nil {
		m.connectionTimeoutTimer.Stop()
		m.connectionTimeoutTimer = nil
	}
}

func (m *KnxNetIpConnection) resetConnection() {
	fmt.Println("Bad connection detected")
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

func (m *KnxNetIpConnection) getNewSequenceCounter() uint8 {
	sequenceCounter := atomic.AddInt32(&m.SequenceCounter, 1)
	if sequenceCounter >= math.MaxUint8 {
		atomic.StoreInt32(&m.SequenceCounter, -1)
		sequenceCounter = -1
	}
	return uint8(sequenceCounter)
}

func (m *KnxNetIpConnection) getNextCounter(targetAddress driverModel.KnxAddress) uint8 {
	m.Lock()
	defer m.Unlock()

	connection, ok := m.DeviceConnections[targetAddress]
	if !ok {
		return 0
	}
	counter := connection.counter
	connection.counter++
	if connection.counter == 16 {
		connection.counter = 0
	}
	return counter
}

func KnxAddressToString(knxAddress *driverModel.KnxAddress) string {
	return strconv.Itoa(int(knxAddress.MainGroup)) + "." + strconv.Itoa(int(knxAddress.MiddleGroup)) + "." + strconv.Itoa(int(knxAddress.SubGroup))
}
