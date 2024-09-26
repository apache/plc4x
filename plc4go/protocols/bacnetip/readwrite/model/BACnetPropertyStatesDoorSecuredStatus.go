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

// BACnetPropertyStatesDoorSecuredStatus is the corresponding interface of BACnetPropertyStatesDoorSecuredStatus
type BACnetPropertyStatesDoorSecuredStatus interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetPropertyStates
	// GetDoorSecuredStatus returns DoorSecuredStatus (property field)
	GetDoorSecuredStatus() BACnetDoorSecuredStatusTagged
	// IsBACnetPropertyStatesDoorSecuredStatus is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetPropertyStatesDoorSecuredStatus()
	// CreateBuilder creates a BACnetPropertyStatesDoorSecuredStatusBuilder
	CreateBACnetPropertyStatesDoorSecuredStatusBuilder() BACnetPropertyStatesDoorSecuredStatusBuilder
}

// _BACnetPropertyStatesDoorSecuredStatus is the data-structure of this message
type _BACnetPropertyStatesDoorSecuredStatus struct {
	BACnetPropertyStatesContract
	DoorSecuredStatus BACnetDoorSecuredStatusTagged
}

var _ BACnetPropertyStatesDoorSecuredStatus = (*_BACnetPropertyStatesDoorSecuredStatus)(nil)
var _ BACnetPropertyStatesRequirements = (*_BACnetPropertyStatesDoorSecuredStatus)(nil)

