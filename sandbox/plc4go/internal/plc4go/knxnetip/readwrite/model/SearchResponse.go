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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
)

// The data-structure of this message
type SearchResponse struct {
	hpaiControlEndpoint HPAIControlEndpoint
	dibDeviceInfo       DIBDeviceInfo
	dibSuppSvcFamilies  DIBSuppSvcFamilies
	KNXNetIPMessage
}

// The corresponding interface
type ISearchResponse interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m SearchResponse) MsgType() uint16 {
	return 0x0202
}

func (m SearchResponse) initialize() spi.Message {
	return m
}

func NewSearchResponse(hpaiControlEndpoint HPAIControlEndpoint, dibDeviceInfo DIBDeviceInfo, dibSuppSvcFamilies DIBSuppSvcFamilies) KNXNetIPMessageInitializer {
	return &SearchResponse{hpaiControlEndpoint: hpaiControlEndpoint, dibDeviceInfo: dibDeviceInfo, dibSuppSvcFamilies: dibSuppSvcFamilies}
}

func (m SearchResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (hpaiControlEndpoint)
	lengthInBits += m.hpaiControlEndpoint.LengthInBits()

	// Simple field (dibDeviceInfo)
	lengthInBits += m.dibDeviceInfo.LengthInBits()

	// Simple field (dibSuppSvcFamilies)
	lengthInBits += m.dibSuppSvcFamilies.LengthInBits()

	return lengthInBits
}

func (m SearchResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func SearchResponseParse(io spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Simple Field (hpaiControlEndpoint)
	_hpaiControlEndpointMessage, _err := HPAIControlEndpointParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'hpaiControlEndpoint'. " + _err.Error())
	}
	var hpaiControlEndpoint HPAIControlEndpoint
	hpaiControlEndpoint, _hpaiControlEndpointOk := _hpaiControlEndpointMessage.(HPAIControlEndpoint)
	if !_hpaiControlEndpointOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiControlEndpointMessage).Name() + " to HPAIControlEndpoint")
	}

	// Simple Field (dibDeviceInfo)
	_dibDeviceInfoMessage, _err := DIBDeviceInfoParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'dibDeviceInfo'. " + _err.Error())
	}
	var dibDeviceInfo DIBDeviceInfo
	dibDeviceInfo, _dibDeviceInfoOk := _dibDeviceInfoMessage.(DIBDeviceInfo)
	if !_dibDeviceInfoOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_dibDeviceInfoMessage).Name() + " to DIBDeviceInfo")
	}

	// Simple Field (dibSuppSvcFamilies)
	_dibSuppSvcFamiliesMessage, _err := DIBSuppSvcFamiliesParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'dibSuppSvcFamilies'. " + _err.Error())
	}
	var dibSuppSvcFamilies DIBSuppSvcFamilies
	dibSuppSvcFamilies, _dibSuppSvcFamiliesOk := _dibSuppSvcFamiliesMessage.(DIBSuppSvcFamilies)
	if !_dibSuppSvcFamiliesOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_dibSuppSvcFamiliesMessage).Name() + " to DIBSuppSvcFamilies")
	}

	// Create the instance
	return NewSearchResponse(hpaiControlEndpoint, dibDeviceInfo, dibSuppSvcFamilies), nil
}

func (m SearchResponse) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(ISearchResponse); ok {

			// Simple Field (hpaiControlEndpoint)
			var hpaiControlEndpoint HPAIControlEndpoint = m.hpaiControlEndpoint
			hpaiControlEndpoint.Serialize(io)

			// Simple Field (dibDeviceInfo)
			var dibDeviceInfo DIBDeviceInfo = m.dibDeviceInfo
			dibDeviceInfo.Serialize(io)

			// Simple Field (dibSuppSvcFamilies)
			var dibSuppSvcFamilies DIBSuppSvcFamilies = m.dibSuppSvcFamilies
			dibSuppSvcFamilies.Serialize(io)
		}
	}
	serializeFunc(m)
}
