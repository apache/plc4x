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

// MediaTransportControlDataFastForward is the corresponding interface of MediaTransportControlDataFastForward
type MediaTransportControlDataFastForward interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	MediaTransportControlData
	// GetOperation returns Operation (property field)
	GetOperation() byte
	// GetIsCeaseFastForward returns IsCeaseFastForward (virtual field)
	GetIsCeaseFastForward() bool
	// GetIs2x returns Is2x (virtual field)
	GetIs2x() bool
	// GetIs4x returns Is4x (virtual field)
	GetIs4x() bool
	// GetIs8x returns Is8x (virtual field)
	GetIs8x() bool
	// GetIs16x returns Is16x (virtual field)
	GetIs16x() bool
	// GetIs32x returns Is32x (virtual field)
	GetIs32x() bool
	// GetIs64x returns Is64x (virtual field)
	GetIs64x() bool
	// GetIsReserved returns IsReserved (virtual field)
	GetIsReserved() bool
	// IsMediaTransportControlDataFastForward is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsMediaTransportControlDataFastForward()
}

// _MediaTransportControlDataFastForward is the data-structure of this message
type _MediaTransportControlDataFastForward struct {
	MediaTransportControlDataContract
	Operation byte
}

var _ MediaTransportControlDataFastForward = (*_MediaTransportControlDataFastForward)(nil)
var _ MediaTransportControlDataRequirements = (*_MediaTransportControlDataFastForward)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_MediaTransportControlDataFastForward) GetParent() MediaTransportControlDataContract {
	return m.MediaTransportControlDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_MediaTransportControlDataFastForward) GetOperation() byte {
	return m.Operation
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_MediaTransportControlDataFastForward) GetIsCeaseFastForward() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetOperation()) == (0x00)))
}

func (m *_MediaTransportControlDataFastForward) GetIs2x() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetOperation()) == (0x02)))
}

func (m *_MediaTransportControlDataFastForward) GetIs4x() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetOperation()) == (0x04)))
}

func (m *_MediaTransportControlDataFastForward) GetIs8x() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetOperation()) == (0x06)))
}

func (m *_MediaTransportControlDataFastForward) GetIs16x() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetOperation()) == (0x08)))
}

func (m *_MediaTransportControlDataFastForward) GetIs32x() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetOperation()) == (0x0A)))
}

func (m *_MediaTransportControlDataFastForward) GetIs64x() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetOperation()) == (0x0C)))
}

