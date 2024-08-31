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
	"errors"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DataWriterSimpleByteArrayAA struct {
	*DataWriterSimpleBase[[]byte]
}

var _ DataWriter[[]byte] = (*DataWriterSimpleByteArrayAA)(nil)

func NewDataWriterSimpleByteArray(WritBuffer utils.WriteBuffer, bitLength uint8) *DataWriterSimpleByteArrayAA {
	return &DataWriterSimpleByteArrayAA{
		DataWriterSimpleBase: NewDataWriterSimpleBase[[]byte](WritBuffer, uint(bitLength)),
	}
}

func (d *DataWriterSimpleByteArrayAA) Write(ctx context.Context, logicalName string, value []byte, args ...utils.WithWriterArgs) error {
	if d.bitLength != 8 {
		return errors.New("ByteArray fields only support bitLength of 8")
	}
	// TODO: Get a WirReaderArgs parameter for the number of bytes ...
	return d.WriteBuffer.WriteByteArray(logicalName, value, args...)
}