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
	"bytes"
	"encoding/hex"
	"fmt"
	_default "github.com/apache/plc4x/plc4go/internal/spi/default"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/internal/spi/interceptors"
	internalModel "github.com/apache/plc4x/plc4go/internal/spi/model"
	"github.com/apache/plc4x/plc4go/internal/spi/transports"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ConnectionMetadata struct {
	KnxMedium         driverModel.KnxMedium
	GatewayName       string
	GatewayKnxAddress string
	ClientKnxAddress  string

	ProjectNumber          uint8
	InstallationNumber     uint8
	DeviceSerialNumber     []byte
	DeviceMulticastAddress []byte
	DeviceMacAddress       []byte
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
		"DeviceSerialNumber":     utils.ByteArrayToString(m.DeviceSerialNumber, " "),
		"DeviceMulticastAddress": utils.ByteArrayToString(m.DeviceSerialNumber, "."),
		"DeviceMacAddress":       utils.ByteArrayToString(m.DeviceSerialNumber, ":"),
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

type Connection struct {
	messageCodec             spi.MessageCodec
	options                  map[string][]string
	fieldHandler             spi.PlcFieldHandler
	valueHandler             spi.PlcValueHandler
	connectionStateTimer     *time.Ticker
	quitConnectionStateTimer chan struct{}
	subscribers              []*Subscriber

	valueCache      map[uint16][]byte
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

	requestInterceptor interceptors.RequestInterceptor
	sync.Mutex

	// indicates if the tunneling requests loop is running
	handleTunnelingRequests bool

	connectionId string
	tracer       *spi.Tracer
}

func (m *Connection) String() string {
	return fmt.Sprintf("knx.Connection{}")
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

func NewConnection(transportInstance transports.TransportInstance, options map[string][]string, fieldHandler spi.PlcFieldHandler) *Connection {
	connection := &Connection{
		options:      options,
		fieldHandler: fieldHandler,
		valueHandler: NewValueHandler(),
		requestInterceptor: interceptors.NewSingleItemRequestInterceptor(
			internalModel.NewDefaultPlcReadRequest,
			internalModel.NewDefaultPlcWriteRequest,
			internalModel.NewDefaultPlcReadResponse,
			internalModel.NewDefaultPlcWriteResponse,
		),
		subscribers:             []*Subscriber{},
		valueCache:              map[uint16][]byte{},
		valueCacheMutex:         sync.RWMutex{},
		metadata:                &ConnectionMetadata{},
		defaultTtl:              time.Second * 10,
		DeviceConnections:       map[driverModel.KnxAddress]*KnxDeviceConnection{},
		handleTunnelingRequests: true,
	}
	connection.connectionTtl = connection.defaultTtl * 2

	if traceEnabledOption, ok := options["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = spi.NewTracer(connection.connectionId)
		}
	}
	// If a building key was provided, save that in a dedicated variable
	if buildingKey, ok := options["buildingKey"]; ok {
		bc, err := hex.DecodeString(buildingKey[0])
		if err == nil {
			connection.buildingKey = bc
		}
	}
	connection.messageCodec = NewMessageCodec(transportInstance, connection.interceptIncomingMessage)
	return connection
}

func (m *Connection) GetConnectionId() string {
	return m.connectionId
}

func (m *Connection) IsTraceEnabled() bool {
	return m.tracer != nil
}

func (m *Connection) GetTracer() *spi.Tracer {
	return m.tracer
}

func (m *Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	result := make(chan plc4go.PlcConnectionConnectResult)
	sendResult := func(connection plc4go.PlcConnection, err error) {
		result <- _default.NewDefaultPlcConnectionConnectResult(connection, err)
	}

	go func() {
		// Open the UDP Connection
		err := m.messageCodec.Connect()
		if err != nil {
			m.doSomethingAndClose(func() { sendResult(nil, errors.Wrap(err, "error opening connection")) })
			return
		}

		// Send a search request before connecting to the device.
		searchResponse, err := m.sendGatewaySearchRequest()
		if err != nil {
			m.doSomethingAndClose(func() { sendResult(nil, errors.Wrap(err, "error discovering device capabilities")) })
			return
		}

		// Save some important information
		m.metadata.KnxMedium = searchResponse.DibDeviceInfo.KnxMedium
		m.metadata.GatewayName = string(bytes.Trim(searchResponse.DibDeviceInfo.DeviceFriendlyName, "\x00"))
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
				m.doSomethingAndClose(func() { sendResult(nil, errors.Wrap(err, "error connecting to device")) })
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
				tunnelConnectionDataBlock := driverModel.CastConnectionResponseDataBlockTunnelConnection(
					connectionResponse.ConnectionResponseDataBlock,
				)
				m.ClientKnxAddress = tunnelConnectionDataBlock.KnxAddress

				// Create a go routine to handle incoming tunneling-requests which haven't been
				// handled by any other handler. This is where usually the GroupValueWrite messages
				// are being handled.
				log.Debug().Msg("Starting tunneling handler")
				go func() {
					defaultIncomingMessageChannel := m.messageCodec.GetDefaultIncomingMessageChannel()
					for m.handleTunnelingRequests {
						incomingMessage := <-defaultIncomingMessageChannel
						tunnelingRequest := driverModel.CastTunnelingRequest(incomingMessage)
						if tunnelingRequest == nil {
							tunnelingResponse := driverModel.CastTunnelingResponse(incomingMessage)
							if tunnelingResponse != nil {
								log.Warn().Msgf("Got an unhandled TunnelingResponse message %v\n", tunnelingResponse)
							} else {
								log.Warn().Msgf("Not a TunnelingRequest or TunnelingResponse message %v\n", incomingMessage)
							}
							continue
						}

						if tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
							log.Warn().Msgf("Not for this connection %v\n", tunnelingRequest)
							continue
						}

						lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
						if lDataInd == nil {
							continue
						}
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
					log.Warn().Msg("Tunneling handler shat down")
				}()

				// Fire the "connected" event
				sendResult(m, nil)
			case driverModel.Status_NO_MORE_CONNECTIONS:
				m.doSomethingAndClose(func() { sendResult(nil, errors.New("no more connections")) })
			default:
				m.doSomethingAndClose(func() { sendResult(nil, errors.Errorf("got a return status of: %s", connectionResponse.Status)) })
			}
		} else {
			m.doSomethingAndClose(func() { sendResult(nil, errors.New("this device doesn't support tunneling")) })
		}
	}()

	return result
}

