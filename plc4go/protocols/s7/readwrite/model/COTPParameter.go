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

// COTPParameter is the corresponding interface of COTPParameter
type COTPParameter interface {
	COTPParameterContract
	COTPParameterRequirements
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// IsCOTPParameter is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCOTPParameter()
}

// COTPParameterContract provides a set of functions which can be overwritten by a sub struct
type COTPParameterContract interface {
	// GetRest() returns a parser argument
	GetRest() uint8
	// IsCOTPParameter is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCOTPParameter()
}

// COTPParameterRequirements provides a set of functions which need to be implemented by a sub struct
type COTPParameterRequirements interface {
	GetLengthInBits(ctx context.Context) uint16
	GetLengthInBytes(ctx context.Context) uint16
	// GetParameterType returns ParameterType (discriminator field)
	GetParameterType() uint8
}

// _COTPParameter is the data-structure of this message
type _COTPParameter struct {
	_SubType COTPParameter

	// Arguments.
	Rest uint8
}

var _ COTPParameterContract = (*_COTPParameter)(nil)

// NewCOTPParameter factory function for _COTPParameter
func NewCOTPParameter(rest uint8) *_COTPParameter {
	return &_COTPParameter{Rest: rest}
}

// Deprecated: use the interface for direct cast
func CastCOTPParameter(structType any) COTPParameter {
	if casted, ok := structType.(COTPParameter); ok {
		return casted
	}
	if casted, ok := structType.(*COTPParameter); ok {
		return *casted
	}
	return nil
}

func (m *_COTPParameter) GetTypeName() string {
	return "COTPParameter"
}

func (m *_COTPParameter) getLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)
	// Discriminator Field (parameterType)
	lengthInBits += 8

	// Implicit Field (parameterLength)
	lengthInBits += 8

	return lengthInBits
}

func (m *_COTPParameter) GetLengthInBytes(ctx context.Context) uint16 {
	return m._SubType.GetLengthInBits(ctx) / 8
}

func COTPParameterParse[T COTPParameter](ctx context.Context, theBytes []byte, rest uint8) (T, error) {
	return COTPParameterParseWithBuffer[T](ctx, utils.NewReadBufferByteBased(theBytes), rest)
}

func COTPParameterParseWithBufferProducer[T COTPParameter](rest uint8) func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
		v, err := COTPParameterParseWithBuffer[T](ctx, readBuffer, rest)
		if err != nil {
			var zero T
			return zero, err
		}
		return v, err
	}
}

func COTPParameterParseWithBuffer[T COTPParameter](ctx context.Context, readBuffer utils.ReadBuffer, rest uint8) (T, error) {
	v, err := (&_COTPParameter{Rest: rest}).parse(ctx, readBuffer, rest)
	if err != nil {
		var zero T
		return zero, err
	}
	return v.(T), err
}

func (m *_COTPParameter) parse(ctx context.Context, readBuffer utils.ReadBuffer, rest uint8) (__cOTPParameter COTPParameter, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("COTPParameter"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for COTPParameter")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	parameterType, err := ReadDiscriminatorField[uint8](ctx, "parameterType", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'parameterType' field"))
	}

	parameterLength, err := ReadImplicitField[uint8](ctx, "parameterLength", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'parameterLength' field"))
	}
	_ = parameterLength

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var _child COTPParameter
	switch {
	case parameterType == 0xC0: // COTPParameterTpduSize
		if _child, err = (&_COTPParameterTpduSize{}).parse(ctx, readBuffer, m, rest); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type COTPParameterTpduSize for type-switch of COTPParameter")
		}
	case parameterType == 0xC1: // COTPParameterCallingTsap
		if _child, err = (&_COTPParameterCallingTsap{}).parse(ctx, readBuffer, m, rest); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type COTPParameterCallingTsap for type-switch of COTPParameter")
		}
	case parameterType == 0xC2: // COTPParameterCalledTsap
		if _child, err = (&_COTPParameterCalledTsap{}).parse(ctx, readBuffer, m, rest); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type COTPParameterCalledTsap for type-switch of COTPParameter")
		}
	case parameterType == 0xC3: // COTPParameterChecksum
		if _child, err = (&_COTPParameterChecksum{}).parse(ctx, readBuffer, m, rest); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type COTPParameterChecksum for type-switch of COTPParameter")
		}
	case parameterType == 0xE0: // COTPParameterDisconnectAdditionalInformation
		if _child, err = (&_COTPParameterDisconnectAdditionalInformation{}).parse(ctx, readBuffer, m, rest); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type COTPParameterDisconnectAdditionalInformation for type-switch of COTPParameter")
		}
	default:
		return nil, errors.Errorf("Unmapped type for parameters [parameterType=%v]", parameterType)
	}

	if closeErr := readBuffer.CloseContext("COTPParameter"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for COTPParameter")
	}

	return _child, nil
}

func (pm *_COTPParameter) serializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child COTPParameter, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("COTPParameter"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for COTPParameter")
	}

	if err := WriteDiscriminatorField(ctx, "parameterType", m.GetParameterType(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'parameterType' field")
	}
	parameterLength := uint8(uint8(uint8(m.GetLengthInBytes(ctx))) - uint8(uint8(2)))
	if err := WriteImplicitField(ctx, "parameterLength", parameterLength, WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'parameterLength' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("COTPParameter"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for COTPParameter")
	}
	return nil
}

////
// Arguments Getter

func (m *_COTPParameter) GetRest() uint8 {
	return m.Rest
}

//
////

func (m *_COTPParameter) IsCOTPParameter() {}
