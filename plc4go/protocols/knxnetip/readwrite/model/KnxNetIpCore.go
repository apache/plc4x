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

// KnxNetIpCore is the corresponding interface of KnxNetIpCore
type KnxNetIpCore interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ServiceId
	// GetVersion returns Version (property field)
	GetVersion() uint8
	// IsKnxNetIpCore is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsKnxNetIpCore()
	// CreateBuilder creates a KnxNetIpCoreBuilder
	CreateKnxNetIpCoreBuilder() KnxNetIpCoreBuilder
}

// _KnxNetIpCore is the data-structure of this message
type _KnxNetIpCore struct {
	ServiceIdContract
	Version uint8
}

var _ KnxNetIpCore = (*_KnxNetIpCore)(nil)
var _ ServiceIdRequirements = (*_KnxNetIpCore)(nil)

// NewKnxNetIpCore factory function for _KnxNetIpCore
func NewKnxNetIpCore(version uint8) *_KnxNetIpCore {
	_result := &_KnxNetIpCore{
		ServiceIdContract: NewServiceId(),
		Version:           version,
	}
	_result.ServiceIdContract.(*_ServiceId)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// KnxNetIpCoreBuilder is a builder for KnxNetIpCore
type KnxNetIpCoreBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(version uint8) KnxNetIpCoreBuilder
	// WithVersion adds Version (property field)
	WithVersion(uint8) KnxNetIpCoreBuilder
	// Build builds the KnxNetIpCore or returns an error if something is wrong
	Build() (KnxNetIpCore, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() KnxNetIpCore
}

// NewKnxNetIpCoreBuilder() creates a KnxNetIpCoreBuilder
func NewKnxNetIpCoreBuilder() KnxNetIpCoreBuilder {
	return &_KnxNetIpCoreBuilder{_KnxNetIpCore: new(_KnxNetIpCore)}
}

type _KnxNetIpCoreBuilder struct {
	*_KnxNetIpCore

	parentBuilder *_ServiceIdBuilder

	err *utils.MultiError
}

var _ (KnxNetIpCoreBuilder) = (*_KnxNetIpCoreBuilder)(nil)

func (b *_KnxNetIpCoreBuilder) setParent(contract ServiceIdContract) {
	b.ServiceIdContract = contract
}

func (b *_KnxNetIpCoreBuilder) WithMandatoryFields(version uint8) KnxNetIpCoreBuilder {
	return b.WithVersion(version)
}

func (b *_KnxNetIpCoreBuilder) WithVersion(version uint8) KnxNetIpCoreBuilder {
	b.Version = version
	return b
}

func (b *_KnxNetIpCoreBuilder) Build() (KnxNetIpCore, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._KnxNetIpCore.deepCopy(), nil
}

func (b *_KnxNetIpCoreBuilder) MustBuild() KnxNetIpCore {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_KnxNetIpCoreBuilder) Done() ServiceIdBuilder {
	return b.parentBuilder
}

func (b *_KnxNetIpCoreBuilder) buildForServiceId() (ServiceId, error) {
	return b.Build()
}

func (b *_KnxNetIpCoreBuilder) DeepCopy() any {
	_copy := b.CreateKnxNetIpCoreBuilder().(*_KnxNetIpCoreBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateKnxNetIpCoreBuilder creates a KnxNetIpCoreBuilder
func (b *_KnxNetIpCore) CreateKnxNetIpCoreBuilder() KnxNetIpCoreBuilder {
	if b == nil {
		return NewKnxNetIpCoreBuilder()
	}
	return &_KnxNetIpCoreBuilder{_KnxNetIpCore: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_KnxNetIpCore) GetServiceType() uint8 {
	return 0x02
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_KnxNetIpCore) GetParent() ServiceIdContract {
	return m.ServiceIdContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_KnxNetIpCore) GetVersion() uint8 {
	return m.Version
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastKnxNetIpCore(structType any) KnxNetIpCore {
	if casted, ok := structType.(KnxNetIpCore); ok {
		return casted
	}
	if casted, ok := structType.(*KnxNetIpCore); ok {
		return *casted
	}
	return nil
}

func (m *_KnxNetIpCore) GetTypeName() string {
	return "KnxNetIpCore"
}

func (m *_KnxNetIpCore) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ServiceIdContract.(*_ServiceId).GetLengthInBits(ctx))

	// Simple field (version)
	lengthInBits += 8

	return lengthInBits
}

func (m *_KnxNetIpCore) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_KnxNetIpCore) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ServiceId) (__knxNetIpCore KnxNetIpCore, err error) {
	m.ServiceIdContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("KnxNetIpCore"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for KnxNetIpCore")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	version, err := ReadSimpleField(ctx, "version", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'version' field"))
	}
	m.Version = version

	if closeErr := readBuffer.CloseContext("KnxNetIpCore"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for KnxNetIpCore")
	}

	return m, nil
}

func (m *_KnxNetIpCore) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_KnxNetIpCore) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("KnxNetIpCore"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for KnxNetIpCore")
		}

		if err := WriteSimpleField[uint8](ctx, "version", m.GetVersion(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'version' field")
		}

		if popErr := writeBuffer.PopContext("KnxNetIpCore"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for KnxNetIpCore")
		}
		return nil
	}
	return m.ServiceIdContract.(*_ServiceId).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_KnxNetIpCore) IsKnxNetIpCore() {}

func (m *_KnxNetIpCore) DeepCopy() any {
	return m.deepCopy()
}

func (m *_KnxNetIpCore) deepCopy() *_KnxNetIpCore {
	if m == nil {
		return nil
	}
	_KnxNetIpCoreCopy := &_KnxNetIpCore{
		m.ServiceIdContract.(*_ServiceId).deepCopy(),
		m.Version,
	}
	m.ServiceIdContract.(*_ServiceId)._SubType = m
	return _KnxNetIpCoreCopy
}

func (m *_KnxNetIpCore) String() string {
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
