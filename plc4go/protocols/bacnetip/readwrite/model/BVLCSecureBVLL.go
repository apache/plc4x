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
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// BVLCSecureBVLL is the corresponding interface of BVLCSecureBVLL
type BVLCSecureBVLL interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BVLC
	// GetSecurityWrapper returns SecurityWrapper (property field)
	GetSecurityWrapper() []byte
	// IsBVLCSecureBVLL is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBVLCSecureBVLL()
	// CreateBuilder creates a BVLCSecureBVLLBuilder
	CreateBVLCSecureBVLLBuilder() BVLCSecureBVLLBuilder
}

// _BVLCSecureBVLL is the data-structure of this message
type _BVLCSecureBVLL struct {
	BVLCContract
	SecurityWrapper []byte

	// Arguments.
	BvlcPayloadLength uint16
}

var _ BVLCSecureBVLL = (*_BVLCSecureBVLL)(nil)
var _ BVLCRequirements = (*_BVLCSecureBVLL)(nil)

// NewBVLCSecureBVLL factory function for _BVLCSecureBVLL
func NewBVLCSecureBVLL(securityWrapper []byte, bvlcPayloadLength uint16) *_BVLCSecureBVLL {
	_result := &_BVLCSecureBVLL{
		BVLCContract:    NewBVLC(),
		SecurityWrapper: securityWrapper,
	}
	_result.BVLCContract.(*_BVLC)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BVLCSecureBVLLBuilder is a builder for BVLCSecureBVLL
type BVLCSecureBVLLBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(securityWrapper []byte) BVLCSecureBVLLBuilder
	// WithSecurityWrapper adds SecurityWrapper (property field)
	WithSecurityWrapper(...byte) BVLCSecureBVLLBuilder
	// Build builds the BVLCSecureBVLL or returns an error if something is wrong
	Build() (BVLCSecureBVLL, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BVLCSecureBVLL
}

// NewBVLCSecureBVLLBuilder() creates a BVLCSecureBVLLBuilder
func NewBVLCSecureBVLLBuilder() BVLCSecureBVLLBuilder {
	return &_BVLCSecureBVLLBuilder{_BVLCSecureBVLL: new(_BVLCSecureBVLL)}
}

type _BVLCSecureBVLLBuilder struct {
	*_BVLCSecureBVLL

	parentBuilder *_BVLCBuilder

	err *utils.MultiError
}

var _ (BVLCSecureBVLLBuilder) = (*_BVLCSecureBVLLBuilder)(nil)

func (b *_BVLCSecureBVLLBuilder) setParent(contract BVLCContract) {
	b.BVLCContract = contract
}

func (b *_BVLCSecureBVLLBuilder) WithMandatoryFields(securityWrapper []byte) BVLCSecureBVLLBuilder {
	return b.WithSecurityWrapper(securityWrapper...)
}

func (b *_BVLCSecureBVLLBuilder) WithSecurityWrapper(securityWrapper ...byte) BVLCSecureBVLLBuilder {
	b.SecurityWrapper = securityWrapper
	return b
}

func (b *_BVLCSecureBVLLBuilder) Build() (BVLCSecureBVLL, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BVLCSecureBVLL.deepCopy(), nil
}

func (b *_BVLCSecureBVLLBuilder) MustBuild() BVLCSecureBVLL {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BVLCSecureBVLLBuilder) Done() BVLCBuilder {
	return b.parentBuilder
}

func (b *_BVLCSecureBVLLBuilder) buildForBVLC() (BVLC, error) {
	return b.Build()
}

func (b *_BVLCSecureBVLLBuilder) DeepCopy() any {
	_copy := b.CreateBVLCSecureBVLLBuilder().(*_BVLCSecureBVLLBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBVLCSecureBVLLBuilder creates a BVLCSecureBVLLBuilder
func (b *_BVLCSecureBVLL) CreateBVLCSecureBVLLBuilder() BVLCSecureBVLLBuilder {
	if b == nil {
		return NewBVLCSecureBVLLBuilder()
	}
	return &_BVLCSecureBVLLBuilder{_BVLCSecureBVLL: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BVLCSecureBVLL) GetBvlcFunction() uint8 {
	return 0x0C
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BVLCSecureBVLL) GetParent() BVLCContract {
	return m.BVLCContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BVLCSecureBVLL) GetSecurityWrapper() []byte {
	return m.SecurityWrapper
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBVLCSecureBVLL(structType any) BVLCSecureBVLL {
	if casted, ok := structType.(BVLCSecureBVLL); ok {
		return casted
	}
	if casted, ok := structType.(*BVLCSecureBVLL); ok {
		return *casted
	}
	return nil
}

func (m *_BVLCSecureBVLL) GetTypeName() string {
	return "BVLCSecureBVLL"
}

func (m *_BVLCSecureBVLL) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BVLCContract.(*_BVLC).GetLengthInBits(ctx))

	// Array field
	if len(m.SecurityWrapper) > 0 {
		lengthInBits += 8 * uint16(len(m.SecurityWrapper))
	}

	return lengthInBits
}

func (m *_BVLCSecureBVLL) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BVLCSecureBVLL) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BVLC, bvlcPayloadLength uint16) (__bVLCSecureBVLL BVLCSecureBVLL, err error) {
	m.BVLCContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BVLCSecureBVLL"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BVLCSecureBVLL")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	securityWrapper, err := readBuffer.ReadByteArray("securityWrapper", int(bvlcPayloadLength), codegen.WithByteOrder(binary.BigEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'securityWrapper' field"))
	}
	m.SecurityWrapper = securityWrapper

	if closeErr := readBuffer.CloseContext("BVLCSecureBVLL"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BVLCSecureBVLL")
	}

	return m, nil
}

func (m *_BVLCSecureBVLL) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))), utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BVLCSecureBVLL) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BVLCSecureBVLL"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BVLCSecureBVLL")
		}

		if err := WriteByteArrayField(ctx, "securityWrapper", m.GetSecurityWrapper(), WriteByteArray(writeBuffer, 8), codegen.WithByteOrder(binary.BigEndian)); err != nil {
			return errors.Wrap(err, "Error serializing 'securityWrapper' field")
		}

		if popErr := writeBuffer.PopContext("BVLCSecureBVLL"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BVLCSecureBVLL")
		}
		return nil
	}
	return m.BVLCContract.(*_BVLC).serializeParent(ctx, writeBuffer, m, ser)
}

////
// Arguments Getter

func (m *_BVLCSecureBVLL) GetBvlcPayloadLength() uint16 {
	return m.BvlcPayloadLength
}

//
////

func (m *_BVLCSecureBVLL) IsBVLCSecureBVLL() {}

func (m *_BVLCSecureBVLL) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BVLCSecureBVLL) deepCopy() *_BVLCSecureBVLL {
	if m == nil {
		return nil
	}
	_BVLCSecureBVLLCopy := &_BVLCSecureBVLL{
		m.BVLCContract.(*_BVLC).deepCopy(),
		utils.DeepCopySlice[byte, byte](m.SecurityWrapper),
		m.BvlcPayloadLength,
	}
	m.BVLCContract.(*_BVLC)._SubType = m
	return _BVLCSecureBVLLCopy
}

func (m *_BVLCSecureBVLL) String() string {
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
