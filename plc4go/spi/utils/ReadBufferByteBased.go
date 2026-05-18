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
	"encoding/binary"
	"math"
	"math/big"
	"math/bits"

	"github.com/pkg/errors"
)

type ReadBufferByteBased interface {
	ReadBuffer
	GetBytes() []byte
	GetTotalBytes() uint64
	PeekByte(offset byte) byte
}

func NewReadBufferByteBased(data []byte, options ...ReadBufferByteBasedOptions) ReadBufferByteBased {
	b := &byteReadBuffer{
		data:      data,
		bits:      NewReadBitBuffer(data),
		pos:       uint64(0),
		byteOrder: binary.BigEndian,
	}
	for _, option := range options {
		option(b)
	}
	return b
}

type ReadBufferByteBasedOptions = func(b *byteReadBuffer)

func WithByteOrderForReadBufferByteBased(byteOrder binary.ByteOrder) ReadBufferByteBasedOptions {
	return func(b *byteReadBuffer) {
		b.byteOrder = byteOrder
	}
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type byteReadBuffer struct {
	data      []byte
	bits      *ReadBitBuffer
	pos       uint64
	byteOrder binary.ByteOrder
}

var _ ReadBuffer = (*byteReadBuffer)(nil)

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (rb *byteReadBuffer) SetByteOrder(byteOrder binary.ByteOrder) {
	rb.byteOrder = byteOrder
}

func (rb *byteReadBuffer) GetByteOrder() binary.ByteOrder {
	return rb.byteOrder
}

func (rb *byteReadBuffer) GetPos() uint16 {
	return uint16(rb.pos / 8)
}

func (rb *byteReadBuffer) Reset(pos uint16) {
	rb.pos = uint64(pos) * 8
	rb.bits.ResetTo(rb.pos)
}

func (rb *byteReadBuffer) GetBytes() []byte {
	return rb.data
}

func (rb *byteReadBuffer) GetTotalBytes() uint64 {
	return uint64(len(rb.data))
}

func (rb *byteReadBuffer) HasMore(bitLength uint8) bool {
	return (rb.pos + uint64(bitLength)) <= (uint64(len(rb.data)) * 8)
}

func (rb *byteReadBuffer) PeekByte(offset uint8) uint8 {
	return rb.data[rb.GetPos()+uint16(offset)]
}

func (rb *byteReadBuffer) PullContext(_ string, _ ...WithReaderArgs) error {
	return nil
}

func (rb *byteReadBuffer) ReadBit(_ string, _ ...WithReaderArgs) (bool, error) {
	rb.pos += 1
	return rb.bits.ReadBool()
}

func (rb *byteReadBuffer) ReadByte(_ string, _ ...WithReaderArgs) (byte, error) {
	rb.pos += 8
	return rb.bits.ReadByte()
}

func (rb *byteReadBuffer) ReadByteArray(_ string, numberOfBytes int, _ ...WithReaderArgs) ([]byte, error) {
	byteArray := make([]byte, numberOfBytes)
	for i := range numberOfBytes {
		rb.pos += 8
		readByte, err := rb.bits.ReadByte()
		if err != nil {
			return nil, err
		}
		byteArray[i] = readByte
	}
	return byteArray, nil
}

func (rb *byteReadBuffer) ReadUint8(_ string, bitLength uint8, _ ...WithReaderArgs) (uint8, error) {
	res, err := rb.bits.ReadBits(bitLength)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
	}
	rb.pos += uint64(bitLength)
	return uint8(res), nil
}

func (rb *byteReadBuffer) ReadUint16(_ string, bitLength uint8, _ ...WithReaderArgs) (uint16, error) {
	res, err := rb.bits.ReadBits(bitLength)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
	}
	rb.pos += uint64(bitLength)
	if rb.byteOrder == binary.LittleEndian {
		return uint16(bits.ReverseBytes64(res) >> (64 - bitLength)), nil
	}
	return uint16(res), nil
}

func (rb *byteReadBuffer) ReadUint32(_ string, bitLength uint8, _ ...WithReaderArgs) (uint32, error) {
	res, err := rb.bits.ReadBits(bitLength)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
	}
	rb.pos += uint64(bitLength)
	if rb.byteOrder == binary.LittleEndian {
		return uint32(bits.ReverseBytes64(res) >> (64 - bitLength)), nil
	}
	return uint32(res), nil
}

func (rb *byteReadBuffer) ReadUint64(_ string, bitLength uint8, _ ...WithReaderArgs) (uint64, error) {
	res, err := rb.bits.ReadBits(bitLength)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
	}
	rb.pos += uint64(bitLength)
	if rb.byteOrder == binary.LittleEndian {
		return bits.ReverseBytes64(res) >> (64 - bitLength), nil
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadInt8(_ string, bitLength uint8, _ ...WithReaderArgs) (int8, error) {
	res, err := rb.bits.ReadBits(bitLength)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
	}
	rb.pos += uint64(bitLength)
	return int8(res), nil
}

func (rb *byteReadBuffer) ReadInt16(_ string, bitLength uint8, _ ...WithReaderArgs) (int16, error) {
	res, err := rb.bits.ReadBits(bitLength)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
	}
	rb.pos += uint64(bitLength)
	if rb.byteOrder == binary.LittleEndian {
		return int16(bits.ReverseBytes64(res) >> (64 - bitLength)), nil
	}
	return int16(res), nil
}

