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

// BACnetConstructedDataMinimumOffTime is the corresponding interface of BACnetConstructedDataMinimumOffTime
type BACnetConstructedDataMinimumOffTime interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetMinimumOffTime returns MinimumOffTime (property field)
	GetMinimumOffTime() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
}

// BACnetConstructedDataMinimumOffTimeExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataMinimumOffTime.
// This is useful for switch cases.
type BACnetConstructedDataMinimumOffTimeExactly interface {
	BACnetConstructedDataMinimumOffTime
	isBACnetConstructedDataMinimumOffTime() bool
}

// _BACnetConstructedDataMinimumOffTime is the data-structure of this message
type _BACnetConstructedDataMinimumOffTime struct {
	*_BACnetConstructedData
	MinimumOffTime BACnetApplicationTagUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataMinimumOffTime) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataMinimumOffTime) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_MINIMUM_OFF_TIME
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataMinimumOffTime) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataMinimumOffTime) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataMinimumOffTime) GetMinimumOffTime() BACnetApplicationTagUnsignedInteger {
	return m.MinimumOffTime
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataMinimumOffTime) GetActualValue() BACnetApplicationTagUnsignedInteger {
	return CastBACnetApplicationTagUnsignedInteger(m.GetMinimumOffTime())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataMinimumOffTime factory function for _BACnetConstructedDataMinimumOffTime
func NewBACnetConstructedDataMinimumOffTime(minimumOffTime BACnetApplicationTagUnsignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataMinimumOffTime {
	_result := &_BACnetConstructedDataMinimumOffTime{
		MinimumOffTime:         minimumOffTime,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataMinimumOffTime(structType interface{}) BACnetConstructedDataMinimumOffTime {
	if casted, ok := structType.(BACnetConstructedDataMinimumOffTime); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataMinimumOffTime); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataMinimumOffTime) GetTypeName() string {
	return "BACnetConstructedDataMinimumOffTime"
}

func (m *_BACnetConstructedDataMinimumOffTime) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataMinimumOffTime) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (minimumOffTime)
	lengthInBits += m.MinimumOffTime.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataMinimumOffTime) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataMinimumOffTimeParse(theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataMinimumOffTime, error) {
	return BACnetConstructedDataMinimumOffTimeParseWithBuffer(utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataMinimumOffTimeParseWithBuffer(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataMinimumOffTime, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataMinimumOffTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataMinimumOffTime")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (minimumOffTime)
	if pullErr := readBuffer.PullContext("minimumOffTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for minimumOffTime")
	}
	_minimumOffTime, _minimumOffTimeErr := BACnetApplicationTagParseWithBuffer(readBuffer)
	if _minimumOffTimeErr != nil {
		return nil, errors.Wrap(_minimumOffTimeErr, "Error parsing 'minimumOffTime' field of BACnetConstructedDataMinimumOffTime")
	}
	minimumOffTime := _minimumOffTime.(BACnetApplicationTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("minimumOffTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for minimumOffTime")
	}

	// Virtual field
	_actualValue := minimumOffTime
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataMinimumOffTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataMinimumOffTime")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataMinimumOffTime{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		MinimumOffTime: minimumOffTime,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataMinimumOffTime) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataMinimumOffTime) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataMinimumOffTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataMinimumOffTime")
		}

		// Simple Field (minimumOffTime)
		if pushErr := writeBuffer.PushContext("minimumOffTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for minimumOffTime")
		}
		_minimumOffTimeErr := writeBuffer.WriteSerializable(m.GetMinimumOffTime())
		if popErr := writeBuffer.PopContext("minimumOffTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for minimumOffTime")
		}
		if _minimumOffTimeErr != nil {
			return errors.Wrap(_minimumOffTimeErr, "Error serializing 'minimumOffTime' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataMinimumOffTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataMinimumOffTime")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataMinimumOffTime) isBACnetConstructedDataMinimumOffTime() bool {
	return true
}

func (m *_BACnetConstructedDataMinimumOffTime) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
