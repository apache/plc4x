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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    "strconv"
)

// Constant values.
const TPKTPacket_PROTOCOLID uint8 = 0x03

// The data-structure of this message
type TPKTPacket struct {
    Payload *COTPPacket
    ITPKTPacket
}

// The corresponding interface
type ITPKTPacket interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

func NewTPKTPacket(payload *COTPPacket) *TPKTPacket {
    return &TPKTPacket{Payload: payload}
}

func CastTPKTPacket(structType interface{}) TPKTPacket {
    castFunc := func(typ interface{}) TPKTPacket {
        if casted, ok := typ.(TPKTPacket); ok {
            return casted
        }
        if casted, ok := typ.(*TPKTPacket); ok {
            return *casted
        }
        return TPKTPacket{}
    }
    return castFunc(structType)
}

func (m *TPKTPacket) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Const Field (protocolId)
    lengthInBits += 8

    // Reserved Field (reserved)
    lengthInBits += 8

    // Implicit Field (len)
    lengthInBits += 16

    // Simple field (payload)
    lengthInBits += m.Payload.LengthInBits()

    return lengthInBits
}

func (m *TPKTPacket) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func TPKTPacketParse(io *utils.ReadBuffer) (*TPKTPacket, error) {

    // Const Field (protocolId)
    protocolId, _protocolIdErr := io.ReadUint8(8)
    if _protocolIdErr != nil {
        return nil, errors.New("Error parsing 'protocolId' field " + _protocolIdErr.Error())
    }
    if protocolId != TPKTPacket_PROTOCOLID {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(TPKTPacket_PROTOCOLID)) + " but got " + strconv.Itoa(int(protocolId)))
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint8(0x00) {
            log.WithFields(log.Fields{
                "expected value": uint8(0x00),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Implicit Field (len) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    len, _lenErr := io.ReadUint16(16)
    if _lenErr != nil {
        return nil, errors.New("Error parsing 'len' field " + _lenErr.Error())
    }

    // Simple Field (payload)
    payload, _payloadErr := COTPPacketParse(io, uint16(len) - uint16(uint16(4)))
    if _payloadErr != nil {
        return nil, errors.New("Error parsing 'payload' field " + _payloadErr.Error())
    }

    // Create the instance
    return NewTPKTPacket(payload), nil
}

func (m *TPKTPacket) Serialize(io utils.WriteBuffer) error {

    // Const Field (protocolId)
    _protocolIdErr := io.WriteUint8(8, 0x03)
    if _protocolIdErr != nil {
        return errors.New("Error serializing 'protocolId' field " + _protocolIdErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(8, uint8(0x00))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

    // Implicit Field (len) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    len := uint16(uint16(m.Payload.LengthInBytes()) + uint16(uint16(4)))
    _lenErr := io.WriteUint16(16, (len))
    if _lenErr != nil {
        return errors.New("Error serializing 'len' field " + _lenErr.Error())
    }

    // Simple Field (payload)
    _payloadErr := m.Payload.Serialize(io)
    if _payloadErr != nil {
        return errors.New("Error serializing 'payload' field " + _payloadErr.Error())
    }

    return nil
}

func (m *TPKTPacket) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    for {
        token, err = d.Token()
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
            case "payload":
                var dt *COTPPacket
                if err := d.DecodeElement(&dt, &tok); err != nil {
                    return err
                }
                m.Payload = dt
            }
        }
    }
}

func (m *TPKTPacket) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := "org.apache.plc4x.java.s7.readwrite.TPKTPacket"
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Payload, xml.StartElement{Name: xml.Name{Local: "payload"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

