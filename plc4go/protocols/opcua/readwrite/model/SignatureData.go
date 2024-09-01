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

// SignatureData is the corresponding interface of SignatureData
type SignatureData interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetAlgorithm returns Algorithm (property field)
	GetAlgorithm() PascalString
	// GetSignature returns Signature (property field)
	GetSignature() PascalByteString
}

// SignatureDataExactly can be used when we want exactly this type and not a type which fulfills SignatureData.
// This is useful for switch cases.
type SignatureDataExactly interface {
	SignatureData
	isSignatureData() bool
}

// _SignatureData is the data-structure of this message
type _SignatureData struct {
	*_ExtensionObjectDefinition
	Algorithm PascalString
	Signature PascalByteString
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_SignatureData) GetIdentifier() string {
	return "458"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SignatureData) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_SignatureData) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_SignatureData) GetAlgorithm() PascalString {
	return m.Algorithm
}

func (m *_SignatureData) GetSignature() PascalByteString {
	return m.Signature
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewSignatureData factory function for _SignatureData
func NewSignatureData(algorithm PascalString, signature PascalByteString) *_SignatureData {
	_result := &_SignatureData{
		Algorithm:                  algorithm,
		Signature:                  signature,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastSignatureData(structType any) SignatureData {
	if casted, ok := structType.(SignatureData); ok {
		return casted
	}
	if casted, ok := structType.(*SignatureData); ok {
		return *casted
	}
	return nil
}

func (m *_SignatureData) GetTypeName() string {
	return "SignatureData"
}

func (m *_SignatureData) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (algorithm)
	lengthInBits += m.Algorithm.GetLengthInBits(ctx)

	// Simple field (signature)
	lengthInBits += m.Signature.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_SignatureData) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func SignatureDataParse(ctx context.Context, theBytes []byte, identifier string) (SignatureData, error) {
	return SignatureDataParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func SignatureDataParseWithBufferProducer(identifier string) func(ctx context.Context, readBuffer utils.ReadBuffer) (SignatureData, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (SignatureData, error) {
		return SignatureDataParseWithBuffer(ctx, readBuffer, identifier)
	}
}

func SignatureDataParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (SignatureData, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SignatureData"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SignatureData")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	algorithm, err := ReadSimpleField[PascalString](ctx, "algorithm", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'algorithm' field"))
	}

	signature, err := ReadSimpleField[PascalByteString](ctx, "signature", ReadComplex[PascalByteString](PascalByteStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'signature' field"))
	}

	if closeErr := readBuffer.CloseContext("SignatureData"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SignatureData")
	}

	// Create a partially initialized instance
	_child := &_SignatureData{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		Algorithm:                  algorithm,
		Signature:                  signature,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_SignatureData) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_SignatureData) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SignatureData"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SignatureData")
		}

		if err := WriteSimpleField[PascalString](ctx, "algorithm", m.GetAlgorithm(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'algorithm' field")
		}

		if err := WriteSimpleField[PascalByteString](ctx, "signature", m.GetSignature(), WriteComplex[PascalByteString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'signature' field")
		}

		if popErr := writeBuffer.PopContext("SignatureData"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SignatureData")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SignatureData) isSignatureData() bool {
	return true
}

func (m *_SignatureData) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
