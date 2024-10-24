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
const ModbusPDUReadDeviceIdentificationResponse_MEITYPE uint8 = 0x0E

// ModbusPDUReadDeviceIdentificationResponse is the corresponding interface of ModbusPDUReadDeviceIdentificationResponse
type ModbusPDUReadDeviceIdentificationResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ModbusPDU
	// GetLevel returns Level (property field)
	GetLevel() ModbusDeviceInformationLevel
	// GetIndividualAccess returns IndividualAccess (property field)
	GetIndividualAccess() bool
	// GetConformityLevel returns ConformityLevel (property field)
	GetConformityLevel() ModbusDeviceInformationConformityLevel
	// GetMoreFollows returns MoreFollows (property field)
	GetMoreFollows() ModbusDeviceInformationMoreFollows
	// GetNextObjectId returns NextObjectId (property field)
	GetNextObjectId() uint8
	// GetObjects returns Objects (property field)
	GetObjects() []ModbusDeviceInformationObject
	// IsModbusPDUReadDeviceIdentificationResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsModbusPDUReadDeviceIdentificationResponse()
	// CreateBuilder creates a ModbusPDUReadDeviceIdentificationResponseBuilder
	CreateModbusPDUReadDeviceIdentificationResponseBuilder() ModbusPDUReadDeviceIdentificationResponseBuilder
}

// _ModbusPDUReadDeviceIdentificationResponse is the data-structure of this message
type _ModbusPDUReadDeviceIdentificationResponse struct {
	ModbusPDUContract
	Level            ModbusDeviceInformationLevel
	IndividualAccess bool
	ConformityLevel  ModbusDeviceInformationConformityLevel
	MoreFollows      ModbusDeviceInformationMoreFollows
	NextObjectId     uint8
	Objects          []ModbusDeviceInformationObject
}

var _ ModbusPDUReadDeviceIdentificationResponse = (*_ModbusPDUReadDeviceIdentificationResponse)(nil)
var _ ModbusPDURequirements = (*_ModbusPDUReadDeviceIdentificationResponse)(nil)

