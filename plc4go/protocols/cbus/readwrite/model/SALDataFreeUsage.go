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

// SALDataFreeUsage is the corresponding interface of SALDataFreeUsage
type SALDataFreeUsage interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	SALData
	// IsSALDataFreeUsage is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsSALDataFreeUsage()
	// CreateBuilder creates a SALDataFreeUsageBuilder
	CreateSALDataFreeUsageBuilder() SALDataFreeUsageBuilder
}

// _SALDataFreeUsage is the data-structure of this message
type _SALDataFreeUsage struct {
	SALDataContract
}

var _ SALDataFreeUsage = (*_SALDataFreeUsage)(nil)
var _ SALDataRequirements = (*_SALDataFreeUsage)(nil)

// NewSALDataFreeUsage factory function for _SALDataFreeUsage
func NewSALDataFreeUsage(salData SALData) *_SALDataFreeUsage {
	_result := &_SALDataFreeUsage{
		SALDataContract: NewSALData(salData),
	}
	_result.SALDataContract.(*_SALData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// SALDataFreeUsageBuilder is a builder for SALDataFreeUsage
type SALDataFreeUsageBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() SALDataFreeUsageBuilder
	// Build builds the SALDataFreeUsage or returns an error if something is wrong
	Build() (SALDataFreeUsage, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() SALDataFreeUsage
}

// NewSALDataFreeUsageBuilder() creates a SALDataFreeUsageBuilder
func NewSALDataFreeUsageBuilder() SALDataFreeUsageBuilder {
	return &_SALDataFreeUsageBuilder{_SALDataFreeUsage: new(_SALDataFreeUsage)}
}

type _SALDataFreeUsageBuilder struct {
	*_SALDataFreeUsage

	parentBuilder *_SALDataBuilder

	err *utils.MultiError
}

var _ (SALDataFreeUsageBuilder) = (*_SALDataFreeUsageBuilder)(nil)

func (b *_SALDataFreeUsageBuilder) setParent(contract SALDataContract) {
	b.SALDataContract = contract
}

func (b *_SALDataFreeUsageBuilder) WithMandatoryFields() SALDataFreeUsageBuilder {
	return b
}

func (b *_SALDataFreeUsageBuilder) Build() (SALDataFreeUsage, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._SALDataFreeUsage.deepCopy(), nil
}

func (b *_SALDataFreeUsageBuilder) MustBuild() SALDataFreeUsage {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_SALDataFreeUsageBuilder) Done() SALDataBuilder {
	return b.parentBuilder
}

func (b *_SALDataFreeUsageBuilder) buildForSALData() (SALData, error) {
	return b.Build()
}

func (b *_SALDataFreeUsageBuilder) DeepCopy() any {
	_copy := b.CreateSALDataFreeUsageBuilder().(*_SALDataFreeUsageBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateSALDataFreeUsageBuilder creates a SALDataFreeUsageBuilder
func (b *_SALDataFreeUsage) CreateSALDataFreeUsageBuilder() SALDataFreeUsageBuilder {
	if b == nil {
		return NewSALDataFreeUsageBuilder()
	}
	return &_SALDataFreeUsageBuilder{_SALDataFreeUsage: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_SALDataFreeUsage) GetApplicationId() ApplicationId {
	return ApplicationId_FREE_USAGE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SALDataFreeUsage) GetParent() SALDataContract {
	return m.SALDataContract
}

// Deprecated: use the interface for direct cast
func CastSALDataFreeUsage(structType any) SALDataFreeUsage {
	if casted, ok := structType.(SALDataFreeUsage); ok {
		return casted
	}
	if casted, ok := structType.(*SALDataFreeUsage); ok {
		return *casted
	}
	return nil
}

func (m *_SALDataFreeUsage) GetTypeName() string {
	return "SALDataFreeUsage"
}

func (m *_SALDataFreeUsage) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.SALDataContract.(*_SALData).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_SALDataFreeUsage) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_SALDataFreeUsage) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_SALData, applicationId ApplicationId) (__sALDataFreeUsage SALDataFreeUsage, err error) {
	m.SALDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SALDataFreeUsage"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SALDataFreeUsage")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Validation
	if !(bool((1) == (2))) {
		return nil, errors.WithStack(utils.ParseValidationError{Message: "FREE_USAGE Not yet implemented"})
	}

	if closeErr := readBuffer.CloseContext("SALDataFreeUsage"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SALDataFreeUsage")
	}

	return m, nil
}

func (m *_SALDataFreeUsage) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_SALDataFreeUsage) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SALDataFreeUsage"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SALDataFreeUsage")
		}

		if popErr := writeBuffer.PopContext("SALDataFreeUsage"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SALDataFreeUsage")
		}
		return nil
	}
	return m.SALDataContract.(*_SALData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SALDataFreeUsage) IsSALDataFreeUsage() {}

func (m *_SALDataFreeUsage) DeepCopy() any {
	return m.deepCopy()
}

func (m *_SALDataFreeUsage) deepCopy() *_SALDataFreeUsage {
	if m == nil {
		return nil
	}
	_SALDataFreeUsageCopy := &_SALDataFreeUsage{
		m.SALDataContract.(*_SALData).deepCopy(),
	}
	m.SALDataContract.(*_SALData)._SubType = m
	return _SALDataFreeUsageCopy
}

func (m *_SALDataFreeUsage) String() string {
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
