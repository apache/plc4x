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

// BACnetPropertyStatesLiftCarDriveStatus is the corresponding interface of BACnetPropertyStatesLiftCarDriveStatus
type BACnetPropertyStatesLiftCarDriveStatus interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetPropertyStates
	// GetLiftCarDriveStatus returns LiftCarDriveStatus (property field)
	GetLiftCarDriveStatus() BACnetLiftCarDriveStatusTagged
	// IsBACnetPropertyStatesLiftCarDriveStatus is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetPropertyStatesLiftCarDriveStatus()
	// CreateBuilder creates a BACnetPropertyStatesLiftCarDriveStatusBuilder
	CreateBACnetPropertyStatesLiftCarDriveStatusBuilder() BACnetPropertyStatesLiftCarDriveStatusBuilder
}

// _BACnetPropertyStatesLiftCarDriveStatus is the data-structure of this message
type _BACnetPropertyStatesLiftCarDriveStatus struct {
	BACnetPropertyStatesContract
	LiftCarDriveStatus BACnetLiftCarDriveStatusTagged
}

var _ BACnetPropertyStatesLiftCarDriveStatus = (*_BACnetPropertyStatesLiftCarDriveStatus)(nil)
var _ BACnetPropertyStatesRequirements = (*_BACnetPropertyStatesLiftCarDriveStatus)(nil)

