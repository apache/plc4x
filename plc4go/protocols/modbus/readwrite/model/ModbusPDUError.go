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

// ModbusPDUError is the corresponding interface of ModbusPDUError
type ModbusPDUError interface {
	utils.LengthAware
	utils.Serializable
	ModbusPDU
	// GetExceptionCode returns ExceptionCode (property field)
	GetExceptionCode() ModbusErrorCode
}

// ModbusPDUErrorExactly can be used when we want exactly this type and not a type which fulfills ModbusPDUError.
// This is useful for switch cases.
type ModbusPDUErrorExactly interface {
	ModbusPDUError
	isModbusPDUError() bool
}

// _ModbusPDUError is the data-structure of this message
type _ModbusPDUError struct {
	*_ModbusPDU
	ExceptionCode ModbusErrorCode
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ModbusPDUError) GetErrorFlag() bool {
	return bool(true)
}

func (m *_ModbusPDUError) GetFunctionFlag() uint8 {
	return 0
}

func (m *_ModbusPDUError) GetResponse() bool {
	return false
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ModbusPDUError) InitializeParent(parent ModbusPDU) {}

func (m *_ModbusPDUError) GetParent() ModbusPDU {
	return m._ModbusPDU
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ModbusPDUError) GetExceptionCode() ModbusErrorCode {
	return m.ExceptionCode
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewModbusPDUError factory function for _ModbusPDUError
func NewModbusPDUError(exceptionCode ModbusErrorCode) *_ModbusPDUError {
	_result := &_ModbusPDUError{
		ExceptionCode: exceptionCode,
		_ModbusPDU:    NewModbusPDU(),
	}
	_result._ModbusPDU._ModbusPDUChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastModbusPDUError(structType interface{}) ModbusPDUError {
	if casted, ok := structType.(ModbusPDUError); ok {
		return casted
	}
	if casted, ok := structType.(*ModbusPDUError); ok {
		return *casted
	}
	return nil
}

func (m *_ModbusPDUError) GetTypeName() string {
	return "ModbusPDUError"
}

func (m *_ModbusPDUError) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_ModbusPDUError) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (exceptionCode)
	lengthInBits += 8

	return lengthInBits
}

func (m *_ModbusPDUError) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func ModbusPDUErrorParse(theBytes []byte, response bool) (ModbusPDUError, error) {
	return ModbusPDUErrorParseWithBuffer(utils.NewReadBufferByteBased(theBytes), response)
}

func ModbusPDUErrorParseWithBuffer(readBuffer utils.ReadBuffer, response bool) (ModbusPDUError, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ModbusPDUError"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ModbusPDUError")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (exceptionCode)
	if pullErr := readBuffer.PullContext("exceptionCode"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for exceptionCode")
	}
	_exceptionCode, _exceptionCodeErr := ModbusErrorCodeParseWithBuffer(readBuffer)
	if _exceptionCodeErr != nil {
		return nil, errors.Wrap(_exceptionCodeErr, "Error parsing 'exceptionCode' field of ModbusPDUError")
	}
	exceptionCode := _exceptionCode
	if closeErr := readBuffer.CloseContext("exceptionCode"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for exceptionCode")
	}

	if closeErr := readBuffer.CloseContext("ModbusPDUError"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ModbusPDUError")
	}

	// Create a partially initialized instance
	_child := &_ModbusPDUError{
		_ModbusPDU:    &_ModbusPDU{},
		ExceptionCode: exceptionCode,
	}
	_child._ModbusPDU._ModbusPDUChildRequirements = _child
	return _child, nil
}

func (m *_ModbusPDUError) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ModbusPDUError) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ModbusPDUError"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ModbusPDUError")
		}

		// Simple Field (exceptionCode)
		if pushErr := writeBuffer.PushContext("exceptionCode"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for exceptionCode")
		}
		_exceptionCodeErr := writeBuffer.WriteSerializable(m.GetExceptionCode())
		if popErr := writeBuffer.PopContext("exceptionCode"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for exceptionCode")
		}
		if _exceptionCodeErr != nil {
			return errors.Wrap(_exceptionCodeErr, "Error serializing 'exceptionCode' field")
		}

		if popErr := writeBuffer.PopContext("ModbusPDUError"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ModbusPDUError")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_ModbusPDUError) isModbusPDUError() bool {
	return true
}

func (m *_ModbusPDUError) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
