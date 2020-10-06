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

type ReadBuffer struct {
	data    []uint8
	bytePos uint8
	bitPos  uint8
}

func (rb ReadBuffer) GetPos() uint16 {
	return 0
}

func (rb ReadBuffer) GetBytes() []uint8 {
	return rb.data
}

func (rb ReadBuffer) GetTotalBytes() uint64 {
	return 0
}

func (rb ReadBuffer) HasMore(bitLength uint8) bool {
	return false
}

func (rb ReadBuffer) PeekByte(offset uint8) uint8 {
	return 0
}

func (rb ReadBuffer) ReadBit() bool {
	return false
}

func (rb ReadBuffer) ReadUint8(bitLength uint8) uint8 {
	return 0
}

func (rb ReadBuffer) ReadUint16(bitLength uint8) uint16 {
	return 0
}

func (rb ReadBuffer) ReadUint32(bitLength uint8) uint32 {
	return 0
}

func (rb ReadBuffer) ReadUint64(bitLength uint8) uint64 {
	return 0
}

func (rb ReadBuffer) ReadInt8(bitLength uint8) int8 {
	return 0
}

func (rb ReadBuffer) ReadInt16(bitLength uint8) int16 {
	return 0
}

func (rb ReadBuffer) ReadInt32(bitLength uint8) int32 {
	return 0
}

func (rb ReadBuffer) ReadInt64(bitLength uint8) int64 {
	return 0
}

func (rb ReadBuffer) ReadFloat32(bitLength uint8) float32 {
	return 0
}

func (rb ReadBuffer) ReadFloat64(bitLength uint8) float64 {
	return 0
}

func (rb ReadBuffer) ReadString(bitLength uint8) string {
	return ""
}
