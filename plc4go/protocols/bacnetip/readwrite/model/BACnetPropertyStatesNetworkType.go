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
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetPropertyStatesNetworkType is the corresponding interface of BACnetPropertyStatesNetworkType
type BACnetPropertyStatesNetworkType interface {
	utils.LengthAware
	utils.Serializable
	BACnetPropertyStates
	// GetNetworkType returns NetworkType (property field)
	GetNetworkType() BACnetNetworkTypeTagged
}

// BACnetPropertyStatesNetworkTypeExactly can be used when we want exactly this type and not a type which fulfills BACnetPropertyStatesNetworkType.
// This is useful for switch cases.
type BACnetPropertyStatesNetworkTypeExactly interface {
	BACnetPropertyStatesNetworkType
	isBACnetPropertyStatesNetworkType() bool
}

// _BACnetPropertyStatesNetworkType is the data-structure of this message
type _BACnetPropertyStatesNetworkType struct {
	*_BACnetPropertyStates
	NetworkType BACnetNetworkTypeTagged
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetPropertyStatesNetworkType) InitializeParent(parent BACnetPropertyStates, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetPropertyStatesNetworkType) GetParent() BACnetPropertyStates {
	return m._BACnetPropertyStates
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPropertyStatesNetworkType) GetNetworkType() BACnetNetworkTypeTagged {
	return m.NetworkType
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetPropertyStatesNetworkType factory function for _BACnetPropertyStatesNetworkType
func NewBACnetPropertyStatesNetworkType(networkType BACnetNetworkTypeTagged, peekedTagHeader BACnetTagHeader) *_BACnetPropertyStatesNetworkType {
	_result := &_BACnetPropertyStatesNetworkType{
		NetworkType:           networkType,
		_BACnetPropertyStates: NewBACnetPropertyStates(peekedTagHeader),
	}
	_result._BACnetPropertyStates._BACnetPropertyStatesChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetPropertyStatesNetworkType(structType interface{}) BACnetPropertyStatesNetworkType {
	if casted, ok := structType.(BACnetPropertyStatesNetworkType); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPropertyStatesNetworkType); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPropertyStatesNetworkType) GetTypeName() string {
	return "BACnetPropertyStatesNetworkType"
}

func (m *_BACnetPropertyStatesNetworkType) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetPropertyStatesNetworkType) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (networkType)
	lengthInBits += m.NetworkType.GetLengthInBits()

	return lengthInBits
}

func (m *_BACnetPropertyStatesNetworkType) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetPropertyStatesNetworkTypeParse(readBuffer utils.ReadBuffer, peekedTagNumber uint8) (BACnetPropertyStatesNetworkType, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPropertyStatesNetworkType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPropertyStatesNetworkType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (networkType)
	if pullErr := readBuffer.PullContext("networkType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for networkType")
	}
	_networkType, _networkTypeErr := BACnetNetworkTypeTaggedParse(readBuffer, uint8(peekedTagNumber), TagClass(TagClass_CONTEXT_SPECIFIC_TAGS))
	if _networkTypeErr != nil {
		return nil, errors.Wrap(_networkTypeErr, "Error parsing 'networkType' field of BACnetPropertyStatesNetworkType")
	}
	networkType := _networkType.(BACnetNetworkTypeTagged)
	if closeErr := readBuffer.CloseContext("networkType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for networkType")
	}

	if closeErr := readBuffer.CloseContext("BACnetPropertyStatesNetworkType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPropertyStatesNetworkType")
	}

	// Create a partially initialized instance
	_child := &_BACnetPropertyStatesNetworkType{
		_BACnetPropertyStates: &_BACnetPropertyStates{},
		NetworkType:           networkType,
	}
	_child._BACnetPropertyStates._BACnetPropertyStatesChildRequirements = _child
	return _child, nil
}

func (m *_BACnetPropertyStatesNetworkType) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPropertyStatesNetworkType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPropertyStatesNetworkType")
		}

		// Simple Field (networkType)
		if pushErr := writeBuffer.PushContext("networkType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for networkType")
		}
		_networkTypeErr := writeBuffer.WriteSerializable(m.GetNetworkType())
		if popErr := writeBuffer.PopContext("networkType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for networkType")
		}
		if _networkTypeErr != nil {
			return errors.Wrap(_networkTypeErr, "Error serializing 'networkType' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPropertyStatesNetworkType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPropertyStatesNetworkType")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetPropertyStatesNetworkType) isBACnetPropertyStatesNetworkType() bool {
	return true
}

func (m *_BACnetPropertyStatesNetworkType) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
