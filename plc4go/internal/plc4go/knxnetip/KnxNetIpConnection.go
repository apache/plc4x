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
    driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
    internalModel "github.com/apache/plc4x/plc4go/internal/plc4go/model"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi"
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/interceptors"
    "github.com/apache/plc4x/plc4go/internal/plc4go/transports"
    "github.com/apache/plc4x/plc4go/internal/plc4go/transports/udp"
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
    "github.com/apache/plc4x/plc4go/pkg/plc4go"
    apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
    "net"
    "strconv"
    "strings"
    "sync"
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
    quitConnectionStateTimer chan struct{}
    subscribers              []*KnxNetIpSubscriber
    leve3AddressCache        map[uint16]*driverModel.KnxGroupAddress3Level
    leve2AddressCache        map[uint16]*driverModel.KnxGroupAddress2Level
    leve1AddressCache        map[uint16]*driverModel.KnxGroupAddressFreeLevel

    valueCache      map[uint16][]int8
    valueCacheMutex sync.RWMutex
    metadata        *ConnectionMetadata

    GatewayKnxAddress      *driverModel.KnxAddress
    ClientKnxAddress       *driverModel.KnxAddress
    CommunicationChannelId uint8

    requestInterceptor internalModel.RequestInterceptor
    plc4go.PlcConnection
}

