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
	"encoding/binary"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// DataReaderSimpleBaseContract is a DataReader minus the Read method.
type DataReaderSimpleBaseContract[T any] interface {
	utils.ByteOrderAware
	codegen.PositionMover
	ContextReader
	GetReadBuffer() utils.ReadBuffer
}

type DataReaderSimpleBase[T any] struct {
	readBuffer utils.ReadBuffer
	bitLength  uint
}

var _ DataReaderSimpleBaseContract[string] = (*DataReaderSimpleBase[string])(nil)

func NewDataReaderSimpleBase[T any](readBuffer utils.ReadBuffer, bitLength uint) *DataReaderSimpleBase[T] {
	return &DataReaderSimpleBase[T]{
		readBuffer: readBuffer,
		bitLength:  bitLength,
	}
}

func (d *DataReaderSimpleBase[T]) GetPos() uint16 {
	return d.readBuffer.GetPos()
}

func (d *DataReaderSimpleBase[T]) SetPos(pos uint16) {
	d.readBuffer.Reset(pos)
}

func (d *DataReaderSimpleBase[T]) GetByteOrder() binary.ByteOrder {
	return d.readBuffer.GetByteOrder()
}

func (d *DataReaderSimpleBase[T]) SetByteOrder(byteOrder binary.ByteOrder) {
	d.readBuffer.SetByteOrder(byteOrder)
}

func (d *DataReaderSimpleBase[T]) PullContext(logicalName string, readerArgs ...utils.WithReaderArgs) error {
	return d.readBuffer.PullContext(logicalName, readerArgs...)
}

func (d *DataReaderSimpleBase[T]) CloseContext(logicalName string, readerArgs ...utils.WithReaderArgs) error {
	return d.readBuffer.CloseContext(logicalName, readerArgs...)
}

func (d *DataReaderSimpleBase[T]) GetReadBuffer() utils.ReadBuffer {
	return d.readBuffer
}
