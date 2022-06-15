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

// BACnetConstructedDataFailedAttemptsTime is the corresponding interface of BACnetConstructedDataFailedAttemptsTime
type BACnetConstructedDataFailedAttemptsTime interface {
	BACnetConstructedData
	// GetFailedAttemptsTime returns FailedAttemptsTime (property field)
	GetFailedAttemptsTime() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetConstructedDataFailedAttemptsTime is the data-structure of this message
type _BACnetConstructedDataFailedAttemptsTime struct {
	*_BACnetConstructedData
	FailedAttemptsTime BACnetApplicationTagUnsignedInteger

	// Arguments.
	TagNumber          uint8
	ArrayIndexArgument BACnetTagPayloadUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataFailedAttemptsTime) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataFailedAttemptsTime) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_FAILED_ATTEMPTS_TIME
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataFailedAttemptsTime) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataFailedAttemptsTime) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataFailedAttemptsTime) GetFailedAttemptsTime() BACnetApplicationTagUnsignedInteger {
	return m.FailedAttemptsTime
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataFailedAttemptsTime) GetActualValue() BACnetApplicationTagUnsignedInteger {
	return CastBACnetApplicationTagUnsignedInteger(m.GetFailedAttemptsTime())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataFailedAttemptsTime factory function for _BACnetConstructedDataFailedAttemptsTime
func NewBACnetConstructedDataFailedAttemptsTime(failedAttemptsTime BACnetApplicationTagUnsignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataFailedAttemptsTime {
	_result := &_BACnetConstructedDataFailedAttemptsTime{
		FailedAttemptsTime:     failedAttemptsTime,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataFailedAttemptsTime(structType interface{}) BACnetConstructedDataFailedAttemptsTime {
	if casted, ok := structType.(BACnetConstructedDataFailedAttemptsTime); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataFailedAttemptsTime); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataFailedAttemptsTime) GetTypeName() string {
	return "BACnetConstructedDataFailedAttemptsTime"
}

func (m *_BACnetConstructedDataFailedAttemptsTime) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataFailedAttemptsTime) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (failedAttemptsTime)
	lengthInBits += m.FailedAttemptsTime.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataFailedAttemptsTime) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataFailedAttemptsTimeParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataFailedAttemptsTime, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataFailedAttemptsTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataFailedAttemptsTime")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (failedAttemptsTime)
	if pullErr := readBuffer.PullContext("failedAttemptsTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for failedAttemptsTime")
	}
	_failedAttemptsTime, _failedAttemptsTimeErr := BACnetApplicationTagParse(readBuffer)
	if _failedAttemptsTimeErr != nil {
		return nil, errors.Wrap(_failedAttemptsTimeErr, "Error parsing 'failedAttemptsTime' field")
	}
	failedAttemptsTime := _failedAttemptsTime.(BACnetApplicationTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("failedAttemptsTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for failedAttemptsTime")
	}

	// Virtual field
	_actualValue := failedAttemptsTime
	actualValue := _actualValue.(BACnetApplicationTagUnsignedInteger)
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataFailedAttemptsTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataFailedAttemptsTime")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataFailedAttemptsTime{
		FailedAttemptsTime:     failedAttemptsTime,
		_BACnetConstructedData: &_BACnetConstructedData{},
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataFailedAttemptsTime) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataFailedAttemptsTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataFailedAttemptsTime")
		}

		// Simple Field (failedAttemptsTime)
		if pushErr := writeBuffer.PushContext("failedAttemptsTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for failedAttemptsTime")
		}
		_failedAttemptsTimeErr := writeBuffer.WriteSerializable(m.GetFailedAttemptsTime())
		if popErr := writeBuffer.PopContext("failedAttemptsTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for failedAttemptsTime")
		}
		if _failedAttemptsTimeErr != nil {
			return errors.Wrap(_failedAttemptsTimeErr, "Error serializing 'failedAttemptsTime' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataFailedAttemptsTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataFailedAttemptsTime")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataFailedAttemptsTime) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
