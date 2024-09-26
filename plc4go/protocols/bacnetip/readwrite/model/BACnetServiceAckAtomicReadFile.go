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

// BACnetServiceAckAtomicReadFile is the corresponding interface of BACnetServiceAckAtomicReadFile
type BACnetServiceAckAtomicReadFile interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetServiceAck
	// GetEndOfFile returns EndOfFile (property field)
	GetEndOfFile() BACnetApplicationTagBoolean
	// GetAccessMethod returns AccessMethod (property field)
	GetAccessMethod() BACnetServiceAckAtomicReadFileStreamOrRecord
	// IsBACnetServiceAckAtomicReadFile is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetServiceAckAtomicReadFile()
	// CreateBuilder creates a BACnetServiceAckAtomicReadFileBuilder
	CreateBACnetServiceAckAtomicReadFileBuilder() BACnetServiceAckAtomicReadFileBuilder
}

// _BACnetServiceAckAtomicReadFile is the data-structure of this message
type _BACnetServiceAckAtomicReadFile struct {
	BACnetServiceAckContract
	EndOfFile    BACnetApplicationTagBoolean
	AccessMethod BACnetServiceAckAtomicReadFileStreamOrRecord
}

var _ BACnetServiceAckAtomicReadFile = (*_BACnetServiceAckAtomicReadFile)(nil)
var _ BACnetServiceAckRequirements = (*_BACnetServiceAckAtomicReadFile)(nil)

// NewBACnetServiceAckAtomicReadFile factory function for _BACnetServiceAckAtomicReadFile
func NewBACnetServiceAckAtomicReadFile(endOfFile BACnetApplicationTagBoolean, accessMethod BACnetServiceAckAtomicReadFileStreamOrRecord, serviceAckLength uint32) *_BACnetServiceAckAtomicReadFile {
	if endOfFile == nil {
		panic("endOfFile of type BACnetApplicationTagBoolean for BACnetServiceAckAtomicReadFile must not be nil")
	}
	if accessMethod == nil {
		panic("accessMethod of type BACnetServiceAckAtomicReadFileStreamOrRecord for BACnetServiceAckAtomicReadFile must not be nil")
	}
	_result := &_BACnetServiceAckAtomicReadFile{
		BACnetServiceAckContract: NewBACnetServiceAck(serviceAckLength),
		EndOfFile:                endOfFile,
		AccessMethod:             accessMethod,
	}
	_result.BACnetServiceAckContract.(*_BACnetServiceAck)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetServiceAckAtomicReadFileBuilder is a builder for BACnetServiceAckAtomicReadFile
type BACnetServiceAckAtomicReadFileBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(endOfFile BACnetApplicationTagBoolean, accessMethod BACnetServiceAckAtomicReadFileStreamOrRecord) BACnetServiceAckAtomicReadFileBuilder
	// WithEndOfFile adds EndOfFile (property field)
	WithEndOfFile(BACnetApplicationTagBoolean) BACnetServiceAckAtomicReadFileBuilder
	// WithEndOfFileBuilder adds EndOfFile (property field) which is build by the builder
	WithEndOfFileBuilder(func(BACnetApplicationTagBooleanBuilder) BACnetApplicationTagBooleanBuilder) BACnetServiceAckAtomicReadFileBuilder
	// WithAccessMethod adds AccessMethod (property field)
	WithAccessMethod(BACnetServiceAckAtomicReadFileStreamOrRecord) BACnetServiceAckAtomicReadFileBuilder
	// WithAccessMethodBuilder adds AccessMethod (property field) which is build by the builder
	WithAccessMethodBuilder(func(BACnetServiceAckAtomicReadFileStreamOrRecordBuilder) BACnetServiceAckAtomicReadFileStreamOrRecordBuilder) BACnetServiceAckAtomicReadFileBuilder
	// Build builds the BACnetServiceAckAtomicReadFile or returns an error if something is wrong
	Build() (BACnetServiceAckAtomicReadFile, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetServiceAckAtomicReadFile
}

// NewBACnetServiceAckAtomicReadFileBuilder() creates a BACnetServiceAckAtomicReadFileBuilder
func NewBACnetServiceAckAtomicReadFileBuilder() BACnetServiceAckAtomicReadFileBuilder {
	return &_BACnetServiceAckAtomicReadFileBuilder{_BACnetServiceAckAtomicReadFile: new(_BACnetServiceAckAtomicReadFile)}
}

type _BACnetServiceAckAtomicReadFileBuilder struct {
	*_BACnetServiceAckAtomicReadFile

	parentBuilder *_BACnetServiceAckBuilder

	err *utils.MultiError
}

var _ (BACnetServiceAckAtomicReadFileBuilder) = (*_BACnetServiceAckAtomicReadFileBuilder)(nil)

func (b *_BACnetServiceAckAtomicReadFileBuilder) setParent(contract BACnetServiceAckContract) {
	b.BACnetServiceAckContract = contract
}

func (b *_BACnetServiceAckAtomicReadFileBuilder) WithMandatoryFields(endOfFile BACnetApplicationTagBoolean, accessMethod BACnetServiceAckAtomicReadFileStreamOrRecord) BACnetServiceAckAtomicReadFileBuilder {
	return b.WithEndOfFile(endOfFile).WithAccessMethod(accessMethod)
}

func (b *_BACnetServiceAckAtomicReadFileBuilder) WithEndOfFile(endOfFile BACnetApplicationTagBoolean) BACnetServiceAckAtomicReadFileBuilder {
	b.EndOfFile = endOfFile
	return b
}

func (b *_BACnetServiceAckAtomicReadFileBuilder) WithEndOfFileBuilder(builderSupplier func(BACnetApplicationTagBooleanBuilder) BACnetApplicationTagBooleanBuilder) BACnetServiceAckAtomicReadFileBuilder {
	builder := builderSupplier(b.EndOfFile.CreateBACnetApplicationTagBooleanBuilder())
	var err error
	b.EndOfFile, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagBooleanBuilder failed"))
	}
	return b
}

