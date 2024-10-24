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

// BACnetConstructedDataFileAll is the corresponding interface of BACnetConstructedDataFileAll
type BACnetConstructedDataFileAll interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// IsBACnetConstructedDataFileAll is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataFileAll()
	// CreateBuilder creates a BACnetConstructedDataFileAllBuilder
	CreateBACnetConstructedDataFileAllBuilder() BACnetConstructedDataFileAllBuilder
}

// _BACnetConstructedDataFileAll is the data-structure of this message
type _BACnetConstructedDataFileAll struct {
	BACnetConstructedDataContract
}

var _ BACnetConstructedDataFileAll = (*_BACnetConstructedDataFileAll)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataFileAll)(nil)

// NewBACnetConstructedDataFileAll factory function for _BACnetConstructedDataFileAll
func NewBACnetConstructedDataFileAll(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataFileAll {
	_result := &_BACnetConstructedDataFileAll{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataFileAllBuilder is a builder for BACnetConstructedDataFileAll
type BACnetConstructedDataFileAllBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() BACnetConstructedDataFileAllBuilder
	// Build builds the BACnetConstructedDataFileAll or returns an error if something is wrong
	Build() (BACnetConstructedDataFileAll, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataFileAll
}

// NewBACnetConstructedDataFileAllBuilder() creates a BACnetConstructedDataFileAllBuilder
func NewBACnetConstructedDataFileAllBuilder() BACnetConstructedDataFileAllBuilder {
	return &_BACnetConstructedDataFileAllBuilder{_BACnetConstructedDataFileAll: new(_BACnetConstructedDataFileAll)}
}

type _BACnetConstructedDataFileAllBuilder struct {
	*_BACnetConstructedDataFileAll

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataFileAllBuilder) = (*_BACnetConstructedDataFileAllBuilder)(nil)

func (b *_BACnetConstructedDataFileAllBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataFileAllBuilder) WithMandatoryFields() BACnetConstructedDataFileAllBuilder {
	return b
}

func (b *_BACnetConstructedDataFileAllBuilder) Build() (BACnetConstructedDataFileAll, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataFileAll.deepCopy(), nil
}

func (b *_BACnetConstructedDataFileAllBuilder) MustBuild() BACnetConstructedDataFileAll {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataFileAllBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataFileAllBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataFileAllBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataFileAllBuilder().(*_BACnetConstructedDataFileAllBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataFileAllBuilder creates a BACnetConstructedDataFileAllBuilder
func (b *_BACnetConstructedDataFileAll) CreateBACnetConstructedDataFileAllBuilder() BACnetConstructedDataFileAllBuilder {
	if b == nil {
		return NewBACnetConstructedDataFileAllBuilder()
	}
	return &_BACnetConstructedDataFileAllBuilder{_BACnetConstructedDataFileAll: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataFileAll) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_FILE
}

func (m *_BACnetConstructedDataFileAll) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_ALL
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataFileAll) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataFileAll(structType any) BACnetConstructedDataFileAll {
	if casted, ok := structType.(BACnetConstructedDataFileAll); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataFileAll); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataFileAll) GetTypeName() string {
	return "BACnetConstructedDataFileAll"
}

func (m *_BACnetConstructedDataFileAll) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_BACnetConstructedDataFileAll) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataFileAll) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataFileAll BACnetConstructedDataFileAll, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataFileAll"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataFileAll")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Validation
	if !(bool((1) == (2))) {
		return nil, errors.WithStack(utils.ParseValidationError{Message: "All should never occur in context of constructed data. If it does please report"})
	}

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataFileAll"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataFileAll")
	}

	return m, nil
}

func (m *_BACnetConstructedDataFileAll) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataFileAll) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataFileAll"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataFileAll")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataFileAll"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataFileAll")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataFileAll) IsBACnetConstructedDataFileAll() {}

func (m *_BACnetConstructedDataFileAll) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataFileAll) deepCopy() *_BACnetConstructedDataFileAll {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataFileAllCopy := &_BACnetConstructedDataFileAll{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataFileAllCopy
}

func (m *_BACnetConstructedDataFileAll) String() string {
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
