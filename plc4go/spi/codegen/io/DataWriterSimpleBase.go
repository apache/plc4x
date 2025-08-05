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

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// DataWriterSimpleBaseContract is a DataWriter minus the Read method.
type DataWriterSimpleBaseContract[T any] interface {
	utils.ByteOrderAware
	ContextWriter
	GetWriteBuffer() utils.WriteBuffer
}

type DataWriterSimpleBase[T any] struct {
	WriteBuffer utils.WriteBuffer
	bitLength   uint
}

var _ DataWriterSimpleBaseContract[string] = (*DataWriterSimpleBase[string])(nil)

func NewDataWriterSimpleBase[T any](WriteBuffer utils.WriteBuffer, bitLength uint) *DataWriterSimpleBase[T] {
	return &DataWriterSimpleBase[T]{
		WriteBuffer: WriteBuffer,
		bitLength:   bitLength,
	}
}

func (d *DataWriterSimpleBase[T]) GetByteOrder() binary.ByteOrder {
	return d.WriteBuffer.GetByteOrder()
}

func (d *DataWriterSimpleBase[T]) SetByteOrder(byteOrder binary.ByteOrder) {
	d.WriteBuffer.SetByteOrder(byteOrder)
}

func (d *DataWriterSimpleBase[T]) PushContext(logicalName string, writerArgs ...utils.WithWriterArgs) error {
	return d.WriteBuffer.PushContext(logicalName, writerArgs...)
}

func (d *DataWriterSimpleBase[T]) PopContext(logicalName string, writerArgs ...utils.WithWriterArgs) error {
	return d.WriteBuffer.PopContext(logicalName, writerArgs...)
}

func (d *DataWriterSimpleBase[T]) GetWriteBuffer() utils.WriteBuffer {
	return d.WriteBuffer
}
