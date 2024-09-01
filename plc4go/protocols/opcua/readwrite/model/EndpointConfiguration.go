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

// EndpointConfiguration is the corresponding interface of EndpointConfiguration
type EndpointConfiguration interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetOperationTimeout returns OperationTimeout (property field)
	GetOperationTimeout() int32
	// GetUseBinaryEncoding returns UseBinaryEncoding (property field)
	GetUseBinaryEncoding() bool
	// GetMaxStringLength returns MaxStringLength (property field)
	GetMaxStringLength() int32
	// GetMaxByteStringLength returns MaxByteStringLength (property field)
	GetMaxByteStringLength() int32
	// GetMaxArrayLength returns MaxArrayLength (property field)
	GetMaxArrayLength() int32
	// GetMaxMessageSize returns MaxMessageSize (property field)
	GetMaxMessageSize() int32
	// GetMaxBufferSize returns MaxBufferSize (property field)
	GetMaxBufferSize() int32
	// GetChannelLifetime returns ChannelLifetime (property field)
	GetChannelLifetime() int32
	// GetSecurityTokenLifetime returns SecurityTokenLifetime (property field)
	GetSecurityTokenLifetime() int32
}

// EndpointConfigurationExactly can be used when we want exactly this type and not a type which fulfills EndpointConfiguration.
// This is useful for switch cases.
type EndpointConfigurationExactly interface {
	EndpointConfiguration
	isEndpointConfiguration() bool
}

// _EndpointConfiguration is the data-structure of this message
type _EndpointConfiguration struct {
	*_ExtensionObjectDefinition
	OperationTimeout      int32
	UseBinaryEncoding     bool
	MaxStringLength       int32
	MaxByteStringLength   int32
	MaxArrayLength        int32
	MaxMessageSize        int32
	MaxBufferSize         int32
	ChannelLifetime       int32
	SecurityTokenLifetime int32
	// Reserved Fields
	reservedField0 *uint8
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_EndpointConfiguration) GetIdentifier() string {
	return "333"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_EndpointConfiguration) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_EndpointConfiguration) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_EndpointConfiguration) GetOperationTimeout() int32 {
	return m.OperationTimeout
}

func (m *_EndpointConfiguration) GetUseBinaryEncoding() bool {
	return m.UseBinaryEncoding
}

func (m *_EndpointConfiguration) GetMaxStringLength() int32 {
	return m.MaxStringLength
}

func (m *_EndpointConfiguration) GetMaxByteStringLength() int32 {
	return m.MaxByteStringLength
}

func (m *_EndpointConfiguration) GetMaxArrayLength() int32 {
	return m.MaxArrayLength
}

func (m *_EndpointConfiguration) GetMaxMessageSize() int32 {
	return m.MaxMessageSize
}

func (m *_EndpointConfiguration) GetMaxBufferSize() int32 {
	return m.MaxBufferSize
}

func (m *_EndpointConfiguration) GetChannelLifetime() int32 {
	return m.ChannelLifetime
}

