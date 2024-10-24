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

// BACnetUnconfirmedServiceRequestWriteGroup is the corresponding interface of BACnetUnconfirmedServiceRequestWriteGroup
type BACnetUnconfirmedServiceRequestWriteGroup interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetUnconfirmedServiceRequest
	// GetGroupNumber returns GroupNumber (property field)
	GetGroupNumber() BACnetContextTagUnsignedInteger
	// GetWritePriority returns WritePriority (property field)
	GetWritePriority() BACnetContextTagUnsignedInteger
	// GetChangeList returns ChangeList (property field)
	GetChangeList() BACnetGroupChannelValueList
	// GetInhibitDelay returns InhibitDelay (property field)
	GetInhibitDelay() BACnetContextTagUnsignedInteger
	// IsBACnetUnconfirmedServiceRequestWriteGroup is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetUnconfirmedServiceRequestWriteGroup()
	// CreateBuilder creates a BACnetUnconfirmedServiceRequestWriteGroupBuilder
	CreateBACnetUnconfirmedServiceRequestWriteGroupBuilder() BACnetUnconfirmedServiceRequestWriteGroupBuilder
}

// _BACnetUnconfirmedServiceRequestWriteGroup is the data-structure of this message
type _BACnetUnconfirmedServiceRequestWriteGroup struct {
	BACnetUnconfirmedServiceRequestContract
	GroupNumber   BACnetContextTagUnsignedInteger
	WritePriority BACnetContextTagUnsignedInteger
	ChangeList    BACnetGroupChannelValueList
	InhibitDelay  BACnetContextTagUnsignedInteger
}

var _ BACnetUnconfirmedServiceRequestWriteGroup = (*_BACnetUnconfirmedServiceRequestWriteGroup)(nil)
var _ BACnetUnconfirmedServiceRequestRequirements = (*_BACnetUnconfirmedServiceRequestWriteGroup)(nil)

