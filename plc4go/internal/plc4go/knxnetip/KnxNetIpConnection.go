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
    "github.com/apache/plc4x/plc4go/pkg/plc4go"
    apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
    "github.com/apache/plc4x/plc4go/pkg/plc4go/values"
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

type InternalKnxNetIpConnection interface {
    Send(request *driverModel.KnxNetIpMessage) error
    SendRequest(request *driverModel.KnxNetIpMessage, expect func(response interface{}) (bool, bool)) (int32, chan interface{})
}

type KnxReadResult struct {
    returnCode apiModel.PlcResponseCode
    value      *values.PlcValue
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
        defaultTtl:         time.Millisecond * 5000,
        DeviceConnections:  map[driverModel.KnxAddress]*KnxDeviceConnection{},
    }
    connection.messageCodec = NewKnxNetIpMessageCodec(transportInstance, connection.interceptIncomingMessage)
    return connection
}

func (m *KnxNetIpConnection) Connect() <-chan plc4go.PlcConnectionConnectResult {
    result := make(chan plc4go.PlcConnectionConnectResult)
    sendResult := func(connection plc4go.PlcConnection, err error) {
        select {
        case result <- plc4go.NewPlcConnectionConnectResult(connection, err):
        }
    }

    go func() {
        err := m.messageCodec.Connect()
        if err != nil {
            sendResult(nil, err)
            return
        }

        searchResponseChannel := m.sendGatewaySearchRequest()
        select {
        case searchResponse := <-searchResponseChannel:
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
                connectionResponseChannel := m.sendGatewayConnectionRequest()
                select {
                case connectionResponse := <-connectionResponseChannel:
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

                        // Start a timer that sends connection-state requests every 60 seconds
                        m.connectionStateTimer = time.NewTicker(60 * time.Second)
                        m.quitConnectionStateTimer = make(chan struct{})
                        go func() {
                            for {
                                select {
                                case <-m.connectionStateTimer.C:
                                    // We're using the connection-state-request as ping operation ...
                                    ping := m.Ping()
                                    select {
                                    case pingResult := <-ping:
                                        if pingResult.Err != nil {
                                            // TODO: Do some error handling here ...
                                            m.connectionStateTimer.Stop()
                                        }
                                    case <-time.After(5 * time.Second):
                                        // Close the connection
                                        m.Close()
                                    }

                                // If externally a request to stop the timer was issued, stop the timer.
                                case <-m.quitConnectionStateTimer:
                                    // TODO: Do some error handling here ...
                                    m.connectionStateTimer.Stop()
                                    return
                                }
                            }
                        }()

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
                                        fmt.Printf("Got an unhandled TunnelingResponse message %v\n", tunnelingResponse)
                                    } else {
                                        fmt.Printf("Not a TunnelingRequest message %v\n", incomingMessage)
                                    }
                                } else {
                                    if tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId != m.CommunicationChannelId {
                                        fmt.Printf("Not for this connection %v\n", tunnelingRequest)
                                        continue
                                    }

                                    lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi)
                                    if lDataInd != nil {
                                        // Get APDU, source and target address
                                        var apdu *driverModel.Apdu
                                        var sourceAddress *driverModel.KnxAddress
                                        //var targetAddress []int8
                                        lDataFrameData := driverModel.CastLDataFrameData(lDataInd.DataFrame)
                                        if lDataFrameData != nil {
                                            apdu = lDataFrameData.Apdu
                                            sourceAddress = lDataFrameData.SourceAddress
                                            //targetAddress = lDataFrameData.DestinationAddress
                                        } else {
                                            lDataFrameDataExt := driverModel.CastLDataFrameDataExt(lDataInd.DataFrame)
                                            if lDataFrameDataExt != nil {
                                                apdu = lDataFrameDataExt.Apdu
                                                sourceAddress = lDataFrameDataExt.SourceAddress
                                                //targetAddress = lDataFrameDataExt.DestinationAddress
                                            }
                                        }

                                        // If this is not an APDU, there is no need to further handle it.
                                        if apdu == nil {
                                            continue
                                        }

                                        // If this is an incoming disconnect request, remove the device
                                        // from the device connections, otherwise handle it as normal
                                        // incoming message.
                                        apduControlContainer := driverModel.CastApduControlContainer(apdu)
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
            }
        }
    }()
    return result
}

