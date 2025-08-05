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
	"encoding/binary"

	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DataWriterEnumDefault[T any, I any] struct {
	enumSerializer func(T) I
	enumNamer      func(T) string
	dataWriter     DataWriter[I]
}

var _ DataWriterEnum[string] = (*DataWriterEnumDefault[string, string])(nil)

func NewDataWriterEnumDefault[T any, I any](enumSerializer func(T) I, enumNamer func(T) string, dataWriter DataWriter[I]) *DataWriterEnumDefault[T, I] {
	return &DataWriterEnumDefault[T, I]{
		enumSerializer: enumSerializer,
		enumNamer:      enumNamer,
		dataWriter:     dataWriter,
	}
}

func (d *DataWriterEnumDefault[T, I]) GetByteOrder() binary.ByteOrder {
	return d.dataWriter.GetByteOrder()
}

func (d *DataWriterEnumDefault[T, I]) SetByteOrder(byteOrder binary.ByteOrder) {
	d.dataWriter.SetByteOrder(byteOrder)
}

func (d *DataWriterEnumDefault[T, I]) PushContext(logicalName string, writerArgs ...utils.WithWriterArgs) error {
	return d.dataWriter.PushContext(logicalName, writerArgs...)
}

func (d *DataWriterEnumDefault[T, I]) PopContext(logicalName string, writerArgs ...utils.WithWriterArgs) error {
	return d.dataWriter.PopContext(logicalName, writerArgs...)
}

func (d *DataWriterEnumDefault[T, I]) Write(ctx context.Context, logicalName string, value T, writerArgs ...utils.WithWriterArgs) error {
	return WriteWithRawWriter(ctx, logicalName, value, d.enumSerializer, d.enumNamer, d.dataWriter, writerArgs...)
}

func WriteWithRawWriter[T any, I any](ctx context.Context, logicalName string, value T, enumSerializer func(T) I, enumNamer func(T) string, rawWriter DataWriter[I], args ...utils.WithWriterArgs) error {
	rawValue := enumSerializer(value)
	return rawWriter.Write(ctx, logicalName, rawValue, append(args, utils.WithAdditionalStringRepresentation(enumNamer(value)))...)
}

func (d *DataWriterEnumDefault[T, I]) GetWriteBuffer() utils.WriteBuffer {
	return d.dataWriter.GetWriteBuffer()
}
