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

// BACnetTagPayloadReal is the corresponding interface of BACnetTagPayloadReal
type BACnetTagPayloadReal interface {
	utils.LengthAware
	utils.Serializable
	// GetValue returns Value (property field)
	GetValue() float32
}

// BACnetTagPayloadRealExactly can be used when we want exactly this type and not a type which fulfills BACnetTagPayloadReal.
// This is useful for switch cases.
type BACnetTagPayloadRealExactly interface {
	BACnetTagPayloadReal
	isBACnetTagPayloadReal() bool
}

// _BACnetTagPayloadReal is the data-structure of this message
type _BACnetTagPayloadReal struct {
	Value float32
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetTagPayloadReal) GetValue() float32 {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetTagPayloadReal factory function for _BACnetTagPayloadReal
func NewBACnetTagPayloadReal(value float32) *_BACnetTagPayloadReal {
	return &_BACnetTagPayloadReal{Value: value}
}

// Deprecated: use the interface for direct cast
func CastBACnetTagPayloadReal(structType interface{}) BACnetTagPayloadReal {
	if casted, ok := structType.(BACnetTagPayloadReal); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetTagPayloadReal); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetTagPayloadReal) GetTypeName() string {
	return "BACnetTagPayloadReal"
}

func (m *_BACnetTagPayloadReal) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetTagPayloadReal) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Simple field (value)
	lengthInBits += 32

	return lengthInBits
}

func (m *_BACnetTagPayloadReal) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetTagPayloadRealParse(readBuffer utils.ReadBuffer) (BACnetTagPayloadReal, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetTagPayloadReal"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetTagPayloadReal")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (value)
	_value, _valueErr := readBuffer.ReadFloat32("value", 32)
	if _valueErr != nil {
		return nil, errors.Wrap(_valueErr, "Error parsing 'value' field of BACnetTagPayloadReal")
	}
	value := _value

	if closeErr := readBuffer.CloseContext("BACnetTagPayloadReal"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetTagPayloadReal")
	}

	// Create the instance
	return &_BACnetTagPayloadReal{
		Value: value,
	}, nil
}

func (m *_BACnetTagPayloadReal) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetTagPayloadReal"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetTagPayloadReal")
	}

	// Simple Field (value)
	value := float32(m.GetValue())
	_valueErr := writeBuffer.WriteFloat32("value", 32, (value))
	if _valueErr != nil {
		return errors.Wrap(_valueErr, "Error serializing 'value' field")
	}

	if popErr := writeBuffer.PopContext("BACnetTagPayloadReal"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetTagPayloadReal")
	}
	return nil
}

func (m *_BACnetTagPayloadReal) isBACnetTagPayloadReal() bool {
	return true
}

func (m *_BACnetTagPayloadReal) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
