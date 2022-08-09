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

// ParameterValueRaw is the corresponding interface of ParameterValueRaw
type ParameterValueRaw interface {
	utils.LengthAware
	utils.Serializable
	ParameterValue
	// GetData returns Data (property field)
	GetData() []byte
}

// ParameterValueRawExactly can be used when we want exactly this type and not a type which fulfills ParameterValueRaw.
// This is useful for switch cases.
type ParameterValueRawExactly interface {
	ParameterValueRaw
	isParameterValueRaw() bool
}

// _ParameterValueRaw is the data-structure of this message
type _ParameterValueRaw struct {
	*_ParameterValue
	Data []byte
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ParameterValueRaw) GetParameterType() ParameterType {
	return 0
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ParameterValueRaw) InitializeParent(parent ParameterValue) {}

func (m *_ParameterValueRaw) GetParent() ParameterValue {
	return m._ParameterValue
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ParameterValueRaw) GetData() []byte {
	return m.Data
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewParameterValueRaw factory function for _ParameterValueRaw
func NewParameterValueRaw(data []byte, numBytes uint8) *_ParameterValueRaw {
	_result := &_ParameterValueRaw{
		Data:            data,
		_ParameterValue: NewParameterValue(numBytes),
	}
	_result._ParameterValue._ParameterValueChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastParameterValueRaw(structType interface{}) ParameterValueRaw {
	if casted, ok := structType.(ParameterValueRaw); ok {
		return casted
	}
	if casted, ok := structType.(*ParameterValueRaw); ok {
		return *casted
	}
	return nil
}

func (m *_ParameterValueRaw) GetTypeName() string {
	return "ParameterValueRaw"
}

func (m *_ParameterValueRaw) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_ParameterValueRaw) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Array field
	if len(m.Data) > 0 {
		lengthInBits += 8 * uint16(len(m.Data))
	}

	return lengthInBits
}

func (m *_ParameterValueRaw) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func ParameterValueRawParse(readBuffer utils.ReadBuffer, parameterType ParameterType, numBytes uint8) (ParameterValueRaw, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ParameterValueRaw"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ParameterValueRaw")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos
	// Byte Array field (data)
	numberOfBytesdata := int(numBytes)
	data, _readArrayErr := readBuffer.ReadByteArray("data", numberOfBytesdata)
	if _readArrayErr != nil {
		return nil, errors.Wrap(_readArrayErr, "Error parsing 'data' field of ParameterValueRaw")
	}

	if closeErr := readBuffer.CloseContext("ParameterValueRaw"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ParameterValueRaw")
	}

	// Create a partially initialized instance
	_child := &_ParameterValueRaw{
		_ParameterValue: &_ParameterValue{
			NumBytes: numBytes,
		},
		Data: data,
	}
	_child._ParameterValue._ParameterValueChildRequirements = _child
	return _child, nil
}

func (m *_ParameterValueRaw) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ParameterValueRaw"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ParameterValueRaw")
		}

		// Array Field (data)
		// Byte Array field (data)
		if err := writeBuffer.WriteByteArray("data", m.GetData()); err != nil {
			return errors.Wrap(err, "Error serializing 'data' field")
		}

		if popErr := writeBuffer.PopContext("ParameterValueRaw"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ParameterValueRaw")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_ParameterValueRaw) isParameterValueRaw() bool {
	return true
}

func (m *_ParameterValueRaw) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
