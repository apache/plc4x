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

// MediaTransportControlDataSetTrack is the corresponding interface of MediaTransportControlDataSetTrack
type MediaTransportControlDataSetTrack interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	MediaTransportControlData
	// GetTrackMSB returns TrackMSB (property field)
	GetTrackMSB() byte
	// GetTrackMMSB returns TrackMMSB (property field)
	GetTrackMMSB() byte
	// GetTrackMLSB returns TrackMLSB (property field)
	GetTrackMLSB() byte
	// GetTrackLSB returns TrackLSB (property field)
	GetTrackLSB() byte
	// IsMediaTransportControlDataSetTrack is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsMediaTransportControlDataSetTrack()
	// CreateBuilder creates a MediaTransportControlDataSetTrackBuilder
	CreateMediaTransportControlDataSetTrackBuilder() MediaTransportControlDataSetTrackBuilder
}

// _MediaTransportControlDataSetTrack is the data-structure of this message
type _MediaTransportControlDataSetTrack struct {
	MediaTransportControlDataContract
	TrackMSB  byte
	TrackMMSB byte
	TrackMLSB byte
	TrackLSB  byte
}

var _ MediaTransportControlDataSetTrack = (*_MediaTransportControlDataSetTrack)(nil)
var _ MediaTransportControlDataRequirements = (*_MediaTransportControlDataSetTrack)(nil)

// NewMediaTransportControlDataSetTrack factory function for _MediaTransportControlDataSetTrack
func NewMediaTransportControlDataSetTrack(commandTypeContainer MediaTransportControlCommandTypeContainer, mediaLinkGroup byte, trackMSB byte, trackMMSB byte, trackMLSB byte, trackLSB byte) *_MediaTransportControlDataSetTrack {
	_result := &_MediaTransportControlDataSetTrack{
		MediaTransportControlDataContract: NewMediaTransportControlData(commandTypeContainer, mediaLinkGroup),
		TrackMSB:                          trackMSB,
		TrackMMSB:                         trackMMSB,
		TrackMLSB:                         trackMLSB,
		TrackLSB:                          trackLSB,
	}
	_result.MediaTransportControlDataContract.(*_MediaTransportControlData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// MediaTransportControlDataSetTrackBuilder is a builder for MediaTransportControlDataSetTrack
type MediaTransportControlDataSetTrackBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(trackMSB byte, trackMMSB byte, trackMLSB byte, trackLSB byte) MediaTransportControlDataSetTrackBuilder
	// WithTrackMSB adds TrackMSB (property field)
	WithTrackMSB(byte) MediaTransportControlDataSetTrackBuilder
	// WithTrackMMSB adds TrackMMSB (property field)
	WithTrackMMSB(byte) MediaTransportControlDataSetTrackBuilder
	// WithTrackMLSB adds TrackMLSB (property field)
	WithTrackMLSB(byte) MediaTransportControlDataSetTrackBuilder
	// WithTrackLSB adds TrackLSB (property field)
	WithTrackLSB(byte) MediaTransportControlDataSetTrackBuilder
	// Build builds the MediaTransportControlDataSetTrack or returns an error if something is wrong
	Build() (MediaTransportControlDataSetTrack, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() MediaTransportControlDataSetTrack
}

// NewMediaTransportControlDataSetTrackBuilder() creates a MediaTransportControlDataSetTrackBuilder
func NewMediaTransportControlDataSetTrackBuilder() MediaTransportControlDataSetTrackBuilder {
	return &_MediaTransportControlDataSetTrackBuilder{_MediaTransportControlDataSetTrack: new(_MediaTransportControlDataSetTrack)}
}

type _MediaTransportControlDataSetTrackBuilder struct {
	*_MediaTransportControlDataSetTrack

	parentBuilder *_MediaTransportControlDataBuilder

	err *utils.MultiError
}

var _ (MediaTransportControlDataSetTrackBuilder) = (*_MediaTransportControlDataSetTrackBuilder)(nil)

func (b *_MediaTransportControlDataSetTrackBuilder) setParent(contract MediaTransportControlDataContract) {
	b.MediaTransportControlDataContract = contract
}

func (b *_MediaTransportControlDataSetTrackBuilder) WithMandatoryFields(trackMSB byte, trackMMSB byte, trackMLSB byte, trackLSB byte) MediaTransportControlDataSetTrackBuilder {
	return b.WithTrackMSB(trackMSB).WithTrackMMSB(trackMMSB).WithTrackMLSB(trackMLSB).WithTrackLSB(trackLSB)
}

func (b *_MediaTransportControlDataSetTrackBuilder) WithTrackMSB(trackMSB byte) MediaTransportControlDataSetTrackBuilder {
	b.TrackMSB = trackMSB
	return b
}

func (b *_MediaTransportControlDataSetTrackBuilder) WithTrackMMSB(trackMMSB byte) MediaTransportControlDataSetTrackBuilder {
	b.TrackMMSB = trackMMSB
	return b
}

func (b *_MediaTransportControlDataSetTrackBuilder) WithTrackMLSB(trackMLSB byte) MediaTransportControlDataSetTrackBuilder {
	b.TrackMLSB = trackMLSB
	return b
}

func (b *_MediaTransportControlDataSetTrackBuilder) WithTrackLSB(trackLSB byte) MediaTransportControlDataSetTrackBuilder {
	b.TrackLSB = trackLSB
	return b
}

func (b *_MediaTransportControlDataSetTrackBuilder) Build() (MediaTransportControlDataSetTrack, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._MediaTransportControlDataSetTrack.deepCopy(), nil
}

func (b *_MediaTransportControlDataSetTrackBuilder) MustBuild() MediaTransportControlDataSetTrack {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_MediaTransportControlDataSetTrackBuilder) Done() MediaTransportControlDataBuilder {
	return b.parentBuilder
}

func (b *_MediaTransportControlDataSetTrackBuilder) buildForMediaTransportControlData() (MediaTransportControlData, error) {
	return b.Build()
}

func (b *_MediaTransportControlDataSetTrackBuilder) DeepCopy() any {
	_copy := b.CreateMediaTransportControlDataSetTrackBuilder().(*_MediaTransportControlDataSetTrackBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateMediaTransportControlDataSetTrackBuilder creates a MediaTransportControlDataSetTrackBuilder
func (b *_MediaTransportControlDataSetTrack) CreateMediaTransportControlDataSetTrackBuilder() MediaTransportControlDataSetTrackBuilder {
	if b == nil {
		return NewMediaTransportControlDataSetTrackBuilder()
	}
	return &_MediaTransportControlDataSetTrackBuilder{_MediaTransportControlDataSetTrack: b.deepCopy()}
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

func (m *_MediaTransportControlDataSetTrack) GetParent() MediaTransportControlDataContract {
	return m.MediaTransportControlDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_MediaTransportControlDataSetTrack) GetTrackMSB() byte {
	return m.TrackMSB
}

func (m *_MediaTransportControlDataSetTrack) GetTrackMMSB() byte {
	return m.TrackMMSB
}

func (m *_MediaTransportControlDataSetTrack) GetTrackMLSB() byte {
	return m.TrackMLSB
}

func (m *_MediaTransportControlDataSetTrack) GetTrackLSB() byte {
	return m.TrackLSB
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastMediaTransportControlDataSetTrack(structType any) MediaTransportControlDataSetTrack {
	if casted, ok := structType.(MediaTransportControlDataSetTrack); ok {
		return casted
	}
	if casted, ok := structType.(*MediaTransportControlDataSetTrack); ok {
		return *casted
	}
	return nil
}

func (m *_MediaTransportControlDataSetTrack) GetTypeName() string {
	return "MediaTransportControlDataSetTrack"
}

func (m *_MediaTransportControlDataSetTrack) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.MediaTransportControlDataContract.(*_MediaTransportControlData).GetLengthInBits(ctx))

	// Simple field (trackMSB)
	lengthInBits += 8

	// Simple field (trackMMSB)
	lengthInBits += 8

	// Simple field (trackMLSB)
	lengthInBits += 8

	// Simple field (trackLSB)
	lengthInBits += 8

	return lengthInBits
}

func (m *_MediaTransportControlDataSetTrack) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_MediaTransportControlDataSetTrack) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_MediaTransportControlData) (__mediaTransportControlDataSetTrack MediaTransportControlDataSetTrack, err error) {
	m.MediaTransportControlDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("MediaTransportControlDataSetTrack"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for MediaTransportControlDataSetTrack")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	trackMSB, err := ReadSimpleField(ctx, "trackMSB", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'trackMSB' field"))
	}
	m.TrackMSB = trackMSB

	trackMMSB, err := ReadSimpleField(ctx, "trackMMSB", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'trackMMSB' field"))
	}
	m.TrackMMSB = trackMMSB

	trackMLSB, err := ReadSimpleField(ctx, "trackMLSB", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'trackMLSB' field"))
	}
	m.TrackMLSB = trackMLSB

	trackLSB, err := ReadSimpleField(ctx, "trackLSB", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'trackLSB' field"))
	}
	m.TrackLSB = trackLSB

	if closeErr := readBuffer.CloseContext("MediaTransportControlDataSetTrack"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for MediaTransportControlDataSetTrack")
	}

	return m, nil
}

