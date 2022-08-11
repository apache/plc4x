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
	"math/big"
)

type ReadBuffer interface {
	// GetPos return the current byte position
	GetPos() uint16
	// Reset sets the position to the supplied byte position
	Reset(pos uint16)
	// HasMore returns true if there are bitLength bits available
	HasMore(bitLength uint8) bool
	// PullContext signals that we expect now a context with the supplied logical name
	PullContext(logicalName string, readerArgs ...WithReaderArgs) error
	ReadBit(logicalName string, readerArgs ...WithReaderArgs) (bool, error)
	ReadByte(logicalName string, readerArgs ...WithReaderArgs) (byte, error)
	ReadByteArray(logicalName string, numberOfBytes int, readerArgs ...WithReaderArgs) ([]byte, error)
	ReadUint8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint8, error)
	ReadUint16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint16, error)
	ReadUint32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint32, error)
	ReadUint64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (uint64, error)
	ReadInt8(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int8, error)
	ReadInt16(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int16, error)
	ReadInt32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int32, error)
	ReadInt64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (int64, error)
	ReadBigInt(logicalName string, bitLength uint64, readerArgs ...WithReaderArgs) (*big.Int, error)
	ReadFloat32(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (float32, error)
	ReadFloat64(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (float64, error)
	ReadBigFloat(logicalName string, bitLength uint8, readerArgs ...WithReaderArgs) (*big.Float, error)
	ReadString(logicalName string, bitLength uint32, readerArgs ...WithReaderArgs) (string, error)
	// CloseContext signals that we expect the end of the context with the supplied logical name
	CloseContext(logicalName string, readerArgs ...WithReaderArgs) error
}

// WithReaderArgs is a marker interface for reader args supplied by the builders
type WithReaderArgs interface {
	isReaderArgs() bool
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type readerArg struct {
}

func (_ readerArg) isReaderArgs() bool {
	return true
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////
