/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetConstructedDataLastUseTime is the corresponding interface of BACnetConstructedDataLastUseTime
type BACnetConstructedDataLastUseTime interface {
	BACnetConstructedData
	// GetLastUseTime returns LastUseTime (property field)
	GetLastUseTime() BACnetDateTime
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetDateTime
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetConstructedDataLastUseTime is the data-structure of this message
type _BACnetConstructedDataLastUseTime struct {
	*_BACnetConstructedData
	LastUseTime BACnetDateTime

	// Arguments.
	TagNumber          uint8
	ArrayIndexArgument BACnetTagPayloadUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataLastUseTime) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataLastUseTime) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_LAST_USE_TIME
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataLastUseTime) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataLastUseTime) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataLastUseTime) GetLastUseTime() BACnetDateTime {
	return m.LastUseTime
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataLastUseTime) GetActualValue() BACnetDateTime {
	return CastBACnetDateTime(m.GetLastUseTime())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataLastUseTime factory function for _BACnetConstructedDataLastUseTime
func NewBACnetConstructedDataLastUseTime(lastUseTime BACnetDateTime, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataLastUseTime {
	_result := &_BACnetConstructedDataLastUseTime{
		LastUseTime:            lastUseTime,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataLastUseTime(structType interface{}) BACnetConstructedDataLastUseTime {
	if casted, ok := structType.(BACnetConstructedDataLastUseTime); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataLastUseTime); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataLastUseTime) GetTypeName() string {
	return "BACnetConstructedDataLastUseTime"
}

func (m *_BACnetConstructedDataLastUseTime) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataLastUseTime) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (lastUseTime)
	lengthInBits += m.LastUseTime.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataLastUseTime) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataLastUseTimeParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataLastUseTime, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataLastUseTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataLastUseTime")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (lastUseTime)
	if pullErr := readBuffer.PullContext("lastUseTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for lastUseTime")
	}
	_lastUseTime, _lastUseTimeErr := BACnetDateTimeParse(readBuffer)
	if _lastUseTimeErr != nil {
		return nil, errors.Wrap(_lastUseTimeErr, "Error parsing 'lastUseTime' field")
	}
	lastUseTime := _lastUseTime.(BACnetDateTime)
	if closeErr := readBuffer.CloseContext("lastUseTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for lastUseTime")
	}

	// Virtual field
	_actualValue := lastUseTime
	actualValue := _actualValue.(BACnetDateTime)
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataLastUseTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataLastUseTime")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataLastUseTime{
		LastUseTime:            lastUseTime,
		_BACnetConstructedData: &_BACnetConstructedData{},
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataLastUseTime) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataLastUseTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataLastUseTime")
		}

		// Simple Field (lastUseTime)
		if pushErr := writeBuffer.PushContext("lastUseTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for lastUseTime")
		}
		_lastUseTimeErr := writeBuffer.WriteSerializable(m.GetLastUseTime())
		if popErr := writeBuffer.PopContext("lastUseTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for lastUseTime")
		}
		if _lastUseTimeErr != nil {
			return errors.Wrap(_lastUseTimeErr, "Error serializing 'lastUseTime' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataLastUseTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataLastUseTime")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataLastUseTime) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
