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

// BACnetConstructedDataNetworkAccessSecurityPolicies is the corresponding interface of BACnetConstructedDataNetworkAccessSecurityPolicies
type BACnetConstructedDataNetworkAccessSecurityPolicies interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetNumberOfDataElements returns NumberOfDataElements (property field)
	GetNumberOfDataElements() BACnetApplicationTagUnsignedInteger
	// GetNetworkAccessSecurityPolicies returns NetworkAccessSecurityPolicies (property field)
	GetNetworkAccessSecurityPolicies() []BACnetNetworkSecurityPolicy
	// GetZero returns Zero (virtual field)
	GetZero() uint64
	// IsBACnetConstructedDataNetworkAccessSecurityPolicies is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataNetworkAccessSecurityPolicies()
	// CreateBuilder creates a BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder
	CreateBACnetConstructedDataNetworkAccessSecurityPoliciesBuilder() BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder
}

// _BACnetConstructedDataNetworkAccessSecurityPolicies is the data-structure of this message
type _BACnetConstructedDataNetworkAccessSecurityPolicies struct {
	BACnetConstructedDataContract
	NumberOfDataElements          BACnetApplicationTagUnsignedInteger
	NetworkAccessSecurityPolicies []BACnetNetworkSecurityPolicy
}

var _ BACnetConstructedDataNetworkAccessSecurityPolicies = (*_BACnetConstructedDataNetworkAccessSecurityPolicies)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataNetworkAccessSecurityPolicies)(nil)

// NewBACnetConstructedDataNetworkAccessSecurityPolicies factory function for _BACnetConstructedDataNetworkAccessSecurityPolicies
func NewBACnetConstructedDataNetworkAccessSecurityPolicies(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, numberOfDataElements BACnetApplicationTagUnsignedInteger, networkAccessSecurityPolicies []BACnetNetworkSecurityPolicy, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataNetworkAccessSecurityPolicies {
	_result := &_BACnetConstructedDataNetworkAccessSecurityPolicies{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		NumberOfDataElements:          numberOfDataElements,
		NetworkAccessSecurityPolicies: networkAccessSecurityPolicies,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder is a builder for BACnetConstructedDataNetworkAccessSecurityPolicies
type BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(networkAccessSecurityPolicies []BACnetNetworkSecurityPolicy) BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder
	// WithNumberOfDataElements adds NumberOfDataElements (property field)
	WithOptionalNumberOfDataElements(BACnetApplicationTagUnsignedInteger) BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder
	// WithOptionalNumberOfDataElementsBuilder adds NumberOfDataElements (property field) which is build by the builder
	WithOptionalNumberOfDataElementsBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder
	// WithNetworkAccessSecurityPolicies adds NetworkAccessSecurityPolicies (property field)
	WithNetworkAccessSecurityPolicies(...BACnetNetworkSecurityPolicy) BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder
	// Build builds the BACnetConstructedDataNetworkAccessSecurityPolicies or returns an error if something is wrong
	Build() (BACnetConstructedDataNetworkAccessSecurityPolicies, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataNetworkAccessSecurityPolicies
}

// NewBACnetConstructedDataNetworkAccessSecurityPoliciesBuilder() creates a BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder
func NewBACnetConstructedDataNetworkAccessSecurityPoliciesBuilder() BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder {
	return &_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder{_BACnetConstructedDataNetworkAccessSecurityPolicies: new(_BACnetConstructedDataNetworkAccessSecurityPolicies)}
}

type _BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder struct {
	*_BACnetConstructedDataNetworkAccessSecurityPolicies

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) = (*_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder)(nil)

func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) WithMandatoryFields(networkAccessSecurityPolicies []BACnetNetworkSecurityPolicy) BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder {
	return b.WithNetworkAccessSecurityPolicies(networkAccessSecurityPolicies...)
}

func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) WithOptionalNumberOfDataElements(numberOfDataElements BACnetApplicationTagUnsignedInteger) BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder {
	b.NumberOfDataElements = numberOfDataElements
	return b
}

func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) WithOptionalNumberOfDataElementsBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder {
	builder := builderSupplier(b.NumberOfDataElements.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.NumberOfDataElements, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) WithNetworkAccessSecurityPolicies(networkAccessSecurityPolicies ...BACnetNetworkSecurityPolicy) BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder {
	b.NetworkAccessSecurityPolicies = networkAccessSecurityPolicies
	return b
}

func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) Build() (BACnetConstructedDataNetworkAccessSecurityPolicies, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataNetworkAccessSecurityPolicies.deepCopy(), nil
}

