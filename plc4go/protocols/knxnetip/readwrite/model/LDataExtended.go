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

// LDataExtended is the corresponding interface of LDataExtended
type LDataExtended interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	LDataFrame
	// GetGroupAddress returns GroupAddress (property field)
	GetGroupAddress() bool
	// GetHopCount returns HopCount (property field)
	GetHopCount() uint8
	// GetExtendedFrameFormat returns ExtendedFrameFormat (property field)
	GetExtendedFrameFormat() uint8
	// GetSourceAddress returns SourceAddress (property field)
	GetSourceAddress() KnxAddress
	// GetDestinationAddress returns DestinationAddress (property field)
	GetDestinationAddress() []byte
	// GetApdu returns Apdu (property field)
	GetApdu() Apdu
	// IsLDataExtended is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsLDataExtended()
	// CreateBuilder creates a LDataExtendedBuilder
	CreateLDataExtendedBuilder() LDataExtendedBuilder
}

// _LDataExtended is the data-structure of this message
type _LDataExtended struct {
	LDataFrameContract
	GroupAddress        bool
	HopCount            uint8
	ExtendedFrameFormat uint8
	SourceAddress       KnxAddress
	DestinationAddress  []byte
	Apdu                Apdu
}

var _ LDataExtended = (*_LDataExtended)(nil)
var _ LDataFrameRequirements = (*_LDataExtended)(nil)

// NewLDataExtended factory function for _LDataExtended
func NewLDataExtended(frameType bool, notRepeated bool, priority CEMIPriority, acknowledgeRequested bool, errorFlag bool, groupAddress bool, hopCount uint8, extendedFrameFormat uint8, sourceAddress KnxAddress, destinationAddress []byte, apdu Apdu) *_LDataExtended {
	if sourceAddress == nil {
		panic("sourceAddress of type KnxAddress for LDataExtended must not be nil")
	}
	if apdu == nil {
		panic("apdu of type Apdu for LDataExtended must not be nil")
	}
	_result := &_LDataExtended{
		LDataFrameContract:  NewLDataFrame(frameType, notRepeated, priority, acknowledgeRequested, errorFlag),
		GroupAddress:        groupAddress,
		HopCount:            hopCount,
		ExtendedFrameFormat: extendedFrameFormat,
		SourceAddress:       sourceAddress,
		DestinationAddress:  destinationAddress,
		Apdu:                apdu,
	}
	_result.LDataFrameContract.(*_LDataFrame)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// LDataExtendedBuilder is a builder for LDataExtended
type LDataExtendedBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(groupAddress bool, hopCount uint8, extendedFrameFormat uint8, sourceAddress KnxAddress, destinationAddress []byte, apdu Apdu) LDataExtendedBuilder
	// WithGroupAddress adds GroupAddress (property field)
	WithGroupAddress(bool) LDataExtendedBuilder
	// WithHopCount adds HopCount (property field)
	WithHopCount(uint8) LDataExtendedBuilder
	// WithExtendedFrameFormat adds ExtendedFrameFormat (property field)
	WithExtendedFrameFormat(uint8) LDataExtendedBuilder
	// WithSourceAddress adds SourceAddress (property field)
	WithSourceAddress(KnxAddress) LDataExtendedBuilder
	// WithSourceAddressBuilder adds SourceAddress (property field) which is build by the builder
	WithSourceAddressBuilder(func(KnxAddressBuilder) KnxAddressBuilder) LDataExtendedBuilder
	// WithDestinationAddress adds DestinationAddress (property field)
	WithDestinationAddress(...byte) LDataExtendedBuilder
	// WithApdu adds Apdu (property field)
	WithApdu(Apdu) LDataExtendedBuilder
	// WithApduBuilder adds Apdu (property field) which is build by the builder
	WithApduBuilder(func(ApduBuilder) ApduBuilder) LDataExtendedBuilder
	// Build builds the LDataExtended or returns an error if something is wrong
	Build() (LDataExtended, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() LDataExtended
}

// NewLDataExtendedBuilder() creates a LDataExtendedBuilder
func NewLDataExtendedBuilder() LDataExtendedBuilder {
	return &_LDataExtendedBuilder{_LDataExtended: new(_LDataExtended)}
}

type _LDataExtendedBuilder struct {
	*_LDataExtended

	parentBuilder *_LDataFrameBuilder

	err *utils.MultiError
}

var _ (LDataExtendedBuilder) = (*_LDataExtendedBuilder)(nil)

func (b *_LDataExtendedBuilder) setParent(contract LDataFrameContract) {
	b.LDataFrameContract = contract
}

func (b *_LDataExtendedBuilder) WithMandatoryFields(groupAddress bool, hopCount uint8, extendedFrameFormat uint8, sourceAddress KnxAddress, destinationAddress []byte, apdu Apdu) LDataExtendedBuilder {
	return b.WithGroupAddress(groupAddress).WithHopCount(hopCount).WithExtendedFrameFormat(extendedFrameFormat).WithSourceAddress(sourceAddress).WithDestinationAddress(destinationAddress...).WithApdu(apdu)
}

func (b *_LDataExtendedBuilder) WithGroupAddress(groupAddress bool) LDataExtendedBuilder {
	b.GroupAddress = groupAddress
	return b
}

func (b *_LDataExtendedBuilder) WithHopCount(hopCount uint8) LDataExtendedBuilder {
	b.HopCount = hopCount
	return b
}

func (b *_LDataExtendedBuilder) WithExtendedFrameFormat(extendedFrameFormat uint8) LDataExtendedBuilder {
	b.ExtendedFrameFormat = extendedFrameFormat
	return b
}

func (b *_LDataExtendedBuilder) WithSourceAddress(sourceAddress KnxAddress) LDataExtendedBuilder {
	b.SourceAddress = sourceAddress
	return b
}

func (b *_LDataExtendedBuilder) WithSourceAddressBuilder(builderSupplier func(KnxAddressBuilder) KnxAddressBuilder) LDataExtendedBuilder {
	builder := builderSupplier(b.SourceAddress.CreateKnxAddressBuilder())
	var err error
	b.SourceAddress, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "KnxAddressBuilder failed"))
	}
	return b
}

