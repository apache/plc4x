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

// AdsReadStateRequest is the corresponding interface of AdsReadStateRequest
type AdsReadStateRequest interface {
	utils.LengthAware
	utils.Serializable
	AmsPacket
}

// AdsReadStateRequestExactly can be used when we want exactly this type and not a type which fulfills AdsReadStateRequest.
// This is useful for switch cases.
type AdsReadStateRequestExactly interface {
	AdsReadStateRequest
	isAdsReadStateRequest() bool
}

// _AdsReadStateRequest is the data-structure of this message
type _AdsReadStateRequest struct {
	*_AmsPacket
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_AdsReadStateRequest) GetCommandId() CommandId {
	return CommandId_ADS_READ_STATE
}

func (m *_AdsReadStateRequest) GetResponse() bool {
	return bool(false)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_AdsReadStateRequest) InitializeParent(parent AmsPacket, targetAmsNetId AmsNetId, targetAmsPort uint16, sourceAmsNetId AmsNetId, sourceAmsPort uint16, errorCode uint32, invokeId uint32) {
	m.TargetAmsNetId = targetAmsNetId
	m.TargetAmsPort = targetAmsPort
	m.SourceAmsNetId = sourceAmsNetId
	m.SourceAmsPort = sourceAmsPort
	m.ErrorCode = errorCode
	m.InvokeId = invokeId
}

func (m *_AdsReadStateRequest) GetParent() AmsPacket {
	return m._AmsPacket
}

// NewAdsReadStateRequest factory function for _AdsReadStateRequest
func NewAdsReadStateRequest(targetAmsNetId AmsNetId, targetAmsPort uint16, sourceAmsNetId AmsNetId, sourceAmsPort uint16, errorCode uint32, invokeId uint32) *_AdsReadStateRequest {
	_result := &_AdsReadStateRequest{
		_AmsPacket: NewAmsPacket(targetAmsNetId, targetAmsPort, sourceAmsNetId, sourceAmsPort, errorCode, invokeId),
	}
	_result._AmsPacket._AmsPacketChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastAdsReadStateRequest(structType interface{}) AdsReadStateRequest {
	if casted, ok := structType.(AdsReadStateRequest); ok {
		return casted
	}
	if casted, ok := structType.(*AdsReadStateRequest); ok {
		return *casted
	}
	return nil
}

func (m *_AdsReadStateRequest) GetTypeName() string {
	return "AdsReadStateRequest"
}

func (m *_AdsReadStateRequest) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_AdsReadStateRequest) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	return lengthInBits
}

func (m *_AdsReadStateRequest) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func AdsReadStateRequestParse(theBytes []byte) (AdsReadStateRequest, error) {
	return AdsReadStateRequestParseWithBuffer(utils.NewReadBufferByteBased(theBytes))
}

func AdsReadStateRequestParseWithBuffer(readBuffer utils.ReadBuffer) (AdsReadStateRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AdsReadStateRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AdsReadStateRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("AdsReadStateRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AdsReadStateRequest")
	}

	// Create a partially initialized instance
	_child := &_AdsReadStateRequest{
		_AmsPacket: &_AmsPacket{},
	}
	_child._AmsPacket._AmsPacketChildRequirements = _child
	return _child, nil
}

func (m *_AdsReadStateRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_AdsReadStateRequest) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("AdsReadStateRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for AdsReadStateRequest")
		}

		if popErr := writeBuffer.PopContext("AdsReadStateRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for AdsReadStateRequest")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_AdsReadStateRequest) isAdsReadStateRequest() bool {
	return true
}

func (m *_AdsReadStateRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