func NewKnxNetIpConnection(messageCodec spi.MessageCodec, options map[string][]string, fieldHandler spi.PlcFieldHandler) *KnxNetIpConnection {
    return &KnxNetIpConnection{
        messageCodec:       messageCodec,
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
        err = m.messageCodec.Send(searchRequest)
        if err != nil {
            ch <- plc4go.NewPlcConnectionConnectResult(m, errors.New(
                "error sending search request"))
            return
        }
        // Register an expected response
        check := func(response interface{}) (bool, bool) {
            searchResponse := driverModel.CastSearchResponse(response)
            return searchResponse != nil, false
        }

        // Create a channel for async execution of the connection
        connectionResult := make(chan error)

        // Register a callback to handle the response
        searchResponseChan := m.messageCodec.Expect(check)
        go func() {
            select {
            case response := <-searchResponseChan:
                searchResponse := driverModel.CastSearchResponse(response)
                // Check if this device supports tunneling services
                supportsTunneling := false
                for _, serviceId := range searchResponse.DibSuppSvcFamilies.ServiceIds {
                    _, ok := serviceId.Child.(*driverModel.KnxNetIpTunneling)
                    if ok {
                        supportsTunneling = true
                        break
                    }
                }
                if supportsTunneling {
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
                    for _, serviceId := range searchResponse.DibSuppSvcFamilies.ServiceIds {
                        m.metadata.SupportedServices = append(m.metadata.SupportedServices, serviceId.Child.GetTypeName())
                    }

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
                    err = m.messageCodec.Send(connectionRequest)
                    if err != nil {
                        // TODO: Different channel ...
                        ch <- plc4go.NewPlcConnectionConnectResult(m, errors.New(
                            "error sending connection request"))
                        return
                    }
                    // Register an expected response
                    check := func(response interface{}) (bool, bool) {
                        connectionResponse := driverModel.CastConnectionResponse(response)
                        return connectionResponse != nil, false
                    }

                    // Register a callback to handle the response
                    connectionResponseChan := m.messageCodec.Expect(check)
                    go func() {
                        response := <-connectionResponseChan
                        connectionResponse := driverModel.CastConnectionResponse(response)
                        // Save the communication channel id
                        m.CommunicationChannelId = connectionResponse.CommunicationChannelId
                        if connectionResponse.Status == driverModel.Status_NO_ERROR {
                            // Register a listener for incoming tunneling requests
                            checkTunnelReq := func(response interface{}) (bool, bool) {
                                tunnelingRequest := driverModel.CastTunnelingRequest(response)
                                return (tunnelingRequest != nil) &&
                                        (tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId == m.CommunicationChannelId),
                                    true
                            }
                            tunnelingRequestChan := m.messageCodec.Expect(checkTunnelReq)
                            go func() {
                                m.handleIncomingTunnelingRequest(tunnelingRequestChan)
                            }()

                            tunnelConnectionDataBlock := connectionResponse.ConnectionResponseDataBlock.Child.(*driverModel.ConnectionResponseDataBlockTunnelConnection)
                            // Save the KNX Address the Gateway assigned to this connection.
                            m.ClientKnxAddress = tunnelConnectionDataBlock.KnxAddress

                            // Fire the "connected" event
                            ch <- plc4go.NewPlcConnectionConnectResult(m, nil)

                            // Start a timer that sends connection-state requests every 60 seconds
                            connectionStateTimer := time.NewTicker(60 * time.Second)
                            m.quitConnectionStateTimer = make(chan struct{})
                            go func() {
                                for {
                                    select {
                                    case <-connectionStateTimer.C:
                                        // We're using the connection-state-request as ping operation ...
                                        ping := m.Ping()
                                        pingResult := <-ping
                                        if pingResult.Err != nil {
                                            // TODO: Do some error handling here ...
                                            connectionStateTimer.Stop()
                                        }
                                    case <-m.quitConnectionStateTimer:
                                        // TODO: Do some error handling here ...
                                        connectionStateTimer.Stop()
                                        return
                                    }
                                }
                            }()
                        } else {
                            ch <- plc4go.NewPlcConnectionConnectResult(m,
                                errors.New("got a connection response with status "+strconv.Itoa(int(connectionResponse.Status))))
                        }
                    }()
                } else {
                    ch <- plc4go.NewPlcConnectionConnectResult(m,
                        errors.New("this connection doesn't support tunneling"))
                }
            case <-time.After(1 * time.Second):
                ch <- plc4go.NewPlcConnectionConnectResult(m, errors.New("request timed out"))
            }
        }()
        // Wait for the connection to be established
        err = <-connectionResult
        ch <- plc4go.NewPlcConnectionConnectResult(m, err)
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
        err := m.messageCodec.Send(connectionStateRequest)
        if err != nil {
            result <- plc4go.NewPlcConnectionPingResult(err)
            return
        }
        // Register an expected response
        check := func(response interface{}) (bool, bool) {
            connectionStateResponse := driverModel.CastConnectionStateResponse(response)
            return connectionStateResponse != nil, false
        }

        // Register a callback to handle the response
        connectionStateResponseChan := m.messageCodec.Expect(check)
        go func() {
            response := <-connectionStateResponseChan
            connectionStateResponse := driverModel.CastConnectionStateResponse(response)
            if connectionStateResponse.Status != driverModel.Status_NO_ERROR {
                result <- plc4go.NewPlcConnectionPingResult(errors.New(
                    "got a failure response code " + strconv.Itoa(int(connectionStateResponse.Status))))
            } else {
                result <- plc4go.NewPlcConnectionPingResult(nil)
            }
        }()
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

func (m *KnxNetIpConnection) castIpToKnxAddress(ip net.IP) *driverModel.IPAddress {
    return driverModel.NewIPAddress(utils.ByteArrayToInt8Array(ip)[len(ip)-4:])
}

func (m *KnxNetIpConnection) handleIncomingTunnelingRequest(tunnelingRequestChan chan interface{}) {
    for {
        msg := <-tunnelingRequestChan
        tunnelingRequest := driverModel.CastTunnelingRequest(msg)
        // Send a response for this message
        tunnelingResponse := driverModel.NewTunnelingResponse(driverModel.NewTunnelingResponseDataBlock(
            m.CommunicationChannelId, tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
            driverModel.Status_NO_ERROR))
        err := m.messageCodec.Send(tunnelingResponse)
        if err != nil {
            // TODO: Somehow react on this ...
            break
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
                        for _, subscriber := range m.subscribers {
                            subscriber.handleValueChange(lDataInd.DataFrame, changed)
                        }
                    }
                }
            }
        }()
    }
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
