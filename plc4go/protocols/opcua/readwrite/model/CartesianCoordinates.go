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

// CartesianCoordinates is the corresponding interface of CartesianCoordinates
type CartesianCoordinates interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// IsCartesianCoordinates is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCartesianCoordinates()
	// CreateBuilder creates a CartesianCoordinatesBuilder
	CreateCartesianCoordinatesBuilder() CartesianCoordinatesBuilder
}

// _CartesianCoordinates is the data-structure of this message
type _CartesianCoordinates struct {
	ExtensionObjectDefinitionContract
}

var _ CartesianCoordinates = (*_CartesianCoordinates)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_CartesianCoordinates)(nil)

// NewCartesianCoordinates factory function for _CartesianCoordinates
func NewCartesianCoordinates() *_CartesianCoordinates {
	_result := &_CartesianCoordinates{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// CartesianCoordinatesBuilder is a builder for CartesianCoordinates
type CartesianCoordinatesBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() CartesianCoordinatesBuilder
	// Build builds the CartesianCoordinates or returns an error if something is wrong
	Build() (CartesianCoordinates, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() CartesianCoordinates
}

// NewCartesianCoordinatesBuilder() creates a CartesianCoordinatesBuilder
func NewCartesianCoordinatesBuilder() CartesianCoordinatesBuilder {
	return &_CartesianCoordinatesBuilder{_CartesianCoordinates: new(_CartesianCoordinates)}
}

type _CartesianCoordinatesBuilder struct {
	*_CartesianCoordinates

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (CartesianCoordinatesBuilder) = (*_CartesianCoordinatesBuilder)(nil)

func (b *_CartesianCoordinatesBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_CartesianCoordinatesBuilder) WithMandatoryFields() CartesianCoordinatesBuilder {
	return b
}

func (b *_CartesianCoordinatesBuilder) Build() (CartesianCoordinates, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._CartesianCoordinates.deepCopy(), nil
}

func (b *_CartesianCoordinatesBuilder) MustBuild() CartesianCoordinates {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_CartesianCoordinatesBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_CartesianCoordinatesBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_CartesianCoordinatesBuilder) DeepCopy() any {
	_copy := b.CreateCartesianCoordinatesBuilder().(*_CartesianCoordinatesBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateCartesianCoordinatesBuilder creates a CartesianCoordinatesBuilder
func (b *_CartesianCoordinates) CreateCartesianCoordinatesBuilder() CartesianCoordinatesBuilder {
	if b == nil {
		return NewCartesianCoordinatesBuilder()
	}
	return &_CartesianCoordinatesBuilder{_CartesianCoordinates: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_CartesianCoordinates) GetExtensionId() int32 {
	return int32(18811)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_CartesianCoordinates) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

// Deprecated: use the interface for direct cast
func CastCartesianCoordinates(structType any) CartesianCoordinates {
	if casted, ok := structType.(CartesianCoordinates); ok {
		return casted
	}
	if casted, ok := structType.(*CartesianCoordinates); ok {
		return *casted
	}
	return nil
}

func (m *_CartesianCoordinates) GetTypeName() string {
	return "CartesianCoordinates"
}

func (m *_CartesianCoordinates) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_CartesianCoordinates) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_CartesianCoordinates) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__cartesianCoordinates CartesianCoordinates, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CartesianCoordinates"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CartesianCoordinates")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("CartesianCoordinates"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CartesianCoordinates")
	}

	return m, nil
}

func (m *_CartesianCoordinates) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CartesianCoordinates) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CartesianCoordinates"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CartesianCoordinates")
		}

		if popErr := writeBuffer.PopContext("CartesianCoordinates"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CartesianCoordinates")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CartesianCoordinates) IsCartesianCoordinates() {}

func (m *_CartesianCoordinates) DeepCopy() any {
	return m.deepCopy()
}

func (m *_CartesianCoordinates) deepCopy() *_CartesianCoordinates {
	if m == nil {
		return nil
	}
	_CartesianCoordinatesCopy := &_CartesianCoordinates{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _CartesianCoordinatesCopy
}

func (m *_CartesianCoordinates) String() string {
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
