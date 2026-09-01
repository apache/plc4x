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

package modbus

import (
	"context"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/spi/values"
)

// A Modbus register is 16 bits wide. DataItem encodes one value at its natural width and knows
// nothing about that, so a value narrower than a register is padded here when it stands alone and
// packed when several of them follow one another. Which half of the register a padded value sits
// in depends on the byte order.

// widthBits is the width of one value of the given type, in bits. A string of the largest
// declarable length is wider than 16 bits can express, and a run of such strings wider than 32,
// so the arithmetic around this is done in 64 bits throughout.
func widthBits(dataType readWriteModel.ModbusDataType, stringLength uint16) uint64 {
	switch dataType {
	case readWriteModel.ModbusDataType_BOOL:
		return 1
	case readWriteModel.ModbusDataType_BYTE,
		readWriteModel.ModbusDataType_SINT,
		readWriteModel.ModbusDataType_USINT,
		readWriteModel.ModbusDataType_CHAR:
		return 8
	case readWriteModel.ModbusDataType_WORD,
		readWriteModel.ModbusDataType_INT,
		readWriteModel.ModbusDataType_UINT,
		readWriteModel.ModbusDataType_WCHAR:
		return 16
	case readWriteModel.ModbusDataType_DWORD,
		readWriteModel.ModbusDataType_DINT,
		readWriteModel.ModbusDataType_UDINT,
		readWriteModel.ModbusDataType_REAL:
		return 32
	case readWriteModel.ModbusDataType_LWORD,
		readWriteModel.ModbusDataType_LINT,
		readWriteModel.ModbusDataType_ULINT,
		readWriteModel.ModbusDataType_LREAL:
		return 64
	case readWriteModel.ModbusDataType_STRING:
		return uint64(stringLength) * 8
	case readWriteModel.ModbusDataType_WSTRING:
		return uint64(stringLength) * 16
	default:
		return 0
	}
}

// paddingBits is what a lone value of this type needs on top of its own width to fill a register.
// CHAR is deliberately left unpadded, as it always has been.
func paddingBits(dataType readWriteModel.ModbusDataType) uint16 {
	switch dataType {
	case readWriteModel.ModbusDataType_BOOL:
		return 15
	case readWriteModel.ModbusDataType_BYTE,
		readWriteModel.ModbusDataType_SINT,
		readWriteModel.ModbusDataType_USINT:
		return 8
	default:
		return 0
	}
}

// leadingPaddingBitsLittleEndian says where the padding sits for a lone value in little endian
// order. A byte simply comes first and is padded behind it, but a BOOL sits between seven bits
// and eight.
func leadingPaddingBitsLittleEndian(dataType readWriteModel.ModbusDataType) uint16 {
	if dataType == readWriteModel.ModbusDataType_BOOL {
		return 7
	}
	return 0
}

// trailingPaddingBits rounds a packed run of values up to a whole register. The width of one
// element has to be the declared one - a STRING(20) is 160 bits wide, not the 8 bits a single
// character takes - or a run of them gets a pad byte that nothing on the read side accounts for.
func trailingPaddingBits(dataType readWriteModel.ModbusDataType, numberOfValues uint16, stringLength uint16) uint16 {
	width := widthBits(dataType, stringLength)
	if width >= 16 {
		return 0
	}
	remainder := uint16(width) * numberOfValues % 16
	if remainder == 0 {
		return 0
	}
	return 16 - remainder
}

// ParseRegisters reads numberOfValues values of the given type out of the register bytes as the
// given byte order lays them out, returning the value itself for a single one and a list for
// several.
func ParseRegisters(ctx context.Context, data []byte, dataType readWriteModel.ModbusDataType, numberOfValues uint16, byteOrder ByteOrder, stringLength uint16) (apiValues.PlcValue, error) {
	// The two byte-swap modes differ from their plain counterparts only in that the two bytes of
	// every register are exchanged, so undo that first and read what is left as usual (plc4j
	// ModbusTcpConnection.toPlcValue).
	if byteOrder.swapsBytes() {
		data = byteSwap(data)
	}
	readBuffer := utils.NewReadBufferByteBased(data, utils.WithByteOrderForReadBufferByteBased(byteOrder.bufferByteOrder()))
	bigEndian := byteOrder.isBigEndian()
	if numberOfValues == 1 {
		padding := paddingBits(dataType)
		if padding == 0 {
			return readWriteModel.DataItemParseWithBuffer(ctx, readBuffer, dataType, stringLength)
		}
		if bigEndian {
			if _, err := readBuffer.ReadUint16("", uint8(padding)); err != nil {
				return nil, errors.Wrap(err, "error reading padding")
			}
			return readWriteModel.DataItemParseWithBuffer(ctx, readBuffer, dataType, stringLength)
		}
		leading := leadingPaddingBitsLittleEndian(dataType)
		if leading > 0 {
			if _, err := readBuffer.ReadUint16("", uint8(leading)); err != nil {
				return nil, errors.Wrap(err, "error reading leading padding")
			}
		}
		value, err := readWriteModel.DataItemParseWithBuffer(ctx, readBuffer, dataType, stringLength)
		if err != nil {
			return nil, err
		}
		if _, err := readBuffer.ReadUint16("", uint8(padding-leading)); err != nil {
			return nil, errors.Wrap(err, "error reading trailing padding")
		}
		return value, nil
	}

	// Several values are packed without padding between them; a trailing pad is left unread.
	parsed := make([]apiValues.PlcValue, 0, numberOfValues)
	for range numberOfValues {
		value, err := readWriteModel.DataItemParseWithBuffer(ctx, readBuffer, dataType, stringLength)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing data item")
		}
		parsed = append(parsed, value)
	}
	return values.NewPlcList(parsed), nil
}

