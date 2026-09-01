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

package slmp

import (
	"encoding/binary"
	"math"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// DataType is a data type an slmp tag can carry, ported from plc4j's SlmpDataType. It maps the type
// to its SLMP word footprint and owns the conversion between little-endian SLMP words and a
// plc4x value in both directions.
//
// The wire layer leaves the read response's and the write request's payload as raw bytes on purpose
// (see the mspec), so typed encoding and decoding lives here rather than in the generated model.
// Word units only - plc4j's first version has no bit-unit support and neither does this one.
type DataType uint8

const (
	// DataTypeWORD is an unsigned 16-bit bit pattern, the type a tag without an explicit one gets.
	DataTypeWORD DataType = iota
	DataTypeINT
	DataTypeUINT
	DataTypeDINT
	DataTypeUDINT
	DataTypeREAL
)

// dataTypeNames is the spelling of each type in a tag address, which is also the plc4x value type
// name it maps to.
var dataTypeNames = map[DataType]string{
	DataTypeWORD:  "WORD",
	DataTypeINT:   "INT",
	DataTypeUINT:  "UINT",
	DataTypeDINT:  "DINT",
	DataTypeUDINT: "UDINT",
	DataTypeREAL:  "REAL",
}

// SupportedDataTypeNames is every type name a tag address may name, in the order plc4j's enum
// declares them. Used for the "supported: ..." part of the parse error.
var SupportedDataTypeNames = []string{"WORD", "INT", "UINT", "DINT", "UDINT", "REAL"}

func (d DataType) String() string {
	if name, ok := dataTypeNames[d]; ok {
		return name
	}
	return "Unknown"
}

// DataTypeByName resolves a type name out of a tag address.
func DataTypeByName(name string) (DataType, bool) {
	for dataType, dataTypeName := range dataTypeNames {
		if dataTypeName == name {
			return dataType, true
		}
	}
	return DataTypeWORD, false
}

// WordsPerElement is how many 16-bit words one element of this type occupies. SLMP Batch
// Read/Write in word units counts points in words, so this is the factor between a tag's quantity
// and the numberOfPoints the frame asks for.
func (d DataType) WordsPerElement() uint16 {
	switch d {
	case DataTypeDINT, DataTypeUDINT, DataTypeREAL:
		return 2
	default:
		return 1
	}
}

// GetValueType is the plc4x value type this type decodes into.
func (d DataType) GetValueType() apiValues.PlcValueType {
	if valueType, ok := apiValues.PlcValueTypeByName(d.String()); ok {
		return valueType
	}
	return apiValues.NULL
}

// Decode reads quantity elements of this type out of the little-endian response payload. A payload
// shorter than required is an error the caller maps to INVALID_DATA, matching how plc4j's
// SlmpDataType.decode returns null for it.
//
// Deliberately tolerant of a payload that is longer than required: a device is free to answer with
// the points that were asked for and nothing about the extra bytes makes the requested ones wrong.
func (d DataType) Decode(responseData []byte, quantity uint16) (apiValues.PlcValue, error) {
	requiredBytes := int(quantity) * int(d.WordsPerElement()) * 2
	if len(responseData) < requiredBytes {
		return nil, errors.Errorf("SLMP response carries %d byte(s) but %d are needed for %d %s element(s)",
			len(responseData), requiredBytes, quantity, d)
	}
	if quantity == 1 {
		return d.decodeOne(responseData), nil
	}
	elements := make([]apiValues.PlcValue, 0, quantity)
	stride := int(d.WordsPerElement()) * 2
	for i := 0; i < int(quantity); i++ {
		elements = append(elements, d.decodeOne(responseData[i*stride:]))
	}
	return spiValues.NewPlcList(elements), nil
}

// decodeOne decodes the element at the head of data, which is guaranteed to be long enough.
func (d DataType) decodeOne(data []byte) apiValues.PlcValue {
	switch d {
	case DataTypeWORD:
		return spiValues.NewPlcWORD(binary.LittleEndian.Uint16(data))
	case DataTypeINT:
		return spiValues.NewPlcINT(int16(binary.LittleEndian.Uint16(data)))
	case DataTypeUINT:
		return spiValues.NewPlcUINT(binary.LittleEndian.Uint16(data))
	case DataTypeDINT:
		return spiValues.NewPlcDINT(int32(binary.LittleEndian.Uint32(data)))
	case DataTypeUDINT:
		return spiValues.NewPlcUDINT(binary.LittleEndian.Uint32(data))
	case DataTypeREAL:
		return spiValues.NewPlcREAL(math.Float32frombits(binary.LittleEndian.Uint32(data)))
	default:
		// Unreachable: every DataType this package hands out is one of the six above.
		return spiValues.NewPlcNULL()
	}
}

// Encode writes quantity elements of this type as little-endian SLMP words. The result is always
// exactly quantity * WordsPerElement * 2 bytes long, which is what keeps the write request's
// announced point count and its payload length in agreement - the generated serializer derives the
// payload length from numberOfPoints, so a disagreement would go out as a frame no device accepts.
//
// A type or arity mismatch is an error the caller maps to INVALID_DATA, symmetric with Decode and
// with plc4j's SlmpDataType.encode returning null for it.
func (d DataType) Encode(value apiValues.PlcValue, quantity uint16) ([]byte, error) {
	stride := int(d.WordsPerElement()) * 2
	data := make([]byte, int(quantity)*stride)
	if value == nil {
		return nil, errors.Errorf("no value to write as %d %s element(s)", quantity, d)
	}
	if quantity == 1 {
		if value.IsList() {
			return nil, errors.Errorf("a scalar %s tag can't be written from a list of %d value(s)", d, len(value.GetList()))
		}
		if err := d.encodeOne(data, value); err != nil {
			return nil, err
		}
		return data, nil
	}
	if !value.IsList() {
		return nil, errors.Errorf("a %s[%d] tag needs a list of %d values", d, quantity, quantity)
	}
	elements := value.GetList()
	if len(elements) != int(quantity) {
		return nil, errors.Errorf("a %s[%d] tag needs exactly %d values but got %d", d, quantity, quantity, len(elements))
	}
	for i, element := range elements {
		if element == nil {
			return nil, errors.Errorf("value %d of a %s[%d] tag is nil", i, d, quantity)
		}
		if err := d.encodeOne(data[i*stride:], element); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// encodeOne writes one element into the head of data, which is guaranteed to be long enough.
func (d DataType) encodeOne(data []byte, value apiValues.PlcValue) error {
	switch d {
	case DataTypeWORD, DataTypeUINT:
		// The full 0..65535 range: a WORD is a bit pattern and a UINT is unsigned, so a value that
		// only fits when read as signed is a mistake worth reporting rather than truncating.
		if !value.IsUint16() {
			return errors.Errorf("%s doesn't fit into an unsigned 16-bit SLMP word", value)
		}
		binary.LittleEndian.PutUint16(data, value.GetUint16())
	case DataTypeINT:
		if !value.IsInt16() {
			return errors.Errorf("%s doesn't fit into a signed 16-bit SLMP word", value)
		}
		binary.LittleEndian.PutUint16(data, uint16(value.GetInt16()))
	case DataTypeDINT:
		if !value.IsInt32() {
			return errors.Errorf("%s doesn't fit into a signed 32-bit SLMP double word", value)
		}
		binary.LittleEndian.PutUint32(data, uint32(value.GetInt32()))
	case DataTypeUDINT:
		if !value.IsUint32() {
			return errors.Errorf("%s doesn't fit into an unsigned 32-bit SLMP double word", value)
		}
		binary.LittleEndian.PutUint32(data, value.GetUint32())
	case DataTypeREAL:
		if !value.IsFloat32() {
			return errors.Errorf("%s isn't a 32-bit float", value)
		}
		binary.LittleEndian.PutUint32(data, math.Float32bits(value.GetFloat32()))
	default:
		// Unreachable: every DataType this package hands out is one of the six above.
		return errors.Errorf("unsupported SLMP data type %d", uint8(d))
	}
	return nil
}