func (rb *byteReadBuffer) ReadInt32(_ string, bitLength uint8, _ ...WithReaderArgs) (int32, error) {
	res, err := rb.bits.ReadBits(bitLength)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
	}
	rb.pos += uint64(bitLength)
	if rb.byteOrder == binary.LittleEndian {
		return int32(bits.ReverseBytes64(res) >> (64 - bitLength)), nil
	}
	return int32(res), nil
}

func (rb *byteReadBuffer) ReadInt64(_ string, bitLength uint8, _ ...WithReaderArgs) (int64, error) {
	res, err := rb.bits.ReadBits(bitLength)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
	}
	rb.pos += uint64(bitLength)
	if rb.byteOrder == binary.LittleEndian {
		return int64(bits.ReverseBytes64(res) >> (64 - bitLength)), nil
	}
	return int64(res), nil
}

func (rb *byteReadBuffer) ReadBigInt(_ string, bitLength uint64, _ ...WithReaderArgs) (*big.Int, error) {
	rawBytes := make([]byte, 0, (bitLength+7)/8)

	fullBytes := bitLength / 8
	remainingBits := uint8(bitLength % 8)

	for range fullBytes {
		b, err := rb.bits.ReadByte()
		if err != nil {
			return nil, errors.Wrapf(err, "error reading big int at pos [%d]bit ([%d]byte)", rb.pos, rb.pos/8)
		}
		rb.pos += 8
		rawBytes = append(rawBytes, b)
	}
	if remainingBits > 0 {
		b, err := rb.bits.ReadBits(remainingBits)
		if err != nil {
			return nil, errors.Wrapf(err, "error reading big int at pos [%d]bit ([%d]byte)", rb.pos, rb.pos/8)
		}
		rb.pos += uint64(remainingBits)
		rawBytes = append(rawBytes, byte(b))
	}

	res := new(big.Int).SetBytes(rawBytes)

	if rb.byteOrder == binary.LittleEndian {
		// rawBytes are in LE stream order; reverse to get BE for big.Int
		for i, j := 0, len(rawBytes)-1; i < j; i, j = i+1, j-1 {
			rawBytes[i], rawBytes[j] = rawBytes[j], rawBytes[i]
		}
		res.SetBytes(rawBytes)
	}

	return res, nil
}

func (rb *byteReadBuffer) ReadFloat32(logicalName string, bitLength uint8, _ ...WithReaderArgs) (float32, error) {
	if bitLength == 32 {
		_uintValue, err := rb.bits.ReadBits(bitLength)
		if err != nil {
			return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
		}
		rb.pos += uint64(bitLength)
		uintValue := uint32(_uintValue)
		if rb.byteOrder == binary.LittleEndian {
			array := make([]byte, 4)
			binary.LittleEndian.PutUint32(array, uintValue)
			uintValue = binary.BigEndian.Uint32(array)
		}
		return math.Float32frombits(uintValue), nil
	} else if bitLength < 32 {
		// TODO: Note ... this is the format as described in the KNX specification
		var err error
		sign, err := rb.ReadBit(logicalName)
		if err != nil {
			return 0.0, errors.Wrap(err, "error reading sign")
		}
		exp, err := rb.ReadInt32(logicalName, 5)
		if err != nil {
			return 0.0, errors.Wrap(err, "error reading exponent")
		}
		mantissa, err := rb.ReadUint32(logicalName, 10)
		// In the mantissa notation actually the first bit is omitted, we need to add it back
		f := (0.01 * float64(mantissa)) * math.Pow(float64(2), float64(exp))
		if sign {
			return -float32(f), nil
		}
		return float32(f), nil
	} else {
		return 0.0, errors.New("too many bits for float32")
	}
}

func (rb *byteReadBuffer) ReadFloat64(_ string, bitLength uint8, _ ...WithReaderArgs) (float64, error) {
	uintValue, err := rb.bits.ReadBits(bitLength)
	if err != nil {
		return 0, errors.Wrapf(err, "error reading %d bits at pos [%d]bit ([%d]byte)", bitLength, rb.pos, rb.pos/8)
	}
	rb.pos += uint64(bitLength)
	if rb.byteOrder == binary.LittleEndian {
		array := make([]byte, 8)
		binary.LittleEndian.PutUint64(array, uintValue)
		uintValue = binary.BigEndian.Uint64(array)
	}
	res := math.Float64frombits(uintValue)
	return res, nil
}

func (rb *byteReadBuffer) ReadBigFloat(logicalName string, bitLength uint8, _ ...WithReaderArgs) (*big.Float, error) {
	readFloat64, err := rb.ReadFloat64(logicalName, bitLength)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading float64")
	}
	return big.NewFloat(readFloat64), nil
}

func (rb *byteReadBuffer) ReadString(logicalName string, bitLength uint32, _ ...WithReaderArgs) (string, error) {
	stringBytes, err := rb.ReadByteArray(logicalName, int(bitLength/8))
	if err != nil {
		return "", errors.Wrap(err, "Error reading big int")
	}
	// TODO: make the null-termination a reader arg
	// End the string at the 0-character.
	for i, value := range stringBytes {
		if value == 0x00 {
			return string(stringBytes[0:i]), nil
		}
	}
	return string(stringBytes), nil
}

func (rb *byteReadBuffer) CloseContext(_ string, _ ...WithReaderArgs) error {
	return nil
}
