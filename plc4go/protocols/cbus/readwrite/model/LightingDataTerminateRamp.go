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

// LightingDataTerminateRamp is the corresponding interface of LightingDataTerminateRamp
type LightingDataTerminateRamp interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	LightingData
	// GetGroup returns Group (property field)
	GetGroup() byte
	// IsLightingDataTerminateRamp is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsLightingDataTerminateRamp()
	// CreateBuilder creates a LightingDataTerminateRampBuilder
	CreateLightingDataTerminateRampBuilder() LightingDataTerminateRampBuilder
}

// _LightingDataTerminateRamp is the data-structure of this message
type _LightingDataTerminateRamp struct {
	LightingDataContract
	Group byte
}

var _ LightingDataTerminateRamp = (*_LightingDataTerminateRamp)(nil)
var _ LightingDataRequirements = (*_LightingDataTerminateRamp)(nil)

// NewLightingDataTerminateRamp factory function for _LightingDataTerminateRamp
func NewLightingDataTerminateRamp(commandTypeContainer LightingCommandTypeContainer, group byte) *_LightingDataTerminateRamp {
	_result := &_LightingDataTerminateRamp{
		LightingDataContract: NewLightingData(commandTypeContainer),
		Group:                group,
	}
	_result.LightingDataContract.(*_LightingData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// LightingDataTerminateRampBuilder is a builder for LightingDataTerminateRamp
type LightingDataTerminateRampBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(group byte) LightingDataTerminateRampBuilder
	// WithGroup adds Group (property field)
	WithGroup(byte) LightingDataTerminateRampBuilder
	// Build builds the LightingDataTerminateRamp or returns an error if something is wrong
	Build() (LightingDataTerminateRamp, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() LightingDataTerminateRamp
}

// NewLightingDataTerminateRampBuilder() creates a LightingDataTerminateRampBuilder
func NewLightingDataTerminateRampBuilder() LightingDataTerminateRampBuilder {
	return &_LightingDataTerminateRampBuilder{_LightingDataTerminateRamp: new(_LightingDataTerminateRamp)}
}

type _LightingDataTerminateRampBuilder struct {
	*_LightingDataTerminateRamp

	parentBuilder *_LightingDataBuilder

	err *utils.MultiError
}

var _ (LightingDataTerminateRampBuilder) = (*_LightingDataTerminateRampBuilder)(nil)

func (b *_LightingDataTerminateRampBuilder) setParent(contract LightingDataContract) {
	b.LightingDataContract = contract
}

func (b *_LightingDataTerminateRampBuilder) WithMandatoryFields(group byte) LightingDataTerminateRampBuilder {
	return b.WithGroup(group)
}

func (b *_LightingDataTerminateRampBuilder) WithGroup(group byte) LightingDataTerminateRampBuilder {
	b.Group = group
	return b
}

func (b *_LightingDataTerminateRampBuilder) Build() (LightingDataTerminateRamp, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._LightingDataTerminateRamp.deepCopy(), nil
}

func (b *_LightingDataTerminateRampBuilder) MustBuild() LightingDataTerminateRamp {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_LightingDataTerminateRampBuilder) Done() LightingDataBuilder {
	return b.parentBuilder
}

func (b *_LightingDataTerminateRampBuilder) buildForLightingData() (LightingData, error) {
	return b.Build()
}

func (b *_LightingDataTerminateRampBuilder) DeepCopy() any {
	_copy := b.CreateLightingDataTerminateRampBuilder().(*_LightingDataTerminateRampBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateLightingDataTerminateRampBuilder creates a LightingDataTerminateRampBuilder
func (b *_LightingDataTerminateRamp) CreateLightingDataTerminateRampBuilder() LightingDataTerminateRampBuilder {
	if b == nil {
		return NewLightingDataTerminateRampBuilder()
	}
	return &_LightingDataTerminateRampBuilder{_LightingDataTerminateRamp: b.deepCopy()}
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

func (m *_LightingDataTerminateRamp) GetParent() LightingDataContract {
	return m.LightingDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_LightingDataTerminateRamp) GetGroup() byte {
	return m.Group
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastLightingDataTerminateRamp(structType any) LightingDataTerminateRamp {
	if casted, ok := structType.(LightingDataTerminateRamp); ok {
		return casted
	}
	if casted, ok := structType.(*LightingDataTerminateRamp); ok {
		return *casted
	}
	return nil
}

func (m *_LightingDataTerminateRamp) GetTypeName() string {
	return "LightingDataTerminateRamp"
}

func (m *_LightingDataTerminateRamp) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.LightingDataContract.(*_LightingData).GetLengthInBits(ctx))

	// Simple field (group)
	lengthInBits += 8

	return lengthInBits
}

func (m *_LightingDataTerminateRamp) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_LightingDataTerminateRamp) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_LightingData) (__lightingDataTerminateRamp LightingDataTerminateRamp, err error) {
	m.LightingDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("LightingDataTerminateRamp"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for LightingDataTerminateRamp")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	group, err := ReadSimpleField(ctx, "group", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'group' field"))
	}
	m.Group = group

	if closeErr := readBuffer.CloseContext("LightingDataTerminateRamp"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for LightingDataTerminateRamp")
	}

	return m, nil
}

func (m *_LightingDataTerminateRamp) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_LightingDataTerminateRamp) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("LightingDataTerminateRamp"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for LightingDataTerminateRamp")
		}

		if err := WriteSimpleField[byte](ctx, "group", m.GetGroup(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'group' field")
		}

		if popErr := writeBuffer.PopContext("LightingDataTerminateRamp"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for LightingDataTerminateRamp")
		}
		return nil
	}
	return m.LightingDataContract.(*_LightingData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_LightingDataTerminateRamp) IsLightingDataTerminateRamp() {}

func (m *_LightingDataTerminateRamp) DeepCopy() any {
	return m.deepCopy()
}

func (m *_LightingDataTerminateRamp) deepCopy() *_LightingDataTerminateRamp {
	if m == nil {
		return nil
	}
	_LightingDataTerminateRampCopy := &_LightingDataTerminateRamp{
		m.LightingDataContract.(*_LightingData).deepCopy(),
		m.Group,
	}
	m.LightingDataContract.(*_LightingData)._SubType = m
	return _LightingDataTerminateRampCopy
}

func (m *_LightingDataTerminateRamp) String() string {
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
