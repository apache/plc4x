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

// BACnetConstructedDataRestoreCompletionTime is the corresponding interface of BACnetConstructedDataRestoreCompletionTime
type BACnetConstructedDataRestoreCompletionTime interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetCompletionTime returns CompletionTime (property field)
	GetCompletionTime() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
}

// BACnetConstructedDataRestoreCompletionTimeExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataRestoreCompletionTime.
// This is useful for switch cases.
type BACnetConstructedDataRestoreCompletionTimeExactly interface {
	BACnetConstructedDataRestoreCompletionTime
	isBACnetConstructedDataRestoreCompletionTime() bool
}

// _BACnetConstructedDataRestoreCompletionTime is the data-structure of this message
type _BACnetConstructedDataRestoreCompletionTime struct {
	*_BACnetConstructedData
	CompletionTime BACnetApplicationTagUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataRestoreCompletionTime) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataRestoreCompletionTime) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_RESTORE_COMPLETION_TIME
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataRestoreCompletionTime) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataRestoreCompletionTime) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataRestoreCompletionTime) GetCompletionTime() BACnetApplicationTagUnsignedInteger {
	return m.CompletionTime
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataRestoreCompletionTime) GetActualValue() BACnetApplicationTagUnsignedInteger {
	return CastBACnetApplicationTagUnsignedInteger(m.GetCompletionTime())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataRestoreCompletionTime factory function for _BACnetConstructedDataRestoreCompletionTime
func NewBACnetConstructedDataRestoreCompletionTime(completionTime BACnetApplicationTagUnsignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataRestoreCompletionTime {
	_result := &_BACnetConstructedDataRestoreCompletionTime{
		CompletionTime:         completionTime,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataRestoreCompletionTime(structType interface{}) BACnetConstructedDataRestoreCompletionTime {
	if casted, ok := structType.(BACnetConstructedDataRestoreCompletionTime); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataRestoreCompletionTime); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataRestoreCompletionTime) GetTypeName() string {
	return "BACnetConstructedDataRestoreCompletionTime"
}

func (m *_BACnetConstructedDataRestoreCompletionTime) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataRestoreCompletionTime) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (completionTime)
	lengthInBits += m.CompletionTime.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataRestoreCompletionTime) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataRestoreCompletionTimeParse(theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataRestoreCompletionTime, error) {
	return BACnetConstructedDataRestoreCompletionTimeParseWithBuffer(utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataRestoreCompletionTimeParseWithBuffer(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataRestoreCompletionTime, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataRestoreCompletionTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataRestoreCompletionTime")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (completionTime)
	if pullErr := readBuffer.PullContext("completionTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for completionTime")
	}
	_completionTime, _completionTimeErr := BACnetApplicationTagParseWithBuffer(readBuffer)
	if _completionTimeErr != nil {
		return nil, errors.Wrap(_completionTimeErr, "Error parsing 'completionTime' field of BACnetConstructedDataRestoreCompletionTime")
	}
	completionTime := _completionTime.(BACnetApplicationTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("completionTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for completionTime")
	}

	// Virtual field
	_actualValue := completionTime
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataRestoreCompletionTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataRestoreCompletionTime")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataRestoreCompletionTime{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		CompletionTime: completionTime,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataRestoreCompletionTime) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataRestoreCompletionTime) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataRestoreCompletionTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataRestoreCompletionTime")
		}

		// Simple Field (completionTime)
		if pushErr := writeBuffer.PushContext("completionTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for completionTime")
		}
		_completionTimeErr := writeBuffer.WriteSerializable(m.GetCompletionTime())
		if popErr := writeBuffer.PopContext("completionTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for completionTime")
		}
		if _completionTimeErr != nil {
			return errors.Wrap(_completionTimeErr, "Error serializing 'completionTime' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataRestoreCompletionTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataRestoreCompletionTime")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataRestoreCompletionTime) isBACnetConstructedDataRestoreCompletionTime() bool {
	return true
}

func (m *_BACnetConstructedDataRestoreCompletionTime) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
