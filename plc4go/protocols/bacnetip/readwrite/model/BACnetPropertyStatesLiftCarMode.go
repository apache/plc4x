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

// BACnetPropertyStatesLiftCarMode is the corresponding interface of BACnetPropertyStatesLiftCarMode
type BACnetPropertyStatesLiftCarMode interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetPropertyStates
	// GetLiftCarMode returns LiftCarMode (property field)
	GetLiftCarMode() BACnetLiftCarModeTagged
	// IsBACnetPropertyStatesLiftCarMode is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetPropertyStatesLiftCarMode()
	// CreateBuilder creates a BACnetPropertyStatesLiftCarModeBuilder
	CreateBACnetPropertyStatesLiftCarModeBuilder() BACnetPropertyStatesLiftCarModeBuilder
}

// _BACnetPropertyStatesLiftCarMode is the data-structure of this message
type _BACnetPropertyStatesLiftCarMode struct {
	BACnetPropertyStatesContract
	LiftCarMode BACnetLiftCarModeTagged
}

var _ BACnetPropertyStatesLiftCarMode = (*_BACnetPropertyStatesLiftCarMode)(nil)
var _ BACnetPropertyStatesRequirements = (*_BACnetPropertyStatesLiftCarMode)(nil)

// NewBACnetPropertyStatesLiftCarMode factory function for _BACnetPropertyStatesLiftCarMode
func NewBACnetPropertyStatesLiftCarMode(peekedTagHeader BACnetTagHeader, liftCarMode BACnetLiftCarModeTagged) *_BACnetPropertyStatesLiftCarMode {
	if liftCarMode == nil {
		panic("liftCarMode of type BACnetLiftCarModeTagged for BACnetPropertyStatesLiftCarMode must not be nil")
	}
	_result := &_BACnetPropertyStatesLiftCarMode{
		BACnetPropertyStatesContract: NewBACnetPropertyStates(peekedTagHeader),
		LiftCarMode:                  liftCarMode,
	}
	_result.BACnetPropertyStatesContract.(*_BACnetPropertyStates)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetPropertyStatesLiftCarModeBuilder is a builder for BACnetPropertyStatesLiftCarMode
type BACnetPropertyStatesLiftCarModeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(liftCarMode BACnetLiftCarModeTagged) BACnetPropertyStatesLiftCarModeBuilder
	// WithLiftCarMode adds LiftCarMode (property field)
	WithLiftCarMode(BACnetLiftCarModeTagged) BACnetPropertyStatesLiftCarModeBuilder
	// WithLiftCarModeBuilder adds LiftCarMode (property field) which is build by the builder
	WithLiftCarModeBuilder(func(BACnetLiftCarModeTaggedBuilder) BACnetLiftCarModeTaggedBuilder) BACnetPropertyStatesLiftCarModeBuilder
	// Build builds the BACnetPropertyStatesLiftCarMode or returns an error if something is wrong
	Build() (BACnetPropertyStatesLiftCarMode, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetPropertyStatesLiftCarMode
}

// NewBACnetPropertyStatesLiftCarModeBuilder() creates a BACnetPropertyStatesLiftCarModeBuilder
func NewBACnetPropertyStatesLiftCarModeBuilder() BACnetPropertyStatesLiftCarModeBuilder {
	return &_BACnetPropertyStatesLiftCarModeBuilder{_BACnetPropertyStatesLiftCarMode: new(_BACnetPropertyStatesLiftCarMode)}
}

type _BACnetPropertyStatesLiftCarModeBuilder struct {
	*_BACnetPropertyStatesLiftCarMode

	parentBuilder *_BACnetPropertyStatesBuilder

	err *utils.MultiError
}

var _ (BACnetPropertyStatesLiftCarModeBuilder) = (*_BACnetPropertyStatesLiftCarModeBuilder)(nil)

func (b *_BACnetPropertyStatesLiftCarModeBuilder) setParent(contract BACnetPropertyStatesContract) {
	b.BACnetPropertyStatesContract = contract
}

func (b *_BACnetPropertyStatesLiftCarModeBuilder) WithMandatoryFields(liftCarMode BACnetLiftCarModeTagged) BACnetPropertyStatesLiftCarModeBuilder {
	return b.WithLiftCarMode(liftCarMode)
}

func (b *_BACnetPropertyStatesLiftCarModeBuilder) WithLiftCarMode(liftCarMode BACnetLiftCarModeTagged) BACnetPropertyStatesLiftCarModeBuilder {
	b.LiftCarMode = liftCarMode
	return b
}

func (b *_BACnetPropertyStatesLiftCarModeBuilder) WithLiftCarModeBuilder(builderSupplier func(BACnetLiftCarModeTaggedBuilder) BACnetLiftCarModeTaggedBuilder) BACnetPropertyStatesLiftCarModeBuilder {
	builder := builderSupplier(b.LiftCarMode.CreateBACnetLiftCarModeTaggedBuilder())
	var err error
	b.LiftCarMode, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetLiftCarModeTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetPropertyStatesLiftCarModeBuilder) Build() (BACnetPropertyStatesLiftCarMode, error) {
	if b.LiftCarMode == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'liftCarMode' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetPropertyStatesLiftCarMode.deepCopy(), nil
}

func (b *_BACnetPropertyStatesLiftCarModeBuilder) MustBuild() BACnetPropertyStatesLiftCarMode {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetPropertyStatesLiftCarModeBuilder) Done() BACnetPropertyStatesBuilder {
	return b.parentBuilder
}

func (b *_BACnetPropertyStatesLiftCarModeBuilder) buildForBACnetPropertyStates() (BACnetPropertyStates, error) {
	return b.Build()
}

func (b *_BACnetPropertyStatesLiftCarModeBuilder) DeepCopy() any {
	_copy := b.CreateBACnetPropertyStatesLiftCarModeBuilder().(*_BACnetPropertyStatesLiftCarModeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetPropertyStatesLiftCarModeBuilder creates a BACnetPropertyStatesLiftCarModeBuilder
func (b *_BACnetPropertyStatesLiftCarMode) CreateBACnetPropertyStatesLiftCarModeBuilder() BACnetPropertyStatesLiftCarModeBuilder {
	if b == nil {
		return NewBACnetPropertyStatesLiftCarModeBuilder()
	}
	return &_BACnetPropertyStatesLiftCarModeBuilder{_BACnetPropertyStatesLiftCarMode: b.deepCopy()}
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

func (m *_BACnetPropertyStatesLiftCarMode) GetParent() BACnetPropertyStatesContract {
	return m.BACnetPropertyStatesContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPropertyStatesLiftCarMode) GetLiftCarMode() BACnetLiftCarModeTagged {
	return m.LiftCarMode
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetPropertyStatesLiftCarMode(structType any) BACnetPropertyStatesLiftCarMode {
	if casted, ok := structType.(BACnetPropertyStatesLiftCarMode); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPropertyStatesLiftCarMode); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPropertyStatesLiftCarMode) GetTypeName() string {
	return "BACnetPropertyStatesLiftCarMode"
}

func (m *_BACnetPropertyStatesLiftCarMode) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).GetLengthInBits(ctx))

	// Simple field (liftCarMode)
	lengthInBits += m.LiftCarMode.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetPropertyStatesLiftCarMode) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetPropertyStatesLiftCarMode) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetPropertyStates, peekedTagNumber uint8) (__bACnetPropertyStatesLiftCarMode BACnetPropertyStatesLiftCarMode, err error) {
	m.BACnetPropertyStatesContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPropertyStatesLiftCarMode"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPropertyStatesLiftCarMode")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	liftCarMode, err := ReadSimpleField[BACnetLiftCarModeTagged](ctx, "liftCarMode", ReadComplex[BACnetLiftCarModeTagged](BACnetLiftCarModeTaggedParseWithBufferProducer((uint8)(peekedTagNumber), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'liftCarMode' field"))
	}
	m.LiftCarMode = liftCarMode

	if closeErr := readBuffer.CloseContext("BACnetPropertyStatesLiftCarMode"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPropertyStatesLiftCarMode")
	}

	return m, nil
}

func (m *_BACnetPropertyStatesLiftCarMode) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetPropertyStatesLiftCarMode) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPropertyStatesLiftCarMode"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPropertyStatesLiftCarMode")
		}

		if err := WriteSimpleField[BACnetLiftCarModeTagged](ctx, "liftCarMode", m.GetLiftCarMode(), WriteComplex[BACnetLiftCarModeTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'liftCarMode' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPropertyStatesLiftCarMode"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPropertyStatesLiftCarMode")
		}
		return nil
	}
	return m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetPropertyStatesLiftCarMode) IsBACnetPropertyStatesLiftCarMode() {}

func (m *_BACnetPropertyStatesLiftCarMode) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetPropertyStatesLiftCarMode) deepCopy() *_BACnetPropertyStatesLiftCarMode {
	if m == nil {
		return nil
	}
	_BACnetPropertyStatesLiftCarModeCopy := &_BACnetPropertyStatesLiftCarMode{
		m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).deepCopy(),
		m.LiftCarMode.DeepCopy().(BACnetLiftCarModeTagged),
	}
	m.BACnetPropertyStatesContract.(*_BACnetPropertyStates)._SubType = m
	return _BACnetPropertyStatesLiftCarModeCopy
}

func (m *_BACnetPropertyStatesLiftCarMode) String() string {
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
