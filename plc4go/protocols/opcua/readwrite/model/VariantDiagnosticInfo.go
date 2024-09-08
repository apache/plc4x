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

// VariantDiagnosticInfo is the corresponding interface of VariantDiagnosticInfo
type VariantDiagnosticInfo interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	Variant
	// GetArrayLength returns ArrayLength (property field)
	GetArrayLength() *int32
	// GetValue returns Value (property field)
	GetValue() []DiagnosticInfo
	// IsVariantDiagnosticInfo is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsVariantDiagnosticInfo()
}

// _VariantDiagnosticInfo is the data-structure of this message
type _VariantDiagnosticInfo struct {
	VariantContract
	ArrayLength *int32
	Value       []DiagnosticInfo
}

var _ VariantDiagnosticInfo = (*_VariantDiagnosticInfo)(nil)
var _ VariantRequirements = (*_VariantDiagnosticInfo)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_VariantDiagnosticInfo) GetVariantType() uint8 {
	return uint8(25)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_VariantDiagnosticInfo) GetParent() VariantContract {
	return m.VariantContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_VariantDiagnosticInfo) GetArrayLength() *int32 {
	return m.ArrayLength
}

func (m *_VariantDiagnosticInfo) GetValue() []DiagnosticInfo {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewVariantDiagnosticInfo factory function for _VariantDiagnosticInfo
func NewVariantDiagnosticInfo(arrayLength *int32, value []DiagnosticInfo, arrayLengthSpecified bool, arrayDimensionsSpecified bool, noOfArrayDimensions *int32, arrayDimensions []bool) *_VariantDiagnosticInfo {
	_result := &_VariantDiagnosticInfo{
		VariantContract: NewVariant(arrayLengthSpecified, arrayDimensionsSpecified, noOfArrayDimensions, arrayDimensions),
		ArrayLength:     arrayLength,
		Value:           value,
	}
	_result.VariantContract.(*_Variant)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastVariantDiagnosticInfo(structType any) VariantDiagnosticInfo {
	if casted, ok := structType.(VariantDiagnosticInfo); ok {
		return casted
	}
	if casted, ok := structType.(*VariantDiagnosticInfo); ok {
		return *casted
	}
	return nil
}

func (m *_VariantDiagnosticInfo) GetTypeName() string {
	return "VariantDiagnosticInfo"
}

func (m *_VariantDiagnosticInfo) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.VariantContract.(*_Variant).getLengthInBits(ctx))

	// Optional Field (arrayLength)
	if m.ArrayLength != nil {
		lengthInBits += 32
	}

	// Array field
	if len(m.Value) > 0 {
		for _curItem, element := range m.Value {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.Value), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_VariantDiagnosticInfo) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_VariantDiagnosticInfo) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_Variant, arrayLengthSpecified bool) (__variantDiagnosticInfo VariantDiagnosticInfo, err error) {
	m.VariantContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("VariantDiagnosticInfo"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for VariantDiagnosticInfo")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	var arrayLength *int32
	arrayLength, err = ReadOptionalField[int32](ctx, "arrayLength", ReadSignedInt(readBuffer, uint8(32)), arrayLengthSpecified)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'arrayLength' field"))
	}
	m.ArrayLength = arrayLength

	value, err := ReadCountArrayField[DiagnosticInfo](ctx, "value", ReadComplex[DiagnosticInfo](DiagnosticInfoParseWithBuffer, readBuffer), uint64(utils.InlineIf(bool((arrayLength) == (nil)), func() any { return int32(int32(1)) }, func() any { return int32((*arrayLength)) }).(int32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'value' field"))
	}
	m.Value = value

	if closeErr := readBuffer.CloseContext("VariantDiagnosticInfo"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for VariantDiagnosticInfo")
	}

	return m, nil
}

func (m *_VariantDiagnosticInfo) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_VariantDiagnosticInfo) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("VariantDiagnosticInfo"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for VariantDiagnosticInfo")
		}

		if err := WriteOptionalField[int32](ctx, "arrayLength", m.GetArrayLength(), WriteSignedInt(writeBuffer, 32), true); err != nil {
			return errors.Wrap(err, "Error serializing 'arrayLength' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "value", m.GetValue(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'value' field")
		}

		if popErr := writeBuffer.PopContext("VariantDiagnosticInfo"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for VariantDiagnosticInfo")
		}
		return nil
	}
	return m.VariantContract.(*_Variant).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_VariantDiagnosticInfo) IsVariantDiagnosticInfo() {}

func (m *_VariantDiagnosticInfo) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
