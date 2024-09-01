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

package io

import (
	"math/big"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func WriteBoolean(writeBuffer utils.WriteBuffer) DataWriter[bool] {
	return NewDataWriterSimpleBoolean(writeBuffer, 1)
}

func WriteUnsignedByte(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[uint8] {
	return NewDataWriterSimpleUnsignedByte(writeBuffer, bitLength)
}

func WriteByte(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[byte] {
	return NewDataWriterSimpleByte(writeBuffer, bitLength)
}

func WriteByteArray(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[[]byte] {
	return NewDataWriterSimpleByteArray(writeBuffer, bitLength)
}

func WriteUnsignedShort(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[uint16] {
	return NewDataWriterSimpleUnsignedShort(writeBuffer, bitLength)
}

func WriteUnsignedInt(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[uint32] {
	return NewDataWriterSimpleUnsignedInt(writeBuffer, bitLength)
}

func WriteUnsignedLong(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[uint64] {
	return NewDataWriterSimpleUnsignedLong(writeBuffer, bitLength)
}

func WriteUnsignedBigInteger(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[*big.Int] {
	return NewDataWriterSimpleUnsignedBigInteger(writeBuffer, bitLength)
}

func WriteSignedByte(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[int8] {
	return NewDataWriterSimpleSignedByte(writeBuffer, bitLength)
}

func WriteSignedShort(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[int16] {
	return NewDataWriterSimpleSignedShort(writeBuffer, bitLength)
}

func WriteSignedInt(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[int32] {
	return NewDataWriterSimpleSignedInt(writeBuffer, bitLength)
}

func WriteSignedLong(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[int64] {
	return NewDataWriterSimpleSignedLong(writeBuffer, bitLength)
}

func WriteSignedBigInteger(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[*big.Int] {
	return NewDataWriterSimpleSignedBigInteger(writeBuffer, bitLength)
}

func WriteFloat(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[float32] {
	return NewDataWriterSimpleFloat(writeBuffer, bitLength)
}

func WriteDouble(writeBuffer utils.WriteBuffer, bitLength uint8) DataWriter[float64] {
	return NewDataWriterSimpleDouble(writeBuffer, bitLength)
}

func WriteString(writeBuffer utils.WriteBuffer, bitLength int32) DataWriter[string] {
	return NewDataWriterSimpleString(writeBuffer, bitLength)
}

func WriteEnum[T any, I any](enumSerializer func(T) I, enumNamer func(T) string, dataWriter DataWriter[I]) *DataWriterEnumDefault[T, I] {
	return NewDataWriterEnumDefault[T, I](enumSerializer, enumNamer, dataWriter)
}

func WriteComplex[T spi.Message](writeBuffer utils.WriteBuffer) *DataWriterComplexDefault[T] {
	return NewDataWriterComplexDefault[T](writeBuffer)
}

func WriteDataIO[T spi.Message](writeBuffer utils.WriteBuffer, serializer func(utils.WriteBuffer, values.PlcValue) error) *DataWriterDataIoDefault {
	return NewDataWriterDataIoDefault(writeBuffer, serializer)
}

func WriteDate(writeBuffer utils.WriteBuffer) DataWriter[time.Time] {
	return NewDataWriterSimpleDate(writeBuffer)
}

func WriteDateTime(writeBuffer utils.WriteBuffer) DataWriter[time.Time] {
	return NewDataWriterSimpleDateTime(writeBuffer)
}

func WriteTime(writeBuffer utils.WriteBuffer) DataWriter[time.Time] {
	return NewDataWriterSimpleTime(writeBuffer)
}
