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

// DF1ResponseMessage is the corresponding interface of DF1ResponseMessage
type DF1ResponseMessage interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetCommandCode returns CommandCode (discriminator field)
	GetCommandCode() uint8
	// GetDestinationAddress returns DestinationAddress (property field)
	GetDestinationAddress() uint8
	// GetSourceAddress returns SourceAddress (property field)
	GetSourceAddress() uint8
	// GetStatus returns Status (property field)
	GetStatus() uint8
	// GetTransactionCounter returns TransactionCounter (property field)
	GetTransactionCounter() uint16
}

// DF1ResponseMessageExactly can be used when we want exactly this type and not a type which fulfills DF1ResponseMessage.
// This is useful for switch cases.
type DF1ResponseMessageExactly interface {
	DF1ResponseMessage
	isDF1ResponseMessage() bool
}

// _DF1ResponseMessage is the data-structure of this message
type _DF1ResponseMessage struct {
	_DF1ResponseMessageChildRequirements
	DestinationAddress uint8
	SourceAddress      uint8
	Status             uint8
	TransactionCounter uint16

	// Arguments.
	PayloadLength uint16
	// Reserved Fields
	reservedField0 *uint8
	reservedField1 *uint8
}

type _DF1ResponseMessageChildRequirements interface {
	utils.Serializable
	GetLengthInBits(ctx context.Context) uint16
	GetCommandCode() uint8
}

type DF1ResponseMessageParent interface {
	SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child DF1ResponseMessage, serializeChildFunction func() error) error
	GetTypeName() string
}

type DF1ResponseMessageChild interface {
	utils.Serializable
	InitializeParent(parent DF1ResponseMessage, destinationAddress uint8, sourceAddress uint8, status uint8, transactionCounter uint16)
	GetParent() *DF1ResponseMessage

	GetTypeName() string
	DF1ResponseMessage
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_DF1ResponseMessage) GetDestinationAddress() uint8 {
	return m.DestinationAddress
}

func (m *_DF1ResponseMessage) GetSourceAddress() uint8 {
	return m.SourceAddress
}

func (m *_DF1ResponseMessage) GetStatus() uint8 {
	return m.Status
}

func (m *_DF1ResponseMessage) GetTransactionCounter() uint16 {
	return m.TransactionCounter
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewDF1ResponseMessage factory function for _DF1ResponseMessage
func NewDF1ResponseMessage(destinationAddress uint8, sourceAddress uint8, status uint8, transactionCounter uint16, payloadLength uint16) *_DF1ResponseMessage {
	return &_DF1ResponseMessage{DestinationAddress: destinationAddress, SourceAddress: sourceAddress, Status: status, TransactionCounter: transactionCounter, PayloadLength: payloadLength}
}

// Deprecated: use the interface for direct cast
func CastDF1ResponseMessage(structType any) DF1ResponseMessage {
	if casted, ok := structType.(DF1ResponseMessage); ok {
		return casted
	}
	if casted, ok := structType.(*DF1ResponseMessage); ok {
		return *casted
	}
	return nil
}

func (m *_DF1ResponseMessage) GetTypeName() string {
	return "DF1ResponseMessage"
}

func (m *_DF1ResponseMessage) GetParentLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Reserved Field (reserved)
	lengthInBits += 8

	// Simple field (destinationAddress)
	lengthInBits += 8

	// Simple field (sourceAddress)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 8
	// Discriminator Field (commandCode)
	lengthInBits += 8

	// Simple field (status)
	lengthInBits += 8

	// Simple field (transactionCounter)
	lengthInBits += 16

	return lengthInBits
}

func (m *_DF1ResponseMessage) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func DF1ResponseMessageParse(ctx context.Context, theBytes []byte, payloadLength uint16) (DF1ResponseMessage, error) {
	return DF1ResponseMessageParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), payloadLength)
}

func DF1ResponseMessageParseWithBufferProducer[T DF1ResponseMessage](payloadLength uint16) func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
		buffer, err := DF1ResponseMessageParseWithBuffer(ctx, readBuffer, payloadLength)
		if err != nil {
			var zero T
			return zero, err
		}
		return buffer.(T), err
	}
}

