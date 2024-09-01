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

// NLMSetMasterKey is the corresponding interface of NLMSetMasterKey
type NLMSetMasterKey interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	NLM
	// GetKey returns Key (property field)
	GetKey() NLMUpdateKeyUpdateKeyEntry
}

// NLMSetMasterKeyExactly can be used when we want exactly this type and not a type which fulfills NLMSetMasterKey.
// This is useful for switch cases.
type NLMSetMasterKeyExactly interface {
	NLMSetMasterKey
	isNLMSetMasterKey() bool
}

// _NLMSetMasterKey is the data-structure of this message
type _NLMSetMasterKey struct {
	*_NLM
	Key NLMUpdateKeyUpdateKeyEntry
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_NLMSetMasterKey) GetMessageType() uint8 {
	return 0x11
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_NLMSetMasterKey) InitializeParent(parent NLM) {}

func (m *_NLMSetMasterKey) GetParent() NLM {
	return m._NLM
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_NLMSetMasterKey) GetKey() NLMUpdateKeyUpdateKeyEntry {
	return m.Key
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewNLMSetMasterKey factory function for _NLMSetMasterKey
func NewNLMSetMasterKey(key NLMUpdateKeyUpdateKeyEntry, apduLength uint16) *_NLMSetMasterKey {
	_result := &_NLMSetMasterKey{
		Key:  key,
		_NLM: NewNLM(apduLength),
	}
	_result._NLM._NLMChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastNLMSetMasterKey(structType any) NLMSetMasterKey {
	if casted, ok := structType.(NLMSetMasterKey); ok {
		return casted
	}
	if casted, ok := structType.(*NLMSetMasterKey); ok {
		return *casted
	}
	return nil
}

func (m *_NLMSetMasterKey) GetTypeName() string {
	return "NLMSetMasterKey"
}

func (m *_NLMSetMasterKey) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (key)
	lengthInBits += m.Key.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_NLMSetMasterKey) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func NLMSetMasterKeyParse(ctx context.Context, theBytes []byte, apduLength uint16) (NLMSetMasterKey, error) {
	return NLMSetMasterKeyParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), apduLength)
}

func NLMSetMasterKeyParseWithBufferProducer(apduLength uint16) func(ctx context.Context, readBuffer utils.ReadBuffer) (NLMSetMasterKey, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (NLMSetMasterKey, error) {
		return NLMSetMasterKeyParseWithBuffer(ctx, readBuffer, apduLength)
	}
}

func NLMSetMasterKeyParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, apduLength uint16) (NLMSetMasterKey, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("NLMSetMasterKey"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for NLMSetMasterKey")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	key, err := ReadSimpleField[NLMUpdateKeyUpdateKeyEntry](ctx, "key", ReadComplex[NLMUpdateKeyUpdateKeyEntry](NLMUpdateKeyUpdateKeyEntryParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'key' field"))
	}

	if closeErr := readBuffer.CloseContext("NLMSetMasterKey"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for NLMSetMasterKey")
	}

	// Create a partially initialized instance
	_child := &_NLMSetMasterKey{
		_NLM: &_NLM{
			ApduLength: apduLength,
		},
		Key: key,
	}
	_child._NLM._NLMChildRequirements = _child
	return _child, nil
}

func (m *_NLMSetMasterKey) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_NLMSetMasterKey) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("NLMSetMasterKey"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for NLMSetMasterKey")
		}

		if err := WriteSimpleField[NLMUpdateKeyUpdateKeyEntry](ctx, "key", m.GetKey(), WriteComplex[NLMUpdateKeyUpdateKeyEntry](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'key' field")
		}

		if popErr := writeBuffer.PopContext("NLMSetMasterKey"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for NLMSetMasterKey")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_NLMSetMasterKey) isNLMSetMasterKey() bool {
	return true
}

func (m *_NLMSetMasterKey) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