func (m *KnxNetIpConnection) Close() <-chan plc4go.PlcConnectionCloseResult {
    result := make(chan plc4go.PlcConnectionCloseResult)
    sendResult := func(connection plc4go.PlcConnection, err error) {
        select {
        case result <- plc4go.NewPlcConnectionCloseResult(connection, err):
        }
    }

    go func() {
        // Stop the connection-state checker.
        m.connectionStateTimer.Stop()

        // Disconnect from all knx devices we are still connected to.
        for address := range m.DeviceConnections {
            disconnects := m.sendDeviceDisconnectionRequest(address)
            select {
            case _ = <-disconnects:
                delete(m.DeviceConnections, address)
            }
        }

        // Send a disconnect request from the gateway.
        disconnectionResponseChannel := m.sendGatewayDisconnectionRequest()
        select {
        case disconnectResponse := <-disconnectionResponseChannel:
            if disconnectResponse.Status == driverModel.Status_NO_ERROR {
                sendResult(m, nil)
            } else {
                sendResult(m, errors.New("got an unexpected response for disconnect "+disconnectResponse.Status.String()))
            }
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
        case <-time.After(time.Second * 5):
            return false
        }
    }
    return false
}

func (m *KnxNetIpConnection) Ping() <-chan plc4go.PlcConnectionPingResult {
    result := make(chan plc4go.PlcConnectionPingResult)
    sendResult := func(err error) {
        select {
        case result <- plc4go.NewPlcConnectionPingResult(err):
        }
    }

    go func() {
        // Send the connection state request
        connectionStateResponseChannel := m.sendConnectionStateRequest()
        select {
        case connectionStateResponse := <-connectionStateResponseChannel:
            if connectionStateResponse.Status != driverModel.Status_NO_ERROR {
                sendResult(errors.New("got a failure response code " + strconv.Itoa(int(connectionStateResponse.Status))))
            } else {
                sendResult(nil)
            }
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
        case _ = <-controlConnectResponseChannel:
            // If the connection request was successful, try to read the device-descriptor
            deviceDescriptorResponses := m.sendDeviceDeviceDescriptorReadRequest(targetAddress)
            select {
            case _ = <-deviceDescriptorResponses:
                // Last, not least, read the max APDU size
                propertyValueResponses := m.sendDevicePropertyReadRequest(targetAddress, 0, 56)
                select {
                case propertyValueResponse := <-propertyValueResponses:
                    // If we were able to read the max APDU size, then use the minimum of
                    // the connection APDU size and the device APDU size, otherwise use the
                    // default APDU Size
                    deviceApduSize := uint16(100)
                    if propertyValueResponse.returnCode == apiModel.PlcResponseCode_OK {
                        deviceApduSize = (*propertyValueResponse.value).GetUint16()
                    }
                    connection := m.DeviceConnections[targetAddress]
                    connection.maxApdu = uint16(math.Min(float64(deviceApduSize), 240))
                    result <- connection
                }
            }
        }
    }()

    return result
}

func (m *KnxNetIpConnection) DisconnectFromDevice(targetAddress driverModel.KnxAddress) <-chan *KnxDeviceConnection {
    result := make(chan *KnxDeviceConnection)

    if connection, ok := m.DeviceConnections[targetAddress]; ok {
        disconnects := m.sendDeviceDisconnectionRequest(targetAddress)
        select {
        case _ = <-disconnects:
            result <- connection
        }
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
            }
        }

        // If we successfully got a connection, read the property
        if connection != nil {
            propertyValueResponses := m.sendDevicePropertyReadRequest(targetAddress, objectId, propertyId)
            select {
            case propertyValueResponse := <-propertyValueResponses:
                result <- propertyValueResponse
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

func (m *KnxNetIpConnection) sendGatewaySearchRequest() chan driverModel.SearchResponse {
    result := make(chan driverModel.SearchResponse)
    localAddress, err := m.getLocalAddress()
    if err != nil {
        close(result)
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
            result <- *searchResponse
            close(result)
            return nil
        },
        m.defaultTtl)
    if err != nil {
        close(result)
    }
    return result
}

func (m *KnxNetIpConnection) sendGatewayConnectionRequest() chan driverModel.ConnectionResponse {
    result := make(chan driverModel.ConnectionResponse)
    localAddress, err := m.getLocalAddress()
    if err != nil {
        close(result)
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
            result <- *connectionResponse
            close(result)
            return nil
        },
        m.defaultTtl)
    if err != nil {
        close(result)
    }
    return result
}

