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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetConstructedDataChangeOfStateTime is the corresponding interface of BACnetConstructedDataChangeOfStateTime
type BACnetConstructedDataChangeOfStateTime interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetChangeOfStateTime returns ChangeOfStateTime (property field)
	GetChangeOfStateTime() BACnetDateTime
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetDateTime
}

// BACnetConstructedDataChangeOfStateTimeExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataChangeOfStateTime.
// This is useful for switch cases.
type BACnetConstructedDataChangeOfStateTimeExactly interface {
	BACnetConstructedDataChangeOfStateTime
	isBACnetConstructedDataChangeOfStateTime() bool
}

// _BACnetConstructedDataChangeOfStateTime is the data-structure of this message
type _BACnetConstructedDataChangeOfStateTime struct {
	*_BACnetConstructedData
	ChangeOfStateTime BACnetDateTime
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataChangeOfStateTime) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataChangeOfStateTime) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_CHANGE_OF_STATE_TIME
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataChangeOfStateTime) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataChangeOfStateTime) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataChangeOfStateTime) GetChangeOfStateTime() BACnetDateTime {
	return m.ChangeOfStateTime
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataChangeOfStateTime) GetActualValue() BACnetDateTime {
	return CastBACnetDateTime(m.GetChangeOfStateTime())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataChangeOfStateTime factory function for _BACnetConstructedDataChangeOfStateTime
func NewBACnetConstructedDataChangeOfStateTime(changeOfStateTime BACnetDateTime, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataChangeOfStateTime {
	_result := &_BACnetConstructedDataChangeOfStateTime{
		ChangeOfStateTime:      changeOfStateTime,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataChangeOfStateTime(structType interface{}) BACnetConstructedDataChangeOfStateTime {
	if casted, ok := structType.(BACnetConstructedDataChangeOfStateTime); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataChangeOfStateTime); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataChangeOfStateTime) GetTypeName() string {
	return "BACnetConstructedDataChangeOfStateTime"
}

func (m *_BACnetConstructedDataChangeOfStateTime) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataChangeOfStateTime) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (changeOfStateTime)
	lengthInBits += m.ChangeOfStateTime.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataChangeOfStateTime) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataChangeOfStateTimeParse(theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataChangeOfStateTime, error) {
	return BACnetConstructedDataChangeOfStateTimeParseWithBuffer(utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataChangeOfStateTimeParseWithBuffer(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataChangeOfStateTime, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataChangeOfStateTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataChangeOfStateTime")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (changeOfStateTime)
	if pullErr := readBuffer.PullContext("changeOfStateTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for changeOfStateTime")
	}
	_changeOfStateTime, _changeOfStateTimeErr := BACnetDateTimeParseWithBuffer(readBuffer)
	if _changeOfStateTimeErr != nil {
		return nil, errors.Wrap(_changeOfStateTimeErr, "Error parsing 'changeOfStateTime' field of BACnetConstructedDataChangeOfStateTime")
	}
	changeOfStateTime := _changeOfStateTime.(BACnetDateTime)
	if closeErr := readBuffer.CloseContext("changeOfStateTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for changeOfStateTime")
	}

	// Virtual field
	_actualValue := changeOfStateTime
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataChangeOfStateTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataChangeOfStateTime")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataChangeOfStateTime{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		ChangeOfStateTime: changeOfStateTime,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataChangeOfStateTime) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataChangeOfStateTime) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataChangeOfStateTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataChangeOfStateTime")
		}

		// Simple Field (changeOfStateTime)
		if pushErr := writeBuffer.PushContext("changeOfStateTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for changeOfStateTime")
		}
		_changeOfStateTimeErr := writeBuffer.WriteSerializable(m.GetChangeOfStateTime())
		if popErr := writeBuffer.PopContext("changeOfStateTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for changeOfStateTime")
		}
		if _changeOfStateTimeErr != nil {
			return errors.Wrap(_changeOfStateTimeErr, "Error serializing 'changeOfStateTime' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataChangeOfStateTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataChangeOfStateTime")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataChangeOfStateTime) isBACnetConstructedDataChangeOfStateTime() bool {
	return true
}

func (m *_BACnetConstructedDataChangeOfStateTime) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
