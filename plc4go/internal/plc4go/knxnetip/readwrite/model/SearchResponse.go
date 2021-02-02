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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    "io"
)

// The data-structure of this message
type SearchResponse struct {
    HpaiControlEndpoint *HPAIControlEndpoint
    DibDeviceInfo *DIBDeviceInfo
    DibSuppSvcFamilies *DIBSuppSvcFamilies
    Parent *KnxNetIpMessage
    ISearchResponse
}

// The corresponding interface
type ISearchResponse interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *SearchResponse) MsgType() uint16 {
    return 0x0202
}


func (m *SearchResponse) InitializeParent(parent *KnxNetIpMessage) {
}

func NewSearchResponse(hpaiControlEndpoint *HPAIControlEndpoint, dibDeviceInfo *DIBDeviceInfo, dibSuppSvcFamilies *DIBSuppSvcFamilies, ) *KnxNetIpMessage {
    child := &SearchResponse{
        HpaiControlEndpoint: hpaiControlEndpoint,
        DibDeviceInfo: dibDeviceInfo,
        DibSuppSvcFamilies: dibSuppSvcFamilies,
        Parent: NewKnxNetIpMessage(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastSearchResponse(structType interface{}) *SearchResponse {
    castFunc := func(typ interface{}) *SearchResponse {
        if casted, ok := typ.(SearchResponse); ok {
            return &casted
        }
        if casted, ok := typ.(*SearchResponse); ok {
            return casted
        }
        if casted, ok := typ.(KnxNetIpMessage); ok {
            return CastSearchResponse(casted.Child)
        }
        if casted, ok := typ.(*KnxNetIpMessage); ok {
            return CastSearchResponse(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *SearchResponse) GetTypeName() string {
    return "SearchResponse"
}

func (m *SearchResponse) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (hpaiControlEndpoint)
    lengthInBits += m.HpaiControlEndpoint.LengthInBits()

    // Simple field (dibDeviceInfo)
    lengthInBits += m.DibDeviceInfo.LengthInBits()

    // Simple field (dibSuppSvcFamilies)
    lengthInBits += m.DibSuppSvcFamilies.LengthInBits()

    return lengthInBits
}

func (m *SearchResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func SearchResponseParse(io *utils.ReadBuffer) (*KnxNetIpMessage, error) {

    // Simple Field (hpaiControlEndpoint)
    hpaiControlEndpoint, _hpaiControlEndpointErr := HPAIControlEndpointParse(io)
    if _hpaiControlEndpointErr != nil {
        return nil, errors.New("Error parsing 'hpaiControlEndpoint' field " + _hpaiControlEndpointErr.Error())
    }

    // Simple Field (dibDeviceInfo)
    dibDeviceInfo, _dibDeviceInfoErr := DIBDeviceInfoParse(io)
    if _dibDeviceInfoErr != nil {
        return nil, errors.New("Error parsing 'dibDeviceInfo' field " + _dibDeviceInfoErr.Error())
    }

    // Simple Field (dibSuppSvcFamilies)
    dibSuppSvcFamilies, _dibSuppSvcFamiliesErr := DIBSuppSvcFamiliesParse(io)
    if _dibSuppSvcFamiliesErr != nil {
        return nil, errors.New("Error parsing 'dibSuppSvcFamilies' field " + _dibSuppSvcFamiliesErr.Error())
    }

    // Create a partially initialized instance
    _child := &SearchResponse{
        HpaiControlEndpoint: hpaiControlEndpoint,
        DibDeviceInfo: dibDeviceInfo,
        DibSuppSvcFamilies: dibSuppSvcFamilies,
        Parent: &KnxNetIpMessage{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *SearchResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (hpaiControlEndpoint)
    _hpaiControlEndpointErr := m.HpaiControlEndpoint.Serialize(io)
    if _hpaiControlEndpointErr != nil {
        return errors.New("Error serializing 'hpaiControlEndpoint' field " + _hpaiControlEndpointErr.Error())
    }

    // Simple Field (dibDeviceInfo)
    _dibDeviceInfoErr := m.DibDeviceInfo.Serialize(io)
    if _dibDeviceInfoErr != nil {
        return errors.New("Error serializing 'dibDeviceInfo' field " + _dibDeviceInfoErr.Error())
    }

    // Simple Field (dibSuppSvcFamilies)
    _dibSuppSvcFamiliesErr := m.DibSuppSvcFamilies.Serialize(io)
    if _dibSuppSvcFamiliesErr != nil {
        return errors.New("Error serializing 'dibSuppSvcFamilies' field " + _dibSuppSvcFamiliesErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *SearchResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
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
            case "dibDeviceInfo":
                var data *DIBDeviceInfo
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.DibDeviceInfo = data
            case "dibSuppSvcFamilies":
                var data *DIBSuppSvcFamilies
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.DibSuppSvcFamilies = data
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

func (m *SearchResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.HpaiControlEndpoint, xml.StartElement{Name: xml.Name{Local: "hpaiControlEndpoint"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DibDeviceInfo, xml.StartElement{Name: xml.Name{Local: "dibDeviceInfo"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.DibSuppSvcFamilies, xml.StartElement{Name: xml.Name{Local: "dibSuppSvcFamilies"}}); err != nil {
        return err
    }
    return nil
}

