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

// StatusRequestBinaryState is the corresponding interface of StatusRequestBinaryState
type StatusRequestBinaryState interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	StatusRequest
	// GetApplication returns Application (property field)
	GetApplication() ApplicationIdContainer
}

// StatusRequestBinaryStateExactly can be used when we want exactly this type and not a type which fulfills StatusRequestBinaryState.
// This is useful for switch cases.
type StatusRequestBinaryStateExactly interface {
	StatusRequestBinaryState
	isStatusRequestBinaryState() bool
}

// _StatusRequestBinaryState is the data-structure of this message
type _StatusRequestBinaryState struct {
	*_StatusRequest
	Application ApplicationIdContainer
	// Reserved Fields
	reservedField0 *byte
	reservedField1 *byte
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_StatusRequestBinaryState) InitializeParent(parent StatusRequest, statusType byte) {
	m.StatusType = statusType
}

func (m *_StatusRequestBinaryState) GetParent() StatusRequest {
	return m._StatusRequest
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_StatusRequestBinaryState) GetApplication() ApplicationIdContainer {
	return m.Application
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewStatusRequestBinaryState factory function for _StatusRequestBinaryState
func NewStatusRequestBinaryState(application ApplicationIdContainer, statusType byte) *_StatusRequestBinaryState {
	_result := &_StatusRequestBinaryState{
		Application:    application,
		_StatusRequest: NewStatusRequest(statusType),
	}
	_result._StatusRequest._StatusRequestChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastStatusRequestBinaryState(structType any) StatusRequestBinaryState {
	if casted, ok := structType.(StatusRequestBinaryState); ok {
		return casted
	}
	if casted, ok := structType.(*StatusRequestBinaryState); ok {
		return *casted
	}
	return nil
}

func (m *_StatusRequestBinaryState) GetTypeName() string {
	return "StatusRequestBinaryState"
}

func (m *_StatusRequestBinaryState) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Reserved Field (reserved)
	lengthInBits += 8

	// Simple field (application)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 8

	return lengthInBits
}

func (m *_StatusRequestBinaryState) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func StatusRequestBinaryStateParse(ctx context.Context, theBytes []byte) (StatusRequestBinaryState, error) {
	return StatusRequestBinaryStateParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func StatusRequestBinaryStateParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (StatusRequestBinaryState, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (StatusRequestBinaryState, error) {
		return StatusRequestBinaryStateParseWithBuffer(ctx, readBuffer)
	}
}

func StatusRequestBinaryStateParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (StatusRequestBinaryState, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("StatusRequestBinaryState"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for StatusRequestBinaryState")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	reservedField0, err := ReadReservedField(ctx, "reserved", ReadByte(readBuffer, 8), byte(0x7A))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}

	application, err := ReadEnumField[ApplicationIdContainer](ctx, "application", "ApplicationIdContainer", ReadEnum(ApplicationIdContainerByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'application' field"))
	}

	reservedField1, err := ReadReservedField(ctx, "reserved", ReadByte(readBuffer, 8), byte(0x00))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}

	if closeErr := readBuffer.CloseContext("StatusRequestBinaryState"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for StatusRequestBinaryState")
	}

	// Create a partially initialized instance
	_child := &_StatusRequestBinaryState{
		_StatusRequest: &_StatusRequest{},
		Application:    application,
		reservedField0: reservedField0,
		reservedField1: reservedField1,
	}
	_child._StatusRequest._StatusRequestChildRequirements = _child
	return _child, nil
}

func (m *_StatusRequestBinaryState) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_StatusRequestBinaryState) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("StatusRequestBinaryState"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for StatusRequestBinaryState")
		}

		if err := WriteReservedField[byte](ctx, "reserved", byte(0x7A), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'reserved' field number 1")
		}

		if err := WriteSimpleEnumField[ApplicationIdContainer](ctx, "application", "ApplicationIdContainer", m.GetApplication(), WriteEnum[ApplicationIdContainer, uint8](ApplicationIdContainer.GetValue, ApplicationIdContainer.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
			return errors.Wrap(err, "Error serializing 'application' field")
		}

		if err := WriteReservedField[byte](ctx, "reserved", byte(0x00), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'reserved' field number 2")
		}

		if popErr := writeBuffer.PopContext("StatusRequestBinaryState"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for StatusRequestBinaryState")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_StatusRequestBinaryState) isStatusRequestBinaryState() bool {
	return true
}

func (m *_StatusRequestBinaryState) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
