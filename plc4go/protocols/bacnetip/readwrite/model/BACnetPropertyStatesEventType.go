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

// BACnetPropertyStatesEventType is the corresponding interface of BACnetPropertyStatesEventType
type BACnetPropertyStatesEventType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetPropertyStates
	// GetEventType returns EventType (property field)
	GetEventType() BACnetEventTypeTagged
}

// BACnetPropertyStatesEventTypeExactly can be used when we want exactly this type and not a type which fulfills BACnetPropertyStatesEventType.
// This is useful for switch cases.
type BACnetPropertyStatesEventTypeExactly interface {
	BACnetPropertyStatesEventType
	isBACnetPropertyStatesEventType() bool
}

// _BACnetPropertyStatesEventType is the data-structure of this message
type _BACnetPropertyStatesEventType struct {
	*_BACnetPropertyStates
	EventType BACnetEventTypeTagged
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetPropertyStatesEventType) InitializeParent(parent BACnetPropertyStates, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetPropertyStatesEventType) GetParent() BACnetPropertyStates {
	return m._BACnetPropertyStates
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPropertyStatesEventType) GetEventType() BACnetEventTypeTagged {
	return m.EventType
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetPropertyStatesEventType factory function for _BACnetPropertyStatesEventType
func NewBACnetPropertyStatesEventType(eventType BACnetEventTypeTagged, peekedTagHeader BACnetTagHeader) *_BACnetPropertyStatesEventType {
	_result := &_BACnetPropertyStatesEventType{
		EventType:             eventType,
		_BACnetPropertyStates: NewBACnetPropertyStates(peekedTagHeader),
	}
	_result._BACnetPropertyStates._BACnetPropertyStatesChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetPropertyStatesEventType(structType any) BACnetPropertyStatesEventType {
	if casted, ok := structType.(BACnetPropertyStatesEventType); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPropertyStatesEventType); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPropertyStatesEventType) GetTypeName() string {
	return "BACnetPropertyStatesEventType"
}

func (m *_BACnetPropertyStatesEventType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (eventType)
	lengthInBits += m.EventType.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetPropertyStatesEventType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetPropertyStatesEventTypeParse(ctx context.Context, theBytes []byte, peekedTagNumber uint8) (BACnetPropertyStatesEventType, error) {
	return BACnetPropertyStatesEventTypeParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), peekedTagNumber)
}

func BACnetPropertyStatesEventTypeParseWithBufferProducer(peekedTagNumber uint8) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetPropertyStatesEventType, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetPropertyStatesEventType, error) {
		return BACnetPropertyStatesEventTypeParseWithBuffer(ctx, readBuffer, peekedTagNumber)
	}
}

func BACnetPropertyStatesEventTypeParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, peekedTagNumber uint8) (BACnetPropertyStatesEventType, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPropertyStatesEventType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPropertyStatesEventType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	eventType, err := ReadSimpleField[BACnetEventTypeTagged](ctx, "eventType", ReadComplex[BACnetEventTypeTagged](BACnetEventTypeTaggedParseWithBufferProducer((uint8)(peekedTagNumber), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'eventType' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetPropertyStatesEventType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPropertyStatesEventType")
	}

	// Create a partially initialized instance
	_child := &_BACnetPropertyStatesEventType{
		_BACnetPropertyStates: &_BACnetPropertyStates{},
		EventType:             eventType,
	}
	_child._BACnetPropertyStates._BACnetPropertyStatesChildRequirements = _child
	return _child, nil
}

func (m *_BACnetPropertyStatesEventType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetPropertyStatesEventType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPropertyStatesEventType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPropertyStatesEventType")
		}

		if err := WriteSimpleField[BACnetEventTypeTagged](ctx, "eventType", m.GetEventType(), WriteComplex[BACnetEventTypeTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'eventType' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPropertyStatesEventType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPropertyStatesEventType")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetPropertyStatesEventType) isBACnetPropertyStatesEventType() bool {
	return true
}

func (m *_BACnetPropertyStatesEventType) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
