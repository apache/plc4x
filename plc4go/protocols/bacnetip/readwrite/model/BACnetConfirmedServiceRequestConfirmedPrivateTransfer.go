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

// BACnetConfirmedServiceRequestConfirmedPrivateTransfer is the corresponding interface of BACnetConfirmedServiceRequestConfirmedPrivateTransfer
type BACnetConfirmedServiceRequestConfirmedPrivateTransfer interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConfirmedServiceRequest
	// GetVendorId returns VendorId (property field)
	GetVendorId() BACnetVendorIdTagged
	// GetServiceNumber returns ServiceNumber (property field)
	GetServiceNumber() BACnetContextTagUnsignedInteger
	// GetServiceParameters returns ServiceParameters (property field)
	GetServiceParameters() BACnetConstructedData
	// IsBACnetConfirmedServiceRequestConfirmedPrivateTransfer is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConfirmedServiceRequestConfirmedPrivateTransfer()
	// CreateBuilder creates a BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
	CreateBACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder() BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
}

// _BACnetConfirmedServiceRequestConfirmedPrivateTransfer is the data-structure of this message
type _BACnetConfirmedServiceRequestConfirmedPrivateTransfer struct {
	BACnetConfirmedServiceRequestContract
	VendorId          BACnetVendorIdTagged
	ServiceNumber     BACnetContextTagUnsignedInteger
	ServiceParameters BACnetConstructedData
}

var _ BACnetConfirmedServiceRequestConfirmedPrivateTransfer = (*_BACnetConfirmedServiceRequestConfirmedPrivateTransfer)(nil)
var _ BACnetConfirmedServiceRequestRequirements = (*_BACnetConfirmedServiceRequestConfirmedPrivateTransfer)(nil)

