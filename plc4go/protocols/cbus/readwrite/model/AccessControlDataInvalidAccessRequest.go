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
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// AccessControlDataInvalidAccessRequest is the corresponding interface of AccessControlDataInvalidAccessRequest
type AccessControlDataInvalidAccessRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	AccessControlData
	// GetAccessControlDirection returns AccessControlDirection (property field)
	GetAccessControlDirection() AccessControlDirection
	// GetData returns Data (property field)
	GetData() []byte
}

// AccessControlDataInvalidAccessRequestExactly can be used when we want exactly this type and not a type which fulfills AccessControlDataInvalidAccessRequest.
// This is useful for switch cases.
type AccessControlDataInvalidAccessRequestExactly interface {
	AccessControlDataInvalidAccessRequest
	isAccessControlDataInvalidAccessRequest() bool
}

// _AccessControlDataInvalidAccessRequest is the data-structure of this message
type _AccessControlDataInvalidAccessRequest struct {
	*_AccessControlData
	AccessControlDirection AccessControlDirection
	Data                   []byte
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_AccessControlDataInvalidAccessRequest) InitializeParent(parent AccessControlData, commandTypeContainer AccessControlCommandTypeContainer, networkId byte, accessPointId byte) {
	m.CommandTypeContainer = commandTypeContainer
	m.NetworkId = networkId
	m.AccessPointId = accessPointId
}

func (m *_AccessControlDataInvalidAccessRequest) GetParent() AccessControlData {
	return m._AccessControlData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AccessControlDataInvalidAccessRequest) GetAccessControlDirection() AccessControlDirection {
	return m.AccessControlDirection
}

func (m *_AccessControlDataInvalidAccessRequest) GetData() []byte {
	return m.Data
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewAccessControlDataInvalidAccessRequest factory function for _AccessControlDataInvalidAccessRequest
func NewAccessControlDataInvalidAccessRequest(accessControlDirection AccessControlDirection, data []byte, commandTypeContainer AccessControlCommandTypeContainer, networkId byte, accessPointId byte) *_AccessControlDataInvalidAccessRequest {
	_result := &_AccessControlDataInvalidAccessRequest{
		AccessControlDirection: accessControlDirection,
		Data:                   data,
		_AccessControlData:     NewAccessControlData(commandTypeContainer, networkId, accessPointId),
	}
	_result._AccessControlData._AccessControlDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastAccessControlDataInvalidAccessRequest(structType any) AccessControlDataInvalidAccessRequest {
	if casted, ok := structType.(AccessControlDataInvalidAccessRequest); ok {
		return casted
	}
	if casted, ok := structType.(*AccessControlDataInvalidAccessRequest); ok {
		return *casted
	}
	return nil
}

func (m *_AccessControlDataInvalidAccessRequest) GetTypeName() string {
	return "AccessControlDataInvalidAccessRequest"
}

func (m *_AccessControlDataInvalidAccessRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (accessControlDirection)
	lengthInBits += 8

	// Array field
	if len(m.Data) > 0 {
		lengthInBits += 8 * uint16(len(m.Data))
	}

	return lengthInBits
}

func (m *_AccessControlDataInvalidAccessRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func AccessControlDataInvalidAccessRequestParse(ctx context.Context, theBytes []byte, commandTypeContainer AccessControlCommandTypeContainer) (AccessControlDataInvalidAccessRequest, error) {
	return AccessControlDataInvalidAccessRequestParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), commandTypeContainer)
}

func AccessControlDataInvalidAccessRequestParseWithBufferProducer(commandTypeContainer AccessControlCommandTypeContainer) func(ctx context.Context, readBuffer utils.ReadBuffer) (AccessControlDataInvalidAccessRequest, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (AccessControlDataInvalidAccessRequest, error) {
		return AccessControlDataInvalidAccessRequestParseWithBuffer(ctx, readBuffer, commandTypeContainer)
	}
}

func AccessControlDataInvalidAccessRequestParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, commandTypeContainer AccessControlCommandTypeContainer) (AccessControlDataInvalidAccessRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AccessControlDataInvalidAccessRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AccessControlDataInvalidAccessRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	accessControlDirection, err := ReadEnumField[AccessControlDirection](ctx, "accessControlDirection", "AccessControlDirection", ReadEnum(AccessControlDirectionByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'accessControlDirection' field"))
	}

	data, err := readBuffer.ReadByteArray("data", int(int32(commandTypeContainer.NumBytes())-int32(int32(3))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'data' field"))
	}

	if closeErr := readBuffer.CloseContext("AccessControlDataInvalidAccessRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AccessControlDataInvalidAccessRequest")
	}

	// Create a partially initialized instance
	_child := &_AccessControlDataInvalidAccessRequest{
		_AccessControlData:     &_AccessControlData{},
		AccessControlDirection: accessControlDirection,
		Data:                   data,
	}
	_child._AccessControlData._AccessControlDataChildRequirements = _child
	return _child, nil
}

func (m *_AccessControlDataInvalidAccessRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_AccessControlDataInvalidAccessRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("AccessControlDataInvalidAccessRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for AccessControlDataInvalidAccessRequest")
		}

		if err := WriteSimpleEnumField[AccessControlDirection](ctx, "accessControlDirection", "AccessControlDirection", m.GetAccessControlDirection(), WriteEnum[AccessControlDirection, uint8](AccessControlDirection.GetValue, AccessControlDirection.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
			return errors.Wrap(err, "Error serializing 'accessControlDirection' field")
		}

		if err := WriteByteArrayField(ctx, "data", m.GetData(), WriteByteArray(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'data' field")
		}

		if popErr := writeBuffer.PopContext("AccessControlDataInvalidAccessRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for AccessControlDataInvalidAccessRequest")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_AccessControlDataInvalidAccessRequest) isAccessControlDataInvalidAccessRequest() bool {
	return true
}

func (m *_AccessControlDataInvalidAccessRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
