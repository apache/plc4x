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
    "strconv"
)

// Constant values.
const KNXNetIPMessage_PROTOCOLVERSION uint8 = 0x10

// The data-structure of this message
type KNXNetIPMessage struct {
    Child IKNXNetIPMessageChild
    IKNXNetIPMessage
    IKNXNetIPMessageParent
}

// The corresponding interface
type IKNXNetIPMessage interface {
    MsgType() uint16
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

type IKNXNetIPMessageParent interface {
    SerializeParent(io utils.WriteBuffer, child IKNXNetIPMessage, serializeChildFunction func() error) error
}

type IKNXNetIPMessageChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *KNXNetIPMessage)
    IKNXNetIPMessage
}

func NewKNXNetIPMessage() *KNXNetIPMessage {
    return &KNXNetIPMessage{}
}

func CastKNXNetIPMessage(structType interface{}) KNXNetIPMessage {
    castFunc := func(typ interface{}) KNXNetIPMessage {
        if casted, ok := typ.(KNXNetIPMessage); ok {
            return casted
        }
        if casted, ok := typ.(*KNXNetIPMessage); ok {
            return *casted
        }
        return KNXNetIPMessage{}
    }
    return castFunc(structType)
}

func (m *KNXNetIPMessage) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Implicit Field (headerLength)
    lengthInBits += 8

    // Const Field (protocolVersion)
    lengthInBits += 8

    // Discriminator Field (msgType)
    lengthInBits += 16

    // Implicit Field (totalLength)
    lengthInBits += 16

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *KNXNetIPMessage) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KNXNetIPMessageParse(io *utils.ReadBuffer) (*KNXNetIPMessage, error) {

    // Implicit Field (headerLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    _, _headerLengthErr := io.ReadUint8(8)
    if _headerLengthErr != nil {
        return nil, errors.New("Error parsing 'headerLength' field " + _headerLengthErr.Error())
    }

    // Const Field (protocolVersion)
    protocolVersion, _protocolVersionErr := io.ReadUint8(8)
    if _protocolVersionErr != nil {
        return nil, errors.New("Error parsing 'protocolVersion' field " + _protocolVersionErr.Error())
    }
    if protocolVersion != KNXNetIPMessage_PROTOCOLVERSION {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(KNXNetIPMessage_PROTOCOLVERSION)) + " but got " + strconv.Itoa(int(protocolVersion)))
    }

    // Discriminator Field (msgType) (Used as input to a switch field)
    msgType, _msgTypeErr := io.ReadUint16(16)
    if _msgTypeErr != nil {
        return nil, errors.New("Error parsing 'msgType' field " + _msgTypeErr.Error())
    }

    // Implicit Field (totalLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    totalLength, _totalLengthErr := io.ReadUint16(16)
    if _totalLengthErr != nil {
        return nil, errors.New("Error parsing 'totalLength' field " + _totalLengthErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *KNXNetIPMessage
    var typeSwitchError error
    switch {
    case msgType == 0x0201:
        _parent, typeSwitchError = SearchRequestParse(io)
    case msgType == 0x0202:
        _parent, typeSwitchError = SearchResponseParse(io)
    case msgType == 0x0203:
        _parent, typeSwitchError = DescriptionRequestParse(io)
    case msgType == 0x0204:
        _parent, typeSwitchError = DescriptionResponseParse(io)
    case msgType == 0x0205:
        _parent, typeSwitchError = ConnectionRequestParse(io)
    case msgType == 0x0206:
        _parent, typeSwitchError = ConnectionResponseParse(io)
    case msgType == 0x0207:
        _parent, typeSwitchError = ConnectionStateRequestParse(io)
    case msgType == 0x0208:
        _parent, typeSwitchError = ConnectionStateResponseParse(io)
    case msgType == 0x0209:
        _parent, typeSwitchError = DisconnectRequestParse(io)
    case msgType == 0x020A:
        _parent, typeSwitchError = DisconnectResponseParse(io)
    case msgType == 0x020B:
        _parent, typeSwitchError = UnknownMessageParse(io, totalLength)
    case msgType == 0x0310:
        _parent, typeSwitchError = DeviceConfigurationRequestParse(io, totalLength)
    case msgType == 0x0311:
        _parent, typeSwitchError = DeviceConfigurationAckParse(io)
    case msgType == 0x0420:
        _parent, typeSwitchError = TunnelingRequestParse(io, totalLength)
    case msgType == 0x0421:
        _parent, typeSwitchError = TunnelingResponseParse(io)
    case msgType == 0x0530:
        _parent, typeSwitchError = RoutingIndicationParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *KNXNetIPMessage) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *KNXNetIPMessage) SerializeParent(io utils.WriteBuffer, child IKNXNetIPMessage, serializeChildFunction func() error) error {

    // Implicit Field (headerLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    headerLength := uint8(uint8(6))
    _headerLengthErr := io.WriteUint8(8, (headerLength))
    if _headerLengthErr != nil {
        return errors.New("Error serializing 'headerLength' field " + _headerLengthErr.Error())
    }

    // Const Field (protocolVersion)
    _protocolVersionErr := io.WriteUint8(8, 0x10)
    if _protocolVersionErr != nil {
        return errors.New("Error serializing 'protocolVersion' field " + _protocolVersionErr.Error())
    }

    // Discriminator Field (msgType) (Used as input to a switch field)
    msgType := uint16(child.MsgType())
    _msgTypeErr := io.WriteUint16(16, (msgType))
    if _msgTypeErr != nil {
        return errors.New("Error serializing 'msgType' field " + _msgTypeErr.Error())
    }

    // Implicit Field (totalLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    totalLength := uint16(uint16(m.LengthInBytes()))
    _totalLengthErr := io.WriteUint16(16, (totalLength))
    if _totalLengthErr != nil {
        return errors.New("Error serializing 'totalLength' field " + _totalLengthErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *KNXNetIPMessage) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            }
        }
    }
}

func (m *KNXNetIPMessage) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.KNXNetIPMessage"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

