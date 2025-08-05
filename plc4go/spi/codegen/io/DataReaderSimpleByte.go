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

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DataReaderSimpleByte struct {
	*DataReaderSimpleBase[byte]
}

var _ DataReader[byte] = (*DataReaderSimpleByte)(nil)

func NewDataReaderSimpleByte(readBuffer utils.ReadBuffer, bitLength uint8) *DataReaderSimpleByte {
	return &DataReaderSimpleByte{
		DataReaderSimpleBase: NewDataReaderSimpleBase[byte](readBuffer, uint(bitLength)),
	}
}

func (d *DataReaderSimpleByte) Read(ctx context.Context, logicalName string, readerArgs ...utils.WithReaderArgs) (byte, error) {
	if d.bitLength != 8 {
		return 0, errors.New("bit length must be 8 bytes")
	}
	return d.readBuffer.ReadByte(logicalName, readerArgs...)
}
