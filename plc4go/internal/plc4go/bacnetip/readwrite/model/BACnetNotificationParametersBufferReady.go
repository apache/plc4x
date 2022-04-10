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
type BACnetNotificationParametersBufferReady struct {
	*BACnetNotificationParameters
	InnerOpeningTag      *BACnetOpeningTag
	BufferProperty       *BACnetDeviceObjectPropertyReference
	PreviousNotification *BACnetContextTagUnsignedInteger
	CurrentNotification  *BACnetContextTagUnsignedInteger
	InnerClosingTag      *BACnetClosingTag

	// Arguments.
	TagNumber  uint8
	ObjectType BACnetObjectType
}

// The corresponding interface
type IBACnetNotificationParametersBufferReady interface {
	IBACnetNotificationParameters
	// GetInnerOpeningTag returns InnerOpeningTag (property field)
	GetInnerOpeningTag() *BACnetOpeningTag
	// GetBufferProperty returns BufferProperty (property field)
	GetBufferProperty() *BACnetDeviceObjectPropertyReference
	// GetPreviousNotification returns PreviousNotification (property field)
	GetPreviousNotification() *BACnetContextTagUnsignedInteger
	// GetCurrentNotification returns CurrentNotification (property field)
	GetCurrentNotification() *BACnetContextTagUnsignedInteger
	// GetInnerClosingTag returns InnerClosingTag (property field)
	GetInnerClosingTag() *BACnetClosingTag
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////
///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *BACnetNotificationParametersBufferReady) InitializeParent(parent *BACnetNotificationParameters, openingTag *BACnetOpeningTag, peekedTagHeader *BACnetTagHeader, closingTag *BACnetClosingTag) {
	m.BACnetNotificationParameters.OpeningTag = openingTag
	m.BACnetNotificationParameters.PeekedTagHeader = peekedTagHeader
	m.BACnetNotificationParameters.ClosingTag = closingTag
}

func (m *BACnetNotificationParametersBufferReady) GetParent() *BACnetNotificationParameters {
	return m.BACnetNotificationParameters
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////
func (m *BACnetNotificationParametersBufferReady) GetInnerOpeningTag() *BACnetOpeningTag {
	return m.InnerOpeningTag
}

func (m *BACnetNotificationParametersBufferReady) GetBufferProperty() *BACnetDeviceObjectPropertyReference {
	return m.BufferProperty
}

func (m *BACnetNotificationParametersBufferReady) GetPreviousNotification() *BACnetContextTagUnsignedInteger {
	return m.PreviousNotification
}

func (m *BACnetNotificationParametersBufferReady) GetCurrentNotification() *BACnetContextTagUnsignedInteger {
	return m.CurrentNotification
}

func (m *BACnetNotificationParametersBufferReady) GetInnerClosingTag() *BACnetClosingTag {
	return m.InnerClosingTag
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetNotificationParametersBufferReady factory function for BACnetNotificationParametersBufferReady
func NewBACnetNotificationParametersBufferReady(innerOpeningTag *BACnetOpeningTag, bufferProperty *BACnetDeviceObjectPropertyReference, previousNotification *BACnetContextTagUnsignedInteger, currentNotification *BACnetContextTagUnsignedInteger, innerClosingTag *BACnetClosingTag, openingTag *BACnetOpeningTag, peekedTagHeader *BACnetTagHeader, closingTag *BACnetClosingTag, tagNumber uint8, objectType BACnetObjectType) *BACnetNotificationParametersBufferReady {
	_result := &BACnetNotificationParametersBufferReady{
		InnerOpeningTag:              innerOpeningTag,
		BufferProperty:               bufferProperty,
		PreviousNotification:         previousNotification,
		CurrentNotification:          currentNotification,
		InnerClosingTag:              innerClosingTag,
		BACnetNotificationParameters: NewBACnetNotificationParameters(openingTag, peekedTagHeader, closingTag, tagNumber, objectType),
	}
	_result.Child = _result
	return _result
}

func CastBACnetNotificationParametersBufferReady(structType interface{}) *BACnetNotificationParametersBufferReady {
	if casted, ok := structType.(BACnetNotificationParametersBufferReady); ok {
		return &casted
	}
	if casted, ok := structType.(*BACnetNotificationParametersBufferReady); ok {
		return casted
	}
	if casted, ok := structType.(BACnetNotificationParameters); ok {
		return CastBACnetNotificationParametersBufferReady(casted.Child)
	}
	if casted, ok := structType.(*BACnetNotificationParameters); ok {
		return CastBACnetNotificationParametersBufferReady(casted.Child)
	}
	return nil
}

func (m *BACnetNotificationParametersBufferReady) GetTypeName() string {
	return "BACnetNotificationParametersBufferReady"
}

func (m *BACnetNotificationParametersBufferReady) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *BACnetNotificationParametersBufferReady) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (innerOpeningTag)
	lengthInBits += m.InnerOpeningTag.GetLengthInBits()

	// Simple field (bufferProperty)
	lengthInBits += m.BufferProperty.GetLengthInBits()

	// Simple field (previousNotification)
	lengthInBits += m.PreviousNotification.GetLengthInBits()

	// Simple field (currentNotification)
	lengthInBits += m.CurrentNotification.GetLengthInBits()

	// Simple field (innerClosingTag)
	lengthInBits += m.InnerClosingTag.GetLengthInBits()

	return lengthInBits
}