func (m *KnxNetIpConnection) sendGatewayDisconnectionRequest() chan driverModel.DisconnectResponse {
    result := make(chan driverModel.DisconnectResponse)
    localAddress, err := m.getLocalAddress()
    if err != nil {
        close(result)
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
            result <- *disconnectResponse
            close(result)
            return nil
        },
        m.defaultTtl)
    if err != nil {
        close(result)
    }
    return result
}

func (m *KnxNetIpConnection) sendConnectionStateRequest() chan driverModel.ConnectionStateResponse {
    result := make(chan driverModel.ConnectionStateResponse)
    localAddress, err := m.getLocalAddress()
    if err != nil {
        close(result)
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
            result <- *connectionStateResponse
            close(result)
            return nil
        },
        m.defaultTtl)
    if err != nil {
        close(result)
    }
    return result
}

func (m *KnxNetIpConnection) sendDeviceConnectionRequest(targetAddress driverModel.KnxAddress) chan driverModel.ApduControlConnect {
    result := make(chan driverModel.ApduControlConnect)

    go func() {
        // Send a connection request to the individual KNX device
        deviceConnectionRequest := driverModel.NewTunnelingRequest(
            driverModel.NewTunnelingRequestDataBlock(0, 0),
            driverModel.NewLDataReq(0, nil,
                driverModel.NewLDataFrameDataExt(false, 6, uint8(0),
                    driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
                    driverModel.NewApduControlContainer(driverModel.NewApduControlConnect(), 0, false, 0),
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
                lDataFrameExt := driverModel.CastLDataFrameDataExt(lDataCon.DataFrame)
                if lDataFrameExt == nil {
                    return false
                }
                // TODO: Check this too: lDataFrameExt.DestinationAddress
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
                    close(result)
                }
                lDataFrameExt := driverModel.CastLDataFrameDataExt(lDataCon.DataFrame)
                apduControlContainer := driverModel.CastApduControlContainer(lDataFrameExt.Apdu)
                apduControlConnect := driverModel.CastApduControlConnect(apduControlContainer.ControlApdu)

                // Create a new connection object and save it in the map
                deviceConnection := &KnxDeviceConnection{
                    counter: 0,
                    maxApdu: 0, // TODO: Initialize this with the default max APDU Size
                }
                m.DeviceConnections[targetAddress] = deviceConnection

                result <- *apduControlConnect
                return nil
            },
            m.defaultTtl)
        if err != nil {
            close(result)
        }
    }()

    return result
}

func (m *KnxNetIpConnection) sendDeviceDisconnectionRequest(targetAddress driverModel.KnxAddress) chan driverModel.ApduControlDisconnect {
    result := make(chan driverModel.ApduControlDisconnect)

    go func() {
        // Send a connection request to the individual KNX device
        deviceDisconnectionRequest := driverModel.NewTunnelingRequest(
            driverModel.NewTunnelingRequestDataBlock(0, 0),
            driverModel.NewLDataReq(0, nil,
                driverModel.NewLDataFrameDataExt(false, 6, uint8(0),
                    driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
                    driverModel.NewApduControlContainer(driverModel.NewApduControlDisconnect(), 0, false, 0),
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
                lDataFrameExt := driverModel.CastLDataFrameDataExt(lDataCon.DataFrame)
                if lDataFrameExt == nil {
                    return false
                }
                // TODO: Check this too: lDataFrameExt.DestinationAddress
                apduControlContainer := driverModel.CastApduControlContainer(lDataFrameExt.Apdu)
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
                    close(result)
                }
                lDataFrameExt := driverModel.CastLDataFrameDataExt(lDataCon.DataFrame)
                apduControlContainer := driverModel.CastApduControlContainer(lDataFrameExt.Apdu)
                apduControlDisconnect := driverModel.CastApduControlDisconnect(apduControlContainer.ControlApdu)

                result <- *apduControlDisconnect
                return nil
            },
            m.defaultTtl)
        if err != nil {
            close(result)
        }
    }()

    return result
}

