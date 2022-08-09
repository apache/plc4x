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

// S7PayloadUserDataItemCpuFunctionReadSzlRequest is the corresponding interface of S7PayloadUserDataItemCpuFunctionReadSzlRequest
type S7PayloadUserDataItemCpuFunctionReadSzlRequest interface {
	utils.LengthAware
	utils.Serializable
	S7PayloadUserDataItem
	// GetSzlId returns SzlId (property field)
	GetSzlId() SzlId
	// GetSzlIndex returns SzlIndex (property field)
	GetSzlIndex() uint16
}

// S7PayloadUserDataItemCpuFunctionReadSzlRequestExactly can be used when we want exactly this type and not a type which fulfills S7PayloadUserDataItemCpuFunctionReadSzlRequest.
// This is useful for switch cases.
type S7PayloadUserDataItemCpuFunctionReadSzlRequestExactly interface {
	S7PayloadUserDataItemCpuFunctionReadSzlRequest
	isS7PayloadUserDataItemCpuFunctionReadSzlRequest() bool
}

// _S7PayloadUserDataItemCpuFunctionReadSzlRequest is the data-structure of this message
type _S7PayloadUserDataItemCpuFunctionReadSzlRequest struct {
	*_S7PayloadUserDataItem
	SzlId    SzlId
	SzlIndex uint16
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetCpuFunctionType() uint8 {
	return 0x04
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetCpuSubfunction() uint8 {
	return 0x01
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetDataLength() uint16 {
	return 0
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) InitializeParent(parent S7PayloadUserDataItem, returnCode DataTransportErrorCode, transportSize DataTransportSize) {
	m.ReturnCode = returnCode
	m.TransportSize = transportSize
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetParent() S7PayloadUserDataItem {
	return m._S7PayloadUserDataItem
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetSzlId() SzlId {
	return m.SzlId
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetSzlIndex() uint16 {
	return m.SzlIndex
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewS7PayloadUserDataItemCpuFunctionReadSzlRequest factory function for _S7PayloadUserDataItemCpuFunctionReadSzlRequest
func NewS7PayloadUserDataItemCpuFunctionReadSzlRequest(szlId SzlId, szlIndex uint16, returnCode DataTransportErrorCode, transportSize DataTransportSize) *_S7PayloadUserDataItemCpuFunctionReadSzlRequest {
	_result := &_S7PayloadUserDataItemCpuFunctionReadSzlRequest{
		SzlId:                  szlId,
		SzlIndex:               szlIndex,
		_S7PayloadUserDataItem: NewS7PayloadUserDataItem(returnCode, transportSize),
	}
	_result._S7PayloadUserDataItem._S7PayloadUserDataItemChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastS7PayloadUserDataItemCpuFunctionReadSzlRequest(structType interface{}) S7PayloadUserDataItemCpuFunctionReadSzlRequest {
	if casted, ok := structType.(S7PayloadUserDataItemCpuFunctionReadSzlRequest); ok {
		return casted
	}
	if casted, ok := structType.(*S7PayloadUserDataItemCpuFunctionReadSzlRequest); ok {
		return *casted
	}
	return nil
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetTypeName() string {
	return "S7PayloadUserDataItemCpuFunctionReadSzlRequest"
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (szlId)
	lengthInBits += m.SzlId.GetLengthInBits()

	// Simple field (szlIndex)
	lengthInBits += 16

	return lengthInBits
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func S7PayloadUserDataItemCpuFunctionReadSzlRequestParse(readBuffer utils.ReadBuffer, cpuFunctionType uint8, cpuSubfunction uint8) (S7PayloadUserDataItemCpuFunctionReadSzlRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("S7PayloadUserDataItemCpuFunctionReadSzlRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for S7PayloadUserDataItemCpuFunctionReadSzlRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (szlId)
	if pullErr := readBuffer.PullContext("szlId"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for szlId")
	}
	_szlId, _szlIdErr := SzlIdParse(readBuffer)
	if _szlIdErr != nil {
		return nil, errors.Wrap(_szlIdErr, "Error parsing 'szlId' field of S7PayloadUserDataItemCpuFunctionReadSzlRequest")
	}
	szlId := _szlId.(SzlId)
	if closeErr := readBuffer.CloseContext("szlId"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for szlId")
	}

	// Simple Field (szlIndex)
	_szlIndex, _szlIndexErr := readBuffer.ReadUint16("szlIndex", 16)
	if _szlIndexErr != nil {
		return nil, errors.Wrap(_szlIndexErr, "Error parsing 'szlIndex' field of S7PayloadUserDataItemCpuFunctionReadSzlRequest")
	}
	szlIndex := _szlIndex

	if closeErr := readBuffer.CloseContext("S7PayloadUserDataItemCpuFunctionReadSzlRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for S7PayloadUserDataItemCpuFunctionReadSzlRequest")
	}

	// Create a partially initialized instance
	_child := &_S7PayloadUserDataItemCpuFunctionReadSzlRequest{
		_S7PayloadUserDataItem: &_S7PayloadUserDataItem{},
		SzlId:                  szlId,
		SzlIndex:               szlIndex,
	}
	_child._S7PayloadUserDataItem._S7PayloadUserDataItemChildRequirements = _child
	return _child, nil
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("S7PayloadUserDataItemCpuFunctionReadSzlRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for S7PayloadUserDataItemCpuFunctionReadSzlRequest")
		}

		// Simple Field (szlId)
		if pushErr := writeBuffer.PushContext("szlId"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for szlId")
		}
		_szlIdErr := writeBuffer.WriteSerializable(m.GetSzlId())
		if popErr := writeBuffer.PopContext("szlId"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for szlId")
		}
		if _szlIdErr != nil {
			return errors.Wrap(_szlIdErr, "Error serializing 'szlId' field")
		}

		// Simple Field (szlIndex)
		szlIndex := uint16(m.GetSzlIndex())
		_szlIndexErr := writeBuffer.WriteUint16("szlIndex", 16, (szlIndex))
		if _szlIndexErr != nil {
			return errors.Wrap(_szlIndexErr, "Error serializing 'szlIndex' field")
		}

		if popErr := writeBuffer.PopContext("S7PayloadUserDataItemCpuFunctionReadSzlRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for S7PayloadUserDataItemCpuFunctionReadSzlRequest")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) isS7PayloadUserDataItemCpuFunctionReadSzlRequest() bool {
	return true
}

func (m *_S7PayloadUserDataItemCpuFunctionReadSzlRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
