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

    GatewayKnxAddress             *driverModel.KnxAddress
    ClientKnxAddress              *driverModel.KnxAddress
    CommunicationChannelId        uint8
    SequenceCounter               int32
    TunnelingRequestExpectationId int32

    requestInterceptor internalModel.RequestInterceptor
    plc4go.PlcConnection
}

type InternalKnxNetIpConnection interface {
    SendRequest(request *driverModel.KnxNetIpMessage, expect func(response interface{}) (bool, bool)) (int32, chan interface{})
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
    }
    connection.messageCodec = NewKnxNetIpMessageCodec(transportInstance, connection.interceptIncomingMessage)
    return connection
}

func (m *KnxNetIpConnection) Connect() <-chan plc4go.PlcConnectionConnectResult {
    ch := make(chan plc4go.PlcConnectionConnectResult)
    go func() {
        err := m.messageCodec.Connect()
        if err != nil {
            ch <- plc4go.NewPlcConnectionConnectResult(m, err)
            return
        }

        transportInstanceExposer, ok := m.messageCodec.(spi.TransportInstanceExposer)
        if !ok {
            ch <- plc4go.NewPlcConnectionConnectResult(m, errors.New(
                "used transport, is not a TransportInstanceExposer"))
            return
        }

        // Prepare a SearchReq
        udpTransportInstance, ok := transportInstanceExposer.GetTransportInstance().(*udp.UdpTransportInstance)
        if !ok {
            ch <- plc4go.NewPlcConnectionConnectResult(m, errors.New(
                "used transport, is not a UdpTransportInstance"))
            return
        }
        localAddress := driverModel.NewIPAddress(utils.ByteArrayToInt8Array(udpTransportInstance.LocalAddress.IP))
        discoveryEndpoint := driverModel.NewHPAIDiscoveryEndpoint(
            driverModel.HostProtocolCode_IPV4_UDP, localAddress, uint16(udpTransportInstance.LocalAddress.Port))
        searchRequest := driverModel.NewSearchRequest(discoveryEndpoint)

        // Send the SearchReq
        err = m.messageCodec.SendRequest(
            searchRequest,
            func(message interface{}) bool {
                searchResponse := driverModel.CastSearchResponse(message)
                return searchResponse != nil
            },
            func(message interface{}) error {
                searchResponse := driverModel.CastSearchResponse(message)

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
                    localAddress := m.castIpToKnxAddress(udpTransportInstance.LocalAddress.IP)
                    connectionRequest := driverModel.NewConnectionRequest(
                        driverModel.NewHPAIDiscoveryEndpoint(driverModel.HostProtocolCode_IPV4_UDP,
                            localAddress, uint16(udpTransportInstance.LocalAddress.Port)),
                        driverModel.NewHPAIDataEndpoint(driverModel.HostProtocolCode_IPV4_UDP,
                            localAddress, uint16(udpTransportInstance.LocalAddress.Port)),
                        driverModel.NewConnectionRequestInformationTunnelConnection(driverModel.KnxLayer_TUNNEL_LINK_LAYER),
                    )

                    // Send the connection request
                    err = m.messageCodec.SendRequest(
                        connectionRequest,
                        func(message interface{}) bool {
                            connectionResponse := driverModel.CastConnectionResponse(message)
                            return connectionResponse != nil
                        },
                        func(message interface{}) error {
                            connectionResponse := driverModel.CastConnectionResponse(message)

                            // Save the communication channel id
                            m.CommunicationChannelId = connectionResponse.CommunicationChannelId

                            // Reset the sequence counter
                            m.SequenceCounter = -1

                            // If the connection was successful, the gateway will now forward any packets
                            // on the KNX bus that are broadcast packets to us, so we have to setup things
                            // to handle these incoming messages.
                            if connectionResponse.Status == driverModel.Status_NO_ERROR {
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
                                                m.handleIncomingTunnelingRequest(tunnelingRequest)
                                            } else {
                                                fmt.Printf("Not a LDataInd message %v\n", tunnelingRequest.Cemi)
                                            }
                                        }
                                    }
                                }()

                                // Save the KNX Address the Gateway assigned to us for this connection.
                                tunnelConnectionDataBlock :=
                                    driverModel.CastConnectionResponseDataBlockTunnelConnection(
                                        connectionResponse.ConnectionResponseDataBlock)
                                m.ClientKnxAddress = tunnelConnectionDataBlock.KnxAddress

                                // Fire the "connected" event
                                ch <- plc4go.NewPlcConnectionConnectResult(m, nil)

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
                                            case <-time.After(1 * time.Second):
                                                m.connectionStateTimer.Stop()
                                            }

                                        // If externally a request to stop the timer was issued, stop the timer.
                                        case <-m.quitConnectionStateTimer:
                                            // TODO: Do some error handling here ...
                                            m.connectionStateTimer.Stop()
                                            return
                                        }
                                    }
                                }()
                            }
                            return nil
                        },
                        time.Second*1)
                }
                return nil
            },
            time.Second*1)
        if err != nil {
            ch <- plc4go.NewPlcConnectionConnectResult(nil, err)
        }
    }()
    return ch
}

