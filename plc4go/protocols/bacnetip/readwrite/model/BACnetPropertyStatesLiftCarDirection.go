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

// BACnetPropertyStatesLiftCarDirection is the corresponding interface of BACnetPropertyStatesLiftCarDirection
type BACnetPropertyStatesLiftCarDirection interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetPropertyStates
	// GetLiftCarDirection returns LiftCarDirection (property field)
	GetLiftCarDirection() BACnetLiftCarDirectionTagged
}

// BACnetPropertyStatesLiftCarDirectionExactly can be used when we want exactly this type and not a type which fulfills BACnetPropertyStatesLiftCarDirection.
// This is useful for switch cases.
type BACnetPropertyStatesLiftCarDirectionExactly interface {
	BACnetPropertyStatesLiftCarDirection
	isBACnetPropertyStatesLiftCarDirection() bool
}

// _BACnetPropertyStatesLiftCarDirection is the data-structure of this message
type _BACnetPropertyStatesLiftCarDirection struct {
	*_BACnetPropertyStates
	LiftCarDirection BACnetLiftCarDirectionTagged
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetPropertyStatesLiftCarDirection) InitializeParent(parent BACnetPropertyStates, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetPropertyStatesLiftCarDirection) GetParent() BACnetPropertyStates {
	return m._BACnetPropertyStates
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPropertyStatesLiftCarDirection) GetLiftCarDirection() BACnetLiftCarDirectionTagged {
	return m.LiftCarDirection
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetPropertyStatesLiftCarDirection factory function for _BACnetPropertyStatesLiftCarDirection
func NewBACnetPropertyStatesLiftCarDirection(liftCarDirection BACnetLiftCarDirectionTagged, peekedTagHeader BACnetTagHeader) *_BACnetPropertyStatesLiftCarDirection {
	_result := &_BACnetPropertyStatesLiftCarDirection{
		LiftCarDirection:      liftCarDirection,
		_BACnetPropertyStates: NewBACnetPropertyStates(peekedTagHeader),
	}
	_result._BACnetPropertyStates._BACnetPropertyStatesChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetPropertyStatesLiftCarDirection(structType any) BACnetPropertyStatesLiftCarDirection {
	if casted, ok := structType.(BACnetPropertyStatesLiftCarDirection); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPropertyStatesLiftCarDirection); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPropertyStatesLiftCarDirection) GetTypeName() string {
	return "BACnetPropertyStatesLiftCarDirection"
}

func (m *_BACnetPropertyStatesLiftCarDirection) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (liftCarDirection)
	lengthInBits += m.LiftCarDirection.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetPropertyStatesLiftCarDirection) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetPropertyStatesLiftCarDirectionParse(ctx context.Context, theBytes []byte, peekedTagNumber uint8) (BACnetPropertyStatesLiftCarDirection, error) {
	return BACnetPropertyStatesLiftCarDirectionParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), peekedTagNumber)
}

func BACnetPropertyStatesLiftCarDirectionParseWithBufferProducer(peekedTagNumber uint8) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetPropertyStatesLiftCarDirection, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetPropertyStatesLiftCarDirection, error) {
		return BACnetPropertyStatesLiftCarDirectionParseWithBuffer(ctx, readBuffer, peekedTagNumber)
	}
}

func BACnetPropertyStatesLiftCarDirectionParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, peekedTagNumber uint8) (BACnetPropertyStatesLiftCarDirection, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPropertyStatesLiftCarDirection"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPropertyStatesLiftCarDirection")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	liftCarDirection, err := ReadSimpleField[BACnetLiftCarDirectionTagged](ctx, "liftCarDirection", ReadComplex[BACnetLiftCarDirectionTagged](BACnetLiftCarDirectionTaggedParseWithBufferProducer((uint8)(peekedTagNumber), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'liftCarDirection' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetPropertyStatesLiftCarDirection"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPropertyStatesLiftCarDirection")
	}

	// Create a partially initialized instance
	_child := &_BACnetPropertyStatesLiftCarDirection{
		_BACnetPropertyStates: &_BACnetPropertyStates{},
		LiftCarDirection:      liftCarDirection,
	}
	_child._BACnetPropertyStates._BACnetPropertyStatesChildRequirements = _child
	return _child, nil
}

func (m *_BACnetPropertyStatesLiftCarDirection) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetPropertyStatesLiftCarDirection) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPropertyStatesLiftCarDirection"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPropertyStatesLiftCarDirection")
		}

		if err := WriteSimpleField[BACnetLiftCarDirectionTagged](ctx, "liftCarDirection", m.GetLiftCarDirection(), WriteComplex[BACnetLiftCarDirectionTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'liftCarDirection' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPropertyStatesLiftCarDirection"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPropertyStatesLiftCarDirection")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetPropertyStatesLiftCarDirection) isBACnetPropertyStatesLiftCarDirection() bool {
	return true
}

func (m *_BACnetPropertyStatesLiftCarDirection) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
