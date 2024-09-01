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

// BACnetPriorityValueReal is the corresponding interface of BACnetPriorityValueReal
type BACnetPriorityValueReal interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetPriorityValue
	// GetRealValue returns RealValue (property field)
	GetRealValue() BACnetApplicationTagReal
}

// BACnetPriorityValueRealExactly can be used when we want exactly this type and not a type which fulfills BACnetPriorityValueReal.
// This is useful for switch cases.
type BACnetPriorityValueRealExactly interface {
	BACnetPriorityValueReal
	isBACnetPriorityValueReal() bool
}

// _BACnetPriorityValueReal is the data-structure of this message
type _BACnetPriorityValueReal struct {
	*_BACnetPriorityValue
	RealValue BACnetApplicationTagReal
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetPriorityValueReal) InitializeParent(parent BACnetPriorityValue, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetPriorityValueReal) GetParent() BACnetPriorityValue {
	return m._BACnetPriorityValue
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPriorityValueReal) GetRealValue() BACnetApplicationTagReal {
	return m.RealValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetPriorityValueReal factory function for _BACnetPriorityValueReal
func NewBACnetPriorityValueReal(realValue BACnetApplicationTagReal, peekedTagHeader BACnetTagHeader, objectTypeArgument BACnetObjectType) *_BACnetPriorityValueReal {
	_result := &_BACnetPriorityValueReal{
		RealValue:            realValue,
		_BACnetPriorityValue: NewBACnetPriorityValue(peekedTagHeader, objectTypeArgument),
	}
	_result._BACnetPriorityValue._BACnetPriorityValueChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetPriorityValueReal(structType any) BACnetPriorityValueReal {
	if casted, ok := structType.(BACnetPriorityValueReal); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPriorityValueReal); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPriorityValueReal) GetTypeName() string {
	return "BACnetPriorityValueReal"
}

func (m *_BACnetPriorityValueReal) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (realValue)
	lengthInBits += m.RealValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetPriorityValueReal) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetPriorityValueRealParse(ctx context.Context, theBytes []byte, objectTypeArgument BACnetObjectType) (BACnetPriorityValueReal, error) {
	return BACnetPriorityValueRealParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), objectTypeArgument)
}

func BACnetPriorityValueRealParseWithBufferProducer(objectTypeArgument BACnetObjectType) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetPriorityValueReal, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetPriorityValueReal, error) {
		return BACnetPriorityValueRealParseWithBuffer(ctx, readBuffer, objectTypeArgument)
	}
}

func BACnetPriorityValueRealParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, objectTypeArgument BACnetObjectType) (BACnetPriorityValueReal, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPriorityValueReal"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPriorityValueReal")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	realValue, err := ReadSimpleField[BACnetApplicationTagReal](ctx, "realValue", ReadComplex[BACnetApplicationTagReal](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagReal](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'realValue' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetPriorityValueReal"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPriorityValueReal")
	}

	// Create a partially initialized instance
	_child := &_BACnetPriorityValueReal{
		_BACnetPriorityValue: &_BACnetPriorityValue{
			ObjectTypeArgument: objectTypeArgument,
		},
		RealValue: realValue,
	}
	_child._BACnetPriorityValue._BACnetPriorityValueChildRequirements = _child
	return _child, nil
}

func (m *_BACnetPriorityValueReal) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetPriorityValueReal) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPriorityValueReal"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPriorityValueReal")
		}

		if err := WriteSimpleField[BACnetApplicationTagReal](ctx, "realValue", m.GetRealValue(), WriteComplex[BACnetApplicationTagReal](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'realValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPriorityValueReal"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPriorityValueReal")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetPriorityValueReal) isBACnetPriorityValueReal() bool {
	return true
}

func (m *_BACnetPriorityValueReal) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
