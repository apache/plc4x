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
type DescriptionResponse struct {
	dibDeviceInfo      IDIBDeviceInfo
	dibSuppSvcFamilies IDIBSuppSvcFamilies
	KNXNetIPMessage
}

// The corresponding interface
type IDescriptionResponse interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m DescriptionResponse) MsgType() uint16 {
	return 0x0204
}

func (m DescriptionResponse) initialize() spi.Message {
	return m
}

func NewDescriptionResponse(dibDeviceInfo IDIBDeviceInfo, dibSuppSvcFamilies IDIBSuppSvcFamilies) KNXNetIPMessageInitializer {
	return &DescriptionResponse{dibDeviceInfo: dibDeviceInfo, dibSuppSvcFamilies: dibSuppSvcFamilies}
}

func CastIDescriptionResponse(structType interface{}) IDescriptionResponse {
	castFunc := func(typ interface{}) IDescriptionResponse {
		if iDescriptionResponse, ok := typ.(IDescriptionResponse); ok {
			return iDescriptionResponse
		}
		return nil
	}
	return castFunc(structType)
}

func CastDescriptionResponse(structType interface{}) DescriptionResponse {
	castFunc := func(typ interface{}) DescriptionResponse {
		if sDescriptionResponse, ok := typ.(DescriptionResponse); ok {
			return sDescriptionResponse
		}
		return DescriptionResponse{}
	}
	return castFunc(structType)
}

func (m DescriptionResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (dibDeviceInfo)
	lengthInBits += m.dibDeviceInfo.LengthInBits()

	// Simple field (dibSuppSvcFamilies)
	lengthInBits += m.dibSuppSvcFamilies.LengthInBits()

	return lengthInBits
}

func (m DescriptionResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DescriptionResponseParse(io *spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Simple Field (dibDeviceInfo)
	_dibDeviceInfoMessage, _err := DIBDeviceInfoParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'dibDeviceInfo'. " + _err.Error())
	}
	var dibDeviceInfo IDIBDeviceInfo
	dibDeviceInfo, _dibDeviceInfoOk := _dibDeviceInfoMessage.(IDIBDeviceInfo)
	if !_dibDeviceInfoOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_dibDeviceInfoMessage).Name() + " to IDIBDeviceInfo")
	}

	// Simple Field (dibSuppSvcFamilies)
	_dibSuppSvcFamiliesMessage, _err := DIBSuppSvcFamiliesParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'dibSuppSvcFamilies'. " + _err.Error())
	}
	var dibSuppSvcFamilies IDIBSuppSvcFamilies
	dibSuppSvcFamilies, _dibSuppSvcFamiliesOk := _dibSuppSvcFamiliesMessage.(IDIBSuppSvcFamilies)
	if !_dibSuppSvcFamiliesOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_dibSuppSvcFamiliesMessage).Name() + " to IDIBSuppSvcFamilies")
	}

	// Create the instance
	return NewDescriptionResponse(dibDeviceInfo, dibSuppSvcFamilies), nil
}

func (m DescriptionResponse) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (dibDeviceInfo)
		dibDeviceInfo := CastIDIBDeviceInfo(m.dibDeviceInfo)
		_dibDeviceInfoErr := dibDeviceInfo.Serialize(io)
		if _dibDeviceInfoErr != nil {
			return errors.New("Error serializing 'dibDeviceInfo' field " + _dibDeviceInfoErr.Error())
		}

		// Simple Field (dibSuppSvcFamilies)
		dibSuppSvcFamilies := CastIDIBSuppSvcFamilies(m.dibSuppSvcFamilies)
		_dibSuppSvcFamiliesErr := dibSuppSvcFamilies.Serialize(io)
		if _dibSuppSvcFamiliesErr != nil {
			return errors.New("Error serializing 'dibSuppSvcFamilies' field " + _dibSuppSvcFamiliesErr.Error())
		}

		return nil
	}
	return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}
