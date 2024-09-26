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

// MediaTransportControlDataShuffleOnOff is the corresponding interface of MediaTransportControlDataShuffleOnOff
type MediaTransportControlDataShuffleOnOff interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	MediaTransportControlData
	// GetState returns State (property field)
	GetState() byte
	// GetIsOff returns IsOff (virtual field)
	GetIsOff() bool
	// GetIsOn returns IsOn (virtual field)
	GetIsOn() bool
	// IsMediaTransportControlDataShuffleOnOff is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsMediaTransportControlDataShuffleOnOff()
	// CreateBuilder creates a MediaTransportControlDataShuffleOnOffBuilder
	CreateMediaTransportControlDataShuffleOnOffBuilder() MediaTransportControlDataShuffleOnOffBuilder
}

// _MediaTransportControlDataShuffleOnOff is the data-structure of this message
type _MediaTransportControlDataShuffleOnOff struct {
	MediaTransportControlDataContract
	State byte
}

var _ MediaTransportControlDataShuffleOnOff = (*_MediaTransportControlDataShuffleOnOff)(nil)
var _ MediaTransportControlDataRequirements = (*_MediaTransportControlDataShuffleOnOff)(nil)

// NewMediaTransportControlDataShuffleOnOff factory function for _MediaTransportControlDataShuffleOnOff
func NewMediaTransportControlDataShuffleOnOff(commandTypeContainer MediaTransportControlCommandTypeContainer, mediaLinkGroup byte, state byte) *_MediaTransportControlDataShuffleOnOff {
	_result := &_MediaTransportControlDataShuffleOnOff{
		MediaTransportControlDataContract: NewMediaTransportControlData(commandTypeContainer, mediaLinkGroup),
		State:                             state,
	}
	_result.MediaTransportControlDataContract.(*_MediaTransportControlData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// MediaTransportControlDataShuffleOnOffBuilder is a builder for MediaTransportControlDataShuffleOnOff
type MediaTransportControlDataShuffleOnOffBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(state byte) MediaTransportControlDataShuffleOnOffBuilder
	// WithState adds State (property field)
	WithState(byte) MediaTransportControlDataShuffleOnOffBuilder
	// Build builds the MediaTransportControlDataShuffleOnOff or returns an error if something is wrong
	Build() (MediaTransportControlDataShuffleOnOff, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() MediaTransportControlDataShuffleOnOff
}

// NewMediaTransportControlDataShuffleOnOffBuilder() creates a MediaTransportControlDataShuffleOnOffBuilder
func NewMediaTransportControlDataShuffleOnOffBuilder() MediaTransportControlDataShuffleOnOffBuilder {
	return &_MediaTransportControlDataShuffleOnOffBuilder{_MediaTransportControlDataShuffleOnOff: new(_MediaTransportControlDataShuffleOnOff)}
}

type _MediaTransportControlDataShuffleOnOffBuilder struct {
	*_MediaTransportControlDataShuffleOnOff

	parentBuilder *_MediaTransportControlDataBuilder

	err *utils.MultiError
}

var _ (MediaTransportControlDataShuffleOnOffBuilder) = (*_MediaTransportControlDataShuffleOnOffBuilder)(nil)

func (b *_MediaTransportControlDataShuffleOnOffBuilder) setParent(contract MediaTransportControlDataContract) {
	b.MediaTransportControlDataContract = contract
}

func (b *_MediaTransportControlDataShuffleOnOffBuilder) WithMandatoryFields(state byte) MediaTransportControlDataShuffleOnOffBuilder {
	return b.WithState(state)
}

func (b *_MediaTransportControlDataShuffleOnOffBuilder) WithState(state byte) MediaTransportControlDataShuffleOnOffBuilder {
	b.State = state
	return b
}

func (b *_MediaTransportControlDataShuffleOnOffBuilder) Build() (MediaTransportControlDataShuffleOnOff, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._MediaTransportControlDataShuffleOnOff.deepCopy(), nil
}

func (b *_MediaTransportControlDataShuffleOnOffBuilder) MustBuild() MediaTransportControlDataShuffleOnOff {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_MediaTransportControlDataShuffleOnOffBuilder) Done() MediaTransportControlDataBuilder {
	return b.parentBuilder
}

func (b *_MediaTransportControlDataShuffleOnOffBuilder) buildForMediaTransportControlData() (MediaTransportControlData, error) {
	return b.Build()
}

func (b *_MediaTransportControlDataShuffleOnOffBuilder) DeepCopy() any {
	_copy := b.CreateMediaTransportControlDataShuffleOnOffBuilder().(*_MediaTransportControlDataShuffleOnOffBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateMediaTransportControlDataShuffleOnOffBuilder creates a MediaTransportControlDataShuffleOnOffBuilder
func (b *_MediaTransportControlDataShuffleOnOff) CreateMediaTransportControlDataShuffleOnOffBuilder() MediaTransportControlDataShuffleOnOffBuilder {
	if b == nil {
		return NewMediaTransportControlDataShuffleOnOffBuilder()
	}
	return &_MediaTransportControlDataShuffleOnOffBuilder{_MediaTransportControlDataShuffleOnOff: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_MediaTransportControlDataShuffleOnOff) GetParent() MediaTransportControlDataContract {
	return m.MediaTransportControlDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_MediaTransportControlDataShuffleOnOff) GetState() byte {
	return m.State
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_MediaTransportControlDataShuffleOnOff) GetIsOff() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetState()) == (0x00)))
}

func (m *_MediaTransportControlDataShuffleOnOff) GetIsOn() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetState()) > (0xFE)))
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastMediaTransportControlDataShuffleOnOff(structType any) MediaTransportControlDataShuffleOnOff {
	if casted, ok := structType.(MediaTransportControlDataShuffleOnOff); ok {
		return casted
	}
	if casted, ok := structType.(*MediaTransportControlDataShuffleOnOff); ok {
		return *casted
	}
	return nil
}

