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

// BACnetConstructedDataLargeAnalogValueAll is the corresponding interface of BACnetConstructedDataLargeAnalogValueAll
type BACnetConstructedDataLargeAnalogValueAll interface {
	BACnetConstructedData
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetConstructedDataLargeAnalogValueAll is the data-structure of this message
type _BACnetConstructedDataLargeAnalogValueAll struct {
	*_BACnetConstructedData

	// Arguments.
	TagNumber          uint8
	ArrayIndexArgument BACnetTagPayloadUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataLargeAnalogValueAll) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_LARGE_ANALOG_VALUE
}

func (m *_BACnetConstructedDataLargeAnalogValueAll) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_ALL
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataLargeAnalogValueAll) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataLargeAnalogValueAll) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

// NewBACnetConstructedDataLargeAnalogValueAll factory function for _BACnetConstructedDataLargeAnalogValueAll
func NewBACnetConstructedDataLargeAnalogValueAll(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataLargeAnalogValueAll {
	_result := &_BACnetConstructedDataLargeAnalogValueAll{
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataLargeAnalogValueAll(structType interface{}) BACnetConstructedDataLargeAnalogValueAll {
	if casted, ok := structType.(BACnetConstructedDataLargeAnalogValueAll); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataLargeAnalogValueAll); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataLargeAnalogValueAll) GetTypeName() string {
	return "BACnetConstructedDataLargeAnalogValueAll"
}

func (m *_BACnetConstructedDataLargeAnalogValueAll) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataLargeAnalogValueAll) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	return lengthInBits
}

func (m *_BACnetConstructedDataLargeAnalogValueAll) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataLargeAnalogValueAllParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataLargeAnalogValueAll, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataLargeAnalogValueAll"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataLargeAnalogValueAll")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Validation
	if !(bool((1) == (2))) {
		return nil, errors.WithStack(utils.ParseValidationError{"All should never occur in context of constructed data. If it does please report"})
	}

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataLargeAnalogValueAll"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataLargeAnalogValueAll")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataLargeAnalogValueAll{
		_BACnetConstructedData: &_BACnetConstructedData{},
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataLargeAnalogValueAll) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataLargeAnalogValueAll"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataLargeAnalogValueAll")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataLargeAnalogValueAll"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataLargeAnalogValueAll")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataLargeAnalogValueAll) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
