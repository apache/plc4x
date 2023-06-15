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
	"context"
	"encoding/hex"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/tracer"
	"github.com/rs/zerolog"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	_default "github.com/apache/plc4x/plc4go/spi/default"

	"github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/interceptors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/transports"
	"github.com/pkg/errors"
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
		"DeviceSerialNumber":     ByteArrayToString(m.DeviceSerialNumber, " "),
		"DeviceMulticastAddress": ByteArrayToString(m.DeviceSerialNumber, "."),
		"DeviceMacAddress":       ByteArrayToString(m.DeviceSerialNumber, ":"),
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
	tagHandler               spi.PlcTagHandler
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

	GatewayKnxAddress             driverModel.KnxAddress
	ClientKnxAddress              driverModel.KnxAddress
	CommunicationChannelId        uint8
	SequenceCounter               int32
	TunnelingRequestExpectationId int32
	DeviceConnections             map[driverModel.KnxAddress]*KnxDeviceConnection

	requestInterceptor interceptors.RequestInterceptor
	sync.Mutex

	// indicates if the tunneling requests loop is running
	handleTunnelingRequests bool

	connectionId string
	tracer       tracer.Tracer

	passLogToModel bool
	log            zerolog.Logger
}

func (m *Connection) String() string {
	return fmt.Sprintf("knx.Connection{}")
}

type KnxReadResult struct {
	value    values.PlcValue
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
	responseMessage spi.Message
	err             error
}

