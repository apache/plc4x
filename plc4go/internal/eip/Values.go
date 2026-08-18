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

package eip

import (
	"encoding/binary"
	"math"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

const (
	stringLenOffset  = 2
	stringDataOffset = 6
)

// elementCount is the number of elements to ask the device for. An array tag
// has to request all of its elements (GH-1008) - otherwise the device returns
// a single element and the decoder has nothing to read the rest from.
func elementCount(tag PlcTag) uint16 {
	if tag.GetElementNb() > 1 {
		return tag.GetElementNb()
	}
	return 1
}

// isFixedSize reports whether the type is stored as a flat sequence of equally
// sized elements.
func isFixedSize(dataType readWriteModel.CIPDataTypeCode) bool {
	switch dataType {
	case readWriteModel.CIPDataTypeCode_SINT,
		readWriteModel.CIPDataTypeCode_INT,
		readWriteModel.CIPDataTypeCode_DINT,
		readWriteModel.CIPDataTypeCode_LINT,
		readWriteModel.CIPDataTypeCode_REAL,
		readWriteModel.CIPDataTypeCode_LREAL,
		readWriteModel.CIPDataTypeCode_BOOL:
		return true
	default:
		return false
	}
}

// parsePlcValue decodes a CIP read payload. CIP data is little-endian on the
// wire regardless of the encapsulation byte order.
func parsePlcValue(tag PlcTag, rawData []byte, dataType readWriteModel.CIPDataTypeCode) (values.PlcValue, error) {
	nb := int(elementCount(tag))
	elementSize := int(dataType.Size())
	if isFixedSize(dataType) && len(rawData) < nb*elementSize {
		// Never read past what the device actually sent - a short reply must not
		// blow up the whole response.
		return nil, errors.Errorf("device returned %d bytes for tag '%s', expected %d for %d element(s) of %s",
			len(rawData), tag.GetTag(), nb*elementSize, nb, dataType)
	}
	if nb > 1 {
		list := make([]values.PlcValue, 0, nb)
		for i := 0; i < nb; i++ {
			offset := i * elementSize
			switch dataType {
			case readWriteModel.CIPDataTypeCode_SINT:
				list = append(list, spiValues.NewPlcSINT(int8(rawData[offset])))
			case readWriteModel.CIPDataTypeCode_INT:
				list = append(list, spiValues.NewPlcINT(int16(binary.LittleEndian.Uint16(rawData[offset:]))))
			case readWriteModel.CIPDataTypeCode_DINT:
				list = append(list, spiValues.NewPlcDINT(int32(binary.LittleEndian.Uint32(rawData[offset:]))))
			case readWriteModel.CIPDataTypeCode_LINT:
				list = append(list, spiValues.NewPlcLINT(int64(binary.LittleEndian.Uint64(rawData[offset:]))))
			case readWriteModel.CIPDataTypeCode_REAL:
				list = append(list, spiValues.NewPlcREAL(math.Float32frombits(binary.LittleEndian.Uint32(rawData[offset:]))))
			case readWriteModel.CIPDataTypeCode_LREAL:
				list = append(list, spiValues.NewPlcLREAL(math.Float64frombits(binary.LittleEndian.Uint64(rawData[offset:]))))
			case readWriteModel.CIPDataTypeCode_BOOL:
				list = append(list, spiValues.NewPlcBOOL(rawData[offset] != 0))
			case readWriteModel.CIPDataTypeCode_STRING, readWriteModel.CIPDataTypeCode_STRUCTURED:
				// STRING/STRUCTURED carry their own length; plc4j decodes a single
				// string here and returns.
				value, err := parseStructured(tag, rawData)
				if err != nil {
					return nil, err
				}
				list = append(list, value)
				return spiValues.NewPlcList(list), nil
			default:
				return nil, errors.Errorf("unsupported type %s", dataType)
			}
		}
		return spiValues.NewPlcList(list), nil
	}
	switch dataType {
	case readWriteModel.CIPDataTypeCode_SINT:
		return spiValues.NewPlcSINT(int8(rawData[0])), nil
	case readWriteModel.CIPDataTypeCode_INT:
		return spiValues.NewPlcINT(int16(binary.LittleEndian.Uint16(rawData))), nil
	case readWriteModel.CIPDataTypeCode_DINT:
		return spiValues.NewPlcDINT(int32(binary.LittleEndian.Uint32(rawData))), nil
	case readWriteModel.CIPDataTypeCode_LINT:
		return spiValues.NewPlcLINT(int64(binary.LittleEndian.Uint64(rawData))), nil
	case readWriteModel.CIPDataTypeCode_REAL:
		return spiValues.NewPlcREAL(math.Float32frombits(binary.LittleEndian.Uint32(rawData))), nil
	case readWriteModel.CIPDataTypeCode_LREAL:
		return spiValues.NewPlcLREAL(math.Float64frombits(binary.LittleEndian.Uint64(rawData))), nil
	case readWriteModel.CIPDataTypeCode_BOOL:
		return spiValues.NewPlcBOOL(rawData[0] != 0), nil
	case readWriteModel.CIPDataTypeCode_STRING, readWriteModel.CIPDataTypeCode_STRUCTURED:
		return parseStructured(tag, rawData)
	default:
		return nil, errors.Errorf("unsupported type %s", dataType)
	}
}

func parseStructured(tag PlcTag, rawData []byte) (values.PlcValue, error) {
	if len(rawData) < stringDataOffset {
		return nil, errors.Errorf("structured payload for tag '%s' too short: %d bytes", tag.GetTag(), len(rawData))
	}
	structuredType := binary.LittleEndian.Uint16(rawData[0:])
	structuredLen := int(binary.LittleEndian.Uint16(rawData[stringLenOffset:]))
	if structuredType != uint16(readWriteModel.CIPStructTypeCode_STRING) {
		return nil, errors.Errorf("unsupported structured type %#x for tag '%s'", structuredType, tag.GetTag())
	}
	if len(rawData) < stringDataOffset+structuredLen {
		return nil, errors.Errorf("string payload for tag '%s' too short: %d bytes for length %d",
			tag.GetTag(), len(rawData), structuredLen)
	}
	return spiValues.NewPlcSTRING(string(rawData[stringDataOffset : stringDataOffset+structuredLen])), nil
}

// encodeValue serializes a single value for a CIP write. Like plc4j it encodes
// one element; array writes carry the element count in the CipWriteRequest.
func encodeValue(value apiValues.PlcValue, dataType readWriteModel.CIPDataTypeCode) ([]byte, error) {
	switch dataType {
	case readWriteModel.CIPDataTypeCode_BOOL:
		boolByte := byte(0)
		if value.GetBool() {
			boolByte = 1
		}
		return []byte{boolByte}, nil
	case readWriteModel.CIPDataTypeCode_SINT:
		return []byte{byte(value.GetInt8())}, nil
	case readWriteModel.CIPDataTypeCode_INT:
		return binary.LittleEndian.AppendUint16(nil, uint16(value.GetInt16())), nil
	case readWriteModel.CIPDataTypeCode_DINT:
		return binary.LittleEndian.AppendUint32(nil, uint32(value.GetInt32())), nil
	case readWriteModel.CIPDataTypeCode_LINT:
		return binary.LittleEndian.AppendUint64(nil, uint64(value.GetInt64())), nil
	case readWriteModel.CIPDataTypeCode_REAL:
		return binary.LittleEndian.AppendUint32(nil, math.Float32bits(value.GetFloat32())), nil
	case readWriteModel.CIPDataTypeCode_LREAL:
		return binary.LittleEndian.AppendUint64(nil, math.Float64bits(value.GetFloat64())), nil
	case readWriteModel.CIPDataTypeCode_STRING, readWriteModel.CIPDataTypeCode_STRUCTURED:
		text := value.GetString()
		raw := binary.LittleEndian.AppendUint32(nil, uint32(len(text)))
		return append(raw, text...), nil
	default:
		return nil, errors.Errorf("unsupported type %s", dataType)
	}
}

// decodeResponseCode converts a CIP status into a PLC4X response code.
func decodeResponseCode(status uint8) apiModel.PlcResponseCode {
	switch status {
	case 0:
		return apiModel.PlcResponseCode_OK
	default:
		return apiModel.PlcResponseCode_INTERNAL_ERROR
	}
}
