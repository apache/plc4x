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

// Constant values.
const CBusCommandDeviceManagement_DELIMITER byte = 0x0

// CBusCommandDeviceManagement is the corresponding interface of CBusCommandDeviceManagement
type CBusCommandDeviceManagement interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	CBusCommand
	// GetParamNo returns ParamNo (property field)
	GetParamNo() Parameter
	// GetParameterValue returns ParameterValue (property field)
	GetParameterValue() byte
}

// CBusCommandDeviceManagementExactly can be used when we want exactly this type and not a type which fulfills CBusCommandDeviceManagement.
// This is useful for switch cases.
type CBusCommandDeviceManagementExactly interface {
	CBusCommandDeviceManagement
	isCBusCommandDeviceManagement() bool
}

// _CBusCommandDeviceManagement is the data-structure of this message
type _CBusCommandDeviceManagement struct {
	*_CBusCommand
	ParamNo        Parameter
	ParameterValue byte
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_CBusCommandDeviceManagement) InitializeParent(parent CBusCommand, header CBusHeader) {
	m.Header = header
}

func (m *_CBusCommandDeviceManagement) GetParent() CBusCommand {
	return m._CBusCommand
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CBusCommandDeviceManagement) GetParamNo() Parameter {
	return m.ParamNo
}

func (m *_CBusCommandDeviceManagement) GetParameterValue() byte {
	return m.ParameterValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for const fields.
///////////////////////

func (m *_CBusCommandDeviceManagement) GetDelimiter() byte {
	return CBusCommandDeviceManagement_DELIMITER
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCBusCommandDeviceManagement factory function for _CBusCommandDeviceManagement
func NewCBusCommandDeviceManagement(paramNo Parameter, parameterValue byte, header CBusHeader, cBusOptions CBusOptions) *_CBusCommandDeviceManagement {
	_result := &_CBusCommandDeviceManagement{
		ParamNo:        paramNo,
		ParameterValue: parameterValue,
		_CBusCommand:   NewCBusCommand(header, cBusOptions),
	}
	_result._CBusCommand._CBusCommandChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastCBusCommandDeviceManagement(structType any) CBusCommandDeviceManagement {
	if casted, ok := structType.(CBusCommandDeviceManagement); ok {
		return casted
	}
	if casted, ok := structType.(*CBusCommandDeviceManagement); ok {
		return *casted
	}
	return nil
}

func (m *_CBusCommandDeviceManagement) GetTypeName() string {
	return "CBusCommandDeviceManagement"
}

func (m *_CBusCommandDeviceManagement) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (paramNo)
	lengthInBits += 8

	// Const Field (delimiter)
	lengthInBits += 8

	// Simple field (parameterValue)
	lengthInBits += 8

	return lengthInBits
}

func (m *_CBusCommandDeviceManagement) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func CBusCommandDeviceManagementParse(ctx context.Context, theBytes []byte, cBusOptions CBusOptions) (CBusCommandDeviceManagement, error) {
	return CBusCommandDeviceManagementParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), cBusOptions)
}

func CBusCommandDeviceManagementParseWithBufferProducer(cBusOptions CBusOptions) func(ctx context.Context, readBuffer utils.ReadBuffer) (CBusCommandDeviceManagement, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (CBusCommandDeviceManagement, error) {
		return CBusCommandDeviceManagementParseWithBuffer(ctx, readBuffer, cBusOptions)
	}
}

func CBusCommandDeviceManagementParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, cBusOptions CBusOptions) (CBusCommandDeviceManagement, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CBusCommandDeviceManagement"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CBusCommandDeviceManagement")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	paramNo, err := ReadEnumField[Parameter](ctx, "paramNo", "Parameter", ReadEnum(ParameterByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'paramNo' field"))
	}

	delimiter, err := ReadConstField[byte](ctx, "delimiter", ReadByte(readBuffer, 8), CBusCommandDeviceManagement_DELIMITER)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'delimiter' field"))
	}
	_ = delimiter

	parameterValue, err := ReadSimpleField(ctx, "parameterValue", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'parameterValue' field"))
	}

	if closeErr := readBuffer.CloseContext("CBusCommandDeviceManagement"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CBusCommandDeviceManagement")
	}

	// Create a partially initialized instance
	_child := &_CBusCommandDeviceManagement{
		_CBusCommand: &_CBusCommand{
			CBusOptions: cBusOptions,
		},
		ParamNo:        paramNo,
		ParameterValue: parameterValue,
	}
	_child._CBusCommand._CBusCommandChildRequirements = _child
	return _child, nil
}

func (m *_CBusCommandDeviceManagement) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CBusCommandDeviceManagement) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CBusCommandDeviceManagement"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CBusCommandDeviceManagement")
		}

		if err := WriteSimpleEnumField[Parameter](ctx, "paramNo", "Parameter", m.GetParamNo(), WriteEnum[Parameter, uint8](Parameter.GetValue, Parameter.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
			return errors.Wrap(err, "Error serializing 'paramNo' field")
		}

		if err := WriteConstField(ctx, "delimiter", CBusCommandDeviceManagement_DELIMITER, WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'delimiter' field")
		}

		if err := WriteSimpleField[byte](ctx, "parameterValue", m.GetParameterValue(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'parameterValue' field")
		}

		if popErr := writeBuffer.PopContext("CBusCommandDeviceManagement"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CBusCommandDeviceManagement")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CBusCommandDeviceManagement) isCBusCommandDeviceManagement() bool {
	return true
}

func (m *_CBusCommandDeviceManagement) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
