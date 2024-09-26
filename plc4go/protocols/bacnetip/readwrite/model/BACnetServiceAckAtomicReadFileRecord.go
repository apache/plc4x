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

// BACnetServiceAckAtomicReadFileRecord is the corresponding interface of BACnetServiceAckAtomicReadFileRecord
type BACnetServiceAckAtomicReadFileRecord interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetServiceAckAtomicReadFileStreamOrRecord
	// GetFileStartRecord returns FileStartRecord (property field)
	GetFileStartRecord() BACnetApplicationTagSignedInteger
	// GetReturnedRecordCount returns ReturnedRecordCount (property field)
	GetReturnedRecordCount() BACnetApplicationTagUnsignedInteger
	// GetFileRecordData returns FileRecordData (property field)
	GetFileRecordData() []BACnetApplicationTagOctetString
	// IsBACnetServiceAckAtomicReadFileRecord is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetServiceAckAtomicReadFileRecord()
	// CreateBuilder creates a BACnetServiceAckAtomicReadFileRecordBuilder
	CreateBACnetServiceAckAtomicReadFileRecordBuilder() BACnetServiceAckAtomicReadFileRecordBuilder
}

// _BACnetServiceAckAtomicReadFileRecord is the data-structure of this message
type _BACnetServiceAckAtomicReadFileRecord struct {
	BACnetServiceAckAtomicReadFileStreamOrRecordContract
	FileStartRecord     BACnetApplicationTagSignedInteger
	ReturnedRecordCount BACnetApplicationTagUnsignedInteger
	FileRecordData      []BACnetApplicationTagOctetString
}

var _ BACnetServiceAckAtomicReadFileRecord = (*_BACnetServiceAckAtomicReadFileRecord)(nil)
var _ BACnetServiceAckAtomicReadFileStreamOrRecordRequirements = (*_BACnetServiceAckAtomicReadFileRecord)(nil)

