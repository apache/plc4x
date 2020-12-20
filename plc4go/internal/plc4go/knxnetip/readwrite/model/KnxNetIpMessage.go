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
    "strconv"
    "reflect"
    "strings"
)

// Constant values.
const KnxNetIpMessage_PROTOCOLVERSION uint8 = 0x10

// The data-structure of this message
type KnxNetIpMessage struct {
    Child IKnxNetIpMessageChild
    IKnxNetIpMessage
    IKnxNetIpMessageParent
}

// The corresponding interface
type IKnxNetIpMessage interface {
    MsgType() uint16
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type IKnxNetIpMessageParent interface {
    SerializeParent(io utils.WriteBuffer, child IKnxNetIpMessage, serializeChildFunction func() error) error
    GetTypeName() string
}

type IKnxNetIpMessageChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *KnxNetIpMessage)
    GetTypeName() string
    IKnxNetIpMessage
}

func NewKnxNetIpMessage() *KnxNetIpMessage {
    return &KnxNetIpMessage{}
}

func CastKnxNetIpMessage(structType interface{}) *KnxNetIpMessage {
    castFunc := func(typ interface{}) *KnxNetIpMessage {
        if casted, ok := typ.(KnxNetIpMessage); ok {
            return &casted
        }
        if casted, ok := typ.(*KnxNetIpMessage); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *KnxNetIpMessage) GetTypeName() string {
    return "KnxNetIpMessage"
}

func (m *KnxNetIpMessage) LengthInBits() uint16 {
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

func (m *KnxNetIpMessage) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxNetIpMessageParse(io *utils.ReadBuffer) (*KnxNetIpMessage, error) {

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
    if protocolVersion != KnxNetIpMessage_PROTOCOLVERSION {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(KnxNetIpMessage_PROTOCOLVERSION)) + " but got " + strconv.Itoa(int(protocolVersion)))
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
    var _parent *KnxNetIpMessage
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

func (m *KnxNetIpMessage) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *KnxNetIpMessage) SerializeParent(io utils.WriteBuffer, child IKnxNetIpMessage, serializeChildFunction func() error) error {

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

func (m *KnxNetIpMessage) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            default:
                switch start.Attr[0].Value {
                    case "org.apache.plc4x.java.knxnetip.readwrite.SearchRequest":
                        var dt *SearchRequest
                        if m.Child != nil {
                            dt = m.Child.(*SearchRequest)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.SearchResponse":
                        var dt *SearchResponse
                        if m.Child != nil {
                            dt = m.Child.(*SearchResponse)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.DescriptionRequest":
                        var dt *DescriptionRequest
                        if m.Child != nil {
                            dt = m.Child.(*DescriptionRequest)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.DescriptionResponse":
                        var dt *DescriptionResponse
                        if m.Child != nil {
                            dt = m.Child.(*DescriptionResponse)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequest":
                        var dt *ConnectionRequest
                        if m.Child != nil {
                            dt = m.Child.(*ConnectionRequest)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.ConnectionResponse":
                        var dt *ConnectionResponse
                        if m.Child != nil {
                            dt = m.Child.(*ConnectionResponse)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.ConnectionStateRequest":
                        var dt *ConnectionStateRequest
                        if m.Child != nil {
                            dt = m.Child.(*ConnectionStateRequest)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.ConnectionStateResponse":
                        var dt *ConnectionStateResponse
                        if m.Child != nil {
                            dt = m.Child.(*ConnectionStateResponse)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.DisconnectRequest":
                        var dt *DisconnectRequest
                        if m.Child != nil {
                            dt = m.Child.(*DisconnectRequest)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.DisconnectResponse":
                        var dt *DisconnectResponse
                        if m.Child != nil {
                            dt = m.Child.(*DisconnectResponse)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.UnknownMessage":
                        var dt *UnknownMessage
                        if m.Child != nil {
                            dt = m.Child.(*UnknownMessage)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.DeviceConfigurationRequest":
                        var dt *DeviceConfigurationRequest
                        if m.Child != nil {
                            dt = m.Child.(*DeviceConfigurationRequest)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.DeviceConfigurationAck":
                        var dt *DeviceConfigurationAck
                        if m.Child != nil {
                            dt = m.Child.(*DeviceConfigurationAck)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.TunnelingRequest":
                        var dt *TunnelingRequest
                        if m.Child != nil {
                            dt = m.Child.(*TunnelingRequest)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.TunnelingResponse":
                        var dt *TunnelingResponse
                        if m.Child != nil {
                            dt = m.Child.(*TunnelingResponse)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.knxnetip.readwrite.RoutingIndication":
                        var dt *RoutingIndication
                        if m.Child != nil {
                            dt = m.Child.(*RoutingIndication)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                }
            }
        }
    }
}

func (m *KnxNetIpMessage) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := reflect.TypeOf(m.Child).String()
    className = "org.apache.plc4x.java.knxnetip.readwrite." + className[strings.LastIndex(className, ".") + 1:]
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    marshaller, ok := m.Child.(xml.Marshaler)
    if !ok {
        return errors.New("child is not castable to Marshaler")
    }
    marshaller.MarshalXML(e, start)
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