func DF1ResponseMessageParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, payloadLength uint16) (DF1ResponseMessage, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("DF1ResponseMessage"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for DF1ResponseMessage")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	reservedField0, err := ReadReservedField(ctx, "reserved", ReadUnsignedByte(readBuffer, uint8(8)), uint8(0x00))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}

	destinationAddress, err := ReadSimpleField(ctx, "destinationAddress", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'destinationAddress' field"))
	}

	sourceAddress, err := ReadSimpleField(ctx, "sourceAddress", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'sourceAddress' field"))
	}

	reservedField1, err := ReadReservedField(ctx, "reserved", ReadUnsignedByte(readBuffer, uint8(8)), uint8(0x00))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}

	commandCode, err := ReadDiscriminatorField[uint8](ctx, "commandCode", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'commandCode' field"))
	}

	status, err := ReadSimpleField(ctx, "status", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'status' field"))
	}

	transactionCounter, err := ReadSimpleField(ctx, "transactionCounter", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'transactionCounter' field"))
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	type DF1ResponseMessageChildSerializeRequirement interface {
		DF1ResponseMessage
		InitializeParent(DF1ResponseMessage, uint8, uint8, uint8, uint16)
		GetParent() DF1ResponseMessage
	}
	var _childTemp any
	var _child DF1ResponseMessageChildSerializeRequirement
	var typeSwitchError error
	switch {
	case commandCode == 0x4F: // DF1CommandResponseMessageProtectedTypedLogicalRead
		_childTemp, typeSwitchError = DF1CommandResponseMessageProtectedTypedLogicalReadParseWithBuffer(ctx, readBuffer, payloadLength)
	default:
		typeSwitchError = errors.Errorf("Unmapped type for parameters [commandCode=%v]", commandCode)
	}
	if typeSwitchError != nil {
		return nil, errors.Wrap(typeSwitchError, "Error parsing sub-type for type-switch of DF1ResponseMessage")
	}
	_child = _childTemp.(DF1ResponseMessageChildSerializeRequirement)

	if closeErr := readBuffer.CloseContext("DF1ResponseMessage"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for DF1ResponseMessage")
	}

	// Finish initializing
	_child.InitializeParent(_child, destinationAddress, sourceAddress, status, transactionCounter)
	_child.GetParent().(*_DF1ResponseMessage).reservedField0 = reservedField0
	_child.GetParent().(*_DF1ResponseMessage).reservedField1 = reservedField1
	return _child, nil
}

func (pm *_DF1ResponseMessage) SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child DF1ResponseMessage, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("DF1ResponseMessage"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for DF1ResponseMessage")
	}

	if err := WriteReservedField[uint8](ctx, "reserved", uint8(0x00), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'reserved' field number 1")
	}

	if err := WriteSimpleField[uint8](ctx, "destinationAddress", m.GetDestinationAddress(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'destinationAddress' field")
	}

	if err := WriteSimpleField[uint8](ctx, "sourceAddress", m.GetSourceAddress(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'sourceAddress' field")
	}

	if err := WriteReservedField[uint8](ctx, "reserved", uint8(0x00), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'reserved' field number 2")
	}

	if err := WriteDiscriminatorField(ctx, "commandCode", m.GetCommandCode(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'commandCode' field")
	}

	if err := WriteSimpleField[uint8](ctx, "status", m.GetStatus(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'status' field")
	}

	if err := WriteSimpleField[uint16](ctx, "transactionCounter", m.GetTransactionCounter(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
		return errors.Wrap(err, "Error serializing 'transactionCounter' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("DF1ResponseMessage"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for DF1ResponseMessage")
	}
	return nil
}

////
// Arguments Getter

func (m *_DF1ResponseMessage) GetPayloadLength() uint16 {
	return m.PayloadLength
}

//
////

func (m *_DF1ResponseMessage) isDF1ResponseMessage() bool {
	return true
}

func (m *_DF1ResponseMessage) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