// SerializeRegisters writes a value, or every element of a list, with the layout ParseRegisters
// expects for the same byte order.
func SerializeRegisters(ctx context.Context, value apiValues.PlcValue, dataType readWriteModel.ModbusDataType, numberOfValues uint16, byteOrder ByteOrder, stringLength uint16) ([]byte, error) {
	data, err := serializeRegisters(ctx, value, dataType, numberOfValues, byteOrder, stringLength)
	if err != nil {
		return nil, err
	}
	if byteOrder.swapsBytes() {
		data = byteSwap(data)
	}
	return data, nil
}

func serializeRegisters(ctx context.Context, value apiValues.PlcValue, dataType readWriteModel.ModbusDataType, numberOfValues uint16, byteOrder ByteOrder, stringLength uint16) ([]byte, error) {
	writeBuffer := utils.NewWriteBufferByteBased(
		utils.WithInitialSizeForByteBasedBuffer(int(lengthInBytes(dataType, numberOfValues, stringLength))),
		utils.WithByteOrderForByteBasedBuffer(byteOrder.bufferByteOrder()),
	)
	bigEndian := byteOrder.isBigEndian()
	if numberOfValues == 1 {
		padding := paddingBits(dataType)
		if padding == 0 {
			if err := readWriteModel.DataItemSerializeWithWriteBuffer(ctx, writeBuffer, value, dataType, stringLength); err != nil {
				return nil, err
			}
			return writeBuffer.GetBytes(), nil
		}
		if bigEndian {
			if err := writeBuffer.WriteUint16("", uint8(padding), 0); err != nil {
				return nil, errors.Wrap(err, "error writing padding")
			}
			if err := readWriteModel.DataItemSerializeWithWriteBuffer(ctx, writeBuffer, value, dataType, stringLength); err != nil {
				return nil, err
			}
			return writeBuffer.GetBytes(), nil
		}
		leading := leadingPaddingBitsLittleEndian(dataType)
		if leading > 0 {
			if err := writeBuffer.WriteUint16("", uint8(leading), 0); err != nil {
				return nil, errors.Wrap(err, "error writing leading padding")
			}
		}
		if err := readWriteModel.DataItemSerializeWithWriteBuffer(ctx, writeBuffer, value, dataType, stringLength); err != nil {
			return nil, err
		}
		if err := writeBuffer.WriteUint16("", uint8(padding-leading), 0); err != nil {
			return nil, errors.Wrap(err, "error writing trailing padding")
		}
		return writeBuffer.GetBytes(), nil
	}

	for i := range numberOfValues {
		element := value
		if value.IsList() {
			element = value.GetList()[i]
		}
		if err := readWriteModel.DataItemSerializeWithWriteBuffer(ctx, writeBuffer, element, dataType, stringLength); err != nil {
			return nil, errors.Wrap(err, "error serializing data item")
		}
	}
	if trailing := trailingPaddingBits(dataType, numberOfValues, stringLength); trailing > 0 {
		if err := writeBuffer.WriteUint16("", uint8(trailing), 0); err != nil {
			return nil, errors.Wrap(err, "error writing trailing padding")
		}
	}
	return writeBuffer.GetBytes(), nil
}

// lengthInBytes is how many bytes SerializeRegisters will write, and equally how many
// ParseRegisters will consume.
func lengthInBytes(dataType readWriteModel.ModbusDataType, numberOfValues uint16, stringLength uint16) uint64 {
	if numberOfValues == 1 {
		return (widthBits(dataType, stringLength) + uint64(paddingBits(dataType)) + 7) / 8
	}
	bits := (widthBits(dataType, stringLength) * uint64(numberOfValues)) + uint64(trailingPaddingBits(dataType, numberOfValues, stringLength))
	return (bits + 7) / 8
}