func (m *KnxNetIpConnection) sendDeviceDeviceDescriptorReadRequest(targetAddress driverModel.KnxAddress) chan driverModel.ApduDataDeviceDescriptorResponse {
    result := make(chan driverModel.ApduDataDeviceDescriptorResponse)

    connection, ok := m.DeviceConnections[targetAddress]
    if !ok {
        close(result)
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
            driverModel.NewLDataFrameDataExt(false, 6, uint8(0),
                driverModel.NewKnxAddress(0, 0, 0),
                KnxAddressToInt8Array(targetAddress),
                driverModel.NewApduDataContainer(
                    driverModel.NewApduDataDeviceDescriptorRead(0), 1, true, counter),
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
            dataFrame := driverModel.CastLDataFrameDataExt(lDataInd.DataFrame)
            if dataFrame == nil {
                return false
            }
            // TODO: Check this too: dataFrame.SourceAddress
            // TODO: Check this too: dataFrame.Apdu.Counter
            dataContainer := driverModel.CastApduDataContainer(dataFrame.Apdu)
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
            dataFrame := driverModel.CastLDataFrameDataExt(lDataInd.DataFrame)
            dataContainer := driverModel.CastApduDataContainer(dataFrame.Apdu)
            deviceDescriptorResponse := driverModel.CastApduDataDeviceDescriptorResponse(dataContainer.DataApdu)

            // Send an ACK and wait for it to be delivered
            go func() {
                ackResults := m.sendDeviceAck(targetAddress, dataFrame.Apdu.Counter)
                select {
                case ackResult := <-ackResults:
                    if !ackResult {
                        fmt.Printf("Error getting Ack")
                    }
                }

                // Save the device-descriptor value
                deviceDescriptor := uint16(deviceDescriptorResponse.Data[0])<<8 | uint16(deviceDescriptorResponse.Data[1])
                connection.deviceDescriptor = deviceDescriptor

                result <- *deviceDescriptorResponse
            }()

            return nil
        },
        time.Second*5)
    if err != nil {
        close(result)
    }
    return result
}

