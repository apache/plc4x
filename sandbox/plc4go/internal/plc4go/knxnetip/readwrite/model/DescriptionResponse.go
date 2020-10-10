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
	dibDeviceInfo      DIBDeviceInfo
	dibSuppSvcFamilies DIBSuppSvcFamilies
	KNXNetIPMessage
}

// The corresponding interface
type IDescriptionResponse interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m DescriptionResponse) MsgType() uint16 {
	return 0x0204
}

func (m DescriptionResponse) initialize() spi.Message {
	return m
}

func NewDescriptionResponse(dibDeviceInfo DIBDeviceInfo, dibSuppSvcFamilies DIBSuppSvcFamilies) KNXNetIPMessageInitializer {
	return &DescriptionResponse{dibDeviceInfo: dibDeviceInfo, dibSuppSvcFamilies: dibSuppSvcFamilies}
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

func DescriptionResponseParse(io spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

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
	return NewDescriptionResponse(dibDeviceInfo, dibSuppSvcFamilies), nil
}

func (m DescriptionResponse) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IDescriptionResponse); ok {

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