func (m *_MediaTransportControlDataFastForward) GetIsReserved() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool(bool(bool(bool(bool(bool(!(m.GetIsCeaseFastForward())) && bool(!(m.GetIs2x()))) && bool(!(m.GetIs4x()))) && bool(!(m.GetIs8x()))) && bool(!(m.GetIs16x()))) && bool(!(m.GetIs32x()))) && bool(!(m.GetIs64x())))
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewMediaTransportControlDataFastForward factory function for _MediaTransportControlDataFastForward
func NewMediaTransportControlDataFastForward(operation byte, commandTypeContainer MediaTransportControlCommandTypeContainer, mediaLinkGroup byte) *_MediaTransportControlDataFastForward {
	_result := &_MediaTransportControlDataFastForward{
		MediaTransportControlDataContract: NewMediaTransportControlData(commandTypeContainer, mediaLinkGroup),
		Operation:                         operation,
	}
	_result.MediaTransportControlDataContract.(*_MediaTransportControlData)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastMediaTransportControlDataFastForward(structType any) MediaTransportControlDataFastForward {
	if casted, ok := structType.(MediaTransportControlDataFastForward); ok {
		return casted
	}
	if casted, ok := structType.(*MediaTransportControlDataFastForward); ok {
		return *casted
	}
	return nil
}

func (m *_MediaTransportControlDataFastForward) GetTypeName() string {
	return "MediaTransportControlDataFastForward"
}

func (m *_MediaTransportControlDataFastForward) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.MediaTransportControlDataContract.(*_MediaTransportControlData).getLengthInBits(ctx))

	// Simple field (operation)
	lengthInBits += 8

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_MediaTransportControlDataFastForward) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_MediaTransportControlDataFastForward) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_MediaTransportControlData) (__mediaTransportControlDataFastForward MediaTransportControlDataFastForward, err error) {
	m.MediaTransportControlDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("MediaTransportControlDataFastForward"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for MediaTransportControlDataFastForward")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	operation, err := ReadSimpleField(ctx, "operation", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'operation' field"))
	}
	m.Operation = operation

	isCeaseFastForward, err := ReadVirtualField[bool](ctx, "isCeaseFastForward", (*bool)(nil), bool((operation) == (0x00)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isCeaseFastForward' field"))
	}
	_ = isCeaseFastForward

	is2x, err := ReadVirtualField[bool](ctx, "is2x", (*bool)(nil), bool((operation) == (0x02)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'is2x' field"))
	}
	_ = is2x

	is4x, err := ReadVirtualField[bool](ctx, "is4x", (*bool)(nil), bool((operation) == (0x04)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'is4x' field"))
	}
	_ = is4x

	is8x, err := ReadVirtualField[bool](ctx, "is8x", (*bool)(nil), bool((operation) == (0x06)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'is8x' field"))
	}
	_ = is8x

	is16x, err := ReadVirtualField[bool](ctx, "is16x", (*bool)(nil), bool((operation) == (0x08)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'is16x' field"))
	}
	_ = is16x

	is32x, err := ReadVirtualField[bool](ctx, "is32x", (*bool)(nil), bool((operation) == (0x0A)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'is32x' field"))
	}
	_ = is32x

	is64x, err := ReadVirtualField[bool](ctx, "is64x", (*bool)(nil), bool((operation) == (0x0C)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'is64x' field"))
	}
	_ = is64x

	isReserved, err := ReadVirtualField[bool](ctx, "isReserved", (*bool)(nil), bool(bool(bool(bool(bool(bool(!(isCeaseFastForward)) && bool(!(is2x))) && bool(!(is4x))) && bool(!(is8x))) && bool(!(is16x))) && bool(!(is32x))) && bool(!(is64x)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isReserved' field"))
	}
	_ = isReserved

	if closeErr := readBuffer.CloseContext("MediaTransportControlDataFastForward"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for MediaTransportControlDataFastForward")
	}

	return m, nil
}

func (m *_MediaTransportControlDataFastForward) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_MediaTransportControlDataFastForward) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("MediaTransportControlDataFastForward"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for MediaTransportControlDataFastForward")
		}

		if err := WriteSimpleField[byte](ctx, "operation", m.GetOperation(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'operation' field")
		}
		// Virtual field
		isCeaseFastForward := m.GetIsCeaseFastForward()
		_ = isCeaseFastForward
		if _isCeaseFastForwardErr := writeBuffer.WriteVirtual(ctx, "isCeaseFastForward", m.GetIsCeaseFastForward()); _isCeaseFastForwardErr != nil {
			return errors.Wrap(_isCeaseFastForwardErr, "Error serializing 'isCeaseFastForward' field")
		}
		// Virtual field
		is2x := m.GetIs2x()
		_ = is2x
		if _is2xErr := writeBuffer.WriteVirtual(ctx, "is2x", m.GetIs2x()); _is2xErr != nil {
			return errors.Wrap(_is2xErr, "Error serializing 'is2x' field")
		}
		// Virtual field
		is4x := m.GetIs4x()
		_ = is4x
		if _is4xErr := writeBuffer.WriteVirtual(ctx, "is4x", m.GetIs4x()); _is4xErr != nil {
			return errors.Wrap(_is4xErr, "Error serializing 'is4x' field")
		}
		// Virtual field
		is8x := m.GetIs8x()
		_ = is8x
		if _is8xErr := writeBuffer.WriteVirtual(ctx, "is8x", m.GetIs8x()); _is8xErr != nil {
			return errors.Wrap(_is8xErr, "Error serializing 'is8x' field")
		}
		// Virtual field
		is16x := m.GetIs16x()
		_ = is16x
		if _is16xErr := writeBuffer.WriteVirtual(ctx, "is16x", m.GetIs16x()); _is16xErr != nil {
			return errors.Wrap(_is16xErr, "Error serializing 'is16x' field")
		}
		// Virtual field
		is32x := m.GetIs32x()
		_ = is32x
		if _is32xErr := writeBuffer.WriteVirtual(ctx, "is32x", m.GetIs32x()); _is32xErr != nil {
			return errors.Wrap(_is32xErr, "Error serializing 'is32x' field")
		}
		// Virtual field
		is64x := m.GetIs64x()
		_ = is64x
		if _is64xErr := writeBuffer.WriteVirtual(ctx, "is64x", m.GetIs64x()); _is64xErr != nil {
			return errors.Wrap(_is64xErr, "Error serializing 'is64x' field")
		}
		// Virtual field
		isReserved := m.GetIsReserved()
		_ = isReserved
		if _isReservedErr := writeBuffer.WriteVirtual(ctx, "isReserved", m.GetIsReserved()); _isReservedErr != nil {
			return errors.Wrap(_isReservedErr, "Error serializing 'isReserved' field")
		}

		if popErr := writeBuffer.PopContext("MediaTransportControlDataFastForward"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for MediaTransportControlDataFastForward")
		}
		return nil
	}
	return m.MediaTransportControlDataContract.(*_MediaTransportControlData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_MediaTransportControlDataFastForward) IsMediaTransportControlDataFastForward() {}

func (m *_MediaTransportControlDataFastForward) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
