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

type WriteBuffer struct {
	data   *bytes.Buffer
	writer *bitio.Writer
}

func NewWriteBuffer() *WriteBuffer {
	data := &bytes.Buffer{}
	writer := bitio.NewWriter(data)
	return &WriteBuffer{
		data:   data,
		writer: writer,
	}
}

func (rb WriteBuffer) GetPos() uint16 {
	return 0
}

func (rb WriteBuffer) GetBytes() []uint8 {
	return rb.data.Bytes()
}

func (rb WriteBuffer) GetTotalBytes() uint64 {
	return uint64(rb.data.Len())
}

func (rb WriteBuffer) WriteBit(value bool) error {
	return rb.writer.WriteBool(value)
}

func (rb WriteBuffer) WriteUint8(bitLength uint8, value uint8) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb WriteBuffer) WriteUint16(bitLength uint8, value uint16) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb WriteBuffer) WriteUint32(bitLength uint8, value uint32) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb WriteBuffer) WriteUint64(bitLength uint8, value uint64) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb WriteBuffer) WriteInt8(bitLength uint8, value int8) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb WriteBuffer) WriteInt16(bitLength uint8, value int16) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb WriteBuffer) WriteInt32(bitLength uint8, value int32) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb WriteBuffer) WriteInt64(bitLength uint8, value int64) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb WriteBuffer) WriteFloat32(bitLength uint8, value float32) error {
	res := math.Float32bits(value)
	return rb.writer.WriteBits(uint64(res), bitLength)
}

func (rb WriteBuffer) WriteFloat64(bitLength uint8, value float64) error {
	res := math.Float64bits(value)
	return rb.writer.WriteBits(res, bitLength)
}

func (rb WriteBuffer) WriteString(bitLength uint8, encoding string, value string) error {
	return errors.New("WriteString is currently not implemented")
}
