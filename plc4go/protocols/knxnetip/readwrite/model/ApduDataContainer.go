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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// ApduDataContainer is the corresponding interface of ApduDataContainer
type ApduDataContainer interface {
	utils.LengthAware
	utils.Serializable
	Apdu
	// GetDataApdu returns DataApdu (property field)
	GetDataApdu() ApduData
}

// ApduDataContainerExactly can be used when we want exactly this type and not a type which fulfills ApduDataContainer.
// This is useful for switch cases.
type ApduDataContainerExactly interface {
	ApduDataContainer
	isApduDataContainer() bool
}

// _ApduDataContainer is the data-structure of this message
type _ApduDataContainer struct {
	*_Apdu
	DataApdu ApduData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ApduDataContainer) GetControl() uint8 {
	return uint8(0)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ApduDataContainer) InitializeParent(parent Apdu, numbered bool, counter uint8) {
	m.Numbered = numbered
	m.Counter = counter
}

func (m *_ApduDataContainer) GetParent() Apdu {
	return m._Apdu
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ApduDataContainer) GetDataApdu() ApduData {
	return m.DataApdu
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewApduDataContainer factory function for _ApduDataContainer
func NewApduDataContainer(dataApdu ApduData, numbered bool, counter uint8, dataLength uint8) *_ApduDataContainer {
	_result := &_ApduDataContainer{
		DataApdu: dataApdu,
		_Apdu:    NewApdu(numbered, counter, dataLength),
	}
	_result._Apdu._ApduChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastApduDataContainer(structType interface{}) ApduDataContainer {
	if casted, ok := structType.(ApduDataContainer); ok {
		return casted
	}
	if casted, ok := structType.(*ApduDataContainer); ok {
		return *casted
	}
	return nil
}

func (m *_ApduDataContainer) GetTypeName() string {
	return "ApduDataContainer"
}

func (m *_ApduDataContainer) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_ApduDataContainer) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (dataApdu)
	lengthInBits += m.DataApdu.GetLengthInBits()

	return lengthInBits
}

func (m *_ApduDataContainer) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func ApduDataContainerParse(theBytes []byte, dataLength uint8) (ApduDataContainer, error) {
	return ApduDataContainerParseWithBuffer(utils.NewReadBufferByteBased(theBytes), dataLength)
}

func ApduDataContainerParseWithBuffer(readBuffer utils.ReadBuffer, dataLength uint8) (ApduDataContainer, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ApduDataContainer"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ApduDataContainer")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (dataApdu)
	if pullErr := readBuffer.PullContext("dataApdu"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for dataApdu")
	}
	_dataApdu, _dataApduErr := ApduDataParseWithBuffer(readBuffer, uint8(dataLength))
	if _dataApduErr != nil {
		return nil, errors.Wrap(_dataApduErr, "Error parsing 'dataApdu' field of ApduDataContainer")
	}
	dataApdu := _dataApdu.(ApduData)
	if closeErr := readBuffer.CloseContext("dataApdu"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for dataApdu")
	}

	if closeErr := readBuffer.CloseContext("ApduDataContainer"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ApduDataContainer")
	}

	// Create a partially initialized instance
	_child := &_ApduDataContainer{
		_Apdu: &_Apdu{
			DataLength: dataLength,
		},
		DataApdu: dataApdu,
	}
	_child._Apdu._ApduChildRequirements = _child
	return _child, nil
}

func (m *_ApduDataContainer) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ApduDataContainer) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduDataContainer"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ApduDataContainer")
		}

		// Simple Field (dataApdu)
		if pushErr := writeBuffer.PushContext("dataApdu"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for dataApdu")
		}
		_dataApduErr := writeBuffer.WriteSerializable(m.GetDataApdu())
		if popErr := writeBuffer.PopContext("dataApdu"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for dataApdu")
		}
		if _dataApduErr != nil {
			return errors.Wrap(_dataApduErr, "Error serializing 'dataApdu' field")
		}

		if popErr := writeBuffer.PopContext("ApduDataContainer"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ApduDataContainer")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_ApduDataContainer) isApduDataContainer() bool {
	return true
}

func (m *_ApduDataContainer) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
