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
package model

import (
    "encoding/xml"
    "errors"
    "io"
    log "github.com/sirupsen/logrus"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    "reflect"
)

// The data-structure of this message
type NPDU struct {
    ProtocolVersionNumber uint8
    MessageTypeFieldPresent bool
    DestinationSpecified bool
    SourceSpecified bool
    ExpectingReply bool
    NetworkPriority uint8
    DestinationNetworkAddress *uint16
    DestinationLength *uint8
    DestinationAddress []uint8
    SourceNetworkAddress *uint16
    SourceLength *uint8
    SourceAddress []uint8
    HopCount *uint8
    Nlm *INLM
    Apdu *IAPDU

}

// The corresponding interface
type INPDU interface {
    spi.Message
    Serialize(io utils.WriteBuffer) error
}


func NewNPDU(protocolVersionNumber uint8, messageTypeFieldPresent bool, destinationSpecified bool, sourceSpecified bool, expectingReply bool, networkPriority uint8, destinationNetworkAddress *uint16, destinationLength *uint8, destinationAddress []uint8, sourceNetworkAddress *uint16, sourceLength *uint8, sourceAddress []uint8, hopCount *uint8, nlm *INLM, apdu *IAPDU) spi.Message {
    return &NPDU{ProtocolVersionNumber: protocolVersionNumber, MessageTypeFieldPresent: messageTypeFieldPresent, DestinationSpecified: destinationSpecified, SourceSpecified: sourceSpecified, ExpectingReply: expectingReply, NetworkPriority: networkPriority, DestinationNetworkAddress: destinationNetworkAddress, DestinationLength: destinationLength, DestinationAddress: destinationAddress, SourceNetworkAddress: sourceNetworkAddress, SourceLength: sourceLength, SourceAddress: sourceAddress, HopCount: hopCount, Nlm: nlm, Apdu: apdu}
}

func CastINPDU(structType interface{}) INPDU {
    castFunc := func(typ interface{}) INPDU {
        if iNPDU, ok := typ.(INPDU); ok {
            return iNPDU
        }
        return nil
    }
    return castFunc(structType)
}

func CastNPDU(structType interface{}) NPDU {
    castFunc := func(typ interface{}) NPDU {
        if sNPDU, ok := typ.(NPDU); ok {
            return sNPDU
        }
        if sNPDU, ok := typ.(*NPDU); ok {
            return *sNPDU
        }
        return NPDU{}
    }
    return castFunc(structType)
}

func (m NPDU) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Simple field (protocolVersionNumber)
    lengthInBits += 8

    // Simple field (messageTypeFieldPresent)
    lengthInBits += 1

    // Reserved Field (reserved)
    lengthInBits += 1

    // Simple field (destinationSpecified)
    lengthInBits += 1

    // Reserved Field (reserved)
    lengthInBits += 1

    // Simple field (sourceSpecified)
    lengthInBits += 1

    // Simple field (expectingReply)
    lengthInBits += 1

    // Simple field (networkPriority)
    lengthInBits += 2

    // Optional Field (destinationNetworkAddress)
    if m.DestinationNetworkAddress != nil {
        lengthInBits += 16
    }

    // Optional Field (destinationLength)
    if m.DestinationLength != nil {
        lengthInBits += 8
    }

    // Array field
    if len(m.DestinationAddress) > 0 {
        lengthInBits += 8 * uint16(len(m.DestinationAddress))
    }

    // Optional Field (sourceNetworkAddress)
    if m.SourceNetworkAddress != nil {
        lengthInBits += 16
    }

    // Optional Field (sourceLength)
    if m.SourceLength != nil {
        lengthInBits += 8
    }

    // Array field
    if len(m.SourceAddress) > 0 {
        lengthInBits += 8 * uint16(len(m.SourceAddress))
    }

    // Optional Field (hopCount)
    if m.HopCount != nil {
        lengthInBits += 8
    }

    // Optional Field (nlm)
    if m.Nlm != nil {
        lengthInBits += (*m.Nlm).LengthInBits()
    }

    // Optional Field (apdu)
    if m.Apdu != nil {
        lengthInBits += (*m.Apdu).LengthInBits()
    }

    return lengthInBits
}

