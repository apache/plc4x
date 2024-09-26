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

// BACnetConstructedDataBBMDBroadcastDistributionTable is the corresponding interface of BACnetConstructedDataBBMDBroadcastDistributionTable
type BACnetConstructedDataBBMDBroadcastDistributionTable interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetBbmdBroadcastDistributionTable returns BbmdBroadcastDistributionTable (property field)
	GetBbmdBroadcastDistributionTable() []BACnetBDTEntry
	// IsBACnetConstructedDataBBMDBroadcastDistributionTable is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataBBMDBroadcastDistributionTable()
	// CreateBuilder creates a BACnetConstructedDataBBMDBroadcastDistributionTableBuilder
	CreateBACnetConstructedDataBBMDBroadcastDistributionTableBuilder() BACnetConstructedDataBBMDBroadcastDistributionTableBuilder
}

// _BACnetConstructedDataBBMDBroadcastDistributionTable is the data-structure of this message
type _BACnetConstructedDataBBMDBroadcastDistributionTable struct {
	BACnetConstructedDataContract
	BbmdBroadcastDistributionTable []BACnetBDTEntry
}

var _ BACnetConstructedDataBBMDBroadcastDistributionTable = (*_BACnetConstructedDataBBMDBroadcastDistributionTable)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataBBMDBroadcastDistributionTable)(nil)

// NewBACnetConstructedDataBBMDBroadcastDistributionTable factory function for _BACnetConstructedDataBBMDBroadcastDistributionTable
func NewBACnetConstructedDataBBMDBroadcastDistributionTable(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, bbmdBroadcastDistributionTable []BACnetBDTEntry, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataBBMDBroadcastDistributionTable {
	_result := &_BACnetConstructedDataBBMDBroadcastDistributionTable{
		BACnetConstructedDataContract:  NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		BbmdBroadcastDistributionTable: bbmdBroadcastDistributionTable,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataBBMDBroadcastDistributionTableBuilder is a builder for BACnetConstructedDataBBMDBroadcastDistributionTable
type BACnetConstructedDataBBMDBroadcastDistributionTableBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(bbmdBroadcastDistributionTable []BACnetBDTEntry) BACnetConstructedDataBBMDBroadcastDistributionTableBuilder
	// WithBbmdBroadcastDistributionTable adds BbmdBroadcastDistributionTable (property field)
	WithBbmdBroadcastDistributionTable(...BACnetBDTEntry) BACnetConstructedDataBBMDBroadcastDistributionTableBuilder
	// Build builds the BACnetConstructedDataBBMDBroadcastDistributionTable or returns an error if something is wrong
	Build() (BACnetConstructedDataBBMDBroadcastDistributionTable, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataBBMDBroadcastDistributionTable
}

// NewBACnetConstructedDataBBMDBroadcastDistributionTableBuilder() creates a BACnetConstructedDataBBMDBroadcastDistributionTableBuilder
func NewBACnetConstructedDataBBMDBroadcastDistributionTableBuilder() BACnetConstructedDataBBMDBroadcastDistributionTableBuilder {
	return &_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder{_BACnetConstructedDataBBMDBroadcastDistributionTable: new(_BACnetConstructedDataBBMDBroadcastDistributionTable)}
}

type _BACnetConstructedDataBBMDBroadcastDistributionTableBuilder struct {
	*_BACnetConstructedDataBBMDBroadcastDistributionTable

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataBBMDBroadcastDistributionTableBuilder) = (*_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder)(nil)

func (b *_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder) WithMandatoryFields(bbmdBroadcastDistributionTable []BACnetBDTEntry) BACnetConstructedDataBBMDBroadcastDistributionTableBuilder {
	return b.WithBbmdBroadcastDistributionTable(bbmdBroadcastDistributionTable...)
}

func (b *_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder) WithBbmdBroadcastDistributionTable(bbmdBroadcastDistributionTable ...BACnetBDTEntry) BACnetConstructedDataBBMDBroadcastDistributionTableBuilder {
	b.BbmdBroadcastDistributionTable = bbmdBroadcastDistributionTable
	return b
}

func (b *_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder) Build() (BACnetConstructedDataBBMDBroadcastDistributionTable, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataBBMDBroadcastDistributionTable.deepCopy(), nil
}

func (b *_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder) MustBuild() BACnetConstructedDataBBMDBroadcastDistributionTable {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataBBMDBroadcastDistributionTableBuilder().(*_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataBBMDBroadcastDistributionTableBuilder creates a BACnetConstructedDataBBMDBroadcastDistributionTableBuilder
func (b *_BACnetConstructedDataBBMDBroadcastDistributionTable) CreateBACnetConstructedDataBBMDBroadcastDistributionTableBuilder() BACnetConstructedDataBBMDBroadcastDistributionTableBuilder {
	if b == nil {
		return NewBACnetConstructedDataBBMDBroadcastDistributionTableBuilder()
	}
	return &_BACnetConstructedDataBBMDBroadcastDistributionTableBuilder{_BACnetConstructedDataBBMDBroadcastDistributionTable: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_BBMD_BROADCAST_DISTRIBUTION_TABLE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) GetBbmdBroadcastDistributionTable() []BACnetBDTEntry {
	return m.BbmdBroadcastDistributionTable
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataBBMDBroadcastDistributionTable(structType any) BACnetConstructedDataBBMDBroadcastDistributionTable {
	if casted, ok := structType.(BACnetConstructedDataBBMDBroadcastDistributionTable); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataBBMDBroadcastDistributionTable); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) GetTypeName() string {
	return "BACnetConstructedDataBBMDBroadcastDistributionTable"
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Array field
	if len(m.BbmdBroadcastDistributionTable) > 0 {
		for _, element := range m.BbmdBroadcastDistributionTable {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	return lengthInBits
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataBBMDBroadcastDistributionTable BACnetConstructedDataBBMDBroadcastDistributionTable, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataBBMDBroadcastDistributionTable"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataBBMDBroadcastDistributionTable")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	bbmdBroadcastDistributionTable, err := ReadTerminatedArrayField[BACnetBDTEntry](ctx, "bbmdBroadcastDistributionTable", ReadComplex[BACnetBDTEntry](BACnetBDTEntryParseWithBuffer, readBuffer), IsBACnetConstructedDataClosingTag(ctx, readBuffer, false, tagNumber))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'bbmdBroadcastDistributionTable' field"))
	}
	m.BbmdBroadcastDistributionTable = bbmdBroadcastDistributionTable

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataBBMDBroadcastDistributionTable"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataBBMDBroadcastDistributionTable")
	}

	return m, nil
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataBBMDBroadcastDistributionTable"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataBBMDBroadcastDistributionTable")
		}

		if err := WriteComplexTypeArrayField(ctx, "bbmdBroadcastDistributionTable", m.GetBbmdBroadcastDistributionTable(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'bbmdBroadcastDistributionTable' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataBBMDBroadcastDistributionTable"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataBBMDBroadcastDistributionTable")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) IsBACnetConstructedDataBBMDBroadcastDistributionTable() {
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) deepCopy() *_BACnetConstructedDataBBMDBroadcastDistributionTable {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataBBMDBroadcastDistributionTableCopy := &_BACnetConstructedDataBBMDBroadcastDistributionTable{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		utils.DeepCopySlice[BACnetBDTEntry, BACnetBDTEntry](m.BbmdBroadcastDistributionTable),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataBBMDBroadcastDistributionTableCopy
}

func (m *_BACnetConstructedDataBBMDBroadcastDistributionTable) String() string {
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
