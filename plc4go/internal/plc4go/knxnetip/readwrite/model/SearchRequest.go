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
type SearchRequest struct {
    HpaiIDiscoveryEndpoint *HPAIDiscoveryEndpoint
    Parent *KnxNetIpMessage
    ISearchRequest
}

// The corresponding interface
type ISearchRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *SearchRequest) MsgType() uint16 {
    return 0x0201
}


func (m *SearchRequest) InitializeParent(parent *KnxNetIpMessage) {
}

func NewSearchRequest(hpaiIDiscoveryEndpoint *HPAIDiscoveryEndpoint, ) *KnxNetIpMessage {
    child := &SearchRequest{
        HpaiIDiscoveryEndpoint: hpaiIDiscoveryEndpoint,
        Parent: NewKnxNetIpMessage(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastSearchRequest(structType interface{}) SearchRequest {
    castFunc := func(typ interface{}) SearchRequest {
        if casted, ok := typ.(SearchRequest); ok {
            return casted
        }
        if casted, ok := typ.(*SearchRequest); ok {
            return *casted
        }
        if casted, ok := typ.(KnxNetIpMessage); ok {
            return CastSearchRequest(casted.Child)
        }
        if casted, ok := typ.(*KnxNetIpMessage); ok {
            return CastSearchRequest(casted.Child)
        }
        return SearchRequest{}
    }
    return castFunc(structType)
}

func (m *SearchRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (hpaiIDiscoveryEndpoint)
    lengthInBits += m.HpaiIDiscoveryEndpoint.LengthInBits()

    return lengthInBits
}

func (m *SearchRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func SearchRequestParse(io *utils.ReadBuffer) (*KnxNetIpMessage, error) {

    // Simple Field (hpaiIDiscoveryEndpoint)
    hpaiIDiscoveryEndpoint, _hpaiIDiscoveryEndpointErr := HPAIDiscoveryEndpointParse(io)
    if _hpaiIDiscoveryEndpointErr != nil {
        return nil, errors.New("Error parsing 'hpaiIDiscoveryEndpoint' field " + _hpaiIDiscoveryEndpointErr.Error())
    }

    // Create a partially initialized instance
    _child := &SearchRequest{
        HpaiIDiscoveryEndpoint: hpaiIDiscoveryEndpoint,
        Parent: &KnxNetIpMessage{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *SearchRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (hpaiIDiscoveryEndpoint)
    _hpaiIDiscoveryEndpointErr := m.HpaiIDiscoveryEndpoint.Serialize(io)
    if _hpaiIDiscoveryEndpointErr != nil {
        return errors.New("Error serializing 'hpaiIDiscoveryEndpoint' field " + _hpaiIDiscoveryEndpointErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *SearchRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "hpaiIDiscoveryEndpoint":
                var data *HPAIDiscoveryEndpoint
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.HpaiIDiscoveryEndpoint = data
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

func (m *SearchRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.HpaiIDiscoveryEndpoint, xml.StartElement{Name: xml.Name{Local: "hpaiIDiscoveryEndpoint"}}); err != nil {
        return err
    }
    return nil
}