// NewBACnetConfirmedServiceRequestConfirmedPrivateTransfer factory function for _BACnetConfirmedServiceRequestConfirmedPrivateTransfer
func NewBACnetConfirmedServiceRequestConfirmedPrivateTransfer(vendorId BACnetVendorIdTagged, serviceNumber BACnetContextTagUnsignedInteger, serviceParameters BACnetConstructedData, serviceRequestLength uint32) *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer {
	if vendorId == nil {
		panic("vendorId of type BACnetVendorIdTagged for BACnetConfirmedServiceRequestConfirmedPrivateTransfer must not be nil")
	}
	if serviceNumber == nil {
		panic("serviceNumber of type BACnetContextTagUnsignedInteger for BACnetConfirmedServiceRequestConfirmedPrivateTransfer must not be nil")
	}
	_result := &_BACnetConfirmedServiceRequestConfirmedPrivateTransfer{
		BACnetConfirmedServiceRequestContract: NewBACnetConfirmedServiceRequest(serviceRequestLength),
		VendorId:                              vendorId,
		ServiceNumber:                         serviceNumber,
		ServiceParameters:                     serviceParameters,
	}
	_result.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder is a builder for BACnetConfirmedServiceRequestConfirmedPrivateTransfer
type BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(vendorId BACnetVendorIdTagged, serviceNumber BACnetContextTagUnsignedInteger) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
	// WithVendorId adds VendorId (property field)
	WithVendorId(BACnetVendorIdTagged) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
	// WithVendorIdBuilder adds VendorId (property field) which is build by the builder
	WithVendorIdBuilder(func(BACnetVendorIdTaggedBuilder) BACnetVendorIdTaggedBuilder) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
	// WithServiceNumber adds ServiceNumber (property field)
	WithServiceNumber(BACnetContextTagUnsignedInteger) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
	// WithServiceNumberBuilder adds ServiceNumber (property field) which is build by the builder
	WithServiceNumberBuilder(func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
	// WithServiceParameters adds ServiceParameters (property field)
	WithOptionalServiceParameters(BACnetConstructedData) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
	// WithOptionalServiceParametersBuilder adds ServiceParameters (property field) which is build by the builder
	WithOptionalServiceParametersBuilder(func(BACnetConstructedDataBuilder) BACnetConstructedDataBuilder) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
	// Build builds the BACnetConfirmedServiceRequestConfirmedPrivateTransfer or returns an error if something is wrong
	Build() (BACnetConfirmedServiceRequestConfirmedPrivateTransfer, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConfirmedServiceRequestConfirmedPrivateTransfer
}

// NewBACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder() creates a BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
func NewBACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder() BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder {
	return &_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder{_BACnetConfirmedServiceRequestConfirmedPrivateTransfer: new(_BACnetConfirmedServiceRequestConfirmedPrivateTransfer)}
}

type _BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder struct {
	*_BACnetConfirmedServiceRequestConfirmedPrivateTransfer

	parentBuilder *_BACnetConfirmedServiceRequestBuilder

	err *utils.MultiError
}

var _ (BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) = (*_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder)(nil)

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) setParent(contract BACnetConfirmedServiceRequestContract) {
	b.BACnetConfirmedServiceRequestContract = contract
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) WithMandatoryFields(vendorId BACnetVendorIdTagged, serviceNumber BACnetContextTagUnsignedInteger) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder {
	return b.WithVendorId(vendorId).WithServiceNumber(serviceNumber)
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) WithVendorId(vendorId BACnetVendorIdTagged) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder {
	b.VendorId = vendorId
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) WithVendorIdBuilder(builderSupplier func(BACnetVendorIdTaggedBuilder) BACnetVendorIdTaggedBuilder) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder {
	builder := builderSupplier(b.VendorId.CreateBACnetVendorIdTaggedBuilder())
	var err error
	b.VendorId, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetVendorIdTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) WithServiceNumber(serviceNumber BACnetContextTagUnsignedInteger) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder {
	b.ServiceNumber = serviceNumber
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) WithServiceNumberBuilder(builderSupplier func(BACnetContextTagUnsignedIntegerBuilder) BACnetContextTagUnsignedIntegerBuilder) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder {
	builder := builderSupplier(b.ServiceNumber.CreateBACnetContextTagUnsignedIntegerBuilder())
	var err error
	b.ServiceNumber, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) WithOptionalServiceParameters(serviceParameters BACnetConstructedData) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder {
	b.ServiceParameters = serviceParameters
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) WithOptionalServiceParametersBuilder(builderSupplier func(BACnetConstructedDataBuilder) BACnetConstructedDataBuilder) BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder {
	builder := builderSupplier(b.ServiceParameters.CreateBACnetConstructedDataBuilder())
	var err error
	b.ServiceParameters, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetConstructedDataBuilder failed"))
	}
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) Build() (BACnetConfirmedServiceRequestConfirmedPrivateTransfer, error) {
	if b.VendorId == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'vendorId' not set"))
	}
	if b.ServiceNumber == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'serviceNumber' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConfirmedServiceRequestConfirmedPrivateTransfer.deepCopy(), nil
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) MustBuild() BACnetConfirmedServiceRequestConfirmedPrivateTransfer {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) Done() BACnetConfirmedServiceRequestBuilder {
	return b.parentBuilder
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) buildForBACnetConfirmedServiceRequest() (BACnetConfirmedServiceRequest, error) {
	return b.Build()
}

