//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package s7

import (
	"encoding/xml"
	"fmt"
	readWrite "github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

func (m PlcField) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	log.Trace().Msg("MarshalXML")
	if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}

	if err := e.EncodeElement(m.MemoryArea.String(), xml.StartElement{Name: xml.Name{Local: "memoryArea"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.BlockNumber, xml.StartElement{Name: xml.Name{Local: "blockNumber"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.ByteOffset, xml.StartElement{Name: xml.Name{Local: "byteOffset"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.BitOffset, xml.StartElement{Name: xml.Name{Local: "bitOffset"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.NumElements, xml.StartElement{Name: xml.Name{Local: "numElements"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(m.Datatype.String(), xml.StartElement{Name: xml.Name{Local: "dataType"}}); err != nil {
		return err
	}

	if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: m.FieldType.GetName()}}); err != nil {
		return err
	}
	return nil
}