// NewBACnetUnconfirmedServiceRequestWriteGroup factory function for _BACnetUnconfirmedServiceRequestWriteGroup
func NewBACnetUnconfirmedServiceRequestWriteGroup(groupNumber BACnetContextTagUnsignedInteger, writePriority BACnetContextTagUnsignedInteger, changeList BACnetGroupChannelValueList, inhibitDelay BACnetContextTagUnsignedInteger, serviceRequestLength uint16) *_BACnetUnconfirmedServiceRequestWriteGroup {
	if groupNumber == nil {
		panic("groupNumber of type BACnetContextTagUnsignedInteger for BACnetUnconfirmedServiceRequestWriteGroup must not be nil")
	}
	if writePriority == nil {
		panic("writePriority of type BACnetContextTagUnsignedInteger for BACnetUnconfirmedServiceRequestWriteGroup must not be nil")
	}
	if changeList == nil {
		panic("changeList of type BACnetGroupChannelValueList for BACnetUnconfirmedServiceRequestWriteGroup must not be nil")
	}
	_result := &_BACnetUnconfirmedServiceRequestWriteGroup{
		BACnetUnconfirmedServiceRequestContract: NewBACnetUnconfirmedServiceRequest(serviceRequestLength),
		GroupNumber:                             groupNumber,
		WritePriority:                           writePriority,
		ChangeList:                              changeList,
		InhibitDelay:                            inhibitDelay,
	}
	_result.BACnetUnconfirmedServiceRequestContract.(*_BACnetUnconfirmedServiceRequest)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetUnconfirmedServiceRequestWriteGroupBuilder is a builder for BACnetUnconfirmedServiceRequestWriteGroup
type BACnetUnconfirmedServiceRequestWriteGroupBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(groupNumber BACnetContextTagUnsignedInteger, writePriority BACnetContextTagUnsignedInteger, changeList BACnetGroupChannelValueList) BACnetUnconfirmedServiceRequestWriteGroupBuilder
	// WithGroupNumber adds GroupNumber (property field)
	WithGroupNumber(BACnetContextTagUnsignedInteger) BACnetUnconfirmedServiceRequestWriteGroupBuilder
	// WithGroupNumberBuilder adds GroupNumber (property field) which is build by the builder
	WithGroupNumberBuilder(func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetUnconfirmedServiceRequestWriteGroupBuilder
	// WithWritePriority adds WritePriority (property field)
	WithWritePriority(BACnetContextTagUnsignedInteger) BACnetUnconfirmedServiceRequestWriteGroupBuilder
	// WithWritePriorityBuilder adds WritePriority (property field) which is build by the builder
	WithWritePriorityBuilder(func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetUnconfirmedServiceRequestWriteGroupBuilder
	// WithChangeList adds ChangeList (property field)
	WithChangeList(BACnetGroupChannelValueList) BACnetUnconfirmedServiceRequestWriteGroupBuilder
	// WithChangeListBuilder adds ChangeList (property field) which is build by the builder
	WithChangeListBuilder(func(BACnetGroupChannelValueListBuilder) BACnetGroupChannelValueListBuilder) BACnetUnconfirmedServiceRequestWriteGroupBuilder
	// WithInhibitDelay adds InhibitDelay (property field)
	WithOptionalInhibitDelay(BACnetContextTagUnsignedInteger) BACnetUnconfirmedServiceRequestWriteGroupBuilder
	// WithOptionalInhibitDelayBuilder adds InhibitDelay (property field) which is build by the builder
	WithOptionalInhibitDelayBuilder(func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetUnconfirmedServiceRequestWriteGroupBuilder
	// Build builds the BACnetUnconfirmedServiceRequestWriteGroup or returns an error if something is wrong
	Build() (BACnetUnconfirmedServiceRequestWriteGroup, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetUnconfirmedServiceRequestWriteGroup
}

// NewBACnetUnconfirmedServiceRequestWriteGroupBuilder() creates a BACnetUnconfirmedServiceRequestWriteGroupBuilder
func NewBACnetUnconfirmedServiceRequestWriteGroupBuilder() BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	return &_BACnetUnconfirmedServiceRequestWriteGroupBuilder{_BACnetUnconfirmedServiceRequestWriteGroup: new(_BACnetUnconfirmedServiceRequestWriteGroup)}
}

type _BACnetUnconfirmedServiceRequestWriteGroupBuilder struct {
	*_BACnetUnconfirmedServiceRequestWriteGroup

	parentBuilder *_BACnetUnconfirmedServiceRequestBuilder

	err *utils.MultiError
}

var _ (BACnetUnconfirmedServiceRequestWriteGroupBuilder) = (*_BACnetUnconfirmedServiceRequestWriteGroupBuilder)(nil)

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) setParent(contract BACnetUnconfirmedServiceRequestContract) {
	b.BACnetUnconfirmedServiceRequestContract = contract
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) WithMandatoryFields(groupNumber BACnetContextTagUnsignedInteger, writePriority BACnetContextTagUnsignedInteger, changeList BACnetGroupChannelValueList) BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	return b.WithGroupNumber(groupNumber).WithWritePriority(writePriority).WithChangeList(changeList)
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) WithGroupNumber(groupNumber BACnetContextTagUnsignedInteger) BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	b.GroupNumber = groupNumber
	return b
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) WithGroupNumberBuilder(builderSupplier func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	builder := builderSupplier(b.GroupNumber.CreateBACnetContextTagUnsignedIntegerBuilder())
	var err error
	b.GroupNumber, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) WithWritePriority(writePriority BACnetContextTagUnsignedInteger) BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	b.WritePriority = writePriority
	return b
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) WithWritePriorityBuilder(builderSupplier func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	builder := builderSupplier(b.WritePriority.CreateBACnetContextTagUnsignedIntegerBuilder())
	var err error
	b.WritePriority, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) WithChangeList(changeList BACnetGroupChannelValueList) BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	b.ChangeList = changeList
	return b
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) WithChangeListBuilder(builderSupplier func(BACnetGroupChannelValueListBuilder) BACnetGroupChannelValueListBuilder) BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	builder := builderSupplier(b.ChangeList.CreateBACnetGroupChannelValueListBuilder())
	var err error
	b.ChangeList, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetGroupChannelValueListBuilder failed"))
	}
	return b
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) WithOptionalInhibitDelay(inhibitDelay BACnetContextTagUnsignedInteger) BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	b.InhibitDelay = inhibitDelay
	return b
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) WithOptionalInhibitDelayBuilder(builderSupplier func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	builder := builderSupplier(b.InhibitDelay.CreateBACnetContextTagUnsignedIntegerBuilder())
	var err error
	b.InhibitDelay, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) Build() (BACnetUnconfirmedServiceRequestWriteGroup, error) {
	if b.GroupNumber == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'groupNumber' not set"))
	}
	if b.WritePriority == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'writePriority' not set"))
	}
	if b.ChangeList == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'changeList' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetUnconfirmedServiceRequestWriteGroup.deepCopy(), nil
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) MustBuild() BACnetUnconfirmedServiceRequestWriteGroup {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) Done() BACnetUnconfirmedServiceRequestBuilder {
	return b.parentBuilder
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) buildForBACnetUnconfirmedServiceRequest() (BACnetUnconfirmedServiceRequest, error) {
	return b.Build()
}

