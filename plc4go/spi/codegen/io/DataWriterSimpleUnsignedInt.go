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

type DataWriterSimpleUnsignedInt struct {
	*DataWriterSimpleBase[uint32]
}

var _ DataWriter[uint32] = (*DataWriterSimpleUnsignedInt)(nil)

func NewDataWriterSimpleUnsignedInt(WritBuffer utils.WriteBuffer, bitLength uint8) *DataWriterSimpleUnsignedInt {
	return &DataWriterSimpleUnsignedInt{
		DataWriterSimpleBase: NewDataWriterSimpleBase[uint32](WritBuffer, uint(bitLength)),
	}
}

func (d *DataWriterSimpleUnsignedInt) Write(ctx context.Context, logicalName string, value uint32, args ...utils.WithWriterArgs) error {
	return d.WriteBuffer.WriteUint32(logicalName, uint8(d.bitLength), value, args...)
}
