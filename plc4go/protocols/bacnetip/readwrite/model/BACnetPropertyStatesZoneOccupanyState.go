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

// BACnetPropertyStatesZoneOccupanyState is the corresponding interface of BACnetPropertyStatesZoneOccupanyState
type BACnetPropertyStatesZoneOccupanyState interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetPropertyStates
	// GetZoneOccupanyState returns ZoneOccupanyState (property field)
	GetZoneOccupanyState() BACnetAccessZoneOccupancyStateTagged
	// IsBACnetPropertyStatesZoneOccupanyState is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetPropertyStatesZoneOccupanyState()
	// CreateBuilder creates a BACnetPropertyStatesZoneOccupanyStateBuilder
	CreateBACnetPropertyStatesZoneOccupanyStateBuilder() BACnetPropertyStatesZoneOccupanyStateBuilder
}

// _BACnetPropertyStatesZoneOccupanyState is the data-structure of this message
type _BACnetPropertyStatesZoneOccupanyState struct {
	BACnetPropertyStatesContract
	ZoneOccupanyState BACnetAccessZoneOccupancyStateTagged
}

var _ BACnetPropertyStatesZoneOccupanyState = (*_BACnetPropertyStatesZoneOccupanyState)(nil)
var _ BACnetPropertyStatesRequirements = (*_BACnetPropertyStatesZoneOccupanyState)(nil)

