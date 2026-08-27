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

// valueReader decodes one element of a fixed-size type at the given offset into the payload.
type valueReader func(rawData []byte, offset int) values.PlcValue

// valueWriter encodes one element of a fixed-size type.
type valueWriter func(value values.PlcValue) []byte

type typeCodec struct {
	read  valueReader
	write valueWriter
}

// fixedSizeCodecs holds the types that are stored as a flat sequence of equally sized elements,
// and how to decode and encode one. It is the single source of truth for that set: it drives the
// short-reply guard, the array decoding, the scalar decoding and the write path alike, so a type
// cannot be added to one of them and silently forgotten in the others.
var fixedSizeCodecs = map[readWriteModel.CIPDataTypeCode]typeCodec{
	readWriteModel.CIPDataTypeCode_BOOL: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcBOOL(rawData[offset] != 0)
		},
		write: func(value values.PlcValue) []byte {
			boolByte := byte(0)
			if value.GetBool() {
				boolByte = 1
			}
			return []byte{boolByte}
		},
	},
	readWriteModel.CIPDataTypeCode_SINT: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcSINT(int8(rawData[offset]))
		},
		write: func(value values.PlcValue) []byte {
			return []byte{byte(value.GetInt8())}
		},
	},
	readWriteModel.CIPDataTypeCode_INT: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcINT(int16(binary.LittleEndian.Uint16(rawData[offset:])))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint16(nil, uint16(value.GetInt16()))
		},
	},
	readWriteModel.CIPDataTypeCode_DINT: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcDINT(int32(binary.LittleEndian.Uint32(rawData[offset:])))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint32(nil, uint32(value.GetInt32()))
		},
	},
	readWriteModel.CIPDataTypeCode_LINT: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcLINT(int64(binary.LittleEndian.Uint64(rawData[offset:])))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint64(nil, uint64(value.GetInt64()))
		},
	},
	readWriteModel.CIPDataTypeCode_USINT: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcUSINT(rawData[offset])
		},
		write: func(value values.PlcValue) []byte {
			return []byte{value.GetUint8()}
		},
	},
	readWriteModel.CIPDataTypeCode_UINT: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcUINT(binary.LittleEndian.Uint16(rawData[offset:]))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint16(nil, value.GetUint16())
		},
	},
	readWriteModel.CIPDataTypeCode_UDINT: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcUDINT(binary.LittleEndian.Uint32(rawData[offset:]))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint32(nil, value.GetUint32())
		},
	},
	readWriteModel.CIPDataTypeCode_ULINT: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcULINT(binary.LittleEndian.Uint64(rawData[offset:]))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint64(nil, value.GetUint64())
		},
	},
	readWriteModel.CIPDataTypeCode_BYTE: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcBYTE(rawData[offset])
		},
		write: func(value values.PlcValue) []byte {
			return []byte{value.GetUint8()}
		},
	},
	readWriteModel.CIPDataTypeCode_WORD: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcWORD(binary.LittleEndian.Uint16(rawData[offset:]))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint16(nil, value.GetUint16())
		},
	},
	readWriteModel.CIPDataTypeCode_DWORD: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcDWORD(binary.LittleEndian.Uint32(rawData[offset:]))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint32(nil, value.GetUint32())
		},
	},
	readWriteModel.CIPDataTypeCode_LWORD: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcLWORD(binary.LittleEndian.Uint64(rawData[offset:]))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint64(nil, value.GetUint64())
		},
	},
	readWriteModel.CIPDataTypeCode_REAL: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcREAL(math.Float32frombits(binary.LittleEndian.Uint32(rawData[offset:])))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint32(nil, math.Float32bits(value.GetFloat32()))
		},
	},
	readWriteModel.CIPDataTypeCode_LREAL: {
		read: func(rawData []byte, offset int) values.PlcValue {
			return spiValues.NewPlcLREAL(math.Float64frombits(binary.LittleEndian.Uint64(rawData[offset:])))
		},
		write: func(value values.PlcValue) []byte {
			return binary.LittleEndian.AppendUint64(nil, math.Float64bits(value.GetFloat64()))
		},
	},
}

