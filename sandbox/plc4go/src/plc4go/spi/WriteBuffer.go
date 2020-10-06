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

type WriteBuffer struct {
	data    []uint8
	bytePos uint8
	bitPos  uint8
}

func (rb WriteBuffer) GetPos() uint16 {
	return 0
}

func (rb WriteBuffer) GetBytes() []uint8 {
	return rb.data
}

func (rb WriteBuffer) GetTotalBytes() uint64 {
	return 0
}

func (rb WriteBuffer) WriteBit(value bool) {
}

func (rb WriteBuffer) WriteUint8(bitLength uint8, value uint8) {
}

func (rb WriteBuffer) WriteUint16(bitLength uint8, value uint16) {
}

func (rb WriteBuffer) WriteUint32(bitLength uint8, value uint32) {
}

func (rb WriteBuffer) WriteUint64(bitLength uint8, value uint64) {
}

func (rb WriteBuffer) WriteInt8(bitLength uint8, value int8) {
}

func (rb WriteBuffer) WriteInt16(bitLength uint8, value int16) {
}

func (rb WriteBuffer) WriteInt32(bitLength uint8, value int32) {
}

func (rb WriteBuffer) WriteInt64(bitLength uint8, value int64) {
}

func (rb WriteBuffer) WriteFloat32(bitLength uint8, value float32) {
}

func (rb WriteBuffer) WriteFloat64(bitLength uint8, value float64) {
}

func (rb WriteBuffer) WriteString(bitLength uint8, value string) {
}