func (b *_LDataExtendedBuilder) WithDestinationAddress(destinationAddress ...byte) LDataExtendedBuilder {
	b.DestinationAddress = destinationAddress
	return b
}

func (b *_LDataExtendedBuilder) WithApdu(apdu Apdu) LDataExtendedBuilder {
	b.Apdu = apdu
	return b
}

func (b *_LDataExtendedBuilder) WithApduBuilder(builderSupplier func(ApduBuilder) ApduBuilder) LDataExtendedBuilder {
	builder := builderSupplier(b.Apdu.CreateApduBuilder())
	var err error
	b.Apdu, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "ApduBuilder failed"))
	}
	return b
}

func (b *_LDataExtendedBuilder) Build() (LDataExtended, error) {
	if b.SourceAddress == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'sourceAddress' not set"))
	}
	if b.Apdu == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'apdu' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._LDataExtended.deepCopy(), nil
}

func (b *_LDataExtendedBuilder) MustBuild() LDataExtended {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_LDataExtendedBuilder) Done() LDataFrameBuilder {
	return b.parentBuilder
}

func (b *_LDataExtendedBuilder) buildForLDataFrame() (LDataFrame, error) {
	return b.Build()
}

func (b *_LDataExtendedBuilder) DeepCopy() any {
	_copy := b.CreateLDataExtendedBuilder().(*_LDataExtendedBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateLDataExtendedBuilder creates a LDataExtendedBuilder
func (b *_LDataExtended) CreateLDataExtendedBuilder() LDataExtendedBuilder {
	if b == nil {
		return NewLDataExtendedBuilder()
	}
	return &_LDataExtendedBuilder{_LDataExtended: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_LDataExtended) GetNotAckFrame() bool {
	return bool(true)
}

func (m *_LDataExtended) GetPolling() bool {
	return bool(false)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_LDataExtended) GetParent() LDataFrameContract {
	return m.LDataFrameContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_LDataExtended) GetGroupAddress() bool {
	return m.GroupAddress
}

func (m *_LDataExtended) GetHopCount() uint8 {
	return m.HopCount
}

func (m *_LDataExtended) GetExtendedFrameFormat() uint8 {
	return m.ExtendedFrameFormat
}

func (m *_LDataExtended) GetSourceAddress() KnxAddress {
	return m.SourceAddress
}

func (m *_LDataExtended) GetDestinationAddress() []byte {
	return m.DestinationAddress
}

func (m *_LDataExtended) GetApdu() Apdu {
	return m.Apdu
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastLDataExtended(structType any) LDataExtended {
	if casted, ok := structType.(LDataExtended); ok {
		return casted
	}
	if casted, ok := structType.(*LDataExtended); ok {
		return *casted
	}
	return nil
}

func (m *_LDataExtended) GetTypeName() string {
	return "LDataExtended"
}

func (m *_LDataExtended) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.LDataFrameContract.(*_LDataFrame).GetLengthInBits(ctx))

	// Simple field (groupAddress)
	lengthInBits += 1

	// Simple field (hopCount)
	lengthInBits += 3

	// Simple field (extendedFrameFormat)
	lengthInBits += 4

	// Simple field (sourceAddress)
	lengthInBits += m.SourceAddress.GetLengthInBits(ctx)

	// Array field
	if len(m.DestinationAddress) > 0 {
		lengthInBits += 8 * uint16(len(m.DestinationAddress))
	}

	// Implicit Field (dataLength)
	lengthInBits += 8

	// Simple field (apdu)
	lengthInBits += m.Apdu.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_LDataExtended) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_LDataExtended) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_LDataFrame) (__lDataExtended LDataExtended, err error) {
	m.LDataFrameContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("LDataExtended"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for LDataExtended")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	groupAddress, err := ReadSimpleField(ctx, "groupAddress", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'groupAddress' field"))
	}
	m.GroupAddress = groupAddress

	hopCount, err := ReadSimpleField(ctx, "hopCount", ReadUnsignedByte(readBuffer, uint8(3)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'hopCount' field"))
	}
	m.HopCount = hopCount

	extendedFrameFormat, err := ReadSimpleField(ctx, "extendedFrameFormat", ReadUnsignedByte(readBuffer, uint8(4)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'extendedFrameFormat' field"))
	}
	m.ExtendedFrameFormat = extendedFrameFormat

	sourceAddress, err := ReadSimpleField[KnxAddress](ctx, "sourceAddress", ReadComplex[KnxAddress](KnxAddressParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'sourceAddress' field"))
	}
	m.SourceAddress = sourceAddress

	destinationAddress, err := readBuffer.ReadByteArray("destinationAddress", int(int32(2)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'destinationAddress' field"))
	}
	m.DestinationAddress = destinationAddress

	dataLength, err := ReadImplicitField[uint8](ctx, "dataLength", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'dataLength' field"))
	}
	_ = dataLength

	apdu, err := ReadSimpleField[Apdu](ctx, "apdu", ReadComplex[Apdu](ApduParseWithBufferProducer[Apdu]((uint8)(dataLength)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'apdu' field"))
	}
	m.Apdu = apdu

	if closeErr := readBuffer.CloseContext("LDataExtended"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for LDataExtended")
	}

	return m, nil
}

func (m *_LDataExtended) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_LDataExtended) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("LDataExtended"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for LDataExtended")
		}

		if err := WriteSimpleField[bool](ctx, "groupAddress", m.GetGroupAddress(), WriteBoolean(writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'groupAddress' field")
		}

		if err := WriteSimpleField[uint8](ctx, "hopCount", m.GetHopCount(), WriteUnsignedByte(writeBuffer, 3)); err != nil {
			return errors.Wrap(err, "Error serializing 'hopCount' field")
		}

		if err := WriteSimpleField[uint8](ctx, "extendedFrameFormat", m.GetExtendedFrameFormat(), WriteUnsignedByte(writeBuffer, 4)); err != nil {
			return errors.Wrap(err, "Error serializing 'extendedFrameFormat' field")
		}

		if err := WriteSimpleField[KnxAddress](ctx, "sourceAddress", m.GetSourceAddress(), WriteComplex[KnxAddress](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'sourceAddress' field")
		}

		if err := WriteByteArrayField(ctx, "destinationAddress", m.GetDestinationAddress(), WriteByteArray(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'destinationAddress' field")
		}
		dataLength := uint8(uint8(m.GetApdu().GetLengthInBytes(ctx)) - uint8(uint8(1)))
		if err := WriteImplicitField(ctx, "dataLength", dataLength, WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'dataLength' field")
		}

		if err := WriteSimpleField[Apdu](ctx, "apdu", m.GetApdu(), WriteComplex[Apdu](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'apdu' field")
		}

		if popErr := writeBuffer.PopContext("LDataExtended"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for LDataExtended")
		}
		return nil
	}
	return m.LDataFrameContract.(*_LDataFrame).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_LDataExtended) IsLDataExtended() {}

func (m *_LDataExtended) DeepCopy() any {
	return m.deepCopy()
}

func (m *_LDataExtended) deepCopy() *_LDataExtended {
	if m == nil {
		return nil
	}
	_LDataExtendedCopy := &_LDataExtended{
		m.LDataFrameContract.(*_LDataFrame).deepCopy(),
		m.GroupAddress,
		m.HopCount,
		m.ExtendedFrameFormat,
		m.SourceAddress.DeepCopy().(KnxAddress),
		utils.DeepCopySlice[byte, byte](m.DestinationAddress),
		m.Apdu.DeepCopy().(Apdu),
	}
	m.LDataFrameContract.(*_LDataFrame)._SubType = m
	return _LDataExtendedCopy
}

func (m *_LDataExtended) String() string {
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
