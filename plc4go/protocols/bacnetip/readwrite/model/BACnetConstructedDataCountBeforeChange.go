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

// BACnetConstructedDataCountBeforeChange is the corresponding interface of BACnetConstructedDataCountBeforeChange
type BACnetConstructedDataCountBeforeChange interface {
	BACnetConstructedData
	// GetCountBeforeChange returns CountBeforeChange (property field)
	GetCountBeforeChange() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetConstructedDataCountBeforeChange is the data-structure of this message
type _BACnetConstructedDataCountBeforeChange struct {
	*_BACnetConstructedData
	CountBeforeChange BACnetApplicationTagUnsignedInteger

	// Arguments.
	TagNumber          uint8
	ArrayIndexArgument BACnetTagPayloadUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataCountBeforeChange) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataCountBeforeChange) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_COUNT_BEFORE_CHANGE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataCountBeforeChange) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataCountBeforeChange) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataCountBeforeChange) GetCountBeforeChange() BACnetApplicationTagUnsignedInteger {
	return m.CountBeforeChange
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataCountBeforeChange) GetActualValue() BACnetApplicationTagUnsignedInteger {
	return CastBACnetApplicationTagUnsignedInteger(m.GetCountBeforeChange())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataCountBeforeChange factory function for _BACnetConstructedDataCountBeforeChange
func NewBACnetConstructedDataCountBeforeChange(countBeforeChange BACnetApplicationTagUnsignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataCountBeforeChange {
	_result := &_BACnetConstructedDataCountBeforeChange{
		CountBeforeChange:      countBeforeChange,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataCountBeforeChange(structType interface{}) BACnetConstructedDataCountBeforeChange {
	if casted, ok := structType.(BACnetConstructedDataCountBeforeChange); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataCountBeforeChange); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataCountBeforeChange) GetTypeName() string {
	return "BACnetConstructedDataCountBeforeChange"
}

func (m *_BACnetConstructedDataCountBeforeChange) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataCountBeforeChange) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (countBeforeChange)
	lengthInBits += m.CountBeforeChange.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataCountBeforeChange) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataCountBeforeChangeParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataCountBeforeChange, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataCountBeforeChange"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataCountBeforeChange")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (countBeforeChange)
	if pullErr := readBuffer.PullContext("countBeforeChange"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for countBeforeChange")
	}
	_countBeforeChange, _countBeforeChangeErr := BACnetApplicationTagParse(readBuffer)
	if _countBeforeChangeErr != nil {
		return nil, errors.Wrap(_countBeforeChangeErr, "Error parsing 'countBeforeChange' field")
	}
	countBeforeChange := _countBeforeChange.(BACnetApplicationTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("countBeforeChange"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for countBeforeChange")
	}

	// Virtual field
	_actualValue := countBeforeChange
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataCountBeforeChange"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataCountBeforeChange")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataCountBeforeChange{
		CountBeforeChange:      countBeforeChange,
		_BACnetConstructedData: &_BACnetConstructedData{},
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataCountBeforeChange) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataCountBeforeChange"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataCountBeforeChange")
		}

		// Simple Field (countBeforeChange)
		if pushErr := writeBuffer.PushContext("countBeforeChange"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for countBeforeChange")
		}
		_countBeforeChangeErr := writeBuffer.WriteSerializable(m.GetCountBeforeChange())
		if popErr := writeBuffer.PopContext("countBeforeChange"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for countBeforeChange")
		}
		if _countBeforeChangeErr != nil {
			return errors.Wrap(_countBeforeChangeErr, "Error serializing 'countBeforeChange' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataCountBeforeChange"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataCountBeforeChange")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataCountBeforeChange) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
