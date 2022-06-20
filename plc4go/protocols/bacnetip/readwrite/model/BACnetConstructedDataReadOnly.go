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

// BACnetConstructedDataReadOnly is the corresponding interface of BACnetConstructedDataReadOnly
type BACnetConstructedDataReadOnly interface {
	BACnetConstructedData
	// GetReadOnly returns ReadOnly (property field)
	GetReadOnly() BACnetApplicationTagBoolean
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagBoolean
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetConstructedDataReadOnly is the data-structure of this message
type _BACnetConstructedDataReadOnly struct {
	*_BACnetConstructedData
	ReadOnly BACnetApplicationTagBoolean

	// Arguments.
	TagNumber          uint8
	ArrayIndexArgument BACnetTagPayloadUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataReadOnly) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataReadOnly) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_READ_ONLY
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataReadOnly) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataReadOnly) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataReadOnly) GetReadOnly() BACnetApplicationTagBoolean {
	return m.ReadOnly
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataReadOnly) GetActualValue() BACnetApplicationTagBoolean {
	return CastBACnetApplicationTagBoolean(m.GetReadOnly())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataReadOnly factory function for _BACnetConstructedDataReadOnly
func NewBACnetConstructedDataReadOnly(readOnly BACnetApplicationTagBoolean, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataReadOnly {
	_result := &_BACnetConstructedDataReadOnly{
		ReadOnly:               readOnly,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataReadOnly(structType interface{}) BACnetConstructedDataReadOnly {
	if casted, ok := structType.(BACnetConstructedDataReadOnly); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataReadOnly); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataReadOnly) GetTypeName() string {
	return "BACnetConstructedDataReadOnly"
}

func (m *_BACnetConstructedDataReadOnly) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataReadOnly) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (readOnly)
	lengthInBits += m.ReadOnly.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataReadOnly) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataReadOnlyParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataReadOnly, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataReadOnly"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataReadOnly")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (readOnly)
	if pullErr := readBuffer.PullContext("readOnly"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for readOnly")
	}
	_readOnly, _readOnlyErr := BACnetApplicationTagParse(readBuffer)
	if _readOnlyErr != nil {
		return nil, errors.Wrap(_readOnlyErr, "Error parsing 'readOnly' field")
	}
	readOnly := _readOnly.(BACnetApplicationTagBoolean)
	if closeErr := readBuffer.CloseContext("readOnly"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for readOnly")
	}

	// Virtual field
	_actualValue := readOnly
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataReadOnly"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataReadOnly")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataReadOnly{
		ReadOnly:               readOnly,
		_BACnetConstructedData: &_BACnetConstructedData{},
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataReadOnly) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataReadOnly"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataReadOnly")
		}

		// Simple Field (readOnly)
		if pushErr := writeBuffer.PushContext("readOnly"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for readOnly")
		}
		_readOnlyErr := writeBuffer.WriteSerializable(m.GetReadOnly())
		if popErr := writeBuffer.PopContext("readOnly"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for readOnly")
		}
		if _readOnlyErr != nil {
			return errors.Wrap(_readOnlyErr, "Error serializing 'readOnly' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataReadOnly"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataReadOnly")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataReadOnly) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
