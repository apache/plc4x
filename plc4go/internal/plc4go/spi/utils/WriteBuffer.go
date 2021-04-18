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
	"encoding/binary"
	"errors"
	"github.com/icza/bitio"
	"math"
	"math/big"
)

type WriteBuffer interface {
	GetPos() uint16
	GetBytes() []uint8
	GetTotalBytes() uint64
	WriteBit(logicalName string, value bool) error
	WriteUint8(logicalName string, bitLength uint8, value uint8) error
	WriteUint16(logicalName string, bitLength uint8, value uint16) error
	WriteUint32(logicalName string, bitLength uint8, value uint32) error
	WriteUint64(logicalName string, bitLength uint8, value uint64) error
	WriteInt8(logicalName string, bitLength uint8, value int8) error
	WriteInt16(logicalName string, bitLength uint8, value int16) error
	WriteInt32(logicalName string, bitLength uint8, value int32) error
	WriteInt64(logicalName string, bitLength uint8, value int64) error
	WriteBigInt(logicalName string, bitLength uint8, value *big.Int) error
	WriteFloat32(logicalName string, bitLength uint8, value float32) error
	WriteFloat64(logicalName string, bitLength uint8, value float64) error
	WriteString(logicalName string, bitLength uint8, encoding string, value string) error
}

func NewWriteBuffer() WriteBuffer {
	data := &bytes.Buffer{}
	writer := bitio.NewWriter(data)
	return &writeBuffer{
		data:      data,
		writer:    writer,
		byteOrder: binary.BigEndian,
	}
}

func NewLittleEndianWriteBuffer() WriteBuffer {
	data := &bytes.Buffer{}
	writer := bitio.NewWriter(data)
	return &writeBuffer{
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

type writeBuffer struct {
	data      *bytes.Buffer
	writer    *bitio.Writer
	byteOrder binary.ByteOrder
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////

func (rb *writeBuffer) GetPos() uint16 {
	return 0
}

func (rb *writeBuffer) GetBytes() []uint8 {
	return rb.data.Bytes()
}

func (rb *writeBuffer) GetTotalBytes() uint64 {
	return uint64(rb.data.Len())
}

func (rb *writeBuffer) WriteBit(_ string, value bool) error {
	return rb.writer.WriteBool(value)
}

func (rb *writeBuffer) WriteUint8(_ string, bitLength uint8, value uint8) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb *writeBuffer) WriteUint16(_ string, bitLength uint8, value uint16) error {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(rb.data, rb.byteOrder, value)
	}
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb *writeBuffer) WriteUint32(_ string, bitLength uint8, value uint32) error {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(rb.data, rb.byteOrder, value)
	}
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb *writeBuffer) WriteUint64(_ string, bitLength uint8, value uint64) error {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(rb.data, rb.byteOrder, value)
	}
	return rb.writer.WriteBits(value, bitLength)
}

func (rb *writeBuffer) WriteInt8(_ string, bitLength uint8, value int8) error {
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb *writeBuffer) WriteInt16(_ string, bitLength uint8, value int16) error {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(rb.data, rb.byteOrder, value)
	}
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb *writeBuffer) WriteInt32(_ string, bitLength uint8, value int32) error {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(rb.data, rb.byteOrder, value)
	}
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb *writeBuffer) WriteInt64(_ string, bitLength uint8, value int64) error {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(rb.data, rb.byteOrder, value)
	}
	return rb.writer.WriteBits(uint64(value), bitLength)
}

func (rb *writeBuffer) WriteBigInt(_ string, bitLength uint8, value *big.Int) error {
	return errors.New("not implemented yet")
}

func (rb *writeBuffer) WriteFloat32(_ string, bitLength uint8, value float32) error {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(rb.data, rb.byteOrder, value)
	}
	res := math.Float32bits(value)
	return rb.writer.WriteBits(uint64(res), bitLength)
}

func (rb *writeBuffer) WriteFloat64(_ string, bitLength uint8, value float64) error {
	if rb.byteOrder == binary.LittleEndian {
		// TODO: indirection till we have a native LE implementation
		// TODO: validate that this produces the desired result
		return binary.Write(rb.data, rb.byteOrder, value)
	}
	res := math.Float64bits(value)
	return rb.writer.WriteBits(res, bitLength)
}

func (rb *writeBuffer) WriteString(_ string, bitLength uint8, encoding string, value string) error {
	return errors.New("WriteString is currently not implemented")
}