// NewBACnetServiceAckAtomicReadFileRecord factory function for _BACnetServiceAckAtomicReadFileRecord
func NewBACnetServiceAckAtomicReadFileRecord(peekedTagHeader BACnetTagHeader, openingTag BACnetOpeningTag, closingTag BACnetClosingTag, fileStartRecord BACnetApplicationTagSignedInteger, returnedRecordCount BACnetApplicationTagUnsignedInteger, fileRecordData []BACnetApplicationTagOctetString) *_BACnetServiceAckAtomicReadFileRecord {
	if fileStartRecord == nil {
		panic("fileStartRecord of type BACnetApplicationTagSignedInteger for BACnetServiceAckAtomicReadFileRecord must not be nil")
	}
	if returnedRecordCount == nil {
		panic("returnedRecordCount of type BACnetApplicationTagUnsignedInteger for BACnetServiceAckAtomicReadFileRecord must not be nil")
	}
	_result := &_BACnetServiceAckAtomicReadFileRecord{
		BACnetServiceAckAtomicReadFileStreamOrRecordContract: NewBACnetServiceAckAtomicReadFileStreamOrRecord(peekedTagHeader, openingTag, closingTag),
		FileStartRecord:     fileStartRecord,
		ReturnedRecordCount: returnedRecordCount,
		FileRecordData:      fileRecordData,
	}
	_result.BACnetServiceAckAtomicReadFileStreamOrRecordContract.(*_BACnetServiceAckAtomicReadFileStreamOrRecord)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetServiceAckAtomicReadFileRecordBuilder is a builder for BACnetServiceAckAtomicReadFileRecord
type BACnetServiceAckAtomicReadFileRecordBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(fileStartRecord BACnetApplicationTagSignedInteger, returnedRecordCount BACnetApplicationTagUnsignedInteger, fileRecordData []BACnetApplicationTagOctetString) BACnetServiceAckAtomicReadFileRecordBuilder
	// WithFileStartRecord adds FileStartRecord (property field)
	WithFileStartRecord(BACnetApplicationTagSignedInteger) BACnetServiceAckAtomicReadFileRecordBuilder
	// WithFileStartRecordBuilder adds FileStartRecord (property field) which is build by the builder
	WithFileStartRecordBuilder(func(BACnetApplicationTagSignedIntegerBuilder) BACnetApplicationTagSignedIntegerBuilder) BACnetServiceAckAtomicReadFileRecordBuilder
	// WithReturnedRecordCount adds ReturnedRecordCount (property field)
	WithReturnedRecordCount(BACnetApplicationTagUnsignedInteger) BACnetServiceAckAtomicReadFileRecordBuilder
	// WithReturnedRecordCountBuilder adds ReturnedRecordCount (property field) which is build by the builder
	WithReturnedRecordCountBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetServiceAckAtomicReadFileRecordBuilder
	// WithFileRecordData adds FileRecordData (property field)
	WithFileRecordData(...BACnetApplicationTagOctetString) BACnetServiceAckAtomicReadFileRecordBuilder
	// Build builds the BACnetServiceAckAtomicReadFileRecord or returns an error if something is wrong
	Build() (BACnetServiceAckAtomicReadFileRecord, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetServiceAckAtomicReadFileRecord
}

// NewBACnetServiceAckAtomicReadFileRecordBuilder() creates a BACnetServiceAckAtomicReadFileRecordBuilder
func NewBACnetServiceAckAtomicReadFileRecordBuilder() BACnetServiceAckAtomicReadFileRecordBuilder {
	return &_BACnetServiceAckAtomicReadFileRecordBuilder{_BACnetServiceAckAtomicReadFileRecord: new(_BACnetServiceAckAtomicReadFileRecord)}
}

type _BACnetServiceAckAtomicReadFileRecordBuilder struct {
	*_BACnetServiceAckAtomicReadFileRecord

	parentBuilder *_BACnetServiceAckAtomicReadFileStreamOrRecordBuilder

	err *utils.MultiError
}

var _ (BACnetServiceAckAtomicReadFileRecordBuilder) = (*_BACnetServiceAckAtomicReadFileRecordBuilder)(nil)

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) setParent(contract BACnetServiceAckAtomicReadFileStreamOrRecordContract) {
	b.BACnetServiceAckAtomicReadFileStreamOrRecordContract = contract
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) WithMandatoryFields(fileStartRecord BACnetApplicationTagSignedInteger, returnedRecordCount BACnetApplicationTagUnsignedInteger, fileRecordData []BACnetApplicationTagOctetString) BACnetServiceAckAtomicReadFileRecordBuilder {
	return b.WithFileStartRecord(fileStartRecord).WithReturnedRecordCount(returnedRecordCount).WithFileRecordData(fileRecordData...)
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) WithFileStartRecord(fileStartRecord BACnetApplicationTagSignedInteger) BACnetServiceAckAtomicReadFileRecordBuilder {
	b.FileStartRecord = fileStartRecord
	return b
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) WithFileStartRecordBuilder(builderSupplier func(BACnetApplicationTagSignedIntegerBuilder) BACnetApplicationTagSignedIntegerBuilder) BACnetServiceAckAtomicReadFileRecordBuilder {
	builder := builderSupplier(b.FileStartRecord.CreateBACnetApplicationTagSignedIntegerBuilder())
	var err error
	b.FileStartRecord, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagSignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) WithReturnedRecordCount(returnedRecordCount BACnetApplicationTagUnsignedInteger) BACnetServiceAckAtomicReadFileRecordBuilder {
	b.ReturnedRecordCount = returnedRecordCount
	return b
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) WithReturnedRecordCountBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetServiceAckAtomicReadFileRecordBuilder {
	builder := builderSupplier(b.ReturnedRecordCount.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.ReturnedRecordCount, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) WithFileRecordData(fileRecordData ...BACnetApplicationTagOctetString) BACnetServiceAckAtomicReadFileRecordBuilder {
	b.FileRecordData = fileRecordData
	return b
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) Build() (BACnetServiceAckAtomicReadFileRecord, error) {
	if b.FileStartRecord == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'fileStartRecord' not set"))
	}
	if b.ReturnedRecordCount == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'returnedRecordCount' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetServiceAckAtomicReadFileRecord.deepCopy(), nil
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) MustBuild() BACnetServiceAckAtomicReadFileRecord {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) Done() BACnetServiceAckAtomicReadFileStreamOrRecordBuilder {
	return b.parentBuilder
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) buildForBACnetServiceAckAtomicReadFileStreamOrRecord() (BACnetServiceAckAtomicReadFileStreamOrRecord, error) {
	return b.Build()
}

