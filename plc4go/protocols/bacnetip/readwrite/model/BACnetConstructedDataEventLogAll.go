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

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetConstructedDataEventLogAll is the corresponding interface of BACnetConstructedDataEventLogAll
type BACnetConstructedDataEventLogAll interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// IsBACnetConstructedDataEventLogAll is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataEventLogAll()
	// CreateBuilder creates a BACnetConstructedDataEventLogAllBuilder
	CreateBACnetConstructedDataEventLogAllBuilder() BACnetConstructedDataEventLogAllBuilder
}

// _BACnetConstructedDataEventLogAll is the data-structure of this message
type _BACnetConstructedDataEventLogAll struct {
	BACnetConstructedDataContract
}

var _ BACnetConstructedDataEventLogAll = (*_BACnetConstructedDataEventLogAll)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataEventLogAll)(nil)

// NewBACnetConstructedDataEventLogAll factory function for _BACnetConstructedDataEventLogAll
func NewBACnetConstructedDataEventLogAll(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataEventLogAll {
	_result := &_BACnetConstructedDataEventLogAll{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataEventLogAllBuilder is a builder for BACnetConstructedDataEventLogAll
type BACnetConstructedDataEventLogAllBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() BACnetConstructedDataEventLogAllBuilder
	// Build builds the BACnetConstructedDataEventLogAll or returns an error if something is wrong
	Build() (BACnetConstructedDataEventLogAll, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataEventLogAll
}

// NewBACnetConstructedDataEventLogAllBuilder() creates a BACnetConstructedDataEventLogAllBuilder
func NewBACnetConstructedDataEventLogAllBuilder() BACnetConstructedDataEventLogAllBuilder {
	return &_BACnetConstructedDataEventLogAllBuilder{_BACnetConstructedDataEventLogAll: new(_BACnetConstructedDataEventLogAll)}
}

type _BACnetConstructedDataEventLogAllBuilder struct {
	*_BACnetConstructedDataEventLogAll

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataEventLogAllBuilder) = (*_BACnetConstructedDataEventLogAllBuilder)(nil)

func (b *_BACnetConstructedDataEventLogAllBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataEventLogAllBuilder) WithMandatoryFields() BACnetConstructedDataEventLogAllBuilder {
	return b
}

func (b *_BACnetConstructedDataEventLogAllBuilder) Build() (BACnetConstructedDataEventLogAll, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataEventLogAll.deepCopy(), nil
}

func (b *_BACnetConstructedDataEventLogAllBuilder) MustBuild() BACnetConstructedDataEventLogAll {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataEventLogAllBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataEventLogAllBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataEventLogAllBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataEventLogAllBuilder().(*_BACnetConstructedDataEventLogAllBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataEventLogAllBuilder creates a BACnetConstructedDataEventLogAllBuilder
func (b *_BACnetConstructedDataEventLogAll) CreateBACnetConstructedDataEventLogAllBuilder() BACnetConstructedDataEventLogAllBuilder {
	if b == nil {
		return NewBACnetConstructedDataEventLogAllBuilder()
	}
	return &_BACnetConstructedDataEventLogAllBuilder{_BACnetConstructedDataEventLogAll: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataEventLogAll) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_EVENT_LOG
}

func (m *_BACnetConstructedDataEventLogAll) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_ALL
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataEventLogAll) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataEventLogAll(structType any) BACnetConstructedDataEventLogAll {
	if casted, ok := structType.(BACnetConstructedDataEventLogAll); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataEventLogAll); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataEventLogAll) GetTypeName() string {
	return "BACnetConstructedDataEventLogAll"
}

func (m *_BACnetConstructedDataEventLogAll) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_BACnetConstructedDataEventLogAll) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataEventLogAll) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataEventLogAll BACnetConstructedDataEventLogAll, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataEventLogAll"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataEventLogAll")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Validation
	if !(bool((1) == (2))) {
		return nil, errors.WithStack(utils.ParseValidationError{Message: "All should never occur in context of constructed data. If it does please report"})
	}

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataEventLogAll"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataEventLogAll")
	}

	return m, nil
}

func (m *_BACnetConstructedDataEventLogAll) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataEventLogAll) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataEventLogAll"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataEventLogAll")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataEventLogAll"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataEventLogAll")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataEventLogAll) IsBACnetConstructedDataEventLogAll() {}

func (m *_BACnetConstructedDataEventLogAll) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataEventLogAll) deepCopy() *_BACnetConstructedDataEventLogAll {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataEventLogAllCopy := &_BACnetConstructedDataEventLogAll{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataEventLogAllCopy
}

func (m *_BACnetConstructedDataEventLogAll) String() string {
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
