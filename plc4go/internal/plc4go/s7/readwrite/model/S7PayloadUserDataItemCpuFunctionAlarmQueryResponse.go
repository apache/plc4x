/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Code generated by code-generation. DO NOT EDIT.

// Constant values.
const S7PayloadUserDataItemCpuFunctionAlarmQueryResponse_FUNCTIONID uint8 = 0x00
const S7PayloadUserDataItemCpuFunctionAlarmQueryResponse_NUMBERMESSAGEOBJ uint8 = 0x01

// The data-structure of this message
type S7PayloadUserDataItemCpuFunctionAlarmQueryResponse struct {
	*S7PayloadUserDataItem
	PudicfReturnCode    DataTransportErrorCode
	PudicftransportSize DataTransportSize
}

// The corresponding interface
type IS7PayloadUserDataItemCpuFunctionAlarmQueryResponse interface {
	IS7PayloadUserDataItem
	// GetPudicfReturnCode returns PudicfReturnCode (property field)
	GetPudicfReturnCode() DataTransportErrorCode
	// GetPudicftransportSize returns PudicftransportSize (property field)
	GetPudicftransportSize() DataTransportSize
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////
func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetCpuFunctionType() uint8 {
	return 0x08
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetCpuSubfunction() uint8 {
	return 0x13
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetDataLength() uint16 {
	return 0
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) InitializeParent(parent *S7PayloadUserDataItem, returnCode DataTransportErrorCode, transportSize DataTransportSize) {
	m.S7PayloadUserDataItem.ReturnCode = returnCode
	m.S7PayloadUserDataItem.TransportSize = transportSize
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetParent() *S7PayloadUserDataItem {
	return m.S7PayloadUserDataItem
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////
func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetPudicfReturnCode() DataTransportErrorCode {
	return m.PudicfReturnCode
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetPudicftransportSize() DataTransportSize {
	return m.PudicftransportSize
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for const fields.
///////////////////////
func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetFunctionId() uint8 {
	return S7PayloadUserDataItemCpuFunctionAlarmQueryResponse_FUNCTIONID
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetNumberMessageObj() uint8 {
	return S7PayloadUserDataItemCpuFunctionAlarmQueryResponse_NUMBERMESSAGEOBJ
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewS7PayloadUserDataItemCpuFunctionAlarmQueryResponse factory function for S7PayloadUserDataItemCpuFunctionAlarmQueryResponse
func NewS7PayloadUserDataItemCpuFunctionAlarmQueryResponse(pudicfReturnCode DataTransportErrorCode, pudicftransportSize DataTransportSize, returnCode DataTransportErrorCode, transportSize DataTransportSize) *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse {
	_result := &S7PayloadUserDataItemCpuFunctionAlarmQueryResponse{
		PudicfReturnCode:      pudicfReturnCode,
		PudicftransportSize:   pudicftransportSize,
		S7PayloadUserDataItem: NewS7PayloadUserDataItem(returnCode, transportSize),
	}
	_result.Child = _result
	return _result
}

func CastS7PayloadUserDataItemCpuFunctionAlarmQueryResponse(structType interface{}) *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse {
	if casted, ok := structType.(S7PayloadUserDataItemCpuFunctionAlarmQueryResponse); ok {
		return &casted
	}
	if casted, ok := structType.(*S7PayloadUserDataItemCpuFunctionAlarmQueryResponse); ok {
		return casted
	}
	if casted, ok := structType.(S7PayloadUserDataItem); ok {
		return CastS7PayloadUserDataItemCpuFunctionAlarmQueryResponse(casted.Child)
	}
	if casted, ok := structType.(*S7PayloadUserDataItem); ok {
		return CastS7PayloadUserDataItemCpuFunctionAlarmQueryResponse(casted.Child)
	}
	return nil
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetTypeName() string {
	return "S7PayloadUserDataItemCpuFunctionAlarmQueryResponse"
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Const Field (functionId)
	lengthInBits += 8

	// Const Field (numberMessageObj)
	lengthInBits += 8

	// Simple field (pudicfReturnCode)
	lengthInBits += 8

	// Simple field (pudicftransportSize)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 8

	return lengthInBits
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func S7PayloadUserDataItemCpuFunctionAlarmQueryResponseParse(readBuffer utils.ReadBuffer, cpuFunctionType uint8, cpuSubfunction uint8) (*S7PayloadUserDataItemCpuFunctionAlarmQueryResponse, error) {
	if pullErr := readBuffer.PullContext("S7PayloadUserDataItemCpuFunctionAlarmQueryResponse"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	// Const Field (functionId)
	functionId, _functionIdErr := readBuffer.ReadUint8("functionId", 8)
	if _functionIdErr != nil {
		return nil, errors.Wrap(_functionIdErr, "Error parsing 'functionId' field")
	}
	if functionId != S7PayloadUserDataItemCpuFunctionAlarmQueryResponse_FUNCTIONID {
		return nil, errors.New("Expected constant value " + fmt.Sprintf("%d", S7PayloadUserDataItemCpuFunctionAlarmQueryResponse_FUNCTIONID) + " but got " + fmt.Sprintf("%d", functionId))
	}

	// Const Field (numberMessageObj)
	numberMessageObj, _numberMessageObjErr := readBuffer.ReadUint8("numberMessageObj", 8)
	if _numberMessageObjErr != nil {
		return nil, errors.Wrap(_numberMessageObjErr, "Error parsing 'numberMessageObj' field")
	}
	if numberMessageObj != S7PayloadUserDataItemCpuFunctionAlarmQueryResponse_NUMBERMESSAGEOBJ {
		return nil, errors.New("Expected constant value " + fmt.Sprintf("%d", S7PayloadUserDataItemCpuFunctionAlarmQueryResponse_NUMBERMESSAGEOBJ) + " but got " + fmt.Sprintf("%d", numberMessageObj))
	}

	// Simple Field (pudicfReturnCode)
	if pullErr := readBuffer.PullContext("pudicfReturnCode"); pullErr != nil {
		return nil, pullErr
	}
	_pudicfReturnCode, _pudicfReturnCodeErr := DataTransportErrorCodeParse(readBuffer)
	if _pudicfReturnCodeErr != nil {
		return nil, errors.Wrap(_pudicfReturnCodeErr, "Error parsing 'pudicfReturnCode' field")
	}
	pudicfReturnCode := _pudicfReturnCode
	if closeErr := readBuffer.CloseContext("pudicfReturnCode"); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (pudicftransportSize)
	if pullErr := readBuffer.PullContext("pudicftransportSize"); pullErr != nil {
		return nil, pullErr
	}
	_pudicftransportSize, _pudicftransportSizeErr := DataTransportSizeParse(readBuffer)
	if _pudicftransportSizeErr != nil {
		return nil, errors.Wrap(_pudicftransportSizeErr, "Error parsing 'pudicftransportSize' field")
	}
	pudicftransportSize := _pudicftransportSize
	if closeErr := readBuffer.CloseContext("pudicftransportSize"); closeErr != nil {
		return nil, closeErr
	}

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := readBuffer.ReadUint8("reserved", 8)
		if _err != nil {
			return nil, errors.Wrap(_err, "Error parsing 'reserved' field")
		}
		if reserved != uint8(0x00) {
			log.Info().Fields(map[string]interface{}{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Msg("Got unexpected response.")
		}
	}

	if closeErr := readBuffer.CloseContext("S7PayloadUserDataItemCpuFunctionAlarmQueryResponse"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &S7PayloadUserDataItemCpuFunctionAlarmQueryResponse{
		PudicfReturnCode:      pudicfReturnCode,
		PudicftransportSize:   pudicftransportSize,
		S7PayloadUserDataItem: &S7PayloadUserDataItem{},
	}
	_child.S7PayloadUserDataItem.Child = _child
	return _child, nil
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("S7PayloadUserDataItemCpuFunctionAlarmQueryResponse"); pushErr != nil {
			return pushErr
		}

		// Const Field (functionId)
		_functionIdErr := writeBuffer.WriteUint8("functionId", 8, 0x00)
		if _functionIdErr != nil {
			return errors.Wrap(_functionIdErr, "Error serializing 'functionId' field")
		}

		// Const Field (numberMessageObj)
		_numberMessageObjErr := writeBuffer.WriteUint8("numberMessageObj", 8, 0x01)
		if _numberMessageObjErr != nil {
			return errors.Wrap(_numberMessageObjErr, "Error serializing 'numberMessageObj' field")
		}

		// Simple Field (pudicfReturnCode)
		if pushErr := writeBuffer.PushContext("pudicfReturnCode"); pushErr != nil {
			return pushErr
		}
		_pudicfReturnCodeErr := m.PudicfReturnCode.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("pudicfReturnCode"); popErr != nil {
			return popErr
		}
		if _pudicfReturnCodeErr != nil {
			return errors.Wrap(_pudicfReturnCodeErr, "Error serializing 'pudicfReturnCode' field")
		}

		// Simple Field (pudicftransportSize)
		if pushErr := writeBuffer.PushContext("pudicftransportSize"); pushErr != nil {
			return pushErr
		}
		_pudicftransportSizeErr := m.PudicftransportSize.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("pudicftransportSize"); popErr != nil {
			return popErr
		}
		if _pudicftransportSizeErr != nil {
			return errors.Wrap(_pudicftransportSizeErr, "Error serializing 'pudicftransportSize' field")
		}

		// Reserved Field (reserved)
		{
			_err := writeBuffer.WriteUint8("reserved", 8, uint8(0x00))
			if _err != nil {
				return errors.Wrap(_err, "Error serializing 'reserved' field")
			}
		}

		if popErr := writeBuffer.PopContext("S7PayloadUserDataItemCpuFunctionAlarmQueryResponse"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmQueryResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