func (m NPDU) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func NPDUParse(io *utils.ReadBuffer, npduLength uint16) (spi.Message, error) {

    // Simple Field (protocolVersionNumber)
    protocolVersionNumber, _protocolVersionNumberErr := io.ReadUint8(8)
    if _protocolVersionNumberErr != nil {
        return nil, errors.New("Error parsing 'protocolVersionNumber' field " + _protocolVersionNumberErr.Error())
    }

    // Simple Field (messageTypeFieldPresent)
    messageTypeFieldPresent, _messageTypeFieldPresentErr := io.ReadBit()
    if _messageTypeFieldPresentErr != nil {
        return nil, errors.New("Error parsing 'messageTypeFieldPresent' field " + _messageTypeFieldPresentErr.Error())
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(1)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint8(0) {
            log.WithFields(log.Fields{
                "expected value": uint8(0),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (destinationSpecified)
    destinationSpecified, _destinationSpecifiedErr := io.ReadBit()
    if _destinationSpecifiedErr != nil {
        return nil, errors.New("Error parsing 'destinationSpecified' field " + _destinationSpecifiedErr.Error())
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(1)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint8(0) {
            log.WithFields(log.Fields{
                "expected value": uint8(0),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Simple Field (sourceSpecified)
    sourceSpecified, _sourceSpecifiedErr := io.ReadBit()
    if _sourceSpecifiedErr != nil {
        return nil, errors.New("Error parsing 'sourceSpecified' field " + _sourceSpecifiedErr.Error())
    }

    // Simple Field (expectingReply)
    expectingReply, _expectingReplyErr := io.ReadBit()
    if _expectingReplyErr != nil {
        return nil, errors.New("Error parsing 'expectingReply' field " + _expectingReplyErr.Error())
    }

    // Simple Field (networkPriority)
    networkPriority, _networkPriorityErr := io.ReadUint8(2)
    if _networkPriorityErr != nil {
        return nil, errors.New("Error parsing 'networkPriority' field " + _networkPriorityErr.Error())
    }

    // Optional Field (destinationNetworkAddress) (Can be skipped, if a given expression evaluates to false)
    var destinationNetworkAddress *uint16 = nil
    if destinationSpecified {
        _val, _err := io.ReadUint16(16)
        if _err != nil {
            return nil, errors.New("Error parsing 'destinationNetworkAddress' field " + _err.Error())
        }

        destinationNetworkAddress = &_val
    }

    // Optional Field (destinationLength) (Can be skipped, if a given expression evaluates to false)
    var destinationLength *uint8 = nil
    if destinationSpecified {
        _val, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'destinationLength' field " + _err.Error())
        }

        destinationLength = &_val
    }

    // Array field (destinationAddress)
    // Count array
    destinationAddress := make([]uint8, utils.InlineIf(destinationSpecified, uint16((*destinationLength)), uint16(uint16(0))))
    for curItem := uint16(0); curItem < uint16(utils.InlineIf(destinationSpecified, uint16((*destinationLength)), uint16(uint16(0)))); curItem++ {

        _item, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'destinationAddress' field " + _err.Error())
        }
        destinationAddress[curItem] = _item
    }

    // Optional Field (sourceNetworkAddress) (Can be skipped, if a given expression evaluates to false)
    var sourceNetworkAddress *uint16 = nil
    if sourceSpecified {
        _val, _err := io.ReadUint16(16)
        if _err != nil {
            return nil, errors.New("Error parsing 'sourceNetworkAddress' field " + _err.Error())
        }

        sourceNetworkAddress = &_val
    }

    // Optional Field (sourceLength) (Can be skipped, if a given expression evaluates to false)
    var sourceLength *uint8 = nil
    if sourceSpecified {
        _val, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'sourceLength' field " + _err.Error())
        }

        sourceLength = &_val
    }

    // Array field (sourceAddress)
    // Count array
    sourceAddress := make([]uint8, utils.InlineIf(sourceSpecified, uint16((*sourceLength)), uint16(uint16(0))))
    for curItem := uint16(0); curItem < uint16(utils.InlineIf(sourceSpecified, uint16((*sourceLength)), uint16(uint16(0)))); curItem++ {

        _item, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'sourceAddress' field " + _err.Error())
        }
        sourceAddress[curItem] = _item
    }

    // Optional Field (hopCount) (Can be skipped, if a given expression evaluates to false)
    var hopCount *uint8 = nil
    if destinationSpecified {
        _val, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'hopCount' field " + _err.Error())
        }

        hopCount = &_val
    }

    // Optional Field (nlm) (Can be skipped, if a given expression evaluates to false)
    var nlm *INLM = nil
    if messageTypeFieldPresent {
        _message, _err := NLMParse(io, uint16(npduLength) - uint16(uint16(uint16(uint16(uint16(uint16(2)) + uint16(uint16(utils.InlineIf(sourceSpecified, uint16(uint16(uint16(3)) + uint16((*sourceLength))), uint16(uint16(0)))))) + uint16(uint16(utils.InlineIf(destinationSpecified, uint16(uint16(uint16(3)) + uint16((*destinationLength))), uint16(uint16(0)))))) + uint16(uint16(utils.InlineIf(bool(bool(destinationSpecified) || bool(sourceSpecified)), uint16(uint16(1)), uint16(uint16(0))))))))
        if _err != nil {
            return nil, errors.New("Error parsing 'nlm' field " + _err.Error())
        }
        var _item INLM
        _item, _ok := _message.(INLM)
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to INLM")
        }
        nlm = &_item
    }

    // Optional Field (apdu) (Can be skipped, if a given expression evaluates to false)
    var apdu *IAPDU = nil
    if !(messageTypeFieldPresent) {
        _message, _err := APDUParse(io, uint16(npduLength) - uint16(uint16(uint16(uint16(uint16(uint16(2)) + uint16(uint16(utils.InlineIf(sourceSpecified, uint16(uint16(uint16(3)) + uint16((*sourceLength))), uint16(uint16(0)))))) + uint16(uint16(utils.InlineIf(destinationSpecified, uint16(uint16(uint16(3)) + uint16((*destinationLength))), uint16(uint16(0)))))) + uint16(uint16(utils.InlineIf(bool(bool(destinationSpecified) || bool(sourceSpecified)), uint16(uint16(1)), uint16(uint16(0))))))))
        if _err != nil {
            return nil, errors.New("Error parsing 'apdu' field " + _err.Error())
        }
        var _item IAPDU
        _item, _ok := _message.(IAPDU)
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to IAPDU")
        }
        apdu = &_item
    }

    // Create the instance
    return NewNPDU(protocolVersionNumber, messageTypeFieldPresent, destinationSpecified, sourceSpecified, expectingReply, networkPriority, destinationNetworkAddress, destinationLength, destinationAddress, sourceNetworkAddress, sourceLength, sourceAddress, hopCount, nlm, apdu), nil
}

func (m NPDU) Serialize(io utils.WriteBuffer) error {

    // Simple Field (protocolVersionNumber)
    protocolVersionNumber := uint8(m.ProtocolVersionNumber)
    _protocolVersionNumberErr := io.WriteUint8(8, (protocolVersionNumber))
    if _protocolVersionNumberErr != nil {
        return errors.New("Error serializing 'protocolVersionNumber' field " + _protocolVersionNumberErr.Error())
    }

    // Simple Field (messageTypeFieldPresent)
    messageTypeFieldPresent := bool(m.MessageTypeFieldPresent)
    _messageTypeFieldPresentErr := io.WriteBit((messageTypeFieldPresent))
    if _messageTypeFieldPresentErr != nil {
        return errors.New("Error serializing 'messageTypeFieldPresent' field " + _messageTypeFieldPresentErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(1, uint8(0))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (destinationSpecified)
    destinationSpecified := bool(m.DestinationSpecified)
    _destinationSpecifiedErr := io.WriteBit((destinationSpecified))
    if _destinationSpecifiedErr != nil {
        return errors.New("Error serializing 'destinationSpecified' field " + _destinationSpecifiedErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(1, uint8(0))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Simple Field (sourceSpecified)
    sourceSpecified := bool(m.SourceSpecified)
    _sourceSpecifiedErr := io.WriteBit((sourceSpecified))
    if _sourceSpecifiedErr != nil {
        return errors.New("Error serializing 'sourceSpecified' field " + _sourceSpecifiedErr.Error())
    }

    // Simple Field (expectingReply)
    expectingReply := bool(m.ExpectingReply)
    _expectingReplyErr := io.WriteBit((expectingReply))
    if _expectingReplyErr != nil {
        return errors.New("Error serializing 'expectingReply' field " + _expectingReplyErr.Error())
    }

    // Simple Field (networkPriority)
    networkPriority := uint8(m.NetworkPriority)
    _networkPriorityErr := io.WriteUint8(2, (networkPriority))
    if _networkPriorityErr != nil {
        return errors.New("Error serializing 'networkPriority' field " + _networkPriorityErr.Error())
    }

    // Optional Field (destinationNetworkAddress) (Can be skipped, if the value is null)
    var destinationNetworkAddress *uint16 = nil
    if m.DestinationNetworkAddress != nil {
        destinationNetworkAddress = m.DestinationNetworkAddress
        _destinationNetworkAddressErr := io.WriteUint16(16, *(destinationNetworkAddress))
        if _destinationNetworkAddressErr != nil {
            return errors.New("Error serializing 'destinationNetworkAddress' field " + _destinationNetworkAddressErr.Error())
        }
    }

    // Optional Field (destinationLength) (Can be skipped, if the value is null)
    var destinationLength *uint8 = nil
    if m.DestinationLength != nil {
        destinationLength = m.DestinationLength
        _destinationLengthErr := io.WriteUint8(8, *(destinationLength))
        if _destinationLengthErr != nil {
            return errors.New("Error serializing 'destinationLength' field " + _destinationLengthErr.Error())
        }
    }

    // Array Field (destinationAddress)
    if m.DestinationAddress != nil {
        for _, _element := range m.DestinationAddress {
            _elementErr := io.WriteUint8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'destinationAddress' field " + _elementErr.Error())
            }
        }
    }

    // Optional Field (sourceNetworkAddress) (Can be skipped, if the value is null)
    var sourceNetworkAddress *uint16 = nil
    if m.SourceNetworkAddress != nil {
        sourceNetworkAddress = m.SourceNetworkAddress
        _sourceNetworkAddressErr := io.WriteUint16(16, *(sourceNetworkAddress))
        if _sourceNetworkAddressErr != nil {
            return errors.New("Error serializing 'sourceNetworkAddress' field " + _sourceNetworkAddressErr.Error())
        }
    }

    // Optional Field (sourceLength) (Can be skipped, if the value is null)
    var sourceLength *uint8 = nil
    if m.SourceLength != nil {
        sourceLength = m.SourceLength
        _sourceLengthErr := io.WriteUint8(8, *(sourceLength))
        if _sourceLengthErr != nil {
            return errors.New("Error serializing 'sourceLength' field " + _sourceLengthErr.Error())
        }
    }

    // Array Field (sourceAddress)
    if m.SourceAddress != nil {
        for _, _element := range m.SourceAddress {
            _elementErr := io.WriteUint8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'sourceAddress' field " + _elementErr.Error())
            }
        }
    }

    // Optional Field (hopCount) (Can be skipped, if the value is null)
    var hopCount *uint8 = nil
    if m.HopCount != nil {
        hopCount = m.HopCount
        _hopCountErr := io.WriteUint8(8, *(hopCount))
        if _hopCountErr != nil {
            return errors.New("Error serializing 'hopCount' field " + _hopCountErr.Error())
        }
    }

    // Optional Field (nlm) (Can be skipped, if the value is null)
    var nlm *INLM = nil
    if m.Nlm != nil {
        nlm = m.Nlm
        _nlmErr := CastINLM(*nlm).Serialize(io)
        if _nlmErr != nil {
            return errors.New("Error serializing 'nlm' field " + _nlmErr.Error())
        }
    }

    // Optional Field (apdu) (Can be skipped, if the value is null)
    var apdu *IAPDU = nil
    if m.Apdu != nil {
        apdu = m.Apdu
        _apduErr := CastIAPDU(*apdu).Serialize(io)
        if _apduErr != nil {
            return errors.New("Error serializing 'apdu' field " + _apduErr.Error())
        }
    }

    return nil
}

func (m *NPDU) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        token, err := d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "protocolVersionNumber":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ProtocolVersionNumber = data
            case "messageTypeFieldPresent":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MessageTypeFieldPresent = data
            case "destinationSpecified":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DestinationSpecified = data
            case "sourceSpecified":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SourceSpecified = data
            case "expectingReply":
                var data bool
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ExpectingReply = data
            case "networkPriority":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.NetworkPriority = data
            case "destinationNetworkAddress":
                var data *uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DestinationNetworkAddress = data
            case "destinationLength":
                var data *uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DestinationLength = data
            case "destinationAddress":
                var data []uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DestinationAddress = data
            case "sourceNetworkAddress":
                var data *uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SourceNetworkAddress = data
            case "sourceLength":
                var data *uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SourceLength = data
            case "sourceAddress":
                var data []uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SourceAddress = data
            case "hopCount":
                var data *uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.HopCount = data
            case "nlm":
                switch tok.Attr[0].Value {
                    case "org.apache.plc4x.java.bacnetip.readwrite.NLMWhoIsRouterToNetwork":
                        var dt *NLMWhoIsRouterToNetwork
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Nlm = dt
                    case "org.apache.plc4x.java.bacnetip.readwrite.NLMIAmRouterToNetwork":
                        var dt *NLMIAmRouterToNetwork
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Nlm = dt
                    }
            case "apdu":
                switch tok.Attr[0].Value {
                    case "org.apache.plc4x.java.bacnetip.readwrite.APDUConfirmedRequest":
                        var dt *APDUConfirmedRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Apdu = dt
                    case "org.apache.plc4x.java.bacnetip.readwrite.APDUUnconfirmedRequest":
                        var dt *APDUUnconfirmedRequest
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Apdu = dt
                    case "org.apache.plc4x.java.bacnetip.readwrite.APDUSimpleAck":
                        var dt *APDUSimpleAck
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Apdu = dt
                    case "org.apache.plc4x.java.bacnetip.readwrite.APDUComplexAck":
                        var dt *APDUComplexAck
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Apdu = dt
                    case "org.apache.plc4x.java.bacnetip.readwrite.APDUSegmentAck":
                        var dt *APDUSegmentAck
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Apdu = dt
                    case "org.apache.plc4x.java.bacnetip.readwrite.APDUError":
                        var dt *APDUError
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Apdu = dt
                    case "org.apache.plc4x.java.bacnetip.readwrite.APDUReject":
                        var dt *APDUReject
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Apdu = dt
                    case "org.apache.plc4x.java.bacnetip.readwrite.APDUAbort":
                        var dt *APDUAbort
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        *m.Apdu = dt
                    }
            }
        }
    }
}