func (b *_BACnetServiceAckAtomicReadFileBuilder) WithAccessMethod(accessMethod BACnetServiceAckAtomicReadFileStreamOrRecord) BACnetServiceAckAtomicReadFileBuilder {
	b.AccessMethod = accessMethod
	return b
}

func (b *_BACnetServiceAckAtomicReadFileBuilder) WithAccessMethodBuilder(builderSupplier func(BACnetServiceAckAtomicReadFileStreamOrRecordBuilder) BACnetServiceAckAtomicReadFileStreamOrRecordBuilder) BACnetServiceAckAtomicReadFileBuilder {
	builder := builderSupplier(b.AccessMethod.CreateBACnetServiceAckAtomicReadFileStreamOrRecordBuilder())
	var err error
	b.AccessMethod, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetServiceAckAtomicReadFileStreamOrRecordBuilder failed"))
	}
	return b
}

func (b *_BACnetServiceAckAtomicReadFileBuilder) Build() (BACnetServiceAckAtomicReadFile, error) {
	if b.EndOfFile == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'endOfFile' not set"))
	}
	if b.AccessMethod == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'accessMethod' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetServiceAckAtomicReadFile.deepCopy(), nil
}

func (b *_BACnetServiceAckAtomicReadFileBuilder) MustBuild() BACnetServiceAckAtomicReadFile {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetServiceAckAtomicReadFileBuilder) Done() BACnetServiceAckBuilder {
	return b.parentBuilder
}

func (b *_BACnetServiceAckAtomicReadFileBuilder) buildForBACnetServiceAck() (BACnetServiceAck, error) {
	return b.Build()
}

