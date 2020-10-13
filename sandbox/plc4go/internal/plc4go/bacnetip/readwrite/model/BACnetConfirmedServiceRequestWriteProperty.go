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
	"strconv"
)

// Constant values.
const BACnetConfirmedServiceRequestWriteProperty_OBJECTIDENTIFIERHEADER uint8 = 0x0C
const BACnetConfirmedServiceRequestWriteProperty_PROPERTYIDENTIFIERHEADER uint8 = 0x03
const BACnetConfirmedServiceRequestWriteProperty_OPENINGTAG uint8 = 0x3E
const BACnetConfirmedServiceRequestWriteProperty_CLOSINGTAG uint8 = 0x3F

// The data-structure of this message
type BACnetConfirmedServiceRequestWriteProperty struct {
	objectType               uint16
	objectInstanceNumber     uint32
	propertyIdentifierLength uint8
	propertyIdentifier       []int8
	value                    IBACnetTag
	priority                 *IBACnetTag
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestWriteProperty interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestWriteProperty) ServiceChoice() uint8 {
	return 0x0F
}

func (m BACnetConfirmedServiceRequestWriteProperty) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestWriteProperty(objectType uint16, objectInstanceNumber uint32, propertyIdentifierLength uint8, propertyIdentifier []int8, value IBACnetTag, priority *IBACnetTag) BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestWriteProperty{objectType: objectType, objectInstanceNumber: objectInstanceNumber, propertyIdentifierLength: propertyIdentifierLength, propertyIdentifier: propertyIdentifier, value: value, priority: priority}
}

func CastIBACnetConfirmedServiceRequestWriteProperty(structType interface{}) IBACnetConfirmedServiceRequestWriteProperty {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestWriteProperty {
		if iBACnetConfirmedServiceRequestWriteProperty, ok := typ.(IBACnetConfirmedServiceRequestWriteProperty); ok {
			return iBACnetConfirmedServiceRequestWriteProperty
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestWriteProperty(structType interface{}) BACnetConfirmedServiceRequestWriteProperty {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestWriteProperty {
		if sBACnetConfirmedServiceRequestWriteProperty, ok := typ.(BACnetConfirmedServiceRequestWriteProperty); ok {
			return sBACnetConfirmedServiceRequestWriteProperty
		}
		return BACnetConfirmedServiceRequestWriteProperty{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestWriteProperty) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

	// Const Field (objectIdentifierHeader)
	lengthInBits += 8

	// Simple field (objectType)
	lengthInBits += 10

	// Simple field (objectInstanceNumber)
	lengthInBits += 22

	// Const Field (propertyIdentifierHeader)
	lengthInBits += 5

	// Simple field (propertyIdentifierLength)
	lengthInBits += 3

	// Array field
	if len(m.propertyIdentifier) > 0 {
		lengthInBits += 8 * uint16(len(m.propertyIdentifier))
	}

	// Const Field (openingTag)
	lengthInBits += 8

	// Simple field (value)
	lengthInBits += m.value.LengthInBits()

	// Const Field (closingTag)
	lengthInBits += 8

	// Optional Field (priority)
	if m.priority != nil {
		lengthInBits += (*m.priority).LengthInBits()
	}

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestWriteProperty) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestWritePropertyParse(io *spi.ReadBuffer, len uint16) (BACnetConfirmedServiceRequestInitializer, error) {
	var startPos = io.GetPos()
	var curPos uint16

	// Const Field (objectIdentifierHeader)
	objectIdentifierHeader, _objectIdentifierHeaderErr := io.ReadUint8(8)
	if _objectIdentifierHeaderErr != nil {
		return nil, errors.New("Error parsing 'objectIdentifierHeader' field " + _objectIdentifierHeaderErr.Error())
	}
	if objectIdentifierHeader != BACnetConfirmedServiceRequestWriteProperty_OBJECTIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestWriteProperty_OBJECTIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(objectIdentifierHeader)))
	}

	// Simple Field (objectType)
	objectType, _objectTypeErr := io.ReadUint16(10)
	if _objectTypeErr != nil {
		return nil, errors.New("Error parsing 'objectType' field " + _objectTypeErr.Error())
	}

	// Simple Field (objectInstanceNumber)
	objectInstanceNumber, _objectInstanceNumberErr := io.ReadUint32(22)
	if _objectInstanceNumberErr != nil {
		return nil, errors.New("Error parsing 'objectInstanceNumber' field " + _objectInstanceNumberErr.Error())
	}

	// Const Field (propertyIdentifierHeader)
	propertyIdentifierHeader, _propertyIdentifierHeaderErr := io.ReadUint8(5)
	if _propertyIdentifierHeaderErr != nil {
		return nil, errors.New("Error parsing 'propertyIdentifierHeader' field " + _propertyIdentifierHeaderErr.Error())
	}
	if propertyIdentifierHeader != BACnetConfirmedServiceRequestWriteProperty_PROPERTYIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestWriteProperty_PROPERTYIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(propertyIdentifierHeader)))
	}

	// Simple Field (propertyIdentifierLength)
	propertyIdentifierLength, _propertyIdentifierLengthErr := io.ReadUint8(3)
	if _propertyIdentifierLengthErr != nil {
		return nil, errors.New("Error parsing 'propertyIdentifierLength' field " + _propertyIdentifierLengthErr.Error())
	}

	// Array field (propertyIdentifier)
	// Count array
	propertyIdentifier := make([]int8, propertyIdentifierLength)
	for curItem := uint16(0); curItem < uint16(propertyIdentifierLength); curItem++ {

		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'propertyIdentifier' field " + _err.Error())
		}
		propertyIdentifier[curItem] = _item
	}

	// Const Field (openingTag)
	openingTag, _openingTagErr := io.ReadUint8(8)
	if _openingTagErr != nil {
		return nil, errors.New("Error parsing 'openingTag' field " + _openingTagErr.Error())
	}
	if openingTag != BACnetConfirmedServiceRequestWriteProperty_OPENINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestWriteProperty_OPENINGTAG)) + " but got " + strconv.Itoa(int(openingTag)))
	}

	// Simple Field (value)
	_valueMessage, _err := BACnetTagParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'value'. " + _err.Error())
	}
	var value IBACnetTag
	value, _valueOk := _valueMessage.(IBACnetTag)
	if !_valueOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_valueMessage).Name() + " to IBACnetTag")
	}

	// Const Field (closingTag)
	closingTag, _closingTagErr := io.ReadUint8(8)
	if _closingTagErr != nil {
		return nil, errors.New("Error parsing 'closingTag' field " + _closingTagErr.Error())
	}
	if closingTag != BACnetConfirmedServiceRequestWriteProperty_CLOSINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestWriteProperty_CLOSINGTAG)) + " but got " + strconv.Itoa(int(closingTag)))
	}

	// Optional Field (priority) (Can be skipped, if a given expression evaluates to false)
	curPos = io.GetPos() - startPos
	var priority *IBACnetTag = nil
	if bool((curPos) < ((len) - (1))) {
		_message, _err := BACnetTagParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'priority' field " + _err.Error())
		}
		var _item IBACnetTag
		_item, _ok := _message.(IBACnetTag)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to IBACnetTag")
		}
		priority = &_item
	}

	// Create the instance
	return NewBACnetConfirmedServiceRequestWriteProperty(objectType, objectInstanceNumber, propertyIdentifierLength, propertyIdentifier, value, priority), nil
}

