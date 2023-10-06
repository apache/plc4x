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

package s7

import (
	"context"
	"encoding/binary"
	"fmt"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcTag interface {
	apiModel.PlcTag
	utils.Serializable

	GetDataType() readWriteModel.TransportSize
	GetNumElements() uint16
	GetBlockNumber() uint16
	GetMemoryArea() readWriteModel.MemoryArea
	GetByteOffset() uint16
	GetBitOffset() uint8
}

type plcTag struct {
	TagType     TagType
	MemoryArea  readWriteModel.MemoryArea
	BlockNumber uint16
	ByteOffset  uint16
	BitOffset   uint8
	NumElements uint16
	Datatype    readWriteModel.TransportSize
}

func NewTag(memoryArea readWriteModel.MemoryArea, blockNumber uint16, byteOffset uint16, bitOffset uint8, numElements uint16, datatype readWriteModel.TransportSize) PlcTag {
	return plcTag{
		TagType:     S7Tag,
		MemoryArea:  memoryArea,
		BlockNumber: blockNumber,
		ByteOffset:  byteOffset,
		BitOffset:   bitOffset,
		NumElements: numElements,
		Datatype:    datatype,
	}
}

type PlcStringTag struct {
	plcTag
	stringLength uint16
}

func NewStringTag(memoryArea readWriteModel.MemoryArea, blockNumber uint16, byteOffset uint16, bitOffset uint8, numElements uint16, stringLength uint16, datatype readWriteModel.TransportSize) PlcStringTag {
	return PlcStringTag{
		plcTag: plcTag{
			TagType:     S7StringTag,
			MemoryArea:  memoryArea,
			BlockNumber: blockNumber,
			ByteOffset:  byteOffset,
			BitOffset:   bitOffset,
			NumElements: numElements,
			Datatype:    datatype,
		},
		stringLength: stringLength,
	}
}

func (m plcTag) GetAddressString() string {
	// TODO: add missing variables like memory area, block number, byte offset, bit offset
	return fmt.Sprintf("%d:%s[%d]", m.TagType, m.Datatype, m.NumElements)
}

func (m plcTag) GetValueType() apiValues.PlcValueType {
	if plcValueByName, ok := apiValues.PlcValueByName(m.Datatype.String()); ok {
		return plcValueByName
	}
	return apiValues.NULL
}

func (m plcTag) GetArrayInfo() []apiModel.ArrayInfo {
	if m.NumElements != 1 {
		return []apiModel.ArrayInfo{
			&spiModel.DefaultArrayInfo{
				LowerBound: 0,
				UpperBound: uint32(m.NumElements),
			},
		}
	}
	return []apiModel.ArrayInfo{}
}

func (m plcTag) GetDataType() readWriteModel.TransportSize {
	return m.Datatype
}

func (m plcTag) GetNumElements() uint16 {
	return m.NumElements
}

func (m plcTag) GetBlockNumber() uint16 {
	return m.BlockNumber
}

func (m plcTag) GetMemoryArea() readWriteModel.MemoryArea {
	return m.MemoryArea
}

func (m plcTag) GetByteOffset() uint16 {
	return m.ByteOffset
}

func (m plcTag) GetBitOffset() uint8 {
	return m.BitOffset
}

func (m plcTag) GetQuantity() uint16 {
	return m.NumElements
}

func (m plcTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m plcTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.TagType.GetName()); err != nil {
		return err
	}

	if err := writeBuffer.WriteString("memoryArea", uint32(len(m.MemoryArea.String())*8), "UTF-8", m.MemoryArea.String()); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint16("blockNumber", 16, m.BlockNumber); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint16("byteOffset", 16, m.ByteOffset); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint8("bitOffset", 8, m.BitOffset); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint16("numElements", 16, m.NumElements); err != nil {
		return err
	}
	if err := writeBuffer.WriteString("dataType", uint32(len(m.Datatype.String())*8), "UTF-8", m.Datatype.String()); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(m.TagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (m plcTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}

func (m PlcStringTag) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcStringTag) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.TagType.GetName()); err != nil {
		return err
	}

	if err := writeBuffer.WriteString("memoryArea", uint32(len(m.MemoryArea.String())*8), "UTF-8", m.MemoryArea.String()); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint16("blockNumber", 16, m.BlockNumber); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint16("byteOffset", 16, m.ByteOffset); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint8("bitOffset", 8, m.BitOffset); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint16("numElements", 16, m.NumElements); err != nil {
		return err
	}
	if err := writeBuffer.WriteUint16("stringLength", 16, m.stringLength); err != nil {
		return err
	}
	if err := writeBuffer.WriteString("dataType", uint32(len(m.Datatype.String())*8), "UTF-8", m.Datatype.String()); err != nil {
		return err
	}

	if err := writeBuffer.PopContext(m.TagType.GetName()); err != nil {
		return err
	}
	return nil
}

func (m PlcStringTag) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
