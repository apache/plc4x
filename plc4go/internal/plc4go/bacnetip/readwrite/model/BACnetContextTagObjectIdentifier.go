/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package model

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// The data-structure of this message
type BACnetContextTagObjectIdentifier struct {
	*BACnetContextTag
	ObjectType     BACnetObjectType
	InstanceNumber uint32
}

// The corresponding interface
type IBACnetContextTagObjectIdentifier interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(writeBuffer utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *BACnetContextTagObjectIdentifier) DataType() BACnetDataType {
	return BACnetDataType_BACNET_OBJECT_IDENTIFIER
}

func (m *BACnetContextTagObjectIdentifier) InitializeParent(parent *BACnetContextTag, tagNumber uint8, tagClass TagClass, lengthValueType uint8, extTagNumber *uint8, extLength *uint8, extExtLength *uint16, extExtExtLength *uint32, actualTagNumber uint8, actualLength uint32) {
	m.TagNumber = tagNumber
	m.TagClass = tagClass
	m.LengthValueType = lengthValueType
	m.ExtTagNumber = extTagNumber
	m.ExtLength = extLength
	m.ExtExtLength = extExtLength
	m.ExtExtExtLength = extExtExtLength
}

func NewBACnetContextTagObjectIdentifier(objectType BACnetObjectType, instanceNumber uint32, tagNumber uint8, tagClass TagClass, lengthValueType uint8, extTagNumber *uint8, extLength *uint8, extExtLength *uint16, extExtExtLength *uint32, actualTagNumber uint8, actualLength uint32) *BACnetContextTag {
	child := &BACnetContextTagObjectIdentifier{
		ObjectType:       objectType,
		InstanceNumber:   instanceNumber,
		BACnetContextTag: NewBACnetContextTag(tagNumber, tagClass, lengthValueType, extTagNumber, extLength, extExtLength, extExtExtLength, actualTagNumber, actualLength),
	}
	child.Child = child
	return child.BACnetContextTag
}

func CastBACnetContextTagObjectIdentifier(structType interface{}) *BACnetContextTagObjectIdentifier {
	castFunc := func(typ interface{}) *BACnetContextTagObjectIdentifier {
		if casted, ok := typ.(BACnetContextTagObjectIdentifier); ok {
			return &casted
		}
		if casted, ok := typ.(*BACnetContextTagObjectIdentifier); ok {
			return casted
		}
		if casted, ok := typ.(BACnetContextTag); ok {
			return CastBACnetContextTagObjectIdentifier(casted.Child)
		}
		if casted, ok := typ.(*BACnetContextTag); ok {
			return CastBACnetContextTagObjectIdentifier(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *BACnetContextTagObjectIdentifier) GetTypeName() string {
	return "BACnetContextTagObjectIdentifier"
}

func (m *BACnetContextTagObjectIdentifier) LengthInBits() uint16 {
	return m.LengthInBitsConditional(false)
}

func (m *BACnetContextTagObjectIdentifier) LengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.ParentLengthInBits())

	// Simple field (objectType)
	lengthInBits += 10

	// Simple field (instanceNumber)
	lengthInBits += 22

	return lengthInBits
}

func (m *BACnetContextTagObjectIdentifier) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetContextTagObjectIdentifierParse(readBuffer utils.ReadBuffer, tagNumberArgument uint8, dataType BACnetDataType) (*BACnetContextTag, error) {
	if pullErr := readBuffer.PullContext("BACnetContextTagObjectIdentifier"); pullErr != nil {
		return nil, pullErr
	}

	// Simple Field (objectType)
	if pullErr := readBuffer.PullContext("objectType"); pullErr != nil {
		return nil, pullErr
	}
	_objectType, _objectTypeErr := BACnetObjectTypeParse(readBuffer)
	if _objectTypeErr != nil {
		return nil, errors.Wrap(_objectTypeErr, "Error parsing 'objectType' field")
	}
	objectType := _objectType
	if closeErr := readBuffer.CloseContext("objectType"); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (instanceNumber)
	_instanceNumber, _instanceNumberErr := readBuffer.ReadUint32("instanceNumber", 22)
	if _instanceNumberErr != nil {
		return nil, errors.Wrap(_instanceNumberErr, "Error parsing 'instanceNumber' field")
	}
	instanceNumber := _instanceNumber

	if closeErr := readBuffer.CloseContext("BACnetContextTagObjectIdentifier"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &BACnetContextTagObjectIdentifier{
		ObjectType:       objectType,
		InstanceNumber:   instanceNumber,
		BACnetContextTag: &BACnetContextTag{},
	}
	_child.BACnetContextTag.Child = _child
	return _child.BACnetContextTag, nil
}

func (m *BACnetContextTagObjectIdentifier) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetContextTagObjectIdentifier"); pushErr != nil {
			return pushErr
		}

		// Simple Field (objectType)
		if pushErr := writeBuffer.PushContext("objectType"); pushErr != nil {
			return pushErr
		}
		_objectTypeErr := m.ObjectType.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("objectType"); popErr != nil {
			return popErr
		}
		if _objectTypeErr != nil {
			return errors.Wrap(_objectTypeErr, "Error serializing 'objectType' field")
		}

		// Simple Field (instanceNumber)
		instanceNumber := uint32(m.InstanceNumber)
		_instanceNumberErr := writeBuffer.WriteUint32("instanceNumber", 22, (instanceNumber))
		if _instanceNumberErr != nil {
			return errors.Wrap(_instanceNumberErr, "Error serializing 'instanceNumber' field")
		}

		if popErr := writeBuffer.PopContext("BACnetContextTagObjectIdentifier"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *BACnetContextTagObjectIdentifier) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	m.Serialize(buffer)
	return buffer.GetBox().String()
}
