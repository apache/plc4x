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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// The data-structure of this message
type S7PayloadUserDataItemCpuFunctionAlarmAck struct {
	*S7PayloadUserDataItem
	FunctionId     uint8
	MessageObjects []*AlarmMessageObjectAckType
}

// The corresponding interface
type IS7PayloadUserDataItemCpuFunctionAlarmAck interface {
	IS7PayloadUserDataItem
	// GetFunctionId returns FunctionId (property field)
	GetFunctionId() uint8
	// GetMessageObjects returns MessageObjects (property field)
	GetMessageObjects() []*AlarmMessageObjectAckType
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
func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetCpuFunctionType() uint8 {
	return 0x04
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetCpuSubfunction() uint8 {
	return 0x0b
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetDataLength() uint16 {
	return 0
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) InitializeParent(parent *S7PayloadUserDataItem, returnCode DataTransportErrorCode, transportSize DataTransportSize) {
	m.S7PayloadUserDataItem.ReturnCode = returnCode
	m.S7PayloadUserDataItem.TransportSize = transportSize
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetParent() *S7PayloadUserDataItem {
	return m.S7PayloadUserDataItem
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////
func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetFunctionId() uint8 {
	return m.FunctionId
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetMessageObjects() []*AlarmMessageObjectAckType {
	return m.MessageObjects
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewS7PayloadUserDataItemCpuFunctionAlarmAck factory function for S7PayloadUserDataItemCpuFunctionAlarmAck
func NewS7PayloadUserDataItemCpuFunctionAlarmAck(functionId uint8, messageObjects []*AlarmMessageObjectAckType, returnCode DataTransportErrorCode, transportSize DataTransportSize) *S7PayloadUserDataItemCpuFunctionAlarmAck {
	_result := &S7PayloadUserDataItemCpuFunctionAlarmAck{
		FunctionId:            functionId,
		MessageObjects:        messageObjects,
		S7PayloadUserDataItem: NewS7PayloadUserDataItem(returnCode, transportSize),
	}
	_result.Child = _result
	return _result
}

func CastS7PayloadUserDataItemCpuFunctionAlarmAck(structType interface{}) *S7PayloadUserDataItemCpuFunctionAlarmAck {
	if casted, ok := structType.(S7PayloadUserDataItemCpuFunctionAlarmAck); ok {
		return &casted
	}
	if casted, ok := structType.(*S7PayloadUserDataItemCpuFunctionAlarmAck); ok {
		return casted
	}
	if casted, ok := structType.(S7PayloadUserDataItem); ok {
		return CastS7PayloadUserDataItemCpuFunctionAlarmAck(casted.Child)
	}
	if casted, ok := structType.(*S7PayloadUserDataItem); ok {
		return CastS7PayloadUserDataItemCpuFunctionAlarmAck(casted.Child)
	}
	return nil
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetTypeName() string {
	return "S7PayloadUserDataItemCpuFunctionAlarmAck"
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (functionId)
	lengthInBits += 8

	// Implicit Field (numberOfObjects)
	lengthInBits += 8

	// Array field
	if len(m.MessageObjects) > 0 {
		for i, element := range m.MessageObjects {
			last := i == len(m.MessageObjects)-1
			lengthInBits += element.GetLengthInBitsConditional(last)
		}
	}

	return lengthInBits
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func S7PayloadUserDataItemCpuFunctionAlarmAckParse(readBuffer utils.ReadBuffer, cpuFunctionType uint8, cpuSubfunction uint8) (*S7PayloadUserDataItemCpuFunctionAlarmAck, error) {
	if pullErr := readBuffer.PullContext("S7PayloadUserDataItemCpuFunctionAlarmAck"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	// Simple Field (functionId)
	_functionId, _functionIdErr := readBuffer.ReadUint8("functionId", 8)
	if _functionIdErr != nil {
		return nil, errors.Wrap(_functionIdErr, "Error parsing 'functionId' field")
	}
	functionId := _functionId

	// Implicit Field (numberOfObjects) (Used for parsing, but its value is not stored as it's implicitly given by the objects content)
	numberOfObjects, _numberOfObjectsErr := readBuffer.ReadUint8("numberOfObjects", 8)
	_ = numberOfObjects
	if _numberOfObjectsErr != nil {
		return nil, errors.Wrap(_numberOfObjectsErr, "Error parsing 'numberOfObjects' field")
	}

	// Array field (messageObjects)
	if pullErr := readBuffer.PullContext("messageObjects", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, pullErr
	}
	// Count array
	messageObjects := make([]*AlarmMessageObjectAckType, numberOfObjects)
	{
		for curItem := uint16(0); curItem < uint16(numberOfObjects); curItem++ {
			_item, _err := AlarmMessageObjectAckTypeParse(readBuffer)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'messageObjects' field")
			}
			messageObjects[curItem] = CastAlarmMessageObjectAckType(_item)
		}
	}
	if closeErr := readBuffer.CloseContext("messageObjects", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, closeErr
	}

	if closeErr := readBuffer.CloseContext("S7PayloadUserDataItemCpuFunctionAlarmAck"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &S7PayloadUserDataItemCpuFunctionAlarmAck{
		FunctionId:            functionId,
		MessageObjects:        messageObjects,
		S7PayloadUserDataItem: &S7PayloadUserDataItem{},
	}
	_child.S7PayloadUserDataItem.Child = _child
	return _child, nil
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("S7PayloadUserDataItemCpuFunctionAlarmAck"); pushErr != nil {
			return pushErr
		}

		// Simple Field (functionId)
		functionId := uint8(m.FunctionId)
		_functionIdErr := writeBuffer.WriteUint8("functionId", 8, (functionId))
		if _functionIdErr != nil {
			return errors.Wrap(_functionIdErr, "Error serializing 'functionId' field")
		}

		// Implicit Field (numberOfObjects) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		numberOfObjects := uint8(uint8(len(m.GetMessageObjects())))
		_numberOfObjectsErr := writeBuffer.WriteUint8("numberOfObjects", 8, (numberOfObjects))
		if _numberOfObjectsErr != nil {
			return errors.Wrap(_numberOfObjectsErr, "Error serializing 'numberOfObjects' field")
		}

		// Array Field (messageObjects)
		if m.MessageObjects != nil {
			if pushErr := writeBuffer.PushContext("messageObjects", utils.WithRenderAsList(true)); pushErr != nil {
				return pushErr
			}
			for _, _element := range m.MessageObjects {
				_elementErr := _element.Serialize(writeBuffer)
				if _elementErr != nil {
					return errors.Wrap(_elementErr, "Error serializing 'messageObjects' field")
				}
			}
			if popErr := writeBuffer.PopContext("messageObjects", utils.WithRenderAsList(true)); popErr != nil {
				return popErr
			}
		}

		if popErr := writeBuffer.PopContext("S7PayloadUserDataItemCpuFunctionAlarmAck"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *S7PayloadUserDataItemCpuFunctionAlarmAck) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
