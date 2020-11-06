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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type COTPPacketConnectionRequest struct {
    DestinationReference uint16
    SourceReference uint16
    ProtocolClass COTPProtocolClass
    Parent *COTPPacket
    ICOTPPacketConnectionRequest
}

// The corresponding interface
type ICOTPPacketConnectionRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *COTPPacketConnectionRequest) TpduCode() uint8 {
    return 0xE0
}


func (m *COTPPacketConnectionRequest) InitializeParent(parent *COTPPacket, parameters []*COTPParameter, payload *S7Message) {
    m.Parent.Parameters = parameters
    m.Parent.Payload = payload
}

func NewCOTPPacketConnectionRequest(destinationReference uint16, sourceReference uint16, protocolClass COTPProtocolClass, parameters []*COTPParameter, payload *S7Message) *COTPPacket {
    child := &COTPPacketConnectionRequest{
        DestinationReference: destinationReference,
        SourceReference: sourceReference,
        ProtocolClass: protocolClass,
        Parent: NewCOTPPacket(parameters, payload),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastCOTPPacketConnectionRequest(structType interface{}) COTPPacketConnectionRequest {
    castFunc := func(typ interface{}) COTPPacketConnectionRequest {
        if casted, ok := typ.(COTPPacketConnectionRequest); ok {
            return casted
        }
        if casted, ok := typ.(*COTPPacketConnectionRequest); ok {
            return *casted
        }
        if casted, ok := typ.(COTPPacket); ok {
            return CastCOTPPacketConnectionRequest(casted.Child)
        }
        if casted, ok := typ.(*COTPPacket); ok {
            return CastCOTPPacketConnectionRequest(casted.Child)
        }
        return COTPPacketConnectionRequest{}
    }
    return castFunc(structType)
}

func (m *COTPPacketConnectionRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (destinationReference)
    lengthInBits += 16

    // Simple field (sourceReference)
    lengthInBits += 16

    // Enum Field (protocolClass)
    lengthInBits += 8

    return lengthInBits
}

func (m *COTPPacketConnectionRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func COTPPacketConnectionRequestParse(io *utils.ReadBuffer) (*COTPPacket, error) {

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

    // Create a partially initialized instance
    _child := &COTPPacketConnectionRequest{
        DestinationReference: destinationReference,
        SourceReference: sourceReference,
        ProtocolClass: protocolClass,
        Parent: &COTPPacket{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *COTPPacketConnectionRequest) Serialize(io utils.WriteBuffer) error {
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
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *COTPPacketConnectionRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
                var data COTPProtocolClass
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ProtocolClass = data
            }
        }
    }
}

func (m *COTPPacketConnectionRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.COTPPacketConnectionRequest"},
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