func (m *Connection) doSomethingAndClose(something func()) {
	something()
	err := m.messageCodec.Disconnect()
	if err != nil {
		log.Warn().Msgf("error closing connection: %s", err)
	}
}

func (m *Connection) BlockingClose() {
	ttlTimer := time.NewTimer(m.defaultTtl)
	closeResults := m.Close()
	select {
	case <-closeResults:
		if !ttlTimer.Stop() {
			<-ttlTimer.C
		}
		return
	case <-ttlTimer.C:
		ttlTimer.Stop()
		return
	}
}

func (m *Connection) Close() <-chan plc4go.PlcConnectionCloseResult {
	result := make(chan plc4go.PlcConnectionCloseResult)

	go func() {
		// Stop the connection-state checker.
		if m.connectionStateTimer != nil {
			m.connectionStateTimer.Stop()
		}

		// Disconnect from all knx devices we are still connected to.
		for targetAddress := range m.DeviceConnections {
			ttlTimer := time.NewTimer(m.defaultTtl)
			disconnects := m.DeviceDisconnect(targetAddress)
			select {
			case _ = <-disconnects:
				if !ttlTimer.Stop() {
					<-ttlTimer.C
				}
			case <-ttlTimer.C:
				ttlTimer.Stop()
				// If we got a timeout here, well just continue the device will just auto disconnect.
				log.Debug().Msgf("Timeout disconnecting from device %s.", KnxAddressToString(&targetAddress))
			}
		}

		// Send a disconnect request from the gateway.
		_, err := m.sendGatewayDisconnectionRequest()
		if err != nil {
			result <- _default.NewDefaultPlcConnectionCloseResult(m, errors.Wrap(err, "got an error while disconnecting"))
		} else {
			result <- _default.NewDefaultPlcConnectionCloseResult(m, nil)
		}
	}()

	return result
}

func (m *Connection) IsConnected() bool {
	if m.messageCodec != nil {
		ttlTimer := time.NewTimer(m.defaultTtl)
		pingChannel := m.Ping()
		select {
		case pingResponse := <-pingChannel:
			if !ttlTimer.Stop() {
				<-ttlTimer.C
			}
			return pingResponse.GetErr() == nil
		case <-ttlTimer.C:
			ttlTimer.Stop()
			m.handleTimeout()
			return false
		}
	}
	return false
}

func (m *Connection) Ping() <-chan plc4go.PlcConnectionPingResult {
	result := make(chan plc4go.PlcConnectionPingResult)

	go func() {
		// Send the connection state request
		_, err := m.sendConnectionStateRequest()
		if err != nil {
			result <- _default.NewDefaultPlcConnectionPingResult(errors.Wrap(err, "got an error"))
		} else {
			result <- _default.NewDefaultPlcConnectionPingResult(nil)
		}
		return
	}()

	return result
}

func (m *Connection) GetMetadata() apiModel.PlcConnectionMetadata {
	return m.metadata
}

func (m *Connection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
	return internalModel.NewDefaultPlcReadRequestBuilder(
		m.fieldHandler, NewReader(m))
}

func (m *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return internalModel.NewDefaultPlcWriteRequestBuilder(
		m.fieldHandler, m.valueHandler, NewWriter(m.messageCodec))
}

func (m *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return internalModel.NewDefaultPlcSubscriptionRequestBuilder(
		m.fieldHandler, m.valueHandler, NewSubscriber(m))
}

func (m *Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	return internalModel.NewDefaultPlcBrowseRequestBuilder(NewBrowser(m, m.messageCodec))
}

func (m *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return nil /*internalModel.NewDefaultPlcUnsubscriptionRequestBuilder(
	  m.fieldHandler, m.valueHandler, NewSubscriber(m.messageCodec))*/
}

func (m *Connection) GetTransportInstance() transports.TransportInstance {
	if mc, ok := m.messageCodec.(spi.TransportInstanceExposer); ok {
		return mc.GetTransportInstance()
	}
	return nil
}

func (m *Connection) GetPlcFieldHandler() spi.PlcFieldHandler {
	return m.fieldHandler
}

func (m *Connection) GetPlcValueHandler() spi.PlcValueHandler {
	return m.valueHandler
}