// NewBACnetPropertyStatesZoneOccupanyState factory function for _BACnetPropertyStatesZoneOccupanyState
func NewBACnetPropertyStatesZoneOccupanyState(peekedTagHeader BACnetTagHeader, zoneOccupanyState BACnetAccessZoneOccupancyStateTagged) *_BACnetPropertyStatesZoneOccupanyState {
	if zoneOccupanyState == nil {
		panic("zoneOccupanyState of type BACnetAccessZoneOccupancyStateTagged for BACnetPropertyStatesZoneOccupanyState must not be nil")
	}
	_result := &_BACnetPropertyStatesZoneOccupanyState{
		BACnetPropertyStatesContract: NewBACnetPropertyStates(peekedTagHeader),
		ZoneOccupanyState:            zoneOccupanyState,
	}
	_result.BACnetPropertyStatesContract.(*_BACnetPropertyStates)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetPropertyStatesZoneOccupanyStateBuilder is a builder for BACnetPropertyStatesZoneOccupanyState
type BACnetPropertyStatesZoneOccupanyStateBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(zoneOccupanyState BACnetAccessZoneOccupancyStateTagged) BACnetPropertyStatesZoneOccupanyStateBuilder
	// WithZoneOccupanyState adds ZoneOccupanyState (property field)
	WithZoneOccupanyState(BACnetAccessZoneOccupancyStateTagged) BACnetPropertyStatesZoneOccupanyStateBuilder
	// WithZoneOccupanyStateBuilder adds ZoneOccupanyState (property field) which is build by the builder
	WithZoneOccupanyStateBuilder(func(BACnetAccessZoneOccupancyStateTaggedBuilder) BACnetAccessZoneOccupancyStateTaggedBuilder) BACnetPropertyStatesZoneOccupanyStateBuilder
	// Build builds the BACnetPropertyStatesZoneOccupanyState or returns an error if something is wrong
	Build() (BACnetPropertyStatesZoneOccupanyState, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetPropertyStatesZoneOccupanyState
}

// NewBACnetPropertyStatesZoneOccupanyStateBuilder() creates a BACnetPropertyStatesZoneOccupanyStateBuilder
func NewBACnetPropertyStatesZoneOccupanyStateBuilder() BACnetPropertyStatesZoneOccupanyStateBuilder {
	return &_BACnetPropertyStatesZoneOccupanyStateBuilder{_BACnetPropertyStatesZoneOccupanyState: new(_BACnetPropertyStatesZoneOccupanyState)}
}

type _BACnetPropertyStatesZoneOccupanyStateBuilder struct {
	*_BACnetPropertyStatesZoneOccupanyState

	parentBuilder *_BACnetPropertyStatesBuilder

	err *utils.MultiError
}

var _ (BACnetPropertyStatesZoneOccupanyStateBuilder) = (*_BACnetPropertyStatesZoneOccupanyStateBuilder)(nil)

func (b *_BACnetPropertyStatesZoneOccupanyStateBuilder) setParent(contract BACnetPropertyStatesContract) {
	b.BACnetPropertyStatesContract = contract
}

func (b *_BACnetPropertyStatesZoneOccupanyStateBuilder) WithMandatoryFields(zoneOccupanyState BACnetAccessZoneOccupancyStateTagged) BACnetPropertyStatesZoneOccupanyStateBuilder {
	return b.WithZoneOccupanyState(zoneOccupanyState)
}

func (b *_BACnetPropertyStatesZoneOccupanyStateBuilder) WithZoneOccupanyState(zoneOccupanyState BACnetAccessZoneOccupancyStateTagged) BACnetPropertyStatesZoneOccupanyStateBuilder {
	b.ZoneOccupanyState = zoneOccupanyState
	return b
}

func (b *_BACnetPropertyStatesZoneOccupanyStateBuilder) WithZoneOccupanyStateBuilder(builderSupplier func(BACnetAccessZoneOccupancyStateTaggedBuilder) BACnetAccessZoneOccupancyStateTaggedBuilder) BACnetPropertyStatesZoneOccupanyStateBuilder {
	builder := builderSupplier(b.ZoneOccupanyState.CreateBACnetAccessZoneOccupancyStateTaggedBuilder())
	var err error
	b.ZoneOccupanyState, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetAccessZoneOccupancyStateTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetPropertyStatesZoneOccupanyStateBuilder) Build() (BACnetPropertyStatesZoneOccupanyState, error) {
	if b.ZoneOccupanyState == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'zoneOccupanyState' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetPropertyStatesZoneOccupanyState.deepCopy(), nil
}

func (b *_BACnetPropertyStatesZoneOccupanyStateBuilder) MustBuild() BACnetPropertyStatesZoneOccupanyState {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetPropertyStatesZoneOccupanyStateBuilder) Done() BACnetPropertyStatesBuilder {
	return b.parentBuilder
}

func (b *_BACnetPropertyStatesZoneOccupanyStateBuilder) buildForBACnetPropertyStates() (BACnetPropertyStates, error) {
	return b.Build()
}

func (b *_BACnetPropertyStatesZoneOccupanyStateBuilder) DeepCopy() any {
	_copy := b.CreateBACnetPropertyStatesZoneOccupanyStateBuilder().(*_BACnetPropertyStatesZoneOccupanyStateBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetPropertyStatesZoneOccupanyStateBuilder creates a BACnetPropertyStatesZoneOccupanyStateBuilder
func (b *_BACnetPropertyStatesZoneOccupanyState) CreateBACnetPropertyStatesZoneOccupanyStateBuilder() BACnetPropertyStatesZoneOccupanyStateBuilder {
	if b == nil {
		return NewBACnetPropertyStatesZoneOccupanyStateBuilder()
	}
	return &_BACnetPropertyStatesZoneOccupanyStateBuilder{_BACnetPropertyStatesZoneOccupanyState: b.deepCopy()}
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

func (m *_BACnetPropertyStatesZoneOccupanyState) GetParent() BACnetPropertyStatesContract {
	return m.BACnetPropertyStatesContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPropertyStatesZoneOccupanyState) GetZoneOccupanyState() BACnetAccessZoneOccupancyStateTagged {
	return m.ZoneOccupanyState
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetPropertyStatesZoneOccupanyState(structType any) BACnetPropertyStatesZoneOccupanyState {
	if casted, ok := structType.(BACnetPropertyStatesZoneOccupanyState); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPropertyStatesZoneOccupanyState); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPropertyStatesZoneOccupanyState) GetTypeName() string {
	return "BACnetPropertyStatesZoneOccupanyState"
}

func (m *_BACnetPropertyStatesZoneOccupanyState) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).GetLengthInBits(ctx))

	// Simple field (zoneOccupanyState)
	lengthInBits += m.ZoneOccupanyState.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetPropertyStatesZoneOccupanyState) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetPropertyStatesZoneOccupanyState) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetPropertyStates, peekedTagNumber uint8) (__bACnetPropertyStatesZoneOccupanyState BACnetPropertyStatesZoneOccupanyState, err error) {
	m.BACnetPropertyStatesContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPropertyStatesZoneOccupanyState"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPropertyStatesZoneOccupanyState")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	zoneOccupanyState, err := ReadSimpleField[BACnetAccessZoneOccupancyStateTagged](ctx, "zoneOccupanyState", ReadComplex[BACnetAccessZoneOccupancyStateTagged](BACnetAccessZoneOccupancyStateTaggedParseWithBufferProducer((uint8)(peekedTagNumber), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'zoneOccupanyState' field"))
	}
	m.ZoneOccupanyState = zoneOccupanyState

	if closeErr := readBuffer.CloseContext("BACnetPropertyStatesZoneOccupanyState"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPropertyStatesZoneOccupanyState")
	}

	return m, nil
}

func (m *_BACnetPropertyStatesZoneOccupanyState) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetPropertyStatesZoneOccupanyState) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPropertyStatesZoneOccupanyState"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPropertyStatesZoneOccupanyState")
		}

		if err := WriteSimpleField[BACnetAccessZoneOccupancyStateTagged](ctx, "zoneOccupanyState", m.GetZoneOccupanyState(), WriteComplex[BACnetAccessZoneOccupancyStateTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'zoneOccupanyState' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPropertyStatesZoneOccupanyState"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPropertyStatesZoneOccupanyState")
		}
		return nil
	}
	return m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetPropertyStatesZoneOccupanyState) IsBACnetPropertyStatesZoneOccupanyState() {}

func (m *_BACnetPropertyStatesZoneOccupanyState) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetPropertyStatesZoneOccupanyState) deepCopy() *_BACnetPropertyStatesZoneOccupanyState {
	if m == nil {
		return nil
	}
	_BACnetPropertyStatesZoneOccupanyStateCopy := &_BACnetPropertyStatesZoneOccupanyState{
		m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).deepCopy(),
		m.ZoneOccupanyState.DeepCopy().(BACnetAccessZoneOccupancyStateTagged),
	}
	m.BACnetPropertyStatesContract.(*_BACnetPropertyStates)._SubType = m
	return _BACnetPropertyStatesZoneOccupanyStateCopy
}

func (m *_BACnetPropertyStatesZoneOccupanyState) String() string {
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
