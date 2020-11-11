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
    "net"
    driverModel "plc4x.apache.org/plc4go/v0/internal/plc4go/knxnetip/readwrite/model"
    internalModel "plc4x.apache.org/plc4go/v0/internal/plc4go/model"
    "plc4x.apache.org/plc4go/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go/v0/internal/plc4go/spi/interceptors"
    "plc4x.apache.org/plc4go/v0/internal/plc4go/transports"
    "plc4x.apache.org/plc4go/v0/internal/plc4go/transports/udp"
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
    "plc4x.apache.org/plc4go/v0/pkg/plc4go"
    apiModel "plc4x.apache.org/plc4go/v0/pkg/plc4go/model"
    "strconv"
    "time"
)

type ConnectionMetadata struct {
    apiModel.PlcConnectionMetadata
}

func (m ConnectionMetadata) CanRead() bool {
    return false
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

    GatewayKnxAddress      *driverModel.KnxAddress
    GatewayName            string
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
        localAddress := driverModel.NewIPAddress(utils.ByteToInt8(udpTransportInstance.LocalAddress.IP))
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
            response := <-searchResponseChan

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
                m.GatewayName = string(bytes.Trim(utils.Int8ToByte(
                    searchResponse.DibDeviceInfo.DeviceFriendlyName), "\x00"))
                m.GatewayKnxAddress = searchResponse.DibDeviceInfo.KnxAddress

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

                        tunnelConnectionDataBlock := connectionResponse.ConnectionResponseDataBlock.Child.(
                        *driverModel.ConnectionResponseDataBlockTunnelConnection)
                        // Save the KNX Address the Gateway assigned to this connection.
                        m.ClientKnxAddress = tunnelConnectionDataBlock.KnxAddress

                        fmt.Printf("Successfully connected to KNXnet/IP Gateway '%s' with KNX address '%d.%d.%d' got assigned client KNX address '%d.%d.%d'",
                            m.GatewayName,
                            m.GatewayKnxAddress.MainGroup, m.GatewayKnxAddress.MiddleGroup, m.GatewayKnxAddress.SubGroup,
                            m.ClientKnxAddress.MainGroup, m.ClientKnxAddress.MiddleGroup, m.ClientKnxAddress.SubGroup)

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
    panic("implement me")
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
    return ConnectionMetadata{}
}

func (m *KnxNetIpConnection) ReadRequestBuilder() apiModel.PlcReadRequestBuilder {
    panic("this connection doesn't support reading")
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
    return driverModel.NewIPAddress(utils.ByteToInt8(ip)[len(ip)-4:])
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

        cemiDataInd := driverModel.CastCEMIDataInd(tunnelingRequest.Cemi.Child)
        if cemiDataInd != nil {
            for _, subscriber := range m.subscribers {
                subscriber.handle(cemiDataInd.CemiDataFrame)
            }
        }
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
