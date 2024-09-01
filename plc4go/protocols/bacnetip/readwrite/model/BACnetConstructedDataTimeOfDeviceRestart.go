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

// BACnetConstructedDataTimeOfDeviceRestart is the corresponding interface of BACnetConstructedDataTimeOfDeviceRestart
type BACnetConstructedDataTimeOfDeviceRestart interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetTimeOfDeviceRestart returns TimeOfDeviceRestart (property field)
	GetTimeOfDeviceRestart() BACnetTimeStamp
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetTimeStamp
}

// BACnetConstructedDataTimeOfDeviceRestartExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataTimeOfDeviceRestart.
// This is useful for switch cases.
type BACnetConstructedDataTimeOfDeviceRestartExactly interface {
	BACnetConstructedDataTimeOfDeviceRestart
	isBACnetConstructedDataTimeOfDeviceRestart() bool
}

// _BACnetConstructedDataTimeOfDeviceRestart is the data-structure of this message
type _BACnetConstructedDataTimeOfDeviceRestart struct {
	*_BACnetConstructedData
	TimeOfDeviceRestart BACnetTimeStamp
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataTimeOfDeviceRestart) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataTimeOfDeviceRestart) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_TIME_OF_DEVICE_RESTART
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataTimeOfDeviceRestart) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataTimeOfDeviceRestart) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataTimeOfDeviceRestart) GetTimeOfDeviceRestart() BACnetTimeStamp {
	return m.TimeOfDeviceRestart
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataTimeOfDeviceRestart) GetActualValue() BACnetTimeStamp {
	ctx := context.Background()
	_ = ctx
	return CastBACnetTimeStamp(m.GetTimeOfDeviceRestart())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataTimeOfDeviceRestart factory function for _BACnetConstructedDataTimeOfDeviceRestart
func NewBACnetConstructedDataTimeOfDeviceRestart(timeOfDeviceRestart BACnetTimeStamp, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataTimeOfDeviceRestart {
	_result := &_BACnetConstructedDataTimeOfDeviceRestart{
		TimeOfDeviceRestart:    timeOfDeviceRestart,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataTimeOfDeviceRestart(structType any) BACnetConstructedDataTimeOfDeviceRestart {
	if casted, ok := structType.(BACnetConstructedDataTimeOfDeviceRestart); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataTimeOfDeviceRestart); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataTimeOfDeviceRestart) GetTypeName() string {
	return "BACnetConstructedDataTimeOfDeviceRestart"
}

func (m *_BACnetConstructedDataTimeOfDeviceRestart) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (timeOfDeviceRestart)
	lengthInBits += m.TimeOfDeviceRestart.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataTimeOfDeviceRestart) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConstructedDataTimeOfDeviceRestartParse(ctx context.Context, theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataTimeOfDeviceRestart, error) {
	return BACnetConstructedDataTimeOfDeviceRestartParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataTimeOfDeviceRestartParseWithBufferProducer(tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConstructedDataTimeOfDeviceRestart, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConstructedDataTimeOfDeviceRestart, error) {
		return BACnetConstructedDataTimeOfDeviceRestartParseWithBuffer(ctx, readBuffer, tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
	}
}

func BACnetConstructedDataTimeOfDeviceRestartParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataTimeOfDeviceRestart, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataTimeOfDeviceRestart"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataTimeOfDeviceRestart")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	timeOfDeviceRestart, err := ReadSimpleField[BACnetTimeStamp](ctx, "timeOfDeviceRestart", ReadComplex[BACnetTimeStamp](BACnetTimeStampParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'timeOfDeviceRestart' field"))
	}

	actualValue, err := ReadVirtualField[BACnetTimeStamp](ctx, "actualValue", (*BACnetTimeStamp)(nil), timeOfDeviceRestart)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataTimeOfDeviceRestart"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataTimeOfDeviceRestart")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataTimeOfDeviceRestart{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		TimeOfDeviceRestart: timeOfDeviceRestart,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataTimeOfDeviceRestart) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataTimeOfDeviceRestart) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataTimeOfDeviceRestart"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataTimeOfDeviceRestart")
		}

		if err := WriteSimpleField[BACnetTimeStamp](ctx, "timeOfDeviceRestart", m.GetTimeOfDeviceRestart(), WriteComplex[BACnetTimeStamp](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'timeOfDeviceRestart' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataTimeOfDeviceRestart"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataTimeOfDeviceRestart")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataTimeOfDeviceRestart) isBACnetConstructedDataTimeOfDeviceRestart() bool {
	return true
}

func (m *_BACnetConstructedDataTimeOfDeviceRestart) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
