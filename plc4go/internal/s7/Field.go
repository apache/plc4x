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
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
)

type PlcField interface {
	model.PlcField
	GetDataType() readWriteModel.TransportSize
	GetNumElements() uint16
	GetBlockNumber() uint16
	GetMemoryArea() readWriteModel.MemoryArea
	GetByteOffset() uint16
	GetBitOffset() uint8
}

type plcField struct {
	FieldType   FieldType
	MemoryArea  readWriteModel.MemoryArea
	BlockNumber uint16
	ByteOffset  uint16
	BitOffset   uint8
	NumElements uint16
	Datatype    readWriteModel.TransportSize
}

func NewField(memoryArea readWriteModel.MemoryArea, blockNumber uint16, byteOffset uint16, bitOffset uint8, numElements uint16, datatype readWriteModel.TransportSize) PlcField {
	return plcField{
		FieldType:   S7Field,
		MemoryArea:  memoryArea,
		BlockNumber: blockNumber,
		ByteOffset:  byteOffset,
		BitOffset:   bitOffset,
		NumElements: numElements,
		Datatype:    datatype,
	}
}

type PlcStringField struct {
	plcField
	stringLength uint16
}

func NewStringField(memoryArea readWriteModel.MemoryArea, blockNumber uint16, byteOffset uint16, bitOffset uint8, numElements uint16, stringLength uint16, datatype readWriteModel.TransportSize) PlcStringField {
	return PlcStringField{
		plcField: plcField{
			FieldType:   S7StringField,
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

func (m plcField) GetAddressString() string {
	// TODO: add missing variables like memory area, block number, byte offset, bit offset
	return fmt.Sprintf("%d:%s[%d]", m.FieldType, m.Datatype, m.NumElements)
}

func (m plcField) GetTypeName() string {
	return m.Datatype.String()
}

func (m plcField) GetDataType() readWriteModel.TransportSize {
	return m.Datatype
}

func (m plcField) GetNumElements() uint16 {
	return m.NumElements
}

func (m plcField) GetBlockNumber() uint16 {
	return m.BlockNumber
}

func (m plcField) GetMemoryArea() readWriteModel.MemoryArea {
	return m.MemoryArea
}

func (m plcField) GetByteOffset() uint16 {
	return m.ByteOffset
}

func (m plcField) GetBitOffset() uint8 {
	return m.BitOffset
}

func (m plcField) GetQuantity() uint16 {
	return m.NumElements
}

func (m plcField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.FieldType.GetName()); err != nil {
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

	if err := writeBuffer.PopContext(m.FieldType.GetName()); err != nil {
		return err
	}
	return nil
}

func (m PlcStringField) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext(m.FieldType.GetName()); err != nil {
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

	if err := writeBuffer.PopContext(m.FieldType.GetName()); err != nil {
		return err
	}
	return nil
}
