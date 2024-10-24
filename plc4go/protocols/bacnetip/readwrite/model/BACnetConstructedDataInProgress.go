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

// BACnetConstructedDataInProgress is the corresponding interface of BACnetConstructedDataInProgress
type BACnetConstructedDataInProgress interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetInProgress returns InProgress (property field)
	GetInProgress() BACnetLightingInProgressTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetLightingInProgressTagged
	// IsBACnetConstructedDataInProgress is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataInProgress()
	// CreateBuilder creates a BACnetConstructedDataInProgressBuilder
	CreateBACnetConstructedDataInProgressBuilder() BACnetConstructedDataInProgressBuilder
}

// _BACnetConstructedDataInProgress is the data-structure of this message
type _BACnetConstructedDataInProgress struct {
	BACnetConstructedDataContract
	InProgress BACnetLightingInProgressTagged
}

var _ BACnetConstructedDataInProgress = (*_BACnetConstructedDataInProgress)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataInProgress)(nil)

// NewBACnetConstructedDataInProgress factory function for _BACnetConstructedDataInProgress
func NewBACnetConstructedDataInProgress(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, inProgress BACnetLightingInProgressTagged, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataInProgress {
	if inProgress == nil {
		panic("inProgress of type BACnetLightingInProgressTagged for BACnetConstructedDataInProgress must not be nil")
	}
	_result := &_BACnetConstructedDataInProgress{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		InProgress:                    inProgress,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataInProgressBuilder is a builder for BACnetConstructedDataInProgress
type BACnetConstructedDataInProgressBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(inProgress BACnetLightingInProgressTagged) BACnetConstructedDataInProgressBuilder
	// WithInProgress adds InProgress (property field)
	WithInProgress(BACnetLightingInProgressTagged) BACnetConstructedDataInProgressBuilder
	// WithInProgressBuilder adds InProgress (property field) which is build by the builder
	WithInProgressBuilder(func(BACnetLightingInProgressTaggedBuilder) BACnetLightingInProgressTaggedBuilder) BACnetConstructedDataInProgressBuilder
	// Build builds the BACnetConstructedDataInProgress or returns an error if something is wrong
	Build() (BACnetConstructedDataInProgress, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataInProgress
}

// NewBACnetConstructedDataInProgressBuilder() creates a BACnetConstructedDataInProgressBuilder
func NewBACnetConstructedDataInProgressBuilder() BACnetConstructedDataInProgressBuilder {
	return &_BACnetConstructedDataInProgressBuilder{_BACnetConstructedDataInProgress: new(_BACnetConstructedDataInProgress)}
}

type _BACnetConstructedDataInProgressBuilder struct {
	*_BACnetConstructedDataInProgress

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataInProgressBuilder) = (*_BACnetConstructedDataInProgressBuilder)(nil)

func (b *_BACnetConstructedDataInProgressBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataInProgressBuilder) WithMandatoryFields(inProgress BACnetLightingInProgressTagged) BACnetConstructedDataInProgressBuilder {
	return b.WithInProgress(inProgress)
}

func (b *_BACnetConstructedDataInProgressBuilder) WithInProgress(inProgress BACnetLightingInProgressTagged) BACnetConstructedDataInProgressBuilder {
	b.InProgress = inProgress
	return b
}

func (b *_BACnetConstructedDataInProgressBuilder) WithInProgressBuilder(builderSupplier func(BACnetLightingInProgressTaggedBuilder) BACnetLightingInProgressTaggedBuilder) BACnetConstructedDataInProgressBuilder {
	builder := builderSupplier(b.InProgress.CreateBACnetLightingInProgressTaggedBuilder())
	var err error
	b.InProgress, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetLightingInProgressTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataInProgressBuilder) Build() (BACnetConstructedDataInProgress, error) {
	if b.InProgress == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'inProgress' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataInProgress.deepCopy(), nil
}

func (b *_BACnetConstructedDataInProgressBuilder) MustBuild() BACnetConstructedDataInProgress {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataInProgressBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataInProgressBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataInProgressBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataInProgressBuilder().(*_BACnetConstructedDataInProgressBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataInProgressBuilder creates a BACnetConstructedDataInProgressBuilder
func (b *_BACnetConstructedDataInProgress) CreateBACnetConstructedDataInProgressBuilder() BACnetConstructedDataInProgressBuilder {
	if b == nil {
		return NewBACnetConstructedDataInProgressBuilder()
	}
	return &_BACnetConstructedDataInProgressBuilder{_BACnetConstructedDataInProgress: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataInProgress) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataInProgress) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_IN_PROGRESS
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataInProgress) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataInProgress) GetInProgress() BACnetLightingInProgressTagged {
	return m.InProgress
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataInProgress) GetActualValue() BACnetLightingInProgressTagged {
	ctx := context.Background()
	_ = ctx
	return CastBACnetLightingInProgressTagged(m.GetInProgress())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataInProgress(structType any) BACnetConstructedDataInProgress {
	if casted, ok := structType.(BACnetConstructedDataInProgress); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataInProgress); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataInProgress) GetTypeName() string {
	return "BACnetConstructedDataInProgress"
}

func (m *_BACnetConstructedDataInProgress) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (inProgress)
	lengthInBits += m.InProgress.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataInProgress) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataInProgress) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataInProgress BACnetConstructedDataInProgress, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataInProgress"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataInProgress")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	inProgress, err := ReadSimpleField[BACnetLightingInProgressTagged](ctx, "inProgress", ReadComplex[BACnetLightingInProgressTagged](BACnetLightingInProgressTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'inProgress' field"))
	}
	m.InProgress = inProgress

	actualValue, err := ReadVirtualField[BACnetLightingInProgressTagged](ctx, "actualValue", (*BACnetLightingInProgressTagged)(nil), inProgress)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataInProgress"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataInProgress")
	}

	return m, nil
}

func (m *_BACnetConstructedDataInProgress) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataInProgress) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataInProgress"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataInProgress")
		}

		if err := WriteSimpleField[BACnetLightingInProgressTagged](ctx, "inProgress", m.GetInProgress(), WriteComplex[BACnetLightingInProgressTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'inProgress' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataInProgress"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataInProgress")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataInProgress) IsBACnetConstructedDataInProgress() {}

func (m *_BACnetConstructedDataInProgress) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataInProgress) deepCopy() *_BACnetConstructedDataInProgress {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataInProgressCopy := &_BACnetConstructedDataInProgress{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.InProgress.DeepCopy().(BACnetLightingInProgressTagged),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataInProgressCopy
}

func (m *_BACnetConstructedDataInProgress) String() string {
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
