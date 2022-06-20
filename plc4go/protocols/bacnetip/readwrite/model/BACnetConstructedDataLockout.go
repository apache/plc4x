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

// BACnetConstructedDataLockout is the corresponding interface of BACnetConstructedDataLockout
type BACnetConstructedDataLockout interface {
	BACnetConstructedData
	// GetLockout returns Lockout (property field)
	GetLockout() BACnetApplicationTagBoolean
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagBoolean
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetConstructedDataLockout is the data-structure of this message
type _BACnetConstructedDataLockout struct {
	*_BACnetConstructedData
	Lockout BACnetApplicationTagBoolean

	// Arguments.
	TagNumber          uint8
	ArrayIndexArgument BACnetTagPayloadUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataLockout) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataLockout) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_LOCKOUT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataLockout) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataLockout) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataLockout) GetLockout() BACnetApplicationTagBoolean {
	return m.Lockout
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataLockout) GetActualValue() BACnetApplicationTagBoolean {
	return CastBACnetApplicationTagBoolean(m.GetLockout())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataLockout factory function for _BACnetConstructedDataLockout
func NewBACnetConstructedDataLockout(lockout BACnetApplicationTagBoolean, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataLockout {
	_result := &_BACnetConstructedDataLockout{
		Lockout:                lockout,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataLockout(structType interface{}) BACnetConstructedDataLockout {
	if casted, ok := structType.(BACnetConstructedDataLockout); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataLockout); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataLockout) GetTypeName() string {
	return "BACnetConstructedDataLockout"
}

func (m *_BACnetConstructedDataLockout) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataLockout) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (lockout)
	lengthInBits += m.Lockout.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataLockout) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataLockoutParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataLockout, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataLockout"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataLockout")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (lockout)
	if pullErr := readBuffer.PullContext("lockout"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for lockout")
	}
	_lockout, _lockoutErr := BACnetApplicationTagParse(readBuffer)
	if _lockoutErr != nil {
		return nil, errors.Wrap(_lockoutErr, "Error parsing 'lockout' field")
	}
	lockout := _lockout.(BACnetApplicationTagBoolean)
	if closeErr := readBuffer.CloseContext("lockout"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for lockout")
	}

	// Virtual field
	_actualValue := lockout
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataLockout"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataLockout")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataLockout{
		Lockout:                lockout,
		_BACnetConstructedData: &_BACnetConstructedData{},
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataLockout) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataLockout"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataLockout")
		}

		// Simple Field (lockout)
		if pushErr := writeBuffer.PushContext("lockout"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for lockout")
		}
		_lockoutErr := writeBuffer.WriteSerializable(m.GetLockout())
		if popErr := writeBuffer.PopContext("lockout"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for lockout")
		}
		if _lockoutErr != nil {
			return errors.Wrap(_lockoutErr, "Error serializing 'lockout' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataLockout"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataLockout")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataLockout) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