func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) MustBuild() BACnetConstructedDataNetworkAccessSecurityPolicies {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataNetworkAccessSecurityPoliciesBuilder().(*_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataNetworkAccessSecurityPoliciesBuilder creates a BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder
func (b *_BACnetConstructedDataNetworkAccessSecurityPolicies) CreateBACnetConstructedDataNetworkAccessSecurityPoliciesBuilder() BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder {
	if b == nil {
		return NewBACnetConstructedDataNetworkAccessSecurityPoliciesBuilder()
	}
	return &_BACnetConstructedDataNetworkAccessSecurityPoliciesBuilder{_BACnetConstructedDataNetworkAccessSecurityPolicies: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_NETWORK_ACCESS_SECURITY_POLICIES
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) GetNumberOfDataElements() BACnetApplicationTagUnsignedInteger {
	return m.NumberOfDataElements
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) GetNetworkAccessSecurityPolicies() []BACnetNetworkSecurityPolicy {
	return m.NetworkAccessSecurityPolicies
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) GetZero() uint64 {
	ctx := context.Background()
	_ = ctx
	numberOfDataElements := m.GetNumberOfDataElements()
	_ = numberOfDataElements
	return uint64(uint64(0))
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataNetworkAccessSecurityPolicies(structType any) BACnetConstructedDataNetworkAccessSecurityPolicies {
	if casted, ok := structType.(BACnetConstructedDataNetworkAccessSecurityPolicies); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataNetworkAccessSecurityPolicies); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) GetTypeName() string {
	return "BACnetConstructedDataNetworkAccessSecurityPolicies"
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// A virtual field doesn't have any in- or output.

	// Optional Field (numberOfDataElements)
	if m.NumberOfDataElements != nil {
		lengthInBits += m.NumberOfDataElements.GetLengthInBits(ctx)
	}

	// Array field
	if len(m.NetworkAccessSecurityPolicies) > 0 {
		for _, element := range m.NetworkAccessSecurityPolicies {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	return lengthInBits
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataNetworkAccessSecurityPolicies BACnetConstructedDataNetworkAccessSecurityPolicies, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataNetworkAccessSecurityPolicies"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataNetworkAccessSecurityPolicies")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	zero, err := ReadVirtualField[uint64](ctx, "zero", (*uint64)(nil), uint64(0))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'zero' field"))
	}
	_ = zero

	var numberOfDataElements BACnetApplicationTagUnsignedInteger
	_numberOfDataElements, err := ReadOptionalField[BACnetApplicationTagUnsignedInteger](ctx, "numberOfDataElements", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer), bool(bool((arrayIndexArgument) != (nil))) && bool(bool((arrayIndexArgument.GetActualValue()) == (zero))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'numberOfDataElements' field"))
	}
	if _numberOfDataElements != nil {
		numberOfDataElements = *_numberOfDataElements
		m.NumberOfDataElements = numberOfDataElements
	}

	networkAccessSecurityPolicies, err := ReadTerminatedArrayField[BACnetNetworkSecurityPolicy](ctx, "networkAccessSecurityPolicies", ReadComplex[BACnetNetworkSecurityPolicy](BACnetNetworkSecurityPolicyParseWithBuffer, readBuffer), IsBACnetConstructedDataClosingTag(ctx, readBuffer, false, tagNumber))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'networkAccessSecurityPolicies' field"))
	}
	m.NetworkAccessSecurityPolicies = networkAccessSecurityPolicies

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataNetworkAccessSecurityPolicies"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataNetworkAccessSecurityPolicies")
	}

	return m, nil
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataNetworkAccessSecurityPolicies"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataNetworkAccessSecurityPolicies")
		}
		// Virtual field
		zero := m.GetZero()
		_ = zero
		if _zeroErr := writeBuffer.WriteVirtual(ctx, "zero", m.GetZero()); _zeroErr != nil {
			return errors.Wrap(_zeroErr, "Error serializing 'zero' field")
		}

		if err := WriteOptionalField[BACnetApplicationTagUnsignedInteger](ctx, "numberOfDataElements", GetRef(m.GetNumberOfDataElements()), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer), true); err != nil {
			return errors.Wrap(err, "Error serializing 'numberOfDataElements' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "networkAccessSecurityPolicies", m.GetNetworkAccessSecurityPolicies(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'networkAccessSecurityPolicies' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataNetworkAccessSecurityPolicies"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataNetworkAccessSecurityPolicies")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) IsBACnetConstructedDataNetworkAccessSecurityPolicies() {
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) deepCopy() *_BACnetConstructedDataNetworkAccessSecurityPolicies {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataNetworkAccessSecurityPoliciesCopy := &_BACnetConstructedDataNetworkAccessSecurityPolicies{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.NumberOfDataElements.DeepCopy().(BACnetApplicationTagUnsignedInteger),
		utils.DeepCopySlice[BACnetNetworkSecurityPolicy, BACnetNetworkSecurityPolicy](m.NetworkAccessSecurityPolicies),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataNetworkAccessSecurityPoliciesCopy
}

func (m *_BACnetConstructedDataNetworkAccessSecurityPolicies) String() string {
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
