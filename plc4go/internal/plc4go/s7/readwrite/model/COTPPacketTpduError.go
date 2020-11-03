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
type COTPPacketTpduError struct {
    DestinationReference uint16
    RejectCause uint8
    COTPPacket
}

// The corresponding interface
type ICOTPPacketTpduError interface {
    ICOTPPacket
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m COTPPacketTpduError) TpduCode() uint8 {
    return 0x70
}

func (m COTPPacketTpduError) initialize(parameters []ICOTPParameter, payload *IS7Message) spi.Message {
    m.Parameters = parameters
    m.Payload = payload
    return m
}

func NewCOTPPacketTpduError(destinationReference uint16, rejectCause uint8) COTPPacketInitializer {
    return &COTPPacketTpduError{DestinationReference: destinationReference, RejectCause: rejectCause}
}

func CastICOTPPacketTpduError(structType interface{}) ICOTPPacketTpduError {
    castFunc := func(typ interface{}) ICOTPPacketTpduError {
        if iCOTPPacketTpduError, ok := typ.(ICOTPPacketTpduError); ok {
            return iCOTPPacketTpduError
        }
        return nil
    }
    return castFunc(structType)
}

func CastCOTPPacketTpduError(structType interface{}) COTPPacketTpduError {
    castFunc := func(typ interface{}) COTPPacketTpduError {
        if sCOTPPacketTpduError, ok := typ.(COTPPacketTpduError); ok {
            return sCOTPPacketTpduError
        }
        if sCOTPPacketTpduError, ok := typ.(*COTPPacketTpduError); ok {
            return *sCOTPPacketTpduError
        }
        return COTPPacketTpduError{}
    }
    return castFunc(structType)
}

func (m COTPPacketTpduError) LengthInBits() uint16 {
    var lengthInBits uint16 = m.COTPPacket.LengthInBits()

    // Simple field (destinationReference)
    lengthInBits += 16

    // Simple field (rejectCause)
    lengthInBits += 8

    return lengthInBits
}

func (m COTPPacketTpduError) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func COTPPacketTpduErrorParse(io *utils.ReadBuffer) (COTPPacketInitializer, error) {

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

    // Create the instance
    return NewCOTPPacketTpduError(destinationReference, rejectCause), nil
}

func (m COTPPacketTpduError) Serialize(io utils.WriteBuffer) error {
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
    return COTPPacketSerialize(io, m.COTPPacket, CastICOTPPacket(m), ser)
}

func (m *COTPPacketTpduError) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "rejectCause":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.RejectCause = data
            }
        }
    }
}

func (m COTPPacketTpduError) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.COTPPacketTpduError"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DestinationReference, xml.StartElement{Name: xml.Name{Local: "destinationReference"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.RejectCause, xml.StartElement{Name: xml.Name{Local: "rejectCause"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