func (m *KnxNetIpConnection) Close() <-chan plc4go.PlcConnectionCloseResult {
    // TODO: Implement ...
    ch := make(chan plc4go.PlcConnectionCloseResult)
    go func() {
        ch <- plc4go.NewPlcConnectionCloseResult(m, nil)
    }()
    return ch
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
    //	diagnosticRequestPdu := driverModel.NewModbusPDUDiagnosticRequest(0, 0x42)
    go func() {
        transportInstanceExposer, ok := m.messageCodec.(spi.TransportInstanceExposer)
        if !ok {
            result <- plc4go.NewPlcConnectionPingResult(errors.New(
                "used transport, is not a TransportInstanceExposer"))
            return
        }

        // Prepare a SearchReq
        udpTransportInstance, ok := transportInstanceExposer.GetTransportInstance().(*udp.UdpTransportInstance)
        if !ok {
            result <- plc4go.NewPlcConnectionPingResult(errors.New(
                "used transport, is not a UdpTransportInstance"))
            return
        }

        localAddress := m.castIpToKnxAddress(udpTransportInstance.LocalAddress.IP)

        connectionStateRequest := driverModel.NewConnectionStateRequest(
            m.CommunicationChannelId,
            driverModel.NewHPAIControlEndpoint(
                driverModel.HostProtocolCode_IPV4_UDP,
                localAddress, uint16(udpTransportInstance.LocalAddress.Port)))

        // Send the connection state request
        err := m.messageCodec.SendRequest(
            connectionStateRequest,
            func(message interface{}) bool {
                fmt.Printf("Check for Connection-State-Response ...")
                connectionStateResponse := driverModel.CastConnectionStateResponse(message)
                fmt.Printf(" Result %t\n", connectionStateResponse != nil)
                return connectionStateResponse != nil
            },
            func(message interface{}) error {
                connectionStateResponse := driverModel.CastConnectionStateResponse(message)
                if connectionStateResponse.Status != driverModel.Status_NO_ERROR {
                    result <- plc4go.NewPlcConnectionPingResult(errors.New(
                        "got a failure response code " + strconv.Itoa(int(connectionStateResponse.Status))))
                } else {
                    result <- plc4go.NewPlcConnectionPingResult(nil)
                }
                return nil
            },
            time.Second*1)
        result <- plc4go.NewPlcConnectionPingResult(err)
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

func (m *KnxNetIpConnection) SendRequest(request *driverModel.KnxNetIpMessage, acceptsMessage spi.AcceptsMessage, handleMessage spi.HandleMessage, ttl time.Duration) error {
    // If this is a tunneling request, we need to update the communicationChannelId and assign a sequenceCounter
    tunnelingRequest := driverModel.CastTunnelingRequest(request)
    if tunnelingRequest != nil {
        tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId = m.CommunicationChannelId
        tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter = m.getNewSequenceCounter()
    }
    return m.messageCodec.SendRequest(request, acceptsMessage, handleMessage, ttl)
}

func (m *KnxNetIpConnection) interceptIncomingMessage(message interface{}) {
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
    // Send an Ack response for this message
    tunnelingResponse := driverModel.NewTunnelingResponse(driverModel.NewTunnelingResponseDataBlock(
        m.CommunicationChannelId, tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
        driverModel.Status_NO_ERROR))
    err := m.messageCodec.Send(tunnelingResponse)
    if err != nil {
        return
    }

    go func() {
        lDataInd := driverModel.CastLDataInd(tunnelingRequest.Cemi.Child)
        if lDataInd != nil {
            var destinationAddress []int8
            var dataFirstByte *int8
            var data []int8
            switch lDataInd.DataFrame.Child.(type) {
            case *driverModel.LDataFrameData:
                dataFrame := driverModel.CastLDataFrameData(lDataInd.DataFrame)
                destinationAddress = dataFrame.DestinationAddress
                dataFirstByte = dataFrame.DataFirstByte
                data = dataFrame.Data
            case *driverModel.LDataFrameDataExt:
                dataFrame := driverModel.CastLDataFrameDataExt(lDataInd.DataFrame)
                destinationAddress = dataFrame.DestinationAddress
                dataFirstByte = dataFrame.DataFirstByte
                data = dataFrame.Data
            }
            if destinationAddress != nil {
                addressData := uint16(destinationAddress[0])<<8 | (uint16(destinationAddress[1]) & 0xFF)
                m.valueCacheMutex.RLock()
                val, ok := m.valueCache[addressData]
                m.valueCacheMutex.RUnlock()
                changed := false
                if dataFirstByte != nil {
                    var payload []int8
                    payload = append(payload, *dataFirstByte)
                    payload = append(payload, data...)
                    if !ok || !m.sliceEqual(val, payload) {
                        m.valueCacheMutex.Lock()
                        m.valueCache[addressData] = payload
                        m.valueCacheMutex.Unlock()
                        // If this is a new value, we have to also provide the 3 different types of addresses.
                        if !ok {
                            arb := utils.NewReadBuffer(utils.Int8ArrayToUint8Array(destinationAddress))
                            if address, err2 := driverModel.KnxGroupAddressParse(arb, 3); err2 == nil {
                                m.leve3AddressCache[addressData] = driverModel.CastKnxGroupAddress3Level(address)
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
