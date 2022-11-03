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
	"bytes"
	"encoding/binary"
	"github.com/icza/bitio"
	"github.com/pkg/errors"
	"math"
	"math/big"
)

type ReadBufferByteBased interface {
	ReadBuffer
	GetBytes() []byte
	GetTotalBytes() uint64
	PeekByte(offset byte) byte
}

func NewReadBufferByteBased(data []byte, options ...ReadBufferByteBasedOptions) ReadBufferByteBased {
	buffer := bytes.NewBuffer(data)
	reader := bitio.NewReader(buffer)
	b := &byteReadBuffer{
		data:      data,
		reader:    reader,
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
	reader    *bitio.Reader
	pos       uint64
	byteOrder binary.ByteOrder
}

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
	rb.pos = uint64(0)
	rb.reader = bitio.NewReader(bytes.NewBuffer(rb.data))
	bytesToSkip := make([]byte, pos)
	_, err := rb.reader.Read(bytesToSkip)
	if err != nil {
		panic(err)
	}
	rb.pos = uint64(pos * 8)
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
	return rb.reader.ReadBool()
}

func (rb *byteReadBuffer) ReadByte(_ string, _ ...WithReaderArgs) (byte, error) {
	rb.pos += 8
	return rb.reader.ReadByte()
}

func (rb *byteReadBuffer) ReadByteArray(_ string, numberOfBytes int, _ ...WithReaderArgs) ([]byte, error) {
	byteArray := make([]byte, numberOfBytes)
	for i := 0; i < numberOfBytes; i++ {
		rb.pos += 8
		readByte, err := rb.reader.ReadByte()
		if err != nil {
			return nil, err
		}
		byteArray[i] = readByte
	}
	return byteArray, nil
}

func (rb *byteReadBuffer) ReadUint8(_ string, bitLength uint8, _ ...WithReaderArgs) (uint8, error) {
	rb.pos += uint64(bitLength)
	res := uint8(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadUint16(logicalName string, bitLength uint8, _ ...WithReaderArgs) (uint16, error) {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		bigInt, err := rb.ReadBigInt(logicalName, uint64(bitLength))
		if err != nil {
			return 0, err
		}
		return uint16(bigInt.Uint64()), nil
	}
	rb.pos += uint64(bitLength)
	res := uint16(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadUint32(logicalName string, bitLength uint8, _ ...WithReaderArgs) (uint32, error) {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		bigInt, err := rb.ReadBigInt(logicalName, uint64(bitLength))
		if err != nil {
			return 0, err
		}
		return uint32(bigInt.Uint64()), nil
	}
	rb.pos += uint64(bitLength)
	res := uint32(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadUint64(logicalName string, bitLength uint8, _ ...WithReaderArgs) (uint64, error) {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		bigInt, err := rb.ReadBigInt(logicalName, uint64(bitLength))
		if err != nil {
			return 0, err
		}
		return bigInt.Uint64(), nil
	}
	rb.pos += uint64(bitLength)
	res := rb.reader.TryReadBits(bitLength)
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadInt8(_ string, bitLength uint8, _ ...WithReaderArgs) (int8, error) {
	rb.pos += uint64(bitLength)
	res := int8(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadInt16(logicalName string, bitLength uint8, _ ...WithReaderArgs) (int16, error) {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		bigInt, err := rb.ReadBigInt(logicalName, uint64(bitLength))
		if err != nil {
			return 0, err
		}
		return int16(bigInt.Int64()), nil
	}
	rb.pos += uint64(bitLength)
	res := int16(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadInt32(logicalName string, bitLength uint8, _ ...WithReaderArgs) (int32, error) {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		bigInt, err := rb.ReadBigInt(logicalName, uint64(bitLength))
		if err != nil {
			return 0, err
		}
		return int32(bigInt.Int64()), nil
	}
	rb.pos += uint64(bitLength)
	res := int32(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadInt64(logicalName string, bitLength uint8, _ ...WithReaderArgs) (int64, error) {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		bigInt, err := rb.ReadBigInt(logicalName, uint64(bitLength))
		if err != nil {
			return 0, err
		}
		return bigInt.Int64(), nil
	}
	rb.pos += uint64(bitLength)
	res := int64(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadBigInt(_ string, bitLength uint64, _ ...WithReaderArgs) (*big.Int, error) {
	// TODO: highly experimental remove this comment when tested or verifyed
	res := big.NewInt(0)

	// TODO: maybe we can use left shift and or of big int
	rawBytes := make([]byte, 0)
	correction := uint8(0)

	for remainingBits := bitLength; remainingBits > 0; {
		// we can max read 64 bit with bitio
		bitToRead := uint8(64)
		if remainingBits < 64 {
			bitToRead = uint8(remainingBits)
		}
		// we now read the bits
		data := rb.reader.TryReadBits(bitToRead)

		// and check for uneven bits for a right shift at the end
		correction = 64 - bitToRead
		data <<= correction

		dataBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(dataBytes, data)
		rawBytes = append(rawBytes, dataBytes...)

		if rb.reader.TryError != nil {
			return big.NewInt(0), rb.reader.TryError
		}
		remainingBits -= uint64(bitToRead)
	}

	res.SetBytes(rawBytes)

	// now we need to shift the last correction to right again
	res.Rsh(res, uint(correction))
	if rb.byteOrder == binary.LittleEndian {
		originalByteLength := len(rawBytes) - int(correction/8)
		resBytes := res.Bytes()
		padding := make([]byte, originalByteLength-len(resBytes))
		resBytes = append(padding, resBytes...)
		if rb.byteOrder == binary.LittleEndian {
			for i, j := 0, len(resBytes)-1; i <= j; i, j = i+1, j-1 {
				resBytes[i], resBytes[j] = resBytes[j], resBytes[i]
			}
		}
		res.SetBytes(resBytes)
	}

	return res, nil
}

func (rb *byteReadBuffer) ReadFloat32(logicalName string, bitLength uint8, _ ...WithReaderArgs) (float32, error) {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		bigInt, err := rb.ReadBigFloat(logicalName, bitLength)
		if err != nil {
			return 0, err
		}
		f, _ := bigInt.Float32()
		return f, nil
	}
	if bitLength == 32 {
		rb.pos += uint64(bitLength)
		uintValue := uint32(rb.reader.TryReadBits(bitLength))
		res := math.Float32frombits(uintValue)
		if rb.reader.TryError != nil {
			return 0, rb.reader.TryError
		}
		return res, nil
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
	rb.pos += uint64(bitLength)
	uintValue := rb.reader.TryReadBits(bitLength)
	if rb.byteOrder == binary.LittleEndian {
		array := make([]byte, 8)
		binary.LittleEndian.PutUint64(array, uintValue)
		uintValue = binary.BigEndian.Uint64(array)
	}
	res := math.Float64frombits(uintValue)
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *byteReadBuffer) ReadBigFloat(logicalName string, bitLength uint8, _ ...WithReaderArgs) (*big.Float, error) {
	readFloat64, err := rb.ReadFloat64(logicalName, bitLength)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading float64")
	}
	return big.NewFloat(readFloat64), nil
}

func (rb *byteReadBuffer) ReadString(logicalName string, bitLength uint32, encoding string, _ ...WithReaderArgs) (string, error) {
	stringBytes, err := rb.ReadByteArray(logicalName, int(bitLength/8))
	if err != nil {
		return "", errors.Wrap(err, "Error reading big int")
	}
	return string(stringBytes), nil
}

func (rb *byteReadBuffer) CloseContext(_ string, _ ...WithReaderArgs) error {
	return nil
}
