/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package utils

import (
	"bytes"
	"encoding/binary"
	"github.com/icza/bitio"
	"github.com/pkg/errors"
	"math"
	"math/big"
)

type WriteBufferByteBased interface {
	WriteBuffer
	GetPos() uint16
	GetBytes() []byte
	GetTotalBytes() uint64
}

func NewWriteBufferByteBased() WriteBufferByteBased {
	data := &bytes.Buffer{}
	writer := bitio.NewWriter(data)
	return &byteWriteBuffer{
		data:      data,
		writer:    writer,
		byteOrder: binary.BigEndian,
	}
}

func NewLittleEndianWriteBufferByteBased() WriteBufferByteBased {
	data := &bytes.Buffer{}
	writer := bitio.NewWriter(data)
	return &byteWriteBuffer{
		data:      data,
		writer:    writer,
		byteOrder: binary.LittleEndian,
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type byteWriteBuffer struct {
	data      *bytes.Buffer
	writer    *bitio.Writer
	byteOrder binary.ByteOrder
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (wb *byteWriteBuffer) PushContext(_ string, _ ...WithWriterArgs) error {
	return nil
}

func (wb *byteWriteBuffer) PopContext(_ string, _ ...WithWriterArgs) error {
	return nil
}

func (wb *byteWriteBuffer) SetByteOrder(byteOrder binary.ByteOrder) {
	wb.byteOrder = byteOrder
}

func (wb *byteWriteBuffer) GetByteOrder() binary.ByteOrder {
	return wb.byteOrder
}

func (wb *byteWriteBuffer) GetPos() uint16 {
	return 0
}

func (wb *byteWriteBuffer) GetBytes() []byte {
	return wb.data.Bytes()
}

func (wb *byteWriteBuffer) GetTotalBytes() uint64 {
	return uint64(wb.data.Len())
}

func (wb *byteWriteBuffer) WriteBit(_ string, value bool, _ ...WithWriterArgs) error {
	return wb.writer.WriteBool(value)
}

func (wb *byteWriteBuffer) WriteByte(_ string, value byte, _ ...WithWriterArgs) error {
	return wb.writer.WriteBits(uint64(value), 8)
}

func (wb *byteWriteBuffer) WriteByteArray(_ string, data []byte, _ ...WithWriterArgs) error {
	for _, dataElement := range data {
		err := wb.writer.WriteBits(uint64(dataElement), 8)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wb *byteWriteBuffer) WriteUint8(_ string, bitLength uint8, value uint8, _ ...WithWriterArgs) error {
	return wb.writer.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteUint16(_ string, bitLength uint8, value uint16, _ ...WithWriterArgs) error {
	if wb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(wb.data, wb.byteOrder, value)
	}
	return wb.writer.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteUint32(_ string, bitLength uint8, value uint32, _ ...WithWriterArgs) error {
	if wb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(wb.data, wb.byteOrder, value)
	}
	return wb.writer.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteUint64(_ string, bitLength uint8, value uint64, _ ...WithWriterArgs) error {
	if wb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(wb.data, wb.byteOrder, value)
	}
	return wb.writer.WriteBits(value, bitLength)
}

func (wb *byteWriteBuffer) WriteInt8(_ string, bitLength uint8, value int8, _ ...WithWriterArgs) error {
	return wb.writer.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteInt16(_ string, bitLength uint8, value int16, _ ...WithWriterArgs) error {
	if wb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(wb.data, wb.byteOrder, value)
	}
	return wb.writer.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteInt32(_ string, bitLength uint8, value int32, _ ...WithWriterArgs) error {
	if wb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(wb.data, wb.byteOrder, value)
	}
	return wb.writer.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteInt64(_ string, bitLength uint8, value int64, _ ...WithWriterArgs) error {
	if wb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(wb.data, wb.byteOrder, value)
	}
	return wb.writer.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteBigInt(_ string, bitLength uint8, value *big.Int, _ ...WithWriterArgs) error {
	return errors.New("not implemented yet")
}

func (wb *byteWriteBuffer) WriteFloat32(_ string, bitLength uint8, value float32, _ ...WithWriterArgs) error {
	if wb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(wb.data, wb.byteOrder, value)
	}
	res := math.Float32bits(value)
	return wb.writer.WriteBits(uint64(res), bitLength)
}

func (wb *byteWriteBuffer) WriteFloat64(_ string, bitLength uint8, value float64, _ ...WithWriterArgs) error {
	if wb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(wb.data, wb.byteOrder, value)
	}
	res := math.Float64bits(value)
	return wb.writer.WriteBits(res, bitLength)
}

func (wb *byteWriteBuffer) WriteBigFloat(_ string, bitLength uint8, value *big.Float, _ ...WithWriterArgs) error {
	return errors.New("not implemented yet")
}

func (wb *byteWriteBuffer) WriteString(_ string, bitLength uint8, encoding string, value string, _ ...WithWriterArgs) error {
	// TODO: the implementation completely ignores encoding for now. Fix this
	for _, theByte := range []byte(value) {
		wb.writer.TryWriteByte(theByte)
	}
	return wb.writer.TryError
}