func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder().(*_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder creates a BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder
func (b *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) CreateBACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder() BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder {
	if b == nil {
		return NewBACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder()
	}
	return &_BACnetConfirmedServiceRequestConfirmedPrivateTransferBuilder{_BACnetConfirmedServiceRequestConfirmedPrivateTransfer: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetServiceChoice() BACnetConfirmedServiceChoice {
	return BACnetConfirmedServiceChoice_CONFIRMED_PRIVATE_TRANSFER
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetParent() BACnetConfirmedServiceRequestContract {
	return m.BACnetConfirmedServiceRequestContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetVendorId() BACnetVendorIdTagged {
	return m.VendorId
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetServiceNumber() BACnetContextTagUnsignedInteger {
	return m.ServiceNumber
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetServiceParameters() BACnetConstructedData {
	return m.ServiceParameters
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConfirmedServiceRequestConfirmedPrivateTransfer(structType any) BACnetConfirmedServiceRequestConfirmedPrivateTransfer {
	if casted, ok := structType.(BACnetConfirmedServiceRequestConfirmedPrivateTransfer); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceRequestConfirmedPrivateTransfer); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetTypeName() string {
	return "BACnetConfirmedServiceRequestConfirmedPrivateTransfer"
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest).GetLengthInBits(ctx))

	// Simple field (vendorId)
	lengthInBits += m.VendorId.GetLengthInBits(ctx)

	// Simple field (serviceNumber)
	lengthInBits += m.ServiceNumber.GetLengthInBits(ctx)

	// Optional Field (serviceParameters)
	if m.ServiceParameters != nil {
		lengthInBits += m.ServiceParameters.GetLengthInBits(ctx)
	}

	return lengthInBits
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConfirmedServiceRequest, serviceRequestLength uint32) (__bACnetConfirmedServiceRequestConfirmedPrivateTransfer BACnetConfirmedServiceRequestConfirmedPrivateTransfer, err error) {
	m.BACnetConfirmedServiceRequestContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceRequestConfirmedPrivateTransfer"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConfirmedServiceRequestConfirmedPrivateTransfer")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	vendorId, err := ReadSimpleField[BACnetVendorIdTagged](ctx, "vendorId", ReadComplex[BACnetVendorIdTagged](BACnetVendorIdTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'vendorId' field"))
	}
	m.VendorId = vendorId

	serviceNumber, err := ReadSimpleField[BACnetContextTagUnsignedInteger](ctx, "serviceNumber", ReadComplex[BACnetContextTagUnsignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagUnsignedInteger]((uint8)(uint8(1)), (BACnetDataType)(BACnetDataType_UNSIGNED_INTEGER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'serviceNumber' field"))
	}
	m.ServiceNumber = serviceNumber

	var serviceParameters BACnetConstructedData
	_serviceParameters, err := ReadOptionalField[BACnetConstructedData](ctx, "serviceParameters", ReadComplex[BACnetConstructedData](BACnetConstructedDataParseWithBufferProducer[BACnetConstructedData]((uint8)(uint8(2)), (BACnetObjectType)(BACnetObjectType_VENDOR_PROPRIETARY_VALUE), (BACnetPropertyIdentifier)(BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE), (BACnetTagPayloadUnsignedInteger)(nil)), readBuffer), true)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'serviceParameters' field"))
	}
	if _serviceParameters != nil {
		serviceParameters = *_serviceParameters
		m.ServiceParameters = serviceParameters
	}

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceRequestConfirmedPrivateTransfer"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConfirmedServiceRequestConfirmedPrivateTransfer")
	}

	return m, nil
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceRequestConfirmedPrivateTransfer"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConfirmedServiceRequestConfirmedPrivateTransfer")
		}

		if err := WriteSimpleField[BACnetVendorIdTagged](ctx, "vendorId", m.GetVendorId(), WriteComplex[BACnetVendorIdTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'vendorId' field")
		}

		if err := WriteSimpleField[BACnetContextTagUnsignedInteger](ctx, "serviceNumber", m.GetServiceNumber(), WriteComplex[BACnetContextTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'serviceNumber' field")
		}

		if err := WriteOptionalField[BACnetConstructedData](ctx, "serviceParameters", GetRef(m.GetServiceParameters()), WriteComplex[BACnetConstructedData](writeBuffer), true); err != nil {
			return errors.Wrap(err, "Error serializing 'serviceParameters' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConfirmedServiceRequestConfirmedPrivateTransfer"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConfirmedServiceRequestConfirmedPrivateTransfer")
		}
		return nil
	}
	return m.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) IsBACnetConfirmedServiceRequestConfirmedPrivateTransfer() {
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) deepCopy() *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer {
	if m == nil {
		return nil
	}
	_BACnetConfirmedServiceRequestConfirmedPrivateTransferCopy := &_BACnetConfirmedServiceRequestConfirmedPrivateTransfer{
		m.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest).deepCopy(),
		m.VendorId.DeepCopy().(BACnetVendorIdTagged),
		m.ServiceNumber.DeepCopy().(BACnetContextTagUnsignedInteger),
		m.ServiceParameters.DeepCopy().(BACnetConstructedData),
	}
	m.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest)._SubType = m
	return _BACnetConfirmedServiceRequestConfirmedPrivateTransferCopy
}

func (m *_BACnetConfirmedServiceRequestConfirmedPrivateTransfer) String() string {
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