// NewBACnetPropertyStatesDoorSecuredStatus factory function for _BACnetPropertyStatesDoorSecuredStatus
func NewBACnetPropertyStatesDoorSecuredStatus(peekedTagHeader BACnetTagHeader, doorSecuredStatus BACnetDoorSecuredStatusTagged) *_BACnetPropertyStatesDoorSecuredStatus {
	if doorSecuredStatus == nil {
		panic("doorSecuredStatus of type BACnetDoorSecuredStatusTagged for BACnetPropertyStatesDoorSecuredStatus must not be nil")
	}
	_result := &_BACnetPropertyStatesDoorSecuredStatus{
		BACnetPropertyStatesContract: NewBACnetPropertyStates(peekedTagHeader),
		DoorSecuredStatus:            doorSecuredStatus,
	}
	_result.BACnetPropertyStatesContract.(*_BACnetPropertyStates)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetPropertyStatesDoorSecuredStatusBuilder is a builder for BACnetPropertyStatesDoorSecuredStatus
type BACnetPropertyStatesDoorSecuredStatusBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(doorSecuredStatus BACnetDoorSecuredStatusTagged) BACnetPropertyStatesDoorSecuredStatusBuilder
	// WithDoorSecuredStatus adds DoorSecuredStatus (property field)
	WithDoorSecuredStatus(BACnetDoorSecuredStatusTagged) BACnetPropertyStatesDoorSecuredStatusBuilder
	// WithDoorSecuredStatusBuilder adds DoorSecuredStatus (property field) which is build by the builder
	WithDoorSecuredStatusBuilder(func(BACnetDoorSecuredStatusTaggedBuilder) BACnetDoorSecuredStatusTaggedBuilder) BACnetPropertyStatesDoorSecuredStatusBuilder
	// Build builds the BACnetPropertyStatesDoorSecuredStatus or returns an error if something is wrong
	Build() (BACnetPropertyStatesDoorSecuredStatus, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetPropertyStatesDoorSecuredStatus
}

// NewBACnetPropertyStatesDoorSecuredStatusBuilder() creates a BACnetPropertyStatesDoorSecuredStatusBuilder
func NewBACnetPropertyStatesDoorSecuredStatusBuilder() BACnetPropertyStatesDoorSecuredStatusBuilder {
	return &_BACnetPropertyStatesDoorSecuredStatusBuilder{_BACnetPropertyStatesDoorSecuredStatus: new(_BACnetPropertyStatesDoorSecuredStatus)}
}

type _BACnetPropertyStatesDoorSecuredStatusBuilder struct {
	*_BACnetPropertyStatesDoorSecuredStatus

	parentBuilder *_BACnetPropertyStatesBuilder

	err *utils.MultiError
}

var _ (BACnetPropertyStatesDoorSecuredStatusBuilder) = (*_BACnetPropertyStatesDoorSecuredStatusBuilder)(nil)

func (b *_BACnetPropertyStatesDoorSecuredStatusBuilder) setParent(contract BACnetPropertyStatesContract) {
	b.BACnetPropertyStatesContract = contract
}

func (b *_BACnetPropertyStatesDoorSecuredStatusBuilder) WithMandatoryFields(doorSecuredStatus BACnetDoorSecuredStatusTagged) BACnetPropertyStatesDoorSecuredStatusBuilder {
	return b.WithDoorSecuredStatus(doorSecuredStatus)
}

func (b *_BACnetPropertyStatesDoorSecuredStatusBuilder) WithDoorSecuredStatus(doorSecuredStatus BACnetDoorSecuredStatusTagged) BACnetPropertyStatesDoorSecuredStatusBuilder {
	b.DoorSecuredStatus = doorSecuredStatus
	return b
}

func (b *_BACnetPropertyStatesDoorSecuredStatusBuilder) WithDoorSecuredStatusBuilder(builderSupplier func(BACnetDoorSecuredStatusTaggedBuilder) BACnetDoorSecuredStatusTaggedBuilder) BACnetPropertyStatesDoorSecuredStatusBuilder {
	builder := builderSupplier(b.DoorSecuredStatus.CreateBACnetDoorSecuredStatusTaggedBuilder())
	var err error
	b.DoorSecuredStatus, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetDoorSecuredStatusTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetPropertyStatesDoorSecuredStatusBuilder) Build() (BACnetPropertyStatesDoorSecuredStatus, error) {
	if b.DoorSecuredStatus == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'doorSecuredStatus' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetPropertyStatesDoorSecuredStatus.deepCopy(), nil
}

func (b *_BACnetPropertyStatesDoorSecuredStatusBuilder) MustBuild() BACnetPropertyStatesDoorSecuredStatus {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetPropertyStatesDoorSecuredStatusBuilder) Done() BACnetPropertyStatesBuilder {
	return b.parentBuilder
}

func (b *_BACnetPropertyStatesDoorSecuredStatusBuilder) buildForBACnetPropertyStates() (BACnetPropertyStates, error) {
	return b.Build()
}

func (b *_BACnetPropertyStatesDoorSecuredStatusBuilder) DeepCopy() any {
	_copy := b.CreateBACnetPropertyStatesDoorSecuredStatusBuilder().(*_BACnetPropertyStatesDoorSecuredStatusBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetPropertyStatesDoorSecuredStatusBuilder creates a BACnetPropertyStatesDoorSecuredStatusBuilder
func (b *_BACnetPropertyStatesDoorSecuredStatus) CreateBACnetPropertyStatesDoorSecuredStatusBuilder() BACnetPropertyStatesDoorSecuredStatusBuilder {
	if b == nil {
		return NewBACnetPropertyStatesDoorSecuredStatusBuilder()
	}
	return &_BACnetPropertyStatesDoorSecuredStatusBuilder{_BACnetPropertyStatesDoorSecuredStatus: b.deepCopy()}
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

func (m *_BACnetPropertyStatesDoorSecuredStatus) GetParent() BACnetPropertyStatesContract {
	return m.BACnetPropertyStatesContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPropertyStatesDoorSecuredStatus) GetDoorSecuredStatus() BACnetDoorSecuredStatusTagged {
	return m.DoorSecuredStatus
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetPropertyStatesDoorSecuredStatus(structType any) BACnetPropertyStatesDoorSecuredStatus {
	if casted, ok := structType.(BACnetPropertyStatesDoorSecuredStatus); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPropertyStatesDoorSecuredStatus); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPropertyStatesDoorSecuredStatus) GetTypeName() string {
	return "BACnetPropertyStatesDoorSecuredStatus"
}

func (m *_BACnetPropertyStatesDoorSecuredStatus) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).GetLengthInBits(ctx))

	// Simple field (doorSecuredStatus)
	lengthInBits += m.DoorSecuredStatus.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetPropertyStatesDoorSecuredStatus) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetPropertyStatesDoorSecuredStatus) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetPropertyStates, peekedTagNumber uint8) (__bACnetPropertyStatesDoorSecuredStatus BACnetPropertyStatesDoorSecuredStatus, err error) {
	m.BACnetPropertyStatesContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPropertyStatesDoorSecuredStatus"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPropertyStatesDoorSecuredStatus")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	doorSecuredStatus, err := ReadSimpleField[BACnetDoorSecuredStatusTagged](ctx, "doorSecuredStatus", ReadComplex[BACnetDoorSecuredStatusTagged](BACnetDoorSecuredStatusTaggedParseWithBufferProducer((uint8)(peekedTagNumber), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'doorSecuredStatus' field"))
	}
	m.DoorSecuredStatus = doorSecuredStatus

	if closeErr := readBuffer.CloseContext("BACnetPropertyStatesDoorSecuredStatus"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPropertyStatesDoorSecuredStatus")
	}

	return m, nil
}

func (m *_BACnetPropertyStatesDoorSecuredStatus) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetPropertyStatesDoorSecuredStatus) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPropertyStatesDoorSecuredStatus"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPropertyStatesDoorSecuredStatus")
		}

		if err := WriteSimpleField[BACnetDoorSecuredStatusTagged](ctx, "doorSecuredStatus", m.GetDoorSecuredStatus(), WriteComplex[BACnetDoorSecuredStatusTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'doorSecuredStatus' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPropertyStatesDoorSecuredStatus"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPropertyStatesDoorSecuredStatus")
		}
		return nil
	}
	return m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetPropertyStatesDoorSecuredStatus) IsBACnetPropertyStatesDoorSecuredStatus() {}

func (m *_BACnetPropertyStatesDoorSecuredStatus) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetPropertyStatesDoorSecuredStatus) deepCopy() *_BACnetPropertyStatesDoorSecuredStatus {
	if m == nil {
		return nil
	}
	_BACnetPropertyStatesDoorSecuredStatusCopy := &_BACnetPropertyStatesDoorSecuredStatus{
		m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).deepCopy(),
		m.DoorSecuredStatus.DeepCopy().(BACnetDoorSecuredStatusTagged),
	}
	m.BACnetPropertyStatesContract.(*_BACnetPropertyStates)._SubType = m
	return _BACnetPropertyStatesDoorSecuredStatusCopy
}

func (m *_BACnetPropertyStatesDoorSecuredStatus) String() string {
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
