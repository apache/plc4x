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

// SALDataMediaTransport is the corresponding interface of SALDataMediaTransport
type SALDataMediaTransport interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	SALData
	// GetMediaTransportControlData returns MediaTransportControlData (property field)
	GetMediaTransportControlData() MediaTransportControlData
	// IsSALDataMediaTransport is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsSALDataMediaTransport()
	// CreateBuilder creates a SALDataMediaTransportBuilder
	CreateSALDataMediaTransportBuilder() SALDataMediaTransportBuilder
}

// _SALDataMediaTransport is the data-structure of this message
type _SALDataMediaTransport struct {
	SALDataContract
	MediaTransportControlData MediaTransportControlData
}

var _ SALDataMediaTransport = (*_SALDataMediaTransport)(nil)
var _ SALDataRequirements = (*_SALDataMediaTransport)(nil)

// NewSALDataMediaTransport factory function for _SALDataMediaTransport
func NewSALDataMediaTransport(salData SALData, mediaTransportControlData MediaTransportControlData) *_SALDataMediaTransport {
	if mediaTransportControlData == nil {
		panic("mediaTransportControlData of type MediaTransportControlData for SALDataMediaTransport must not be nil")
	}
	_result := &_SALDataMediaTransport{
		SALDataContract:           NewSALData(salData),
		MediaTransportControlData: mediaTransportControlData,
	}
	_result.SALDataContract.(*_SALData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// SALDataMediaTransportBuilder is a builder for SALDataMediaTransport
type SALDataMediaTransportBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(mediaTransportControlData MediaTransportControlData) SALDataMediaTransportBuilder
	// WithMediaTransportControlData adds MediaTransportControlData (property field)
	WithMediaTransportControlData(MediaTransportControlData) SALDataMediaTransportBuilder
	// WithMediaTransportControlDataBuilder adds MediaTransportControlData (property field) which is build by the builder
	WithMediaTransportControlDataBuilder(func(MediaTransportControlDataBuilder) MediaTransportControlDataBuilder) SALDataMediaTransportBuilder
	// Build builds the SALDataMediaTransport or returns an error if something is wrong
	Build() (SALDataMediaTransport, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() SALDataMediaTransport
}

// NewSALDataMediaTransportBuilder() creates a SALDataMediaTransportBuilder
func NewSALDataMediaTransportBuilder() SALDataMediaTransportBuilder {
	return &_SALDataMediaTransportBuilder{_SALDataMediaTransport: new(_SALDataMediaTransport)}
}

type _SALDataMediaTransportBuilder struct {
	*_SALDataMediaTransport

	parentBuilder *_SALDataBuilder

	err *utils.MultiError
}

var _ (SALDataMediaTransportBuilder) = (*_SALDataMediaTransportBuilder)(nil)

func (b *_SALDataMediaTransportBuilder) setParent(contract SALDataContract) {
	b.SALDataContract = contract
}

func (b *_SALDataMediaTransportBuilder) WithMandatoryFields(mediaTransportControlData MediaTransportControlData) SALDataMediaTransportBuilder {
	return b.WithMediaTransportControlData(mediaTransportControlData)
}

func (b *_SALDataMediaTransportBuilder) WithMediaTransportControlData(mediaTransportControlData MediaTransportControlData) SALDataMediaTransportBuilder {
	b.MediaTransportControlData = mediaTransportControlData
	return b
}

func (b *_SALDataMediaTransportBuilder) WithMediaTransportControlDataBuilder(builderSupplier func(MediaTransportControlDataBuilder) MediaTransportControlDataBuilder) SALDataMediaTransportBuilder {
	builder := builderSupplier(b.MediaTransportControlData.CreateMediaTransportControlDataBuilder())
	var err error
	b.MediaTransportControlData, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "MediaTransportControlDataBuilder failed"))
	}
	return b
}