func (m *KnxNetIpConnection) sendDevicePropertyReadRequest(targetAddress driverModel.KnxAddress, objectId uint8, propertyId uint8) chan KnxReadResult {
    result := make(chan KnxReadResult)

    connection, ok := m.DeviceConnections[targetAddress]
    if !ok {
        close(result)
        result <- KnxReadResult{
            returnCode: apiModel.PlcResponseCode_INTERNAL_ERROR,
        }
    }

    // Next, read the device descriptor so we know how we have to communicate with the device.
    counter := connection.counter
    connection.counter++

    // Send the property read request and wait for a confirmation that this property is readable.
    propertyReadRequest := driverModel.NewTunnelingRequest(
        driverModel.NewTunnelingRequestDataBlock(0, 0),
        driverModel.NewLDataReq(0, nil,
            driverModel.NewLDataFrameDataExt(false, 6, 0,
                driverModel.NewKnxAddress(0, 0, 0),
                KnxAddressToInt8Array(targetAddress),
                driverModel.NewApduDataContainer(
                    driverModel.NewApduDataOther(
                        // TODO: The counter should be incremented per KNX individual address
                        driverModel.NewApduDataExtPropertyValueRead(objectId, propertyId, 1, 1)),
                    5, true, counter),
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
            var apdu *driverModel.Apdu
            dataFrameExt := driverModel.CastLDataFrameDataExt(lDataInd.DataFrame)
            if dataFrameExt != nil {
                apdu = dataFrameExt.Apdu
            } else {
                dataFrame := driverModel.CastLDataFrameData(lDataInd.DataFrame)
                if dataFrame != nil {
                    apdu = dataFrame.Apdu
                }
            }
            if apdu == nil {
                return false
            }
            // TODO: Check this too: dataFrame.SourceAddress
            // TODO: Check this too: dataFrame.Apdu.Counter
            dataContainer := driverModel.CastApduDataContainer(apdu)
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
            var apdu *driverModel.Apdu
            dataFrameExt := driverModel.CastLDataFrameDataExt(lDataInd.DataFrame)
            if dataFrameExt != nil {
                apdu = dataFrameExt.Apdu
            } else {
                dataFrame := driverModel.CastLDataFrameData(lDataInd.DataFrame)
                if dataFrame != nil {
                    apdu = dataFrame.Apdu
                }
            }
            dataContainer := driverModel.CastApduDataContainer(apdu)
            dataApduOther := driverModel.CastApduDataOther(dataContainer.DataApdu)
            propertyValueResponse := driverModel.CastApduDataExtPropertyValueResponse(dataApduOther.ExtendedApdu)

            // Send an ACK and wait for it to be delivered
            go func() {
                ackResults := m.sendDeviceAck(targetAddress, apdu.Counter)
                select {
                case ackResult := <-ackResults:
                    if !ackResult {
                        fmt.Printf("Error getting Ack")
                    }
                }

                // If the count is 0, then this property doesn't exist or the user has no permission to read it.
                if propertyValueResponse.Count == 0 {
                    result <- KnxReadResult{
                        returnCode: apiModel.PlcResponseCode_NOT_FOUND,
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

                    // Parse the data according to the property type information
                    propertyDataType := property.PropertyDataType()
                    dataLength := uint8(len(propertyValueResponse.Data))
                    data := propertyValueResponse.Data
                    rb := utils.NewReadBuffer(data)
                    plcValue, err := driverModel.KnxPropertyParse(rb, propertyDataType, dataLength)

                    // Return the result
                    if err != nil {
                        result <- KnxReadResult{
                            returnCode: apiModel.PlcResponseCode_INTERNAL_ERROR,
                            value:      nil,
                        }
                    } else {
                        result <- KnxReadResult{
                            returnCode: apiModel.PlcResponseCode_OK,
                            value:      &plcValue,
                        }
                    }
                }

            }()
            return nil
        },
        time.Second*5)

    if err != nil {
        close(result)
    }
    return result
}

func (m *KnxNetIpConnection) sendDeviceAck(targetAddress driverModel.KnxAddress, counter uint8) chan bool {
    result := make(chan bool)

    ack := driverModel.NewTunnelingRequest(
        driverModel.NewTunnelingRequestDataBlock(0, 0),
        driverModel.NewLDataReq(0, nil,
            driverModel.NewLDataFrameDataExt(false, 6, uint8(0),
                driverModel.NewKnxAddress(0, 0, 0), KnxAddressToInt8Array(targetAddress),
                driverModel.NewApduControlContainer(driverModel.NewApduControlAck(), 0, true, counter),
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
            dataFrame := driverModel.CastLDataFrameDataExt(lDataCon.DataFrame)
            if dataFrame == nil {
                return false
            }
            // TODO: Check this too: dataFrame.SourceAddress (This should match the targetAddress)
            // TODO: Check this too: dataFrame.DestinationAddress (This should match our KNX Address)
            // TODO: Check this too: dataFrame.Apdu.Counter
            controlContainer := driverModel.CastApduControlContainer(dataFrame.Apdu)
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
            result <- true
            return nil
        },
        time.Second*5)

    if err != nil {
        result <- false
    }

    return result
}

func (m *KnxNetIpConnection) sendRequest(request *driverModel.KnxNetIpMessage, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, ttl time.Duration) error {
    // If this is a tunneling request, we need to update the communicationChannelId and assign a sequenceCounter
    tunnelingRequest := driverModel.CastTunnelingRequest(request)
    if tunnelingRequest != nil {
        tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId = m.CommunicationChannelId
        tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter = m.getNewSequenceCounter()
    }
    return m.messageCodec.SendRequest(request, acceptsMessage, handleMessage, ttl)
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
            var apdu *driverModel.Apdu
            switch lDataInd.DataFrame.Child.(type) {
            case *driverModel.LDataFrameData:
                dataFrame := driverModel.CastLDataFrameData(lDataInd.DataFrame)
                destinationAddress = dataFrame.DestinationAddress
                apdu = dataFrame.Apdu
            case *driverModel.LDataFrameDataExt:
                dataFrame := driverModel.CastLDataFrameDataExt(lDataInd.DataFrame)
                destinationAddress = dataFrame.DestinationAddress
                apdu = dataFrame.Apdu
            }
            container := driverModel.CastApduDataContainer(apdu)
            if container == nil {
                return
            }
            groupValueWrite := driverModel.CastApduDataGroupValueWrite(container.DataApdu)
            if groupValueWrite == nil {
                return
            }
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
