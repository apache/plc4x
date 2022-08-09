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

// BACnetChannelValueInteger is the corresponding interface of BACnetChannelValueInteger
type BACnetChannelValueInteger interface {
	utils.LengthAware
	utils.Serializable
	BACnetChannelValue
	// GetIntegerValue returns IntegerValue (property field)
	GetIntegerValue() BACnetApplicationTagSignedInteger
}

// BACnetChannelValueIntegerExactly can be used when we want exactly this type and not a type which fulfills BACnetChannelValueInteger.
// This is useful for switch cases.
type BACnetChannelValueIntegerExactly interface {
	BACnetChannelValueInteger
	isBACnetChannelValueInteger() bool
}

// _BACnetChannelValueInteger is the data-structure of this message
type _BACnetChannelValueInteger struct {
	*_BACnetChannelValue
	IntegerValue BACnetApplicationTagSignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetChannelValueInteger) InitializeParent(parent BACnetChannelValue, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetChannelValueInteger) GetParent() BACnetChannelValue {
	return m._BACnetChannelValue
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetChannelValueInteger) GetIntegerValue() BACnetApplicationTagSignedInteger {
	return m.IntegerValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetChannelValueInteger factory function for _BACnetChannelValueInteger
func NewBACnetChannelValueInteger(integerValue BACnetApplicationTagSignedInteger, peekedTagHeader BACnetTagHeader) *_BACnetChannelValueInteger {
	_result := &_BACnetChannelValueInteger{
		IntegerValue:        integerValue,
		_BACnetChannelValue: NewBACnetChannelValue(peekedTagHeader),
	}
	_result._BACnetChannelValue._BACnetChannelValueChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetChannelValueInteger(structType interface{}) BACnetChannelValueInteger {
	if casted, ok := structType.(BACnetChannelValueInteger); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetChannelValueInteger); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetChannelValueInteger) GetTypeName() string {
	return "BACnetChannelValueInteger"
}

func (m *_BACnetChannelValueInteger) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetChannelValueInteger) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (integerValue)
	lengthInBits += m.IntegerValue.GetLengthInBits()

	return lengthInBits
}

func (m *_BACnetChannelValueInteger) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetChannelValueIntegerParse(readBuffer utils.ReadBuffer) (BACnetChannelValueInteger, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetChannelValueInteger"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetChannelValueInteger")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (integerValue)
	if pullErr := readBuffer.PullContext("integerValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for integerValue")
	}
	_integerValue, _integerValueErr := BACnetApplicationTagParse(readBuffer)
	if _integerValueErr != nil {
		return nil, errors.Wrap(_integerValueErr, "Error parsing 'integerValue' field of BACnetChannelValueInteger")
	}
	integerValue := _integerValue.(BACnetApplicationTagSignedInteger)
	if closeErr := readBuffer.CloseContext("integerValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for integerValue")
	}

	if closeErr := readBuffer.CloseContext("BACnetChannelValueInteger"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetChannelValueInteger")
	}

	// Create a partially initialized instance
	_child := &_BACnetChannelValueInteger{
		_BACnetChannelValue: &_BACnetChannelValue{},
		IntegerValue:        integerValue,
	}
	_child._BACnetChannelValue._BACnetChannelValueChildRequirements = _child
	return _child, nil
}

func (m *_BACnetChannelValueInteger) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetChannelValueInteger"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetChannelValueInteger")
		}

		// Simple Field (integerValue)
		if pushErr := writeBuffer.PushContext("integerValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for integerValue")
		}
		_integerValueErr := writeBuffer.WriteSerializable(m.GetIntegerValue())
		if popErr := writeBuffer.PopContext("integerValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for integerValue")
		}
		if _integerValueErr != nil {
			return errors.Wrap(_integerValueErr, "Error serializing 'integerValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetChannelValueInteger"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetChannelValueInteger")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetChannelValueInteger) isBACnetChannelValueInteger() bool {
	return true
}

func (m *_BACnetChannelValueInteger) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