func (b *_SALDataMediaTransportBuilder) Build() (SALDataMediaTransport, error) {
	if b.MediaTransportControlData == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'mediaTransportControlData' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._SALDataMediaTransport.deepCopy(), nil
}

func (b *_SALDataMediaTransportBuilder) MustBuild() SALDataMediaTransport {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_SALDataMediaTransportBuilder) Done() SALDataBuilder {
	return b.parentBuilder
}

func (b *_SALDataMediaTransportBuilder) buildForSALData() (SALData, error) {
	return b.Build()
}

func (b *_SALDataMediaTransportBuilder) DeepCopy() any {
	_copy := b.CreateSALDataMediaTransportBuilder().(*_SALDataMediaTransportBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateSALDataMediaTransportBuilder creates a SALDataMediaTransportBuilder
func (b *_SALDataMediaTransport) CreateSALDataMediaTransportBuilder() SALDataMediaTransportBuilder {
	if b == nil {
		return NewSALDataMediaTransportBuilder()
	}
	return &_SALDataMediaTransportBuilder{_SALDataMediaTransport: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_SALDataMediaTransport) GetApplicationId() ApplicationId {
	return ApplicationId_MEDIA_TRANSPORT_CONTROL
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SALDataMediaTransport) GetParent() SALDataContract {
	return m.SALDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_SALDataMediaTransport) GetMediaTransportControlData() MediaTransportControlData {
	return m.MediaTransportControlData
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastSALDataMediaTransport(structType any) SALDataMediaTransport {
	if casted, ok := structType.(SALDataMediaTransport); ok {
		return casted
	}
	if casted, ok := structType.(*SALDataMediaTransport); ok {
		return *casted
	}
	return nil
}

func (m *_SALDataMediaTransport) GetTypeName() string {
	return "SALDataMediaTransport"
}

func (m *_SALDataMediaTransport) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.SALDataContract.(*_SALData).GetLengthInBits(ctx))

	// Simple field (mediaTransportControlData)
	lengthInBits += m.MediaTransportControlData.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_SALDataMediaTransport) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_SALDataMediaTransport) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_SALData, applicationId ApplicationId) (__sALDataMediaTransport SALDataMediaTransport, err error) {
	m.SALDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SALDataMediaTransport"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SALDataMediaTransport")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	mediaTransportControlData, err := ReadSimpleField[MediaTransportControlData](ctx, "mediaTransportControlData", ReadComplex[MediaTransportControlData](MediaTransportControlDataParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'mediaTransportControlData' field"))
	}
	m.MediaTransportControlData = mediaTransportControlData

	if closeErr := readBuffer.CloseContext("SALDataMediaTransport"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SALDataMediaTransport")
	}

	return m, nil
}

func (m *_SALDataMediaTransport) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_SALDataMediaTransport) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SALDataMediaTransport"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SALDataMediaTransport")
		}

		if err := WriteSimpleField[MediaTransportControlData](ctx, "mediaTransportControlData", m.GetMediaTransportControlData(), WriteComplex[MediaTransportControlData](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'mediaTransportControlData' field")
		}

		if popErr := writeBuffer.PopContext("SALDataMediaTransport"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SALDataMediaTransport")
		}
		return nil
	}
	return m.SALDataContract.(*_SALData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SALDataMediaTransport) IsSALDataMediaTransport() {}

func (m *_SALDataMediaTransport) DeepCopy() any {
	return m.deepCopy()
}

func (m *_SALDataMediaTransport) deepCopy() *_SALDataMediaTransport {
	if m == nil {
		return nil
	}
	_SALDataMediaTransportCopy := &_SALDataMediaTransport{
		m.SALDataContract.(*_SALData).deepCopy(),
		m.MediaTransportControlData.DeepCopy().(MediaTransportControlData),
	}
	m.SALDataContract.(*_SALData)._SubType = m
	return _SALDataMediaTransportCopy
}

func (m *_SALDataMediaTransport) String() string {
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
