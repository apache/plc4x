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

// COTPParameterDisconnectAdditionalInformation is the corresponding interface of COTPParameterDisconnectAdditionalInformation
type COTPParameterDisconnectAdditionalInformation interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	COTPParameter
	// GetData returns Data (property field)
	GetData() []byte
	// IsCOTPParameterDisconnectAdditionalInformation is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCOTPParameterDisconnectAdditionalInformation()
	// CreateBuilder creates a COTPParameterDisconnectAdditionalInformationBuilder
	CreateCOTPParameterDisconnectAdditionalInformationBuilder() COTPParameterDisconnectAdditionalInformationBuilder
}

// _COTPParameterDisconnectAdditionalInformation is the data-structure of this message
type _COTPParameterDisconnectAdditionalInformation struct {
	COTPParameterContract
	Data []byte
}

var _ COTPParameterDisconnectAdditionalInformation = (*_COTPParameterDisconnectAdditionalInformation)(nil)
var _ COTPParameterRequirements = (*_COTPParameterDisconnectAdditionalInformation)(nil)

// NewCOTPParameterDisconnectAdditionalInformation factory function for _COTPParameterDisconnectAdditionalInformation
func NewCOTPParameterDisconnectAdditionalInformation(data []byte, rest uint8) *_COTPParameterDisconnectAdditionalInformation {
	_result := &_COTPParameterDisconnectAdditionalInformation{
		COTPParameterContract: NewCOTPParameter(rest),
		Data:                  data,
	}
	_result.COTPParameterContract.(*_COTPParameter)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// COTPParameterDisconnectAdditionalInformationBuilder is a builder for COTPParameterDisconnectAdditionalInformation
type COTPParameterDisconnectAdditionalInformationBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(data []byte) COTPParameterDisconnectAdditionalInformationBuilder
	// WithData adds Data (property field)
	WithData(...byte) COTPParameterDisconnectAdditionalInformationBuilder
	// Build builds the COTPParameterDisconnectAdditionalInformation or returns an error if something is wrong
	Build() (COTPParameterDisconnectAdditionalInformation, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() COTPParameterDisconnectAdditionalInformation
}

// NewCOTPParameterDisconnectAdditionalInformationBuilder() creates a COTPParameterDisconnectAdditionalInformationBuilder
func NewCOTPParameterDisconnectAdditionalInformationBuilder() COTPParameterDisconnectAdditionalInformationBuilder {
	return &_COTPParameterDisconnectAdditionalInformationBuilder{_COTPParameterDisconnectAdditionalInformation: new(_COTPParameterDisconnectAdditionalInformation)}
}

type _COTPParameterDisconnectAdditionalInformationBuilder struct {
	*_COTPParameterDisconnectAdditionalInformation

	parentBuilder *_COTPParameterBuilder

	err *utils.MultiError
}

var _ (COTPParameterDisconnectAdditionalInformationBuilder) = (*_COTPParameterDisconnectAdditionalInformationBuilder)(nil)

func (b *_COTPParameterDisconnectAdditionalInformationBuilder) setParent(contract COTPParameterContract) {
	b.COTPParameterContract = contract
}

func (b *_COTPParameterDisconnectAdditionalInformationBuilder) WithMandatoryFields(data []byte) COTPParameterDisconnectAdditionalInformationBuilder {
	return b.WithData(data...)
}

func (b *_COTPParameterDisconnectAdditionalInformationBuilder) WithData(data ...byte) COTPParameterDisconnectAdditionalInformationBuilder {
	b.Data = data
	return b
}

func (b *_COTPParameterDisconnectAdditionalInformationBuilder) Build() (COTPParameterDisconnectAdditionalInformation, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._COTPParameterDisconnectAdditionalInformation.deepCopy(), nil
}

func (b *_COTPParameterDisconnectAdditionalInformationBuilder) MustBuild() COTPParameterDisconnectAdditionalInformation {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_COTPParameterDisconnectAdditionalInformationBuilder) Done() COTPParameterBuilder {
	return b.parentBuilder
}

func (b *_COTPParameterDisconnectAdditionalInformationBuilder) buildForCOTPParameter() (COTPParameter, error) {
	return b.Build()
}

func (b *_COTPParameterDisconnectAdditionalInformationBuilder) DeepCopy() any {
	_copy := b.CreateCOTPParameterDisconnectAdditionalInformationBuilder().(*_COTPParameterDisconnectAdditionalInformationBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateCOTPParameterDisconnectAdditionalInformationBuilder creates a COTPParameterDisconnectAdditionalInformationBuilder
func (b *_COTPParameterDisconnectAdditionalInformation) CreateCOTPParameterDisconnectAdditionalInformationBuilder() COTPParameterDisconnectAdditionalInformationBuilder {
	if b == nil {
		return NewCOTPParameterDisconnectAdditionalInformationBuilder()
	}
	return &_COTPParameterDisconnectAdditionalInformationBuilder{_COTPParameterDisconnectAdditionalInformation: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_COTPParameterDisconnectAdditionalInformation) GetParameterType() uint8 {
	return 0xE0
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_COTPParameterDisconnectAdditionalInformation) GetParent() COTPParameterContract {
	return m.COTPParameterContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_COTPParameterDisconnectAdditionalInformation) GetData() []byte {
	return m.Data
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastCOTPParameterDisconnectAdditionalInformation(structType any) COTPParameterDisconnectAdditionalInformation {
	if casted, ok := structType.(COTPParameterDisconnectAdditionalInformation); ok {
		return casted
	}
	if casted, ok := structType.(*COTPParameterDisconnectAdditionalInformation); ok {
		return *casted
	}
	return nil
}

func (m *_COTPParameterDisconnectAdditionalInformation) GetTypeName() string {
	return "COTPParameterDisconnectAdditionalInformation"
}

func (m *_COTPParameterDisconnectAdditionalInformation) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.COTPParameterContract.(*_COTPParameter).GetLengthInBits(ctx))

	// Array field
	if len(m.Data) > 0 {
		lengthInBits += 8 * uint16(len(m.Data))
	}

	return lengthInBits
}

func (m *_COTPParameterDisconnectAdditionalInformation) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_COTPParameterDisconnectAdditionalInformation) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_COTPParameter, rest uint8) (__cOTPParameterDisconnectAdditionalInformation COTPParameterDisconnectAdditionalInformation, err error) {
	m.COTPParameterContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("COTPParameterDisconnectAdditionalInformation"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for COTPParameterDisconnectAdditionalInformation")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	data, err := readBuffer.ReadByteArray("data", int(rest))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'data' field"))
	}
	m.Data = data

	if closeErr := readBuffer.CloseContext("COTPParameterDisconnectAdditionalInformation"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for COTPParameterDisconnectAdditionalInformation")
	}

	return m, nil
}

