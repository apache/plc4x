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
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

// The data-structure of this message
type ConnectionStateResponse struct {
    CommunicationChannelId uint8
    Status Status
    Parent *KnxNetIpMessage
    IConnectionStateResponse
}

// The corresponding interface
type IConnectionStateResponse interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ConnectionStateResponse) MsgType() uint16 {
    return 0x0208
}


func (m *ConnectionStateResponse) InitializeParent(parent *KnxNetIpMessage) {
}

func NewConnectionStateResponse(communicationChannelId uint8, status Status, ) *KnxNetIpMessage {
    child := &ConnectionStateResponse{
        CommunicationChannelId: communicationChannelId,
        Status: status,
        Parent: NewKnxNetIpMessage(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastConnectionStateResponse(structType interface{}) *ConnectionStateResponse {
    castFunc := func(typ interface{}) *ConnectionStateResponse {
        if casted, ok := typ.(ConnectionStateResponse); ok {
            return &casted
        }
        if casted, ok := typ.(*ConnectionStateResponse); ok {
            return casted
        }
        if casted, ok := typ.(KnxNetIpMessage); ok {
            return CastConnectionStateResponse(casted.Child)
        }
        if casted, ok := typ.(*KnxNetIpMessage); ok {
            return CastConnectionStateResponse(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *ConnectionStateResponse) GetTypeName() string {
    return "ConnectionStateResponse"
}

func (m *ConnectionStateResponse) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (communicationChannelId)
    lengthInBits += 8

    // Simple field (status)
    lengthInBits += 8

    return lengthInBits
}

func (m *ConnectionStateResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ConnectionStateResponseParse(io *utils.ReadBuffer) (*KnxNetIpMessage, error) {

    // Simple Field (communicationChannelId)
    communicationChannelId, _communicationChannelIdErr := io.ReadUint8(8)
    if _communicationChannelIdErr != nil {
        return nil, errors.New("Error parsing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
    }

    // Simple Field (status)
    status, _statusErr := StatusParse(io)
    if _statusErr != nil {
        return nil, errors.New("Error parsing 'status' field " + _statusErr.Error())
    }

    // Create a partially initialized instance
    _child := &ConnectionStateResponse{
        CommunicationChannelId: communicationChannelId,
        Status: status,
        Parent: &KnxNetIpMessage{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *ConnectionStateResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (communicationChannelId)
    communicationChannelId := uint8(m.CommunicationChannelId)
    _communicationChannelIdErr := io.WriteUint8(8, (communicationChannelId))
    if _communicationChannelIdErr != nil {
        return errors.New("Error serializing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
    }

    // Simple Field (status)
    _statusErr := m.Status.Serialize(io)
    if _statusErr != nil {
        return errors.New("Error serializing 'status' field " + _statusErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *ConnectionStateResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
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
            case "status":
                var data Status
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Status = data
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

func (m *ConnectionStateResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.CommunicationChannelId, xml.StartElement{Name: xml.Name{Local: "communicationChannelId"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Status, xml.StartElement{Name: xml.Name{Local: "status"}}); err != nil {
        return err
    }
    return nil
}

