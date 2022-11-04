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

// BACnetFaultParameterFaultExtendedParametersEntryDouble is the corresponding interface of BACnetFaultParameterFaultExtendedParametersEntryDouble
type BACnetFaultParameterFaultExtendedParametersEntryDouble interface {
	utils.LengthAware
	utils.Serializable
	BACnetFaultParameterFaultExtendedParametersEntry
	// GetDoubleValue returns DoubleValue (property field)
	GetDoubleValue() BACnetApplicationTagDouble
}

// BACnetFaultParameterFaultExtendedParametersEntryDoubleExactly can be used when we want exactly this type and not a type which fulfills BACnetFaultParameterFaultExtendedParametersEntryDouble.
// This is useful for switch cases.
type BACnetFaultParameterFaultExtendedParametersEntryDoubleExactly interface {
	BACnetFaultParameterFaultExtendedParametersEntryDouble
	isBACnetFaultParameterFaultExtendedParametersEntryDouble() bool
}

// _BACnetFaultParameterFaultExtendedParametersEntryDouble is the data-structure of this message
type _BACnetFaultParameterFaultExtendedParametersEntryDouble struct {
	*_BACnetFaultParameterFaultExtendedParametersEntry
	DoubleValue BACnetApplicationTagDouble
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) InitializeParent(parent BACnetFaultParameterFaultExtendedParametersEntry, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) GetParent() BACnetFaultParameterFaultExtendedParametersEntry {
	return m._BACnetFaultParameterFaultExtendedParametersEntry
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) GetDoubleValue() BACnetApplicationTagDouble {
	return m.DoubleValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetFaultParameterFaultExtendedParametersEntryDouble factory function for _BACnetFaultParameterFaultExtendedParametersEntryDouble
func NewBACnetFaultParameterFaultExtendedParametersEntryDouble(doubleValue BACnetApplicationTagDouble, peekedTagHeader BACnetTagHeader) *_BACnetFaultParameterFaultExtendedParametersEntryDouble {
	_result := &_BACnetFaultParameterFaultExtendedParametersEntryDouble{
		DoubleValue: doubleValue,
		_BACnetFaultParameterFaultExtendedParametersEntry: NewBACnetFaultParameterFaultExtendedParametersEntry(peekedTagHeader),
	}
	_result._BACnetFaultParameterFaultExtendedParametersEntry._BACnetFaultParameterFaultExtendedParametersEntryChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetFaultParameterFaultExtendedParametersEntryDouble(structType interface{}) BACnetFaultParameterFaultExtendedParametersEntryDouble {
	if casted, ok := structType.(BACnetFaultParameterFaultExtendedParametersEntryDouble); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetFaultParameterFaultExtendedParametersEntryDouble); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) GetTypeName() string {
	return "BACnetFaultParameterFaultExtendedParametersEntryDouble"
}

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (doubleValue)
	lengthInBits += m.DoubleValue.GetLengthInBits()

	return lengthInBits
}

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetFaultParameterFaultExtendedParametersEntryDoubleParse(theBytes []byte) (BACnetFaultParameterFaultExtendedParametersEntryDouble, error) {
	return BACnetFaultParameterFaultExtendedParametersEntryDoubleParseWithBuffer(utils.NewReadBufferByteBased(theBytes))
}

func BACnetFaultParameterFaultExtendedParametersEntryDoubleParseWithBuffer(readBuffer utils.ReadBuffer) (BACnetFaultParameterFaultExtendedParametersEntryDouble, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetFaultParameterFaultExtendedParametersEntryDouble"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetFaultParameterFaultExtendedParametersEntryDouble")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (doubleValue)
	if pullErr := readBuffer.PullContext("doubleValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for doubleValue")
	}
	_doubleValue, _doubleValueErr := BACnetApplicationTagParseWithBuffer(readBuffer)
	if _doubleValueErr != nil {
		return nil, errors.Wrap(_doubleValueErr, "Error parsing 'doubleValue' field of BACnetFaultParameterFaultExtendedParametersEntryDouble")
	}
	doubleValue := _doubleValue.(BACnetApplicationTagDouble)
	if closeErr := readBuffer.CloseContext("doubleValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for doubleValue")
	}

	if closeErr := readBuffer.CloseContext("BACnetFaultParameterFaultExtendedParametersEntryDouble"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetFaultParameterFaultExtendedParametersEntryDouble")
	}

	// Create a partially initialized instance
	_child := &_BACnetFaultParameterFaultExtendedParametersEntryDouble{
		_BACnetFaultParameterFaultExtendedParametersEntry: &_BACnetFaultParameterFaultExtendedParametersEntry{},
		DoubleValue: doubleValue,
	}
	_child._BACnetFaultParameterFaultExtendedParametersEntry._BACnetFaultParameterFaultExtendedParametersEntryChildRequirements = _child
	return _child, nil
}

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetFaultParameterFaultExtendedParametersEntryDouble"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetFaultParameterFaultExtendedParametersEntryDouble")
		}

		// Simple Field (doubleValue)
		if pushErr := writeBuffer.PushContext("doubleValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for doubleValue")
		}
		_doubleValueErr := writeBuffer.WriteSerializable(m.GetDoubleValue())
		if popErr := writeBuffer.PopContext("doubleValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for doubleValue")
		}
		if _doubleValueErr != nil {
			return errors.Wrap(_doubleValueErr, "Error serializing 'doubleValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetFaultParameterFaultExtendedParametersEntryDouble"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetFaultParameterFaultExtendedParametersEntryDouble")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) isBACnetFaultParameterFaultExtendedParametersEntryDouble() bool {
	return true
}

func (m *_BACnetFaultParameterFaultExtendedParametersEntryDouble) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
