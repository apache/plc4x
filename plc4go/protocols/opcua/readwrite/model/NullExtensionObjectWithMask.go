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
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// NullExtensionObjectWithMask is the corresponding interface of NullExtensionObjectWithMask
type NullExtensionObjectWithMask interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectWithMask
	// GetBody returns Body (virtual field)
	GetBody() ExtensionObjectDefinition
	// IsNullExtensionObjectWithMask is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsNullExtensionObjectWithMask()
	// CreateBuilder creates a NullExtensionObjectWithMaskBuilder
	CreateNullExtensionObjectWithMaskBuilder() NullExtensionObjectWithMaskBuilder
}

// _NullExtensionObjectWithMask is the data-structure of this message
type _NullExtensionObjectWithMask struct {
	ExtensionObjectWithMaskContract
}

var _ NullExtensionObjectWithMask = (*_NullExtensionObjectWithMask)(nil)
var _ ExtensionObjectWithMaskRequirements = (*_NullExtensionObjectWithMask)(nil)

// NewNullExtensionObjectWithMask factory function for _NullExtensionObjectWithMask
func NewNullExtensionObjectWithMask(typeId ExpandedNodeId, encodingMask ExtensionObjectEncodingMask, extensionId int32, includeEncodingMask bool) *_NullExtensionObjectWithMask {
	_result := &_NullExtensionObjectWithMask{
		ExtensionObjectWithMaskContract: NewExtensionObjectWithMask(typeId, encodingMask, extensionId),
	}
	_result.ExtensionObjectWithMaskContract.(*_ExtensionObjectWithMask)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// NullExtensionObjectWithMaskBuilder is a builder for NullExtensionObjectWithMask
type NullExtensionObjectWithMaskBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() NullExtensionObjectWithMaskBuilder
	// Build builds the NullExtensionObjectWithMask or returns an error if something is wrong
	Build() (NullExtensionObjectWithMask, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() NullExtensionObjectWithMask
}

// NewNullExtensionObjectWithMaskBuilder() creates a NullExtensionObjectWithMaskBuilder
func NewNullExtensionObjectWithMaskBuilder() NullExtensionObjectWithMaskBuilder {
	return &_NullExtensionObjectWithMaskBuilder{_NullExtensionObjectWithMask: new(_NullExtensionObjectWithMask)}
}

type _NullExtensionObjectWithMaskBuilder struct {
	*_NullExtensionObjectWithMask

	parentBuilder *_ExtensionObjectWithMaskBuilder

	err *utils.MultiError
}

var _ (NullExtensionObjectWithMaskBuilder) = (*_NullExtensionObjectWithMaskBuilder)(nil)

func (b *_NullExtensionObjectWithMaskBuilder) setParent(contract ExtensionObjectWithMaskContract) {
	b.ExtensionObjectWithMaskContract = contract
}

func (b *_NullExtensionObjectWithMaskBuilder) WithMandatoryFields() NullExtensionObjectWithMaskBuilder {
	return b
}

func (b *_NullExtensionObjectWithMaskBuilder) Build() (NullExtensionObjectWithMask, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._NullExtensionObjectWithMask.deepCopy(), nil
}

func (b *_NullExtensionObjectWithMaskBuilder) MustBuild() NullExtensionObjectWithMask {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_NullExtensionObjectWithMaskBuilder) Done() ExtensionObjectWithMaskBuilder {
	return b.parentBuilder
}

func (b *_NullExtensionObjectWithMaskBuilder) buildForExtensionObjectWithMask() (ExtensionObjectWithMask, error) {
	return b.Build()
}

func (b *_NullExtensionObjectWithMaskBuilder) DeepCopy() any {
	_copy := b.CreateNullExtensionObjectWithMaskBuilder().(*_NullExtensionObjectWithMaskBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateNullExtensionObjectWithMaskBuilder creates a NullExtensionObjectWithMaskBuilder
func (b *_NullExtensionObjectWithMask) CreateNullExtensionObjectWithMaskBuilder() NullExtensionObjectWithMaskBuilder {
	if b == nil {
		return NewNullExtensionObjectWithMaskBuilder()
	}
	return &_NullExtensionObjectWithMaskBuilder{_NullExtensionObjectWithMask: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_NullExtensionObjectWithMask) GetEncodingMaskXmlBody() bool {
	return bool(false)
}

func (m *_NullExtensionObjectWithMask) GetEncodingMaskBinaryBody() bool {
	return bool(false)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_NullExtensionObjectWithMask) GetParent() ExtensionObjectWithMaskContract {
	return m.ExtensionObjectWithMaskContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_NullExtensionObjectWithMask) GetBody() ExtensionObjectDefinition {
	ctx := context.Background()
	_ = ctx
	return CastExtensionObjectDefinition(nil)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastNullExtensionObjectWithMask(structType any) NullExtensionObjectWithMask {
	if casted, ok := structType.(NullExtensionObjectWithMask); ok {
		return casted
	}
	if casted, ok := structType.(*NullExtensionObjectWithMask); ok {
		return *casted
	}
	return nil
}

func (m *_NullExtensionObjectWithMask) GetTypeName() string {
	return "NullExtensionObjectWithMask"
}

func (m *_NullExtensionObjectWithMask) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectWithMaskContract.(*_ExtensionObjectWithMask).GetLengthInBits(ctx))

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_NullExtensionObjectWithMask) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_NullExtensionObjectWithMask) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectWithMask, extensionId int32, includeEncodingMask bool) (__nullExtensionObjectWithMask NullExtensionObjectWithMask, err error) {
	m.ExtensionObjectWithMaskContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("NullExtensionObjectWithMask"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for NullExtensionObjectWithMask")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	body, err := ReadVirtualField[ExtensionObjectDefinition](ctx, "body", (*ExtensionObjectDefinition)(nil), nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'body' field"))
	}
	_ = body

	if closeErr := readBuffer.CloseContext("NullExtensionObjectWithMask"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for NullExtensionObjectWithMask")
	}

	return m, nil
}

func (m *_NullExtensionObjectWithMask) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_NullExtensionObjectWithMask) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("NullExtensionObjectWithMask"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for NullExtensionObjectWithMask")
		}
		// Virtual field
		body := m.GetBody()
		_ = body
		if _bodyErr := writeBuffer.WriteVirtual(ctx, "body", m.GetBody()); _bodyErr != nil {
			return errors.Wrap(_bodyErr, "Error serializing 'body' field")
		}

		if popErr := writeBuffer.PopContext("NullExtensionObjectWithMask"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for NullExtensionObjectWithMask")
		}
		return nil
	}
	return m.ExtensionObjectWithMaskContract.(*_ExtensionObjectWithMask).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_NullExtensionObjectWithMask) IsNullExtensionObjectWithMask() {}

func (m *_NullExtensionObjectWithMask) DeepCopy() any {
	return m.deepCopy()
}

func (m *_NullExtensionObjectWithMask) deepCopy() *_NullExtensionObjectWithMask {
	if m == nil {
		return nil
	}
	_NullExtensionObjectWithMaskCopy := &_NullExtensionObjectWithMask{
		m.ExtensionObjectWithMaskContract.(*_ExtensionObjectWithMask).deepCopy(),
	}
	m.ExtensionObjectWithMaskContract.(*_ExtensionObjectWithMask)._SubType = m
	return _NullExtensionObjectWithMaskCopy
}

func (m *_NullExtensionObjectWithMask) String() string {
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
