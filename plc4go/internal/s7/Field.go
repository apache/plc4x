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
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	readWrite "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/pkg/errors"
)

type S7PlcField interface {
	GetDataType() readWrite.TransportSize
	GetNumElements() uint16
	GetBlockNumber() uint16
	GetMemoryArea() readWrite.MemoryArea
	GetByteOffset() uint16
	GetBitOffset() uint8
}

type PlcField struct {
	FieldType   FieldType
	MemoryArea  readWrite.MemoryArea
	BlockNumber uint16
	ByteOffset  uint16
	BitOffset   uint8
	NumElements uint16
	Datatype    readWrite.TransportSize
}

func NewField(memoryArea readWrite.MemoryArea, blockNumber uint16, byteOffset uint16, bitOffset uint8, numElements uint16, datatype readWrite.TransportSize) PlcField {
	return PlcField{
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
	PlcField
	stringLength uint16
}

func NewStringField(memoryArea readWrite.MemoryArea, blockNumber uint16, byteOffset uint16, bitOffset uint8, numElements uint16, stringLength uint16, datatype readWrite.TransportSize) PlcStringField {
	return PlcStringField{
		PlcField: PlcField{
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

func (m PlcField) GetAddressString() string {
	// TODO: add missing variables like memory area, block number, byte offset, bit offset
	return fmt.Sprintf("%d:%s[%d]", m.FieldType, m.Datatype, m.NumElements)
}

func (m PlcField) GetTypeName() string {
	return m.Datatype.String()
}

func (m PlcField) GetDataType() readWrite.TransportSize {
	return m.Datatype
}

func (m PlcField) GetNumElements() uint16 {
	return m.NumElements
}

func (m PlcField) GetBlockNumber() uint16 {
	return m.BlockNumber
}

func (m PlcField) GetMemoryArea() readWrite.MemoryArea {
	return m.MemoryArea
}

func (m PlcField) GetByteOffset() uint16 {
	return m.ByteOffset
}

func (m PlcField) GetBitOffset() uint8 {
	return m.BitOffset
}

func (m PlcField) GetQuantity() uint16 {
	return m.NumElements
}

func CastTos7FieldFromPlcField(plcField model.PlcField) (PlcField, error) {
	if s7Field, ok := plcField.(PlcField); ok {
		return s7Field, nil
	}
	return PlcField{}, errors.New("couldn't cast to s7PlcField")
}

func (m PlcField) Serialize(writeBuffer utils.WriteBuffer) error {
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
	if err := writeBuffer.WriteUint16("numElements", 16, m.stringLength); err != nil {
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
