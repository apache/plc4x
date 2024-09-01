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

// KnxGroupAddressFreeLevel is the corresponding interface of KnxGroupAddressFreeLevel
type KnxGroupAddressFreeLevel interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	KnxGroupAddress
	// GetSubGroup returns SubGroup (property field)
	GetSubGroup() uint16
}

// KnxGroupAddressFreeLevelExactly can be used when we want exactly this type and not a type which fulfills KnxGroupAddressFreeLevel.
// This is useful for switch cases.
type KnxGroupAddressFreeLevelExactly interface {
	KnxGroupAddressFreeLevel
	isKnxGroupAddressFreeLevel() bool
}

// _KnxGroupAddressFreeLevel is the data-structure of this message
type _KnxGroupAddressFreeLevel struct {
	*_KnxGroupAddress
	SubGroup uint16
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_KnxGroupAddressFreeLevel) GetNumLevels() uint8 {
	return uint8(1)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_KnxGroupAddressFreeLevel) InitializeParent(parent KnxGroupAddress) {}

func (m *_KnxGroupAddressFreeLevel) GetParent() KnxGroupAddress {
	return m._KnxGroupAddress
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_KnxGroupAddressFreeLevel) GetSubGroup() uint16 {
	return m.SubGroup
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewKnxGroupAddressFreeLevel factory function for _KnxGroupAddressFreeLevel
func NewKnxGroupAddressFreeLevel(subGroup uint16) *_KnxGroupAddressFreeLevel {
	_result := &_KnxGroupAddressFreeLevel{
		SubGroup:         subGroup,
		_KnxGroupAddress: NewKnxGroupAddress(),
	}
	_result._KnxGroupAddress._KnxGroupAddressChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastKnxGroupAddressFreeLevel(structType any) KnxGroupAddressFreeLevel {
	if casted, ok := structType.(KnxGroupAddressFreeLevel); ok {
		return casted
	}
	if casted, ok := structType.(*KnxGroupAddressFreeLevel); ok {
		return *casted
	}
	return nil
}

func (m *_KnxGroupAddressFreeLevel) GetTypeName() string {
	return "KnxGroupAddressFreeLevel"
}

func (m *_KnxGroupAddressFreeLevel) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (subGroup)
	lengthInBits += 16

	return lengthInBits
}

func (m *_KnxGroupAddressFreeLevel) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func KnxGroupAddressFreeLevelParse(ctx context.Context, theBytes []byte, numLevels uint8) (KnxGroupAddressFreeLevel, error) {
	return KnxGroupAddressFreeLevelParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), numLevels)
}

func KnxGroupAddressFreeLevelParseWithBufferProducer(numLevels uint8) func(ctx context.Context, readBuffer utils.ReadBuffer) (KnxGroupAddressFreeLevel, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (KnxGroupAddressFreeLevel, error) {
		return KnxGroupAddressFreeLevelParseWithBuffer(ctx, readBuffer, numLevels)
	}
}

func KnxGroupAddressFreeLevelParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, numLevels uint8) (KnxGroupAddressFreeLevel, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("KnxGroupAddressFreeLevel"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for KnxGroupAddressFreeLevel")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	subGroup, err := ReadSimpleField(ctx, "subGroup", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'subGroup' field"))
	}

	if closeErr := readBuffer.CloseContext("KnxGroupAddressFreeLevel"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for KnxGroupAddressFreeLevel")
	}

	// Create a partially initialized instance
	_child := &_KnxGroupAddressFreeLevel{
		_KnxGroupAddress: &_KnxGroupAddress{},
		SubGroup:         subGroup,
	}
	_child._KnxGroupAddress._KnxGroupAddressChildRequirements = _child
	return _child, nil
}

func (m *_KnxGroupAddressFreeLevel) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_KnxGroupAddressFreeLevel) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("KnxGroupAddressFreeLevel"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for KnxGroupAddressFreeLevel")
		}

		if err := WriteSimpleField[uint16](ctx, "subGroup", m.GetSubGroup(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'subGroup' field")
		}

		if popErr := writeBuffer.PopContext("KnxGroupAddressFreeLevel"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for KnxGroupAddressFreeLevel")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_KnxGroupAddressFreeLevel) isKnxGroupAddressFreeLevel() bool {
	return true
}

func (m *_KnxGroupAddressFreeLevel) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
