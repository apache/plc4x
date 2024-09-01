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

// BACnetLandingDoorStatusLandingDoorsListEntry is the corresponding interface of BACnetLandingDoorStatusLandingDoorsListEntry
type BACnetLandingDoorStatusLandingDoorsListEntry interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetFloorNumber returns FloorNumber (property field)
	GetFloorNumber() BACnetContextTagUnsignedInteger
	// GetDoorStatus returns DoorStatus (property field)
	GetDoorStatus() BACnetDoorStatusTagged
}

// BACnetLandingDoorStatusLandingDoorsListEntryExactly can be used when we want exactly this type and not a type which fulfills BACnetLandingDoorStatusLandingDoorsListEntry.
// This is useful for switch cases.
type BACnetLandingDoorStatusLandingDoorsListEntryExactly interface {
	BACnetLandingDoorStatusLandingDoorsListEntry
	isBACnetLandingDoorStatusLandingDoorsListEntry() bool
}

// _BACnetLandingDoorStatusLandingDoorsListEntry is the data-structure of this message
type _BACnetLandingDoorStatusLandingDoorsListEntry struct {
	FloorNumber BACnetContextTagUnsignedInteger
	DoorStatus  BACnetDoorStatusTagged
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetLandingDoorStatusLandingDoorsListEntry) GetFloorNumber() BACnetContextTagUnsignedInteger {
	return m.FloorNumber
}

func (m *_BACnetLandingDoorStatusLandingDoorsListEntry) GetDoorStatus() BACnetDoorStatusTagged {
	return m.DoorStatus
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetLandingDoorStatusLandingDoorsListEntry factory function for _BACnetLandingDoorStatusLandingDoorsListEntry
func NewBACnetLandingDoorStatusLandingDoorsListEntry(floorNumber BACnetContextTagUnsignedInteger, doorStatus BACnetDoorStatusTagged) *_BACnetLandingDoorStatusLandingDoorsListEntry {
	return &_BACnetLandingDoorStatusLandingDoorsListEntry{FloorNumber: floorNumber, DoorStatus: doorStatus}
}

// Deprecated: use the interface for direct cast
func CastBACnetLandingDoorStatusLandingDoorsListEntry(structType any) BACnetLandingDoorStatusLandingDoorsListEntry {
	if casted, ok := structType.(BACnetLandingDoorStatusLandingDoorsListEntry); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetLandingDoorStatusLandingDoorsListEntry); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetLandingDoorStatusLandingDoorsListEntry) GetTypeName() string {
	return "BACnetLandingDoorStatusLandingDoorsListEntry"
}

func (m *_BACnetLandingDoorStatusLandingDoorsListEntry) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (floorNumber)
	lengthInBits += m.FloorNumber.GetLengthInBits(ctx)

	// Simple field (doorStatus)
	lengthInBits += m.DoorStatus.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetLandingDoorStatusLandingDoorsListEntry) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetLandingDoorStatusLandingDoorsListEntryParse(ctx context.Context, theBytes []byte) (BACnetLandingDoorStatusLandingDoorsListEntry, error) {
	return BACnetLandingDoorStatusLandingDoorsListEntryParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func BACnetLandingDoorStatusLandingDoorsListEntryParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLandingDoorStatusLandingDoorsListEntry, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLandingDoorStatusLandingDoorsListEntry, error) {
		return BACnetLandingDoorStatusLandingDoorsListEntryParseWithBuffer(ctx, readBuffer)
	}
}

func BACnetLandingDoorStatusLandingDoorsListEntryParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLandingDoorStatusLandingDoorsListEntry, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetLandingDoorStatusLandingDoorsListEntry"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetLandingDoorStatusLandingDoorsListEntry")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	floorNumber, err := ReadSimpleField[BACnetContextTagUnsignedInteger](ctx, "floorNumber", ReadComplex[BACnetContextTagUnsignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagUnsignedInteger]((uint8)(uint8(0)), (BACnetDataType)(BACnetDataType_UNSIGNED_INTEGER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'floorNumber' field"))
	}

	doorStatus, err := ReadSimpleField[BACnetDoorStatusTagged](ctx, "doorStatus", ReadComplex[BACnetDoorStatusTagged](BACnetDoorStatusTaggedParseWithBufferProducer((uint8)(uint8(1)), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'doorStatus' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetLandingDoorStatusLandingDoorsListEntry"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetLandingDoorStatusLandingDoorsListEntry")
	}

	// Create the instance
	return &_BACnetLandingDoorStatusLandingDoorsListEntry{
		FloorNumber: floorNumber,
		DoorStatus:  doorStatus,
	}, nil
}

func (m *_BACnetLandingDoorStatusLandingDoorsListEntry) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetLandingDoorStatusLandingDoorsListEntry) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("BACnetLandingDoorStatusLandingDoorsListEntry"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetLandingDoorStatusLandingDoorsListEntry")
	}

	if err := WriteSimpleField[BACnetContextTagUnsignedInteger](ctx, "floorNumber", m.GetFloorNumber(), WriteComplex[BACnetContextTagUnsignedInteger](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'floorNumber' field")
	}

	if err := WriteSimpleField[BACnetDoorStatusTagged](ctx, "doorStatus", m.GetDoorStatus(), WriteComplex[BACnetDoorStatusTagged](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'doorStatus' field")
	}

	if popErr := writeBuffer.PopContext("BACnetLandingDoorStatusLandingDoorsListEntry"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetLandingDoorStatusLandingDoorsListEntry")
	}
	return nil
}

func (m *_BACnetLandingDoorStatusLandingDoorsListEntry) isBACnetLandingDoorStatusLandingDoorsListEntry() bool {
	return true
}

func (m *_BACnetLandingDoorStatusLandingDoorsListEntry) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
