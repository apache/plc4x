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
type TunnelingResponseDataBlock struct {
    CommunicationChannelId uint8
    SequenceCounter uint8
    Status Status
    ITunnelingResponseDataBlock
}

// The corresponding interface
type ITunnelingResponseDataBlock interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

func NewTunnelingResponseDataBlock(communicationChannelId uint8, sequenceCounter uint8, status Status) *TunnelingResponseDataBlock {
    return &TunnelingResponseDataBlock{CommunicationChannelId: communicationChannelId, SequenceCounter: sequenceCounter, Status: status}
}

func CastTunnelingResponseDataBlock(structType interface{}) TunnelingResponseDataBlock {
    castFunc := func(typ interface{}) TunnelingResponseDataBlock {
        if casted, ok := typ.(TunnelingResponseDataBlock); ok {
            return casted
        }
        if casted, ok := typ.(*TunnelingResponseDataBlock); ok {
            return *casted
        }
        return TunnelingResponseDataBlock{}
    }
    return castFunc(structType)
}

func (m *TunnelingResponseDataBlock) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Implicit Field (structureLength)
    lengthInBits += 8

    // Simple field (communicationChannelId)
    lengthInBits += 8

    // Simple field (sequenceCounter)
    lengthInBits += 8

    // Enum Field (status)
    lengthInBits += 8

    return lengthInBits
}

func (m *TunnelingResponseDataBlock) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func TunnelingResponseDataBlockParse(io *utils.ReadBuffer) (*TunnelingResponseDataBlock, error) {

    // Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    _, _structureLengthErr := io.ReadUint8(8)
    if _structureLengthErr != nil {
        return nil, errors.New("Error parsing 'structureLength' field " + _structureLengthErr.Error())
    }

    // Simple Field (communicationChannelId)
    communicationChannelId, _communicationChannelIdErr := io.ReadUint8(8)
    if _communicationChannelIdErr != nil {
        return nil, errors.New("Error parsing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
    }

    // Simple Field (sequenceCounter)
    sequenceCounter, _sequenceCounterErr := io.ReadUint8(8)
    if _sequenceCounterErr != nil {
        return nil, errors.New("Error parsing 'sequenceCounter' field " + _sequenceCounterErr.Error())
    }

    // Enum field (status)
    status, _statusErr := StatusParse(io)
    if _statusErr != nil {
        return nil, errors.New("Error parsing 'status' field " + _statusErr.Error())
    }

    // Create the instance
    return NewTunnelingResponseDataBlock(communicationChannelId, sequenceCounter, status), nil
}

func (m *TunnelingResponseDataBlock) Serialize(io utils.WriteBuffer) error {

    // Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    structureLength := uint8(uint8(m.LengthInBytes()))
    _structureLengthErr := io.WriteUint8(8, (structureLength))
    if _structureLengthErr != nil {
        return errors.New("Error serializing 'structureLength' field " + _structureLengthErr.Error())
    }

    // Simple Field (communicationChannelId)
    communicationChannelId := uint8(m.CommunicationChannelId)
    _communicationChannelIdErr := io.WriteUint8(8, (communicationChannelId))
    if _communicationChannelIdErr != nil {
        return errors.New("Error serializing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
    }

    // Simple Field (sequenceCounter)
    sequenceCounter := uint8(m.SequenceCounter)
    _sequenceCounterErr := io.WriteUint8(8, (sequenceCounter))
    if _sequenceCounterErr != nil {
        return errors.New("Error serializing 'sequenceCounter' field " + _sequenceCounterErr.Error())
    }

    // Enum field (status)
    status := CastStatus(m.Status)
    _statusErr := status.Serialize(io)
    if _statusErr != nil {
        return errors.New("Error serializing 'status' field " + _statusErr.Error())
    }

    return nil
}

func (m *TunnelingResponseDataBlock) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "communicationChannelId":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.CommunicationChannelId = data
            case "sequenceCounter":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SequenceCounter = data
            case "status":
                var data Status
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Status = data
            }
        }
    }
}

func (m *TunnelingResponseDataBlock) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := "org.apache.plc4x.java.knxnetip.readwrite.TunnelingResponseDataBlock"
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.CommunicationChannelId, xml.StartElement{Name: xml.Name{Local: "communicationChannelId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SequenceCounter, xml.StartElement{Name: xml.Name{Local: "sequenceCounter"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Status, xml.StartElement{Name: xml.Name{Local: "status"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

