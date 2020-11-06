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
type DescriptionRequest struct {
    HpaiControlEndpoint *HPAIControlEndpoint
    Parent *KNXNetIPMessage
    IDescriptionRequest
}

// The corresponding interface
type IDescriptionRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *DescriptionRequest) MsgType() uint16 {
    return 0x0203
}


func (m *DescriptionRequest) InitializeParent(parent *KNXNetIPMessage) {
}

func NewDescriptionRequest(hpaiControlEndpoint *HPAIControlEndpoint, ) *KNXNetIPMessage {
    child := &DescriptionRequest{
        HpaiControlEndpoint: hpaiControlEndpoint,
        Parent: NewKNXNetIPMessage(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastDescriptionRequest(structType interface{}) DescriptionRequest {
    castFunc := func(typ interface{}) DescriptionRequest {
        if casted, ok := typ.(DescriptionRequest); ok {
            return casted
        }
        if casted, ok := typ.(*DescriptionRequest); ok {
            return *casted
        }
        if casted, ok := typ.(KNXNetIPMessage); ok {
            return CastDescriptionRequest(casted.Child)
        }
        if casted, ok := typ.(*KNXNetIPMessage); ok {
            return CastDescriptionRequest(casted.Child)
        }
        return DescriptionRequest{}
    }
    return castFunc(structType)
}

func (m *DescriptionRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (hpaiControlEndpoint)
    lengthInBits += m.HpaiControlEndpoint.LengthInBits()

    return lengthInBits
}

func (m *DescriptionRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func DescriptionRequestParse(io *utils.ReadBuffer) (*KNXNetIPMessage, error) {

    // Simple Field (hpaiControlEndpoint)
    hpaiControlEndpoint, _hpaiControlEndpointErr := HPAIControlEndpointParse(io)
    if _hpaiControlEndpointErr != nil {
        return nil, errors.New("Error parsing 'hpaiControlEndpoint' field " + _hpaiControlEndpointErr.Error())
    }

    // Create a partially initialized instance
    _child := &DescriptionRequest{
        HpaiControlEndpoint: hpaiControlEndpoint,
        Parent: &KNXNetIPMessage{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *DescriptionRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (hpaiControlEndpoint)
    _hpaiControlEndpointErr := m.HpaiControlEndpoint.Serialize(io)
    if _hpaiControlEndpointErr != nil {
        return errors.New("Error serializing 'hpaiControlEndpoint' field " + _hpaiControlEndpointErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *DescriptionRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "hpaiControlEndpoint":
                var data *HPAIControlEndpoint
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.HpaiControlEndpoint = data
            }
        }
    }
}

func (m *DescriptionRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.DescriptionRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.HpaiControlEndpoint, xml.StartElement{Name: xml.Name{Local: "hpaiControlEndpoint"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

