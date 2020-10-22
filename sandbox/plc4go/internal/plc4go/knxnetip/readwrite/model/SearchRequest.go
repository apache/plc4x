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
    "errors"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "reflect"
)

// The data-structure of this message
type SearchRequest struct {
    HpaiIDiscoveryEndpoint IHPAIDiscoveryEndpoint
    KNXNetIPMessage
}

// The corresponding interface
type ISearchRequest interface {
    IKNXNetIPMessage
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m SearchRequest) MsgType() uint16 {
    return 0x0201
}

func (m SearchRequest) initialize() spi.Message {
    return m
}

func NewSearchRequest(hpaiIDiscoveryEndpoint IHPAIDiscoveryEndpoint) KNXNetIPMessageInitializer {
    return &SearchRequest{HpaiIDiscoveryEndpoint: hpaiIDiscoveryEndpoint}
}

func CastISearchRequest(structType interface{}) ISearchRequest {
    castFunc := func(typ interface{}) ISearchRequest {
        if iSearchRequest, ok := typ.(ISearchRequest); ok {
            return iSearchRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastSearchRequest(structType interface{}) SearchRequest {
    castFunc := func(typ interface{}) SearchRequest {
        if sSearchRequest, ok := typ.(SearchRequest); ok {
            return sSearchRequest
        }
        if sSearchRequest, ok := typ.(*SearchRequest); ok {
            return *sSearchRequest
        }
        return SearchRequest{}
    }
    return castFunc(structType)
}

func (m SearchRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

    // Simple field (hpaiIDiscoveryEndpoint)
    lengthInBits += m.HpaiIDiscoveryEndpoint.LengthInBits()

    return lengthInBits
}

func (m SearchRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func SearchRequestParse(io *spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

    // Simple Field (hpaiIDiscoveryEndpoint)
    _hpaiIDiscoveryEndpointMessage, _err := HPAIDiscoveryEndpointParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'hpaiIDiscoveryEndpoint'. " + _err.Error())
    }
    var hpaiIDiscoveryEndpoint IHPAIDiscoveryEndpoint
    hpaiIDiscoveryEndpoint, _hpaiIDiscoveryEndpointOk := _hpaiIDiscoveryEndpointMessage.(IHPAIDiscoveryEndpoint)
    if !_hpaiIDiscoveryEndpointOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiIDiscoveryEndpointMessage).Name() + " to IHPAIDiscoveryEndpoint")
    }

    // Create the instance
    return NewSearchRequest(hpaiIDiscoveryEndpoint), nil
}

func (m SearchRequest) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

    // Simple Field (hpaiIDiscoveryEndpoint)
    hpaiIDiscoveryEndpoint := CastIHPAIDiscoveryEndpoint(m.HpaiIDiscoveryEndpoint)
    _hpaiIDiscoveryEndpointErr := hpaiIDiscoveryEndpoint.Serialize(io)
    if _hpaiIDiscoveryEndpointErr != nil {
        return errors.New("Error serializing 'hpaiIDiscoveryEndpoint' field " + _hpaiIDiscoveryEndpointErr.Error())
    }

        return nil
    }
    return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}
