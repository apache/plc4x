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

// CBusHeader is the corresponding interface of CBusHeader
type CBusHeader interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetPriorityClass returns PriorityClass (property field)
	GetPriorityClass() PriorityClass
	// GetDp returns Dp (property field)
	GetDp() bool
	// GetRc returns Rc (property field)
	GetRc() uint8
	// GetDestinationAddressType returns DestinationAddressType (property field)
	GetDestinationAddressType() DestinationAddressType
}

// CBusHeaderExactly can be used when we want exactly this type and not a type which fulfills CBusHeader.
// This is useful for switch cases.
type CBusHeaderExactly interface {
	CBusHeader
	isCBusHeader() bool
}

// _CBusHeader is the data-structure of this message
type _CBusHeader struct {
	PriorityClass          PriorityClass
	Dp                     bool
	Rc                     uint8
	DestinationAddressType DestinationAddressType
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CBusHeader) GetPriorityClass() PriorityClass {
	return m.PriorityClass
}

func (m *_CBusHeader) GetDp() bool {
	return m.Dp
}

func (m *_CBusHeader) GetRc() uint8 {
	return m.Rc
}

func (m *_CBusHeader) GetDestinationAddressType() DestinationAddressType {
	return m.DestinationAddressType
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCBusHeader factory function for _CBusHeader
func NewCBusHeader(priorityClass PriorityClass, dp bool, rc uint8, destinationAddressType DestinationAddressType) *_CBusHeader {
	return &_CBusHeader{PriorityClass: priorityClass, Dp: dp, Rc: rc, DestinationAddressType: destinationAddressType}
}

// Deprecated: use the interface for direct cast
func CastCBusHeader(structType any) CBusHeader {
	if casted, ok := structType.(CBusHeader); ok {
		return casted
	}
	if casted, ok := structType.(*CBusHeader); ok {
		return *casted
	}
	return nil
}

func (m *_CBusHeader) GetTypeName() string {
	return "CBusHeader"
}

func (m *_CBusHeader) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (priorityClass)
	lengthInBits += 2

	// Simple field (dp)
	lengthInBits += 1

	// Simple field (rc)
	lengthInBits += 2

	// Simple field (destinationAddressType)
	lengthInBits += 3

	return lengthInBits
}

func (m *_CBusHeader) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func CBusHeaderParse(ctx context.Context, theBytes []byte) (CBusHeader, error) {
	return CBusHeaderParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func CBusHeaderParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (CBusHeader, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (CBusHeader, error) {
		return CBusHeaderParseWithBuffer(ctx, readBuffer)
	}
}

func CBusHeaderParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (CBusHeader, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CBusHeader"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CBusHeader")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	priorityClass, err := ReadEnumField[PriorityClass](ctx, "priorityClass", "PriorityClass", ReadEnum(PriorityClassByValue, ReadUnsignedByte(readBuffer, uint8(2))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'priorityClass' field"))
	}

	dp, err := ReadSimpleField(ctx, "dp", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'dp' field"))
	}

	rc, err := ReadSimpleField(ctx, "rc", ReadUnsignedByte(readBuffer, uint8(2)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'rc' field"))
	}

	destinationAddressType, err := ReadEnumField[DestinationAddressType](ctx, "destinationAddressType", "DestinationAddressType", ReadEnum(DestinationAddressTypeByValue, ReadUnsignedByte(readBuffer, uint8(3))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'destinationAddressType' field"))
	}

	if closeErr := readBuffer.CloseContext("CBusHeader"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CBusHeader")
	}

	// Create the instance
	return &_CBusHeader{
		PriorityClass:          priorityClass,
		Dp:                     dp,
		Rc:                     rc,
		DestinationAddressType: destinationAddressType,
	}, nil
}

func (m *_CBusHeader) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CBusHeader) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("CBusHeader"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for CBusHeader")
	}

	if err := WriteSimpleEnumField[PriorityClass](ctx, "priorityClass", "PriorityClass", m.GetPriorityClass(), WriteEnum[PriorityClass, uint8](PriorityClass.GetValue, PriorityClass.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 2))); err != nil {
		return errors.Wrap(err, "Error serializing 'priorityClass' field")
	}

	if err := WriteSimpleField[bool](ctx, "dp", m.GetDp(), WriteBoolean(writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'dp' field")
	}

	if err := WriteSimpleField[uint8](ctx, "rc", m.GetRc(), WriteUnsignedByte(writeBuffer, 2)); err != nil {
		return errors.Wrap(err, "Error serializing 'rc' field")
	}

	if err := WriteSimpleEnumField[DestinationAddressType](ctx, "destinationAddressType", "DestinationAddressType", m.GetDestinationAddressType(), WriteEnum[DestinationAddressType, uint8](DestinationAddressType.GetValue, DestinationAddressType.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 3))); err != nil {
		return errors.Wrap(err, "Error serializing 'destinationAddressType' field")
	}

	if popErr := writeBuffer.PopContext("CBusHeader"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for CBusHeader")
	}
	return nil
}

func (m *_CBusHeader) isCBusHeader() bool {
	return true
}

func (m *_CBusHeader) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
