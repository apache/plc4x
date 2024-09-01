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

// BACnetLandingCallStatusCommandDirection is the corresponding interface of BACnetLandingCallStatusCommandDirection
type BACnetLandingCallStatusCommandDirection interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetLandingCallStatusCommand
	// GetDirection returns Direction (property field)
	GetDirection() BACnetLiftCarDirectionTagged
}

// BACnetLandingCallStatusCommandDirectionExactly can be used when we want exactly this type and not a type which fulfills BACnetLandingCallStatusCommandDirection.
// This is useful for switch cases.
type BACnetLandingCallStatusCommandDirectionExactly interface {
	BACnetLandingCallStatusCommandDirection
	isBACnetLandingCallStatusCommandDirection() bool
}

// _BACnetLandingCallStatusCommandDirection is the data-structure of this message
type _BACnetLandingCallStatusCommandDirection struct {
	*_BACnetLandingCallStatusCommand
	Direction BACnetLiftCarDirectionTagged
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetLandingCallStatusCommandDirection) InitializeParent(parent BACnetLandingCallStatusCommand, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetLandingCallStatusCommandDirection) GetParent() BACnetLandingCallStatusCommand {
	return m._BACnetLandingCallStatusCommand
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetLandingCallStatusCommandDirection) GetDirection() BACnetLiftCarDirectionTagged {
	return m.Direction
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetLandingCallStatusCommandDirection factory function for _BACnetLandingCallStatusCommandDirection
func NewBACnetLandingCallStatusCommandDirection(direction BACnetLiftCarDirectionTagged, peekedTagHeader BACnetTagHeader) *_BACnetLandingCallStatusCommandDirection {
	_result := &_BACnetLandingCallStatusCommandDirection{
		Direction:                       direction,
		_BACnetLandingCallStatusCommand: NewBACnetLandingCallStatusCommand(peekedTagHeader),
	}
	_result._BACnetLandingCallStatusCommand._BACnetLandingCallStatusCommandChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetLandingCallStatusCommandDirection(structType any) BACnetLandingCallStatusCommandDirection {
	if casted, ok := structType.(BACnetLandingCallStatusCommandDirection); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetLandingCallStatusCommandDirection); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetLandingCallStatusCommandDirection) GetTypeName() string {
	return "BACnetLandingCallStatusCommandDirection"
}

func (m *_BACnetLandingCallStatusCommandDirection) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (direction)
	lengthInBits += m.Direction.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetLandingCallStatusCommandDirection) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetLandingCallStatusCommandDirectionParse(ctx context.Context, theBytes []byte) (BACnetLandingCallStatusCommandDirection, error) {
	return BACnetLandingCallStatusCommandDirectionParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func BACnetLandingCallStatusCommandDirectionParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLandingCallStatusCommandDirection, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLandingCallStatusCommandDirection, error) {
		return BACnetLandingCallStatusCommandDirectionParseWithBuffer(ctx, readBuffer)
	}
}

func BACnetLandingCallStatusCommandDirectionParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLandingCallStatusCommandDirection, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetLandingCallStatusCommandDirection"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetLandingCallStatusCommandDirection")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	direction, err := ReadSimpleField[BACnetLiftCarDirectionTagged](ctx, "direction", ReadComplex[BACnetLiftCarDirectionTagged](BACnetLiftCarDirectionTaggedParseWithBufferProducer((uint8)(uint8(1)), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'direction' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetLandingCallStatusCommandDirection"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetLandingCallStatusCommandDirection")
	}

	// Create a partially initialized instance
	_child := &_BACnetLandingCallStatusCommandDirection{
		_BACnetLandingCallStatusCommand: &_BACnetLandingCallStatusCommand{},
		Direction:                       direction,
	}
	_child._BACnetLandingCallStatusCommand._BACnetLandingCallStatusCommandChildRequirements = _child
	return _child, nil
}

func (m *_BACnetLandingCallStatusCommandDirection) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetLandingCallStatusCommandDirection) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetLandingCallStatusCommandDirection"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetLandingCallStatusCommandDirection")
		}

		if err := WriteSimpleField[BACnetLiftCarDirectionTagged](ctx, "direction", m.GetDirection(), WriteComplex[BACnetLiftCarDirectionTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'direction' field")
		}

		if popErr := writeBuffer.PopContext("BACnetLandingCallStatusCommandDirection"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetLandingCallStatusCommandDirection")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetLandingCallStatusCommandDirection) isBACnetLandingCallStatusCommandDirection() bool {
	return true
}

func (m *_BACnetLandingCallStatusCommandDirection) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