// NewModbusPDUReadDeviceIdentificationResponse factory function for _ModbusPDUReadDeviceIdentificationResponse
func NewModbusPDUReadDeviceIdentificationResponse(level ModbusDeviceInformationLevel, individualAccess bool, conformityLevel ModbusDeviceInformationConformityLevel, moreFollows ModbusDeviceInformationMoreFollows, nextObjectId uint8, objects []ModbusDeviceInformationObject) *_ModbusPDUReadDeviceIdentificationResponse {
	_result := &_ModbusPDUReadDeviceIdentificationResponse{
		ModbusPDUContract: NewModbusPDU(),
		Level:             level,
		IndividualAccess:  individualAccess,
		ConformityLevel:   conformityLevel,
		MoreFollows:       moreFollows,
		NextObjectId:      nextObjectId,
		Objects:           objects,
	}
	_result.ModbusPDUContract.(*_ModbusPDU)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ModbusPDUReadDeviceIdentificationResponseBuilder is a builder for ModbusPDUReadDeviceIdentificationResponse
type ModbusPDUReadDeviceIdentificationResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(level ModbusDeviceInformationLevel, individualAccess bool, conformityLevel ModbusDeviceInformationConformityLevel, moreFollows ModbusDeviceInformationMoreFollows, nextObjectId uint8, objects []ModbusDeviceInformationObject) ModbusPDUReadDeviceIdentificationResponseBuilder
	// WithLevel adds Level (property field)
	WithLevel(ModbusDeviceInformationLevel) ModbusPDUReadDeviceIdentificationResponseBuilder
	// WithIndividualAccess adds IndividualAccess (property field)
	WithIndividualAccess(bool) ModbusPDUReadDeviceIdentificationResponseBuilder
	// WithConformityLevel adds ConformityLevel (property field)
	WithConformityLevel(ModbusDeviceInformationConformityLevel) ModbusPDUReadDeviceIdentificationResponseBuilder
	// WithMoreFollows adds MoreFollows (property field)
	WithMoreFollows(ModbusDeviceInformationMoreFollows) ModbusPDUReadDeviceIdentificationResponseBuilder
	// WithNextObjectId adds NextObjectId (property field)
	WithNextObjectId(uint8) ModbusPDUReadDeviceIdentificationResponseBuilder
	// WithObjects adds Objects (property field)
	WithObjects(...ModbusDeviceInformationObject) ModbusPDUReadDeviceIdentificationResponseBuilder
	// Build builds the ModbusPDUReadDeviceIdentificationResponse or returns an error if something is wrong
	Build() (ModbusPDUReadDeviceIdentificationResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ModbusPDUReadDeviceIdentificationResponse
}

// NewModbusPDUReadDeviceIdentificationResponseBuilder() creates a ModbusPDUReadDeviceIdentificationResponseBuilder
func NewModbusPDUReadDeviceIdentificationResponseBuilder() ModbusPDUReadDeviceIdentificationResponseBuilder {
	return &_ModbusPDUReadDeviceIdentificationResponseBuilder{_ModbusPDUReadDeviceIdentificationResponse: new(_ModbusPDUReadDeviceIdentificationResponse)}
}

type _ModbusPDUReadDeviceIdentificationResponseBuilder struct {
	*_ModbusPDUReadDeviceIdentificationResponse

	parentBuilder *_ModbusPDUBuilder

	err *utils.MultiError
}

var _ (ModbusPDUReadDeviceIdentificationResponseBuilder) = (*_ModbusPDUReadDeviceIdentificationResponseBuilder)(nil)

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) setParent(contract ModbusPDUContract) {
	b.ModbusPDUContract = contract
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) WithMandatoryFields(level ModbusDeviceInformationLevel, individualAccess bool, conformityLevel ModbusDeviceInformationConformityLevel, moreFollows ModbusDeviceInformationMoreFollows, nextObjectId uint8, objects []ModbusDeviceInformationObject) ModbusPDUReadDeviceIdentificationResponseBuilder {
	return b.WithLevel(level).WithIndividualAccess(individualAccess).WithConformityLevel(conformityLevel).WithMoreFollows(moreFollows).WithNextObjectId(nextObjectId).WithObjects(objects...)
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) WithLevel(level ModbusDeviceInformationLevel) ModbusPDUReadDeviceIdentificationResponseBuilder {
	b.Level = level
	return b
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) WithIndividualAccess(individualAccess bool) ModbusPDUReadDeviceIdentificationResponseBuilder {
	b.IndividualAccess = individualAccess
	return b
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) WithConformityLevel(conformityLevel ModbusDeviceInformationConformityLevel) ModbusPDUReadDeviceIdentificationResponseBuilder {
	b.ConformityLevel = conformityLevel
	return b
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) WithMoreFollows(moreFollows ModbusDeviceInformationMoreFollows) ModbusPDUReadDeviceIdentificationResponseBuilder {
	b.MoreFollows = moreFollows
	return b
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) WithNextObjectId(nextObjectId uint8) ModbusPDUReadDeviceIdentificationResponseBuilder {
	b.NextObjectId = nextObjectId
	return b
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) WithObjects(objects ...ModbusDeviceInformationObject) ModbusPDUReadDeviceIdentificationResponseBuilder {
	b.Objects = objects
	return b
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) Build() (ModbusPDUReadDeviceIdentificationResponse, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ModbusPDUReadDeviceIdentificationResponse.deepCopy(), nil
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) MustBuild() ModbusPDUReadDeviceIdentificationResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) Done() ModbusPDUBuilder {
	return b.parentBuilder
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) buildForModbusPDU() (ModbusPDU, error) {
	return b.Build()
}