func (m *_MediaTransportControlDataSetTrack) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_MediaTransportControlDataSetTrack) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("MediaTransportControlDataSetTrack"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for MediaTransportControlDataSetTrack")
		}

		if err := WriteSimpleField[byte](ctx, "trackMSB", m.GetTrackMSB(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'trackMSB' field")
		}

		if err := WriteSimpleField[byte](ctx, "trackMMSB", m.GetTrackMMSB(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'trackMMSB' field")
		}

		if err := WriteSimpleField[byte](ctx, "trackMLSB", m.GetTrackMLSB(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'trackMLSB' field")
		}

		if err := WriteSimpleField[byte](ctx, "trackLSB", m.GetTrackLSB(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'trackLSB' field")
		}

		if popErr := writeBuffer.PopContext("MediaTransportControlDataSetTrack"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for MediaTransportControlDataSetTrack")
		}
		return nil
	}
	return m.MediaTransportControlDataContract.(*_MediaTransportControlData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_MediaTransportControlDataSetTrack) IsMediaTransportControlDataSetTrack() {}

func (m *_MediaTransportControlDataSetTrack) DeepCopy() any {
	return m.deepCopy()
}

func (m *_MediaTransportControlDataSetTrack) deepCopy() *_MediaTransportControlDataSetTrack {
	if m == nil {
		return nil
	}
	_MediaTransportControlDataSetTrackCopy := &_MediaTransportControlDataSetTrack{
		m.MediaTransportControlDataContract.(*_MediaTransportControlData).deepCopy(),
		m.TrackMSB,
		m.TrackMMSB,
		m.TrackMLSB,
		m.TrackLSB,
	}
	m.MediaTransportControlDataContract.(*_MediaTransportControlData)._SubType = m
	return _MediaTransportControlDataSetTrackCopy
}

func (m *_MediaTransportControlDataSetTrack) String() string {
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
