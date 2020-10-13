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
package spi

import (
	"bytes"
	"github.com/icza/bitio"
)

type ReadBuffer struct {
	buffer *bytes.Buffer
	reader *bitio.Reader
	pos    uint64
}

func ReadBufferNew(data []uint8) *ReadBuffer {
	buffer := bytes.NewBuffer(data)
	reader := bitio.NewReader(buffer)
	return &ReadBuffer{
		buffer: buffer,
		reader: reader,
		pos:    uint64(0),
	}
}

func (rb ReadBuffer) GetPos() uint16 {
	return uint16(rb.pos / 8)
}

func (rb ReadBuffer) GetBytes() []uint8 {
	return rb.buffer.Bytes()
}

func (rb ReadBuffer) GetTotalBytes() uint64 {
	return uint64(rb.buffer.Len())
}

func (rb ReadBuffer) HasMore(bitLength uint8) bool {
	return false
}

func (rb ReadBuffer) PeekByte(offset uint8) uint8 {
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

func (rb *ReadBuffer) ReadFloat32(bitLength uint8) (float32, error) {
	rb.pos += uint64(bitLength)
	res := float32(rb.reader.TryReadBits(bitLength))
	if rb.reader.TryError != nil {
		return 0, rb.reader.TryError
	}
	return res, nil
}

func (rb *ReadBuffer) ReadFloat64(bitLength uint8) (float64, error) {
	rb.pos += uint64(bitLength)
	res := float64(rb.reader.TryReadBits(bitLength))
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
