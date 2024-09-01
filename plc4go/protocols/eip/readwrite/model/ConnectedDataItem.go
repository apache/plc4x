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

// ConnectedDataItem is the corresponding interface of ConnectedDataItem
type ConnectedDataItem interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	TypeId
	// GetSequenceCount returns SequenceCount (property field)
	GetSequenceCount() uint16
	// GetService returns Service (property field)
	GetService() CipService
}

// ConnectedDataItemExactly can be used when we want exactly this type and not a type which fulfills ConnectedDataItem.
// This is useful for switch cases.
type ConnectedDataItemExactly interface {
	ConnectedDataItem
	isConnectedDataItem() bool
}

// _ConnectedDataItem is the data-structure of this message
type _ConnectedDataItem struct {
	*_TypeId
	SequenceCount uint16
	Service       CipService
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ConnectedDataItem) GetId() uint16 {
	return 0x00B1
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ConnectedDataItem) InitializeParent(parent TypeId) {}

func (m *_ConnectedDataItem) GetParent() TypeId {
	return m._TypeId
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ConnectedDataItem) GetSequenceCount() uint16 {
	return m.SequenceCount
}

func (m *_ConnectedDataItem) GetService() CipService {
	return m.Service
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewConnectedDataItem factory function for _ConnectedDataItem
func NewConnectedDataItem(sequenceCount uint16, service CipService) *_ConnectedDataItem {
	_result := &_ConnectedDataItem{
		SequenceCount: sequenceCount,
		Service:       service,
		_TypeId:       NewTypeId(),
	}
	_result._TypeId._TypeIdChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastConnectedDataItem(structType any) ConnectedDataItem {
	if casted, ok := structType.(ConnectedDataItem); ok {
		return casted
	}
	if casted, ok := structType.(*ConnectedDataItem); ok {
		return *casted
	}
	return nil
}

func (m *_ConnectedDataItem) GetTypeName() string {
	return "ConnectedDataItem"
}

func (m *_ConnectedDataItem) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Implicit Field (packetSize)
	lengthInBits += 16

	// Simple field (sequenceCount)
	lengthInBits += 16

	// Simple field (service)
	lengthInBits += m.Service.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_ConnectedDataItem) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func ConnectedDataItemParse(ctx context.Context, theBytes []byte) (ConnectedDataItem, error) {
	return ConnectedDataItemParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func ConnectedDataItemParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (ConnectedDataItem, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (ConnectedDataItem, error) {
		return ConnectedDataItemParseWithBuffer(ctx, readBuffer)
	}
}

func ConnectedDataItemParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (ConnectedDataItem, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ConnectedDataItem"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ConnectedDataItem")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	packetSize, err := ReadImplicitField[uint16](ctx, "packetSize", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'packetSize' field"))
	}
	_ = packetSize

	sequenceCount, err := ReadSimpleField(ctx, "sequenceCount", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'sequenceCount' field"))
	}

	service, err := ReadSimpleField[CipService](ctx, "service", ReadComplex[CipService](CipServiceParseWithBufferProducer[CipService]((bool)(bool(true)), (uint16)(uint16(packetSize)-uint16(uint16(2)))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'service' field"))
	}

	if closeErr := readBuffer.CloseContext("ConnectedDataItem"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ConnectedDataItem")
	}

	// Create a partially initialized instance
	_child := &_ConnectedDataItem{
		_TypeId:       &_TypeId{},
		SequenceCount: sequenceCount,
		Service:       service,
	}
	_child._TypeId._TypeIdChildRequirements = _child
	return _child, nil
}

func (m *_ConnectedDataItem) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ConnectedDataItem) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ConnectedDataItem"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ConnectedDataItem")
		}
		packetSize := uint16(uint16(m.GetService().GetLengthInBytes(ctx)) + uint16(uint16(2)))
		if err := WriteImplicitField(ctx, "packetSize", packetSize, WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'packetSize' field")
		}

		if err := WriteSimpleField[uint16](ctx, "sequenceCount", m.GetSequenceCount(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'sequenceCount' field")
		}

		if err := WriteSimpleField[CipService](ctx, "service", m.GetService(), WriteComplex[CipService](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'service' field")
		}

		if popErr := writeBuffer.PopContext("ConnectedDataItem"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ConnectedDataItem")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ConnectedDataItem) isConnectedDataItem() bool {
	return true
}

func (m *_ConnectedDataItem) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
