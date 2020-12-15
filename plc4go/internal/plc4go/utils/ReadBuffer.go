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
package utils

import (
	"bytes"
	"errors"
	"github.com/icza/bitio"
	"math"
)

type ReadBuffer struct {
	data   []uint8
	reader *bitio.Reader
	pos    uint64
}

func NewReadBuffer(data []uint8) *ReadBuffer {
	buffer := bytes.NewBuffer(data)
	reader := bitio.NewReader(buffer)
	return &ReadBuffer{
		data:   data,
		reader: reader,
		pos:    uint64(0),
	}
}

func (rb *ReadBuffer) Reset() {
	rb.pos = uint64(0)
	rb.reader = bitio.NewReader(bytes.NewBuffer(rb.data))
}

func (rb *ReadBuffer) GetPos() uint16 {
	return uint16(rb.pos / 8)
}

func (rb *ReadBuffer) GetBytes() []uint8 {
	return rb.data
}

func (rb *ReadBuffer) GetTotalBytes() uint64 {
	return uint64(len(rb.data))
}

func (rb *ReadBuffer) HasMore(bitLength uint8) bool {
	return false
}

func (rb *ReadBuffer) PeekByte(offset uint8) uint8 {
	return 0
}

func (rb *ReadBuffer) ReadBit() (bool, error) {
	rb.pos += 1
	return rb.reader.ReadBool()
}

func (rb *ReadBuffer) ReadUint8(bitLength uint8) (uint8, error) {
	rb.pos += uint64(bitLength)
	res := uint8(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadUint16(bitLength uint8) (uint16, error) {
	rb.pos += uint64(bitLength)
	res := uint16(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadUint32(bitLength uint8) (uint32, error) {
	rb.pos += uint64(bitLength)
	res := uint32(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadUint64(bitLength uint8) (uint64, error) {
	rb.pos += uint64(bitLength)
	res := uint64(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadInt8(bitLength uint8) (int8, error) {
	rb.pos += uint64(bitLength)
	res := int8(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadInt16(bitLength uint8) (int16, error) {
	rb.pos += uint64(bitLength)
	res := int16(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadInt32(bitLength uint8) (int32, error) {
	rb.pos += uint64(bitLength)
	res := int32(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadInt64(bitLength uint8) (int64, error) {
	rb.pos += uint64(bitLength)
	res := int64(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadFloat32(signed bool, exponentBitLength uint8, mantissaBitLength uint8) (float32, error) {
	bitLength := exponentBitLength + mantissaBitLength
	if signed {
		bitLength++
	}
	if signed && exponentBitLength == 8 && mantissaBitLength == 23 {
		rb.pos += uint64(bitLength)
		uintValue := uint32(rb.reader.TryReadBits(bitLength))
		res := math.Float32frombits(uintValue)
		if rb.reader.TryError != nil {
			return 0, rb.reader.TryError
		}
		return res, nil
	} else if bitLength < 32 {
		// TODO: Note ... this is the format as described in the KNX specification
		sign := true
		var err error
		if signed {
			sign, err = rb.ReadBit()
			if err != nil {
				return 0.0, errors.New("error reading sign")
			}
		}
		exp, err := rb.ReadInt32(exponentBitLength)
		if err != nil {
			return 0.0, errors.New("error reading exponent")
		}
		mantissa, err := rb.ReadUint32(mantissaBitLength)
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

func (rb *ReadBuffer) ReadFloat64(signed bool, exponentBitLength uint8, mantissaBitLength uint8) (float64, error) {
	bitLength := 1 + exponentBitLength + mantissaBitLength
	rb.pos += uint64(bitLength)
	uintValue := rb.reader.TryReadBits(bitLength)
	res := math.Float64frombits(uintValue)
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadString(bitLength uint8) (string, error) {
	rb.pos += uint64(bitLength)
	res := string(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return "", rb.reader.TryError
	}
	return res, nil
}