func (b *_ModbusPDUReadDeviceIdentificationResponseBuilder) DeepCopy() any {
	_copy := b.CreateModbusPDUReadDeviceIdentificationResponseBuilder().(*_ModbusPDUReadDeviceIdentificationResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateModbusPDUReadDeviceIdentificationResponseBuilder creates a ModbusPDUReadDeviceIdentificationResponseBuilder
func (b *_ModbusPDUReadDeviceIdentificationResponse) CreateModbusPDUReadDeviceIdentificationResponseBuilder() ModbusPDUReadDeviceIdentificationResponseBuilder {
	if b == nil {
		return NewModbusPDUReadDeviceIdentificationResponseBuilder()
	}
	return &_ModbusPDUReadDeviceIdentificationResponseBuilder{_ModbusPDUReadDeviceIdentificationResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetErrorFlag() bool {
	return bool(false)
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetFunctionFlag() uint8 {
	return 0x2B
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetResponse() bool {
	return bool(true)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetParent() ModbusPDUContract {
	return m.ModbusPDUContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetLevel() ModbusDeviceInformationLevel {
	return m.Level
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetIndividualAccess() bool {
	return m.IndividualAccess
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetConformityLevel() ModbusDeviceInformationConformityLevel {
	return m.ConformityLevel
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetMoreFollows() ModbusDeviceInformationMoreFollows {
	return m.MoreFollows
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetNextObjectId() uint8 {
	return m.NextObjectId
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetObjects() []ModbusDeviceInformationObject {
	return m.Objects
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for const fields.
///////////////////////

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetMeiType() uint8 {
	return ModbusPDUReadDeviceIdentificationResponse_MEITYPE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastModbusPDUReadDeviceIdentificationResponse(structType any) ModbusPDUReadDeviceIdentificationResponse {
	if casted, ok := structType.(ModbusPDUReadDeviceIdentificationResponse); ok {
		return casted
	}
	if casted, ok := structType.(*ModbusPDUReadDeviceIdentificationResponse); ok {
		return *casted
	}
	return nil
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetTypeName() string {
	return "ModbusPDUReadDeviceIdentificationResponse"
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ModbusPDUContract.(*_ModbusPDU).GetLengthInBits(ctx))

	// Const Field (meiType)
	lengthInBits += 8

	// Simple field (level)
	lengthInBits += 8

	// Simple field (individualAccess)
	lengthInBits += 1

	// Simple field (conformityLevel)
	lengthInBits += 7

	// Simple field (moreFollows)
	lengthInBits += 8

	// Simple field (nextObjectId)
	lengthInBits += 8

	// Implicit Field (numberOfObjects)
	lengthInBits += 8

	// Array field
	if len(m.Objects) > 0 {
		for _curItem, element := range m.Objects {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.Objects), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ModbusPDU, response bool) (__modbusPDUReadDeviceIdentificationResponse ModbusPDUReadDeviceIdentificationResponse, err error) {
	m.ModbusPDUContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ModbusPDUReadDeviceIdentificationResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ModbusPDUReadDeviceIdentificationResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	meiType, err := ReadConstField[uint8](ctx, "meiType", ReadUnsignedByte(readBuffer, uint8(8)), ModbusPDUReadDeviceIdentificationResponse_MEITYPE)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'meiType' field"))
	}
	_ = meiType

	level, err := ReadEnumField[ModbusDeviceInformationLevel](ctx, "level", "ModbusDeviceInformationLevel", ReadEnum(ModbusDeviceInformationLevelByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'level' field"))
	}
	m.Level = level

	individualAccess, err := ReadSimpleField(ctx, "individualAccess", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'individualAccess' field"))
	}
	m.IndividualAccess = individualAccess

	conformityLevel, err := ReadEnumField[ModbusDeviceInformationConformityLevel](ctx, "conformityLevel", "ModbusDeviceInformationConformityLevel", ReadEnum(ModbusDeviceInformationConformityLevelByValue, ReadUnsignedByte(readBuffer, uint8(7))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'conformityLevel' field"))
	}
	m.ConformityLevel = conformityLevel

	moreFollows, err := ReadEnumField[ModbusDeviceInformationMoreFollows](ctx, "moreFollows", "ModbusDeviceInformationMoreFollows", ReadEnum(ModbusDeviceInformationMoreFollowsByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'moreFollows' field"))
	}
	m.MoreFollows = moreFollows

	nextObjectId, err := ReadSimpleField(ctx, "nextObjectId", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'nextObjectId' field"))
	}
	m.NextObjectId = nextObjectId

	numberOfObjects, err := ReadImplicitField[uint8](ctx, "numberOfObjects", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'numberOfObjects' field"))
	}
	_ = numberOfObjects

	objects, err := ReadCountArrayField[ModbusDeviceInformationObject](ctx, "objects", ReadComplex[ModbusDeviceInformationObject](ModbusDeviceInformationObjectParseWithBuffer, readBuffer), uint64(numberOfObjects))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'objects' field"))
	}
	m.Objects = objects

	if closeErr := readBuffer.CloseContext("ModbusPDUReadDeviceIdentificationResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ModbusPDUReadDeviceIdentificationResponse")
	}

	return m, nil
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ModbusPDUReadDeviceIdentificationResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ModbusPDUReadDeviceIdentificationResponse")
		}

		if err := WriteConstField(ctx, "meiType", ModbusPDUReadDeviceIdentificationResponse_MEITYPE, WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'meiType' field")
		}

		if err := WriteSimpleEnumField[ModbusDeviceInformationLevel](ctx, "level", "ModbusDeviceInformationLevel", m.GetLevel(), WriteEnum[ModbusDeviceInformationLevel, uint8](ModbusDeviceInformationLevel.GetValue, ModbusDeviceInformationLevel.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
			return errors.Wrap(err, "Error serializing 'level' field")
		}

		if err := WriteSimpleField[bool](ctx, "individualAccess", m.GetIndividualAccess(), WriteBoolean(writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'individualAccess' field")
		}

		if err := WriteSimpleEnumField[ModbusDeviceInformationConformityLevel](ctx, "conformityLevel", "ModbusDeviceInformationConformityLevel", m.GetConformityLevel(), WriteEnum[ModbusDeviceInformationConformityLevel, uint8](ModbusDeviceInformationConformityLevel.GetValue, ModbusDeviceInformationConformityLevel.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 7))); err != nil {
			return errors.Wrap(err, "Error serializing 'conformityLevel' field")
		}

		if err := WriteSimpleEnumField[ModbusDeviceInformationMoreFollows](ctx, "moreFollows", "ModbusDeviceInformationMoreFollows", m.GetMoreFollows(), WriteEnum[ModbusDeviceInformationMoreFollows, uint8](ModbusDeviceInformationMoreFollows.GetValue, ModbusDeviceInformationMoreFollows.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
			return errors.Wrap(err, "Error serializing 'moreFollows' field")
		}

		if err := WriteSimpleField[uint8](ctx, "nextObjectId", m.GetNextObjectId(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'nextObjectId' field")
		}
		numberOfObjects := uint8(uint8(len(m.GetObjects())))
		if err := WriteImplicitField(ctx, "numberOfObjects", numberOfObjects, WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'numberOfObjects' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "objects", m.GetObjects(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'objects' field")
		}

		if popErr := writeBuffer.PopContext("ModbusPDUReadDeviceIdentificationResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ModbusPDUReadDeviceIdentificationResponse")
		}
		return nil
	}
	return m.ModbusPDUContract.(*_ModbusPDU).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) IsModbusPDUReadDeviceIdentificationResponse() {}

func (m *_ModbusPDUReadDeviceIdentificationResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) deepCopy() *_ModbusPDUReadDeviceIdentificationResponse {
	if m == nil {
		return nil
	}
	_ModbusPDUReadDeviceIdentificationResponseCopy := &_ModbusPDUReadDeviceIdentificationResponse{
		m.ModbusPDUContract.(*_ModbusPDU).deepCopy(),
		m.Level,
		m.IndividualAccess,
		m.ConformityLevel,
		m.MoreFollows,
		m.NextObjectId,
		utils.DeepCopySlice[ModbusDeviceInformationObject, ModbusDeviceInformationObject](m.Objects),
	}
	m.ModbusPDUContract.(*_ModbusPDU)._SubType = m
	return _ModbusPDUReadDeviceIdentificationResponseCopy
}

func (m *_ModbusPDUReadDeviceIdentificationResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	wb := utils.NewWriteBufferBoxBased(
		utils.WithWriteBufferBoxBasedMergeSingleBoxes(),
		utils.WithWriteBufferBoxBasedOmitEmptyBoxes(),
		utils.WithWriteBufferBoxBasedPrintPosLengthFooter(),
	)
	if err := wb.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return wb.GetBox().String()
}