func (b *_BACnetServiceAckAtomicReadFileBuilder) DeepCopy() any {
	_copy := b.CreateBACnetServiceAckAtomicReadFileBuilder().(*_BACnetServiceAckAtomicReadFileBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetServiceAckAtomicReadFileBuilder creates a BACnetServiceAckAtomicReadFileBuilder
func (b *_BACnetServiceAckAtomicReadFile) CreateBACnetServiceAckAtomicReadFileBuilder() BACnetServiceAckAtomicReadFileBuilder {
	if b == nil {
		return NewBACnetServiceAckAtomicReadFileBuilder()
	}
	return &_BACnetServiceAckAtomicReadFileBuilder{_BACnetServiceAckAtomicReadFile: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetServiceAckAtomicReadFile) GetServiceChoice() BACnetConfirmedServiceChoice {
	return BACnetConfirmedServiceChoice_ATOMIC_READ_FILE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetServiceAckAtomicReadFile) GetParent() BACnetServiceAckContract {
	return m.BACnetServiceAckContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetServiceAckAtomicReadFile) GetEndOfFile() BACnetApplicationTagBoolean {
	return m.EndOfFile
}

func (m *_BACnetServiceAckAtomicReadFile) GetAccessMethod() BACnetServiceAckAtomicReadFileStreamOrRecord {
	return m.AccessMethod
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetServiceAckAtomicReadFile(structType any) BACnetServiceAckAtomicReadFile {
	if casted, ok := structType.(BACnetServiceAckAtomicReadFile); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetServiceAckAtomicReadFile); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetServiceAckAtomicReadFile) GetTypeName() string {
	return "BACnetServiceAckAtomicReadFile"
}

func (m *_BACnetServiceAckAtomicReadFile) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetServiceAckContract.(*_BACnetServiceAck).GetLengthInBits(ctx))

	// Simple field (endOfFile)
	lengthInBits += m.EndOfFile.GetLengthInBits(ctx)

	// Simple field (accessMethod)
	lengthInBits += m.AccessMethod.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetServiceAckAtomicReadFile) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetServiceAckAtomicReadFile) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetServiceAck, serviceAckLength uint32) (__bACnetServiceAckAtomicReadFile BACnetServiceAckAtomicReadFile, err error) {
	m.BACnetServiceAckContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetServiceAckAtomicReadFile"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetServiceAckAtomicReadFile")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	endOfFile, err := ReadSimpleField[BACnetApplicationTagBoolean](ctx, "endOfFile", ReadComplex[BACnetApplicationTagBoolean](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagBoolean](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'endOfFile' field"))
	}
	m.EndOfFile = endOfFile

	accessMethod, err := ReadSimpleField[BACnetServiceAckAtomicReadFileStreamOrRecord](ctx, "accessMethod", ReadComplex[BACnetServiceAckAtomicReadFileStreamOrRecord](BACnetServiceAckAtomicReadFileStreamOrRecordParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'accessMethod' field"))
	}
	m.AccessMethod = accessMethod

	if closeErr := readBuffer.CloseContext("BACnetServiceAckAtomicReadFile"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetServiceAckAtomicReadFile")
	}

	return m, nil
}

func (m *_BACnetServiceAckAtomicReadFile) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetServiceAckAtomicReadFile) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetServiceAckAtomicReadFile"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetServiceAckAtomicReadFile")
		}

		if err := WriteSimpleField[BACnetApplicationTagBoolean](ctx, "endOfFile", m.GetEndOfFile(), WriteComplex[BACnetApplicationTagBoolean](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'endOfFile' field")
		}

		if err := WriteSimpleField[BACnetServiceAckAtomicReadFileStreamOrRecord](ctx, "accessMethod", m.GetAccessMethod(), WriteComplex[BACnetServiceAckAtomicReadFileStreamOrRecord](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'accessMethod' field")
		}

		if popErr := writeBuffer.PopContext("BACnetServiceAckAtomicReadFile"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetServiceAckAtomicReadFile")
		}
		return nil
	}
	return m.BACnetServiceAckContract.(*_BACnetServiceAck).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetServiceAckAtomicReadFile) IsBACnetServiceAckAtomicReadFile() {}

func (m *_BACnetServiceAckAtomicReadFile) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetServiceAckAtomicReadFile) deepCopy() *_BACnetServiceAckAtomicReadFile {
	if m == nil {
		return nil
	}
	_BACnetServiceAckAtomicReadFileCopy := &_BACnetServiceAckAtomicReadFile{
		m.BACnetServiceAckContract.(*_BACnetServiceAck).deepCopy(),
		m.EndOfFile.DeepCopy().(BACnetApplicationTagBoolean),
		m.AccessMethod.DeepCopy().(BACnetServiceAckAtomicReadFileStreamOrRecord),
	}
	m.BACnetServiceAckContract.(*_BACnetServiceAck)._SubType = m
	return _BACnetServiceAckAtomicReadFileCopy
}

func (m *_BACnetServiceAckAtomicReadFile) String() string {
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