func (b *_BACnetUnconfirmedServiceRequestWriteGroupBuilder) DeepCopy() any {
	_copy := b.CreateBACnetUnconfirmedServiceRequestWriteGroupBuilder().(*_BACnetUnconfirmedServiceRequestWriteGroupBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetUnconfirmedServiceRequestWriteGroupBuilder creates a BACnetUnconfirmedServiceRequestWriteGroupBuilder
func (b *_BACnetUnconfirmedServiceRequestWriteGroup) CreateBACnetUnconfirmedServiceRequestWriteGroupBuilder() BACnetUnconfirmedServiceRequestWriteGroupBuilder {
	if b == nil {
		return NewBACnetUnconfirmedServiceRequestWriteGroupBuilder()
	}
	return &_BACnetUnconfirmedServiceRequestWriteGroupBuilder{_BACnetUnconfirmedServiceRequestWriteGroup: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) GetServiceChoice() BACnetUnconfirmedServiceChoice {
	return BACnetUnconfirmedServiceChoice_WRITE_GROUP
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) GetParent() BACnetUnconfirmedServiceRequestContract {
	return m.BACnetUnconfirmedServiceRequestContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) GetGroupNumber() BACnetContextTagUnsignedInteger {
	return m.GroupNumber
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) GetWritePriority() BACnetContextTagUnsignedInteger {
	return m.WritePriority
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) GetChangeList() BACnetGroupChannelValueList {
	return m.ChangeList
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) GetInhibitDelay() BACnetContextTagUnsignedInteger {
	return m.InhibitDelay
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetUnconfirmedServiceRequestWriteGroup(structType any) BACnetUnconfirmedServiceRequestWriteGroup {
	if casted, ok := structType.(BACnetUnconfirmedServiceRequestWriteGroup); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetUnconfirmedServiceRequestWriteGroup); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) GetTypeName() string {
	return "BACnetUnconfirmedServiceRequestWriteGroup"
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetUnconfirmedServiceRequestContract.(*_BACnetUnconfirmedServiceRequest).GetLengthInBits(ctx))

	// Simple field (groupNumber)
	lengthInBits += m.GroupNumber.GetLengthInBits(ctx)

	// Simple field (writePriority)
	lengthInBits += m.WritePriority.GetLengthInBits(ctx)

	// Simple field (changeList)
	lengthInBits += m.ChangeList.GetLengthInBits(ctx)

	// Optional Field (inhibitDelay)
	if m.InhibitDelay != nil {
		lengthInBits += m.InhibitDelay.GetLengthInBits(ctx)
	}

	return lengthInBits
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetUnconfirmedServiceRequest, serviceRequestLength uint16) (__bACnetUnconfirmedServiceRequestWriteGroup BACnetUnconfirmedServiceRequestWriteGroup, err error) {
	m.BACnetUnconfirmedServiceRequestContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetUnconfirmedServiceRequestWriteGroup"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetUnconfirmedServiceRequestWriteGroup")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	groupNumber, err := ReadSimpleField[BACnetContextTagUnsignedInteger](ctx, "groupNumber", ReadComplex[BACnetContextTagUnsignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagUnsignedInteger]((uint8)(uint8(0)), (BACnetDataType)(BACnetDataType_UNSIGNED_INTEGER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'groupNumber' field"))
	}
	m.GroupNumber = groupNumber

	writePriority, err := ReadSimpleField[BACnetContextTagUnsignedInteger](ctx, "writePriority", ReadComplex[BACnetContextTagUnsignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagUnsignedInteger]((uint8)(uint8(1)), (BACnetDataType)(BACnetDataType_UNSIGNED_INTEGER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'writePriority' field"))
	}
	m.WritePriority = writePriority

	changeList, err := ReadSimpleField[BACnetGroupChannelValueList](ctx, "changeList", ReadComplex[BACnetGroupChannelValueList](BACnetGroupChannelValueListParseWithBufferProducer((uint8)(uint8(2))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'changeList' field"))
	}
	m.ChangeList = changeList

	var inhibitDelay BACnetContextTagUnsignedInteger
	_inhibitDelay, err := ReadOptionalField[BACnetContextTagUnsignedInteger](ctx, "inhibitDelay", ReadComplex[BACnetContextTagUnsignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagUnsignedInteger]((uint8)(uint8(3)), (BACnetDataType)(BACnetDataType_UNSIGNED_INTEGER)), readBuffer), true)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'inhibitDelay' field"))
	}
	if _inhibitDelay != nil {
		inhibitDelay = *_inhibitDelay
		m.InhibitDelay = inhibitDelay
	}

	if closeErr := readBuffer.CloseContext("BACnetUnconfirmedServiceRequestWriteGroup"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetUnconfirmedServiceRequestWriteGroup")
	}

	return m, nil
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetUnconfirmedServiceRequestWriteGroup"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetUnconfirmedServiceRequestWriteGroup")
		}

		if err := WriteSimpleField[BACnetContextTagUnsignedInteger](ctx, "groupNumber", m.GetGroupNumber(), WriteComplex[BACnetContextTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'groupNumber' field")
		}

		if err := WriteSimpleField[BACnetContextTagUnsignedInteger](ctx, "writePriority", m.GetWritePriority(), WriteComplex[BACnetContextTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'writePriority' field")
		}

		if err := WriteSimpleField[BACnetGroupChannelValueList](ctx, "changeList", m.GetChangeList(), WriteComplex[BACnetGroupChannelValueList](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'changeList' field")
		}

		if err := WriteOptionalField[BACnetContextTagUnsignedInteger](ctx, "inhibitDelay", GetRef(m.GetInhibitDelay()), WriteComplex[BACnetContextTagUnsignedInteger](writeBuffer), true); err != nil {
			return errors.Wrap(err, "Error serializing 'inhibitDelay' field")
		}

		if popErr := writeBuffer.PopContext("BACnetUnconfirmedServiceRequestWriteGroup"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetUnconfirmedServiceRequestWriteGroup")
		}
		return nil
	}
	return m.BACnetUnconfirmedServiceRequestContract.(*_BACnetUnconfirmedServiceRequest).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) IsBACnetUnconfirmedServiceRequestWriteGroup() {}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) deepCopy() *_BACnetUnconfirmedServiceRequestWriteGroup {
	if m == nil {
		return nil
	}
	_BACnetUnconfirmedServiceRequestWriteGroupCopy := &_BACnetUnconfirmedServiceRequestWriteGroup{
		m.BACnetUnconfirmedServiceRequestContract.(*_BACnetUnconfirmedServiceRequest).deepCopy(),
		m.GroupNumber.DeepCopy().(BACnetContextTagUnsignedInteger),
		m.WritePriority.DeepCopy().(BACnetContextTagUnsignedInteger),
		m.ChangeList.DeepCopy().(BACnetGroupChannelValueList),
		m.InhibitDelay.DeepCopy().(BACnetContextTagUnsignedInteger),
	}
	m.BACnetUnconfirmedServiceRequestContract.(*_BACnetUnconfirmedServiceRequest)._SubType = m
	return _BACnetUnconfirmedServiceRequestWriteGroupCopy
}

func (m *_BACnetUnconfirmedServiceRequestWriteGroup) String() string {
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