func (m *_COTPParameterDisconnectAdditionalInformation) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_COTPParameterDisconnectAdditionalInformation) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("COTPParameterDisconnectAdditionalInformation"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for COTPParameterDisconnectAdditionalInformation")
		}

		if err := WriteByteArrayField(ctx, "data", m.GetData(), WriteByteArray(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'data' field")
		}

		if popErr := writeBuffer.PopContext("COTPParameterDisconnectAdditionalInformation"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for COTPParameterDisconnectAdditionalInformation")
		}
		return nil
	}
	return m.COTPParameterContract.(*_COTPParameter).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_COTPParameterDisconnectAdditionalInformation) IsCOTPParameterDisconnectAdditionalInformation() {
}

func (m *_COTPParameterDisconnectAdditionalInformation) DeepCopy() any {
	return m.deepCopy()
}

func (m *_COTPParameterDisconnectAdditionalInformation) deepCopy() *_COTPParameterDisconnectAdditionalInformation {
	if m == nil {
		return nil
	}
	_COTPParameterDisconnectAdditionalInformationCopy := &_COTPParameterDisconnectAdditionalInformation{
		m.COTPParameterContract.(*_COTPParameter).deepCopy(),
		utils.DeepCopySlice[byte, byte](m.Data),
	}
	m.COTPParameterContract.(*_COTPParameter)._SubType = m
	return _COTPParameterDisconnectAdditionalInformationCopy
}

func (m *_COTPParameterDisconnectAdditionalInformation) String() string {
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
