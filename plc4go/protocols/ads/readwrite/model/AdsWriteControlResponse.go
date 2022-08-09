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

// AdsWriteControlResponse is the corresponding interface of AdsWriteControlResponse
type AdsWriteControlResponse interface {
	utils.LengthAware
	utils.Serializable
	AdsData
	// GetResult returns Result (property field)
	GetResult() ReturnCode
}

// AdsWriteControlResponseExactly can be used when we want exactly this type and not a type which fulfills AdsWriteControlResponse.
// This is useful for switch cases.
type AdsWriteControlResponseExactly interface {
	AdsWriteControlResponse
	isAdsWriteControlResponse() bool
}

// _AdsWriteControlResponse is the data-structure of this message
type _AdsWriteControlResponse struct {
	*_AdsData
	Result ReturnCode
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_AdsWriteControlResponse) GetCommandId() CommandId {
	return CommandId_ADS_WRITE_CONTROL
}

func (m *_AdsWriteControlResponse) GetResponse() bool {
	return bool(true)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_AdsWriteControlResponse) InitializeParent(parent AdsData) {}

func (m *_AdsWriteControlResponse) GetParent() AdsData {
	return m._AdsData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AdsWriteControlResponse) GetResult() ReturnCode {
	return m.Result
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewAdsWriteControlResponse factory function for _AdsWriteControlResponse
func NewAdsWriteControlResponse(result ReturnCode) *_AdsWriteControlResponse {
	_result := &_AdsWriteControlResponse{
		Result:   result,
		_AdsData: NewAdsData(),
	}
	_result._AdsData._AdsDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastAdsWriteControlResponse(structType interface{}) AdsWriteControlResponse {
	if casted, ok := structType.(AdsWriteControlResponse); ok {
		return casted
	}
	if casted, ok := structType.(*AdsWriteControlResponse); ok {
		return *casted
	}
	return nil
}

func (m *_AdsWriteControlResponse) GetTypeName() string {
	return "AdsWriteControlResponse"
}

func (m *_AdsWriteControlResponse) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_AdsWriteControlResponse) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (result)
	lengthInBits += 32

	return lengthInBits
}

func (m *_AdsWriteControlResponse) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func AdsWriteControlResponseParse(readBuffer utils.ReadBuffer, commandId CommandId, response bool) (AdsWriteControlResponse, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AdsWriteControlResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AdsWriteControlResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (result)
	if pullErr := readBuffer.PullContext("result"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for result")
	}
	_result, _resultErr := ReturnCodeParse(readBuffer)
	if _resultErr != nil {
		return nil, errors.Wrap(_resultErr, "Error parsing 'result' field of AdsWriteControlResponse")
	}
	result := _result
	if closeErr := readBuffer.CloseContext("result"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for result")
	}

	if closeErr := readBuffer.CloseContext("AdsWriteControlResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AdsWriteControlResponse")
	}

	// Create a partially initialized instance
	_child := &_AdsWriteControlResponse{
		_AdsData: &_AdsData{},
		Result:   result,
	}
	_child._AdsData._AdsDataChildRequirements = _child
	return _child, nil
}

func (m *_AdsWriteControlResponse) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("AdsWriteControlResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for AdsWriteControlResponse")
		}

		// Simple Field (result)
		if pushErr := writeBuffer.PushContext("result"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for result")
		}
		_resultErr := writeBuffer.WriteSerializable(m.GetResult())
		if popErr := writeBuffer.PopContext("result"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for result")
		}
		if _resultErr != nil {
			return errors.Wrap(_resultErr, "Error serializing 'result' field")
		}

		if popErr := writeBuffer.PopContext("AdsWriteControlResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for AdsWriteControlResponse")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_AdsWriteControlResponse) isAdsWriteControlResponse() bool {
	return true
}

func (m *_AdsWriteControlResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
