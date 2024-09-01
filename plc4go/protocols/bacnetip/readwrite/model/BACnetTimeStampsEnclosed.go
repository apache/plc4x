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

// BACnetTimeStampsEnclosed is the corresponding interface of BACnetTimeStampsEnclosed
type BACnetTimeStampsEnclosed interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetOpeningTag returns OpeningTag (property field)
	GetOpeningTag() BACnetOpeningTag
	// GetTimestamps returns Timestamps (property field)
	GetTimestamps() []BACnetTimeStamp
	// GetClosingTag returns ClosingTag (property field)
	GetClosingTag() BACnetClosingTag
}

// BACnetTimeStampsEnclosedExactly can be used when we want exactly this type and not a type which fulfills BACnetTimeStampsEnclosed.
// This is useful for switch cases.
type BACnetTimeStampsEnclosedExactly interface {
	BACnetTimeStampsEnclosed
	isBACnetTimeStampsEnclosed() bool
}

// _BACnetTimeStampsEnclosed is the data-structure of this message
type _BACnetTimeStampsEnclosed struct {
	OpeningTag BACnetOpeningTag
	Timestamps []BACnetTimeStamp
	ClosingTag BACnetClosingTag

	// Arguments.
	TagNumber uint8
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetTimeStampsEnclosed) GetOpeningTag() BACnetOpeningTag {
	return m.OpeningTag
}

func (m *_BACnetTimeStampsEnclosed) GetTimestamps() []BACnetTimeStamp {
	return m.Timestamps
}

func (m *_BACnetTimeStampsEnclosed) GetClosingTag() BACnetClosingTag {
	return m.ClosingTag
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetTimeStampsEnclosed factory function for _BACnetTimeStampsEnclosed
func NewBACnetTimeStampsEnclosed(openingTag BACnetOpeningTag, timestamps []BACnetTimeStamp, closingTag BACnetClosingTag, tagNumber uint8) *_BACnetTimeStampsEnclosed {
	return &_BACnetTimeStampsEnclosed{OpeningTag: openingTag, Timestamps: timestamps, ClosingTag: closingTag, TagNumber: tagNumber}
}

// Deprecated: use the interface for direct cast
func CastBACnetTimeStampsEnclosed(structType any) BACnetTimeStampsEnclosed {
	if casted, ok := structType.(BACnetTimeStampsEnclosed); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetTimeStampsEnclosed); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetTimeStampsEnclosed) GetTypeName() string {
	return "BACnetTimeStampsEnclosed"
}

func (m *_BACnetTimeStampsEnclosed) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (openingTag)
	lengthInBits += m.OpeningTag.GetLengthInBits(ctx)

	// Array field
	if len(m.Timestamps) > 0 {
		for _, element := range m.Timestamps {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	// Simple field (closingTag)
	lengthInBits += m.ClosingTag.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetTimeStampsEnclosed) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetTimeStampsEnclosedParse(ctx context.Context, theBytes []byte, tagNumber uint8) (BACnetTimeStampsEnclosed, error) {
	return BACnetTimeStampsEnclosedParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), tagNumber)
}

func BACnetTimeStampsEnclosedParseWithBufferProducer(tagNumber uint8) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetTimeStampsEnclosed, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetTimeStampsEnclosed, error) {
		return BACnetTimeStampsEnclosedParseWithBuffer(ctx, readBuffer, tagNumber)
	}
}

func BACnetTimeStampsEnclosedParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, tagNumber uint8) (BACnetTimeStampsEnclosed, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetTimeStampsEnclosed"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetTimeStampsEnclosed")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	openingTag, err := ReadSimpleField[BACnetOpeningTag](ctx, "openingTag", ReadComplex[BACnetOpeningTag](BACnetOpeningTagParseWithBufferProducer((uint8)(tagNumber)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'openingTag' field"))
	}

	timestamps, err := ReadTerminatedArrayField[BACnetTimeStamp](ctx, "timestamps", ReadComplex[BACnetTimeStamp](BACnetTimeStampParseWithBuffer, readBuffer), IsBACnetConstructedDataClosingTag(ctx, readBuffer, false, tagNumber))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'timestamps' field"))
	}

	closingTag, err := ReadSimpleField[BACnetClosingTag](ctx, "closingTag", ReadComplex[BACnetClosingTag](BACnetClosingTagParseWithBufferProducer((uint8)(tagNumber)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'closingTag' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetTimeStampsEnclosed"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetTimeStampsEnclosed")
	}

	// Create the instance
	return &_BACnetTimeStampsEnclosed{
		TagNumber:  tagNumber,
		OpeningTag: openingTag,
		Timestamps: timestamps,
		ClosingTag: closingTag,
	}, nil
}

func (m *_BACnetTimeStampsEnclosed) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetTimeStampsEnclosed) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("BACnetTimeStampsEnclosed"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetTimeStampsEnclosed")
	}

	if err := WriteSimpleField[BACnetOpeningTag](ctx, "openingTag", m.GetOpeningTag(), WriteComplex[BACnetOpeningTag](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'openingTag' field")
	}

	if err := WriteComplexTypeArrayField(ctx, "timestamps", m.GetTimestamps(), writeBuffer); err != nil {
		return errors.Wrap(err, "Error serializing 'timestamps' field")
	}

	if err := WriteSimpleField[BACnetClosingTag](ctx, "closingTag", m.GetClosingTag(), WriteComplex[BACnetClosingTag](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'closingTag' field")
	}

	if popErr := writeBuffer.PopContext("BACnetTimeStampsEnclosed"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetTimeStampsEnclosed")
	}
	return nil
}

////
// Arguments Getter

func (m *_BACnetTimeStampsEnclosed) GetTagNumber() uint8 {
	return m.TagNumber
}

//
////

func (m *_BACnetTimeStampsEnclosed) isBACnetTimeStampsEnclosed() bool {
	return true
}

func (m *_BACnetTimeStampsEnclosed) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
