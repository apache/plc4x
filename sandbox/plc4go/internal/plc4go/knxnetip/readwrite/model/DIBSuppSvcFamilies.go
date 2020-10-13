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
type DIBSuppSvcFamilies struct {
	descriptionType uint8
	serviceIds      []IServiceId
}

// The corresponding interface
type IDIBSuppSvcFamilies interface {
	spi.Message
	Serialize(io spi.WriteBuffer) error
}

func NewDIBSuppSvcFamilies(descriptionType uint8, serviceIds []IServiceId) spi.Message {
	return &DIBSuppSvcFamilies{descriptionType: descriptionType, serviceIds: serviceIds}
}

func CastIDIBSuppSvcFamilies(structType interface{}) IDIBSuppSvcFamilies {
	castFunc := func(typ interface{}) IDIBSuppSvcFamilies {
		if iDIBSuppSvcFamilies, ok := typ.(IDIBSuppSvcFamilies); ok {
			return iDIBSuppSvcFamilies
		}
		return nil
	}
	return castFunc(structType)
}

func CastDIBSuppSvcFamilies(structType interface{}) DIBSuppSvcFamilies {
	castFunc := func(typ interface{}) DIBSuppSvcFamilies {
		if sDIBSuppSvcFamilies, ok := typ.(DIBSuppSvcFamilies); ok {
			return sDIBSuppSvcFamilies
		}
		return DIBSuppSvcFamilies{}
	}
	return castFunc(structType)
}

func (m DIBSuppSvcFamilies) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Implicit Field (structureLength)
	lengthInBits += 8

	// Simple field (descriptionType)
	lengthInBits += 8

	// Array field
	if len(m.serviceIds) > 0 {
		for _, element := range m.serviceIds {
			lengthInBits += element.LengthInBits()
		}
	}

	return lengthInBits
}

func (m DIBSuppSvcFamilies) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DIBSuppSvcFamiliesParse(io *spi.ReadBuffer) (spi.Message, error) {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _structureLengthErr := io.ReadUint8(8)
	if _structureLengthErr != nil {
		return nil, errors.New("Error parsing 'structureLength' field " + _structureLengthErr.Error())
	}

	// Simple Field (descriptionType)
	descriptionType, _descriptionTypeErr := io.ReadUint8(8)
	if _descriptionTypeErr != nil {
		return nil, errors.New("Error parsing 'descriptionType' field " + _descriptionTypeErr.Error())
	}

	// Array field (serviceIds)
	// Count array
	serviceIds := make([]IServiceId, uint16(3))
	for curItem := uint16(0); curItem < uint16(uint16(3)); curItem++ {

		_message, _err := ServiceIdParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'serviceIds' field " + _err.Error())
		}
		var _item IServiceId
		_item, _ok := _message.(IServiceId)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ServiceId")
		}
		serviceIds[curItem] = _item
	}

	// Create the instance
	return NewDIBSuppSvcFamilies(descriptionType, serviceIds), nil
}

func (m DIBSuppSvcFamilies) Serialize(io spi.WriteBuffer) error {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	structureLength := uint8(uint8(m.LengthInBytes()))
	_structureLengthErr := io.WriteUint8(8, (structureLength))
	if _structureLengthErr != nil {
		return errors.New("Error serializing 'structureLength' field " + _structureLengthErr.Error())
	}

	// Simple Field (descriptionType)
	descriptionType := uint8(m.descriptionType)
	_descriptionTypeErr := io.WriteUint8(8, (descriptionType))
	if _descriptionTypeErr != nil {
		return errors.New("Error serializing 'descriptionType' field " + _descriptionTypeErr.Error())
	}

	// Array Field (serviceIds)
	if m.serviceIds != nil {
		for _, _element := range m.serviceIds {
			_elementErr := _element.Serialize(io)
			if _elementErr != nil {
				return errors.New("Error serializing 'serviceIds' field " + _elementErr.Error())
			}
		}
	}

	return nil
}
