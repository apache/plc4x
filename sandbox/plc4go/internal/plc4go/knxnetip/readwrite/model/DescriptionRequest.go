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
type DescriptionRequest struct {
	hpaiControlEndpoint IHPAIControlEndpoint
	KNXNetIPMessage
}

// The corresponding interface
type IDescriptionRequest interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m DescriptionRequest) MsgType() uint16 {
	return 0x0203
}

func (m DescriptionRequest) initialize() spi.Message {
	return m
}

func NewDescriptionRequest(hpaiControlEndpoint IHPAIControlEndpoint) KNXNetIPMessageInitializer {
	return &DescriptionRequest{hpaiControlEndpoint: hpaiControlEndpoint}
}

func CastIDescriptionRequest(structType interface{}) IDescriptionRequest {
	castFunc := func(typ interface{}) IDescriptionRequest {
		if iDescriptionRequest, ok := typ.(IDescriptionRequest); ok {
			return iDescriptionRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastDescriptionRequest(structType interface{}) DescriptionRequest {
	castFunc := func(typ interface{}) DescriptionRequest {
		if sDescriptionRequest, ok := typ.(DescriptionRequest); ok {
			return sDescriptionRequest
		}
		return DescriptionRequest{}
	}
	return castFunc(structType)
}

func (m DescriptionRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (hpaiControlEndpoint)
	lengthInBits += m.hpaiControlEndpoint.LengthInBits()

	return lengthInBits
}

func (m DescriptionRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DescriptionRequestParse(io *spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Simple Field (hpaiControlEndpoint)
	_hpaiControlEndpointMessage, _err := HPAIControlEndpointParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'hpaiControlEndpoint'. " + _err.Error())
	}
	var hpaiControlEndpoint IHPAIControlEndpoint
	hpaiControlEndpoint, _hpaiControlEndpointOk := _hpaiControlEndpointMessage.(IHPAIControlEndpoint)
	if !_hpaiControlEndpointOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiControlEndpointMessage).Name() + " to IHPAIControlEndpoint")
	}

	// Create the instance
	return NewDescriptionRequest(hpaiControlEndpoint), nil
}

func (m DescriptionRequest) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (hpaiControlEndpoint)
		hpaiControlEndpoint := CastIHPAIControlEndpoint(m.hpaiControlEndpoint)
		_hpaiControlEndpointErr := hpaiControlEndpoint.Serialize(io)
		if _hpaiControlEndpointErr != nil {
			return errors.New("Error serializing 'hpaiControlEndpoint' field " + _hpaiControlEndpointErr.Error())
		}

		return nil
	}
	return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}
