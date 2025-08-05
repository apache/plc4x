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

type DataReaderSimpleBoolean struct {
	*DataReaderSimpleBase[bool]
}

var _ DataReader[bool] = (*DataReaderSimpleBoolean)(nil)

func NewDataReaderSimpleBoolean(readBuffer utils.ReadBuffer) *DataReaderSimpleBoolean {
	return &DataReaderSimpleBoolean{
		DataReaderSimpleBase: NewDataReaderSimpleBase[bool](readBuffer, 1),
	}
}

func (d *DataReaderSimpleBoolean) Read(ctx context.Context, logicalName string, readerArgs ...utils.WithReaderArgs) (bool, error) {
	return d.readBuffer.ReadBit(logicalName, readerArgs...)
}