func (m NPDU) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.NPDU"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ProtocolVersionNumber, xml.StartElement{Name: xml.Name{Local: "protocolVersionNumber"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MessageTypeFieldPresent, xml.StartElement{Name: xml.Name{Local: "messageTypeFieldPresent"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DestinationSpecified, xml.StartElement{Name: xml.Name{Local: "destinationSpecified"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SourceSpecified, xml.StartElement{Name: xml.Name{Local: "sourceSpecified"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ExpectingReply, xml.StartElement{Name: xml.Name{Local: "expectingReply"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.NetworkPriority, xml.StartElement{Name: xml.Name{Local: "networkPriority"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DestinationNetworkAddress, xml.StartElement{Name: xml.Name{Local: "destinationNetworkAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DestinationLength, xml.StartElement{Name: xml.Name{Local: "destinationLength"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "destinationAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DestinationAddress, xml.StartElement{Name: xml.Name{Local: "destinationAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "destinationAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SourceNetworkAddress, xml.StartElement{Name: xml.Name{Local: "sourceNetworkAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SourceLength, xml.StartElement{Name: xml.Name{Local: "sourceLength"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "sourceAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SourceAddress, xml.StartElement{Name: xml.Name{Local: "sourceAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "sourceAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.HopCount, xml.StartElement{Name: xml.Name{Local: "hopCount"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Nlm, xml.StartElement{Name: xml.Name{Local: "nlm"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Apdu, xml.StartElement{Name: xml.Name{Local: "apdu"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