func (m *_EndpointConfiguration) GetSecurityTokenLifetime() int32 {
	return m.SecurityTokenLifetime
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewEndpointConfiguration factory function for _EndpointConfiguration
func NewEndpointConfiguration(operationTimeout int32, useBinaryEncoding bool, maxStringLength int32, maxByteStringLength int32, maxArrayLength int32, maxMessageSize int32, maxBufferSize int32, channelLifetime int32, securityTokenLifetime int32) *_EndpointConfiguration {
	_result := &_EndpointConfiguration{
		OperationTimeout:           operationTimeout,
		UseBinaryEncoding:          useBinaryEncoding,
		MaxStringLength:            maxStringLength,
		MaxByteStringLength:        maxByteStringLength,
		MaxArrayLength:             maxArrayLength,
		MaxMessageSize:             maxMessageSize,
		MaxBufferSize:              maxBufferSize,
		ChannelLifetime:            channelLifetime,
		SecurityTokenLifetime:      securityTokenLifetime,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastEndpointConfiguration(structType any) EndpointConfiguration {
	if casted, ok := structType.(EndpointConfiguration); ok {
		return casted
	}
	if casted, ok := structType.(*EndpointConfiguration); ok {
		return *casted
	}
	return nil
}

func (m *_EndpointConfiguration) GetTypeName() string {
	return "EndpointConfiguration"
}

func (m *_EndpointConfiguration) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (operationTimeout)
	lengthInBits += 32

	// Reserved Field (reserved)
	lengthInBits += 7

	// Simple field (useBinaryEncoding)
	lengthInBits += 1

	// Simple field (maxStringLength)
	lengthInBits += 32

	// Simple field (maxByteStringLength)
	lengthInBits += 32

	// Simple field (maxArrayLength)
	lengthInBits += 32

	// Simple field (maxMessageSize)
	lengthInBits += 32

	// Simple field (maxBufferSize)
	lengthInBits += 32

	// Simple field (channelLifetime)
	lengthInBits += 32

	// Simple field (securityTokenLifetime)
	lengthInBits += 32

	return lengthInBits
}

func (m *_EndpointConfiguration) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func EndpointConfigurationParse(ctx context.Context, theBytes []byte, identifier string) (EndpointConfiguration, error) {
	return EndpointConfigurationParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func EndpointConfigurationParseWithBufferProducer(identifier string) func(ctx context.Context, readBuffer utils.ReadBuffer) (EndpointConfiguration, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (EndpointConfiguration, error) {
		return EndpointConfigurationParseWithBuffer(ctx, readBuffer, identifier)
	}
}

func EndpointConfigurationParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (EndpointConfiguration, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("EndpointConfiguration"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for EndpointConfiguration")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	operationTimeout, err := ReadSimpleField(ctx, "operationTimeout", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'operationTimeout' field"))
	}

	reservedField0, err := ReadReservedField(ctx, "reserved", ReadUnsignedByte(readBuffer, uint8(7)), uint8(0x00))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}

	useBinaryEncoding, err := ReadSimpleField(ctx, "useBinaryEncoding", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'useBinaryEncoding' field"))
	}

	maxStringLength, err := ReadSimpleField(ctx, "maxStringLength", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'maxStringLength' field"))
	}

	maxByteStringLength, err := ReadSimpleField(ctx, "maxByteStringLength", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'maxByteStringLength' field"))
	}

	maxArrayLength, err := ReadSimpleField(ctx, "maxArrayLength", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'maxArrayLength' field"))
	}

	maxMessageSize, err := ReadSimpleField(ctx, "maxMessageSize", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'maxMessageSize' field"))
	}

	maxBufferSize, err := ReadSimpleField(ctx, "maxBufferSize", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'maxBufferSize' field"))
	}

	channelLifetime, err := ReadSimpleField(ctx, "channelLifetime", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'channelLifetime' field"))
	}

	securityTokenLifetime, err := ReadSimpleField(ctx, "securityTokenLifetime", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'securityTokenLifetime' field"))
	}

	if closeErr := readBuffer.CloseContext("EndpointConfiguration"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for EndpointConfiguration")
	}

	// Create a partially initialized instance
	_child := &_EndpointConfiguration{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		OperationTimeout:           operationTimeout,
		UseBinaryEncoding:          useBinaryEncoding,
		MaxStringLength:            maxStringLength,
		MaxByteStringLength:        maxByteStringLength,
		MaxArrayLength:             maxArrayLength,
		MaxMessageSize:             maxMessageSize,
		MaxBufferSize:              maxBufferSize,
		ChannelLifetime:            channelLifetime,
		SecurityTokenLifetime:      securityTokenLifetime,
		reservedField0:             reservedField0,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_EndpointConfiguration) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_EndpointConfiguration) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("EndpointConfiguration"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for EndpointConfiguration")
		}

		if err := WriteSimpleField[int32](ctx, "operationTimeout", m.GetOperationTimeout(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'operationTimeout' field")
		}

		if err := WriteReservedField[uint8](ctx, "reserved", uint8(0x00), WriteUnsignedByte(writeBuffer, 7)); err != nil {
			return errors.Wrap(err, "Error serializing 'reserved' field number 1")
		}

		if err := WriteSimpleField[bool](ctx, "useBinaryEncoding", m.GetUseBinaryEncoding(), WriteBoolean(writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'useBinaryEncoding' field")
		}

		if err := WriteSimpleField[int32](ctx, "maxStringLength", m.GetMaxStringLength(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'maxStringLength' field")
		}

		if err := WriteSimpleField[int32](ctx, "maxByteStringLength", m.GetMaxByteStringLength(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'maxByteStringLength' field")
		}

		if err := WriteSimpleField[int32](ctx, "maxArrayLength", m.GetMaxArrayLength(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'maxArrayLength' field")
		}

		if err := WriteSimpleField[int32](ctx, "maxMessageSize", m.GetMaxMessageSize(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'maxMessageSize' field")
		}

		if err := WriteSimpleField[int32](ctx, "maxBufferSize", m.GetMaxBufferSize(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'maxBufferSize' field")
		}

		if err := WriteSimpleField[int32](ctx, "channelLifetime", m.GetChannelLifetime(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'channelLifetime' field")
		}

		if err := WriteSimpleField[int32](ctx, "securityTokenLifetime", m.GetSecurityTokenLifetime(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'securityTokenLifetime' field")
		}

		if popErr := writeBuffer.PopContext("EndpointConfiguration"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for EndpointConfiguration")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_EndpointConfiguration) isEndpointConfiguration() bool {
	return true
}

func (m *_EndpointConfiguration) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
