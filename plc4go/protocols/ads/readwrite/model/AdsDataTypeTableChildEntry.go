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
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// Constant values.
const AdsDataTypeTableChildEntry_PROPERTYNAMETERMINATOR uint8 = 0x00
const AdsDataTypeTableChildEntry_DATATYPENAMETERMINATOR uint8 = 0x00
const AdsDataTypeTableChildEntry_COMMENTTERMINATOR uint8 = 0x00

// AdsDataTypeTableChildEntry is the corresponding interface of AdsDataTypeTableChildEntry
type AdsDataTypeTableChildEntry interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetEntryLength returns EntryLength (property field)
	GetEntryLength() uint32
	// GetVersion returns Version (property field)
	GetVersion() uint32
	// GetHashValue returns HashValue (property field)
	GetHashValue() uint32
	// GetTypeHashValue returns TypeHashValue (property field)
	GetTypeHashValue() uint32
	// GetSize returns Size (property field)
	GetSize() uint32
	// GetOffset returns Offset (property field)
	GetOffset() uint32
	// GetDataType returns DataType (property field)
	GetDataType() uint32
	// GetFlags returns Flags (property field)
	GetFlags() uint32
	// GetArrayDimensions returns ArrayDimensions (property field)
	GetArrayDimensions() uint16
	// GetNumChildren returns NumChildren (property field)
	GetNumChildren() uint16
	// GetPropertyName returns PropertyName (property field)
	GetPropertyName() string
	// GetDataTypeName returns DataTypeName (property field)
	GetDataTypeName() string
	// GetComment returns Comment (property field)
	GetComment() string
	// GetArrayInfo returns ArrayInfo (property field)
	GetArrayInfo() []AdsDataTypeArrayInfo
	// GetChildren returns Children (property field)
	GetChildren() []AdsDataTypeTableEntry
	// GetRest returns Rest (property field)
	GetRest() []byte
}

// AdsDataTypeTableChildEntryExactly can be used when we want exactly this type and not a type which fulfills AdsDataTypeTableChildEntry.
// This is useful for switch cases.
type AdsDataTypeTableChildEntryExactly interface {
	AdsDataTypeTableChildEntry
	isAdsDataTypeTableChildEntry() bool
}

// _AdsDataTypeTableChildEntry is the data-structure of this message
type _AdsDataTypeTableChildEntry struct {
	EntryLength     uint32
	Version         uint32
	HashValue       uint32
	TypeHashValue   uint32
	Size            uint32
	Offset          uint32
	DataType        uint32
	Flags           uint32
	ArrayDimensions uint16
	NumChildren     uint16
	PropertyName    string
	DataTypeName    string
	Comment         string
	ArrayInfo       []AdsDataTypeArrayInfo
	Children        []AdsDataTypeTableEntry
	Rest            []byte
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AdsDataTypeTableChildEntry) GetEntryLength() uint32 {
	return m.EntryLength
}

func (m *_AdsDataTypeTableChildEntry) GetVersion() uint32 {
	return m.Version
}

func (m *_AdsDataTypeTableChildEntry) GetHashValue() uint32 {
	return m.HashValue
}

func (m *_AdsDataTypeTableChildEntry) GetTypeHashValue() uint32 {
	return m.TypeHashValue
}

func (m *_AdsDataTypeTableChildEntry) GetSize() uint32 {
	return m.Size
}

func (m *_AdsDataTypeTableChildEntry) GetOffset() uint32 {
	return m.Offset
}

func (m *_AdsDataTypeTableChildEntry) GetDataType() uint32 {
	return m.DataType
}

func (m *_AdsDataTypeTableChildEntry) GetFlags() uint32 {
	return m.Flags
}

func (m *_AdsDataTypeTableChildEntry) GetArrayDimensions() uint16 {
	return m.ArrayDimensions
}