func (m BACnetConfirmedServiceRequestWriteProperty) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Const Field (objectIdentifierHeader)
		_objectIdentifierHeaderErr := io.WriteUint8(8, 0x0C)
		if _objectIdentifierHeaderErr != nil {
			return errors.New("Error serializing 'objectIdentifierHeader' field " + _objectIdentifierHeaderErr.Error())
		}

		// Simple Field (objectType)
		objectType := uint16(m.objectType)
		_objectTypeErr := io.WriteUint16(10, (objectType))
		if _objectTypeErr != nil {
			return errors.New("Error serializing 'objectType' field " + _objectTypeErr.Error())
		}

		// Simple Field (objectInstanceNumber)
		objectInstanceNumber := uint32(m.objectInstanceNumber)
		_objectInstanceNumberErr := io.WriteUint32(22, (objectInstanceNumber))
		if _objectInstanceNumberErr != nil {
			return errors.New("Error serializing 'objectInstanceNumber' field " + _objectInstanceNumberErr.Error())
		}

		// Const Field (propertyIdentifierHeader)
		_propertyIdentifierHeaderErr := io.WriteUint8(5, 0x03)
		if _propertyIdentifierHeaderErr != nil {
			return errors.New("Error serializing 'propertyIdentifierHeader' field " + _propertyIdentifierHeaderErr.Error())
		}

		// Simple Field (propertyIdentifierLength)
		propertyIdentifierLength := uint8(m.propertyIdentifierLength)
		_propertyIdentifierLengthErr := io.WriteUint8(3, (propertyIdentifierLength))
		if _propertyIdentifierLengthErr != nil {
			return errors.New("Error serializing 'propertyIdentifierLength' field " + _propertyIdentifierLengthErr.Error())
		}

		// Array Field (propertyIdentifier)
		if m.propertyIdentifier != nil {
			for _, _element := range m.propertyIdentifier {
				_elementErr := io.WriteInt8(8, _element)
				if _elementErr != nil {
					return errors.New("Error serializing 'propertyIdentifier' field " + _elementErr.Error())
				}
			}
		}

		// Const Field (openingTag)
		_openingTagErr := io.WriteUint8(8, 0x3E)
		if _openingTagErr != nil {
			return errors.New("Error serializing 'openingTag' field " + _openingTagErr.Error())
		}

		// Simple Field (value)
		value := CastIBACnetTag(m.value)
		_valueErr := value.Serialize(io)
		if _valueErr != nil {
			return errors.New("Error serializing 'value' field " + _valueErr.Error())
		}

		// Const Field (closingTag)
		_closingTagErr := io.WriteUint8(8, 0x3F)
		if _closingTagErr != nil {
			return errors.New("Error serializing 'closingTag' field " + _closingTagErr.Error())
		}

		// Optional Field (priority) (Can be skipped, if the value is null)
		var priority *IBACnetTag = nil
		if m.priority != nil {
			priority = m.priority
			_priorityErr := CastIBACnetTag(*priority).Serialize(io)
			if _priorityErr != nil {
				return errors.New("Error serializing 'priority' field " + _priorityErr.Error())
			}
		}

		return nil
	}
	return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
