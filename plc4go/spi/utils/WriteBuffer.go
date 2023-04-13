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
	"context"
	"math/big"
)

type WriteBuffer interface {
	PositionAware
	// PushContext signals opening context with the supplied logical name
	PushContext(logicalName string, writerArgs ...WithWriterArgs) error
	WriteBit(logicalName string, value bool, writerArgs ...WithWriterArgs) error
	WriteByte(logicalName string, value byte, writerArgs ...WithWriterArgs) error
	WriteByteArray(logicalName string, data []byte, writerArgs ...WithWriterArgs) error
	WriteUint8(logicalName string, bitLength uint8, value uint8, writerArgs ...WithWriterArgs) error
	WriteUint16(logicalName string, bitLength uint8, value uint16, writerArgs ...WithWriterArgs) error
	WriteUint32(logicalName string, bitLength uint8, value uint32, writerArgs ...WithWriterArgs) error
	WriteUint64(logicalName string, bitLength uint8, value uint64, writerArgs ...WithWriterArgs) error
	WriteInt8(logicalName string, bitLength uint8, value int8, writerArgs ...WithWriterArgs) error
	WriteInt16(logicalName string, bitLength uint8, value int16, writerArgs ...WithWriterArgs) error
	WriteInt32(logicalName string, bitLength uint8, value int32, writerArgs ...WithWriterArgs) error
	WriteInt64(logicalName string, bitLength uint8, value int64, writerArgs ...WithWriterArgs) error
	WriteBigInt(logicalName string, bitLength uint8, value *big.Int, writerArgs ...WithWriterArgs) error
	WriteFloat32(logicalName string, bitLength uint8, value float32, writerArgs ...WithWriterArgs) error
	WriteFloat64(logicalName string, bitLength uint8, value float64, writerArgs ...WithWriterArgs) error
	WriteBigFloat(logicalName string, bitLength uint8, value *big.Float, writerArgs ...WithWriterArgs) error
	WriteString(logicalName string, bitLength uint32, encoding string, value string, writerArgs ...WithWriterArgs) error
	WriteVirtual(ctx context.Context, logicalName string, value interface{}, writerArgs ...WithWriterArgs) error
	WriteSerializable(ctx context.Context, serializable Serializable) error
	// PopContext signals work done with the context with the supplied logical name
	PopContext(logicalName string, writerArgs ...WithWriterArgs) error
}

// WithWriterArgs is a marker interface for writer args supplied by the builders like WithAdditionalStringRepresentation
type WithWriterArgs interface {
	isWriterArgs() bool
}

///////////////////////////////////////
///////////////////////////////////////
//
// Internal section
//

type writerArg struct {
}

func (_ writerArg) isWriterArgs() bool {
	return true
}

//
// Internal section
//
///////////////////////////////////////
///////////////////////////////////////