func (b *_BACnetServiceAckAtomicReadFileRecordBuilder) DeepCopy() any {
	_copy := b.CreateBACnetServiceAckAtomicReadFileRecordBuilder().(*_BACnetServiceAckAtomicReadFileRecordBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetServiceAckAtomicReadFileRecordBuilder creates a BACnetServiceAckAtomicReadFileRecordBuilder
func (b *_BACnetServiceAckAtomicReadFileRecord) CreateBACnetServiceAckAtomicReadFileRecordBuilder() BACnetServiceAckAtomicReadFileRecordBuilder {
	if b == nil {
		return NewBACnetServiceAckAtomicReadFileRecordBuilder()
	}
	return &_BACnetServiceAckAtomicReadFileRecordBuilder{_BACnetServiceAckAtomicReadFileRecord: b.deepCopy()}
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

func (m *_BACnetServiceAckAtomicReadFileRecord) GetParent() BACnetServiceAckAtomicReadFileStreamOrRecordContract {
	return m.BACnetServiceAckAtomicReadFileStreamOrRecordContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetServiceAckAtomicReadFileRecord) GetFileStartRecord() BACnetApplicationTagSignedInteger {
	return m.FileStartRecord
}

func (m *_BACnetServiceAckAtomicReadFileRecord) GetReturnedRecordCount() BACnetApplicationTagUnsignedInteger {
	return m.ReturnedRecordCount
}

func (m *_BACnetServiceAckAtomicReadFileRecord) GetFileRecordData() []BACnetApplicationTagOctetString {
	return m.FileRecordData
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetServiceAckAtomicReadFileRecord(structType any) BACnetServiceAckAtomicReadFileRecord {
	if casted, ok := structType.(BACnetServiceAckAtomicReadFileRecord); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetServiceAckAtomicReadFileRecord); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetServiceAckAtomicReadFileRecord) GetTypeName() string {
	return "BACnetServiceAckAtomicReadFileRecord"
}

func (m *_BACnetServiceAckAtomicReadFileRecord) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetServiceAckAtomicReadFileStreamOrRecordContract.(*_BACnetServiceAckAtomicReadFileStreamOrRecord).GetLengthInBits(ctx))

	// Simple field (fileStartRecord)
	lengthInBits += m.FileStartRecord.GetLengthInBits(ctx)

	// Simple field (returnedRecordCount)
	lengthInBits += m.ReturnedRecordCount.GetLengthInBits(ctx)

	// Array field
	if len(m.FileRecordData) > 0 {
		for _curItem, element := range m.FileRecordData {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.FileRecordData), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_BACnetServiceAckAtomicReadFileRecord) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetServiceAckAtomicReadFileRecord) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetServiceAckAtomicReadFileStreamOrRecord) (__bACnetServiceAckAtomicReadFileRecord BACnetServiceAckAtomicReadFileRecord, err error) {
	m.BACnetServiceAckAtomicReadFileStreamOrRecordContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetServiceAckAtomicReadFileRecord"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetServiceAckAtomicReadFileRecord")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	fileStartRecord, err := ReadSimpleField[BACnetApplicationTagSignedInteger](ctx, "fileStartRecord", ReadComplex[BACnetApplicationTagSignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagSignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'fileStartRecord' field"))
	}
	m.FileStartRecord = fileStartRecord

	returnedRecordCount, err := ReadSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "returnedRecordCount", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'returnedRecordCount' field"))
	}
	m.ReturnedRecordCount = returnedRecordCount

	fileRecordData, err := ReadCountArrayField[BACnetApplicationTagOctetString](ctx, "fileRecordData", ReadComplex[BACnetApplicationTagOctetString](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagOctetString](), readBuffer), uint64(returnedRecordCount.GetPayload().GetActualValue()))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'fileRecordData' field"))
	}
	m.FileRecordData = fileRecordData

	if closeErr := readBuffer.CloseContext("BACnetServiceAckAtomicReadFileRecord"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetServiceAckAtomicReadFileRecord")
	}

	return m, nil
}

func (m *_BACnetServiceAckAtomicReadFileRecord) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetServiceAckAtomicReadFileRecord) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetServiceAckAtomicReadFileRecord"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetServiceAckAtomicReadFileRecord")
		}

		if err := WriteSimpleField[BACnetApplicationTagSignedInteger](ctx, "fileStartRecord", m.GetFileStartRecord(), WriteComplex[BACnetApplicationTagSignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'fileStartRecord' field")
		}

		if err := WriteSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "returnedRecordCount", m.GetReturnedRecordCount(), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'returnedRecordCount' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "fileRecordData", m.GetFileRecordData(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'fileRecordData' field")
		}

		if popErr := writeBuffer.PopContext("BACnetServiceAckAtomicReadFileRecord"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetServiceAckAtomicReadFileRecord")
		}
		return nil
	}
	return m.BACnetServiceAckAtomicReadFileStreamOrRecordContract.(*_BACnetServiceAckAtomicReadFileStreamOrRecord).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetServiceAckAtomicReadFileRecord) IsBACnetServiceAckAtomicReadFileRecord() {}

func (m *_BACnetServiceAckAtomicReadFileRecord) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetServiceAckAtomicReadFileRecord) deepCopy() *_BACnetServiceAckAtomicReadFileRecord {
	if m == nil {
		return nil
	}
	_BACnetServiceAckAtomicReadFileRecordCopy := &_BACnetServiceAckAtomicReadFileRecord{
		m.BACnetServiceAckAtomicReadFileStreamOrRecordContract.(*_BACnetServiceAckAtomicReadFileStreamOrRecord).deepCopy(),
		m.FileStartRecord.DeepCopy().(BACnetApplicationTagSignedInteger),
		m.ReturnedRecordCount.DeepCopy().(BACnetApplicationTagUnsignedInteger),
		utils.DeepCopySlice[BACnetApplicationTagOctetString, BACnetApplicationTagOctetString](m.FileRecordData),
	}
	m.BACnetServiceAckAtomicReadFileStreamOrRecordContract.(*_BACnetServiceAckAtomicReadFileStreamOrRecord)._SubType = m
	return _BACnetServiceAckAtomicReadFileRecordCopy
}

func (m *_BACnetServiceAckAtomicReadFileRecord) String() string {
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
