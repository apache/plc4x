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

// BACnetConstructedDataOccupancyState is the corresponding interface of BACnetConstructedDataOccupancyState
type BACnetConstructedDataOccupancyState interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetOccupancyState returns OccupancyState (property field)
	GetOccupancyState() BACnetAccessZoneOccupancyStateTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetAccessZoneOccupancyStateTagged
	// IsBACnetConstructedDataOccupancyState is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataOccupancyState()
}

// _BACnetConstructedDataOccupancyState is the data-structure of this message
type _BACnetConstructedDataOccupancyState struct {
	BACnetConstructedDataContract
	OccupancyState BACnetAccessZoneOccupancyStateTagged
}

var _ BACnetConstructedDataOccupancyState = (*_BACnetConstructedDataOccupancyState)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataOccupancyState)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataOccupancyState) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataOccupancyState) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_OCCUPANCY_STATE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataOccupancyState) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataOccupancyState) GetOccupancyState() BACnetAccessZoneOccupancyStateTagged {
	return m.OccupancyState
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataOccupancyState) GetActualValue() BACnetAccessZoneOccupancyStateTagged {
	ctx := context.Background()
	_ = ctx
	return CastBACnetAccessZoneOccupancyStateTagged(m.GetOccupancyState())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataOccupancyState factory function for _BACnetConstructedDataOccupancyState
func NewBACnetConstructedDataOccupancyState(occupancyState BACnetAccessZoneOccupancyStateTagged, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataOccupancyState {
	if occupancyState == nil {
		panic("occupancyState of type BACnetAccessZoneOccupancyStateTagged for BACnetConstructedDataOccupancyState must not be nil")
	}
	_result := &_BACnetConstructedDataOccupancyState{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		OccupancyState:                occupancyState,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataOccupancyState(structType any) BACnetConstructedDataOccupancyState {
	if casted, ok := structType.(BACnetConstructedDataOccupancyState); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataOccupancyState); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataOccupancyState) GetTypeName() string {
	return "BACnetConstructedDataOccupancyState"
}

func (m *_BACnetConstructedDataOccupancyState) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).getLengthInBits(ctx))

	// Simple field (occupancyState)
	lengthInBits += m.OccupancyState.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataOccupancyState) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataOccupancyState) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataOccupancyState BACnetConstructedDataOccupancyState, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataOccupancyState"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataOccupancyState")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	occupancyState, err := ReadSimpleField[BACnetAccessZoneOccupancyStateTagged](ctx, "occupancyState", ReadComplex[BACnetAccessZoneOccupancyStateTagged](BACnetAccessZoneOccupancyStateTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'occupancyState' field"))
	}
	m.OccupancyState = occupancyState

	actualValue, err := ReadVirtualField[BACnetAccessZoneOccupancyStateTagged](ctx, "actualValue", (*BACnetAccessZoneOccupancyStateTagged)(nil), occupancyState)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataOccupancyState"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataOccupancyState")
	}

	return m, nil
}

func (m *_BACnetConstructedDataOccupancyState) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataOccupancyState) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataOccupancyState"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataOccupancyState")
		}

		if err := WriteSimpleField[BACnetAccessZoneOccupancyStateTagged](ctx, "occupancyState", m.GetOccupancyState(), WriteComplex[BACnetAccessZoneOccupancyStateTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'occupancyState' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataOccupancyState"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataOccupancyState")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataOccupancyState) IsBACnetConstructedDataOccupancyState() {}

func (m *_BACnetConstructedDataOccupancyState) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
