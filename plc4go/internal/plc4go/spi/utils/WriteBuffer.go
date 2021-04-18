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
	"math/big"
)

type WriteBuffer interface {
	PushContext(logicalName string) error
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
	PopContext(logicalName string) error
}