// NewBACnetPropertyStatesLiftCarDriveStatus factory function for _BACnetPropertyStatesLiftCarDriveStatus
func NewBACnetPropertyStatesLiftCarDriveStatus(peekedTagHeader BACnetTagHeader, liftCarDriveStatus BACnetLiftCarDriveStatusTagged) *_BACnetPropertyStatesLiftCarDriveStatus {
	if liftCarDriveStatus == nil {
		panic("liftCarDriveStatus of type BACnetLiftCarDriveStatusTagged for BACnetPropertyStatesLiftCarDriveStatus must not be nil")
	}
	_result := &_BACnetPropertyStatesLiftCarDriveStatus{
		BACnetPropertyStatesContract: NewBACnetPropertyStates(peekedTagHeader),
		LiftCarDriveStatus:           liftCarDriveStatus,
	}
	_result.BACnetPropertyStatesContract.(*_BACnetPropertyStates)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetPropertyStatesLiftCarDriveStatusBuilder is a builder for BACnetPropertyStatesLiftCarDriveStatus
type BACnetPropertyStatesLiftCarDriveStatusBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(liftCarDriveStatus BACnetLiftCarDriveStatusTagged) BACnetPropertyStatesLiftCarDriveStatusBuilder
	// WithLiftCarDriveStatus adds LiftCarDriveStatus (property field)
	WithLiftCarDriveStatus(BACnetLiftCarDriveStatusTagged) BACnetPropertyStatesLiftCarDriveStatusBuilder
	// WithLiftCarDriveStatusBuilder adds LiftCarDriveStatus (property field) which is build by the builder
	WithLiftCarDriveStatusBuilder(func(BACnetLiftCarDriveStatusTaggedBuilder) BACnetLiftCarDriveStatusTaggedBuilder) BACnetPropertyStatesLiftCarDriveStatusBuilder
	// Build builds the BACnetPropertyStatesLiftCarDriveStatus or returns an error if something is wrong
	Build() (BACnetPropertyStatesLiftCarDriveStatus, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetPropertyStatesLiftCarDriveStatus
}

// NewBACnetPropertyStatesLiftCarDriveStatusBuilder() creates a BACnetPropertyStatesLiftCarDriveStatusBuilder
func NewBACnetPropertyStatesLiftCarDriveStatusBuilder() BACnetPropertyStatesLiftCarDriveStatusBuilder {
	return &_BACnetPropertyStatesLiftCarDriveStatusBuilder{_BACnetPropertyStatesLiftCarDriveStatus: new(_BACnetPropertyStatesLiftCarDriveStatus)}
}

type _BACnetPropertyStatesLiftCarDriveStatusBuilder struct {
	*_BACnetPropertyStatesLiftCarDriveStatus

	parentBuilder *_BACnetPropertyStatesBuilder

	err *utils.MultiError
}

var _ (BACnetPropertyStatesLiftCarDriveStatusBuilder) = (*_BACnetPropertyStatesLiftCarDriveStatusBuilder)(nil)

func (b *_BACnetPropertyStatesLiftCarDriveStatusBuilder) setParent(contract BACnetPropertyStatesContract) {
	b.BACnetPropertyStatesContract = contract
}

func (b *_BACnetPropertyStatesLiftCarDriveStatusBuilder) WithMandatoryFields(liftCarDriveStatus BACnetLiftCarDriveStatusTagged) BACnetPropertyStatesLiftCarDriveStatusBuilder {
	return b.WithLiftCarDriveStatus(liftCarDriveStatus)
}

func (b *_BACnetPropertyStatesLiftCarDriveStatusBuilder) WithLiftCarDriveStatus(liftCarDriveStatus BACnetLiftCarDriveStatusTagged) BACnetPropertyStatesLiftCarDriveStatusBuilder {
	b.LiftCarDriveStatus = liftCarDriveStatus
	return b
}

func (b *_BACnetPropertyStatesLiftCarDriveStatusBuilder) WithLiftCarDriveStatusBuilder(builderSupplier func(BACnetLiftCarDriveStatusTaggedBuilder) BACnetLiftCarDriveStatusTaggedBuilder) BACnetPropertyStatesLiftCarDriveStatusBuilder {
	builder := builderSupplier(b.LiftCarDriveStatus.CreateBACnetLiftCarDriveStatusTaggedBuilder())
	var err error
	b.LiftCarDriveStatus, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetLiftCarDriveStatusTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetPropertyStatesLiftCarDriveStatusBuilder) Build() (BACnetPropertyStatesLiftCarDriveStatus, error) {
	if b.LiftCarDriveStatus == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'liftCarDriveStatus' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetPropertyStatesLiftCarDriveStatus.deepCopy(), nil
}

func (b *_BACnetPropertyStatesLiftCarDriveStatusBuilder) MustBuild() BACnetPropertyStatesLiftCarDriveStatus {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetPropertyStatesLiftCarDriveStatusBuilder) Done() BACnetPropertyStatesBuilder {
	return b.parentBuilder
}

func (b *_BACnetPropertyStatesLiftCarDriveStatusBuilder) buildForBACnetPropertyStates() (BACnetPropertyStates, error) {
	return b.Build()
}

func (b *_BACnetPropertyStatesLiftCarDriveStatusBuilder) DeepCopy() any {
	_copy := b.CreateBACnetPropertyStatesLiftCarDriveStatusBuilder().(*_BACnetPropertyStatesLiftCarDriveStatusBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetPropertyStatesLiftCarDriveStatusBuilder creates a BACnetPropertyStatesLiftCarDriveStatusBuilder
func (b *_BACnetPropertyStatesLiftCarDriveStatus) CreateBACnetPropertyStatesLiftCarDriveStatusBuilder() BACnetPropertyStatesLiftCarDriveStatusBuilder {
	if b == nil {
		return NewBACnetPropertyStatesLiftCarDriveStatusBuilder()
	}
	return &_BACnetPropertyStatesLiftCarDriveStatusBuilder{_BACnetPropertyStatesLiftCarDriveStatus: b.deepCopy()}
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

func (m *_BACnetPropertyStatesLiftCarDriveStatus) GetParent() BACnetPropertyStatesContract {
	return m.BACnetPropertyStatesContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPropertyStatesLiftCarDriveStatus) GetLiftCarDriveStatus() BACnetLiftCarDriveStatusTagged {
	return m.LiftCarDriveStatus
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetPropertyStatesLiftCarDriveStatus(structType any) BACnetPropertyStatesLiftCarDriveStatus {
	if casted, ok := structType.(BACnetPropertyStatesLiftCarDriveStatus); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPropertyStatesLiftCarDriveStatus); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) GetTypeName() string {
	return "BACnetPropertyStatesLiftCarDriveStatus"
}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).GetLengthInBits(ctx))

	// Simple field (liftCarDriveStatus)
	lengthInBits += m.LiftCarDriveStatus.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetPropertyStates, peekedTagNumber uint8) (__bACnetPropertyStatesLiftCarDriveStatus BACnetPropertyStatesLiftCarDriveStatus, err error) {
	m.BACnetPropertyStatesContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPropertyStatesLiftCarDriveStatus"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPropertyStatesLiftCarDriveStatus")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	liftCarDriveStatus, err := ReadSimpleField[BACnetLiftCarDriveStatusTagged](ctx, "liftCarDriveStatus", ReadComplex[BACnetLiftCarDriveStatusTagged](BACnetLiftCarDriveStatusTaggedParseWithBufferProducer((uint8)(peekedTagNumber), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'liftCarDriveStatus' field"))
	}
	m.LiftCarDriveStatus = liftCarDriveStatus

	if closeErr := readBuffer.CloseContext("BACnetPropertyStatesLiftCarDriveStatus"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPropertyStatesLiftCarDriveStatus")
	}

	return m, nil
}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPropertyStatesLiftCarDriveStatus"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPropertyStatesLiftCarDriveStatus")
		}

		if err := WriteSimpleField[BACnetLiftCarDriveStatusTagged](ctx, "liftCarDriveStatus", m.GetLiftCarDriveStatus(), WriteComplex[BACnetLiftCarDriveStatusTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'liftCarDriveStatus' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPropertyStatesLiftCarDriveStatus"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPropertyStatesLiftCarDriveStatus")
		}
		return nil
	}
	return m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) IsBACnetPropertyStatesLiftCarDriveStatus() {}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) deepCopy() *_BACnetPropertyStatesLiftCarDriveStatus {
	if m == nil {
		return nil
	}
	_BACnetPropertyStatesLiftCarDriveStatusCopy := &_BACnetPropertyStatesLiftCarDriveStatus{
		m.BACnetPropertyStatesContract.(*_BACnetPropertyStates).deepCopy(),
		m.LiftCarDriveStatus.DeepCopy().(BACnetLiftCarDriveStatusTagged),
	}
	m.BACnetPropertyStatesContract.(*_BACnetPropertyStates)._SubType = m
	return _BACnetPropertyStatesLiftCarDriveStatusCopy
}

func (m *_BACnetPropertyStatesLiftCarDriveStatus) String() string {
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
