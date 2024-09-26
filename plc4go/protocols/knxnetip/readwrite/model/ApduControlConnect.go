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

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// ApduControlConnect is the corresponding interface of ApduControlConnect
type ApduControlConnect interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ApduControl
	// IsApduControlConnect is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsApduControlConnect()
	// CreateBuilder creates a ApduControlConnectBuilder
	CreateApduControlConnectBuilder() ApduControlConnectBuilder
}

// _ApduControlConnect is the data-structure of this message
type _ApduControlConnect struct {
	ApduControlContract
}

var _ ApduControlConnect = (*_ApduControlConnect)(nil)
var _ ApduControlRequirements = (*_ApduControlConnect)(nil)

// NewApduControlConnect factory function for _ApduControlConnect
func NewApduControlConnect() *_ApduControlConnect {
	_result := &_ApduControlConnect{
		ApduControlContract: NewApduControl(),
	}
	_result.ApduControlContract.(*_ApduControl)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ApduControlConnectBuilder is a builder for ApduControlConnect
type ApduControlConnectBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() ApduControlConnectBuilder
	// Build builds the ApduControlConnect or returns an error if something is wrong
	Build() (ApduControlConnect, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ApduControlConnect
}

// NewApduControlConnectBuilder() creates a ApduControlConnectBuilder
func NewApduControlConnectBuilder() ApduControlConnectBuilder {
	return &_ApduControlConnectBuilder{_ApduControlConnect: new(_ApduControlConnect)}
}

type _ApduControlConnectBuilder struct {
	*_ApduControlConnect

	parentBuilder *_ApduControlBuilder

	err *utils.MultiError
}

var _ (ApduControlConnectBuilder) = (*_ApduControlConnectBuilder)(nil)

func (b *_ApduControlConnectBuilder) setParent(contract ApduControlContract) {
	b.ApduControlContract = contract
}

func (b *_ApduControlConnectBuilder) WithMandatoryFields() ApduControlConnectBuilder {
	return b
}

func (b *_ApduControlConnectBuilder) Build() (ApduControlConnect, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ApduControlConnect.deepCopy(), nil
}

func (b *_ApduControlConnectBuilder) MustBuild() ApduControlConnect {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ApduControlConnectBuilder) Done() ApduControlBuilder {
	return b.parentBuilder
}

func (b *_ApduControlConnectBuilder) buildForApduControl() (ApduControl, error) {
	return b.Build()
}

func (b *_ApduControlConnectBuilder) DeepCopy() any {
	_copy := b.CreateApduControlConnectBuilder().(*_ApduControlConnectBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateApduControlConnectBuilder creates a ApduControlConnectBuilder
func (b *_ApduControlConnect) CreateApduControlConnectBuilder() ApduControlConnectBuilder {
	if b == nil {
		return NewApduControlConnectBuilder()
	}
	return &_ApduControlConnectBuilder{_ApduControlConnect: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ApduControlConnect) GetControlType() uint8 {
	return 0x0
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ApduControlConnect) GetParent() ApduControlContract {
	return m.ApduControlContract
}

// Deprecated: use the interface for direct cast
func CastApduControlConnect(structType any) ApduControlConnect {
	if casted, ok := structType.(ApduControlConnect); ok {
		return casted
	}
	if casted, ok := structType.(*ApduControlConnect); ok {
		return *casted
	}
	return nil
}

func (m *_ApduControlConnect) GetTypeName() string {
	return "ApduControlConnect"
}

func (m *_ApduControlConnect) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ApduControlContract.(*_ApduControl).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_ApduControlConnect) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ApduControlConnect) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ApduControl) (__apduControlConnect ApduControlConnect, err error) {
	m.ApduControlContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ApduControlConnect"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ApduControlConnect")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("ApduControlConnect"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ApduControlConnect")
	}

	return m, nil
}

func (m *_ApduControlConnect) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ApduControlConnect) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduControlConnect"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ApduControlConnect")
		}

		if popErr := writeBuffer.PopContext("ApduControlConnect"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ApduControlConnect")
		}
		return nil
	}
	return m.ApduControlContract.(*_ApduControl).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ApduControlConnect) IsApduControlConnect() {}

func (m *_ApduControlConnect) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ApduControlConnect) deepCopy() *_ApduControlConnect {
	if m == nil {
		return nil
	}
	_ApduControlConnectCopy := &_ApduControlConnect{
		m.ApduControlContract.(*_ApduControl).deepCopy(),
	}
	m.ApduControlContract.(*_ApduControl)._SubType = m
	return _ApduControlConnectCopy
}

func (m *_ApduControlConnect) String() string {
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