func (m *_MediaTransportControlDataShuffleOnOff) GetTypeName() string {
	return "MediaTransportControlDataShuffleOnOff"
}

func (m *_MediaTransportControlDataShuffleOnOff) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.MediaTransportControlDataContract.(*_MediaTransportControlData).GetLengthInBits(ctx))

	// Simple field (state)
	lengthInBits += 8

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_MediaTransportControlDataShuffleOnOff) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_MediaTransportControlDataShuffleOnOff) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_MediaTransportControlData) (__mediaTransportControlDataShuffleOnOff MediaTransportControlDataShuffleOnOff, err error) {
	m.MediaTransportControlDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("MediaTransportControlDataShuffleOnOff"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for MediaTransportControlDataShuffleOnOff")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	state, err := ReadSimpleField(ctx, "state", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'state' field"))
	}
	m.State = state

	isOff, err := ReadVirtualField[bool](ctx, "isOff", (*bool)(nil), bool((state) == (0x00)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isOff' field"))
	}
	_ = isOff

	isOn, err := ReadVirtualField[bool](ctx, "isOn", (*bool)(nil), bool((state) > (0xFE)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isOn' field"))
	}
	_ = isOn

	if closeErr := readBuffer.CloseContext("MediaTransportControlDataShuffleOnOff"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for MediaTransportControlDataShuffleOnOff")
	}

	return m, nil
}

func (m *_MediaTransportControlDataShuffleOnOff) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_MediaTransportControlDataShuffleOnOff) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("MediaTransportControlDataShuffleOnOff"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for MediaTransportControlDataShuffleOnOff")
		}

		if err := WriteSimpleField[byte](ctx, "state", m.GetState(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'state' field")
		}
		// Virtual field
		isOff := m.GetIsOff()
		_ = isOff
		if _isOffErr := writeBuffer.WriteVirtual(ctx, "isOff", m.GetIsOff()); _isOffErr != nil {
			return errors.Wrap(_isOffErr, "Error serializing 'isOff' field")
		}
		// Virtual field
		isOn := m.GetIsOn()
		_ = isOn
		if _isOnErr := writeBuffer.WriteVirtual(ctx, "isOn", m.GetIsOn()); _isOnErr != nil {
			return errors.Wrap(_isOnErr, "Error serializing 'isOn' field")
		}

		if popErr := writeBuffer.PopContext("MediaTransportControlDataShuffleOnOff"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for MediaTransportControlDataShuffleOnOff")
		}
		return nil
	}
	return m.MediaTransportControlDataContract.(*_MediaTransportControlData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_MediaTransportControlDataShuffleOnOff) IsMediaTransportControlDataShuffleOnOff() {}

func (m *_MediaTransportControlDataShuffleOnOff) DeepCopy() any {
	return m.deepCopy()
}

func (m *_MediaTransportControlDataShuffleOnOff) deepCopy() *_MediaTransportControlDataShuffleOnOff {
	if m == nil {
		return nil
	}
	_MediaTransportControlDataShuffleOnOffCopy := &_MediaTransportControlDataShuffleOnOff{
		m.MediaTransportControlDataContract.(*_MediaTransportControlData).deepCopy(),
		m.State,
	}
	m.MediaTransportControlDataContract.(*_MediaTransportControlData)._SubType = m
	return _MediaTransportControlDataShuffleOnOffCopy
}

func (m *_MediaTransportControlDataShuffleOnOff) String() string {
	if m == nil {
		return "<nil>"
	}
	wb := utils.NewWriteBufferBoxBased(
		utils.WithWriteBufferBoxBasedMergeSingleBoxes(),
		utils.WithWriteBufferBoxBasedOmitEmptyBoxes(),
		utils.WithWriteBufferBoxBasedPrintPosLengthFooter(),
	)
	if err := wb.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return wb.GetBox().String()
}
