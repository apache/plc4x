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
    "reflect"
    "strings"
)

// The data-structure of this message
type NLM struct {
    VendorId *uint16
    Child INLMChild
    INLM
    INLMParent
}

// The corresponding interface
type INLM interface {
    MessageType() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type INLMParent interface {
    SerializeParent(io utils.WriteBuffer, child INLM, serializeChildFunction func() error) error
}

type INLMChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *NLM, vendorId *uint16)
    INLM
}

func NewNLM(vendorId *uint16) *NLM {
    return &NLM{VendorId: vendorId}
}

func CastNLM(structType interface{}) *NLM {
    castFunc := func(typ interface{}) *NLM {
        if casted, ok := typ.(NLM); ok {
            return &casted
        }
        if casted, ok := typ.(*NLM); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *NLM) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (messageType)
    lengthInBits += 8

    // Optional Field (vendorId)
    if m.VendorId != nil {
        lengthInBits += 16
    }

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *NLM) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func NLMParse(io *utils.ReadBuffer, apduLength uint16) (*NLM, error) {

    // Discriminator Field (messageType) (Used as input to a switch field)
    messageType, _messageTypeErr := io.ReadUint8(8)
    if _messageTypeErr != nil {
        return nil, errors.New("Error parsing 'messageType' field " + _messageTypeErr.Error())
    }

    // Optional Field (vendorId) (Can be skipped, if a given expression evaluates to false)
    var vendorId *uint16 = nil
    if bool(bool(bool((messageType) >= ((128))))) && bool(bool(bool((messageType) <= ((255))))) {
        _val, _err := io.ReadUint16(16)
        if _err != nil {
            return nil, errors.New("Error parsing 'vendorId' field " + _err.Error())
        }

        vendorId = &_val
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *NLM
    var typeSwitchError error
    switch {
    case messageType == 0x0:
        _parent, typeSwitchError = NLMWhoIsRouterToNetworkParse(io, apduLength, messageType)
    case messageType == 0x1:
        _parent, typeSwitchError = NLMIAmRouterToNetworkParse(io, apduLength, messageType)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent, vendorId)
    return _parent, nil
}

func (m *NLM) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *NLM) SerializeParent(io utils.WriteBuffer, child INLM, serializeChildFunction func() error) error {

    // Discriminator Field (messageType) (Used as input to a switch field)
    messageType := uint8(child.MessageType())
    _messageTypeErr := io.WriteUint8(8, (messageType))
    if _messageTypeErr != nil {
        return errors.New("Error serializing 'messageType' field " + _messageTypeErr.Error())
    }

    // Optional Field (vendorId) (Can be skipped, if the value is null)
    var vendorId *uint16 = nil
    if m.VendorId != nil {
        vendorId = m.VendorId
        _vendorIdErr := io.WriteUint16(16, *(vendorId))
        if _vendorIdErr != nil {
            return errors.New("Error serializing 'vendorId' field " + _vendorIdErr.Error())
        }
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *NLM) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "vendorId":
                var data *uint16
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.VendorId = data
                default:
                    switch start.Attr[0].Value {
                        case "org.apache.plc4x.java.bacnetip.readwrite.NLMWhoIsRouterToNetwork":
                            var dt *NLMWhoIsRouterToNetwork
                            if m.Child != nil {
                                dt = m.Child.(*NLMWhoIsRouterToNetwork)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.NLMIAmRouterToNetwork":
                            var dt *NLMIAmRouterToNetwork
                            if m.Child != nil {
                                dt = m.Child.(*NLMIAmRouterToNetwork)
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

func (m *NLM) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := reflect.TypeOf(m.Child).String()
    className = "org.apache.plc4x.java.bacnetip.readwrite." + className[strings.LastIndex(className, ".") + 1:]
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.VendorId, xml.StartElement{Name: xml.Name{Local: "vendorId"}}); err != nil {
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

