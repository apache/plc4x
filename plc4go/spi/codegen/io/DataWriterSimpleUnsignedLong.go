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

	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DataWriterSimpleUnsignedLong struct {
	*DataWriterSimpleBase[uint64]
}

var _ DataWriter[uint64] = (*DataWriterSimpleUnsignedLong)(nil)

func NewDataWriterSimpleUnsignedLong(WritBuffer utils.WriteBuffer, bitLength uint8) *DataWriterSimpleUnsignedLong {
	return &DataWriterSimpleUnsignedLong{
		DataWriterSimpleBase: NewDataWriterSimpleBase[uint64](WritBuffer, uint(bitLength)),
	}
}

func (d *DataWriterSimpleUnsignedLong) Write(ctx context.Context, logicalName string, value uint64, args ...utils.WithWriterArgs) error {
	return d.WriteBuffer.WriteUint64(logicalName, uint8(d.bitLength), value, args...)
}
