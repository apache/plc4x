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

// BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime is the corresponding interface of BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime
type BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetNotificationParametersChangeOfDiscreteValueNewValue
	// GetDateTimeValue returns DateTimeValue (property field)
	GetDateTimeValue() BACnetDateTimeEnclosed
	// IsBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime()
	// CreateBuilder creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder
	CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder
}

// _BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime is the data-structure of this message
type _BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime struct {
	BACnetNotificationParametersChangeOfDiscreteValueNewValueContract
	DateTimeValue BACnetDateTimeEnclosed
}

var _ BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime)(nil)
var _ BACnetNotificationParametersChangeOfDiscreteValueNewValueRequirements = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime)(nil)

// NewBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime factory function for _BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime
func NewBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, dateTimeValue BACnetDateTimeEnclosed, tagNumber uint8) *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime {
	if dateTimeValue == nil {
		panic("dateTimeValue of type BACnetDateTimeEnclosed for BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime must not be nil")
	}
	_result := &_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime{
		BACnetNotificationParametersChangeOfDiscreteValueNewValueContract: NewBACnetNotificationParametersChangeOfDiscreteValueNewValue(openingTag, peekedTagHeader, closingTag, tagNumber),
		DateTimeValue: dateTimeValue,
	}
	_result.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder is a builder for BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime
type BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(dateTimeValue BACnetDateTimeEnclosed) BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder
	// WithDateTimeValue adds DateTimeValue (property field)
	WithDateTimeValue(BACnetDateTimeEnclosed) BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder
	// WithDateTimeValueBuilder adds DateTimeValue (property field) which is build by the builder
	WithDateTimeValueBuilder(func(BACnetDateTimeEnclosedBuilder) BACnetDateTimeEnclosedBuilder) BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder
	// Build builds the BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime or returns an error if something is wrong
	Build() (BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime
}

// NewBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder() creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder
func NewBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder {
	return &_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder{_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime: new(_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime)}
}

type _BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder struct {
	*_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime

	parentBuilder *_BACnetNotificationParametersChangeOfDiscreteValueNewValueBuilder

	err *utils.MultiError
}

var _ (BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder)(nil)

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) setParent(contract BACnetNotificationParametersChangeOfDiscreteValueNewValueContract) {
	b.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract = contract
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) WithMandatoryFields(dateTimeValue BACnetDateTimeEnclosed) BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder {
	return b.WithDateTimeValue(dateTimeValue)
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) WithDateTimeValue(dateTimeValue BACnetDateTimeEnclosed) BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder {
	b.DateTimeValue = dateTimeValue
	return b
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) WithDateTimeValueBuilder(builderSupplier func(BACnetDateTimeEnclosedBuilder) BACnetDateTimeEnclosedBuilder) BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder {
	builder := builderSupplier(b.DateTimeValue.CreateBACnetDateTimeEnclosedBuilder())
	var err error
	b.DateTimeValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetDateTimeEnclosedBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) Build() (BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime, error) {
	if b.DateTimeValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'dateTimeValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime.deepCopy(), nil
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) MustBuild() BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) Done() BACnetNotificationParametersChangeOfDiscreteValueNewValueBuilder {
	return b.parentBuilder
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) buildForBACnetNotificationParametersChangeOfDiscreteValueNewValue() (BACnetNotificationParametersChangeOfDiscreteValueNewValue, error) {
	return b.Build()
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder) DeepCopy() any {
	_copy := b.CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder().(*_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder
func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder {
	if b == nil {
		return NewBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder()
	}
	return &_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeBuilder{_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime: b.deepCopy()}
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

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) GetParent() BACnetNotificationParametersChangeOfDiscreteValueNewValueContract {
	return m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) GetDateTimeValue() BACnetDateTimeEnclosed {
	return m.DateTimeValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime(structType any) BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime {
	if casted, ok := structType.(BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) GetTypeName() string {
	return "BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime"
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).GetLengthInBits(ctx))

	// Simple field (dateTimeValue)
	lengthInBits += m.DateTimeValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetNotificationParametersChangeOfDiscreteValueNewValue, tagNumber uint8) (__bACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime, err error) {
	m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	dateTimeValue, err := ReadSimpleField[BACnetDateTimeEnclosed](ctx, "dateTimeValue", ReadComplex[BACnetDateTimeEnclosed](BACnetDateTimeEnclosedParseWithBufferProducer((uint8)(uint8(0))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'dateTimeValue' field"))
	}
	m.DateTimeValue = dateTimeValue

	if closeErr := readBuffer.CloseContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime")
	}

	return m, nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime")
		}

		if err := WriteSimpleField[BACnetDateTimeEnclosed](ctx, "dateTimeValue", m.GetDateTimeValue(), WriteComplex[BACnetDateTimeEnclosed](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'dateTimeValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime")
		}
		return nil
	}
	return m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) IsBACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime() {
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) deepCopy() *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime {
	if m == nil {
		return nil
	}
	_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeCopy := &_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime{
		m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).deepCopy(),
		m.DateTimeValue.DeepCopy().(BACnetDateTimeEnclosed),
	}
	m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue)._SubType = m
	return _BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetimeCopy
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueDatetime) String() string {
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
