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

package utils

import (
	"context"
	"encoding/binary"
	"math"
	"math/big"
	"math/bits"
	"regexp"
	"strings"
	"unicode/utf16"

	"github.com/pkg/errors"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^A-Z0-9]+`)

type WriteBufferByteBased interface {
	WriteBuffer
	GetPos() uint16
	GetBytes() []byte
	GetTotalBytes() uint64
}

func NewWriteBufferByteBased(options ...WriteBufferByteBasedOptions) WriteBufferByteBased {
	b := &byteWriteBuffer{
		bits:      NewWriteBitBuffer(0),
		byteOrder: binary.BigEndian,
	}
	for _, option := range options {
		option(b)
	}
	return b
}

type WriteBufferByteBasedOptions = func(b *byteWriteBuffer)

func WithInitialSizeForByteBasedBuffer(length int) WriteBufferByteBasedOptions {
	return func(b *byteWriteBuffer) {
		if cap(b.bits.buf) < length {
			b.bits.buf = make([]byte, 0, length)
		}
	}
}

func WithByteOrderForByteBasedBuffer(byteOrder binary.ByteOrder) WriteBufferByteBasedOptions {
	return func(b *byteWriteBuffer) {
		b.byteOrder = byteOrder
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type byteWriteBuffer struct {
	BufferCommons
	bits      *WriteBitBuffer
	byteOrder binary.ByteOrder
	pos       uint
}

var _ WriteBuffer = (*byteWriteBuffer)(nil)

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (wb *byteWriteBuffer) SetByteOrder(byteOrder binary.ByteOrder) {
	wb.byteOrder = byteOrder
}

func (wb *byteWriteBuffer) GetByteOrder() binary.ByteOrder {
	return wb.byteOrder
}

func (wb *byteWriteBuffer) PushContext(_ string, _ ...WithWriterArgs) error {
	return nil
}

func (wb *byteWriteBuffer) GetPos() uint16 {
	return uint16(wb.pos / 8)
}

func (wb *byteWriteBuffer) GetBytes() []byte {
	return wb.bits.Bytes()
}

func (wb *byteWriteBuffer) GetTotalBytes() uint64 {
	return uint64(wb.bits.ByteLen())
}

func (wb *byteWriteBuffer) WriteBit(_ string, value bool, _ ...WithWriterArgs) error {
	wb.move(1)
	return wb.bits.WriteBool(value)
}

func (wb *byteWriteBuffer) WriteByte(_ string, value byte, _ ...WithWriterArgs) error {
	wb.move(8)
	return wb.bits.WriteBits(uint64(value), 8)
}

func (wb *byteWriteBuffer) WriteByteArray(_ string, data []byte, _ ...WithWriterArgs) error {
	for _, dataElement := range data {
		if err := wb.bits.WriteBits(uint64(dataElement), 8); err != nil {
			return err
		}
	}
	wb.move(uint(len(data) * 8))
	return nil
}

func (wb *byteWriteBuffer) WriteUint8(_ string, bitLength uint8, value uint8, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	return wb.bits.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteUint16(_ string, bitLength uint8, value uint16, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	if wb.byteOrder == binary.LittleEndian {
		reverseValue := bits.ReverseBytes64(uint64(value)) >> (64 - bitLength)
		return wb.bits.WriteBits(reverseValue, bitLength)
	}
	return wb.bits.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteUint32(_ string, bitLength uint8, value uint32, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	if wb.byteOrder == binary.LittleEndian {
		reverseValue := bits.ReverseBytes64(uint64(value)) >> (64 - bitLength)
		return wb.bits.WriteBits(reverseValue, bitLength)
	}
	return wb.bits.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteUint64(_ string, bitLength uint8, value uint64, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	if wb.byteOrder == binary.LittleEndian {
		reverseValue := bits.ReverseBytes64(value) >> (64 - bitLength)
		return wb.bits.WriteBits(reverseValue, bitLength)
	}
	return wb.bits.WriteBits(value, bitLength)
}

func (wb *byteWriteBuffer) WriteInt8(_ string, bitLength uint8, value int8, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	return wb.bits.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteInt16(_ string, bitLength uint8, value int16, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	if wb.byteOrder == binary.LittleEndian {
		reverseValue := bits.ReverseBytes64(uint64(value)) >> (64 - bitLength)
		return wb.bits.WriteBits(reverseValue, bitLength)
	}
	return wb.bits.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteInt32(_ string, bitLength uint8, value int32, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	if wb.byteOrder == binary.LittleEndian {
		reverseValue := bits.ReverseBytes64(uint64(value)) >> (64 - bitLength)
		return wb.bits.WriteBits(reverseValue, bitLength)
	}
	return wb.bits.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteInt64(_ string, bitLength uint8, value int64, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	if wb.byteOrder == binary.LittleEndian {
		reverseValue := bits.ReverseBytes64(uint64(value)) >> (64 - bitLength)
		return wb.bits.WriteBits(reverseValue, bitLength)
	}
	return wb.bits.WriteBits(uint64(value), bitLength)
}

func (wb *byteWriteBuffer) WriteBigInt(_ string, bitLength uint8, value *big.Int, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	return errors.New("not implemented yet")
}

func (wb *byteWriteBuffer) WriteFloat32(_ string, bitLength uint8, value float32, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	res := math.Float32bits(value)
	if wb.byteOrder == binary.LittleEndian {
		res = bits.ReverseBytes32(res)
	}
	return wb.bits.WriteBits(uint64(res), bitLength)
}

func (wb *byteWriteBuffer) WriteFloat64(_ string, bitLength uint8, value float64, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	res := math.Float64bits(value)
	if wb.byteOrder == binary.LittleEndian {
		res = bits.ReverseBytes64(res)
	}
	return wb.bits.WriteBits(res, bitLength)
}

func (wb *byteWriteBuffer) WriteBigFloat(_ string, bitLength uint8, value *big.Float, _ ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	return errors.New("not implemented yet")
}

func (wb *byteWriteBuffer) WriteString(_ string, bitLength uint32, value string, writerArgs ...WithWriterArgs) error {
	wb.move(uint(bitLength))
	encoding := nonAlphanumericRegex.ReplaceAllLiteralString(strings.ToUpper(wb.ExtractEncoding(UpcastWriterArgs(writerArgs...)...)), "")
	remainingBits := int64(bitLength) // int64 so subtraction doesn't wrap on underflow
	switch encoding {
	case "UTF8":
		for _, b := range []byte(value) {
			wb.bits.TryWriteByte(b)
			remainingBits -= 8
		}
	case "UTF16":
		fallthrough
	case "UTF16BE":
		for _, u := range utf16.Encode([]rune(value)) {
			wb.bits.TryWriteByte(byte(u >> 8))
			wb.bits.TryWriteByte(byte(u))
			remainingBits -= 16
		}
	case "UTF16LE":
		for _, u := range utf16.Encode([]rune(value)) {
			wb.bits.TryWriteByte(byte(u))
			wb.bits.TryWriteByte(byte(u >> 8))
			remainingBits -= 16
		}
	}
	// Fill remaining allocated space with zero bytes
	for range remainingBits / 8 {
		wb.bits.TryWriteByte(0x00)
	}
	return wb.bits.GetTryError()
}

func (wb *byteWriteBuffer) WriteVirtual(ctx context.Context, logicalName string, value any, writerArgs ...WithWriterArgs) error {
	// NO-OP
	return nil
}

func (wb *byteWriteBuffer) WriteSerializable(ctx context.Context, serializable Serializable) error {
	if serializable == nil {
		return nil
	}
	return serializable.SerializeWithWriteBuffer(ctx, wb)
}

func (wb *byteWriteBuffer) PopContext(_ string, _ ...WithWriterArgs) error {
	return nil
}

func (wb *byteWriteBuffer) move(bits uint) {
	wb.pos += bits
}