func (m *_AdsDataTypeTableChildEntry) GetNumChildren() uint16 {
	return m.NumChildren
}

func (m *_AdsDataTypeTableChildEntry) GetPropertyName() string {
	return m.PropertyName
}

func (m *_AdsDataTypeTableChildEntry) GetDataTypeName() string {
	return m.DataTypeName
}

func (m *_AdsDataTypeTableChildEntry) GetComment() string {
	return m.Comment
}

func (m *_AdsDataTypeTableChildEntry) GetArrayInfo() []AdsDataTypeArrayInfo {
	return m.ArrayInfo
}

func (m *_AdsDataTypeTableChildEntry) GetChildren() []AdsDataTypeTableEntry {
	return m.Children
}

func (m *_AdsDataTypeTableChildEntry) GetRest() []byte {
	return m.Rest
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for const fields.
///////////////////////

func (m *_AdsDataTypeTableChildEntry) GetPropertyNameTerminator() uint8 {
	return AdsDataTypeTableChildEntry_PROPERTYNAMETERMINATOR
}

func (m *_AdsDataTypeTableChildEntry) GetDataTypeNameTerminator() uint8 {
	return AdsDataTypeTableChildEntry_DATATYPENAMETERMINATOR
}

func (m *_AdsDataTypeTableChildEntry) GetCommentTerminator() uint8 {
	return AdsDataTypeTableChildEntry_COMMENTTERMINATOR
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewAdsDataTypeTableChildEntry factory function for _AdsDataTypeTableChildEntry
func NewAdsDataTypeTableChildEntry(entryLength uint32, version uint32, hashValue uint32, typeHashValue uint32, size uint32, offset uint32, dataType uint32, flags uint32, arrayDimensions uint16, numChildren uint16, propertyName string, dataTypeName string, comment string, arrayInfo []AdsDataTypeArrayInfo, children []AdsDataTypeTableEntry, rest []byte) *_AdsDataTypeTableChildEntry {
	return &_AdsDataTypeTableChildEntry{EntryLength: entryLength, Version: version, HashValue: hashValue, TypeHashValue: typeHashValue, Size: size, Offset: offset, DataType: dataType, Flags: flags, ArrayDimensions: arrayDimensions, NumChildren: numChildren, PropertyName: propertyName, DataTypeName: dataTypeName, Comment: comment, ArrayInfo: arrayInfo, Children: children, Rest: rest}
}

// Deprecated: use the interface for direct cast
func CastAdsDataTypeTableChildEntry(structType any) AdsDataTypeTableChildEntry {
	if casted, ok := structType.(AdsDataTypeTableChildEntry); ok {
		return casted
	}
	if casted, ok := structType.(*AdsDataTypeTableChildEntry); ok {
		return *casted
	}
	return nil
}

func (m *_AdsDataTypeTableChildEntry) GetTypeName() string {
	return "AdsDataTypeTableChildEntry"
}

func (m *_AdsDataTypeTableChildEntry) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (entryLength)
	lengthInBits += 32

	// Simple field (version)
	lengthInBits += 32

	// Simple field (hashValue)
	lengthInBits += 32

	// Simple field (typeHashValue)
	lengthInBits += 32

	// Simple field (size)
	lengthInBits += 32

	// Simple field (offset)
	lengthInBits += 32

	// Simple field (dataType)
	lengthInBits += 32

	// Simple field (flags)
	lengthInBits += 32

	// Implicit Field (propertyNameLength)
	lengthInBits += 16

	// Implicit Field (dataTypeNameLength)
	lengthInBits += 16

	// Implicit Field (commentLength)
	lengthInBits += 16

	// Simple field (arrayDimensions)
	lengthInBits += 16

	// Simple field (numChildren)
	lengthInBits += 16

	// Simple field (propertyName)
	lengthInBits += uint16(int32(uint16(len(m.GetPropertyName()))) * int32(int32(8)))

	// Const Field (propertyNameTerminator)
	lengthInBits += 8

	// Simple field (dataTypeName)
	lengthInBits += uint16(int32(uint16(len(m.GetDataTypeName()))) * int32(int32(8)))

	// Const Field (dataTypeNameTerminator)
	lengthInBits += 8

	// Simple field (comment)
	lengthInBits += uint16(int32(uint16(len(m.GetComment()))) * int32(int32(8)))

	// Const Field (commentTerminator)
	lengthInBits += 8

	// Array field
	if len(m.ArrayInfo) > 0 {
		for _curItem, element := range m.ArrayInfo {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.ArrayInfo), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Array field
	if len(m.Children) > 0 {
		for _curItem, element := range m.Children {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.Children), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Array field
	if len(m.Rest) > 0 {
		lengthInBits += 8 * uint16(len(m.Rest))
	}

	return lengthInBits
}

func (m *_AdsDataTypeTableChildEntry) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func AdsDataTypeTableChildEntryParse(ctx context.Context, theBytes []byte) (AdsDataTypeTableChildEntry, error) {
	return AdsDataTypeTableChildEntryParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian)))
}

func AdsDataTypeTableChildEntryParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (AdsDataTypeTableChildEntry, error) {
	positionAware := readBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pullErr := readBuffer.PullContext("AdsDataTypeTableChildEntry"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AdsDataTypeTableChildEntry")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos
	var startPos = positionAware.GetPos()
	_ = startPos

	// Simple Field (entryLength)
	_entryLength, _entryLengthErr := /*TODO: migrate me*/ readBuffer.ReadUint32("entryLength", 32)
	if _entryLengthErr != nil {
		return nil, errors.Wrap(_entryLengthErr, "Error parsing 'entryLength' field of AdsDataTypeTableChildEntry")
	}
	entryLength := _entryLength

	// Simple Field (version)
	_version, _versionErr := /*TODO: migrate me*/ readBuffer.ReadUint32("version", 32)
	if _versionErr != nil {
		return nil, errors.Wrap(_versionErr, "Error parsing 'version' field of AdsDataTypeTableChildEntry")
	}
	version := _version

	// Simple Field (hashValue)
	_hashValue, _hashValueErr := /*TODO: migrate me*/ readBuffer.ReadUint32("hashValue", 32)
	if _hashValueErr != nil {
		return nil, errors.Wrap(_hashValueErr, "Error parsing 'hashValue' field of AdsDataTypeTableChildEntry")
	}
	hashValue := _hashValue

	// Simple Field (typeHashValue)
	_typeHashValue, _typeHashValueErr := /*TODO: migrate me*/ readBuffer.ReadUint32("typeHashValue", 32)
	if _typeHashValueErr != nil {
		return nil, errors.Wrap(_typeHashValueErr, "Error parsing 'typeHashValue' field of AdsDataTypeTableChildEntry")
	}
	typeHashValue := _typeHashValue

	// Simple Field (size)
	_size, _sizeErr := /*TODO: migrate me*/ readBuffer.ReadUint32("size", 32)
	if _sizeErr != nil {
		return nil, errors.Wrap(_sizeErr, "Error parsing 'size' field of AdsDataTypeTableChildEntry")
	}
	size := _size

	// Simple Field (offset)
	_offset, _offsetErr := /*TODO: migrate me*/ readBuffer.ReadUint32("offset", 32)
	if _offsetErr != nil {
		return nil, errors.Wrap(_offsetErr, "Error parsing 'offset' field of AdsDataTypeTableChildEntry")
	}
	offset := _offset

	// Simple Field (dataType)
	_dataType, _dataTypeErr := /*TODO: migrate me*/ readBuffer.ReadUint32("dataType", 32)
	if _dataTypeErr != nil {
		return nil, errors.Wrap(_dataTypeErr, "Error parsing 'dataType' field of AdsDataTypeTableChildEntry")
	}
	dataType := _dataType

	// Simple Field (flags)
	_flags, _flagsErr := /*TODO: migrate me*/ readBuffer.ReadUint32("flags", 32)
	if _flagsErr != nil {
		return nil, errors.Wrap(_flagsErr, "Error parsing 'flags' field of AdsDataTypeTableChildEntry")
	}
	flags := _flags

	propertyNameLength, err := ReadImplicitField[uint16](ctx, "propertyNameLength", ReadUnsignedShort(readBuffer, 16), codegen.WithByteOrder(binary.LittleEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'propertyNameLength' field"))
	}

	dataTypeNameLength, err := ReadImplicitField[uint16](ctx, "dataTypeNameLength", ReadUnsignedShort(readBuffer, 16), codegen.WithByteOrder(binary.LittleEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'dataTypeNameLength' field"))
	}

	commentLength, err := ReadImplicitField[uint16](ctx, "commentLength", ReadUnsignedShort(readBuffer, 16), codegen.WithByteOrder(binary.LittleEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'commentLength' field"))
	}

	// Simple Field (arrayDimensions)
	_arrayDimensions, _arrayDimensionsErr := /*TODO: migrate me*/ readBuffer.ReadUint16("arrayDimensions", 16)
	if _arrayDimensionsErr != nil {
		return nil, errors.Wrap(_arrayDimensionsErr, "Error parsing 'arrayDimensions' field of AdsDataTypeTableChildEntry")
	}
	arrayDimensions := _arrayDimensions

	// Simple Field (numChildren)
	_numChildren, _numChildrenErr := /*TODO: migrate me*/ readBuffer.ReadUint16("numChildren", 16)
	if _numChildrenErr != nil {
		return nil, errors.Wrap(_numChildrenErr, "Error parsing 'numChildren' field of AdsDataTypeTableChildEntry")
	}
	numChildren := _numChildren

	// Simple Field (propertyName)
	_propertyName, _propertyNameErr := /*TODO: migrate me*/ readBuffer.ReadString("propertyName", uint32((propertyNameLength)*(8)), utils.WithEncoding("UTF-8"))
	if _propertyNameErr != nil {
		return nil, errors.Wrap(_propertyNameErr, "Error parsing 'propertyName' field of AdsDataTypeTableChildEntry")
	}
	propertyName := _propertyName

	propertyNameTerminator, err := ReadConstField[uint8](ctx, "propertyNameTerminator", ReadUnsignedByte(readBuffer, 8), AdsDataTypeTableChildEntry_PROPERTYNAMETERMINATOR, codegen.WithByteOrder(binary.LittleEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'propertyNameTerminator' field"))
	}
	_ = propertyNameTerminator

	// Simple Field (dataTypeName)
	_dataTypeName, _dataTypeNameErr := /*TODO: migrate me*/ readBuffer.ReadString("dataTypeName", uint32((dataTypeNameLength)*(8)), utils.WithEncoding("UTF-8"))
	if _dataTypeNameErr != nil {
		return nil, errors.Wrap(_dataTypeNameErr, "Error parsing 'dataTypeName' field of AdsDataTypeTableChildEntry")
	}
	dataTypeName := _dataTypeName

	dataTypeNameTerminator, err := ReadConstField[uint8](ctx, "dataTypeNameTerminator", ReadUnsignedByte(readBuffer, 8), AdsDataTypeTableChildEntry_DATATYPENAMETERMINATOR, codegen.WithByteOrder(binary.LittleEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'dataTypeNameTerminator' field"))
	}
	_ = dataTypeNameTerminator

	// Simple Field (comment)
	_comment, _commentErr := /*TODO: migrate me*/ readBuffer.ReadString("comment", uint32((commentLength)*(8)), utils.WithEncoding("UTF-8"))
	if _commentErr != nil {
		return nil, errors.Wrap(_commentErr, "Error parsing 'comment' field of AdsDataTypeTableChildEntry")
	}
	comment := _comment

	commentTerminator, err := ReadConstField[uint8](ctx, "commentTerminator", ReadUnsignedByte(readBuffer, 8), AdsDataTypeTableChildEntry_COMMENTTERMINATOR, codegen.WithByteOrder(binary.LittleEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'commentTerminator' field"))
	}
	_ = commentTerminator

	arrayInfo, err := ReadCountArrayField[AdsDataTypeArrayInfo](ctx, "arrayInfo", ReadComplex[AdsDataTypeArrayInfo](AdsDataTypeArrayInfoParseWithBuffer, readBuffer), uint64(arrayDimensions), codegen.WithByteOrder(binary.LittleEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'arrayInfo' field"))
	}

	children, err := ReadCountArrayField[AdsDataTypeTableEntry](ctx, "children", ReadComplex[AdsDataTypeTableEntry](AdsDataTypeTableEntryParseWithBuffer, readBuffer), uint64(numChildren), codegen.WithByteOrder(binary.LittleEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'children' field"))
	}

	rest, err := readBuffer.ReadByteArray("rest", int(int32(entryLength)-int32((positionAware.GetPos()-startPos))), codegen.WithByteOrder(binary.LittleEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'rest' field"))
	}

	if closeErr := readBuffer.CloseContext("AdsDataTypeTableChildEntry"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AdsDataTypeTableChildEntry")
	}

	// Create the instance
	return &_AdsDataTypeTableChildEntry{
		EntryLength:     entryLength,
		Version:         version,
		HashValue:       hashValue,
		TypeHashValue:   typeHashValue,
		Size:            size,
		Offset:          offset,
		DataType:        dataType,
		Flags:           flags,
		ArrayDimensions: arrayDimensions,
		NumChildren:     numChildren,
		PropertyName:    propertyName,
		DataTypeName:    dataTypeName,
		Comment:         comment,
		ArrayInfo:       arrayInfo,
		Children:        children,
		Rest:            rest,
	}, nil
}

func (m *_AdsDataTypeTableChildEntry) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))), utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_AdsDataTypeTableChildEntry) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("AdsDataTypeTableChildEntry"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for AdsDataTypeTableChildEntry")
	}

	// Simple Field (entryLength)
	entryLength := uint32(m.GetEntryLength())
	_entryLengthErr := /*TODO: migrate me*/ writeBuffer.WriteUint32("entryLength", 32, uint32((entryLength)))
	if _entryLengthErr != nil {
		return errors.Wrap(_entryLengthErr, "Error serializing 'entryLength' field")
	}

	// Simple Field (version)
	version := uint32(m.GetVersion())
	_versionErr := /*TODO: migrate me*/ writeBuffer.WriteUint32("version", 32, uint32((version)))
	if _versionErr != nil {
		return errors.Wrap(_versionErr, "Error serializing 'version' field")
	}

	// Simple Field (hashValue)
	hashValue := uint32(m.GetHashValue())
	_hashValueErr := /*TODO: migrate me*/ writeBuffer.WriteUint32("hashValue", 32, uint32((hashValue)))
	if _hashValueErr != nil {
		return errors.Wrap(_hashValueErr, "Error serializing 'hashValue' field")
	}

	// Simple Field (typeHashValue)
	typeHashValue := uint32(m.GetTypeHashValue())
	_typeHashValueErr := /*TODO: migrate me*/ writeBuffer.WriteUint32("typeHashValue", 32, uint32((typeHashValue)))
	if _typeHashValueErr != nil {
		return errors.Wrap(_typeHashValueErr, "Error serializing 'typeHashValue' field")
	}

	// Simple Field (size)
	size := uint32(m.GetSize())
	_sizeErr := /*TODO: migrate me*/ writeBuffer.WriteUint32("size", 32, uint32((size)))
	if _sizeErr != nil {
		return errors.Wrap(_sizeErr, "Error serializing 'size' field")
	}

	// Simple Field (offset)
	offset := uint32(m.GetOffset())
	_offsetErr := /*TODO: migrate me*/ writeBuffer.WriteUint32("offset", 32, uint32((offset)))
	if _offsetErr != nil {
		return errors.Wrap(_offsetErr, "Error serializing 'offset' field")
	}

	// Simple Field (dataType)
	dataType := uint32(m.GetDataType())
	_dataTypeErr := /*TODO: migrate me*/ writeBuffer.WriteUint32("dataType", 32, uint32((dataType)))
	if _dataTypeErr != nil {
		return errors.Wrap(_dataTypeErr, "Error serializing 'dataType' field")
	}

	// Simple Field (flags)
	flags := uint32(m.GetFlags())
	_flagsErr := /*TODO: migrate me*/ writeBuffer.WriteUint32("flags", 32, uint32((flags)))
	if _flagsErr != nil {
		return errors.Wrap(_flagsErr, "Error serializing 'flags' field")
	}

	// Implicit Field (propertyNameLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	propertyNameLength := uint16(uint16(len(m.GetPropertyName())))
	_propertyNameLengthErr := /*TODO: migrate me*/ writeBuffer.WriteUint16("propertyNameLength", 16, uint16((propertyNameLength)))
	if _propertyNameLengthErr != nil {
		return errors.Wrap(_propertyNameLengthErr, "Error serializing 'propertyNameLength' field")
	}

	// Implicit Field (dataTypeNameLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	dataTypeNameLength := uint16(uint16(len(m.GetDataTypeName())))
	_dataTypeNameLengthErr := /*TODO: migrate me*/ writeBuffer.WriteUint16("dataTypeNameLength", 16, uint16((dataTypeNameLength)))
	if _dataTypeNameLengthErr != nil {
		return errors.Wrap(_dataTypeNameLengthErr, "Error serializing 'dataTypeNameLength' field")
	}

	// Implicit Field (commentLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	commentLength := uint16(uint16(len(m.GetComment())))
	_commentLengthErr := /*TODO: migrate me*/ writeBuffer.WriteUint16("commentLength", 16, uint16((commentLength)))
	if _commentLengthErr != nil {
		return errors.Wrap(_commentLengthErr, "Error serializing 'commentLength' field")
	}

	// Simple Field (arrayDimensions)
	arrayDimensions := uint16(m.GetArrayDimensions())
	_arrayDimensionsErr := /*TODO: migrate me*/ writeBuffer.WriteUint16("arrayDimensions", 16, uint16((arrayDimensions)))
	if _arrayDimensionsErr != nil {
		return errors.Wrap(_arrayDimensionsErr, "Error serializing 'arrayDimensions' field")
	}

	// Simple Field (numChildren)
	numChildren := uint16(m.GetNumChildren())
	_numChildrenErr := /*TODO: migrate me*/ writeBuffer.WriteUint16("numChildren", 16, uint16((numChildren)))
	if _numChildrenErr != nil {
		return errors.Wrap(_numChildrenErr, "Error serializing 'numChildren' field")
	}

	// Simple Field (propertyName)
	propertyName := string(m.GetPropertyName())
	_propertyNameErr := /*TODO: migrate me*/ writeBuffer.WriteString("propertyName", uint32((uint16(len(m.GetPropertyName())))*(8)), (propertyName), utils.WithEncoding("UTF-8)"))
	if _propertyNameErr != nil {
		return errors.Wrap(_propertyNameErr, "Error serializing 'propertyName' field")
	}

	// Const Field (propertyNameTerminator)
	_propertyNameTerminatorErr := /*TODO: migrate me*/ /*TODO: migrate me*/ writeBuffer.WriteUint8("propertyNameTerminator", 8, uint8(0x00))
	if _propertyNameTerminatorErr != nil {
		return errors.Wrap(_propertyNameTerminatorErr, "Error serializing 'propertyNameTerminator' field")
	}

	// Simple Field (dataTypeName)
	dataTypeName := string(m.GetDataTypeName())
	_dataTypeNameErr := /*TODO: migrate me*/ writeBuffer.WriteString("dataTypeName", uint32((uint16(len(m.GetDataTypeName())))*(8)), (dataTypeName), utils.WithEncoding("UTF-8)"))
	if _dataTypeNameErr != nil {
		return errors.Wrap(_dataTypeNameErr, "Error serializing 'dataTypeName' field")
	}

	// Const Field (dataTypeNameTerminator)
	_dataTypeNameTerminatorErr := /*TODO: migrate me*/ /*TODO: migrate me*/ writeBuffer.WriteUint8("dataTypeNameTerminator", 8, uint8(0x00))
	if _dataTypeNameTerminatorErr != nil {
		return errors.Wrap(_dataTypeNameTerminatorErr, "Error serializing 'dataTypeNameTerminator' field")
	}

	// Simple Field (comment)
	comment := string(m.GetComment())
	_commentErr := /*TODO: migrate me*/ writeBuffer.WriteString("comment", uint32((uint16(len(m.GetComment())))*(8)), (comment), utils.WithEncoding("UTF-8)"))
	if _commentErr != nil {
		return errors.Wrap(_commentErr, "Error serializing 'comment' field")
	}

	// Const Field (commentTerminator)
	_commentTerminatorErr := /*TODO: migrate me*/ /*TODO: migrate me*/ writeBuffer.WriteUint8("commentTerminator", 8, uint8(0x00))
	if _commentTerminatorErr != nil {
		return errors.Wrap(_commentTerminatorErr, "Error serializing 'commentTerminator' field")
	}

	// Array Field (arrayInfo)
	if pushErr := writeBuffer.PushContext("arrayInfo", utils.WithRenderAsList(true)); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for arrayInfo")
	}
	for _curItem, _element := range m.GetArrayInfo() {
		_ = _curItem
		arrayCtx := utils.CreateArrayContext(ctx, len(m.GetArrayInfo()), _curItem)
		_ = arrayCtx
		_elementErr := writeBuffer.WriteSerializable(arrayCtx, _element)
		if _elementErr != nil {
			return errors.Wrap(_elementErr, "Error serializing 'arrayInfo' field")
		}
	}
	if popErr := writeBuffer.PopContext("arrayInfo", utils.WithRenderAsList(true)); popErr != nil {
		return errors.Wrap(popErr, "Error popping for arrayInfo")
	}

	// Array Field (children)
	if pushErr := writeBuffer.PushContext("children", utils.WithRenderAsList(true)); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for children")
	}
	for _curItem, _element := range m.GetChildren() {
		_ = _curItem
		arrayCtx := utils.CreateArrayContext(ctx, len(m.GetChildren()), _curItem)
		_ = arrayCtx
		_elementErr := writeBuffer.WriteSerializable(arrayCtx, _element)
		if _elementErr != nil {
			return errors.Wrap(_elementErr, "Error serializing 'children' field")
		}
	}
	if popErr := writeBuffer.PopContext("children", utils.WithRenderAsList(true)); popErr != nil {
		return errors.Wrap(popErr, "Error popping for children")
	}

	// Array Field (rest)
	// Byte Array field (rest)
	if err := writeBuffer.WriteByteArray("rest", m.GetRest()); err != nil {
		return errors.Wrap(err, "Error serializing 'rest' field")
	}

	if popErr := writeBuffer.PopContext("AdsDataTypeTableChildEntry"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for AdsDataTypeTableChildEntry")
	}
	return nil
}

func (m *_AdsDataTypeTableChildEntry) isAdsDataTypeTableChildEntry() bool {
	return true
}

func (m *_AdsDataTypeTableChildEntry) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
