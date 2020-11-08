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
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
)

// The data-structure of this message
type COTPPacketTpduError struct {
    DestinationReference uint16
    RejectCause uint8
    Parent *COTPPacket
    ICOTPPacketTpduError
}

// The corresponding interface
type ICOTPPacketTpduError interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *COTPPacketTpduError) TpduCode() uint8 {
    return 0x70
}


func (m *COTPPacketTpduError) InitializeParent(parent *COTPPacket, parameters []*COTPParameter, payload *S7Message) {
    m.Parent.Parameters = parameters
    m.Parent.Payload = payload
}

func NewCOTPPacketTpduError(destinationReference uint16, rejectCause uint8, parameters []*COTPParameter, payload *S7Message) *COTPPacket {
    child := &COTPPacketTpduError{
        DestinationReference: destinationReference,
        RejectCause: rejectCause,
        Parent: NewCOTPPacket(parameters, payload),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastCOTPPacketTpduError(structType interface{}) COTPPacketTpduError {
    castFunc := func(typ interface{}) COTPPacketTpduError {
        if casted, ok := typ.(COTPPacketTpduError); ok {
            return casted
        }
        if casted, ok := typ.(*COTPPacketTpduError); ok {
            return *casted
        }
        if casted, ok := typ.(COTPPacket); ok {
            return CastCOTPPacketTpduError(casted.Child)
        }
        if casted, ok := typ.(*COTPPacket); ok {
            return CastCOTPPacketTpduError(casted.Child)
        }
        return COTPPacketTpduError{}
    }
    return castFunc(structType)
}

func (m *COTPPacketTpduError) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (destinationReference)
    lengthInBits += 16

    // Simple field (rejectCause)
    lengthInBits += 8

    return lengthInBits
}

func (m *COTPPacketTpduError) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func COTPPacketTpduErrorParse(io *utils.ReadBuffer) (*COTPPacket, error) {

    // Simple Field (destinationReference)
    destinationReference, _destinationReferenceErr := io.ReadUint16(16)
    if _destinationReferenceErr != nil {
        return nil, errors.New("Error parsing 'destinationReference' field " + _destinationReferenceErr.Error())
    }

    // Simple Field (rejectCause)
    rejectCause, _rejectCauseErr := io.ReadUint8(8)
    if _rejectCauseErr != nil {
        return nil, errors.New("Error parsing 'rejectCause' field " + _rejectCauseErr.Error())
    }

    // Create a partially initialized instance
    _child := &COTPPacketTpduError{
        DestinationReference: destinationReference,
        RejectCause: rejectCause,
        Parent: &COTPPacket{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *COTPPacketTpduError) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (destinationReference)
    destinationReference := uint16(m.DestinationReference)
    _destinationReferenceErr := io.WriteUint16(16, (destinationReference))
    if _destinationReferenceErr != nil {
        return errors.New("Error serializing 'destinationReference' field " + _destinationReferenceErr.Error())
    }

    // Simple Field (rejectCause)
    rejectCause := uint8(m.RejectCause)
    _rejectCauseErr := io.WriteUint8(8, (rejectCause))
    if _rejectCauseErr != nil {
        return errors.New("Error serializing 'rejectCause' field " + _rejectCauseErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *COTPPacketTpduError) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
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
            case "rejectCause":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.RejectCause = data
            }
        }
        token, err = d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
    }
}

func (m *COTPPacketTpduError) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.DestinationReference, xml.StartElement{Name: xml.Name{Local: "destinationReference"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.RejectCause, xml.StartElement{Name: xml.Name{Local: "rejectCause"}}); err != nil {
        return err
    }
    return nil
}

