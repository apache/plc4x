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

	"github.com/pkg/errors"

	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DataWriterComplexDefault[T spi.Message] struct {
	writeBuffer utils.WriteBuffer
}

var _ DataWriterComplex[spi.Message] = (*DataWriterComplexDefault[spi.Message])(nil)

func NewDataWriterComplexDefault[T spi.Message](writeBuffer utils.WriteBuffer) *DataWriterComplexDefault[T] {
	return &DataWriterComplexDefault[T]{
		writeBuffer: writeBuffer,
	}
}

func (d *DataWriterComplexDefault[T]) GetByteOrder() binary.ByteOrder {
	return d.writeBuffer.GetByteOrder()
}

func (d *DataWriterComplexDefault[T]) SetByteOrder(byteOrder binary.ByteOrder) {
	d.writeBuffer.SetByteOrder(byteOrder)
}

func (d *DataWriterComplexDefault[T]) PushContext(logicalName string, writerArgs ...utils.WithWriterArgs) error {
	return d.writeBuffer.PushContext(logicalName, writerArgs...)
}

func (d *DataWriterComplexDefault[T]) PopContext(logicalName string, writerArgs ...utils.WithWriterArgs) error {
	return d.writeBuffer.PopContext(logicalName, writerArgs...)
}

func (d *DataWriterComplexDefault[T]) Write(ctx context.Context, logicalName string, value T, writerArgs ...utils.WithWriterArgs) error {
	if logicalName != "" {
		err := d.writeBuffer.PushContext(logicalName, writerArgs...)
		if err != nil {
			return errors.Wrap(err, "error pushing context")
		}
	}
	err := d.writeBuffer.WriteSerializable(ctx, value)
	if err != nil {
		return errors.Wrap(err, "error getting value")
	}
	if logicalName != "" {
		err := d.writeBuffer.PopContext(logicalName, writerArgs...)
		if err != nil {
			return errors.Wrap(err, "error popping context")
		}
	}
	return nil
}

func (d *DataWriterComplexDefault[T]) GetWriteBuffer() utils.WriteBuffer {
	return d.writeBuffer
}
