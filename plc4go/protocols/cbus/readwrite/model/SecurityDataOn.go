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

// SecurityDataOn is the corresponding interface of SecurityDataOn
type SecurityDataOn interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	SecurityData
	// GetData returns Data (property field)
	GetData() []byte
	// IsSecurityDataOn is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsSecurityDataOn()
	// CreateBuilder creates a SecurityDataOnBuilder
	CreateSecurityDataOnBuilder() SecurityDataOnBuilder
}

// _SecurityDataOn is the data-structure of this message
type _SecurityDataOn struct {
	SecurityDataContract
	Data []byte
}

var _ SecurityDataOn = (*_SecurityDataOn)(nil)
var _ SecurityDataRequirements = (*_SecurityDataOn)(nil)

// NewSecurityDataOn factory function for _SecurityDataOn
func NewSecurityDataOn(commandTypeContainer SecurityCommandTypeContainer, argument byte, data []byte) *_SecurityDataOn {
	_result := &_SecurityDataOn{
		SecurityDataContract: NewSecurityData(commandTypeContainer, argument),
		Data:                 data,
	}
	_result.SecurityDataContract.(*_SecurityData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// SecurityDataOnBuilder is a builder for SecurityDataOn
type SecurityDataOnBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(data []byte) SecurityDataOnBuilder
	// WithData adds Data (property field)
	WithData(...byte) SecurityDataOnBuilder
	// Build builds the SecurityDataOn or returns an error if something is wrong
	Build() (SecurityDataOn, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() SecurityDataOn
}

// NewSecurityDataOnBuilder() creates a SecurityDataOnBuilder
func NewSecurityDataOnBuilder() SecurityDataOnBuilder {
	return &_SecurityDataOnBuilder{_SecurityDataOn: new(_SecurityDataOn)}
}

type _SecurityDataOnBuilder struct {
	*_SecurityDataOn

	parentBuilder *_SecurityDataBuilder

	err *utils.MultiError
}

var _ (SecurityDataOnBuilder) = (*_SecurityDataOnBuilder)(nil)

func (b *_SecurityDataOnBuilder) setParent(contract SecurityDataContract) {
	b.SecurityDataContract = contract
}

func (b *_SecurityDataOnBuilder) WithMandatoryFields(data []byte) SecurityDataOnBuilder {
	return b.WithData(data...)
}

func (b *_SecurityDataOnBuilder) WithData(data ...byte) SecurityDataOnBuilder {
	b.Data = data
	return b
}

func (b *_SecurityDataOnBuilder) Build() (SecurityDataOn, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._SecurityDataOn.deepCopy(), nil
}

func (b *_SecurityDataOnBuilder) MustBuild() SecurityDataOn {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_SecurityDataOnBuilder) Done() SecurityDataBuilder {
	return b.parentBuilder
}

func (b *_SecurityDataOnBuilder) buildForSecurityData() (SecurityData, error) {
	return b.Build()
}

func (b *_SecurityDataOnBuilder) DeepCopy() any {
	_copy := b.CreateSecurityDataOnBuilder().(*_SecurityDataOnBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateSecurityDataOnBuilder creates a SecurityDataOnBuilder
func (b *_SecurityDataOn) CreateSecurityDataOnBuilder() SecurityDataOnBuilder {
	if b == nil {
		return NewSecurityDataOnBuilder()
	}
	return &_SecurityDataOnBuilder{_SecurityDataOn: b.deepCopy()}
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

func (m *_SecurityDataOn) GetParent() SecurityDataContract {
	return m.SecurityDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_SecurityDataOn) GetData() []byte {
	return m.Data
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastSecurityDataOn(structType any) SecurityDataOn {
	if casted, ok := structType.(SecurityDataOn); ok {
		return casted
	}
	if casted, ok := structType.(*SecurityDataOn); ok {
		return *casted
	}
	return nil
}

func (m *_SecurityDataOn) GetTypeName() string {
	return "SecurityDataOn"
}

func (m *_SecurityDataOn) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.SecurityDataContract.(*_SecurityData).GetLengthInBits(ctx))

	// Array field
	if len(m.Data) > 0 {
		lengthInBits += 8 * uint16(len(m.Data))
	}

	return lengthInBits
}

func (m *_SecurityDataOn) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_SecurityDataOn) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_SecurityData, commandTypeContainer SecurityCommandTypeContainer) (__securityDataOn SecurityDataOn, err error) {
	m.SecurityDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SecurityDataOn"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SecurityDataOn")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	data, err := readBuffer.ReadByteArray("data", int(int32(commandTypeContainer.NumBytes())-int32(int32(1))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'data' field"))
	}
	m.Data = data

	if closeErr := readBuffer.CloseContext("SecurityDataOn"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SecurityDataOn")
	}

	return m, nil
}

func (m *_SecurityDataOn) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_SecurityDataOn) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SecurityDataOn"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SecurityDataOn")
		}

		if err := WriteByteArrayField(ctx, "data", m.GetData(), WriteByteArray(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'data' field")
		}

		if popErr := writeBuffer.PopContext("SecurityDataOn"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SecurityDataOn")
		}
		return nil
	}
	return m.SecurityDataContract.(*_SecurityData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SecurityDataOn) IsSecurityDataOn() {}

func (m *_SecurityDataOn) DeepCopy() any {
	return m.deepCopy()
}

func (m *_SecurityDataOn) deepCopy() *_SecurityDataOn {
	if m == nil {
		return nil
	}
	_SecurityDataOnCopy := &_SecurityDataOn{
		m.SecurityDataContract.(*_SecurityData).deepCopy(),
		utils.DeepCopySlice[byte, byte](m.Data),
	}
	m.SecurityDataContract.(*_SecurityData)._SubType = m
	return _SecurityDataOnCopy
}

func (m *_SecurityDataOn) String() string {
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
