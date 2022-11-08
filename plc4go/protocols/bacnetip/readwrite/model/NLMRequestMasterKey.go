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
	"encoding/binary"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// NLMRequestMasterKey is the corresponding interface of NLMRequestMasterKey
type NLMRequestMasterKey interface {
	utils.LengthAware
	utils.Serializable
	NLM
	// GetNumberOfSupportedKeyAlgorithms returns NumberOfSupportedKeyAlgorithms (property field)
	GetNumberOfSupportedKeyAlgorithms() uint8
	// GetEncryptionAndSignatureAlgorithms returns EncryptionAndSignatureAlgorithms (property field)
	GetEncryptionAndSignatureAlgorithms() []byte
}

// NLMRequestMasterKeyExactly can be used when we want exactly this type and not a type which fulfills NLMRequestMasterKey.
// This is useful for switch cases.
type NLMRequestMasterKeyExactly interface {
	NLMRequestMasterKey
	isNLMRequestMasterKey() bool
}

// _NLMRequestMasterKey is the data-structure of this message
type _NLMRequestMasterKey struct {
	*_NLM
	NumberOfSupportedKeyAlgorithms   uint8
	EncryptionAndSignatureAlgorithms []byte
}

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

func (m *_NLMRequestMasterKey) InitializeParent(parent NLM) {}

func (m *_NLMRequestMasterKey) GetParent() NLM {
	return m._NLM
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
		NumberOfSupportedKeyAlgorithms:   numberOfSupportedKeyAlgorithms,
		EncryptionAndSignatureAlgorithms: encryptionAndSignatureAlgorithms,
		_NLM:                             NewNLM(apduLength),
	}
	_result._NLM._NLMChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastNLMRequestMasterKey(structType interface{}) NLMRequestMasterKey {
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

func (m *_NLMRequestMasterKey) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_NLMRequestMasterKey) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (numberOfSupportedKeyAlgorithms)
	lengthInBits += 8

	// Array field
	if len(m.EncryptionAndSignatureAlgorithms) > 0 {
		lengthInBits += 8 * uint16(len(m.EncryptionAndSignatureAlgorithms))
	}

	return lengthInBits
}

func (m *_NLMRequestMasterKey) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func NLMRequestMasterKeyParse(theBytes []byte, apduLength uint16) (NLMRequestMasterKey, error) {
	return NLMRequestMasterKeyParseWithBuffer(utils.NewReadBufferByteBased(theBytes, utils.WithByteOrderForReadBufferByteBased(binary.BigEndian)), apduLength) // TODO: get endianness from mspec
}

func NLMRequestMasterKeyParseWithBuffer(readBuffer utils.ReadBuffer, apduLength uint16) (NLMRequestMasterKey, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("NLMRequestMasterKey"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for NLMRequestMasterKey")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (numberOfSupportedKeyAlgorithms)
	_numberOfSupportedKeyAlgorithms, _numberOfSupportedKeyAlgorithmsErr := readBuffer.ReadUint8("numberOfSupportedKeyAlgorithms", 8)
	if _numberOfSupportedKeyAlgorithmsErr != nil {
		return nil, errors.Wrap(_numberOfSupportedKeyAlgorithmsErr, "Error parsing 'numberOfSupportedKeyAlgorithms' field of NLMRequestMasterKey")
	}
	numberOfSupportedKeyAlgorithms := _numberOfSupportedKeyAlgorithms
	// Byte Array field (encryptionAndSignatureAlgorithms)
	numberOfBytesencryptionAndSignatureAlgorithms := int(uint16(apduLength) - uint16(uint16(2)))
	encryptionAndSignatureAlgorithms, _readArrayErr := readBuffer.ReadByteArray("encryptionAndSignatureAlgorithms", numberOfBytesencryptionAndSignatureAlgorithms)
	if _readArrayErr != nil {
		return nil, errors.Wrap(_readArrayErr, "Error parsing 'encryptionAndSignatureAlgorithms' field of NLMRequestMasterKey")
	}

	if closeErr := readBuffer.CloseContext("NLMRequestMasterKey"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for NLMRequestMasterKey")
	}

	// Create a partially initialized instance
	_child := &_NLMRequestMasterKey{
		_NLM: &_NLM{
			ApduLength: apduLength,
		},
		NumberOfSupportedKeyAlgorithms:   numberOfSupportedKeyAlgorithms,
		EncryptionAndSignatureAlgorithms: encryptionAndSignatureAlgorithms,
	}
	_child._NLM._NLMChildRequirements = _child
	return _child, nil
}

func (m *_NLMRequestMasterKey) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian), utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes()))) // TODO: get endianness from mspec
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_NLMRequestMasterKey) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("NLMRequestMasterKey"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for NLMRequestMasterKey")
		}

		// Simple Field (numberOfSupportedKeyAlgorithms)
		numberOfSupportedKeyAlgorithms := uint8(m.GetNumberOfSupportedKeyAlgorithms())
		_numberOfSupportedKeyAlgorithmsErr := writeBuffer.WriteUint8("numberOfSupportedKeyAlgorithms", 8, (numberOfSupportedKeyAlgorithms))
		if _numberOfSupportedKeyAlgorithmsErr != nil {
			return errors.Wrap(_numberOfSupportedKeyAlgorithmsErr, "Error serializing 'numberOfSupportedKeyAlgorithms' field")
		}

		// Array Field (encryptionAndSignatureAlgorithms)
		// Byte Array field (encryptionAndSignatureAlgorithms)
		if err := writeBuffer.WriteByteArray("encryptionAndSignatureAlgorithms", m.GetEncryptionAndSignatureAlgorithms()); err != nil {
			return errors.Wrap(err, "Error serializing 'encryptionAndSignatureAlgorithms' field")
		}

		if popErr := writeBuffer.PopContext("NLMRequestMasterKey"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for NLMRequestMasterKey")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_NLMRequestMasterKey) isNLMRequestMasterKey() bool {
	return true
}

func (m *_NLMRequestMasterKey) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
