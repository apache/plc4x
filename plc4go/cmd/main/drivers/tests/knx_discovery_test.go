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
package tests

import (
    "encoding/hex"
    "errors"
    "fmt"
    "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
    log "github.com/sirupsen/logrus"
    "net"
    "strconv"
    "strings"
    "testing"
)

func TestKnxAutoDiscovery(t *testing.T) {
    interfaces, err := net.Interfaces()
    if err != nil {
        log.Errorf("Error getting interefaces: %s", err.Error())
        t.Fail()
    }

    for _, interf := range interfaces {
        addrs, err := interf.Addrs()
        if err != nil {
            log.Errorf("Error getting addresses of interface: %s. Got error: %s", interf.Name, err.Error())
            t.Fail()
        }
        for _, addr := range addrs {
            var ipv4Addr net.IP
            switch addr.(type) {
            // If the device is configured to communicate with a subnet
            case *net.IPNet:
                ipv4Addr = addr.(*net.IPNet).IP.To4()

            // If the device is configured for a point-to-point connection
            case *net.IPAddr:
                ipv4Addr = addr.(*net.IPAddr).IP.To4()
            }

            // Only if this is an IPv4 address, will we open a port for it.
            if ipv4Addr != nil {
                // Open a listening port on a random free port number
                udpIpv4Addr := &net.UDPAddr{IP: ipv4Addr, Port: 0}

                udpSocket, err := net.ListenUDP("udp4", udpIpv4Addr)
                if err != nil {
                    log.Warnf("Error creating listening port for KNX on address %s", ipv4Addr.String())
                    continue
                }

                go func() {
                    buf := make([]byte, 1024)

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Send the Search Request using the current network device
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    localIp := udpSocket.LocalAddr().(*net.UDPAddr).IP
                    localPort := udpSocket.LocalAddr().(*net.UDPAddr).Port

                    // Prepare the discovery packet data
                    searchRequestMessage := model.NewSearchRequest(model.NewHPAIDiscoveryEndpoint(
                        model.HostProtocolCode_IPV4_UDP,
                        model.NewIPAddress(utils.ByteArrayToInt8Array(localIp.To4())),
                        uint16(localPort)))
                    writeBuffer := utils.NewWriteBuffer()
                    err := searchRequestMessage.Serialize(*writeBuffer)
                    if err != nil {
                        panic("Failed preparing search request.")
                    }

                    // This is the multicast address and port KNX devices are supposed to listen to.
                    destination := &net.UDPAddr{IP: net.IPv4(224, 0, 23, 12), Port: 3671}

                    // Send the message
                    _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), destination)
                    if err != nil {
                        panic("Failed sending search request.")
                    }

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Wait for a response from a device that supports tunneling
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    var searchResponse *model.SearchResponse
                    var gatewayAddr *net.UDPAddr
                    for searchResponse == nil {
                        // Read a new packet from the socket
                        _, src, err := udpSocket.ReadFromUDP(buf)
                        if err != nil {
                            panic("Error reading from KNX UDP socket")
                        }

                        readBuffer := utils.NewReadBuffer(buf)
                        knxMessage, err := model.KnxNetIpMessageParse(readBuffer)
                        if err != nil {
                            hexEncodedPayload := hex.EncodeToString(buf)
                            panic(fmt.Sprintf("Error decoding incoming KNX message from %v with payload %s", src, hexEncodedPayload))
                        }

                        switch knxMessage.Child.(type) {
                        // If this is a search response and the current device supports tunneling,
                        // we've found what we're looking for.
                        case *model.SearchResponse:
                            curSearchResponse := model.CastSearchResponse(knxMessage)
                            for _, serviceId := range curSearchResponse.DibSuppSvcFamilies.ServiceIds {
                                switch (*serviceId).Child.(type) {
                                case *model.KnxNetIpTunneling:
                                    searchResponse = curSearchResponse
                                    gatewayAddr = src
                                    break
                                }
                            }

                        // Just ACK any incoming tunneling requests
                        case *model.TunnelingRequest:
                            tunnelingRequest := model.CastTunnelingRequest(knxMessage)

                            tunnelingResponse := model.NewTunnelingResponse(
                                model.NewTunnelingResponseDataBlock(
                                    tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId,
                                    tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
                                    model.Status_NO_ERROR))
                            writeBuffer := utils.NewWriteBuffer()
                            err = tunnelingResponse.Serialize(*writeBuffer)
                            if err != nil {
                                panic("Failed preparing tunneling response.")
                            }

                            // Send the message
                            _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), src)
                            if err != nil {
                                panic("Failed sending tunneling response.")
                            }
                        }
                    }

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Send a connection request to the device we just found
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    connectionRequestMessage := model.NewConnectionRequest(
                        model.NewHPAIDiscoveryEndpoint(
                            model.HostProtocolCode_IPV4_UDP,
                            model.NewIPAddress(utils.ByteArrayToInt8Array(localIp.To4())),
                            uint16(localPort)),
                        model.NewHPAIDataEndpoint(
                            model.HostProtocolCode_IPV4_UDP,
                            model.NewIPAddress(utils.ByteArrayToInt8Array(localIp.To4())),
                            uint16(localPort)),
                        model.NewConnectionRequestInformationTunnelConnection(model.KnxLayer_TUNNEL_LINK_LAYER))
                    writeBuffer = utils.NewWriteBuffer()
                    err = connectionRequestMessage.Serialize(*writeBuffer)
                    if err != nil {
                        panic("Failed preparing connection request.")
                    }

                    // Send the message
                    _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), gatewayAddr)

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Wait for a connection response
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    var communicationChannelId uint8
                    for communicationChannelId == 0 {
                        // Read a new packet from the socket
                        _, src, err := udpSocket.ReadFromUDP(buf)
                        if err != nil {
                            panic("Error reading from KNX UDP socket")
                        }

                        readBuffer := utils.NewReadBuffer(buf)
                        knxMessage, err := model.KnxNetIpMessageParse(readBuffer)
                        if err != nil {
                            hexEncodedPayload := hex.EncodeToString(buf)
                            panic(fmt.Sprintf("Error decoding incoming KNX message from %v with payload %s", src, hexEncodedPayload))
                        }

                        switch knxMessage.Child.(type) {
                        case *model.ConnectionResponse:
                            connectionResponse := model.CastConnectionResponse(knxMessage)
                            if connectionResponse.Status == model.Status_NO_ERROR {
                                communicationChannelId = connectionResponse.CommunicationChannelId
                            } else {
                                panic("Got an error while connecting")
                            }

                        // Just ACK any incoming tunneling requests
                        case *model.TunnelingRequest:
                            tunnelingRequest := model.CastTunnelingRequest(knxMessage)

                            tunnelingResponse := model.NewTunnelingResponse(
                                model.NewTunnelingResponseDataBlock(
                                    tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId,
                                    tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
                                    model.Status_NO_ERROR))
                            writeBuffer := utils.NewWriteBuffer()
                            err = tunnelingResponse.Serialize(*writeBuffer)
                            if err != nil {
                                panic("Failed preparing tunneling response.")
                            }

                            // Send the message
                            _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), src)
                            if err != nil {
                                panic("Failed sending tunneling response.")
                            }
                        }
                    }

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Send a config connection request
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    /*configConnectionRequestMessage := model.NewConnectionRequest(
                        model.NewHPAIDiscoveryEndpoint(
                            model.HostProtocolCode_IPV4_UDP,
                            model.NewIPAddress(utils.ByteArrayToInt8Array(localIp.To4())),
                            uint16(localPort)),
                        model.NewHPAIDataEndpoint(
                            model.HostProtocolCode_IPV4_UDP,
                            model.NewIPAddress(utils.ByteArrayToInt8Array(localIp.To4())),
                            uint16(localPort)),
                        model.NewConnectionRequestInformationDeviceManagement())
                    writeBuffer = utils.NewWriteBuffer()
                    err = configConnectionRequestMessage.Serialize(*writeBuffer)
                    if err != nil {
                        panic("Failed preparing connection request.")
                    }

                    // Send the message
                    _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), gatewayAddr)

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Wait for a config connection response
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    var configCommunicationChannelId uint8
                    for configCommunicationChannelId == 0 {
                        // Read a new packet from the socket
                        _, src, err := udpSocket.ReadFromUDP(buf)
                        if err != nil {
                            panic("Error reading from KNX UDP socket")
                        }

                        readBuffer := utils.NewReadBuffer(buf)
                        knxMessage, err := model.KnxNetIpMessageParse(readBuffer)
                        if err != nil {
                            hexEncodedPayload := hex.EncodeToString(buf)
                            panic(fmt.Sprintf("Error decoding incoming KNX message from %v with payload %s", src, hexEncodedPayload))
                        }

                        switch knxMessage.Child.(type) {
                        case *model.ConnectionResponse:
                            connectionResponse := model.CastConnectionResponse(knxMessage)
                            if connectionResponse.Status == model.Status_NO_ERROR {
                                configCommunicationChannelId = connectionResponse.CommunicationChannelId
                            } else {
                                panic("Got an error while connecting")
                            }

                        // Just ACK any incoming tunneling requests
                        case *model.TunnelingRequest:
                            tunnelingRequest := model.CastTunnelingRequest(knxMessage)

                            tunnelingResponse := model.NewTunnelingResponse(
                                model.NewTunnelingResponseDataBlock(
                                    tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId,
                                    tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
                                    model.Status_NO_ERROR))
                            writeBuffer := utils.NewWriteBuffer()
                            err = tunnelingResponse.Serialize(*writeBuffer)
                            if err != nil {
                                panic("Failed preparing tunneling response.")
                            }

                            // Send the message
                            _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), src)
                            if err != nil {
                                panic("Failed sending tunneling response.")
                            }
                        }
                    }

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Send a prop-read request for PID_MEDIUM_TYPE
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    mediumTypePropRequestMessage := model.NewDeviceConfigurationRequest(
                        model.NewDeviceConfigurationRequestDataBlock(configCommunicationChannelId, uint8(0)),
                        model.NewMPropReadReq(uint16(8), uint8(1), uint8(51), uint8(1), uint16(1)))
                    writeBuffer = utils.NewWriteBuffer()
                    err = mediumTypePropRequestMessage.Serialize(*writeBuffer)
                    if err != nil {
                        panic("Failed preparing config disconnect request.")
                    }

                    // Send the message
                    _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), gatewayAddr)

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Wait for a prop read response
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    var mediumType uint16
                    for mediumType == 0 {
                        // Read a new packet from the socket
                        _, src, err := udpSocket.ReadFromUDP(buf)
                        if err != nil {
                            panic("Error reading from KNX UDP socket")
                        }

                        readBuffer := utils.NewReadBuffer(buf)
                        knxMessage, err := model.KnxNetIpMessageParse(readBuffer)
                        if err != nil {
                            hexEncodedPayload := hex.EncodeToString(buf)
                            panic(fmt.Sprintf("Error decoding incoming KNX message from %v with payload %s", src, hexEncodedPayload))
                        }

                        switch knxMessage.Child.(type) {
                        case *model.DeviceConfigurationRequest:
                            deviceConfigurationRequest := model.CastDeviceConfigurationRequest(knxMessage)
                            switch deviceConfigurationRequest.Cemi.Child.(type) {
                            case *model.MPropReadCon:
                                readCon := model.CastMPropReadCon(deviceConfigurationRequest.Cemi)
                                // TODO: This should be renamed to "Data"
                                mediumType = readCon.Unknown

                                // Send and ACK for this response
                                tunnelingResponse := model.NewTunnelingResponse(
                                    model.NewTunnelingResponseDataBlock(
                                        deviceConfigurationRequest.DeviceConfigurationRequestDataBlock.CommunicationChannelId,
                                        deviceConfigurationRequest.DeviceConfigurationRequestDataBlock.SequenceCounter,
                                        model.Status_NO_ERROR))
                                writeBuffer := utils.NewWriteBuffer()
                                err = tunnelingResponse.Serialize(*writeBuffer)
                                if err != nil {
                                    panic("Failed preparing tunneling response.")
                                }

                                // Send the message
                                _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), src)
                                if err != nil {
                                    panic("Failed sending tunneling response.")
                                }
                            }

                        // Just ACK any incoming tunneling requests
                        case *model.TunnelingRequest:
                            tunnelingRequest := model.CastTunnelingRequest(knxMessage)

                            tunnelingResponse := model.NewTunnelingResponse(
                                model.NewTunnelingResponseDataBlock(
                                    tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId,
                                    tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
                                    model.Status_NO_ERROR))
                            writeBuffer := utils.NewWriteBuffer()
                            err = tunnelingResponse.Serialize(*writeBuffer)
                            if err != nil {
                                panic("Failed preparing tunneling response.")
                            }

                            // Send the message
                            _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), src)
                            if err != nil {
                                panic("Failed sending tunneling response.")
                            }
                        }
                    }

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Send a prop-read request for PID_MAX_APDULENGTH
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    maxApdulengthPropRequestMessage := model.NewDeviceConfigurationRequest(
                        model.NewDeviceConfigurationRequestDataBlock(configCommunicationChannelId, uint8(1)),
                        model.NewMPropReadReq(uint16(0), uint8(1), uint8(56), uint8(1), uint16(1)))
                    writeBuffer = utils.NewWriteBuffer()
                    err = maxApdulengthPropRequestMessage.Serialize(*writeBuffer)
                    if err != nil {
                        panic("Failed preparing config disconnect request.")
                    }

                    // Send the message
                    _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), gatewayAddr)

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Wait for a prop read response
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    var maxApdulength uint16
                    for maxApdulength == 0 {
                        // Read a new packet from the socket
                        _, src, err := udpSocket.ReadFromUDP(buf)
                        if err != nil {
                            panic("Error reading from KNX UDP socket")
                        }

                        readBuffer := utils.NewReadBuffer(buf)
                        knxMessage, err := model.KnxNetIpMessageParse(readBuffer)
                        if err != nil {
                            hexEncodedPayload := hex.EncodeToString(buf)
                            panic(fmt.Sprintf("Error decoding incoming KNX message from %v with payload %s", src, hexEncodedPayload))
                        }

                        switch knxMessage.Child.(type) {
                        case *model.DeviceConfigurationRequest:
                            deviceConfigurationRequest := model.CastDeviceConfigurationRequest(knxMessage)
                            switch deviceConfigurationRequest.Cemi.Child.(type) {
                            case *model.MPropReadCon:
                                readCon := model.CastMPropReadCon(deviceConfigurationRequest.Cemi)
                                // TODO: This should be renamed to "Data"
                                maxApdulength = readCon.Unknown

                                // Send and ACK for this response
                                tunnelingResponse := model.NewTunnelingResponse(
                                    model.NewTunnelingResponseDataBlock(
                                        deviceConfigurationRequest.DeviceConfigurationRequestDataBlock.CommunicationChannelId,
                                        deviceConfigurationRequest.DeviceConfigurationRequestDataBlock.SequenceCounter,
                                        model.Status_NO_ERROR))
                                writeBuffer := utils.NewWriteBuffer()
                                err = tunnelingResponse.Serialize(*writeBuffer)
                                if err != nil {
                                    panic("Failed preparing tunneling response.")
                                }

                                // Send the message
                                _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), src)
                                if err != nil {
                                    panic("Failed sending tunneling response.")
                                }
                            }

                        // Just ACK any incoming tunneling requests
                        case *model.TunnelingRequest:
                            tunnelingRequest := model.CastTunnelingRequest(knxMessage)

                            tunnelingResponse := model.NewTunnelingResponse(
                                model.NewTunnelingResponseDataBlock(
                                    tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId,
                                    tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
                                    model.Status_NO_ERROR))
                            writeBuffer := utils.NewWriteBuffer()
                            err = tunnelingResponse.Serialize(*writeBuffer)
                            if err != nil {
                                panic("Failed preparing tunneling response.")
                            }

                            // Send the message
                            _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), src)
                            if err != nil {
                                panic("Failed sending tunneling response.")
                            }
                        }
                    }

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Send a config connection disconnect request
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    configConnectionDisconnectRequestMessage := model.NewDisconnectRequest(
                        configCommunicationChannelId,
                        model.NewHPAIControlEndpoint(
                            model.HostProtocolCode_IPV4_UDP,
                            model.NewIPAddress(utils.ByteArrayToInt8Array(localIp.To4())),
                            uint16(localPort)))
                    writeBuffer = utils.NewWriteBuffer()
                    err = configConnectionDisconnectRequestMessage.Serialize(*writeBuffer)
                    if err != nil {
                        panic("Failed preparing config disconnect request.")
                    }

                    // Send the message
                    _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), gatewayAddr)

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Wait for a config disconnect response
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    configDisconnected := false
                    for !configDisconnected {
                        // Read a new packet from the socket
                        _, src, err := udpSocket.ReadFromUDP(buf)
                        if err != nil {
                            panic("Error reading from KNX UDP socket")
                        }

                        readBuffer := utils.NewReadBuffer(buf)
                        knxMessage, err := model.KnxNetIpMessageParse(readBuffer)
                        if err != nil {
                            hexEncodedPayload := hex.EncodeToString(buf)
                            panic(fmt.Sprintf("Error decoding incoming KNX message from %v with payload %s", src, hexEncodedPayload))
                        }

                        switch knxMessage.Child.(type) {
                        case *model.DisconnectResponse:
                            disconnectResponse := model.CastDisconnectResponse(knxMessage)
                            if disconnectResponse.Status == model.Status_NO_ERROR {
                                configDisconnected = true
                            } else {
                                panic("Got an error while connecting")
                            }

                        // Just ACK any incoming tunneling requests
                        case *model.TunnelingRequest:
                            tunnelingRequest := model.CastTunnelingRequest(knxMessage)

                            tunnelingResponse := model.NewTunnelingResponse(
                                model.NewTunnelingResponseDataBlock(
                                    tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId,
                                    tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
                                    model.Status_NO_ERROR))
                            writeBuffer := utils.NewWriteBuffer()
                            err = tunnelingResponse.Serialize(*writeBuffer)
                            if err != nil {
                                panic("Failed preparing tunneling response.")
                            }

                            // Send the message
                            _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), src)
                            if err != nil {
                                panic("Failed sending tunneling response.")
                            }
                        }
                    }*/

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Send a device connection request to KNX address 1.1.10
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    targetKnxAddress, err := ParseKnxAddressString("1.1.10")
                    if err != nil {
                        panic("Failed preparing discovery request.")
                    }
                    sourceAddress := &model.KnxAddress{
                        MainGroup:   0,
                        MiddleGroup: 0,
                        SubGroup:    0,
                    }
                    controlType := model.ControlType_CONNECT
                    deviceConnectionRequest := model.NewTunnelingRequest(
                        model.NewTunnelingRequestDataBlock(communicationChannelId, 0),
                        model.NewLDataReq(0, nil,
                            model.NewLDataFrameDataExt(false, 6, uint8(0),
                                sourceAddress, targetKnxAddress, uint8(0), true, false,
                                uint8(0), &controlType, nil, nil, nil, nil,
                                false, model.CEMIPriority_SYSTEM, false, false)))
                    writeBuffer = utils.NewWriteBuffer()
                    err = deviceConnectionRequest.Serialize(*writeBuffer)
                    if err != nil {
                        panic("Failed preparing device connection request.")
                    }

                    // Send the message
                    _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), gatewayAddr)
                    if err != nil {
                        panic("Failed sending device connection request.")
                    }

                    ////////////////////////////////////////////////////////////////////////////////////////////////////
                    // Wait for a device connection response
                    ////////////////////////////////////////////////////////////////////////////////////////////////////

                    for {
                        // Read a new packet from the socket
                        _, src, err := udpSocket.ReadFromUDP(buf)
                        if err != nil {
                            panic("Error reading from KNX UDP socket")
                        }

                        readBuffer := utils.NewReadBuffer(buf)
                        knxMessage, err := model.KnxNetIpMessageParse(readBuffer)
                        if err != nil {
                            hexEncodedPayload := hex.EncodeToString(buf)
                            panic(fmt.Sprintf("Error decoding incoming KNX message from %v with payload %s", src, hexEncodedPayload))
                        }

                        switch knxMessage.Child.(type) {
                        case *model.TunnelingRequest:
                            tunnelingRequest := model.CastTunnelingRequest(knxMessage)

                            tunnelingResponse := model.NewTunnelingResponse(
                                model.NewTunnelingResponseDataBlock(
                                    tunnelingRequest.TunnelingRequestDataBlock.CommunicationChannelId,
                                    tunnelingRequest.TunnelingRequestDataBlock.SequenceCounter,
                                    model.Status_NO_ERROR))
                            writeBuffer := utils.NewWriteBuffer()
                            err = tunnelingResponse.Serialize(*writeBuffer)
                            if err != nil {
                                panic("Failed preparing tunneling response.")
                            }

                            // Send the message
                            _, err = udpSocket.WriteTo(writeBuffer.GetBytes(), src)
                            if err != nil {
                                panic("Failed sending tunneling response.")
                            }
                        }
                    }
                }()
            }
        }
    }
}

func ParseKnxAddressString(knxAddressString string) ([]int8, error) {
    if strings.Count(knxAddressString, ".") != 2 {
        return nil, errors.New("this is not a valid knx address")
    }

    split := strings.Split(knxAddressString, ".")

    mainSegment, err := strconv.Atoi(split[0])
    if err != nil {
        return nil, errors.New("this is not a valid knx address")
    }
    if mainSegment < 0 || mainSegment > 15 {
        return nil, errors.New("this is not a valid knx address")
    }

    middleSegment, err := strconv.Atoi(split[1])
    if err != nil {
        return nil, errors.New("this is not a valid knx address")
    }
    if middleSegment < 0 || middleSegment > 15 {
        return nil, errors.New("this is not a valid knx address")
    }

    subSegment, err := strconv.Atoi(split[2])
    if err != nil {
        return nil, errors.New("this is not a valid knx address")
    }
    if subSegment < 0 || subSegment > 255 {
        return nil, errors.New("this is not a valid knx address")
    }

    return []int8{
        int8(mainSegment << 4 | middleSegment),
        int8(subSegment),
    }, nil
}