// isFixedSize reports whether the type is stored as a flat sequence of equally sized elements.
func isFixedSize(dataType readWriteModel.CIPDataTypeCode) bool {
	_, found := fixedSizeCodecs[dataType]
	return found
}

// parsePlcValue decodes a CIP read payload. CIP data is little-endian on the
// wire regardless of the encapsulation byte order.
func parsePlcValue(tag PlcTag, rawData []byte, dataType readWriteModel.CIPDataTypeCode) (values.PlcValue, error) {
	nb := int(elementCount(tag))
	codec, fixedSize := fixedSizeCodecs[dataType]
	if !fixedSize {
		// STRING and STRUCTURED carry their own length rather than being laid out element by
		// element, so only the first one can be decoded from the reply.
		value, err := parseStructured(tag, rawData, dataType)
		if err != nil || nb == 1 {
			return value, err
		}
		return spiValues.NewPlcList([]values.PlcValue{value}), nil
	}

	// Never read past what the device actually sent - a short reply must not blow up the whole
	// response.
	elementSize := int(dataType.Size())
	if len(rawData) < nb*elementSize {
		return nil, errors.Errorf("device returned %d bytes for tag '%s', expected %d for %d element(s) of %s",
			len(rawData), tag.GetTag(), nb*elementSize, nb, dataType)
	}
	if nb == 1 {
		return codec.read(rawData, 0), nil
	}
	list := make([]values.PlcValue, 0, nb)
	for i := range nb {
		list = append(list, codec.read(rawData, i*elementSize))
	}
	return spiValues.NewPlcList(list), nil
}

// parseStructured decodes a STRING or STRUCTURED payload, which carries its own length. Any
// other type is not decodable here.
func parseStructured(tag PlcTag, rawData []byte, dataType readWriteModel.CIPDataTypeCode) (values.PlcValue, error) {
	if dataType != readWriteModel.CIPDataTypeCode_STRING && dataType != readWriteModel.CIPDataTypeCode_STRUCTURED {
		return nil, errors.Errorf("unsupported type %s", dataType)
	}
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
func encodeValue(value values.PlcValue, dataType readWriteModel.CIPDataTypeCode) ([]byte, error) {
	if codec, fixedSize := fixedSizeCodecs[dataType]; fixedSize {
		return codec.write(value), nil
	}
	switch dataType {
	case readWriteModel.CIPDataTypeCode_STRING, readWriteModel.CIPDataTypeCode_STRUCTURED:
		return encodeString(value.GetString(), dataType)
	default:
		return nil, errors.Errorf("unsupported type %s", dataType)
	}
}

// encodeString writes the structure the read path parses back: a 2-byte structure handle, a
// 4-byte length and then the characters, padded out to the size the CipWriteRequest serializer
// emits for the type (its 'count' expression is dataType.size * elementNb). The length is the
// number of bytes, not characters - they differ for any non-ASCII text, and a length that
// disagrees with the bytes that follow reads back truncated. Text that does not fit cannot be
// sent, so it is reported rather than silently losing the tail.
func encodeString(text string, dataType readWriteModel.CIPDataTypeCode) ([]byte, error) {
	capacity := int(dataType.Size()) - stringDataOffset
	if len(text) > capacity {
		return nil, errors.Errorf("a %s holds at most %d bytes of text, but %d were given",
			dataType, max(capacity, 0), len(text))
	}
	raw := make([]byte, dataType.Size())
	binary.LittleEndian.PutUint16(raw, uint16(readWriteModel.CIPStructTypeCode_STRING))
	binary.LittleEndian.PutUint32(raw[stringLenOffset:], uint32(len(text)))
	copy(raw[stringDataOffset:], text)
	return raw, nil
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
