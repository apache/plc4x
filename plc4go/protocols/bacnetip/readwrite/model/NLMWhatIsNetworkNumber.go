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

// NLMWhatIsNetworkNumber is the corresponding interface of NLMWhatIsNetworkNumber
type NLMWhatIsNetworkNumber interface {
	utils.LengthAware
	utils.Serializable
	NLM
}

// NLMWhatIsNetworkNumberExactly can be used when we want exactly this type and not a type which fulfills NLMWhatIsNetworkNumber.
// This is useful for switch cases.
type NLMWhatIsNetworkNumberExactly interface {
	NLMWhatIsNetworkNumber
	isNLMWhatIsNetworkNumber() bool
}

// _NLMWhatIsNetworkNumber is the data-structure of this message
type _NLMWhatIsNetworkNumber struct {
	*_NLM
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_NLMWhatIsNetworkNumber) GetMessageType() uint8 {
	return 0x12
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_NLMWhatIsNetworkNumber) InitializeParent(parent NLM) {}

func (m *_NLMWhatIsNetworkNumber) GetParent() NLM {
	return m._NLM
}

// NewNLMWhatIsNetworkNumber factory function for _NLMWhatIsNetworkNumber
func NewNLMWhatIsNetworkNumber(apduLength uint16) *_NLMWhatIsNetworkNumber {
	_result := &_NLMWhatIsNetworkNumber{
		_NLM: NewNLM(apduLength),
	}
	_result._NLM._NLMChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastNLMWhatIsNetworkNumber(structType interface{}) NLMWhatIsNetworkNumber {
	if casted, ok := structType.(NLMWhatIsNetworkNumber); ok {
		return casted
	}
	if casted, ok := structType.(*NLMWhatIsNetworkNumber); ok {
		return *casted
	}
	return nil
}

func (m *_NLMWhatIsNetworkNumber) GetTypeName() string {
	return "NLMWhatIsNetworkNumber"
}

func (m *_NLMWhatIsNetworkNumber) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_NLMWhatIsNetworkNumber) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	return lengthInBits
}

func (m *_NLMWhatIsNetworkNumber) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func NLMWhatIsNetworkNumberParse(theBytes []byte, apduLength uint16) (NLMWhatIsNetworkNumber, error) {
	return NLMWhatIsNetworkNumberParseWithBuffer(utils.NewReadBufferByteBased(theBytes, utils.WithByteOrderForReadBufferByteBased(binary.BigEndian)), apduLength) // TODO: get endianness from mspec
}

func NLMWhatIsNetworkNumberParseWithBuffer(readBuffer utils.ReadBuffer, apduLength uint16) (NLMWhatIsNetworkNumber, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("NLMWhatIsNetworkNumber"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for NLMWhatIsNetworkNumber")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("NLMWhatIsNetworkNumber"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for NLMWhatIsNetworkNumber")
	}

	// Create a partially initialized instance
	_child := &_NLMWhatIsNetworkNumber{
		_NLM: &_NLM{
			ApduLength: apduLength,
		},
	}
	_child._NLM._NLMChildRequirements = _child
	return _child, nil
}

func (m *_NLMWhatIsNetworkNumber) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian), utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes()))) // TODO: get endianness from mspec
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_NLMWhatIsNetworkNumber) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("NLMWhatIsNetworkNumber"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for NLMWhatIsNetworkNumber")
		}

		if popErr := writeBuffer.PopContext("NLMWhatIsNetworkNumber"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for NLMWhatIsNetworkNumber")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_NLMWhatIsNetworkNumber) isNLMWhatIsNetworkNumber() bool {
	return true
}

func (m *_NLMWhatIsNetworkNumber) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
