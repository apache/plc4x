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

// NLMRequestMasterKey is the corresponding interface of NLMRequestMasterKey
type NLMRequestMasterKey interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	NLM
	// GetNumberOfSupportedKeyAlgorithms returns NumberOfSupportedKeyAlgorithms (property field)
	GetNumberOfSupportedKeyAlgorithms() uint8
	// GetEncryptionAndSignatureAlgorithms returns EncryptionAndSignatureAlgorithms (property field)
	GetEncryptionAndSignatureAlgorithms() []byte
	// IsNLMRequestMasterKey is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsNLMRequestMasterKey()
}

// _NLMRequestMasterKey is the data-structure of this message
type _NLMRequestMasterKey struct {
	NLMContract
	NumberOfSupportedKeyAlgorithms   uint8
	EncryptionAndSignatureAlgorithms []byte
}

var _ NLMRequestMasterKey = (*_NLMRequestMasterKey)(nil)
var _ NLMRequirements = (*_NLMRequestMasterKey)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_NLMRequestMasterKey) GetMessageType() uint8 {
	return 0x10
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_NLMRequestMasterKey) GetParent() NLMContract {
	return m.NLMContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_NLMRequestMasterKey) GetNumberOfSupportedKeyAlgorithms() uint8 {
	return m.NumberOfSupportedKeyAlgorithms
}

func (m *_NLMRequestMasterKey) GetEncryptionAndSignatureAlgorithms() []byte {
	return m.EncryptionAndSignatureAlgorithms
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewNLMRequestMasterKey factory function for _NLMRequestMasterKey
func NewNLMRequestMasterKey(numberOfSupportedKeyAlgorithms uint8, encryptionAndSignatureAlgorithms []byte, apduLength uint16) *_NLMRequestMasterKey {
	_result := &_NLMRequestMasterKey{
		NLMContract:                      NewNLM(apduLength),
		NumberOfSupportedKeyAlgorithms:   numberOfSupportedKeyAlgorithms,
		EncryptionAndSignatureAlgorithms: encryptionAndSignatureAlgorithms,
	}
	_result.NLMContract.(*_NLM)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastNLMRequestMasterKey(structType any) NLMRequestMasterKey {
	if casted, ok := structType.(NLMRequestMasterKey); ok {
		return casted
	}
	if casted, ok := structType.(*NLMRequestMasterKey); ok {
		return *casted
	}
	return nil
}

func (m *_NLMRequestMasterKey) GetTypeName() string {
	return "NLMRequestMasterKey"
}

func (m *_NLMRequestMasterKey) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.NLMContract.(*_NLM).getLengthInBits(ctx))

	// Simple field (numberOfSupportedKeyAlgorithms)
	lengthInBits += 8

	// Array field
	if len(m.EncryptionAndSignatureAlgorithms) > 0 {
		lengthInBits += 8 * uint16(len(m.EncryptionAndSignatureAlgorithms))
	}

	return lengthInBits
}

func (m *_NLMRequestMasterKey) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_NLMRequestMasterKey) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_NLM, apduLength uint16) (__nLMRequestMasterKey NLMRequestMasterKey, err error) {
	m.NLMContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("NLMRequestMasterKey"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for NLMRequestMasterKey")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	numberOfSupportedKeyAlgorithms, err := ReadSimpleField(ctx, "numberOfSupportedKeyAlgorithms", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'numberOfSupportedKeyAlgorithms' field"))
	}
	m.NumberOfSupportedKeyAlgorithms = numberOfSupportedKeyAlgorithms

	encryptionAndSignatureAlgorithms, err := readBuffer.ReadByteArray("encryptionAndSignatureAlgorithms", int(int32(apduLength)-int32(int32(2))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'encryptionAndSignatureAlgorithms' field"))
	}
	m.EncryptionAndSignatureAlgorithms = encryptionAndSignatureAlgorithms

	if closeErr := readBuffer.CloseContext("NLMRequestMasterKey"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for NLMRequestMasterKey")
	}

	return m, nil
}

func (m *_NLMRequestMasterKey) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_NLMRequestMasterKey) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("NLMRequestMasterKey"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for NLMRequestMasterKey")
		}

		if err := WriteSimpleField[uint8](ctx, "numberOfSupportedKeyAlgorithms", m.GetNumberOfSupportedKeyAlgorithms(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'numberOfSupportedKeyAlgorithms' field")
		}

		if err := WriteByteArrayField(ctx, "encryptionAndSignatureAlgorithms", m.GetEncryptionAndSignatureAlgorithms(), WriteByteArray(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'encryptionAndSignatureAlgorithms' field")
		}

		if popErr := writeBuffer.PopContext("NLMRequestMasterKey"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for NLMRequestMasterKey")
		}
		return nil
	}
	return m.NLMContract.(*_NLM).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_NLMRequestMasterKey) IsNLMRequestMasterKey() {}

func (m *_NLMRequestMasterKey) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
