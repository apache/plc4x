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

// MeteringDataMeasureOil is the corresponding interface of MeteringDataMeasureOil
type MeteringDataMeasureOil interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	MeteringData
	// IsMeteringDataMeasureOil is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsMeteringDataMeasureOil()
	// CreateBuilder creates a MeteringDataMeasureOilBuilder
	CreateMeteringDataMeasureOilBuilder() MeteringDataMeasureOilBuilder
}

// _MeteringDataMeasureOil is the data-structure of this message
type _MeteringDataMeasureOil struct {
	MeteringDataContract
}

var _ MeteringDataMeasureOil = (*_MeteringDataMeasureOil)(nil)
var _ MeteringDataRequirements = (*_MeteringDataMeasureOil)(nil)

// NewMeteringDataMeasureOil factory function for _MeteringDataMeasureOil
func NewMeteringDataMeasureOil(commandTypeContainer MeteringCommandTypeContainer, argument byte) *_MeteringDataMeasureOil {
	_result := &_MeteringDataMeasureOil{
		MeteringDataContract: NewMeteringData(commandTypeContainer, argument),
	}
	_result.MeteringDataContract.(*_MeteringData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// MeteringDataMeasureOilBuilder is a builder for MeteringDataMeasureOil
type MeteringDataMeasureOilBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() MeteringDataMeasureOilBuilder
	// Build builds the MeteringDataMeasureOil or returns an error if something is wrong
	Build() (MeteringDataMeasureOil, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() MeteringDataMeasureOil
}

// NewMeteringDataMeasureOilBuilder() creates a MeteringDataMeasureOilBuilder
func NewMeteringDataMeasureOilBuilder() MeteringDataMeasureOilBuilder {
	return &_MeteringDataMeasureOilBuilder{_MeteringDataMeasureOil: new(_MeteringDataMeasureOil)}
}

type _MeteringDataMeasureOilBuilder struct {
	*_MeteringDataMeasureOil

	parentBuilder *_MeteringDataBuilder

	err *utils.MultiError
}

var _ (MeteringDataMeasureOilBuilder) = (*_MeteringDataMeasureOilBuilder)(nil)

func (b *_MeteringDataMeasureOilBuilder) setParent(contract MeteringDataContract) {
	b.MeteringDataContract = contract
}

func (b *_MeteringDataMeasureOilBuilder) WithMandatoryFields() MeteringDataMeasureOilBuilder {
	return b
}

func (b *_MeteringDataMeasureOilBuilder) Build() (MeteringDataMeasureOil, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._MeteringDataMeasureOil.deepCopy(), nil
}

func (b *_MeteringDataMeasureOilBuilder) MustBuild() MeteringDataMeasureOil {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_MeteringDataMeasureOilBuilder) Done() MeteringDataBuilder {
	return b.parentBuilder
}

func (b *_MeteringDataMeasureOilBuilder) buildForMeteringData() (MeteringData, error) {
	return b.Build()
}

func (b *_MeteringDataMeasureOilBuilder) DeepCopy() any {
	_copy := b.CreateMeteringDataMeasureOilBuilder().(*_MeteringDataMeasureOilBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateMeteringDataMeasureOilBuilder creates a MeteringDataMeasureOilBuilder
func (b *_MeteringDataMeasureOil) CreateMeteringDataMeasureOilBuilder() MeteringDataMeasureOilBuilder {
	if b == nil {
		return NewMeteringDataMeasureOilBuilder()
	}
	return &_MeteringDataMeasureOilBuilder{_MeteringDataMeasureOil: b.deepCopy()}
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

func (m *_MeteringDataMeasureOil) GetParent() MeteringDataContract {
	return m.MeteringDataContract
}

// Deprecated: use the interface for direct cast
func CastMeteringDataMeasureOil(structType any) MeteringDataMeasureOil {
	if casted, ok := structType.(MeteringDataMeasureOil); ok {
		return casted
	}
	if casted, ok := structType.(*MeteringDataMeasureOil); ok {
		return *casted
	}
	return nil
}

func (m *_MeteringDataMeasureOil) GetTypeName() string {
	return "MeteringDataMeasureOil"
}

func (m *_MeteringDataMeasureOil) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.MeteringDataContract.(*_MeteringData).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_MeteringDataMeasureOil) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_MeteringDataMeasureOil) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_MeteringData) (__meteringDataMeasureOil MeteringDataMeasureOil, err error) {
	m.MeteringDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("MeteringDataMeasureOil"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for MeteringDataMeasureOil")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("MeteringDataMeasureOil"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for MeteringDataMeasureOil")
	}

	return m, nil
}

func (m *_MeteringDataMeasureOil) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_MeteringDataMeasureOil) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("MeteringDataMeasureOil"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for MeteringDataMeasureOil")
		}

		if popErr := writeBuffer.PopContext("MeteringDataMeasureOil"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for MeteringDataMeasureOil")
		}
		return nil
	}
	return m.MeteringDataContract.(*_MeteringData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_MeteringDataMeasureOil) IsMeteringDataMeasureOil() {}

func (m *_MeteringDataMeasureOil) DeepCopy() any {
	return m.deepCopy()
}

func (m *_MeteringDataMeasureOil) deepCopy() *_MeteringDataMeasureOil {
	if m == nil {
		return nil
	}
	_MeteringDataMeasureOilCopy := &_MeteringDataMeasureOil{
		m.MeteringDataContract.(*_MeteringData).deepCopy(),
	}
	m.MeteringDataContract.(*_MeteringData)._SubType = m
	return _MeteringDataMeasureOilCopy
}

func (m *_MeteringDataMeasureOil) String() string {
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
