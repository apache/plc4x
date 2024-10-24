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

// CALDataReset is the corresponding interface of CALDataReset
type CALDataReset interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	CALData
	// IsCALDataReset is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCALDataReset()
	// CreateBuilder creates a CALDataResetBuilder
	CreateCALDataResetBuilder() CALDataResetBuilder
}

// _CALDataReset is the data-structure of this message
type _CALDataReset struct {
	CALDataContract
}

var _ CALDataReset = (*_CALDataReset)(nil)
var _ CALDataRequirements = (*_CALDataReset)(nil)

// NewCALDataReset factory function for _CALDataReset
func NewCALDataReset(commandTypeContainer CALCommandTypeContainer, additionalData CALData, requestContext RequestContext) *_CALDataReset {
	_result := &_CALDataReset{
		CALDataContract: NewCALData(commandTypeContainer, additionalData, requestContext),
	}
	_result.CALDataContract.(*_CALData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// CALDataResetBuilder is a builder for CALDataReset
type CALDataResetBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() CALDataResetBuilder
	// Build builds the CALDataReset or returns an error if something is wrong
	Build() (CALDataReset, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() CALDataReset
}

// NewCALDataResetBuilder() creates a CALDataResetBuilder
func NewCALDataResetBuilder() CALDataResetBuilder {
	return &_CALDataResetBuilder{_CALDataReset: new(_CALDataReset)}
}

type _CALDataResetBuilder struct {
	*_CALDataReset

	parentBuilder *_CALDataBuilder

	err *utils.MultiError
}

var _ (CALDataResetBuilder) = (*_CALDataResetBuilder)(nil)

func (b *_CALDataResetBuilder) setParent(contract CALDataContract) {
	b.CALDataContract = contract
}

func (b *_CALDataResetBuilder) WithMandatoryFields() CALDataResetBuilder {
	return b
}

func (b *_CALDataResetBuilder) Build() (CALDataReset, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._CALDataReset.deepCopy(), nil
}

func (b *_CALDataResetBuilder) MustBuild() CALDataReset {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_CALDataResetBuilder) Done() CALDataBuilder {
	return b.parentBuilder
}

func (b *_CALDataResetBuilder) buildForCALData() (CALData, error) {
	return b.Build()
}

func (b *_CALDataResetBuilder) DeepCopy() any {
	_copy := b.CreateCALDataResetBuilder().(*_CALDataResetBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateCALDataResetBuilder creates a CALDataResetBuilder
func (b *_CALDataReset) CreateCALDataResetBuilder() CALDataResetBuilder {
	if b == nil {
		return NewCALDataResetBuilder()
	}
	return &_CALDataResetBuilder{_CALDataReset: b.deepCopy()}
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

func (m *_CALDataReset) GetParent() CALDataContract {
	return m.CALDataContract
}

// Deprecated: use the interface for direct cast
func CastCALDataReset(structType any) CALDataReset {
	if casted, ok := structType.(CALDataReset); ok {
		return casted
	}
	if casted, ok := structType.(*CALDataReset); ok {
		return *casted
	}
	return nil
}

func (m *_CALDataReset) GetTypeName() string {
	return "CALDataReset"
}

func (m *_CALDataReset) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.CALDataContract.(*_CALData).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_CALDataReset) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_CALDataReset) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_CALData, requestContext RequestContext) (__cALDataReset CALDataReset, err error) {
	m.CALDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CALDataReset"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CALDataReset")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("CALDataReset"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CALDataReset")
	}

	return m, nil
}

func (m *_CALDataReset) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CALDataReset) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CALDataReset"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CALDataReset")
		}

		if popErr := writeBuffer.PopContext("CALDataReset"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CALDataReset")
		}
		return nil
	}
	return m.CALDataContract.(*_CALData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CALDataReset) IsCALDataReset() {}

func (m *_CALDataReset) DeepCopy() any {
	return m.deepCopy()
}

func (m *_CALDataReset) deepCopy() *_CALDataReset {
	if m == nil {
		return nil
	}
	_CALDataResetCopy := &_CALDataReset{
		m.CALDataContract.(*_CALData).deepCopy(),
	}
	m.CALDataContract.(*_CALData)._SubType = m
	return _CALDataResetCopy
}

func (m *_CALDataReset) String() string {
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
