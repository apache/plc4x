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

// BACnetShedLevelPercent is the corresponding interface of BACnetShedLevelPercent
type BACnetShedLevelPercent interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetShedLevel
	// GetPercent returns Percent (property field)
	GetPercent() BACnetContextTagUnsignedInteger
	// IsBACnetShedLevelPercent is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetShedLevelPercent()
	// CreateBuilder creates a BACnetShedLevelPercentBuilder
	CreateBACnetShedLevelPercentBuilder() BACnetShedLevelPercentBuilder
}

// _BACnetShedLevelPercent is the data-structure of this message
type _BACnetShedLevelPercent struct {
	BACnetShedLevelContract
	Percent BACnetContextTagUnsignedInteger
}

var _ BACnetShedLevelPercent = (*_BACnetShedLevelPercent)(nil)
var _ BACnetShedLevelRequirements = (*_BACnetShedLevelPercent)(nil)

// NewBACnetShedLevelPercent factory function for _BACnetShedLevelPercent
func NewBACnetShedLevelPercent(peekedTagHeader BACnetTagHeader, percent BACnetContextTagUnsignedInteger) *_BACnetShedLevelPercent {
	if percent == nil {
		panic("percent of type BACnetContextTagUnsignedInteger for BACnetShedLevelPercent must not be nil")
	}
	_result := &_BACnetShedLevelPercent{
		BACnetShedLevelContract: NewBACnetShedLevel(peekedTagHeader),
		Percent:                 percent,
	}
	_result.BACnetShedLevelContract.(*_BACnetShedLevel)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetShedLevelPercentBuilder is a builder for BACnetShedLevelPercent
type BACnetShedLevelPercentBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(percent BACnetContextTagUnsignedInteger) BACnetShedLevelPercentBuilder
	// WithPercent adds Percent (property field)
	WithPercent(BACnetContextTagUnsignedInteger) BACnetShedLevelPercentBuilder
	// WithPercentBuilder adds Percent (property field) which is build by the builder
	WithPercentBuilder(func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetShedLevelPercentBuilder
	// Build builds the BACnetShedLevelPercent or returns an error if something is wrong
	Build() (BACnetShedLevelPercent, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetShedLevelPercent
}

// NewBACnetShedLevelPercentBuilder() creates a BACnetShedLevelPercentBuilder
func NewBACnetShedLevelPercentBuilder() BACnetShedLevelPercentBuilder {
	return &_BACnetShedLevelPercentBuilder{_BACnetShedLevelPercent: new(_BACnetShedLevelPercent)}
}

type _BACnetShedLevelPercentBuilder struct {
	*_BACnetShedLevelPercent

	parentBuilder *_BACnetShedLevelBuilder

	err *utils.MultiError
}

var _ (BACnetShedLevelPercentBuilder) = (*_BACnetShedLevelPercentBuilder)(nil)

func (b *_BACnetShedLevelPercentBuilder) setParent(contract BACnetShedLevelContract) {
	b.BACnetShedLevelContract = contract
}

func (b *_BACnetShedLevelPercentBuilder) WithMandatoryFields(percent BACnetContextTagUnsignedInteger) BACnetShedLevelPercentBuilder {
	return b.WithPercent(percent)
}

func (b *_BACnetShedLevelPercentBuilder) WithPercent(percent BACnetContextTagUnsignedInteger) BACnetShedLevelPercentBuilder {
	b.Percent = percent
	return b
}

func (b *_BACnetShedLevelPercentBuilder) WithPercentBuilder(builderSupplier func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetShedLevelPercentBuilder {
	builder := builderSupplier(b.Percent.CreateBACnetContextTagUnsignedIntegerBuilder())
	var err error
	b.Percent, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetShedLevelPercentBuilder) Build() (BACnetShedLevelPercent, error) {
	if b.Percent == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'percent' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetShedLevelPercent.deepCopy(), nil
}

func (b *_BACnetShedLevelPercentBuilder) MustBuild() BACnetShedLevelPercent {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetShedLevelPercentBuilder) Done() BACnetShedLevelBuilder {
	return b.parentBuilder
}

func (b *_BACnetShedLevelPercentBuilder) buildForBACnetShedLevel() (BACnetShedLevel, error) {
	return b.Build()
}

func (b *_BACnetShedLevelPercentBuilder) DeepCopy() any {
	_copy := b.CreateBACnetShedLevelPercentBuilder().(*_BACnetShedLevelPercentBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetShedLevelPercentBuilder creates a BACnetShedLevelPercentBuilder
func (b *_BACnetShedLevelPercent) CreateBACnetShedLevelPercentBuilder() BACnetShedLevelPercentBuilder {
	if b == nil {
		return NewBACnetShedLevelPercentBuilder()
	}
	return &_BACnetShedLevelPercentBuilder{_BACnetShedLevelPercent: b.deepCopy()}
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

func (m *_BACnetShedLevelPercent) GetParent() BACnetShedLevelContract {
	return m.BACnetShedLevelContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetShedLevelPercent) GetPercent() BACnetContextTagUnsignedInteger {
	return m.Percent
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetShedLevelPercent(structType any) BACnetShedLevelPercent {
	if casted, ok := structType.(BACnetShedLevelPercent); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetShedLevelPercent); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetShedLevelPercent) GetTypeName() string {
	return "BACnetShedLevelPercent"
}

func (m *_BACnetShedLevelPercent) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetShedLevelContract.(*_BACnetShedLevel).GetLengthInBits(ctx))

	// Simple field (percent)
	lengthInBits += m.Percent.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetShedLevelPercent) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetShedLevelPercent) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetShedLevel) (__bACnetShedLevelPercent BACnetShedLevelPercent, err error) {
	m.BACnetShedLevelContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetShedLevelPercent"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetShedLevelPercent")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	percent, err := ReadSimpleField[BACnetContextTagUnsignedInteger](ctx, "percent", ReadComplex[BACnetContextTagUnsignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagUnsignedInteger]((uint8)(uint8(0)), (BACnetDataType)(BACnetDataType_UNSIGNED_INTEGER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'percent' field"))
	}
	m.Percent = percent

	if closeErr := readBuffer.CloseContext("BACnetShedLevelPercent"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetShedLevelPercent")
	}

	return m, nil
}

func (m *_BACnetShedLevelPercent) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetShedLevelPercent) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetShedLevelPercent"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetShedLevelPercent")
		}

		if err := WriteSimpleField[BACnetContextTagUnsignedInteger](ctx, "percent", m.GetPercent(), WriteComplex[BACnetContextTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'percent' field")
		}

		if popErr := writeBuffer.PopContext("BACnetShedLevelPercent"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetShedLevelPercent")
		}
		return nil
	}
	return m.BACnetShedLevelContract.(*_BACnetShedLevel).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetShedLevelPercent) IsBACnetShedLevelPercent() {}

func (m *_BACnetShedLevelPercent) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetShedLevelPercent) deepCopy() *_BACnetShedLevelPercent {
	if m == nil {
		return nil
	}
	_BACnetShedLevelPercentCopy := &_BACnetShedLevelPercent{
		m.BACnetShedLevelContract.(*_BACnetShedLevel).deepCopy(),
		m.Percent.DeepCopy().(BACnetContextTagUnsignedInteger),
	}
	m.BACnetShedLevelContract.(*_BACnetShedLevel)._SubType = m
	return _BACnetShedLevelPercentCopy
}

func (m *_BACnetShedLevelPercent) String() string {
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
