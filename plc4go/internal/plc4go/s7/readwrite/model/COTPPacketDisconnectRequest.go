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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type COTPPacketDisconnectRequest struct {
    DestinationReference uint16
    SourceReference uint16
    ProtocolClass ICOTPProtocolClass
    COTPPacket
}

// The corresponding interface
type ICOTPPacketDisconnectRequest interface {
    ICOTPPacket
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m COTPPacketDisconnectRequest) TpduCode() uint8 {
    return 0x80
}

func (m COTPPacketDisconnectRequest) initialize(parameters []ICOTPParameter, payload *IS7Message) spi.Message {
    m.Parameters = parameters
    m.Payload = payload
    return m
}

func NewCOTPPacketDisconnectRequest(destinationReference uint16, sourceReference uint16, protocolClass ICOTPProtocolClass) COTPPacketInitializer {
    return &COTPPacketDisconnectRequest{DestinationReference: destinationReference, SourceReference: sourceReference, ProtocolClass: protocolClass}
}

func CastICOTPPacketDisconnectRequest(structType interface{}) ICOTPPacketDisconnectRequest {
    castFunc := func(typ interface{}) ICOTPPacketDisconnectRequest {
        if iCOTPPacketDisconnectRequest, ok := typ.(ICOTPPacketDisconnectRequest); ok {
            return iCOTPPacketDisconnectRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastCOTPPacketDisconnectRequest(structType interface{}) COTPPacketDisconnectRequest {
    castFunc := func(typ interface{}) COTPPacketDisconnectRequest {
        if sCOTPPacketDisconnectRequest, ok := typ.(COTPPacketDisconnectRequest); ok {
            return sCOTPPacketDisconnectRequest
        }
        if sCOTPPacketDisconnectRequest, ok := typ.(*COTPPacketDisconnectRequest); ok {
            return *sCOTPPacketDisconnectRequest
        }
        return COTPPacketDisconnectRequest{}
    }
    return castFunc(structType)
}

func (m COTPPacketDisconnectRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.COTPPacket.LengthInBits()

    // Simple field (destinationReference)
    lengthInBits += 16

    // Simple field (sourceReference)
    lengthInBits += 16

    // Enum Field (protocolClass)
    lengthInBits += 8

    return lengthInBits
}

func (m COTPPacketDisconnectRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func COTPPacketDisconnectRequestParse(io *utils.ReadBuffer) (COTPPacketInitializer, error) {

    // Simple Field (destinationReference)
    destinationReference, _destinationReferenceErr := io.ReadUint16(16)
    if _destinationReferenceErr != nil {
        return nil, errors.New("Error parsing 'destinationReference' field " + _destinationReferenceErr.Error())
    }

    // Simple Field (sourceReference)
    sourceReference, _sourceReferenceErr := io.ReadUint16(16)
    if _sourceReferenceErr != nil {
        return nil, errors.New("Error parsing 'sourceReference' field " + _sourceReferenceErr.Error())
    }

    // Enum field (protocolClass)
    protocolClass, _protocolClassErr := COTPProtocolClassParse(io)
    if _protocolClassErr != nil {
        return nil, errors.New("Error parsing 'protocolClass' field " + _protocolClassErr.Error())
    }

    // Create the instance
    return NewCOTPPacketDisconnectRequest(destinationReference, sourceReference, protocolClass), nil
}

func (m COTPPacketDisconnectRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (destinationReference)
    destinationReference := uint16(m.DestinationReference)
    _destinationReferenceErr := io.WriteUint16(16, (destinationReference))
    if _destinationReferenceErr != nil {
        return errors.New("Error serializing 'destinationReference' field " + _destinationReferenceErr.Error())
    }

    // Simple Field (sourceReference)
    sourceReference := uint16(m.SourceReference)
    _sourceReferenceErr := io.WriteUint16(16, (sourceReference))
    if _sourceReferenceErr != nil {
        return errors.New("Error serializing 'sourceReference' field " + _sourceReferenceErr.Error())
    }

    // Enum field (protocolClass)
    protocolClass := CastCOTPProtocolClass(m.ProtocolClass)
    _protocolClassErr := protocolClass.Serialize(io)
    if _protocolClassErr != nil {
        return errors.New("Error serializing 'protocolClass' field " + _protocolClassErr.Error())
    }

        return nil
    }
    return COTPPacketSerialize(io, m.COTPPacket, CastICOTPPacket(m), ser)
}

func (m *COTPPacketDisconnectRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "destinationReference":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.DestinationReference = data
            case "sourceReference":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SourceReference = data
            case "protocolClass":
                var data *COTPProtocolClass
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ProtocolClass = data
            }
        }
    }
}

func (m COTPPacketDisconnectRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.COTPPacketDisconnectRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DestinationReference, xml.StartElement{Name: xml.Name{Local: "destinationReference"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SourceReference, xml.StartElement{Name: xml.Name{Local: "sourceReference"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ProtocolClass, xml.StartElement{Name: xml.Name{Local: "protocolClass"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

