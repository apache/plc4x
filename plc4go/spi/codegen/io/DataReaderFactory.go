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
	"context"
	"math/big"
	"time"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

func ReadBoolean(readBuffer utils.ReadBuffer) DataReader[bool] {
	return NewDataReaderSimpleBoolean(readBuffer)
}

func ReadUnsignedByte(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[byte] {
	return NewDataReaderSimpleUnsignedByte(readBuffer, bitLength)
}

func ReadByte(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[byte] {
	return NewDataReaderSimpleByte(readBuffer, bitLength)
}

func ReadUnsignedShort(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[uint16] {
	return NewDataReaderSimpleUnsignedShort(readBuffer, bitLength)
}

func ReadUnsignedInt(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[uint32] {
	return NewDataReaderSimpleUnsignedInt(readBuffer, bitLength)
}

func ReadUnsignedLong(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[uint64] {
	return NewDataReaderSimpleUnsignedLong(readBuffer, bitLength)
}

func ReadUnsignedBigInteger(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[*big.Int] {
	return NewDataReaderSimpleUnsignedBigInteger(readBuffer, bitLength)
}

func ReadSignedByte(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[int8] {
	return NewDataReaderSimpleSignedByte(readBuffer, bitLength)
}

func ReadSignedShort(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[int16] {
	return NewDataReaderSimpleSignedShort(readBuffer, bitLength)
}

func ReadSignedInt(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[int32] {
	return NewDataReaderSimpleSignedInt(readBuffer, bitLength)
}

func ReadSignedLong(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[int64] {
	return NewDataReaderSimpleSignedLong(readBuffer, bitLength)
}

func ReadSignedBigInteger(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[*big.Int] {
	return NewDataReaderSimpleSignedBigInteger(readBuffer, bitLength)
}

func ReadFloat(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[float32] {
	return NewDataReaderSimpleFloat(readBuffer, bitLength)
}

func ReadDouble(readBuffer utils.ReadBuffer, bitLength uint8) DataReader[float64] {
	return NewDataReaderSimpleDouble(readBuffer, bitLength)
}

func ReadString(readBuffer utils.ReadBuffer, bitLength uint32) DataReader[string] {
	return NewDataReaderSimpleString(readBuffer, bitLength)
}

func ReadEnum[T any, I any](enumResolver func(I) (T, bool), dataReader DataReader[I]) *DataReaderEnumDefault[T, I] {
	return NewDataReaderEnumDefault[T, I](enumResolver, dataReader)
}

func ReadComplex[T any](complexSupplier func(context.Context, utils.ReadBuffer) (T, error), readBuffer utils.ReadBuffer) *DataReaderComplexDefault[T] {
	return NewDataReaderComplexDefault[T](complexSupplier, readBuffer)
}

func ReadDate(readBuffer utils.ReadBuffer) DataReader[time.Time] {
	return NewDataReaderSimpleDate(readBuffer)
}

func ReadDateTime(readBuffer utils.ReadBuffer) DataReader[time.Time] {
	return NewDataReaderSimpleDateTime(readBuffer)
}

func ReadTime(readBuffer utils.ReadBuffer) DataReader[time.Time] {
	return NewDataReaderSimpleTime(readBuffer)
}
