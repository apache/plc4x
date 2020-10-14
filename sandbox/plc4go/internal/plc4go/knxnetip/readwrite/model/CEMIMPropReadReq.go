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
)

// The data-structure of this message
type CEMIMPropReadReq struct {
	InterfaceObjectType uint16
	ObjectInstance      uint8
	PropertyId          uint8
	NumberOfElements    uint8
	StartIndex          uint16
	CEMI
}

// The corresponding interface
type ICEMIMPropReadReq interface {
	ICEMI
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m CEMIMPropReadReq) MessageCode() uint8 {
	return 0xFC
}

func (m CEMIMPropReadReq) initialize() spi.Message {
	return m
}

func NewCEMIMPropReadReq(interfaceObjectType uint16, objectInstance uint8, propertyId uint8, numberOfElements uint8, startIndex uint16) CEMIInitializer {
	return &CEMIMPropReadReq{InterfaceObjectType: interfaceObjectType, ObjectInstance: objectInstance, PropertyId: propertyId, NumberOfElements: numberOfElements, StartIndex: startIndex}
}

func CastICEMIMPropReadReq(structType interface{}) ICEMIMPropReadReq {
	castFunc := func(typ interface{}) ICEMIMPropReadReq {
		if iCEMIMPropReadReq, ok := typ.(ICEMIMPropReadReq); ok {
			return iCEMIMPropReadReq
		}
		return nil
	}
	return castFunc(structType)
}

func CastCEMIMPropReadReq(structType interface{}) CEMIMPropReadReq {
	castFunc := func(typ interface{}) CEMIMPropReadReq {
		if sCEMIMPropReadReq, ok := typ.(CEMIMPropReadReq); ok {
			return sCEMIMPropReadReq
		}
		return CEMIMPropReadReq{}
	}
	return castFunc(structType)
}

func (m CEMIMPropReadReq) LengthInBits() uint16 {
	var lengthInBits uint16 = m.CEMI.LengthInBits()

	// Simple field (interfaceObjectType)
	lengthInBits += 16

	// Simple field (objectInstance)
	lengthInBits += 8

	// Simple field (propertyId)
	lengthInBits += 8

	// Simple field (numberOfElements)
	lengthInBits += 4

	// Simple field (startIndex)
	lengthInBits += 12

	return lengthInBits
}

func (m CEMIMPropReadReq) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CEMIMPropReadReqParse(io *spi.ReadBuffer) (CEMIInitializer, error) {

	// Simple Field (interfaceObjectType)
	interfaceObjectType, _interfaceObjectTypeErr := io.ReadUint16(16)
	if _interfaceObjectTypeErr != nil {
		return nil, errors.New("Error parsing 'interfaceObjectType' field " + _interfaceObjectTypeErr.Error())
	}

	// Simple Field (objectInstance)
	objectInstance, _objectInstanceErr := io.ReadUint8(8)
	if _objectInstanceErr != nil {
		return nil, errors.New("Error parsing 'objectInstance' field " + _objectInstanceErr.Error())
	}

	// Simple Field (propertyId)
	propertyId, _propertyIdErr := io.ReadUint8(8)
	if _propertyIdErr != nil {
		return nil, errors.New("Error parsing 'propertyId' field " + _propertyIdErr.Error())
	}

	// Simple Field (numberOfElements)
	numberOfElements, _numberOfElementsErr := io.ReadUint8(4)
	if _numberOfElementsErr != nil {
		return nil, errors.New("Error parsing 'numberOfElements' field " + _numberOfElementsErr.Error())
	}

	// Simple Field (startIndex)
	startIndex, _startIndexErr := io.ReadUint16(12)
	if _startIndexErr != nil {
		return nil, errors.New("Error parsing 'startIndex' field " + _startIndexErr.Error())
	}

	// Create the instance
	return NewCEMIMPropReadReq(interfaceObjectType, objectInstance, propertyId, numberOfElements, startIndex), nil
}

func (m CEMIMPropReadReq) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (interfaceObjectType)
		interfaceObjectType := uint16(m.InterfaceObjectType)
		_interfaceObjectTypeErr := io.WriteUint16(16, (interfaceObjectType))
		if _interfaceObjectTypeErr != nil {
			return errors.New("Error serializing 'interfaceObjectType' field " + _interfaceObjectTypeErr.Error())
		}

		// Simple Field (objectInstance)
		objectInstance := uint8(m.ObjectInstance)
		_objectInstanceErr := io.WriteUint8(8, (objectInstance))
		if _objectInstanceErr != nil {
			return errors.New("Error serializing 'objectInstance' field " + _objectInstanceErr.Error())
		}

		// Simple Field (propertyId)
		propertyId := uint8(m.PropertyId)
		_propertyIdErr := io.WriteUint8(8, (propertyId))
		if _propertyIdErr != nil {
			return errors.New("Error serializing 'propertyId' field " + _propertyIdErr.Error())
		}

		// Simple Field (numberOfElements)
		numberOfElements := uint8(m.NumberOfElements)
		_numberOfElementsErr := io.WriteUint8(4, (numberOfElements))
		if _numberOfElementsErr != nil {
			return errors.New("Error serializing 'numberOfElements' field " + _numberOfElementsErr.Error())
		}

		// Simple Field (startIndex)
		startIndex := uint16(m.StartIndex)
		_startIndexErr := io.WriteUint16(12, (startIndex))
		if _startIndexErr != nil {
			return errors.New("Error serializing 'startIndex' field " + _startIndexErr.Error())
		}

		return nil
	}
	return CEMISerialize(io, m.CEMI, CastICEMI(m), ser)
}
