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
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetVMACEntry is the corresponding interface of BACnetVMACEntry
type BACnetVMACEntry interface {
	utils.LengthAware
	utils.Serializable
	// GetVirtualMacAddress returns VirtualMacAddress (property field)
	GetVirtualMacAddress() BACnetContextTagOctetString
	// GetNativeMacAddress returns NativeMacAddress (property field)
	GetNativeMacAddress() BACnetContextTagOctetString
}

// BACnetVMACEntryExactly can be used when we want exactly this type and not a type which fulfills BACnetVMACEntry.
// This is useful for switch cases.
type BACnetVMACEntryExactly interface {
	BACnetVMACEntry
	isBACnetVMACEntry() bool
}

// _BACnetVMACEntry is the data-structure of this message
type _BACnetVMACEntry struct {
	VirtualMacAddress BACnetContextTagOctetString
	NativeMacAddress  BACnetContextTagOctetString
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetVMACEntry) GetVirtualMacAddress() BACnetContextTagOctetString {
	return m.VirtualMacAddress
}

func (m *_BACnetVMACEntry) GetNativeMacAddress() BACnetContextTagOctetString {
	return m.NativeMacAddress
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetVMACEntry factory function for _BACnetVMACEntry
func NewBACnetVMACEntry(virtualMacAddress BACnetContextTagOctetString, nativeMacAddress BACnetContextTagOctetString) *_BACnetVMACEntry {
	return &_BACnetVMACEntry{VirtualMacAddress: virtualMacAddress, NativeMacAddress: nativeMacAddress}
}

// Deprecated: use the interface for direct cast
func CastBACnetVMACEntry(structType interface{}) BACnetVMACEntry {
	if casted, ok := structType.(BACnetVMACEntry); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetVMACEntry); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetVMACEntry) GetTypeName() string {
	return "BACnetVMACEntry"
}

func (m *_BACnetVMACEntry) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetVMACEntry) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Optional Field (virtualMacAddress)
	if m.VirtualMacAddress != nil {
		lengthInBits += m.VirtualMacAddress.GetLengthInBits()
	}

	// Optional Field (nativeMacAddress)
	if m.NativeMacAddress != nil {
		lengthInBits += m.NativeMacAddress.GetLengthInBits()
	}

	return lengthInBits
}

func (m *_BACnetVMACEntry) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetVMACEntryParse(readBuffer utils.ReadBuffer) (BACnetVMACEntry, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetVMACEntry"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetVMACEntry")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Optional Field (virtualMacAddress) (Can be skipped, if a given expression evaluates to false)
	var virtualMacAddress BACnetContextTagOctetString = nil
	{
		currentPos = positionAware.GetPos()
		if pullErr := readBuffer.PullContext("virtualMacAddress"); pullErr != nil {
			return nil, errors.Wrap(pullErr, "Error pulling for virtualMacAddress")
		}
		_val, _err := BACnetContextTagParse(readBuffer, uint8(0), BACnetDataType_OCTET_STRING)
		switch {
		case errors.Is(_err, utils.ParseAssertError{}) || errors.Is(_err, io.EOF):
			log.Debug().Err(_err).Msg("Resetting position because optional threw an error")
			readBuffer.Reset(currentPos)
		case _err != nil:
			return nil, errors.Wrap(_err, "Error parsing 'virtualMacAddress' field of BACnetVMACEntry")
		default:
			virtualMacAddress = _val.(BACnetContextTagOctetString)
			if closeErr := readBuffer.CloseContext("virtualMacAddress"); closeErr != nil {
				return nil, errors.Wrap(closeErr, "Error closing for virtualMacAddress")
			}
		}
	}

	// Optional Field (nativeMacAddress) (Can be skipped, if a given expression evaluates to false)
	var nativeMacAddress BACnetContextTagOctetString = nil
	{
		currentPos = positionAware.GetPos()
		if pullErr := readBuffer.PullContext("nativeMacAddress"); pullErr != nil {
			return nil, errors.Wrap(pullErr, "Error pulling for nativeMacAddress")
		}
		_val, _err := BACnetContextTagParse(readBuffer, uint8(1), BACnetDataType_OCTET_STRING)
		switch {
		case errors.Is(_err, utils.ParseAssertError{}) || errors.Is(_err, io.EOF):
			log.Debug().Err(_err).Msg("Resetting position because optional threw an error")
			readBuffer.Reset(currentPos)
		case _err != nil:
			return nil, errors.Wrap(_err, "Error parsing 'nativeMacAddress' field of BACnetVMACEntry")
		default:
			nativeMacAddress = _val.(BACnetContextTagOctetString)
			if closeErr := readBuffer.CloseContext("nativeMacAddress"); closeErr != nil {
				return nil, errors.Wrap(closeErr, "Error closing for nativeMacAddress")
			}
		}
	}

	if closeErr := readBuffer.CloseContext("BACnetVMACEntry"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetVMACEntry")
	}

	// Create the instance
	return &_BACnetVMACEntry{
		VirtualMacAddress: virtualMacAddress,
		NativeMacAddress:  nativeMacAddress,
	}, nil
}

func (m *_BACnetVMACEntry) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetVMACEntry"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetVMACEntry")
	}

	// Optional Field (virtualMacAddress) (Can be skipped, if the value is null)
	var virtualMacAddress BACnetContextTagOctetString = nil
	if m.GetVirtualMacAddress() != nil {
		if pushErr := writeBuffer.PushContext("virtualMacAddress"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for virtualMacAddress")
		}
		virtualMacAddress = m.GetVirtualMacAddress()
		_virtualMacAddressErr := writeBuffer.WriteSerializable(virtualMacAddress)
		if popErr := writeBuffer.PopContext("virtualMacAddress"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for virtualMacAddress")
		}
		if _virtualMacAddressErr != nil {
			return errors.Wrap(_virtualMacAddressErr, "Error serializing 'virtualMacAddress' field")
		}
	}

	// Optional Field (nativeMacAddress) (Can be skipped, if the value is null)
	var nativeMacAddress BACnetContextTagOctetString = nil
	if m.GetNativeMacAddress() != nil {
		if pushErr := writeBuffer.PushContext("nativeMacAddress"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for nativeMacAddress")
		}
		nativeMacAddress = m.GetNativeMacAddress()
		_nativeMacAddressErr := writeBuffer.WriteSerializable(nativeMacAddress)
		if popErr := writeBuffer.PopContext("nativeMacAddress"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for nativeMacAddress")
		}
		if _nativeMacAddressErr != nil {
			return errors.Wrap(_nativeMacAddressErr, "Error serializing 'nativeMacAddress' field")
		}
	}

	if popErr := writeBuffer.PopContext("BACnetVMACEntry"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetVMACEntry")
	}
	return nil
}

func (m *_BACnetVMACEntry) isBACnetVMACEntry() bool {
	return true
}

func (m *_BACnetVMACEntry) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