func (m *BACnetNotificationParametersBufferReady) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetNotificationParametersBufferReadyParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectType BACnetObjectType, peekedTagNumber uint8) (*BACnetNotificationParametersBufferReady, error) {
	if pullErr := readBuffer.PullContext("BACnetNotificationParametersBufferReady"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	// Simple Field (innerOpeningTag)
	if pullErr := readBuffer.PullContext("innerOpeningTag"); pullErr != nil {
		return nil, pullErr
	}
	_innerOpeningTag, _innerOpeningTagErr := BACnetContextTagParse(readBuffer, uint8(peekedTagNumber), BACnetDataType(BACnetDataType_OPENING_TAG))
	if _innerOpeningTagErr != nil {
		return nil, errors.Wrap(_innerOpeningTagErr, "Error parsing 'innerOpeningTag' field")
	}
	innerOpeningTag := CastBACnetOpeningTag(_innerOpeningTag)
	if closeErr := readBuffer.CloseContext("innerOpeningTag"); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (bufferProperty)
	if pullErr := readBuffer.PullContext("bufferProperty"); pullErr != nil {
		return nil, pullErr
	}
	_bufferProperty, _bufferPropertyErr := BACnetDeviceObjectPropertyReferenceParse(readBuffer, uint8(uint8(0)))
	if _bufferPropertyErr != nil {
		return nil, errors.Wrap(_bufferPropertyErr, "Error parsing 'bufferProperty' field")
	}
	bufferProperty := CastBACnetDeviceObjectPropertyReference(_bufferProperty)
	if closeErr := readBuffer.CloseContext("bufferProperty"); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (previousNotification)
	if pullErr := readBuffer.PullContext("previousNotification"); pullErr != nil {
		return nil, pullErr
	}
	_previousNotification, _previousNotificationErr := BACnetContextTagParse(readBuffer, uint8(uint8(1)), BACnetDataType(BACnetDataType_UNSIGNED_INTEGER))
	if _previousNotificationErr != nil {
		return nil, errors.Wrap(_previousNotificationErr, "Error parsing 'previousNotification' field")
	}
	previousNotification := CastBACnetContextTagUnsignedInteger(_previousNotification)
	if closeErr := readBuffer.CloseContext("previousNotification"); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (currentNotification)
	if pullErr := readBuffer.PullContext("currentNotification"); pullErr != nil {
		return nil, pullErr
	}
	_currentNotification, _currentNotificationErr := BACnetContextTagParse(readBuffer, uint8(uint8(2)), BACnetDataType(BACnetDataType_UNSIGNED_INTEGER))
	if _currentNotificationErr != nil {
		return nil, errors.Wrap(_currentNotificationErr, "Error parsing 'currentNotification' field")
	}
	currentNotification := CastBACnetContextTagUnsignedInteger(_currentNotification)
	if closeErr := readBuffer.CloseContext("currentNotification"); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (innerClosingTag)
	if pullErr := readBuffer.PullContext("innerClosingTag"); pullErr != nil {
		return nil, pullErr
	}
	_innerClosingTag, _innerClosingTagErr := BACnetContextTagParse(readBuffer, uint8(peekedTagNumber), BACnetDataType(BACnetDataType_CLOSING_TAG))
	if _innerClosingTagErr != nil {
		return nil, errors.Wrap(_innerClosingTagErr, "Error parsing 'innerClosingTag' field")
	}
	innerClosingTag := CastBACnetClosingTag(_innerClosingTag)
	if closeErr := readBuffer.CloseContext("innerClosingTag"); closeErr != nil {
		return nil, closeErr
	}

	if closeErr := readBuffer.CloseContext("BACnetNotificationParametersBufferReady"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &BACnetNotificationParametersBufferReady{
		InnerOpeningTag:              CastBACnetOpeningTag(innerOpeningTag),
		BufferProperty:               CastBACnetDeviceObjectPropertyReference(bufferProperty),
		PreviousNotification:         CastBACnetContextTagUnsignedInteger(previousNotification),
		CurrentNotification:          CastBACnetContextTagUnsignedInteger(currentNotification),
		InnerClosingTag:              CastBACnetClosingTag(innerClosingTag),
		BACnetNotificationParameters: &BACnetNotificationParameters{},
	}
	_child.BACnetNotificationParameters.Child = _child
	return _child, nil
}

func (m *BACnetNotificationParametersBufferReady) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetNotificationParametersBufferReady"); pushErr != nil {
			return pushErr
		}

		// Simple Field (innerOpeningTag)
		if pushErr := writeBuffer.PushContext("innerOpeningTag"); pushErr != nil {
			return pushErr
		}
		_innerOpeningTagErr := m.InnerOpeningTag.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("innerOpeningTag"); popErr != nil {
			return popErr
		}
		if _innerOpeningTagErr != nil {
			return errors.Wrap(_innerOpeningTagErr, "Error serializing 'innerOpeningTag' field")
		}

		// Simple Field (bufferProperty)
		if pushErr := writeBuffer.PushContext("bufferProperty"); pushErr != nil {
			return pushErr
		}
		_bufferPropertyErr := m.BufferProperty.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("bufferProperty"); popErr != nil {
			return popErr
		}
		if _bufferPropertyErr != nil {
			return errors.Wrap(_bufferPropertyErr, "Error serializing 'bufferProperty' field")
		}

		// Simple Field (previousNotification)
		if pushErr := writeBuffer.PushContext("previousNotification"); pushErr != nil {
			return pushErr
		}
		_previousNotificationErr := m.PreviousNotification.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("previousNotification"); popErr != nil {
			return popErr
		}
		if _previousNotificationErr != nil {
			return errors.Wrap(_previousNotificationErr, "Error serializing 'previousNotification' field")
		}

		// Simple Field (currentNotification)
		if pushErr := writeBuffer.PushContext("currentNotification"); pushErr != nil {
			return pushErr
		}
		_currentNotificationErr := m.CurrentNotification.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("currentNotification"); popErr != nil {
			return popErr
		}
		if _currentNotificationErr != nil {
			return errors.Wrap(_currentNotificationErr, "Error serializing 'currentNotification' field")
		}

		// Simple Field (innerClosingTag)
		if pushErr := writeBuffer.PushContext("innerClosingTag"); pushErr != nil {
			return pushErr
		}
		_innerClosingTagErr := m.InnerClosingTag.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("innerClosingTag"); popErr != nil {
			return popErr
		}
		if _innerClosingTagErr != nil {
			return errors.Wrap(_innerClosingTagErr, "Error serializing 'innerClosingTag' field")
		}

		if popErr := writeBuffer.PopContext("BACnetNotificationParametersBufferReady"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *BACnetNotificationParametersBufferReady) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