func NewConnection(transportInstance transports.TransportInstance, connectionOptions map[string][]string, tagHandler spi.PlcTagHandler, _options ...options.WithOption) *Connection {
	connection := &Connection{
		options:      connectionOptions,
		tagHandler:   tagHandler,
		valueHandler: NewValueHandler(),
		requestInterceptor: interceptors.NewSingleItemRequestInterceptor(
			spiModel.NewDefaultPlcReadRequest,
			spiModel.NewDefaultPlcWriteRequest,
			spiModel.NewDefaultPlcReadResponse,
			spiModel.NewDefaultPlcWriteResponse,
			_options...,
		),
		subscribers:             []*Subscriber{},
		valueCache:              map[uint16][]byte{},
		valueCacheMutex:         sync.RWMutex{},
		metadata:                &ConnectionMetadata{},
		defaultTtl:              time.Second * 10,
		DeviceConnections:       map[driverModel.KnxAddress]*KnxDeviceConnection{},
		handleTunnelingRequests: true,
		passLogToModel:          options.ExtractPassLoggerToModel(_options...),
		log:                     options.ExtractCustomLogger(_options...),
	}
	connection.connectionTtl = connection.defaultTtl * 2

	if traceEnabledOption, ok := connectionOptions["traceEnabled"]; ok {
		if len(traceEnabledOption) == 1 {
			connection.tracer = tracer.NewTracer(connection.connectionId, _options...)
		}
	}
	// If a building key was provided, save that in a dedicated variable
	if buildingKey, ok := connectionOptions["buildingKey"]; ok {
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

func (m *Connection) GetTracer() tracer.Tracer {
	return m.tracer
}

func (m *Connection) Connect() <-chan plc4go.PlcConnectionConnectResult {
	return m.ConnectWithContext(context.Background())
}

func (m *Connection) ConnectWithContext(ctx context.Context) <-chan plc4go.PlcConnectionConnectResult {
	result := make(chan plc4go.PlcConnectionConnectResult, 1)
	sendResult := func(connection plc4go.PlcConnection, err error) {
		result <- _default.NewDefaultPlcConnectionConnectResult(connection, err)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		// Open the UDP Connection
		err := m.messageCodec.Connect()
		if err != nil {
			m.doSomethingAndClose(func() { sendResult(nil, errors.Wrap(err, "error opening connection")) })
			return
		}

		// Send a search request before connecting to the device.
		searchResponse, err := m.sendGatewaySearchRequest(ctx)
		if err != nil {
			m.doSomethingAndClose(func() { sendResult(nil, errors.Wrap(err, "error discovering device capabilities")) })
			return
		}

		// Save some important information
		dibDeviceInfo := searchResponse.GetDibDeviceInfo()
		m.metadata.KnxMedium = dibDeviceInfo.GetKnxMedium()
		m.metadata.GatewayName = string(bytes.Trim(dibDeviceInfo.GetDeviceFriendlyName(), "\x00"))
		m.GatewayKnxAddress = dibDeviceInfo.GetKnxAddress()
		m.metadata.GatewayKnxAddress = KnxAddressToString(m.GatewayKnxAddress)
		m.metadata.ProjectNumber = dibDeviceInfo.GetProjectInstallationIdentifier().GetProjectNumber()
		m.metadata.InstallationNumber = dibDeviceInfo.GetProjectInstallationIdentifier().GetInstallationNumber()
		m.metadata.DeviceSerialNumber = dibDeviceInfo.GetKnxNetIpDeviceSerialNumber()
		m.metadata.DeviceMulticastAddress = dibDeviceInfo.GetKnxNetIpDeviceMulticastAddress().GetAddr()
		m.metadata.DeviceMacAddress = dibDeviceInfo.GetKnxNetIpDeviceMacAddress().GetAddr()
		m.metadata.SupportedServices = []string{}
		supportsTunneling := false
		for _, serviceId := range searchResponse.GetDibSuppSvcFamilies().GetServiceIds() {
			m.metadata.SupportedServices = append(m.metadata.SupportedServices, serviceId.(interface{ GetTypeName() string }).GetTypeName())
			// If this is an instance of the "tunneling", service, this connection supports tunneling
			_, ok := serviceId.(driverModel.KnxNetIpTunneling)
			if ok {
				supportsTunneling = true
				break
			}
		}

		// If the current device supports tunneling, create a tunneling connection.
		// Via this connection we then get access to the entire KNX network this Gateway is connected to.
		if supportsTunneling {
			// As soon as we got a successful search-response back, send a connection request.
			connectionResponse, err := m.sendGatewayConnectionRequest(ctx)
			if err != nil {
				m.doSomethingAndClose(func() { sendResult(nil, errors.Wrap(err, "error connecting to device")) })
				return
			}

			// Save the communication channel id
			m.CommunicationChannelId = connectionResponse.GetCommunicationChannelId()

			// Reset the sequence counter
			m.SequenceCounter = -1

			// If the connection was successful, the gateway will now forward any packets
			// on the KNX bus that are broadcast packets to us, so we have to setup things
			// to handle these incoming messages.
			switch connectionResponse.GetStatus() {
			case driverModel.Status_NO_ERROR:
				// Save the KNX Address the Gateway assigned to us for this connection.
				tunnelConnectionDataBlock := connectionResponse.GetConnectionResponseDataBlock().(driverModel.ConnectionResponseDataBlockTunnelConnection)
				m.ClientKnxAddress = tunnelConnectionDataBlock.GetKnxAddress()

				// Create a go routine to handle incoming tunneling-requests which haven't been
				// handled by any other handler. This is where usually the GroupValueWrite messages
				// are being handled.
				m.log.Debug().Msg("Starting tunneling handler")
				go func() {
					defer func() {
						if err := recover(); err != nil {
							m.log.Error().Msgf("panic-ed %v. Stack: %s", err, debug.Stack())
						}
					}()
					defaultIncomingMessageChannel := m.messageCodec.GetDefaultIncomingMessageChannel()
					for m.handleTunnelingRequests {
						incomingMessage := <-defaultIncomingMessageChannel
						tunnelingRequest, ok := incomingMessage.(driverModel.TunnelingRequestExactly)
						if !ok {
							tunnelingResponse, ok := incomingMessage.(driverModel.TunnelingResponseExactly)
							if ok {
								m.log.Warn().Msgf("Got an unhandled TunnelingResponse message %v\n", tunnelingResponse)
							} else {
								m.log.Warn().Msgf("Not a TunnelingRequest or TunnelingResponse message %v\n", incomingMessage)
							}
							continue
						}

						if tunnelingRequest.GetTunnelingRequestDataBlock().GetCommunicationChannelId() != m.CommunicationChannelId {
							m.log.Warn().Msgf("Not for this connection %v\n", tunnelingRequest)
							continue
						}

						lDataInd, ok := tunnelingRequest.GetCemi().(driverModel.LDataIndExactly)
						if !ok {
							continue
						}
						// Get APDU, source and target address
						lDataFrameData := lDataInd.GetDataFrame().(driverModel.LDataExtended)
						sourceAddress := lDataFrameData.GetSourceAddress()

						// If this is not an APDU, there is no need to further handle it.
						if lDataFrameData.GetApdu() == nil {
							continue
						}

						// If this is an incoming disconnect request, remove the device
						// from the device connections, otherwise handle it as normal
						// incoming message.
						apduControlContainer, ok := lDataFrameData.GetApdu().(driverModel.ApduControlContainerExactly)
						if ok {
							_, ok := apduControlContainer.GetControlApdu().(driverModel.ApduControlDisconnectExactly)
							if ok {
								if m.DeviceConnections[sourceAddress] != nil /* && m.ClientKnxAddress == Int8ArrayToKnxAddress(targetAddress)*/ {
									// Remove the connection
									delete(m.DeviceConnections, sourceAddress)
								}
							}
						} else {
							m.handleIncomingTunnelingRequest(ctx, tunnelingRequest)
						}
					}
					m.log.Warn().Msg("Tunneling handler shat down")
				}()

				// Fire the "connected" event
				sendResult(m, nil)
			case driverModel.Status_NO_MORE_CONNECTIONS:
				m.doSomethingAndClose(func() { sendResult(nil, errors.New("no more connections")) })
			default:
				m.doSomethingAndClose(func() { sendResult(nil, errors.Errorf("got a return status of: %s", connectionResponse.GetStatus())) })
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
		m.log.Warn().Msgf("error closing connection: %s", err)
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
	// TODO: use proper context
	ctx := context.TODO()
	result := make(chan plc4go.PlcConnectionCloseResult, 1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- _default.NewDefaultPlcConnectionConnectResult(nil, errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		// Stop the connection-state checker.
		if m.connectionStateTimer != nil {
			m.connectionStateTimer.Stop()
		}

		// Disconnect from all knx devices we are still connected to.
		for targetAddress := range m.DeviceConnections {
			ttlTimer := time.NewTimer(m.defaultTtl)
			disconnects := m.DeviceDisconnect(ctx, targetAddress)
			select {
			case _ = <-disconnects:
				if !ttlTimer.Stop() {
					<-ttlTimer.C
				}
			case <-ttlTimer.C:
				ttlTimer.Stop()
				// If we got a timeout here, well just continue the device will just auto disconnect.
				m.log.Debug().Msgf("Timeout disconnecting from device %s.", KnxAddressToString(targetAddress))
			}
		}

		// Send a disconnect request from the gateway.
		_, err := m.sendGatewayDisconnectionRequest(ctx)
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
	// TODO: use proper context
	ctx := context.TODO()
	result := make(chan plc4go.PlcConnectionPingResult, 1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				result <- _default.NewDefaultPlcConnectionPingResult(errors.Errorf("panic-ed %v. Stack: %s", err, debug.Stack()))
			}
		}()
		// Send the connection state request
		_, err := m.sendConnectionStateRequest(ctx)
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
	return spiModel.NewDefaultPlcReadRequestBuilder(
		m.tagHandler, NewReader(m))
}

func (m *Connection) WriteRequestBuilder() apiModel.PlcWriteRequestBuilder {
	return spiModel.NewDefaultPlcWriteRequestBuilder(
		m.tagHandler, m.valueHandler, NewWriter(m.messageCodec))
}

func (m *Connection) SubscriptionRequestBuilder() apiModel.PlcSubscriptionRequestBuilder {
	return spiModel.NewDefaultPlcSubscriptionRequestBuilder(
		m.tagHandler, m.valueHandler, NewSubscriber(m, options.WithCustomLogger(m.log)))
}

func (m *Connection) BrowseRequestBuilder() apiModel.PlcBrowseRequestBuilder {
	return spiModel.NewDefaultPlcBrowseRequestBuilder(m.tagHandler, NewBrowser(m, m.messageCodec))
}

func (m *Connection) UnsubscriptionRequestBuilder() apiModel.PlcUnsubscriptionRequestBuilder {
	return nil /*spiModel.NewDefaultPlcUnsubscriptionRequestBuilder(
	  m.tagHandler, m.valueHandler, NewSubscriber(m.messageCodec))*/
}

func (m *Connection) GetTransportInstance() transports.TransportInstance {
	if mc, ok := m.messageCodec.(spi.TransportInstanceExposer); ok {
		return mc.GetTransportInstance()
	}
	return nil
}

func (m *Connection) GetPlcTagHandler() spi.PlcTagHandler {
	return m.tagHandler
}

func (m *Connection) GetPlcValueHandler() spi.PlcValueHandler {
	return m.valueHandler
}

func ByteArrayToString(data []byte, separator string) string {
	var sb strings.Builder
	if data != nil {
		for i, element := range data {
			sb.WriteString(strconv.Itoa(int(element)))
			if i < (len(data) - 1) {
				sb.WriteString(separator)
			}
		}
	}
	return sb.String()
}
